package tenant

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles tenant data operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new tenant repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Save creates a new tenant in the database
func (r *Repository) Save(tenant *Tenant) (*Tenant, error) {
	if err := r.db.Create(tenant).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

// FindByID retrieves a tenant by ID
func (r *Repository) FindByID(id uuid.UUID) (*Tenant, error) {
	var tenant Tenant
	if err := r.db.First(&tenant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// FindBySubdomain retrieves a tenant by subdomain
func (r *Repository) FindBySubdomain(subdomain string) (*Tenant, error) {
	var tenant Tenant
	if err := r.db.First(&tenant, "subdomain = ?", subdomain).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// FindByCustomDomain retrieves a tenant by custom domain
func (r *Repository) FindByCustomDomain(domain string) (*Tenant, error) {
	var tenant Tenant
	if err := r.db.First(&tenant, "custom_domain = ?", domain).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// Update updates an existing tenant
func (r *Repository) Update(tenant *Tenant) (*Tenant, error) {
	if err := r.db.Save(tenant).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

// UpdateStatus updates only the status field
func (r *Repository) UpdateStatus(id uuid.UUID, status Status) error {
	return r.db.Model(&Tenant{}).Where("id = ?", id).Update("status", status).Error
}

// SubdomainExists checks if a subdomain is already taken
func (r *Repository) SubdomainExists(subdomain string) (bool, error) {
	var count int64
	err := r.db.Model(&Tenant{}).Where("subdomain = ?", subdomain).Count(&count).Error
	return count > 0, err
}

// CustomDomainExists checks if a custom domain is already taken
func (r *Repository) CustomDomainExists(domain string) (bool, error) {
	var count int64
	err := r.db.Model(&Tenant{}).Where("custom_domain = ?", domain).Count(&count).Error
	return count > 0, err
}

// List retrieves tenants with pagination
func (r *Repository) List(offset, limit int) ([]*Tenant, int64, error) {
	var tenants []*Tenant
	var total int64

	// Get total count
	if err := r.db.Model(&Tenant{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}

// ListByStatus retrieves tenants by status with pagination
func (r *Repository) ListByStatus(status Status, offset, limit int) ([]*Tenant, int64, error) {
	var tenants []*Tenant
	var total int64

	query := r.db.Model(&Tenant{}).Where("status = ?", status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}

// ListByPlan retrieves tenants by plan with pagination
func (r *Repository) ListByPlan(plan Plan, offset, limit int) ([]*Tenant, int64, error) {
	var tenants []*Tenant
	var total int64

	query := r.db.Model(&Tenant{}).Where("plan = ?", plan)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}

// Delete soft deletes a tenant
func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Tenant{}, "id = ?", id).Error
}

// GetActiveCount returns the count of active tenants
func (r *Repository) GetActiveCount() (int64, error) {
	var count int64
	err := r.db.Model(&Tenant{}).Where("status = ?", StatusActive).Count(&count).Error
	return count, err
}

// GetCountByPlan returns the count of tenants by plan
func (r *Repository) GetCountByPlan(plan Plan) (int64, error) {
	var count int64
	err := r.db.Model(&Tenant{}).Where("plan = ?", plan).Count(&count).Error
	return count, err
}

// Search searches tenants by name or subdomain
func (r *Repository) Search(query string, offset, limit int) ([]*Tenant, int64, error) {
	var tenants []*Tenant
	var total int64

	searchQuery := r.db.Model(&Tenant{}).Where(
		"name ILIKE ? OR subdomain ILIKE ? OR email ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%",
	)

	// Get total count
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := searchQuery.Offset(offset).Limit(limit).Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}

// GetExpiringTrials returns tenants whose trials are expiring soon
func (r *Repository) GetExpiringTrials(days int) ([]*Tenant, error) {
	// TODO: Implement trial logic when trial system is added
	var tenants []*Tenant
	return tenants, nil
}

// BulkUpdateStatus updates status for multiple tenants
func (r *Repository) BulkUpdateStatus(ids []uuid.UUID, status Status) error {
	return r.db.Model(&Tenant{}).Where("id IN ?", ids).Update("status", status).Error
}

// GetTenantStats returns statistics for a specific tenant
func (r *Repository) GetTenantStats(tenantID uuid.UUID) (*TenantStatsResponse, error) {
	// TODO: This would require joins with products, orders, etc.
	// For now, return basic stats
	stats := &TenantStatsResponse{
		TenantID:     tenantID.String(),
		ProductCount: 0,
		OrderCount:   0,
		Revenue:      0,
		StorageUsed:  0,
		BandwidthUsed: 0,
	}
	
	// You can add actual database queries here when product/order modules are integrated
	// Example:
	// r.db.Model(&Product{}).Where("tenant_id = ?", tenantID).Count(&stats.ProductCount)
	// r.db.Model(&Order{}).Where("tenant_id = ?", tenantID).Count(&stats.OrderCount)
	
	return stats, nil
}

// UpdateCustomDomain updates the custom domain for a tenant
func (r *Repository) UpdateCustomDomain(id uuid.UUID, domain string) error {
	return r.db.Model(&Tenant{}).Where("id = ?", id).Update("custom_domain", domain).Error
}

// TODO: Add more repository methods as needed
// - UpdateUsageMetrics(id uuid.UUID, storage, bandwidth int) error
// - GetTenantsByDateRange(start, end time.Time) ([]*Tenant, error)
// - GetRevenueStats() (*RevenueStats, error)
