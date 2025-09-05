package order

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// CreateOrderItem represents an item to be added to an order
type CreateOrderItem struct {
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	ProductName string    `json:"product_name"`
	ProductSKU  string    `json:"product_sku"`
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

// Note: DTOs removed - using domain models directly

// CreateOrder creates a new order
func (s *Service) CreateOrder(ctx context.Context, tenantID uuid.UUID, order *Order) (*Order, error) {
	// Validate order
	if err := s.ValidateOrder(order); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Set order defaults
	order.ID = uuid.New()
	order.TenantID = tenantID
	order.OrderNumber = s.generateOrderNumber(tenantID)
	order.Status = StatusPending
	order.PaymentStatus = PaymentPending
	order.FulfillmentStatus = FulfillmentPending
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Set default currency if not provided
	if order.Currency == "" {
		order.Currency = "BDT"
	}

	// Set billing address to shipping address if not provided
	if order.BillingAddress.Address1 == "" {
		order.BillingAddress = order.ShippingAddress
	}

	// Create order in database
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Process order items
	subtotal := 0.0

	for i, item := range order.Items {
		// Get product details from product service
		product, err := s.productService.GetProduct(tenantID, item.ProductID.String())
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		// Check if product is active
		if product.Status != "active" {
			tx.Rollback()
			return nil, fmt.Errorf("product %s is not available for purchase", product.Name)
		}

		// Check inventory availability
		if product.Inventory < item.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient inventory for product %s. Available: %d, Requested: %d", product.Name, product.Inventory, item.Quantity)
		}

		// Reserve inventory for this item
		if err := s.inventoryService.ReserveStock(ctx, tenantID, item.ProductID, item.Quantity); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to reserve inventory for product %s: %w", product.Name, err)
		}

		// Update order item with product details
		order.Items[i].ID = uuid.New()
		order.Items[i].OrderID = order.ID
		order.Items[i].ProductName = product.Name
		order.Items[i].ProductSKU = product.SKU
		order.Items[i].UnitPrice = product.Price
		order.Items[i].CreatedAt = time.Now()
		order.Items[i].UpdatedAt = time.Now()

		order.Items[i].UpdateTotal()
		subtotal += order.Items[i].TotalPrice

		if err := tx.Create(&order.Items[i]).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Calculate totals
	order.SubtotalAmount = subtotal
	order.TaxAmount = s.calculateTax(order.SubtotalAmount, order.ShippingAddress.Country)
	order.ShippingAmount = s.calculateShipping(order.SubtotalAmount, order.ShippingAddress.Country)
	order.DiscountAmount = 0.0

	order.CalculateTotal()

	// Update order with calculated amounts
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update order totals: %w", err)
	}

	// Create payment if payment gateway is specified
	if order.PaymentGateway != "" {
		_, err := s.createPayment(tenantID, order.UserID, order)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create payment: %w", err)
		}
		// Store payment information in order
		// order.PaymentIntentID = &paymentResp.PaymentID // Field doesn't exist in Order struct
		if err := tx.Save(order).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update order with payment info: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit order creation: %w", err)
	}

	// Send order confirmation notification
	go s.sendOrderConfirmationNotification(ctx, order)

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
func (s *Service) UpdateOrderStatus(ctx context.Context, tenantID, orderID uuid.UUID, status OrderStatus, trackingNumber, trackingURL, notes string) (*Order, error) {
	// Get existing order
	order, err := s.repository.GetOrderByID(tenantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Check if status change is valid
	if !s.isValidStatusTransition(order.Status, status) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", order.Status, status)
	}

	// Update order status
	oldStatus := order.Status
	order.Status = status
	order.UpdatedAt = time.Now()

	// Update tracking information if provided
	if trackingNumber != "" {
		order.TrackingNumber = trackingNumber
	}
	if trackingURL != "" {
		order.TrackingURL = trackingURL
	}

	// Update fulfillment status based on order status
	switch status {
	case StatusProcessing:
		order.FulfillmentStatus = FulfillmentPacked
	case StatusShipped:
		order.FulfillmentStatus = FulfillmentShipped
	case StatusDelivered:
		order.FulfillmentStatus = FulfillmentDelivered
	case StatusCancelled:
		order.FulfillmentStatus = FulfillmentPending
	}

	// Save updated order
	if _, err := s.repository.UpdateOrder(order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Create order history entry
	history := &OrderHistory{
		ID:          uuid.New(),
		OrderID:     order.ID,
		TenantID:    order.TenantID,
		FromStatus:  oldStatus,
		ToStatus:    status,
		Action:      "status_changed",
		Description: fmt.Sprintf("Order status changed from %s to %s", oldStatus, status),
		Notes:       notes,
		CreatedAt:   time.Now(),
	}

	if _, err := s.repository.CreateOrderHistory(history); err != nil {
		// Log error but don't fail the status update
		// s.logger.Error("Failed to create order history", "error", err)
	}

	// Send notifications based on status
	switch status {
	case StatusShipped:
		go s.sendOrderShippingNotification(ctx, order)
	case StatusCancelled:
		go s.sendOrderCancellationNotification(ctx, order)
	}

	return order, nil
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
		if err := s.inventoryService.RestoreStock(context.Background(), tenantID, item.ProductID, item.Quantity); err != nil {
			// Log error but don't fail the cancellation
			// In production, this should be logged properly
			fmt.Printf("Warning: failed to restore inventory for product %s: %v\n", item.ProductID, err)
		}
	}

	// TODO: Handle refund if payment was processed

	// Send order cancellation notification
	go s.sendOrderCancellationNotification(context.Background(), order)

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
func (s *Service) RefundOrder(ctx context.Context, tenantID, orderID uuid.UUID, paymentID string, amount float64, reason string) (*Payment, error) {
	// Get order
	order, err := s.repository.GetOrderByID(tenantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Check if order can be refunded
	if !order.IsRefundable() {
		return nil, fmt.Errorf("order %s is not refundable", order.OrderNumber)
	}

	// Process refund through gateway
	if err := s.paymentService.RefundPayment(ctx, tenantID, paymentID, amount, reason); err != nil {
		return nil, fmt.Errorf("failed to process refund: %w", err)
	}

	// Update order status
	order.PaymentStatus = PaymentRefunded
	order.Status = StatusCancelled
	order.UpdatedAt = time.Now()

	if _, err := s.repository.UpdateOrder(order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Create order history entry
	history := &OrderHistory{
		ID:          uuid.New(),
		OrderID:     order.ID,
		TenantID:    order.TenantID,
		ToStatus:    order.Status,
		Action:      "refunded",
		Description: fmt.Sprintf("Order refunded - Amount: %.2f, Reason: %s", amount, reason),
		CreatedAt:   time.Now(),
	}

	if _, err := s.repository.CreateOrderHistory(history); err != nil {
		// Log error but don't fail the refund
		// s.logger.Error("Failed to create order history", "error", err)
	}

	// Create a payment response for the refund
	payment := &Payment{
		ID:       uuid.New(),
		TenantID: tenantID,
		OrderID:  orderID,
		Amount:   amount,
		Status:   "refunded",
	}

	return payment, nil
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
	// Convert filters to OrderFilter and get orders
	filter := OrderFilter{} // TODO: Convert map filters to OrderFilter
	orders, _, err := s.repository.ListOrders(tenantID, filter, 0, 1000) // Get up to 1000 orders for export
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch orders: %w", err)
	}

	// Convert []*Order to []Order for export functions
	orderList := make([]Order, len(orders))
	for i, order := range orders {
		orderList[i] = *order
	}

	switch format {
	case "csv":
		return s.exportOrdersToCSV(orderList)
	case "excel":
		return s.exportOrdersToExcel(orderList)
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
	if addr.Address1 == "" {
		return ""
	}
	return fmt.Sprintf("%s, %s, %s, %s %s, %s",
		addr.Address1, addr.City, addr.State, addr.PostalCode, addr.Country, addr.Phone)
}

// ImportOrders imports orders from CSV format
func (s *Service) ImportOrders(ctx context.Context, tenantID uuid.UUID, file io.Reader) (totalRecords, successfulImports, failedImports int, errors []string, err error) {
	records, err := s.importOrdersFromCSV(file)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	totalRecords = len(records)
	errors = []string{}

	for i, record := range records {
		order, err := s.parseOrderFromCSVRecord(record)
		if err != nil {
			failedImports++
			errors = append(errors, fmt.Sprintf("Row %d: %v", i+1, err))
			continue
		}

		order.TenantID = tenantID
		order.ID = uuid.New()
		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()

		if _, err := s.repository.CreateOrder(order); err != nil {
			failedImports++
			errors = append(errors, fmt.Sprintf("Row %d: Failed to save order: %v", i+1, err))
			continue
		}

		successfulImports++
	}

	return totalRecords, successfulImports, failedImports, errors, nil
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
func (s *Service) importOrdersFromCSV(file io.Reader) ([][]string, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	// Skip header row if present
	if len(records) > 0 {
		return records[1:], nil
	}

	return records, nil
}

// parseOrderFromCSVRecord parses an order from a CSV record
func (s *Service) parseOrderFromCSVRecord(record []string) (*Order, error) {
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
		return Address{Address1: addrStr}
	}

	return Address{
		Address1:   parts[0],
		City:       parts[1],
		State:      parts[2],
		PostalCode: strings.Fields(parts[3])[0], // Extract postal code
		Country:    strings.Fields(parts[3])[1], // Extract country
		Phone:      parts[5],
	}
}

// sendOrderConfirmationNotification sends order confirmation notification
func (s *Service) sendOrderConfirmationNotification(ctx context.Context, order *Order) error {
	return s.notificationService.SendEmail(ctx, order.TenantID, []string{order.CustomerEmail}, 
		fmt.Sprintf("Order Confirmation - %s", order.OrderNumber),
		"order_confirmation",
		"text/html",
		map[string]interface{}{
			"order":        order,
			"customer":     order.CustomerEmail,
			"order_number": order.OrderNumber,
			"total":        order.TotalAmount,
		},
		"order_confirmation")
}

// sendOrderCancellationNotification sends order cancellation notification
func (s *Service) sendOrderCancellationNotification(ctx context.Context, order *Order) error {
	return s.notificationService.SendEmail(ctx, order.TenantID, []string{order.CustomerEmail},
		fmt.Sprintf("Order Cancelled - %s", order.OrderNumber),
		"order_cancellation",
		"text/html",
		map[string]interface{}{
			"order":        order,
			"customer":     order.CustomerEmail,
			"order_number": order.OrderNumber,
		},
		"order_cancellation")
}

// sendOrderShippingNotification sends order shipping notification
func (s *Service) sendOrderShippingNotification(ctx context.Context, order *Order) error {
	return s.notificationService.SendEmail(ctx, order.TenantID, []string{order.CustomerEmail},
		fmt.Sprintf("Order Shipped - %s", order.OrderNumber),
		"order_shipping",
		"text/html",
		map[string]interface{}{
			"order":          order,
			"customer":       order.CustomerEmail,
			"order_number":   order.OrderNumber,
			"tracking_number": order.TrackingNumber,
			"tracking_url":   order.TrackingURL,
		},
		"order_shipping")
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
	validation, err := s.discountService.ValidateDiscountCode(ctx, tenantID, couponCode, &userID, order.CustomerEmail, order.SubtotalAmount, totalQuantity, productIDs, []string{})
	if err != nil {
		return 0, fmt.Errorf("failed to validate discount code: %w", err)
	}

	if !validation.Valid {
		return 0, fmt.Errorf("discount code is not valid: %s", validation.Message)
	}

	// Apply discount
	application, err := s.discountService.ApplyDiscount(ctx, tenantID, couponCode, order.ID, &userID, order.CustomerEmail, order.SubtotalAmount, totalQuantity, productIDs, []string{}, "", "")
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
	// Create payment through payment service
	paymentResp, err := s.paymentService.CreatePayment(context.Background(), tenantID, order.ID.String(), order.TotalAmount, order.Currency, order.PaymentGateway, order.PaymentMethod, order.CustomerEmail, order.CustomerPhone, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return paymentResp, nil
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
func (s *Service) BulkUpdateOrders(ctx context.Context, tenantID uuid.UUID, orderIDs []uuid.UUID, action string, data map[string]interface{}) (int, int, []string, error) {
	if len(orderIDs) == 0 {
		return 0, 0, nil, fmt.Errorf("no order IDs provided")
	}

	if err := s.validateBulkUpdateAction(action, data); err != nil {
		return 0, 0, nil, fmt.Errorf("invalid bulk update request: %w", err)
	}

	var successfulUpdates, failedUpdates int
	var errors []string

	// Process each order
	for _, orderID := range orderIDs {
		var err error
		switch action {
		case "update_status":
			err = s.bulkUpdateOrderStatus(ctx, tenantID, orderID, data)
		case "cancel":
			err = s.bulkCancelOrder(ctx, tenantID, orderID)
		case "refund":
			err = s.bulkRefundOrder(ctx, tenantID, orderID, data)
		default:
			err = fmt.Errorf("unsupported action: %s", action)
		}

		if err != nil {
			failedUpdates++
			errors = append(errors, fmt.Sprintf("Order %s: %v", orderID, err))
		} else {
			successfulUpdates++
		}
	}

	return successfulUpdates, failedUpdates, errors, nil
}

// validateBulkUpdateAction validates the bulk update action and data
func (s *Service) validateBulkUpdateAction(action string, data map[string]interface{}) error {
	switch action {
	case "update_status":
		if status, ok := data["status"]; !ok || status == "" {
			return fmt.Errorf("status is required for update_status action")
		}
	case "cancel":
		// No additional validation needed
	case "refund":
		if amount, ok := data["amount"]; !ok {
			return fmt.Errorf("amount is required for refund action")
		} else if amountFloat, ok := amount.(float64); !ok || amountFloat <= 0 {
			return fmt.Errorf("amount must be a positive number")
		}
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
	return nil
}

// bulkUpdateOrderStatus updates the status of a single order in bulk operation
func (s *Service) bulkUpdateOrderStatus(ctx context.Context, tenantID, orderID uuid.UUID, data map[string]interface{}) error {
	status, ok := data["status"].(string)
	if !ok {
		return fmt.Errorf("invalid status type")
	}

	trackingNumber, _ := data["tracking_number"].(string)
	trackingURL, _ := data["tracking_url"].(string)
	notes, _ := data["notes"].(string)

	_, err := s.UpdateOrderStatus(ctx, tenantID, orderID, OrderStatus(status), trackingNumber, trackingURL, notes)
	return err
}

// bulkCancelOrder cancels a single order in bulk operation
func (s *Service) bulkCancelOrder(ctx context.Context, tenantID, orderID uuid.UUID) error {
	_, err := s.CancelOrder(tenantID, orderID.String(), "Bulk cancellation")
	return err
}

// bulkRefundOrder refunds a single order in bulk operation
func (s *Service) bulkRefundOrder(ctx context.Context, tenantID, orderID uuid.UUID, data map[string]interface{}) error {
	amount, ok := data["amount"].(float64)
	if !ok {
		return fmt.Errorf("invalid amount type")
	}

	reason, _ := data["reason"].(string)
	paymentID, _ := data["payment_id"].(string)

	_, err := s.RefundOrder(ctx, tenantID, orderID, paymentID, amount, reason)
	return err
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
		FromStatus:    fromStatus,
		ToStatus:      toStatus,
		Action:        "status_change",
		Description:   fmt.Sprintf("Status changed from %s to %s", fromStatus, toStatus),
		Reason:        reason,
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
		FromPaymentStatus:   fromPaymentStatus,
		ToPaymentStatus:     toPaymentStatus,
		Action:              "payment_status_change",
		Description:         fmt.Sprintf("Payment status changed from %s to %s", fromPaymentStatus, toPaymentStatus),
		Reason:              reason,
		ChangedBy:           &changedBy,
		ChangedByType:       "system",
		CreatedAt:           time.Now(),
	}
	
	_, err := s.repository.CreateOrderHistory(entry)
	return err
}

// UpdateOrder updates an existing order
func (s *Service) UpdateOrder(ctx context.Context, tenantID uuid.UUID, orderID string, req *UpdateOrderRequest) (*Order, error) {
	// Get existing order
	order, err := s.GetOrder(tenantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	
	// Check if order is editable
	if !order.IsEditable() {
		return nil, fmt.Errorf("order cannot be edited in current status: %s", order.Status)
	}
	
	// Update fields if provided
	if req.CustomerEmail != nil {
		order.CustomerEmail = *req.CustomerEmail
	}
	if req.CustomerPhone != nil {
		order.CustomerPhone = *req.CustomerPhone
	}
	if req.ShippingAddress != nil {
		order.ShippingAddress = *req.ShippingAddress
	}
	if req.BillingAddress != nil {
		order.BillingAddress = *req.BillingAddress
	}
	if req.Notes != nil {
		order.Notes = *req.Notes
	}
	
	order.UpdatedAt = time.Now()
	
	// Update in repository
	updatedOrder, err := s.repository.UpdateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	
	// Add history entry
	s.AddOrderHistoryEntry(order.ID, "order_updated", "Order details updated", uuid.Nil, "user", nil)
	
	return updatedOrder, nil
}

// DeleteOrder deletes a single order
func (s *Service) DeleteOrder(ctx context.Context, tenantID uuid.UUID, orderID uuid.UUID) error {
	// Get order to validate it exists and belongs to tenant
	order, err := s.repository.GetOrderByID(tenantID, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}
	
	// Check if order can be deleted (only pending/draft orders)
	if order.Status != OrderStatusPending && order.Status != OrderStatusDraft {
		return fmt.Errorf("cannot delete order in status %s", order.Status)
	}
	
	// Add history entry before deletion
	s.AddOrderHistoryEntry(order.ID, "order_deleted", "Order deleted", uuid.Nil, "user", nil)
	
	// Delete order
	if err := s.repository.DeleteOrder(tenantID, order.ID); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	
	return nil
}

// BulkDeleteOrders deletes multiple orders
func (s *Service) BulkDeleteOrders(ctx context.Context, tenantID uuid.UUID, orderIDs []string, reason string) (int, int, []string, error) {
	if len(orderIDs) == 0 {
		return 0, 0, nil, fmt.Errorf("no order IDs provided")
	}
	
	var successfulDeletes, failedDeletes int
	var errors []string
	
	// Process each order
	for _, orderID := range orderIDs {
		// Get order to validate it exists and belongs to tenant
		order, err := s.GetOrder(tenantID, orderID)
		if err != nil {
			failedDeletes++
			errors = append(errors, fmt.Sprintf("Order %s: %v", orderID, err))
			continue
		}
		
		// Check if order can be deleted (only pending/draft orders)
		if order.Status != OrderStatusPending && order.Status != OrderStatusDraft {
			failedDeletes++
			errors = append(errors, fmt.Sprintf("Order %s: cannot delete order in status %s", orderID, order.Status))
			continue
		}
		
		// Add history entry before deletion
		s.AddOrderHistoryEntry(order.ID, "order_deleted", fmt.Sprintf("Order deleted: %s", reason), uuid.Nil, "user", nil)
		
		// Delete order
		if err := s.repository.DeleteOrder(tenantID, order.ID); err != nil {
			failedDeletes++
			errors = append(errors, fmt.Sprintf("Order %s: failed to delete: %v", orderID, err))
			continue
		}
		
		successfulDeletes++
	}
	
	return successfulDeletes, failedDeletes, errors, nil
}
