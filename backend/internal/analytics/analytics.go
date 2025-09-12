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

// Additional business logic methods

// CalculateConversionRate calculates the conversion rate for analytics stats
func (as *AnalyticsStats) CalculateConversionRate() float64 {
	if as.UniqueVisitors == 0 {
		return 0
	}
	return (float64(as.Orders) / float64(as.UniqueVisitors)) * 100
}

// UpdateConversionRate updates the conversion rate based on current data
func (as *AnalyticsStats) UpdateConversionRate() {
	as.ConversionRate = as.CalculateConversionRate()
}

// CalculateAvgOrderValue calculates the average order value
func (as *AnalyticsStats) CalculateAvgOrderValue() float64 {
	if as.Orders == 0 {
		return 0
	}
	return as.Revenue / float64(as.Orders)
}

// UpdateAvgOrderValue updates the average order value based on current data
func (as *AnalyticsStats) UpdateAvgOrderValue() {
	as.AvgOrderValue = as.CalculateAvgOrderValue()
}

// CalculateBounceRate calculates bounce rate from page views
func (as *AnalyticsStats) CalculateBounceRate(bounces int64) {
	if as.PageViews == 0 {
		as.BounceRate = 0
		return
	}
	as.BounceRate = (float64(bounces) / float64(as.PageViews)) * 100
}

// IsHighPerforming checks if the stats indicate high performance
func (as *AnalyticsStats) IsHighPerforming() bool {
	return as.ConversionRate > 2.0 && as.BounceRate < 50.0
}

// GetRevenueGrowth calculates revenue growth compared to previous period
func (as *AnalyticsStats) GetRevenueGrowth(previousRevenue float64) float64 {
	if previousRevenue == 0 {
		return 0
	}
	return ((as.Revenue - previousRevenue) / previousRevenue) * 100
}

// Product view business logic methods

// IsEngaged checks if the product view shows user engagement
func (pv *ProductView) IsEngaged() bool {
	if pv.DurationSeconds == nil {
		return false
	}
	return *pv.DurationSeconds > 60 // More than 1 minute is considered engaged
}

// GetEngagementLevel returns the engagement level based on duration
func (pv *ProductView) GetEngagementLevel() string {
	if pv.DurationSeconds == nil {
		return "unknown"
	}
	
	duration := *pv.DurationSeconds
	switch {
	case duration < 10:
		return "low"
	case duration < 60:
		return "medium"
	case duration < 300:
		return "high"
	default:
		return "very_high"
	}
}

// Purchase business logic methods

// IsHighValue checks if the purchase is considered high value
func (p *Purchase) IsHighValue(threshold float64) bool {
	return p.TotalAmount >= threshold
}

// GetOrderSize returns the order size category
func (p *Purchase) GetOrderSize() string {
	switch {
	case p.ItemCount == 1:
		return "single"
	case p.ItemCount <= 5:
		return "small"
	case p.ItemCount <= 15:
		return "medium"
	default:
		return "large"
	}
}

// GetAverageItemValue calculates the average value per item
func (p *Purchase) GetAverageItemValue() float64 {
	if p.ItemCount == 0 {
		return 0
	}
	return p.TotalAmount / float64(p.ItemCount)
}

// Analytics event business logic methods

// GetEventCategory returns the category of the event
func (e *AnalyticsEvent) GetEventCategory() string {
	switch e.EventType {
	case "page_view":
		return "engagement"
	case "product_view":
		return "product"
	case "purchase", "add_to_cart", "checkout":
		return "ecommerce"
	case "signup", "login", "logout":
		return "user"
	default:
		return "other"
	}
}

// HasUTMParameters checks if the event has UTM tracking parameters
func (e *AnalyticsEvent) HasUTMParameters() bool {
	return e.UTMSource != "" || e.UTMMedium != "" || e.UTMCampaign != ""
}

// GetTrafficSource returns the traffic source based on referrer and UTM
func (e *AnalyticsEvent) GetTrafficSource() string {
	if e.UTMSource != "" {
		return e.UTMSource
	}
	
	if e.Referrer == "" {
		return "direct"
	}
	
	// Simple referrer categorization
	if contains(e.Referrer, "google") {
		return "google"
	}
	if contains(e.Referrer, "facebook") {
		return "facebook"
	}
	if contains(e.Referrer, "twitter") {
		return "twitter"
	}
	
	return "referral"
}

// IsFromMobile checks if the event is from a mobile device
func (e *AnalyticsEvent) IsFromMobile() bool {
	return contains(e.UserAgent, "Mobile") || contains(e.UserAgent, "Android") || contains(e.UserAgent, "iPhone")
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
			(s[:len(substr)] == substr || 
			 s[len(s)-len(substr):] == substr || 
			 indexOfSubstring(s, substr) >= 0)))
}

// Helper function to find substring index
func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
