package tenant

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

// Request/Response DTOs
type CreateTenantRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Subdomain   string `json:"subdomain" validate:"required,min=3,max=50,alphanum"`
	Description string `json:"description,omitempty" validate:"max=500"`
	Phone       string `json:"phone,omitempty" validate:"max=20"`
	Email       string `json:"email,omitempty" validate:"email"`
	Address     string `json:"address,omitempty" validate:"max=255"`
}

type UpdateTenantRequest struct {
	Name         string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description  string `json:"description,omitempty" validate:"max=500"`
	Phone        string `json:"phone,omitempty" validate:"max=20"`
	Email        string `json:"email,omitempty" validate:"email"`
	Address      string `json:"address,omitempty" validate:"max=255"`
	CustomDomain string `json:"custom_domain,omitempty" validate:"omitempty,fqdn"`
}

type UpdatePlanRequest struct {
	Plan Plan `json:"plan" validate:"required"`
}

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

// TODO: Add more service methods
// - ValidateCustomDomain(domain string) error
// - GetUsageStats(tenantID uuid.UUID) (*UsageStats, error)
// - CalculateBilling(tenantID uuid.UUID, month int, year int) (*BillingInfo, error)
// - SendWelcomeEmail(tenant *Tenant) error
