package entities

import (
	"time"
	"github.com/google/uuid"
)

// TenantType defines the tenant deployment type
type TenantType string

const (
	TenantTypeShared     TenantType = "shared"     // Uses shared database with tenant_id
	TenantTypeDedicated  TenantType = "dedicated"  // Uses dedicated database
)

// TenantStatus defines the tenant status
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusInactive  TenantStatus = "inactive"
	TenantStatusSuspended TenantStatus = "suspended"
)

// SubscriptionPlan defines the subscription plan types
type SubscriptionPlan string

const (
	PlanBasic        SubscriptionPlan = "basic"        // Shared database, up to 1k products
	PlanProfessional SubscriptionPlan = "professional" // Shared database, up to 10k products  
	PlanEnterprise   SubscriptionPlan = "enterprise"   // Dedicated database, unlimited
)

// Tenant represents a merchant tenant in the system
type Tenant struct {
	ID               uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name             string           `json:"name" gorm:"not null"`
	Subdomain        string           `json:"subdomain" gorm:"uniqueIndex;not null"`
	CustomDomain     *string          `json:"custom_domain,omitempty"`
	Plan             SubscriptionPlan `json:"plan" gorm:"not null;default:'basic'"`
	Status           TenantStatus     `json:"status" gorm:"not null;default:'active'"`
	Type             TenantType       `json:"type" gorm:"not null;default:'shared'"`
	DatabaseName     *string          `json:"database_name,omitempty"` // Only for dedicated tenants
	ProductCount     int              `json:"product_count" gorm:"default:0"`
	MonthlyRequests  int              `json:"monthly_requests" gorm:"default:0"`
	Settings         TenantSettings   `json:"settings" gorm:"type:jsonb"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	DeletedAt        *time.Time       `json:"deleted_at,omitempty" gorm:"index"`
}

// TenantSettings contains tenant-specific configurations
type TenantSettings struct {
	Theme          string            `json:"theme"`
	Logo           string            `json:"logo"`
	PrimaryColor   string            `json:"primary_color"`
	SecondaryColor string            `json:"secondary_color"`
	Currency       string            `json:"currency"`
	Timezone       string            `json:"timezone"`
	Language       string            `json:"language"`
	Analytics      AnalyticsSettings `json:"analytics"`
	Email          EmailSettings     `json:"email"`
	Payment        PaymentSettings   `json:"payment"`
}

// AnalyticsSettings contains analytics configuration
type AnalyticsSettings struct {
	GoogleAnalyticsID string `json:"google_analytics_id"`
	FacebookPixelID   string `json:"facebook_pixel_id"`
}

// EmailSettings contains email configuration  
type EmailSettings struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	Provider  string `json:"provider"`
}

// PaymentSettings contains payment configuration
type PaymentSettings struct {
	StripePublishableKey string `json:"stripe_publishable_key"`
	PaypalClientID       string `json:"paypal_client_id"`
}

// Business methods
func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}

// ShouldUseDedicatedDatabase determines if tenant should use dedicated database
func (t *Tenant) ShouldUseDedicatedDatabase() bool {
	return t.Plan == PlanEnterprise || t.ProductCount >= 10000
}

// CanAccessFeature checks if tenant can access a specific feature based on plan
func (t *Tenant) CanAccessFeature(feature string) bool {
	features := map[SubscriptionPlan][]string{
		PlanBasic:        {"basic_analytics", "email_support"},
		PlanProfessional: {"basic_analytics", "advanced_analytics", "email_support", "priority_support", "custom_domain"},
		PlanEnterprise:   {"basic_analytics", "advanced_analytics", "email_support", "priority_support", "custom_domain", "dedicated_database", "api_access"},
	}

	for _, f := range features[t.Plan] {
		if f == feature {
			return true
		}
	}
	return false
}

// UpdatePlan updates the tenant's subscription plan
func (t *Tenant) UpdatePlan(newPlan SubscriptionPlan) error {
	t.Plan = newPlan
	
	// Update tenant type based on plan
	if newPlan == PlanEnterprise {
		t.Type = TenantTypeDedicated
	} else {
		t.Type = TenantTypeShared
	}
	
	t.UpdatedAt = time.Now()
	return nil
}