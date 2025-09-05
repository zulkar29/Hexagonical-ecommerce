package wishlist

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for wishlist business logic
type Service interface {
	// Wishlist operations
	CreateWishlist(ctx context.Context, req CreateWishlistRequest) (*WishlistResponse, error)
	GetWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID) (*WishlistResponse, error)
	GetWishlistByShareToken(ctx context.Context, shareToken string) (*WishlistResponse, error)
	UpdateWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID, req UpdateWishlistRequest) (*WishlistResponse, error)
	DeleteWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID) error
	ListWishlists(ctx context.Context, tenantID uuid.UUID, filter WishlistFilter, limit, offset int) ([]WishlistResponse, int64, error)
	GetCustomerWishlists(ctx context.Context, tenantID, customerID uuid.UUID) ([]WishlistResponse, error)
	GetDefaultWishlist(ctx context.Context, tenantID, customerID uuid.UUID) (*WishlistResponse, error)
	
	// Wishlist item operations
	AddItem(ctx context.Context, tenantID uuid.UUID, req AddItemRequest) (*WishlistItemResponse, error)
	UpdateItem(ctx context.Context, tenantID, itemID uuid.UUID, req UpdateItemRequest) (*WishlistItemResponse, error)
	RemoveItem(ctx context.Context, tenantID, itemID uuid.UUID) error
	ListItems(ctx context.Context, tenantID uuid.UUID, filter WishlistItemFilter, limit, offset int) ([]WishlistItemResponse, int64, error)
	GetItem(ctx context.Context, tenantID, itemID uuid.UUID) (*WishlistItemResponse, error)
	
	// Wishlist management
	SetDefaultWishlist(ctx context.Context, tenantID, customerID, wishlistID uuid.UUID) error
	ClearWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID) error
	MoveItem(ctx context.Context, tenantID, itemID, targetWishlistID uuid.UUID) error
	CopyItem(ctx context.Context, tenantID, itemID, targetWishlistID uuid.UUID) error
	MergeWishlists(ctx context.Context, tenantID, sourceWishlistID, targetWishlistID uuid.UUID) error
	ShareWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID, isPublic bool) (string, error)
	
	// Bulk operations
	BulkAddItems(ctx context.Context, tenantID uuid.UUID, items []AddItemRequest) ([]WishlistItemResponse, error)
	BulkRemoveItems(ctx context.Context, tenantID uuid.UUID, itemIDs []uuid.UUID) error
	BulkUpdateItemPriority(ctx context.Context, tenantID uuid.UUID, updates map[uuid.UUID]int) error
	ReorderItems(ctx context.Context, tenantID, wishlistID uuid.UUID, itemOrder []uuid.UUID) error
	
	// Analytics and statistics
	GetWishlistStats(ctx context.Context, tenantID uuid.UUID) (*WishlistStats, error)
	GetMostWishedProducts(ctx context.Context, tenantID uuid.UUID, limit int) ([]ProductWishCount, error)
	GetCustomerActivity(ctx context.Context, tenantID, customerID uuid.UUID, days int) ([]WishlistActivity, error)
	GetPopularWishlists(ctx context.Context, tenantID uuid.UUID, limit int) ([]WishlistResponse, error)
	
	// Maintenance operations
	CleanupEmptyWishlists(ctx context.Context, tenantID uuid.UUID, olderThanDays int) (int64, error)
	CleanupOrphanedItems(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	repo Repository
	// External service interfaces would go here
	// productService ProductService
	// userService    UserService
}

// NewService creates a new wishlist service
func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

// Wishlist operations

// CreateWishlist creates a new wishlist
func (s *ServiceImpl) CreateWishlist(ctx context.Context, req CreateWishlistRequest) (*WishlistResponse, error) {
	// Validate request
	if err := s.validateCreateWishlistRequest(req); err != nil {
		return nil, err
	}
	
	// Check if name already exists for customer
	exists, err := s.repo.ExistsByName(ctx, req.TenantID, req.CustomerID, req.Name, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check wishlist name: %w", err)
	}
	if exists {
		return nil, ErrWishlistNameExists
	}
	
	// Check wishlist limit for customer
	count, err := s.repo.CountWishlistsByCustomer(ctx, req.TenantID, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to count customer wishlists: %w", err)
	}
	if count >= MaxWishlistsPerCustomer {
		return nil, ErrWishlistLimitExceeded
	}
	
	// Create wishlist
	wishlist := &Wishlist{
		ID:          uuid.New(),
		TenantID:    req.TenantID,
		CustomerID:  req.CustomerID,
		Name:        req.Name,
		Description: req.Description,
		IsDefault:   req.IsDefault,
		IsPublic:    req.IsPublic,
		ItemCount:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// If this is set as default, ensure no other default exists
	if req.IsDefault {
		if err := s.repo.SetDefaultWishlist(ctx, req.TenantID, req.CustomerID, wishlist.ID); err != nil {
			return nil, fmt.Errorf("failed to set default wishlist: %w", err)
		}
	}
	
	// Generate share token if public
	if req.IsPublic {
		wishlist.GenerateShareToken()
	}
	
	// Save wishlist
	if err := s.repo.Save(ctx, wishlist); err != nil {
		return nil, fmt.Errorf("failed to save wishlist: %w", err)
	}
	
	return s.buildWishlistResponse(wishlist), nil
}

// GetWishlist retrieves a wishlist by ID
func (s *ServiceImpl) GetWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID) (*WishlistResponse, error) {
	wishlist, err := s.repo.FindByID(ctx, tenantID, wishlistID)
	if err != nil {
		return nil, err
	}
	
	return s.buildWishlistResponse(wishlist), nil
}

// GetWishlistByShareToken retrieves a public wishlist by share token
func (s *ServiceImpl) GetWishlistByShareToken(ctx context.Context, shareToken string) (*WishlistResponse, error) {
	wishlist, err := s.repo.FindByShareToken(ctx, shareToken)
	if err != nil {
		return nil, err
	}
	
	return s.buildWishlistResponse(wishlist), nil
}

// UpdateWishlist updates an existing wishlist
func (s *ServiceImpl) UpdateWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID, req UpdateWishlistRequest) (*WishlistResponse, error) {
	// Validate request
	if err := s.validateUpdateWishlistRequest(req); err != nil {
		return nil, err
	}
	
	// Get existing wishlist
	wishlist, err := s.repo.FindByID(ctx, tenantID, wishlistID)
	if err != nil {
		return nil, err
	}
	
	// Check if name already exists for customer (excluding current wishlist)
	if req.Name != nil && *req.Name != wishlist.Name {
		exists, err := s.repo.ExistsByName(ctx, tenantID, wishlist.CustomerID, *req.Name, &wishlistID)
		if err != nil {
			return nil, fmt.Errorf("failed to check wishlist name: %w", err)
		}
		if exists {
			return nil, ErrWishlistNameExists
		}
		wishlist.Name = *req.Name
	}
	
	// Update fields
	if req.Description != nil {
		wishlist.Description = *req.Description
	}
	
	if req.IsDefault != nil && *req.IsDefault != wishlist.IsDefault {
		if *req.IsDefault {
			if err := s.repo.SetDefaultWishlist(ctx, tenantID, wishlist.CustomerID, wishlistID); err != nil {
				return nil, fmt.Errorf("failed to set default wishlist: %w", err)
			}
		}
		wishlist.IsDefault = *req.IsDefault
	}
	
	if req.IsPublic != nil {
		if *req.IsPublic && !wishlist.IsPublic {
			// Making public - generate share token
			wishlist.MakePublic()
		} else if !*req.IsPublic && wishlist.IsPublic {
			// Making private - clear share token
			wishlist.MakePrivate()
		}
	}
	
	wishlist.UpdatedAt = time.Now()
	
	// Save changes
	if err := s.repo.Update(ctx, wishlist); err != nil {
		return nil, fmt.Errorf("failed to update wishlist: %w", err)
	}
	
	return s.buildWishlistResponse(wishlist), nil
}

// DeleteWishlist deletes a wishlist
func (s *ServiceImpl) DeleteWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID) error {
	// Get wishlist to check if it can be deleted
	wishlist, err := s.repo.FindByID(ctx, tenantID, wishlistID)
	if err != nil {
		return err
	}
	
	// Check if wishlist can be deleted
	if !wishlist.CanDelete() {
		return ErrCannotDeleteDefaultWishlist
	}
	
	// Delete wishlist (items will be deleted by cascade)
	if err := s.repo.Delete(ctx, tenantID, wishlistID); err != nil {
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}
	
	return nil
}

// ListWishlists returns paginated wishlists
func (s *ServiceImpl) ListWishlists(ctx context.Context, tenantID uuid.UUID, filter WishlistFilter, limit, offset int) ([]WishlistResponse, int64, error) {
	// Validate pagination
	if limit <= 0 || limit > MaxPageSize {
		limit = DefaultPageSize
	}
	if offset < 0 {
		offset = 0
	}
	
	// Get wishlists
	wishlists, err := s.repo.List(ctx, tenantID, filter, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list wishlists: %w", err)
	}
	
	// Get total count
	total, err := s.repo.Count(ctx, tenantID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count wishlists: %w", err)
	}
	
	// Build responses
	responses := make([]WishlistResponse, len(wishlists))
	for i, wishlist := range wishlists {
		responses[i] = *s.buildWishlistResponse(&wishlist)
	}
	
	return responses, total, nil
}

// GetCustomerWishlists returns all wishlists for a customer
func (s *ServiceImpl) GetCustomerWishlists(ctx context.Context, tenantID, customerID uuid.UUID) ([]WishlistResponse, error) {
	wishlists, err := s.repo.FindByCustomerID(ctx, tenantID, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer wishlists: %w", err)
	}
	
	responses := make([]WishlistResponse, len(wishlists))
	for i, wishlist := range wishlists {
		responses[i] = *s.buildWishlistResponse(&wishlist)
	}
	
	return responses, nil
}

// GetDefaultWishlist returns the default wishlist for a customer
func (s *ServiceImpl) GetDefaultWishlist(ctx context.Context, tenantID, customerID uuid.UUID) (*WishlistResponse, error) {
	wishlist, err := s.repo.FindDefaultByCustomerID(ctx, tenantID, customerID)
	if err != nil {
		return nil, err
	}
	
	return s.buildWishlistResponse(wishlist), nil
}

// Wishlist item operations

// AddItem adds an item to a wishlist
func (s *ServiceImpl) AddItem(ctx context.Context, tenantID uuid.UUID, req AddItemRequest) (*WishlistItemResponse, error) {
	// Validate request
	if err := s.validateAddItemRequest(req); err != nil {
		return nil, err
	}
	
	// Get wishlist to validate
	wishlist, err := s.repo.FindByID(ctx, tenantID, req.WishlistID)
	if err != nil {
		return nil, err
	}
	
	// Check if wishlist can accept more items
	if !wishlist.CanAddItem() {
		return nil, ErrWishlistFull
	}
	
	// Check if item already exists
	existingItem, err := s.repo.FindItemByProductAndWishlist(ctx, tenantID, req.WishlistID, req.ProductID, req.VariantID)
	if err != nil && err != ErrWishlistItemNotFound {
		return nil, fmt.Errorf("failed to check existing item: %w", err)
	}
	
	if existingItem != nil {
		// Update existing item quantity
		existingItem.Quantity += req.Quantity
		if req.Notes != "" {
			existingItem.Notes = req.Notes
		}
		if req.Priority > 0 {
			existingItem.Priority = req.Priority
		}
		existingItem.UpdatedAt = time.Now()
		
		if err := s.repo.UpdateItem(ctx, existingItem); err != nil {
			return nil, fmt.Errorf("failed to update existing item: %w", err)
		}
		
		return s.buildWishlistItemResponse(existingItem), nil
	}
	
	// Create new item
	item := &WishlistItem{
		ID:         uuid.New(),
		TenantID:   tenantID,
		WishlistID: req.WishlistID,
		ProductID:  req.ProductID,
		VariantID:  req.VariantID,
		Quantity:   req.Quantity,
		Notes:      req.Notes,
		Priority:   req.Priority,
		AddedAt:    time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	// Save item
	if err := s.repo.SaveItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save wishlist item: %w", err)
	}
	
	return s.buildWishlistItemResponse(item), nil
}

// UpdateItem updates a wishlist item
func (s *ServiceImpl) UpdateItem(ctx context.Context, tenantID, itemID uuid.UUID, req UpdateItemRequest) (*WishlistItemResponse, error) {
	// Validate request
	if err := s.validateUpdateItemRequest(req); err != nil {
		return nil, err
	}
	
	// Get existing item
	item, err := s.repo.FindItemByID(ctx, tenantID, itemID)
	if err != nil {
		return nil, err
	}
	
	// Update fields
	if req.Quantity != nil {
		item.UpdateQuantity(*req.Quantity)
	}
	
	if req.Notes != nil {
		item.Notes = *req.Notes
	}
	
	if req.Priority != nil {
		item.Priority = *req.Priority
	}
	
	item.UpdatedAt = time.Now()
	
	// Save changes
	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update wishlist item: %w", err)
	}
	
	return s.buildWishlistItemResponse(item), nil
}

// RemoveItem removes an item from a wishlist
func (s *ServiceImpl) RemoveItem(ctx context.Context, tenantID, itemID uuid.UUID) error {
	// Check if item exists
	_, err := s.repo.FindItemByID(ctx, tenantID, itemID)
	if err != nil {
		return err
	}
	
	// Delete item
	if err := s.repo.DeleteItem(ctx, tenantID, itemID); err != nil {
		return fmt.Errorf("failed to delete wishlist item: %w", err)
	}
	
	return nil
}

// ListItems returns paginated wishlist items
func (s *ServiceImpl) ListItems(ctx context.Context, tenantID uuid.UUID, filter WishlistItemFilter, limit, offset int) ([]WishlistItemResponse, int64, error) {
	// Validate pagination
	if limit <= 0 || limit > MaxPageSize {
		limit = DefaultPageSize
	}
	if offset < 0 {
		offset = 0
	}
	
	// Get items
	items, err := s.repo.ListItems(ctx, tenantID, filter, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list wishlist items: %w", err)
	}
	
	// Get total count
	total, err := s.repo.CountItems(ctx, tenantID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count wishlist items: %w", err)
	}
	
	// Build responses
	responses := make([]WishlistItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.buildWishlistItemResponse(&item)
	}
	
	return responses, total, nil
}

// GetItem retrieves a wishlist item by ID
func (s *ServiceImpl) GetItem(ctx context.Context, tenantID, itemID uuid.UUID) (*WishlistItemResponse, error) {
	item, err := s.repo.FindItemByID(ctx, tenantID, itemID)
	if err != nil {
		return nil, err
	}
	
	return s.buildWishlistItemResponse(item), nil
}

// Wishlist management operations

// SetDefaultWishlist sets a wishlist as default for a customer
func (s *ServiceImpl) SetDefaultWishlist(ctx context.Context, tenantID, customerID, wishlistID uuid.UUID) error {
	// Verify wishlist exists and belongs to customer
	wishlist, err := s.repo.FindByID(ctx, tenantID, wishlistID)
	if err != nil {
		return err
	}
	
	if wishlist.CustomerID != customerID {
		return ErrWishlistNotFound
	}
	
	// Set as default
	if err := s.repo.SetDefaultWishlist(ctx, tenantID, customerID, wishlistID); err != nil {
		return fmt.Errorf("failed to set default wishlist: %w", err)
	}
	
	return nil
}

// ClearWishlist removes all items from a wishlist
func (s *ServiceImpl) ClearWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID) error {
	// Verify wishlist exists
	_, err := s.repo.FindByID(ctx, tenantID, wishlistID)
	if err != nil {
		return err
	}
	
	// Clear items
	if err := s.repo.ClearItems(ctx, tenantID, wishlistID); err != nil {
		return fmt.Errorf("failed to clear wishlist items: %w", err)
	}
	
	return nil
}

// MoveItem moves an item to another wishlist
func (s *ServiceImpl) MoveItem(ctx context.Context, tenantID, itemID, targetWishlistID uuid.UUID) error {
	// Verify item exists
	item, err := s.repo.FindItemByID(ctx, tenantID, itemID)
	if err != nil {
		return err
	}
	
	// Verify target wishlist exists
	targetWishlist, err := s.repo.FindByID(ctx, tenantID, targetWishlistID)
	if err != nil {
		return err
	}
	
	// Check if target wishlist can accept more items
	if !targetWishlist.CanAddItem() {
		return ErrWishlistFull
	}
	
	// Check if item already exists in target wishlist
	existingItem, err := s.repo.FindItemByProductAndWishlist(ctx, tenantID, targetWishlistID, item.ProductID, item.VariantID)
	if err != nil && err != ErrWishlistItemNotFound {
		return fmt.Errorf("failed to check existing item: %w", err)
	}
	
	if existingItem != nil {
		// Merge quantities and delete original
		existingItem.Quantity += item.Quantity
		existingItem.UpdatedAt = time.Now()
		
		if err := s.repo.UpdateItem(ctx, existingItem); err != nil {
			return fmt.Errorf("failed to update existing item: %w", err)
		}
		
		if err := s.repo.DeleteItem(ctx, tenantID, itemID); err != nil {
			return fmt.Errorf("failed to delete original item: %w", err)
		}
		
		return nil
	}
	
	// Move item
	if err := s.repo.MoveItem(ctx, tenantID, itemID, targetWishlistID); err != nil {
		return fmt.Errorf("failed to move item: %w", err)
	}
	
	return nil
}

// CopyItem copies an item to another wishlist
func (s *ServiceImpl) CopyItem(ctx context.Context, tenantID, itemID, targetWishlistID uuid.UUID) error {
	// Verify item exists
	item, err := s.repo.FindItemByID(ctx, tenantID, itemID)
	if err != nil {
		return err
	}
	
	// Verify target wishlist exists
	targetWishlist, err := s.repo.FindByID(ctx, tenantID, targetWishlistID)
	if err != nil {
		return err
	}
	
	// Check if target wishlist can accept more items
	if !targetWishlist.CanAddItem() {
		return ErrWishlistFull
	}
	
	// Check if item already exists in target wishlist
	existingItem, err := s.repo.FindItemByProductAndWishlist(ctx, tenantID, targetWishlistID, item.ProductID, item.VariantID)
	if err != nil && err != ErrWishlistItemNotFound {
		return fmt.Errorf("failed to check existing item: %w", err)
	}
	
	if existingItem != nil {
		// Update existing item quantity
		existingItem.Quantity += item.Quantity
		existingItem.UpdatedAt = time.Now()
		
		if err := s.repo.UpdateItem(ctx, existingItem); err != nil {
			return fmt.Errorf("failed to update existing item: %w", err)
		}
		
		return nil
	}
	
	// Copy item
	if err := s.repo.CopyItem(ctx, tenantID, itemID, targetWishlistID); err != nil {
		return fmt.Errorf("failed to copy item: %w", err)
	}
	
	return nil
}

// MergeWishlists merges source wishlist into target wishlist
func (s *ServiceImpl) MergeWishlists(ctx context.Context, tenantID, sourceWishlistID, targetWishlistID uuid.UUID) error {
	// Verify both wishlists exist
	sourceWishlist, err := s.repo.FindByID(ctx, tenantID, sourceWishlistID)
	if err != nil {
		return err
	}
	
	targetWishlist, err := s.repo.FindByID(ctx, tenantID, targetWishlistID)
	if err != nil {
		return err
	}
	
	// Check if source can be deleted
	if !sourceWishlist.CanDelete() {
		return ErrCannotDeleteDefaultWishlist
	}
	
	// Check if wishlists belong to same customer
	if sourceWishlist.CustomerID != targetWishlist.CustomerID {
		return ErrWishlistNotFound
	}
	
	// Merge wishlists
	if err := s.repo.MergeWishlists(ctx, tenantID, sourceWishlistID, targetWishlistID); err != nil {
		return fmt.Errorf("failed to merge wishlists: %w", err)
	}
	
	return nil
}

// ShareWishlist makes a wishlist public/private and returns share URL
func (s *ServiceImpl) ShareWishlist(ctx context.Context, tenantID, wishlistID uuid.UUID, isPublic bool) (string, error) {
	// Get wishlist
	wishlist, err := s.repo.FindByID(ctx, tenantID, wishlistID)
	if err != nil {
		return "", err
	}
	
	if isPublic {
		// Make public and generate share token
		wishlist.MakePublic()
	} else {
		// Make private and clear share token
		wishlist.MakePrivate()
	}
	
	// Save changes
	if err := s.repo.Update(ctx, wishlist); err != nil {
		return "", fmt.Errorf("failed to update wishlist: %w", err)
	}
	
	return wishlist.GetShareURL("https://example.com"), nil
}

// Bulk operations

// BulkAddItems adds multiple items to wishlists
func (s *ServiceImpl) BulkAddItems(ctx context.Context, tenantID uuid.UUID, requests []AddItemRequest) ([]WishlistItemResponse, error) {
	// Validate all requests
	for _, req := range requests {
		if err := s.validateAddItemRequest(req); err != nil {
			return nil, err
		}
	}
	
	// Create items
	items := make([]WishlistItem, 0, len(requests))
	for _, req := range requests {
		item := WishlistItem{
			ID:         uuid.New(),
			TenantID:   tenantID,
			WishlistID: req.WishlistID,
			ProductID:  req.ProductID,
			VariantID:  req.VariantID,
			Quantity:   req.Quantity,
			Notes:      req.Notes,
			Priority:   req.Priority,
			AddedAt:    time.Now(),
			UpdatedAt:  time.Now(),
		}
		items = append(items, item)
	}
	
	// Save all items
	if err := s.repo.BulkAddItems(ctx, items); err != nil {
		return nil, fmt.Errorf("failed to bulk add items: %w", err)
	}
	
	// Build responses
	responses := make([]WishlistItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.buildWishlistItemResponse(&item)
	}
	
	return responses, nil
}

// BulkRemoveItems removes multiple items from wishlists
func (s *ServiceImpl) BulkRemoveItems(ctx context.Context, tenantID uuid.UUID, itemIDs []uuid.UUID) error {
	if len(itemIDs) == 0 {
		return nil
	}
	
	// Delete items
	if err := s.repo.BulkDeleteItems(ctx, tenantID, itemIDs); err != nil {
		return fmt.Errorf("failed to bulk delete items: %w", err)
	}
	
	return nil
}

// BulkUpdateItemPriority updates priority for multiple items
func (s *ServiceImpl) BulkUpdateItemPriority(ctx context.Context, tenantID uuid.UUID, updates map[uuid.UUID]int) error {
	if len(updates) == 0 {
		return nil
	}
	
	// Update priorities
	if err := s.repo.BulkUpdateItemPriority(ctx, tenantID, updates); err != nil {
		return fmt.Errorf("failed to bulk update item priorities: %w", err)
	}
	
	return nil
}

// ReorderItems reorders items in a wishlist
func (s *ServiceImpl) ReorderItems(ctx context.Context, tenantID, wishlistID uuid.UUID, itemOrder []uuid.UUID) error {
	// Verify wishlist exists
	_, err := s.repo.FindByID(ctx, tenantID, wishlistID)
	if err != nil {
		return err
	}
	
	// Create priority updates based on order
	updates := make(map[uuid.UUID]int)
	for i, itemID := range itemOrder {
		updates[itemID] = len(itemOrder) - i // Higher index = higher priority
	}
	
	// Update priorities
	if err := s.repo.BulkUpdateItemPriority(ctx, tenantID, updates); err != nil {
		return fmt.Errorf("failed to reorder items: %w", err)
	}
	
	return nil
}

// Analytics and statistics

// GetWishlistStats returns wishlist statistics
func (s *ServiceImpl) GetWishlistStats(ctx context.Context, tenantID uuid.UUID) (*WishlistStats, error) {
	stats, err := s.repo.GetStats(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist stats: %w", err)
	}
	
	return stats, nil
}

// GetMostWishedProducts returns the most wished products
func (s *ServiceImpl) GetMostWishedProducts(ctx context.Context, tenantID uuid.UUID, limit int) ([]ProductWishCount, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	products, err := s.repo.GetMostWishedProducts(ctx, tenantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most wished products: %w", err)
	}
	
	return products, nil
}

// GetCustomerActivity returns customer wishlist activity
func (s *ServiceImpl) GetCustomerActivity(ctx context.Context, tenantID, customerID uuid.UUID, days int) ([]WishlistActivity, error) {
	if days <= 0 || days > 365 {
		days = 30
	}
	
	activity, err := s.repo.GetCustomerWishlistActivity(ctx, tenantID, customerID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer activity: %w", err)
	}
	
	return activity, nil
}

// GetPopularWishlists returns popular public wishlists
func (s *ServiceImpl) GetPopularWishlists(ctx context.Context, tenantID uuid.UUID, limit int) ([]WishlistResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	wishlists, err := s.repo.GetPopularWishlists(ctx, tenantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular wishlists: %w", err)
	}
	
	responses := make([]WishlistResponse, len(wishlists))
	for i, wishlist := range wishlists {
		responses[i] = *s.buildWishlistResponse(&wishlist)
	}
	
	return responses, nil
}

// Maintenance operations

// CleanupEmptyWishlists removes empty wishlists older than specified days
func (s *ServiceImpl) CleanupEmptyWishlists(ctx context.Context, tenantID uuid.UUID, olderThanDays int) (int64, error) {
	if olderThanDays <= 0 {
		olderThanDays = 30
	}
	
	count, err := s.repo.CleanupEmptyWishlists(ctx, tenantID, olderThanDays)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup empty wishlists: %w", err)
	}
	
	return count, nil
}

// CleanupOrphanedItems removes items for non-existent products
func (s *ServiceImpl) CleanupOrphanedItems(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	count, err := s.repo.CleanupOrphanedItems(ctx, tenantID)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup orphaned items: %w", err)
	}
	
	return count, nil
}

// Helper methods

// validateCreateWishlistRequest validates create wishlist request
func (s *ServiceImpl) validateCreateWishlistRequest(req CreateWishlistRequest) error {
	if req.TenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	
	if req.CustomerID == uuid.Nil {
		return ErrInvalidCustomerID
	}
	
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidWishlistName
	}
	
	if len(req.Name) > MaxWishlistNameLength {
		return ErrWishlistNameTooLong
	}
	
	if len(req.Description) > MaxWishlistDescriptionLength {
		return ErrWishlistDescriptionTooLong
	}
	
	return nil
}

// validateUpdateWishlistRequest validates update wishlist request
func (s *ServiceImpl) validateUpdateWishlistRequest(req UpdateWishlistRequest) error {
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return ErrInvalidWishlistName
		}
		
		if len(*req.Name) > MaxWishlistNameLength {
			return ErrWishlistNameTooLong
		}
	}
	
	if req.Description != nil && len(*req.Description) > MaxWishlistDescriptionLength {
		return ErrWishlistDescriptionTooLong
	}
	
	return nil
}

// validateAddItemRequest validates add item request
func (s *ServiceImpl) validateAddItemRequest(req AddItemRequest) error {
	if req.WishlistID == uuid.Nil {
		return ErrInvalidWishlistID
	}
	
	if req.ProductID == uuid.Nil {
		return ErrInvalidProductID
	}
	
	if req.Quantity <= 0 || req.Quantity > MaxItemQuantity {
		return ErrInvalidQuantity
	}
	
	if len(req.Notes) > MaxItemNotesLength {
		return ErrItemNotesTooLong
	}
	
	if req.Priority < 0 || req.Priority > MaxItemPriority {
		return ErrInvalidPriority
	}
	
	return nil
}

// validateUpdateItemRequest validates update item request
func (s *ServiceImpl) validateUpdateItemRequest(req UpdateItemRequest) error {
	if req.Quantity != nil && (*req.Quantity <= 0 || *req.Quantity > MaxItemQuantity) {
		return ErrInvalidQuantity
	}
	
	if req.Notes != nil && len(*req.Notes) > MaxItemNotesLength {
		return ErrItemNotesTooLong
	}
	
	if req.Priority != nil && (*req.Priority < 0 || *req.Priority > MaxItemPriority) {
		return ErrInvalidPriority
	}
	
	return nil
}

// buildWishlistResponse builds a wishlist response
func (s *ServiceImpl) buildWishlistResponse(wishlist *Wishlist) *WishlistResponse {
	return &WishlistResponse{
		Wishlist: wishlist,
		ShareURL: wishlist.GetShareURL("https://example.com"),
	}
}

// buildWishlistItemResponse builds a wishlist item response
func (s *ServiceImpl) buildWishlistItemResponse(item *WishlistItem) *WishlistItemResponse {
	return &WishlistItemResponse{
		WishlistItem: item,
		DisplayName: item.GetDisplayName(),
		CurrentPrice: item.GetPrice(),
		ComparePrice: item.GetComparePrice(),
		DiscountPercentage: item.GetDiscountPercentage(),
		IsAvailable: item.IsAvailable(),
	}
}