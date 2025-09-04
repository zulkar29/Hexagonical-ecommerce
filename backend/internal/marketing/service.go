package marketing

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service defines the marketing service interface
type Service interface {
	// Campaign operations
	CreateCampaign(ctx context.Context, req CreateCampaignRequest) (*Campaign, error)
	GetCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) (*Campaign, error)
	GetCampaigns(ctx context.Context, tenantID uuid.UUID, filter CampaignFilter) ([]Campaign, error)
	UpdateCampaign(ctx context.Context, tenantID, campaignID uuid.UUID, req UpdateCampaignRequest) (*Campaign, error)
	DeleteCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error
	
	// Campaign execution
	ScheduleCampaign(ctx context.Context, tenantID, campaignID uuid.UUID, scheduledAt time.Time) error
	StartCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error
	PauseCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error
	StopCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error
	
	// Email operations
	GetCampaignEmails(ctx context.Context, tenantID, campaignID uuid.UUID, filter EmailFilter) ([]CampaignEmail, error)
	TrackEmailOpen(ctx context.Context, emailID uuid.UUID) error
	TrackEmailClick(ctx context.Context, emailID uuid.UUID) error
	
	// Template operations
	CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*EmailTemplate, error)
	GetTemplate(ctx context.Context, tenantID, templateID uuid.UUID) (*EmailTemplate, error)
	GetTemplates(ctx context.Context, tenantID uuid.UUID, filter TemplateFilter) ([]EmailTemplate, error)
	UpdateTemplate(ctx context.Context, tenantID, templateID uuid.UUID, req UpdateTemplateRequest) (*EmailTemplate, error)
	DeleteTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error
	
	// Segment operations
	CreateSegment(ctx context.Context, req CreateSegmentRequest) (*CustomerSegment, error)
	GetSegment(ctx context.Context, tenantID, segmentID uuid.UUID) (*CustomerSegment, error)
	GetSegments(ctx context.Context, tenantID uuid.UUID) ([]CustomerSegment, error)
	UpdateSegment(ctx context.Context, tenantID, segmentID uuid.UUID, req UpdateSegmentRequest) (*CustomerSegment, error)
	DeleteSegment(ctx context.Context, tenantID, segmentID uuid.UUID) error
	RefreshSegment(ctx context.Context, tenantID, segmentID uuid.UUID) error
	
	// Newsletter operations
	Subscribe(ctx context.Context, req SubscribeRequest) (*NewsletterSubscriber, error)
	Unsubscribe(ctx context.Context, tenantID uuid.UUID, email string) error
	GetSubscriber(ctx context.Context, tenantID uuid.UUID, email string) (*NewsletterSubscriber, error)
	GetSubscribers(ctx context.Context, tenantID uuid.UUID, filter SubscriberFilter) ([]NewsletterSubscriber, error)
	
	// Abandoned cart operations
	CreateAbandonedCart(ctx context.Context, req CreateAbandonedCartRequest) (*AbandonedCart, error)
	GetAbandonedCarts(ctx context.Context, tenantID uuid.UUID, filter AbandonedCartFilter) ([]AbandonedCart, error)
	MarkCartRecovered(ctx context.Context, tenantID, cartID uuid.UUID, recoveredValue float64) error
	SendAbandonedCartEmail(ctx context.Context, tenantID, abandonedCartID uuid.UUID) error
	
	// Settings operations
	GetSettings(ctx context.Context, tenantID uuid.UUID) (*MarketingSettings, error)
	UpdateSettings(ctx context.Context, tenantID uuid.UUID, req UpdateSettingsRequest) (*MarketingSettings, error)
	
	// Analytics
	GetCampaignStats(ctx context.Context, tenantID, campaignID uuid.UUID) (*CampaignStats, error)
	GetMarketingOverview(ctx context.Context, tenantID uuid.UUID, period string) (*MarketingOverview, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new marketing service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Request/Response DTOs
type CreateCampaignRequest struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description"`
	Type        CampaignType   `json:"type" validate:"required"`
	Subject     string         `json:"subject"`
	Content     string         `json:"content" validate:"required"`
	PreviewText string         `json:"preview_text"`
	SegmentType SegmentType    `json:"segment_type"`
	SegmentRules string        `json:"segment_rules"`
	FromName    string         `json:"from_name"`
	FromEmail   string         `json:"from_email"`
	ReplyToEmail string        `json:"reply_to_email"`
	ScheduledAt *time.Time     `json:"scheduled_at"`
}

type UpdateCampaignRequest struct {
	Name        *string        `json:"name"`
	Description *string        `json:"description"`
	Subject     *string        `json:"subject"`
	Content     *string        `json:"content"`
	PreviewText *string        `json:"preview_text"`
	SegmentType *SegmentType   `json:"segment_type"`
	SegmentRules *string       `json:"segment_rules"`
	FromName    *string        `json:"from_name"`
	FromEmail   *string        `json:"from_email"`
	ReplyToEmail *string       `json:"reply_to_email"`
	ScheduledAt *time.Time     `json:"scheduled_at"`
}

type CampaignFilter struct {
	Type      []CampaignType   `json:"type"`
	Status    []CampaignStatus `json:"status"`
	Search    string           `json:"search"`
	StartDate *time.Time       `json:"start_date"`
	EndDate   *time.Time       `json:"end_date"`
	Page      int              `json:"page"`
	Limit     int              `json:"limit"`
}

type EmailFilter struct {
	Status    []EmailStatus `json:"status"`
	Search    string        `json:"search"`
	StartDate *time.Time    `json:"start_date"`
	EndDate   *time.Time    `json:"end_date"`
	Page      int           `json:"page"`
	Limit     int           `json:"limit"`
}

type CreateTemplateRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Type        CampaignType    `json:"type"`
	Subject     string          `json:"subject" validate:"required"`
	Content     string          `json:"content" validate:"required"`
	PreviewText string          `json:"preview_text"`
	DesignJSON  string          `json:"design_json"`
	FromName    string          `json:"from_name"`
	FromEmail   string          `json:"from_email"`
	ReplyToEmail string         `json:"reply_to_email"`
}

type UpdateTemplateRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Category    *string  `json:"category"`
	Subject     *string  `json:"subject"`
	Content     *string  `json:"content"`
	PreviewText *string  `json:"preview_text"`
	DesignJSON  *string  `json:"design_json"`
	IsActive    *bool    `json:"is_active"`
	FromName    *string  `json:"from_name"`
	FromEmail   *string  `json:"from_email"`
	ReplyToEmail *string `json:"reply_to_email"`
}

type TemplateFilter struct {
	Category  string          `json:"category"`
	Type      []CampaignType  `json:"type"`
	IsActive  *bool           `json:"is_active"`
	Search    string          `json:"search"`
	Page      int             `json:"page"`
	Limit     int             `json:"limit"`
}

type CreateSegmentRequest struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description"`
	Type        SegmentType `json:"type" validate:"required"`
	Rules       string      `json:"rules"`
	AutoUpdate  bool        `json:"auto_update"`
}

type UpdateSegmentRequest struct {
	Name        *string      `json:"name"`
	Description *string      `json:"description"`
	Rules       *string      `json:"rules"`
	IsActive    *bool        `json:"is_active"`
	AutoUpdate  *bool        `json:"auto_update"`
}

type SubscribeRequest struct {
	Email           string    `json:"email" validate:"required,email"`
	Name            string    `json:"name"`
	Preferences     string    `json:"preferences"`
	Tags            []string  `json:"tags"`
	SourceURL       string    `json:"source_url"`
	SourceCampaign  string    `json:"source_campaign"`
	IPAddress       string    `json:"ip_address"`
	UserAgent       string    `json:"user_agent"`
	DoubleOptIn     bool      `json:"double_opt_in"`
}

type SubscriberFilter struct {
	Status    []string  `json:"status"`
	Tags      []string  `json:"tags"`
	Search    string    `json:"search"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Page      int       `json:"page"`
	Limit     int       `json:"limit"`
}

type CreateAbandonedCartRequest struct {
	CartID        uuid.UUID `json:"cart_id" validate:"required"`
	CustomerID    *uuid.UUID `json:"customer_id"`
	CustomerEmail string    `json:"customer_email" validate:"required,email"`
	CustomerName  string    `json:"customer_name"`
	CartValue     float64   `json:"cart_value" validate:"required,min=0"`
	ItemCount     int       `json:"item_count" validate:"required,min=1"`
	Items         string    `json:"items" validate:"required"`
}

type AbandonedCartFilter struct {
	IsRecovered *bool      `json:"is_recovered"`
	MinValue    *float64   `json:"min_value"`
	MaxValue    *float64   `json:"max_value"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Search      string     `json:"search"`
	Page        int        `json:"page"`
	Limit       int        `json:"limit"`
}

type UpdateSettingsRequest struct {
	FromName              *string `json:"from_name"`
	FromEmail             *string `json:"from_email"`
	ReplyToEmail          *string `json:"reply_to_email"`
	EmailProvider         *string `json:"email_provider"`
	SMTPHost              *string `json:"smtp_host"`
	SMTPPort              *int    `json:"smtp_port"`
	SMTPUsername          *string `json:"smtp_username"`
	SMTPPassword          *string `json:"smtp_password"`
	SendGridAPIKey        *string `json:"sendgrid_api_key"`
	MailgunAPIKey         *string `json:"mailgun_api_key"`
	MailgunDomain         *string `json:"mailgun_domain"`
	AbandonedCartEnabled  *bool   `json:"abandoned_cart_enabled"`
	AbandonedCartDelay    *int    `json:"abandoned_cart_delay"`
	WelcomeEmailEnabled   *bool   `json:"welcome_email_enabled"`
	WelcomeEmailDelay     *int    `json:"welcome_email_delay"`
	TrackingEnabled       *bool   `json:"tracking_enabled"`
	DoubleOptIn           *bool   `json:"double_opt_in"`
	UnsubscribeFooter     *string `json:"unsubscribe_footer"`
}

type CampaignStats struct {
	CampaignID     uuid.UUID `json:"campaign_id"`
	SentCount      int       `json:"sent_count"`
	DeliveredCount int       `json:"delivered_count"`
	OpenedCount    int       `json:"opened_count"`
	ClickedCount   int       `json:"clicked_count"`
	BouncedCount   int       `json:"bounced_count"`
	UnsubscribedCount int    `json:"unsubscribed_count"`
	OpenRate       float64   `json:"open_rate"`
	ClickRate      float64   `json:"click_rate"`
	BounceRate     float64   `json:"bounce_rate"`
}

type MarketingOverview struct {
	TotalCampaigns      int     `json:"total_campaigns"`
	ActiveCampaigns     int     `json:"active_campaigns"`
	TotalSubscribers    int     `json:"total_subscribers"`
	TotalEmailsSent     int     `json:"total_emails_sent"`
	AverageOpenRate     float64 `json:"average_open_rate"`
	AverageClickRate    float64 `json:"average_click_rate"`
	AbandonedCartsCount int     `json:"abandoned_carts_count"`
	RecoveredCartsCount int     `json:"recovered_carts_count"`
	RecoveryRate        float64 `json:"recovery_rate"`
}

// Implementation methods (TODO: implement business logic)
func (s *service) CreateCampaign(ctx context.Context, req CreateCampaignRequest) (*Campaign, error) {
	campaign := &Campaign{
		ID:          uuid.New(),
		TenantID:    req.TenantID,
		Name:        req.Name,
		Subject:     req.Subject,
		Content:     req.Content,
		TemplateID:  req.TemplateID,
		SegmentID:   req.SegmentID,
		Type:        req.Type,
		Status:      StatusDraft,
		ScheduledAt: req.ScheduledAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.CreateCampaign(ctx, campaign)
}

func (s *service) GetCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) (*Campaign, error) {
	return s.repo.GetCampaign(ctx, tenantID, campaignID)
}

func (s *service) GetCampaigns(ctx context.Context, tenantID uuid.UUID, filter CampaignFilter) ([]Campaign, error) {
	return s.repo.GetCampaigns(ctx, tenantID, filter)
}

func (s *service) UpdateCampaign(ctx context.Context, tenantID, campaignID uuid.UUID, req UpdateCampaignRequest) (*Campaign, error) {
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.TemplateID != nil {
		updates["template_id"] = *req.TemplateID
	}
	if req.SegmentID != nil {
		updates["segment_id"] = *req.SegmentID
	}
	if req.ScheduledAt != nil {
		updates["scheduled_at"] = *req.ScheduledAt
	}

	err := s.repo.UpdateCampaign(ctx, tenantID, campaignID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetCampaign(ctx, tenantID, campaignID)
}

func (s *service) DeleteCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error {
	return s.repo.DeleteCampaign(ctx, tenantID, campaignID)
}

func (s *service) ScheduleCampaign(ctx context.Context, tenantID, campaignID uuid.UUID, scheduledAt time.Time) error {
	return s.repo.UpdateCampaign(ctx, tenantID, campaignID, map[string]interface{}{
		"status": StatusScheduled,
		"scheduled_at": scheduledAt,
		"updated_at": time.Now(),
	})
}

func (s *service) StartCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error {
	now := time.Now()
	return s.repo.UpdateCampaign(ctx, tenantID, campaignID, map[string]interface{}{
		"status": StatusRunning,
		"started_at": &now,
		"updated_at": now,
	})
}

func (s *service) PauseCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error {
	return s.repo.UpdateCampaign(ctx, tenantID, campaignID, map[string]interface{}{
		"status": StatusPaused,
		"updated_at": time.Now(),
	})
}

func (s *service) StopCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error {
	now := time.Now()
	return s.repo.UpdateCampaign(ctx, tenantID, campaignID, map[string]interface{}{
		"status": StatusCompleted,
		"completed_at": &now,
		"updated_at": now,
	})
}

func (s *service) GetCampaignEmails(ctx context.Context, tenantID, campaignID uuid.UUID, filter EmailFilter) ([]CampaignEmail, error) {
	return s.repo.GetCampaignEmails(ctx, tenantID, campaignID, filter)
}

func (s *service) TrackEmailOpen(ctx context.Context, emailID uuid.UUID) error {
	now := time.Now()
	return s.repo.UpdateCampaignEmail(ctx, emailID, map[string]interface{}{
		"opened_at": &now,
		"status": "opened",
	})
}

func (s *service) TrackEmailClick(ctx context.Context, emailID uuid.UUID) error {
	now := time.Now()
	return s.repo.UpdateCampaignEmail(ctx, emailID, map[string]interface{}{
		"clicked_at": &now,
		"status": "clicked",
	})
}

func (s *service) CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*EmailTemplate, error) {
	template := &EmailTemplate{
		ID:        uuid.New(),
		TenantID:  req.TenantID,
		Name:      req.Name,
		Subject:   req.Subject,
		Content:   req.Content,
		Type:      req.Type,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateTemplate(ctx, template)
}

func (s *service) GetTemplate(ctx context.Context, tenantID, templateID uuid.UUID) (*EmailTemplate, error) {
	return s.repo.GetTemplate(ctx, tenantID, templateID)
}

func (s *service) GetTemplates(ctx context.Context, tenantID uuid.UUID, filter TemplateFilter) ([]EmailTemplate, error) {
	return s.repo.GetTemplates(ctx, tenantID, filter)
}

func (s *service) UpdateTemplate(ctx context.Context, tenantID, templateID uuid.UUID, req UpdateTemplateRequest) (*EmailTemplate, error) {
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	err := s.repo.UpdateTemplate(ctx, tenantID, templateID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetTemplate(ctx, tenantID, templateID)
}

func (s *service) DeleteTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error {
	return s.repo.DeleteTemplate(ctx, tenantID, templateID)
}

func (s *service) CreateSegment(ctx context.Context, req CreateSegmentRequest) (*CustomerSegment, error) {
	segment := &CustomerSegment{
		ID:           uuid.New(),
		TenantID:     req.TenantID,
		Name:         req.Name,
		Description:  req.Description,
		Conditions:   req.Conditions,
		IsActive:     true,
		CustomerCount: 0, // Will be calculated on refresh
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.repo.CreateSegment(ctx, segment)
}

func (s *service) GetSegment(ctx context.Context, tenantID, segmentID uuid.UUID) (*CustomerSegment, error) {
	return s.repo.GetSegment(ctx, tenantID, segmentID)
}

func (s *service) GetSegments(ctx context.Context, tenantID uuid.UUID) ([]CustomerSegment, error) {
	return s.repo.GetSegments(ctx, tenantID)
}

func (s *service) UpdateSegment(ctx context.Context, tenantID, segmentID uuid.UUID, req UpdateSegmentRequest) (*CustomerSegment, error) {
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Conditions != nil {
		updates["conditions"] = *req.Conditions
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	err := s.repo.UpdateSegment(ctx, tenantID, segmentID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetSegment(ctx, tenantID, segmentID)
}

func (s *service) DeleteSegment(ctx context.Context, tenantID, segmentID uuid.UUID) error {
	return s.repo.DeleteSegment(ctx, tenantID, segmentID)
}

func (s *service) RefreshSegment(ctx context.Context, tenantID, segmentID uuid.UUID) error {
	// Get segment to access conditions
	segment, err := s.repo.GetSegment(ctx, tenantID, segmentID)
	if err != nil {
		return err
	}

	// Calculate customer count based on segment conditions
	count, err := s.repo.GetSegmentCustomerCount(ctx, tenantID, segment.Conditions)
	if err != nil {
		return err
	}

	// Update segment with new customer count
	return s.repo.UpdateSegment(ctx, tenantID, segmentID, map[string]interface{}{
		"customer_count": count,
		"updated_at": time.Now(),
	})
}

func (s *service) Subscribe(ctx context.Context, req SubscribeRequest) (*NewsletterSubscriber, error) {
	subscriber := &NewsletterSubscriber{
		ID:           uuid.New(),
		TenantID:     req.TenantID,
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Status:       "active",
		Source:       req.Source,
		SubscribedAt: time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.repo.CreateSubscriber(ctx, subscriber)
}

func (s *service) Unsubscribe(ctx context.Context, tenantID uuid.UUID, email string) error {
	// TODO: Implement unsubscribe logic
	now := time.Now()
	return s.repo.UpdateSubscriber(ctx, tenantID, email, map[string]interface{}{
		"status": "unsubscribed",
		"unsubscribed_at": &now,
	})
}

func (s *service) GetSubscriber(ctx context.Context, tenantID uuid.UUID, email string) (*NewsletterSubscriber, error) {
	return s.repo.GetSubscriber(ctx, tenantID, email)
}

func (s *service) GetSubscribers(ctx context.Context, tenantID uuid.UUID, filter SubscriberFilter) ([]NewsletterSubscriber, error) {
	return s.repo.GetSubscribers(ctx, tenantID, filter)
}

func (s *service) CreateAbandonedCart(ctx context.Context, req CreateAbandonedCartRequest) (*AbandonedCart, error) {
	abandonedCart := &AbandonedCart{
		ID:            uuid.New(),
		TenantID:      req.TenantID,
		CartID:        req.CartID,
		CustomerID:    req.CustomerID,
		CustomerEmail: req.CustomerEmail,
		CustomerName:  req.CustomerName,
		CartValue:     req.CartValue,
		ItemCount:     req.ItemCount,
		Items:         req.Items,
		IsRecovered:   false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return s.repo.CreateAbandonedCart(ctx, abandonedCart)
}

func (s *service) GetAbandonedCarts(ctx context.Context, tenantID uuid.UUID, filter AbandonedCartFilter) ([]AbandonedCart, error) {
	return s.repo.GetAbandonedCarts(ctx, tenantID, filter)
}

func (s *service) MarkCartRecovered(ctx context.Context, tenantID, cartID uuid.UUID, recoveredValue float64) error {
	// TODO: Implement cart recovery logic
	now := time.Now()
	return s.repo.UpdateAbandonedCart(ctx, tenantID, cartID, map[string]interface{}{
		"is_recovered": true,
		"recovered_at": &now,
		"recovered_value": recoveredValue,
	})
}

func (s *service) SendAbandonedCartEmail(ctx context.Context, tenantID, abandonedCartID uuid.UUID) error {
	// Get abandoned cart details
	cart, err := s.repo.GetAbandonedCart(ctx, tenantID, abandonedCartID)
	if err != nil {
		return err
	}

	// TODO: Integrate with email service to send abandoned cart email
	// For now, just update the email sent timestamp
	now := time.Now()
	return s.repo.UpdateAbandonedCart(ctx, tenantID, abandonedCartID, map[string]interface{}{
		"email_sent_at": &now,
		"updated_at": now,
	})
}

func (s *service) GetSettings(ctx context.Context, tenantID uuid.UUID) (*MarketingSettings, error) {
	return s.repo.GetSettings(ctx, tenantID)
}

func (s *service) UpdateSettings(ctx context.Context, tenantID uuid.UUID, req UpdateSettingsRequest) (*MarketingSettings, error) {
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	if req.FromName != nil {
		updates["from_name"] = *req.FromName
	}
	if req.FromEmail != nil {
		updates["from_email"] = *req.FromEmail
	}
	if req.ReplyToEmail != nil {
		updates["reply_to_email"] = *req.ReplyToEmail
	}
	if req.EmailProvider != nil {
		updates["email_provider"] = *req.EmailProvider
	}
	if req.SMTPHost != nil {
		updates["smtp_host"] = *req.SMTPHost
	}
	if req.SMTPPort != nil {
		updates["smtp_port"] = *req.SMTPPort
	}
	if req.SMTPUsername != nil {
		updates["smtp_username"] = *req.SMTPUsername
	}
	if req.SMTPPassword != nil {
		updates["smtp_password"] = *req.SMTPPassword
	}
	if req.SendGridAPIKey != nil {
		updates["sendgrid_api_key"] = *req.SendGridAPIKey
	}
	if req.MailgunAPIKey != nil {
		updates["mailgun_api_key"] = *req.MailgunAPIKey
	}
	if req.MailgunDomain != nil {
		updates["mailgun_domain"] = *req.MailgunDomain
	}
	if req.AbandonedCartEnabled != nil {
		updates["abandoned_cart_enabled"] = *req.AbandonedCartEnabled
	}
	if req.AbandonedCartDelay != nil {
		updates["abandoned_cart_delay"] = *req.AbandonedCartDelay
	}
	if req.WelcomeEmailEnabled != nil {
		updates["welcome_email_enabled"] = *req.WelcomeEmailEnabled
	}
	if req.WelcomeEmailDelay != nil {
		updates["welcome_email_delay"] = *req.WelcomeEmailDelay
	}
	if req.TrackingEnabled != nil {
		updates["tracking_enabled"] = *req.TrackingEnabled
	}
	if req.DoubleOptIn != nil {
		updates["double_opt_in"] = *req.DoubleOptIn
	}
	if req.UnsubscribeFooter != nil {
		updates["unsubscribe_footer"] = *req.UnsubscribeFooter
	}

	err := s.repo.UpdateSettings(ctx, tenantID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetSettings(ctx, tenantID)
}

func (s *service) GetCampaignStats(ctx context.Context, tenantID, campaignID uuid.UUID) (*CampaignStats, error) {
	return s.repo.GetCampaignStats(ctx, tenantID, campaignID)
}

func (s *service) GetMarketingOverview(ctx context.Context, tenantID uuid.UUID, period string) (*MarketingOverview, error) {
	return s.repo.GetMarketingOverview(ctx, tenantID, period)
}