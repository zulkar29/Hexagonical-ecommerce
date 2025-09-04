package tenant

import (
	"errors"
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

// Request/Response structs

// CreateTenantRequest represents the request to create a new tenant
type CreateTenantRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Subdomain   string `json:"subdomain" validate:"required,min=3,max=50,alphanum"`
	Description string `json:"description,omitempty" validate:"max=500"`
	Phone       string `json:"phone,omitempty" validate:"max=20"`
	Email       string `json:"email,omitempty" validate:"email"`
	Address     string `json:"address,omitempty" validate:"max=255"`
}

// UpdateTenantRequest represents the request to update tenant information
type UpdateTenantRequest struct {
	Name         string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description  string `json:"description,omitempty" validate:"max=500"`
	Phone        string `json:"phone,omitempty" validate:"max=20"`
	Email        string `json:"email,omitempty" validate:"email"`
	Address      string `json:"address,omitempty" validate:"max=255"`
	CustomDomain string `json:"custom_domain,omitempty" validate:"omitempty,fqdn"`
	Logo         string `json:"logo,omitempty"`
	Currency     string `json:"currency,omitempty" validate:"omitempty,len=3"`
	Language     string `json:"language,omitempty" validate:"omitempty,len=2"`
	Timezone     string `json:"timezone,omitempty"`
}

// UpdatePlanRequest represents the request to update subscription plan
type UpdatePlanRequest struct {
	Plan Plan `json:"plan" validate:"required"`
}

// TenantStatsResponse represents tenant statistics
type TenantStatsResponse struct {
	TenantID        string  `json:"tenant_id"`
	ProductCount    int64   `json:"product_count"`
	OrderCount      int64   `json:"order_count"`
	CustomerCount   int64   `json:"customer_count"`
	Revenue         float64 `json:"revenue"`
	StorageUsed     int64   `json:"storage_used_mb"`
	BandwidthUsed   int64   `json:"bandwidth_used_mb"`
	StorageLimit    int     `json:"storage_limit_mb"`
	BandwidthLimit  int     `json:"bandwidth_limit_mb"`
	ProductLimit    int     `json:"product_limit"`
	PlanFeatures    []string `json:"plan_features"`
}

// TenantFilter represents filtering options for tenant listing
type TenantFilter struct {
	Status   *Status `json:"status,omitempty"`
	Plan     *Plan   `json:"plan,omitempty"`
	Search   string  `json:"search,omitempty"`
	DateFrom string  `json:"date_from,omitempty"`
	DateTo   string  `json:"date_to,omitempty"`
}

// Additional business logic methods

// ValidateBusinessInfo validates that all required business information is complete
func (t *Tenant) ValidateBusinessInfo() error {
	if t.Name == "" {
		return errors.New("business name is required")
	}
	if t.Email == "" {
		return errors.New("business email is required")
	}
	if t.Phone == "" {
		return errors.New("business phone is required")
	}
	if t.Address == "" {
		return errors.New("business address is required")
	}
	return nil
}

// CanUpgradeTo checks if tenant can upgrade to a new plan
func (t *Tenant) CanUpgradeTo(newPlan Plan) bool {
	// Define upgrade paths
	upgradePaths := map[Plan][]Plan{
		PlanStarter:    {PlanPro, PlanPremium, PlanEnterprise},
		PlanPro:        {PlanPremium, PlanEnterprise},
		PlanPremium:    {PlanEnterprise},
		PlanEnterprise: {}, // Cannot upgrade from enterprise
	}
	
	allowedUpgrades, exists := upgradePaths[t.Plan]
	if !exists {
		return false
	}
	
	for _, allowed := range allowedUpgrades {
		if allowed == newPlan {
			return true
		}
	}
	return false
}

// GetFeatureList returns available features for the current plan
func (t *Tenant) GetFeatureList() []string {
	features := map[Plan][]string{
		PlanStarter: {
			"Basic storefront",
			"Up to 100 products",
			"Standard templates",
			"Basic analytics",
			"Email support",
			"SSL certificate",
			"Payment gateway integration",
		},
		PlanPro: {
			"All Starter features",
			"Up to 1,000 products",
			"Advanced analytics",
			"Custom domain",
			"Priority support",
			"Inventory management",
			"Discount codes",
			"Abandoned cart recovery",
		},
		PlanPremium: {
			"All Pro features",
			"Up to 5,000 products",
			"Multi-language support",
			"Advanced SEO tools",
			"API access",
			"Custom integrations",
			"Advanced reporting",
			"24/7 phone support",
		},
		PlanEnterprise: {
			"All Premium features",
			"Unlimited products",
			"White-label solution",
			"Dedicated account manager",
			"Custom development",
			"SLA guarantee",
			"Advanced security",
			"Multi-store management",
		},
	}
	
	if planFeatures, exists := features[t.Plan]; exists {
		return planFeatures
	}
	return []string{}
}

// IsTrialExpired checks if trial period has expired (if applicable)
func (t *Tenant) IsTrialExpired() bool {
	// TODO: Implement trial logic when trial system is added
	// For now, assume no trials
	return false
}

// GetBandwidthLimit returns bandwidth limit in MB based on plan
func (t *Tenant) GetBandwidthLimit() int {
	limits := map[Plan]int{
		PlanStarter:    10240,  // 10GB
		PlanPro:        51200,  // 50GB
		PlanPremium:    102400, // 100GB
		PlanEnterprise: 512000, // 500GB
	}
	
	if limit, exists := limits[t.Plan]; exists {
		return limit
	}
	return 10240 // default
}

// GetDisplayName returns a formatted display name for the tenant
func (t *Tenant) GetDisplayName() string {
	if t.Name != "" {
		return t.Name
	}
	return t.Subdomain
}

// IsBusinessInfoComplete checks if all business information is provided
func (t *Tenant) IsBusinessInfoComplete() bool {
	return t.Name != "" && t.Email != "" && t.Phone != "" && t.Address != ""
}

// CanAccessFeature checks if tenant plan allows access to a specific feature
func (t *Tenant) CanAccessFeature(feature string) bool {
	premiumFeatures := []string{"custom_domain", "api_access", "advanced_analytics", "multi_language"}
	enterpriseFeatures := []string{"white_label", "custom_development", "dedicated_support"}
	
	// Check if it's a premium feature
	for _, f := range premiumFeatures {
		if f == feature {
			return t.Plan == PlanPremium || t.Plan == PlanEnterprise
		}
	}
	
	// Check if it's an enterprise feature
	for _, f := range enterpriseFeatures {
		if f == feature {
			return t.Plan == PlanEnterprise
		}
	}
	
	// Basic features available to all plans
	return true
}
