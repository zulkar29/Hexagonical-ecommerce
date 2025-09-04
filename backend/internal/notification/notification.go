package notification

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: Implement notification entities
// This will handle:
// - Email notifications
// - SMS notifications  
// - Push notifications
// - In-app notifications
// - Notification templates

type Notification struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID     uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	UserID       *uuid.UUID `json:"user_id" gorm:"type:uuid;index"` // Nullable for system notifications
	Type         string    `json:"type" gorm:"size:50;not null"`   // email, sms, push, in_app
	Channel      string    `json:"channel" gorm:"size:50;not null"` // order_confirmation, password_reset, marketing, etc.
	Subject      string    `json:"subject" gorm:"size:255"`
	Content      string    `json:"content" gorm:"type:text;not null"`
	Recipient    string    `json:"recipient" gorm:"size:255;not null"` // email address, phone number, device token
	Status       string    `json:"status" gorm:"size:50;not null;default:'pending'"` // pending, sent, delivered, failed, read
	Priority     string    `json:"priority" gorm:"size:20;not null;default:'normal'"` // low, normal, high, urgent
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
