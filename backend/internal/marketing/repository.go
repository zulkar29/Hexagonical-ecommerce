package marketing

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sharedErrors "ecommerce-saas/internal/shared/errors"
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
	GetSegmentCustomerCount(ctx context.Context, tenantID uuid.UUID, rules string) (int, error)

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
	GetMarketingOverview(ctx context.Context, tenantID uuid.UUID, period string) (*MarketingOverview, error)
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
	if err := r.db.WithContext(ctx).Create(campaign).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create campaign", err)
	}
	return nil
}

func (r *repository) GetCampaignByID(ctx context.Context, tenantID, campaignID uuid.UUID) (*Campaign, error) {
	var campaign Campaign
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", campaignID, tenantID).
		Preload("Emails").
		First(&campaign).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Campaign not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get campaign", err)
	}
	return &campaign, nil
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
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get campaigns", err)
	}
	return campaigns, nil
}

func (r *repository) UpdateCampaign(ctx context.Context, tenantID, campaignID uuid.UUID, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).
		Model(&Campaign{}).
		Where("id = ? AND tenant_id = ?", campaignID, tenantID).
		Updates(updates).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to update campaign", err)
	}
	return nil
}

func (r *repository) DeleteCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete campaign emails first
		if err := tx.Where("campaign_id = ?", campaignID).Delete(&CampaignEmail{}).Error; err != nil {
			return sharedErrors.NewInternalError("Failed to delete campaign emails", err)
		}

		// Delete campaign
		if err := tx.Where("id = ? AND tenant_id = ?", campaignID, tenantID).
			Delete(&Campaign{}).Error; err != nil {
			return sharedErrors.NewInternalError("Failed to delete campaign", err)
		}
		return nil
	})
}

// Email operations
func (r *repository) CreateCampaignEmail(ctx context.Context, email *CampaignEmail) error {
	if err := r.db.WithContext(ctx).Create(email).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create campaign email", err)
	}
	return nil
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
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Campaign email not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get campaign email", err)
	}
	return &email, nil
}

func (r *repository) UpdateCampaignEmail(ctx context.Context, emailID uuid.UUID, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).
		Model(&CampaignEmail{}).
		Where("id = ?", emailID).
		Updates(updates).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to update campaign email", err)
	}
	return nil
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
			return sharedErrors.NewInternalError("Failed to bulk create campaign emails", err)
		}
	}

	return nil
}

// Template operations
func (r *repository) CreateTemplate(ctx context.Context, template *EmailTemplate) error {
	if err := r.db.WithContext(ctx).Create(template).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create template", err)
	}
	return nil
}

func (r *repository) GetTemplateByID(ctx context.Context, tenantID, templateID uuid.UUID) (*EmailTemplate, error) {
	var template EmailTemplate
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", templateID, tenantID).
		First(&template).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Template not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get template", err)
	}
	return &template, nil
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
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get templates", err)
	}
	return templates, nil
}

func (r *repository) UpdateTemplate(ctx context.Context, tenantID, templateID uuid.UUID, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).
		Model(&EmailTemplate{}).
		Where("id = ? AND tenant_id = ?", templateID, tenantID).
		Updates(updates).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to update template", err)
	}
	return nil
}

func (r *repository) DeleteTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", templateID, tenantID).
		Delete(&EmailTemplate{}).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to delete template", err)
	}
	return nil
}

// Segment operations
func (r *repository) CreateSegment(ctx context.Context, segment *CustomerSegment) error {
	if err := r.db.WithContext(ctx).Create(segment).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create segment", err)
	}
	return nil
}

func (r *repository) GetSegmentByID(ctx context.Context, tenantID, segmentID uuid.UUID) (*CustomerSegment, error) {
	var segment CustomerSegment
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", segmentID, tenantID).
		First(&segment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Segment not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get segment", err)
	}
	return &segment, nil
}

func (r *repository) GetSegments(ctx context.Context, tenantID uuid.UUID) ([]CustomerSegment, error) {
	var segments []CustomerSegment
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&segments).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get segments", err)
	}
	return segments, nil
}

func (r *repository) UpdateSegment(ctx context.Context, tenantID, segmentID uuid.UUID, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).
		Model(&CustomerSegment{}).
		Where("id = ? AND tenant_id = ?", segmentID, tenantID).
		Updates(updates).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to update segment", err)
	}
	return nil
}

func (r *repository) DeleteSegment(ctx context.Context, tenantID, segmentID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", segmentID, tenantID).
		Delete(&CustomerSegment{}).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to delete segment", err)
	}
	return nil
}

func (r *repository) GetSegmentCustomerCount(ctx context.Context, tenantID uuid.UUID, rules string) (int, error) {
	// This is a simplified implementation - in a real system, you would parse the rules
	// and build a dynamic query based on customer attributes
	// For now, we'll return a placeholder count since the Customer model isn't defined yet
	// TODO: Implement proper customer segmentation logic when Customer model is available
	return 0, nil
}

// Newsletter operations
func (r *repository) CreateSubscriber(ctx context.Context, subscriber *NewsletterSubscriber) error {
	if err := r.db.WithContext(ctx).Create(subscriber).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create subscriber", err)
	}
	return nil
}

func (r *repository) GetSubscriberByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*NewsletterSubscriber, error) {
	var subscriber NewsletterSubscriber
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		First(&subscriber).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Subscriber not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get subscriber", err)
	}
	return &subscriber, nil
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
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get subscribers", err)
	}
	return subscribers, nil
}

func (r *repository) UpdateSubscriber(ctx context.Context, tenantID uuid.UUID, email string, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).
		Model(&NewsletterSubscriber{}).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		Updates(updates).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to update subscriber", err)
	}
	return nil
}

func (r *repository) DeleteSubscriber(ctx context.Context, tenantID uuid.UUID, email string) error {
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		Delete(&NewsletterSubscriber{}).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to delete subscriber", err)
	}
	return nil
}

// Abandoned cart operations
func (r *repository) CreateAbandonedCart(ctx context.Context, cart *AbandonedCart) error {
	if err := r.db.WithContext(ctx).Create(cart).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create abandoned cart", err)
	}
	return nil
}

func (r *repository) GetAbandonedCartByID(ctx context.Context, tenantID, cartID uuid.UUID) (*AbandonedCart, error) {
	var cart AbandonedCart
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND tenant_id = ?", cartID, tenantID).
		First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Abandoned cart not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get abandoned cart", err)
	}
	return &cart, nil
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
	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get abandoned carts", err)
	}
	return carts, nil
}

func (r *repository) UpdateAbandonedCart(ctx context.Context, tenantID, cartID uuid.UUID, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).
		Model(&AbandonedCart{}).
		Where("cart_id = ? AND tenant_id = ?", cartID, tenantID).
		Updates(updates).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to update abandoned cart", err)
	}
	return nil
}

func (r *repository) DeleteAbandonedCart(ctx context.Context, tenantID, cartID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("cart_id = ? AND tenant_id = ?", cartID, tenantID).
		Delete(&AbandonedCart{}).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to delete abandoned cart", err)
	}
	return nil
}

// Settings operations
func (r *repository) GetSettings(ctx context.Context, tenantID uuid.UUID) (*MarketingSettings, error) {
	var settings MarketingSettings
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Marketing settings not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get marketing settings", err)
	}
	return &settings, nil
}

func (r *repository) CreateSettings(ctx context.Context, settings *MarketingSettings) error {
	if err := r.db.WithContext(ctx).Create(settings).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create marketing settings", err)
	}
	return nil
}

func (r *repository) UpdateSettings(ctx context.Context, tenantID uuid.UUID, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).
		Model(&MarketingSettings{}).
		Where("tenant_id = ?", tenantID).
		Updates(updates).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to update marketing settings", err)
	}
	return nil
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
		CampaignID:        campaignID,
		SentCount:         campaign.SentCount,
		DeliveredCount:    campaign.DeliveredCount,
		OpenedCount:       campaign.OpenedCount,
		ClickedCount:      campaign.ClickedCount,
		BouncedCount:      campaign.BouncedCount,
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
	if err != nil {
		return 0, sharedErrors.NewInternalError("Failed to get subscriber count", err)
	}
	return count, nil
}

func (r *repository) GetAbandonedCartStats(ctx context.Context, tenantID uuid.UUID) (total int64, recovered int64, err error) {
	// Get total abandoned carts
	if err := r.db.WithContext(ctx).Model(&AbandonedCart{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return 0, 0, sharedErrors.NewInternalError("failed to count abandoned carts", err)
	}
	
	// Get recovered carts
	if err := r.db.WithContext(ctx).Model(&AbandonedCart{}).Where("tenant_id = ? AND is_recovered = ?", tenantID, true).Count(&recovered).Error; err != nil {
		return 0, 0, sharedErrors.NewInternalError("failed to count recovered carts", err)
	}
	
	return total, recovered, nil
}

func (r *repository) GetMarketingOverview(ctx context.Context, tenantID uuid.UUID, period string) (*MarketingOverview, error) {
	var overview MarketingOverview
	
	// Get campaign statistics
	var totalCampaigns, activeCampaigns int64
	if err := r.db.WithContext(ctx).Model(&Campaign{}).Where("tenant_id = ?", tenantID).Count(&totalCampaigns).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to count total campaigns", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&Campaign{}).Where("tenant_id = ? AND status IN ?", tenantID, []string{"running", "scheduled"}).Count(&activeCampaigns).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to count active campaigns", err)
	}
	
	// Get subscriber statistics
	var totalSubscribers, activeSubscribers int64
	if err := r.db.WithContext(ctx).Model(&NewsletterSubscriber{}).Where("tenant_id = ?", tenantID).Count(&totalSubscribers).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to count total subscribers", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&NewsletterSubscriber{}).Where("tenant_id = ? AND status = ?", tenantID, "active").Count(&activeSubscribers).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to count active subscribers", err)
	}
	
	// Get email statistics from campaigns
	var emailsSent, emailsOpened, emailsClicked int64
	if err := r.db.WithContext(ctx).Model(&Campaign{}).Where("tenant_id = ?", tenantID).Select("COALESCE(SUM(sent_count), 0)").Scan(&emailsSent).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to calculate emails sent", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&Campaign{}).Where("tenant_id = ?", tenantID).Select("COALESCE(SUM(opened_count), 0)").Scan(&emailsOpened).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to calculate emails opened", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&Campaign{}).Where("tenant_id = ?", tenantID).Select("COALESCE(SUM(clicked_count), 0)").Scan(&emailsClicked).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to calculate emails clicked", err)
	}
	
	// Get abandoned cart statistics
	var abandonedCarts, recoveredCarts int64
	if err := r.db.WithContext(ctx).Model(&AbandonedCart{}).Where("tenant_id = ?", tenantID).Count(&abandonedCarts).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to count abandoned carts", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&AbandonedCart{}).Where("tenant_id = ? AND is_recovered = ?", tenantID, true).Count(&recoveredCarts).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to count recovered carts", err)
	}
	
	// Calculate rates
	openRate := float64(0)
	clickRate := float64(0)
	recoveryRate := float64(0)
	
	if emailsSent > 0 {
		openRate = float64(emailsOpened) / float64(emailsSent) * 100
		clickRate = float64(emailsClicked) / float64(emailsSent) * 100
	}
	
	if abandonedCarts > 0 {
		recoveryRate = float64(recoveredCarts) / float64(abandonedCarts) * 100
	}
	
	overview = MarketingOverview{
		TotalCampaigns:    int(totalCampaigns),
		ActiveCampaigns:   int(activeCampaigns),
		TotalSubscribers:  int(totalSubscribers),
		ActiveSubscribers: int(activeSubscribers),
		EmailsSent:        int(emailsSent),
		EmailsOpened:      int(emailsOpened),
		EmailsClicked:     int(emailsClicked),
		OpenRate:          openRate,
		ClickRate:         clickRate,
		AbandonedCarts:    int(abandonedCarts),
		RecoveredCarts:    int(recoveredCarts),
		RecoveryRate:      recoveryRate,
	}
	
	return &overview, nil
}
