package analytics

import (
	"time"

	"github.com/google/uuid"
)

// AnalyticsEvent represents an analytics event
type AnalyticsEvent struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID  uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Event details
	EventType   string                 `json:"event_type" gorm:"not null;index"` // page_view, product_view, purchase, etc.
	EventName   string                 `json:"event_name" gorm:"not null"`
	Properties  map[string]interface{} `json:"properties" gorm:"serializer:json"`
	
	// User context
	UserID      *uuid.UUID `json:"user_id,omitempty" gorm:"index"`
	SessionID   string     `json:"session_id,omitempty" gorm:"index"`
	AnonymousID string     `json:"anonymous_id,omitempty" gorm:"index"`
	
	// Request context
	IPAddress   string `json:"ip_address,omitempty"`
	UserAgent   string `json:"user_agent,omitempty"`
	Referrer    string `json:"referrer,omitempty"`
	UTMSource   string `json:"utm_source,omitempty"`
	UTMMedium   string `json:"utm_medium,omitempty"`
	UTMCampaign string `json:"utm_campaign,omitempty"`
	
	// Timestamps
	Timestamp time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// PageView represents a page view event
type PageView struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Page details
	URL      string `json:"url" gorm:"not null"`
	Path     string `json:"path" gorm:"not null;index"`
	Title    string `json:"title,omitempty"`
	
	// User context
	UserID      *uuid.UUID `json:"user_id,omitempty" gorm:"index"`
	SessionID   string     `json:"session_id,omitempty" gorm:"index"`
	AnonymousID string     `json:"anonymous_id,omitempty" gorm:"index"`
	
	// Context
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Referrer  string `json:"referrer,omitempty"`
	
	// Duration (filled when user leaves page)
	DurationSeconds *int `json:"duration_seconds,omitempty"`
	
	Timestamp time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductView represents a product view event
type ProductView struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID  uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null;index"`
	
	// User context
	UserID      *uuid.UUID `json:"user_id,omitempty" gorm:"index"`
	SessionID   string     `json:"session_id,omitempty" gorm:"index"`
	AnonymousID string     `json:"anonymous_id,omitempty" gorm:"index"`
	
	// Context
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Referrer  string `json:"referrer,omitempty"`
	
	// Duration (filled when user leaves product page)
	DurationSeconds *int `json:"duration_seconds,omitempty"`
	
	Timestamp time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// Purchase represents a purchase event
type Purchase struct {
	ID      uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	OrderID  uuid.UUID `json:"order_id" gorm:"not null;index"`
	
	// User context
	UserID      *uuid.UUID `json:"user_id,omitempty" gorm:"index"`
	SessionID   string     `json:"session_id,omitempty" gorm:"index"`
	AnonymousID string     `json:"anonymous_id,omitempty" gorm:"index"`
	
	// Purchase details
	TotalAmount    float64 `json:"total_amount" gorm:"not null"`
	Currency       string  `json:"currency" gorm:"default:BDT"`
	ItemCount      int     `json:"item_count" gorm:"not null"`
	PaymentMethod  string  `json:"payment_method,omitempty"`
	
	// Context
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	
	Timestamp time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// AnalyticsStats represents aggregated analytics data
type AnalyticsStats struct {
	TenantID uuid.UUID `json:"tenant_id"`
	Date     time.Time `json:"date"`
	
	// Traffic metrics
	PageViews       int64 `json:"page_views"`
	UniqueVisitors  int64 `json:"unique_visitors"`
	Sessions        int64 `json:"sessions"`
	BounceRate      float64 `json:"bounce_rate"`
	AvgSessionTime  float64 `json:"avg_session_time"`
	
	// E-commerce metrics
	ProductViews    int64   `json:"product_views"`
	Orders          int64   `json:"orders"`
	Revenue         float64 `json:"revenue"`
	ConversionRate  float64 `json:"conversion_rate"`
	AvgOrderValue   float64 `json:"avg_order_value"`
	
	// Top performers
	TopPages     []string `json:"top_pages"`
	TopProducts  []string `json:"top_products"`
	TopReferrers []string `json:"top_referrers"`
}

// Business Logic Methods

// GetEventKey returns a unique key for the event type
func (e *AnalyticsEvent) GetEventKey() string {
	return e.EventType + ":" + e.EventName
}

// IsConversion checks if the event represents a conversion
func (e *AnalyticsEvent) IsConversion() bool {
	conversionEvents := []string{"purchase", "signup", "subscription", "lead"}
	for _, ce := range conversionEvents {
		if e.EventType == ce {
			return true
		}
	}
	return false
}

// GetValue returns the monetary value of the event if applicable
func (e *AnalyticsEvent) GetValue() float64 {
	if value, exists := e.Properties["value"]; exists {
		if v, ok := value.(float64); ok {
			return v
		}
	}
	return 0
}

// IsBounce checks if the page view is a bounce (short duration)
func (pv *PageView) IsBounce() bool {
	if pv.DurationSeconds == nil {
		return false
	}
	return *pv.DurationSeconds < 30 // Less than 30 seconds is considered a bounce
}

// GetDuration returns the duration in seconds, or 0 if not set
func (pv *PageView) GetDuration() int {
	if pv.DurationSeconds == nil {
		return 0
	}
	return *pv.DurationSeconds
}

// TODO: Add more business logic methods
// - CalculateConversionRate() float64
// - GetTopReferrers(limit int) []string
// - GetUserJourney(userID uuid.UUID) []*AnalyticsEvent
// - CalculateRetentionRate(days int) float64
