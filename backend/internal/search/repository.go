package search

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new search repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Search performs global search across all content types
func (r *repository) Search(ctx context.Context, tenantID uuid.UUID, query *SearchQuery) ([]*SearchResult, int64, error) {
	var results []*SearchResult
	var total int64

	// Build search query based on type
	switch strings.ToLower(query.Type) {
	case "product", "":
		// Search products
		productResults, productTotal, err := r.searchProducts(ctx, tenantID, query)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, productResults...)
		total += productTotal

		if query.Type == "" {
			// Also search categories and pages if global search
			categoryResults, categoryTotal, err := r.searchCategories(ctx, tenantID, query)
			if err == nil {
				results = append(results, categoryResults...)
				total += categoryTotal
			}
		}

	case "category":
		// Search categories
		categoryResults, categoryTotal, err := r.searchCategories(ctx, tenantID, query)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, categoryResults...)
		total += categoryTotal

	case "page":
		// Search content pages
		pageResults, pageTotal, err := r.searchPages(ctx, tenantID, query)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, pageResults...)
		total += pageTotal

	default:
		return nil, 0, fmt.Errorf("unsupported search type: %s", query.Type)
	}

	// Apply sorting and pagination
	results = r.sortResults(results, query.SortBy)
	results = r.paginateResults(results, query.Offset, query.Limit)

	return results, total, nil
}

// SearchProducts performs product-specific search with advanced filters
func (r *repository) SearchProducts(ctx context.Context, tenantID uuid.UUID, req *ProductSearchRequest) ([]*SearchResult, int64, error) {
	var results []*SearchResult

	// Build query
	query := r.db.WithContext(ctx).Table("products p").
		Select("p.id, p.name as title, p.description, p.price, p.image_url, p.created_at, p.updated_at").
		Where("p.tenant_id = ? AND p.deleted_at IS NULL", tenantID)

	// Add search conditions
	if req.Query != "" {
		searchTerm := "%" + strings.ToLower(req.Query) + "%"
		query = query.Where("(LOWER(p.name) LIKE ? OR LOWER(p.description) LIKE ? OR LOWER(p.sku) LIKE ?)", 
			searchTerm, searchTerm, searchTerm)
	}

	// Add filters
	if req.CategoryID != "" {
		query = query.Where("p.category_id = ?", req.CategoryID)
	}

	if req.BrandID != "" {
		query = query.Where("p.brand_id = ?", req.BrandID)
	}

	if req.MinPrice != nil {
		query = query.Where("p.price >= ?", *req.MinPrice)
	}

	if req.MaxPrice != nil {
		query = query.Where("p.price <= ?", *req.MaxPrice)
	}

	if req.InStock != nil && *req.InStock {
		query = query.Where("p.stock_quantity > 0")
	}

	if req.OnSale != nil && *req.OnSale {
		query = query.Where("p.sale_price IS NOT NULL AND p.sale_price < p.price")
	}

	if req.Rating != nil {
		query = query.Where("p.average_rating >= ?", *req.Rating)
	}

	// Add tag filters
	if len(req.Tags) > 0 {
		query = query.Joins("JOIN product_tags pt ON p.id = pt.product_id").
			Joins("JOIN tags t ON pt.tag_id = t.id").
			Where("t.name IN ?", req.Tags)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Add sorting
	switch req.SortBy {
	case "price_asc":
		query = query.Order("p.price ASC")
	case "price_desc":
		query = query.Order("p.price DESC")
	case "newest":
		query = query.Order("p.created_at DESC")
	case "rating":
		query = query.Order("p.average_rating DESC")
	default:
		// Relevance sorting (simplified)
		query = query.Order("p.name ASC")
	}

	// Add pagination
	query = query.Offset(req.Offset).Limit(req.Limit)

	// Execute query
	var products []struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Price       float64   `json:"price"`
		ImageURL    string    `json:"image_url"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	if err := query.Scan(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	// Convert to search results
	for _, product := range products {
		result := &SearchResult{
			ID:          product.ID,
			Type:        "product",
			Title:       product.Title,
			Description: product.Description,
			URL:         fmt.Sprintf("/products/%s", product.ID),
			ImageURL:    product.ImageURL,
			Price:       &product.Price,
			Score:       1.0, // Simplified scoring
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
		results = append(results, result)
	}

	return results, total, nil
}

// GetSuggestions returns search suggestions/autocomplete
func (r *repository) GetSuggestions(ctx context.Context, tenantID uuid.UUID, query string, searchType string, limit int) ([]Suggestion, error) {
	var suggestions []Suggestion

	searchTerm := strings.ToLower(query) + "%"

	switch searchType {
	case "product", "":
		// Get product suggestions
		var productNames []string
		err := r.db.WithContext(ctx).Table("products").
			Select("DISTINCT name").
			Where("tenant_id = ? AND deleted_at IS NULL AND LOWER(name) LIKE ?", tenantID, searchTerm).
			Limit(limit).
			Pluck("name", &productNames)
		if err != nil {
			return nil, fmt.Errorf("failed to get product suggestions: %w", err)
		}

		for i, name := range productNames {
			suggestions = append(suggestions, Suggestion{
				Text:  name,
				Type:  "product",
				Score: float64(limit - i), // Simple scoring
			})
		}

	case "category":
		// Get category suggestions
		var categoryNames []string
		err := r.db.WithContext(ctx).Table("categories").
			Select("DISTINCT name").
			Where("tenant_id = ? AND deleted_at IS NULL AND LOWER(name) LIKE ?", tenantID, searchTerm).
			Limit(limit).
			Pluck("name", &categoryNames)
		if err != nil {
			return nil, fmt.Errorf("failed to get category suggestions: %w", err)
		}

		for i, name := range categoryNames {
			suggestions = append(suggestions, Suggestion{
				Text:  name,
				Type:  "category",
				Score: float64(limit - i),
			})
		}
	}

	return suggestions, nil
}

// LogSearch logs search activity
func (r *repository) LogSearch(ctx context.Context, log *SearchLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetSearchAnalytics returns search analytics and metrics
func (r *repository) GetSearchAnalytics(ctx context.Context, tenantID uuid.UUID, req *SearchAnalyticsRequest) (*SearchAnalyticsResponse, error) {
	response := &SearchAnalyticsResponse{
		Period: fmt.Sprintf("%s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02")),
	}

	// Get total searches
	var totalSearches int64
	err := r.db.WithContext(ctx).Model(&SearchLog{}).
		Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, req.StartDate, req.EndDate).
		Count(&totalSearches).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get total searches: %w", err)
	}
	response.TotalSearches = totalSearches

	// Get top queries
	var topQueries []QueryStat
	err = r.db.WithContext(ctx).Table("search_logs").
		Select("query, COUNT(*) as count, AVG(results) as result_count").
		Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, req.StartDate, req.EndDate).
		Group("query").
		Order("count DESC").
		Limit(req.Limit).
		Scan(&topQueries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get top queries: %w", err)
	}
	response.TopQueries = topQueries

	// Get no results queries
	var noResults []QueryStat
	err = r.db.WithContext(ctx).Table("search_logs").
		Select("query, COUNT(*) as count").
		Where("tenant_id = ? AND created_at BETWEEN ? AND ? AND results = 0", tenantID, req.StartDate, req.EndDate).
		Group("query").
		Order("count DESC").
		Limit(req.Limit).
		Scan(&noResults).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get no results queries: %w", err)
	}
	response.NoResults = noResults

	// Calculate metrics
	var uniqueQueries int64
	err = r.db.WithContext(ctx).Table("search_logs").
		Select("COUNT(DISTINCT query)").
		Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, req.StartDate, req.EndDate).
		Scan(&uniqueQueries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get unique queries: %w", err)
	}

	var avgResults float64
	err = r.db.WithContext(ctx).Table("search_logs").
		Select("AVG(results)").
		Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, req.StartDate, req.EndDate).
		Scan(&avgResults).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get average results: %w", err)
	}

	var noResultsCount int64
	err = r.db.WithContext(ctx).Model(&SearchLog{}).
		Where("tenant_id = ? AND created_at BETWEEN ? AND ? AND results = 0", tenantID, req.StartDate, req.EndDate).
		Count(&noResultsCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get no results count: %w", err)
	}

	noResultsRate := float64(0)
	if totalSearches > 0 {
		noResultsRate = float64(noResultsCount) / float64(totalSearches) * 100
	}

	response.Metrics = SearchMetrics{
		TotalSearches:    totalSearches,
		UniqueQueries:    uniqueQueries,
		AverageResults:   avgResults,
		NoResultsRate:    noResultsRate,
		ClickThroughRate: 0, // Would need click tracking
		AveragePosition:  0, // Would need position tracking
	}

	return response, nil
}

// GetFilters returns available search filters
func (r *repository) GetFilters(ctx context.Context, tenantID uuid.UUID, searchType string) ([]Filter, error) {
	var filters []Filter

	// This is a simplified implementation
	// In a real system, you'd have a filters table
	switch searchType {
	case "product":
		// Category filter
		var categories []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		err := r.db.WithContext(ctx).Table("categories").
			Select("id, name").
			Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
			Scan(&categories).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get categories: %w", err)
		}

		var categoryValues []FilterValue
		for _, cat := range categories {
			categoryValues = append(categoryValues, FilterValue{
				Value: cat.ID,
				Label: cat.Name,
				Count: 0, // Would need to count products
			})
		}

		filters = append(filters, Filter{
			ID:       "category",
			Type:     "category",
			Name:     "Category",
			Values:   categoryValues,
			IsActive: true,
		})

		// Price range filter
		filters = append(filters, Filter{
			ID:   "price",
			Type: "range",
			Name: "Price Range",
			Values: []FilterValue{
				{Value: "0-50", Label: "Under $50", Count: 0},
				{Value: "50-100", Label: "$50 - $100", Count: 0},
				{Value: "100-200", Label: "$100 - $200", Count: 0},
				{Value: "200+", Label: "Over $200", Count: 0},
			},
			IsActive: true,
		})
	}

	return filters, nil
}

// ManageFilter manages search filters (create, update, delete)
func (r *repository) ManageFilter(ctx context.Context, tenantID uuid.UUID, req *FilterRequest) error {
	// This would implement filter management
	// For now, return success as filters are hardcoded
	return nil
}

// Helper methods

func (r *repository) searchProducts(ctx context.Context, tenantID uuid.UUID, query *SearchQuery) ([]*SearchResult, int64, error) {
	// Convert to ProductSearchRequest
	productReq := &ProductSearchRequest{
		Query:  query.Query,
		Offset: query.Offset,
		Limit:  query.Limit,
		SortBy: query.SortBy,
	}
	return r.SearchProducts(ctx, tenantID, productReq)
}

func (r *repository) searchCategories(ctx context.Context, tenantID uuid.UUID, query *SearchQuery) ([]*SearchResult, int64, error) {
	var results []*SearchResult

	searchTerm := "%" + strings.ToLower(query.Query) + "%"

	var categories []struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	err := r.db.WithContext(ctx).Table("categories").
		Select("id, name, description, created_at, updated_at").
		Where("tenant_id = ? AND deleted_at IS NULL AND (LOWER(name) LIKE ? OR LOWER(description) LIKE ?)", 
			tenantID, searchTerm, searchTerm).
		Offset(query.Offset).Limit(query.Limit).
		Scan(&categories).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search categories: %w", err)
	}

	for _, category := range categories {
		result := &SearchResult{
			ID:          category.ID,
			Type:        "category",
			Title:       category.Name,
			Description: category.Description,
			URL:         fmt.Sprintf("/categories/%s", category.ID),
			Score:       1.0,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		}
		results = append(results, result)
	}

	return results, int64(len(results)), nil
}

func (r *repository) searchPages(ctx context.Context, tenantID uuid.UUID, query *SearchQuery) ([]*SearchResult, int64, error) {
	// This would search content pages
	// For now, return empty results
	return []*SearchResult{}, 0, nil
}

func (r *repository) sortResults(results []*SearchResult, sortBy string) []*SearchResult {
	// Implement sorting logic
	// For now, return as-is
	return results
}

func (r *repository) paginateResults(results []*SearchResult, offset, limit int) []*SearchResult {
	if offset >= len(results) {
		return []*SearchResult{}
	}

	end := offset + limit
	if end > len(results) {
		end = len(results)
	}

	return results[offset:end]
}