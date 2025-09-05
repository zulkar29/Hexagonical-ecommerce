package admin

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the admin repository interface
type Repository interface {
	// Staff management
	GetStaff(ctx context.Context, tenantID *uuid.UUID, role, status string) ([]*Staff, error)
	CreateStaff(ctx context.Context, staff *Staff) (*Staff, error)
	UpdateStaff(ctx context.Context, tenantID *uuid.UUID, id uuid.UUID, req StaffRequest) (*Staff, error)
	DeleteStaff(ctx context.Context, tenantID *uuid.UUID, id uuid.UUID) error
	AssignRolesToStaff(ctx context.Context, tenantID *uuid.UUID, staffID uuid.UUID, roles []string) error
	UpdateStaffStatus(ctx context.Context, tenantID *uuid.UUID, staffID uuid.UUID, status string) error
	
	// Role management
	GetRoles(ctx context.Context, tenantID *uuid.UUID, includePermissions bool) ([]*Role, error)
	CreateRole(ctx context.Context, role *Role, permissions []string) (*Role, error)
	UpdateRole(ctx context.Context, tenantID *uuid.UUID, id uuid.UUID, req RoleRequest) (*Role, error)
	DeleteRole(ctx context.Context, tenantID *uuid.UUID, id uuid.UUID) error
	AssignPermissionsToRole(ctx context.Context, tenantID *uuid.UUID, roleID uuid.UUID, permissions []string) error
	
	// Activity logs
	GetActivityLogs(ctx context.Context, tenantID *uuid.UUID, filter ActivityLogFilter) ([]*ActivityLog, error)
	CreateActivityLog(ctx context.Context, log *ActivityLog) error
}

// RepositoryImpl implements the admin repository using GORM
type RepositoryImpl struct {
	db *gorm.DB
}

// NewRepository creates a new admin repository
func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

// GetStaff retrieves staff members with optional filtering
func (r *RepositoryImpl) GetStaff(ctx context.Context, tenantID *uuid.UUID, role, status string) ([]*Staff, error) {
	var staff []*Staff
	
	query := r.db.WithContext(ctx)
	
	// Apply tenant filter if provided
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	// Apply role filter if provided
	if role != "" {
		query = query.Where("role = ?", role)
	}
	
	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	// TODO: This should query from the users table with admin/staff roles
	// For now, we'll create a placeholder query
	// In reality, this would be something like:
	// query = query.Table("users").Where("role IN ?", []string{"admin", "staff", "manager"})
	
	err := query.Find(&staff).Error
	return staff, err
}

// CreateStaff creates a new staff member
func (r *RepositoryImpl) CreateStaff(ctx context.Context, staff *Staff) (*Staff, error) {
	// TODO: This should create a user record with staff role
	// For now, we'll use a placeholder implementation
	
	err := r.db.WithContext(ctx).Create(staff).Error
	if err != nil {
		return nil, err
	}
	
	return staff, nil
}

// UpdateStaff updates an existing staff member
func (r *RepositoryImpl) UpdateStaff(ctx context.Context, tenantID *uuid.UUID, id uuid.UUID, req StaffRequest) (*Staff, error) {
	var staff Staff
	
	// Find the staff member
	query := r.db.WithContext(ctx).Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	err := query.First(&staff).Error
	if err != nil {
		return nil, err
	}
	
	// Update fields
	staff.Email = req.Email
	staff.FirstName = req.FirstName
	staff.LastName = req.LastName
	staff.Role = req.Role
	staff.UpdatedAt = time.Now()
	
	err = r.db.WithContext(ctx).Save(&staff).Error
	if err != nil {
		return nil, err
	}
	
	return &staff, nil
}

// DeleteStaff deletes a staff member (soft delete)
func (r *RepositoryImpl) DeleteStaff(ctx context.Context, tenantID *uuid.UUID, id uuid.UUID) error {
	query := r.db.WithContext(ctx).Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	// Soft delete
	return query.Delete(&Staff{}).Error
}

// AssignRolesToStaff assigns roles to a staff member
func (r *RepositoryImpl) AssignRolesToStaff(ctx context.Context, tenantID *uuid.UUID, staffID uuid.UUID, roles []string) error {
	// TODO: Implement role assignment
	// This would typically involve:
	// 1. Clear existing role assignments
	// 2. Create new role assignments
	// 3. Update user role field
	
	// Placeholder implementation
	return nil
}

// UpdateStaffStatus updates the status of a staff member
func (r *RepositoryImpl) UpdateStaffStatus(ctx context.Context, tenantID *uuid.UUID, staffID uuid.UUID, status string) error {
	query := r.db.WithContext(ctx).Where("id = ?", staffID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	return query.Update("status", status).Error
}

// GetRoles retrieves roles with optional permissions
func (r *RepositoryImpl) GetRoles(ctx context.Context, tenantID *uuid.UUID, includePermissions bool) ([]*Role, error) {
	var roles []*Role
	
	query := r.db.WithContext(ctx)
	
	// Apply tenant filter if provided
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	// Include permissions if requested
	if includePermissions {
		query = query.Preload("Permissions")
	}
	
	err := query.Find(&roles).Error
	return roles, err
}

// CreateRole creates a new role with permissions
func (r *RepositoryImpl) CreateRole(ctx context.Context, role *Role, permissions []string) (*Role, error) {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// Create the role
	err := tx.Create(role).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// TODO: Assign permissions to the role
	// This would involve creating role_permission records
	
	return role, tx.Commit().Error
}

// UpdateRole updates an existing role
func (r *RepositoryImpl) UpdateRole(ctx context.Context, tenantID *uuid.UUID, id uuid.UUID, req RoleRequest) (*Role, error) {
	var role Role
	
	// Find the role
	query := r.db.WithContext(ctx).Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	err := query.First(&role).Error
	if err != nil {
		return nil, err
	}
	
	// Update fields
	role.Name = req.Name
	role.Description = req.Description
	role.UpdatedAt = time.Now()
	
	err = r.db.WithContext(ctx).Save(&role).Error
	if err != nil {
		return nil, err
	}
	
	// TODO: Update permissions if provided
	
	return &role, nil
}

// DeleteRole deletes a role
func (r *RepositoryImpl) DeleteRole(ctx context.Context, tenantID *uuid.UUID, id uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// TODO: Check if role is in use
	// TODO: Delete role permissions first
	
	// Delete the role
	query := tx.Where("id = ?", id)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	err := query.Delete(&Role{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit().Error
}

// AssignPermissionsToRole assigns permissions to a role
func (r *RepositoryImpl) AssignPermissionsToRole(ctx context.Context, tenantID *uuid.UUID, roleID uuid.UUID, permissions []string) error {
	// TODO: Implement permission assignment
	// This would typically involve:
	// 1. Clear existing permission assignments for the role
	// 2. Create new permission assignments
	
	// Placeholder implementation
	return nil
}

// GetActivityLogs retrieves activity logs with filtering
func (r *RepositoryImpl) GetActivityLogs(ctx context.Context, tenantID *uuid.UUID, filter ActivityLogFilter) ([]*ActivityLog, error) {
	var logs []*ActivityLog
	
	query := r.db.WithContext(ctx)
	
	// Apply tenant filter if provided
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	
	// Apply user filter if provided
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	
	// Apply action filter if provided
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	
	// Apply date filters if provided
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	
	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	
	// Order by created_at descending
	query = query.Order("created_at DESC")
	
	err := query.Find(&logs).Error
	return logs, err
}

// CreateActivityLog creates a new activity log entry
func (r *RepositoryImpl) CreateActivityLog(ctx context.Context, log *ActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}