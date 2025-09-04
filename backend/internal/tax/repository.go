package tax

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for tax data operations
type Repository interface {
	// Tax Rule operations
	CreateTaxRule(ctx context.Context, rule *TaxRule) error
	GetTaxRule(ctx context.Context, tenantID, ruleID uuid.UUID) (*TaxRule, error)
	GetTaxRuleByCode(ctx context.Context, tenantID uuid.UUID, code string) (*TaxRule, error)
	UpdateTaxRule(ctx context.Context, rule *TaxRule) error
	DeleteTaxRule(ctx context.Context, tenantID, ruleID uuid.UUID) error
	ListTaxRules(ctx context.Context, tenantID uuid.UUID, filter TaxRuleFilter, page, pageSize int) ([]*TaxRule, int64, error)
	GetActiveTaxRules(ctx context.Context, tenantID uuid.UUID, date time.Time) ([]*TaxRule, error)
	GetApplicableTaxRules(ctx context.Context, tenantID uuid.UUID, req TaxCalculationRequest) ([]*TaxRule, error)
	
	// Tax Rate operations
	CreateTaxRate(ctx context.Context, rate *TaxRate) error
	GetTaxRate(ctx context.Context, tenantID, rateID uuid.UUID) (*TaxRate, error)
	UpdateTaxRate(ctx context.Context, rate *TaxRate) error
	DeleteTaxRate(ctx context.Context, tenantID, rateID uuid.UUID) error
	ListTaxRates(ctx context.Context, tenantID uuid.UUID, filter TaxRateFilter, page, pageSize int) ([]*TaxRate, int64, error)
	GetTaxRatesByRule(ctx context.Context, tenantID, ruleID uuid.UUID) ([]*TaxRate, error)
	GetActiveTaxRates(ctx context.Context, tenantID uuid.UUID, date time.Time) ([]*TaxRate, error)
	
	// Tax Calculation operations
	CreateTax(ctx context.Context, tax *Tax) error
	GetTax(ctx context.Context, tenantID, taxID uuid.UUID) (*Tax, error)
	ListTaxes(ctx context.Context, tenantID uuid.UUID, filter TaxFilter, page, pageSize int) ([]*Tax, int64, error)
	GetTaxesByOrder(ctx context.Context, tenantID, orderID uuid.UUID) ([]*Tax, error)
	GetTaxesByProduct(ctx context.Context, tenantID, productID uuid.UUID) ([]*Tax, error)
	GetTaxesByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*Tax, error)
	
	// Tax Rule Application operations
	CreateTaxRuleApplication(ctx context.Context, application *TaxRuleApplication) error
	GetTaxRuleApplications(ctx context.Context, tenantID, taxID uuid.UUID) ([]*TaxRuleApplication, error)
	
	// Bulk operations
	BulkCreateTaxRules(ctx context.Context, rules []*TaxRule) error
	BulkUpdateTaxRuleStatus(ctx context.Context, tenantID uuid.UUID, ruleIDs []uuid.UUID, status string) error
	BulkDeleteTaxRules(ctx context.Context, tenantID uuid.UUID, ruleIDs []uuid.UUID) error
	BulkCreateTaxRates(ctx context.Context, rates []*TaxRate) error
	BulkUpdateTaxRateStatus(ctx context.Context, tenantID uuid.UUID, rateIDs []uuid.UUID, isActive bool) error
	
	// Validation operations
	CheckTaxRuleCodeExists(ctx context.Context, tenantID uuid.UUID, code string, excludeID *uuid.UUID) (bool, error)
	ValidateTaxRuleReferences(ctx context.Context, tenantID uuid.UUID, productIDs, categoryIDs, customerIDs []uuid.UUID) error
	
	// Statistics and analytics
	GetTaxStats(ctx context.Context, tenantID uuid.UUID) (*TaxStats, error)
	GetTaxStatsByLocation(ctx context.Context, tenantID uuid.UUID, limit int) ([]*TaxByLocation, error)
	GetTaxStatsByType(ctx context.Context, tenantID uuid.UUID) ([]*TaxByType, error)
	GetTaxTrends(ctx context.Context, tenantID uuid.UUID, days int) (map[string]interface{}, error)
	
	// Maintenance operations
	CleanupExpiredRules(ctx context.Context, tenantID uuid.UUID) (int64, error)
	CleanupExpiredRates(ctx context.Context, tenantID uuid.UUID) (int64, error)
	ArchiveOldTaxCalculations(ctx context.Context, tenantID uuid.UUID, olderThan time.Time) (int64, error)
}

// GormRepository implements Repository using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM repository
func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// Tax Rule operations

// CreateTaxRule creates a new tax rule
func (r *GormRepository) CreateTaxRule(ctx context.Context, rule *TaxRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

// GetTaxRule retrieves a tax rule by ID
func (r *GormRepository) GetTaxRule(ctx context.Context, tenantID, ruleID uuid.UUID) (*TaxRule, error) {
	var rule TaxRule
	err := r.db.WithContext(ctx).
		Preload("Rates").
		Where("tenant_id = ? AND id = ?", tenantID, ruleID).
		First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// GetTaxRuleByCode retrieves a tax rule by code
func (r *GormRepository) GetTaxRuleByCode(ctx context.Context, tenantID uuid.UUID, code string) (*TaxRule, error) {
	var rule TaxRule
	err := r.db.WithContext(ctx).
		Preload("Rates").
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdateTaxRule updates a tax rule
func (r *GormRepository) UpdateTaxRule(ctx context.Context, rule *TaxRule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

// DeleteTaxRule soft deletes a tax rule
func (r *GormRepository) DeleteTaxRule(ctx context.Context, tenantID, ruleID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, ruleID).
		Delete(&TaxRule{}).Error
}

// ListTaxRules retrieves tax rules with filtering and pagination
func (r *GormRepository) ListTaxRules(ctx context.Context, tenantID uuid.UUID, filter TaxRuleFilter, page, pageSize int) ([]*TaxRule, int64, error) {
	query := r.db.WithContext(ctx).Model(&TaxRule{}).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.TaxType != "" {
		query = query.Where("tax_type = ?", filter.TaxType)
	}
	if filter.Country != "" {
		query = query.Where("countries @> ?", fmt.Sprintf(`["%s"]`, filter.Country))
	}
	if filter.State != "" {
		query = query.Where("states @> ?", fmt.Sprintf(`["%s"]`, filter.State))
	}
	if filter.City != "" {
		query = query.Where("cities @> ?", fmt.Sprintf(`["%s"]`, filter.City))
	}
	if filter.ProductID != nil {
		query = query.Where("product_ids @> ?", fmt.Sprintf(`["%s"]`, filter.ProductID.String()))
	}
	if filter.CategoryID != nil {
		query = query.Where("category_ids @> ?", fmt.Sprintf(`["%s"]`, filter.CategoryID.String()))
	}
	if filter.CustomerID != nil {
		query = query.Where("customer_ids @> ?", fmt.Sprintf(`["%s"]`, filter.CustomerID.String()))
	}
	if filter.ValidDate != nil {
		query = query.Where("(valid_from IS NULL OR valid_from <= ?) AND (valid_to IS NULL OR valid_to >= ?)", 
			filter.ValidDate, filter.ValidDate)
	}
	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?", 
			searchTerm, searchTerm, searchTerm)
	}
	
	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination and get results
	var rules []*TaxRule
	offset := (page - 1) * pageSize
	err := query.Preload("Rates").
		Order("priority DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&rules).Error
	
	return rules, total, err
}

// GetActiveTaxRules retrieves active tax rules for a specific date
func (r *GormRepository) GetActiveTaxRules(ctx context.Context, tenantID uuid.UUID, date time.Time) ([]*TaxRule, error) {
	var rules []*TaxRule
	err := r.db.WithContext(ctx).
		Preload("Rates").
		Where("tenant_id = ? AND status = ?", tenantID, StatusActive).
		Where("(valid_from IS NULL OR valid_from <= ?) AND (valid_to IS NULL OR valid_to >= ?)", date, date).
		Order("priority DESC, created_at ASC").
		Find(&rules).Error
	return rules, err
}

// GetApplicableTaxRules retrieves tax rules applicable to a calculation request
func (r *GormRepository) GetApplicableTaxRules(ctx context.Context, tenantID uuid.UUID, req TaxCalculationRequest) ([]*TaxRule, error) {
	date := time.Now()
	if req.Date != nil {
		date = *req.Date
	}
	
	query := r.db.WithContext(ctx).
		Preload("Rates").
		Where("tenant_id = ? AND status = ?", tenantID, StatusActive).
		Where("(valid_from IS NULL OR valid_from <= ?) AND (valid_to IS NULL OR valid_to >= ?)", date, date).
		Where("(min_amount IS NULL OR min_amount <= ?) AND (max_amount IS NULL OR max_amount >= ?)", req.Amount, req.Amount)
	
	// Location filtering
	if req.Country != "" {
		query = query.Where("countries IS NULL OR countries = '[]' OR countries @> ?", fmt.Sprintf(`["%s"]`, req.Country))
	}
	if req.State != "" {
		query = query.Where("states IS NULL OR states = '[]' OR states @> ?", fmt.Sprintf(`["%s"]`, req.State))
	}
	if req.City != "" {
		query = query.Where("cities IS NULL OR cities = '[]' OR cities @> ?", fmt.Sprintf(`["%s"]`, req.City))
	}
	
	// Product filtering
	if req.ProductID != nil {
		query = query.Where("product_ids IS NULL OR product_ids = '[]' OR product_ids @> ?", fmt.Sprintf(`["%s"]`, req.ProductID.String()))
	}
	
	// Customer filtering
	if req.CustomerID != nil {
		query = query.Where("customer_ids IS NULL OR customer_ids = '[]' OR customer_ids @> ?", fmt.Sprintf(`["%s"]`, req.CustomerID.String()))
	}
	
	var rules []*TaxRule
	err := query.Order("priority DESC, created_at ASC").Find(&rules).Error
	return rules, err
}

// Tax Rate operations

// CreateTaxRate creates a new tax rate
func (r *GormRepository) CreateTaxRate(ctx context.Context, rate *TaxRate) error {
	return r.db.WithContext(ctx).Create(rate).Error
}

// GetTaxRate retrieves a tax rate by ID
func (r *GormRepository) GetTaxRate(ctx context.Context, tenantID, rateID uuid.UUID) (*TaxRate, error) {
	var rate TaxRate
	err := r.db.WithContext(ctx).
		Preload("Rule").
		Where("tenant_id = ? AND id = ?", tenantID, rateID).
		First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

// UpdateTaxRate updates a tax rate
func (r *GormRepository) UpdateTaxRate(ctx context.Context, rate *TaxRate) error {
	return r.db.WithContext(ctx).Save(rate).Error
}

// DeleteTaxRate soft deletes a tax rate
func (r *GormRepository) DeleteTaxRate(ctx context.Context, tenantID, rateID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, rateID).
		Delete(&TaxRate{}).Error
}

// ListTaxRates retrieves tax rates with filtering and pagination
func (r *GormRepository) ListTaxRates(ctx context.Context, tenantID uuid.UUID, filter TaxRateFilter, page, pageSize int) ([]*TaxRate, int64, error) {
	query := r.db.WithContext(ctx).Model(&TaxRate{}).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filter.RuleID != nil {
		query = query.Where("rule_id = ?", *filter.RuleID)
	}
	if filter.TaxType != "" {
		query = query.Where("tax_type = ?", filter.TaxType)
	}
	if filter.Country != "" {
		query = query.Where("country = ?", filter.Country)
	}
	if filter.State != "" {
		query = query.Where("state = ?", filter.State)
	}
	if filter.City != "" {
		query = query.Where("city = ?", filter.City)
	}
	if filter.PostalCode != "" {
		query = query.Where("postal_code = ?", filter.PostalCode)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.ValidDate != nil {
		query = query.Where("(valid_from IS NULL OR valid_from <= ?) AND (valid_to IS NULL OR valid_to >= ?)", 
			filter.ValidDate, filter.ValidDate)
	}
	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}
	
	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination and get results
	var rates []*TaxRate
	offset := (page - 1) * pageSize
	err := query.Preload("Rule").
		Order("country, state, city, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&rates).Error
	
	return rates, total, err
}

// GetTaxRatesByRule retrieves tax rates for a specific rule
func (r *GormRepository) GetTaxRatesByRule(ctx context.Context, tenantID, ruleID uuid.UUID) ([]*TaxRate, error) {
	var rates []*TaxRate
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND rule_id = ?", tenantID, ruleID).
		Order("country, state, city, created_at DESC").
		Find(&rates).Error
	return rates, err
}

// GetActiveTaxRates retrieves active tax rates for a specific date
func (r *GormRepository) GetActiveTaxRates(ctx context.Context, tenantID uuid.UUID, date time.Time) ([]*TaxRate, error) {
	var rates []*TaxRate
	err := r.db.WithContext(ctx).
		Preload("Rule").
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Where("(valid_from IS NULL OR valid_from <= ?) AND (valid_to IS NULL OR valid_to >= ?)", date, date).
		Order("country, state, city, created_at DESC").
		Find(&rates).Error
	return rates, err
}

// Tax Calculation operations

// CreateTax creates a new tax calculation
func (r *GormRepository) CreateTax(ctx context.Context, tax *Tax) error {
	return r.db.WithContext(ctx).Create(tax).Error
}

// GetTax retrieves a tax calculation by ID
func (r *GormRepository) GetTax(ctx context.Context, tenantID, taxID uuid.UUID) (*Tax, error) {
	var tax Tax
	err := r.db.WithContext(ctx).
		Preload("AppliedRules").
		Preload("AppliedRules.Rule").
		Where("tenant_id = ? AND id = ?", tenantID, taxID).
		First(&tax).Error
	if err != nil {
		return nil, err
	}
	return &tax, nil
}

// ListTaxes retrieves tax calculations with filtering and pagination
func (r *GormRepository) ListTaxes(ctx context.Context, tenantID uuid.UUID, filter TaxFilter, page, pageSize int) ([]*Tax, int64, error) {
	query := r.db.WithContext(ctx).Model(&Tax{}).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filter.OrderID != nil {
		query = query.Where("order_id = ?", *filter.OrderID)
	}
	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.Country != "" {
		query = query.Where("country = ?", filter.Country)
	}
	if filter.State != "" {
		query = query.Where("state = ?", filter.State)
	}
	if filter.City != "" {
		query = query.Where("city = ?", filter.City)
	}
	if filter.TaxType != "" {
		query = query.Where("tax_type = ?", filter.TaxType)
	}
	if filter.Method != "" {
		query = query.Where("method = ?", filter.Method)
	}
	if filter.DateFrom != nil {
		query = query.Where("calculated_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("calculated_at <= ?", *filter.DateTo)
	}
	
	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination and get results
	var taxes []*Tax
	offset := (page - 1) * pageSize
	err := query.Preload("AppliedRules").
		Order("calculated_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&taxes).Error
	
	return taxes, total, err
}

// GetTaxesByOrder retrieves tax calculations for a specific order
func (r *GormRepository) GetTaxesByOrder(ctx context.Context, tenantID, orderID uuid.UUID) ([]*Tax, error) {
	var taxes []*Tax
	err := r.db.WithContext(ctx).
		Preload("AppliedRules").
		Where("tenant_id = ? AND order_id = ?", tenantID, orderID).
		Order("calculated_at DESC").
		Find(&taxes).Error
	return taxes, err
}

// GetTaxesByProduct retrieves tax calculations for a specific product
func (r *GormRepository) GetTaxesByProduct(ctx context.Context, tenantID, productID uuid.UUID) ([]*Tax, error) {
	var taxes []*Tax
	err := r.db.WithContext(ctx).
		Preload("AppliedRules").
		Where("tenant_id = ? AND product_id = ?", tenantID, productID).
		Order("calculated_at DESC").
		Find(&taxes).Error
	return taxes, err
}

// GetTaxesByCustomer retrieves tax calculations for a specific customer
func (r *GormRepository) GetTaxesByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*Tax, error) {
	var taxes []*Tax
	err := r.db.WithContext(ctx).
		Preload("AppliedRules").
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Order("calculated_at DESC").
		Find(&taxes).Error
	return taxes, err
}

// Tax Rule Application operations

// CreateTaxRuleApplication creates a new tax rule application
func (r *GormRepository) CreateTaxRuleApplication(ctx context.Context, application *TaxRuleApplication) error {
	return r.db.WithContext(ctx).Create(application).Error
}

// GetTaxRuleApplications retrieves tax rule applications for a tax calculation
func (r *GormRepository) GetTaxRuleApplications(ctx context.Context, tenantID, taxID uuid.UUID) ([]*TaxRuleApplication, error) {
	var applications []*TaxRuleApplication
	err := r.db.WithContext(ctx).
		Preload("Rule").
		Where("tenant_id = ? AND tax_id = ?", tenantID, taxID).
		Order("priority DESC, created_at ASC").
		Find(&applications).Error
	return applications, err
}

// Bulk operations

// BulkCreateTaxRules creates multiple tax rules
func (r *GormRepository) BulkCreateTaxRules(ctx context.Context, rules []*TaxRule) error {
	if len(rules) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(rules, 100).Error
}

// BulkUpdateTaxRuleStatus updates status for multiple tax rules
func (r *GormRepository) BulkUpdateTaxRuleStatus(ctx context.Context, tenantID uuid.UUID, ruleIDs []uuid.UUID, status string) error {
	if len(ruleIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&TaxRule{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ruleIDs).
		Update("status", status).Error
}

// BulkDeleteTaxRules soft deletes multiple tax rules
func (r *GormRepository) BulkDeleteTaxRules(ctx context.Context, tenantID uuid.UUID, ruleIDs []uuid.UUID) error {
	if len(ruleIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id IN ?", tenantID, ruleIDs).
		Delete(&TaxRule{}).Error
}

// BulkCreateTaxRates creates multiple tax rates
func (r *GormRepository) BulkCreateTaxRates(ctx context.Context, rates []*TaxRate) error {
	if len(rates) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(rates, 100).Error
}

// BulkUpdateTaxRateStatus updates status for multiple tax rates
func (r *GormRepository) BulkUpdateTaxRateStatus(ctx context.Context, tenantID uuid.UUID, rateIDs []uuid.UUID, isActive bool) error {
	if len(rateIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&TaxRate{}).
		Where("tenant_id = ? AND id IN ?", tenantID, rateIDs).
		Update("is_active", isActive).Error
}

// Validation operations

// CheckTaxRuleCodeExists checks if a tax rule code already exists
func (r *GormRepository) CheckTaxRuleCodeExists(ctx context.Context, tenantID uuid.UUID, code string, excludeID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).Model(&TaxRule{}).Where("tenant_id = ? AND code = ?", tenantID, code)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// ValidateTaxRuleReferences validates that referenced entities exist
func (r *GormRepository) ValidateTaxRuleReferences(ctx context.Context, tenantID uuid.UUID, productIDs, categoryIDs, customerIDs []uuid.UUID) error {
	// This would typically validate against actual product, category, and customer tables
	// For now, we'll just return nil as the validation would depend on other modules
	return nil
}

// Statistics and analytics

// GetTaxStats retrieves tax statistics
func (r *GormRepository) GetTaxStats(ctx context.Context, tenantID uuid.UUID) (*TaxStats, error) {
	stats := &TaxStats{}
	
	// Count tax rules
	r.db.WithContext(ctx).Model(&TaxRule{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalRules)
	r.db.WithContext(ctx).Model(&TaxRule{}).Where("tenant_id = ? AND status = ?", tenantID, StatusActive).Count(&stats.ActiveRules)
	r.db.WithContext(ctx).Model(&TaxRule{}).Where("tenant_id = ? AND status = ?", tenantID, StatusInactive).Count(&stats.InactiveRules)
	
	// Count tax rates
	r.db.WithContext(ctx).Model(&TaxRate{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalRates)
	r.db.WithContext(ctx).Model(&TaxRate{}).Where("tenant_id = ? AND is_active = ?", tenantID, true).Count(&stats.ActiveRates)
	
	// Count tax calculations
	r.db.WithContext(ctx).Model(&Tax{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalCalculations)
	
	// Calculate total tax amount and average rate
	var result struct {
		TotalTax float64
		AvgRate  float64
	}
	r.db.WithContext(ctx).Model(&Tax{}).
		Select("COALESCE(SUM(tax_amount), 0) as total_tax, COALESCE(AVG(tax_rate), 0) as avg_rate").
		Where("tenant_id = ?", tenantID).
		Scan(&result)
	
	stats.TotalTaxAmount = result.TotalTax
	stats.AverageTaxRate = result.AvgRate
	
	return stats, nil
}

// GetTaxStatsByLocation retrieves tax statistics by location
func (r *GormRepository) GetTaxStatsByLocation(ctx context.Context, tenantID uuid.UUID, limit int) ([]*TaxByLocation, error) {
	var stats []*TaxByLocation
	err := r.db.WithContext(ctx).Model(&Tax{}).
		Select("country, state, city, COUNT(*) as calculations, SUM(tax_amount) as total_tax, AVG(tax_rate) as average_rate").
		Where("tenant_id = ?", tenantID).
		Group("country, state, city").
		Order("total_tax DESC").
		Limit(limit).
		Scan(&stats).Error
	return stats, err
}

// GetTaxStatsByType retrieves tax statistics by type
func (r *GormRepository) GetTaxStatsByType(ctx context.Context, tenantID uuid.UUID) ([]*TaxByType, error) {
	var stats []*TaxByType
	err := r.db.WithContext(ctx).Model(&Tax{}).
		Select("tax_type, COUNT(*) as calculations, SUM(tax_amount) as total_tax, AVG(tax_rate) as average_rate").
		Where("tenant_id = ?", tenantID).
		Group("tax_type").
		Order("total_tax DESC").
		Scan(&stats).Error
	return stats, err
}

// GetTaxTrends retrieves tax trends over time
func (r *GormRepository) GetTaxTrends(ctx context.Context, tenantID uuid.UUID, days int) (map[string]interface{}, error) {
	sinceDate := time.Now().AddDate(0, 0, -days)
	
	var dailyStats []struct {
		Date         string  `json:"date"`
		Calculations int64   `json:"calculations"`
		TotalTax     float64 `json:"total_tax"`
		AvgRate      float64 `json:"avg_rate"`
	}
	
	err := r.db.WithContext(ctx).Model(&Tax{}).
		Select("DATE(calculated_at) as date, COUNT(*) as calculations, SUM(tax_amount) as total_tax, AVG(tax_rate) as avg_rate").
		Where("tenant_id = ? AND calculated_at >= ?", tenantID, sinceDate).
		Group("DATE(calculated_at)").
		Order("date").
		Scan(&dailyStats).Error
	
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"daily_stats": dailyStats,
		"period_days": days,
		"since_date":  sinceDate,
	}, nil
}

// Maintenance operations

// CleanupExpiredRules soft deletes expired tax rules
func (r *GormRepository) CleanupExpiredRules(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&TaxRule{}).
		Where("tenant_id = ? AND valid_to IS NOT NULL AND valid_to < ? AND status != ?", tenantID, now, StatusArchived).
		Update("status", StatusArchived)
	return result.RowsAffected, result.Error
}

// CleanupExpiredRates deactivates expired tax rates
func (r *GormRepository) CleanupExpiredRates(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&TaxRate{}).
		Where("tenant_id = ? AND valid_to IS NOT NULL AND valid_to < ? AND is_active = ?", tenantID, now, true).
		Update("is_active", false)
	return result.RowsAffected, result.Error
}

// ArchiveOldTaxCalculations archives old tax calculations
func (r *GormRepository) ArchiveOldTaxCalculations(ctx context.Context, tenantID uuid.UUID, olderThan time.Time) (int64, error) {
	// In a real implementation, this might move records to an archive table
	// For now, we'll just count how many would be archived
	var count int64
	err := r.db.WithContext(ctx).Model(&Tax{}).
		Where("tenant_id = ? AND calculated_at < ?", tenantID, olderThan).
		Count(&count).Error
	return count, err
}