package product

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// InventoryService provides inventory management functionality
// This is a wrapper around product repository methods to satisfy order module dependencies
type InventoryService struct {
	repo Repository
}

// NewInventoryService creates a new inventory service
func NewInventoryService(repo Repository) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

// ReserveStock reserves inventory for an order
// In a simple implementation, this decrements the inventory quantity
func (s *InventoryService) ReserveStock(tenantID, productID uuid.UUID, quantity int) error {
	// Get current product
	product, err := s.repo.FindProductByID(tenantID, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Check if we can decrement inventory
	if err := product.DecrementInventory(quantity); err != nil {
		return fmt.Errorf("cannot reserve stock: %w", err)
	}

	// Update inventory in database
	newQuantity := product.InventoryQuantity
	if err := s.repo.UpdateInventory(tenantID, productID, newQuantity); err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	return nil
}

// RestoreStock restores inventory (e.g., when an order is cancelled)
func (s *InventoryService) RestoreStock(tenantID, productID uuid.UUID, quantity int) error {
	// Get current product
	product, err := s.repo.FindProductByID(tenantID, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Only restore if tracking quantity
	if product.TrackQuantity {
		newQuantity := product.InventoryQuantity + quantity
		if err := s.repo.UpdateInventory(tenantID, productID, newQuantity); err != nil {
			return fmt.Errorf("failed to restore inventory: %w", err)
		}
	}

	return nil
}

// UpdateInventory directly updates inventory quantity
func (s *InventoryService) UpdateInventory(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID, quantity int) error {
	if err := s.repo.UpdateInventory(tenantID, productID, quantity); err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	return nil
}