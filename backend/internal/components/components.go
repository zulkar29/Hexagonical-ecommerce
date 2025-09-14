package components

import (
	"time"

	"github.com/google/uuid"
)

// ComponentType represents the type of component
type ComponentType string

// ComponentStatus represents the status of a component
type ComponentStatus string

// ComponentCategory represents the category of component
type ComponentCategory string

const (
	// Theme Templates
	TypeThemeTemplate ComponentType = "theme-template"

	// Layout Components
	TypeHeader     ComponentType = "header"
	TypeFooter     ComponentType = "footer"
	TypeNavigation ComponentType = "navigation"
	TypeSidebar    ComponentType = "sidebar"
	TypeGrid       ComponentType = "grid"
	TypeContainer  ComponentType = "container"

	// Content Blocks
	TypeHero        ComponentType = "hero"
	TypeBanner      ComponentType = "banner"
	TypeTextBlock   ComponentType = "text-block"
	TypeImageBlock  ComponentType = "image-block"
	TypeVideoBlock  ComponentType = "video-block"
	TypeTestimonial ComponentType = "testimonial"
	TypeFeatureList ComponentType = "feature-list"
	TypeCTA         ComponentType = "cta"

	// Product Widgets
	TypeProduct         ComponentType = "product"
	TypeProductGrid     ComponentType = "product-grid"
	TypeProductCarousel ComponentType = "product-carousel"
	TypeProductCard     ComponentType = "product-card"
	TypeProductFilter   ComponentType = "product-filter"
	TypeShoppingCart    ComponentType = "shopping-cart"
	TypeWishlist        ComponentType = "wishlist"

	// Media Gallery
	TypeImageGallery ComponentType = "image-gallery"
	TypeSlider       ComponentType = "slider"
	TypeCarousel     ComponentType = "carousel"
	TypeVideoGallery ComponentType = "video-gallery"

	// Interactive Elements
	TypeForm             ComponentType = "form"
	TypeContactForm      ComponentType = "contact-form"
	TypeNewsletterSignup ComponentType = "newsletter-signup"
	TypeSearchBar        ComponentType = "search-bar"
	TypeModal            ComponentType = "modal"
	TypeAccordion        ComponentType = "accordion"
	TypeTabs             ComponentType = "tabs"

	// Social & Marketing
	TypeSocialProof ComponentType = "social-proof"
	TypePromotional ComponentType = "promotional"
	TypeSocialMedia ComponentType = "social-media"
	TypeReviews     ComponentType = "reviews"
	TypeBlog        ComponentType = "blog"

	// E-commerce Specific
	TypeCheckout         ComponentType = "checkout"
	TypePayment          ComponentType = "payment"
	TypeShipping         ComponentType = "shipping"
	TypeOrderTracking    ComponentType = "order-tracking"
	TypeAccountDashboard ComponentType = "account-dashboard"
)

const (
	StatusDraft    ComponentStatus = "draft"
	StatusActive   ComponentStatus = "active"
	StatusInactive ComponentStatus = "inactive"
	StatusArchived ComponentStatus = "archived"
)

const (
	CategoryThemeTemplates      ComponentCategory = "theme-templates"
	CategoryLayoutComponents    ComponentCategory = "layout-components"
	CategoryContentBlocks       ComponentCategory = "content-blocks"
	CategoryProductWidgets      ComponentCategory = "product-widgets"
	CategoryMediaGallery        ComponentCategory = "media-gallery"
	CategoryInteractiveElements ComponentCategory = "interactive-elements"
	CategorySocialMarketing     ComponentCategory = "social-marketing"
	CategoryEcommerceSpecific   ComponentCategory = "ecommerce-specific"
	CategorySettings            ComponentCategory = "settings"
)

// Component represents a customizable component template
type Component struct {
	ID          uuid.UUID       `json:"id" gorm:"primarykey"`
	TenantID    uuid.UUID       `json:"tenant_id" gorm:"not null;index"`
	Name        string          `json:"name" gorm:"not null"`
	Slug        string          `json:"slug" gorm:"not null;index"`
	Type        ComponentType   `json:"type" gorm:"not null"`
	Status      ComponentStatus `json:"status" gorm:"default:draft"`
	Description string          `json:"description,omitempty"`

	// Component structure and styling
	HTML       string                 `json:"html" gorm:"type:text"`
	CSS        string                 `json:"css" gorm:"type:text"`
	JavaScript string                 `json:"javascript,omitempty" gorm:"type:text"`
	Config     map[string]interface{} `json:"config" gorm:"serializer:json"`

	// Customization options
	Customizable bool                   `json:"customizable" gorm:"default:true"`
	Options      map[string]interface{} `json:"options" gorm:"serializer:json"`

	// Preview and thumbnail
	Thumbnail  string `json:"thumbnail,omitempty"`
	PreviewURL string `json:"preview_url,omitempty"`

	// Metadata
	Tags     []string          `json:"tags,omitempty" gorm:"serializer:json"`
	Category ComponentCategory `json:"category,omitempty"`
	Version  string            `json:"version" gorm:"default:1.0.0"`

	// Usage tracking
	UsageCount int  `json:"usage_count" gorm:"default:0"`
	IsFeatured bool `json:"is_featured" gorm:"default:false"`
	IsDefault  bool `json:"is_default" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Instances []ComponentInstance `json:"instances,omitempty" gorm:"foreignKey:ComponentID"`
}

// ComponentInstance represents an instance of a component used in a specific page/theme
type ComponentInstance struct {
	ID          uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	ComponentID uuid.UUID `json:"component_id" gorm:"not null;index"`
	ThemeID     uuid.UUID `json:"theme_id" gorm:"not null;index"`

	// Instance-specific customizations
	CustomHTML string                 `json:"custom_html,omitempty" gorm:"type:text"`
	CustomCSS  string                 `json:"custom_css,omitempty" gorm:"type:text"`
	CustomJS   string                 `json:"custom_js,omitempty" gorm:"type:text"`
	Settings   map[string]interface{} `json:"settings" gorm:"serializer:json"`

	// Layout and positioning
	Position  int    `json:"position" gorm:"default:0"`
	Zone      string `json:"zone" gorm:"default:main"` // header, footer, main, sidebar
	IsVisible bool   `json:"is_visible" gorm:"default:true"`

	// Responsive settings
	Breakpoints map[string]interface{} `json:"breakpoints" gorm:"serializer:json"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Component *Component `json:"component,omitempty" gorm:"foreignKey:ComponentID"`
	Theme     *Theme     `json:"theme,omitempty" gorm:"foreignKey:ThemeID"`
}

// Theme represents a collection of component instances forming a complete design
type Theme struct {
	ID          uuid.UUID       `json:"id" gorm:"primarykey"`
	TenantID    uuid.UUID       `json:"tenant_id" gorm:"not null;index"`
	Name        string          `json:"name" gorm:"not null"`
	Slug        string          `json:"slug" gorm:"not null;index"`
	Description string          `json:"description,omitempty"`
	Status      ComponentStatus `json:"status" gorm:"default:draft"`

	// Theme settings
	GlobalCSS string                 `json:"global_css,omitempty" gorm:"type:text"`
	GlobalJS  string                 `json:"global_js,omitempty" gorm:"type:text"`
	Settings  map[string]interface{} `json:"settings" gorm:"serializer:json"`

	// Theme metadata
	Thumbnail  string   `json:"thumbnail,omitempty"`
	PreviewURL string   `json:"preview_url,omitempty"`
	Tags       []string `json:"tags,omitempty" gorm:"serializer:json"`
	Version    string   `json:"version" gorm:"default:1.0.0"`

	// Publishing
	IsActive    bool       `json:"is_active" gorm:"default:false"`
	IsDefault   bool       `json:"is_default" gorm:"default:false"`
	PublishedAt *time.Time `json:"published_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Instances []ComponentInstance `json:"instances,omitempty" gorm:"foreignKey:ThemeID"`
}

// ComponentTemplate represents predefined component templates
type ComponentTemplate struct {
	ID          uuid.UUID     `json:"id" gorm:"primarykey"`
	Name        string        `json:"name" gorm:"not null"`
	Slug        string        `json:"slug" gorm:"unique;not null"`
	Type        ComponentType `json:"type" gorm:"not null"`
	Description string        `json:"description,omitempty"`

	// Template content
	HTML       string                 `json:"html" gorm:"type:text"`
	CSS        string                 `json:"css" gorm:"type:text"`
	JavaScript string                 `json:"javascript,omitempty" gorm:"type:text"`
	Config     map[string]interface{} `json:"config" gorm:"serializer:json"`

	// Template metadata
	Thumbnail  string            `json:"thumbnail,omitempty"`
	Tags       []string          `json:"tags,omitempty" gorm:"serializer:json"`
	Category   ComponentCategory `json:"category,omitempty"`
	IsFree     bool              `json:"is_free" gorm:"default:true"`
	IsFeatured bool              `json:"is_featured" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ThemeTemplate represents complete theme templates with multiple components
type ThemeTemplate struct {
	ID          uuid.UUID `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"unique;not null"`
	Description string    `json:"description,omitempty"`

	// Theme template content
	Components []ThemeTemplateComponent `json:"components" gorm:"foreignKey:ThemeTemplateID"`
	GlobalCSS  string                   `json:"global_css,omitempty" gorm:"type:text"`
	GlobalJS   string                   `json:"global_js,omitempty" gorm:"type:text"`
	Settings   map[string]interface{}   `json:"settings" gorm:"serializer:json"`

	// Template metadata
	Thumbnail  string   `json:"thumbnail,omitempty"`
	PreviewURL string   `json:"preview_url,omitempty"`
	Tags       []string `json:"tags,omitempty" gorm:"serializer:json"`
	Category   string   `json:"category,omitempty"`
	IsFree     bool     `json:"is_free" gorm:"default:true"`
	IsFeatured bool     `json:"is_featured" gorm:"default:false"`
	IsDefault  bool     `json:"is_default" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ThemeTemplateComponent represents a component within a theme template
type ThemeTemplateComponent struct {
	ID                  uuid.UUID `json:"id" gorm:"primarykey"`
	ThemeTemplateID     uuid.UUID `json:"theme_template_id" gorm:"not null;index"`
	ComponentTemplateID uuid.UUID `json:"component_template_id" gorm:"not null;index"`

	// Component positioning and configuration
	Zone      string                 `json:"zone" gorm:"default:main"` // header, footer, main, sidebar
	Position  int                    `json:"position" gorm:"default:0"`
	Settings  map[string]interface{} `json:"settings" gorm:"serializer:json"`
	CustomCSS string                 `json:"custom_css,omitempty" gorm:"type:text"`
	CustomJS  string                 `json:"custom_js,omitempty" gorm:"type:text"`
	IsVisible bool                   `json:"is_visible" gorm:"default:true"`

	// Responsive settings
	Breakpoints map[string]interface{} `json:"breakpoints" gorm:"serializer:json"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	ThemeTemplate     *ThemeTemplate     `json:"theme_template,omitempty" gorm:"foreignKey:ThemeTemplateID"`
	ComponentTemplate *ComponentTemplate `json:"component_template,omitempty" gorm:"foreignKey:ComponentTemplateID"`
}

// Business Logic Methods

// IsActive checks if component is active
func (c *Component) IsActive() bool {
	return c.Status == StatusActive
}

// CanBeCustomized checks if component allows customization
func (c *Component) CanBeCustomized() bool {
	return c.Customizable
}

// IncrementUsage increments the usage count
func (c *Component) IncrementUsage() {
	c.UsageCount++
}

// IsPublished checks if theme is published and active
func (t *Theme) IsPublished() bool {
	return t.IsActive && t.PublishedAt != nil
}

// Publish marks theme as published
func (t *Theme) Publish() {
	now := time.Now()
	t.IsActive = true
	t.PublishedAt = &now
	t.Status = StatusActive
}

// Unpublish marks theme as unpublished
func (t *Theme) Unpublish() {
	t.IsActive = false
	t.PublishedAt = nil
	t.Status = StatusInactive
}

// GetVisibleInstances returns only visible component instances
func (t *Theme) GetVisibleInstances() []ComponentInstance {
	var visible []ComponentInstance
	for _, instance := range t.Instances {
		if instance.IsVisible {
			visible = append(visible, instance)
		}
	}
	return visible
}

// ComponentStats represents component statistics
type ComponentStats struct {
	TotalComponents    int64       `json:"total_components"`
	ActiveComponents   int64       `json:"active_components"`
	DraftComponents    int64       `json:"draft_components"`
	TotalThemes        int64       `json:"total_themes"`
	ActiveThemes       int64       `json:"active_themes"`
	TotalInstances     int64       `json:"total_instances"`
	MostUsedComponents []Component `json:"most_used_components"`
	PopularTypes       []struct {
		Type  ComponentType `json:"type"`
		Count int64         `json:"count"`
	} `json:"popular_types"`
}

// Filter types
type ComponentFilter struct {
	Type      ComponentType     `json:"type"`
	Category  ComponentCategory `json:"category"`
	Status    ComponentStatus   `json:"status"`
	Search    string            `json:"search"`
	SortBy    string            `json:"sort_by"`
	SortOrder string            `json:"sort_order"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
}

type ComponentInstanceFilter struct {
	ComponentID *uuid.UUID      `json:"component_id"`
	ThemeID     *uuid.UUID      `json:"theme_id"`
	Status      ComponentStatus `json:"status"`
	Zone        string          `json:"zone"`
	Visible     *bool           `json:"visible"`
	Page        int             `json:"page"`
	Limit       int             `json:"limit"`
	SortBy      string          `json:"sort_by"`
	SortOrder   string          `json:"sort_order"`
}

type ThemeFilter struct {
	Category  ComponentCategory `json:"category"`
	Status    ComponentStatus   `json:"status"`
	Active    *bool             `json:"active"`
	Search    string            `json:"search"`
	SortBy    string            `json:"sort_by"`
	SortOrder string            `json:"sort_order"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
}

type TemplateFilters struct {
	Type      ComponentType     `json:"type"`
	Category  ComponentCategory `json:"category"`
	Search    string            `json:"search"`
	Free      *bool             `json:"free"`
	Featured  *bool             `json:"featured"`
	SortBy    string            `json:"sort_by"`
	SortOrder string            `json:"sort_order"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
}

type ComponentTemplateFilter struct {
	Search     string   `json:"search"`
	Category   string   `json:"category"`
	Type       string   `json:"type"`
	Tags       []string `json:"tags"`
	IsFree     *bool    `json:"is_free"`
	IsFeatured *bool    `json:"is_featured"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	Offset     int      `json:"offset"`
	SortBy     string   `json:"sort_by"`
	SortOrder  string   `json:"sort_order"`
}

type ThemeTemplateFilter struct {
	Search     string   `json:"search"`
	Category   string   `json:"category"`
	Type       string   `json:"type"`
	Tags       []string `json:"tags"`
	IsFree     *bool    `json:"is_free"`
	IsFeatured *bool    `json:"is_featured"`
	IsDefault  *bool    `json:"is_default"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	Offset     int      `json:"offset"`
	SortBy     string   `json:"sort_by"`
	SortOrder  string   `json:"sort_order"`
}

type ComponentFilters struct {
	Type      ComponentType     `json:"type"`
	Category  ComponentCategory `json:"category"`
	Status    ComponentStatus   `json:"status"`
	Search    string            `json:"search"`
	Tags      []string          `json:"tags"`
	Featured  *bool             `json:"featured"`
	SortBy    string            `json:"sort_by"`
	SortOrder string            `json:"sort_order"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
}

type ThemeFilters struct {
	Category  ComponentCategory `json:"category"`
	Status    ComponentStatus   `json:"status"`
	Active    *bool             `json:"active"`
	Search    string            `json:"search"`
	SortBy    string            `json:"sort_by"`
	SortOrder string            `json:"sort_order"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
}
