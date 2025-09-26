package social

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new social repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Social integrations

// CreateIntegration creates a new social integration
func (r *repository) CreateIntegration(ctx context.Context, integration *SocialIntegration) error {
	return r.db.WithContext(ctx).Create(integration).Error
}

// GetIntegrationByTenant gets integration by tenant and platform
func (r *repository) GetIntegrationByTenant(ctx context.Context, tenantID uuid.UUID, platform Platform) (*SocialIntegration, error) {
	var integration SocialIntegration
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND platform = ?", tenantID, platform).
		First(&integration).Error
	if err != nil {
		return nil, err
	}
	return &integration, nil
}

// GetIntegrationsByTenant gets all integrations for a tenant
func (r *repository) GetIntegrationsByTenant(ctx context.Context, tenantID uuid.UUID) ([]SocialIntegration, error) {
	var integrations []SocialIntegration
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Find(&integrations).Error
	return integrations, err
}

// UpdateIntegration updates a social integration
func (r *repository) UpdateIntegration(ctx context.Context, integration *SocialIntegration) error {
	return r.db.WithContext(ctx).Save(integration).Error
}

// DeleteIntegration deletes a social integration
func (r *repository) DeleteIntegration(ctx context.Context, tenantID uuid.UUID, platform Platform) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND platform = ?", tenantID, platform).
		Delete(&SocialIntegration{}).Error
}

// Social products

// CreateSocialProduct creates a new social product
func (r *repository) CreateSocialProduct(ctx context.Context, socialProduct *SocialProduct) error {
	return r.db.WithContext(ctx).Create(socialProduct).Error
}

// GetSocialProduct gets a social product by tenant, product, and platform
func (r *repository) GetSocialProduct(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID, platform Platform) (*SocialProduct, error) {
	var socialProduct SocialProduct
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND product_id = ? AND platform = ?", tenantID, productID, platform).
		First(&socialProduct).Error
	if err != nil {
		return nil, err
	}
	return &socialProduct, nil
}

// GetSocialProductsByTenant gets all social products for a tenant, optionally filtered by platform
func (r *repository) GetSocialProductsByTenant(ctx context.Context, tenantID uuid.UUID, platform *Platform) ([]SocialProduct, error) {
	var socialProducts []SocialProduct
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	if platform != nil {
		query = query.Where("platform = ?", *platform)
	}

	err := query.Find(&socialProducts).Error
	return socialProducts, err
}

// UpdateSocialProduct updates a social product
func (r *repository) UpdateSocialProduct(ctx context.Context, socialProduct *SocialProduct) error {
	return r.db.WithContext(ctx).Save(socialProduct).Error
}

// DeleteSocialProduct deletes a social product
func (r *repository) DeleteSocialProduct(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID, platform Platform) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND product_id = ? AND platform = ?", tenantID, productID, platform).
		Delete(&SocialProduct{}).Error
}

// Analytics

// CreateAnalytics creates analytics record
func (r *repository) CreateAnalytics(ctx context.Context, analytics *SocialAnalytics) error {
	return r.db.WithContext(ctx).Create(analytics).Error
}

// GetAnalytics gets analytics data for a tenant, platform, and date range
func (r *repository) GetAnalytics(ctx context.Context, tenantID uuid.UUID, platform Platform, dateFrom, dateTo time.Time) ([]SocialAnalytics, error) {
	var analytics []SocialAnalytics
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND platform = ? AND date >= ? AND date <= ?", tenantID, platform, dateFrom, dateTo).
		Order("date ASC").
		Find(&analytics).Error
	return analytics, err
}

// Additional helper methods

// GetActiveIntegrations gets all active integrations that need syncing
func (r *repository) GetActiveIntegrations(ctx context.Context) ([]SocialIntegration, error) {
	var integrations []SocialIntegration
	err := r.db.WithContext(ctx).
		Where("status = ? AND auto_sync = true AND (next_sync_at IS NULL OR next_sync_at <= ?)", StatusConnected, time.Now()).
		Find(&integrations).Error
	return integrations, err
}

// GetProductsNeedingSync gets products that need syncing for a platform
func (r *repository) GetProductsNeedingSync(ctx context.Context, tenantID uuid.UUID, platform Platform) ([]SocialProduct, error) {
	var socialProducts []SocialProduct
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND platform = ? AND is_enabled = true AND (status = ? OR status = ?)",
			tenantID, platform, SyncPending, SyncFailed).
		Where("sync_attempts < ?", 3).
		Find(&socialProducts).Error
	return socialProducts, err
}

// GetSyncStats gets sync statistics for a tenant and platform
func (r *repository) GetSyncStats(ctx context.Context, tenantID uuid.UUID, platform Platform) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Count by status
	statusCounts := []struct {
		Status SyncStatus
		Count  int64
	}{}

	err := r.db.WithContext(ctx).
		Model(&SocialProduct{}).
		Select("status, count(*) as count").
		Where("tenant_id = ? AND platform = ?", tenantID, platform).
		Group("status").
		Find(&statusCounts).Error

	if err != nil {
		return stats, err
	}

	for _, sc := range statusCounts {
		stats[string(sc.Status)] = sc.Count
	}

	// Get total count
	var total int64
	err = r.db.WithContext(ctx).
		Model(&SocialProduct{}).
		Where("tenant_id = ? AND platform = ?", tenantID, platform).
		Count(&total).Error

	stats["total"] = total

	return stats, err
}

// UpdateProductSyncStatus bulk updates sync status
func (r *repository) UpdateProductSyncStatus(ctx context.Context, tenantID uuid.UUID, platform Platform, productIDs []uuid.UUID, status SyncStatus) error {
	return r.db.WithContext(ctx).
		Model(&SocialProduct{}).
		Where("tenant_id = ? AND platform = ? AND product_id IN ?", tenantID, platform, productIDs).
		Update("status", status).Error
}

// DeleteSocialProductsByProduct deletes all social products for a given product (when product is deleted)
func (r *repository) DeleteSocialProductsByProduct(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND product_id = ?", tenantID, productID).
		Delete(&SocialProduct{}).Error
}

// GetIntegrationStats gets integration statistics for dashboard
func (r *repository) GetIntegrationStats(ctx context.Context, tenantID uuid.UUID) (map[Platform]map[string]interface{}, error) {
	integrations, err := r.GetIntegrationsByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	stats := make(map[Platform]map[string]interface{})

	for _, integration := range integrations {
		platformStats := make(map[string]interface{})

		// Get product counts
		socialProducts, _ := r.GetSocialProductsByTenant(ctx, tenantID, &integration.Platform)

		var active, synced, pending, failed int
		for _, sp := range socialProducts {
			if sp.IsEnabled {
				active++
			}
			switch sp.Status {
			case SyncCompleted:
				synced++
			case SyncPending, SyncInProgress:
				pending++
			case SyncFailed:
				failed++
			}
		}

		platformStats["connected"] = integration.IsConnected()
		platformStats["total_products"] = len(socialProducts)
		platformStats["active_products"] = active
		platformStats["synced_products"] = synced
		platformStats["pending_products"] = pending
		platformStats["failed_products"] = failed
		platformStats["last_sync"] = integration.LastSyncAt
		platformStats["next_sync"] = integration.NextSyncAt
		platformStats["error_count"] = integration.ErrorCount

		stats[integration.Platform] = platformStats
	}

	return stats, nil
}