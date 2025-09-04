package support

import (
	"time"

	"github.com/google/uuid"
)

// TicketStatus represents the status of a support ticket
type TicketStatus string

// TicketPriority represents the priority of a support ticket
type TicketPriority string

// TicketCategory represents the category of a support ticket
type TicketCategory string

const (
	StatusOpen       TicketStatus = "open"
	StatusInProgress TicketStatus = "in_progress"
	StatusWaiting    TicketStatus = "waiting"
	StatusResolved   TicketStatus = "resolved"
	StatusClosed     TicketStatus = "closed"
)

const (
	PriorityLow      TicketPriority = "low"
	PriorityMedium   TicketPriority = "medium"
	PriorityHigh     TicketPriority = "high"
	PriorityCritical TicketPriority = "critical"
)

const (
	CategoryGeneral    TicketCategory = "general"
	CategoryTechnical  TicketCategory = "technical"
	CategoryBilling    TicketCategory = "billing"
	CategoryBugReport  TicketCategory = "bug_report"
	CategoryFeature    TicketCategory = "feature_request"
)

// Ticket represents a support ticket
type Ticket struct {
	ID       uuid.UUID      `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID      `json:"tenant_id" gorm:"not null;index"`
	UserID   uuid.UUID      `json:"user_id" gorm:"not null;index"`
	
	// Ticket details
	Subject     string         `json:"subject" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Status      TicketStatus   `json:"status" gorm:"default:open"`
	Priority    TicketPriority `json:"priority" gorm:"default:medium"`
	Category    TicketCategory `json:"category" gorm:"default:general"`
	
	// Assignment
	AssignedToID *uuid.UUID `json:"assigned_to_id,omitempty" gorm:"index"`
	
	// Customer info
	CustomerEmail string `json:"customer_email" gorm:"not null"`
	CustomerName  string `json:"customer_name,omitempty"`
	
	// Metadata
	Tags        []string `json:"tags,omitempty" gorm:"serializer:json"`
	Attachments []string `json:"attachments,omitempty" gorm:"serializer:json"`
	
	// Timestamps
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	
	// Relations
	Messages []TicketMessage `json:"messages,omitempty" gorm:"foreignKey:TicketID"`
}

// TicketMessage represents a message in a support ticket
type TicketMessage struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TicketID uuid.UUID `json:"ticket_id" gorm:"not null;index"`
	UserID   *uuid.UUID `json:"user_id,omitempty" gorm:"index"` // nil for customer messages
	
	// Message details
	Content     string   `json:"content" gorm:"type:text;not null"`
	IsInternal  bool     `json:"is_internal" gorm:"default:false"`
	Attachments []string `json:"attachments,omitempty" gorm:"serializer:json"`
	
	// Sender info
	SenderName  string `json:"sender_name" gorm:"not null"`
	SenderEmail string `json:"sender_email" gorm:"not null"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FAQ represents frequently asked questions
type FAQ struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// FAQ details
	Question    string   `json:"question" gorm:"not null"`
	Answer      string   `json:"answer" gorm:"type:text;not null"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty" gorm:"serializer:json"`
	IsPublished bool     `json:"is_published" gorm:"default:true"`
	
	// Metadata
	ViewCount int `json:"view_count" gorm:"default:0"`
	Order     int `json:"order" gorm:"default:0"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// KnowledgeBase represents knowledge base articles
type KnowledgeBase struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Article details
	Title       string   `json:"title" gorm:"not null"`
	Content     string   `json:"content" gorm:"type:text;not null"`
	Excerpt     string   `json:"excerpt,omitempty"`
	Slug        string   `json:"slug" gorm:"unique;not null"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty" gorm:"serializer:json"`
	IsPublished bool     `json:"is_published" gorm:"default:true"`
	
	// SEO
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
	
	// Metadata
	ViewCount int `json:"view_count" gorm:"default:0"`
	Order     int `json:"order" gorm:"default:0"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SupportSettings represents support module settings
type SupportSettings struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"unique;not null"`
	
	// General settings
	SupportEmail       string `json:"support_email" gorm:"not null"`
	AutoReplyEnabled   bool   `json:"auto_reply_enabled" gorm:"default:true"`
	AutoReplyMessage   string `json:"auto_reply_message,omitempty"`
	BusinessHours      string `json:"business_hours,omitempty"`
	
	// Notification settings
	EmailNotifications bool `json:"email_notifications" gorm:"default:true"`
	SlackWebhookURL    string `json:"slack_webhook_url,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}