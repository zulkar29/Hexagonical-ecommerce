package tenant

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

// Service handles tenant business logic
type Service struct {
	repo      *Repository
	validator *validator.Validate
}

// NewService creates a new tenant service
func NewService(repo *Repository) *Service {
	return &Service{
		repo:      repo,
		validator: validator.New(),
	}
}

// CreateTenant creates a new tenant
func (s *Service) CreateTenant(req CreateTenantRequest) (*Tenant, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Normalize subdomain
	req.Subdomain = strings.ToLower(strings.TrimSpace(req.Subdomain))

	// Check subdomain availability
	if exists, err := s.repo.SubdomainExists(req.Subdomain); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.New("subdomain already taken")
	}

	// Validate subdomain format (additional business rules)
	if err := s.validateSubdomain(req.Subdomain); err != nil {
		return nil, err
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
		Currency:    "BDT",
		Language:    "bn",
		Timezone:    "Asia/Dhaka",
	}

	return s.repo.Save(tenant)
}

// GetTenant retrieves a tenant by ID
func (s *Service) GetTenant(id string) (*Tenant, error) {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid tenant ID")
	}

	return s.repo.FindByID(tenantID)
}

// GetTenantBySubdomain retrieves a tenant by subdomain
func (s *Service) GetTenantBySubdomain(subdomain string) (*Tenant, error) {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	if subdomain == "" {
		return nil, errors.New("subdomain is required")
	}

	return s.repo.FindBySubdomain(subdomain)
}

// UpdateTenant updates tenant information
func (s *Service) UpdateTenant(id string, req UpdateTenantRequest) (*Tenant, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid tenant ID")
	}

	// Get existing tenant
	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid tenant ID")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return nil, err
	}

	// Validate plan change
	if !s.canChangeToPlan(tenant.Plan, req.Plan) {
		return nil, errors.New("plan change not allowed")
	}

	tenant.Plan = req.Plan
	
	// Update limits based on new plan
	tenant.ProductLimit = s.getProductLimitForPlan(req.Plan)
	tenant.StorageLimit = tenant.GetStorageLimit()

	return s.repo.Update(tenant)
}

// ListTenants returns a paginated list of tenants
func (s *Service) ListTenants(offset, limit int) ([]*Tenant, int64, error) {
	return s.repo.List(offset, limit)
}

// DeactivateTenant deactivates a tenant
func (s *Service) DeactivateTenant(id string) error {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid tenant ID")
	}

	return s.repo.UpdateStatus(tenantID, StatusInactive)
}

// ActivateTenant activates a tenant
func (s *Service) ActivateTenant(id string) error {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid tenant ID")
	}

	return s.repo.UpdateStatus(tenantID, StatusActive)
}

// GetTenantStats retrieves comprehensive statistics for a tenant
func (s *Service) GetTenantStats(id string) (*TenantStatsResponse, error) {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid tenant ID")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return nil, err
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
		return errors.New("invalid tenant ID")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return err
	}

	// Check if tenant plan supports custom domains
	if !tenant.CanUsePremiumFeatures() {
		return errors.New("custom domain requires premium or enterprise plan")
	}

	// Normalize domain
	domain = strings.ToLower(strings.TrimSpace(domain))

	// Check if domain is already taken
	if exists, err := s.repo.CustomDomainExists(domain); err != nil {
		return err
	} else if exists {
		return errors.New("domain already in use")
	}

	// TODO: Add DNS validation logic here
	// - Check if domain points to our servers
	// - Validate SSL certificate
	// - Verify ownership

	tenant.CustomDomain = domain
	_, err = s.repo.Update(tenant)
	return err
}

// SuspendTenant suspends a tenant (e.g., for non-payment)
func (s *Service) SuspendTenant(id string, reason string) error {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid tenant ID")
	}

	// TODO: Log suspension reason and send notification
	return s.repo.UpdateStatus(tenantID, StatusSuspended)
}

// GetPlanUpgradeOptions returns available upgrade options for a tenant
func (s *Service) GetPlanUpgradeOptions(id string) ([]Plan, error) {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid tenant ID")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return nil, err
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
		return false, err
	}

	exists, err := s.repo.SubdomainExists(subdomain)
	return !exists, err
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
		return errors.New("invalid tenant ID")
	}

	tenant, err := s.repo.FindByID(tenantID)
	if err != nil {
		return err
	}

	// Update limits based on current plan
	tenant.ProductLimit = s.getProductLimitForPlan(tenant.Plan)
	tenant.StorageLimit = tenant.GetStorageLimit()
	tenant.BandwidthLimit = tenant.GetBandwidthLimit()
	tenant.UpdatedAt = time.Now()

	_, err = s.repo.Update(tenant)
	return err
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
	// Business rules for plan changes
	// TODO: Implement plan change validation logic
	return true
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

// TODO: Add integration methods when other modules are ready
// - GetProductUsage(tenantID uuid.UUID) (int, error)
// - GetOrderMetrics(tenantID uuid.UUID) (*OrderMetrics, error)
// - GetStorageUsage(tenantID uuid.UUID) (int64, error)
// - GetBandwidthUsage(tenantID uuid.UUID) (int64, error)
// - SendWelcomeEmail(tenant *Tenant) error
// - SendPlanUpgradeNotification(tenant *Tenant) error
// - CalculateMonthlyBill(tenantID uuid.UUID) (*BillingInfo, error)
