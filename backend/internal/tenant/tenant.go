package tenant

import (
	"time"

	"github.com/google/uuid"
)

// Status represents the tenant status
type Status string

// Plan represents the subscription plan
type Plan string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusSuspended Status = "suspended"
)

const (
	PlanStarter    Plan = "starter"     // ৳1,990
	PlanPro        Plan = "professional" // ৳4,990
	PlanPremium    Plan = "premium"     // ৳7,990
	PlanEnterprise Plan = "enterprise"  // ৳12,990
)

// Tenant represents a store/tenant in the system
type Tenant struct {
	ID          uuid.UUID `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"not null"`
	Subdomain   string    `json:"subdomain" gorm:"unique;not null"`
	CustomDomain string    `json:"custom_domain,omitempty"`
	Status      Status    `json:"status" gorm:"default:active"`
	Plan        Plan      `json:"plan" gorm:"default:starter"`
	
	// Business Information
	Description  string `json:"description,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	Address      string `json:"address,omitempty"`
	Logo         string `json:"logo,omitempty"`
	
	// Settings
	Currency     string `json:"currency" gorm:"default:BDT"`
	Language     string `json:"language" gorm:"default:bn"`
	Timezone     string `json:"timezone" gorm:"default:Asia/Dhaka"`
	
	// Limits based on plan
	ProductLimit    int `json:"product_limit" gorm:"default:100"`
	StorageLimit    int `json:"storage_limit" gorm:"default:1024"` // MB
	BandwidthLimit  int `json:"bandwidth_limit" gorm:"default:10240"` // MB
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Business Logic Methods

// IsActive checks if the tenant is active
func (t *Tenant) IsActive() bool {
	return t.Status == StatusActive
}

// CanCreateProducts checks if tenant can create more products based on plan limits
func (t *Tenant) CanCreateProducts(currentCount int) bool {
	limits := map[Plan]int{
		PlanStarter:    100,
		PlanPro:        1000,
		PlanPremium:    5000,
		PlanEnterprise: -1, // unlimited
	}
	
	limit, exists := limits[t.Plan]
	if !exists {
		return false
	}
	
	return limit == -1 || currentCount < limit
}

// GetStorageLimit returns storage limit in MB based on plan
func (t *Tenant) GetStorageLimit() int {
	limits := map[Plan]int{
		PlanStarter:    1024,   // 1GB
		PlanPro:        5120,   // 5GB
		PlanPremium:    10240,  // 10GB
		PlanEnterprise: 51200,  // 50GB
	}
	
	if limit, exists := limits[t.Plan]; exists {
		return limit
	}
	return 1024 // default
}

// GetMonthlyPrice returns the monthly price for the plan in BDT
func (t *Tenant) GetMonthlyPrice() int {
	prices := map[Plan]int{
		PlanStarter:    1990,
		PlanPro:        4990,
		PlanPremium:    7990,
		PlanEnterprise: 12990,
	}
	
	if price, exists := prices[t.Plan]; exists {
		return price
	}
	return 0
}

// HasCustomDomain checks if tenant has a custom domain configured
func (t *Tenant) HasCustomDomain() bool {
	return t.CustomDomain != ""
}

// GetDomain returns the primary domain (custom if available, otherwise subdomain)
func (t *Tenant) GetDomain() string {
	if t.HasCustomDomain() {
		return t.CustomDomain
	}
	return t.Subdomain + ".esass.com" // TODO: Use actual domain from config
}

// CanUsePremiumFeatures checks if tenant plan allows premium features
func (t *Tenant) CanUsePremiumFeatures() bool {
	return t.Plan == PlanPremium || t.Plan == PlanEnterprise
}

// CanUseAdvancedAnalytics checks if tenant can access advanced analytics
func (t *Tenant) CanUseAdvancedAnalytics() bool {
	return t.Plan == PlanPro || t.Plan == PlanPremium || t.Plan == PlanEnterprise
}

// TODO: Add more business logic methods as needed
// - ValidateBusinessInfo()
// - CanUpgradeTo(newPlan Plan) bool
// - GetFeatureList() []string
// - IsTrialExpired() bool
