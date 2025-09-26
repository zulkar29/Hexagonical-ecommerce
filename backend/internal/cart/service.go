package cart

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"ecommerce-saas/internal/product"
	"ecommerce-saas/internal/discount"
	"ecommerce-saas/internal/shipping"
	sharedErrors "ecommerce-saas/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Request/Response Structures
type CreateCartRequest struct {
	CustomerID *uuid.UUID `json:"customer_id,omitempty"`
	SessionID  string     `json:"session_id,omitempty"`
	Currency   string     `json:"currency" validate:"required,len=3"`
	Notes      string     `json:"notes,omitempty" validate:"max=500"`
}

type AddItemRequest struct {
	ProductID      uuid.UUID         `json:"product_id" validate:"required"`
	Quantity       int               `json:"quantity" validate:"required,min=1"`
	Customizations map[string]any    `json:"customizations,omitempty"`
	Notes          string            `json:"notes,omitempty"`
}

type UpdateItemRequest struct {
	Quantity       *int           `json:"quantity,omitempty" validate:"omitempty,min=1,max=100"`
	Customizations map[string]any `json:"customizations,omitempty"`
	Notes          string         `json:"notes,omitempty" validate:"max=200"`
}

type ApplyCouponRequest struct {
	CouponCode string `json:"coupon_code" validate:"required,min=1,max=50"`
}

type UpdateAddressRequest struct {
	ShippingAddress *Address `json:"shipping_address,omitempty"`
	BillingAddress  *Address `json:"billing_address,omitempty"`
}

type UpdateShippingRequest struct {
	ShippingMethodID *uuid.UUID `json:"shipping_method_id,omitempty"`
}

type UpdateCartRequest struct {
	ShippingAddress  *Address   `json:"shipping_address,omitempty"`
	BillingAddress   *Address   `json:"billing_address,omitempty"`
	ShippingMethodID *uuid.UUID `json:"shipping_method_id,omitempty"`
	CouponCode       *string    `json:"coupon_code,omitempty"`
	Notes            *string    `json:"notes,omitempty" validate:"omitempty,max=500"`
}

type EstimateRequest struct {
	ShippingAddress *Address   `json:"shipping_address" validate:"required"`
	ShippingMethodID *uuid.UUID `json:"shipping_method_id,omitempty"`
}

type EstimateResponse struct {
	ShippingMethods []ShippingMethod `json:"shipping_methods"`
	Subtotal        float64          `json:"subtotal"`
	Total           float64          `json:"total"`
}



type GuestCheckoutRequest struct {
	SessionID       string  `json:"session_id" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	ShippingAddress Address `json:"shipping_address" validate:"required"`
	BillingAddress  Address `json:"billing_address" validate:"required"`
	ShippingMethodID uuid.UUID `json:"shipping_method_id" validate:"required"`
	PaymentMethodID string  `json:"payment_method_id" validate:"required"`
}

type GuestCheckoutResponse struct {
	OrderID     uuid.UUID `json:"order_id"`
	OrderNumber string    `json:"order_number"`
	Total       float64   `json:"total"`
	Status      string    `json:"status"`
}

type OrderConversionResult struct {
	OrderID      uuid.UUID `json:"order_id"`
	OrderNumber  string    `json:"order_number"`
	Total        float64   `json:"total"`
	ItemCount    int       `json:"item_count"`
	Success      bool      `json:"success"`
	Error        string    `json:"error,omitempty"`
}

// Response structures
type CartResponse struct {
	*Cart
	ItemCount       int     `json:"item_count"`
	UniqueItemCount int     `json:"unique_item_count"`
	SavingsAmount   float64 `json:"savings_amount"`
}

type CartSummary struct {
	ID              uuid.UUID `json:"id"`
	ItemCount       int       `json:"item_count"`
	UniqueItemCount int       `json:"unique_item_count"`
	Subtotal        float64   `json:"subtotal"`
	Total           float64   `json:"total"`
	Currency        string    `json:"currency"`
	Status          CartStatus `json:"status"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Service interfaces
type ProductService interface {
	GetProduct(tenantID uuid.UUID, productID string) (*product.Product, error)
	CheckAvailability(tenantID, productID uuid.UUID, variantID *uuid.UUID, quantity int) (bool, error)
	ReserveStock(tenantID, productID uuid.UUID, quantity int) error
	RestoreStock(tenantID, productID uuid.UUID, quantity int) error
}

type DiscountService interface {
	ValidateDiscountCode(ctx context.Context, req discount.ValidateDiscountRequest) (*discount.DiscountValidation, error)
	ApplyDiscount(ctx context.Context, req discount.ApplyDiscountRequest) (*discount.DiscountApplication, error)
	GetDiscountByCode(ctx context.Context, tenantID uuid.UUID, code string) (*discount.Discount, error)
}



type ShippingService interface {
	CreateShippingZone(tenantID uuid.UUID, req shipping.CreateShippingZoneRequest) (*shipping.ShippingZone, error)
	GetShippingZones(tenantID uuid.UUID) ([]shipping.ShippingZone, error)
}

type OrderService interface {
	CreateOrderFromCart(ctx context.Context, tenantID, cartID uuid.UUID) (*OrderFromCartResult, error)
}

// OrderFromCartResult represents the result of creating an order from a cart
type OrderFromCartResult struct {
	OrderID      uuid.UUID `json:"order_id"`
	OrderNumber  string    `json:"order_number"`
	Total        float64   `json:"total"`
	ItemCount    int       `json:"item_count"`
}

// External service data structures


type CouponInfo struct {
	Code         string    `json:"code"`
	DiscountType string    `json:"discount_type"` // percentage, fixed
	Value        float64   `json:"value"`
	MinAmount    float64   `json:"min_amount"`
	MaxDiscount  float64   `json:"max_discount"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type ShippingMethod struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost"`
	EstimatedDays int     `json:"estimated_days"`
}

// ServiceInterface defines the cart service interface
type ServiceInterface interface {
	CreateCart(tenantID uuid.UUID, req CreateCartRequest) (*CartResponse, error)
	GetCartByID(tenantID, cartID uuid.UUID) (*CartResponse, error)
	GetCartByCustomer(tenantID, customerID uuid.UUID) (*CartResponse, error)
	GetCartBySession(tenantID uuid.UUID, sessionID string) (*CartResponse, error)
	AddItem(tenantID, cartID uuid.UUID, req AddItemRequest) (*CartResponse, error)
	UpdateItem(tenantID, cartID, itemID uuid.UUID, req UpdateItemRequest) (*CartItem, error)
	RemoveItem(tenantID, cartID, itemID uuid.UUID) error
	UpdateCart(tenantID, cartID uuid.UUID, req UpdateCartRequest) (*CartResponse, error)
	ApplyCoupon(tenantID, cartID uuid.UUID, req ApplyCouponRequest) (*CartResponse, error)
	RemoveCoupon(tenantID, cartID uuid.UUID) (*CartResponse, error)
	GetCartSummary(tenantID, cartID uuid.UUID) (*CartSummary, error)
	GetEstimates(tenantID, cartID uuid.UUID, req EstimateRequest) (*EstimateResponse, error)
	ClearCart(tenantID, cartID uuid.UUID) error
	DeleteCart(tenantID, cartID uuid.UUID) error
	ConvertToOrder(tenantID, cartID uuid.UUID) (*OrderConversionResult, error)
	ProcessGuestCheckout(tenantID uuid.UUID, req GuestCheckoutRequest) (*GuestCheckoutResponse, error)
	ListCarts(tenantID uuid.UUID, filter CartListFilter, offset, limit int) ([]*CartResponse, int64, error)
	GetCartStats(tenantID uuid.UUID) (*CartStats, error)
}

// Service type alias for interface compatibility
type Service = ServiceInterface

// CartService handles cart business logic
type CartService struct {
	repo            Repository
	validator       *validator.Validate
	productService  ProductService
	discountService DiscountService
	shippingService ShippingService
	orderService    OrderService
	cartExpiration  time.Duration
}

// NewCartService creates a new cart service implementation
func NewCartService(repo Repository, productService ProductService, discountService DiscountService, shippingService ShippingService, orderService OrderService) *CartService {
	return &CartService{
		repo:            repo,
		validator:       validator.New(),
		productService:  productService,
		discountService: discountService,
		shippingService: shippingService,
		orderService:    orderService,
		cartExpiration:  24 * time.Hour * 30, // 30 days default
	}
}

// NewService creates a new cart service (interface compatibility)
func NewService(repo Repository, productService ProductService, discountService DiscountService, shippingService ShippingService, orderService OrderService) Service {
	return NewCartService(repo, productService, discountService, shippingService, orderService)
}

// CreateCart creates a new cart
func (s *CartService) CreateCart(tenantID uuid.UUID, req CreateCartRequest) (*CartResponse, error) {
	// Validate request
	if validateErr := s.validator.Struct(req); validateErr != nil {
		return nil, sharedErrors.NewValidationError(validateErr.Error(), nil)
	}

	// Validate that either customer ID or session ID is provided
	if req.CustomerID == nil && req.SessionID == "" {
		return nil, sharedErrors.NewValidationError("either customer_id or session_id must be provided", nil)
	}

	// Check if cart already exists
	if req.CustomerID != nil {
		if existingCart, err := s.repo.FindCartByCustomerID(tenantID, *req.CustomerID); err == nil {
			return s.buildCartResponse(existingCart), nil
		}
	} else if req.SessionID != "" {
		if existingCart, err := s.repo.FindCartBySessionID(tenantID, req.SessionID); err == nil {
			return s.buildCartResponse(existingCart), nil
		}
	}

	// Create new cart
	cart := &Cart{
		ID:         uuid.New(),
		TenantID:   tenantID,
		CustomerID: req.CustomerID,
		SessionID:  req.SessionID,
		Status:     StatusActive,
		Currency:   strings.ToUpper(req.Currency),
		Notes:      strings.TrimSpace(req.Notes),
		Items:      []CartItem{},
	}

	// Set expiration
	cart.SetExpiration(s.cartExpiration)

	savedCart, err := s.repo.SaveCart(cart)
	if err != nil {
		return nil, fmt.Errorf("failed to save cart: %w", err)
	}

	return s.buildCartResponse(savedCart), nil
}

// GetCart retrieves a cart by ID
func (s *CartService) GetCart(tenantID, cartID uuid.UUID) (*CartResponse, error) {
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Check if cart is expired and update status
	if cart.IsExpired() && cart.Status == StatusActive {
		cart.Status = StatusExpired
		if _, err := s.repo.UpdateCart(cart); err != nil {
			log.Printf("Failed to update cart status: %v", err)
		}
		return nil, ErrCartExpired
	}

	return s.buildCartResponse(cart), nil
}

// GetCartByID retrieves a cart by ID (alias for GetCart)
func (s *CartService) GetCartByID(tenantID, cartID uuid.UUID) (*CartResponse, error) {
	return s.GetCart(tenantID, cartID)
}

// GetCartByCustomer retrieves active cart for a customer
func (s *CartService) GetCartByCustomer(tenantID, customerID uuid.UUID) (*CartResponse, error) {
	cart, err := s.repo.FindCartByCustomerID(tenantID, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Check if cart is expired
	if cart.IsExpired() && cart.Status == StatusActive {
		cart.Status = StatusExpired
		if _, err := s.repo.UpdateCart(cart); err != nil {
			log.Printf("Failed to update cart status: %v", err)
		}
		return nil, ErrCartExpired
	}

	return s.buildCartResponse(cart), nil
}

// GetCartBySession retrieves cart for a guest session
func (s *CartService) GetCartBySession(tenantID uuid.UUID, sessionID string) (*CartResponse, error) {
	cart, err := s.repo.FindCartBySessionID(tenantID, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Check if cart is expired
	if cart.IsExpired() && cart.Status == StatusActive {
		cart.Status = StatusExpired
		if _, err := s.repo.UpdateCart(cart); err != nil {
			log.Printf("Failed to update cart status to expired: %v", err)
		}
		return nil, ErrCartExpired
	}

	return s.buildCartResponse(cart), nil
}

// AddItem adds an item to the cart
func (s *CartService) AddItem(tenantID, cartID uuid.UUID, req AddItemRequest) (*CartResponse, error) {
	// Validate request
	if validateErr := s.validator.Struct(req); validateErr != nil {
		return nil, sharedErrors.NewValidationError(validateErr.Error(), nil)
	}

	// Get cart
	cart, cartErr := s.repo.FindCartByID(tenantID, cartID)
	if cartErr != nil {
		if errors.Is(cartErr, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", cartErr)
	}

	// Check if cart can be modified
	if modifyErr := cart.CanModify(); modifyErr != nil {
		return nil, modifyErr
	}

	// Get product information
	product, err := s.productService.GetProduct(tenantID, req.ProductID.String())
	if err != nil {
		return nil, ErrProductNotFound
	}

	// Check if product is available
	if product.Status != "active" {
		return nil, errors.New("product is not available")
	}

	// Check inventory
	available, err := s.productService.CheckAvailability(tenantID, req.ProductID, nil, req.Quantity)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, ErrInsufficientStock
	}

	// Check if item already exists in cart
	existingItem := cart.FindItem(req.ProductID)
	if existingItem != nil {
		// Update existing item quantity
		newQuantity := existingItem.Quantity + req.Quantity
		
		// Check inventory for new total quantity
		available, availErr := s.productService.CheckAvailability(tenantID, req.ProductID, nil, newQuantity)
		if availErr != nil {
			return nil, availErr
		}
		if !available {
			return nil, ErrInsufficientStock
		}
		
		existingItem.Quantity = newQuantity
		existingItem.Customizations = req.Customizations
		existingItem.Notes = strings.TrimSpace(req.Notes)
		existingItem.CalculateLineTotal()
		
		if _, updateErr := s.repo.UpdateCartItem(existingItem); updateErr != nil {
			return nil, updateErr
		}
	} else {
		// Create new cart item
		price := product.Price
		comparePrice := product.ComparePrice
		sku := product.SKU
		image := product.FeaturedImage
		
		item := &CartItem{
			ID:             uuid.New(),
			CartID:         cartID,
			ProductID:      req.ProductID,
			ProductName:    product.Name,
			ProductSlug:    product.Slug,
			SKU:            sku,
			Price:          price,
			ComparePrice:   comparePrice,
			Image:          image,
			Quantity:       req.Quantity,
			Customizations: req.Customizations,
			Notes:          strings.TrimSpace(req.Notes),
		}
		
		item.CalculateLineTotal()
		
		if _, addErr := s.repo.AddCartItem(item); addErr != nil {
			return nil, addErr
		}
		
		cart.Items = append(cart.Items, *item)
	}

	// Recalculate cart totals
	if recalcErr := s.recalculateCart(cart); recalcErr != nil {
		return nil, recalcErr
	}

	// Update cart
	updatedCart, err := s.repo.UpdateCart(cart)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	return s.buildCartResponse(updatedCart), nil
}

// UpdateItem updates a cart item
func (s *CartService) UpdateItem(tenantID, cartID, itemID uuid.UUID, req UpdateItemRequest) (*CartItem, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, sharedErrors.NewValidationError(err.Error(), nil)
	}

	// Get cart
	cart, cartErr := s.repo.FindCartByID(tenantID, cartID)
	if cartErr != nil {
		if errors.Is(cartErr, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", cartErr)
	}

	// Check if cart can be modified
	if modifyErr := cart.CanModify(); modifyErr != nil {
		return nil, modifyErr
	}

	// Find cart item
	item, itemErr := s.repo.FindCartItem(tenantID, cartID, itemID)
	if itemErr != nil {
		return nil, ErrItemNotFound
	}

	// Update quantity if provided
	if req.Quantity != nil {
		// Check inventory
		available, availErr := s.productService.CheckAvailability(tenantID, item.ProductID, nil, *req.Quantity)
		if availErr != nil {
			return nil, availErr
		}
		if !available {
			return nil, ErrInsufficientStock
		}
		
		item.Quantity = *req.Quantity
		item.CalculateLineTotal()
	}

	// Update customizations and notes
	if req.Customizations != nil {
		item.Customizations = req.Customizations
	}
	item.Notes = strings.TrimSpace(req.Notes)

	// Update item
	updatedItem, updateErr := s.repo.UpdateCartItem(item)
	if updateErr != nil {
		return nil, updateErr
	}

	// Reload cart with updated items
	cart, reloadErr := s.repo.FindCartByID(tenantID, cartID)
	if reloadErr != nil {
		return nil, reloadErr
	}

	// Recalculate cart totals
	if recalcErr := s.recalculateCart(cart); recalcErr != nil {
		return nil, recalcErr
	}

	// Update cart
	_, finalUpdateErr := s.repo.UpdateCart(cart)
	if finalUpdateErr != nil {
		return nil, finalUpdateErr
	}

	return updatedItem, nil
}

// RemoveItem removes an item from the cart
func (s *CartService) RemoveItem(tenantID, cartID, itemID uuid.UUID) error {
	// Get cart
	cart, cartErr := s.repo.FindCartByID(tenantID, cartID)
	if cartErr != nil {
		if errors.Is(cartErr, gorm.ErrRecordNotFound) {
			return sharedErrors.NewNotFoundError("cart not found")
		}
		return fmt.Errorf("failed to get cart: %w", cartErr)
	}

	// Check if cart can be modified
	if modifyErr := cart.CanModify(); modifyErr != nil {
		return sharedErrors.NewBadRequestError(modifyErr.Error())
	}

	// Remove item
	if removeErr := s.repo.RemoveCartItem(tenantID, cartID, itemID); removeErr != nil {
		return sharedErrors.NewNotFoundError("item not found in cart")
	}

	// Reload cart
	cart, reloadErr := s.repo.FindCartByID(tenantID, cartID)
	if reloadErr != nil {
		return reloadErr
	}

	// Recalculate cart totals
	if recalcErr := s.recalculateCart(cart); recalcErr != nil {
		return recalcErr
	}

	// Update cart
	_, updateErr := s.repo.UpdateCart(cart)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

// ClearCart removes all items from the cart
func (s *CartService) ClearCart(tenantID, cartID uuid.UUID) error {
	// Get cart
	cart, cartErr := s.repo.FindCartByID(tenantID, cartID)
	if cartErr != nil {
		if errors.Is(cartErr, gorm.ErrRecordNotFound) {
			return sharedErrors.NewNotFoundError("cart not found")
		}
		return fmt.Errorf("failed to get cart: %w", cartErr)
	}

	// Check if cart can be modified
	if modifyErr := cart.CanModify(); modifyErr != nil {
		return sharedErrors.NewBadRequestError(modifyErr.Error())
	}

	// Clear all items
	if clearErr := s.repo.ClearCartItems(tenantID, cartID); clearErr != nil {
		return clearErr
	}

	// Reset cart totals
	cart.Items = []CartItem{}
	cart.Subtotal = 0
	cart.ShippingCost = 0
	cart.DiscountAmount = 0
	cart.Total = 0
	cart.CouponCode = ""
	cart.DiscountID = nil

	// Update cart
	_, updateErr := s.repo.UpdateCart(cart)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

// ApplyCoupon applies a coupon to the cart
func (s *CartService) ApplyCoupon(tenantID, cartID uuid.UUID, req ApplyCouponRequest) (*CartResponse, error) {
	// Validate request
	if validateErr := s.validator.Struct(req); validateErr != nil {
		return nil, sharedErrors.NewValidationError(validateErr.Error(), nil)
	}

	// Get cart
	cart, cartErr := s.repo.FindCartByID(tenantID, cartID)
	if cartErr != nil {
		if errors.Is(cartErr, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", cartErr)
	}

	// Check if cart can be modified
	if modifyErr := cart.CanModify(); modifyErr != nil {
		return nil, sharedErrors.NewBadRequestError(modifyErr.Error())
	}

	// Apply coupon (simplified - actual discount calculation would be done during order processing)
	cart.CouponCode = req.CouponCode
	cart.DiscountAmount = 0 // Will be calculated during order processing

	// Recalculate totals
	cart.UpdateTotals()

	// Update cart
	updatedCart, updateErr := s.repo.UpdateCart(cart)
	if updateErr != nil {
		return nil, updateErr
	}

	return s.buildCartResponse(updatedCart), nil
}

// RemoveCoupon removes coupon from the cart
func (s *CartService) RemoveCoupon(tenantID, cartID uuid.UUID) (*CartResponse, error) {
	// Get cart
	cart, cartErr := s.repo.FindCartByID(tenantID, cartID)
	if cartErr != nil {
		if errors.Is(cartErr, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", cartErr)
	}

	// Check if cart can be modified
	if modifyErr := cart.CanModify(); modifyErr != nil {
		return nil, sharedErrors.NewBadRequestError(modifyErr.Error())
	}

	// Remove coupon
	cart.CouponCode = ""
	cart.DiscountAmount = 0
	cart.DiscountID = nil

	// Recalculate totals
	cart.UpdateTotals()

	// Update cart
	updatedCart, updateErr := s.repo.UpdateCart(cart)
	if updateErr != nil {
		return nil, updateErr
	}

	return s.buildCartResponse(updatedCart), nil
}

// UpdateAddress updates shipping/billing address
func (s *CartService) UpdateAddress(tenantID, cartID uuid.UUID, req UpdateAddressRequest) (*CartResponse, error) {
	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Check if cart can be modified
	if modifyErr := cart.CanModify(); modifyErr != nil {
		return nil, sharedErrors.NewBadRequestError(modifyErr.Error())
	}

	// Update addresses
	if req.ShippingAddress != nil {
		cart.ShippingAddress = req.ShippingAddress
	}
	if req.BillingAddress != nil {
		cart.BillingAddress = req.BillingAddress
	}

	// Recalculate shipping if address changed
	if recalcErr := s.recalculateCart(cart); recalcErr != nil {
		return nil, recalcErr
	}

	// Update cart
	updatedCart, updateErr := s.repo.UpdateCart(cart)
	if updateErr != nil {
		return nil, updateErr
	}

	return s.buildCartResponse(updatedCart), nil
}

// UpdateShipping updates shipping method
func (s *CartService) UpdateShipping(tenantID, cartID uuid.UUID, req UpdateShippingRequest) (*CartResponse, error) {
	// Validate request
	if validateErr := s.validator.Struct(req); validateErr != nil {
		return nil, sharedErrors.NewValidationError(validateErr.Error(), nil)
	}

	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Check if cart can be modified
	if modifyErr := cart.CanModify(); modifyErr != nil {
		return nil, sharedErrors.NewBadRequestError(modifyErr.Error())
	}

	// Update shipping method
	cart.ShippingMethodID = req.ShippingMethodID

	// Shipping cost will be calculated during order processing
	cart.ShippingCost = 0

	// Recalculate totals
	cart.UpdateTotals()

	// Update cart
	updatedCart, updateErr := s.repo.UpdateCart(cart)
	if updateErr != nil {
		return nil, updateErr
	}

	return s.buildCartResponse(updatedCart), nil
}

// MergeGuestCart merges guest cart to customer cart
func (s *CartService) MergeGuestCart(tenantID uuid.UUID, sessionID string, customerID uuid.UUID) (*CartResponse, error) {
	if mergeErr := s.repo.MergeGuestCartToCustomer(tenantID, sessionID, customerID); mergeErr != nil {
		return nil, mergeErr
	}

	// Get the merged cart
	cart, err := s.repo.FindCartByCustomerID(tenantID, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Recalculate totals
	if recalcErr := s.recalculateCart(cart); recalcErr != nil {
		return nil, recalcErr
	}

	// Update cart
	updatedCart, updateErr := s.repo.UpdateCart(cart)
	if updateErr != nil {
		return nil, updateErr
	}

	return s.buildCartResponse(updatedCart), nil
}

// AbandonCart marks cart as abandoned
func (s *CartService) AbandonCart(tenantID, cartID uuid.UUID) error {
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedErrors.NewNotFoundError("cart not found")
		}
		return fmt.Errorf("failed to get cart: %w", err)
	}

	if cart.Status == StatusActive {
		cart.MarkAsAbandoned()
		_, updateErr := s.repo.UpdateCart(cart)
		if updateErr != nil {
			return fmt.Errorf("failed to update cart: %w", updateErr)
		}
	}

	return nil
}

// ConvertCart marks cart as converted to order
func (s *CartService) ConvertCart(tenantID, cartID uuid.UUID) error {
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedErrors.NewNotFoundError("cart not found")
		}
		return fmt.Errorf("failed to get cart: %w", err)
	}

	if cart.Status == StatusActive {
		cart.MarkAsConverted()
		_, updateErr := s.repo.UpdateCart(cart)
		if updateErr != nil {
			return fmt.Errorf("failed to update cart: %w", updateErr)
		}
	}

	return nil
}

// DeleteCart soft deletes a cart
func (s *CartService) DeleteCart(tenantID, cartID uuid.UUID) error {
	return s.repo.DeleteCart(tenantID, cartID)
}

// ConvertToOrder converts a cart to an order
func (s *CartService) ConvertToOrder(tenantID, cartID uuid.UUID) (*OrderConversionResult, error) {
	// Get the cart first
	cart, err := s.repo.GetCartByID(tenantID, cartID)
	if err != nil {
		return &OrderConversionResult{
			Success: false,
			Error:   fmt.Sprintf("Cart not found: %v", err),
		}, nil
	}

	// Validate cart has items
	if len(cart.Items) == 0 {
		return &OrderConversionResult{
			Success: false,
			Error:   "Cart is empty",
		}, nil
	}

	// Validate cart has required information for order creation
	if cart.ShippingAddress.Address1 == "" {
		return &OrderConversionResult{
			Success: false,
			Error:   "Shipping address is required",
		}, nil
	}

	// Check inventory availability for all items
	for _, item := range cart.Items {
		available, err := s.productService.CheckAvailability(tenantID, item.ProductID, item.VariantID, item.Quantity)
		if err != nil {
			return &OrderConversionResult{
				Success: false,
				Error:   fmt.Sprintf("Failed to check availability for product %s: %v", item.ProductID, err),
			}, nil
		}
		if !available {
			return &OrderConversionResult{
				Success: false,
				Error:   fmt.Sprintf("Product %s is not available in requested quantity", item.ProductName),
			}, nil
		}
	}

	// Convert cart to order using order service
	ctx := context.Background()
	orderResult, err := s.orderService.CreateOrderFromCart(ctx, tenantID, cartID)
	if err != nil {
		return &OrderConversionResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to create order: %v", err),
		}, nil
	}

	// Clear the cart after successful order creation
	if err := s.ClearCart(tenantID, cartID); err != nil {
		log.Printf("Warning: Failed to clear cart %s after order creation: %v", cartID, err)
	}

	return &OrderConversionResult{
		OrderID:     orderResult.OrderID,
		OrderNumber: orderResult.OrderNumber,
		Total:       orderResult.Total,
		ItemCount:   orderResult.ItemCount,
		Success:     true,
	}, nil
}

// GetCartSummary returns a summary of the cart
func (s *CartService) GetCartSummary(tenantID, cartID uuid.UUID) (*CartSummary, error) {
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &CartSummary{
		ID:              cart.ID,
		ItemCount:       cart.GetItemCount(),
		UniqueItemCount: cart.GetUniqueItemCount(),
		Subtotal:        cart.Subtotal,
		Total:           cart.Total,
		Currency:        cart.Currency,
		Status:          cart.Status,
		UpdatedAt:       cart.UpdatedAt,
	}, nil
}

// ListCarts returns paginated carts
func (s *CartService) ListCarts(tenantID uuid.UUID, filter CartListFilter, offset, limit int) ([]*CartResponse, int64, error) {
	carts, total, err := s.repo.ListCarts(tenantID, filter, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*CartResponse, len(carts))
	for i, cart := range carts {
		responses[i] = s.buildCartResponse(cart)
	}

	return responses, total, nil
}

// GetCartStats returns cart statistics
func (s *CartService) GetCartStats(tenantID uuid.UUID) (*CartStats, error) {
	return s.repo.GetCartStats(tenantID)
}

// CleanupExpiredCarts marks expired carts as expired
func (s *CartService) CleanupExpiredCarts(tenantID uuid.UUID) error {
	return s.repo.CleanupExpiredCarts(tenantID)
}

// Helper methods

// recalculateCart recalculates all cart totals
func (s *CartService) recalculateCart(cart *Cart) error {
	// Calculate subtotal
	cart.UpdateTotals()

	// Shipping and discount calculations will be done during order processing
	// Cart serves as a temporary storage for items and basic information
	cart.ShippingCost = 0
	if cart.CouponCode == "" {
		cart.DiscountAmount = 0
	}

	// Update final total
	cart.UpdateTotals()

	return nil
}

// UpdateCart updates cart properties
func (s *CartService) UpdateCart(tenantID, cartID uuid.UUID, req UpdateCartRequest) (*CartResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, sharedErrors.NewValidationError(err.Error(), nil)
	}

	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Check if cart can be modified
	if tempErr := cart.CanModify(); tempErr != nil {
		return nil, tempErr
	}

	// Update addresses
	if req.ShippingAddress != nil {
		cart.ShippingAddress = req.ShippingAddress
	}
	if req.BillingAddress != nil {
		cart.BillingAddress = req.BillingAddress
	}

	// Update shipping method
	if req.ShippingMethodID != nil {
		cart.ShippingMethodID = req.ShippingMethodID
	}

	// Update notes
	if req.Notes != nil {
		cart.Notes = strings.TrimSpace(*req.Notes)
	}

	// Handle coupon
	if req.CouponCode != nil {
		if *req.CouponCode == "" {
			// Remove coupon
			cart.CouponCode = ""
			cart.DiscountAmount = 0
		} else {
			// Apply coupon (validation will be done during order processing)
			cart.CouponCode = *req.CouponCode
		}
	}

	// Recalculate cart totals
	if recalcErr := s.recalculateCart(cart); recalcErr != nil {
		return nil, recalcErr
	}

	// Update cart
	updatedCart, err := s.repo.UpdateCart(cart)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	return s.buildCartResponse(updatedCart), nil
}

// GetEstimates calculates shipping estimates
func (s *CartService) GetEstimates(tenantID, cartID uuid.UUID, req EstimateRequest) (*EstimateResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, sharedErrors.NewValidationError(err.Error(), nil)
	}

	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Create temporary cart with shipping address for calculations
	tempCart := *cart
	tempCart.ShippingAddress = req.ShippingAddress

	// Estimates will be calculated during order processing
	// Return basic estimates for now
	var shippingMethods []*ShippingMethod
	shippingCost := 0.0

	// Calculate shipping cost based on the provided address
	if tempCart.ShippingAddress != nil {
		// Use the shipping address for location-based calculations
		// This ensures the ShippingAddress field is actually utilized
		shippingCost = s.calculateShippingCostByAddress(tempCart.ShippingAddress)
	}

	// Initialize with basic shipping method if none available
	if len(shippingMethods) == 0 {
		shippingMethods = []*ShippingMethod{
			{
				ID:            uuid.New(),
				Name:          "Standard Shipping",
				Description:   "5-7 business days",
				Cost:          shippingCost,
				EstimatedDays: 7,
			},
		}
	}

	// Calculate total using tempCart which includes the provided shipping address
	total := tempCart.Subtotal + shippingCost - tempCart.DiscountAmount

	// Convert shipping methods to response format
	responseShippingMethods := make([]ShippingMethod, len(shippingMethods))
	for i, method := range shippingMethods {
		if method != nil {
			responseShippingMethods[i] = *method
		}
	}

	return &EstimateResponse{
		ShippingMethods: responseShippingMethods,
		Subtotal:        tempCart.Subtotal,
		Total:           total,
	}, nil
}

// ProcessGuestCheckout processes guest checkout
func (s *CartService) ProcessGuestCheckout(tenantID uuid.UUID, req GuestCheckoutRequest) (*GuestCheckoutResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, sharedErrors.NewValidationError(err.Error(), nil)
	}

	// Get guest cart
	cart, err := s.repo.FindCartBySessionID(tenantID, req.SessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Check if cart can be modified
	if tempErr := cart.CanModify(); tempErr != nil {
		return nil, sharedErrors.NewBadRequestError(tempErr.Error())
	}

	// Check if cart has items
	if len(cart.Items) == 0 {
		return nil, sharedErrors.NewBadRequestError("cart is empty")
	}

	// Update cart with checkout information
	cart.ShippingAddress = &req.ShippingAddress
	cart.BillingAddress = &req.BillingAddress
	cart.ShippingMethodID = &req.ShippingMethodID

	// Recalculate totals
	if recalcErr := s.recalculateCart(cart); recalcErr != nil {
		return nil, recalcErr
	}

	// TODO: Integrate with order service to create order
	// TODO: Integrate with payment service to process payment
	// For now, return a mock response
	orderID := uuid.New()
	orderNumber := "ORD-" + orderID.String()[:8]

	// Mark cart as converted
	cart.MarkAsConverted()
	_, err = s.repo.UpdateCart(cart)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	return &GuestCheckoutResponse{
		OrderID:     orderID,
		OrderNumber: orderNumber,
		Total:       cart.Total,
		Status:      "pending",
	}, nil
}

// buildCartResponse builds a cart response with additional calculated fields
func (s *CartService) buildCartResponse(cart *Cart) *CartResponse {
	savingsAmount := 0.0
	if cart.Items != nil {
		for _, item := range cart.Items {
			savingsAmount += item.GetDiscountAmount()
		}
	}

	return &CartResponse{
		Cart:            cart,
		ItemCount:       cart.GetItemCount(),
		UniqueItemCount: cart.GetUniqueItemCount(),
		SavingsAmount:   savingsAmount,
	}
}

// calculateShippingCostByAddress calculates shipping cost based on address
func (s *CartService) calculateShippingCostByAddress(address *Address) float64 {
	// Simple calculation based on address - in a real implementation
	// this would integrate with shipping providers and calculate based on
	// distance, weight, dimensions, etc.
	if address == nil {
		return 0.0
	}
	
	// Basic shipping cost calculation - could be enhanced based on:
	// - address.Country, address.State for different rates
	// - address.PostalCode for zone-based pricing
	// - integration with shipping providers
	return 5.99 // Default shipping cost
}
