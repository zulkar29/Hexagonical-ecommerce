package contact

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service defines the contact service interface
type Service interface {
	// Contact operations
	CreateContact(ctx context.Context, req CreateContactRequest) (*Contact, error)
	GetContact(ctx context.Context, tenantID, contactID uuid.UUID) (*Contact, error)
	GetContacts(ctx context.Context, tenantID uuid.UUID, filter ContactFilter) ([]Contact, error)
	UpdateContact(ctx context.Context, tenantID, contactID uuid.UUID, req UpdateContactRequest) (*Contact, error)
	DeleteContact(ctx context.Context, tenantID, contactID uuid.UUID) error
	
	// Contact status management
	MarkAsRead(ctx context.Context, tenantID, contactID uuid.UUID, userID uuid.UUID) error
	AssignContact(ctx context.Context, tenantID, contactID uuid.UUID, assigneeID uuid.UUID) error
	UpdateContactStatus(ctx context.Context, tenantID, contactID uuid.UUID, status ContactStatus, userID uuid.UUID) error
	ResolveContact(ctx context.Context, tenantID, contactID uuid.UUID, userID uuid.UUID) error
	
	// Contact replies
	AddReply(ctx context.Context, req AddReplyRequest) (*ContactReply, error)
	GetReplies(ctx context.Context, tenantID, contactID uuid.UUID) ([]ContactReply, error)
	UpdateReply(ctx context.Context, tenantID, replyID uuid.UUID, req UpdateReplyRequest) (*ContactReply, error)
	DeleteReply(ctx context.Context, tenantID, replyID uuid.UUID) error
	
	// Contact forms
	CreateContactForm(ctx context.Context, req CreateContactFormRequest) (*ContactForm, error)
	GetContactForm(ctx context.Context, tenantID uuid.UUID, formKey string) (*ContactForm, error)
	GetContactForms(ctx context.Context, tenantID uuid.UUID) ([]ContactForm, error)
	UpdateContactForm(ctx context.Context, tenantID, formID uuid.UUID, req UpdateContactFormRequest) (*ContactForm, error)
	DeleteContactForm(ctx context.Context, tenantID, formID uuid.UUID) error
	SubmitContactForm(ctx context.Context, formKey string, req SubmitFormRequest) (*Contact, error)
	ActivateContactForm(ctx context.Context, tenantID, formID uuid.UUID) error
	DeactivateContactForm(ctx context.Context, tenantID, formID uuid.UUID) error
	
	// Templates
	CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*ContactTemplate, error)
	GetTemplate(ctx context.Context, tenantID, templateID uuid.UUID) (*ContactTemplate, error)
	GetTemplates(ctx context.Context, tenantID uuid.UUID, filter TemplateFilter) ([]ContactTemplate, error)
	UpdateTemplate(ctx context.Context, tenantID, templateID uuid.UUID, req UpdateTemplateRequest) (*ContactTemplate, error)
	DeleteTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error
	ActivateContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error
	DeactivateContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error
	
	// Settings
	GetSettings(ctx context.Context, tenantID uuid.UUID) (*ContactSettings, error)
	UpdateSettings(ctx context.Context, tenantID uuid.UUID, req UpdateSettingsRequest) (*ContactSettings, error)
	
	// Analytics and reporting
	GetContactStats(ctx context.Context, tenantID uuid.UUID, period string) (*ContactStats, error)
	GetContactTrends(ctx context.Context, tenantID uuid.UUID, period string) (*ContactTrends, error)
	GetResponseTimeAnalytics(ctx context.Context, tenantID uuid.UUID, period string) (*ResponseAnalytics, error)
	GetOverdueContacts(ctx context.Context, tenantID uuid.UUID) ([]Contact, error)
	GetAgentPerformance(ctx context.Context, tenantID uuid.UUID, period string, agentID *uuid.UUID) (*AgentPerformance, error)
	GetCustomerSatisfaction(ctx context.Context, tenantID uuid.UUID, period string) (*CustomerSatisfaction, error)
	GetResolutionTimeAnalytics(ctx context.Context, tenantID uuid.UUID, period string) (*ResolutionTimeAnalytics, error)
	
	// Email operations
	SendReply(ctx context.Context, tenantID, contactID, replyID uuid.UUID) error
	SendAutoReply(ctx context.Context, contactID uuid.UUID) error
	
	// Bulk operations
	BulkUpdateStatus(ctx context.Context, tenantID uuid.UUID, contactIDs []uuid.UUID, status ContactStatus, userID uuid.UUID) error
	BulkAssign(ctx context.Context, tenantID uuid.UUID, contactIDs []uuid.UUID, assigneeID uuid.UUID) error
	BulkDelete(ctx context.Context, tenantID uuid.UUID, contactIDs []uuid.UUID) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new contact service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Request/Response DTOs
type CreateContactRequest struct {
	Name        string         `json:"name" validate:"required"`
	Email       string         `json:"email" validate:"required,email"`
	Phone       string         `json:"phone"`
	Company     string         `json:"company"`
	Subject     string         `json:"subject" validate:"required"`
	Message     string         `json:"message" validate:"required"`
	Type        ContactType    `json:"type"`
	Priority    ContactPriority `json:"priority"`
	CustomerID  *uuid.UUID     `json:"customer_id"`
	OrderID     *uuid.UUID     `json:"order_id"`
	Source      string         `json:"source"`
	Tags        []string       `json:"tags"`
	Attachments []string       `json:"attachments"`
	IPAddress   string         `json:"ip_address"`
	UserAgent   string         `json:"user_agent"`
	ReferrerURL string         `json:"referrer_url"`
	PageURL     string         `json:"page_url"`
}

type UpdateContactRequest struct {
	Subject       *string         `json:"subject"`
	Message       *string         `json:"message"`
	Type          *ContactType    `json:"type"`
	Priority      *ContactPriority `json:"priority"`
	Status        *ContactStatus  `json:"status"`
	AssignedToID  *uuid.UUID      `json:"assigned_to_id"`
	Tags          []string        `json:"tags"`
	InternalNotes *string         `json:"internal_notes"`
	CustomerSatisfactionRating *int `json:"customer_satisfaction_rating"`
}

type ContactFilter struct {
	Status       []ContactStatus   `json:"status"`
	Type         []ContactType     `json:"type"`
	Priority     []ContactPriority `json:"priority"`
	AssignedToID *uuid.UUID        `json:"assigned_to_id"`
	CustomerID   *uuid.UUID        `json:"customer_id"`
	Source       []string          `json:"source"`
	Tags         []string          `json:"tags"`
	Search       string            `json:"search"`
	IsOverdue    *bool             `json:"is_overdue"`
	StartDate    *time.Time        `json:"start_date"`
	EndDate      *time.Time        `json:"end_date"`
	SortBy       string            `json:"sort_by"`
	SortOrder    string            `json:"sort_order"`
	Page         int               `json:"page"`
	Limit        int               `json:"limit"`
}

type AddReplyRequest struct {
	ContactID   uuid.UUID  `json:"contact_id" validate:"required"`
	UserID      *uuid.UUID `json:"user_id"`
	CustomerID  *uuid.UUID `json:"customer_id"`
	AuthorName  string     `json:"author_name" validate:"required"`
	AuthorEmail string     `json:"author_email" validate:"required,email"`
	Subject     string     `json:"subject"`
	Content     string     `json:"content" validate:"required"`
	ContentType string     `json:"content_type"`
	IsInternal  bool       `json:"is_internal"`
	IsStaff     bool       `json:"is_staff"`
	Attachments []string   `json:"attachments"`
	SendEmail   bool       `json:"send_email"`
}

type UpdateReplyRequest struct {
	Subject     *string  `json:"subject"`
	Content     *string  `json:"content"`
	ContentType *string  `json:"content_type"`
	IsInternal  *bool    `json:"is_internal"`
	Attachments []string `json:"attachments"`
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

type SubmitFormRequest struct {
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

type CreateTemplateRequest struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description"`
	Type        ContactType `json:"type"`
	Category    string      `json:"category"`
	Subject     string      `json:"subject" validate:"required"`
	Content     string      `json:"content" validate:"required"`
	ContentType string      `json:"content_type"`
	Variables   string      `json:"variables"`
	IsDefault   bool        `json:"is_default"`
}

type UpdateTemplateRequest struct {
	Name        *string      `json:"name"`
	Description *string      `json:"description"`
	Type        *ContactType `json:"type"`
	Category    *string      `json:"category"`
	Subject     *string      `json:"subject"`
	Content     *string      `json:"content"`
	ContentType *string      `json:"content_type"`
	Variables   *string      `json:"variables"`
	IsActive    *bool        `json:"is_active"`
	IsDefault   *bool        `json:"is_default"`
}

type TemplateFilter struct {
	Type      []ContactType `json:"type"`
	Category  []string      `json:"category"`
	IsActive  *bool         `json:"is_active"`
	IsDefault *bool         `json:"is_default"`
	Search    string        `json:"search"`
	Page      int           `json:"page"`
	Limit     int           `json:"limit"`
}

type UpdateSettingsRequest struct {
	ContactEmail             *string    `json:"contact_email"`
	SupportEmail             *string    `json:"support_email"`
	SalesEmail               *string    `json:"sales_email"`
	TechnicalEmail           *string    `json:"technical_email"`
	BusinessHours            *string    `json:"business_hours"`
	Timezone                 *string    `json:"timezone"`
	AutoReplyEnabled         *bool      `json:"auto_reply_enabled"`
	AutoAssignEnabled        *bool      `json:"auto_assign_enabled"`
	DefaultAssigneeID        *uuid.UUID `json:"default_assignee_id"`
	SLAResponseTime          *int       `json:"sla_response_time"`
	SLAResolutionTime        *int       `json:"sla_resolution_time"`
	EmailNotifications       *bool      `json:"email_notifications"`
	NotifyOnNewContact       *bool      `json:"notify_on_new_contact"`
	NotifyOnAssignment       *bool      `json:"notify_on_assignment"`
	NotifyOnStatusChange     *bool      `json:"notify_on_status_change"`
	SlackWebhookURL          *string    `json:"slack_webhook_url"`
	AllowAnonymousContact    *bool      `json:"allow_anonymous_contact"`
	RequirePhoneNumber       *bool      `json:"require_phone_number"`
	RequireCompany           *bool      `json:"require_company"`
	EnableSpamFilter         *bool      `json:"enable_spam_filter"`
	SpamKeywords             *string    `json:"spam_keywords"`
	BlockedDomains           *string    `json:"blocked_domains"`
	MaxDailySubmissions      *int       `json:"max_daily_submissions"`
	CRMIntegrationEnabled    *bool      `json:"crm_integration_enabled"`
	CRMType                  *string    `json:"crm_type"`
	CRMAPIKey                *string    `json:"crm_api_key"`
	DataRetentionDays        *int       `json:"data_retention_days"`
	ConsentRequired          *bool      `json:"consent_required"`
	ConsentText              *string    `json:"consent_text"`
}

// Analytics DTOs
type ContactStats struct {
	TotalContacts       int                     `json:"total_contacts"`
	NewContacts         int                     `json:"new_contacts"`
	ResolvedContacts    int                     `json:"resolved_contacts"`
	OverdueContacts     int                     `json:"overdue_contacts"`
	AvgResponseTime     string                  `json:"avg_response_time"`
	AvgResolutionTime   string                  `json:"avg_resolution_time"`
	ContactsByStatus    map[ContactStatus]int   `json:"contacts_by_status"`
	ContactsByType      map[ContactType]int     `json:"contacts_by_type"`
	ContactsByPriority  map[ContactPriority]int `json:"contacts_by_priority"`
	SatisfactionRating  float64                 `json:"satisfaction_rating"`
	ResponseRate        float64                 `json:"response_rate"`
	ResolutionRate      float64                 `json:"resolution_rate"`
}

type ContactTrends struct {
	Period           string             `json:"period"`
	TotalContacts    int                `json:"total_contacts"`
	DailyContacts    []DailyContactCount `json:"daily_contacts"`
	TypeDistribution map[ContactType]int `json:"type_distribution"`
	SourceDistribution map[string]int   `json:"source_distribution"`
	ResponseTrend    []ResponseTimePoint `json:"response_trend"`
	ResolutionTrend  []ResolutionTimePoint `json:"resolution_trend"`
}

type DailyContactCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ResponseTimePoint struct {
	Date            string  `json:"date"`
	AvgResponseTime float64 `json:"avg_response_time"`
}

type ResolutionTimePoint struct {
	Date              string  `json:"date"`
	AvgResolutionTime float64 `json:"avg_resolution_time"`
}

type ResponseAnalytics struct {
	TotalContacts       int     `json:"total_contacts"`
	ContactsWithReply   int     `json:"contacts_with_reply"`
	AvgResponseTime     float64 `json:"avg_response_time"`     // hours
	MedianResponseTime  float64 `json:"median_response_time"`  // hours
	SLABreaches         int     `json:"sla_breaches"`
	ResponseRate        float64 `json:"response_rate"`         // percentage
	Within1Hour         int     `json:"within_1_hour"`
	Within4Hours        int     `json:"within_4_hours"`
	Within24Hours       int     `json:"within_24_hours"`
	Over24Hours         int     `json:"over_24_hours"`
}

type AgentPerformance struct {
	AgentID             *uuid.UUID `json:"agent_id"`
	AgentName           string     `json:"agent_name"`
	TotalContacts       int        `json:"total_contacts"`
	ResolvedContacts    int        `json:"resolved_contacts"`
	AvgResponseTime     float64    `json:"avg_response_time"`     // hours
	AvgResolutionTime   float64    `json:"avg_resolution_time"`   // hours
	ResolutionRate      float64    `json:"resolution_rate"`       // percentage
	CustomerSatisfaction float64   `json:"customer_satisfaction"`  // average rating
	SLACompliance       float64    `json:"sla_compliance"`        // percentage
	Workload            int        `json:"workload"`              // current assigned contacts
}

type CustomerSatisfaction struct {
	TotalRatings        int                    `json:"total_ratings"`
	AverageRating       float64                `json:"average_rating"`
	RatingDistribution  map[int]int            `json:"rating_distribution"`  // rating -> count
	SatisfactionTrend   []SatisfactionPoint    `json:"satisfaction_trend"`
	ByContactType       map[ContactType]float64 `json:"by_contact_type"`
	ByAgent             map[uuid.UUID]float64   `json:"by_agent"`
	PositiveRatings     int                    `json:"positive_ratings"`     // 4-5 stars
	NegativeRatings     int                    `json:"negative_ratings"`     // 1-2 stars
	NeutralRatings      int                    `json:"neutral_ratings"`      // 3 stars
}

type SatisfactionPoint struct {
	Date   string  `json:"date"`
	Rating float64 `json:"rating"`
}

type ResolutionTimeAnalytics struct {
	TotalResolved       int                     `json:"total_resolved"`
	AvgResolutionTime   float64                 `json:"avg_resolution_time"`   // hours
	MedianResolutionTime float64                `json:"median_resolution_time"` // hours
	ResolutionTrend     []ResolutionTimePoint   `json:"resolution_trend"`
	ByPriority          map[ContactPriority]float64 `json:"by_priority"`
	ByType              map[ContactType]float64     `json:"by_type"`
	Within1Day          int                     `json:"within_1_day"`
	Within3Days         int                     `json:"within_3_days"`
	Within1Week         int                     `json:"within_1_week"`
	Over1Week           int                     `json:"over_1_week"`
	SLACompliance       float64                 `json:"sla_compliance"`        // percentage
}

// Implementation methods (TODO: implement business logic)
func (s *service) CreateContact(ctx context.Context, req CreateContactRequest) (*Contact, error) {
	contact := &Contact{
		ID:          uuid.New(),
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Company:     req.Company,
		Subject:     req.Subject,
		Message:     req.Message,
		Type:        req.Type,
		Priority:    req.Priority,
		Status:      StatusNew,
		CustomerID:  req.CustomerID,
		OrderID:     req.OrderID,
		Source:      req.Source,
		Tags:        req.Tags,
		Attachments: req.Attachments,
		IPAddress:   req.IPAddress,
		UserAgent:   req.UserAgent,
		ReferrerURL: req.ReferrerURL,
		PageURL:     req.PageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.CreateContact(ctx, contact)
}

func (s *service) GetContact(ctx context.Context, tenantID, contactID uuid.UUID) (*Contact, error) {
	return s.repo.GetContactByID(ctx, tenantID, contactID)
}

func (s *service) GetContacts(ctx context.Context, tenantID uuid.UUID, filter ContactFilter) ([]Contact, error) {
	return s.repo.GetContacts(ctx, tenantID, filter)
}

func (s *service) UpdateContact(ctx context.Context, tenantID, contactID uuid.UUID, req UpdateContactRequest) (*Contact, error) {
	updates := make(map[string]interface{})

	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}
	if req.Message != nil {
		updates["message"] = *req.Message
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.AssignedToID != nil {
		updates["assigned_to_id"] = *req.AssignedToID
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.InternalNotes != nil {
		updates["internal_notes"] = *req.InternalNotes
	}
	if req.CustomerSatisfactionRating != nil {
		updates["customer_satisfaction_rating"] = *req.CustomerSatisfactionRating
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateContact(ctx, tenantID, contactID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetContactByID(ctx, tenantID, contactID)
}

func (s *service) DeleteContact(ctx context.Context, tenantID, contactID uuid.UUID) error {
	return s.repo.DeleteContact(ctx, tenantID, contactID)
}

func (s *service) MarkAsRead(ctx context.Context, tenantID, contactID uuid.UUID, userID uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":  StatusRead,
		"read_at": &now,
	}
	
	return s.repo.UpdateContact(ctx, tenantID, contactID, updates)
}

func (s *service) AssignContact(ctx context.Context, tenantID, contactID uuid.UUID, assigneeID uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"assigned_to_id": assigneeID,
		"assigned_at":    &now,
	}
	
	return s.repo.UpdateContact(ctx, tenantID, contactID, updates)
}

func (s *service) UpdateContactStatus(ctx context.Context, tenantID, contactID uuid.UUID, status ContactStatus, userID uuid.UUID) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	if status == StatusResolved {
		now := time.Now()
		updates["resolved_by"] = userID
		updates["resolved_at"] = &now
	}
	
	return s.repo.UpdateContact(ctx, tenantID, contactID, updates)
}

func (s *service) ResolveContact(ctx context.Context, tenantID, contactID uuid.UUID, userID uuid.UUID) error {
	return s.UpdateContactStatus(ctx, tenantID, contactID, StatusResolved, userID)
}

func (s *service) AddReply(ctx context.Context, req AddReplyRequest) (*ContactReply, error) {
	reply := &ContactReply{
		ID:          uuid.New(),
		ContactID:   req.ContactID,
		UserID:      req.UserID,
		CustomerID:  req.CustomerID,
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
		Subject:     req.Subject,
		Content:     req.Content,
		ContentType: req.ContentType,
		IsInternal:  req.IsInternal,
		IsStaff:     req.IsStaff,
		Attachments: req.Attachments,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if reply.ContentType == "" {
		reply.ContentType = "text/plain"
	}

	return s.repo.CreateReply(ctx, reply)
}

func (s *service) GetReplies(ctx context.Context, tenantID, contactID uuid.UUID) ([]ContactReply, error) {
	return s.repo.GetRepliesByContactID(ctx, tenantID, contactID)
}

func (s *service) UpdateReply(ctx context.Context, tenantID, replyID uuid.UUID, req UpdateReplyRequest) (*ContactReply, error) {
	updates := make(map[string]interface{})

	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.ContentType != nil {
		updates["content_type"] = *req.ContentType
	}
	if req.IsInternal != nil {
		updates["is_internal"] = *req.IsInternal
	}
	if req.Attachments != nil {
		updates["attachments"] = req.Attachments
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateReply(ctx, tenantID, replyID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetReplyByID(ctx, tenantID, replyID)
}

func (s *service) DeleteReply(ctx context.Context, tenantID, replyID uuid.UUID) error {
	return s.repo.DeleteReply(ctx, tenantID, replyID)
}

func (s *service) CreateContactForm(ctx context.Context, req CreateContactFormRequest) (*ContactForm, error) {
	form := &ContactForm{
		ID:                uuid.New(),
		FormKey:           generateFormKey(),
		Name:              req.Name,
		Title:             req.Title,
		Description:       req.Description,
		Fields:            req.Fields,
		Settings:          req.Settings,
		DefaultType:       req.DefaultType,
		DefaultPriority:   req.DefaultPriority,
		DefaultAssignee:   req.DefaultAssignee,
		RequireAuth:       req.RequireAuth,
		AllowAttachments:  req.AllowAttachments,
		MaxAttachments:    req.MaxAttachments,
		AutoReply:         req.AutoReply,
		AutoReplySubject:  req.AutoReplySubject,
		AutoReplyMessage:  req.AutoReplyMessage,
		EnableCaptcha:     req.EnableCaptcha,
		EnableRateLimit:   req.EnableRateLimit,
		RateLimitWindow:   req.RateLimitWindow,
		RateLimitRequests: req.RateLimitRequests,
		IsActive:          true,
		IsPublic:          req.IsPublic,
		SubmissionCount:   0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return s.repo.CreateContactForm(ctx, form)
}

// generateFormKey creates a unique form key
func generateFormKey() string {
	return fmt.Sprintf("form_%s", uuid.New().String()[:8])
}

func (s *service) GetContactForm(ctx context.Context, tenantID uuid.UUID, formKey string) (*ContactForm, error) {
	return s.repo.GetContactFormByKey(ctx, tenantID, formKey)
}

func (s *service) GetContactForms(ctx context.Context, tenantID uuid.UUID) ([]ContactForm, error) {
	return s.repo.GetContactForms(ctx, tenantID)
}

func (s *service) UpdateContactForm(ctx context.Context, tenantID, formID uuid.UUID, req UpdateContactFormRequest) (*ContactForm, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Fields != nil {
		updates["fields"] = *req.Fields
	}
	if req.Settings != nil {
		updates["settings"] = *req.Settings
	}
	if req.DefaultType != nil {
		updates["default_type"] = *req.DefaultType
	}
	if req.DefaultPriority != nil {
		updates["default_priority"] = *req.DefaultPriority
	}
	if req.DefaultAssignee != nil {
		updates["default_assignee"] = *req.DefaultAssignee
	}
	if req.RequireAuth != nil {
		updates["require_auth"] = *req.RequireAuth
	}
	if req.AllowAttachments != nil {
		updates["allow_attachments"] = *req.AllowAttachments
	}
	if req.MaxAttachments != nil {
		updates["max_attachments"] = *req.MaxAttachments
	}
	if req.AutoReply != nil {
		updates["auto_reply"] = *req.AutoReply
	}
	if req.AutoReplySubject != nil {
		updates["auto_reply_subject"] = *req.AutoReplySubject
	}
	if req.AutoReplyMessage != nil {
		updates["auto_reply_message"] = *req.AutoReplyMessage
	}
	if req.EnableCaptcha != nil {
		updates["enable_captcha"] = *req.EnableCaptcha
	}
	if req.EnableRateLimit != nil {
		updates["enable_rate_limit"] = *req.EnableRateLimit
	}
	if req.RateLimitWindow != nil {
		updates["rate_limit_window"] = *req.RateLimitWindow
	}
	if req.RateLimitRequests != nil {
		updates["rate_limit_requests"] = *req.RateLimitRequests
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateContactForm(ctx, tenantID, formID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetContactFormByID(ctx, tenantID, formID)
}

func (s *service) DeleteContactForm(ctx context.Context, tenantID, formID uuid.UUID) error {
	return s.repo.DeleteContactForm(ctx, tenantID, formID)
}

func (s *service) SubmitContactForm(ctx context.Context, formKey string, req SubmitFormRequest) (*Contact, error) {
	// Get the form configuration
	form, err := s.repo.GetContactFormByKey(ctx, uuid.Nil, formKey)
	if err != nil {
		return nil, fmt.Errorf("form not found: %w", err)
	}

	if !form.IsActive {
		return nil, fmt.Errorf("form is not active")
	}

	// Create contact from form submission
	contact := &Contact{
		ID:          uuid.New(),
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Company:     req.Company,
		Subject:     req.Subject,
		Message:     req.Message,
		Type:        form.DefaultType,
		Priority:    form.DefaultPriority,
		Status:      StatusNew,
		AssignedToID: form.DefaultAssignee,
		Source:      "contact_form",
		Attachments: req.Attachments,
		IPAddress:   req.IPAddress,
		UserAgent:   req.UserAgent,
		ReferrerURL: req.ReferrerURL,
		PageURL:     req.PageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create the contact
	createdContact, err := s.repo.CreateContact(ctx, contact)
	if err != nil {
		return nil, err
	}

	// Increment form submission count
	s.repo.IncrementFormSubmissions(ctx, form.TenantID, form.ID)

	return createdContact, nil
}

func (s *service) CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*ContactTemplate, error) {
	template := &ContactTemplate{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Category:    req.Category,
		Subject:     req.Subject,
		Content:     req.Content,
		ContentType: req.ContentType,
		Variables:   req.Variables,
		IsActive:    true,
		IsDefault:   req.IsDefault,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if template.ContentType == "" {
		template.ContentType = "text/html"
	}

	return s.repo.CreateTemplate(ctx, template)
}

func (s *service) GetTemplate(ctx context.Context, tenantID, templateID uuid.UUID) (*ContactTemplate, error) {
	return s.repo.GetTemplateByID(ctx, tenantID, templateID)
}

func (s *service) GetTemplates(ctx context.Context, tenantID uuid.UUID, filter TemplateFilter) ([]ContactTemplate, error) {
	return s.repo.GetTemplates(ctx, tenantID, filter)
}

func (s *service) UpdateTemplate(ctx context.Context, tenantID, templateID uuid.UUID, req UpdateTemplateRequest) (*ContactTemplate, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.ContentType != nil {
		updates["content_type"] = *req.ContentType
	}
	if req.Variables != nil {
		updates["variables"] = *req.Variables
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateTemplate(ctx, tenantID, templateID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetTemplateByID(ctx, tenantID, templateID)
}

func (s *service) DeleteTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error {
	return s.repo.DeleteTemplate(ctx, tenantID, templateID)
}

func (s *service) GetSettings(ctx context.Context, tenantID uuid.UUID) (*ContactSettings, error) {
	return s.repo.GetSettings(ctx, tenantID)
}

func (s *service) UpdateSettings(ctx context.Context, tenantID uuid.UUID, req UpdateSettingsRequest) (*ContactSettings, error) {
	updates := make(map[string]interface{})

	if req.ContactEmail != nil {
		updates["contact_email"] = *req.ContactEmail
	}
	if req.SupportEmail != nil {
		updates["support_email"] = *req.SupportEmail
	}
	if req.SalesEmail != nil {
		updates["sales_email"] = *req.SalesEmail
	}
	if req.TechnicalEmail != nil {
		updates["technical_email"] = *req.TechnicalEmail
	}
	if req.BusinessHours != nil {
		updates["business_hours"] = *req.BusinessHours
	}
	if req.Timezone != nil {
		updates["timezone"] = *req.Timezone
	}
	if req.AutoReplyEnabled != nil {
		updates["auto_reply_enabled"] = *req.AutoReplyEnabled
	}
	if req.AutoAssignEnabled != nil {
		updates["auto_assign_enabled"] = *req.AutoAssignEnabled
	}
	if req.DefaultAssigneeID != nil {
		updates["default_assignee_id"] = *req.DefaultAssigneeID
	}
	if req.SLAResponseTime != nil {
		updates["sla_response_time"] = *req.SLAResponseTime
	}
	if req.SLAResolutionTime != nil {
		updates["sla_resolution_time"] = *req.SLAResolutionTime
	}
	if req.EmailNotifications != nil {
		updates["email_notifications"] = *req.EmailNotifications
	}
	if req.NotifyOnNewContact != nil {
		updates["notify_on_new_contact"] = *req.NotifyOnNewContact
	}
	if req.NotifyOnAssignment != nil {
		updates["notify_on_assignment"] = *req.NotifyOnAssignment
	}
	if req.NotifyOnStatusChange != nil {
		updates["notify_on_status_change"] = *req.NotifyOnStatusChange
	}
	if req.SlackWebhookURL != nil {
		updates["slack_webhook_url"] = *req.SlackWebhookURL
	}
	if req.AllowAnonymousContact != nil {
		updates["allow_anonymous_contact"] = *req.AllowAnonymousContact
	}
	if req.RequirePhoneNumber != nil {
		updates["require_phone_number"] = *req.RequirePhoneNumber
	}
	if req.RequireCompany != nil {
		updates["require_company"] = *req.RequireCompany
	}
	if req.EnableSpamFilter != nil {
		updates["enable_spam_filter"] = *req.EnableSpamFilter
	}
	if req.SpamKeywords != nil {
		updates["spam_keywords"] = *req.SpamKeywords
	}
	if req.BlockedDomains != nil {
		updates["blocked_domains"] = *req.BlockedDomains
	}
	if req.MaxDailySubmissions != nil {
		updates["max_daily_submissions"] = *req.MaxDailySubmissions
	}
	if req.CRMIntegrationEnabled != nil {
		updates["crm_integration_enabled"] = *req.CRMIntegrationEnabled
	}
	if req.CRMType != nil {
		updates["crm_type"] = *req.CRMType
	}
	if req.CRMAPIKey != nil {
		updates["crm_api_key"] = *req.CRMAPIKey
	}
	if req.DataRetentionDays != nil {
		updates["data_retention_days"] = *req.DataRetentionDays
	}
	if req.ConsentRequired != nil {
		updates["consent_required"] = *req.ConsentRequired
	}
	if req.ConsentText != nil {
		updates["consent_text"] = *req.ConsentText
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateSettings(ctx, tenantID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetSettings(ctx, tenantID)
}

func (s *service) GetContactStats(ctx context.Context, tenantID uuid.UUID, period string) (*ContactStats, error) {
	// Calculate date range based on period
	var startDate, endDate time.Time
	now := time.Now()

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "month":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
		endDate = now
	case "year":
		startDate = now.AddDate(-1, 0, 0)
		endDate = now
	default:
		// Default to last month
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	}

	return s.repo.GetContactStats(ctx, tenantID, startDate, endDate)
}

func (s *service) GetContactTrends(ctx context.Context, tenantID uuid.UUID, period string) (*ContactTrends, error) {
	// Calculate date range based on period
	var startDate, endDate time.Time
	now := time.Now()

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "month":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
		endDate = now
	case "year":
		startDate = now.AddDate(-1, 0, 0)
		endDate = now
	default:
		// Default to last month
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	}

	return s.repo.GetContactTrends(ctx, tenantID, startDate, endDate)
}

func (s *service) GetResponseTimeAnalytics(ctx context.Context, tenantID uuid.UUID, period string) (*ResponseAnalytics, error) {
	// Calculate date range based on period
	var startDate, endDate time.Time
	now := time.Now()

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "month":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
		endDate = now
	case "year":
		startDate = now.AddDate(-1, 0, 0)
		endDate = now
	default:
		// Default to last month
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	}

	return s.repo.GetResponseTimeAnalytics(ctx, tenantID, startDate, endDate)
}

func (s *service) GetOverdueContacts(ctx context.Context, tenantID uuid.UUID) ([]Contact, error) {
	// Get settings to determine SLA response time
	settings, err := s.repo.GetSettings(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Calculate cutoff time based on SLA response time (default to 24 hours if not set)
	slaHours := 24
	if settings.SLAResponseTime > 0 {
		slaHours = settings.SLAResponseTime
	}
	cutoffTime := time.Now().Add(-time.Duration(slaHours) * time.Hour)

	return s.repo.GetOverdueContacts(ctx, tenantID, cutoffTime)
}

func (s *service) SendReply(ctx context.Context, tenantID, contactID, replyID uuid.UUID) error {
	// Get the reply details
	reply, err := s.repo.GetReplyByID(ctx, tenantID, replyID)
	if err != nil {
		return err
	}

	// Get the contact details
	contact, err := s.repo.GetContactByID(ctx, tenantID, contactID)
	if err != nil {
		return err
	}

	// Get settings for email configuration
	settings, err := s.repo.GetSettings(ctx, tenantID)
	if err != nil {
		return err
	}

	// TODO: Implement actual email sending logic using SMTP or email service
	// This would typically involve:
	// 1. Format the email with reply content
	// 2. Set appropriate headers (Reply-To, References, etc.)
	// 3. Send via configured email service
	// 4. Update reply status to sent
	// 5. Log the email activity

	// For now, just mark the reply as sent
	updates := map[string]interface{}{
		"sent_at": time.Now(),
		"updated_at": time.Now(),
	}

	err = s.repo.UpdateReply(ctx, tenantID, replyID, updates)
	if err != nil {
		return err
	}

	// Update contact last activity
	contactUpdates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	return s.repo.UpdateContact(ctx, tenantID, contactID, contactUpdates)
}

func (s *service) SendAutoReply(ctx context.Context, contactID uuid.UUID) error {
	// Get the contact details to determine tenant
	contact, err := s.repo.GetContactByID(ctx, uuid.Nil, contactID)
	if err != nil {
		return err
	}

	// Get settings to check if auto-reply is enabled
	settings, err := s.repo.GetSettings(ctx, contact.TenantID)
	if err != nil {
		return err
	}

	if !settings.AutoReplyEnabled {
		return nil // Auto-reply is disabled
	}

	// Create auto-reply message
	autoReply := &ContactReply{
		ID: uuid.New(),
		ContactID: contactID,
		UserID: nil, // System generated
		AuthorName: "Support Team",
		AuthorEmail: settings.SupportEmail,
		Subject: "Re: " + contact.Subject,
		Content: "Thank you for contacting us. We have received your message and will respond as soon as possible.",
		ContentType: "text/plain",
		IsInternal: false,
		IsStaff: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the auto-reply
	_, err = s.repo.CreateReply(ctx, autoReply)
	if err != nil {
		return err
	}

	// TODO: Implement actual email sending logic
	// This would send the auto-reply email to the contact's email address

	// Update contact last activity
	contactUpdates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	return s.repo.UpdateContact(ctx, contact.TenantID, contactID, contactUpdates)
}

func (s *service) BulkUpdateStatus(ctx context.Context, tenantID uuid.UUID, contactIDs []uuid.UUID, status ContactStatus, userID uuid.UUID) error {
	updates := map[string]interface{}{
		"status": status,
		"updated_at": time.Now(),
	}

	if status == StatusResolved {
		now := time.Now()
		updates["resolved_by"] = userID
		updates["resolved_at"] = &now
	}

	return s.repo.BulkUpdateContacts(ctx, tenantID, contactIDs, updates)
}

func (s *service) BulkAssign(ctx context.Context, tenantID uuid.UUID, contactIDs []uuid.UUID, assigneeID uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"assigned_to_id": assigneeID,
		"assigned_at": &now,
		"updated_at": now,
	}

	return s.repo.BulkUpdateContacts(ctx, tenantID, contactIDs, updates)
}

func (s *service) BulkDelete(ctx context.Context, tenantID uuid.UUID, contactIDs []uuid.UUID) error {
	return s.repo.BulkDeleteContacts(ctx, tenantID, contactIDs)
}

// Form activation/deactivation methods
func (s *service) ActivateContactForm(ctx context.Context, tenantID, formID uuid.UUID) error {
	updates := map[string]interface{}{
		"is_active": true,
		"updated_at": time.Now(),
	}
	return s.repo.UpdateContactForm(ctx, tenantID, formID, updates)
}

func (s *service) DeactivateContactForm(ctx context.Context, tenantID, formID uuid.UUID) error {
	updates := map[string]interface{}{
		"is_active": false,
		"updated_at": time.Now(),
	}
	return s.repo.UpdateContactForm(ctx, tenantID, formID, updates)
}

// Template activation/deactivation methods
func (s *service) ActivateContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error {
	updates := map[string]interface{}{
		"is_active": true,
		"updated_at": time.Now(),
	}
	return s.repo.UpdateContactTemplate(ctx, tenantID, templateID, updates)
}

func (s *service) DeactivateContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error {
	updates := map[string]interface{}{
		"is_active": false,
		"updated_at": time.Now(),
	}
	return s.repo.UpdateContactTemplate(ctx, tenantID, templateID, updates)
}

// Advanced analytics methods
func (s *service) GetAgentPerformance(ctx context.Context, tenantID uuid.UUID, period string, agentID *uuid.UUID) (*AgentPerformance, error) {
	// Calculate date range based on period
	var startDate, endDate time.Time
	now := time.Now()

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "month":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
		endDate = now
	case "year":
		startDate = now.AddDate(-1, 0, 0)
		endDate = now
	default:
		// Default to last month
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	}

	return s.repo.GetAgentPerformance(ctx, tenantID, startDate, endDate, agentID)
}

func (s *service) GetCustomerSatisfaction(ctx context.Context, tenantID uuid.UUID, period string) (*CustomerSatisfaction, error) {
	// Calculate date range based on period
	var startDate, endDate time.Time
	now := time.Now()

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "month":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
		endDate = now
	case "year":
		startDate = now.AddDate(-1, 0, 0)
		endDate = now
	default:
		// Default to last month
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	}

	return s.repo.GetCustomerSatisfaction(ctx, tenantID, startDate, endDate)
}

func (s *service) GetResolutionTimeAnalytics(ctx context.Context, tenantID uuid.UUID, period string) (*ResolutionTimeAnalytics, error) {
	// Calculate date range based on period
	var startDate, endDate time.Time
	now := time.Now()

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "month":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
		endDate = now
	case "year":
		startDate = now.AddDate(-1, 0, 0)
		endDate = now
	default:
		// Default to last month
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	}

	return s.repo.GetResolutionTimeAnalytics(ctx, tenantID, startDate, endDate)
}