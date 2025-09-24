package tenant

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"ecommerce-saas/internal/shared/config"
)

// Service handles tenant business logic
type Service struct {
	repo      Repository
	validator *validator.Validate
	config    *config.Config
}

// NewService creates a new tenant service
func NewService(repo Repository, config *config.Config) *Service {
	return &Service{
		repo:      repo,
		validator: validator.New(),
		config:    config,
	}
}

// CreateTenant creates a new tenant
func (s *Service) CreateTenant(req CreateTenantRequest) (*Tenant, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, NewValidationError("Invalid request data", "")
	}

	// Normalize subdomain
	req.Subdomain = strings.ToLower(strings.TrimSpace(req.Subdomain))

	// Check subdomain availability
	if exists, err := s.repo.SubdomainExists(req.Subdomain); err != nil {
		return nil, NewInternalError("Failed to check subdomain availability")
	} else if exists {
		return nil, NewConflictError("Subdomain already taken")
	}

	// Validate subdomain format (additional business rules)
	if err := s.validateSubdomain(req.Subdomain); err != nil {
		return nil, NewValidationError(err.Error(), "subdomain")
	}

	// Create tenant
	tenant := &Tenant{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(req.Name),
		Subdomain:   req.Subdomain,
		Description: strings.TrimSpace(req.Description),
		Phone:       strings.TrimSpace(req.Phone),
		Email:       strings.TrimSpace(req.Email),
		Address:     strings.TrimSpace(req.Address),
		Status:      StatusActive,
		Plan:        PlanStarter,
		Currency:    s.getDefaultCurrency(),
		Language:    "bn",
		Timezone:    "Asia/Dhaka",
	}

	return s.repo.Save(tenant)
}

// GetTenant retrieves a tenant by ID
func (s *Service) GetTenant(id string) (*Tenant, error) {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, NewValidationError("Invalid tenant ID format", "id")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return nil, NewNotFoundError("Tenant not found")
	}
	return tenant, nil
}

// GetTenantBySubdomain retrieves a tenant by subdomain
func (s *Service) GetTenantBySubdomain(subdomain string) (*Tenant, error) {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	if subdomain == "" {
		return nil, NewValidationError("Subdomain is required", "subdomain")
	}

	tenant, err := s.repo.FindBySubdomain(subdomain)
	if err != nil {
		return nil, NewNotFoundError("Tenant not found")
	}
	return tenant, nil
}

// GetTenantByCustomDomain retrieves a tenant by custom domain
func (s *Service) GetTenantByCustomDomain(domain string) (*Tenant, error) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	if domain == "" {
		return nil, NewValidationError("Domain is required", "domain")
	}

	tenant, err := s.repo.FindByCustomDomain(domain)
	if err != nil {
		return nil, NewNotFoundError("Tenant not found")
	}
	return tenant, nil
}

// UpdateTenant updates tenant information
func (s *Service) UpdateTenant(id string, req UpdateTenantRequest) (*Tenant, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, NewValidationError("Invalid request data", "")
	}

	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, NewValidationError("Invalid tenant ID format", "id")
	}

	// Get existing tenant
	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return nil, NewNotFoundError("Tenant not found")
	}

	// Update fields if provided
	if req.Name != "" {
		tenant.Name = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		tenant.Description = strings.TrimSpace(req.Description)
	}
	if req.Phone != "" {
		tenant.Phone = strings.TrimSpace(req.Phone)
	}
	if req.Email != "" {
		tenant.Email = strings.TrimSpace(req.Email)
	}
	if req.Address != "" {
		tenant.Address = strings.TrimSpace(req.Address)
	}
	if req.CustomDomain != "" {
		// TODO: Validate domain ownership before setting
		tenant.CustomDomain = strings.ToLower(strings.TrimSpace(req.CustomDomain))
	}

	return s.repo.Update(tenant)
}

// UpdatePlan updates tenant subscription plan
func (s *Service) UpdatePlan(id string, req UpdatePlanRequest) (*Tenant, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, NewValidationError("Invalid plan update request", "")
	}

	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, NewValidationError("Invalid tenant ID format", "id")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return nil, NewNotFoundError("Tenant not found")
	}

	// Validate plan change
	if !s.canChangeToPlan(tenant.Plan, req.Plan) {
		return nil, NewValidationError("Plan change not allowed", "plan")
	}

	tenant.Plan = req.Plan

	// Update limits based on new plan
	tenant.ProductLimit = s.getProductLimitForPlan(req.Plan)
	tenant.StorageLimit = tenant.GetStorageLimit()

	updatedTenant, err := s.repo.Update(tenant)
	if err != nil {
		return nil, NewInternalError("Failed to update tenant plan")
	}
	return updatedTenant, nil
}

// ListTenants returns a paginated list of tenants
func (s *Service) ListTenants(offset, limit int) ([]*Tenant, int64, error) {
	return s.repo.List(offset, limit)
}

// DeactivateTenant deactivates a tenant
func (s *Service) DeactivateTenant(id string) error {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return NewValidationError("Invalid tenant ID format", "id")
	}

	if err := s.repo.UpdateStatus(tenantID, StatusInactive); err != nil {
		return NewInternalError("Failed to deactivate tenant")
	}
	return nil
}

// ActivateTenant activates a tenant
func (s *Service) ActivateTenant(id string) error {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return NewValidationError("Invalid tenant ID format", "id")
	}

	if err := s.repo.UpdateStatus(tenantID, StatusActive); err != nil {
		return NewInternalError("Failed to activate tenant")
	}
	return nil
}

// GetTenantStats retrieves comprehensive statistics for a tenant
func (s *Service) GetTenantStats(id string) (*TenantStatsResponse, error) {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, NewValidationError("Invalid tenant ID format", "id")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return nil, NewNotFoundError("Tenant not found")
	}

	// TODO: Integrate with other modules to get actual stats
	// For now, return basic tenant information
	stats := &TenantStatsResponse{
		TenantID:       id,
		ProductCount:   0, // TODO: Get from product service
		OrderCount:     0, // TODO: Get from order service
		CustomerCount:  0, // TODO: Get from user service
		Revenue:        0, // TODO: Calculate from orders
		StorageUsed:    0, // TODO: Calculate storage usage
		BandwidthUsed:  0, // TODO: Calculate bandwidth usage
		StorageLimit:   tenant.GetStorageLimit(),
		BandwidthLimit: tenant.GetBandwidthLimit(),
		ProductLimit:   tenant.ProductLimit,
		PlanFeatures:   tenant.GetFeatureList(),
	}

	return stats, nil
}

// ListTenantsWithFilter returns a filtered and paginated list of tenants
func (s *Service) ListTenantsWithFilter(filter TenantFilter, offset, limit int) ([]*Tenant, int64, error) {
	// If no specific filters, use regular list
	if filter.Status == nil && filter.Plan == nil && filter.Search == "" {
		return s.repo.List(offset, limit)
	}

	// Apply status filter
	if filter.Status != nil {
		return s.repo.ListByStatus(*filter.Status, offset, limit)
	}

	// Apply plan filter
	if filter.Plan != nil {
		return s.repo.ListByPlan(*filter.Plan, offset, limit)
	}

	// Apply search filter
	if filter.Search != "" {
		return s.repo.Search(filter.Search, offset, limit)
	}

	return s.repo.List(offset, limit)
}

// ValidateCustomDomain validates and sets a custom domain for a tenant
func (s *Service) ValidateCustomDomain(id, domain string) error {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return NewValidationError("Invalid tenant ID format", "id")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return NewNotFoundError("Tenant not found")
	}

	// Check if tenant plan supports custom domains
	if !tenant.CanUsePremiumFeatures() {
		return NewValidationError("Custom domain requires premium or enterprise plan", "plan")
	}

	// Normalize domain
	domain = strings.ToLower(strings.TrimSpace(domain))

	// Validate domain format
	if err := s.validateDomainFormat(domain); err != nil {
		return err
	}

	// Check if domain is already taken
	if exists, err := s.repo.CustomDomainExists(domain); err != nil {
		return NewInternalError("Failed to check domain availability")
	} else if exists {
		return NewConflictError("Domain already in use")
	}

	// Perform DNS validation
	if err := s.validateDNSRecord(domain); err != nil {
		return NewValidationError(err.Error(), "domain")
	}

	// Verify HTTP accessibility
	if err := s.verifyDomainAccessibility(domain); err != nil {
		return NewValidationError(err.Error(), "domain")
	}

	tenant.CustomDomain = domain
	if _, err = s.repo.Update(tenant); err != nil {
		return NewInternalError("Failed to update tenant domain")
	}
	return nil
}

// SuspendTenant suspends a tenant (e.g., for non-payment)
func (s *Service) SuspendTenant(id string, reason string) error {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return NewValidationError("Invalid tenant ID format", "id")
	}

	// TODO: Log suspension reason and send notification
	if err := s.repo.UpdateStatus(tenantID, StatusSuspended); err != nil {
		return NewInternalError("Failed to suspend tenant")
	}
	return nil
}

// GetPlanUpgradeOptions returns available upgrade options for a tenant
func (s *Service) GetPlanUpgradeOptions(id string) ([]Plan, error) {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, NewValidationError("Invalid tenant ID format", "id")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return nil, NewNotFoundError("Tenant not found")
	}

	// Get available upgrade options
	var options []Plan
	allPlans := []Plan{PlanStarter, PlanPro, PlanPremium, PlanEnterprise}

	for _, plan := range allPlans {
		if tenant.CanUpgradeTo(plan) {
			options = append(options, plan)
		}
	}

	return options, nil
}

// InitializeTenantDefaults sets up default data for a new tenant
func (s *Service) InitializeTenantDefaults(tenantID uuid.UUID) error {
	// TODO: This could initialize:
	// - Default product categories
	// - Default pages (About, Contact, etc.)
	// - Default email templates
	// - Default payment settings
	// - Sample products (optional)

	// For now, just return success
	return nil
}

// CheckSubdomainAvailability checks if a subdomain is available
func (s *Service) CheckSubdomainAvailability(subdomain string) (bool, error) {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))

	// Validate subdomain format
	if err := s.validateSubdomain(subdomain); err != nil {
		return false, NewValidationError(err.Error(), "subdomain")
	}

	exists, err := s.repo.SubdomainExists(subdomain)
	if err != nil {
		return false, NewInternalError("Failed to check subdomain availability")
	}
	return !exists, nil
}

// GetTenantsByPlan returns tenants for a specific plan
func (s *Service) GetTenantsByPlan(plan Plan, offset, limit int) ([]*Tenant, int64, error) {
	return s.repo.ListByPlan(plan, offset, limit)
}

// GetActiveTenantsCount returns the count of active tenants
func (s *Service) GetActiveTenantsCount() (int64, error) {
	return s.repo.GetActiveCount()
}

// UpdateTenantLimits updates tenant limits based on plan
func (s *Service) UpdateTenantLimits(id string) error {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return NewValidationError("Invalid tenant ID format", "id")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return NewNotFoundError("Tenant not found")
	}

	// Update limits based on current plan
	tenant.ProductLimit = s.getProductLimitForPlan(tenant.Plan)
	tenant.StorageLimit = tenant.GetStorageLimit()
	tenant.BandwidthLimit = tenant.GetBandwidthLimit()
	tenant.UpdatedAt = time.Now()

	if _, err = s.repo.Update(tenant); err != nil {
		return NewInternalError("Failed to update tenant limits")
	}
	return nil
}

// Private helper methods

func (s *Service) validateSubdomain(subdomain string) error {
	// Reserved subdomains
	reserved := []string{"www", "api", "admin", "app", "mail", "ftp", "blog", "shop", "store"}
	for _, r := range reserved {
		if subdomain == r {
			return errors.New("subdomain is reserved")
		}
	}

	// Additional validation rules
	if len(subdomain) < 3 {
		return errors.New("subdomain must be at least 3 characters")
	}

	if strings.Contains(subdomain, "--") {
		return errors.New("subdomain cannot contain consecutive hyphens")
	}

	return nil
}

func (s *Service) canChangeToPlan(currentPlan, newPlan Plan) bool {
	// Cannot change to the same plan
	if currentPlan == newPlan {
		return false
	}

	// Define plan hierarchy for validation
	planHierarchy := map[Plan]int{
		PlanStarter:    1,
		PlanPro:        2,
		PlanPremium:    3,
		PlanEnterprise: 4,
	}

	currentLevel, currentExists := planHierarchy[currentPlan]
	newLevel, newExists := planHierarchy[newPlan]

	// Both plans must be valid
	if !currentExists || !newExists {
		return false
	}

	// Allow upgrades (moving to higher tier)
	if newLevel > currentLevel {
		return true
	}

	// Allow downgrades with restrictions
	if newLevel < currentLevel {
		// Enterprise can downgrade to any plan
		if currentPlan == PlanEnterprise {
			return true
		}
		// Premium can downgrade to Pro or Starter
		if currentPlan == PlanPremium && (newPlan == PlanPro || newPlan == PlanStarter) {
			return true
		}
		// Pro can downgrade to Starter
		if currentPlan == PlanPro && newPlan == PlanStarter {
			return true
		}
		// Starter cannot downgrade (lowest tier)
		return false
	}

	return false
}

func (s *Service) getProductLimitForPlan(plan Plan) int {
	limits := map[Plan]int{
		PlanStarter:    100,
		PlanPro:        1000,
		PlanPremium:    5000,
		PlanEnterprise: -1, // unlimited
	}

	if limit, exists := limits[plan]; exists {
		return limit
	}
	return 100 // default
}

// DNS validation helper methods

func (s *Service) validateDomainFormat(domain string) error {
	// Basic domain format validation
	if len(domain) == 0 {
		return errors.New("domain cannot be empty")
	}

	if len(domain) > 253 {
		return errors.New("domain name too long")
	}

	// Check for valid characters and format
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return errors.New("domain cannot start or end with a dot")
	}

	if strings.Contains(domain, "..") {
		return errors.New("domain cannot contain consecutive dots")
	}

	// Prevent using platform domain as custom domain
	if strings.HasSuffix(domain, "."+s.getBaseDomain()) {
		return errors.New("cannot use platform subdomain as custom domain")
	}

	return nil
}

func (s *Service) validateDNSRecord(domain string) error {
	baseDomain := s.getBaseDomain()

	// Lookup CNAME record
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return errors.New("DNS lookup failed: ensure CNAME record points to " + baseDomain)
	}

	// Verify CNAME points to our domain
	if !strings.HasSuffix(strings.TrimSuffix(cname, "."), baseDomain) {
		return errors.New("CNAME record must point to " + baseDomain)
	}

	return nil
}

func (s *Service) verifyDomainAccessibility(domain string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test HTTP accessibility
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://"+domain, nil)
	if err != nil {
		return errors.New("failed to create HTTP request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.New("domain is not accessible via HTTP: " + err.Error())
	}
	defer resp.Body.Close()

	// Accept any HTTP status as long as domain is reachable
	return nil
}

func (s *Service) getBaseDomain() string {
	// Get base domain from config
	if s.config != nil && s.config.App.Domain != "" {
		return s.config.App.Domain
	}
	// Fallback to environment variable
	baseDomain := os.Getenv("DOMAIN")
	if baseDomain == "" {
		// No hardcoded fallback - require proper configuration
		panic("Base domain must be configured via config.App.Domain or DOMAIN environment variable")
	}
	return baseDomain
}

func (s *Service) getDefaultCurrency() string {
	// Get default currency from config
	if s.config != nil && s.config.App.Currency != "" {
		return s.config.App.Currency
	}
	// Fallback to environment variable
	currency := os.Getenv("DEFAULT_CURRENCY")
	if currency == "" {
		// Use default USD if not configured
		return "USD"
	}
	return currency
}

// Trial Management Methods

// StartTrial initiates a trial period for a tenant
func (s *Service) StartTrial(tenantID uuid.UUID, trialDurationDays int) error {
	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return err
	}

	if tenant.IsTrialActive {
		return NewConflictError("Tenant already has an active trial")
	}

	tenant.StartTrial(trialDurationDays)
	tenant.UpdatedAt = time.Now()

	_, err = s.repo.Update(tenant)
	if err != nil {
		return NewInternalError("Failed to start trial")
	}

	return nil
}

// EndTrial ends the trial period for a tenant
func (s *Service) EndTrial(tenantID uuid.UUID) error {
	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return err
	}

	tenant.EndTrial()
	tenant.UpdatedAt = time.Now()

	_, err = s.repo.Update(tenant)
	if err != nil {
		return NewInternalError("Failed to end trial")
	}

	return nil
}

// GetExpiringTrials returns tenants whose trials are expiring within the specified days
func (s *Service) GetExpiringTrials(days int) ([]*Tenant, error) {
	if days < 0 {
		return nil, NewValidationError("Days must be non-negative", "days")
	}

	tenants, err := s.repo.GetExpiringTrials(days)
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

// GetTrialStatus returns the trial status for a tenant
func (s *Service) GetTrialStatus(tenantID uuid.UUID) (*TrialStatusResponse, error) {
	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return nil, err
	}

	status := &TrialStatusResponse{
		TenantID:        tenant.ID.String(),
		IsTrialActive:   tenant.IsTrialActive,
		IsTrialExpired:  tenant.IsTrialExpired(),
		IsInTrialPeriod: tenant.IsInTrialPeriod(),
		DaysRemaining:   tenant.GetTrialDaysRemaining(),
	}

	if tenant.TrialStartDate != nil {
		status.TrialStartDate = tenant.TrialStartDate
	}
	if tenant.TrialEndDate != nil {
		status.TrialEndDate = tenant.TrialEndDate
	}

	return status, nil
}

// ProcessExpiredTrials processes all expired trials and updates tenant status
func (s *Service) ProcessExpiredTrials() error {
	// Get all tenants with expired trials
	expiredTrials, err := s.repo.GetExpiringTrials(0) // 0 days means already expired
	if err != nil {
		return NewInternalError("Failed to get expired trials")
	}

	for _, tenant := range expiredTrials {
		if tenant.IsTrialExpired() {
			// End the trial
			tenant.EndTrial()
			
			// Optionally change status to suspended or require payment
			// This depends on business logic
			if tenant.Plan == PlanStarter {
				tenant.Status = StatusSuspended
			}
			
			tenant.UpdatedAt = time.Now()
			
			if _, updateErr := s.repo.Update(tenant); updateErr != nil {
				// Log error but continue processing other tenants
				continue
			}
		}
	}

	return nil
}

// TODO: Add integration methods when other modules are ready
// - GetProductUsage(tenantID uuid.UUID) (int, error)
// - GetOrderMetrics(tenantID uuid.UUID) (*OrderMetrics, error)
// - GetStorageUsage(tenantID uuid.UUID) (int64, error)
// - GetBandwidthUsage(tenantID uuid.UUID) (int64, error)
// - SendWelcomeEmail(tenant *Tenant) error
// - SendPlanUpgradeNotification(tenant *Tenant) error
// - CalculateMonthlyBill(tenantID uuid.UUID) (*BillingInfo, error)
// - SendTrialExpirationNotification(tenant *Tenant) error
