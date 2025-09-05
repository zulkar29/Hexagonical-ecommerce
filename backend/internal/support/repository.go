package support

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the support repository interface
type Repository interface {
	// Ticket operations
	CreateTicket(ctx context.Context, ticket *Ticket) error
	GetTicket(ctx context.Context, tenantID, ticketID uuid.UUID) (*Ticket, error)
	GetTicketByID(ctx context.Context, tenantID, ticketID uuid.UUID) (*Ticket, error)
	GetTickets(ctx context.Context, tenantID uuid.UUID, filter TicketFilter) ([]Ticket, error)
	UpdateTicket(ctx context.Context, tenantID, ticketID uuid.UUID, updates map[string]interface{}) error
	DeleteTicket(ctx context.Context, tenantID, ticketID uuid.UUID) error
	
	// Ticket message operations
	CreateMessage(ctx context.Context, message *TicketMessage) error
	CreateTicketMessage(ctx context.Context, message *TicketMessage) error
	GetMessagesByTicketID(ctx context.Context, tenantID, ticketID uuid.UUID) ([]TicketMessage, error)
	
	// FAQ operations
	CreateFAQ(ctx context.Context, faq *FAQ) error
	GetFAQByID(ctx context.Context, tenantID, faqID uuid.UUID) (*FAQ, error)
	GetFAQs(ctx context.Context, tenantID uuid.UUID, filter FAQFilter) ([]FAQ, error)
	UpdateFAQ(ctx context.Context, tenantID, faqID uuid.UUID, updates map[string]interface{}) error
	DeleteFAQ(ctx context.Context, tenantID, faqID uuid.UUID) error
	
	// Knowledge base operations
	CreateArticle(ctx context.Context, article *KnowledgeBase) error
	GetArticleBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*KnowledgeBase, error)
	GetArticleByID(ctx context.Context, tenantID, articleID uuid.UUID) (*KnowledgeBase, error)
	GetArticles(ctx context.Context, tenantID uuid.UUID, filter ArticleFilter) ([]KnowledgeBase, error)
	UpdateArticle(ctx context.Context, tenantID, articleID uuid.UUID, updates map[string]interface{}) error
	DeleteArticle(ctx context.Context, tenantID, articleID uuid.UUID) error
	
	// Settings operations
	GetSettings(ctx context.Context, tenantID uuid.UUID) (*SupportSettings, error)
	CreateSettings(ctx context.Context, settings *SupportSettings) error
	UpdateSettings(ctx context.Context, tenantID uuid.UUID, updates map[string]interface{}) error
	
	// Analytics
	GetTicketCount(ctx context.Context, tenantID uuid.UUID, status *TicketStatus) (int64, error)
	GetTicketsByStatusCount(ctx context.Context, tenantID uuid.UUID) (map[TicketStatus]int, error)
	GetTicketsByPriorityCount(ctx context.Context, tenantID uuid.UUID) (map[TicketPriority]int, error)
	GetTicketStats(ctx context.Context, tenantID uuid.UUID) (*TicketStats, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new support repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Ticket operations
func (r *repository) CreateTicket(ctx context.Context, ticket *Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *repository) GetTicketByID(ctx context.Context, tenantID, ticketID uuid.UUID) (*Ticket, error) {
	var ticket Ticket
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", ticketID, tenantID).
		Preload("Messages").
		First(&ticket).Error
	return &ticket, err
}

func (r *repository) GetTickets(ctx context.Context, tenantID uuid.UUID, filter TicketFilter) ([]Ticket, error) {
	var tickets []Ticket
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}
	
	if len(filter.Priority) > 0 {
		query = query.Where("priority IN ?", filter.Priority)
	}
	
	if len(filter.Category) > 0 {
		query = query.Where("category IN ?", filter.Category)
	}
	
	if filter.AssignedToID != nil {
		if *filter.AssignedToID == uuid.Nil {
			query = query.Where("assigned_to_id IS NULL")
		} else {
			query = query.Where("assigned_to_id = ?", *filter.AssignedToID)
		}
	}
	
	if filter.Search != "" {
		query = query.Where("subject ILIKE ? OR description ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

func (r *repository) UpdateTicket(ctx context.Context, tenantID, ticketID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&Ticket{}).
		Where("id = ? AND tenant_id = ?", ticketID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteTicket(ctx context.Context, tenantID, ticketID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", ticketID, tenantID).
		Delete(&Ticket{}).Error
}

// Ticket message operations
func (r *repository) CreateMessage(ctx context.Context, message *TicketMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *repository) GetMessagesByTicketID(ctx context.Context, tenantID, ticketID uuid.UUID) ([]TicketMessage, error) {
	var messages []TicketMessage
	
	// Verify ticket belongs to tenant first
	var count int64
	r.db.WithContext(ctx).Model(&Ticket{}).
		Where("id = ? AND tenant_id = ?", ticketID, tenantID).
		Count(&count)
	
	if count == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	
	err := r.db.WithContext(ctx).
		Where("ticket_id = ?", ticketID).
		Order("created_at ASC").
		Find(&messages).Error
	
	return messages, err
}

// FAQ operations
func (r *repository) CreateFAQ(ctx context.Context, faq *FAQ) error {
	return r.db.WithContext(ctx).Create(faq).Error
}

func (r *repository) GetFAQByID(ctx context.Context, tenantID, faqID uuid.UUID) (*FAQ, error) {
	var faq FAQ
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", faqID, tenantID).
		First(&faq).Error
	return &faq, err
}

func (r *repository) GetFAQs(ctx context.Context, tenantID uuid.UUID, filter FAQFilter) ([]FAQ, error) {
	var faqs []FAQ
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	
	if filter.IsPublished != nil {
		query = query.Where("is_published = ?", *filter.IsPublished)
	}
	
	if filter.Search != "" {
		query = query.Where("question ILIKE ? OR answer ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Order("`order` ASC, created_at DESC").Find(&faqs).Error
	return faqs, err
}

func (r *repository) UpdateFAQ(ctx context.Context, tenantID, faqID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&FAQ{}).
		Where("id = ? AND tenant_id = ?", faqID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteFAQ(ctx context.Context, tenantID, faqID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", faqID, tenantID).
		Delete(&FAQ{}).Error
}

// Knowledge base operations
func (r *repository) CreateArticle(ctx context.Context, article *KnowledgeBase) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *repository) GetArticleBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*KnowledgeBase, error) {
	var article KnowledgeBase
	err := r.db.WithContext(ctx).
		Where("slug = ? AND tenant_id = ?", slug, tenantID).
		First(&article).Error
	return &article, err
}

func (r *repository) GetArticleByID(ctx context.Context, tenantID, articleID uuid.UUID) (*KnowledgeBase, error) {
	var article KnowledgeBase
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", articleID, tenantID).
		First(&article).Error
	return &article, err
}

func (r *repository) GetArticles(ctx context.Context, tenantID uuid.UUID, filter ArticleFilter) ([]KnowledgeBase, error) {
	var articles []KnowledgeBase
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	
	if filter.IsPublished != nil {
		query = query.Where("is_published = ?", *filter.IsPublished)
	}
	
	if filter.Search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Order("`order` ASC, created_at DESC").Find(&articles).Error
	return articles, err
}

func (r *repository) UpdateArticle(ctx context.Context, tenantID, articleID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&KnowledgeBase{}).
		Where("id = ? AND tenant_id = ?", articleID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteArticle(ctx context.Context, tenantID, articleID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", articleID, tenantID).
		Delete(&KnowledgeBase{}).Error
}

// Settings operations
func (r *repository) GetSettings(ctx context.Context, tenantID uuid.UUID) (*SupportSettings, error) {
	var settings SupportSettings
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&settings).Error
	return &settings, err
}

func (r *repository) CreateSettings(ctx context.Context, settings *SupportSettings) error {
	return r.db.WithContext(ctx).Create(settings).Error
}

func (r *repository) UpdateSettings(ctx context.Context, tenantID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&SupportSettings{}).
		Where("tenant_id = ?", tenantID).
		Updates(updates).Error
}

// Analytics
func (r *repository) GetTicketCount(ctx context.Context, tenantID uuid.UUID, status *TicketStatus) (int64, error) {
	query := r.db.WithContext(ctx).Model(&Ticket{}).Where("tenant_id = ?", tenantID)
	
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (r *repository) GetTicketsByStatusCount(ctx context.Context, tenantID uuid.UUID) (map[TicketStatus]int, error) {
	type statusCount struct {
		Status TicketStatus `json:"status"`
		Count  int          `json:"count"`
	}
	
	var results []statusCount
	err := r.db.WithContext(ctx).
		Model(&Ticket{}).
		Select("status, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("status").
		Find(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	counts := make(map[TicketStatus]int)
	for _, result := range results {
		counts[result.Status] = result.Count
	}
	
	return counts, nil
}

func (r *repository) GetTicketsByPriorityCount(ctx context.Context, tenantID uuid.UUID) (map[TicketPriority]int, error) {
	type priorityCount struct {
		Priority TicketPriority `json:"priority"`
		Count    int            `json:"count"`
	}
	
	var results []priorityCount
	err := r.db.WithContext(ctx).
		Model(&Ticket{}).
		Select("priority, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("priority").
		Find(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	counts := make(map[TicketPriority]int)
	for _, result := range results {
		counts[result.Priority] = result.Count
	}
	
	return counts, nil
}

// GetTicket is an alias for GetTicketByID
func (r *repository) GetTicket(ctx context.Context, tenantID, ticketID uuid.UUID) (*Ticket, error) {
	return r.GetTicketByID(ctx, tenantID, ticketID)
}

// CreateTicketMessage is an alias for CreateMessage
func (r *repository) CreateTicketMessage(ctx context.Context, message *TicketMessage) error {
	return r.CreateMessage(ctx, message)
}

// GetTicketStats returns ticket statistics
func (r *repository) GetTicketStats(ctx context.Context, tenantID uuid.UUID) (*TicketStats, error) {
	stats := &TicketStats{}
	
	// Get total tickets
	totalCount, err := r.GetTicketCount(ctx, tenantID, nil)
	if err != nil {
		return nil, err
	}
	stats.TotalTickets = int(totalCount)
	
	// Get tickets by status
	statusOpen := StatusOpen
	openCount, err := r.GetTicketCount(ctx, tenantID, &statusOpen)
	if err != nil {
		return nil, err
	}
	stats.OpenTickets = int(openCount)
	
	statusResolved := StatusResolved
	resolvedCount, err := r.GetTicketCount(ctx, tenantID, &statusResolved)
	if err != nil {
		return nil, err
	}
	stats.ResolvedTickets = int(resolvedCount)
	
	statusClosed := StatusClosed
	closedCount, err := r.GetTicketCount(ctx, tenantID, &statusClosed)
	if err != nil {
		return nil, err
	}
	stats.ClosedTickets = int(closedCount)
	
	return stats, nil
}