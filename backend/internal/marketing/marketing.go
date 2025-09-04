package marketing

import (
	"time"

	"github.com/google/uuid"
)

// CampaignType represents the type of marketing campaign
type CampaignType string

// CampaignStatus represents the status of a campaign
type CampaignStatus string

// SegmentType represents customer segment criteria
type SegmentType string

// EmailStatus represents the status of an email
type EmailStatus string

const (
	CampaignEmail       CampaignType = "email"
	CampaignSMS         CampaignType = "sms"
	CampaignPushNotif   CampaignType = "push_notification"
	CampaignAbandonedCart CampaignType = "abandoned_cart"
	CampaignWelcome     CampaignType = "welcome"
)

const (
	StatusDraft     CampaignStatus = "draft"
	StatusScheduled CampaignStatus = "scheduled"
	StatusRunning   CampaignStatus = "running"
	StatusPaused    CampaignStatus = "paused"
	StatusCompleted CampaignStatus = "completed"
	StatusCancelled CampaignStatus = "cancelled"
)

const (
	SegmentAll          SegmentType = "all_customers"
	SegmentNewCustomers SegmentType = "new_customers"
	SegmentReturning    SegmentType = "returning_customers"
	SegmentHighValue    SegmentType = "high_value"
	SegmentInactive     SegmentType = "inactive"
	SegmentCustom       SegmentType = "custom"
)

const (
	EmailPending    EmailStatus = "pending"
	EmailSent       EmailStatus = "sent"
	EmailDelivered  EmailStatus = "delivered"
	EmailOpened     EmailStatus = "opened"
	EmailClicked    EmailStatus = "clicked"
	EmailBounced    EmailStatus = "bounced"
	EmailFailed     EmailStatus = "failed"
	EmailUnsubscribed EmailStatus = "unsubscribed"
)

// Campaign represents a marketing campaign
type Campaign struct {
	ID       uuid.UUID      `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID      `json:"tenant_id" gorm:"not null;index"`
	
	// Campaign details
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Type        CampaignType   `json:"type" gorm:"not null"`
	Status      CampaignStatus `json:"status" gorm:"default:draft"`
	
	// Content
	Subject     string `json:"subject,omitempty"`
	Content     string `json:"content" gorm:"type:text"`
	PreviewText string `json:"preview_text,omitempty"`
	
	// Targeting
	SegmentType    SegmentType `json:"segment_type" gorm:"default:all_customers"`
	SegmentRules   string      `json:"segment_rules,omitempty" gorm:"type:json"` // JSON criteria
	TargetCount    int         `json:"target_count" gorm:"default:0"`
	
	// Scheduling
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	// Settings
	FromName     string `json:"from_name,omitempty"`
	FromEmail    string `json:"from_email,omitempty"`
	ReplyToEmail string `json:"reply_to_email,omitempty"`
	
	// Analytics
	SentCount      int `json:"sent_count" gorm:"default:0"`
	DeliveredCount int `json:"delivered_count" gorm:"default:0"`
	OpenedCount    int `json:"opened_count" gorm:"default:0"`
	ClickedCount   int `json:"clicked_count" gorm:"default:0"`
	BouncedCount   int `json:"bounced_count" gorm:"default:0"`
	UnsubscribedCount int `json:"unsubscribed_count" gorm:"default:0"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Emails []CampaignEmail `json:"emails,omitempty" gorm:"foreignKey:CampaignID"`
}

// CampaignEmail represents an individual email sent as part of a campaign
type CampaignEmail struct {
	ID         uuid.UUID   `json:"id" gorm:"primarykey"`
	CampaignID uuid.UUID   `json:"campaign_id" gorm:"not null;index"`
	CustomerID *uuid.UUID  `json:"customer_id,omitempty" gorm:"index"`
	
	// Recipient details
	RecipientEmail string `json:"recipient_email" gorm:"not null"`
	RecipientName  string `json:"recipient_name,omitempty"`
	
	// Email content (may be personalized)
	Subject     string `json:"subject"`
	Content     string `json:"content" gorm:"type:text"`
	PreviewText string `json:"preview_text,omitempty"`
	
	// Status tracking
	Status      EmailStatus `json:"status" gorm:"default:pending"`
	SentAt      *time.Time  `json:"sent_at,omitempty"`
	DeliveredAt *time.Time  `json:"delivered_at,omitempty"`
	OpenedAt    *time.Time  `json:"opened_at,omitempty"`
	ClickedAt   *time.Time  `json:"clicked_at,omitempty"`
	BouncedAt   *time.Time  `json:"bounced_at,omitempty"`
	
	// Tracking
	OpenCount  int `json:"open_count" gorm:"default:0"`
	ClickCount int `json:"click_count" gorm:"default:0"`
	
	// External provider info
	ExternalID  string `json:"external_id,omitempty"`
	ProviderLog string `json:"provider_log,omitempty" gorm:"type:text"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CustomerSegment represents a customer segment for targeting
type CustomerSegment struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Segment details
	Name        string      `json:"name" gorm:"not null"`
	Description string      `json:"description,omitempty"`
	Type        SegmentType `json:"type" gorm:"not null"`
	
	// Segment criteria (JSON)
	Rules       string `json:"rules" gorm:"type:json"` // Complex filtering rules
	CustomerCount int  `json:"customer_count" gorm:"default:0"`
	
	// Settings
	IsActive    bool `json:"is_active" gorm:"default:true"`
	AutoUpdate  bool `json:"auto_update" gorm:"default:true"`
	LastUpdated *time.Time `json:"last_updated,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EmailTemplate represents reusable email templates
type EmailTemplate struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Template details
	Name        string          `json:"name" gorm:"not null"`
	Description string          `json:"description,omitempty"`
	Category    string          `json:"category,omitempty"`
	Type        CampaignType    `json:"type" gorm:"default:email"`
	
	// Content
	Subject     string `json:"subject" gorm:"not null"`
	Content     string `json:"content" gorm:"type:text;not null"`
	PreviewText string `json:"preview_text,omitempty"`
	
	// Design
	DesignJSON string `json:"design_json,omitempty" gorm:"type:json"` // Template builder data
	
	// Settings
	IsActive      bool   `json:"is_active" gorm:"default:true"`
	FromName      string `json:"from_name,omitempty"`
	FromEmail     string `json:"from_email,omitempty"`
	ReplyToEmail  string `json:"reply_to_email,omitempty"`
	
	// Usage stats
	UsageCount int `json:"usage_count" gorm:"default:0"`
	LastUsed   *time.Time `json:"last_used,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AbandonedCart represents abandoned cart recovery campaigns
type AbandonedCart struct {
	ID         uuid.UUID  `json:"id" gorm:"primarykey"`
	TenantID   uuid.UUID  `json:"tenant_id" gorm:"not null;index"`
	CustomerID *uuid.UUID `json:"customer_id,omitempty" gorm:"index"`
	
	// Cart details
	CartID        uuid.UUID `json:"cart_id" gorm:"not null;index"`
	CustomerEmail string    `json:"customer_email" gorm:"not null"`
	CustomerName  string    `json:"customer_name,omitempty"`
	CartValue     float64   `json:"cart_value" gorm:"not null"`
	ItemCount     int       `json:"item_count" gorm:"not null"`
	
	// Cart items (simplified JSON)
	Items string `json:"items" gorm:"type:json"`
	
	// Recovery status
	IsRecovered     bool       `json:"is_recovered" gorm:"default:false"`
	RecoveredAt     *time.Time `json:"recovered_at,omitempty"`
	RecoveredValue  float64    `json:"recovered_value" gorm:"default:0"`
	
	// Email campaign tracking
	EmailsSent      int        `json:"emails_sent" gorm:"default:0"`
	LastEmailSent   *time.Time `json:"last_email_sent,omitempty"`
	EmailsOpened    int        `json:"emails_opened" gorm:"default:0"`
	EmailsClicked   int        `json:"emails_clicked" gorm:"default:0"`
	
	// Timestamps
	AbandonedAt time.Time `json:"abandoned_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewsletterSubscriber represents newsletter subscriptions
type NewsletterSubscriber struct {
	ID       uuid.UUID  `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID  `json:"tenant_id" gorm:"not null;index"`
	CustomerID *uuid.UUID `json:"customer_id,omitempty" gorm:"index"`
	
	// Subscriber details
	Email        string  `json:"email" gorm:"not null;uniqueIndex:idx_tenant_email"`
	Name         string  `json:"name,omitempty"`
	Status       string  `json:"status" gorm:"default:active"` // active, unsubscribed, bounced
	
	// Preferences
	Preferences  string  `json:"preferences,omitempty" gorm:"type:json"`
	Tags         []string `json:"tags,omitempty" gorm:"serializer:json"`
	
	// Tracking
	SubscribedAt   time.Time  `json:"subscribed_at"`
	UnsubscribedAt *time.Time `json:"unsubscribed_at,omitempty"`
	ConfirmedAt    *time.Time `json:"confirmed_at,omitempty"`
	
	// Source tracking
	SourceURL      string `json:"source_url,omitempty"`
	SourceCampaign string `json:"source_campaign,omitempty"`
	IPAddress      string `json:"ip_address,omitempty"`
	UserAgent      string `json:"user_agent,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MarketingSettings represents marketing module settings
type MarketingSettings struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"unique;not null"`
	
	// Email settings
	FromName              string `json:"from_name" gorm:"not null"`
	FromEmail             string `json:"from_email" gorm:"not null"`
	ReplyToEmail          string `json:"reply_to_email,omitempty"`
	UnsubscribeURL        string `json:"unsubscribe_url,omitempty"`
	
	// SMTP/Email provider settings
	EmailProvider         string `json:"email_provider" gorm:"default:smtp"` // smtp, sendgrid, mailgun, etc.
	SMTPHost             string `json:"smtp_host,omitempty"`
	SMTPPort             int    `json:"smtp_port,omitempty"`
	SMTPUsername         string `json:"smtp_username,omitempty"`
	SMTPPassword         string `json:"smtp_password,omitempty"`
	SendGridAPIKey       string `json:"sendgrid_api_key,omitempty"`
	MailgunAPIKey        string `json:"mailgun_api_key,omitempty"`
	MailgunDomain        string `json:"mailgun_domain,omitempty"`
	
	// Abandoned cart settings
	AbandonedCartEnabled  bool `json:"abandoned_cart_enabled" gorm:"default:true"`
	AbandonedCartDelay    int  `json:"abandoned_cart_delay" gorm:"default:60"`  // minutes
	AbandonedCartTemplate uuid.UUID `json:"abandoned_cart_template,omitempty"`
	
	// Welcome email settings
	WelcomeEmailEnabled   bool       `json:"welcome_email_enabled" gorm:"default:true"`
	WelcomeEmailDelay     int        `json:"welcome_email_delay" gorm:"default:0"` // minutes
	WelcomeEmailTemplate  uuid.UUID  `json:"welcome_email_template,omitempty"`
	
	// General settings
	TrackingEnabled       bool   `json:"tracking_enabled" gorm:"default:true"`
	DoubleOptIn           bool   `json:"double_opt_in" gorm:"default:false"`
	UnsubscribeFooter     string `json:"unsubscribe_footer,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}