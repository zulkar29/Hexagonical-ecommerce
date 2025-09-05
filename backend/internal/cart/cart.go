package cart

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CartStatus represents the status of a cart
type CartStatus string

const (
	StatusActive    CartStatus = "active"
	StatusAbandoned CartStatus = "abandoned"
	StatusConverted CartStatus = "converted"
	StatusExpired   CartStatus = "expired"
)

// Cart represents a shopping cart in the system
type Cart struct {
	ID         uuid.UUID   `json:"id" gorm:"primarykey"`
	TenantID   uuid.UUID   `json:"tenant_id" gorm:"not null;index"`
	CustomerID *uuid.UUID  `json:"customer_id,omitempty" gorm:"index"` // Nullable for guest carts
	SessionID  string      `json:"session_id,omitempty" gorm:"index"` // For guest carts
	Status     CartStatus  `json:"status" gorm:"default:active"`
	
	// Cart totals
	Subtotal     float64 `json:"subtotal" gorm:"default:0"`
	TaxAmount    float64 `json:"tax_amount" gorm:"default:0"`
	ShippingCost float64 `json:"shipping_cost" gorm:"default:0"`
	DiscountAmount float64 `json:"discount_amount" gorm:"default:0"`
	Total        float64 `json:"total" gorm:"default:0"`
	
	// Applied discounts and coupons
	CouponCode   string     `json:"coupon_code,omitempty"`
	DiscountID   *uuid.UUID `json:"discount_id,omitempty" gorm:"index"`
	
	// Shipping information
	ShippingMethodID *uuid.UUID `json:"shipping_method_id,omitempty"`
	ShippingAddress  *Address   `json:"shipping_address,omitempty" gorm:"embedded;embeddedPrefix:shipping_"`
	BillingAddress   *Address   `json:"billing_address,omitempty" gorm:"embedded;embeddedPrefix:billing_"`
	
	// Cart metadata
	Currency     string `json:"currency" gorm:"default:USD"`
	Notes        string `json:"notes,omitempty"`
	AbandonedAt  *time.Time `json:"abandoned_at,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Items []CartItem `json:"items,omitempty" gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	CartID    uuid.UUID `json:"cart_id" gorm:"not null;index"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null;index"`
	VariantID *uuid.UUID `json:"variant_id,omitempty" gorm:"index"`
	
	// Item details (snapshot at time of adding)
	ProductName  string  `json:"product_name" gorm:"not null"`
	ProductSlug  string  `json:"product_slug"`
	VariantName  string  `json:"variant_name,omitempty"`
	SKU          string  `json:"sku,omitempty"`
	Price        float64 `json:"price" gorm:"not null"`
	ComparePrice float64 `json:"compare_price,omitempty"`
	Image        string  `json:"image,omitempty"`
	
	// Quantity and totals
	Quantity   int     `json:"quantity" gorm:"not null;default:1"`
	LineTotal  float64 `json:"line_total" gorm:"not null"`
	
	// Item metadata
	Customizations map[string]interface{} `json:"customizations,omitempty" gorm:"serializer:json"`
	Notes          string                 `json:"notes,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Address represents shipping/billing address
type Address struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Company     string `json:"company,omitempty"`
	Address1    string `json:"address1,omitempty"`
	Address2    string `json:"address2,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Country     string `json:"country,omitempty"`
	Phone       string `json:"phone,omitempty"`
}

// Business Logic Errors
var (
	ErrCartNotFound     = errors.New("cart not found")
	ErrCartExpired      = errors.New("cart has expired")
	ErrCartConverted    = errors.New("cart has already been converted to order")
	ErrCartNotModifiable = errors.New("cart cannot be modified")
	ErrItemNotFound     = errors.New("cart item not found")
	ErrInvalidQuantity  = errors.New("invalid quantity")
	ErrProductNotFound  = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidCoupon    = errors.New("invalid or expired coupon")
)

// Business Logic Methods for Cart

// IsActive checks if cart is active and can be modified
func (c *Cart) IsActive() bool {
	return c.Status == StatusActive
}

// IsExpired checks if cart has expired
func (c *Cart) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// CanModify checks if cart can be modified
func (c *Cart) CanModify() error {
	if c.Status == StatusConverted {
		return ErrCartConverted
	}
	if c.IsExpired() {
		return ErrCartExpired
	}
	return nil
}

// GetItemCount returns total number of items in cart
func (c *Cart) GetItemCount() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}

// GetUniqueItemCount returns number of unique items in cart
func (c *Cart) GetUniqueItemCount() int {
	return len(c.Items)
}

// FindItem finds a cart item by product and variant ID
func (c *Cart) FindItem(productID uuid.UUID, variantID *uuid.UUID) *CartItem {
	for i := range c.Items {
		item := &c.Items[i]
		if item.ProductID == productID {
			// Check variant match
			if variantID == nil && item.VariantID == nil {
				return item
			}
			if variantID != nil && item.VariantID != nil && *variantID == *item.VariantID {
				return item
			}
		}
	}
	return nil
}

// CalculateSubtotal calculates cart subtotal from items
func (c *Cart) CalculateSubtotal() float64 {
	subtotal := 0.0
	for _, item := range c.Items {
		subtotal += item.LineTotal
	}
	return subtotal
}

// CalculateTotal calculates final cart total
func (c *Cart) CalculateTotal() float64 {
	return c.Subtotal + c.TaxAmount + c.ShippingCost - c.DiscountAmount
}

// UpdateTotals recalculates all cart totals
func (c *Cart) UpdateTotals() {
	c.Subtotal = c.CalculateSubtotal()
	c.Total = c.CalculateTotal()
}

// MarkAsAbandoned marks cart as abandoned
func (c *Cart) MarkAsAbandoned() {
	c.Status = StatusAbandoned
	now := time.Now()
	c.AbandonedAt = &now
}

// MarkAsConverted marks cart as converted to order
func (c *Cart) MarkAsConverted() {
	c.Status = StatusConverted
}

// SetExpiration sets cart expiration time
func (c *Cart) SetExpiration(duration time.Duration) {
	expiresAt := time.Now().Add(duration)
	c.ExpiresAt = &expiresAt
}

// Business Logic Methods for CartItem

// CalculateLineTotal calculates line total for the item
func (ci *CartItem) CalculateLineTotal() {
	ci.LineTotal = ci.Price * float64(ci.Quantity)
}

// UpdateQuantity updates item quantity and recalculates line total
func (ci *CartItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	ci.Quantity = quantity
	ci.CalculateLineTotal()
	return nil
}

// GetDiscountAmount calculates discount amount for this item
func (ci *CartItem) GetDiscountAmount() float64 {
	if ci.ComparePrice <= 0 || ci.Price >= ci.ComparePrice {
		return 0
	}
	return (ci.ComparePrice - ci.Price) * float64(ci.Quantity)
}

// GetDiscountPercentage calculates discount percentage for this item
func (ci *CartItem) GetDiscountPercentage() float64 {
	if ci.ComparePrice <= 0 || ci.Price >= ci.ComparePrice {
		return 0
	}
	return ((ci.ComparePrice - ci.Price) / ci.ComparePrice) * 100
}

// Helper functions

// GetUUIDValue safely gets UUID value from pointer
func GetUUIDValue(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}

// GetStringPointer returns pointer to string
func GetStringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}