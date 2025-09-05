package search

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler handles search-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new search handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers search routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	search := router.Group("/search")
	{
		// Global search
		search.GET("", h.Search)                    // GET /search
		// Product search
		search.GET("/products", h.SearchProducts)   // GET /search/products
		// Search suggestions
		search.GET("/suggestions", h.GetSuggestions) // GET /search/suggestions
		// Search analytics
		search.GET("/analytics", h.GetAnalytics)    // GET /search/analytics
		// Search filters
		search.GET("/filters", h.GetFilters)        // GET /search/filters
		search.POST("/filters", h.ManageFilters)    // POST /search/filters
	}
}

// Search performs global search across all content types
// GET /search?q=query&type=product&categories=cat1,cat2&sort_by=relevance&offset=0&limit=20
func (h *Handler) Search(c *gin.Context) {
	// Parse query parameters
	req := &SearchQuery{
		Query:  c.Query("q"),
		Type:   c.Query("type"),
		SortBy: c.Query("sort_by"),
	}

	// Parse categories
	if categories := c.Query("categories"); categories != "" {
		req.Categories = strings.Split(categories, ",")
	}

	// Parse tags
	if tags := c.Query("tags"); tags != "" {
		req.Tags = strings.Split(tags, ",")
	}

	// Parse price filters
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			req.MinPrice = &price
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			req.MaxPrice = &price
		}
	}

	// Parse in_stock filter
	if inStock := c.Query("in_stock"); inStock != "" {
		if stock, err := strconv.ParseBool(inStock); err == nil {
			req.InStock = &stock
		}
	}

	// Parse pagination
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			req.Offset = o
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			req.Limit = l
		}
	}

	// Validate required fields
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	// Perform search
	response, err := h.service.Search(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SearchProducts performs product-specific search with advanced filters
// GET /search/products?q=query&category_id=uuid&brand_id=uuid&tags=tag1,tag2&min_price=10&max_price=100&in_stock=true&on_sale=true&min_rating=4&sort_by=price_asc&offset=0&limit=20&include_facets=true
func (h *Handler) SearchProducts(c *gin.Context) {
	// Parse query parameters
	req := &ProductSearchRequest{
		Query:      c.Query("q"),
		CategoryID: c.Query("category_id"),
		BrandID:    c.Query("brand_id"),
		SortBy:     c.Query("sort_by"),
	}

	// Parse tags
	if tags := c.Query("tags"); tags != "" {
		req.Tags = strings.Split(tags, ",")
	}

	// Parse price filters
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			req.MinPrice = &price
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			req.MaxPrice = &price
		}
	}

	// Parse boolean filters
	if inStock := c.Query("in_stock"); inStock != "" {
		if stock, err := strconv.ParseBool(inStock); err == nil {
			req.InStock = &stock
		}
	}
	if onSale := c.Query("on_sale"); onSale != "" {
		if sale, err := strconv.ParseBool(onSale); err == nil {
			req.OnSale = &sale
		}
	}
	if includeFacets := c.Query("include_facets"); includeFacets != "" {
		if facets, err := strconv.ParseBool(includeFacets); err == nil {
			req.IncludeFacets = facets
		}
	}

	// Parse rating filter
	if rating := c.Query("min_rating"); rating != "" {
		if r, err := strconv.ParseFloat(rating, 64); err == nil {
			req.Rating = &r
		}
	}

	// Parse pagination
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			req.Offset = o
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			req.Limit = l
		}
	}

	// Validate required fields
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	// Perform product search
	response, err := h.service.SearchProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSuggestions returns search suggestions/autocomplete
// GET /search/suggestions?q=query&type=product&limit=10
func (h *Handler) GetSuggestions(c *gin.Context) {
	// Parse query parameters
	req := &SuggestionRequest{
		Query: c.Query("q"),
		Type:  c.Query("type"),
	}

	// Parse limit
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			req.Limit = l
		}
	}

	// Validate required fields
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	// Get suggestions
	response, err := h.service.GetSuggestions(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAnalytics returns search analytics and metrics
// GET /search/analytics?type=queries&start_date=2024-01-01&end_date=2024-01-31&limit=50
func (h *Handler) GetAnalytics(c *gin.Context) {
	// Check if user has admin access
	userRole, exists := c.Get("user_role")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Parse query parameters
	req := &SearchAnalyticsRequest{
		Type: c.Query("type"),
	}

	// Parse dates
	if startDate := c.Query("start_date"); startDate != "" {
		if date, err := parseDate(startDate); err == nil {
			req.StartDate = &date
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if date, err := parseDate(endDate); err == nil {
			req.EndDate = &date
		}
	}

	// Parse limit
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			req.Limit = l
		}
	}

	// Get analytics
	response, err := h.service.GetAnalytics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetFilters returns available search filters
// GET /search/filters?type=product
func (h *Handler) GetFilters(c *gin.Context) {
	searchType := c.Query("type")
	if searchType == "" {
		searchType = "product" // Default to product filters
	}

	// Get filters
	response, err := h.service.GetFilters(c.Request.Context(), searchType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ManageFilters manages search filters (create, update, delete)
// POST /search/filters
func (h *Handler) ManageFilters(c *gin.Context) {
	// Check if user has admin access
	userRole, exists := c.Get("user_role")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Parse request body
	var req FilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Manage filters
	response, err := h.service.ManageFilters(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Helper functions

func parseDate(dateStr string) (time.Time, error) {
	// Try different date formats
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}