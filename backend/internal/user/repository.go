package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sharedErrors "ecommerce-saas/internal/shared/errors"
)

// Repository defines user repository interface
type Repository interface {
	// User CRUD operations
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*User, error) // Alias for GetUserByID
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	ListUsers(ctx context.Context, tenantID *uuid.UUID, filter UserFilter, offset, limit int) ([]*User, int64, error)
	UpdateUserByAdmin(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*User, error)
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]User, error)
	GetUsersWithFilters(ctx context.Context, filters interface{}) ([]User, error)
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) error
	UpdatePhoneVerification(ctx context.Context, phone string, verified bool) error
	UpdateUser2FA(ctx context.Context, userID uuid.UUID, enabled bool) error

	// Permission operations
	GetUserPermissions(userID uuid.UUID) ([]*Permission, error)
	CheckUserPermission(userID uuid.UUID, resource, action string) (bool, error)

	// Note: User preferences functionality has been removed due to JSONB complexity
}

// repository implements Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// User CRUD operations

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, sharedErrors.NewConflictError("User with this email already exists")
		}
		return nil, sharedErrors.NewInternalError("Failed to create user", err)
	}
	return user, nil
}

func (r *repository) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("User")
		}
		return nil, sharedErrors.NewInternalError("Failed to retrieve user", err)
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewNotFoundError("User")
		}
		return nil, sharedErrors.NewInternalError("Failed to retrieve user by email", err)
	}
	return &user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) (*User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&User{}, userID).Error
}

func (r *repository) ListUsers(ctx context.Context, tenantID *uuid.UUID, filter UserFilter, offset, limit int) ([]*User, int64, error) {
	var users []*User
	var total int64

	query := r.db.WithContext(ctx).Model(&User{})

	// Apply filters
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", *filter.TenantID)
	}
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *repository) UpdateUserByAdmin(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*User, error) {
	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return nil, err
	}

	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]User, error) {
	var users []User
	err := r.db.WithContext(ctx).Where("id IN ?", userIDs).Find(&users).Error
	return users, err
}

func (r *repository) GetUsersWithFilters(ctx context.Context, filters interface{}) ([]User, error) {
	var users []User
	query := r.db.WithContext(ctx)

	// Apply filters based on the filters interface
	if filterMap, ok := filters.(map[string]interface{}); ok {
		for key, value := range filterMap {
			switch key {
			case "status":
				query = query.Where("status = ?", value)
			case "created_after":
				query = query.Where("created_at > ?", value)
			case "created_before":
				query = query.Where("created_at < ?", value)
			}
		}
	}

	err := query.Find(&users).Error
	return users, err
}

func (r *repository) UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("status", status).Error
}

func (r *repository) UpdatePhoneVerification(ctx context.Context, phone string, verified bool) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("phone = ?", phone).Update("phone_verified", verified).Error
}

// Permission operations

func (r *repository) GetUserPermissions(userID uuid.UUID) ([]*Permission, error) {
	var permissions []*Permission

	// Get user role
	var user User
	if err := r.db.Select("role").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// Get permissions for the role
	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role = ?", user.Role).
		Find(&permissions).Error

	return permissions, err
}

func (r *repository) CheckUserPermission(userID uuid.UUID, resource, action string) (bool, error) {
	var count int64

	// Get user role
	var user User
	if err := r.db.Select("role").Where("id = ?", userID).First(&user).Error; err != nil {
		return false, err
	}

	// Check if role has permission
	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role = ? AND permissions.resource = ? AND permissions.action = ?",
			user.Role, resource, action).
		Count(&count).Error

	return count > 0, err
}

// Note: User preferences functionality has been removed due to JSONB complexity
// If needed in the future, consider using a separate notification_preferences table

// GetByID is an alias for GetUserByID for compatibility
func (r *repository) GetByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	return r.GetUserByID(ctx, userID)
}

// GetUserByIDForSecurity returns user as interface for security module compatibility
func (r *repository) GetUserByIDForSecurity(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser2FA updates the two-factor authentication status for a user
func (r *repository) UpdateUser2FA(ctx context.Context, userID uuid.UUID, enabled bool) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("two_factor_enabled", enabled).Error
}
