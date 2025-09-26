package social

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service handles social commerce business logic
type Service struct {
	repo      Repository
	validator interface{}
}

// NewService creates a new social service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Repository interface for social data access
type Repository interface {
	// Social integrations
	CreateIntegration(ctx context.Context, integration *SocialIntegration) error
	GetIntegrationByTenant(ctx context.Context, tenantID uuid.UUID, platform Platform) (*SocialIntegration, error)
	GetIntegrationsByTenant(ctx context.Context, tenantID uuid.UUID) ([]SocialIntegration, error)
	UpdateIntegration(ctx context.Context, integration *SocialIntegration) error
	DeleteIntegration(ctx context.Context, tenantID uuid.UUID, platform Platform) error

	// Social products
	CreateSocialProduct(ctx context.Context, socialProduct *SocialProduct) error
	GetSocialProduct(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID, platform Platform) (*SocialProduct, error)
	GetSocialProductsByTenant(ctx context.Context, tenantID uuid.UUID, platform *Platform) ([]SocialProduct, error)
	UpdateSocialProduct(ctx context.Context, socialProduct *SocialProduct) error
	DeleteSocialProduct(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID, platform Platform) error

	// Analytics
	CreateAnalytics(ctx context.Context, analytics *SocialAnalytics) error
	GetAnalytics(ctx context.Context, tenantID uuid.UUID, platform Platform, dateFrom, dateTo time.Time) ([]SocialAnalytics, error)
}

// ConnectPlatform connects a social media platform
func (s *Service) ConnectPlatform(ctx context.Context, tenantID uuid.UUID, req *ConnectPlatformRequest) (*SocialIntegration, error) {
	// Check if integration already exists
	existing, err := s.repo.GetIntegrationByTenant(ctx, tenantID, req.Platform)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing integration: %w", err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(req.ExpiresIn) * time.Second)

	if existing != nil {
		// Update existing integration
		existing.Status = StatusConnected
		existing.AccessToken = req.AccessToken
		existing.RefreshToken = req.RefreshToken
		existing.ExpiresAt = expiresAt
		existing.BusinessAccountID = req.BusinessAccountID
		existing.ClearError()

		if err := s.repo.UpdateIntegration(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to update integration: %w", err)
		}
		return existing, nil
	}

	// Create new integration
	integration := &SocialIntegration{
		ID:                uuid.New(),
		TenantID:          tenantID,
		Platform:          req.Platform,
		Status:            StatusConnected,
		AccessToken:       req.AccessToken,
		RefreshToken:      req.RefreshToken,
		ExpiresAt:         expiresAt,
		BusinessAccountID: req.BusinessAccountID,
		AutoSync:          true,
		SyncFrequency:     24, // Default to 24 hours
	}

	integration.ScheduleNextSync()

	if err := s.repo.CreateIntegration(ctx, integration); err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}

	return integration, nil
}

// DisconnectPlatform disconnects a social media platform
func (s *Service) DisconnectPlatform(ctx context.Context, tenantID uuid.UUID, platform Platform) error {
	integration, err := s.repo.GetIntegrationByTenant(ctx, tenantID, platform)
	if err != nil {
		return fmt.Errorf("integration not found: %w", err)
	}

	integration.Status = StatusDisconnected
	integration.AccessToken = ""
	integration.RefreshToken = ""

	if err := s.repo.UpdateIntegration(ctx, integration); err != nil {
		return fmt.Errorf("failed to disconnect integration: %w", err)
	}

	return nil
}

// GetIntegrationStatus gets the status of social platform integrations
func (s *Service) GetIntegrationStatus(ctx context.Context, tenantID uuid.UUID) ([]IntegrationStatusResponse, error) {
	integrations, err := s.repo.GetIntegrationsByTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integrations: %w", err)
	}

	var responses []IntegrationStatusResponse
	for _, integration := range integrations {
		// Get product counts for this platform
		socialProducts, err := s.repo.GetSocialProductsByTenant(ctx, tenantID, &integration.Platform)
		if err != nil {
			continue // Skip on error
		}

		var syncedCount, pendingCount, errorCount int
		for _, sp := range socialProducts {
			switch sp.Status {
			case SyncCompleted:
				syncedCount++
			case SyncPending, SyncInProgress:
				pendingCount++
			case SyncFailed:
				errorCount++
			}
		}

		response := IntegrationStatusResponse{
			Platform:     integration.Platform,
			Status:       integration.Status,
			Username:     integration.PlatformUsername,
			IsConnected:  integration.IsConnected(),
			LastSyncAt:   integration.LastSyncAt,
			NextSyncAt:   integration.NextSyncAt,
			ProductCount: len(socialProducts),
			SyncedCount:  syncedCount,
			PendingCount: pendingCount,
			ErrorCount:   errorCount,
			LastError:    integration.LastError,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// SyncProducts syncs products to social platforms
func (s *Service) SyncProducts(ctx context.Context, tenantID uuid.UUID, req *SyncProductsRequest) error {
	// Get integration for the platform
	integration, err := s.repo.GetIntegrationByTenant(ctx, tenantID, req.Platform)
	if err != nil {
		return fmt.Errorf("platform not connected: %w", err)
	}

	if !integration.CanSync() {
		return fmt.Errorf("platform cannot sync: status=%s, expired=%v", integration.Status, integration.NeedsReauth())
	}

	// Get products to sync
	var socialProducts []SocialProduct
	if len(req.ProductIDs) > 0 {
		// Sync specific products
		for _, productID := range req.ProductIDs {
			sp, err := s.repo.GetSocialProduct(ctx, tenantID, productID, req.Platform)
			if err != nil && err != gorm.ErrRecordNotFound {
				continue // Skip on error
			}

			if sp == nil {
				// Create new social product entry
				sp = &SocialProduct{
					ID:        uuid.New(),
					TenantID:  tenantID,
					ProductID: productID,
					Platform:  req.Platform,
					Status:    SyncPending,
					IsEnabled: true,
				}
				if err := s.repo.CreateSocialProduct(ctx, sp); err != nil {
					continue // Skip on error
				}
			}

			socialProducts = append(socialProducts, *sp)
		}
	} else {
		// Sync all products
		socialProducts, err = s.repo.GetSocialProductsByTenant(ctx, tenantID, &req.Platform)
		if err != nil {
			return fmt.Errorf("failed to get social products: %w", err)
		}
	}

	// TODO: Implement actual sync logic with platform APIs
	// This would involve:
	// 1. Get product details from product service
	// 2. Format for platform API
	// 3. Call platform API to create/update products
	// 4. Update sync status based on API response

	// For now, mark as completed (placeholder)
	for _, sp := range socialProducts {
		if req.Force || sp.CanSync() {
			sp.RecordSyncAttempt()
			sp.RecordSyncSuccess() // TODO: Replace with actual sync logic
			s.repo.UpdateSocialProduct(ctx, &sp)
		}
	}

	// Update integration sync time
	now := time.Now()
	integration.LastSyncAt = &now
	integration.ScheduleNextSync()
	s.repo.UpdateIntegration(ctx, integration)

	return nil
}

// GetSyncStatus gets the sync status of products
func (s *Service) GetSyncStatus(ctx context.Context, tenantID uuid.UUID, platform *Platform) ([]ProductSyncStatusResponse, error) {
	socialProducts, err := s.repo.GetSocialProductsByTenant(ctx, tenantID, platform)
	if err != nil {
		return nil, fmt.Errorf("failed to get social products: %w", err)
	}

	var responses []ProductSyncStatusResponse
	for _, sp := range socialProducts {
		// TODO: Get product name from product service
		response := ProductSyncStatusResponse{
			ProductID:         sp.ProductID,
			ProductName:       "Product " + sp.ProductID.String(), // TODO: Get real name
			Platform:          sp.Platform,
			Status:            sp.Status,
			PlatformProductID: sp.PlatformProductID,
			LastSyncAt:        sp.LastSyncAt,
			LastSyncError:     sp.LastSyncError,
			IsEnabled:         sp.IsEnabled,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// UpdateSocialProduct updates social product settings
func (s *Service) UpdateSocialProduct(ctx context.Context, tenantID, productID uuid.UUID, platform Platform, req *UpdateSocialProductRequest) (*SocialProduct, error) {
	socialProduct, err := s.repo.GetSocialProduct(ctx, tenantID, productID, platform)
	if err != nil {
		return nil, fmt.Errorf("social product not found: %w", err)
	}

	// Update fields
	if req.IsEnabled != nil {
		socialProduct.IsEnabled = *req.IsEnabled
	}
	if req.CustomTitle != nil {
		socialProduct.CustomTitle = *req.CustomTitle
	}
	if req.CustomDescription != nil {
		socialProduct.CustomDescription = *req.CustomDescription
	}
	if req.Tags != nil {
		socialProduct.Tags = req.Tags
	}
	if req.InstagramSettings != nil {
		socialProduct.InstagramSettings = req.InstagramSettings
	}
	if req.FacebookSettings != nil {
		socialProduct.FacebookSettings = req.FacebookSettings
	}
	if req.TikTokSettings != nil {
		socialProduct.TikTokSettings = req.TikTokSettings
	}

	// Mark for re-sync if enabled
	if socialProduct.IsEnabled && socialProduct.Status == SyncCompleted {
		socialProduct.Status = SyncPending
	}

	if err := s.repo.UpdateSocialProduct(ctx, socialProduct); err != nil {
		return nil, fmt.Errorf("failed to update social product: %w", err)
	}

	return socialProduct, nil
}

// GetAnalytics gets social commerce analytics
func (s *Service) GetAnalytics(ctx context.Context, tenantID uuid.UUID, platform Platform, dateFrom, dateTo time.Time) (*SocialAnalyticsResponse, error) {
	analytics, err := s.repo.GetAnalytics(ctx, tenantID, platform, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics: %w", err)
	}

	// Aggregate totals
	var response SocialAnalyticsResponse
	response.Platform = platform
	response.DateFrom = dateFrom.Format("2006-01-02")
	response.DateTo = dateTo.Format("2006-01-02")
	response.DailyBreakdown = analytics

	for _, a := range analytics {
		response.TotalImpressions += a.Impressions
		response.TotalClicks += a.Clicks
		response.TotalProductViews += a.ProductViews
		response.TotalAddToCarts += a.AddToCarts
		response.TotalPurchases += a.Purchases
		response.TotalRevenue += a.Revenue
	}

	// Calculate averages
	if len(analytics) > 0 {
		response.AverageCTR = response.AverageCTR / float64(len(analytics))
		response.AverageConversion = response.AverageConversion / float64(len(analytics))
	}

	return &response, nil
}

// Helper method to validate platform
func (s *Service) validatePlatform(platform Platform) error {
	switch platform {
	case PlatformInstagram, PlatformFacebook, PlatformTikTok:
		return nil
	default:
		return fmt.Errorf("unsupported platform: %s", platform)
	}
}