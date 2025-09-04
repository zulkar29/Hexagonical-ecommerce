package order

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository interface defines order data operations
type Repository interface {
	// Order operations
	CreateOrder(order *Order) (*Order, error)
	GetOrderByID(tenantID, orderID uuid.UUID) (*Order, error)
	GetOrderByNumber(tenantID uuid.UUID, orderNumber string) (*Order, error)
	UpdateOrder(order *Order) (*Order, error)
	ListOrders(tenantID uuid.UUID, filter OrderFilter, offset, limit int) ([]*Order, int64, error)
	DeleteOrder(tenantID, orderID uuid.UUID) error
	
	// Order statistics
	GetOrderStats(tenantID uuid.UUID) (map[string]interface{}, error)
	GetOrdersByCustomer(tenantID, customerID uuid.UUID) ([]*Order, error)
	GetOrdersByDateRange(tenantID uuid.UUID, start, end time.Time) ([]*Order, error)
	GetTopCustomers(tenantID uuid.UUID, limit int) ([]map[string]interface{}, error)
	
	// Order item operations
	CreateOrderItem(item *OrderItem) (*OrderItem, error)
	UpdateOrderItem(item *OrderItem) (*OrderItem, error)
	DeleteOrderItem(orderID, itemID uuid.UUID) error
	
	// Order history operations
	CreateOrderHistory(history *OrderHistory) (*OrderHistory, error)
	GetOrderHistory(tenantID, orderID uuid.UUID) ([]*OrderHistory, error)
	GetOrderTimeline(tenantID, orderID uuid.UUID) ([]*OrderHistory, error)
	
	// Utility operations
	GetLowStockAlert(tenantID uuid.UUID, threshold int) ([]*OrderItem, error)
}

// repository handles order data operations
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new order repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// OrderFilter represents filters for order queries
type OrderFilter struct {
	Status            *OrderStatus       `json:"status,omitempty"`
	PaymentStatus     *PaymentStatus     `json:"payment_status,omitempty"`
	FulfillmentStatus *FulfillmentStatus `json:"fulfillment_status,omitempty"`
	CustomerEmail     string             `json:"customer_email,omitempty"`
	OrderNumber       string             `json:"order_number,omitempty"`
	DateFrom          *time.Time         `json:"date_from,omitempty"`
	DateTo            *time.Time         `json:"date_to,omitempty"`
	MinAmount         *float64           `json:"min_amount,omitempty"`
	MaxAmount         *float64           `json:"max_amount,omitempty"`
	Search            string             `json:"search,omitempty"`
}

// CreateOrder saves a new order to the database
func (r *repository) CreateOrder(order *Order) (*Order, error) {
	if err := r.db.Create(order).Error; err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return order, nil
}

// GetOrderByID retrieves an order by ID with items preloaded
func (r *repository) GetOrderByID(tenantID, orderID uuid.UUID) (*Order, error) {
	var order Order
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, orderID).
		Preload("Items").
		First(&order).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	return &order, nil
}

// GetOrderByNumber retrieves an order by order number
func (r *repository) GetOrderByNumber(tenantID uuid.UUID, orderNumber string) (*Order, error) {
	var order Order
	err := r.db.Where("tenant_id = ? AND order_number = ?", tenantID, orderNumber).
		Preload("Items").
		First(&order).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	return &order, nil
}

// UpdateOrder updates an existing order
func (r *repository) UpdateOrder(order *Order) (*Order, error) {
	if err := r.db.Save(order).Error; err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return order, nil
}

// ListOrders retrieves orders with filtering and pagination
func (r *repository) ListOrders(tenantID uuid.UUID, filter OrderFilter, offset, limit int) ([]*Order, int64, error) {
	query := r.db.Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	
	if filter.PaymentStatus != nil {
		query = query.Where("payment_status = ?", *filter.PaymentStatus)
	}
	
	if filter.FulfillmentStatus != nil {
		query = query.Where("fulfillment_status = ?", *filter.FulfillmentStatus)
	}
	
	if filter.CustomerEmail != "" {
		query = query.Where("customer_email ILIKE ?", "%"+filter.CustomerEmail+"%")
	}
	
	if filter.OrderNumber != "" {
		query = query.Where("order_number ILIKE ?", "%"+filter.OrderNumber+"%")
	}
	
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	
	if filter.MinAmount != nil {
		query = query.Where("total_amount >= ?", *filter.MinAmount)
	}
	
	if filter.MaxAmount != nil {
		query = query.Where("total_amount <= ?", *filter.MaxAmount)
	}
	
	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where(
			"LOWER(order_number) LIKE ? OR LOWER(customer_email) LIKE ? OR LOWER(shipping_first_name) LIKE ? OR LOWER(shipping_last_name) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}
	
	// Count total records
	var total int64
	if err := query.Model(&Order{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}
	
	// Get orders with pagination
	var orders []*Order
	err := query.Preload("Items").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	
	return orders, total, nil
}

// GetOrderStats retrieves order statistics for a tenant
func (r *repository) GetOrderStats(tenantID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total orders count
	var totalOrders int64
	if err := r.db.Model(&Order{}).Where("tenant_id = ?", tenantID).Count(&totalOrders).Error; err != nil {
		return nil, fmt.Errorf("failed to count total orders: %w", err)
	}
	stats["total_orders"] = totalOrders
	
	// Total revenue
	var totalRevenue float64
	if err := r.db.Model(&Order{}).
		Where("tenant_id = ? AND payment_status = ?", tenantID, PaymentPaid).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalRevenue).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total revenue: %w", err)
	}
	stats["total_revenue"] = totalRevenue
	
	// Average order value
	var avgOrderValue float64
	if totalOrders > 0 {
		avgOrderValue = totalRevenue / float64(totalOrders)
	}
	stats["average_order_value"] = avgOrderValue
	
	// Orders by status
	var statusCounts []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	if err := r.db.Model(&Order{}).
		Where("tenant_id = ?", tenantID).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get orders by status: %w", err)
	}
	stats["orders_by_status"] = statusCounts
	
	// Recent orders (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var recentOrders int64
	if err := r.db.Model(&Order{}).
		Where("tenant_id = ? AND created_at >= ?", tenantID, thirtyDaysAgo).
		Count(&recentOrders).Error; err != nil {
		return nil, fmt.Errorf("failed to count recent orders: %w", err)
	}
	stats["recent_orders"] = recentOrders
	
	return stats, nil
}

// GetOrdersByCustomer retrieves orders for a specific customer
func (r *repository) GetOrdersByCustomer(tenantID, customerID uuid.UUID) ([]*Order, error) {
	var orders []*Order
	err := r.db.Where("tenant_id = ? AND user_id = ?", tenantID, customerID).
		Preload("Items").
		Order("created_at DESC").
		Find(&orders).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get customer orders: %w", err)
	}
	
	return orders, nil
}

// GetOrdersByDateRange retrieves orders within a date range
func (r *repository) GetOrdersByDateRange(tenantID uuid.UUID, start, end time.Time) ([]*Order, error) {
	var orders []*Order
	err := r.db.Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, start, end).
		Preload("Items").
		Order("created_at DESC").
		Find(&orders).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by date range: %w", err)
	}
	
	return orders, nil
}

// DeleteOrder soft deletes an order
func (r *repository) DeleteOrder(tenantID, orderID uuid.UUID) error {
	result := r.db.Where("tenant_id = ? AND id = ?", tenantID, orderID).Delete(&Order{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete order: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	
	return nil
}

// CreateOrderItem adds an item to an order
func (r *repository) CreateOrderItem(item *OrderItem) (*OrderItem, error) {
	if err := r.db.Create(item).Error; err != nil {
		return nil, fmt.Errorf("failed to create order item: %w", err)
	}
	return item, nil
}

// UpdateOrderItem updates an order item
func (r *repository) UpdateOrderItem(item *OrderItem) (*OrderItem, error) {
	if err := r.db.Save(item).Error; err != nil {
		return nil, fmt.Errorf("failed to update order item: %w", err)
	}
	return item, nil
}

// DeleteOrderItem removes an item from an order
func (r *repository) DeleteOrderItem(orderID, itemID uuid.UUID) error {
	result := r.db.Where("order_id = ? AND id = ?", orderID, itemID).Delete(&OrderItem{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete order item: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("order item not found")
	}
	
	return nil
}

// GetTopCustomers retrieves top customers by order value
func (r *repository) GetTopCustomers(tenantID uuid.UUID, limit int) ([]map[string]interface{}, error) {
	var customers []map[string]interface{}
	
	err := r.db.Model(&Order{}).
		Select("user_id, customer_email, COUNT(*) as order_count, SUM(total_amount) as total_spent").
		Where("tenant_id = ? AND payment_status = ?", tenantID, PaymentPaid).
		Group("user_id, customer_email").
		Order("total_spent DESC").
		Limit(limit).
		Scan(&customers).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get top customers: %w", err)
	}
	
	return customers, nil
}

// GetLowStockAlert retrieves orders with items that have low stock
func (r *repository) GetLowStockAlert(tenantID uuid.UUID, threshold int) ([]*OrderItem, error) {
	var items []*OrderItem
	
	// This would need to join with product inventory
	// For now, return empty slice as products module handles inventory
	return items, nil
}

// CreateOrderHistory creates a new order history entry
func (r *repository) CreateOrderHistory(history *OrderHistory) (*OrderHistory, error) {
	if err := r.db.Create(history).Error; err != nil {
		return nil, fmt.Errorf("failed to create order history: %w", err)
	}
	return history, nil
}

// GetOrderHistory retrieves all history entries for an order
func (r *repository) GetOrderHistory(tenantID, orderID uuid.UUID) ([]*OrderHistory, error) {
	var history []*OrderHistory
	err := r.db.Where("order_id = ?", orderID).
		Joins("JOIN orders ON order_histories.order_id = orders.id").
		Where("orders.tenant_id = ?", tenantID).
		Order("order_histories.created_at ASC").
		Find(&history).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}
	
	return history, nil
}

// GetOrderTimeline retrieves order timeline (same as history but with different semantic meaning)
func (r *repository) GetOrderTimeline(tenantID, orderID uuid.UUID) ([]*OrderHistory, error) {
	return r.GetOrderHistory(tenantID, orderID)
}
