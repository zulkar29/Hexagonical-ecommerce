package webhook

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Webhook Endpoint Repository Methods

func (r *Repository) CreateEndpoint(endpoint *WebhookEndpoint) (*WebhookEndpoint, error) {
	if err := r.db.Create(endpoint).Error; err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (r *Repository) UpdateEndpoint(endpoint *WebhookEndpoint) (*WebhookEndpoint, error) {
	if err := r.db.Save(endpoint).Error; err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (r *Repository) DeleteEndpoint(tenantID uuid.UUID, endpointID uuid.UUID) error {
	return r.db.Where("tenant_id = ? AND id = ?", tenantID, endpointID).Delete(&WebhookEndpoint{}).Error
}

func (r *Repository) GetEndpointByID(tenantID uuid.UUID, endpointID uuid.UUID) (*WebhookEndpoint, error) {
	var endpoint WebhookEndpoint
	if err := r.db.Where("tenant_id = ? AND id = ?", tenantID, endpointID).First(&endpoint).Error; err != nil {
		return nil, err
	}
	return &endpoint, nil
}

func (r *Repository) GetEndpoints(tenantID uuid.UUID) ([]*WebhookEndpoint, error) {
	var endpoints []*WebhookEndpoint
	if err := r.db.Where("tenant_id = ?", tenantID).Find(&endpoints).Error; err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (r *Repository) GetActiveEndpoints(tenantID uuid.UUID) ([]*WebhookEndpoint, error) {
	var endpoints []*WebhookEndpoint
	if err := r.db.Where("tenant_id = ? AND is_active = ?", tenantID, true).Find(&endpoints).Error; err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (r *Repository) GetEndpointsByEvent(tenantID uuid.UUID, event WebhookEvent) ([]*WebhookEndpoint, error) {
	var endpoints []*WebhookEndpoint
	if err := r.db.Where("tenant_id = ? AND is_active = ? AND JSON_CONTAINS(events, ?)", tenantID, true, string(event)).Find(&endpoints).Error; err != nil {
		return nil, err
	}
	return endpoints, nil
}

// Webhook Delivery Repository Methods

func (r *Repository) CreateDelivery(delivery *WebhookDelivery) (*WebhookDelivery, error) {
	if err := r.db.Create(delivery).Error; err != nil {
		return nil, err
	}
	return delivery, nil
}

func (r *Repository) UpdateDelivery(delivery *WebhookDelivery) (*WebhookDelivery, error) {
	if err := r.db.Save(delivery).Error; err != nil {
		return nil, err
	}
	return delivery, nil
}

func (r *Repository) GetDeliveryByID(tenantID uuid.UUID, deliveryID uuid.UUID) (*WebhookDelivery, error) {
	var delivery WebhookDelivery
	if err := r.db.Where("tenant_id = ? AND id = ?", tenantID, deliveryID).First(&delivery).Error; err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (r *Repository) GetDeliveries(tenantID uuid.UUID, endpointID uuid.UUID, limit int, offset int) ([]*WebhookDelivery, error) {
	var deliveries []*WebhookDelivery
	query := r.db.Where("tenant_id = ?", tenantID)
	if endpointID != uuid.Nil {
		query = query.Where("endpoint_id = ?", endpointID)
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}

func (r *Repository) GetFailedDeliveries(tenantID uuid.UUID, limit int) ([]*WebhookDelivery, error) {
	var deliveries []*WebhookDelivery
	if err := r.db.Where("tenant_id = ? AND status = ?", tenantID, "failed").Order("created_at DESC").Limit(limit).Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}

func (r *Repository) GetPendingRetries(limit int) ([]*WebhookDelivery, error) {
	var deliveries []*WebhookDelivery
	if err := r.db.Where("status = ? AND next_retry_at <= ?", "failed", time.Now()).Order("next_retry_at ASC").Limit(limit).Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}

func (r *Repository) GetDeliveriesByEvent(tenantID uuid.UUID, event WebhookEvent, limit int, offset int) ([]*WebhookDelivery, error) {
	var deliveries []*WebhookDelivery
	if err := r.db.Where("tenant_id = ? AND event = ?", tenantID, event).Order("created_at DESC").Limit(limit).Offset(offset).Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}

// Incoming Webhook Repository Methods

func (r *Repository) CreateIncomingWebhook(webhook *WebhookIncoming) (*WebhookIncoming, error) {
	if err := r.db.Create(webhook).Error; err != nil {
		return nil, err
	}
	return webhook, nil
}

func (r *Repository) UpdateIncomingWebhook(webhook *WebhookIncoming) (*WebhookIncoming, error) {
	if err := r.db.Save(webhook).Error; err != nil {
		return nil, err
	}
	return webhook, nil
}

func (r *Repository) GetIncomingWebhookByID(id uuid.UUID) (*WebhookIncoming, error) {
	var webhook WebhookIncoming
	if err := r.db.Where("id = ?", id).First(&webhook).Error; err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (r *Repository) GetIncomingWebhooks(tenantID uuid.UUID, provider WebhookProvider, limit int, offset int) ([]*WebhookIncoming, error) {
	var webhooks []*WebhookIncoming
	query := r.db.Where("tenant_id = ?", tenantID)
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (r *Repository) GetUnprocessedIncomingWebhooks(limit int) ([]*WebhookIncoming, error) {
	var webhooks []*WebhookIncoming
	if err := r.db.Where("processed = ?", false).Order("created_at ASC").Limit(limit).Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (r *Repository) MarkIncomingWebhookProcessed(id uuid.UUID) error {
	return r.db.Model(&WebhookIncoming{}).Where("id = ?", id).Update("processed", true).Error
}

// Rate Limiting Repository Methods

func (r *Repository) CreateRateLimit(rateLimit *WebhookRateLimit) (*WebhookRateLimit, error) {
	if err := r.db.Create(rateLimit).Error; err != nil {
		return nil, err
	}
	return rateLimit, nil
}

func (r *Repository) UpdateRateLimit(rateLimit *WebhookRateLimit) (*WebhookRateLimit, error) {
	if err := r.db.Save(rateLimit).Error; err != nil {
		return nil, err
	}
	return rateLimit, nil
}

func (r *Repository) GetRateLimit(tenantID uuid.UUID, endpointID uuid.UUID, windowStart time.Time) (*WebhookRateLimit, error) {
	var rateLimit WebhookRateLimit
	if err := r.db.Where("tenant_id = ? AND endpoint_id = ? AND window_start = ?", tenantID, endpointID, windowStart).First(&rateLimit).Error; err != nil {
		return nil, err
	}
	return &rateLimit, nil
}

func (r *Repository) CleanupExpiredRateLimits() error {
	return r.db.Where("window_start < ?", time.Now().Add(-24*time.Hour)).Delete(&WebhookRateLimit{}).Error
}

// Analytics Types
type WebhookStats struct {
	TotalDeliveries    int     `json:"total_deliveries"`
	SuccessfulDeliveries int   `json:"successful_deliveries"`
	FailedDeliveries   int     `json:"failed_deliveries"`
	SuccessRate        float64 `json:"success_rate"`
	AverageResponseTime int    `json:"average_response_time"`
}

type EndpointHealth struct {
	EndpointID         uuid.UUID `json:"endpoint_id"`
	TotalDeliveries    int       `json:"total_deliveries"`
	SuccessfulDeliveries int     `json:"successful_deliveries"`
	FailedDeliveries   int       `json:"failed_deliveries"`
	SuccessRate        float64   `json:"success_rate"`
	AverageResponseTime int      `json:"average_response_time"`
	LastDeliveryAt     *time.Time `json:"last_delivery_at"`
	IsHealthy          bool      `json:"is_healthy"`
}

type FailureAnalysis struct {
	TotalFailures      int                    `json:"total_failures"`
	FailuresByStatus   map[int]int           `json:"failures_by_status"`
	FailuresByEndpoint map[uuid.UUID]int     `json:"failures_by_endpoint"`
	CommonErrors       []string              `json:"common_errors"`
}

// Analytics and Monitoring Repository Methods

func (r *Repository) GetDeliveryStats(tenantID uuid.UUID, startDate, endDate time.Time) (*WebhookStats, error) {
	var stats WebhookStats
	
	// Get total deliveries
	if err := r.db.Model(&WebhookDelivery{}).Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, startDate, endDate).Count(&stats.TotalDeliveries).Error; err != nil {
		return nil, err
	}
	
	// Get successful deliveries
	if err := r.db.Model(&WebhookDelivery{}).Where("tenant_id = ? AND status = ? AND created_at BETWEEN ? AND ?", tenantID, StatusDelivered, startDate, endDate).Count(&stats.SuccessfulDeliveries).Error; err != nil {
		return nil, err
	}
	
	// Get failed deliveries
	if err := r.db.Model(&WebhookDelivery{}).Where("tenant_id = ? AND status = ? AND created_at BETWEEN ? AND ?", tenantID, StatusFailed, startDate, endDate).Count(&stats.FailedDeliveries).Error; err != nil {
		return nil, err
	}
	
	// Calculate success rate
	if stats.TotalDeliveries > 0 {
		stats.SuccessRate = float64(stats.SuccessfulDeliveries) / float64(stats.TotalDeliveries) * 100
	}
	
	// Get average response time
	var avgResponseTime float64
	if err := r.db.Model(&WebhookDelivery{}).Where("tenant_id = ? AND status = ? AND created_at BETWEEN ? AND ?", tenantID, StatusDelivered, startDate, endDate).Select("AVG(response_time)").Scan(&avgResponseTime).Error; err != nil {
		return nil, err
	}
	stats.AverageResponseTime = int(avgResponseTime)
	
	return &stats, nil
}

func (r *Repository) GetEndpointHealth(tenantID uuid.UUID, endpointID uuid.UUID, days int) (*EndpointHealth, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	endDate := time.Now()
	
	health := &EndpointHealth{
		EndpointID: endpointID,
	}
	
	// Get total deliveries for endpoint
	if err := r.db.Model(&WebhookDelivery{}).Where("tenant_id = ? AND endpoint_id = ? AND created_at BETWEEN ? AND ?", tenantID, endpointID, startDate, endDate).Count(&health.TotalDeliveries).Error; err != nil {
		return nil, err
	}
	
	// Get successful deliveries
	if err := r.db.Model(&WebhookDelivery{}).Where("tenant_id = ? AND endpoint_id = ? AND status = ? AND created_at BETWEEN ? AND ?", tenantID, endpointID, StatusDelivered, startDate, endDate).Count(&health.SuccessfulDeliveries).Error; err != nil {
		return nil, err
	}
	
	// Get failed deliveries
	if err := r.db.Model(&WebhookDelivery{}).Where("tenant_id = ? AND endpoint_id = ? AND status = ? AND created_at BETWEEN ? AND ?", tenantID, endpointID, StatusFailed, startDate, endDate).Count(&health.FailedDeliveries).Error; err != nil {
		return nil, err
	}
	
	// Calculate success rate
	if health.TotalDeliveries > 0 {
		health.SuccessRate = float64(health.SuccessfulDeliveries) / float64(health.TotalDeliveries) * 100
		health.IsHealthy = health.SuccessRate >= 95.0 // Consider healthy if 95%+ success rate
	}
	
	// Get average response time
	var avgResponseTime float64
	if err := r.db.Model(&WebhookDelivery{}).Where("tenant_id = ? AND endpoint_id = ? AND status = ? AND created_at BETWEEN ? AND ?", tenantID, endpointID, StatusDelivered, startDate, endDate).Select("AVG(response_time)").Scan(&avgResponseTime).Error; err != nil {
		return nil, err
	}
	health.AverageResponseTime = int(avgResponseTime)
	
	// Get last delivery time
	var lastDelivery WebhookDelivery
	if err := r.db.Where("tenant_id = ? AND endpoint_id = ?", tenantID, endpointID).Order("created_at DESC").First(&lastDelivery).Error; err == nil {
		health.LastDeliveryAt = &lastDelivery.CreatedAt
	}
	
	return health, nil
}

func (r *Repository) GetEventDistribution(tenantID uuid.UUID, startDate, endDate time.Time) (map[WebhookEvent]int, error) {
	var results []struct {
		Event WebhookEvent `json:"event"`
		Count int         `json:"count"`
	}
	
	if err := r.db.Model(&WebhookDelivery{}).Select("event, COUNT(*) as count").Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, startDate, endDate).Group("event").Scan(&results).Error; err != nil {
		return nil, err
	}
	
	distribution := make(map[WebhookEvent]int)
	for _, result := range results {
		distribution[result.Event] = result.Count
	}
	
	return distribution, nil
}

func (r *Repository) GetProviderStats(tenantID uuid.UUID, startDate, endDate time.Time) (map[WebhookProvider]int, error) {
	var results []struct {
		Provider WebhookProvider `json:"provider"`
		Count    int             `json:"count"`
	}
	
	if err := r.db.Model(&WebhookIncoming{}).Select("provider, COUNT(*) as count").Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, startDate, endDate).Group("provider").Scan(&results).Error; err != nil {
		return nil, err
	}
	
	stats := make(map[WebhookProvider]int)
	for _, result := range results {
		stats[result.Provider] = result.Count
	}
	
	return stats, nil
}

func (r *Repository) GetFailureAnalysis(tenantID uuid.UUID, startDate, endDate time.Time) (*FailureAnalysis, error) {
	analysis := &FailureAnalysis{
		FailuresByStatus:   make(map[int]int),
		FailuresByEndpoint: make(map[uuid.UUID]int),
	}
	
	// Get total failures
	if err := r.db.Model(&WebhookDelivery{}).Where("tenant_id = ? AND status = ? AND created_at BETWEEN ? AND ?", tenantID, StatusFailed, startDate, endDate).Count(&analysis.TotalFailures).Error; err != nil {
		return nil, err
	}
	
	// Get failures by status code
	var statusResults []struct {
		ResponseStatus int `json:"response_status"`
		Count         int `json:"count"`
	}
	if err := r.db.Model(&WebhookDelivery{}).Select("response_status, COUNT(*) as count").Where("tenant_id = ? AND status = ? AND created_at BETWEEN ? AND ?", tenantID, StatusFailed, startDate, endDate).Group("response_status").Scan(&statusResults).Error; err != nil {
		return nil, err
	}
	for _, result := range statusResults {
		analysis.FailuresByStatus[result.ResponseStatus] = result.Count
	}
	
	// Get failures by endpoint
	var endpointResults []struct {
		EndpointID uuid.UUID `json:"endpoint_id"`
		Count      int       `json:"count"`
	}
	if err := r.db.Model(&WebhookDelivery{}).Select("endpoint_id, COUNT(*) as count").Where("tenant_id = ? AND status = ? AND created_at BETWEEN ? AND ?", tenantID, StatusFailed, startDate, endDate).Group("endpoint_id").Scan(&endpointResults).Error; err != nil {
		return nil, err
	}
	for _, result := range endpointResults {
		analysis.FailuresByEndpoint[result.EndpointID] = result.Count
	}
	
	// Get common error messages
	var errorMessages []string
	if err := r.db.Model(&WebhookDelivery{}).Select("DISTINCT error_message").Where("tenant_id = ? AND status = ? AND error_message != '' AND created_at BETWEEN ? AND ?", tenantID, StatusFailed, startDate, endDate).Limit(10).Pluck("error_message", &errorMessages).Error; err != nil {
		return nil, err
	}
	analysis.CommonErrors = errorMessages
	
	return analysis, nil
}

// Cleanup and Maintenance Repository Methods

type DiskUsage struct {
	DeliveriesSize       int64 `json:"deliveries_size"`
	IncomingWebhooksSize int64 `json:"incoming_webhooks_size"`
	TotalSize           int64 `json:"total_size"`
	RecordCount         int64 `json:"record_count"`
}

func (r *Repository) CleanupOldDeliveries(tenantID uuid.UUID, olderThan time.Time) error {
	return r.db.Where("tenant_id = ? AND created_at < ?", tenantID, olderThan).Delete(&WebhookDelivery{}).Error
}

func (r *Repository) CleanupOldIncomingWebhooks(tenantID uuid.UUID, olderThan time.Time) error {
	return r.db.Where("tenant_id = ? AND created_at < ?", tenantID, olderThan).Delete(&WebhookIncoming{}).Error
}

func (r *Repository) ArchiveDeliveries(tenantID uuid.UUID, olderThan time.Time) error {
	// Create archive table if not exists
	if err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS webhook_deliveries_archive (
			LIKE webhook_deliveries INCLUDING ALL
		)
	`).Error; err != nil {
		return err
	}
	
	// Move old records to archive
	if err := r.db.Exec(`
		INSERT INTO webhook_deliveries_archive 
		SELECT * FROM webhook_deliveries 
		WHERE tenant_id = ? AND created_at < ?
	`, tenantID, olderThan).Error; err != nil {
		return err
	}
	
	// Delete from main table
	return r.CleanupOldDeliveries(tenantID, olderThan)
}

func (r *Repository) GetDiskUsage(tenantID uuid.UUID) (*DiskUsage, error) {
	usage := &DiskUsage{}
	
	// Count webhook deliveries
	var deliveryCount int64
	if err := r.db.Model(&WebhookDelivery{}).Where("tenant_id = ?", tenantID).Count(&deliveryCount).Error; err != nil {
		return nil, err
	}
	
	// Count incoming webhooks
	var incomingCount int64
	if err := r.db.Model(&WebhookIncoming{}).Where("tenant_id = ?", tenantID).Count(&incomingCount).Error; err != nil {
		return nil, err
	}
	
	usage.RecordCount = deliveryCount + incomingCount
	
	// Estimate sizes (rough calculation based on average record size)
	usage.DeliveriesSize = deliveryCount * 1024 // Assume 1KB per delivery record
	usage.IncomingWebhooksSize = incomingCount * 512 // Assume 512B per incoming webhook
	usage.TotalSize = usage.DeliveriesSize + usage.IncomingWebhooksSize
	
	return usage, nil
}

func (r *Repository) OptimizeDatabase() error {
	// Analyze tables for better query performance
	if err := r.db.Exec("ANALYZE webhook_endpoints, webhook_deliveries, webhook_incoming, webhook_rate_limits").Error; err != nil {
		return err
	}
	
	// Vacuum tables to reclaim space (PostgreSQL specific)
	if err := r.db.Exec("VACUUM ANALYZE webhook_endpoints, webhook_deliveries, webhook_incoming, webhook_rate_limits").Error; err != nil {
		return err
	}
	
	return nil
}
