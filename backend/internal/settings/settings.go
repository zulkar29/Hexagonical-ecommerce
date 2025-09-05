package settings

import (
	"time"
)

// Setting represents a configuration setting
type Setting struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TenantID  uint      `json:"tenant_id" gorm:"not null;index"`
	Section   string    `json:"section" gorm:"not null;index"` // general, seo, appearance, integrations
	Key       string    `json:"key" gorm:"not null;index"`
	Value     string    `json:"value" gorm:"type:text"`
	Type      string    `json:"type" gorm:"default:'string'"` // string, number, boolean, json, file
	IsPublic  bool      `json:"is_public" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName returns the table name for Setting model
func (Setting) TableName() string {
	return "settings"
}

// GetSettingsRequest represents request for getting settings
type GetSettingsRequest struct {
	Section string `form:"section" json:"section"` // general, seo, appearance, integrations
}

// GetSettingsResponse represents response for getting settings
type GetSettingsResponse struct {
	Settings map[string]interface{} `json:"settings"`
	Sections []string               `json:"sections,omitempty"`
}

// UpdateSettingsRequest represents request for updating settings
type UpdateSettingsRequest struct {
	Settings map[string]interface{} `json:"settings" binding:"required"`
}

// UpdateSettingsResponse represents response for updating settings
type UpdateSettingsResponse struct {
	Message  string                 `json:"message"`
	Settings map[string]interface{} `json:"settings"`
	Updated  []string               `json:"updated"`
}

// PublicSettingsResponse represents response for public settings
type PublicSettingsResponse struct {
	StoreName   string                 `json:"store_name"`
	Logo        string                 `json:"logo"`
	Contact     map[string]interface{} `json:"contact"`
	Theme       map[string]interface{} `json:"theme"`
	SEO         map[string]interface{} `json:"seo"`
	SocialLinks map[string]interface{} `json:"social_links"`
}

// SettingSection represents available setting sections
type SettingSection struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	Keys        []string `json:"keys"`
}

// DefaultSettings contains default setting values
var DefaultSettings = map[string]map[string]interface{}{
	"general": {
		"store_name":        "My Store",
		"store_description": "Welcome to our store",
		"store_email":       "contact@store.com",
		"store_phone":       "",
		"store_address":     "",
		"currency":          "USD",
		"timezone":          "UTC",
		"language":          "en",
	},
	"seo": {
		"meta_title":       "My Store",
		"meta_description": "Welcome to our store",
		"meta_keywords":    "",
		"google_analytics": "",
		"facebook_pixel":   "",
		"google_tag_manager": "",
	},
	"appearance": {
		"logo":              "",
		"favicon":           "",
		"primary_color":     "#007bff",
		"secondary_color":   "#6c757d",
		"font_family":       "Inter",
		"theme_mode":        "light",
	},
	"integrations": {
		"stripe_publishable_key": "",
		"paypal_client_id":       "",
		"mailchimp_api_key":      "",
		"facebook_app_id":        "",
		"google_client_id":       "",
	},
}

// PublicSettingKeys defines which settings are publicly accessible
var PublicSettingKeys = map[string][]string{
	"general": {"store_name", "store_description", "store_email", "store_phone", "store_address", "currency", "language"},
	"seo":     {"meta_title", "meta_description", "meta_keywords"},
	"appearance": {"logo", "favicon", "primary_color", "secondary_color", "font_family", "theme_mode"},
}