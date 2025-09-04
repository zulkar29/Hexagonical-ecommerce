package shipping

import (
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

// Shipping Zone Repository Methods

func (r *Repository) CreateShippingZone(zone *ShippingZone) (*ShippingZone, error) {
	if err := r.db.Create(zone).Error; err != nil {
		return nil, err
	}
	return zone, nil
}

func (r *Repository) UpdateShippingZone(zone *ShippingZone) (*ShippingZone, error) {
	if err := r.db.Save(zone).Error; err != nil {
		return nil, err
	}
	return zone, nil
}

func (r *Repository) DeleteShippingZone(tenantID, zoneID uuid.UUID) error {
	return r.db.Where("tenant_id = ? AND id = ?", tenantID, zoneID).Delete(&ShippingZone{}).Error
}

func (r *Repository) GetShippingZone(tenantID, zoneID uuid.UUID) (*ShippingZone, error) {
	var zone ShippingZone
	err := r.db.Preload("Countries").Preload("Rates").
		Where("tenant_id = ? AND id = ?", tenantID, zoneID).
		First(&zone).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *Repository) GetShippingZones(tenantID uuid.UUID) ([]ShippingZone, error) {
	var zones []ShippingZone
	err := r.db.Preload("Countries").Preload("Rates").
		Where("tenant_id = ?", tenantID).
		Order("is_default DESC, priority ASC, name ASC").
		Find(&zones).Error
	return zones, err
}

func (r *Repository) GetShippingZonesForDestination(tenantID uuid.UUID, country, state, city string) ([]ShippingZone, error) {
	var zones []ShippingZone
	
	// Subquery to find zone IDs that match the destination
	subQuery := r.db.Model(&ShippingZoneCountry{}).
		Select("zone_id").
		Where("country = ?", country)
	
	if state != "" {
		subQuery = subQuery.Where("states IS NULL OR states LIKE ?", "%"+state+"%")
	}
	if city != "" {
		subQuery = subQuery.Where("cities IS NULL OR cities LIKE ?", "%"+city+"%")
	}

	err := r.db.Preload("Countries").Preload("Rates").
		Where("tenant_id = ? AND is_active = ? AND id IN (?)", tenantID, true, subQuery).
		Order("is_default DESC, priority ASC").
		Find(&zones).Error

	return zones, err
}

// Shipping Rate Repository Methods

func (r *Repository) CreateShippingRate(rate *ShippingRate) (*ShippingRate, error) {
	if err := r.db.Create(rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (r *Repository) UpdateShippingRate(rate *ShippingRate) (*ShippingRate, error) {
	if err := r.db.Save(rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (r *Repository) DeleteShippingRate(tenantID, rateID uuid.UUID) error {
	return r.db.Where("tenant_id = ? AND id = ?", tenantID, rateID).Delete(&ShippingRate{}).Error
}

func (r *Repository) GetShippingRate(rateID uuid.UUID) (*ShippingRate, error) {
	var rate ShippingRate
	err := r.db.Where("id = ?", rateID).First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *Repository) GetShippingRatesForZone(zoneID uuid.UUID) ([]ShippingRate, error) {
	var rates []ShippingRate
	err := r.db.Where("zone_id = ? AND is_active = ?", zoneID, true).
		Order("priority ASC, base_rate ASC").
		Find(&rates).Error
	return rates, err
}

func (r *Repository) GetShippingRates(tenantID, zoneID uuid.UUID) ([]ShippingRate, error) {
	var rates []ShippingRate
	err := r.db.Where("tenant_id = ? AND zone_id = ?", tenantID, zoneID).
		Order("priority ASC, base_rate ASC").
		Find(&rates).Error
	return rates, err
}

// Shipping Label Repository Methods

func (r *Repository) CreateShippingLabel(label *ShippingLabel) (*ShippingLabel, error) {
	if err := r.db.Create(label).Error; err != nil {
		return nil, err
	}
	return label, nil
}

func (r *Repository) UpdateShippingLabel(label *ShippingLabel) (*ShippingLabel, error) {
	if err := r.db.Save(label).Error; err != nil {
		return nil, err
	}
	return label, nil
}

func (r *Repository) GetShippingLabel(tenantID, labelID uuid.UUID) (*ShippingLabel, error) {
	var label ShippingLabel
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, labelID).First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *Repository) GetShippingLabelByTrackingNumber(trackingNumber string) (*ShippingLabel, error) {
	var label ShippingLabel
	err := r.db.Where("tracking_number = ?", trackingNumber).First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *Repository) GetShippingLabelsByOrder(tenantID, orderID uuid.UUID) ([]ShippingLabel, error) {
	var labels []ShippingLabel
	err := r.db.Where("tenant_id = ? AND order_id = ?", tenantID, orderID).
		Order("created_at DESC").
		Find(&labels).Error
	return labels, err
}

func (r *Repository) GetShippingLabels(tenantID uuid.UUID, offset, limit int) ([]ShippingLabel, int64, error) {
	var labels []ShippingLabel
	var total int64
	
	// Get total count
	r.db.Model(&ShippingLabel{}).Where("tenant_id = ?", tenantID).Count(&total)
	
	// Get paginated results
	err := r.db.Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&labels).Error
		
	return labels, total, err
}

// Shipping Tracking Repository Methods

func (r *Repository) CreateShippingTracking(tracking *ShippingTracking) (*ShippingTracking, error) {
	if err := r.db.Create(tracking).Error; err != nil {
		return nil, err
	}
	return tracking, nil
}

func (r *Repository) GetTrackingHistory(labelID uuid.UUID) ([]ShippingTracking, error) {
	var tracking []ShippingTracking
	err := r.db.Where("label_id = ?", labelID).
		Order("timestamp DESC").
		Find(&tracking).Error
	return tracking, err
}

func (r *Repository) GetLatestShippingTracking(labelID uuid.UUID) (*ShippingTracking, error) {
	var tracking ShippingTracking
	err := r.db.Where("label_id = ?", labelID).
		Order("timestamp DESC").
		First(&tracking).Error
	if err != nil {
		return nil, err
	}
	return &tracking, nil
}

// Shipping Provider Repository Methods

func (r *Repository) CreateOrUpdateProvider(provider *ShippingProvider) error {
	// First try to find existing provider
	var existing ShippingProvider
	err := r.db.Where("tenant_id = ? AND provider = ?", provider.TenantID, provider.Provider).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new provider
		return r.db.Create(provider).Error
	} else if err != nil {
		return err
	}
	
	// Update existing provider
	existing.Name = provider.Name
	existing.IsActive = provider.IsActive
	existing.APIKey = provider.APIKey
	existing.APISecret = provider.APISecret
	existing.SandboxMode = provider.SandboxMode
	existing.Settings = provider.Settings
	
	return r.db.Save(&existing).Error
}

func (r *Repository) GetShippingProvider(tenantID uuid.UUID, providerName string) (*ShippingProvider, error) {
	var provider ShippingProvider
	err := r.db.Where("tenant_id = ? AND provider = ?", tenantID, providerName).First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *Repository) GetShippingProviders(tenantID uuid.UUID) ([]ShippingProvider, error) {
	var providers []ShippingProvider
	err := r.db.Where("tenant_id = ?", tenantID).
		Order("name ASC").
		Find(&providers).Error
	return providers, err
}

func (r *Repository) GetActiveShippingProviders(tenantID uuid.UUID) ([]ShippingProvider, error) {
	var providers []ShippingProvider
	err := r.db.Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Order("name ASC").
		Find(&providers).Error
	return providers, err
}

// Statistics Repository Methods

func (r *Repository) GetShippingStats(tenantID uuid.UUID) (*ShippingStats, error) {
	stats := &ShippingStats{}
	
	// Get basic label counts
	r.db.Model(&ShippingLabel{}).Where("tenant_id = ?", tenantID).Count(&[]int64{int64(stats.TotalLabels)}[0])
	r.db.Model(&ShippingLabel{}).Where("tenant_id = ? AND status NOT IN (?)", tenantID, []string{"cancelled", "delivered"}).Count(&[]int64{int64(stats.ActiveLabels)}[0])
	r.db.Model(&ShippingLabel{}).Where("tenant_id = ? AND status = ?", tenantID, "delivered").Count(&[]int64{int64(stats.DeliveredLabels)}[0])
	r.db.Model(&ShippingLabel{}).Where("tenant_id = ? AND status = ?", tenantID, "cancelled").Count(&[]int64{int64(stats.CancelledLabels)}[0])
	
	// Get total cost
	var totalCostResult struct {
		TotalCost float64 `json:"total_cost"`
	}
	r.db.Model(&ShippingLabel{}).
		Where("tenant_id = ? AND status != ?", tenantID, "cancelled").
		Select("SUM(cost) as total_cost").
		Scan(&totalCostResult)
	stats.TotalCost = totalCostResult.TotalCost
	
	// Get provider stats
	var providerResults []struct {
		Provider      string  `json:"provider"`
		TotalLabels   int     `json:"total_labels"`
		DeliveredRate float64 `json:"delivered_rate"`
		TotalCost     float64 `json:"total_cost"`
	}
	
	r.db.Model(&ShippingLabel{}).
		Where("tenant_id = ?", tenantID).
		Select("provider, COUNT(*) as total_labels, AVG(CASE WHEN status = 'delivered' THEN 1.0 ELSE 0.0 END) as delivered_rate, SUM(cost) as total_cost").
		Group("provider").
		Scan(&providerResults)
	
	for _, result := range providerResults {
		providerStat := ProviderStats{
			Provider:      ShippingProvider(result.Provider),
			TotalLabels:   result.TotalLabels,
			DeliveredRate: result.DeliveredRate * 100, // Convert to percentage
			TotalCost:     result.TotalCost,
		}
		stats.ProviderStats = append(stats.ProviderStats, providerStat)
	}
	
	// Get monthly stats for last 12 months
	var monthlyResults []struct {
		Month         string  `json:"month"`
		TotalLabels   int     `json:"total_labels"`
		TotalCost     float64 `json:"total_cost"`
		DeliveredRate float64 `json:"delivered_rate"`
	}
	
	r.db.Model(&ShippingLabel{}).
		Where("tenant_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 12 MONTH)", tenantID).
		Select("DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as total_labels, SUM(cost) as total_cost, AVG(CASE WHEN status = 'delivered' THEN 1.0 ELSE 0.0 END) as delivered_rate").
		Group("DATE_FORMAT(created_at, '%Y-%m')").
		Order("month DESC").
		Scan(&monthlyResults)
	
	for _, result := range monthlyResults {
		monthlyStat := MonthlyShippingStats{
			Month:         result.Month,
			TotalLabels:   result.TotalLabels,
			TotalCost:     result.TotalCost,
			DeliveredRate: result.DeliveredRate * 100, // Convert to percentage
		}
		stats.MonthlyStats = append(stats.MonthlyStats, monthlyStat)
	}
	
	return stats, nil
}

// Additional utility methods

func (r *Repository) GetShippingZoneByDefault(tenantID uuid.UUID) (*ShippingZone, error) {
	var zone ShippingZone
	err := r.db.Preload("Countries").Preload("Rates").
		Where("tenant_id = ? AND is_default = ?", tenantID, true).
		First(&zone).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *Repository) GetShippingRatesByProvider(tenantID uuid.UUID, provider ShippingProvider) ([]ShippingRate, error) {
	var rates []ShippingRate
	err := r.db.Where("tenant_id = ? AND provider = ? AND is_active = ?", tenantID, provider, true).
		Order("base_rate ASC").
		Find(&rates).Error
	return rates, err
}

func (r *Repository) GetRecentShippingLabels(tenantID uuid.UUID, limit int) ([]ShippingLabel, error) {
	var labels []ShippingLabel
	err := r.db.Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Find(&labels).Error
	return labels, err
}

func (r *Repository) GetPendingShipments(tenantID uuid.UUID) ([]ShippingLabel, error) {
	var labels []ShippingLabel
	err := r.db.Where("tenant_id = ? AND status IN (?)", tenantID, []string{"created", "printed", "shipped"}).
		Order("created_at ASC").
		Find(&labels).Error
	return labels, err
}