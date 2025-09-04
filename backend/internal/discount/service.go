package discount

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service defines the discount service interface
type Service interface {
	// Discount/Coupon operations
	CreateDiscount(ctx context.Context, req CreateDiscountRequest) (*Discount, error)
	GetDiscount(ctx context.Context, tenantID, discountID uuid.UUID) (*Discount, error)
	GetDiscountByCode(ctx context.Context, tenantID uuid.UUID, code string) (*Discount, error)
	GetDiscounts(ctx context.Context, tenantID uuid.UUID, filter DiscountFilter) ([]Discount, error)
	UpdateDiscount(ctx context.Context, tenantID, discountID uuid.UUID, req UpdateDiscountRequest) (*Discount, error)
	DeleteDiscount(ctx context.Context, tenantID, discountID uuid.UUID) error
	
	// Discount validation and application
	ValidateDiscountCode(ctx context.Context, req ValidateDiscountRequest) (*DiscountValidation, error)
	ApplyDiscount(ctx context.Context, req ApplyDiscountRequest) (*DiscountApplication, error)
	RemoveDiscount(ctx context.Context, tenantID uuid.UUID, orderID uuid.UUID) error
	
	// Discount usage tracking
	RecordDiscountUsage(ctx context.Context, usage *DiscountUsage) error
	GetDiscountUsage(ctx context.Context, tenantID, discountID uuid.UUID, filter UsageFilter) ([]DiscountUsage, error)
	GetCustomerDiscountUsage(ctx context.Context, tenantID uuid.UUID, customerEmail string, discountID uuid.UUID) (int, error)
	
	// Gift card operations
	CreateGiftCard(ctx context.Context, req CreateGiftCardRequest) (*GiftCard, error)
	GetGiftCard(ctx context.Context, tenantID uuid.UUID, code string) (*GiftCard, error)
	GetGiftCards(ctx context.Context, tenantID uuid.UUID, filter GiftCardFilter) ([]GiftCard, error)
	UpdateGiftCard(ctx context.Context, tenantID, giftCardID uuid.UUID, req UpdateGiftCardRequest) (*GiftCard, error)
	DeleteGiftCard(ctx context.Context, tenantID, giftCardID uuid.UUID) error
	
	// Gift card usage
	ValidateGiftCard(ctx context.Context, tenantID uuid.UUID, code string) (*GiftCardValidation, error)
	UseGiftCard(ctx context.Context, req UseGiftCardRequest) (*GiftCardTransaction, error)
	RefillGiftCard(ctx context.Context, req RefillGiftCardRequest) (*GiftCardTransaction, error)
	GetGiftCardTransactions(ctx context.Context, tenantID, giftCardID uuid.UUID) ([]GiftCardTransaction, error)
	
	// Store credit operations
	GetStoreCredit(ctx context.Context, tenantID, customerID uuid.UUID) (*StoreCredit, error)
	AddStoreCredit(ctx context.Context, req AddStoreCreditRequest) (*StoreCreditTransaction, error)
	UseStoreCredit(ctx context.Context, req UseStoreCreditRequest) (*StoreCreditTransaction, error)
	GetStoreCreditTransactions(ctx context.Context, tenantID, customerID uuid.UUID, filter StoreCreditFilter) ([]StoreCreditTransaction, error)
	
	// Analytics and reporting
	GetDiscountStats(ctx context.Context, tenantID uuid.UUID, period string) (*DiscountStats, error)
	GetTopDiscounts(ctx context.Context, tenantID uuid.UUID, limit int) ([]DiscountPerformance, error)
	GetDiscountRevenue(ctx context.Context, tenantID uuid.UUID, period string) (*RevenueImpact, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new discount service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Request/Response DTOs
type CreateDiscountRequest struct {
	Code                   string             `json:"code" validate:"required"`
	Title                  string             `json:"title" validate:"required"`
	Description            string             `json:"description"`
	Type                   DiscountType       `json:"type" validate:"required"`
	Value                  float64            `json:"value" validate:"required,gt=0"`
	Currency               string             `json:"currency"`
	MinOrderAmount         *float64           `json:"min_order_amount"`
	MinItemQuantity        *int               `json:"min_item_quantity"`
	Target                 DiscountTarget     `json:"target"`
	TargetProductIDs       []string           `json:"target_product_ids"`
	TargetCategoryIDs      []string           `json:"target_category_ids"`
	TargetCollectionIDs    []string           `json:"target_collection_ids"`
	ExcludeProductIDs      []string           `json:"exclude_product_ids"`
	UsageLimit             *int               `json:"usage_limit"`
	UsageLimitType         DiscountUsageLimit `json:"usage_limit_type"`
	CustomerUsageLimit     *int               `json:"customer_usage_limit"`
	StartsAt               *time.Time         `json:"starts_at"`
	ExpiresAt              *time.Time         `json:"expires_at"`
	CustomerEligibility    string             `json:"customer_eligibility"`
	EligibleCustomerIDs    []string           `json:"eligible_customer_ids"`
	EligibleCustomerGroups []string           `json:"eligible_customer_groups"`
	BuyQuantity            *int               `json:"buy_quantity"`
	GetQuantity            *int               `json:"get_quantity"`
	GetValue               *float64           `json:"get_value"`
	Stackable              bool               `json:"stackable"`
	StackableWith          []string           `json:"stackable_with"`
	ExclusiveGroup         string             `json:"exclusive_group"`
	ApplyOnce              bool               `json:"apply_once"`
	ShowInStorefront       bool               `json:"show_in_storefront"`
	RequiresCode           bool               `json:"requires_code"`
}

type UpdateDiscountRequest struct {
	Title                  *string            `json:"title"`
	Description            *string            `json:"description"`
	Value                  *float64           `json:"value"`
	MinOrderAmount         *float64           `json:"min_order_amount"`
	MinItemQuantity        *int               `json:"min_item_quantity"`
	Target                 *DiscountTarget    `json:"target"`
	TargetProductIDs       []string           `json:"target_product_ids"`
	TargetCategoryIDs      []string           `json:"target_category_ids"`
	TargetCollectionIDs    []string           `json:"target_collection_ids"`
	ExcludeProductIDs      []string           `json:"exclude_product_ids"`
	UsageLimit             *int               `json:"usage_limit"`
	CustomerUsageLimit     *int               `json:"customer_usage_limit"`
	StartsAt               *time.Time         `json:"starts_at"`
	ExpiresAt              *time.Time         `json:"expires_at"`
	Status                 *DiscountStatus    `json:"status"`
	CustomerEligibility    *string            `json:"customer_eligibility"`
	EligibleCustomerIDs    []string           `json:"eligible_customer_ids"`
	EligibleCustomerGroups []string           `json:"eligible_customer_groups"`
	Stackable              *bool              `json:"stackable"`
	ShowInStorefront       *bool              `json:"show_in_storefront"`
}

type DiscountFilter struct {
	Status         []DiscountStatus `json:"status"`
	Type           []DiscountType   `json:"type"`
	Target         []DiscountTarget `json:"target"`
	Search         string           `json:"search"`
	IsExpired      *bool            `json:"is_expired"`
	IsActive       *bool            `json:"is_active"`
	CreatedBy      *uuid.UUID       `json:"created_by"`
	StartDate      *time.Time       `json:"start_date"`
	EndDate        *time.Time       `json:"end_date"`
	SortBy         string           `json:"sort_by"`
	SortOrder      string           `json:"sort_order"`
	Page           int              `json:"page"`
	Limit          int              `json:"limit"`
}

type ValidateDiscountRequest struct {
	Code           string     `json:"code" validate:"required"`
	CustomerID     *uuid.UUID `json:"customer_id"`
	CustomerEmail  string     `json:"customer_email"`
	OrderAmount    float64    `json:"order_amount" validate:"required,gt=0"`
	ItemQuantity   int        `json:"item_quantity" validate:"required,gt=0"`
	ProductIDs     []string   `json:"product_ids"`
	CategoryIDs    []string   `json:"category_ids"`
	AppliedDiscounts []string `json:"applied_discounts"` // Currently applied discount codes
}

type DiscountValidation struct {
	Valid           bool    `json:"valid"`
	Discount        *Discount `json:"discount,omitempty"`
	DiscountAmount  float64 `json:"discount_amount"`
	Message         string  `json:"message"`
	RemainingUsage  *int    `json:"remaining_usage,omitempty"`
	CanStack        bool    `json:"can_stack"`
}

type ApplyDiscountRequest struct {
	TenantID       uuid.UUID  `json:"tenant_id" validate:"required"`
	Code           string     `json:"code" validate:"required"`
	OrderID        uuid.UUID  `json:"order_id" validate:"required"`
	CustomerID     *uuid.UUID `json:"customer_id"`
	CustomerEmail  string     `json:"customer_email" validate:"required"`
	OrderAmount    float64    `json:"order_amount" validate:"required,gt=0"`
	ItemQuantity   int        `json:"item_quantity" validate:"required,gt=0"`
	ProductIDs     []string   `json:"product_ids"`
	CategoryIDs    []string   `json:"category_ids"`
	IPAddress      string     `json:"ip_address"`
	UserAgent      string     `json:"user_agent"`
}

type DiscountApplication struct {
	Applied        bool      `json:"applied"`
	Discount       *Discount `json:"discount,omitempty"`
	DiscountAmount float64   `json:"discount_amount"`
	Usage          *DiscountUsage `json:"usage,omitempty"`
	Message        string    `json:"message"`
}

type UsageFilter struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}

// Gift card DTOs
type CreateGiftCardRequest struct {
	Code           string     `json:"code"`
	InitialValue   float64    `json:"initial_value" validate:"required,gt=0"`
	Currency       string     `json:"currency" validate:"required"`
	RecipientName  string     `json:"recipient_name"`
	RecipientEmail string     `json:"recipient_email"`
	Message        string     `json:"message"`
	PurchasedBy    *uuid.UUID `json:"purchased_by"`
	ExpiresAt      *time.Time `json:"expires_at"`
	IsRefillable   bool       `json:"is_refillable"`
	Notes          string     `json:"notes"`
}

type UpdateGiftCardRequest struct {
	Status         *DiscountStatus `json:"status"`
	RecipientName  *string         `json:"recipient_name"`
	RecipientEmail *string         `json:"recipient_email"`
	Message        *string         `json:"message"`
	ExpiresAt      *time.Time      `json:"expires_at"`
	IsRefillable   *bool           `json:"is_refillable"`
	Notes          *string         `json:"notes"`
}

type GiftCardFilter struct {
	Status     []DiscountStatus `json:"status"`
	Search     string           `json:"search"`
	IsExpired  *bool            `json:"is_expired"`
	StartDate  *time.Time       `json:"start_date"`
	EndDate    *time.Time       `json:"end_date"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
}

type GiftCardValidation struct {
	Valid         bool      `json:"valid"`
	GiftCard      *GiftCard `json:"gift_card,omitempty"`
	AvailableAmount float64 `json:"available_amount"`
	Message       string    `json:"message"`
}

type UseGiftCardRequest struct {
	TenantID      uuid.UUID  `json:"tenant_id" validate:"required"`
	Code          string     `json:"code" validate:"required"`
	Amount        float64    `json:"amount" validate:"required,gt=0"`
	OrderID       uuid.UUID  `json:"order_id" validate:"required"`
	OrderNumber   string     `json:"order_number" validate:"required"`
	CustomerID    *uuid.UUID `json:"customer_id"`
	CustomerEmail string     `json:"customer_email" validate:"required"`
}

type RefillGiftCardRequest struct {
	TenantID    uuid.UUID  `json:"tenant_id" validate:"required"`
	GiftCardID  uuid.UUID  `json:"gift_card_id" validate:"required"`
	Amount      float64    `json:"amount" validate:"required,gt=0"`
	Description string     `json:"description"`
	ProcessedBy uuid.UUID  `json:"processed_by" validate:"required"`
}

// Store credit DTOs
type AddStoreCreditRequest struct {
	TenantID      uuid.UUID  `json:"tenant_id" validate:"required"`
	CustomerID    uuid.UUID  `json:"customer_id" validate:"required"`
	Amount        float64    `json:"amount" validate:"required,gt=0"`
	Currency      string     `json:"currency" validate:"required"`
	Description   string     `json:"description"`
	RefundID      *uuid.UUID `json:"refund_id"`
	ReturnID      *uuid.UUID `json:"return_id"`
	ProcessedBy   *uuid.UUID `json:"processed_by"`
	ExpiresAt     *time.Time `json:"expires_at"`
}

type UseStoreCreditRequest struct {
	TenantID      uuid.UUID  `json:"tenant_id" validate:"required"`
	CustomerID    uuid.UUID  `json:"customer_id" validate:"required"`
	Amount        float64    `json:"amount" validate:"required,gt=0"`
	OrderID       uuid.UUID  `json:"order_id" validate:"required"`
	OrderNumber   string     `json:"order_number" validate:"required"`
	Description   string     `json:"description"`
}

type StoreCreditFilter struct {
	Type      []string   `json:"type"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}

// Analytics DTOs
type DiscountStats struct {
	TotalDiscounts    int                        `json:"total_discounts"`
	ActiveDiscounts   int                        `json:"active_discounts"`
	TotalUsage        int                        `json:"total_usage"`
	TotalSavings      float64                    `json:"total_savings"`
	DiscountsByType   map[DiscountType]int       `json:"discounts_by_type"`
	DiscountsByStatus map[DiscountStatus]int     `json:"discounts_by_status"`
	TopDiscounts      []DiscountPerformance      `json:"top_discounts"`
	RevenueImpact     *RevenueImpact            `json:"revenue_impact"`
}

type DiscountPerformance struct {
	DiscountID     uuid.UUID `json:"discount_id"`
	Code           string    `json:"code"`
	Title          string    `json:"title"`
	UsageCount     int       `json:"usage_count"`
	TotalSavings   float64   `json:"total_savings"`
	Revenue        float64   `json:"revenue"`
	ConversionRate float64   `json:"conversion_rate"`
}

type RevenueImpact struct {
	TotalOrders       int     `json:"total_orders"`
	OrdersWithDiscount int    `json:"orders_with_discount"`
	DiscountRate      float64 `json:"discount_rate"`
	AverageOrderValue float64 `json:"average_order_value"`
	AverageDiscount   float64 `json:"average_discount"`
	RevenueWithDiscount float64 `json:"revenue_with_discount"`
	RevenueWithoutDiscount float64 `json:"revenue_without_discount"`
}

// Implementation methods (TODO: implement business logic)
func (s *service) CreateDiscount(ctx context.Context, req CreateDiscountRequest) (*Discount, error) {
	// TODO: Implement discount creation with validation
	return nil, fmt.Errorf("TODO: implement CreateDiscount")
}

func (s *service) GetDiscount(ctx context.Context, tenantID, discountID uuid.UUID) (*Discount, error) {
	return s.repo.GetDiscountByID(ctx, tenantID, discountID)
}

func (s *service) GetDiscountByCode(ctx context.Context, tenantID uuid.UUID, code string) (*Discount, error) {
	// Normalize code
	code = strings.ToUpper(strings.TrimSpace(code))
	return s.repo.GetDiscountByCode(ctx, tenantID, code)
}

func (s *service) GetDiscounts(ctx context.Context, tenantID uuid.UUID, filter DiscountFilter) ([]Discount, error) {
	return s.repo.GetDiscounts(ctx, tenantID, filter)
}

func (s *service) UpdateDiscount(ctx context.Context, tenantID, discountID uuid.UUID, req UpdateDiscountRequest) (*Discount, error) {
	// TODO: Implement discount update with validation
	return nil, fmt.Errorf("TODO: implement UpdateDiscount")
}

func (s *service) DeleteDiscount(ctx context.Context, tenantID, discountID uuid.UUID) error {
	// TODO: Check if discount has been used before allowing deletion
	return s.repo.DeleteDiscount(ctx, tenantID, discountID)
}

func (s *service) ValidateDiscountCode(ctx context.Context, req ValidateDiscountRequest) (*DiscountValidation, error) {
	// Get discount by code
	discount, err := s.repo.GetDiscountByCode(ctx, uuid.Nil, req.Code) // TODO: Add tenantID to request
	if err != nil {
		return &DiscountValidation{
			Valid:   false,
			Message: "Invalid discount code",
		}, nil
	}

	// Check if discount is active
	if !discount.IsActive() {
		return &DiscountValidation{
			Valid:   false,
			Message: "Discount code is not active",
		}, nil
	}

	// Check customer eligibility
	if req.CustomerID != nil && !discount.CanUseDiscount(req.CustomerID, req.CustomerEmail, 0) {
		return &DiscountValidation{
			Valid:   false,
			Message: "You are not eligible for this discount",
		}, nil
	}

	// Calculate discount amount
	discountAmount, err := discount.CalculateDiscount(req.OrderAmount, req.ItemQuantity)
	if err != nil {
		return &DiscountValidation{
			Valid:   false,
			Message: err.Error(),
		}, nil
	}

	// Check remaining usage
	var remainingUsage *int
	if discount.UsageLimit != nil {
		remaining := *discount.UsageLimit - discount.UsageCount
		remainingUsage = &remaining
		if remaining <= 0 {
			return &DiscountValidation{
				Valid:   false,
				Message: "Discount usage limit exceeded",
			}, nil
		}
	}

	return &DiscountValidation{
		Valid:          true,
		Discount:       discount,
		DiscountAmount: discountAmount,
		Message:        "Discount code is valid",
		RemainingUsage: remainingUsage,
		CanStack:       discount.Stackable,
	}, nil
}

func (s *service) ApplyDiscount(ctx context.Context, req ApplyDiscountRequest) (*DiscountApplication, error) {
	// First validate the discount code
	validationReq := ValidateDiscountRequest{
		Code:          req.Code,
		CustomerID:    req.CustomerID,
		CustomerEmail: req.CustomerEmail,
		OrderAmount:   req.OrderAmount,
		ItemQuantity:  req.ItemQuantity,
		ProductIDs:    req.ProductIDs,
		CategoryIDs:   req.CategoryIDs,
	}

	validation, err := s.ValidateDiscountCode(ctx, validationReq)
	if err != nil {
		return nil, fmt.Errorf("failed to validate discount: %w", err)
	}

	if !validation.Valid {
		return &DiscountApplication{
			Applied: false,
			Message: validation.Message,
		}, nil
	}

	// Create discount usage record
	usage := &DiscountUsage{
		ID:            uuid.New(),
		DiscountID:    validation.Discount.ID,
		TenantID:      req.TenantID,
		OrderID:       req.OrderID,
		OrderNumber:   "", // Will be set by order service
		CustomerID:    req.CustomerID,
		CustomerEmail: req.CustomerEmail,
		DiscountAmount: validation.DiscountAmount,
		IPAddress:     req.IPAddress,
		UserAgent:     req.UserAgent,
		CreatedAt:     time.Now(),
	}

	// Record the usage
	if err := s.RecordDiscountUsage(ctx, usage); err != nil {
		return nil, fmt.Errorf("failed to record discount usage: %w", err)
	}

	// Update discount usage count
	if err := s.repo.IncrementUsageCount(ctx, req.TenantID, validation.Discount.ID); err != nil {
		return nil, fmt.Errorf("failed to update discount usage count: %w", err)
	}

	return &DiscountApplication{
		Applied:        true,
		Discount:       validation.Discount,
		DiscountAmount: validation.DiscountAmount,
		Usage:          usage,
		Message:        "Discount applied successfully",
	}, nil
}

func (s *service) RemoveDiscount(ctx context.Context, tenantID uuid.UUID, orderID uuid.UUID) error {
	// TODO: Implement discount removal and update usage counts
	return fmt.Errorf("TODO: implement RemoveDiscount")
}

func (s *service) RecordDiscountUsage(ctx context.Context, usage *DiscountUsage) error {
	return s.repo.CreateDiscountUsage(ctx, usage)
}

func (s *service) GetDiscountUsage(ctx context.Context, tenantID, discountID uuid.UUID, filter UsageFilter) ([]DiscountUsage, error) {
	return s.repo.GetDiscountUsage(ctx, tenantID, discountID, filter)
}

func (s *service) GetCustomerDiscountUsage(ctx context.Context, tenantID uuid.UUID, customerEmail string, discountID uuid.UUID) (int, error) {
	return s.repo.GetCustomerDiscountUsageCount(ctx, tenantID, customerEmail, discountID)
}

// Gift card methods
func (s *service) CreateGiftCard(ctx context.Context, req CreateGiftCardRequest) (*GiftCard, error) {
	// TODO: Implement gift card creation with code generation
	return nil, fmt.Errorf("TODO: implement CreateGiftCard")
}

func (s *service) GetGiftCard(ctx context.Context, tenantID uuid.UUID, code string) (*GiftCard, error) {
	return s.repo.GetGiftCardByCode(ctx, tenantID, code)
}

func (s *service) GetGiftCards(ctx context.Context, tenantID uuid.UUID, filter GiftCardFilter) ([]GiftCard, error) {
	return s.repo.GetGiftCards(ctx, tenantID, filter)
}

func (s *service) UpdateGiftCard(ctx context.Context, tenantID, giftCardID uuid.UUID, req UpdateGiftCardRequest) (*GiftCard, error) {
	// TODO: Implement gift card update
	return nil, fmt.Errorf("TODO: implement UpdateGiftCard")
}

func (s *service) DeleteGiftCard(ctx context.Context, tenantID, giftCardID uuid.UUID) error {
	return s.repo.DeleteGiftCard(ctx, tenantID, giftCardID)
}

func (s *service) ValidateGiftCard(ctx context.Context, tenantID uuid.UUID, code string) (*GiftCardValidation, error) {
	// TODO: Implement gift card validation
	return nil, fmt.Errorf("TODO: implement ValidateGiftCard")
}

func (s *service) UseGiftCard(ctx context.Context, req UseGiftCardRequest) (*GiftCardTransaction, error) {
	// TODO: Implement gift card usage with balance updates
	return nil, fmt.Errorf("TODO: implement UseGiftCard")
}

func (s *service) RefillGiftCard(ctx context.Context, req RefillGiftCardRequest) (*GiftCardTransaction, error) {
	// TODO: Implement gift card refill
	return nil, fmt.Errorf("TODO: implement RefillGiftCard")
}

func (s *service) GetGiftCardTransactions(ctx context.Context, tenantID, giftCardID uuid.UUID) ([]GiftCardTransaction, error) {
	return s.repo.GetGiftCardTransactions(ctx, tenantID, giftCardID)
}

// Store credit methods
func (s *service) GetStoreCredit(ctx context.Context, tenantID, customerID uuid.UUID) (*StoreCredit, error) {
	return s.repo.GetStoreCredit(ctx, tenantID, customerID)
}

func (s *service) AddStoreCredit(ctx context.Context, req AddStoreCreditRequest) (*StoreCreditTransaction, error) {
	// TODO: Implement store credit addition with transaction recording
	return nil, fmt.Errorf("TODO: implement AddStoreCredit")
}

func (s *service) UseStoreCredit(ctx context.Context, req UseStoreCreditRequest) (*StoreCreditTransaction, error) {
	// TODO: Implement store credit usage with balance updates
	return nil, fmt.Errorf("TODO: implement UseStoreCredit")
}

func (s *service) GetStoreCreditTransactions(ctx context.Context, tenantID, customerID uuid.UUID, filter StoreCreditFilter) ([]StoreCreditTransaction, error) {
	return s.repo.GetStoreCreditTransactions(ctx, tenantID, customerID, filter)
}

// Analytics methods
func (s *service) GetDiscountStats(ctx context.Context, tenantID uuid.UUID, period string) (*DiscountStats, error) {
	// TODO: Implement discount analytics
	return nil, fmt.Errorf("TODO: implement GetDiscountStats")
}

func (s *service) GetTopDiscounts(ctx context.Context, tenantID uuid.UUID, limit int) ([]DiscountPerformance, error) {
	// TODO: Implement top performing discounts query
	return nil, fmt.Errorf("TODO: implement GetTopDiscounts")
}

func (s *service) GetDiscountRevenue(ctx context.Context, tenantID uuid.UUID, period string) (*RevenueImpact, error) {
	// TODO: Implement revenue impact analysis
	return nil, fmt.Errorf("TODO: implement GetDiscountRevenue")
}