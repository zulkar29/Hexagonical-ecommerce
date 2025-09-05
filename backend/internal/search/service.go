package search

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service defines the search service interface
type Service interface {
	// Global search across all content types
	Search(ctx context.Context, req *SearchQuery) (*SearchResponse, error)
	
	// Product-specific search with advanced filters
	SearchProducts(ctx context.Context, req *ProductSearchRequest) (*SearchResponse, error)
	
	// Get search suggestions/autocomplete
	GetSuggestions(ctx context.Context, req *SuggestionRequest) (*SuggestionResponse, error)
	
	// Get search analytics and metrics
	GetAnalytics(ctx context.Context, req *SearchAnalyticsRequest) (*SearchAnalyticsResponse, error)
	
	// Manage search filters
	ManageFilters(ctx context.Context, req *FilterRequest) (*FilterResponse, error)
	
	// Get available filters
	GetFilters(ctx context.Context, searchType string) (*FilterResponse, error)
}

// Repository defines the search repository interface
type Repository interface {
	Search(ctx context.Context, tenantID uuid.UUID, query *SearchQuery) ([]*SearchResult, int64, error)
	SearchProducts(ctx context.Context, tenantID uuid.UUID, req *ProductSearchRequest) ([]*SearchResult, int64, error)
	GetSuggestions(ctx context.Context, tenantID uuid.UUID, query string, searchType string, limit int) ([]Suggestion, error)
	LogSearch(ctx context.Context, log *SearchLog) error
	GetSearchAnalytics(ctx context.Context, tenantID uuid.UUID, req *SearchAnalyticsRequest) (*SearchAnalyticsResponse, error)
	GetFilters(ctx context.Context, tenantID uuid.UUID, searchType string) ([]Filter, error)
	ManageFilter(ctx context.Context, tenantID uuid.UUID, req *FilterRequest) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new search service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Search performs global search across all content types
func (s *service) Search(ctx context.Context, req *SearchQuery) (*SearchResponse, error) {
	// Get tenant ID from context
	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("tenant ID not found in context")
	}

	// Validate request
	if err := s.validateSearchQuery(req); err != nil {
		return nil, err
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "relevance"
	}

	// Perform search
	results, total, err := s.repo.Search(ctx, tenantID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	// Log search activity
	userID := s.getUserIDFromContext(ctx)
	sessionID := s.getSessionIDFromContext(ctx)
	log := &SearchLog{
		TenantID:  tenantID,
		UserID:    userID,
		SessionID: sessionID,
		Query:     req.Query,
		Type:      req.Type,
		Results:   total,
		CreatedAt: time.Now(),
	}
	if err := s.repo.LogSearch(ctx, log); err != nil {
		// Log error but don't fail the search
		fmt.Printf("Failed to log search: %v", err)
	}

	// Get suggestions if no results
	var suggestions []string
	if total == 0 {
		suggestionResults, err := s.repo.GetSuggestions(ctx, tenantID, req.Query, req.Type, 5)
		if err == nil {
			for _, suggestion := range suggestionResults {
				suggestions = append(suggestions, suggestion.Text)
			}
		}
	}

	return &SearchResponse{
		Results:     results,
		Total:       total,
		Offset:      req.Offset,
		Limit:       req.Limit,
		Query:       req.Query,
		Suggestions: suggestions,
	}, nil
}

// SearchProducts performs product-specific search with advanced filters
func (s *service) SearchProducts(ctx context.Context, req *ProductSearchRequest) (*SearchResponse, error) {
	// Get tenant ID from context
	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("tenant ID not found in context")
	}

	// Validate request
	if err := s.validateProductSearchRequest(req); err != nil {
		return nil, err
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "relevance"
	}

	// Perform product search
	results, total, err := s.repo.SearchProducts(ctx, tenantID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	// Log search activity
	userID := s.getUserIDFromContext(ctx)
	sessionID := s.getSessionIDFromContext(ctx)
	log := &SearchLog{
		TenantID:  tenantID,
		UserID:    userID,
		SessionID: sessionID,
		Query:     req.Query,
		Type:      "product",
		Results:   total,
		CreatedAt: time.Now(),
	}
	if err := s.repo.LogSearch(ctx, log); err != nil {
		fmt.Printf("Failed to log search: %v", err)
	}

	// Get facets if requested
	var facets map[string]interface{}
	if req.IncludeFacets {
		facets = make(map[string]interface{})
		// This would be implemented based on the search results
		// For now, return empty facets
	}

	return &SearchResponse{
		Results: results,
		Total:   total,
		Offset:  req.Offset,
		Limit:   req.Limit,
		Query:   req.Query,
		Facets:  facets,
	}, nil
}

// GetSuggestions returns search suggestions/autocomplete
func (s *service) GetSuggestions(ctx context.Context, req *SuggestionRequest) (*SuggestionResponse, error) {
	// Get tenant ID from context
	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("tenant ID not found in context")
	}

	// Validate request
	if req.Query == "" {
		return nil, fmt.Errorf("query is required")
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 10
	}

	// Get suggestions
	suggestions, err := s.repo.GetSuggestions(ctx, tenantID, req.Query, req.Type, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions: %w", err)
	}

	return &SuggestionResponse{
		Suggestions: suggestions,
		Query:       req.Query,
	}, nil
}

// GetAnalytics returns search analytics and metrics
func (s *service) GetAnalytics(ctx context.Context, req *SearchAnalyticsRequest) (*SearchAnalyticsResponse, error) {
	// Get tenant ID from context
	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("tenant ID not found in context")
	}

	// Set defaults
	if req.StartDate == nil {
		start := time.Now().AddDate(0, 0, -30) // Last 30 days
		req.StartDate = &start
	}
	if req.EndDate == nil {
		end := time.Now()
		req.EndDate = &end
	}
	if req.Limit == 0 {
		req.Limit = 50
	}

	// Get analytics
	analytics, err := s.repo.GetSearchAnalytics(ctx, tenantID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get search analytics: %w", err)
	}

	return analytics, nil
}

// ManageFilters manages search filters (create, update, delete)
func (s *service) ManageFilters(ctx context.Context, req *FilterRequest) (*FilterResponse, error) {
	// Get tenant ID from context
	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("tenant ID not found in context")
	}

	// Validate request
	if req.Type == "" {
		return nil, fmt.Errorf("filter type is required")
	}

	// Manage filter
	if err := s.repo.ManageFilter(ctx, tenantID, req); err != nil {
		return nil, fmt.Errorf("failed to manage filter: %w", err)
	}

	// Return updated filters
	return s.GetFilters(ctx, req.Type)
}

// GetFilters returns available search filters
func (s *service) GetFilters(ctx context.Context, searchType string) (*FilterResponse, error) {
	// Get tenant ID from context
	tenantID, ok := ctx.Value("tenant_id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("tenant ID not found in context")
	}

	// Get filters
	filters, err := s.repo.GetFilters(ctx, tenantID, searchType)
	if err != nil {
		return nil, fmt.Errorf("failed to get filters: %w", err)
	}

	return &FilterResponse{
		Filters: filters,
	}, nil
}

// Helper methods

func (s *service) validateSearchQuery(req *SearchQuery) error {
	if req.Query == "" {
		return fmt.Errorf("query is required")
	}
	if req.Limit < 0 || req.Limit > 100 {
		return fmt.Errorf("limit must be between 0 and 100")
	}
	if req.Offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}
	return nil
}

func (s *service) validateProductSearchRequest(req *ProductSearchRequest) error {
	if req.Query == "" {
		return fmt.Errorf("query is required")
	}
	if req.Limit < 0 || req.Limit > 100 {
		return fmt.Errorf("limit must be between 0 and 100")
	}
	if req.Offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}
	if req.MinPrice != nil && req.MaxPrice != nil && *req.MinPrice > *req.MaxPrice {
		return fmt.Errorf("min_price cannot be greater than max_price")
	}
	return nil
}

func (s *service) getUserIDFromContext(ctx context.Context) *string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return &userID
	}
	return nil
}

func (s *service) getSessionIDFromContext(ctx context.Context) string {
	if sessionID, ok := ctx.Value("session_id").(string); ok {
		return sessionID
	}
	return "anonymous"
}