package contact

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ContactStatus represents the status of a contact message
type ContactStatus string

// ContactPriority represents the priority level of a contact message
type ContactPriority string

// ContactType represents different types of contact messages
type ContactType string

const (
	StatusNew        ContactStatus = "new"
	StatusRead       ContactStatus = "read"
	StatusInProgress ContactStatus = "in_progress"
	StatusResolved   ContactStatus = "resolved"
	StatusClosed     ContactStatus = "closed"
	StatusSpam       ContactStatus = "spam"
)

const (
	PriorityLow      ContactPriority = "low"
	PriorityMedium   ContactPriority = "medium"
	PriorityHigh     ContactPriority = "high"
	PriorityCritical ContactPriority = "critical"
)

const (
	TypeGeneral      ContactType = "general"
	TypeSupport      ContactType = "support"
	TypeSales        ContactType = "sales"
	TypeTechnical    ContactType = "technical"
	TypeComplaint    ContactType = "complaint"
	TypeFeatureReq   ContactType = "feature_request"
	TypeBugReport    ContactType = "bug_report"
	TypePartnership  ContactType = "partnership"
	TypeMedia        ContactType = "media"
	TypeOther        ContactType = "other"
)

// ContactFormType represents the type of contact form
type ContactFormType string

const (
	FormTypeGeneral   ContactFormType = "general"
	FormTypeSupport   ContactFormType = "support"
	FormTypeSales     ContactFormType = "sales"
	FormTypeFeedback  ContactFormType = "feedback"
)

// ContactTemplateType represents the type of contact template
type ContactTemplateType string

const (
	TemplateTypeAutoReply   ContactTemplateType = "auto_reply"
	TemplateTypeFollowUp    ContactTemplateType = "follow_up"
	TemplateTypeResolution  ContactTemplateType = "resolution"
	TemplateTypeEscalation  ContactTemplateType = "escalation"
)

// AnalyticsPeriod represents the period for analytics
type AnalyticsPeriod string

const (
	AnalyticsPeriodDay     AnalyticsPeriod = "day"
	AnalyticsPeriodWeek    AnalyticsPeriod = "week"
	AnalyticsPeriodMonth   AnalyticsPeriod = "month"
	AnalyticsPeriodQuarter AnalyticsPeriod = "quarter"
	AnalyticsPeriodYear    AnalyticsPeriod = "year"
)

// ContactMetrics represents comprehensive contact metrics
type ContactMetrics struct {
	TenantID           uuid.UUID `json:"tenant_id"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	TotalContacts      int       `json:"total_contacts"`
	NewContacts        int       `json:"new_contacts"`
	ResolvedContacts   int       `json:"resolved_contacts"`
	AvgResponseTime    float64   `json:"avg_response_time"`
	AvgResolutionTime  float64   `json:"avg_resolution_time"`
	SatisfactionScore  float64   `json:"satisfaction_score"`
	TypeDistribution   map[ContactType]int `json:"type_distribution"`
	StatusDistribution map[ContactStatus]int `json:"status_distribution"`
}

// ContactFormFilter represents filter options for contact forms
type ContactFormFilter struct {
	Type     []ContactFormType `json:"type"`
	IsActive *bool             `json:"is_active"`
	IsPublic *bool             `json:"is_public"`
	Search   string            `json:"search"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	SortBy   string            `json:"sort_by"`
	SortDesc bool              `json:"sort_desc"`
}

// ContactTemplateFilter represents filter options for contact templates
type ContactTemplateFilter struct {
	Type     []ContactTemplateType `json:"type"`
	Category []string              `json:"category"`
	IsActive *bool                 `json:"is_active"`
	Search   string                `json:"search"`
	Limit    int                   `json:"limit"`
	Offset   int                   `json:"offset"`
	SortBy   string                `json:"sort_by"`
	SortDesc bool                  `json:"sort_desc"`
}

// Additional request types for handlers
type BulkUpdateContactsRequest struct {
	ContactIDs []uuid.UUID       `json:"contact_ids" validate:"required"`
	Status     *ContactStatus    `json:"status"`
	Priority   *ContactPriority  `json:"priority"`
	AssigneeID *uuid.UUID        `json:"assignee_id"`
	Tags       []string          `json:"tags"`
}

type ExportContactsRequest struct {
	Format    string        `json:"format" validate:"required,oneof=csv xlsx json"`
	Filter    ContactFilter `json:"filter"`
	Fields    []string      `json:"fields"`
	DateRange *struct {
		From time.Time `json:"from"`
		To   time.Time `json:"to"`
	} `json:"date_range"`
}

type UpdateContactStatusRequest struct {
	Status ContactStatus `json:"status" validate:"required"`
	Reason string        `json:"reason"`
}

type AssignContactRequest struct {
	AssigneeID uuid.UUID `json:"assignee_id" validate:"required"`
	Reason     string    `json:"reason"`
}

type UpdateContactPriorityRequest struct {
	Priority ContactPriority `json:"priority" validate:"required"`
	Reason   string          `json:"reason"`
}

type AddContactTagsRequest struct {
	Tags []string `json:"tags" validate:"required,min=1"`
}

type RemoveContactTagsRequest struct {
	Tags []string `json:"tags" validate:"required,min=1"`
}

type CreateContactReplyRequest struct {
	Content     string    `json:"content" validate:"required"`
	IsInternal  bool      `json:"is_internal"`
	AuthorID    uuid.UUID `json:"author_id" validate:"required"`
	Attachments []string  `json:"attachments"`
}

type SubmitContactFormRequest struct {
	Name        string            `json:"name" validate:"required"`
	Email       string            `json:"email" validate:"required,email"`
	Phone       string            `json:"phone"`
	Company     string            `json:"company"`
	Subject     string            `json:"subject" validate:"required"`
	Message     string            `json:"message" validate:"required"`
	CustomFields map[string]interface{} `json:"custom_fields"`
	Attachments []string          `json:"attachments"`
	ConsentGiven bool             `json:"consent_given"`
	CaptchaToken string           `json:"captcha_token"`
	IPAddress   string            `json:"ip_address"`
	UserAgent   string            `json:"user_agent"`
	ReferrerURL string            `json:"referrer_url"`
	PageURL     string            `json:"page_url"`
}

type CreateContactFormRequest struct {
	Name            string          `json:"name" validate:"required"`
	Title           string          `json:"title" validate:"required"`
	Description     string          `json:"description"`
	Fields          string          `json:"fields" validate:"required"`
	Settings        string          `json:"settings"`
	DefaultType     ContactType     `json:"default_type"`
	DefaultPriority ContactPriority `json:"default_priority"`
	DefaultAssignee *uuid.UUID      `json:"default_assignee"`
	RequireAuth     bool            `json:"require_auth"`
	AllowAttachments bool           `json:"allow_attachments"`
	MaxAttachments  int             `json:"max_attachments"`
	AutoReply       bool            `json:"auto_reply"`
	AutoReplySubject string         `json:"auto_reply_subject"`
	AutoReplyMessage string         `json:"auto_reply_message"`
	EnableCaptcha   bool            `json:"enable_captcha"`
	EnableRateLimit bool            `json:"enable_rate_limit"`
	RateLimitWindow int             `json:"rate_limit_window"`
	RateLimitRequests int           `json:"rate_limit_requests"`
	IsPublic        bool            `json:"is_public"`
}

type UpdateContactFormRequest struct {
	Name            *string         `json:"name"`
	Title           *string         `json:"title"`
	Description     *string         `json:"description"`
	Fields          *string         `json:"fields"`
	Settings        *string         `json:"settings"`
	DefaultType     *ContactType    `json:"default_type"`
	DefaultPriority *ContactPriority `json:"default_priority"`
	DefaultAssignee *uuid.UUID      `json:"default_assignee"`
	RequireAuth     *bool           `json:"require_auth"`
	AllowAttachments *bool          `json:"allow_attachments"`
	MaxAttachments  *int            `json:"max_attachments"`
	AutoReply       *bool           `json:"auto_reply"`
	AutoReplySubject *string        `json:"auto_reply_subject"`
	AutoReplyMessage *string        `json:"auto_reply_message"`
	EnableCaptcha   *bool           `json:"enable_captcha"`
	EnableRateLimit *bool           `json:"enable_rate_limit"`
	RateLimitWindow *int            `json:"rate_limit_window"`
	RateLimitRequests *int          `json:"rate_limit_requests"`
	IsActive        *bool           `json:"is_active"`
	IsPublic        *bool           `json:"is_public"`
}

type CreateContactTemplateRequest struct {
	Name        string              `json:"name" validate:"required"`
	Description string              `json:"description"`
	Type        ContactTemplateType `json:"type"`
	Category    string              `json:"category"`
	Subject     string              `json:"subject" validate:"required"`
	Content     string              `json:"content" validate:"required"`
	ContentType string              `json:"content_type"`
	Variables   string              `json:"variables"`
	IsDefault   bool                `json:"is_default"`
}

type UpdateContactTemplateRequest struct {
	Name        *string              `json:"name"`
	Description *string              `json:"description"`
	Type        *ContactTemplateType `json:"type"`
	Category    *string              `json:"category"`
	Subject     *string              `json:"subject"`
	Content     *string              `json:"content"`
	ContentType *string              `json:"content_type"`
	Variables   *string              `json:"variables"`
	IsActive    *bool                `json:"is_active"`
	IsDefault   *bool                `json:"is_default"`
}

type UpdateContactSettingsRequest struct {
	AutoAssignment       *bool   `json:"auto_assignment"`
	DefaultAssigneeID    *uuid.UUID `json:"default_assignee_id"`
	AutoReply            *bool   `json:"auto_reply"`
	AutoReplySubject     *string `json:"auto_reply_subject"`
	AutoReplyMessage     *string `json:"auto_reply_message"`
	BusinessHoursEnabled *bool   `json:"business_hours_enabled"`
	BusinessHoursStart   *string `json:"business_hours_start"`
	BusinessHoursEnd     *string `json:"business_hours_end"`
	BusinessDays         []string `json:"business_days"`
	Timezone             *string `json:"timezone"`
	SLAEnabled           *bool   `json:"sla_enabled"`
	SLAResponseTime      *int    `json:"sla_response_time"`
	SLAResolutionTime    *int    `json:"sla_resolution_time"`
	EmailNotifications   *bool   `json:"email_notifications"`
	SlackNotifications   *bool   `json:"slack_notifications"`
	SlackWebhookURL      *string `json:"slack_webhook_url"`
	CustomFields         map[string]interface{} `json:"custom_fields"`
}

// Contact represents a contact form submission or inquiry
type Contact struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Contact information
	Name    string `json:"name" gorm:"not null"`
	Email   string `json:"email" gorm:"not null"`
	Phone   string `json:"phone,omitempty"`
	Company string `json:"company,omitempty"`
	
	// Message details
	Subject  string      `json:"subject" gorm:"not null"`
	Message  string      `json:"message" gorm:"type:text;not null"`
	Type     ContactType `json:"type" gorm:"default:general"`
	Priority ContactPriority `json:"priority" gorm:"default:medium"`
	Status   ContactStatus `json:"status" gorm:"default:new"`
	
	// Customer information (if registered user)
	CustomerID *uuid.UUID `json:"customer_id,omitempty" gorm:"index"`
	OrderID    *uuid.UUID `json:"order_id,omitempty" gorm:"index"` // Related order if applicable
	
	// Assignment and handling
	AssignedToID *uuid.UUID `json:"assigned_to_id,omitempty" gorm:"index"`
	AssignedAt   *time.Time `json:"assigned_at,omitempty"`
	ResolvedBy   *uuid.UUID `json:"resolved_by,omitempty" gorm:"index"`
	ResolvedAt   *time.Time `json:"resolved_at,omitempty"`
	
	// Additional information
	Source       string   `json:"source,omitempty"`        // contact_form, email, chat, phone, etc.
	Tags         []string `json:"tags,omitempty" gorm:"serializer:json"`
	Attachments  []string `json:"attachments,omitempty" gorm:"serializer:json"`
	
	// Metadata
	IPAddress    string `json:"ip_address,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	ReferrerURL  string `json:"referrer_url,omitempty"`
	PageURL      string `json:"page_url,omitempty"`
	
	// Internal notes
	InternalNotes string `json:"internal_notes,omitempty" gorm:"type:text"`
	
	// Response tracking
	ResponseTime    *int    `json:"response_time,omitempty"` // Minutes to first response
	ResolutionTime  *int    `json:"resolution_time,omitempty"` // Minutes to resolution
	CustomerSatisfactionRating *int `json:"customer_satisfaction_rating,omitempty"` // 1-5 scale
	
	// Timestamps
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	FirstReplyAt *time.Time `json:"first_reply_at,omitempty"`
	
	// Relations
	Replies []ContactReply `json:"replies,omitempty" gorm:"foreignKey:ContactID"`
}

// ContactReply represents a reply to a contact message
type ContactReply struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	ContactID uuid.UUID `json:"contact_id" gorm:"not null;index"`
	
	// Reply author information
	UserID      *uuid.UUID `json:"user_id,omitempty" gorm:"index"` // Staff member
	CustomerID  *uuid.UUID `json:"customer_id,omitempty" gorm:"index"` // Customer reply
	AuthorName  string     `json:"author_name" gorm:"not null"`
	AuthorEmail string     `json:"author_email" gorm:"not null"`
	IsInternal  bool       `json:"is_internal" gorm:"default:false"` // Internal note vs customer-facing reply
	IsStaff     bool       `json:"is_staff" gorm:"default:false"`
	
	// Reply content
	Subject     string   `json:"subject,omitempty"`
	Content     string   `json:"content" gorm:"type:text;not null"`
	ContentType string   `json:"content_type" gorm:"default:text"` // text, html
	Attachments []string `json:"attachments,omitempty" gorm:"serializer:json"`
	
	// Delivery tracking
	SentViaEmail bool       `json:"sent_via_email" gorm:"default:false"`
	EmailSentAt  *time.Time `json:"email_sent_at,omitempty"`
	EmailStatus  string     `json:"email_status,omitempty"` // sent, delivered, bounced, failed
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ContactForm represents a contact form configuration
type ContactForm struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Form configuration
	Name        string `json:"name" gorm:"not null"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description,omitempty"`
	FormKey     string `json:"form_key" gorm:"unique;not null"` // URL-safe identifier
	
	// Form fields configuration (JSON)
	Fields      string `json:"fields" gorm:"type:json;not null"` // Dynamic form fields
	Settings    string `json:"settings,omitempty" gorm:"type:json"` // Form behavior settings
	
	// Default values
	DefaultType     ContactType     `json:"default_type" gorm:"default:general"`
	DefaultPriority ContactPriority `json:"default_priority" gorm:"default:medium"`
	DefaultAssignee *uuid.UUID      `json:"default_assignee,omitempty" gorm:"index"`
	
	// Behavior settings
	RequireAuth        bool   `json:"require_auth" gorm:"default:false"`
	AllowAttachments   bool   `json:"allow_attachments" gorm:"default:true"`
	MaxAttachments     int    `json:"max_attachments" gorm:"default:5"`
	AutoReply          bool   `json:"auto_reply" gorm:"default:true"`
	AutoReplySubject   string `json:"auto_reply_subject,omitempty"`
	AutoReplyMessage   string `json:"auto_reply_message,omitempty" gorm:"type:text"`
	
	// Spam protection
	EnableCaptcha      bool `json:"enable_captcha" gorm:"default:true"`
	EnableRateLimit    bool `json:"enable_rate_limit" gorm:"default:true"`
	RateLimitWindow    int  `json:"rate_limit_window" gorm:"default:60"`    // seconds
	RateLimitRequests  int  `json:"rate_limit_requests" gorm:"default:5"`   // requests per window
	
	// Status
	IsActive    bool `json:"is_active" gorm:"default:true"`
	IsPublic    bool `json:"is_public" gorm:"default:true"`
	
	// Usage statistics
	SubmissionCount int `json:"submission_count" gorm:"default:0"`
	SpamCount       int `json:"spam_count" gorm:"default:0"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ContactTemplate represents email templates for contact responses
type ContactTemplate struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	
	// Template information
	Name        string              `json:"name" gorm:"not null"`
	Description string              `json:"description,omitempty"`
	Type        ContactTemplateType `json:"type" gorm:"default:auto_reply"`
	Category    string              `json:"category,omitempty"` // auto_reply, follow_up, resolution, etc.
	
	// Template content
	Subject     string `json:"subject" gorm:"not null"`
	Content     string `json:"content" gorm:"type:text;not null"`
	ContentType string `json:"content_type" gorm:"default:text"` // text, html
	
	// Template variables (JSON list of available variables)
	Variables   string `json:"variables,omitempty" gorm:"type:json"`
	
	// Usage settings
	IsActive    bool `json:"is_active" gorm:"default:true"`
	IsDefault   bool `json:"is_default" gorm:"default:false"`
	UsageCount  int  `json:"usage_count" gorm:"default:0"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ContactSettings represents contact module settings for a tenant
type ContactSettings struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"unique;not null"`
	
	// General settings
	ContactEmail       string `json:"contact_email" gorm:"not null"`
	SupportEmail       string `json:"support_email,omitempty"`
	SalesEmail         string `json:"sales_email,omitempty"`
	TechnicalEmail     string `json:"technical_email,omitempty"`
	
	// Business hours
	BusinessHours      string `json:"business_hours,omitempty"`
	Timezone           string `json:"timezone" gorm:"default:Asia/Dhaka"`
	
	// Response settings
	AutoReplyEnabled   bool   `json:"auto_reply_enabled" gorm:"default:true"`
	AutoAssignEnabled  bool   `json:"auto_assign_enabled" gorm:"default:false"`
	DefaultAssigneeID  *uuid.UUID `json:"default_assignee_id,omitempty" gorm:"index"`
	
	// SLA settings (Service Level Agreement)
	SLAResponseTime    int    `json:"sla_response_time" gorm:"default:24"`     // hours
	SLAResolutionTime  int    `json:"sla_resolution_time" gorm:"default:72"`   // hours
	
	// Notification settings
	EmailNotifications       bool   `json:"email_notifications" gorm:"default:true"`
	NotifyOnNewContact       bool   `json:"notify_on_new_contact" gorm:"default:true"`
	NotifyOnAssignment       bool   `json:"notify_on_assignment" gorm:"default:true"`
	NotifyOnStatusChange     bool   `json:"notify_on_status_change" gorm:"default:false"`
	SlackWebhookURL          string `json:"slack_webhook_url,omitempty"`
	
	// Contact form settings
	AllowAnonymousContact    bool   `json:"allow_anonymous_contact" gorm:"default:true"`
	RequirePhoneNumber       bool   `json:"require_phone_number" gorm:"default:false"`
	RequireCompany           bool   `json:"require_company" gorm:"default:false"`
	
	// Spam protection
	EnableSpamFilter         bool   `json:"enable_spam_filter" gorm:"default:true"`
	SpamKeywords             string `json:"spam_keywords,omitempty" gorm:"type:text"`
	BlockedDomains           string `json:"blocked_domains,omitempty" gorm:"type:text"`
	MaxDailySubmissions      int    `json:"max_daily_submissions" gorm:"default:10"` // per IP
	
	// Integration settings
	CRMIntegrationEnabled    bool   `json:"crm_integration_enabled" gorm:"default:false"`
	CRMType                  string `json:"crm_type,omitempty"` // hubspot, salesforce, etc.
	CRMAPIKey                string `json:"crm_api_key,omitempty"`
	
	// GDPR compliance
	DataRetentionDays        int    `json:"data_retention_days" gorm:"default:365"`
	ConsentRequired          bool   `json:"consent_required" gorm:"default:true"`
	ConsentText              string `json:"consent_text,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ContactAnalytics represents contact analytics data
type ContactAnalytics struct {
	ID       uuid.UUID `json:"id" gorm:"primarykey"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"not null;index"`
	Date     time.Time `json:"date" gorm:"not null;index"`
	
	// Daily metrics
	NewContacts           int `json:"new_contacts" gorm:"default:0"`
	ResolvedContacts      int `json:"resolved_contacts" gorm:"default:0"`
	AvgResponseTimeHours  int `json:"avg_response_time_hours" gorm:"default:0"`
	AvgResolutionTimeHours int `json:"avg_resolution_time_hours" gorm:"default:0"`
	
	// Contact type breakdown
	GeneralContacts      int `json:"general_contacts" gorm:"default:0"`
	SupportContacts      int `json:"support_contacts" gorm:"default:0"`
	SalesContacts        int `json:"sales_contacts" gorm:"default:0"`
	TechnicalContacts    int `json:"technical_contacts" gorm:"default:0"`
	ComplaintContacts    int `json:"complaint_contacts" gorm:"default:0"`
	
	// Satisfaction metrics
	SatisfactionRatings   int     `json:"satisfaction_ratings" gorm:"default:0"`
	AvgSatisfactionScore  float64 `json:"avg_satisfaction_score" gorm:"default:0"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Business Logic Methods

// IsOverdue checks if a contact is overdue based on SLA
func (c *Contact) IsOverdue(slaHours int) bool {
	if c.Status == StatusResolved || c.Status == StatusClosed {
		return false
	}
	
	overdueTime := c.CreatedAt.Add(time.Duration(slaHours) * time.Hour)
	return time.Now().After(overdueTime)
}

// GetResponseTime calculates response time in minutes
func (c *Contact) GetResponseTime() *int {
	if c.FirstReplyAt == nil {
		return nil
	}
	
	minutes := int(c.FirstReplyAt.Sub(c.CreatedAt).Minutes())
	return &minutes
}

// GetResolutionTime calculates resolution time in minutes
func (c *Contact) GetResolutionTime() *int {
	if c.ResolvedAt == nil {
		return nil
	}
	
	minutes := int(c.ResolvedAt.Sub(c.CreatedAt).Minutes())
	return &minutes
}

// CanBeResolved checks if contact can be marked as resolved
func (c *Contact) CanBeResolved() bool {
	return c.Status == StatusInProgress || c.Status == StatusRead
}

// MarkAsRead marks the contact as read
func (c *Contact) MarkAsRead() {
	if c.Status == StatusNew {
		c.Status = StatusRead
		now := time.Now()
		c.ReadAt = &now
	}
}

// MarkAsResolved marks the contact as resolved
func (c *Contact) MarkAsResolved(resolvedBy uuid.UUID) {
	c.Status = StatusResolved
	c.ResolvedBy = &resolvedBy
	now := time.Now()
	c.ResolvedAt = &now
	
	// Calculate resolution time
	if c.ResolutionTime == nil {
		resolutionTime := int(now.Sub(c.CreatedAt).Minutes())
		c.ResolutionTime = &resolutionTime
	}
}

// Contact form methods

// GenerateFormKey generates a URL-safe form key
func (cf *ContactForm) GenerateFormKey() {
	// TODO: Implement unique key generation based on tenant and form name
	cf.FormKey = "contact-form-" + cf.Name
}

// IncrementSubmissionCount increments the submission counter
func (cf *ContactForm) IncrementSubmissionCount() {
	cf.SubmissionCount++
}

// IncrementSpamCount increments the spam counter
func (cf *ContactForm) IncrementSpamCount() {
	cf.SpamCount++
}

// GetSpamRate calculates spam rate percentage
func (cf *ContactForm) GetSpamRate() float64 {
	if cf.SubmissionCount == 0 {
		return 0
	}
	return float64(cf.SpamCount) / float64(cf.SubmissionCount) * 100
}

// Contact template methods

// IncrementUsageCount increments template usage counter
func (ct *ContactTemplate) IncrementUsageCount() {
	ct.UsageCount++
}

// ReplaceVariables replaces template variables with actual values
func (ct *ContactTemplate) ReplaceVariables(variables map[string]string) (string, string) {
	subject := ct.Subject
	content := ct.Content
	
	// TODO: Implement template variable replacement
	for key, value := range variables {
		// Replace {{key}} with value in both subject and content
		placeholder := "{{" + key + "}}"
		subject = strings.Replace(subject, placeholder, value, -1)
		content = strings.Replace(content, placeholder, value, -1)
	}
	
	return subject, content
}

// Validation methods

// ValidateContact validates contact data
func (c *Contact) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	
	if c.Email == "" {
		return errors.New("email is required")
	}
	
	if c.Subject == "" {
		return errors.New("subject is required")
	}
	
	if c.Message == "" {
		return errors.New("message is required")
	}
	
	return nil
}

// ValidateContactForm validates contact form configuration
func (cf *ContactForm) Validate() error {
	if cf.Name == "" {
		return errors.New("form name is required")
	}
	
	if cf.Title == "" {
		return errors.New("form title is required")
	}
	
	if cf.Fields == "" {
		return errors.New("form fields configuration is required")
	}
	
	return nil
}