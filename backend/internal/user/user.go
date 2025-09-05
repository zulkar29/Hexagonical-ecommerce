package user

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents user roles in the system
type UserRole string

// UserStatus represents user account status
type UserStatus string

const (
	RoleCustomer UserRole = "customer"
	RoleMerchant UserRole = "merchant"
	RoleAdmin    UserRole = "admin"
	RoleSuper    UserRole = "super_admin"
)

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
	StatusPending   UserStatus = "pending"
)

// User represents a user in the system
type User struct {
	ID       uuid.UUID  `json:"id" gorm:"primarykey"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty" gorm:"index"` // Null for super admins
	
	// Basic information
	Email     string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password  string `json:"-" gorm:"not null" validate:"required,min=8"` // Hidden in JSON
	FirstName string `json:"first_name" gorm:"not null" validate:"required"`
	LastName  string `json:"last_name" gorm:"not null" validate:"required"`
	Phone     string `json:"phone,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	
	// Account details
	Role   UserRole   `json:"role" gorm:"default:customer"`
	Status UserStatus `json:"status" gorm:"default:active"`
	
	// Verification
	EmailVerified     bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty"`
	PhoneVerified     bool       `json:"phone_verified" gorm:"default:false"`
	PhoneVerifiedAt   *time.Time `json:"phone_verified_at,omitempty"`
	
	// Security
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty"`
	TwoFactorEnabled  bool       `json:"two_factor_enabled" gorm:"default:false"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Sessions []UserSession `json:"sessions,omitempty" gorm:"foreignKey:UserID"`
}

// UserSession represents an active user session
type UserSession struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Permission represents a system permission
type Permission struct {
	ID          uuid.UUID `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description,omitempty"`
	Resource    string    `json:"resource" gorm:"not null"` // e.g., "products", "orders"
	Action      string    `json:"action" gorm:"not null"`   // e.g., "create", "read", "update", "delete"
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RolePermission links roles to permissions
type RolePermission struct {
	ID           uuid.UUID `json:"id" gorm:"primarykey"`
	Role         UserRole  `json:"role" gorm:"not null"`
	PermissionID uuid.UUID `json:"permission_id" gorm:"not null"`
	
	CreatedAt time.Time `json:"created_at"`
}

// Business Logic Methods for User

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// IsAdmin checks if the user has admin privileges
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuper
}

// IsMerchant checks if the user is a merchant
func (u *User) IsMerchant() bool {
	return u.Role == RoleMerchant
}

// IsCustomer checks if the user is a customer
func (u *User) IsCustomer() bool {
	return u.Role == RoleCustomer
}

// IsSuperAdmin checks if the user is a super admin
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuper
}

// CanAccessTenant checks if user can access a specific tenant
func (u *User) CanAccessTenant(tenantID uuid.UUID) bool {
	// Super admins can access all tenants
	if u.IsSuperAdmin() {
		return true
	}
	
	// Users can only access their own tenant
	return u.TenantID != nil && *u.TenantID == tenantID
}

// NeedsEmailVerification checks if email verification is required
func (u *User) NeedsEmailVerification() bool {
	return !u.EmailVerified
}

// NeedsPhoneVerification checks if phone verification is required
func (u *User) NeedsPhoneVerification() bool {
	return u.Phone != "" && !u.PhoneVerified
}

// HasRecentLogin checks if user has logged in recently (within days)
func (u *User) HasRecentLogin(days int) bool {
	if u.LastLoginAt == nil {
		return false
	}
	threshold := time.Now().AddDate(0, 0, -days)
	return u.LastLoginAt.After(threshold)
}

// ShouldChangePassword checks if password change is required
func (u *User) ShouldChangePassword(maxDays int) bool {
	if u.PasswordChangedAt == nil {
		// If never changed, use creation date
		threshold := time.Now().AddDate(0, 0, -maxDays)
		return u.CreatedAt.Before(threshold)
	}
	
	threshold := time.Now().AddDate(0, 0, -maxDays)
	return u.PasswordChangedAt.Before(threshold)
}

// Business Logic Methods for UserSession

// IsExpired checks if the session has expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if the session is valid (active and not expired)
func (s *UserSession) IsValid() bool {
	return s.IsActive && !s.IsExpired()
}

// Business Logic Methods for Permission

// GetPermissionKey returns a unique key for the permission
func (p *Permission) GetPermissionKey() string {
	return p.Resource + ":" + p.Action
}

// External Integration Structures

// LoginResponse represents login response with tokens
type LoginResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// TokenResponse represents token refresh response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}



// UserFilter represents user listing filters
type UserFilter struct {
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
	Role     UserRole   `json:"role,omitempty"`
	Status   UserStatus `json:"status,omitempty"`
	Search   string     `json:"search,omitempty"` // Search in name or email
}
