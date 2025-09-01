package entities

import (
	"time"
	"github.com/google/uuid"
)

// ProductStatus defines product status types
type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusDraft    ProductStatus = "draft"
)

// Product represents a product in the store
type Product struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID     `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string        `json:"name" gorm:"not null"`
	Description string        `json:"description"`
	SKU         string        `json:"sku" gorm:"uniqueIndex:idx_tenant_sku"`
	Price       float64       `json:"price" gorm:"not null"`
	ComparePrice *float64     `json:"compare_price,omitempty"`
	CostPrice   *float64      `json:"cost_price,omitempty"`
	Status      ProductStatus `json:"status" gorm:"not null;default:'draft'"`
	Vendor      string        `json:"vendor"`
	Tags        []string      `json:"tags" gorm:"type:jsonb"`
	Images      []ProductImage `json:"images" gorm:"foreignKey:ProductID"`
	Variants    []ProductVariant `json:"variants" gorm:"foreignKey:ProductID"`
	Categories  []ProductCategory `json:"categories" gorm:"many2many:product_categories;"`
	Inventory   ProductInventory `json:"inventory" gorm:"embedded"`
	SEO         ProductSEO       `json:"seo" gorm:"embedded"`
	Attributes  map[string]interface{} `json:"attributes" gorm:"type:jsonb"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   *time.Time    `json:"deleted_at,omitempty" gorm:"index"`
}

// ProductImage represents a product image
type ProductImage struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	URL       string    `json:"url" gorm:"not null"`
	AltText   string    `json:"alt_text"`
	Position  int       `json:"position" gorm:"default:0"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductVariant represents a product variant (size, color, etc.)
type ProductVariant struct {
	ID           uuid.UUID           `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID    uuid.UUID           `json:"product_id" gorm:"type:uuid;not null"`
	Title        string              `json:"title" gorm:"not null"`
	SKU          string              `json:"sku" gorm:"uniqueIndex"`
	Price        *float64            `json:"price,omitempty"`
	ComparePrice *float64            `json:"compare_price,omitempty"`
	CostPrice    *float64            `json:"cost_price,omitempty"`
	Weight       *float64            `json:"weight,omitempty"`
	Dimensions   ProductDimensions   `json:"dimensions" gorm:"embedded"`
	ImageID      *uuid.UUID          `json:"image_id,omitempty" gorm:"type:uuid"`
	Options      map[string]string   `json:"options" gorm:"type:jsonb"`
	Inventory    ProductInventory    `json:"inventory" gorm:"embedded;embeddedPrefix:variant_"`
	Position     int                 `json:"position" gorm:"default:0"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// ProductDimensions represents product dimensions
type ProductDimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"` // cm, in
}

// ProductInventory represents inventory management data
type ProductInventory struct {
	TrackQuantity    bool   `json:"track_quantity" gorm:"default:true"`
	Quantity         int    `json:"quantity" gorm:"default:0"`
	ReservedQuantity int    `json:"reserved_quantity" gorm:"default:0"`
	LowStockAlert    int    `json:"low_stock_alert" gorm:"default:5"`
	AllowBackorder   bool   `json:"allow_backorder" gorm:"default:false"`
	Location         string `json:"location"`
}

// ProductSEO represents SEO-related product data
type ProductSEO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Slug        string `json:"slug" gorm:"uniqueIndex:idx_tenant_product_slug"`
}

// ProductCategory represents a product category
type ProductCategory struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Slug        string     `json:"slug" gorm:"uniqueIndex:idx_tenant_category_slug"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" gorm:"type:uuid"`
	Image       string     `json:"image"`
	Position    int        `json:"position" gorm:"default:0"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Business methods
func (p *Product) IsActive() bool {
	return p.Status == ProductStatusActive
}

func (p *Product) IsAvailable() bool {
	return p.IsActive() && (p.Inventory.Quantity > 0 || p.Inventory.AllowBackorder)
}

func (p *Product) UpdateInventory(quantity int) {
	p.Inventory.Quantity = quantity
	p.UpdatedAt = time.Now()
}

func (p *Product) GetMainImage() *ProductImage {
	for _, img := range p.Images {
		if img.Position == 0 {
			return &img
		}
	}
	if len(p.Images) > 0 {
		return &p.Images[0]
	}
	return nil
}

func (p *Product) GetFinalPrice() float64 {
	return p.Price
}

func (p *Product) HasDiscount() bool {
	return p.ComparePrice != nil && *p.ComparePrice > p.Price
}