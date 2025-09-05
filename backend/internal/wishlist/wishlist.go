package wishlist

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Wishlist represents a customer's wishlist
type Wishlist struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID   uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	CustomerID uuid.UUID `json:"customer_id" gorm:"type:uuid;not null;index"`
	Name       string    `json:"name" gorm:"size:255;not null"`
	Description string   `json:"description" gorm:"type:text"`
	IsDefault  bool      `json:"is_default" gorm:"default:false;index"`
	IsPublic   bool      `json:"is_public" gorm:"default:false"`
	ShareToken string    `json:"share_token,omitempty" gorm:"size:64;unique;index"`
	ItemCount  int       `json:"item_count" gorm:"default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relations
	Items []WishlistItem `json:"items,omitempty" gorm:"foreignKey:WishlistID;constraint:OnDelete:CASCADE"`
}

// WishlistItem represents an item in a wishlist
type WishlistItem struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID   uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	WishlistID uuid.UUID `json:"wishlist_id" gorm:"type:uuid;not null;index"`
	ProductID  uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	VariantID  *uuid.UUID `json:"variant_id,omitempty" gorm:"type:uuid;index"`
	Quantity   int       `json:"quantity" gorm:"default:1;check:quantity > 0"`
	Notes      string    `json:"notes" gorm:"type:text"`
	Priority   int       `json:"priority" gorm:"default:0"`
	AddedAt    time.Time `json:"added_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relations
	Wishlist *Wishlist `json:"wishlist,omitempty" gorm:"foreignKey:WishlistID"`
	Product  *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Variant  *ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
}

// Product represents a simplified product for wishlist relations
type Product struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	ComparePrice *float64 `json:"compare_price,omitempty"`
	IsAvailable bool      `json:"is_available"`
	Status      string    `json:"status"`
}

// ProductVariant represents a simplified product variant for wishlist relations
type ProductVariant struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	ProductID   uuid.UUID `json:"product_id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	ComparePrice *float64 `json:"compare_price,omitempty"`
	IsAvailable bool      `json:"is_available"`
	Inventory   int       `json:"inventory"`
}

// Business logic errors
var (
	ErrWishlistNotFound              = errors.New("wishlist not found")
	ErrWishlistItemNotFound          = errors.New("wishlist item not found")
	ErrWishlistNotOwned              = errors.New("wishlist not owned by user")
	ErrWishlistItemNotOwned          = errors.New("wishlist item not owned by user")
	ErrWishlistNameExists            = errors.New("wishlist name already exists")
	ErrWishlistLimitExceeded         = errors.New("wishlist limit exceeded")
	ErrWishlistFull                  = errors.New("wishlist is full")
	ErrCannotDeleteDefaultWishlist   = errors.New("cannot delete default wishlist")
	ErrInvalidTenantID               = errors.New("invalid tenant ID")
	ErrInvalidCustomerID             = errors.New("invalid customer ID")
	ErrInvalidWishlistID             = errors.New("invalid wishlist ID")
	ErrInvalidProductID              = errors.New("invalid product ID")
	ErrInvalidWishlistName           = errors.New("invalid wishlist name")
	ErrWishlistNameTooLong           = errors.New("wishlist name too long")
	ErrWishlistDescriptionTooLong    = errors.New("wishlist description too long")
	ErrInvalidQuantity               = errors.New("invalid quantity")
	ErrInvalidPriority               = errors.New("invalid priority")
	ErrItemNotesTooLong              = errors.New("item notes too long")
)

// Constants
const (
	MaxWishlistsPerCustomer = 10
	MaxItemsPerWishlist     = 100
	DefaultWishlistName     = "My Wishlist"
	ShareTokenLength        = 32
	
	// Pagination constants
	MaxPageSize     = 100
	DefaultPageSize = 20
	
	// Validation constants
	MaxWishlistNameLength        = 255
	MaxWishlistDescriptionLength = 1000
	MaxItemQuantity              = 100
	MaxItemNotesLength           = 500
	MaxItemPriority              = 10
)

// Wishlist business logic methods

// TableName returns the table name for GORM
func (w *Wishlist) TableName() string {
	return "wishlists"
}

// IsEmpty checks if the wishlist has no items
func (w *Wishlist) IsEmpty() bool {
	return w.ItemCount == 0
}

// CanAddItem checks if an item can be added to the wishlist
func (w *Wishlist) CanAddItem() bool {
	return w.ItemCount < MaxItemsPerWishlist
}

// CanDelete checks if the wishlist can be deleted
func (w *Wishlist) CanDelete() bool {
	// Default wishlists cannot be deleted
	return !w.IsDefault
}

// HasItem checks if a product/variant is already in the wishlist
func (w *Wishlist) HasItem(productID uuid.UUID, variantID *uuid.UUID) bool {
	for _, item := range w.Items {
		if item.ProductID == productID {
			if variantID == nil && item.VariantID == nil {
				return true
			}
			if variantID != nil && item.VariantID != nil && *variantID == *item.VariantID {
				return true
			}
		}
	}
	return false
}

// GetShareURL returns the public share URL for the wishlist
func (w *Wishlist) GetShareURL(baseURL string) string {
	if !w.IsPublic || w.ShareToken == "" {
		return ""
	}
	return baseURL + "/wishlists/shared/" + w.ShareToken
}

// GenerateShareToken generates a new share token for the wishlist
func (w *Wishlist) GenerateShareToken() {
	w.ShareToken = generateRandomToken(ShareTokenLength)
}

// MakePublic makes the wishlist public and generates a share token
func (w *Wishlist) MakePublic() {
	w.IsPublic = true
	if w.ShareToken == "" {
		w.GenerateShareToken()
	}
}

// MakePrivate makes the wishlist private and clears the share token
func (w *Wishlist) MakePrivate() {
	w.IsPublic = false
	w.ShareToken = ""
}

// UpdateItemCount updates the item count based on current items
func (w *Wishlist) UpdateItemCount() {
	w.ItemCount = len(w.Items)
}

// GORM hooks

// BeforeCreate is called before creating a wishlist
func (w *Wishlist) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	
	// Set default name if empty
	if w.Name == "" {
		w.Name = DefaultWishlistName
	}
	
	// Generate share token if public
	if w.IsPublic && w.ShareToken == "" {
		w.GenerateShareToken()
	}
	
	return nil
}

// BeforeUpdate is called before updating a wishlist
func (w *Wishlist) BeforeUpdate(tx *gorm.DB) error {
	// Generate share token if making public
	if w.IsPublic && w.ShareToken == "" {
		w.GenerateShareToken()
	}
	
	// Clear share token if making private
	if !w.IsPublic {
		w.ShareToken = ""
	}
	
	return nil
}

// WishlistItem business logic methods

// TableName returns the table name for GORM
func (wi *WishlistItem) TableName() string {
	return "wishlist_items"
}

// IsValid checks if the wishlist item is valid
func (wi *WishlistItem) IsValid() bool {
	return wi.ProductID != uuid.Nil && wi.Quantity > 0
}

// UpdateQuantity updates the item quantity with validation
func (wi *WishlistItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	
	wi.Quantity = quantity
	wi.UpdatedAt = time.Now()
	return nil
}

// GetDisplayName returns the display name for the item
func (wi *WishlistItem) GetDisplayName() string {
	if wi.Product == nil {
		return "Unknown Product"
	}
	
	name := wi.Product.Name
	if wi.Variant != nil && wi.Variant.Name != "" {
		name += " - " + wi.Variant.Name
	}
	
	return name
}

// GetPrice returns the current price for the item
func (wi *WishlistItem) GetPrice() float64 {
	if wi.Variant != nil {
		return wi.Variant.Price
	}
	
	if wi.Product != nil {
		return wi.Product.Price
	}
	
	return 0
}

// GetComparePrice returns the compare price for the item
func (wi *WishlistItem) GetComparePrice() *float64 {
	if wi.Variant != nil {
		return wi.Variant.ComparePrice
	}
	
	if wi.Product != nil {
		return wi.Product.ComparePrice
	}
	
	return nil
}

// IsAvailable checks if the item is available for purchase
func (wi *WishlistItem) IsAvailable() bool {
	if wi.Product == nil || !wi.Product.IsAvailable {
		return false
	}
	
	if wi.Variant != nil {
		return wi.Variant.IsAvailable && wi.Variant.Inventory >= wi.Quantity
	}
	
	return true
}

// HasDiscount checks if the item has a discount
func (wi *WishlistItem) HasDiscount() bool {
	comparePrice := wi.GetComparePrice()
	return comparePrice != nil && *comparePrice > wi.GetPrice()
}

// GetDiscountPercentage returns the discount percentage
func (wi *WishlistItem) GetDiscountPercentage() float64 {
	comparePrice := wi.GetComparePrice()
	if comparePrice == nil || *comparePrice <= wi.GetPrice() {
		return 0
	}
	
	currentPrice := wi.GetPrice()
	return (((*comparePrice) - currentPrice) / (*comparePrice)) * 100
}

// GORM hooks

// BeforeCreate is called before creating a wishlist item
func (wi *WishlistItem) BeforeCreate(tx *gorm.DB) error {
	if wi.ID == uuid.Nil {
		wi.ID = uuid.New()
	}
	
	if wi.Quantity <= 0 {
		wi.Quantity = 1
	}
	
	return nil
}

// Request and Response structures

// CreateWishlistRequest represents a request to create a wishlist
type CreateWishlistRequest struct {
	TenantID    uuid.UUID `json:"-"`
	CustomerID  uuid.UUID `json:"-"`
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
	IsDefault   bool   `json:"is_default"`
	IsPublic    bool   `json:"is_public"`
}

// UpdateWishlistRequest represents a request to update a wishlist
type UpdateWishlistRequest struct {
	Name        *string `json:"name" validate:"max=255"`
	Description *string `json:"description" validate:"max=1000"`
	IsDefault   *bool   `json:"is_default"`
	IsPublic    *bool   `json:"is_public"`
}

// AddItemRequest represents a request to add an item to a wishlist
type AddItemRequest struct {
	WishlistID uuid.UUID  `json:"-"`
	ProductID uuid.UUID  `json:"product_id" validate:"required"`
	VariantID *uuid.UUID `json:"variant_id"`
	Quantity  int        `json:"quantity" validate:"min=1,max=100"`
	Notes     string     `json:"notes" validate:"max=500"`
	Priority  int        `json:"priority" validate:"min=0,max=10"`
}

// UpdateItemRequest represents a request to update a wishlist item
type UpdateItemRequest struct {
	Quantity *int    `json:"quantity" validate:"min=1,max=100"`
	Notes    *string `json:"notes" validate:"max=500"`
	Priority *int    `json:"priority" validate:"min=0,max=10"`
}

// WishlistResponse represents a wishlist response
type WishlistResponse struct {
	*Wishlist
	ShareURL string `json:"share_url,omitempty"`
}

// WishlistItemResponse represents a wishlist item response
type WishlistItemResponse struct {
	*WishlistItem
	DisplayName        string  `json:"display_name"`
	CurrentPrice       float64 `json:"current_price"`
	ComparePrice       *float64 `json:"compare_price,omitempty"`
	DiscountPercentage float64 `json:"discount_percentage"`
	IsAvailable        bool    `json:"is_available"`
}

// WishlistFilter represents filters for listing wishlists
type WishlistFilter struct {
	CustomerID *uuid.UUID `json:"customer_id"`
	Name       string     `json:"name"`
	IsDefault  *bool      `json:"is_default"`
	IsPublic   *bool      `json:"is_public"`
	IsEmpty    *bool      `json:"is_empty"`
}

// WishlistItemFilter represents filters for listing wishlist items
type WishlistItemFilter struct {
	WishlistID *uuid.UUID `json:"wishlist_id"`
	ProductID  *uuid.UUID `json:"product_id"`
	VariantID  *uuid.UUID `json:"variant_id"`
	IsAvailable *bool     `json:"is_available"`
	HasDiscount *bool     `json:"has_discount"`
	MinPriority *int      `json:"min_priority"`
	MaxPriority *int      `json:"max_priority"`
}

// WishlistStats represents wishlist statistics
type WishlistStats struct {
	TotalWishlists     int64   `json:"total_wishlists"`
	TotalItems         int64   `json:"total_items"`
	AverageItemsPerWishlist float64 `json:"average_items_per_wishlist"`
	PublicWishlists    int64   `json:"public_wishlists"`
	PrivateWishlists   int64   `json:"private_wishlists"`
	EmptyWishlists     int64   `json:"empty_wishlists"`
	MostWishedProducts []ProductWishCount `json:"most_wished_products"`
}

// ProductWishCount represents product wish count statistics
type ProductWishCount struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	WishCount   int64     `json:"wish_count"`
}

// Utility functions

// generateRandomToken generates a random token of specified length
func generateRandomToken(length int) string {
	// This is a simplified implementation
	// In production, use a cryptographically secure random generator
	return uuid.New().String()[:length]
}