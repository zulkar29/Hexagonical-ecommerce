package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	sharedErrors "ecommerce-saas/internal/shared/errors"
)

// CartRepository interface to access cart data
type CartRepository interface {
	GetCartByID(tenantID, cartID uuid.UUID) (*Cart, error)
}

// Cart represents cart data needed for order conversion
type Cart struct {
	ID                 uuid.UUID     `json:"id"`
	TenantID           uuid.UUID     `json:"tenant_id"`
	CustomerID         *uuid.UUID    `json:"customer_id"`
	SessionID          string        `json:"session_id"`
	Items              []CartItem    `json:"items"`
	ShippingAddress    Address       `json:"shipping_address"`
	BillingAddress     Address       `json:"billing_address"`
	ShippingMethodID   *uuid.UUID    `json:"shipping_method_id"`
	CouponCode         string        `json:"coupon_code"`
	DiscountAmount     float64       `json:"discount_amount"`
	ShippingCost       float64       `json:"shipping_cost"`
	TaxAmount          float64       `json:"tax_amount"`
	Total              float64       `json:"total"`
	Currency           string        `json:"currency"`
	Status             string        `json:"status"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}

// CartItem represents a cart item
type CartItem struct {
	ID               uuid.UUID          `json:"id"`
	CartID           uuid.UUID          `json:"cart_id"`
	ProductID        uuid.UUID          `json:"product_id"`
	VariantID        *uuid.UUID         `json:"variant_id"`
	ProductName      string             `json:"product_name"`
	ProductSKU       string             `json:"product_sku"`
	Quantity         int                `json:"quantity"`
	UnitPrice        float64            `json:"unit_price"`
	Total            float64            `json:"total"`
	Customizations   map[string]interface{} `json:"customizations"`
	Notes            string             `json:"notes"`
}

// OrderFromCartResult represents the result of cart to order conversion
type OrderFromCartResult struct {
	OrderID      uuid.UUID `json:"order_id"`
	OrderNumber  string    `json:"order_number"`
	Total        float64   `json:"total"`
	ItemCount    int       `json:"item_count"`
}

// CreateOrderFromCart creates an order from a cart
func (s *Service) CreateOrderFromCart(ctx context.Context, tenantID, cartID uuid.UUID) (*OrderFromCartResult, error) {
	// This would need a cart repository to get cart data
	// For now, we'll create a basic order structure
	// In a real implementation, you'd inject CartRepository in the service

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the order
	order := &Order{
		ID:                uuid.New(),
		TenantID:          tenantID,
		OrderNumber:       s.generateOrderNumber(tenantID),
		Status:            StatusPending,
		PaymentStatus:     PaymentPending,
		FulfillmentStatus: FulfillmentPending,
		Currency:          "BDT", // Default currency
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// In a real implementation, you would:
	// 1. Get cart data from CartRepository
	// 2. Convert cart items to order items
	// 3. Apply discounts and calculate totals
	// 4. Reserve inventory for each item
	// 5. Set shipping and billing addresses
	// 6. Calculate taxes and shipping

	// For now, create a placeholder implementation
	order.Total = 0.0           // Would be calculated from cart
	order.SubTotal = 0.0        // Would be calculated from cart items
	order.ShippingCost = 0.0    // Would be from cart shipping calculation
	order.TaxAmount = 0.0       // Would be calculated based on address
	order.DiscountAmount = 0.0  // Would be from applied coupons

	// Save order to database
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit order creation: %w", err)
	}

	return &OrderFromCartResult{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		Total:       order.Total,
		ItemCount:   len(order.Items), // Would be from converted cart items
	}, nil
}

// generateOrderNumber generates a unique order number for the tenant
func (s *Service) generateOrderNumber(tenantID uuid.UUID) string {
	// Simple implementation - in production you'd want more sophisticated numbering
	timestamp := time.Now().Unix()
	return fmt.Sprintf("ORD-%d-%d", timestamp, rand.Intn(1000))
}

// CreateOrderFromCartWithItems creates an order from cart data with full implementation
func (s *Service) CreateOrderFromCartWithItems(ctx context.Context, tenantID uuid.UUID, cart *Cart) (*OrderFromCartResult, error) {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the order
	order := &Order{
		ID:                uuid.New(),
		TenantID:          tenantID,
		CustomerID:        cart.CustomerID,
		OrderNumber:       s.generateOrderNumber(tenantID),
		Status:            StatusPending,
		PaymentStatus:     PaymentPending,
		FulfillmentStatus: FulfillmentPending,
		Currency:          cart.Currency,
		ShippingAddress:   cart.ShippingAddress,
		BillingAddress:    cart.BillingAddress,
		SubTotal:          0.0, // Will be calculated from items
		ShippingCost:      cart.ShippingCost,
		TaxAmount:         cart.TaxAmount,
		DiscountAmount:    cart.DiscountAmount,
		Total:             cart.Total,
		CouponCode:        cart.CouponCode,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Convert cart items to order items
	orderItems := make([]OrderItem, 0, len(cart.Items))
	subtotal := 0.0

	for _, cartItem := range cart.Items {
		// Reserve stock for this item
		if err := s.productService.ReserveStock(tenantID, cartItem.ProductID, cartItem.Quantity); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to reserve stock for product %s: %w", cartItem.ProductID, err)
		}

		orderItem := OrderItem{
			ID:             uuid.New(),
			OrderID:        order.ID,
			ProductID:      cartItem.ProductID,
			VariantID:      cartItem.VariantID,
			ProductName:    cartItem.ProductName,
			ProductSKU:     cartItem.ProductSKU,
			Quantity:       cartItem.Quantity,
			UnitPrice:      cartItem.UnitPrice,
			Total:          cartItem.Total,
			Customizations: cartItem.Customizations,
			Notes:          cartItem.Notes,
		}

		orderItems = append(orderItems, orderItem)
		subtotal += cartItem.Total
	}

	order.Items = orderItems
	order.SubTotal = subtotal

	// Save order to database
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()

		// Restore inventory for all items if order creation fails
		for _, item := range orderItems {
			if restoreErr := s.productService.RestoreStock(tenantID, item.ProductID, item.Quantity); restoreErr != nil {
				log.Printf("Warning: Failed to restore stock for product %s: %v", item.ProductID, restoreErr)
			}
		}

		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		// Restore inventory for all items if commit fails
		for _, item := range orderItems {
			if restoreErr := s.productService.RestoreStock(tenantID, item.ProductID, item.Quantity); restoreErr != nil {
				log.Printf("Warning: Failed to restore stock for product %s: %v", item.ProductID, restoreErr)
			}
		}
		return nil, fmt.Errorf("failed to commit order creation: %w", err)
	}

	// Send order confirmation notification
	if s.notificationService != nil {
		// This would be implemented in the notification service integration
		log.Printf("Order confirmation would be sent for order %s", order.OrderNumber)
	}

	return &OrderFromCartResult{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		Total:       order.Total,
		ItemCount:   len(order.Items),
	}, nil
}