package returns

import (
	"time"

	"github.com/google/uuid"
)

// ReturnStatus represents the status of a return request
type ReturnStatus string

// ReturnType represents the type of return (refund or exchange)
type ReturnType string

const (
	StatusPending   ReturnStatus = "pending"
	StatusApproved  ReturnStatus = "approved"
	StatusRejected  ReturnStatus = "rejected"
	StatusProcessing ReturnStatus = "processing"
	StatusCompleted ReturnStatus = "completed"
	StatusCancelled ReturnStatus = "cancelled"
)

const (
	TypeRefund   ReturnType = "refund"
	TypeExchange ReturnType = "exchange"
)

// Return represents a return request in the system
type Return struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Order and customer information
	OrderID    uuid.UUID `json:"order_id" gorm:"not null;index"`
	CustomerID uuid.UUID `json:"customer_id" gorm:"not null;index"`
	
	// Return details
	ReturnNumber string       `json:"return_number" gorm:"unique;not null"`
	Status       ReturnStatus `json:"status" gorm:"default:pending"`
	Type         ReturnType   `json:"type" gorm:"not null"`
	
	// Return reason
	ReasonID    *uuid.UUID `json:"reason_id,omitempty" gorm:"index"`
	ReasonText  string     `json:"reason_text,omitempty"`
	Description string     `json:"description,omitempty"`
	
	// Financial details
	RefundAmount    float64 `json:"refund_amount" gorm:"default:0"`
	RestockingFee   float64 `json:"restocking_fee" gorm:"default:0"`
	ShippingRefund  float64 `json:"shipping_refund" gorm:"default:0"`
	TotalRefund     float64 `json:"total_refund" gorm:"default:0"`
	Currency        string  `json:"currency" gorm:"default:BDT"`
	
	// Exchange details (for exchange type)
	ExchangeOrderID *uuid.UUID `json:"exchange_order_id,omitempty" gorm:"index"`
	ExchangeAmount  float64    `json:"exchange_amount" gorm:"default:0"`
	
	// Shipping information
	ReturnShippingLabelURL string `json:"return_shipping_label_url,omitempty"`
	TrackingNumber         string `json:"tracking_number,omitempty"`
	TrackingURL            string `json:"tracking_url,omitempty"`
	
	// Processing information
	ProcessedBy   *uuid.UUID `json:"processed_by,omitempty" gorm:"index"`
	ProcessedAt   *time.Time `json:"processed_at,omitempty"`
	ApprovedBy    *uuid.UUID `json:"approved_by,omitempty" gorm:"index"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	RejectedBy    *uuid.UUID `json:"rejected_by,omitempty" gorm:"index"`
	RejectedAt    *time.Time `json:"rejected_at,omitempty"`
	RejectionReason string   `json:"rejection_reason,omitempty"`
	
	// Additional data
	Notes    string                 `json:"notes,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Items  []ReturnItem `json:"items,omitempty" gorm:"foreignKey:ReturnID"`
	Reason *ReturnReason `json:"reason,omitempty" gorm:"foreignKey:ReasonID"`
}

// ReturnItem represents an item being returned
type ReturnItem struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	ReturnID uuid.UUID `json:"return_id" gorm:"not null;index"`
	
	// Order item reference
	OrderItemID uuid.UUID `json:"order_item_id" gorm:"not null;index"`
	ProductID   uuid.UUID `json:"product_id" gorm:"not null;index"`
	VariantID   *uuid.UUID `json:"variant_id,omitempty" gorm:"index"`
	
	// Product details (snapshot)
	ProductName string  `json:"product_name" gorm:"not null"`
	ProductSKU  string  `json:"product_sku,omitempty"`
	VariantName string  `json:"variant_name,omitempty"`
	UnitPrice   float64 `json:"unit_price" gorm:"not null"`
	
	// Return details
	QuantityOrdered  int     `json:"quantity_ordered" gorm:"not null"`
	QuantityReturned int     `json:"quantity_returned" gorm:"not null"`
	RefundAmount     float64 `json:"refund_amount" gorm:"not null"`
	
	// Item condition
	Condition   string `json:"condition,omitempty"` // new, used, damaged
	ConditionNotes string `json:"condition_notes,omitempty"`
	
	// Exchange details (for exchange items)
	ExchangeProductID *uuid.UUID `json:"exchange_product_id,omitempty" gorm:"index"`
	ExchangeVariantID *uuid.UUID `json:"exchange_variant_id,omitempty" gorm:"index"`
	ExchangeQuantity  int        `json:"exchange_quantity" gorm:"default:0"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReturnReason represents predefined return reasons
type ReturnReason struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"` // defective, wrong_item, not_as_described, etc.
	
	// Settings
	IsActive           bool    `json:"is_active" gorm:"default:true"`
	RequiresApproval   bool    `json:"requires_approval" gorm:"default:true"`
	AllowsExchange     bool    `json:"allows_exchange" gorm:"default:true"`
	RestockingFeeRate  float64 `json:"restocking_fee_rate" gorm:"default:0"`
	MaxReturnDays      int     `json:"max_return_days" gorm:"default:30"`
	
	// Display settings
	DisplayOrder int    `json:"display_order" gorm:"default:0"`
	Color        string `json:"color,omitempty"`
	Icon         string `json:"icon,omitempty"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Business Logic Methods for Return

// IsEditable checks if the return can be modified
func (r *Return) IsEditable() bool {
	return r.Status == StatusPending
}

// IsCancellable checks if the return can be cancelled
func (r *Return) IsCancellable() bool {
	return r.Status == StatusPending || r.Status == StatusApproved
}

// IsProcessable checks if the return can be processed
func (r *Return) IsProcessable() bool {
	return r.Status == StatusApproved
}

// CanBeApproved checks if the return can be approved
func (r *Return) CanBeApproved() bool {
	return r.Status == StatusPending
}

// CanBeRejected checks if the return can be rejected
func (r *Return) CanBeRejected() bool {
	return r.Status == StatusPending
}

// CalculateTotalRefund calculates the total refund amount
func (r *Return) CalculateTotalRefund() {
	r.TotalRefund = r.RefundAmount + r.ShippingRefund - r.RestockingFee
	if r.TotalRefund < 0 {
		r.TotalRefund = 0
	}
}

// GenerateReturnNumber generates a unique return number
func (r *Return) GenerateReturnNumber() string {
	if r.ReturnNumber != "" {
		return r.ReturnNumber
	}
	// This would typically be done at the service level
	return "RET-" + r.ID.String()[:8]
}

// GetTotalItemsCount returns the total number of items being returned
func (r *Return) GetTotalItemsCount() int {
	count := 0
	for _, item := range r.Items {
		count += item.QuantityReturned
	}
	return count
}

// IsExchange checks if this is an exchange return
func (r *Return) IsExchange() bool {
	return r.Type == TypeExchange
}

// IsRefund checks if this is a refund return
func (r *Return) IsRefund() bool {
	return r.Type == TypeRefund
}

// GetReturnAge returns the age of the return in days
func (r *Return) GetReturnAge() int {
	return int(time.Since(r.CreatedAt).Hours() / 24)
}

// IsExpired checks if the return has expired based on reason's max return days
func (r *Return) IsExpired() bool {
	if r.Reason == nil {
		return false
	}
	return r.GetReturnAge() > r.Reason.MaxReturnDays
}

// GetReturnSummary returns a summary of the return
func (r *Return) GetReturnSummary() map[string]interface{} {
	return map[string]interface{}{
		"return_number":   r.ReturnNumber,
		"status":          r.Status,
		"type":            r.Type,
		"total_refund":    r.TotalRefund,
		"currency":        r.Currency,
		"items_count":     r.GetTotalItemsCount(),
		"created_at":      r.CreatedAt,
		"is_exchange":     r.IsExchange(),
		"is_expired":      r.IsExpired(),
		"can_be_approved": r.CanBeApproved(),
	}
}

// Business Logic Methods for ReturnItem

// CalculateRefundAmount calculates the refund amount for this item
func (ri *ReturnItem) CalculateRefundAmount() {
	ri.RefundAmount = ri.UnitPrice * float64(ri.QuantityReturned)
}

// GetReturnPercentage returns the percentage of quantity being returned
func (ri *ReturnItem) GetReturnPercentage() float64 {
	if ri.QuantityOrdered == 0 {
		return 0
	}
	return (float64(ri.QuantityReturned) / float64(ri.QuantityOrdered)) * 100
}

// IsPartialReturn checks if this is a partial return
func (ri *ReturnItem) IsPartialReturn() bool {
	return ri.QuantityReturned < ri.QuantityOrdered
}

// IsFullReturn checks if this is a full return
func (ri *ReturnItem) IsFullReturn() bool {
	return ri.QuantityReturned == ri.QuantityOrdered
}

// HasExchange checks if this item has an exchange
func (ri *ReturnItem) HasExchange() bool {
	return ri.ExchangeProductID != nil && ri.ExchangeQuantity > 0
}

// Business Logic Methods for ReturnReason

// IsValidForReturn checks if this reason is valid for creating returns
func (rr *ReturnReason) IsValidForReturn() bool {
	return rr.IsActive
}

// CalculateRestockingFee calculates the restocking fee for an amount
func (rr *ReturnReason) CalculateRestockingFee(amount float64) float64 {
	return amount * rr.RestockingFeeRate
}

// GetMaxReturnDate returns the maximum return date from a given order date
func (rr *ReturnReason) GetMaxReturnDate(orderDate time.Time) time.Time {
	return orderDate.AddDate(0, 0, rr.MaxReturnDays)
}

// IsReturnAllowed checks if return is allowed for a given order date
func (rr *ReturnReason) IsReturnAllowed(orderDate time.Time) bool {
	return time.Since(orderDate).Hours()/24 <= float64(rr.MaxReturnDays)
}