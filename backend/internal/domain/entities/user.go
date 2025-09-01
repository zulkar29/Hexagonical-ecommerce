package entities

// TODO: Implement User entity during development
// This entity represents users within a tenant context

import (
	"time"
	// TODO: Add imports during implementation
	// "github.com/google/uuid"
	// "ecommerce-saas/internal/domain/valueobjects"
)

// User represents a user in the system (merchant staff or customer)
type User struct {
	// TODO: Implement actual fields
	// ID        uuid.UUID               `json:"id"`
	// TenantID  uuid.UUID               `json:"tenant_id"`
	// Email     valueobjects.Email      `json:"email"`
	// Password  valueobjects.Password   `json:"-"` // Never serialize password
	// FirstName string                  `json:"first_name"`
	// LastName  string                  `json:"last_name"`
	// Role      valueobjects.UserRole   `json:"role"`
	// Status    valueobjects.UserStatus `json:"status"`
	// Profile   UserProfile             `json:"profile"`
	// CreatedAt time.Time               `json:"created_at"`
	// UpdatedAt time.Time               `json:"updated_at"`
	// DeletedAt *time.Time              `json:"deleted_at,omitempty"`
}

// UserProfile contains additional user information
type UserProfile struct {
	// TODO: Implement actual profile fields
	// Avatar      string    `json:"avatar"`
	// Phone       string    `json:"phone"`
	// Timezone    string    `json:"timezone"`
	// Language    string    `json:"language"`
	// LastLoginAt time.Time `json:"last_login_at"`
}

// Customer represents a customer user (extends User for storefront)
type Customer struct {
	User
	// TODO: Add customer-specific fields
	// Addresses    []Address             `json:"addresses"`
	// Preferences  CustomerPreferences   `json:"preferences"`
	// LoyaltyPoints int                  `json:"loyalty_points"`
}

// CustomerPreferences contains customer shopping preferences
type CustomerPreferences struct {
	// TODO: Implement customer preferences
	// Newsletter     bool   `json:"newsletter"`
	// SMSNotifications bool `json:"sms_notifications"`
	// Currency       string `json:"preferred_currency"`
	// Language       string `json:"preferred_language"`
}

// Business methods (to be implemented)
// func (u *User) ChangePassword(newPassword string) error
// func (u *User) HasPermission(permission string) bool
// func (u *User) IsActive() bool
// func (u *User) BelongsToTenant(tenantID uuid.UUID) bool