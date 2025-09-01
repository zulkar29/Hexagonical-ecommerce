package entities

// TODO: Implement Store entity during development
// This entity represents an e-commerce store within a tenant

import (
	"time"
	// TODO: Add imports during implementation
	// "github.com/google/uuid"
	// "ecommerce-saas/internal/domain/valueobjects"
)

// Store represents an e-commerce store
type Store struct {
	// TODO: Implement actual fields
	// ID          uuid.UUID                `json:"id"`
	// TenantID    uuid.UUID                `json:"tenant_id"`
	// Name        string                   `json:"name"`
	// Description string                   `json:"description"`
	// Logo        string                   `json:"logo"`
	// Domain      valueobjects.StoreDomain `json:"domain"`
	// Settings    StoreSettings            `json:"settings"`
	// Status      valueobjects.StoreStatus `json:"status"`
	// CreatedAt   time.Time                `json:"created_at"`
	// UpdatedAt   time.Time                `json:"updated_at"`
	// DeletedAt   *time.Time               `json:"deleted_at,omitempty"`
}

// StoreSettings contains store-specific configurations
type StoreSettings struct {
	// TODO: Implement actual settings
	// General      GeneralSettings      `json:"general"`
	// Theme        ThemeSettings        `json:"theme"`
	// SEO          SEOSettings          `json:"seo"`
	// Shipping     ShippingSettings     `json:"shipping"`
	// Tax          TaxSettings          `json:"tax"`
	// Checkout     CheckoutSettings     `json:"checkout"`
	// Notifications NotificationSettings `json:"notifications"`
}

// GeneralSettings contains general store settings
type GeneralSettings struct {
	// TODO: Implement general settings
	// Currency         string `json:"currency"`
	// Timezone         string `json:"timezone"`
	// WeightUnit       string `json:"weight_unit"`
	// DimensionUnit    string `json:"dimension_unit"`
	// CustomerAccounts bool   `json:"customer_accounts"`
	// GuestCheckout    bool   `json:"guest_checkout"`
}

// ThemeSettings contains theme and appearance settings
type ThemeSettings struct {
	// TODO: Implement theme settings
	// ThemeID       string            `json:"theme_id"`
	// PrimaryColor  string            `json:"primary_color"`
	// SecondaryColor string           `json:"secondary_color"`
	// FontFamily    string            `json:"font_family"`
	// CustomCSS     string            `json:"custom_css"`
	// CustomJS      string            `json:"custom_js"`
	// Logo          string            `json:"logo"`
	// Favicon       string            `json:"favicon"`
}

// SEOSettings contains search engine optimization settings
type SEOSettings struct {
	// TODO: Implement SEO settings
	// MetaTitle       string `json:"meta_title"`
	// MetaDescription string `json:"meta_description"`
	// MetaKeywords    string `json:"meta_keywords"`
	// GoogleAnalytics string `json:"google_analytics"`
	// FacebookPixel   string `json:"facebook_pixel"`
	// SitemapEnabled  bool   `json:"sitemap_enabled"`
}

// Business methods (to be implemented)
// func (s *Store) UpdateSettings(settings StoreSettings) error
// func (s *Store) IsActive() bool
// func (s *Store) GetTheme() ThemeSettings
// func (s *Store) SetCustomDomain(domain string) error