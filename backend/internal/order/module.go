package order

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	ValidateDiscountCode(ctx context.Context, tenantID uuid.UUID, code string, customerID *uuid.UUID, customerEmail string, orderAmount float64, itemQuantity int, productIDs []string, categoryIDs []string) (*DiscountValidation, error)
	ApplyDiscount(ctx context.Context, tenantID uuid.UUID, code string, orderID uuid.UUID, customerID *uuid.UUID, customerEmail string, orderAmount float64, itemQuantity int, productIDs []string, categoryIDs []string, ipAddress string, userAgent string) (*DiscountApplication, error)
	RemoveDiscount(ctx context.Context, tenantID uuid.UUID, orderID uuid.UUID) error
}

// PaymentService interface for payment operations
type PaymentService interface {
	CreatePayment(ctx context.Context, tenantID uuid.UUID, orderID string, amount float64, currency string, gateway string, paymentMethodID string, customerEmail string, customerPhone string, returnURL string) (*CreatePaymentResponse, error)
	ProcessPayment(ctx context.Context, tenantID uuid.UUID, paymentID string, gateway string, gatewayResponse map[string]interface{}) error
	RefundPayment(ctx context.Context, tenantID uuid.UUID, paymentID string, amount float64, reason string) error
}

// InventoryService interface for inventory management
type InventoryService interface {
	ReserveStock(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID, quantity int) error
	RestoreStock(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID, quantity int) error
	UpdateInventory(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID, quantity int) error
}

// NotificationService interface for notification operations
type NotificationService interface {
	SendNotification(ctx context.Context, tenantID uuid.UUID, notificationType string, channel string, recipients []string, subject string, content string, userID string, priority string, variables map[string]interface{}, templateID string, scheduledAt *time.Time) (*SendNotificationResponse, error)
	SendEmail(ctx context.Context, tenantID uuid.UUID, to []string, subject string, content string, contentType string, variables map[string]interface{}, templateID string) error
	SendSMS(ctx context.Context, tenantID uuid.UUID, to []string, message string, variables map[string]interface{}, templateID string) error
}

// Notification related structs
type SendNotificationRequest struct {
	Type        string                 `json:"type" validate:"required,oneof=email sms push in_app"`
	Channel     string                 `json:"channel" validate:"required"`
	Recipients  []string               `json:"recipients" validate:"required,min=1"`
	Subject     string                 `json:"subject,omitempty"`
	Content     string                 `json:"content" validate:"required"`
	UserID      string                 `json:"user_id,omitempty"`
	Priority    string                 `json:"priority,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	TemplateID  string                 `json:"template_id,omitempty"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
}

type SendNotificationResponse struct {
	NotificationIDs []string `json:"notification_ids"`
	Status          string   `json:"status"`
	Message         string   `json:"message"`
}

type SendEmailRequest struct {
	To          []string               `json:"to" validate:"required,min=1"`
	Subject     string                 `json:"subject" validate:"required"`
	Content     string                 `json:"content" validate:"required"`
	ContentType string                 `json:"content_type,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	TemplateID  string                 `json:"template_id,omitempty"`
}

type SendSMSRequest struct {
	To         []string               `json:"to" validate:"required,min=1"`
	Message    string                 `json:"message" validate:"required,max=160"`
	Variables  map[string]interface{} `json:"variables,omitempty"`
	TemplateID string                 `json:"template_id,omitempty"`
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

// Module represents the order module
type Module struct {
	Handler    *Handler
	Service    *Service
	Repository Repository
}

// NewModule creates a new order module with all dependencies
func NewModule(db *gorm.DB, productService ProductService, discountService DiscountService, paymentService PaymentService, inventoryService InventoryService, notificationService NotificationService) *Module {
	repository := NewRepository(db)
	service := NewService(repository, db, productService, discountService, paymentService, inventoryService, notificationService)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}

// RegisterRoutes registers all order routes with the router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.Handler.RegisterRoutes(router)
}

// Migrate runs database migrations for order module
func (m *Module) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Order{},
		&OrderItem{},
	)
}

// GetService returns the order service for integration with other modules
func (m *Module) GetService() *Service {
	return m.Service
}

// GetRepository returns the order repository for direct database access
func (m *Module) GetRepository() Repository {
	return m.Repository
}
