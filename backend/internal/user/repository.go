package user

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines user repository interface
type Repository interface {
	// User CRUD operations
	Create(user *User) error
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	List(tenantID *uuid.UUID, filter UserFilter, offset, limit int) ([]*User, int64, error)
	UpdateUserByAdmin(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*User, error)
	GetUsersForExport(ctx context.Context, filters map[string]string) ([]*User, error)

	// Session operations
	CreateSession(session *UserSession) error
	GetSessionByToken(token string) (*UserSession, error)
	UpdateSession(session *UserSession) error
	InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error
	CleanupExpiredSessions() error

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

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) GetByID(id uuid.UUID) (*User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *repository) List(tenantID *uuid.UUID, filter UserFilter, offset, limit int) ([]*User, int64, error) {
	var users []*User
	var total int64

	query := r.db.Model(&User{})

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
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// Session operations

func (r *repository) CreateSession(session *UserSession) error {
	return r.db.Create(session).Error
}

func (r *repository) GetSessionByToken(token string) (*UserSession, error) {
	var session UserSession
	err := r.db.Where("token = ? AND is_active = true", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repository) UpdateSession(session *UserSession) error {
	return r.db.Save(session).Error
}

// Removed duplicate method - using context-aware version below

func (r *repository) CleanupExpiredSessions() error {
	return r.db.Where("expires_at < NOW() OR is_active = false").
		Delete(&UserSession{}).Error
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

// Additional context-aware methods

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&User{}, userID).Error
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

func (r *repository) GetUsersForExport(ctx context.Context, filters map[string]string) ([]*User, error) {
	var users []*User
	query := r.db.WithContext(ctx).Model(&User{})
	
	if role := filters["role"]; role != "" {
		query = query.Where("role = ?", role)
	}
	if status := filters["status"]; status != "" {
		query = query.Where("status = ?", status)
	}
	if search := filters["search"]; search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}
	
	err := query.Find(&users).Error
	return users, err
}

// Context-aware session invalidation
func (r *repository) InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&UserSession{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error
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
