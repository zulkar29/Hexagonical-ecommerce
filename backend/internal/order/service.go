package order

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service handles order business logic
type Service struct {
	repo *Repository
	db   *gorm.DB
}

// NewService creates a new order service
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
		db:   repo.db,
	}
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	CustomerEmail    string              `json:"customer_email" validate:"required,email"`
	CustomerPhone    string              `json:"customer_phone,omitempty"`
	ShippingAddress  Address             `json:"shipping_address" validate:"required"`
	BillingAddress   *Address            `json:"billing_address,omitempty"`
	Items            []CreateOrderItem   `json:"items" validate:"required,min=1"`
	PaymentMethod    string              `json:"payment_method,omitempty"`
	PaymentGateway   string              `json:"payment_gateway,omitempty"`
	CouponCode       string              `json:"coupon_code,omitempty"`
	Notes            string              `json:"notes,omitempty"`
	Currency         string              `json:"currency,omitempty"`
}

// CreateOrderItem represents an item in the order creation request
type CreateOrderItem struct {
	ProductID uuid.UUID  `json:"product_id" validate:"required"`
	VariantID *uuid.UUID `json:"variant_id,omitempty"`
	Quantity  int        `json:"quantity" validate:"required,min=1"`
}

// UpdateOrderStatusRequest represents the request to update order status
type UpdateOrderStatusRequest struct {
	Status        OrderStatus `json:"status" validate:"required"`
	TrackingNumber *string    `json:"tracking_number,omitempty"`
	TrackingURL    *string    `json:"tracking_url,omitempty"`
	Notes         string      `json:"notes,omitempty"`
}

// CreateOrder creates a new order
func (s *Service) CreateOrder(tenantID, userID uuid.UUID, req CreateOrderRequest) (*Order, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generate order number
	orderNumber := s.generateOrderNumber(tenantID)

	// Create order
	order := &Order{
		ID:              uuid.New(),
		TenantID:        tenantID,
		UserID:          userID,
		OrderNumber:     orderNumber,
		Status:          StatusPending,
		CustomerEmail:   req.CustomerEmail,
		CustomerPhone:   req.CustomerPhone,
		ShippingAddress: req.ShippingAddress,
		PaymentStatus:   PaymentPending,
		PaymentMethod:   req.PaymentMethod,
		PaymentGateway:  req.PaymentGateway,
		FulfillmentStatus: FulfillmentPending,
		Currency:        "BDT",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Set billing address (default to shipping if not provided)
	if req.BillingAddress != nil {
		order.BillingAddress = *req.BillingAddress
	} else {
		order.BillingAddress = req.ShippingAddress
	}

	// Set currency if provided
	if req.Currency != "" {
		order.Currency = req.Currency
	}

	// Create order in database
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Process order items
	var orderItems []OrderItem
	subtotal := 0.0

	for _, itemReq := range req.Items {
		// Get product details (you'd need to integrate with product service)
		item := OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   itemReq.ProductID,
			VariantID:   itemReq.VariantID,
			Quantity:    itemReq.Quantity,
			// These would be fetched from product service
			ProductName: "Product Name", // TODO: Get from product service
			ProductSKU:  "SKU123",       // TODO: Get from product service
			UnitPrice:   29.99,          // TODO: Get from product service
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		item.UpdateTotal()
		subtotal += item.TotalPrice

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}

		orderItems = append(orderItems, item)
	}

	// Calculate totals
	order.SubtotalAmount = subtotal
	order.TaxAmount = s.calculateTax(order.SubtotalAmount, order.ShippingAddress.Country)
	order.ShippingAmount = s.calculateShipping(order.SubtotalAmount, order.ShippingAddress.Country)
	order.DiscountAmount = 0.0 // TODO: Apply coupon if provided

	order.CalculateTotal()

	// Update order with calculated amounts
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update order totals: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit order creation: %w", err)
	}

	// Load items for response
	order.Items = orderItems

	return order, nil
}

// GetOrder retrieves an order by ID
func (s *Service) GetOrder(tenantID uuid.UUID, orderID string) (*Order, error) {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	return s.repo.GetOrderByID(tenantID, id)
}

// GetOrderByNumber retrieves an order by order number
func (s *Service) GetOrderByNumber(tenantID uuid.UUID, orderNumber string) (*Order, error) {
	return s.repo.GetOrderByNumber(tenantID, orderNumber)
}

// UpdateOrderStatus updates the status of an order
func (s *Service) UpdateOrderStatus(tenantID uuid.UUID, orderID string, req UpdateOrderStatusRequest) (*Order, error) {
	order, err := s.GetOrder(tenantID, orderID)
	if err != nil {
		return nil, err
	}

	// Validate status transition
	if !s.isValidStatusTransition(order.Status, req.Status) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", order.Status, req.Status)
	}

	// Update status
	order.Status = req.Status
	order.UpdatedAt = time.Now()

	// Update fulfillment status based on order status
	switch req.Status {
	case StatusProcessing:
		order.FulfillmentStatus = FulfillmentPicked
	case StatusShipped:
		order.FulfillmentStatus = FulfillmentShipped
		if req.TrackingNumber != nil {
			order.TrackingNumber = *req.TrackingNumber
		}
		if req.TrackingURL != nil {
			order.TrackingURL = *req.TrackingURL
		}
		now := time.Now()
		order.ShippedAt = &now
	case StatusDelivered:
		order.FulfillmentStatus = FulfillmentDelivered
		now := time.Now()
		order.DeliveredAt = &now
	}

	return s.repo.UpdateOrder(order)
}

// CancelOrder cancels an order
func (s *Service) CancelOrder(tenantID uuid.UUID, orderID string, reason string) (*Order, error) {
	order, err := s.GetOrder(tenantID, orderID)
	if err != nil {
		return nil, err
	}

	if !order.IsCancellable() {
		return nil, fmt.Errorf("order cannot be cancelled in current status: %s", order.Status)
	}

	order.Status = StatusCancelled
	order.UpdatedAt = time.Now()

	// TODO: Handle refund if payment was processed
	// TODO: Restore inventory
	// TODO: Send cancellation email

	return s.repo.UpdateOrder(order)
}

// ListOrders retrieves orders with filtering and pagination
func (s *Service) ListOrders(tenantID uuid.UUID, filter OrderFilter, page, limit int) ([]*Order, int64, error) {
	offset := (page - 1) * limit
	return s.repo.ListOrders(tenantID, filter, offset, limit)
}

// GetOrderStats retrieves order statistics
func (s *Service) GetOrderStats(tenantID uuid.UUID) (map[string]interface{}, error) {
	return s.repo.GetOrderStats(tenantID)
}

// ProcessPayment processes payment for an order
func (s *Service) ProcessPayment(tenantID uuid.UUID, orderID string) (*Order, error) {
	order, err := s.GetOrder(tenantID, orderID)
	if err != nil {
		return nil, err
	}

	if order.PaymentStatus != PaymentPending {
		return nil, fmt.Errorf("order payment is not pending")
	}

	// TODO: Integrate with payment gateway
	// For now, simulate successful payment
	order.PaymentStatus = PaymentPaid
	order.Status = StatusConfirmed
	order.UpdatedAt = time.Now()

	return s.repo.UpdateOrder(order)
}

// RefundOrder processes a refund for an order
func (s *Service) RefundOrder(tenantID uuid.UUID, orderID string, amount float64) (*Order, error) {
	order, err := s.GetOrder(tenantID, orderID)
	if err != nil {
		return nil, err
	}

	if !order.IsRefundable() {
		return nil, fmt.Errorf("order is not eligible for refund")
	}

	if amount > order.TotalAmount {
		return nil, fmt.Errorf("refund amount cannot exceed order total")
	}

	// TODO: Process refund through payment gateway
	// For now, update status
	order.PaymentStatus = PaymentRefunded
	order.UpdatedAt = time.Now()

	return s.repo.UpdateOrder(order)
}

// GetCustomerOrders retrieves orders for a specific customer
func (s *Service) GetCustomerOrders(tenantID, customerID uuid.UUID) ([]*Order, error) {
	return s.repo.GetOrdersByCustomer(tenantID, customerID)
}

// TrackOrder provides tracking information for an order
func (s *Service) TrackOrder(tenantID uuid.UUID, orderNumber string) (map[string]interface{}, error) {
	order, err := s.repo.GetOrderByNumber(tenantID, orderNumber)
	if err != nil {
		return nil, err
	}

	tracking := map[string]interface{}{
		"order_number":       order.OrderNumber,
		"status":            order.Status,
		"fulfillment_status": order.FulfillmentStatus,
		"tracking_number":   order.TrackingNumber,
		"tracking_url":      order.TrackingURL,
		"shipped_at":        order.ShippedAt,
		"delivered_at":      order.DeliveredAt,
		"estimated_delivery": s.calculateEstimatedDelivery(order),
	}

	return tracking, nil
}

// Helper methods

// generateOrderNumber generates a unique order number
func (s *Service) generateOrderNumber(tenantID uuid.UUID) string {
	timestamp := time.Now().Unix()
	random := rand.Intn(1000)
	return fmt.Sprintf("ORD-%d-%03d", timestamp, random)
}

// calculateTax calculates tax amount based on location
func (s *Service) calculateTax(subtotal float64, country string) float64 {
	// Bangladesh VAT is typically 15%
	if country == "BD" {
		return subtotal * 0.15
	}
	// Default no tax for other countries
	return 0.0
}

// calculateShipping calculates shipping cost
func (s *Service) calculateShipping(subtotal float64, country string) float64 {
	// Free shipping for orders over 1000 BDT
	if subtotal >= 1000 {
		return 0.0
	}
	
	// Standard shipping rates
	if country == "BD" {
		return 60.0 // 60 BDT for Bangladesh
	}
	
	return 200.0 // International shipping
}

// isValidStatusTransition checks if status transition is valid
func (s *Service) isValidStatusTransition(from, to OrderStatus) bool {
	validTransitions := map[OrderStatus][]OrderStatus{
		StatusPending:    {StatusConfirmed, StatusCancelled},
		StatusConfirmed:  {StatusProcessing, StatusCancelled},
		StatusProcessing: {StatusShipped, StatusCancelled},
		StatusShipped:    {StatusDelivered, StatusReturned},
		StatusDelivered:  {StatusReturned},
		StatusCancelled:  {}, // No transitions from cancelled
		StatusReturned:   {}, // No transitions from returned
	}

	allowedStatuses, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == to {
			return true
		}
	}

	return false
}

// calculateEstimatedDelivery calculates estimated delivery date
func (s *Service) calculateEstimatedDelivery(order *Order) *time.Time {
	if order.Status == StatusDelivered {
		return order.DeliveredAt
	}

	if order.ShippedAt != nil {
		// Estimate 3-5 days from ship date for domestic, 7-14 for international
		var days int
		if order.ShippingAddress.Country == "BD" {
			days = 3
		} else {
			days = 7
		}
		estimated := order.ShippedAt.AddDate(0, 0, days)
		return &estimated
	}

	// If not shipped yet, estimate from creation date
	var days int
	if order.ShippingAddress.Country == "BD" {
		days = 5 // 2 days processing + 3 days shipping
	} else {
		days = 10 // 3 days processing + 7 days shipping
	}
	estimated := order.CreatedAt.AddDate(0, 0, days)
	return &estimated
}

// ValidateOrder validates order data
func (s *Service) ValidateOrder(order *Order) error {
	if order.CustomerEmail == "" {
		return fmt.Errorf("customer email is required")
	}

	if !order.ShippingAddress.IsComplete() {
		return fmt.Errorf("shipping address is incomplete")
	}

	if len(order.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}

	if order.TotalAmount <= 0 {
		return fmt.Errorf("order total must be greater than zero")
	}

	return nil
}
