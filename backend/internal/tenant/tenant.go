package tenant

import (
	"errors"
	"os"
	"time"

	"ecommerce-saas/internal/shared/constants"
	"github.com/google/uuid"
)

// Plan represents the subscription plan
type Plan string

// Use shared status constants
type Status = constants.TenantStatus
const (
	StatusActive    = constants.TenantStatusActive
	StatusInactive  = constants.TenantStatusInactive
	StatusSuspended = constants.TenantStatusSuspended
)

const (
	PlanFree         Plan = "free"         // ৳0
	PlanStarter      Plan = "starter"      // ৳1,990
	PlanProfessional Plan = "professional" // ৳4,990
	PlanPro          Plan = "pro"          // ৳7,990
	PlanEnterprise   Plan = "enterprise"   // ৳12,990
)

// Tenant represents a store/tenant in the system
type Tenant struct {
	ID           uuid.UUID `json:"id" gorm:"primarykey"`
	Name         string    `json:"name" gorm:"not null"`
	Subdomain    string    `json:"subdomain" gorm:"unique;not null"`
	CustomDomain string    `json:"custom_domain,omitempty"`
	Status       Status    `json:"status" gorm:"default:active"`
	Plan         Plan      `json:"plan" gorm:"default:starter"`

	// Trial Information
	TrialStartDate *time.Time `json:"trial_start_date,omitempty"`
	TrialEndDate   *time.Time `json:"trial_end_date,omitempty"`
	IsTrialActive  bool       `json:"is_trial_active" gorm:"default:false"`

	// Business Information
	Description string `json:"description,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	Address     string `json:"address,omitempty"`
	Logo        string `json:"logo,omitempty"`

	// Settings
	Currency string `json:"currency" gorm:"default:BDT"`
	Language string `json:"language" gorm:"default:bn"`
	Timezone string `json:"timezone" gorm:"default:Asia/Dhaka"`

	// Limits based on plan
	ProductLimit   int `json:"product_limit" gorm:"default:100"`
	StorageLimit   int `json:"storage_limit" gorm:"default:1024"`    // MB
	BandwidthLimit int `json:"bandwidth_limit" gorm:"default:10240"` // MB

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
		PlanFree:         10,
		PlanStarter:      500,
		PlanProfessional: 2000,
		PlanPro:          10000,
		PlanEnterprise:   -1, // unlimited
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
		PlanFree:         1024,   // 1GB
		PlanStarter:      5120,   // 5GB
		PlanProfessional: 20480,  // 20GB
		PlanPro:          102400, // 100GB
		PlanEnterprise:   -1,     // unlimited
	}

	if limit, exists := limits[t.Plan]; exists {
		return limit
	}
	return 1024 // default
}

// GetMonthlyPrice returns the monthly price for the plan in BDT
func (t *Tenant) GetMonthlyPrice() int {
	prices := map[Plan]int{
		PlanFree:         0,
		PlanStarter:      1990,
		PlanProfessional: 4990,
		PlanPro:          7990,
		PlanEnterprise:   12990,
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
// Deprecated: Use GetDomainWithBase(baseDomain string) for configurable domains
// This method is kept for backward compatibility but should be avoided
func (t *Tenant) GetDomain() string {
	if t.HasCustomDomain() {
		return t.CustomDomain
	}
	// Use environment variable - no hardcoded fallback
	baseDomain := os.Getenv("ESASS_DOMAIN")
	if baseDomain == "" {
		baseDomain = os.Getenv("DOMAIN")
	}
	if baseDomain == "" {
		panic("Base domain must be configured via ESASS_DOMAIN or DOMAIN environment variable")
	}
	return t.Subdomain + "." + baseDomain
}

// GetDomainWithBase returns the primary domain with configurable base domain
func (t *Tenant) GetDomainWithBase(baseDomain string) string {
	if t.HasCustomDomain() {
		return t.CustomDomain
	}
	return t.Subdomain + "." + baseDomain
}

// CanUsePremiumFeatures checks if tenant plan allows premium features
func (t *Tenant) CanUsePremiumFeatures() bool {
	return t.Plan == PlanPro || t.Plan == PlanEnterprise
}

// CanUseAdvancedAnalytics checks if tenant can access advanced analytics
func (t *Tenant) CanUseAdvancedAnalytics() bool {
	return t.Plan == PlanProfessional || t.Plan == PlanPro || t.Plan == PlanEnterprise
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
	TenantID       string   `json:"tenant_id"`
	ProductCount   int64    `json:"product_count"`
	OrderCount     int64    `json:"order_count"`
	CustomerCount  int64    `json:"customer_count"`
	Revenue        float64  `json:"revenue"`
	StorageUsed    int64    `json:"storage_used_mb"`
	BandwidthUsed  int64    `json:"bandwidth_used_mb"`
	StorageLimit   int      `json:"storage_limit_mb"`
	BandwidthLimit int      `json:"bandwidth_limit_mb"`
	ProductLimit   int      `json:"product_limit"`
	PlanFeatures   []string `json:"plan_features"`
}

// TrialStatusResponse represents trial status information
type TrialStatusResponse struct {
	TenantID        string     `json:"tenant_id"`
	IsTrialActive   bool       `json:"is_trial_active"`
	IsTrialExpired  bool       `json:"is_trial_expired"`
	IsInTrialPeriod bool       `json:"is_in_trial_period"`
	DaysRemaining   int        `json:"days_remaining"`
	TrialStartDate  *time.Time `json:"trial_start_date,omitempty"`
	TrialEndDate    *time.Time `json:"trial_end_date,omitempty"`
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
		PlanFree:         {PlanStarter, PlanProfessional, PlanPro, PlanEnterprise},
		PlanStarter:      {PlanProfessional, PlanPro, PlanEnterprise},
		PlanProfessional: {PlanPro, PlanEnterprise},
		PlanPro:          {PlanEnterprise},
		PlanEnterprise:   {}, // Cannot upgrade from enterprise
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
		PlanFree: {
			"Up to 10 products",
			"Basic storefront",
			"Community support",
			"Platform branding",
			"1GB storage",
			"Limited features",
		},
		PlanStarter: {
			"Up to 500 products",
			"1 staff account",
			"Basic themes",
			"Email support",
			"5GB storage",
			"Standard features",
		},
		PlanProfessional: {
			"All Starter features",
			"Up to 2,000 products",
			"3 staff accounts",
			"Premium themes",
			"Priority support",
			"20GB storage",
			"Advanced analytics",
			"Email marketing",
		},
		PlanPro: {
			"All Professional features",
			"Up to 10,000 products",
			"10 staff accounts",
			"Custom themes",
			"Advanced analytics",
			"100GB storage",
			"Abandoned cart recovery",
			"Advanced email marketing",
			"Basic API access",
		},
		PlanEnterprise: {
			"All Pro features",
			"Unlimited products",
			"Unlimited staff accounts",
			"White-label options",
			"24/7 dedicated support",
			"Unlimited storage",
			"Full API access",
			"Advanced integrations",
			"Priority database performance",
			"Custom development support",
			"Priority feature requests",
		},
	}

	if planFeatures, exists := features[t.Plan]; exists {
		return planFeatures
	}
	return []string{}
}

// IsTrialExpired checks if trial period has expired (if applicable)
func (t *Tenant) IsTrialExpired() bool {
	if !t.IsTrialActive || t.TrialEndDate == nil {
		return false
	}
	return time.Now().After(*t.TrialEndDate)
}

// StartTrial initiates a trial period for the tenant
func (t *Tenant) StartTrial(trialDurationDays int) {
	now := time.Now()
	t.TrialStartDate = &now
	trialEnd := now.AddDate(0, 0, trialDurationDays)
	t.TrialEndDate = &trialEnd
	t.IsTrialActive = true
}

// EndTrial ends the trial period for the tenant
func (t *Tenant) EndTrial() {
	t.IsTrialActive = false
}

// GetTrialDaysRemaining returns the number of days remaining in the trial
func (t *Tenant) GetTrialDaysRemaining() int {
	if !t.IsTrialActive || t.TrialEndDate == nil {
		return 0
	}
	if t.IsTrialExpired() {
		return 0
	}
	duration := time.Until(*t.TrialEndDate)
	return int(duration.Hours() / 24)
}

// IsInTrialPeriod checks if tenant is currently in trial period
func (t *Tenant) IsInTrialPeriod() bool {
	return t.IsTrialActive && !t.IsTrialExpired()
}

// GetBandwidthLimit returns bandwidth limit in MB based on plan
func (t *Tenant) GetBandwidthLimit() int {
	limits := map[Plan]int{
		PlanFree:         10240,   // 10GB/month
		PlanStarter:      102400,  // 100GB/month
		PlanProfessional: 512000,  // 500GB/month
		PlanPro:          2097152, // 2TB/month
		PlanEnterprise:   -1,      // unlimited
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
	professionalFeatures := []string{"advanced_analytics", "email_marketing", "priority_support"}
	proFeatures := []string{"custom_domain", "api_access", "abandoned_cart_recovery", "advanced_email_marketing"}
	enterpriseFeatures := []string{"white_label", "custom_development", "dedicated_support", "unlimited_storage"}

	// Check if it's a professional feature
	for _, f := range professionalFeatures {
		if f == feature {
			return t.Plan == PlanProfessional || t.Plan == PlanPro || t.Plan == PlanEnterprise
		}
	}

	// Check if it's a pro feature
	for _, f := range proFeatures {
		if f == feature {
			return t.Plan == PlanPro || t.Plan == PlanEnterprise
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
