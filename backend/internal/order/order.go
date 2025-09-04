package order

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

// PaymentStatus represents the payment status
type PaymentStatus string

// FulfillmentStatus represents the fulfillment status
type FulfillmentStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusProcessing OrderStatus = "processing"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
	StatusReturned  OrderStatus = "returned"
)

const (
	PaymentPending    PaymentStatus = "pending"
	PaymentAuthorized PaymentStatus = "authorized"
	PaymentPaid       PaymentStatus = "paid"
	PaymentFailed     PaymentStatus = "failed"
	PaymentRefunded   PaymentStatus = "refunded"
)

const (
	FulfillmentPending   FulfillmentStatus = "pending"
	FulfillmentPicked    FulfillmentStatus = "picked"
	FulfillmentPacked    FulfillmentStatus = "packed"
	FulfillmentShipped   FulfillmentStatus = "shipped"
	FulfillmentDelivered FulfillmentStatus = "delivered"
)

// Order represents an order in the system
type Order struct {
	ID       uuid.UUID   `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID   `json:"tenant_id" gorm:"not null;index"`
	UserID   uuid.UUID   `json:"user_id" gorm:"not null;index"`
	
	// Order details
	OrderNumber string      `json:"order_number" gorm:"unique;not null"`
	Status      OrderStatus `json:"status" gorm:"default:pending"`
	
	// Customer information
	CustomerEmail string `json:"customer_email" gorm:"not null"`
	CustomerPhone string `json:"customer_phone,omitempty"`
	
	// Shipping address
	ShippingAddress Address `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`
	
	// Billing address (optional, defaults to shipping)
	BillingAddress Address `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`
	
	// Financial details
	SubtotalAmount float64 `json:"subtotal_amount" gorm:"not null"`
	TaxAmount      float64 `json:"tax_amount" gorm:"default:0"`
	ShippingAmount float64 `json:"shipping_amount" gorm:"default:0"`
	DiscountAmount float64 `json:"discount_amount" gorm:"default:0"`
	TotalAmount    float64 `json:"total_amount" gorm:"not null"`
	Currency       string  `json:"currency" gorm:"default:BDT"`
	
	// Payment information
	PaymentStatus  PaymentStatus `json:"payment_status" gorm:"default:pending"`
	PaymentMethod  string        `json:"payment_method,omitempty"`
	PaymentGateway string        `json:"payment_gateway,omitempty"`
	TransactionID  string        `json:"transaction_id,omitempty"`
	
	// Fulfillment information
	FulfillmentStatus FulfillmentStatus `json:"fulfillment_status" gorm:"default:pending"`
	TrackingNumber    string            `json:"tracking_number,omitempty"`
	TrackingURL       string            `json:"tracking_url,omitempty"`
	
	// Timestamps
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ShippedAt   *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	
	// Relations
	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	OrderID   uuid.UUID `json:"order_id" gorm:"not null;index"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null;index"`
	VariantID *uuid.UUID `json:"variant_id,omitempty" gorm:"index"`
	
	// Product details (snapshot at time of order)
	ProductName  string  `json:"product_name" gorm:"not null"`
	ProductSKU   string  `json:"product_sku,omitempty"`
	VariantName  string  `json:"variant_name,omitempty"`
	UnitPrice    float64 `json:"unit_price" gorm:"not null"`
	Quantity     int     `json:"quantity" gorm:"not null"`
	TotalPrice   float64 `json:"total_price" gorm:"not null"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Address represents a shipping or billing address
type Address struct {
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Company   string `json:"company,omitempty"`
	Address1  string `json:"address1" gorm:"not null"`
	Address2  string `json:"address2,omitempty"`
	City      string `json:"city" gorm:"not null"`
	State     string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country   string `json:"country" gorm:"default:BD"`
	Phone     string `json:"phone,omitempty"`
}

// Business Logic Methods for Order

// IsEditable checks if the order can be modified
func (o *Order) IsEditable() bool {
	return o.Status == StatusPending || o.Status == StatusConfirmed
}

// IsCancellable checks if the order can be cancelled
func (o *Order) IsCancellable() bool {
	return o.Status != StatusCancelled && 
		   o.Status != StatusDelivered && 
		   o.Status != StatusReturned
}

// IsRefundable checks if the order can be refunded
func (o *Order) IsRefundable() bool {
	return o.PaymentStatus == PaymentPaid &&
		   (o.Status == StatusCancelled || o.Status == StatusReturned)
}

// CalculateTotal recalculates the total amount
func (o *Order) CalculateTotal() {
	o.TotalAmount = o.SubtotalAmount + o.TaxAmount + o.ShippingAmount - o.DiscountAmount
	if o.TotalAmount < 0 {
		o.TotalAmount = 0
	}
}

// GetFullName returns the customer's full name from shipping address
func (o *Order) GetFullName() string {
	return o.ShippingAddress.FirstName + " " + o.ShippingAddress.LastName
}

// GetShippingAddress returns formatted shipping address
func (o *Order) GetShippingAddress() string {
	addr := o.ShippingAddress.Address1
	if o.ShippingAddress.Address2 != "" {
		addr += ", " + o.ShippingAddress.Address2
	}
	addr += ", " + o.ShippingAddress.City
	if o.ShippingAddress.State != "" {
		addr += ", " + o.ShippingAddress.State
	}
	if o.ShippingAddress.PostalCode != "" {
		addr += " " + o.ShippingAddress.PostalCode
	}
	addr += ", " + o.ShippingAddress.Country
	return addr
}

// HasShipped checks if the order has been shipped
func (o *Order) HasShipped() bool {
	return o.FulfillmentStatus == FulfillmentShipped || 
		   o.FulfillmentStatus == FulfillmentDelivered
}

// IsDelivered checks if the order has been delivered
func (o *Order) IsDelivered() bool {
	return o.Status == StatusDelivered
}

// GetItemCount returns the total number of items in the order
func (o *Order) GetItemCount() int {
	count := 0
	for _, item := range o.Items {
		count += item.Quantity
	}
	return count
}

// Business Logic Methods for OrderItem

// GetLineTotal calculates the total for this line item
func (oi *OrderItem) GetLineTotal() float64 {
	return oi.UnitPrice * float64(oi.Quantity)
}

// UpdateTotal updates the total price based on unit price and quantity
func (oi *OrderItem) UpdateTotal() {
	oi.TotalPrice = oi.GetLineTotal()
}

// Business Logic Methods for Address

// GetFullName returns the full name from address
func (a *Address) GetFullName() string {
	return a.FirstName + " " + a.LastName
}

// IsComplete checks if the address has all required fields
func (a *Address) IsComplete() bool {
	return a.FirstName != "" &&
		   a.LastName != "" &&
		   a.Address1 != "" &&
		   a.City != "" &&
		   a.Country != ""
}

// GetFormattedAddress returns a formatted address string
func (a *Address) GetFormattedAddress() string {
	addr := a.Address1
	if a.Address2 != "" {
		addr += ", " + a.Address2
	}
	addr += ", " + a.City
	if a.State != "" {
		addr += ", " + a.State
	}
	if a.PostalCode != "" {
		addr += " " + a.PostalCode
	}
	addr += ", " + a.Country
	return addr
}

// TODO: Add more business logic methods
// - GenerateOrderNumber() string
// - ValidateOrder() error
// - CalculateTax() float64
// - CalculateShipping() float64
// - ApplyDiscount(code string) error
// - ProcessPayment() error
// - SendConfirmationEmail() error

// GenerateOrderNumber generates a unique order number
func (o *Order) GenerateOrderNumber() string {
	if o.OrderNumber != "" {
		return o.OrderNumber
	}
	// This would typically be done at the service level
	return "ORD-" + o.ID.String()[:8]
}

// ValidateOrder validates the order data
func (o *Order) ValidateOrder() error {
	if o.CustomerEmail == "" {
		return fmt.Errorf("customer email is required")
	}
	
	if !o.ShippingAddress.IsComplete() {
		return fmt.Errorf("shipping address is incomplete")
	}
	
	if len(o.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}
	
	if o.TotalAmount <= 0 {
		return fmt.Errorf("order total must be greater than zero")
	}
	
	return nil
}

// CalculateTaxAmount calculates tax based on shipping location
func (o *Order) CalculateTaxAmount() float64 {
	// Bangladesh VAT is typically 15%
	if o.ShippingAddress.Country == "BD" {
		return o.SubtotalAmount * 0.15
	}
	// No tax for other countries in this example
	return 0.0
}

// CalculateShippingAmount calculates shipping cost
func (o *Order) CalculateShippingAmount() float64 {
	// Free shipping for orders over 1000 BDT
	if o.SubtotalAmount >= 1000 {
		return 0.0
	}
	
	// Standard shipping rates
	if o.ShippingAddress.Country == "BD" {
		return 60.0 // 60 BDT for Bangladesh
	}
	
	return 200.0 // International shipping
}

// ApplyDiscount applies a discount to the order
func (o *Order) ApplyDiscount(discountAmount float64) {
	if discountAmount > 0 && discountAmount <= o.SubtotalAmount {
		o.DiscountAmount = discountAmount
		o.CalculateTotal()
	}
}

// GetPaymentDue returns the amount due for payment
func (o *Order) GetPaymentDue() float64 {
	if o.PaymentStatus == PaymentPaid {
		return 0.0
	}
	return o.TotalAmount
}

// GetRefundableAmount returns the amount that can be refunded
func (o *Order) GetRefundableAmount() float64 {
	if o.PaymentStatus != PaymentPaid {
		return 0.0
	}
	return o.TotalAmount
}

// GetOrderAge returns the age of the order in days
func (o *Order) GetOrderAge() int {
	return int(time.Since(o.CreatedAt).Hours() / 24)
}

// IsExpired checks if the order has expired (pending for too long)
func (o *Order) IsExpired() bool {
	if o.Status != StatusPending {
		return false
	}
	// Orders expire after 24 hours if not confirmed
	return time.Since(o.CreatedAt) > 24*time.Hour
}

// GetOrderSummary returns a summary of the order
func (o *Order) GetOrderSummary() map[string]interface{} {
	return map[string]interface{}{
		"order_number":     o.OrderNumber,
		"status":          o.Status,
		"customer_email":  o.CustomerEmail,
		"total_amount":    o.TotalAmount,
		"currency":        o.Currency,
		"item_count":      o.GetItemCount(),
		"created_at":      o.CreatedAt,
		"is_paid":        o.PaymentStatus == PaymentPaid,
		"is_shipped":     o.HasShipped(),
		"is_delivered":   o.IsDelivered(),
	}
}

// CanBeModified checks if the order can be modified
func (o *Order) CanBeModified() bool {
	return o.Status == StatusPending && o.PaymentStatus == PaymentPending
}

// RequiresAction checks if the order requires immediate action
func (o *Order) RequiresAction() bool {
	return (o.Status == StatusPending && o.IsExpired()) ||
		   (o.Status == StatusConfirmed && o.PaymentStatus == PaymentPending) ||
		   (o.Status == StatusProcessing && o.FulfillmentStatus == FulfillmentPending)
}
