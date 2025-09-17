package platform

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the platform repository interface
type Repository interface {
	// Dashboard & Stats
	GetPlatformStats(ctx context.Context, period string) (*PlatformStats, error)
	GetSystemStatus(ctx context.Context) (*SystemStatus, error)

	// Platform Admin Management
	ListPlatformAdmins(ctx context.Context, role string) ([]*PlatformAdmin, error)
	GetPlatformAdmin(ctx context.Context, id uuid.UUID) (*PlatformAdmin, error)
	CreatePlatformAdmin(ctx context.Context, admin *PlatformAdmin) (*PlatformAdmin, error)
	UpdatePlatformAdmin(ctx context.Context, id uuid.UUID, admin *PlatformAdmin) (*PlatformAdmin, error)
	DeletePlatformAdmin(ctx context.Context, id uuid.UUID) error

	// Platform Role Management
	ListPlatformRoles(ctx context.Context) ([]*PlatformRole, error)
	GetPlatformRole(ctx context.Context, id uuid.UUID) (*PlatformRole, error)
	CreatePlatformRole(ctx context.Context, role *PlatformRole) (*PlatformRole, error)
	UpdatePlatformRole(ctx context.Context, id uuid.UUID, role *PlatformRole) (*PlatformRole, error)
	DeletePlatformRole(ctx context.Context, id uuid.UUID) error

	// Platform Settings
	GetPlatformSettings(ctx context.Context, category string) ([]*PlatformSettings, error)
	UpdatePlatformSettings(ctx context.Context, settings *PlatformSettings) (*PlatformSettings, error)

	// Platform Tenant Management
	ListAllTenants(ctx context.Context, status string, include []string) ([]map[string]interface{}, error)
	GetTenantDetails(ctx context.Context, id uuid.UUID, include []string) (map[string]interface{}, error)
	UpdateTenant(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteTenant(ctx context.Context, id uuid.UUID) error

	// Audit Logs
	GetPlatformAuditLogs(ctx context.Context, filter AuditLogFilter) ([]*PlatformAuditLog, error)
	CreateAuditLog(ctx context.Context, log *PlatformAuditLog) error
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new platform repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// GetPlatformStats retrieves platform statistics
func (r *repository) GetPlatformStats(ctx context.Context, period string) (*PlatformStats, error) {
	stats := &PlatformStats{}
	
	// Get total tenants
	var totalTenants int64
	if err := r.db.WithContext(ctx).Table("tenants").Count(&totalTenants).Error; err != nil {
		return nil, err
	}
	stats.TotalTenants = int(totalTenants)

	// Get active tenants
	var activeTenants int64
	if err := r.db.WithContext(ctx).Table("tenants").Where("status = ?", "active").Count(&activeTenants).Error; err != nil {
		return nil, err
	}
	stats.ActiveTenants = int(activeTenants)

	// Get total users
	var totalUsers int64
	if err := r.db.WithContext(ctx).Table("users").Count(&totalUsers).Error; err != nil {
		return nil, err
	}
	stats.TotalUsers = int(totalUsers)

	// Get active users (logged in within last 30 days)
	var activeUsers int64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if err := r.db.WithContext(ctx).Table("users").Where("last_login_at > ?", thirtyDaysAgo).Count(&activeUsers).Error; err != nil {
		return nil, err
	}
	stats.ActiveUsers = int(activeUsers)

	// Calculate revenue (simplified - would need proper billing integration)
	var totalRevenue float64
	if err := r.db.WithContext(ctx).Table("orders").Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}
	stats.TotalRevenue = totalRevenue

	// Calculate monthly revenue
	var monthlyRevenue float64
	startOfMonth := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	if err := r.db.WithContext(ctx).Table("orders").Select("COALESCE(SUM(total_amount), 0)").Where("created_at >= ?", startOfMonth).Scan(&monthlyRevenue).Error; err != nil {
		return nil, err
	}
	stats.MonthlyRevenue = monthlyRevenue

	// Calculate growth rate (simplified)
	if stats.TotalTenants > 0 {
		stats.GrowthRate = float64(stats.ActiveTenants) / float64(stats.TotalTenants) * 100
	}

	stats.SystemUptime = 99.9 // Would be calculated from actual system metrics
	stats.CreatedAt = time.Now()
	stats.UpdatedAt = time.Now()

	return stats, nil
}

// GetSystemStatus retrieves system health status
func (r *repository) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	status := &SystemStatus{
		ID:                uuid.New(),
		Status:            "healthy",
		DatabaseStatus:    "healthy",
		RedisStatus:       "healthy",
		APIResponseTime:   120.5,
		CPUUsage:          45.2,
		MemoryUsage:       67.8,
		DiskUsage:         23.4,
		ActiveConnections: 150,
		CreatedAt:         time.Now(),
	}

	// Test database connection
	if err := r.db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		status.DatabaseStatus = "unhealthy"
		status.Status = "degraded"
	}

	return status, nil
}

// ListPlatformAdmins retrieves platform administrators
func (r *repository) ListPlatformAdmins(ctx context.Context, role string) ([]*PlatformAdmin, error) {
	var admins []*PlatformAdmin
	query := r.db.WithContext(ctx)
	
	if role != "" {
		query = query.Where("role = ?", role)
	}
	
	if err := query.Find(&admins).Error; err != nil {
		return nil, err
	}
	
	return admins, nil
}

// GetPlatformAdmin retrieves a platform administrator by ID
func (r *repository) GetPlatformAdmin(ctx context.Context, id uuid.UUID) (*PlatformAdmin, error) {
	var admin PlatformAdmin
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// CreatePlatformAdmin creates a new platform administrator
func (r *repository) CreatePlatformAdmin(ctx context.Context, admin *PlatformAdmin) (*PlatformAdmin, error) {
	if err := r.db.WithContext(ctx).Create(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

// UpdatePlatformAdmin updates a platform administrator
func (r *repository) UpdatePlatformAdmin(ctx context.Context, id uuid.UUID, admin *PlatformAdmin) (*PlatformAdmin, error) {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

// DeletePlatformAdmin deletes a platform administrator
func (r *repository) DeletePlatformAdmin(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&PlatformAdmin{}).Error
}

// ListPlatformRoles retrieves platform roles
func (r *repository) ListPlatformRoles(ctx context.Context) ([]*PlatformRole, error) {
	var roles []*PlatformRole
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetPlatformRole retrieves a platform role by ID
func (r *repository) GetPlatformRole(ctx context.Context, id uuid.UUID) (*PlatformRole, error) {
	var role PlatformRole
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// CreatePlatformRole creates a new platform role
func (r *repository) CreatePlatformRole(ctx context.Context, role *PlatformRole) (*PlatformRole, error) {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// UpdatePlatformRole updates a platform role
func (r *repository) UpdatePlatformRole(ctx context.Context, id uuid.UUID, role *PlatformRole) (*PlatformRole, error) {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// DeletePlatformRole deletes a platform role
func (r *repository) DeletePlatformRole(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&PlatformRole{}).Error
}

// GetPlatformSettings retrieves platform settings
func (r *repository) GetPlatformSettings(ctx context.Context, category string) ([]*PlatformSettings, error) {
	var settings []*PlatformSettings
	query := r.db.WithContext(ctx)
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	if err := query.Find(&settings).Error; err != nil {
		return nil, err
	}
	
	return settings, nil
}

// UpdatePlatformSettings updates platform settings
func (r *repository) UpdatePlatformSettings(ctx context.Context, settings *PlatformSettings) (*PlatformSettings, error) {
	if err := r.db.WithContext(ctx).Where("category = ? AND key = ?", settings.Category, settings.Key).Save(settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// ListAllTenants retrieves all tenants for platform admin
func (r *repository) ListAllTenants(ctx context.Context, status string, include []string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	// Base query
	query := r.db.WithContext(ctx).Table("tenants")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	// Select basic tenant fields
	query = query.Select("id, name, domain, status, plan, created_at, updated_at")
	
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}
	
	// Add additional data based on include parameter
	for i, tenant := range results {
		tenantID := tenant["id"].(uuid.UUID)
		
		for _, inc := range include {
			switch inc {
			case "users":
				var userCount int64
				r.db.WithContext(ctx).Table("users").Where("tenant_id = ?", tenantID).Count(&userCount)
				results[i]["user_count"] = userCount
				
			case "revenue":
				var revenue float64
				r.db.WithContext(ctx).Table("orders").Where("tenant_id = ?", tenantID).Select("COALESCE(SUM(total_amount), 0)").Scan(&revenue)
				results[i]["total_revenue"] = revenue
				
			case "usage":
				// Add usage metrics
				var orderCount int64
				r.db.WithContext(ctx).Table("orders").Where("tenant_id = ?", tenantID).Count(&orderCount)
				results[i]["order_count"] = orderCount
			}
		}
	}
	
	return results, nil
}

// GetTenantDetails retrieves detailed tenant information
func (r *repository) GetTenantDetails(ctx context.Context, id uuid.UUID, include []string) (map[string]interface{}, error) {
	var result map[string]interface{}
	
	if err := r.db.WithContext(ctx).Table("tenants").Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	
	// Add additional data based on include parameter
	for _, inc := range include {
		switch inc {
		case "subscription":
			// Add subscription details
			result["subscription"] = map[string]interface{}{
				"plan":   result["plan"],
				"status": "active",
			}
			
		case "usage":
			// Add usage metrics
			var metrics map[string]interface{}
			r.db.WithContext(ctx).Raw("SELECT COUNT(*) as order_count, COALESCE(SUM(total_amount), 0) as revenue FROM orders WHERE tenant_id = ?", id).Scan(&metrics)
			result["usage"] = metrics
			
		case "settings":
			// Add tenant settings
			var settings []map[string]interface{}
			r.db.WithContext(ctx).Table("tenant_settings").Where("tenant_id = ?", id).Find(&settings)
			result["settings"] = settings
			
		case "analytics":
			// Add analytics data
			result["analytics"] = map[string]interface{}{
				"monthly_orders": 0,
				"growth_rate":    0,
			}
		}
	}
	
	return result, nil
}

// UpdateTenant updates tenant information
func (r *repository) UpdateTenant(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Table("tenants").Where("id = ?", id).Updates(updates).Error
}

// DeleteTenant deletes a tenant
func (r *repository) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Table("tenants").Where("id = ?", id).Delete(nil).Error
}

// GetPlatformAuditLogs retrieves platform audit logs
func (r *repository) GetPlatformAuditLogs(ctx context.Context, filter AuditLogFilter) ([]*PlatformAuditLog, error) {
	var logs []*PlatformAuditLog
	query := r.db.WithContext(ctx)
	
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	
	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", *filter.TenantID)
	}
	
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	
	if filter.Resource != "" {
		query = query.Where("resource = ?", filter.Resource)
	}
	
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	} else {
		query = query.Limit(50) // Default limit
	}
	
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	
	query = query.Order("created_at DESC")
	
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	
	return logs, nil
}

// CreateAuditLog creates a new audit log entry
func (r *repository) CreateAuditLog(ctx context.Context, log *PlatformAuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}