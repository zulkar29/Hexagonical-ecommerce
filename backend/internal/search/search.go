package search

import (
	"time"

	"github.com/google/uuid"
)

// SearchResult represents a search result item
type SearchResult struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // product, category, page, etc.
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	URL         string                 `json:"url"`
	ImageURL    string                 `json:"image_url,omitempty"`
	Price       *float64               `json:"price,omitempty"`
	Score       float64                `json:"score"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// SearchQuery represents search parameters
type SearchQuery struct {
	Query      string   `json:"query" validate:"required"`
	Type       string   `json:"type,omitempty"`       // product, category, page, all
	Categories []string `json:"categories,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	MinPrice   *float64 `json:"min_price,omitempty"`
	MaxPrice   *float64 `json:"max_price,omitempty"`
	InStock    *bool    `json:"in_stock,omitempty"`
	SortBy     string   `json:"sort_by,omitempty"`    // relevance, price_asc, price_desc, newest
	Offset     int      `json:"offset,omitempty"`
	Limit      int      `json:"limit,omitempty"`
}

// SearchResponse represents search results
type SearchResponse struct {
	Results     []*SearchResult        `json:"results"`
	Total       int64                  `json:"total"`
	Offset      int                    `json:"offset"`
	Limit       int                    `json:"limit"`
	Query       string                 `json:"query"`
	Suggestions []string               `json:"suggestions,omitempty"`
	Filters     map[string]interface{} `json:"filters,omitempty"`
	Facets      map[string]interface{} `json:"facets,omitempty"`
}

// ProductSearchRequest represents product-specific search
type ProductSearchRequest struct {
	Query       string   `json:"query" validate:"required"`
	CategoryID  string   `json:"category_id,omitempty"`
	BrandID     string   `json:"brand_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	MinPrice    *float64 `json:"min_price,omitempty"`
	MaxPrice    *float64 `json:"max_price,omitempty"`
	InStock     *bool    `json:"in_stock,omitempty"`
	OnSale      *bool    `json:"on_sale,omitempty"`
	Rating      *float64 `json:"min_rating,omitempty"`
	SortBy      string   `json:"sort_by,omitempty"`
	Offset      int      `json:"offset,omitempty"`
	Limit       int      `json:"limit,omitempty"`
	IncludeFacets bool   `json:"include_facets,omitempty"`
}

// SuggestionRequest represents search suggestion parameters
type SuggestionRequest struct {
	Query string `json:"query" validate:"required"`
	Type  string `json:"type,omitempty"` // product, category, brand
	Limit int    `json:"limit,omitempty"`
}

// SuggestionResponse represents search suggestions
type SuggestionResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
	Query       string       `json:"query"`
}

// Suggestion represents a single search suggestion
type Suggestion struct {
	Text        string  `json:"text"`
	Type        string  `json:"type"`
	Score       float64 `json:"score"`
	ResultCount int64   `json:"result_count,omitempty"`
}

// SearchAnalyticsRequest represents analytics parameters
type SearchAnalyticsRequest struct {
	Type      string     `json:"type,omitempty"`       // queries, results, trends
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Limit     int        `json:"limit,omitempty"`
}

// SearchAnalyticsResponse represents search analytics
type SearchAnalyticsResponse struct {
	TopQueries    []QueryStat    `json:"top_queries,omitempty"`
	NoResults     []QueryStat    `json:"no_results,omitempty"`
	Trends        []TrendData    `json:"trends,omitempty"`
	TotalSearches int64          `json:"total_searches"`
	Period        string         `json:"period"`
	Metrics       SearchMetrics  `json:"metrics"`
}

// QueryStat represents query statistics
type QueryStat struct {
	Query       string  `json:"query"`
	Count       int64   `json:"count"`
	ResultCount int64   `json:"result_count"`
	ClickRate   float64 `json:"click_rate"`
}

// TrendData represents search trend data
type TrendData struct {
	Date   time.Time `json:"date"`
	Query  string    `json:"query"`
	Count  int64     `json:"count"`
	Change float64   `json:"change_percent"`
}

// SearchMetrics represents overall search metrics
type SearchMetrics struct {
	TotalSearches     int64   `json:"total_searches"`
	UniqueQueries     int64   `json:"unique_queries"`
	AverageResults    float64 `json:"average_results"`
	NoResultsRate     float64 `json:"no_results_rate"`
	ClickThroughRate  float64 `json:"click_through_rate"`
	AveragePosition   float64 `json:"average_position"`
}

// FilterRequest represents filter management
type FilterRequest struct {
	Type       string                 `json:"type" validate:"required"` // category, brand, price, rating
	Operation  string                 `json:"operation,omitempty"`      // create, update, delete
	Name       string                 `json:"name,omitempty"`
	Values     []string               `json:"values,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
	IsActive   *bool                  `json:"is_active,omitempty"`
}

// FilterResponse represents available filters
type FilterResponse struct {
	Filters []Filter `json:"filters"`
}

// Filter represents a search filter
type Filter struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	Name     string      `json:"name"`
	Values   []FilterValue `json:"values"`
	IsActive bool        `json:"is_active"`
}

// FilterValue represents a filter option
type FilterValue struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int64  `json:"count"`
}

// SearchLog represents search activity logging
type SearchLog struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID  uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	UserID    *string   `json:"user_id,omitempty" gorm:"index"`
	SessionID string    `json:"session_id" gorm:"index"`
	Query     string    `json:"query" gorm:"not null"`
	Type      string    `json:"type" gorm:"default:'global'"`
	Results   int64     `json:"results"`
	Clicked   bool      `json:"clicked" gorm:"default:false"`
	Position  *int      `json:"position,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (SearchLog) TableName() string {
	return "search_logs"
}