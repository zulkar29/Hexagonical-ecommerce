package user

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository is an alias for RepositoryInterface for backward compatibility
type Repository = RepositoryInterface

// RepositoryInterface defines user repository interface
type RepositoryInterface interface {
	// User CRUD operations
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	ListUsers(ctx context.Context, tenantID *uuid.UUID, filter UserFilter, offset, limit int) ([]*User, int64, error)
	UpdateUserByAdmin(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*User, error)
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]User, error)
	GetUsersWithFilters(ctx context.Context, filters interface{}) ([]User, error)
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) error
	UpdatePhoneVerification(ctx context.Context, phone string, verified bool) error

	// Session operations
	CreateSession(session *UserSession) error
	GetSessionByToken(tenantID *uuid.UUID, token string) (*UserSession, error)
	UpdateSession(session *UserSession) error
	InvalidateUserSessions(ctx context.Context, tenantID *uuid.UUID, userID uuid.UUID) error
	CleanupExpiredSessions(tenantID *uuid.UUID) error

	// Permission operations
	GetUserPermissions(userID uuid.UUID) ([]*Permission, error)
	CheckUserPermission(userID uuid.UUID, resource, action string) (bool, error)

	// Preferences operations
	GetUserPreferences(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	UpdateUserPreferences(ctx context.Context, userID uuid.UUID, preferences map[string]interface{}) (map[string]interface{}, error)
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
		return nil, err
	}
	return user, nil
}

func (r *repository) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
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

// Session operations

func (r *repository) CreateSession(session *UserSession) error {
	return r.db.Create(session).Error
}

func (r *repository) GetSessionByToken(tenantID *uuid.UUID, token string) (*UserSession, error) {
	var session UserSession
	query := r.db.Where("token = ? AND is_active = true", token)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repository) UpdateSession(session *UserSession) error {
	return r.db.Save(session).Error
}

func (r *repository) InvalidateUserSessions(ctx context.Context, tenantID *uuid.UUID, userID uuid.UUID) error {
	query := r.db.WithContext(ctx).Model(&UserSession{}).Where("user_id = ?", userID)
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	return query.Update("is_active", false).Error
}

func (r *repository) CleanupExpiredSessions(tenantID *uuid.UUID) error {
	query := r.db.Where("expires_at < NOW() OR is_active = false")
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	return query.Delete(&UserSession{}).Error
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

// Preferences operations

func (r *repository) GetUserPreferences(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	var user User
	if err := r.db.WithContext(ctx).Select("preferences").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	
	if user.Preferences == nil {
		return make(map[string]interface{}), nil
	}
	
	return user.Preferences, nil
}

func (r *repository) UpdateUserPreferences(ctx context.Context, userID uuid.UUID, preferences map[string]interface{}) (map[string]interface{}, error) {
	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("preferences", preferences).Error; err != nil {
		return nil, err
	}
	
	return preferences, nil
}
