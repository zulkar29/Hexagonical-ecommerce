package order

import (
	"fmt"
	"time"
	"mime/multipart"

	"github.com/google/uuid"
	"ecommerce-saas/internal/user"
)

// OrderStatus represents the status of an order
type OrderStatus string

// PaymentStatus represents the payment status
type PaymentStatus string

// FulfillmentStatus represents the fulfillment status
type FulfillmentStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusDraft     OrderStatus = "draft"
	StatusConfirmed OrderStatus = "confirmed"
	StatusProcessing OrderStatus = "processing"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCompleted OrderStatus = "completed"
	StatusCancelled OrderStatus = "cancelled"
	StatusReturned  OrderStatus = "returned"
	
	// Aliases for backward compatibility
	OrderStatusPending = StatusPending
	OrderStatusDraft   = StatusDraft
	OrderStatusConfirmed = StatusConfirmed
	OrderStatusProcessing = StatusProcessing
	OrderStatusShipped = StatusShipped
	OrderStatusDelivered = StatusDelivered
	OrderStatusCompleted = StatusCompleted
	OrderStatusCancelled = StatusCancelled
	OrderStatusReturned = StatusReturned
)

const (
	PaymentPending       PaymentStatus = "pending"
	PaymentAuthorized    PaymentStatus = "authorized"
	PaymentPaid          PaymentStatus = "paid"
	PaymentFailed        PaymentStatus = "failed"
	PaymentRefunded      PaymentStatus = "refunded"
	PaymentRefundPending PaymentStatus = "refund_pending"
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
	
	// Shipping information (for backward compatibility)
	ShippingFirstName string `json:"shipping_first_name,omitempty"`
	ShippingLastName  string `json:"shipping_last_name,omitempty"`
	ShippingPhone     string `json:"shipping_phone,omitempty"`
	
	// Shipping address
	ShippingAddress Address `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`
	
	// Billing address (optional, defaults to shipping)
	BillingAddress Address `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`
	
	// Financial details
	SubtotalAmount float64 `json:"subtotal_amount" gorm:"not null"`
	ShippingAmount float64 `json:"shipping_amount" gorm:"default:0"`
	DiscountAmount float64 `json:"discount_amount" gorm:"default:0"`
	TaxAmount      float64 `json:"tax_amount" gorm:"default:0"`
	TotalAmount    float64 `json:"total_amount" gorm:"not null"`
	Currency       string  `json:"currency" gorm:"default:BDT"`
	
	// Payment information
	PaymentStatus  PaymentStatus `json:"payment_status" gorm:"default:pending"`
	PaymentMethod  string        `json:"payment_method,omitempty"`
	PaymentGateway string        `json:"payment_gateway,omitempty"`
	PaymentID      *uuid.UUID    `json:"payment_id,omitempty" gorm:"index"`
	TransactionID  string        `json:"transaction_id,omitempty"` 
	
	// Fulfillment information
	FulfillmentStatus FulfillmentStatus `json:"fulfillment_status" gorm:"default:pending"`
	TrackingNumber    string            `json:"tracking_number,omitempty"`
	TrackingURL       string            `json:"tracking_url,omitempty"`
	
	// Additional information
	Notes string `json:"notes,omitempty"`
	
	// Timestamps
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ShippedAt   *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	
	// Relations
	Items   []OrderItem    `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	History []OrderHistory `json:"history,omitempty" gorm:"foreignKey:OrderID"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID  uuid.UUID `json:"tenant_id" gorm:"not null;index"`
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

// OrderHistory represents the history/timeline of order status changes
type OrderHistory struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	OrderID  uuid.UUID `json:"order_id" gorm:"not null;index"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Status information
	FromStatus OrderStatus `json:"from_status,omitempty"`
	ToStatus   OrderStatus `json:"to_status" gorm:"not null"`
	
	// Payment status changes
	FromPaymentStatus PaymentStatus `json:"from_payment_status,omitempty"`
	ToPaymentStatus   PaymentStatus `json:"to_payment_status,omitempty"`
	
	// Fulfillment status changes
	FromFulfillmentStatus FulfillmentStatus `json:"from_fulfillment_status,omitempty"`
	ToFulfillmentStatus   FulfillmentStatus `json:"to_fulfillment_status,omitempty"`
	
	// Change details
	Action      string `json:"action" gorm:"not null"` // created, status_changed, cancelled, refunded, etc.
	Description string `json:"description,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Notes       string `json:"notes,omitempty"`
	
	// User who made the change
	ChangedBy     *uuid.UUID `json:"changed_by,omitempty" gorm:"index"`
	ChangedByType string     `json:"changed_by_type,omitempty"` // customer, admin, system
	
	// Additional data (JSON)
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	
	// Relations
	Order *Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// OrderDispute represents customer disputes for orders
type OrderDispute struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	OrderID     uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
	CustomerID  uuid.UUID `json:"customer_id" gorm:"type:uuid;not null;index"`
	Reason      string    `json:"reason" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"not null;default:'pending'"`
	Resolution  string    `json:"resolution" gorm:"type:text"`
	Evidence    map[string]interface{} `json:"evidence" gorm:"type:jsonb"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`

	// Relationships
	Order    *Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Customer *user.User  `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
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
	o.TotalAmount = o.SubtotalAmount + o.ShippingAmount - o.DiscountAmount
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

// GetTimeline returns the order timeline sorted by creation date
func (o *Order) GetTimeline() []OrderHistory {
	if len(o.History) == 0 {
		return []OrderHistory{}
	}
	
	// History is already sorted by CreatedAt in repository queries
	return o.History
}

// GetLatestHistoryEntry returns the most recent history entry
func (o *Order) GetLatestHistoryEntry() *OrderHistory {
	if len(o.History) == 0 {
		return nil
	}
	return &o.History[len(o.History)-1]
}

// Business Logic Methods for OrderHistory

// IsStatusChange checks if this history entry represents a status change
func (oh *OrderHistory) IsStatusChange() bool {
	return oh.FromStatus != oh.ToStatus && oh.ToStatus != ""
}

// IsPaymentStatusChange checks if this history entry represents a payment status change
func (oh *OrderHistory) IsPaymentStatusChange() bool {
	return oh.FromPaymentStatus != oh.ToPaymentStatus && oh.ToPaymentStatus != ""
}

// IsFulfillmentStatusChange checks if this history entry represents a fulfillment status change
func (oh *OrderHistory) IsFulfillmentStatusChange() bool {
	return oh.FromFulfillmentStatus != oh.ToFulfillmentStatus && oh.ToFulfillmentStatus != ""
}

// GetChangeDescription returns a human-readable description of the change
func (oh *OrderHistory) GetChangeDescription() string {
	if oh.Description != "" {
		return oh.Description
	}
	
	// Generate description based on action
	switch oh.Action {
	case "created":
		return "Order was created"
	case "status_changed":
		if oh.IsStatusChange() {
			return fmt.Sprintf("Status changed from %s to %s", oh.FromStatus, oh.ToStatus)
		}
	case "payment_status_changed":
		if oh.IsPaymentStatusChange() {
			return fmt.Sprintf("Payment status changed from %s to %s", oh.FromPaymentStatus, oh.ToPaymentStatus)
		}
	case "fulfillment_status_changed":
		if oh.IsFulfillmentStatusChange() {
			return fmt.Sprintf("Fulfillment status changed from %s to %s", oh.FromFulfillmentStatus, oh.ToFulfillmentStatus)
		}
	case "cancelled":
		return "Order was cancelled"
	case "refunded":
		return "Order was refunded"
	case "shipped":
		return "Order was shipped"
	case "delivered":
		return "Order was delivered"
	default:
		return oh.Action
	}
	
	return oh.Action
}

// GetChangeSummary returns a summary of all changes in this history entry
func (oh *OrderHistory) GetChangeSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"action":      oh.Action,
		"description": oh.GetChangeDescription(),
		"created_at":  oh.CreatedAt,
		"changed_by":  oh.ChangedBy,
		"changed_by_type": oh.ChangedByType,
	}
	
	if oh.Reason != "" {
		summary["reason"] = oh.Reason
	}
	
	if oh.Notes != "" {
		summary["notes"] = oh.Notes
	}
	
	if oh.IsStatusChange() {
		summary["status_change"] = map[string]interface{}{
			"from": oh.FromStatus,
			"to":   oh.ToStatus,
		}
	}
	
	if oh.IsPaymentStatusChange() {
		summary["payment_status_change"] = map[string]interface{}{
			"from": oh.FromPaymentStatus,
			"to":   oh.ToPaymentStatus,
		}
	}
	
	if oh.IsFulfillmentStatusChange() {
		summary["fulfillment_status_change"] = map[string]interface{}{
			"from": oh.FromFulfillmentStatus,
			"to":   oh.ToFulfillmentStatus,
		}
	}
	
	return summary
}

// Business Logic Methods for OrderDispute
func (od *OrderDispute) CanBeResolved() bool {
	return od.Status == "pending" || od.Status == "escalated"
}

func (od *OrderDispute) IsResolved() bool {
	return od.Status == "resolved" || od.Status == "closed"
}

func (od *OrderDispute) GetAge() time.Duration {
	return time.Since(od.CreatedAt)
}

func (od *OrderDispute) AddEvidence(key string, value interface{}) {
	if od.Evidence == nil {
		od.Evidence = make(map[string]interface{})
	}
	od.Evidence[key] = value
}

func (od *OrderDispute) Resolve(resolution string) {
	od.Status = "resolved"
	od.Resolution = resolution
	now := time.Now()
	od.ResolvedAt = &now
	od.UpdatedAt = now
}

func (od *OrderDispute) Escalate() {
	od.Status = "escalated"
	od.UpdatedAt = time.Now()
}

func (od *OrderDispute) Close() {
	od.Status = "closed"
	now := time.Now()
	od.ResolvedAt = &now
	od.UpdatedAt = now
}

// Request/Response Types for API handlers

// UpdateOrderRequest represents a request to update an order
type UpdateOrderRequest struct {
	CustomerEmail   *string  `json:"customer_email,omitempty"`
	CustomerPhone   *string  `json:"customer_phone,omitempty"`
	ShippingAddress *Address `json:"shipping_address,omitempty"`
	BillingAddress  *Address `json:"billing_address,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
}

// CancelOrderRequest represents a request to cancel an order
type CancelOrderRequest struct {
	Reason string `json:"reason" validate:"required"`
}



// RefundOrderRequest represents a request to refund an order
type RefundOrderRequest struct {
	PaymentID string  `json:"payment_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
	Reason    string  `json:"reason" validate:"required"`
}

// UpdateOrderStatusRequest represents a request to update order status
type UpdateOrderStatusRequest struct {
	Status         OrderStatus `json:"status" validate:"required"`
	TrackingNumber *string     `json:"tracking_number,omitempty"`
	TrackingURL    *string     `json:"tracking_url,omitempty"`
	Notes          *string     `json:"notes,omitempty"`
}

// ImportOrdersRequest represents a request to import orders
type ImportOrdersRequest struct {
	File   *multipart.FileHeader `form:"file" validate:"required"`
	Format string               `form:"format" validate:"required,oneof=csv excel"`
}

// BulkDeleteOrdersRequest represents a request to bulk delete orders
type BulkDeleteOrdersRequest struct {
	OrderIDs []string `json:"order_ids" validate:"required,min=1"`
	Reason   string   `json:"reason,omitempty"`
}

// Dispute Request/Response Types
type CreateDisputeRequest struct {
	Reason      string                 `json:"reason" binding:"required"`
	Description string                 `json:"description" binding:"required"`
	Evidence    map[string]interface{} `json:"evidence,omitempty"`
}

type UpdateDisputeRequest struct {
	Action     string                 `json:"action" binding:"required"` // resolve, escalate, close, add_evidence
	Resolution string                 `json:"resolution,omitempty"`
	Evidence   map[string]interface{} `json:"evidence,omitempty"`
}

type DisputeFilter struct {
	Status     string    `json:"status,omitempty"`
	CustomerID uuid.UUID `json:"customer_id,omitempty"`
	OrderID    uuid.UUID `json:"order_id,omitempty"`
	DateFrom   time.Time `json:"date_from,omitempty"`
	DateTo     time.Time `json:"date_to,omitempty"`
	Page       int       `json:"page,omitempty"`
	Limit      int       `json:"limit,omitempty"`
}
