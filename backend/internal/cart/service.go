package cart

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

// Request/Response Structures
type CreateCartRequest struct {
	CustomerID *uuid.UUID `json:"customer_id,omitempty"`
	SessionID  string     `json:"session_id,omitempty"`
	Currency   string     `json:"currency" validate:"required,len=3"`
	Notes      string     `json:"notes,omitempty" validate:"max=500"`
}

type AddItemRequest struct {
	ProductID      uuid.UUID              `json:"product_id" validate:"required"`
	VariantID      *uuid.UUID             `json:"variant_id,omitempty"`
	Quantity       int                    `json:"quantity" validate:"required,min=1,max=100"`
	Customizations map[string]interface{} `json:"customizations,omitempty"`
	Notes          string                 `json:"notes,omitempty" validate:"max=200"`
}

type UpdateItemRequest struct {
	Quantity       *int                   `json:"quantity,omitempty" validate:"omitempty,min=1,max=100"`
	Customizations map[string]interface{} `json:"customizations,omitempty"`
	Notes          string                 `json:"notes,omitempty" validate:"max=200"`
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
	Taxes           TaxEstimate      `json:"taxes"`
	Subtotal        float64          `json:"subtotal"`
	Total           float64          `json:"total"`
}

type TaxEstimate struct {
	Amount float64 `json:"amount"`
	Rate   float64 `json:"rate"`
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
	GetProduct(tenantID uuid.UUID, productID string) (*ProductInfo, error)
	GetProductVariant(tenantID, productID, variantID uuid.UUID) (*VariantInfo, error)
	CheckInventory(tenantID, productID uuid.UUID, variantID *uuid.UUID, quantity int) (bool, error)
	ReserveInventory(tenantID, productID uuid.UUID, variantID *uuid.UUID, quantity int) error
	ReleaseInventory(tenantID, productID uuid.UUID, variantID *uuid.UUID, quantity int) error
}

type DiscountService interface {
	ValidateCoupon(tenantID uuid.UUID, couponCode string, cartTotal float64) (*CouponInfo, error)
	CalculateDiscount(tenantID uuid.UUID, cart *Cart, couponCode string) (float64, error)
}

type TaxService interface {
	CalculateTax(tenantID uuid.UUID, cart *Cart) (float64, error)
}

type ShippingService interface {
	CalculateShipping(tenantID uuid.UUID, cart *Cart, methodID uuid.UUID) (float64, error)
	GetAvailableShippingMethods(tenantID uuid.UUID, cart *Cart) ([]*ShippingMethod, error)
}

// External service data structures
type ProductInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Price       float64   `json:"price"`
	ComparePrice float64  `json:"compare_price"`
	Image       string    `json:"image"`
	SKU         string    `json:"sku"`
	IsAvailable bool      `json:"is_available"`
}

type VariantInfo struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
	SKU   string    `json:"sku"`
	Image string    `json:"image"`
}

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
	taxService      TaxService
	shippingService ShippingService
	cartExpiration  time.Duration
}

// NewCartService creates a new cart service implementation
func NewCartService(repo Repository, productService ProductService, discountService DiscountService, taxService TaxService, shippingService ShippingService) *CartService {
	return &CartService{
		repo:            repo,
		validator:       validator.New(),
		productService:  productService,
		discountService: discountService,
		taxService:      taxService,
		shippingService: shippingService,
		cartExpiration:  24 * time.Hour * 30, // 30 days default
	}
}

// NewService creates a new cart service (interface compatibility)
func NewService(repo Repository, productService ProductService, discountService DiscountService, taxService TaxService, shippingService ShippingService) Service {
	return NewCartService(repo, productService, discountService, taxService, shippingService)
}

// CreateCart creates a new cart
func (s *CartService) CreateCart(tenantID uuid.UUID, req CreateCartRequest) (*CartResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Validate that either customer ID or session ID is provided
	if req.CustomerID == nil && req.SessionID == "" {
		return nil, errors.New("either customer_id or session_id must be provided")
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
		return nil, err
	}

	return s.buildCartResponse(savedCart), nil
}

// GetCart retrieves a cart by ID
func (s *CartService) GetCart(tenantID, cartID uuid.UUID) (*CartResponse, error) {
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Check if cart is expired and update status
	if cart.IsExpired() && cart.Status == StatusActive {
		cart.Status = StatusExpired
		s.repo.UpdateCart(cart)
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
		return nil, ErrCartNotFound
	}

	// Check if cart is expired
	if cart.IsExpired() && cart.Status == StatusActive {
		cart.Status = StatusExpired
		s.repo.UpdateCart(cart)
		return nil, ErrCartExpired
	}

	return s.buildCartResponse(cart), nil
}

// GetCartBySession retrieves cart for a guest session
func (s *CartService) GetCartBySession(tenantID uuid.UUID, sessionID string) (*CartResponse, error) {
	cart, err := s.repo.FindCartBySessionID(tenantID, sessionID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Check if cart is expired
	if cart.IsExpired() && cart.Status == StatusActive {
		cart.Status = StatusExpired
		s.repo.UpdateCart(cart)
		return nil, ErrCartExpired
	}

	return s.buildCartResponse(cart), nil
}

// AddItem adds an item to the cart
func (s *CartService) AddItem(tenantID, cartID uuid.UUID, req AddItemRequest) (*CartResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Check if cart can be modified
	if err := cart.CanModify(); err != nil {
		return nil, err
	}

	// Get product information
	product, err := s.productService.GetProduct(tenantID, req.ProductID.String())
	if err != nil {
		return nil, ErrProductNotFound
	}

	if !product.IsAvailable {
		return nil, errors.New("product is not available")
	}

	// Get variant information if specified
	var variant *VariantInfo
	if req.VariantID != nil {
		variant, err = s.productService.GetProductVariant(tenantID, req.ProductID, *req.VariantID)
		if err != nil {
			return nil, errors.New("product variant not found")
		}
	}

	// Check inventory
	available, err := s.productService.CheckInventory(tenantID, req.ProductID, req.VariantID, req.Quantity)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, ErrInsufficientStock
	}

	// Check if item already exists in cart
	existingItem := cart.FindItem(req.ProductID, req.VariantID)
	if existingItem != nil {
		// Update existing item quantity
		newQuantity := existingItem.Quantity + req.Quantity
		
		// Check inventory for new total quantity
		available, err := s.productService.CheckInventory(tenantID, req.ProductID, req.VariantID, newQuantity)
		if err != nil {
			return nil, err
		}
		if !available {
			return nil, ErrInsufficientStock
		}
		
		existingItem.Quantity = newQuantity
		existingItem.Customizations = req.Customizations
		existingItem.Notes = strings.TrimSpace(req.Notes)
		existingItem.CalculateLineTotal()
		
		if _, err := s.repo.UpdateCartItem(existingItem); err != nil {
			return nil, err
		}
	} else {
		// Create new cart item
		price := product.Price
		comparePrice := product.ComparePrice
		sku := product.SKU
		image := product.Image
		variantName := ""
		
		if variant != nil {
			price = variant.Price
			sku = variant.SKU
			variantName = variant.Name
			if variant.Image != "" {
				image = variant.Image
			}
		}
		
		item := &CartItem{
			ID:             uuid.New(),
			CartID:         cartID,
			ProductID:      req.ProductID,
			VariantID:      req.VariantID,
			ProductName:    product.Name,
			ProductSlug:    product.Slug,
			VariantName:    variantName,
			SKU:            sku,
			Price:          price,
			ComparePrice:   comparePrice,
			Image:          image,
			Quantity:       req.Quantity,
			Customizations: req.Customizations,
			Notes:          strings.TrimSpace(req.Notes),
		}
		
		item.CalculateLineTotal()
		
		if _, err := s.repo.AddCartItem(item); err != nil {
			return nil, err
		}
		
		cart.Items = append(cart.Items, *item)
	}

	// Recalculate cart totals
	if err := s.recalculateCart(cart); err != nil {
		return nil, err
	}

	// Update cart
	updatedCart, err := s.repo.UpdateCart(cart)
	if err != nil {
		return nil, err
	}

	return s.buildCartResponse(updatedCart), nil
}

// UpdateItem updates a cart item
func (s *CartService) UpdateItem(tenantID, cartID, itemID uuid.UUID, req UpdateItemRequest) (*CartItem, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Check if cart can be modified
	if err := cart.CanModify(); err != nil {
		return nil, err
	}

	// Find cart item
	item, err := s.repo.FindCartItem(tenantID, cartID, itemID)
	if err != nil {
		return nil, ErrItemNotFound
	}

	// Update quantity if provided
	if req.Quantity != nil {
		// Check inventory
		available, err := s.productService.CheckInventory(tenantID, item.ProductID, item.VariantID, *req.Quantity)
		if err != nil {
			return nil, err
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
	updatedItem, err := s.repo.UpdateCartItem(item)
	if err != nil {
		return nil, err
	}

	// Reload cart with updated items
	cart, err = s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, err
	}

	// Recalculate cart totals
	if err := s.recalculateCart(cart); err != nil {
		return nil, err
	}

	// Update cart
	_, err = s.repo.UpdateCart(cart)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

// RemoveItem removes an item from the cart
func (s *CartService) RemoveItem(tenantID, cartID, itemID uuid.UUID) error {
	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return ErrCartNotFound
	}

	// Check if cart can be modified
	if err := cart.CanModify(); err != nil {
		return err
	}

	// Remove item
	if err := s.repo.RemoveCartItem(tenantID, cartID, itemID); err != nil {
		return ErrItemNotFound
	}

	// Reload cart
	cart, err = s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return err
	}

	// Recalculate cart totals
	if err := s.recalculateCart(cart); err != nil {
		return err
	}

	// Update cart
	_, err = s.repo.UpdateCart(cart)
	if err != nil {
		return err
	}

	return nil
}

// ClearCart removes all items from the cart
func (s *CartService) ClearCart(tenantID, cartID uuid.UUID) error {
	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return ErrCartNotFound
	}

	// Check if cart can be modified
	if err := cart.CanModify(); err != nil {
		return err
	}

	// Clear all items
	if err := s.repo.ClearCartItems(tenantID, cartID); err != nil {
		return err
	}

	// Reset cart totals
	cart.Items = []CartItem{}
	cart.Subtotal = 0
	cart.TaxAmount = 0
	cart.ShippingCost = 0
	cart.DiscountAmount = 0
	cart.Total = 0
	cart.CouponCode = ""
	cart.DiscountID = nil

	// Update cart
	_, err = s.repo.UpdateCart(cart)
	if err != nil {
		return err
	}

	return nil
}

// ApplyCoupon applies a coupon to the cart
func (s *CartService) ApplyCoupon(tenantID, cartID uuid.UUID, req ApplyCouponRequest) (*CartResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Check if cart can be modified
	if err := cart.CanModify(); err != nil {
		return nil, err
	}

	// Validate coupon
	coupon, err := s.discountService.ValidateCoupon(tenantID, req.CouponCode, cart.Subtotal)
	if err != nil {
		return nil, ErrInvalidCoupon
	}

	// Calculate discount
	discountAmount, err := s.discountService.CalculateDiscount(tenantID, cart, req.CouponCode)
	if err != nil {
		return nil, err
	}

	// Apply coupon
	cart.CouponCode = coupon.Code
	cart.DiscountAmount = discountAmount

	// Recalculate totals
	cart.UpdateTotals()

	// Update cart
	updatedCart, err := s.repo.UpdateCart(cart)
	if err != nil {
		return nil, err
	}

	return s.buildCartResponse(updatedCart), nil
}

// RemoveCoupon removes coupon from the cart
func (s *CartService) RemoveCoupon(tenantID, cartID uuid.UUID) (*CartResponse, error) {
	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Check if cart can be modified
	if err := cart.CanModify(); err != nil {
		return nil, err
	}

	// Remove coupon
	cart.CouponCode = ""
	cart.DiscountAmount = 0
	cart.DiscountID = nil

	// Recalculate totals
	cart.UpdateTotals()

	// Update cart
	updatedCart, err := s.repo.UpdateCart(cart)
	if err != nil {
		return nil, err
	}

	return s.buildCartResponse(updatedCart), nil
}

// UpdateAddress updates shipping/billing address
func (s *CartService) UpdateAddress(tenantID, cartID uuid.UUID, req UpdateAddressRequest) (*CartResponse, error) {
	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Check if cart can be modified
	if err := cart.CanModify(); err != nil {
		return nil, err
	}

	// Update addresses
	if req.ShippingAddress != nil {
		cart.ShippingAddress = req.ShippingAddress
	}
	if req.BillingAddress != nil {
		cart.BillingAddress = req.BillingAddress
	}

	// Recalculate tax and shipping if address changed
	if err := s.recalculateCart(cart); err != nil {
		return nil, err
	}

	// Update cart
	updatedCart, err := s.repo.UpdateCart(cart)
	if err != nil {
		return nil, err
	}

	return s.buildCartResponse(updatedCart), nil
}

// UpdateShipping updates shipping method
func (s *CartService) UpdateShipping(tenantID, cartID uuid.UUID, req UpdateShippingRequest) (*CartResponse, error) {
	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Check if cart can be modified
	if err := cart.CanModify(); err != nil {
		return nil, err
	}

	// Update shipping method
	cart.ShippingMethodID = req.ShippingMethodID

	// Recalculate shipping cost
	if req.ShippingMethodID != nil {
		shippingCost, err := s.shippingService.CalculateShipping(tenantID, cart, *req.ShippingMethodID)
		if err != nil {
			return nil, err
		}
		cart.ShippingCost = shippingCost
	} else {
		cart.ShippingCost = 0
	}

	// Recalculate totals
	cart.UpdateTotals()

	// Update cart
	updatedCart, err := s.repo.UpdateCart(cart)
	if err != nil {
		return nil, err
	}

	return s.buildCartResponse(updatedCart), nil
}

// MergeGuestCart merges guest cart to customer cart
func (s *CartService) MergeGuestCart(tenantID uuid.UUID, sessionID string, customerID uuid.UUID) (*CartResponse, error) {
	if err := s.repo.MergeGuestCartToCustomer(tenantID, sessionID, customerID); err != nil {
		return nil, err
	}

	// Get the merged cart
	cart, err := s.repo.FindCartByCustomerID(tenantID, customerID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Recalculate totals
	if err := s.recalculateCart(cart); err != nil {
		return nil, err
	}

	// Update cart
	updatedCart, err := s.repo.UpdateCart(cart)
	if err != nil {
		return nil, err
	}

	return s.buildCartResponse(updatedCart), nil
}

// AbandonCart marks cart as abandoned
func (s *CartService) AbandonCart(tenantID, cartID uuid.UUID) error {
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return ErrCartNotFound
	}

	if cart.Status == StatusActive {
		cart.MarkAsAbandoned()
		_, err = s.repo.UpdateCart(cart)
	}

	return err
}

// ConvertCart marks cart as converted to order
func (s *CartService) ConvertCart(tenantID, cartID uuid.UUID) error {
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return ErrCartNotFound
	}

	if cart.Status == StatusActive {
		cart.MarkAsConverted()
		_, err = s.repo.UpdateCart(cart)
	}

	return err
}

// DeleteCart soft deletes a cart
func (s *CartService) DeleteCart(tenantID, cartID uuid.UUID) error {
	return s.repo.DeleteCart(tenantID, cartID)
}

// GetCartSummary returns a summary of the cart
func (s *CartService) GetCartSummary(tenantID, cartID uuid.UUID) (*CartSummary, error) {
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, ErrCartNotFound
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

	// Calculate tax if tax service is available and address is provided
	if s.taxService != nil && cart.ShippingAddress != nil {
		taxAmount, err := s.taxService.CalculateTax(cart.TenantID, cart)
		if err == nil {
			cart.TaxAmount = taxAmount
		}
	}

	// Calculate shipping if shipping service is available and method is selected
	if s.shippingService != nil && cart.ShippingMethodID != nil {
		shippingCost, err := s.shippingService.CalculateShipping(cart.TenantID, cart, *cart.ShippingMethodID)
		if err == nil {
			cart.ShippingCost = shippingCost
		}
	}

	// Recalculate discount if coupon is applied
	if s.discountService != nil && cart.CouponCode != "" {
		discountAmount, err := s.discountService.CalculateDiscount(cart.TenantID, cart, cart.CouponCode)
		if err == nil {
			cart.DiscountAmount = discountAmount
		}
	}

	// Update final total
	cart.UpdateTotals()

	return nil
}

// UpdateCart updates cart properties
func (s *CartService) UpdateCart(tenantID, cartID uuid.UUID, req UpdateCartRequest) (*CartResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Check if cart can be modified
	if err := cart.CanModify(); err != nil {
		return nil, err
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
			// Apply coupon
			if s.discountService != nil {
				_, err := s.discountService.ValidateCoupon(tenantID, *req.CouponCode, cart.Subtotal)
				if err != nil {
					return nil, ErrInvalidCoupon
				}
				cart.CouponCode = *req.CouponCode
			}
		}
	}

	// Recalculate cart totals
	if err := s.recalculateCart(cart); err != nil {
		return nil, err
	}

	// Update cart
	updatedCart, err := s.repo.UpdateCart(cart)
	if err != nil {
		return nil, err
	}

	return s.buildCartResponse(updatedCart), nil
}

// GetEstimates calculates shipping and tax estimates
func (s *CartService) GetEstimates(tenantID, cartID uuid.UUID, req EstimateRequest) (*EstimateResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Get cart
	cart, err := s.repo.FindCartByID(tenantID, cartID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Create temporary cart with shipping address for calculations
	tempCart := *cart
	tempCart.ShippingAddress = req.ShippingAddress

	// Get available shipping methods
	var shippingMethods []*ShippingMethod
	if s.shippingService != nil {
		methods, err := s.shippingService.GetAvailableShippingMethods(tenantID, &tempCart)
		if err == nil {
			shippingMethods = methods
		}
	}

	// Calculate shipping cost for specific method if provided
	shippingCost := 0.0
	if req.ShippingMethodID != nil && s.shippingService != nil {
		cost, err := s.shippingService.CalculateShipping(tenantID, &tempCart, *req.ShippingMethodID)
		if err == nil {
			shippingCost = cost
		}
	}

	// Calculate tax estimate
	taxEstimate := TaxEstimate{Amount: 0, Rate: 0}
	if s.taxService != nil {
		taxAmount, err := s.taxService.CalculateTax(tenantID, &tempCart)
		if err == nil {
			taxEstimate.Amount = taxAmount
			if cart.Subtotal > 0 {
				taxEstimate.Rate = taxAmount / cart.Subtotal
			}
		}
	}

	// Calculate total
	total := cart.Subtotal + shippingCost + taxEstimate.Amount - cart.DiscountAmount

	// Convert shipping methods to response format
	responseShippingMethods := make([]ShippingMethod, len(shippingMethods))
	for i, method := range shippingMethods {
		responseShippingMethods[i] = *method
	}

	return &EstimateResponse{
		ShippingMethods: responseShippingMethods,
		Taxes:           taxEstimate,
		Subtotal:        cart.Subtotal,
		Total:           total,
	}, nil
}

// ProcessGuestCheckout processes guest checkout
func (s *CartService) ProcessGuestCheckout(tenantID uuid.UUID, req GuestCheckoutRequest) (*GuestCheckoutResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Get guest cart
	cart, err := s.repo.FindCartBySessionID(tenantID, req.SessionID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Check if cart can be modified
	if err := cart.CanModify(); err != nil {
		return nil, err
	}

	// Check if cart has items
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Update cart with checkout information
	cart.ShippingAddress = &req.ShippingAddress
	cart.BillingAddress = &req.BillingAddress
	cart.ShippingMethodID = &req.ShippingMethodID

	// Recalculate totals
	if err := s.recalculateCart(cart); err != nil {
		return nil, err
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
		return nil, err
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
	for _, item := range cart.Items {
		savingsAmount += item.GetDiscountAmount()
	}

	return &CartResponse{
		Cart:            cart,
		ItemCount:       cart.GetItemCount(),
		UniqueItemCount: cart.GetUniqueItemCount(),
		SavingsAmount:   savingsAmount,
	}
}