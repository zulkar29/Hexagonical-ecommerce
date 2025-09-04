package order

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ProductService interface for product operations
type ProductService interface {
	GetProduct(tenantID uuid.UUID, id string) (*Product, error)
	GetProductBySlug(tenantID uuid.UUID, slug string) (*Product, error)
}

// Product represents a product for order integration
type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	Inventory   int       `json:"inventory"`
}

// DiscountService interface for discount operations
type DiscountService interface {
	ValidateDiscountCode(ctx context.Context, req ValidateDiscountRequest) (*DiscountValidation, error)
	ApplyDiscount(ctx context.Context, req ApplyDiscountRequest) (*DiscountApplication, error)
	RemoveDiscount(ctx context.Context, tenantID uuid.UUID, orderID uuid.UUID) error
}

// PaymentService interface for payment operations
type PaymentService interface {
	CreatePayment(tenantID uuid.UUID, req *CreatePaymentRequest) (*CreatePaymentResponse, error)
	ProcessPayment(tenantID uuid.UUID, req *ProcessPaymentRequest) error
	RefundPayment(tenantID uuid.UUID, req *RefundPaymentRequest) error
	GetPayment(tenantID uuid.UUID, paymentID string) (*Payment, error)
}

// Payment related structs
type CreatePaymentRequest struct {
	OrderID         string  `json:"order_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,min=0.01"`
	Currency        string  `json:"currency" validate:"required,len=3"`
	Gateway         string  `json:"gateway" validate:"required"`
	PaymentMethodID string  `json:"payment_method_id,omitempty"`
	CustomerEmail   string  `json:"customer_email" validate:"required,email"`
	CustomerPhone   string  `json:"customer_phone,omitempty"`
	ReturnURL       string  `json:"return_url,omitempty"`
}

type CreatePaymentResponse struct {
	PaymentID      string `json:"payment_id"`
	Status         string `json:"status"`
	PaymentURL     string `json:"payment_url,omitempty"`
	SessionKey     string `json:"session_key,omitempty"`
	GatewayPageURL string `json:"gateway_page_url,omitempty"`
}

type ProcessPaymentRequest struct {
	PaymentID       string                 `json:"payment_id" validate:"required"`
	Gateway         string                 `json:"gateway" validate:"required"`
	GatewayResponse map[string]interface{} `json:"gateway_response"`
}

type RefundPaymentRequest struct {
	PaymentID string  `json:"payment_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,min=0.01"`
	Reason    string  `json:"reason,omitempty"`
}

type Payment struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	OrderID           uuid.UUID  `json:"order_id"`
	UserID            uuid.UUID  `json:"user_id"`
	PaymentIntentID   string     `json:"payment_intent_id"`
	PaymentMethodID   string     `json:"payment_method_id"`
	Amount            float64    `json:"amount"`
	Currency          string     `json:"currency"`
	Status            string     `json:"status"`
	Gateway           string     `json:"gateway"`
	GatewayResponse   string     `json:"gateway_response"`
	FailureReason     string     `json:"failure_reason"`
	RefundedAmount    float64    `json:"refunded_amount"`
	RefundedAt        *time.Time `json:"refunded_at"`
	ProcessedAt       *time.Time `json:"processed_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// Discount-related structs for order integration
type ValidateDiscountRequest struct {
	Code           string     `json:"code"`
	CustomerID     *uuid.UUID `json:"customer_id"`
	CustomerEmail  string     `json:"customer_email"`
	OrderAmount    float64    `json:"order_amount"`
	ItemQuantity   int        `json:"item_quantity"`
	ProductIDs     []string   `json:"product_ids"`
	CategoryIDs    []string   `json:"category_ids"`
}

type DiscountValidation struct {
	Valid           bool    `json:"valid"`
	DiscountAmount  float64 `json:"discount_amount"`
	Message         string  `json:"message"`
	CanStack        bool    `json:"can_stack"`
}

type ApplyDiscountRequest struct {
	TenantID       uuid.UUID  `json:"tenant_id"`
	Code           string     `json:"code"`
	OrderID        uuid.UUID  `json:"order_id"`
	CustomerID     *uuid.UUID `json:"customer_id"`
	CustomerEmail  string     `json:"customer_email"`
	OrderAmount    float64    `json:"order_amount"`
	ItemQuantity   int        `json:"item_quantity"`
	ProductIDs     []string   `json:"product_ids"`
	CategoryIDs    []string   `json:"category_ids"`
	IPAddress      string     `json:"ip_address"`
	UserAgent      string     `json:"user_agent"`
}

type DiscountApplication struct {
	Applied        bool    `json:"applied"`
	DiscountAmount float64 `json:"discount_amount"`
	Message        string  `json:"message"`
}

// Service handles order business logic
type Service struct {
	repository          Repository
	db                  *gorm.DB
	productService      ProductService
	discountService     DiscountService
	paymentService      PaymentService
	inventoryService    InventoryService
	notificationService NotificationService
}

// NewService creates a new order service
func NewService(repo Repository, db *gorm.DB, productService ProductService, discountService DiscountService, paymentService PaymentService, inventoryService InventoryService, notificationService NotificationService) *Service {
	return &Service{
		repository:          repo,
		db:                  db,
		productService:      productService,
		discountService:     discountService,
		paymentService:      paymentService,
		inventoryService:    inventoryService,
		notificationService: notificationService,
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
		// Get product details from product service
		product, err := s.productService.GetProduct(tenantID, itemReq.ProductID.String())
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to get product %s: %w", itemReq.ProductID, err)
		}

		// Check if product is active
		if product.Status != "active" {
			tx.Rollback()
			return nil, fmt.Errorf("product %s is not available for purchase", product.Name)
		}

		// Check inventory availability
		if product.Inventory < itemReq.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient inventory for product %s. Available: %d, Requested: %d", product.Name, product.Inventory, itemReq.Quantity)
		}

		// Reserve inventory for this item
		if err := s.inventoryService.ReserveStock(tenantID, itemReq.ProductID, itemReq.Quantity); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to reserve inventory for product %s: %w", product.Name, err)
		}

		item := OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   itemReq.ProductID,
			VariantID:   itemReq.VariantID,
			Quantity:    itemReq.Quantity,
			ProductName: product.Name,
			ProductSKU:  product.SKU,
			UnitPrice:   product.Price,
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
	order.DiscountAmount = 0.0

	// Apply discount if coupon code is provided
	if req.CouponCode != "" {
		discountAmount, err := s.applyDiscount(tenantID, userID, order, req.CouponCode, req.Items)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to apply discount: %w", err)
		}
		order.DiscountAmount = discountAmount
	}

	order.CalculateTotal()

	// Update order with calculated amounts
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update order totals: %w", err)
	}

	// Create payment if payment gateway is specified
	if req.PaymentGateway != "" {
		paymentResp, err := s.createPayment(tenantID, userID, order)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create payment: %w", err)
		}
		// Store payment information in order
		order.PaymentIntentID = &paymentResp.PaymentID
		if err := tx.Save(order).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update order with payment info: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit order creation: %w", err)
	}

	// Load items for response
	order.Items = orderItems

	// Send order confirmation notification
	go s.sendOrderConfirmationNotification(tenantID, order)

	return order, nil
}

// GetOrder retrieves an order by ID
func (s *Service) GetOrder(tenantID uuid.UUID, orderID string) (*Order, error) {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	return s.repository.GetOrderByID(tenantID, id)
}

// GetOrderByNumber retrieves an order by order number
func (s *Service) GetOrderByNumber(tenantID uuid.UUID, orderNumber string) (*Order, error) {
	return s.repository.GetOrderByNumber(tenantID, orderNumber)
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
		// Send shipping notification
		go s.sendOrderShippingNotification(tenantID, order)
	case StatusDelivered:
		order.FulfillmentStatus = FulfillmentDelivered
		now := time.Now()
		order.DeliveredAt = &now
	}

	return s.repository.UpdateOrder(order)
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

	// Restore inventory for cancelled order
	for _, item := range order.Items {
		if err := s.inventoryService.RestoreStock(tenantID, item.ProductID, item.Quantity); err != nil {
			// Log error but don't fail the cancellation
			// In production, this should be logged properly
			fmt.Printf("Warning: failed to restore inventory for product %s: %v\n", item.ProductID, err)
		}
	}

	// TODO: Handle refund if payment was processed

	// Send order cancellation notification
	go s.sendOrderCancellationNotification(tenantID, order, reason)

	return s.repository.UpdateOrder(order)
}

// ListOrders retrieves orders with filtering and pagination
func (s *Service) ListOrders(tenantID uuid.UUID, filter OrderFilter, page, limit int) ([]*Order, int64, error) {
	offset := (page - 1) * limit
	return s.repository.ListOrders(tenantID, filter, offset, limit)
}

// GetOrderStats retrieves order statistics
func (s *Service) GetOrderStats(tenantID uuid.UUID) (map[string]interface{}, error) {
	return s.repository.GetOrderStats(tenantID)
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

	return s.repository.UpdateOrder(order)
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
	return s.repository.GetOrdersByCustomer(tenantID, customerID)
}

// TrackOrder provides tracking information for an order
func (s *Service) TrackOrder(tenantID uuid.UUID, orderNumber string) (map[string]interface{}, error) {
	order, err := s.repository.GetOrderByNumber(tenantID, orderNumber)
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

// ExportOrders exports orders to CSV format
func (s *Service) ExportOrders(tenantID uuid.UUID, format string, filters map[string]interface{}) ([]byte, string, error) {
	// Get orders based on filters
	orders, err := s.repository.ListOrders(tenantID, filters)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch orders: %w", err)
	}

	switch format {
	case "csv":
		return s.exportOrdersToCSV(orders)
	case "excel":
		return s.exportOrdersToExcel(orders)
	default:
		return nil, "", fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportOrdersToCSV exports orders to CSV format
func (s *Service) exportOrdersToCSV(orders []Order) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV header
	header := []string{
		"Order Number", "Customer Email", "Customer Phone", "Status", "Payment Status",
		"Subtotal", "Tax", "Shipping", "Discount", "Total", "Currency",
		"Payment Gateway", "Payment Method", "Shipping Address", "Billing Address",
		"Created At", "Updated At",
	}
	if err := writer.Write(header); err != nil {
		return nil, "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write order data
	for _, order := range orders {
		record := []string{
			order.OrderNumber,
			order.CustomerEmail,
			order.CustomerPhone,
			string(order.Status),
			string(order.PaymentStatus),
			fmt.Sprintf("%.2f", order.SubtotalAmount),
			fmt.Sprintf("%.2f", order.TaxAmount),
			fmt.Sprintf("%.2f", order.ShippingAmount),
			fmt.Sprintf("%.2f", order.DiscountAmount),
			fmt.Sprintf("%.2f", order.TotalAmount),
			order.Currency,
			order.PaymentGateway,
			order.PaymentMethod,
			s.formatAddress(order.ShippingAddress),
			s.formatAddress(order.BillingAddress),
			order.CreatedAt.Format("2006-01-02 15:04:05"),
			order.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return nil, "", fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", fmt.Errorf("CSV writer error: %w", err)
	}

	filename := fmt.Sprintf("orders_export_%s.csv", time.Now().Format("20060102_150405"))
	return buf.Bytes(), filename, nil
}

// exportOrdersToExcel exports orders to Excel format
func (s *Service) exportOrdersToExcel(orders []Order) ([]byte, string, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: failed to close Excel file: %v\n", err)
		}
	}()

	sheetName := "Orders"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Excel sheet: %w", err)
	}

	// Set headers
	headers := []string{
		"Order Number", "Customer Email", "Customer Phone", "Status", "Payment Status",
		"Subtotal", "Tax", "Shipping", "Discount", "Total", "Currency",
		"Payment Gateway", "Payment Method", "Shipping Address", "Billing Address",
		"Created At", "Updated At",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		if err := file.SetCellValue(sheetName, cell, header); err != nil {
			return nil, "", fmt.Errorf("failed to set Excel header: %w", err)
		}
	}

	// Set data
	for i, order := range orders {
		row := i + 2 // Start from row 2 (after header)
		data := []interface{}{
			order.OrderNumber,
			order.CustomerEmail,
			order.CustomerPhone,
			string(order.Status),
			string(order.PaymentStatus),
			order.SubtotalAmount,
			order.TaxAmount,
			order.ShippingAmount,
			order.DiscountAmount,
			order.TotalAmount,
			order.Currency,
			order.PaymentGateway,
			order.PaymentMethod,
			s.formatAddress(order.ShippingAddress),
			s.formatAddress(order.BillingAddress),
			order.CreatedAt.Format("2006-01-02 15:04:05"),
			order.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		for j, value := range data {
			cell := fmt.Sprintf("%c%d", 'A'+j, row)
			if err := file.SetCellValue(sheetName, cell, value); err != nil {
				return nil, "", fmt.Errorf("failed to set Excel cell value: %w", err)
			}
		}
	}

	file.SetActiveSheet(index)

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		return nil, "", fmt.Errorf("failed to write Excel file: %w", err)
	}

	filename := fmt.Sprintf("orders_export_%s.xlsx", time.Now().Format("20060102_150405"))
	return buf.Bytes(), filename, nil
}

// formatAddress formats an address for export
func (s *Service) formatAddress(addr Address) string {
	if addr.Street == "" {
		return ""
	}
	return fmt.Sprintf("%s, %s, %s, %s %s, %s",
		addr.Street, addr.City, addr.State, addr.PostalCode, addr.Country, addr.Phone)
}

// ImportOrders imports orders from CSV format
func (s *Service) ImportOrders(tenantID uuid.UUID, data []byte, format string) (*ImportResult, error) {
	switch format {
	case "csv":
		return s.importOrdersFromCSV(tenantID, data)
	default:
		return nil, fmt.Errorf("unsupported import format: %s", format)
	}
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	TotalRecords    int      `json:"total_records"`
	SuccessCount    int      `json:"success_count"`
	ErrorCount      int      `json:"error_count"`
	Errors          []string `json:"errors"`
	ImportedOrderIDs []string `json:"imported_order_ids"`
}

// importOrdersFromCSV imports orders from CSV data
func (s *Service) importOrdersFromCSV(tenantID uuid.UUID, data []byte) (*ImportResult, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV must contain at least a header and one data row")
	}

	result := &ImportResult{
		TotalRecords:     len(records) - 1, // Exclude header
		Errors:           []string{},
		ImportedOrderIDs: []string{},
	}

	// Skip header row
	for i, record := range records[1:] {
		rowNum := i + 2 // Account for header and 0-based index
		
		if len(record) < 17 {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: insufficient columns", rowNum))
			result.ErrorCount++
			continue
		}

		// Parse order data from CSV record
		order, err := s.parseOrderFromCSVRecord(tenantID, record, rowNum)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: %v", rowNum, err))
			result.ErrorCount++
			continue
		}

		// Create order
		if err := s.repository.CreateOrder(order); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: failed to create order: %v", rowNum, err))
			result.ErrorCount++
			continue
		}

		result.ImportedOrderIDs = append(result.ImportedOrderIDs, order.ID.String())
		result.SuccessCount++
	}

	return result, nil
}

// parseOrderFromCSVRecord parses an order from a CSV record
func (s *Service) parseOrderFromCSVRecord(tenantID uuid.UUID, record []string, rowNum int) (*Order, error) {
	// Parse amounts
	subtotal, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid subtotal amount: %v", err)
	}

	tax, err := strconv.ParseFloat(record[6], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid tax amount: %v", err)
	}

	shipping, err := strconv.ParseFloat(record[7], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid shipping amount: %v", err)
	}

	discount, err := strconv.ParseFloat(record[8], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid discount amount: %v", err)
	}

	total, err := strconv.ParseFloat(record[9], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid total amount: %v", err)
	}

	// Parse dates
	createdAt, err := time.Parse("2006-01-02 15:04:05", record[15])
	if err != nil {
		return nil, fmt.Errorf("invalid created_at date: %v", err)
	}

	updatedAt, err := time.Parse("2006-01-02 15:04:05", record[16])
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at date: %v", err)
	}

	// Parse addresses
	shippingAddr := s.parseAddressFromString(record[13])
	billingAddr := s.parseAddressFromString(record[14])

	order := &Order{
		ID:               uuid.New(),
		TenantID:         tenantID,
		OrderNumber:      record[0],
		CustomerEmail:    record[1],
		CustomerPhone:    record[2],
		Status:           OrderStatus(record[3]),
		PaymentStatus:    PaymentStatus(record[4]),
		SubtotalAmount:   subtotal,
		TaxAmount:        tax,
		ShippingAmount:   shipping,
		DiscountAmount:   discount,
		TotalAmount:      total,
		Currency:         record[10],
		PaymentGateway:   record[11],
		PaymentMethod:    record[12],
		ShippingAddress:  shippingAddr,
		BillingAddress:   billingAddr,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		Items:            []OrderItem{}, // Items would need separate import
	}

	return order, nil
}

// parseAddressFromString parses an address from a formatted string
func (s *Service) parseAddressFromString(addrStr string) Address {
	if addrStr == "" {
		return Address{}
	}

	// Simple parsing - in production, you might want more sophisticated parsing
	parts := strings.Split(addrStr, ", ")
	if len(parts) < 6 {
		return Address{Street: addrStr}
	}

	return Address{
		Street:     parts[0],
		City:       parts[1],
		State:      parts[2],
		PostalCode: strings.Fields(parts[3])[0], // Extract postal code
		Country:    strings.Fields(parts[3])[1], // Extract country
		Phone:      parts[5],
	}
}

// sendOrderConfirmationNotification sends order confirmation email
func (s *Service) sendOrderConfirmationNotification(tenantID uuid.UUID, order *Order) {
	if order.CustomerEmail == "" {
		return
	}

	// Prepare email variables
	variables := map[string]interface{}{
		"order_number":    order.OrderNumber,
		"customer_name":   order.CustomerName,
		"total_amount":    order.TotalAmount,
		"currency":        order.Currency,
		"order_date":      order.CreatedAt.Format("January 2, 2006"),
		"items":           order.Items,
		"shipping_address": order.ShippingAddress,
		"billing_address":  order.BillingAddress,
	}

	// Send email notification
	emailReq := &SendEmailRequest{
		To:         []string{order.CustomerEmail},
		Subject:    fmt.Sprintf("Order Confirmation - %s", order.OrderNumber),
		TemplateID: "order_confirmation",
		Variables:  variables,
	}

	if err := s.notificationService.SendEmail(tenantID, emailReq); err != nil {
		// Log error but don't fail the order creation
		fmt.Printf("Warning: failed to send order confirmation email: %v\n", err)
	}
}

// sendOrderCancellationNotification sends order cancellation email
func (s *Service) sendOrderCancellationNotification(tenantID uuid.UUID, order *Order, reason string) {
	if order.CustomerEmail == "" {
		return
	}

	// Prepare email variables
	variables := map[string]interface{}{
		"order_number":      order.OrderNumber,
		"customer_name":     order.CustomerName,
		"total_amount":      order.TotalAmount,
		"currency":          order.Currency,
		"cancellation_date": time.Now().Format("January 2, 2006"),
		"reason":            reason,
		"items":             order.Items,
	}

	// Send email notification
	emailReq := &SendEmailRequest{
		To:         []string{order.CustomerEmail},
		Subject:    fmt.Sprintf("Order Cancelled - %s", order.OrderNumber),
		TemplateID: "order_cancellation",
		Variables:  variables,
	}

	if err := s.notificationService.SendEmail(tenantID, emailReq); err != nil {
		// Log error but don't fail the cancellation
		fmt.Printf("Warning: failed to send order cancellation email: %v\n", err)
	}
}

// sendOrderShippingNotification sends order shipping email
func (s *Service) sendOrderShippingNotification(tenantID uuid.UUID, order *Order) {
	if order.CustomerEmail == "" {
		return
	}

	// Prepare email variables
	variables := map[string]interface{}{
		"order_number":     order.OrderNumber,
		"customer_name":    order.CustomerName,
		"tracking_number":  order.TrackingNumber,
		"tracking_url":     order.TrackingURL,
		"shipped_date":     time.Now().Format("January 2, 2006"),
		"items":            order.Items,
		"shipping_address": order.ShippingAddress,
	}

	// Send email notification
	emailReq := &SendEmailRequest{
		To:         []string{order.CustomerEmail},
		Subject:    fmt.Sprintf("Order Shipped - %s", order.OrderNumber),
		TemplateID: "order_shipped",
		Variables:  variables,
	}

	if err := s.notificationService.SendEmail(tenantID, emailReq); err != nil {
		// Log error but don't fail the status update
		fmt.Printf("Warning: failed to send order shipping email: %v\n", err)
	}
}

// applyDiscount validates and applies a discount to an order
func (s *Service) applyDiscount(tenantID, userID uuid.UUID, order *Order, couponCode string, items []CreateOrderItem) (float64, error) {
	ctx := context.Background()

	// Prepare product IDs for discount validation
	productIDs := make([]string, len(items))
	for i, item := range items {
		productIDs[i] = item.ProductID.String()
	}

	// Calculate total quantity
	totalQuantity := 0
	for _, item := range items {
		totalQuantity += item.Quantity
	}

	// Validate discount code
	validateReq := ValidateDiscountRequest{
		Code:          couponCode,
		CustomerID:    &userID,
		CustomerEmail: order.CustomerEmail,
		OrderAmount:   order.SubtotalAmount,
		ItemQuantity:  totalQuantity,
		ProductIDs:    productIDs,
		CategoryIDs:   []string{}, // TODO: Add category IDs if needed
	}

	validation, err := s.discountService.ValidateDiscountCode(ctx, validateReq)
	if err != nil {
		return 0, fmt.Errorf("failed to validate discount code: %w", err)
	}

	if !validation.Valid {
		return 0, fmt.Errorf("discount code is not valid: %s", validation.Message)
	}

	// Apply discount
	applyReq := ApplyDiscountRequest{
		TenantID:      tenantID,
		Code:          couponCode,
		OrderID:       order.ID,
		CustomerID:    &userID,
		CustomerEmail: order.CustomerEmail,
		OrderAmount:   order.SubtotalAmount,
		ItemQuantity:  totalQuantity,
		ProductIDs:    productIDs,
		CategoryIDs:   []string{}, // TODO: Add category IDs if needed
		IPAddress:     "",          // TODO: Extract from request context
		UserAgent:     "",          // TODO: Extract from request context
	}

	application, err := s.discountService.ApplyDiscount(ctx, applyReq)
	if err != nil {
		return 0, fmt.Errorf("failed to apply discount: %w", err)
	}

	if !application.Applied {
		return 0, fmt.Errorf("discount could not be applied: %s", application.Message)
	}

	return application.DiscountAmount, nil
}

// createPayment creates a payment for the order
func (s *Service) createPayment(tenantID, userID uuid.UUID, order *Order) (*CreatePaymentResponse, error) {
	// Prepare payment request
	paymentReq := &CreatePaymentRequest{
		OrderID:         order.ID.String(),
		Amount:          order.TotalAmount,
		Currency:        order.Currency,
		Gateway:         order.PaymentGateway,
		PaymentMethodID: order.PaymentMethod,
		CustomerEmail:   order.CustomerEmail,
		CustomerPhone:   order.CustomerPhone,
		ReturnURL:       "", // TODO: Configure return URL from settings
	}

	// Create payment through payment service
	paymentResp, err := s.paymentService.CreatePayment(tenantID, paymentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return paymentResp, nil
}

// ProcessPayment processes a payment for an order
func (s *Service) ProcessPayment(tenantID uuid.UUID, orderID string, req ProcessPaymentRequest) error {
	// Get order
	order, err := s.GetOrder(tenantID, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Validate order status
	if order.PaymentStatus == PaymentSucceeded {
		return fmt.Errorf("payment already processed for order %s", order.OrderNumber)
	}

	// Process payment through payment service
	err = s.paymentService.ProcessPayment(tenantID, &req)
	if err != nil {
		return fmt.Errorf("failed to process payment: %w", err)
	}

	// Update order payment status
	order.PaymentStatus = PaymentSucceeded
	order.UpdatedAt = time.Now()

	// If order is pending, move to confirmed
	if order.Status == StatusPending {
		order.Status = StatusConfirmed
	}

	// Save order
	if err := s.repository.UpdateOrder(order); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

// RefundOrder processes a refund for an order
func (s *Service) RefundOrder(tenantID uuid.UUID, orderID string, amount float64, reason string) error {
	// Get order
	order, err := s.GetOrder(tenantID, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Validate order can be refunded
	if order.PaymentStatus != PaymentSucceeded {
		return fmt.Errorf("cannot refund order with payment status: %s", order.PaymentStatus)
	}

	if order.PaymentIntentID == nil {
		return fmt.Errorf("no payment found for order %s", order.OrderNumber)
	}

	// Validate refund amount
	if amount <= 0 {
		return fmt.Errorf("refund amount must be greater than zero")
	}

	if amount > order.TotalAmount {
		return fmt.Errorf("refund amount cannot exceed order total")
	}

	// Process refund through payment service
	refundReq := &RefundPaymentRequest{
		PaymentID: *order.PaymentIntentID,
		Amount:    amount,
		Reason:    reason,
	}

	err = s.paymentService.RefundPayment(tenantID, refundReq)
	if err != nil {
		return fmt.Errorf("failed to process refund: %w", err)
	}

	// Update order status
	if amount >= order.TotalAmount {
		// Full refund
		order.PaymentStatus = PaymentRefunded
		order.Status = StatusCancelled
	} else {
		// Partial refund
		order.PaymentStatus = PaymentPartiallyRefunded
	}

	order.UpdatedAt = time.Now()

	// Save order
	if err := s.repository.UpdateOrder(order); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

// BulkUpdateOrdersRequest represents a bulk update request
type BulkUpdateOrdersRequest struct {
	OrderIDs []string `json:"order_ids" validate:"required,min=1"`
	Action   string   `json:"action" validate:"required,oneof=status cancel refund"`
	
	// For status updates
	Status         *OrderStatus `json:"status,omitempty"`
	TrackingNumber *string      `json:"tracking_number,omitempty"`
	
	// For cancellations
	CancelReason *string `json:"cancel_reason,omitempty"`
	
	// For refunds
	RefundAmount *float64 `json:"refund_amount,omitempty"`
	RefundReason *string  `json:"refund_reason,omitempty"`
}

// BulkUpdateResult represents the result of a bulk operation
type BulkUpdateResult struct {
	Total     int                    `json:"total"`
	Succeeded int                    `json:"succeeded"`
	Failed    int                    `json:"failed"`
	Errors    map[string]string      `json:"errors,omitempty"`
	Results   map[string]interface{} `json:"results,omitempty"`
}

// BulkUpdateOrders performs bulk operations on multiple orders
func (s *Service) BulkUpdateOrders(tenantID uuid.UUID, req BulkUpdateOrdersRequest) (*BulkUpdateResult, error) {
	result := &BulkUpdateResult{
		Total:   len(req.OrderIDs),
		Errors:  make(map[string]string),
		Results: make(map[string]interface{}),
	}

	// Validate request based on action
	if err := s.validateBulkUpdateRequest(req); err != nil {
		return nil, err
	}

	// Process each order
	for _, orderID := range req.OrderIDs {
		var err error
		
		switch req.Action {
		case "status":
			err = s.bulkUpdateOrderStatus(tenantID, orderID, *req.Status, req.TrackingNumber)
		case "cancel":
			err = s.bulkCancelOrder(tenantID, orderID, req.CancelReason)
		case "refund":
			err = s.bulkRefundOrder(tenantID, orderID, req.RefundAmount, req.RefundReason)
		default:
			err = fmt.Errorf("unsupported action: %s", req.Action)
		}

		if err != nil {
			result.Failed++
			result.Errors[orderID] = err.Error()
		} else {
			result.Succeeded++
			result.Results[orderID] = "success"
		}
	}

	return result, nil
}

// validateBulkUpdateRequest validates the bulk update request
func (s *Service) validateBulkUpdateRequest(req BulkUpdateOrdersRequest) error {
	switch req.Action {
	case "status":
		if req.Status == nil {
			return fmt.Errorf("status is required for status action")
		}
	case "cancel":
		// Cancel reason is optional
	case "refund":
		if req.RefundAmount == nil || *req.RefundAmount <= 0 {
			return fmt.Errorf("valid refund amount is required for refund action")
		}
	default:
		return fmt.Errorf("unsupported action: %s", req.Action)
	}
	return nil
}

// bulkUpdateOrderStatus updates order status in bulk
func (s *Service) bulkUpdateOrderStatus(tenantID uuid.UUID, orderID string, status OrderStatus, trackingNumber *string) error {
	updateReq := UpdateOrderStatusRequest{
		Status:         status,
		TrackingNumber: trackingNumber,
	}
	
	_, err := s.UpdateOrderStatus(tenantID, orderID, updateReq)
	return err
}

// bulkCancelOrder cancels an order in bulk
func (s *Service) bulkCancelOrder(tenantID uuid.UUID, orderID string, reason *string) error {
	cancelReason := "Bulk cancellation"
	if reason != nil {
		cancelReason = *reason
	}
	
	_, err := s.CancelOrder(tenantID, orderID, cancelReason)
	return err
}

// bulkRefundOrder processes a refund in bulk
func (s *Service) bulkRefundOrder(tenantID uuid.UUID, orderID string, amount *float64, reason *string) error {
	refundReason := "Bulk refund"
	if reason != nil {
		refundReason = *reason
	}
	
	refundAmount := *amount
	if amount == nil {
		// Get order to determine full refund amount
		order, err := s.GetOrder(tenantID, orderID)
		if err != nil {
			return err
		}
		refundAmount = order.TotalAmount
	}
	
	return s.RefundOrder(tenantID, orderID, refundAmount, refundReason)
}

// CreateOrderHistory creates a new order history entry
func (s *Service) CreateOrderHistory(tenantID uuid.UUID, orderID string, entry *OrderHistory) (*OrderHistory, error) {
	// Validate order exists and belongs to tenant
	order, err := s.GetOrder(tenantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	
	// Set order ID and timestamp
	entry.OrderID = order.ID
	entry.CreatedAt = time.Now()
	
	// Create history entry
	history, err := s.repository.CreateOrderHistory(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to create order history: %w", err)
	}
	
	return history, nil
}

// GetOrderTimeline retrieves the complete timeline/history for an order
func (s *Service) GetOrderTimeline(tenantID uuid.UUID, orderID string) ([]*OrderHistory, error) {
	// Validate order exists and belongs to tenant
	order, err := s.GetOrder(tenantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	
	// Get order timeline
	timeline, err := s.repository.GetOrderTimeline(tenantID, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order timeline: %w", err)
	}
	
	return timeline, nil
}

// AddOrderHistoryEntry is a helper method to add history entries during order operations
func (s *Service) AddOrderHistoryEntry(orderID uuid.UUID, action, description string, changedBy uuid.UUID, changedByType string, metadata map[string]interface{}) error {
	entry := &OrderHistory{
		OrderID:       orderID,
		Action:        action,
		Description:   description,
		ChangedBy:     &changedBy,
		ChangedByType: changedByType,
		Metadata:      metadata,
		CreatedAt:     time.Now(),
	}
	
	_, err := s.repository.CreateOrderHistory(entry)
	return err
}

// AddOrderStatusChangeHistory adds a history entry for status changes
func (s *Service) AddOrderStatusChangeHistory(orderID uuid.UUID, fromStatus, toStatus OrderStatus, changedBy uuid.UUID, reason string) error {
	entry := &OrderHistory{
		OrderID:       orderID,
		FromStatus:    &fromStatus,
		ToStatus:      &toStatus,
		Action:        "status_change",
		Description:   fmt.Sprintf("Status changed from %s to %s", fromStatus, toStatus),
		Reason:        &reason,
		ChangedBy:     &changedBy,
		ChangedByType: "user",
		CreatedAt:     time.Now(),
	}
	
	_, err := s.repository.CreateOrderHistory(entry)
	return err
}

// AddOrderPaymentChangeHistory adds a history entry for payment status changes
func (s *Service) AddOrderPaymentChangeHistory(orderID uuid.UUID, fromPaymentStatus, toPaymentStatus PaymentStatus, changedBy uuid.UUID, reason string) error {
	entry := &OrderHistory{
		OrderID:             orderID,
		FromPaymentStatus:   &fromPaymentStatus,
		ToPaymentStatus:     &toPaymentStatus,
		Action:              "payment_status_change",
		Description:         fmt.Sprintf("Payment status changed from %s to %s", fromPaymentStatus, toPaymentStatus),
		Reason:              &reason,
		ChangedBy:           &changedBy,
		ChangedByType:       "system",
		CreatedAt:           time.Now(),
	}
	
	_, err := s.repository.CreateOrderHistory(entry)
	return err
}
