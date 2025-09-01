package entities

// TODO: Implement Order entity during development
// Keep it simple - core order functionality only

import (
	"time"
)

// Order represents a customer order
type Order struct {
	// TODO: Implement basic fields only
	// ID          string    `json:"id"`
	// TenantID    string    `json:"tenant_id"`
	// CustomerID  string    `json:"customer_id"`
	// Items       []OrderItem `json:"items"`
	// Total       float64   `json:"total"`
	// Status      string    `json:"status"` // pending, paid, shipped, delivered
	// CreatedAt   time.Time `json:"created_at"`
	// UpdatedAt   time.Time `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	// TODO: Keep simple
	// ProductID string  `json:"product_id"`
	// Name      string  `json:"name"`
	// Price     float64 `json:"price"`
	// Quantity  int     `json:"quantity"`
}

// Business methods (keep minimal)
// func (o *Order) CalculateTotal() float64
// func (o *Order) UpdateStatus(status string) error