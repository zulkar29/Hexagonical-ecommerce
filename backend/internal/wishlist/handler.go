package wishlist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler handles HTTP requests for wishlist operations
type Handler struct {
	service Service
}

// NewHandler creates a new wishlist handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers wishlist routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Wishlist routes
	router.HandleFunc("/wishlists", h.CreateWishlist).Methods("POST")
	router.HandleFunc("/wishlists", h.ListWishlists).Methods("GET")
	router.HandleFunc("/wishlists/{wishlistID}", h.GetWishlist).Methods("GET")
	router.HandleFunc("/wishlists/{wishlistID}", h.UpdateWishlist).Methods("PUT")
	router.HandleFunc("/wishlists/{wishlistID}", h.DeleteWishlist).Methods("DELETE")
	router.HandleFunc("/wishlists/share/{shareToken}", h.GetWishlistByShareToken).Methods("GET")
	router.HandleFunc("/wishlists/{wishlistID}/share", h.ShareWishlist).Methods("POST")
	
	// Customer wishlist routes
	router.HandleFunc("/customers/{customerID}/wishlists", h.GetCustomerWishlists).Methods("GET")
	router.HandleFunc("/customers/{customerID}/wishlists/default", h.GetDefaultWishlist).Methods("GET")
	router.HandleFunc("/customers/{customerID}/wishlists/{wishlistID}/default", h.SetDefaultWishlist).Methods("POST")
	
	// Wishlist item routes
	router.HandleFunc("/wishlists/{wishlistID}/items", h.AddItem).Methods("POST")
	router.HandleFunc("/wishlists/{wishlistID}/items", h.ListWishlistItems).Methods("GET")
	router.HandleFunc("/wishlists/{wishlistID}/items/bulk", h.BulkAddItems).Methods("POST")
	router.HandleFunc("/wishlists/{wishlistID}/items/reorder", h.ReorderItems).Methods("POST")
	router.HandleFunc("/wishlists/{wishlistID}/clear", h.ClearWishlist).Methods("POST")
	
	router.HandleFunc("/wishlist-items", h.ListItems).Methods("GET")
	router.HandleFunc("/wishlist-items/{itemID}", h.GetItem).Methods("GET")
	router.HandleFunc("/wishlist-items/{itemID}", h.UpdateItem).Methods("PUT")
	router.HandleFunc("/wishlist-items/{itemID}", h.RemoveItem).Methods("DELETE")
	router.HandleFunc("/wishlist-items/{itemID}/move", h.MoveItem).Methods("POST")
	router.HandleFunc("/wishlist-items/{itemID}/copy", h.CopyItem).Methods("POST")
	
	// Bulk operations
	router.HandleFunc("/wishlist-items/bulk", h.BulkRemoveItems).Methods("DELETE")
	router.HandleFunc("/wishlist-items/bulk/priority", h.BulkUpdateItemPriority).Methods("PUT")
	
	// Wishlist management
	router.HandleFunc("/wishlists/{sourceID}/merge/{targetID}", h.MergeWishlists).Methods("POST")
	
	// Analytics and statistics
	router.HandleFunc("/wishlists/stats", h.GetWishlistStats).Methods("GET")
	router.HandleFunc("/wishlists/popular", h.GetPopularWishlists).Methods("GET")
	router.HandleFunc("/products/most-wished", h.GetMostWishedProducts).Methods("GET")
	router.HandleFunc("/customers/{customerID}/wishlist-activity", h.GetCustomerActivity).Methods("GET")
	
	// Maintenance operations
	router.HandleFunc("/wishlists/cleanup/empty", h.CleanupEmptyWishlists).Methods("POST")
	router.HandleFunc("/wishlist-items/cleanup/orphaned", h.CleanupOrphanedItems).Methods("POST")
}

// Wishlist operations

// CreateWishlist creates a new wishlist
func (h *Handler) CreateWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	var req CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	req.TenantID = tenantID
	
	response, err := h.service.CreateWishlist(ctx, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusCreated, response)
}

// GetWishlist retrieves a wishlist by ID
func (h *Handler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	wishlistID, err := h.getUUIDParam(r, "wishlistID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}
	
	response, err := h.service.GetWishlist(ctx, tenantID, wishlistID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// GetWishlistByShareToken retrieves a public wishlist by share token
func (h *Handler) GetWishlistByShareToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	shareToken := mux.Vars(r)["shareToken"]
	if shareToken == "" {
		h.writeError(w, http.StatusBadRequest, "Share token is required", nil)
		return
	}
	
	response, err := h.service.GetWishlistByShareToken(ctx, shareToken)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// UpdateWishlist updates an existing wishlist
func (h *Handler) UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	wishlistID, err := h.getUUIDParam(r, "wishlistID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}
	
	var req UpdateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	response, err := h.service.UpdateWishlist(ctx, tenantID, wishlistID, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// DeleteWishlist deletes a wishlist
func (h *Handler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	wishlistID, err := h.getUUIDParam(r, "wishlistID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}
	
	if err := h.service.DeleteWishlist(ctx, tenantID, wishlistID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// ListWishlists returns paginated wishlists
func (h *Handler) ListWishlists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	// Parse query parameters
	filter := h.parseWishlistFilter(r)
	limit, offset := h.parsePagination(r)
	
	wishlists, total, err := h.service.ListWishlists(ctx, tenantID, filter, limit, offset)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	response := map[string]interface{}{
		"wishlists": wishlists,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// ShareWishlist makes a wishlist public/private
func (h *Handler) ShareWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	wishlistID, err := h.getUUIDParam(r, "wishlistID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}
	
	var req struct {
		IsPublic bool `json:"is_public"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	shareURL, err := h.service.ShareWishlist(ctx, tenantID, wishlistID, req.IsPublic)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	response := map[string]interface{}{
		"is_public": req.IsPublic,
		"share_url": shareURL,
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// Customer wishlist operations

// GetCustomerWishlists returns all wishlists for a customer
func (h *Handler) GetCustomerWishlists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	customerID, err := h.getUUIDParam(r, "customerID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid customer ID", err)
		return
	}
	
	wishlists, err := h.service.GetCustomerWishlists(ctx, tenantID, customerID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, wishlists)
}

// GetDefaultWishlist returns the default wishlist for a customer
func (h *Handler) GetDefaultWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	customerID, err := h.getUUIDParam(r, "customerID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid customer ID", err)
		return
	}
	
	wishlist, err := h.service.GetDefaultWishlist(ctx, tenantID, customerID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, wishlist)
}

// SetDefaultWishlist sets a wishlist as default for a customer
func (h *Handler) SetDefaultWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	customerID, err := h.getUUIDParam(r, "customerID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid customer ID", err)
		return
	}
	
	wishlistID, err := h.getUUIDParam(r, "wishlistID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}
	
	if err := h.service.SetDefaultWishlist(ctx, tenantID, customerID, wishlistID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// Wishlist item operations

// AddItem adds an item to a wishlist
func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	wishlistID, err := h.getUUIDParam(r, "wishlistID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}
	
	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	req.WishlistID = wishlistID
	
	response, err := h.service.AddItem(ctx, tenantID, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusCreated, response)
}

// GetItem retrieves a wishlist item by ID
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	itemID, err := h.getUUIDParam(r, "itemID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}
	
	response, err := h.service.GetItem(ctx, tenantID, itemID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// UpdateItem updates a wishlist item
func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	itemID, err := h.getUUIDParam(r, "itemID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}
	
	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	response, err := h.service.UpdateItem(ctx, tenantID, itemID, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// RemoveItem removes an item from a wishlist
func (h *Handler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	itemID, err := h.getUUIDParam(r, "itemID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}
	
	if err := h.service.RemoveItem(ctx, tenantID, itemID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// ListItems returns paginated wishlist items
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	// Parse query parameters
	filter := h.parseWishlistItemFilter(r)
	limit, offset := h.parsePagination(r)
	
	items, total, err := h.service.ListItems(ctx, tenantID, filter, limit, offset)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	response := map[string]interface{}{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// ListWishlistItems returns items for a specific wishlist
func (h *Handler) ListWishlistItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	wishlistID, err := h.getUUIDParam(r, "wishlistID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}
	
	// Parse query parameters
	filter := h.parseWishlistItemFilter(r)
	filter.WishlistID = &wishlistID
	limit, offset := h.parsePagination(r)
	
	items, total, err := h.service.ListItems(ctx, tenantID, filter, limit, offset)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	response := map[string]interface{}{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// MoveItem moves an item to another wishlist
func (h *Handler) MoveItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	itemID, err := h.getUUIDParam(r, "itemID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}
	
	var req struct {
		TargetWishlistID uuid.UUID `json:"target_wishlist_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	if err := h.service.MoveItem(ctx, tenantID, itemID, req.TargetWishlistID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// CopyItem copies an item to another wishlist
func (h *Handler) CopyItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	itemID, err := h.getUUIDParam(r, "itemID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}
	
	var req struct {
		TargetWishlistID uuid.UUID `json:"target_wishlist_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	if err := h.service.CopyItem(ctx, tenantID, itemID, req.TargetWishlistID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// ClearWishlist removes all items from a wishlist
func (h *Handler) ClearWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	wishlistID, err := h.getUUIDParam(r, "wishlistID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}
	
	if err := h.service.ClearWishlist(ctx, tenantID, wishlistID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// Bulk operations

// BulkAddItems adds multiple items to a wishlist
func (h *Handler) BulkAddItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	wishlistID, err := h.getUUIDParam(r, "wishlistID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}
	
	var req struct {
		Items []AddItemRequest `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	// Set wishlist ID for all items
	for i := range req.Items {
		req.Items[i].WishlistID = wishlistID
	}
	
	responses, err := h.service.BulkAddItems(ctx, tenantID, req.Items)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusCreated, responses)
}

// BulkRemoveItems removes multiple items from wishlists
func (h *Handler) BulkRemoveItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	var req struct {
		ItemIDs []uuid.UUID `json:"item_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	if err := h.service.BulkRemoveItems(ctx, tenantID, req.ItemIDs); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// BulkUpdateItemPriority updates priority for multiple items
func (h *Handler) BulkUpdateItemPriority(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	var req struct {
		Updates map[uuid.UUID]int `json:"updates"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	if err := h.service.BulkUpdateItemPriority(ctx, tenantID, req.Updates); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// ReorderItems reorders items in a wishlist
func (h *Handler) ReorderItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	wishlistID, err := h.getUUIDParam(r, "wishlistID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}
	
	var req struct {
		ItemOrder []uuid.UUID `json:"item_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	if err := h.service.ReorderItems(ctx, tenantID, wishlistID, req.ItemOrder); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// Wishlist management

// MergeWishlists merges source wishlist into target wishlist
func (h *Handler) MergeWishlists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	sourceID, err := h.getUUIDParam(r, "sourceID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid source wishlist ID", err)
		return
	}
	
	targetID, err := h.getUUIDParam(r, "targetID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid target wishlist ID", err)
		return
	}
	
	if err := h.service.MergeWishlists(ctx, tenantID, sourceID, targetID); err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// Analytics and statistics

// GetWishlistStats returns wishlist statistics
func (h *Handler) GetWishlistStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	stats, err := h.service.GetWishlistStats(ctx, tenantID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, stats)
}

// GetMostWishedProducts returns the most wished products
func (h *Handler) GetMostWishedProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	limit := h.getIntParam(r, "limit", 10)
	
	products, err := h.service.GetMostWishedProducts(ctx, tenantID, limit)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, products)
}

// GetCustomerActivity returns customer wishlist activity
func (h *Handler) GetCustomerActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	customerID, err := h.getUUIDParam(r, "customerID")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid customer ID", err)
		return
	}
	
	days := h.getIntParam(r, "days", 30)
	
	activity, err := h.service.GetCustomerActivity(ctx, tenantID, customerID, days)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, activity)
}

// GetPopularWishlists returns popular public wishlists
func (h *Handler) GetPopularWishlists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	limit := h.getIntParam(r, "limit", 10)
	
	wishlists, err := h.service.GetPopularWishlists(ctx, tenantID, limit)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.writeJSON(w, http.StatusOK, wishlists)
}

// Maintenance operations

// CleanupEmptyWishlists removes empty wishlists older than specified days
func (h *Handler) CleanupEmptyWishlists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	days := h.getIntParam(r, "days", 30)
	
	count, err := h.service.CleanupEmptyWishlists(ctx, tenantID, days)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	response := map[string]interface{}{
		"deleted_count": count,
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// CleanupOrphanedItems removes items for non-existent products
func (h *Handler) CleanupOrphanedItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := h.getTenantID(r)
	
	count, err := h.service.CleanupOrphanedItems(ctx, tenantID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	response := map[string]interface{}{
		"deleted_count": count,
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

// Helper methods

// getTenantID extracts tenant ID from request context or headers
func (h *Handler) getTenantID(r *http.Request) uuid.UUID {
	// This would typically come from JWT token or request context
	// For now, return a placeholder
	return uuid.New()
}

// getUUIDParam extracts UUID parameter from URL
func (h *Handler) getUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	value := mux.Vars(r)[param]
	return uuid.Parse(value)
}

// getIntParam extracts integer parameter from query string
func (h *Handler) getIntParam(r *http.Request, param string, defaultValue int) int {
	value := r.URL.Query().Get(param)
	if value == "" {
		return defaultValue
	}
	
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	
	return defaultValue
}

// parsePagination extracts pagination parameters
func (h *Handler) parsePagination(r *http.Request) (limit, offset int) {
	limit = h.getIntParam(r, "limit", 20)
	offset = h.getIntParam(r, "offset", 0)
	
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	
	return limit, offset
}

// parseWishlistFilter extracts wishlist filter parameters
func (h *Handler) parseWishlistFilter(r *http.Request) WishlistFilter {
	filter := WishlistFilter{}
	
	if customerIDStr := r.URL.Query().Get("customer_id"); customerIDStr != "" {
		if customerID, err := uuid.Parse(customerIDStr); err == nil {
			filter.CustomerID = &customerID
		}
	}
	
	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = name
	}
	
	if isDefaultStr := r.URL.Query().Get("is_default"); isDefaultStr != "" {
		if isDefault, err := strconv.ParseBool(isDefaultStr); err == nil {
			filter.IsDefault = &isDefault
		}
	}
	
	if isPublicStr := r.URL.Query().Get("is_public"); isPublicStr != "" {
		if isPublic, err := strconv.ParseBool(isPublicStr); err == nil {
			filter.IsPublic = &isPublic
		}
	}
	
	if isEmptyStr := r.URL.Query().Get("is_empty"); isEmptyStr != "" {
		if isEmpty, err := strconv.ParseBool(isEmptyStr); err == nil {
			filter.IsEmpty = &isEmpty
		}
	}
	
	return filter
}

// parseWishlistItemFilter extracts wishlist item filter parameters
func (h *Handler) parseWishlistItemFilter(r *http.Request) WishlistItemFilter {
	filter := WishlistItemFilter{}
	
	if wishlistIDStr := r.URL.Query().Get("wishlist_id"); wishlistIDStr != "" {
		if wishlistID, err := uuid.Parse(wishlistIDStr); err == nil {
			filter.WishlistID = &wishlistID
		}
	}
	
	if productIDStr := r.URL.Query().Get("product_id"); productIDStr != "" {
		if productID, err := uuid.Parse(productIDStr); err == nil {
			filter.ProductID = &productID
		}
	}
	
	if variantIDStr := r.URL.Query().Get("variant_id"); variantIDStr != "" {
		if variantID, err := uuid.Parse(variantIDStr); err == nil {
			filter.VariantID = &variantID
		}
	}
	
	if minPriorityStr := r.URL.Query().Get("min_priority"); minPriorityStr != "" {
		if minPriority, err := strconv.Atoi(minPriorityStr); err == nil {
			filter.MinPriority = &minPriority
		}
	}
	
	if maxPriorityStr := r.URL.Query().Get("max_priority"); maxPriorityStr != "" {
		if maxPriority, err := strconv.Atoi(maxPriorityStr); err == nil {
			filter.MaxPriority = &maxPriority
		}
	}
	
	return filter
}

// writeJSON writes JSON response
func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes error response
func (h *Handler) writeError(w http.ResponseWriter, status int, message string, err error) {
	response := map[string]interface{}{
		"error": message,
	}
	
	if err != nil {
		response["details"] = err.Error()
	}
	
	h.writeJSON(w, status, response)
}

// handleServiceError handles service layer errors
func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case ErrWishlistNotFound:
		h.writeError(w, http.StatusNotFound, "Wishlist not found", err)
	case ErrWishlistItemNotFound:
		h.writeError(w, http.StatusNotFound, "Wishlist item not found", err)
	case ErrWishlistNameExists:
		h.writeError(w, http.StatusConflict, "Wishlist name already exists", err)
	case ErrWishlistLimitExceeded:
		h.writeError(w, http.StatusBadRequest, "Wishlist limit exceeded", err)
	case ErrWishlistFull:
		h.writeError(w, http.StatusBadRequest, "Wishlist is full", err)
	case ErrCannotDeleteDefaultWishlist:
		h.writeError(w, http.StatusBadRequest, "Cannot delete default wishlist", err)
	case ErrInvalidTenantID, ErrInvalidCustomerID, ErrInvalidWishlistID, ErrInvalidProductID:
		h.writeError(w, http.StatusBadRequest, "Invalid ID provided", err)
	case ErrInvalidWishlistName, ErrWishlistNameTooLong, ErrWishlistDescriptionTooLong:
		h.writeError(w, http.StatusBadRequest, "Invalid wishlist data", err)
	case ErrInvalidQuantity, ErrInvalidPriority, ErrItemNotesTooLong:
		h.writeError(w, http.StatusBadRequest, "Invalid item data", err)
	default:
		h.writeError(w, http.StatusInternalServerError, "Internal server error", err)
	}
}