package discount

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DiscountType represents different types of discounts
type DiscountType string

// DiscountStatus represents the status of a discount
type DiscountStatus string

// DiscountTarget represents what the discount applies to
type DiscountTarget string

// DiscountUsageLimit represents usage limit types
type DiscountUsageLimit string

const (
	TypePercentage DiscountType = "percentage"
	TypeFixed      DiscountType = "fixed"
	TypeFreeShipping DiscountType = "free_shipping"
	TypeBuyXGetY   DiscountType = "buy_x_get_y"
)

const (
	StatusActive   DiscountStatus = "active"
	StatusInactive DiscountStatus = "inactive"
	StatusExpired  DiscountStatus = "expired"
	StatusDraft    DiscountStatus = "draft"
)

const (
	TargetOrder      DiscountTarget = "order"
	TargetProduct    DiscountTarget = "product"
	TargetCategory   DiscountTarget = "category"
	TargetCollection DiscountTarget = "collection"
	TargetShipping   DiscountTarget = "shipping"
)

const (
	UsageLimitNone       DiscountUsageLimit = "none"
	UsageLimitPerCustomer DiscountUsageLimit = "per_customer"
	UsageLimitTotal      DiscountUsageLimit = "total"
)

// Discount represents a discount/coupon code
type Discount struct {
	ID       uuid.UUID      `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID      `json:"tenant_id" gorm:"not null;index"`
	
	// Basic information
	Code        string         `json:"code" gorm:"unique;not null"` // Coupon code
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description,omitempty"`
	Type        DiscountType   `json:"type" gorm:"not null"`
	Status      DiscountStatus `json:"status" gorm:"default:draft"`
	
	// Discount value
	Value    float64 `json:"value" gorm:"not null"` // Percentage (0-100) or fixed amount
	Currency string  `json:"currency,omitempty"`    // For fixed amount discounts
	
	// Minimum requirements
	MinOrderAmount   *float64 `json:"min_order_amount,omitempty"`
	MinItemQuantity  *int     `json:"min_item_quantity,omitempty"`
	
	// Target restrictions
	Target              DiscountTarget `json:"target" gorm:"default:order"`
	TargetProductIDs    []string       `json:"target_product_ids,omitempty" gorm:"serializer:json"`
	TargetCategoryIDs   []string       `json:"target_category_ids,omitempty" gorm:"serializer:json"`
	TargetCollectionIDs []string       `json:"target_collection_ids,omitempty" gorm:"serializer:json"`
	ExcludeProductIDs   []string       `json:"exclude_product_ids,omitempty" gorm:"serializer:json"`
	
	// Usage limits
	UsageLimit       *int               `json:"usage_limit,omitempty"`       // Total usage limit
	UsageLimitType   DiscountUsageLimit `json:"usage_limit_type" gorm:"default:none"`
	CustomerUsageLimit *int             `json:"customer_usage_limit,omitempty"` // Per customer limit
	UsageCount       int                `json:"usage_count" gorm:"default:0"`
	
	// Time restrictions
	StartsAt  *time.Time `json:"starts_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	
	// Customer restrictions
	CustomerEligibility string   `json:"customer_eligibility" gorm:"default:all"` // all, specific, groups
	EligibleCustomerIDs []string `json:"eligible_customer_ids,omitempty" gorm:"serializer:json"`
	EligibleCustomerGroups []string `json:"eligible_customer_groups,omitempty" gorm:"serializer:json"`
	
	// Buy X Get Y settings (for BOGO offers)
	BuyQuantity *int     `json:"buy_quantity,omitempty"`  // Number of items to buy
	GetQuantity *int     `json:"get_quantity,omitempty"`  // Number of items to get free/discounted
	GetValue    *float64 `json:"get_value,omitempty"`     // Discount on "get" items (0-100 for percentage)
	
	// Stackability
	Stackable       bool     `json:"stackable" gorm:"default:false"`           // Can combine with other discounts
	StackableWith   []string `json:"stackable_with,omitempty" gorm:"serializer:json"` // Specific discount IDs it can stack with
	ExclusiveGroup  string   `json:"exclusive_group,omitempty"`                // Mutually exclusive group
	
	// Settings
	ApplyOnce        bool `json:"apply_once" gorm:"default:false"`              // Apply only once per order (for percentage)
	ShowInStorefront bool `json:"show_in_storefront" gorm:"default:false"`      // Display publicly
	RequiresCode     bool `json:"requires_code" gorm:"default:true"`            // Automatic discount if false
	
	// Tracking
	CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Usages []DiscountUsage `json:"usages,omitempty" gorm:"foreignKey:DiscountID"`
}

// DiscountUsage tracks individual discount uses
type DiscountUsage struct {
	ID         uuid.UUID  `json:"id" gorm:"primarykey"`
	DiscountID uuid.UUID  `json:"discount_id" gorm:"not null;index"`
	TenantID   uuid.UUID  `json:"tenant_id" gorm:"not null;index"`
	
	// Order information
	OrderID       uuid.UUID `json:"order_id" gorm:"not null;index"`
	OrderNumber   string    `json:"order_number" gorm:"not null"`
	
	// Customer information
	CustomerID    *uuid.UUID `json:"customer_id,omitempty" gorm:"index"`
	CustomerEmail string     `json:"customer_email" gorm:"not null"`
	
	// Usage details
	DiscountAmount float64 `json:"discount_amount" gorm:"not null"` // Actual amount discounted
	OrderAmount    float64 `json:"order_amount" gorm:"not null"`    // Order total when discount was applied
	
	// Context
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	
	// Timestamps
	UsedAt    time.Time `json:"used_at"`
	CreatedAt time.Time `json:"created_at"`
}

// GiftCard represents a gift card
type GiftCard struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Gift card information
	Code         string         `json:"code" gorm:"unique;not null"`
	Status       DiscountStatus `json:"status" gorm:"default:active"`
	InitialValue float64        `json:"initial_value" gorm:"not null"`
	CurrentValue float64        `json:"current_value" gorm:"not null"`
	Currency     string         `json:"currency" gorm:"not null"`
	
	// Recipient information
	RecipientName  string `json:"recipient_name,omitempty"`
	RecipientEmail string `json:"recipient_email,omitempty"`
	Message        string `json:"message,omitempty"`
	
	// Purchaser information
	PurchasedBy    *uuid.UUID `json:"purchased_by,omitempty" gorm:"index"`
	PurchaseOrderID *uuid.UUID `json:"purchase_order_id,omitempty" gorm:"index"`
	
	// Restrictions
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	
	// Settings
	IsRefillable bool `json:"is_refillable" gorm:"default:false"`
	Notes        string `json:"notes,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Transactions []GiftCardTransaction `json:"transactions,omitempty" gorm:"foreignKey:GiftCardID"`
}

// GiftCardTransaction represents gift card usage/refill
type GiftCardTransaction struct {
	ID         uuid.UUID `json:"id" gorm:"primarykey"`
	GiftCardID uuid.UUID `json:"gift_card_id" gorm:"not null;index"`
	TenantID   uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Transaction details
	Type        string    `json:"type" gorm:"not null"` // usage, refill, initial
	Amount      float64   `json:"amount" gorm:"not null"`
	Balance     float64   `json:"balance" gorm:"not null"` // Balance after transaction
	Description string    `json:"description,omitempty"`
	
	// Order information (for usage transactions)
	OrderID     *uuid.UUID `json:"order_id,omitempty" gorm:"index"`
	OrderNumber string     `json:"order_number,omitempty"`
	
	// User information
	CustomerID    *uuid.UUID `json:"customer_id,omitempty" gorm:"index"`
	CustomerEmail string     `json:"customer_email,omitempty"`
	
	// Admin information (for refills)
	ProcessedBy *uuid.UUID `json:"processed_by,omitempty" gorm:"index"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
}

// StoreCredit represents store credit for customers
type StoreCredit struct {
	ID         uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID   uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	CustomerID uuid.UUID `json:"customer_id" gorm:"not null;index"`
	
	// Credit information
	CurrentBalance float64 `json:"current_balance" gorm:"not null"`
	Currency       string  `json:"currency" gorm:"not null"`
	
	// Settings
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Transactions []StoreCreditTransaction `json:"transactions,omitempty" gorm:"foreignKey:StoreCreditID"`
}

// StoreCreditTransaction represents store credit usage/addition
type StoreCreditTransaction struct {
	ID            uuid.UUID `json:"id" gorm:"primarykey"`
	StoreCreditID uuid.UUID `json:"store_credit_id" gorm:"not null;index"`
	TenantID      uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	CustomerID    uuid.UUID `json:"customer_id" gorm:"not null;index"`
	
	// Transaction details
	Type        string  `json:"type" gorm:"not null"` // usage, addition, refund, admin_adjustment
	Amount      float64 `json:"amount" gorm:"not null"`
	Balance     float64 `json:"balance" gorm:"not null"` // Balance after transaction
	Description string  `json:"description,omitempty"`
	
	// Order information
	OrderID     *uuid.UUID `json:"order_id,omitempty" gorm:"index"`
	OrderNumber string     `json:"order_number,omitempty"`
	
	// Reference information
	RefundID   *uuid.UUID `json:"refund_id,omitempty" gorm:"index"`
	ReturnID   *uuid.UUID `json:"return_id,omitempty" gorm:"index"`
	
	// Admin information
	ProcessedBy *uuid.UUID `json:"processed_by,omitempty" gorm:"index"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
}

// Business Logic Methods

// IsActive checks if the discount is currently active
func (d *Discount) IsActive() bool {
	if d.Status != StatusActive {
		return false
	}
	
	now := time.Now()
	
	// Check start date
	if d.StartsAt != nil && now.Before(*d.StartsAt) {
		return false
	}
	
	// Check expiry date
	if d.ExpiresAt != nil && now.After(*d.ExpiresAt) {
		return false
	}
	
	// Check usage limit
	if d.UsageLimit != nil && d.UsageCount >= *d.UsageLimit {
		return false
	}
	
	return true
}

// CanUseDiscount checks if a customer can use this discount
func (d *Discount) CanUseDiscount(customerID *uuid.UUID, customerEmail string, customerUsageCount int) bool {
	if !d.IsActive() {
		return false
	}
	
	// Check customer eligibility
	if d.CustomerEligibility == "specific" {
		if customerID == nil {
			return false
		}
		
		found := false
		for _, id := range d.EligibleCustomerIDs {
			if eligibleID, err := uuid.Parse(id); err == nil && *customerID == eligibleID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check per-customer usage limit
	if d.CustomerUsageLimit != nil && customerUsageCount >= *d.CustomerUsageLimit {
		return false
	}
	
	return true
}

// CalculateDiscount calculates the discount amount for an order
func (d *Discount) CalculateDiscount(orderAmount float64, itemQuantity int) (float64, error) {
	if !d.IsActive() {
		return 0, errors.New("discount is not active")
	}
	
	// Check minimum requirements
	if d.MinOrderAmount != nil && orderAmount < *d.MinOrderAmount {
		return 0, errors.New("order amount does not meet minimum requirement")
	}
	
	if d.MinItemQuantity != nil && itemQuantity < *d.MinItemQuantity {
		return 0, errors.New("order quantity does not meet minimum requirement")
	}
	
	var discountAmount float64
	
	switch d.Type {
	case TypePercentage:
		discountAmount = orderAmount * (d.Value / 100)
	case TypeFixed:
		discountAmount = d.Value
		if discountAmount > orderAmount {
			discountAmount = orderAmount // Can't discount more than order total
		}
	case TypeFreeShipping:
		// This would need shipping amount to calculate
		discountAmount = 0 // Handled separately in shipping calculation
	case TypeBuyXGetY:
		// Complex BOGO calculation - would need item details
		discountAmount = 0 // TODO: Implement BOGO logic
	}
	
	return discountAmount, nil
}

// IsExpired checks if the discount has expired
func (d *Discount) IsExpired() bool {
	if d.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*d.ExpiresAt)
}

// GetRemainingUsage returns remaining usage count
func (d *Discount) GetRemainingUsage() *int {
	if d.UsageLimit == nil {
		return nil
	}
	
	remaining := *d.UsageLimit - d.UsageCount
	if remaining < 0 {
		remaining = 0
	}
	return &remaining
}

// Gift card methods

// IsValid checks if gift card is valid for use
func (gc *GiftCard) IsValid() bool {
	if gc.Status != StatusActive {
		return false
	}
	
	if gc.CurrentValue <= 0 {
		return false
	}
	
	if gc.ExpiresAt != nil && time.Now().After(*gc.ExpiresAt) {
		return false
	}
	
	return true
}

// CanUseAmount checks if gift card has sufficient balance
func (gc *GiftCard) CanUseAmount(amount float64) bool {
	return gc.IsValid() && gc.CurrentValue >= amount
}

// UseAmount deducts amount from gift card (returns actual amount used)
func (gc *GiftCard) UseAmount(amount float64) float64 {
	if !gc.IsValid() {
		return 0
	}
	
	usedAmount := amount
	if usedAmount > gc.CurrentValue {
		usedAmount = gc.CurrentValue
	}
	
	gc.CurrentValue -= usedAmount
	return usedAmount
}

// Store credit methods

// CanUseAmount checks if store credit has sufficient balance
func (sc *StoreCredit) CanUseAmount(amount float64) bool {
	if sc.ExpiresAt != nil && time.Now().After(*sc.ExpiresAt) {
		return false
	}
	
	return sc.CurrentBalance >= amount
}

// UseAmount deducts amount from store credit (returns actual amount used)
func (sc *StoreCredit) UseAmount(amount float64) float64 {
	if !sc.CanUseAmount(amount) {
		return 0
	}
	
	usedAmount := amount
	if usedAmount > sc.CurrentBalance {
		usedAmount = sc.CurrentBalance
	}
	
	sc.CurrentBalance -= usedAmount
	return usedAmount
}

// AddAmount adds amount to store credit
func (sc *StoreCredit) AddAmount(amount float64) {
	if amount > 0 {
		sc.CurrentBalance += amount
	}
}

// Validation methods

// ValidateDiscount validates discount data
func (d *Discount) Validate() error {
	if d.Code == "" {
		return errors.New("discount code is required")
	}
	
	// Normalize code to uppercase
	d.Code = strings.ToUpper(strings.TrimSpace(d.Code))
	
	if d.Title == "" {
		return errors.New("discount title is required")
	}
	
	if d.Value <= 0 {
		return errors.New("discount value must be positive")
	}
	
	// Validate percentage values
	if d.Type == TypePercentage && d.Value > 100 {
		return errors.New("percentage discount cannot exceed 100%")
	}
	
	// Validate buy X get Y values
	if d.Type == TypeBuyXGetY {
		if d.BuyQuantity == nil || *d.BuyQuantity <= 0 {
			return errors.New("buy quantity must be positive for BOGO offers")
		}
		if d.GetQuantity == nil || *d.GetQuantity <= 0 {
			return errors.New("get quantity must be positive for BOGO offers")
		}
	}
	
	// Validate date range
	if d.StartsAt != nil && d.ExpiresAt != nil && d.StartsAt.After(*d.ExpiresAt) {
		return errors.New("start date cannot be after expiry date")
	}
	
	return nil
}

// ValidateGiftCard validates gift card data
func (gc *GiftCard) Validate() error {
	if gc.Code == "" {
		return errors.New("gift card code is required")
	}
	
	if gc.InitialValue <= 0 {
		return errors.New("gift card initial value must be positive")
	}
	
	if gc.Currency == "" {
		return errors.New("gift card currency is required")
	}
	
	return nil
}