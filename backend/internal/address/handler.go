package address

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler handles HTTP requests for address operations
type Handler struct {
	service Service
}

// NewHandler creates a new address handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers address routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Address CRUD operations
	router.HandleFunc("/addresses", h.CreateAddress).Methods("POST")
	router.HandleFunc("/addresses", h.ListAddresses).Methods("GET")
	router.HandleFunc("/addresses/{id}", h.GetAddress).Methods("GET")
	router.HandleFunc("/addresses/{id}", h.UpdateAddress).Methods("PUT")
	router.HandleFunc("/addresses/{id}", h.DeleteAddress).Methods("DELETE")
	
	// Customer address operations
	router.HandleFunc("/customers/{customerId}/addresses", h.GetCustomerAddresses).Methods("GET")
	router.HandleFunc("/customers/{customerId}/addresses/default", h.GetDefaultAddress).Methods("GET")
	router.HandleFunc("/customers/{customerId}/addresses/{id}/default", h.SetDefaultAddress).Methods("PUT")
	router.HandleFunc("/customers/{customerId}/addresses/default", h.UnsetDefaultAddresses).Methods("DELETE")
	
	// Address validation operations
	router.HandleFunc("/addresses/{id}/validate", h.ValidateAddress).Methods("POST")
	router.HandleFunc("/addresses/{id}/validation", h.GetAddressValidation).Methods("GET")
	router.HandleFunc("/address-validations", h.ListAddressValidations).Methods("GET")
	
	// Bulk operations
	router.HandleFunc("/addresses/bulk", h.BulkCreateAddresses).Methods("POST")
	router.HandleFunc("/addresses/bulk", h.BulkUpdateAddresses).Methods("PUT")
	router.HandleFunc("/addresses/bulk", h.BulkDeleteAddresses).Methods("DELETE")
	
	// Statistics and analytics
	router.HandleFunc("/addresses/stats", h.GetAddressStats).Methods("GET")
	router.HandleFunc("/addresses/stats/country", h.GetAddressesByCountry).Methods("GET")
	router.HandleFunc("/addresses/stats/type", h.GetAddressesByType).Methods("GET")
	router.HandleFunc("/addresses/recent", h.GetRecentAddresses).Methods("GET")
	
	// Maintenance operations
	router.HandleFunc("/addresses/cleanup/unvalidated", h.CleanupUnvalidatedAddresses).Methods("DELETE")
	router.HandleFunc("/address-validations/cleanup/orphaned", h.CleanupOrphanedValidations).Methods("DELETE")
	
	// Utility operations
	router.HandleFunc("/addresses/normalize", h.NormalizeAddress).Methods("POST")
	router.HandleFunc("/addresses/suggest", h.SuggestAddresses).Methods("POST")
}

// Address CRUD operations

// CreateAddress creates a new address
func (h *Handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	var req CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	response, err := h.service.CreateAddress(r.Context(), tenantID, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusCreated, response)
}

// GetAddress retrieves an address by ID
func (h *Handler) GetAddress(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	addressID, err := h.getUUIDParam(r, "id")
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid address ID", err)
		return
	}
	
	response, err := h.service.GetAddress(r.Context(), tenantID, addressID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, response)
}

// UpdateAddress updates an existing address
func (h *Handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	addressID, err := h.getUUIDParam(r, "id")
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid address ID", err)
		return
	}
	
	var req UpdateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	response, err := h.service.UpdateAddress(r.Context(), tenantID, addressID, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, response)
}

// DeleteAddress soft deletes an address
func (h *Handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	addressID, err := h.getUUIDParam(r, "id")
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid address ID", err)
		return
	}
	
	if err := h.service.DeleteAddress(r.Context(), tenantID, addressID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// ListAddresses retrieves addresses with filtering and pagination
func (h *Handler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Parse query parameters
	filter := h.parseAddressFilter(r)
	limit, offset := h.parsePagination(r)
	
	response, err := h.service.ListAddresses(r.Context(), tenantID, filter, limit, offset)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, response)
}

// Customer address operations

// GetCustomerAddresses retrieves all addresses for a customer
func (h *Handler) GetCustomerAddresses(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	customerID, err := h.getUUIDParam(r, "customerId")
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid customer ID", err)
		return
	}
	
	responses, err := h.service.GetCustomerAddresses(r.Context(), tenantID, customerID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"addresses": responses,
	})
}

// GetDefaultAddress retrieves the default address for a customer and type
func (h *Handler) GetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	customerID, err := h.getUUIDParam(r, "customerId")
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid customer ID", err)
		return
	}
	
	addressType := r.URL.Query().Get("type")
	
	response, err := h.service.GetDefaultAddress(r.Context(), tenantID, customerID, addressType)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, response)
}

// SetDefaultAddress sets an address as default for a customer
func (h *Handler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	customerID, err := h.getUUIDParam(r, "customerId")
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid customer ID", err)
		return
	}
	
	addressID, err := h.getUUIDParam(r, "id")
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid address ID", err)
		return
	}
	
	if err := h.service.SetDefaultAddress(r.Context(), tenantID, customerID, addressID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// UnsetDefaultAddresses unsets default addresses for a customer and type
func (h *Handler) UnsetDefaultAddresses(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	customerID, err := h.getUUIDParam(r, "customerId")
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid customer ID", err)
		return
	}
	
	addressType := r.URL.Query().Get("type")
	
	if err := h.service.UnsetDefaultAddresses(r.Context(), tenantID, customerID, addressType); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// Address validation operations

// ValidateAddress validates an address using external service
func (h *Handler) ValidateAddress(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	addressID, err := h.getUUIDParam(r, "id")
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid address ID", err)
		return
	}
	
	var req ValidateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	response, err := h.service.ValidateAddress(r.Context(), tenantID, addressID, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetAddressValidation retrieves the latest validation for an address
func (h *Handler) GetAddressValidation(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	addressID, err := h.getUUIDParam(r, "id")
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid address ID", err)
		return
	}
	
	response, err := h.service.GetAddressValidation(r.Context(), tenantID, addressID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, response)
}

// ListAddressValidations retrieves validations with pagination
func (h *Handler) ListAddressValidations(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	limit, offset := h.parsePagination(r)
	
	response, err := h.service.ListAddressValidations(r.Context(), tenantID, limit, offset)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, response)
}

// Bulk operations

// BulkCreateAddresses creates multiple addresses
func (h *Handler) BulkCreateAddresses(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	var req struct {
		Addresses []CreateAddressRequest `json:"addresses"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	responses, err := h.service.BulkCreateAddresses(r.Context(), tenantID, req.Addresses)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"addresses": responses,
		"count":     len(responses),
	})
}

// BulkUpdateAddresses updates multiple addresses
func (h *Handler) BulkUpdateAddresses(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	var req struct {
		Updates []BulkUpdateAddressRequest `json:"updates"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	responses, err := h.service.BulkUpdateAddresses(r.Context(), tenantID, req.Updates)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"addresses": responses,
		"count":     len(responses),
	})
}

// BulkDeleteAddresses soft deletes multiple addresses
func (h *Handler) BulkDeleteAddresses(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	var req struct {
		AddressIDs []uuid.UUID `json:"address_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	if err := h.service.BulkDeleteAddresses(r.Context(), tenantID, req.AddressIDs); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"deleted_count": len(req.AddressIDs),
	})
}

// Statistics and analytics

// GetAddressStats retrieves address statistics
func (h *Handler) GetAddressStats(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	stats, err := h.service.GetAddressStats(r.Context(), tenantID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, stats)
}

// GetAddressesByCountry retrieves address count by country
func (h *Handler) GetAddressesByCountry(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	stats, err := h.service.GetAddressesByCountry(r.Context(), tenantID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"countries": stats,
	})
}

// GetAddressesByType retrieves address count by type
func (h *Handler) GetAddressesByType(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	stats, err := h.service.GetAddressesByType(r.Context(), tenantID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"types": stats,
	})
}

// GetRecentAddresses retrieves recent addresses
func (h *Handler) GetRecentAddresses(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	days := 30 // Default
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}
	
	addresses, err := h.service.GetRecentAddresses(r.Context(), tenantID, days)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"addresses": addresses,
		"days":      days,
	})
}

// Maintenance operations

// CleanupUnvalidatedAddresses removes old unvalidated addresses
func (h *Handler) CleanupUnvalidatedAddresses(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	days := 90 // Default
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}
	
	deleted, err := h.service.CleanupUnvalidatedAddresses(r.Context(), tenantID, days)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"deleted_count": deleted,
		"days":         days,
	})
}

// CleanupOrphanedValidations removes validations for non-existent addresses
func (h *Handler) CleanupOrphanedValidations(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	deleted, err := h.service.CleanupOrphanedValidations(r.Context(), tenantID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"deleted_count": deleted,
	})
}

// Utility operations

// NormalizeAddress normalizes address data
func (h *Handler) NormalizeAddress(w http.ResponseWriter, r *http.Request) {
	var req NormalizeAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	response, err := h.service.NormalizeAddress(r.Context(), req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, response)
}

// SuggestAddresses provides address suggestions
func (h *Handler) SuggestAddresses(w http.ResponseWriter, r *http.Request) {
	var req AddressSuggestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	response, err := h.service.SuggestAddresses(r.Context(), req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, response)
}

// Helper methods

// getTenantID extracts tenant ID from request context or headers
func (h *Handler) getTenantID(r *http.Request) (uuid.UUID, error) {
	// This would typically come from JWT token or request context
	// For now, we'll get it from header
	tenantIDStr := r.Header.Get("X-Tenant-ID")
	if tenantIDStr == "" {
		return uuid.Nil, fmt.Errorf("tenant ID is required")
	}
	
	return uuid.Parse(tenantIDStr)
}

// getUUIDParam extracts UUID parameter from URL
func (h *Handler) getUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	vars := mux.Vars(r)
	idStr, exists := vars[param]
	if !exists {
		return uuid.Nil, fmt.Errorf("%s parameter is required", param)
	}
	
	return uuid.Parse(idStr)
}

// parseAddressFilter parses address filter from query parameters
func (h *Handler) parseAddressFilter(r *http.Request) AddressFilter {
	filter := AddressFilter{}
	
	if customerIDStr := r.URL.Query().Get("customer_id"); customerIDStr != "" {
		if customerID, err := uuid.Parse(customerIDStr); err == nil {
			filter.CustomerID = &customerID
		}
	}
	
	if addressType := r.URL.Query().Get("type"); addressType != "" {
		filter.Type = addressType
	}
	
	if label := r.URL.Query().Get("label"); label != "" {
		filter.Label = label
	}
	
	if country := r.URL.Query().Get("country"); country != "" {
		filter.Country = country
	}
	
	if state := r.URL.Query().Get("state"); state != "" {
		filter.State = state
	}
	
	if city := r.URL.Query().Get("city"); city != "" {
		filter.City = city
	}
	
	if isDefaultStr := r.URL.Query().Get("is_default"); isDefaultStr != "" {
		if isDefault, err := strconv.ParseBool(isDefaultStr); err == nil {
			filter.IsDefault = &isDefault
		}
	}
	
	if isValidatedStr := r.URL.Query().Get("is_validated"); isValidatedStr != "" {
		if isValidated, err := strconv.ParseBool(isValidatedStr); err == nil {
			filter.IsValidated = &isValidated
		}
	}
	
	if search := r.URL.Query().Get("search"); search != "" {
		filter.Search = search
	}
	
	return filter
}

// parsePagination parses pagination parameters
func (h *Handler) parsePagination(r *http.Request) (limit, offset int) {
	limit = DefaultPageSize
	offset = 0
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > MaxPageSize {
				limit = MaxPageSize
			}
		}
	}
	
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	return limit, offset
}

// writeJSONResponse writes JSON response
func (h *Handler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeErrorResponse writes error response
func (h *Handler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	errorResponse := map[string]interface{}{
		"error":   message,
		"status":  statusCode,
		"details": err.Error(),
	}
	
	h.writeJSONResponse(w, statusCode, errorResponse)
}

// handleServiceError handles service layer errors
func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case err == ErrAddressNotFound:
		h.writeErrorResponse(w, http.StatusNotFound, "Address not found", err)
	case err == ErrValidationNotFound:
		h.writeErrorResponse(w, http.StatusNotFound, "Validation not found", err)
	case err == ErrTooManyAddresses:
		h.writeErrorResponse(w, http.StatusBadRequest, "Too many addresses for customer", err)
	case err == ErrBulkSizeExceeded:
		h.writeErrorResponse(w, http.StatusBadRequest, "Bulk operation size exceeded", err)
	case err == ErrInvalidCustomerID:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid customer ID", err)
	case err == ErrInvalidAddressType:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid address type", err)
	case err == ErrInvalidFirstName:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid first name", err)
	case err == ErrInvalidLastName:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid last name", err)
	case err == ErrInvalidAddress:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid address", err)
	case err == ErrInvalidCity:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid city", err)
	case err == ErrInvalidState:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid state", err)
	case err == ErrInvalidPostalCode:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid postal code", err)
	case err == ErrInvalidCountry:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid country", err)
	case strings.Contains(err.Error(), "validation failed"):
		h.writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
	default:
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error", err)
	}
}