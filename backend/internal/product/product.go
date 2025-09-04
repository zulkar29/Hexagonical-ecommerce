package product

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ProductStatus represents the status of a product
type ProductStatus string

// ProductType represents the type of product
type ProductType string

const (
	StatusDraft     ProductStatus = "draft"
	StatusActive    ProductStatus = "active"
	StatusInactive  ProductStatus = "inactive"
	StatusArchived  ProductStatus = "archived"
)

const (
	TypePhysical ProductType = "physical"
	TypeDigital  ProductType = "digital"
	TypeService  ProductType = "service"
)

// Product represents a product in the system
type Product struct {
	ID          uuid.UUID     `json:"id" gorm:"primarykey"`
	TenantID    uuid.UUID     `json:"tenant_id" gorm:"not null;index"`
	Name        string        `json:"name" gorm:"not null"`
	Slug        string        `json:"slug" gorm:"unique;not null"`
	Description string        `json:"description"`
	Type        ProductType   `json:"type" gorm:"default:physical"`
	Status      ProductStatus `json:"status" gorm:"default:draft"`
	
	// Pricing
	Price         float64 `json:"price" gorm:"not null"`
	ComparePrice  float64 `json:"compare_price,omitempty"` // Original price for discount display
	CostPrice     float64 `json:"cost_price,omitempty"`    // For profit calculations
	
	// Inventory
	SKU               string `json:"sku,omitempty" gorm:"index"`
	Barcode           string `json:"barcode,omitempty"`
	InventoryQuantity int    `json:"inventory_quantity" gorm:"default:0"`
	TrackQuantity     bool   `json:"track_quantity" gorm:"default:true"`
	AllowBackorder    bool   `json:"allow_backorder" gorm:"default:false"`
	
	// Physical properties
	Weight float64 `json:"weight,omitempty"` // in grams
	Length float64 `json:"length,omitempty"` // in cm
	Width  float64 `json:"width,omitempty"`  // in cm
	Height float64 `json:"height,omitempty"` // in cm
	
	// SEO
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
	MetaKeywords    string `json:"meta_keywords,omitempty"`
	
	// Images
	FeaturedImage string   `json:"featured_image,omitempty"`
	Images        []string `json:"images,omitempty" gorm:"serializer:json"`
	
	// Categories and tags
	CategoryID uuid.UUID `json:"category_id,omitempty" gorm:"index"`
	Tags       []string  `json:"tags,omitempty" gorm:"serializer:json"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations (will be loaded separately)
	Variants []ProductVariant `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
	Category *Category        `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// ProductVariant represents product variations (size, color, etc.)
type ProductVariant struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null;index"`
	Name      string    `json:"name" gorm:"not null"` // e.g., "Size: Large, Color: Red"
	SKU       string    `json:"sku,omitempty" gorm:"index"`
	Price     float64   `json:"price"` // Override product price if different
	
	// Variant-specific inventory
	InventoryQuantity int  `json:"inventory_quantity" gorm:"default:0"`
	AllowBackorder    bool `json:"allow_backorder" gorm:"default:false"`
	
	// Variant options (e.g., size: "Large", color: "Red")
	Options map[string]string `json:"options" gorm:"serializer:json"`
	
	// Images specific to this variant
	Images []string `json:"images,omitempty" gorm:"serializer:json"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Category represents product categories
type Category struct {
	ID          uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"not null"`
	Description string    `json:"description,omitempty"`
	Image       string    `json:"image,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" gorm:"index"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	
	// SEO
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Products []Product  `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

// Business Logic Methods for Product

// IsAvailable checks if product is available for purchase
func (p *Product) IsAvailable() bool {
	if p.Status != StatusActive {
		return false
	}
	
	if p.TrackQuantity && p.InventoryQuantity <= 0 && !p.AllowBackorder {
		return false
	}
	
	return true
}

// IsInStock checks if product has inventory
func (p *Product) IsInStock() bool {
	if !p.TrackQuantity {
		return true
	}
	return p.InventoryQuantity > 0
}

// GetDiscountPercentage calculates discount percentage if compare price is set
func (p *Product) GetDiscountPercentage() float64 {
	if p.ComparePrice <= 0 || p.Price >= p.ComparePrice {
		return 0
	}
	return ((p.ComparePrice - p.Price) / p.ComparePrice) * 100
}

// GetProfitMargin calculates profit margin if cost price is set
func (p *Product) GetProfitMargin() float64 {
	if p.CostPrice <= 0 {
		return 0
	}
	return ((p.Price - p.CostPrice) / p.Price) * 100
}

// CanDecrementInventory checks if inventory can be decremented
func (p *Product) CanDecrementInventory(quantity int) bool {
	if !p.TrackQuantity {
		return true
	}
	
	if p.InventoryQuantity >= quantity {
		return true
	}
	
	return p.AllowBackorder
}

// DecrementInventory reduces inventory quantity
func (p *Product) DecrementInventory(quantity int) error {
	if !p.CanDecrementInventory(quantity) {
		return ErrInsufficientInventory
	}
	
	if p.TrackQuantity {
		p.InventoryQuantity -= quantity
	}
	
	return nil
}

// IncrementInventory increases inventory quantity
func (p *Product) IncrementInventory(quantity int) {
	if p.TrackQuantity {
		p.InventoryQuantity += quantity
	}
}

// GetMainImage returns the featured image or first image
func (p *Product) GetMainImage() string {
	if p.FeaturedImage != "" {
		return p.FeaturedImage
	}
	
	if len(p.Images) > 0 {
		return p.Images[0]
	}
	
	return ""
}

// HasVariants checks if product has variants
func (p *Product) HasVariants() bool {
	return len(p.Variants) > 0
}

// GetMinPrice returns the minimum price (considering variants)
func (p *Product) GetMinPrice() float64 {
	minPrice := p.Price
	
	for _, variant := range p.Variants {
		if variant.Price > 0 && variant.Price < minPrice {
			minPrice = variant.Price
		}
	}
	
	return minPrice
}

// GetMaxPrice returns the maximum price (considering variants)
func (p *Product) GetMaxPrice() float64 {
	maxPrice := p.Price
	
	for _, variant := range p.Variants {
		if variant.Price > maxPrice {
			maxPrice = variant.Price
		}
	}
	
	return maxPrice
}

// Business Logic Methods for ProductVariant

// IsAvailable checks if variant is available
func (v *ProductVariant) IsAvailable() bool {
	if v.InventoryQuantity <= 0 && !v.AllowBackorder {
		return false
	}
	return true
}

// GetEffectivePrice returns variant price or falls back to product price
func (v *ProductVariant) GetEffectivePrice(productPrice float64) float64 {
	if v.Price > 0 {
		return v.Price
	}
	return productPrice
}

// Business Logic Methods for Category

// IsRootCategory checks if category is at root level
func (c *Category) IsRootCategory() bool {
	return c.ParentID == nil
}

// HasChildren checks if category has subcategories
func (c *Category) HasChildren() bool {
	return len(c.Children) > 0
}

// Custom errors
var (
	ErrInsufficientInventory = errors.New("insufficient inventory")
	ErrProductNotFound       = errors.New("product not found")
	ErrCategoryNotFound      = errors.New("category not found")
)

// TODO: Add more business logic methods
// - ValidateProductData() error
// - GenerateSlug() string
// - CalculateShippingWeight() float64
// - GetSEOTitle() string
// - GetSEODescription() string
