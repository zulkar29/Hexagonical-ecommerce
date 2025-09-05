package address

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// 📍 CORE ADDRESS ENDPOINTS (5)
	addresses := router.Group("/addresses")
	{
		addresses.POST("", h.CreateAddress)                    // CreateAddress
		addresses.GET("", h.ListAddresses)                     // ListAddresses (with filtering, stats, recent)
		addresses.GET("/:id", h.GetAddress)                    // GetAddress
		addresses.PUT("/:id", h.UpdateAddress)                 // UpdateAddress (handles validation, default setting)
		addresses.DELETE("/:id", h.DeleteAddress)              // DeleteAddress
	}
	
	// 📍 ADDRESS OPERATIONS (4)
	addresses.POST("/bulk", h.BulkOperations)               // BulkOperations (create/update/delete via operation type)
	addresses.POST("/normalize", h.NormalizeAddress)        // NormalizeAddress
	addresses.POST("/suggest", h.SuggestAddresses)          // SuggestAddresses
	addresses.DELETE("/cleanup", h.CleanupOperations)       // CleanupOperations (unvalidated/orphaned via type)
	
	// 📍 ADDRESS VALIDATIONS (2)
	validations := router.Group("/address-validations")
	{
		validations.GET("", h.ListAddressValidations)        // ListAddressValidations
		validations.DELETE("/cleanup", h.CleanupOperations)  // CleanupOperations (orphaned validations)
	}
}

// Address CRUD operations

// CreateAddress creates a new address
func (h *Handler) CreateAddress(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	response, err := h.service.CreateAddress(ctx, tenantID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, response)
}

// GetAddress retrieves an address by ID
func (h *Handler) GetAddress(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	addressID, err := h.getUUIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID", "details": err.Error()})
		return
	}
	
	response, err := h.service.GetAddress(ctx, tenantID, addressID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// UpdateAddress updates an existing address
func (h *Handler) UpdateAddress(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	addressID, err := h.getUUIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID", "details": err.Error()})
		return
	}
	
	var req UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	response, err := h.service.UpdateAddress(ctx, tenantID, addressID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// DeleteAddress soft deletes an address
func (h *Handler) DeleteAddress(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	addressID, err := h.getUUIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID", "details": err.Error()})
		return
	}
	
	if err := h.service.DeleteAddress(ctx, tenantID, addressID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ListAddresses handles GET /addresses with consolidated functionality
// Supports: filtering, stats, recent addresses, customer addresses
func (h *Handler) ListAddresses(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Check for special operations via query parameters
	operationType := c.Query("type")
	customerID := c.Query("customer_id")
	
	// Handle stats operations
	if operationType == "stats" {
		category := c.Query("category")
		switch category {
		case "country":
			h.GetAddressesByCountry(c)
			return
		case "type":
			h.GetAddressesByType(c)
			return
		default:
			h.GetAddressStats(c)
			return
		}
	}
	
	// Handle recent addresses
	if operationType == "recent" {
		h.GetRecentAddresses(c)
		return
	}
	
	// Handle customer-specific addresses
	if customerID != "" {
		if customerUUID, err := uuid.Parse(customerID); err == nil {
			c.Set("customerId", customerUUID.String())
			h.GetCustomerAddresses(c)
			return
		}
	}
	
	// Parse query parameters for standard listing
	filter := h.parseAddressFilter(c)
	limit, offset := h.parsePagination(c)
	
	response, err := h.service.ListAddresses(ctx, tenantID, filter, limit, offset)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// Customer address operations

// GetCustomerAddresses retrieves all addresses for a customer
func (h *Handler) GetCustomerAddresses(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	customerID, err := h.getUUIDParam(c, "customerId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
		return
	}
	
	responses, err := h.service.GetCustomerAddresses(ctx, tenantID, customerID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"addresses": responses,
	})
}

// GetDefaultAddress retrieves the default address for a customer and type
func (h *Handler) GetDefaultAddress(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	customerID, err := h.getUUIDParam(c, "customerId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
		return
	}
	
	addressType := c.Query("type")
	
	response, err := h.service.GetDefaultAddress(ctx, tenantID, customerID, addressType)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// SetDefaultAddress sets an address as default for a customer
func (h *Handler) SetDefaultAddress(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	customerID, err := h.getUUIDParam(c, "customerId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
		return
	}
	
	addressID, err := h.getUUIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID", "details": err.Error()})
		return
	}
	
	if err := h.service.SetDefaultAddress(ctx, tenantID, customerID, addressID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// UnsetDefaultAddresses unsets default addresses for a customer and type
func (h *Handler) UnsetDefaultAddresses(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	customerID, err := h.getUUIDParam(c, "customerId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
		return
	}
	
	addressType := c.Query("type")
	
	if err := h.service.UnsetDefaultAddresses(ctx, tenantID, customerID, addressType); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.Status(http.StatusNoContent)
}

// Address validation operations

// ValidateAddress validates an address using external service
func (h *Handler) ValidateAddress(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	addressID, err := h.getUUIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID", "details": err.Error()})
		return
	}
	
	var req ValidateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	response, err := h.service.ValidateAddress(ctx, tenantID, addressID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// GetAddressValidation retrieves the latest validation for an address
func (h *Handler) GetAddressValidation(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	addressID, err := h.getUUIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID", "details": err.Error()})
		return
	}
	
	response, err := h.service.GetAddressValidation(ctx, tenantID, addressID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// ListAddressValidations retrieves validations with pagination
func (h *Handler) ListAddressValidations(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	limit, offset := h.parsePagination(c)
	
	response, err := h.service.ListAddressValidations(ctx, tenantID, limit, offset)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// Bulk operations

// BulkOperations handles consolidated bulk operations (create/update/delete)
func (h *Handler) BulkOperations(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Determine operation type from query parameter or request body
	operationType := c.Query("operation")
	if operationType == "" {
		// Default to create for backward compatibility
		operationType = "create"
	}
	
	switch operationType {
	case "create":
		h.handleBulkCreate(c, ctx, tenantID)
	case "update":
		h.handleBulkUpdate(c, ctx, tenantID)
	case "delete":
		h.handleBulkDelete(c, ctx, tenantID)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type. Supported: create, update, delete"})
	}
}

// handleBulkCreate handles bulk address creation
func (h *Handler) handleBulkCreate(c *gin.Context, ctx context.Context, tenantID uuid.UUID) {
	var req struct {
		Addresses []CreateAddressRequest `json:"addresses"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	responses, err := h.service.BulkCreateAddresses(ctx, tenantID, req.Addresses)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"addresses": responses,
		"count":     len(responses),
	})
}

// handleBulkUpdate handles bulk address updates
func (h *Handler) handleBulkUpdate(c *gin.Context, ctx context.Context, tenantID uuid.UUID) {
	var req struct {
		Updates []BulkUpdateAddressRequest `json:"updates"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	responses, err := h.service.BulkUpdateAddresses(ctx, tenantID, req.Updates)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"addresses": responses,
		"count":     len(responses),
	})
}

// handleBulkDelete handles bulk address deletion
func (h *Handler) handleBulkDelete(c *gin.Context, ctx context.Context, tenantID uuid.UUID) {
	var req struct {
		AddressIDs []uuid.UUID `json:"address_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	if err := h.service.BulkDeleteAddresses(ctx, tenantID, req.AddressIDs); err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"deleted_count": len(req.AddressIDs),
	})
}

// CleanupOperations handles consolidated cleanup operations
func (h *Handler) CleanupOperations(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	// Determine cleanup type from query parameter
	cleanupType := c.Query("type")
	if cleanupType == "" {
		cleanupType = "unvalidated" // Default
	}
	
	switch cleanupType {
	case "unvalidated":
		h.handleCleanupUnvalidated(c, ctx, tenantID)
	case "orphaned":
		h.handleCleanupOrphaned(c, ctx, tenantID)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cleanup type. Supported: unvalidated, orphaned"})
	}
}

// handleCleanupUnvalidated handles cleanup of unvalidated addresses
func (h *Handler) handleCleanupUnvalidated(c *gin.Context, ctx context.Context, tenantID uuid.UUID) {
	days := 90 // Default
	if daysStr := c.Query("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}
	
	deleted, err := h.service.CleanupUnvalidatedAddresses(ctx, tenantID, days)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"deleted_count": deleted,
		"days":         days,
		"type":         "unvalidated",
	})
}

// handleCleanupOrphaned handles cleanup of orphaned validations
func (h *Handler) handleCleanupOrphaned(c *gin.Context, ctx context.Context, tenantID uuid.UUID) {
	deleted, err := h.service.CleanupOrphanedValidations(ctx, tenantID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"deleted_count": deleted,
		"type":         "orphaned",
	})
}











// Utility operations

// NormalizeAddress normalizes address data
func (h *Handler) NormalizeAddress(c *gin.Context) {
	ctx := c.Request.Context()
	var req NormalizeAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	response, err := h.service.NormalizeAddress(ctx, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// SuggestAddresses provides address suggestions
func (h *Handler) SuggestAddresses(c *gin.Context) {
	ctx := c.Request.Context()
	var req AddressSuggestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	
	response, err := h.service.SuggestAddresses(ctx, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// Statistics handler methods

// GetAddressStats retrieves address statistics
func (h *Handler) GetAddressStats(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	stats, err := h.service.GetAddressStats(ctx, tenantID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// GetAddressesByCountry retrieves address count by country
func (h *Handler) GetAddressesByCountry(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	addressesByCountry, err := h.service.GetAddressesByCountry(ctx, tenantID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"addresses_by_country": addressesByCountry})
}

// GetAddressesByType retrieves address count by type
func (h *Handler) GetAddressesByType(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	addressesByType, err := h.service.GetAddressesByType(ctx, tenantID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"addresses_by_type": addressesByType})
}

// GetRecentAddresses retrieves recent addresses
func (h *Handler) GetRecentAddresses(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}
	
	days := 30 // Default
	if daysStr := c.Query("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}
	
	addresses, err := h.service.GetRecentAddresses(ctx, tenantID, days)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"addresses": addresses, "days": days})
}

// Helper methods

// getTenantID extracts tenant ID from request context or headers
func (h *Handler) getTenantID(c *gin.Context) (uuid.UUID, error) {
	// This would typically come from JWT token or request context
	// For now, we'll get it from header
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		return uuid.Nil, fmt.Errorf("tenant ID is required")
	}
	
	return uuid.Parse(tenantIDStr)
}

// getUUIDParam extracts UUID parameter from URL
func (h *Handler) getUUIDParam(c *gin.Context, param string) (uuid.UUID, error) {
	idStr := c.Param(param)
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("%s parameter is required", param)
	}
	
	return uuid.Parse(idStr)
}

// parseAddressFilter parses address filter from query parameters
func (h *Handler) parseAddressFilter(c *gin.Context) AddressFilter {
	filter := AddressFilter{}
	
	if customerIDStr := c.Query("customer_id"); customerIDStr != "" {
		if customerID, err := uuid.Parse(customerIDStr); err == nil {
			filter.CustomerID = &customerID
		}
	}
	
	if addressType := c.Query("type"); addressType != "" {
		filter.Type = addressType
	}
	
	if label := c.Query("label"); label != "" {
		filter.Label = label
	}
	
	if country := c.Query("country"); country != "" {
		filter.Country = country
	}
	
	if state := c.Query("state"); state != "" {
		filter.State = state
	}
	
	if city := c.Query("city"); city != "" {
		filter.City = city
	}
	
	if isDefaultStr := c.Query("is_default"); isDefaultStr != "" {
		if isDefault, err := strconv.ParseBool(isDefaultStr); err == nil {
			filter.IsDefault = &isDefault
		}
	}
	
	if isValidatedStr := c.Query("is_validated"); isValidatedStr != "" {
		if isValidated, err := strconv.ParseBool(isValidatedStr); err == nil {
			filter.IsValidated = &isValidated
		}
	}
	
	if search := c.Query("search"); search != "" {
		filter.Search = search
	}
	
	return filter
}

// parsePagination parses pagination parameters
func (h *Handler) parsePagination(c *gin.Context) (limit, offset int) {
	limit = DefaultPageSize
	offset = 0
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > MaxPageSize {
				limit = MaxPageSize
			}
		}
	}
	
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	return limit, offset
}


// handleServiceError handles service layer errors
func (h *Handler) handleServiceError(c *gin.Context, err error) {
	switch {
	case err == ErrAddressNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found", "details": err.Error()})
	case err == ErrValidationNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Validation not found", "details": err.Error()})
	case err == ErrTooManyAddresses:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many addresses for customer", "details": err.Error()})
	case err == ErrBulkSizeExceeded:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulk operation size exceeded", "details": err.Error()})
	case err == ErrInvalidCustomerID:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
	case err == ErrInvalidAddressType:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address type", "details": err.Error()})
	case err == ErrInvalidFirstName:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid first name", "details": err.Error()})
	case err == ErrInvalidLastName:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid last name", "details": err.Error()})
	case err == ErrInvalidAddress:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address", "details": err.Error()})
	case err == ErrInvalidCity:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city", "details": err.Error()})
	case err == ErrInvalidState:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state", "details": err.Error()})
	case err == ErrInvalidPostalCode:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid postal code", "details": err.Error()})
	case err == ErrInvalidCountry:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country", "details": err.Error()})
	case strings.Contains(err.Error(), "validation failed"):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
	}
}