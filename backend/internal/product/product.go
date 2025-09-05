package product

import (
	"errors"
	"fmt"
	"strings"
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
	
	// Aliases for backward compatibility
	ProductStatusDraft     = StatusDraft
	ProductStatusActive    = StatusActive
	ProductStatusInactive  = StatusInactive
	ProductStatusArchived  = StatusArchived
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
	ID           uuid.UUID `json:"id" gorm:"primarykey"`
	ProductID    uuid.UUID `json:"product_id" gorm:"not null;index"`
	Name         string    `json:"name" gorm:"not null"` // e.g., "Size: Large, Color: Red"
	SKU          string    `json:"sku,omitempty" gorm:"index"`
	Barcode      string    `json:"barcode,omitempty"`
	Price        float64   `json:"price"` // Override product price if different
	ComparePrice float64   `json:"compare_price,omitempty"` // Original price for discount display
	CostPrice    float64   `json:"cost_price,omitempty"` // For profit calculations
	
	// Variant-specific inventory
	InventoryQuantity int  `json:"inventory_quantity" gorm:"default:0"`
	TrackQuantity     bool `json:"track_quantity" gorm:"default:true"`
	AllowBackorder    bool `json:"allow_backorder" gorm:"default:false"`
	
	// Physical properties
	Weight float64 `json:"weight,omitempty"` // in grams
	Length float64 `json:"length,omitempty"` // in cm
	Width  float64 `json:"width,omitempty"`  // in cm
	Height float64 `json:"height,omitempty"` // in cm
	
	// Variant options (e.g., size: "Large", color: "Red")
	Options map[string]string `json:"options" gorm:"serializer:json"`
	
	// Images specific to this variant
	Images []string `json:"images,omitempty" gorm:"serializer:json"`
	Image  string   `json:"image,omitempty"` // Primary variant image
	
	// Default variant flag
	IsDefault bool `json:"is_default" gorm:"default:false"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductStats represents product statistics
type ProductStats struct {
	TotalProducts    int64   `json:"total_products"`
	ActiveProducts   int64   `json:"active_products"`
	DraftProducts    int64   `json:"draft_products"`
	OutOfStock       int64   `json:"out_of_stock"`
	LowStock         int64   `json:"low_stock"`
	TotalCategories  int64   `json:"total_categories"`
	TotalValue       float64 `json:"total_value"`
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
	if p.Status != ProductStatusActive {
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

// ValidateProductData validates product business rules
func (p *Product) ValidateProductData() error {
	if p.Name == "" {
		return errors.New("product name is required")
	}
	
	if p.Price < 0 {
		return errors.New("product price cannot be negative")
	}
	
	if p.ComparePrice > 0 && p.Price >= p.ComparePrice {
		return errors.New("compare price must be higher than selling price")
	}
	
	if p.CostPrice > 0 && p.CostPrice > p.Price {
		return errors.New("cost price should not exceed selling price")
	}
	
	if p.TrackQuantity && p.InventoryQuantity < 0 {
		return errors.New("inventory quantity cannot be negative")
	}
	
	return nil
}

// CalculateShippingWeight returns total shipping weight
func (p *Product) CalculateShippingWeight() float64 {
	if p.Weight > 0 {
		return p.Weight
	}
	// Default weight for digital products
	if p.Type == TypeDigital {
		return 0
	}
	// Default weight if not specified (in grams)
	return 100
}

// GetSEOTitle returns optimized title for SEO
func (p *Product) GetSEOTitle() string {
	if p.MetaTitle != "" {
		return p.MetaTitle
	}
	return p.Name
}

// GetSEODescription returns optimized description for SEO
func (p *Product) GetSEODescription() string {
	if p.MetaDescription != "" {
		return p.MetaDescription
	}
	
	// Generate from description if available
	if p.Description != "" {
		desc := p.Description
		if len(desc) > 155 {
			desc = desc[:152] + "..."
		}
		return desc
	}
	
	return "Buy " + p.Name + " at the best price"
}

// GetDimensions returns formatted dimensions
func (p *Product) GetDimensions() string {
	if p.Length > 0 && p.Width > 0 && p.Height > 0 {
		return fmt.Sprintf("%.2f x %.2f x %.2f cm", p.Length, p.Width, p.Height)
	}
	return ""
}

// IsDigital checks if product is digital
func (p *Product) IsDigital() bool {
	return p.Type == TypeDigital
}

// IsPhysical checks if product is physical
func (p *Product) IsPhysical() bool {
	return p.Type == TypePhysical
}

// IsService checks if product is a service
func (p *Product) IsService() bool {
	return p.Type == TypeService
}

// GetInventoryStatus returns inventory status string
func (p *Product) GetInventoryStatus() string {
	if !p.TrackQuantity {
		return "unlimited"
	}
	
	if p.InventoryQuantity <= 0 {
		if p.AllowBackorder {
			return "backorder"
		}
		return "out_of_stock"
	}
	
	if p.InventoryQuantity < 10 {
		return "low_stock"
	}
	
	return "in_stock"
}

// CanPurchase checks if product can be purchased
func (p *Product) CanPurchase(quantity int) bool {
	if p.Status != ProductStatusActive {
		return false
	}
	
	if !p.TrackQuantity {
		return true
	}
	
	if p.InventoryQuantity >= quantity {
		return true
	}
	
	return p.AllowBackorder
}

// GetVariantByOptions finds variant by option values
func (p *Product) GetVariantByOptions(options map[string]string) *ProductVariant {
	for _, variant := range p.Variants {
		match := true
		for key, value := range options {
			if variant.Options[key] != value {
				match = false
				break
			}
		}
		if match {
			return &variant
		}
	}
	return nil
}

// GetAvailableOptions returns all available option combinations
func (p *Product) GetAvailableOptions() map[string][]string {
	options := make(map[string][]string)
	
	for _, variant := range p.Variants {
		for key, value := range variant.Options {
			if !contains(options[key], value) {
				options[key] = append(options[key], value)
			}
		}
	}
	
	return options
}

// Enhanced ProductVariant methods

// GetDisplayName returns formatted variant name
func (v *ProductVariant) GetDisplayName() string {
	if v.Name != "" {
		return v.Name
	}
	
	// Generate from options
	var parts []string
	for key, value := range v.Options {
		parts = append(parts, fmt.Sprintf("%s: %s", key, value))
	}
	
	if len(parts) > 0 {
		return strings.Join(parts, ", ")
	}
	
	return "Default"
}

// Enhanced Category methods

// GetFullPath returns full category path (e.g., "Electronics > Smartphones > Android")
func (c *Category) GetFullPath() string {
	if c.Parent == nil {
		return c.Name
	}
	return c.Parent.GetFullPath() + " > " + c.Name
}

// GetLevel returns category depth level (0 for root)
func (c *Category) GetLevel() int {
	if c.Parent == nil {
		return 0
	}
	return c.Parent.GetLevel() + 1
}

// Helper functions

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
