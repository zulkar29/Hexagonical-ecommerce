package category

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler handles HTTP requests for category operations
type Handler struct {
	service Service
}

// NewHandler creates a new category handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateCategory handles POST /categories
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID from context
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Parse request body
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	// Create category
	category, err := h.service.CreateCategory(ctx, tenantID, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusCreated, category)
}

// GetCategory handles GET /categories/{id}
func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}
	
	// Get category
	category, err := h.service.GetCategory(ctx, tenantID, categoryID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, category)
}

// GetCategoryBySlug handles GET /categories/slug/{slug}
func (h *Handler) GetCategoryBySlug(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Extract slug
	slug := mux.Vars(r)["slug"]
	if slug == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Slug is required", nil)
		return
	}
	
	// Get category by slug
	category, err := h.service.GetCategoryBySlug(ctx, tenantID, slug)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, category)
}

// UpdateCategory handles PUT /categories/{id}
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}
	
	// Parse request body
	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	// Update category
	category, err := h.service.UpdateCategory(ctx, tenantID, categoryID, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, category)
}

// DeleteCategory handles DELETE /categories/{id}
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}
	
	// Delete category
	if err := h.service.DeleteCategory(ctx, tenantID, categoryID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// ListCategories handles GET /categories
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Parse query parameters
	filter := h.parseCategoryFilter(r)
	limit, offset := h.parsePagination(r)
	
	// List categories
	categories, total, err := h.service.ListCategories(ctx, tenantID, filter, limit, offset)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	// Build response
	response := map[string]interface{}{
		"categories": categories,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	}
	
	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetCategoryTree handles GET /categories/tree
func (h *Handler) GetCategoryTree(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Parse parent ID if provided
	var parentID *uuid.UUID
	if parentIDStr := r.URL.Query().Get("parent_id"); parentIDStr != "" {
		parsedID, err := uuid.Parse(parentIDStr)
		if err != nil {
			h.writeErrorResponse(w, http.StatusBadRequest, "Invalid parent ID", err)
			return
		}
		parentID = &parsedID
	}
	
	// Get category tree
	tree, err := h.service.GetCategoryTree(ctx, tenantID, parentID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, tree)
}

// GetCategoryPath handles GET /categories/{id}/path
func (h *Handler) GetCategoryPath(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}
	
	// Get category path
	path, err := h.service.GetCategoryPath(ctx, tenantID, categoryID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, path)
}

// MoveCategory handles PUT /categories/{id}/move
func (h *Handler) MoveCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}
	
	// Parse request body
	var req struct {
		ParentID *uuid.UUID `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	// Move category
	if err := h.service.MoveCategory(ctx, tenantID, categoryID, req.ParentID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// ReorderCategories handles PUT /categories/reorder
func (h *Handler) ReorderCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Parse request body
	var req map[uuid.UUID]int
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	// Reorder categories
	if err := h.service.ReorderCategories(ctx, tenantID, req); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// BulkUpdateStatus handles PUT /categories/bulk/status
func (h *Handler) BulkUpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Parse request body
	var req struct {
		CategoryIDs []uuid.UUID   `json:"category_ids"`
		Status      CategoryStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	// Bulk update status
	if err := h.service.BulkUpdateStatus(ctx, tenantID, req.CategoryIDs, req.Status); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// AddProductToCategory handles POST /categories/{id}/products/{product_id}
func (h *Handler) AddProductToCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}
	
	// Extract product ID
	productID, err := h.extractProductID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid product ID", err)
		return
	}
	
	// Add product to category
	if err := h.service.AddProductToCategory(ctx, tenantID, categoryID, productID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// RemoveProductFromCategory handles DELETE /categories/{id}/products/{product_id}
func (h *Handler) RemoveProductFromCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}
	
	// Extract product ID
	productID, err := h.extractProductID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid product ID", err)
		return
	}
	
	// Remove product from category
	if err := h.service.RemoveProductFromCategory(ctx, tenantID, categoryID, productID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// GetCategoryProducts handles GET /categories/{id}/products
func (h *Handler) GetCategoryProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Extract category ID
	categoryID, err := h.extractCategoryID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}
	
	// Parse pagination
	limit, offset := h.parsePagination(r)
	
	// Get category products
	products, err := h.service.GetCategoryProducts(ctx, tenantID, categoryID, limit, offset)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, products)
}

// GetCategoryStats handles GET /categories/stats
func (h *Handler) GetCategoryStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Get category stats
	stats, err := h.service.GetCategoryStats(ctx, tenantID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, stats)
}

// GetFeaturedCategories handles GET /categories/featured
func (h *Handler) GetFeaturedCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Parse limit
	limit := 10 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	// Get featured categories
	categories, err := h.service.GetFeaturedCategories(ctx, tenantID, limit)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, categories)
}

// GetPopularCategories handles GET /categories/popular
func (h *Handler) GetPopularCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract tenant ID
	tenantID, err := h.extractTenantID(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid tenant ID", err)
		return
	}
	
	// Parse limit
	limit := 10 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	// Get popular categories
	categories, err := h.service.GetPopularCategories(ctx, tenantID, limit)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSONResponse(w, http.StatusOK, categories)
}

// Helper methods

// extractTenantID extracts tenant ID from request context or headers
func (h *Handler) extractTenantID(r *http.Request) (uuid.UUID, error) {
	// Try to get from context first (set by middleware)
	if tenantID := r.Context().Value("tenant_id"); tenantID != nil {
		if id, ok := tenantID.(uuid.UUID); ok {
			return id, nil
		}
	}
	
	// Try to get from header
	tenantIDStr := r.Header.Get("X-Tenant-ID")
	if tenantIDStr == "" {
		return uuid.Nil, fmt.Errorf("tenant ID not found")
	}
	
	return uuid.Parse(tenantIDStr)
}

// extractCategoryID extracts category ID from URL path
func (h *Handler) extractCategoryID(r *http.Request) (uuid.UUID, error) {
	categoryIDStr := mux.Vars(r)["id"]
	if categoryIDStr == "" {
		return uuid.Nil, fmt.Errorf("category ID not found")
	}
	
	return uuid.Parse(categoryIDStr)
}

// extractProductID extracts product ID from URL path
func (h *Handler) extractProductID(r *http.Request) (uuid.UUID, error) {
	productIDStr := mux.Vars(r)["product_id"]
	if productIDStr == "" {
		return uuid.Nil, fmt.Errorf("product ID not found")
	}
	
	return uuid.Parse(productIDStr)
}

// parseCategoryFilter parses category filter from query parameters
func (h *Handler) parseCategoryFilter(r *http.Request) CategoryFilter {
	filter := CategoryFilter{}
	
	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = name
	}
	
	if slug := r.URL.Query().Get("slug"); slug != "" {
		filter.Slug = slug
	}
	
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = CategoryStatus(status)
	}
	
	if parentIDStr := r.URL.Query().Get("parent_id"); parentIDStr != "" {
		if parentID, err := uuid.Parse(parentIDStr); err == nil {
			filter.ParentID = &parentID
		}
	}
	
	if levelStr := r.URL.Query().Get("level"); levelStr != "" {
		if level, err := strconv.Atoi(levelStr); err == nil {
			filter.Level = &level
		}
	}
	
	if featuredStr := r.URL.Query().Get("is_featured"); featuredStr != "" {
		if featured, err := strconv.ParseBool(featuredStr); err == nil {
			filter.IsFeatured = &featured
		}
	}
	
	if menuStr := r.URL.Query().Get("show_in_menu"); menuStr != "" {
		if menu, err := strconv.ParseBool(menuStr); err == nil {
			filter.ShowInMenu = &menu
		}
	}
	
	return filter
}

// parsePagination parses pagination parameters from query string
func (h *Handler) parsePagination(r *http.Request) (limit, offset int) {
	limit = 20 // default
	offset = 0 // default
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100 // max limit
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

// writeJSONResponse writes a JSON response
func (h *Handler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error but don't expose it to client
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// writeErrorResponse writes an error response
func (h *Handler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	errorResponse := map[string]interface{}{
		"error":   message,
		"status":  statusCode,
	}
	
	// Include error details in development mode
	if err != nil {
		// In production, you might want to log this instead
		errorResponse["details"] = err.Error()
	}
	
	h.writeJSONResponse(w, statusCode, errorResponse)
}

// handleServiceError handles service layer errors and maps them to HTTP status codes
func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case ErrCategoryNotFound:
		h.writeErrorResponse(w, http.StatusNotFound, "Category not found", err)
	case ErrCategoryExists:
		h.writeErrorResponse(w, http.StatusConflict, "Category already exists", err)
	case ErrInvalidParent:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid parent category", err)
	case ErrCategoryHasChildren:
		h.writeErrorResponse(w, http.StatusBadRequest, "Category has children", err)
	case ErrCategoryHasProducts:
		h.writeErrorResponse(w, http.StatusBadRequest, "Category has products", err)
	case ErrMaxDepthExceeded:
		h.writeErrorResponse(w, http.StatusBadRequest, "Maximum category depth exceeded", err)
	case ErrInvalidSlug:
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid slug", err)
	default:
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error", err)
	}
}