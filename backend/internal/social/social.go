package social

import (
	"time"

	"github.com/google/uuid"
)

// Platform represents social media platforms
type Platform string

// IntegrationStatus represents the status of social integration
type IntegrationStatus string

// SyncStatus represents the sync status of products
type SyncStatus string

const (
	PlatformInstagram Platform = "instagram"
	PlatformFacebook  Platform = "facebook"
	PlatformTikTok    Platform = "tiktok"
)

const (
	StatusPending    IntegrationStatus = "pending"
	StatusConnected  IntegrationStatus = "connected"
	StatusDisconnected IntegrationStatus = "disconnected"
	StatusError      IntegrationStatus = "error"
)

const (
	SyncPending   SyncStatus = "pending"
	SyncInProgress SyncStatus = "in_progress"
	SyncCompleted SyncStatus = "completed"
	SyncFailed    SyncStatus = "failed"
)

// SocialIntegration represents a social media platform integration
type SocialIntegration struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`

	// Platform details
	Platform Platform          `json:"platform" gorm:"not null"`
	Status   IntegrationStatus `json:"status" gorm:"default:pending"`

	// Authentication details
	AccessToken  string    `json:"-" gorm:"type:text"` // Hidden from JSON
	RefreshToken string    `json:"-" gorm:"type:text"` // Hidden from JSON
	ExpiresAt    time.Time `json:"expires_at"`

	// Platform-specific data
	PlatformUserID   string `json:"platform_user_id"`
	PlatformUsername string `json:"platform_username"`
	BusinessAccountID string `json:"business_account_id,omitempty"`
	CatalogID        string `json:"catalog_id,omitempty"`

	// Configuration
	AutoSync       bool `json:"auto_sync" gorm:"default:true"`
	SyncFrequency  int  `json:"sync_frequency" gorm:"default:24"` // hours
	LastSyncAt     *time.Time `json:"last_sync_at,omitempty"`
	NextSyncAt     *time.Time `json:"next_sync_at,omitempty"`

	// Error tracking
	LastError    string    `json:"last_error,omitempty" gorm:"type:text"`
	ErrorCount   int       `json:"error_count" gorm:"default:0"`
	LastErrorAt  *time.Time `json:"last_error_at,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SocialProduct represents a product synced to social platforms
type SocialProduct struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID  uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null;index"`

	// Platform details
	Platform        Platform   `json:"platform" gorm:"not null"`
	PlatformProductID string   `json:"platform_product_id,omitempty"`
	Status          SyncStatus `json:"status" gorm:"default:pending"`

	// Sync settings
	IsEnabled     bool   `json:"is_enabled" gorm:"default:true"`
	CustomTitle   string `json:"custom_title,omitempty"`
	CustomDescription string `json:"custom_description,omitempty" gorm:"type:text"`
	Tags          []string `json:"tags,omitempty" gorm:"serializer:json"`

	// Platform-specific settings
	InstagramSettings map[string]interface{} `json:"instagram_settings,omitempty" gorm:"serializer:json"`
	FacebookSettings  map[string]interface{} `json:"facebook_settings,omitempty" gorm:"serializer:json"`
	TikTokSettings    map[string]interface{} `json:"tiktok_settings,omitempty" gorm:"serializer:json"`

	// Sync tracking
	LastSyncAt    *time.Time `json:"last_sync_at,omitempty"`
	LastSyncError string     `json:"last_sync_error,omitempty" gorm:"type:text"`
	SyncAttempts  int        `json:"sync_attempts" gorm:"default:0"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SocialAnalytics represents social commerce analytics
type SocialAnalytics struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	Date     time.Time `json:"date" gorm:"not null;index"`

	// Platform breakdown
	Platform Platform `json:"platform" gorm:"not null"`

	// Metrics
	Impressions    int     `json:"impressions" gorm:"default:0"`
	Clicks         int     `json:"clicks" gorm:"default:0"`
	ProductViews   int     `json:"product_views" gorm:"default:0"`
	AddToCarts     int     `json:"add_to_carts" gorm:"default:0"`
	Purchases      int     `json:"purchases" gorm:"default:0"`
	Revenue        float64 `json:"revenue" gorm:"default:0"`
	CTR            float64 `json:"ctr" gorm:"default:0"` // Click-through rate
	ConversionRate float64 `json:"conversion_rate" gorm:"default:0"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Business Logic Methods

// IsConnected checks if the integration is active
func (si *SocialIntegration) IsConnected() bool {
	return si.Status == StatusConnected
}

// NeedsReauth checks if token needs refresh
func (si *SocialIntegration) NeedsReauth() bool {
	return time.Now().After(si.ExpiresAt)
}

// CanSync checks if sync is allowed
func (si *SocialIntegration) CanSync() bool {
	return si.IsConnected() && !si.NeedsReauth() && si.AutoSync
}

// ShouldSync checks if it's time for next sync
func (si *SocialIntegration) ShouldSync() bool {
	if !si.CanSync() {
		return false
	}

	if si.NextSyncAt == nil {
		return true
	}

	return time.Now().After(*si.NextSyncAt)
}

// ScheduleNextSync calculates and sets next sync time
func (si *SocialIntegration) ScheduleNextSync() {
	next := time.Now().Add(time.Duration(si.SyncFrequency) * time.Hour)
	si.NextSyncAt = &next
}

// RecordError records sync error
func (si *SocialIntegration) RecordError(err string) {
	si.LastError = err
	si.ErrorCount++
	now := time.Now()
	si.LastErrorAt = &now
}

// ClearError clears error state
func (si *SocialIntegration) ClearError() {
	si.LastError = ""
	si.ErrorCount = 0
	si.LastErrorAt = nil
}

// SocialProduct methods

// IsActive checks if product sync is active
func (sp *SocialProduct) IsActive() bool {
	return sp.IsEnabled && sp.Status != SyncFailed
}

// CanSync checks if product can be synced
func (sp *SocialProduct) CanSync() bool {
	return sp.IsEnabled && sp.SyncAttempts < 3
}

// RecordSyncAttempt records a sync attempt
func (sp *SocialProduct) RecordSyncAttempt() {
	sp.SyncAttempts++
	now := time.Now()
	sp.LastSyncAt = &now
}

// RecordSyncSuccess marks sync as successful
func (sp *SocialProduct) RecordSyncSuccess() {
	sp.Status = SyncCompleted
	sp.LastSyncError = ""
	sp.SyncAttempts = 0
	now := time.Now()
	sp.LastSyncAt = &now
}

// RecordSyncError records sync error
func (sp *SocialProduct) RecordSyncError(err string) {
	sp.Status = SyncFailed
	sp.LastSyncError = err
	sp.SyncAttempts++
}

// GetDisplayTitle returns custom title or falls back to product title
func (sp *SocialProduct) GetDisplayTitle(productTitle string) string {
	if sp.CustomTitle != "" {
		return sp.CustomTitle
	}
	return productTitle
}

// GetDisplayDescription returns custom description or falls back to product description
func (sp *SocialProduct) GetDisplayDescription(productDescription string) string {
	if sp.CustomDescription != "" {
		return sp.CustomDescription
	}
	return productDescription
}

// Request/Response structures

// ConnectPlatformRequest represents a request to connect a social platform
type ConnectPlatformRequest struct {
	Platform      Platform `json:"platform" validate:"required"`
	AccessToken   string   `json:"access_token" validate:"required"`
	RefreshToken  string   `json:"refresh_token,omitempty"`
	ExpiresIn     int      `json:"expires_in"` // seconds
	BusinessAccountID string `json:"business_account_id,omitempty"`
}

// SyncProductsRequest represents a request to sync products
type SyncProductsRequest struct {
	Platform   Platform    `json:"platform"`
	ProductIDs []uuid.UUID `json:"product_ids,omitempty"` // if empty, sync all
	Force      bool        `json:"force" gorm:"default:false"`
}

// UpdateSocialProductRequest represents a request to update social product settings
type UpdateSocialProductRequest struct {
	IsEnabled         *bool                   `json:"is_enabled"`
	CustomTitle       *string                 `json:"custom_title"`
	CustomDescription *string                 `json:"custom_description"`
	Tags              []string                `json:"tags"`
	InstagramSettings map[string]interface{}  `json:"instagram_settings"`
	FacebookSettings  map[string]interface{}  `json:"facebook_settings"`
	TikTokSettings    map[string]interface{}  `json:"tiktok_settings"`
}

// SocialAnalyticsResponse represents analytics response
type SocialAnalyticsResponse struct {
	Platform       Platform `json:"platform"`
	DateFrom       string   `json:"date_from"`
	DateTo         string   `json:"date_to"`
	TotalImpressions    int     `json:"total_impressions"`
	TotalClicks         int     `json:"total_clicks"`
	TotalProductViews   int     `json:"total_product_views"`
	TotalAddToCarts     int     `json:"total_add_to_carts"`
	TotalPurchases      int     `json:"total_purchases"`
	TotalRevenue        float64 `json:"total_revenue"`
	AverageCTR          float64 `json:"average_ctr"`
	AverageConversion   float64 `json:"average_conversion"`
	DailyBreakdown      []SocialAnalytics `json:"daily_breakdown"`
}

// IntegrationStatusResponse represents integration status
type IntegrationStatusResponse struct {
	Platform        Platform          `json:"platform"`
	Status          IntegrationStatus `json:"status"`
	Username        string            `json:"username"`
	IsConnected     bool              `json:"is_connected"`
	LastSyncAt      *time.Time        `json:"last_sync_at"`
	NextSyncAt      *time.Time        `json:"next_sync_at"`
	ProductCount    int               `json:"product_count"`
	SyncedCount     int               `json:"synced_count"`
	PendingCount    int               `json:"pending_count"`
	ErrorCount      int               `json:"error_count"`
	LastError       string            `json:"last_error,omitempty"`
}

// ProductSyncStatusResponse represents product sync status
type ProductSyncStatusResponse struct {
	ProductID         uuid.UUID  `json:"product_id"`
	ProductName       string     `json:"product_name"`
	Platform          Platform   `json:"platform"`
	Status            SyncStatus `json:"status"`
	PlatformProductID string     `json:"platform_product_id,omitempty"`
	LastSyncAt        *time.Time `json:"last_sync_at"`
	LastSyncError     string     `json:"last_sync_error,omitempty"`
	IsEnabled         bool       `json:"is_enabled"`
}