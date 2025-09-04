package category

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryStatus represents the status of a category
type CategoryStatus string

const (
	StatusActive   CategoryStatus = "active"
	StatusInactive CategoryStatus = "inactive"
	StatusArchived CategoryStatus = "archived"
)

// Category represents a product category with hierarchical structure
type Category struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Basic information
	Name        string `json:"name" gorm:"not null"`
	Slug        string `json:"slug" gorm:"not null;index"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	Icon        string `json:"icon,omitempty"`
	
	// Hierarchy
	ParentID *uuid.UUID `json:"parent_id,omitempty" gorm:"index"`
	Level    int         `json:"level" gorm:"default:0"`
	Path     string      `json:"path" gorm:"index"` // e.g., "/electronics/computers/laptops"
	
	// Display and ordering
	SortOrder int            `json:"sort_order" gorm:"default:0"`
	Status    CategoryStatus `json:"status" gorm:"default:active"`
	
	// SEO
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
	MetaKeywords    string `json:"meta_keywords,omitempty"`
	
	// Features
	IsFeatured    bool `json:"is_featured" gorm:"default:false"`
	ShowInMenu    bool `json:"show_in_menu" gorm:"default:true"`
	ProductCount  int  `json:"product_count" gorm:"default:0"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Parent   *Category   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category  `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Products []Product   `json:"products,omitempty" gorm:"many2many:product_categories;"`
}

// Product represents a simplified product structure for category relations
type Product struct {
	ID         uuid.UUID   `json:"id" gorm:"primarykey"`
	TenantID   uuid.UUID   `json:"tenant_id"`
	Name       string      `json:"name"`
	Slug       string      `json:"slug"`
	Status     string      `json:"status"`
	Categories []Category  `json:"categories,omitempty" gorm:"many2many:product_categories;"`
}

// Business Logic Errors
var (
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryExists       = errors.New("category already exists")
	ErrInvalidParent        = errors.New("invalid parent category")
	ErrCircularReference    = errors.New("circular reference detected")
	ErrCategoryHasProducts  = errors.New("category has associated products")
	ErrCategoryHasChildren  = errors.New("category has child categories")
	ErrMaxDepthExceeded     = errors.New("maximum category depth exceeded")
	ErrInvalidSlug          = errors.New("invalid category slug")
)

// Business Logic Methods for Category

// IsActive checks if the category is active
func (c *Category) IsActive() bool {
	return c.Status == StatusActive
}

// IsRoot checks if the category is a root category (no parent)
func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}

// HasChildren checks if the category has child categories
func (c *Category) HasChildren() bool {
	return len(c.Children) > 0
}

// HasProducts checks if the category has associated products
func (c *Category) HasProducts() bool {
	return c.ProductCount > 0
}

// CanDelete checks if the category can be safely deleted
func (c *Category) CanDelete() bool {
	return !c.HasChildren() && !c.HasProducts()
}

// GetAncestors returns the path segments as ancestor names
func (c *Category) GetAncestors() []string {
	if c.Path == "" {
		return []string{}
	}
	
	// Remove leading slash and split
	path := strings.TrimPrefix(c.Path, "/")
	if path == "" {
		return []string{}
	}
	
	return strings.Split(path, "/")
}

// GetFullPath returns the full category path including current category
func (c *Category) GetFullPath() string {
	if c.Path == "" {
		return "/" + c.Slug
	}
	return c.Path + "/" + c.Slug
}

// UpdatePath updates the category path based on parent
func (c *Category) UpdatePath(parent *Category) {
	if parent == nil {
		c.Path = ""
		c.Level = 0
	} else {
		c.Path = parent.GetFullPath()
		c.Level = parent.Level + 1
	}
}

// ValidateDepth checks if the category depth is within limits
func (c *Category) ValidateDepth(maxDepth int) error {
	if c.Level >= maxDepth {
		return ErrMaxDepthExceeded
	}
	return nil
}

// GenerateSlug generates a URL-friendly slug from the name
func (c *Category) GenerateSlug() {
	if c.Slug == "" {
		c.Slug = strings.ToLower(strings.ReplaceAll(c.Name, " ", "-"))
		// Remove special characters and clean up
		c.Slug = strings.ReplaceAll(c.Slug, "&", "and")
		c.Slug = strings.ReplaceAll(c.Slug, "'", "")
	}
}

// UpdateProductCount updates the product count for the category
func (c *Category) UpdateProductCount(count int) {
	c.ProductCount = count
}

// SetFeatured sets the featured status
func (c *Category) SetFeatured(featured bool) {
	c.IsFeatured = featured
}

// SetMenuVisibility sets the menu visibility
func (c *Category) SetMenuVisibility(visible bool) {
	c.ShowInMenu = visible
}

// Activate activates the category
func (c *Category) Activate() {
	c.Status = StatusActive
}

// Deactivate deactivates the category
func (c *Category) Deactivate() {
	c.Status = StatusInactive
}

// Archive archives the category
func (c *Category) Archive() {
	c.Status = StatusArchived
}

// BeforeCreate GORM hook
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	
	// Generate slug if not provided
	c.GenerateSlug()
	
	return nil
}

// BeforeUpdate GORM hook
func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	// Regenerate slug if name changed
	if c.Slug == "" {
		c.GenerateSlug()
	}
	
	return nil
}

// Request/Response Structures

// CreateCategoryRequest represents category creation request
type CreateCategoryRequest struct {
	Name            string     `json:"name" validate:"required,min=1,max=100"`
	Slug            string     `json:"slug,omitempty" validate:"omitempty,min=1,max=100"`
	Description     string     `json:"description,omitempty" validate:"omitempty,max=500"`
	Image           string     `json:"image,omitempty"`
	Icon            string     `json:"icon,omitempty"`
	ParentID        *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder       int        `json:"sort_order,omitempty"`
	MetaTitle       string     `json:"meta_title,omitempty" validate:"omitempty,max=60"`
	MetaDescription string     `json:"meta_description,omitempty" validate:"omitempty,max=160"`
	MetaKeywords    string     `json:"meta_keywords,omitempty" validate:"omitempty,max=255"`
	IsFeatured      bool       `json:"is_featured,omitempty"`
	ShowInMenu      bool       `json:"show_in_menu,omitempty"`
}

// UpdateCategoryRequest represents category update request
type UpdateCategoryRequest struct {
	Name            string         `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Slug            string         `json:"slug,omitempty" validate:"omitempty,min=1,max=100"`
	Description     string         `json:"description,omitempty" validate:"omitempty,max=500"`
	Image           string         `json:"image,omitempty"`
	Icon            string         `json:"icon,omitempty"`
	ParentID        *uuid.UUID     `json:"parent_id,omitempty"`
	SortOrder       *int           `json:"sort_order,omitempty"`
	Status          CategoryStatus `json:"status,omitempty"`
	MetaTitle       string         `json:"meta_title,omitempty" validate:"omitempty,max=60"`
	MetaDescription string         `json:"meta_description,omitempty" validate:"omitempty,max=160"`
	MetaKeywords    string         `json:"meta_keywords,omitempty" validate:"omitempty,max=255"`
	IsFeatured      *bool          `json:"is_featured,omitempty"`
	ShowInMenu      *bool          `json:"show_in_menu,omitempty"`
}

// CategoryResponse represents category response
type CategoryResponse struct {
	*Category
	ChildrenCount int `json:"children_count"`
	ProductCount  int `json:"product_count"`
}

// CategoryTreeResponse represents hierarchical category tree
type CategoryTreeResponse struct {
	*Category
	Children []CategoryTreeResponse `json:"children,omitempty"`
}

// CategoryFilter represents category listing filters
type CategoryFilter struct {
	ParentID   *uuid.UUID     `json:"parent_id,omitempty"`
	Status     CategoryStatus `json:"status,omitempty"`
	Level      *int           `json:"level,omitempty"`
	IsFeatured *bool          `json:"is_featured,omitempty"`
	ShowInMenu *bool          `json:"show_in_menu,omitempty"`
	Search     string         `json:"search,omitempty"`
}

// CategoryStats represents category statistics
type CategoryStats struct {
	TotalCategories   int `json:"total_categories"`
	ActiveCategories  int `json:"active_categories"`
	RootCategories    int `json:"root_categories"`
	FeaturedCategories int `json:"featured_categories"`
	MaxDepth          int `json:"max_depth"`
	AvgProductsPerCategory float64 `json:"avg_products_per_category"`
}