package support

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service defines the support service interface
type Service interface {
	// Ticket operations
	CreateTicket(ctx context.Context, req CreateTicketRequest) (*Ticket, error)
	GetTicket(ctx context.Context, tenantID, ticketID uuid.UUID) (*Ticket, error)
	GetTickets(ctx context.Context, tenantID uuid.UUID, filter TicketFilter) ([]Ticket, error)
	UpdateTicket(ctx context.Context, tenantID, ticketID uuid.UUID, req UpdateTicketRequest) (*Ticket, error)
	DeleteTicket(ctx context.Context, tenantID, ticketID uuid.UUID) error
	AssignTicket(ctx context.Context, tenantID, ticketID, userID uuid.UUID) error
	ResolveTicket(ctx context.Context, tenantID, ticketID uuid.UUID) error
	
	// Ticket message operations
	AddMessage(ctx context.Context, req AddMessageRequest) (*TicketMessage, error)
	GetMessages(ctx context.Context, tenantID, ticketID uuid.UUID) ([]TicketMessage, error)
	
	// FAQ operations
	CreateFAQ(ctx context.Context, req CreateFAQRequest) (*FAQ, error)
	GetFAQ(ctx context.Context, tenantID, faqID uuid.UUID) (*FAQ, error)
	GetFAQs(ctx context.Context, tenantID uuid.UUID, filter FAQFilter) ([]FAQ, error)
	UpdateFAQ(ctx context.Context, tenantID, faqID uuid.UUID, req UpdateFAQRequest) (*FAQ, error)
	DeleteFAQ(ctx context.Context, tenantID, faqID uuid.UUID) error
	
	// Knowledge base operations
	CreateArticle(ctx context.Context, req CreateArticleRequest) (*KnowledgeBase, error)
	GetArticle(ctx context.Context, tenantID uuid.UUID, slug string) (*KnowledgeBase, error)
	GetArticles(ctx context.Context, tenantID uuid.UUID, filter ArticleFilter) ([]KnowledgeBase, error)
	UpdateArticle(ctx context.Context, tenantID, articleID uuid.UUID, req UpdateArticleRequest) (*KnowledgeBase, error)
	DeleteArticle(ctx context.Context, tenantID, articleID uuid.UUID) error
	
	// Settings operations
	GetSettings(ctx context.Context, tenantID uuid.UUID) (*SupportSettings, error)
	UpdateSettings(ctx context.Context, tenantID uuid.UUID, req UpdateSettingsRequest) (*SupportSettings, error)
	
	// Analytics
	GetTicketStats(ctx context.Context, tenantID uuid.UUID, period string) (*TicketStats, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new support service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Request/Response DTOs
type CreateTicketRequest struct {
	Subject       string         `json:"subject" validate:"required"`
	Description   string         `json:"description" validate:"required"`
	Priority      TicketPriority `json:"priority"`
	Category      TicketCategory `json:"category"`
	CustomerEmail string         `json:"customer_email" validate:"required,email"`
	CustomerName  string         `json:"customer_name"`
	Tags          []string       `json:"tags"`
	Attachments   []string       `json:"attachments"`
}

type UpdateTicketRequest struct {
	Subject      *string         `json:"subject"`
	Description  *string         `json:"description"`
	Status       *TicketStatus   `json:"status"`
	Priority     *TicketPriority `json:"priority"`
	Category     *TicketCategory `json:"category"`
	AssignedToID *uuid.UUID      `json:"assigned_to_id"`
	Tags         []string        `json:"tags"`
}

type TicketFilter struct {
	Status       []TicketStatus   `json:"status"`
	Priority     []TicketPriority `json:"priority"`
	Category     []TicketCategory `json:"category"`
	AssignedToID *uuid.UUID       `json:"assigned_to_id"`
	Search       string           `json:"search"`
	Page         int              `json:"page"`
	Limit        int              `json:"limit"`
}

type AddMessageRequest struct {
	TicketID    uuid.UUID  `json:"ticket_id" validate:"required"`
	UserID      *uuid.UUID `json:"user_id"` // nil for customer messages
	Content     string     `json:"content" validate:"required"`
	IsInternal  bool       `json:"is_internal"`
	SenderName  string     `json:"sender_name" validate:"required"`
	SenderEmail string     `json:"sender_email" validate:"required,email"`
	Attachments []string   `json:"attachments"`
}

type CreateFAQRequest struct {
	Question    string   `json:"question" validate:"required"`
	Answer      string   `json:"answer" validate:"required"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	IsPublished bool     `json:"is_published"`
	Order       int      `json:"order"`
}

type UpdateFAQRequest struct {
	Question    *string  `json:"question"`
	Answer      *string  `json:"answer"`
	Category    *string  `json:"category"`
	Tags        []string `json:"tags"`
	IsPublished *bool    `json:"is_published"`
	Order       *int     `json:"order"`
}

type FAQFilter struct {
	Category    string `json:"category"`
	IsPublished *bool  `json:"is_published"`
	Search      string `json:"search"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

type CreateArticleRequest struct {
	Title           string   `json:"title" validate:"required"`
	Content         string   `json:"content" validate:"required"`
	Excerpt         string   `json:"excerpt"`
	Slug            string   `json:"slug" validate:"required"`
	Category        string   `json:"category"`
	Tags            []string `json:"tags"`
	IsPublished     bool     `json:"is_published"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
	Order           int      `json:"order"`
}

type UpdateArticleRequest struct {
	Title           *string  `json:"title"`
	Content         *string  `json:"content"`
	Excerpt         *string  `json:"excerpt"`
	Slug            *string  `json:"slug"`
	Category        *string  `json:"category"`
	Tags            []string `json:"tags"`
	IsPublished     *bool    `json:"is_published"`
	MetaTitle       *string  `json:"meta_title"`
	MetaDescription *string  `json:"meta_description"`
	Order           *int     `json:"order"`
}

type ArticleFilter struct {
	Category    string `json:"category"`
	IsPublished *bool  `json:"is_published"`
	Search      string `json:"search"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

type UpdateSettingsRequest struct {
	SupportEmail       *string `json:"support_email"`
	AutoReplyEnabled   *bool   `json:"auto_reply_enabled"`
	AutoReplyMessage   *string `json:"auto_reply_message"`
	BusinessHours      *string `json:"business_hours"`
	EmailNotifications *bool   `json:"email_notifications"`
	SlackWebhookURL    *string `json:"slack_webhook_url"`
}

type TicketStats struct {
	TotalTickets   int                        `json:"total_tickets"`
	OpenTickets    int                        `json:"open_tickets"`
	ResolvedTickets int                       `json:"resolved_tickets"`
	AvgResponseTime string                    `json:"avg_response_time"`
	TicketsByStatus map[TicketStatus]int      `json:"tickets_by_status"`
	TicketsByPriority map[TicketPriority]int  `json:"tickets_by_priority"`
}

// Implementation methods (TODO: implement business logic)
func (s *service) CreateTicket(ctx context.Context, req CreateTicketRequest) (*Ticket, error) {
	ticket := &Ticket{
		ID:            uuid.New(),
		Subject:       req.Subject,
		Description:   req.Description,
		Status:        StatusOpen,
		Priority:      req.Priority,
		Category:      req.Category,
		CustomerEmail: req.CustomerEmail,
		CustomerName:  req.CustomerName,
		Tags:          req.Tags,
		Attachments:   req.Attachments,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return s.repo.CreateTicket(ctx, ticket)
}

func (s *service) GetTicket(ctx context.Context, tenantID, ticketID uuid.UUID) (*Ticket, error) {
	return s.repo.GetTicket(ctx, tenantID, ticketID)
}

func (s *service) GetTickets(ctx context.Context, tenantID uuid.UUID, filter TicketFilter) ([]Ticket, error) {
	// Set default pagination
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Page == 0 {
		filter.Page = 1
	}

	return s.repo.GetTickets(ctx, tenantID, filter)
}

func (s *service) UpdateTicket(ctx context.Context, tenantID, ticketID uuid.UUID, req UpdateTicketRequest) (*Ticket, error) {
	updates := make(map[string]interface{})

	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.AssignedToID != nil {
		updates["assigned_to_id"] = *req.AssignedToID
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateTicket(ctx, tenantID, ticketID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetTicket(ctx, tenantID, ticketID)
}

func (s *service) DeleteTicket(ctx context.Context, tenantID, ticketID uuid.UUID) error {
	return s.repo.DeleteTicket(ctx, tenantID, ticketID)
}

func (s *service) AssignTicket(ctx context.Context, tenantID, ticketID, userID uuid.UUID) error {
	updates := map[string]interface{}{
		"assigned_to_id": userID,
		"updated_at":     time.Now(),
	}

	return s.repo.UpdateTicket(ctx, tenantID, ticketID, updates)
}

func (s *service) ResolveTicket(ctx context.Context, tenantID, ticketID uuid.UUID) error {
	// TODO: Implement ticket resolution logic
	now := time.Now()
	return s.repo.UpdateTicket(ctx, tenantID, ticketID, map[string]interface{}{
		"status": StatusResolved,
		"resolved_at": &now,
	})
}

func (s *service) AddMessage(ctx context.Context, req AddMessageRequest) (*TicketMessage, error) {
	message := &TicketMessage{
		ID:          uuid.New(),
		TicketID:    req.TicketID,
		UserID:      req.UserID,
		Content:     req.Content,
		IsInternal:  req.IsInternal,
		SenderName:  req.SenderName,
		SenderEmail: req.SenderEmail,
		Attachments: req.Attachments,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreateTicketMessage(ctx, message)
}

func (s *service) GetMessages(ctx context.Context, tenantID, ticketID uuid.UUID) ([]TicketMessage, error) {
	return s.repo.GetTicketMessages(ctx, tenantID, ticketID)
}

func (s *service) CreateFAQ(ctx context.Context, req CreateFAQRequest) (*FAQ, error) {
	faq := &FAQ{
		ID:          uuid.New(),
		Question:    req.Question,
		Answer:      req.Answer,
		Category:    req.Category,
		Tags:        req.Tags,
		IsPublished: req.IsPublished,
		Order:       req.Order,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.CreateFAQ(ctx, faq)
}

func (s *service) GetFAQ(ctx context.Context, tenantID, faqID uuid.UUID) (*FAQ, error) {
	return s.repo.GetFAQ(ctx, tenantID, faqID)
}

func (s *service) GetFAQs(ctx context.Context, tenantID uuid.UUID, filter FAQFilter) ([]FAQ, error) {
	// Set default pagination
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Page == 0 {
		filter.Page = 1
	}

	return s.repo.GetFAQs(ctx, tenantID, filter)
}

func (s *service) UpdateFAQ(ctx context.Context, tenantID, faqID uuid.UUID, req UpdateFAQRequest) (*FAQ, error) {
	updates := make(map[string]interface{})

	if req.Question != nil {
		updates["question"] = *req.Question
	}
	if req.Answer != nil {
		updates["answer"] = *req.Answer
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.IsPublished != nil {
		updates["is_published"] = *req.IsPublished
	}
	if req.Order != nil {
		updates["order"] = *req.Order
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateFAQ(ctx, tenantID, faqID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetFAQ(ctx, tenantID, faqID)
}

func (s *service) DeleteFAQ(ctx context.Context, tenantID, faqID uuid.UUID) error {
	return s.repo.DeleteFAQ(ctx, tenantID, faqID)
}

func (s *service) CreateArticle(ctx context.Context, req CreateArticleRequest) (*KnowledgeBase, error) {
	article := &KnowledgeBase{
		ID:              uuid.New(),
		Title:           req.Title,
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		Slug:            req.Slug,
		Category:        req.Category,
		Tags:            req.Tags,
		IsPublished:     req.IsPublished,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		Order:           req.Order,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return s.repo.CreateArticle(ctx, article)
}

func (s *service) GetArticle(ctx context.Context, tenantID uuid.UUID, slug string) (*KnowledgeBase, error) {
	return s.repo.GetArticleBySlug(ctx, tenantID, slug)
}

func (s *service) GetArticles(ctx context.Context, tenantID uuid.UUID, filter ArticleFilter) ([]KnowledgeBase, error) {
	// Set default pagination
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Page == 0 {
		filter.Page = 1
	}

	return s.repo.GetArticles(ctx, tenantID, filter)
}

func (s *service) UpdateArticle(ctx context.Context, tenantID, articleID uuid.UUID, req UpdateArticleRequest) (*KnowledgeBase, error) {
	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Excerpt != nil {
		updates["excerpt"] = *req.Excerpt
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.IsPublished != nil {
		updates["is_published"] = *req.IsPublished
	}
	if req.MetaTitle != nil {
		updates["meta_title"] = *req.MetaTitle
	}
	if req.MetaDescription != nil {
		updates["meta_description"] = *req.MetaDescription
	}
	if req.Order != nil {
		updates["order"] = *req.Order
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateArticle(ctx, tenantID, articleID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetArticle(ctx, tenantID, articleID)
}

func (s *service) DeleteArticle(ctx context.Context, tenantID, articleID uuid.UUID) error {
	return s.repo.DeleteArticle(ctx, tenantID, articleID)
}

func (s *service) GetSettings(ctx context.Context, tenantID uuid.UUID) (*SupportSettings, error) {
	return s.repo.GetSettings(ctx, tenantID)
}

func (s *service) UpdateSettings(ctx context.Context, tenantID uuid.UUID, req UpdateSettingsRequest) (*SupportSettings, error) {
	updates := make(map[string]interface{})

	if req.SupportEmail != nil {
		updates["support_email"] = *req.SupportEmail
	}
	if req.AutoReplyEnabled != nil {
		updates["auto_reply_enabled"] = *req.AutoReplyEnabled
	}
	if req.AutoReplyMessage != nil {
		updates["auto_reply_message"] = *req.AutoReplyMessage
	}
	if req.BusinessHours != nil {
		updates["business_hours"] = *req.BusinessHours
	}
	if req.EmailNotifications != nil {
		updates["email_notifications"] = *req.EmailNotifications
	}
	if req.SlackWebhookURL != nil {
		updates["slack_webhook_url"] = *req.SlackWebhookURL
	}

	updates["updated_at"] = time.Now()

	err := s.repo.UpdateSettings(ctx, tenantID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetSettings(ctx, tenantID)
}

func (s *service) GetTicketStats(ctx context.Context, tenantID uuid.UUID, period string) (*TicketStats, error) {
	// Calculate date range based on period
	now := time.Now()
	var startDate time.Time

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
	case "year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		// Default to last month
		startDate = now.AddDate(0, -1, 0)
	}

	return s.repo.GetTicketStats(ctx, tenantID, startDate, now)
}