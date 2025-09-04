package webhook

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: Implement webhook entities
// This will handle:
// - Webhook endpoint management
// - Event delivery and retry logic
// - Third-party integrations (payment gateways, shipping providers)
// - API rate limiting and security

type WebhookEvent string
type WebhookStatus string
type WebhookProvider string

const (
	// E-commerce Events
	EventOrderCreated    WebhookEvent = "order.created"
	EventOrderUpdated    WebhookEvent = "order.updated" 
	EventOrderCancelled  WebhookEvent = "order.cancelled"
	EventOrderFulfilled  WebhookEvent = "order.fulfilled"
	
	// Payment Events
	EventPaymentCreated  WebhookEvent = "payment.created"
	EventPaymentSucceded WebhookEvent = "payment.succeeded"
	EventPaymentFailed   WebhookEvent = "payment.failed"
	EventPaymentRefunded WebhookEvent = "payment.refunded"
	
	// Product Events
	EventProductCreated  WebhookEvent = "product.created"
	EventProductUpdated  WebhookEvent = "product.updated"
	EventProductDeleted  WebhookEvent = "product.deleted"
	EventInventoryLow    WebhookEvent = "inventory.low"
	
	// Customer Events
	EventCustomerCreated WebhookEvent = "customer.created"
	EventCustomerUpdated WebhookEvent = "customer.updated"
	
	// Shipping Events
	EventShipmentCreated WebhookEvent = "shipment.created"
	EventShipmentShipped WebhookEvent = "shipment.shipped"
	EventShipmentDelivered WebhookEvent = "shipment.delivered"
)

const (
	StatusPending    WebhookStatus = "pending"
	StatusDelivering WebhookStatus = "delivering"
	StatusDelivered  WebhookStatus = "delivered"
	StatusFailed     WebhookStatus = "failed"
	StatusDisabled   WebhookStatus = "disabled"
)

const (
	// Payment Providers
	ProviderStripe   WebhookProvider = "stripe"
	ProviderPayPal   WebhookProvider = "paypal"
	ProviderBkash    WebhookProvider = "bkash"
	ProviderNagad    WebhookProvider = "nagad"
	ProviderRocket   WebhookProvider = "rocket"
	
	// Shipping Providers
	ProviderPathao   WebhookProvider = "pathao"
	ProviderRedX     WebhookProvider = "redx"
	ProviderPaperfly WebhookProvider = "paperfly"
	ProviderDHL      WebhookProvider = "dhl"
	ProviderFedEx    WebhookProvider = "fedex"
	
	// Integration Providers
	ProviderShopify     WebhookProvider = "shopify"
	ProviderWooCommerce WebhookProvider = "woocommerce"
	ProviderMailchimp   WebhookProvider = "mailchimp"
	ProviderGoogle      WebhookProvider = "google_analytics"
	ProviderFacebook    WebhookProvider = "facebook_pixel"
)

type WebhookEndpoint struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID        uuid.UUID       `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name            string          `json:"name" gorm:"size:100;not null"`
	URL             string          `json:"url" gorm:"size:500;not null"`
	Description     string          `json:"description" gorm:"type:text"`
	Events          []WebhookEvent  `json:"events" gorm:"serializer:json"`
	IsActive        bool            `json:"is_active" gorm:"default:true"`
	Secret          string          `json:"secret" gorm:"size:255"` // For signature verification
	RetryPolicy     string          `json:"retry_policy" gorm:"size:50;default:'exponential'"` // exponential, linear, none
	MaxRetries      int             `json:"max_retries" gorm:"default:3"`
	TimeoutSeconds  int             `json:"timeout_seconds" gorm:"default:30"`
	Headers         map[string]string `json:"headers" gorm:"serializer:json"`
	LastDeliveryAt  *time.Time      `json:"last_delivery_at"`
	LastStatus      WebhookStatus   `json:"last_status" gorm:"default:'pending'"`
	FailureCount    int             `json:"failure_count" gorm:"default:0"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `json:"deleted_at" gorm:"index"`
}

type WebhookDelivery struct {
	ID             uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID       uuid.UUID       `json:"tenant_id" gorm:"type:uuid;not null;index"`
	EndpointID     uuid.UUID       `json:"endpoint_id" gorm:"type:uuid;not null;index"`
	Event          WebhookEvent    `json:"event" gorm:"size:50;not null"`
	EventID        uuid.UUID       `json:"event_id" gorm:"type:uuid;not null"` // Reference to source event
	Status         WebhookStatus   `json:"status" gorm:"default:'pending'"`
	AttemptCount   int             `json:"attempt_count" gorm:"default:0"`
	MaxAttempts    int             `json:"max_attempts" gorm:"default:3"`
	
	// Request details
	RequestURL     string          `json:"request_url" gorm:"size:500;not null"`
	RequestHeaders map[string]string `json:"request_headers" gorm:"serializer:json"`
	RequestBody    string          `json:"request_body" gorm:"type:text"`
	
	// Response details
	ResponseStatus   int            `json:"response_status"`
	ResponseHeaders  map[string]string `json:"response_headers" gorm:"serializer:json"`
	ResponseBody     string         `json:"response_body" gorm:"type:text"`
	ResponseTime     int            `json:"response_time"` // milliseconds
	
	// Delivery tracking
	NextRetryAt    *time.Time      `json:"next_retry_at"`
	DeliveredAt    *time.Time      `json:"delivered_at"`
	FailedAt       *time.Time      `json:"failed_at"`
	ErrorMessage   string          `json:"error_message" gorm:"type:text"`
	
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type WebhookIncoming struct {
	ID             uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID       uuid.UUID       `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Provider       WebhookProvider `json:"provider" gorm:"size:50;not null"`
	Event          string          `json:"event" gorm:"size:100;not null"`
	EventID        string          `json:"event_id" gorm:"size:255"` // External event ID
	Signature      string          `json:"signature" gorm:"size:500"` // Webhook signature
	IsVerified     bool            `json:"is_verified" gorm:"default:false"`
	IsProcessed    bool            `json:"is_processed" gorm:"default:false"`
	
	// Request data
	Headers        map[string]string `json:"headers" gorm:"serializer:json"`
	Body           string          `json:"body" gorm:"type:text"`
	IPAddress      string          `json:"ip_address" gorm:"size:45"`
	UserAgent      string          `json:"user_agent" gorm:"size:500"`
	
	// Processing results
	ProcessedAt    *time.Time      `json:"processed_at"`
	ProcessingError string         `json:"processing_error" gorm:"type:text"`
	
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type WebhookRateLimit struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID     uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	EndpointID   uuid.UUID `json:"endpoint_id" gorm:"type:uuid;not null;index"`
	WindowStart  time.Time `json:"window_start" gorm:"not null"`
	WindowEnd    time.Time `json:"window_end" gorm:"not null"`
	RequestCount int       `json:"request_count" gorm:"default:0"`
	Limit        int       `json:"limit" gorm:"default:100"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Business Logic Methods

// IsDelivered checks if webhook was successfully delivered
func (wd *WebhookDelivery) IsDelivered() bool {
	return wd.Status == StatusDelivered && wd.DeliveredAt != nil
}

// IsFailed checks if webhook delivery failed permanently
func (wd *WebhookDelivery) IsFailed() bool {
	return wd.Status == StatusFailed || wd.AttemptCount >= wd.MaxAttempts
}

// ShouldRetry checks if webhook should be retried
func (wd *WebhookDelivery) ShouldRetry() bool {
	return wd.Status == StatusFailed && 
		   wd.AttemptCount < wd.MaxAttempts && 
		   wd.NextRetryAt != nil && 
		   wd.NextRetryAt.Before(time.Now())
}

// GetNextRetryDelay calculates next retry delay using exponential backoff
func (wd *WebhookDelivery) GetNextRetryDelay() time.Duration {
	baseDelay := 1 * time.Minute
	multiplier := time.Duration(wd.AttemptCount * wd.AttemptCount) // Exponential backoff
	return baseDelay * multiplier
}

// IsActive checks if webhook endpoint is active and healthy
func (we *WebhookEndpoint) IsActive() bool {
	return we.IsActive && we.FailureCount < 10 // Disable after 10 consecutive failures
}

// SupportsEvent checks if endpoint supports a specific event
func (we *WebhookEndpoint) SupportsEvent(event WebhookEvent) bool {
	for _, supportedEvent := range we.Events {
		if supportedEvent == event {
			return true
		}
	}
	return false
}

// IsRateLimited checks if endpoint is rate limited
func (wrl *WebhookRateLimit) IsRateLimited() bool {
	return wrl.RequestCount >= wrl.Limit && time.Now().Before(wrl.WindowEnd)
}

// TODO: Add more business logic methods
// - GenerateSignature(payload, secret) string
// - VerifySignature(payload, signature, secret) bool
// - CalculateRetryDelay(attemptCount int) time.Duration
// - ShouldDisableEndpoint(failureCount int) bool
