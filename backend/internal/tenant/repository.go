package tenant

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sharedErrors "ecommerce-saas/internal/shared/errors"
)

// RepositoryImpl handles tenant data operations
type RepositoryImpl struct {
	db *gorm.DB
}

// NewRepository creates a new tenant repository
func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

// Save creates a new tenant in the database
func (r *RepositoryImpl) Save(tenant *Tenant) (*Tenant, error) {
	if err := r.db.Create(tenant).Error; err != nil {
		// Check for unique constraint violations
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			if strings.Contains(err.Error(), "subdomain") {
				return nil, sharedErrors.NewConflictError("Subdomain already exists")
			}
			if strings.Contains(err.Error(), "custom_domain") {
				return nil, sharedErrors.NewConflictError("Custom domain already exists")
			}
			return nil, sharedErrors.NewConflictError("Tenant already exists")
		}
		return nil, sharedErrors.NewInternalError("Failed to create tenant", err)
	}
	return tenant, nil
}

// FindByID retrieves a tenant by ID
func (r *RepositoryImpl) FindByID(id uuid.UUID) (*Tenant, error) {
	var tenant Tenant
	if err := r.db.First(&tenant, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("Tenant")
		}
		return nil, sharedErrors.NewInternalError("Failed to retrieve tenant", err)
	}
	return &tenant, nil
}

// FindBySubdomain retrieves a tenant by subdomain
func (r *RepositoryImpl) FindBySubdomain(subdomain string) (*Tenant, error) {
	var tenant Tenant
	if err := r.db.First(&tenant, "subdomain = ?", subdomain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("Tenant")
		}
		return nil, sharedErrors.NewInternalError("Failed to retrieve tenant by subdomain", err)
	}
	return &tenant, nil
}

// FindByCustomDomain retrieves a tenant by custom domain
func (r *RepositoryImpl) FindByCustomDomain(domain string) (*Tenant, error) {
	var tenant Tenant
	if err := r.db.First(&tenant, "custom_domain = ?", domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewNotFoundError("Tenant not found")
		}
		return nil, NewInternalError("Failed to retrieve tenant by custom domain")
	}
	return &tenant, nil
}

// Update updates an existing tenant
func (r *RepositoryImpl) Update(tenant *Tenant) (*Tenant, error) {
	if err := r.db.Save(tenant).Error; err != nil {
		// Check for unique constraint violations
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			if strings.Contains(err.Error(), "subdomain") {
				return nil, NewConflictError("Subdomain already exists")
			}
			if strings.Contains(err.Error(), "custom_domain") {
				return nil, NewConflictError("Custom domain already exists")
			}
		}
		return nil, NewInternalError("Failed to update tenant")
	}
	return tenant, nil
}

// UpdateStatus updates only the status field
func (r *RepositoryImpl) UpdateStatus(id uuid.UUID, status Status) error {
	result := r.db.Model(&Tenant{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return NewInternalError("Failed to update tenant status")
	}
	if result.RowsAffected == 0 {
		return NewNotFoundError("Tenant not found")
	}
	return nil
}

// SubdomainExists checks if a subdomain is already taken
func (r *RepositoryImpl) SubdomainExists(subdomain string) (bool, error) {
	var count int64
	if err := r.db.Model(&Tenant{}).Where("subdomain = ?", subdomain).Count(&count).Error; err != nil {
		return false, NewInternalError("Failed to check subdomain existence")
	}
	return count > 0, nil
}

// CustomDomainExists checks if a custom domain is already taken
func (r *RepositoryImpl) CustomDomainExists(domain string) (bool, error) {
	var count int64
	if err := r.db.Model(&Tenant{}).Where("custom_domain = ?", domain).Count(&count).Error; err != nil {
		return false, NewInternalError("Failed to check custom domain existence")
	}
	return count > 0, nil
}

// List retrieves tenants with pagination
func (r *RepositoryImpl) List(offset, limit int) ([]*Tenant, int64, error) {
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
func (r *RepositoryImpl) ListByStatus(status Status, offset, limit int) ([]*Tenant, int64, error) {
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
func (r *RepositoryImpl) ListByPlan(plan Plan, offset, limit int) ([]*Tenant, int64, error) {
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
func (r *RepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&Tenant{}, "id = ?", id).Error
}

// GetActiveCount returns the count of active tenants
func (r *RepositoryImpl) GetActiveCount() (int64, error) {
	var count int64
	err := r.db.Model(&Tenant{}).Where("status = ?", StatusActive).Count(&count).Error
	return count, err
}

// GetCountByPlan returns the count of tenants by plan
func (r *RepositoryImpl) GetCountByPlan(plan Plan) (int64, error) {
	var count int64
	err := r.db.Model(&Tenant{}).Where("plan = ?", plan).Count(&count).Error
	return count, err
}

// Search searches tenants by name or subdomain
func (r *RepositoryImpl) Search(query string, offset, limit int) ([]*Tenant, int64, error) {
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
func (r *RepositoryImpl) GetExpiringTrials(days int) ([]*Tenant, error) {
	var tenants []*Tenant
	expiryDate := time.Now().AddDate(0, 0, days)
	
	err := r.db.Where("is_trial_active = ? AND trial_end_date <= ? AND trial_end_date > ?", 
		true, expiryDate, time.Now()).Find(&tenants).Error
	if err != nil {
		return nil, NewInternalError("failed to get expiring trials")
	}
	
	return tenants, nil
}

// BulkUpdateStatus updates status for multiple tenants
func (r *RepositoryImpl) BulkUpdateStatus(ids []uuid.UUID, status Status) error {
	return r.db.Model(&Tenant{}).Where("id IN ?", ids).Update("status", status).Error
}

// GetTenantStats returns statistics for a specific tenant
func (r *RepositoryImpl) GetTenantStats(tenantID uuid.UUID) (*TenantStatsResponse, error) {
	// TODO: This would require joins with products, orders, etc.
	// For now, return basic stats
	stats := &TenantStatsResponse{
		TenantID:      tenantID.String(),
		ProductCount:  0,
		OrderCount:    0,
		Revenue:       0,
		StorageUsed:   0,
		BandwidthUsed: 0,
	}

	// You can add actual database queries here when product/order modules are integrated
	// Example:
	// r.db.Model(&Product{}).Where("tenant_id = ?", tenantID).Count(&stats.ProductCount)
	// r.db.Model(&Order{}).Where("tenant_id = ?", tenantID).Count(&stats.OrderCount)

	return stats, nil
}

// UpdateCustomDomain updates the custom domain for a tenant
func (r *RepositoryImpl) UpdateCustomDomain(id uuid.UUID, domain string) error {
	return r.db.Model(&Tenant{}).Where("id = ?", id).Update("custom_domain", domain).Error
}

// TODO: Add more repository methods as needed
// - UpdateUsageMetrics(id uuid.UUID, storage, bandwidth int) error
// - GetTenantsByDateRange(start, end time.Time) ([]*Tenant, error)
// - GetRevenueStats() (*RevenueStats, error)
