package components

import (
	"time"

	"github.com/google/uuid"
	"ecommerce-saas/internal/shared/constants"
)

// ComponentType represents the type of component
type ComponentType string

// ComponentStatus represents the status of a component
type ComponentStatus = constants.ComponentStatus

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

// Use shared status constants
const (
	StatusDraft    = constants.ComponentStatusDraft
	StatusActive   = constants.ComponentStatusActive
	StatusInactive = constants.ComponentStatusInactive
	StatusArchived = constants.ComponentStatusArchived
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
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	TenantID    uuid.UUID       `json:"tenant_id" gorm:"type:uuid;not null;index"`
	TemplateID  *uuid.UUID      `json:"template_id,omitempty" gorm:"index"`
	Name        string          `json:"name" gorm:"not null"`
	Slug        string          `json:"slug" gorm:"not null;index"`
	Type        ComponentType   `json:"type" gorm:"not null"`
	Status      ComponentStatus `json:"status" gorm:"default:draft"`
	Description string          `json:"description,omitempty"`
	Category    ComponentCategory `json:"category,omitempty"`

	// Component structure and styling
	HTML       string                 `json:"html" gorm:"column:html_content;type:text"`
	CSS        string                 `json:"css" gorm:"column:css_content;type:text"`
	JavaScript string                 `json:"javascript,omitempty" gorm:"column:js_content;type:text"`
	Config     map[string]interface{} `json:"config" gorm:"serializer:json"`

	// Metadata
	Tags         []string `json:"tags,omitempty" gorm:"serializer:json"`
	IsFeatured   bool     `json:"is_featured" gorm:"default:false"`
	PreviewImage string   `json:"preview_image,omitempty" gorm:"column:preview_image"`

	// SEO
	MetaTitle       string `json:"meta_title,omitempty" gorm:"column:meta_title"`
	MetaDescription string `json:"meta_description,omitempty" gorm:"column:meta_description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Instances []ComponentInstance `json:"instances,omitempty" gorm:"foreignKey:ComponentID"`
}

// ComponentInstance represents an instance of a component used in a specific page/theme
type ComponentInstance struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ComponentID uuid.UUID `json:"component_id" gorm:"type:uuid;not null;index"`
	ThemeID     uuid.UUID `json:"theme_id" gorm:"type:uuid;not null;index"`

	// Instance-specific customizations
	CustomCSS  string                 `json:"custom_css,omitempty" gorm:"column:custom_css;type:text"`
	CustomJS   string                 `json:"custom_js,omitempty" gorm:"column:custom_js;type:text"`
	Settings   map[string]interface{} `json:"settings" gorm:"column:instance_config;serializer:json"`

	// Layout and positioning
	Position  string `json:"position" gorm:"column:position;not null"`
	SortOrder int    `json:"sort_order" gorm:"column:sort_order;default:0"`
	IsVisible bool   `json:"is_visible" gorm:"column:is_visible;default:true"`

	// Visibility rules
	VisibilityRules map[string]interface{} `json:"visibility_rules" gorm:"column:visibility_rules;serializer:json"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Component *Component `json:"component,omitempty" gorm:"foreignKey:ComponentID"`
	Theme     *Theme     `json:"theme,omitempty" gorm:"foreignKey:ThemeID"`
}

// Theme represents a collection of component instances forming a complete design
type Theme struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	TenantID    uuid.UUID       `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string          `json:"name" gorm:"not null"`
	Slug        string          `json:"slug" gorm:"not null;index"`
	Description string          `json:"description,omitempty"`
	Status      ComponentStatus `json:"status" gorm:"default:draft"`

	// Theme settings
	GlobalStyles map[string]interface{} `json:"global_styles,omitempty" gorm:"column:global_styles;serializer:json"`
	LayoutConfig map[string]interface{} `json:"layout_config,omitempty" gorm:"column:layout_config;serializer:json"`

	// Theme metadata
	PreviewImage string `json:"preview_image,omitempty" gorm:"column:preview_image"`

	// Publishing
	IsActive bool `json:"is_active" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Instances []ComponentInstance `json:"instances,omitempty" gorm:"foreignKey:ThemeID"`
}

// ComponentTemplate represents predefined component templates
type ComponentTemplate struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	TenantID    uuid.UUID     `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string        `json:"name" gorm:"not null"`
	Slug        string        `json:"slug" gorm:"not null"`
	Type        ComponentType `json:"type" gorm:"not null"`
	Description string        `json:"description,omitempty"`

	// Template content
	HTML       string                 `json:"html" gorm:"column:html_template;type:text"`
	CSS        string                 `json:"css" gorm:"column:css_template;type:text"`
	JavaScript string                 `json:"javascript,omitempty" gorm:"column:js_template;type:text"`
	ConfigSchema map[string]interface{} `json:"config_schema" gorm:"column:config_schema;serializer:json"`
	DefaultConfig map[string]interface{} `json:"default_config" gorm:"column:default_config;serializer:json"`

	// Template metadata
	PreviewImage string            `json:"preview_image,omitempty" gorm:"column:preview_image"`
	Tags       []string          `json:"tags,omitempty" gorm:"serializer:json"`
	Category   ComponentCategory `json:"category,omitempty"`
	IsFree     bool              `json:"is_free" gorm:"default:true"`
	IsFeatured bool              `json:"is_featured" gorm:"default:false"`
	Price      float64           `json:"price" gorm:"default:0.00"`

	// SEO
	MetaTitle       string `json:"meta_title,omitempty" gorm:"column:meta_title"`
	MetaDescription string `json:"meta_description,omitempty" gorm:"column:meta_description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ThemeTemplate represents complete theme templates with multiple components
type ThemeTemplate struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"not null"`
	Description string    `json:"description,omitempty"`

	// Theme template content
	Components []ThemeTemplateComponent `json:"components" gorm:"foreignKey:ThemeTemplateID"`
	GlobalCSS  string                   `json:"global_css,omitempty" gorm:"column:global_css;type:text"`
	GlobalJS   string                   `json:"global_js,omitempty" gorm:"column:global_js;type:text"`
	Settings   map[string]interface{}   `json:"settings" gorm:"serializer:json"`

	// Template metadata
	Thumbnail  string   `json:"thumbnail,omitempty"`
	PreviewURL string   `json:"preview_url,omitempty" gorm:"column:preview_url"`
	Tags       []string `json:"tags,omitempty" gorm:"serializer:json"`
	Category   string   `json:"category,omitempty"`
	IsFree     bool     `json:"is_free" gorm:"column:is_free;default:true"`
	IsFeatured bool     `json:"is_featured" gorm:"column:is_featured;default:false"`
	IsDefault  bool     `json:"is_default" gorm:"column:is_default;default:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ThemeTemplateComponent represents a component within a theme template
type ThemeTemplateComponent struct {
	ID                  uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	TenantID            uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ThemeTemplateID     uuid.UUID `json:"theme_template_id" gorm:"type:uuid;not null;index"`
	ComponentTemplateID uuid.UUID `json:"component_template_id" gorm:"type:uuid;not null;index"`

	// Component positioning and configuration
	Zone      string                 `json:"zone" gorm:"default:main"` // header, footer, main, sidebar
	Position  int                    `json:"position" gorm:"default:0"`
	Settings  map[string]interface{} `json:"settings" gorm:"serializer:json"`
	CustomCSS string                 `json:"custom_css,omitempty" gorm:"column:custom_css;type:text"`
	CustomJS  string                 `json:"custom_js,omitempty" gorm:"column:custom_js;type:text"`
	IsVisible bool                   `json:"is_visible" gorm:"column:is_visible;default:true"`

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



// IsPublished checks if theme is published and active
func (t *Theme) IsPublished() bool {
	return t.IsActive && t.Status == StatusActive
}

// Publish marks theme as published
func (t *Theme) Publish() {
	t.IsActive = true
	t.Status = StatusActive
}

// Unpublish marks theme as unpublished
func (t *Theme) Unpublish() {
	t.IsActive = false
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
