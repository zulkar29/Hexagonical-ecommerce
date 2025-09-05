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
	CreateContact(ctx context.Context, tenantID uuid.UUID, req CreateContactRequest) (*Contact, error)
	GetContactByID(ctx context.Context, tenantID, contactID uuid.UUID) (*Contact, error)
	ListContacts(ctx context.Context, tenantID uuid.UUID, filter ContactFilter) ([]Contact, int64, error)
	UpdateContact(ctx context.Context, tenantID, contactID uuid.UUID, req UpdateContactRequest) (*Contact, error)
	DeleteContact(ctx context.Context, tenantID, contactID uuid.UUID) error
	BulkUpdateContacts(ctx context.Context, tenantID uuid.UUID, req BulkUpdateContactsRequest) (interface{}, error)
	ExportContacts(ctx context.Context, tenantID uuid.UUID, req ExportContactsRequest) (interface{}, error)
	
	// Contact status management
	MarkAsRead(ctx context.Context, tenantID, contactID uuid.UUID, userID uuid.UUID) error
	AssignContact(ctx context.Context, tenantID, contactID uuid.UUID, req AssignContactRequest) error
	UpdateContactStatus(ctx context.Context, tenantID, contactID uuid.UUID, req UpdateContactStatusRequest) error
	UpdateContactPriority(ctx context.Context, tenantID, contactID uuid.UUID, req UpdateContactPriorityRequest) error
	AddContactTags(ctx context.Context, tenantID, contactID uuid.UUID, req AddContactTagsRequest) error
	RemoveContactTags(ctx context.Context, tenantID, contactID uuid.UUID, req RemoveContactTagsRequest) error
	ResolveContact(ctx context.Context, tenantID, contactID uuid.UUID, userID uuid.UUID) error
	
	// Contact replies
	CreateContactReply(ctx context.Context, tenantID, contactID uuid.UUID, req CreateContactReplyRequest) (*ContactReply, error)
	ListContactReplies(ctx context.Context, tenantID, contactID uuid.UUID) ([]ContactReply, error)
	UpdateReply(ctx context.Context, tenantID, replyID uuid.UUID, req UpdateReplyRequest) (*ContactReply, error)
	DeleteContactReply(ctx context.Context, tenantID, contactID, replyID uuid.UUID) error
	
	// Contact forms
	CreateContactForm(ctx context.Context, tenantID uuid.UUID, req CreateContactFormRequest) (*ContactForm, error)
	ListContactForms(ctx context.Context, tenantID uuid.UUID) ([]ContactForm, error)
	GetContactFormByID(ctx context.Context, tenantID, formID uuid.UUID) (*ContactForm, error)
	GetPublicContactForm(ctx context.Context, formType string) (*ContactForm, error)
	UpdateContactForm(ctx context.Context, tenantID, formID uuid.UUID, req UpdateContactFormRequest) (*ContactForm, error)
	DeleteContactForm(ctx context.Context, tenantID, formID uuid.UUID) error
	SubmitPublicContactForm(ctx context.Context, formType string, req SubmitContactFormRequest) (*Contact, error)
	ActivateContactForm(ctx context.Context, tenantID, formID uuid.UUID) error
	DeactivateContactForm(ctx context.Context, tenantID, formID uuid.UUID) error
	
	// Templates
	CreateContactTemplate(ctx context.Context, tenantID uuid.UUID, req CreateContactTemplateRequest) (*ContactTemplate, error)
	ListContactTemplates(ctx context.Context, tenantID uuid.UUID) ([]ContactTemplate, error)
	GetContactTemplateByID(ctx context.Context, tenantID, templateID uuid.UUID) (*ContactTemplate, error)
	UpdateContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID, req UpdateContactTemplateRequest) (*ContactTemplate, error)
	DeleteContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error
	ActivateContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error
	DeactivateContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error
	
	// Settings
	GetContactSettings(ctx context.Context, tenantID uuid.UUID) (*ContactSettings, error)
	UpdateContactSettings(ctx context.Context, tenantID uuid.UUID, req UpdateContactSettingsRequest) (*ContactSettings, error)
	
	// Analytics and reporting
	GetContactAnalytics(ctx context.Context, tenantID uuid.UUID, period AnalyticsPeriod) (interface{}, error)
	GetContactMetrics(ctx context.Context, tenantID uuid.UUID) (*ContactMetrics, error)
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
	Offset       int               `json:"offset"`
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
func (s *service) CreateContact(ctx context.Context, tenantID uuid.UUID, req CreateContactRequest) (*Contact, error) {
	contact := &Contact{
		ID:          uuid.New(),
		TenantID:    tenantID,
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

	err := s.repo.CreateContact(ctx, contact)
	if err != nil {
		return nil, err
	}
	
	return contact, nil
}

func (s *service) GetContactByID(ctx context.Context, tenantID, contactID uuid.UUID) (*Contact, error) {
	return s.repo.GetContactByID(ctx, tenantID, contactID)
}

func (s *service) ListContacts(ctx context.Context, tenantID uuid.UUID, filter ContactFilter) ([]Contact, int64, error) {
	contacts, total, err := s.repo.ListContacts(ctx, tenantID, filter)
	if err != nil {
		return nil, 0, err
	}
	
	// Convert []*Contact to []Contact
	result := make([]Contact, len(contacts))
	for i, contact := range contacts {
		result[i] = *contact
	}
	
	return result, total, nil
}

func (s *service) UpdateContact(ctx context.Context, tenantID, contactID uuid.UUID, req UpdateContactRequest) (*Contact, error) {
	// Get existing contact
	contact, err := s.repo.GetContactByID(ctx, tenantID, contactID)
	if err != nil {
		return nil, err
	}
	
	// Update fields
	if req.Subject != nil {
		contact.Subject = *req.Subject
	}
	if req.Message != nil {
		contact.Message = *req.Message
	}
	if req.Type != nil {
		contact.Type = *req.Type
	}
	if req.Priority != nil {
		contact.Priority = *req.Priority
	}
	if req.Status != nil {
		contact.Status = *req.Status
	}
	if req.AssignedToID != nil {
		contact.AssignedToID = req.AssignedToID
	}
	if req.InternalNotes != nil {
		contact.InternalNotes = *req.InternalNotes
	}
	if req.CustomerSatisfactionRating != nil {
		contact.CustomerSatisfactionRating = req.CustomerSatisfactionRating
	}
	if len(req.Tags) > 0 {
		contact.Tags = req.Tags
	}
	
	contact.UpdatedAt = time.Now()
	
	// Save to repository
	err = s.repo.UpdateContact(ctx, contact)
	if err != nil {
		return nil, err
	}
	
	return contact, nil
}

func (s *service) DeleteContact(ctx context.Context, tenantID, contactID uuid.UUID) error {
	return s.repo.DeleteContact(ctx, tenantID, contactID)
}

func (s *service) MarkAsRead(ctx context.Context, tenantID, contactID uuid.UUID, userID uuid.UUID) error {
	// Get existing contact
	contact, err := s.repo.GetContactByID(ctx, tenantID, contactID)
	if err != nil {
		return err
	}
	
	now := time.Now()
	contact.Status = StatusRead
	contact.ReadAt = &now
	contact.UpdatedAt = now
	
	err = s.repo.UpdateContact(ctx, contact)
	return err
}

func (s *service) AssignContact(ctx context.Context, tenantID, contactID uuid.UUID, req AssignContactRequest) error {
	// Get existing contact
	contact, err := s.repo.GetContactByID(ctx, tenantID, contactID)
	if err != nil {
		return err
	}
	
	// Update assignment
	contact.AssignedToID = &req.AssigneeID
	contact.UpdatedAt = time.Now()
	
	return s.repo.UpdateContact(ctx, contact)
}

func (s *service) UpdateContactStatus(ctx context.Context, tenantID, contactID uuid.UUID, req UpdateContactStatusRequest) error {
	// Get existing contact
	contact, err := s.repo.GetContactByID(ctx, tenantID, contactID)
	if err != nil {
		return err
	}
	
	contact.Status = req.Status
	contact.UpdatedAt = time.Now()
	
	if req.Status == StatusResolved {
		now := time.Now()
		contact.ResolvedAt = &now
	}
	
	err = s.repo.UpdateContact(ctx, contact)
	return err
}

func (s *service) ResolveContact(ctx context.Context, tenantID, contactID uuid.UUID, userID uuid.UUID) error {
	req := UpdateContactStatusRequest{
		Status: StatusResolved,
	}
	return s.UpdateContactStatus(ctx, tenantID, contactID, req)
}

func (s *service) AddContactTags(ctx context.Context, tenantID, contactID uuid.UUID, req AddContactTagsRequest) error {
	// Get existing contact
	contact, err := s.repo.GetContactByID(ctx, tenantID, contactID)
	if err != nil {
		return err
	}
	
	// Append new tags
	existingTags := make(map[string]bool)
	for _, tag := range contact.Tags {
		existingTags[tag] = true
	}
	
	for _, tag := range req.Tags {
		if !existingTags[tag] {
			contact.Tags = append(contact.Tags, tag)
		}
	}
	
	contact.UpdatedAt = time.Now()
	return s.repo.UpdateContact(ctx, contact)
}

func (s *service) RemoveContactTags(ctx context.Context, tenantID, contactID uuid.UUID, req RemoveContactTagsRequest) error {
	// Get existing contact
	contact, err := s.repo.GetContactByID(ctx, tenantID, contactID)
	if err != nil {
		return err
	}
	
	// Remove specified tags
	removeMap := make(map[string]bool)
	for _, tag := range req.Tags {
		removeMap[tag] = true
	}
	
	var newTags []string
	for _, tag := range contact.Tags {
		if !removeMap[tag] {
			newTags = append(newTags, tag)
		}
	}
	
	contact.Tags = newTags
	contact.UpdatedAt = time.Now()
	return s.repo.UpdateContact(ctx, contact)
}

func (s *service) UpdateContactPriority(ctx context.Context, tenantID, contactID uuid.UUID, req UpdateContactPriorityRequest) error {
	// Get existing contact
	contact, err := s.repo.GetContactByID(ctx, tenantID, contactID)
	if err != nil {
		return err
	}
	
	// Update priority
	contact.Priority = req.Priority
	contact.UpdatedAt = time.Now()
	
	return s.repo.UpdateContact(ctx, contact)
}

func (s *service) BulkUpdateContacts(ctx context.Context, tenantID uuid.UUID, req BulkUpdateContactsRequest) (interface{}, error) {
	// TODO: Implement bulk update functionality
	return nil, fmt.Errorf("bulk update contacts not implemented")
}

func (s *service) ExportContacts(ctx context.Context, tenantID uuid.UUID, req ExportContactsRequest) (interface{}, error) {
	// TODO: Implement export functionality
	return nil, fmt.Errorf("export contacts not implemented")
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

	err := s.repo.CreateContactReply(ctx, reply)
	if err != nil {
		return nil, err
	}
	
	return reply, nil
}

func (s *service) CreateContactReply(ctx context.Context, tenantID, contactID uuid.UUID, req CreateContactReplyRequest) (*ContactReply, error) {
	reply := &ContactReply{
		ID:          uuid.New(),
		ContactID:   contactID,
		UserID:      &req.AuthorID,
		AuthorName:  "Staff User", // Default staff name
		AuthorEmail: "support@company.com", // Default staff email
		Subject:     "Re: Contact Reply",
		Content:     req.Content,
		ContentType: "text/plain",
		IsInternal:  req.IsInternal,
		IsStaff:     true,
		Attachments: req.Attachments,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.CreateContactReply(ctx, reply)
	if err != nil {
		return nil, err
	}
	
	return reply, nil
}

func (s *service) ListContactReplies(ctx context.Context, tenantID, contactID uuid.UUID) ([]ContactReply, error) {
	replies, err := s.repo.ListContactReplies(ctx, tenantID, contactID)
	if err != nil {
		return nil, err
	}
	
	// Convert []*ContactReply to []ContactReply
	result := make([]ContactReply, len(replies))
	for i, reply := range replies {
		result[i] = *reply
	}
	
	return result, nil
}

func (s *service) DeleteContactReply(ctx context.Context, tenantID, contactID, replyID uuid.UUID) error {
	return s.repo.DeleteContactReply(ctx, tenantID, replyID)
}

func (s *service) UpdateReply(ctx context.Context, tenantID, replyID uuid.UUID, req UpdateReplyRequest) (*ContactReply, error) {
	// TODO: Implement reply update functionality
	// This would require adding UpdateContactReply and GetContactReplyByID methods to the repository
	return nil, fmt.Errorf("update reply not implemented")
}

func (s *service) DeleteReply(ctx context.Context, tenantID, replyID uuid.UUID) error {
	return s.repo.DeleteContactReply(ctx, tenantID, replyID)
}

func (s *service) CreateContactForm(ctx context.Context, tenantID uuid.UUID, req CreateContactFormRequest) (*ContactForm, error) {
	form := &ContactForm{
		ID:                uuid.New(),
		TenantID:          tenantID,
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

	err := s.repo.CreateContactForm(ctx, form)
	if err != nil {
		return nil, err
	}
	
	return form, nil
}

// generateFormKey creates a unique form key
func generateFormKey() string {
	return fmt.Sprintf("form_%s", uuid.New().String()[:8])
}

func (s *service) GetPublicContactForm(ctx context.Context, formType string) (*ContactForm, error) {
	// TODO: Implement public form retrieval by type
	// This would require adding GetContactFormByType method to repository
	return nil, fmt.Errorf("get public contact form not implemented")
}

func (s *service) ListContactForms(ctx context.Context, tenantID uuid.UUID) ([]ContactForm, error) {
	filter := ContactFormFilter{} // Empty filter to get all forms
	forms, _, err := s.repo.ListContactForms(ctx, tenantID, filter)
	if err != nil {
		return nil, err
	}
	
	// Convert []*ContactForm to []ContactForm
	result := make([]ContactForm, len(forms))
	for i, form := range forms {
		result[i] = *form
	}
	
	return result, nil
}

func (s *service) GetContactFormByID(ctx context.Context, tenantID, formID uuid.UUID) (*ContactForm, error) {
	return s.repo.GetContactFormByID(ctx, tenantID, formID)
}

func (s *service) UpdateContactForm(ctx context.Context, tenantID, formID uuid.UUID, req UpdateContactFormRequest) (*ContactForm, error) {
	// Get existing form
	form, err := s.repo.GetContactFormByID(ctx, tenantID, formID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		form.Name = *req.Name
	}
	if req.Title != nil {
		form.Title = *req.Title
	}
	if req.Description != nil {
		form.Description = *req.Description
	}
	if req.Fields != nil {
		form.Fields = *req.Fields
	}
	if req.Settings != nil {
		form.Settings = *req.Settings
	}
	if req.DefaultType != nil {
		form.DefaultType = *req.DefaultType
	}
	if req.DefaultPriority != nil {
		form.DefaultPriority = *req.DefaultPriority
	}
	if req.DefaultAssignee != nil {
		form.DefaultAssignee = req.DefaultAssignee
	}
	if req.RequireAuth != nil {
		form.RequireAuth = *req.RequireAuth
	}
	if req.AllowAttachments != nil {
		form.AllowAttachments = *req.AllowAttachments
	}
	if req.MaxAttachments != nil {
		form.MaxAttachments = *req.MaxAttachments
	}
	if req.AutoReply != nil {
		form.AutoReply = *req.AutoReply
	}
	if req.AutoReplySubject != nil {
		form.AutoReplySubject = *req.AutoReplySubject
	}
	if req.AutoReplyMessage != nil {
		form.AutoReplyMessage = *req.AutoReplyMessage
	}
	if req.EnableCaptcha != nil {
		form.EnableCaptcha = *req.EnableCaptcha
	}
	if req.EnableRateLimit != nil {
		form.EnableRateLimit = *req.EnableRateLimit
	}
	if req.RateLimitWindow != nil {
		form.RateLimitWindow = *req.RateLimitWindow
	}
	if req.RateLimitRequests != nil {
		form.RateLimitRequests = *req.RateLimitRequests
	}
	if req.IsActive != nil {
		form.IsActive = *req.IsActive
	}
	if req.IsPublic != nil {
		form.IsPublic = *req.IsPublic
	}

	form.UpdatedAt = time.Now()

	err = s.repo.UpdateContactForm(ctx, form)
	if err != nil {
		return nil, err
	}

	return form, nil
}

func (s *service) DeleteContactForm(ctx context.Context, tenantID, formID uuid.UUID) error {
	return s.repo.DeleteContactForm(ctx, tenantID, formID)
}

func (s *service) SubmitPublicContactForm(ctx context.Context, formType string, req SubmitContactFormRequest) (*Contact, error) {
	// For now, create a simple contact without form validation
	// In a real implementation, you would look up the form by type
	
	contact := &Contact{
		ID:          uuid.New(),
		TenantID:    uuid.New(), // TODO: Get from form configuration
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Company:     req.Company,
		Subject:     req.Subject,
		Message:     req.Message,
		Type:        TypeGeneral,
		Priority:    PriorityMedium,
		Status:      StatusNew,
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
	err := s.repo.CreateContact(ctx, contact)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *service) CreateContactTemplate(ctx context.Context, tenantID uuid.UUID, req CreateContactTemplateRequest) (*ContactTemplate, error) {
	template := &ContactTemplate{
		ID:          uuid.New(),
		TenantID:    tenantID,
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

	err := s.repo.CreateContactTemplate(ctx, template)
	if err != nil {
		return nil, err
	}
	
	return template, nil
}

func (s *service) GetContactTemplateByID(ctx context.Context, tenantID, templateID uuid.UUID) (*ContactTemplate, error) {
	return s.repo.GetContactTemplateByID(ctx, tenantID, templateID)
}

func (s *service) ListContactTemplates(ctx context.Context, tenantID uuid.UUID) ([]ContactTemplate, error) {
	filter := ContactTemplateFilter{} // Empty filter to get all templates
	templates, _, err := s.repo.ListContactTemplates(ctx, tenantID, filter)
	if err != nil {
		return nil, err
	}
	
	// Convert []*ContactTemplate to []ContactTemplate
	result := make([]ContactTemplate, len(templates))
	for i, template := range templates {
		result[i] = *template
	}
	
	return result, nil
}

func (s *service) UpdateContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID, req UpdateContactTemplateRequest) (*ContactTemplate, error) {
	// Get existing template
	template, err := s.repo.GetContactTemplateByID(ctx, tenantID, templateID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		template.Name = *req.Name
	}
	if req.Description != nil {
		template.Description = *req.Description
	}
	if req.Type != nil {
		template.Type = *req.Type
	}
	if req.Category != nil {
		template.Category = *req.Category
	}
	if req.Subject != nil {
		template.Subject = *req.Subject
	}
	if req.Content != nil {
		template.Content = *req.Content
	}
	if req.ContentType != nil {
		template.ContentType = *req.ContentType
	}
	if req.Variables != nil {
		template.Variables = *req.Variables
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}
	if req.IsDefault != nil {
		template.IsDefault = *req.IsDefault
	}

	template.UpdatedAt = time.Now()

	err = s.repo.UpdateContactTemplate(ctx, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (s *service) DeleteContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error {
	return s.repo.DeleteContactTemplate(ctx, tenantID, templateID)
}

func (s *service) GetContactSettings(ctx context.Context, tenantID uuid.UUID) (*ContactSettings, error) {
	return s.repo.GetContactSettings(ctx, tenantID)
}

func (s *service) UpdateContactSettings(ctx context.Context, tenantID uuid.UUID, req UpdateContactSettingsRequest) (*ContactSettings, error) {
	// Get existing settings or create default
	settings, err := s.repo.GetContactSettings(ctx, tenantID)
	if err != nil {
		// Create default settings if none exist
		settings = &ContactSettings{
			ID:                       uuid.New(),
			TenantID:                 tenantID,
			ContactEmail:            "contact@example.com",
			AutoReplyEnabled:        true,
			AutoAssignEnabled:       false,
			SLAResponseTime:         24,
			SLAResolutionTime:       72,
			EmailNotifications:      true,
			NotifyOnNewContact:      true,
			NotifyOnAssignment:      true,
			NotifyOnStatusChange:    false,
			AllowAnonymousContact:   true,
			RequirePhoneNumber:      false,
			RequireCompany:          false,
			EnableSpamFilter:        true,
			MaxDailySubmissions:     10,
			CRMIntegrationEnabled:   false,
			DataRetentionDays:       365,
			ConsentRequired:         true,
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
		}
	}

	// Update fields based on request
	if req.AutoAssignment != nil {
		settings.AutoAssignEnabled = *req.AutoAssignment
	}
	if req.DefaultAssigneeID != nil {
		settings.DefaultAssigneeID = req.DefaultAssigneeID
	}
	if req.AutoReply != nil {
		settings.AutoReplyEnabled = *req.AutoReply
	}
	if req.BusinessHoursEnabled != nil {
		// Map to existing field or add new field as needed
	}
	if req.Timezone != nil {
		settings.Timezone = *req.Timezone
	}
	if req.SLAEnabled != nil {
		// This field might need to be added to the settings struct
	}
	if req.SLAResponseTime != nil {
		settings.SLAResponseTime = *req.SLAResponseTime
	}
	if req.SLAResolutionTime != nil {
		settings.SLAResolutionTime = *req.SLAResolutionTime
	}
	if req.EmailNotifications != nil {
		settings.EmailNotifications = *req.EmailNotifications
	}
	if req.SlackNotifications != nil {
		// Map to slack webhook or similar field
	}
	if req.SlackWebhookURL != nil {
		settings.SlackWebhookURL = *req.SlackWebhookURL
	}

	settings.UpdatedAt = time.Now()

	err = s.repo.UpdateContactSettings(ctx, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *service) GetContactStats(ctx context.Context, tenantID uuid.UUID, period string) (*ContactStats, error) {
	// Calculate date range based on period
	_ = period // Use period parameter to avoid unused warning

	// Period-based analytics logic would go here

	// For now, return empty stats - implement analytics in repository later
	return &ContactStats{
		TotalContacts:      0,
		NewContacts:        0,
		ResolvedContacts:   0,
		OverdueContacts:    0,
		ContactsByStatus:   make(map[ContactStatus]int),
		ContactsByType:     make(map[ContactType]int),
		ContactsByPriority: make(map[ContactPriority]int),
	}, nil
}

func (s *service) GetContactTrends(ctx context.Context, tenantID uuid.UUID, period string) (*ContactTrends, error) {
	// Calculate date range based on period
	_ = period // Use period parameter to avoid unused warning

	// Period-based analytics logic would go here

	// Return empty trends - implement analytics later
	return &ContactTrends{
		Period:             period,
		TotalContacts:      0,
		DailyContacts:      []DailyContactCount{},
		TypeDistribution:   make(map[ContactType]int),
		SourceDistribution: make(map[string]int),
		ResponseTrend:      []ResponseTimePoint{},
		ResolutionTrend:    []ResolutionTimePoint{},
	}, nil
}

func (s *service) GetResponseTimeAnalytics(ctx context.Context, tenantID uuid.UUID, period string) (*ResponseAnalytics, error) {
	// Calculate date range based on period
	_ = period // Use period parameter to avoid unused warning

	// Period-based analytics logic would go here

	// Return empty analytics - implement later
	return &ResponseAnalytics{
		TotalContacts:      0,
		ContactsWithReply:  0,
		AvgResponseTime:    0,
		MedianResponseTime: 0,
		SLABreaches:        0,
		ResponseRate:       0,
		Within1Hour:        0,
		Within4Hours:       0,
		Within24Hours:      0,
		Over24Hours:        0,
	}, nil
}

func (s *service) GetOverdueContacts(ctx context.Context, tenantID uuid.UUID) ([]Contact, error) {
	// For now, return empty slice - implement later
	return []Contact{}, nil
}

func (s *service) SendReply(ctx context.Context, tenantID, contactID, replyID uuid.UUID) error {
	// TODO: Implement email sending functionality
	// For now, just return nil
	return nil
}

func (s *service) SendAutoReply(ctx context.Context, contactID uuid.UUID) error {
	// TODO: Implement auto-reply functionality
	return nil
}

func (s *service) BulkUpdateStatus(ctx context.Context, tenantID uuid.UUID, contactIDs []uuid.UUID, status ContactStatus, userID uuid.UUID) error {
	// TODO: Implement bulk status update
	return nil
}

func (s *service) BulkAssign(ctx context.Context, tenantID uuid.UUID, contactIDs []uuid.UUID, assigneeID uuid.UUID) error {
	// TODO: Implement bulk assignment
	return nil
}

func (s *service) BulkDelete(ctx context.Context, tenantID uuid.UUID, contactIDs []uuid.UUID) error {
	// TODO: Implement bulk deletion
	return nil
}

// Form activation/deactivation methods
func (s *service) ActivateContactForm(ctx context.Context, tenantID, formID uuid.UUID) error {
	form, err := s.repo.GetContactFormByID(ctx, tenantID, formID)
	if err != nil {
		return err
	}
	form.IsActive = true
	form.UpdatedAt = time.Now()
	return s.repo.UpdateContactForm(ctx, form)
}

func (s *service) DeactivateContactForm(ctx context.Context, tenantID, formID uuid.UUID) error {
	form, err := s.repo.GetContactFormByID(ctx, tenantID, formID)
	if err != nil {
		return err
	}
	form.IsActive = false
	form.UpdatedAt = time.Now()
	return s.repo.UpdateContactForm(ctx, form)
}

// Template activation/deactivation methods
func (s *service) ActivateContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error {
	template, err := s.repo.GetContactTemplateByID(ctx, tenantID, templateID)
	if err != nil {
		return err
	}
	template.IsActive = true
	template.UpdatedAt = time.Now()
	return s.repo.UpdateContactTemplate(ctx, template)
}

func (s *service) DeactivateContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error {
	template, err := s.repo.GetContactTemplateByID(ctx, tenantID, templateID)
	if err != nil {
		return err
	}
	template.IsActive = false
	template.UpdatedAt = time.Now()
	return s.repo.UpdateContactTemplate(ctx, template)
}

// Advanced analytics methods
func (s *service) GetAgentPerformance(ctx context.Context, tenantID uuid.UUID, period string, agentID *uuid.UUID) (*AgentPerformance, error) {
	// Calculate date range based on period
	_ = period // Use period parameter to avoid unused warning

	// Period-based analytics logic would go here

	// Return empty performance data - implement analytics later
	return &AgentPerformance{
		AgentID:             agentID,
		TotalContacts:       0,
		ResolvedContacts:    0,
		AvgResponseTime:     0,
		AvgResolutionTime:   0,
		CustomerSatisfaction: 0,
		ResolutionRate:      0,
		SLACompliance:       0,
		Workload:            0,
	}, nil
}

func (s *service) GetCustomerSatisfaction(ctx context.Context, tenantID uuid.UUID, period string) (*CustomerSatisfaction, error) {
	// Calculate date range based on period
	_ = period // Use period parameter to avoid unused warning

	// Period-based analytics logic would go here

	// Return empty satisfaction data - implement analytics later
	return &CustomerSatisfaction{
		TotalRatings:        0,
		AverageRating:       0,
		RatingDistribution:  make(map[int]int),
		SatisfactionTrend:   []SatisfactionPoint{},
		ByContactType:       make(map[ContactType]float64),
	}, nil
}

func (s *service) GetResolutionTimeAnalytics(ctx context.Context, tenantID uuid.UUID, period string) (*ResolutionTimeAnalytics, error) {
	// Calculate date range based on period
	_ = period // Use period parameter to avoid unused warning

	// Period-based analytics logic would go here

	// Return empty resolution time analytics - implement later
	return &ResolutionTimeAnalytics{
		TotalResolved:       0,
		AvgResolutionTime:   0,
		MedianResolutionTime: 0,
	}, nil
}

func (s *service) GetContactAnalytics(ctx context.Context, tenantID uuid.UUID, period AnalyticsPeriod) (interface{}, error) {
	// Return basic analytics data
	return map[string]interface{}{
		"period": period,
		"stats": map[string]int{
			"total_contacts": 0,
			"resolved": 0,
			"pending": 0,
		},
	}, nil
}

func (s *service) GetContactMetrics(ctx context.Context, tenantID uuid.UUID) (*ContactMetrics, error) {
	// Return empty metrics for now
	return &ContactMetrics{
		TenantID:          tenantID,
		StartDate:         time.Now().AddDate(0, 0, -30), // Last 30 days
		EndDate:           time.Now(),
		TotalContacts:     0,
		NewContacts:       0,
		ResolvedContacts:  0,
		AvgResponseTime:   0,
		AvgResolutionTime: 0,
		SatisfactionScore: 0,
		TypeDistribution:  make(map[ContactType]int),
	}, nil
}