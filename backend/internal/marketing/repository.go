package marketing

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the marketing repository interface
type Repository interface {
	// Campaign operations
	CreateCampaign(ctx context.Context, campaign *Campaign) error
	GetCampaignByID(ctx context.Context, tenantID, campaignID uuid.UUID) (*Campaign, error)
	GetCampaigns(ctx context.Context, tenantID uuid.UUID, filter CampaignFilter) ([]Campaign, error)
	UpdateCampaign(ctx context.Context, tenantID, campaignID uuid.UUID, updates map[string]interface{}) error
	DeleteCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error
	
	// Email operations
	CreateCampaignEmail(ctx context.Context, email *CampaignEmail) error
	GetCampaignEmails(ctx context.Context, tenantID, campaignID uuid.UUID, filter EmailFilter) ([]CampaignEmail, error)
	GetCampaignEmailByID(ctx context.Context, emailID uuid.UUID) (*CampaignEmail, error)
	UpdateCampaignEmail(ctx context.Context, emailID uuid.UUID, updates map[string]interface{}) error
	BulkCreateCampaignEmails(ctx context.Context, emails []CampaignEmail) error
	
	// Template operations
	CreateTemplate(ctx context.Context, template *EmailTemplate) error
	GetTemplateByID(ctx context.Context, tenantID, templateID uuid.UUID) (*EmailTemplate, error)
	GetTemplates(ctx context.Context, tenantID uuid.UUID, filter TemplateFilter) ([]EmailTemplate, error)
	UpdateTemplate(ctx context.Context, tenantID, templateID uuid.UUID, updates map[string]interface{}) error
	DeleteTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error
	
	// Segment operations
	CreateSegment(ctx context.Context, segment *CustomerSegment) error
	GetSegmentByID(ctx context.Context, tenantID, segmentID uuid.UUID) (*CustomerSegment, error)
	GetSegments(ctx context.Context, tenantID uuid.UUID) ([]CustomerSegment, error)
	UpdateSegment(ctx context.Context, tenantID, segmentID uuid.UUID, updates map[string]interface{}) error
	DeleteSegment(ctx context.Context, tenantID, segmentID uuid.UUID) error
	
	// Newsletter operations
	CreateSubscriber(ctx context.Context, subscriber *NewsletterSubscriber) error
	GetSubscriberByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*NewsletterSubscriber, error)
	GetSubscribers(ctx context.Context, tenantID uuid.UUID, filter SubscriberFilter) ([]NewsletterSubscriber, error)
	UpdateSubscriber(ctx context.Context, tenantID uuid.UUID, email string, updates map[string]interface{}) error
	DeleteSubscriber(ctx context.Context, tenantID uuid.UUID, email string) error
	
	// Abandoned cart operations
	CreateAbandonedCart(ctx context.Context, cart *AbandonedCart) error
	GetAbandonedCartByID(ctx context.Context, tenantID, cartID uuid.UUID) (*AbandonedCart, error)
	GetAbandonedCarts(ctx context.Context, tenantID uuid.UUID, filter AbandonedCartFilter) ([]AbandonedCart, error)
	UpdateAbandonedCart(ctx context.Context, tenantID, cartID uuid.UUID, updates map[string]interface{}) error
	DeleteAbandonedCart(ctx context.Context, tenantID, cartID uuid.UUID) error
	
	// Settings operations
	GetSettings(ctx context.Context, tenantID uuid.UUID) (*MarketingSettings, error)
	CreateSettings(ctx context.Context, settings *MarketingSettings) error
	UpdateSettings(ctx context.Context, tenantID uuid.UUID, updates map[string]interface{}) error
	
	// Analytics queries
	GetCampaignEmailStats(ctx context.Context, tenantID, campaignID uuid.UUID) (*CampaignStats, error)
	GetSubscriberCount(ctx context.Context, tenantID uuid.UUID, status string) (int64, error)
	GetAbandonedCartStats(ctx context.Context, tenantID uuid.UUID) (total int64, recovered int64, err error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new marketing repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Campaign operations
func (r *repository) CreateCampaign(ctx context.Context, campaign *Campaign) error {
	return r.db.WithContext(ctx).Create(campaign).Error
}

func (r *repository) GetCampaignByID(ctx context.Context, tenantID, campaignID uuid.UUID) (*Campaign, error) {
	var campaign Campaign
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", campaignID, tenantID).
		Preload("Emails").
		First(&campaign).Error
	return &campaign, err
}

func (r *repository) GetCampaigns(ctx context.Context, tenantID uuid.UUID, filter CampaignFilter) ([]Campaign, error) {
	var campaigns []Campaign
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if len(filter.Type) > 0 {
		query = query.Where("type IN ?", filter.Type)
	}
	
	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}
	
	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Order("created_at DESC").Find(&campaigns).Error
	return campaigns, err
}

func (r *repository) UpdateCampaign(ctx context.Context, tenantID, campaignID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&Campaign{}).
		Where("id = ? AND tenant_id = ?", campaignID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete campaign emails first
		if err := tx.Where("campaign_id = ?", campaignID).Delete(&CampaignEmail{}).Error; err != nil {
			return err
		}
		
		// Delete campaign
		return tx.Where("id = ? AND tenant_id = ?", campaignID, tenantID).
			Delete(&Campaign{}).Error
	})
}

// Email operations
func (r *repository) CreateCampaignEmail(ctx context.Context, email *CampaignEmail) error {
	return r.db.WithContext(ctx).Create(email).Error
}

func (r *repository) GetCampaignEmails(ctx context.Context, tenantID, campaignID uuid.UUID, filter EmailFilter) ([]CampaignEmail, error) {
	var emails []CampaignEmail
	
	// Verify campaign belongs to tenant first
	var count int64
	r.db.WithContext(ctx).Model(&Campaign{}).
		Where("id = ? AND tenant_id = ?", campaignID, tenantID).
		Count(&count)
	
	if count == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	
	query := r.db.WithContext(ctx).Where("campaign_id = ?", campaignID)
	
	// Apply filters
	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}
	
	if filter.Search != "" {
		query = query.Where("recipient_email ILIKE ? OR recipient_name ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Order("created_at DESC").Find(&emails).Error
	return emails, err
}

func (r *repository) GetCampaignEmailByID(ctx context.Context, emailID uuid.UUID) (*CampaignEmail, error) {
	var email CampaignEmail
	err := r.db.WithContext(ctx).
		Where("id = ?", emailID).
		First(&email).Error
	return &email, err
}

func (r *repository) UpdateCampaignEmail(ctx context.Context, emailID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&CampaignEmail{}).
		Where("id = ?", emailID).
		Updates(updates).Error
}

func (r *repository) BulkCreateCampaignEmails(ctx context.Context, emails []CampaignEmail) error {
	if len(emails) == 0 {
		return nil
	}
	
	// Use batch insert for better performance
	batchSize := 1000
	for i := 0; i < len(emails); i += batchSize {
		end := i + batchSize
		if end > len(emails) {
			end = len(emails)
		}
		
		if err := r.db.WithContext(ctx).Create(emails[i:end]).Error; err != nil {
			return err
		}
	}
	
	return nil
}

// Template operations
func (r *repository) CreateTemplate(ctx context.Context, template *EmailTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *repository) GetTemplateByID(ctx context.Context, tenantID, templateID uuid.UUID) (*EmailTemplate, error) {
	var template EmailTemplate
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", templateID, tenantID).
		First(&template).Error
	return &template, err
}

func (r *repository) GetTemplates(ctx context.Context, tenantID uuid.UUID, filter TemplateFilter) ([]EmailTemplate, error) {
	var templates []EmailTemplate
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	
	if len(filter.Type) > 0 {
		query = query.Where("type IN ?", filter.Type)
	}
	
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	
	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Order("created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *repository) UpdateTemplate(ctx context.Context, tenantID, templateID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&EmailTemplate{}).
		Where("id = ? AND tenant_id = ?", templateID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", templateID, tenantID).
		Delete(&EmailTemplate{}).Error
}

// Segment operations
func (r *repository) CreateSegment(ctx context.Context, segment *CustomerSegment) error {
	return r.db.WithContext(ctx).Create(segment).Error
}

func (r *repository) GetSegmentByID(ctx context.Context, tenantID, segmentID uuid.UUID) (*CustomerSegment, error) {
	var segment CustomerSegment
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", segmentID, tenantID).
		First(&segment).Error
	return &segment, err
}

func (r *repository) GetSegments(ctx context.Context, tenantID uuid.UUID) ([]CustomerSegment, error) {
	var segments []CustomerSegment
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&segments).Error
	return segments, err
}

func (r *repository) UpdateSegment(ctx context.Context, tenantID, segmentID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&CustomerSegment{}).
		Where("id = ? AND tenant_id = ?", segmentID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteSegment(ctx context.Context, tenantID, segmentID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", segmentID, tenantID).
		Delete(&CustomerSegment{}).Error
}

// Newsletter operations
func (r *repository) CreateSubscriber(ctx context.Context, subscriber *NewsletterSubscriber) error {
	return r.db.WithContext(ctx).Create(subscriber).Error
}

func (r *repository) GetSubscriberByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*NewsletterSubscriber, error) {
	var subscriber NewsletterSubscriber
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		First(&subscriber).Error
	return &subscriber, err
}

func (r *repository) GetSubscribers(ctx context.Context, tenantID uuid.UUID, filter SubscriberFilter) ([]NewsletterSubscriber, error) {
	var subscribers []NewsletterSubscriber
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}
	
	if len(filter.Tags) > 0 {
		query = query.Where("tags && ?", filter.Tags)
	}
	
	if filter.Search != "" {
		query = query.Where("email ILIKE ? OR name ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	if filter.StartDate != nil {
		query = query.Where("subscribed_at >= ?", filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("subscribed_at <= ?", filter.EndDate)
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Order("subscribed_at DESC").Find(&subscribers).Error
	return subscribers, err
}

func (r *repository) UpdateSubscriber(ctx context.Context, tenantID uuid.UUID, email string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&NewsletterSubscriber{}).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		Updates(updates).Error
}

func (r *repository) DeleteSubscriber(ctx context.Context, tenantID uuid.UUID, email string) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		Delete(&NewsletterSubscriber{}).Error
}

// Abandoned cart operations
func (r *repository) CreateAbandonedCart(ctx context.Context, cart *AbandonedCart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *repository) GetAbandonedCartByID(ctx context.Context, tenantID, cartID uuid.UUID) (*AbandonedCart, error) {
	var cart AbandonedCart
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND tenant_id = ?", cartID, tenantID).
		First(&cart).Error
	return &cart, err
}

func (r *repository) GetAbandonedCarts(ctx context.Context, tenantID uuid.UUID, filter AbandonedCartFilter) ([]AbandonedCart, error) {
	var carts []AbandonedCart
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filter.IsRecovered != nil {
		query = query.Where("is_recovered = ?", *filter.IsRecovered)
	}
	
	if filter.MinValue != nil {
		query = query.Where("cart_value >= ?", *filter.MinValue)
	}
	
	if filter.MaxValue != nil {
		query = query.Where("cart_value <= ?", *filter.MaxValue)
	}
	
	if filter.Search != "" {
		query = query.Where("customer_email ILIKE ? OR customer_name ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	if filter.StartDate != nil {
		query = query.Where("abandoned_at >= ?", filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("abandoned_at <= ?", filter.EndDate)
	}
	
	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}
	
	err := query.Order("abandoned_at DESC").Find(&carts).Error
	return carts, err
}

func (r *repository) UpdateAbandonedCart(ctx context.Context, tenantID, cartID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&AbandonedCart{}).
		Where("cart_id = ? AND tenant_id = ?", cartID, tenantID).
		Updates(updates).Error
}

func (r *repository) DeleteAbandonedCart(ctx context.Context, tenantID, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("cart_id = ? AND tenant_id = ?", cartID, tenantID).
		Delete(&AbandonedCart{}).Error
}

// Settings operations
func (r *repository) GetSettings(ctx context.Context, tenantID uuid.UUID) (*MarketingSettings, error) {
	var settings MarketingSettings
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&settings).Error
	return &settings, err
}

func (r *repository) CreateSettings(ctx context.Context, settings *MarketingSettings) error {
	return r.db.WithContext(ctx).Create(settings).Error
}

func (r *repository) UpdateSettings(ctx context.Context, tenantID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&MarketingSettings{}).
		Where("tenant_id = ?", tenantID).
		Updates(updates).Error
}

// Analytics queries
func (r *repository) GetCampaignEmailStats(ctx context.Context, tenantID, campaignID uuid.UUID) (*CampaignStats, error) {
	// Verify campaign belongs to tenant
	var campaign Campaign
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", campaignID, tenantID).
		First(&campaign).Error
	
	if err != nil {
		return nil, err
	}
	
	stats := &CampaignStats{
		CampaignID:     campaignID,
		SentCount:      campaign.SentCount,
		DeliveredCount: campaign.DeliveredCount,
		OpenedCount:    campaign.OpenedCount,
		ClickedCount:   campaign.ClickedCount,
		BouncedCount:   campaign.BouncedCount,
		UnsubscribedCount: campaign.UnsubscribedCount,
	}
	
	// Calculate rates
	if stats.SentCount > 0 {
		stats.OpenRate = float64(stats.OpenedCount) / float64(stats.SentCount) * 100
		stats.ClickRate = float64(stats.ClickedCount) / float64(stats.SentCount) * 100
		stats.BounceRate = float64(stats.BouncedCount) / float64(stats.SentCount) * 100
	}
	
	return stats, nil
}

func (r *repository) GetSubscriberCount(ctx context.Context, tenantID uuid.UUID, status string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&NewsletterSubscriber{}).Where("tenant_id = ?", tenantID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Count(&count).Error
	return count, err
}

func (r *repository) GetAbandonedCartStats(ctx context.Context, tenantID uuid.UUID) (total int64, recovered int64, err error) {
	// Get total abandoned carts
	err = r.db.WithContext(ctx).
		Model(&AbandonedCart{}).
		Where("tenant_id = ?", tenantID).
		Count(&total).Error
	
	if err != nil {
		return 0, 0, err
	}
	
	// Get recovered carts
	err = r.db.WithContext(ctx).
		Model(&AbandonedCart{}).
		Where("tenant_id = ? AND is_recovered = ?", tenantID, true).
		Count(&recovered).Error
	
	return total, recovered, err
}