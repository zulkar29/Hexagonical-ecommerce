package notification

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification types
const (
	TypeEmail = "email"
	TypeSMS   = "sms"
	TypePush  = "push"
	TypeInApp = "in_app"
)

// Notification channels
const (
	ChannelOrderConfirmation = "order_confirmation"
	ChannelPasswordReset     = "password_reset"
	ChannelPaymentSuccess    = "payment_success"
	ChannelPaymentFailed     = "payment_failed"
	ChannelWelcome           = "welcome"
	ChannelMarketing         = "marketing"
	ChannelAbandonedCart     = "abandoned_cart"
	ChannelShippingUpdate    = "shipping_update"
)

// Notification statuses
const (
	StatusPending   = "pending"
	StatusSent      = "sent"
	StatusDelivered = "delivered"
	StatusFailed    = "failed"
	StatusRead      = "read"
)

// Priority levels
const (
	PriorityLow    = "low"
	PriorityNormal = "normal"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"
)

type Notification struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID     uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	UserID       *uuid.UUID `json:"user_id" gorm:"type:uuid;index"` // Nullable for system notifications
	Type         string    `json:"type" gorm:"size:50;not null"`   // email, sms, push, in_app
	Channel      string    `json:"channel" gorm:"size:50;not null"` // order_confirmation, password_reset, marketing, etc.
	Subject      string    `json:"subject" gorm:"size:255"`
	Content      string    `json:"content" gorm:"type:text;not null"`
	Recipient    string    `json:"recipient" gorm:"size:255;not null"` // email address, phone number, device token
	Status       string    `json:"status" gorm:"size:50;not null;default:'pending'"`
	Priority     string    `json:"priority" gorm:"size:20;not null;default:'normal'"`
	ScheduledAt  *time.Time `json:"scheduled_at"` // For scheduled notifications
	SentAt       *time.Time `json:"sent_at"`
	DeliveredAt  *time.Time `json:"delivered_at"`
	ReadAt       *time.Time `json:"read_at"`
	FailedAt     *time.Time `json:"failed_at"`
	FailureReason string    `json:"failure_reason" gorm:"type:text"`
	RetryCount   int       `json:"retry_count" gorm:"default:0"`
	MaxRetries   int       `json:"max_retries" gorm:"default:3"`
	Metadata     string    `json:"metadata" gorm:"type:json"` // Additional data like action buttons, tracking info
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type NotificationTemplate struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Type        string    `json:"type" gorm:"size:50;not null"`        // email, sms, push, in_app
	Channel     string    `json:"channel" gorm:"size:50;not null"`     // order_confirmation, password_reset, etc.
	Subject     string    `json:"subject" gorm:"size:255"`             // For email templates
	Content     string    `json:"content" gorm:"type:text;not null"`   // Template content with placeholders
	Variables   string    `json:"variables" gorm:"type:json"`          // List of available variables
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	IsDefault   bool      `json:"is_default" gorm:"default:false"`     // System default template
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type NotificationPreference struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID     uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Channel      string    `json:"channel" gorm:"size:50;not null"`        // order_updates, marketing, etc.
	EmailEnabled bool      `json:"email_enabled" gorm:"default:true"`
	SMSEnabled   bool      `json:"sms_enabled" gorm:"default:false"`
	PushEnabled  bool      `json:"push_enabled" gorm:"default:true"`
	InAppEnabled bool      `json:"in_app_enabled" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type NotificationLog struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NotificationID uuid.UUID `json:"notification_id" gorm:"type:uuid;not null;index"`
	Event          string    `json:"event" gorm:"size:50;not null"`  // sent, delivered, opened, clicked, failed
	EventData      string    `json:"event_data" gorm:"type:json"`    // Additional event details
	CreatedAt      time.Time `json:"created_at"`
}

// Email/SMS Gateway Configurations
type EmailProvider struct {
	Name       string `json:"name"`       // sendgrid, mailgun, ses
	APIKey     string `json:"api_key"`
	APISecret  string `json:"api_secret,omitempty"`
	FromEmail  string `json:"from_email"`
	FromName   string `json:"from_name"`
	ReplyToEmail string `json:"reply_to_email,omitempty"`
}

type SMSProvider struct {
	Name       string `json:"name"`       // twilio, nexmo, local_sms_bd
	APIKey     string `json:"api_key"`
	APISecret  string `json:"api_secret"`
	FromNumber string `json:"from_number"`
	WebhookURL string `json:"webhook_url,omitempty"`
}

// Bangladesh SMS Gateway specific (for local providers)
type BDSMSGatewayRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Number   string `json:"number"`
	Message  string `json:"message"`
	Type     string `json:"type"` // text, unicode
}

type BDSMSGatewayResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	SMSID   string `json:"sms_id"`
	Balance string `json:"balance"`
}

// Request/Response Structures
type SendNotificationRequest struct {
	Type        string                 `json:"type" validate:"required,oneof=email sms push in_app"`
	Channel     string                 `json:"channel" validate:"required"`
	Recipients  []string               `json:"recipients" validate:"required,min=1"`
	Subject     string                 `json:"subject,omitempty"`
	Content     string                 `json:"content,omitempty"`
	TemplateID  string                 `json:"template_id,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	Priority    string                 `json:"priority,omitempty" validate:"omitempty,oneof=low normal high urgent"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
}

type SendNotificationResponse struct {
	NotificationIDs []string `json:"notification_ids"`
	Status          string   `json:"status"`
	Message         string   `json:"message"`
}

type SendEmailRequest struct {
	To        []string               `json:"to" validate:"required,min=1"`
	Subject   string                 `json:"subject" validate:"required"`
	Content   string                 `json:"content" validate:"required"`
	ContentType string               `json:"content_type,omitempty"` // html, text
	Variables map[string]interface{} `json:"variables,omitempty"`
	TemplateID string                `json:"template_id,omitempty"`
}

type SendSMSRequest struct {
	To        []string               `json:"to" validate:"required,min=1"`
	Message   string                 `json:"message" validate:"required,max=160"`
	Variables map[string]interface{} `json:"variables,omitempty"`
	TemplateID string                `json:"template_id,omitempty"`
}

type NotificationStatsResponse struct {
	TotalSent      int64   `json:"total_sent"`
	TotalDelivered int64   `json:"total_delivered"`
	TotalFailed    int64   `json:"total_failed"`
	DeliveryRate   float64 `json:"delivery_rate"`
	FailureRate    float64 `json:"failure_rate"`
	EmailStats     struct {
		Sent      int64 `json:"sent"`
		Delivered int64 `json:"delivered"`
		Opened    int64 `json:"opened"`
		Clicked   int64 `json:"clicked"`
	} `json:"email_stats"`
	SMSStats struct {
		Sent      int64 `json:"sent"`
		Delivered int64 `json:"delivered"`
		Failed    int64 `json:"failed"`
	} `json:"sms_stats"`
}

type CreateTemplateRequest struct {
	Name      string                 `json:"name" validate:"required"`
	Type      string                 `json:"type" validate:"required,oneof=email sms push in_app"`
	Channel   string                 `json:"channel" validate:"required"`
	Subject   string                 `json:"subject,omitempty"`
	Content   string                 `json:"content" validate:"required"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type UpdateTemplateRequest struct {
	Name      string                 `json:"name,omitempty"`
	Subject   string                 `json:"subject,omitempty"`
	Content   string                 `json:"content,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
	IsActive  *bool                  `json:"is_active,omitempty"`
}

type NotificationPreferenceRequest struct {
	Channel      string `json:"channel" validate:"required"`
	EmailEnabled *bool  `json:"email_enabled,omitempty"`
	SMSEnabled   *bool  `json:"sms_enabled,omitempty"`
	PushEnabled  *bool  `json:"push_enabled,omitempty"`
	InAppEnabled *bool  `json:"in_app_enabled,omitempty"`
}
