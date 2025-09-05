package tenant

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

// RepositoryInterface defines the contract for tenant data operations
type RepositoryInterface interface {
	Save(tenant *Tenant) (*Tenant, error)
	FindByID(id uuid.UUID) (*Tenant, error)
	FindBySubdomain(subdomain string) (*Tenant, error)
	FindByCustomDomain(domain string) (*Tenant, error)
	Update(tenant *Tenant) (*Tenant, error)
	UpdateStatus(id uuid.UUID, status Status) error
	SubdomainExists(subdomain string) (bool, error)
	CustomDomainExists(domain string) (bool, error)
	List(offset, limit int) ([]*Tenant, int64, error)
	ListByStatus(status Status, offset, limit int) ([]*Tenant, int64, error)
	ListByPlan(plan Plan, offset, limit int) ([]*Tenant, int64, error)
	Delete(id uuid.UUID) error
	GetActiveCount() (int64, error)
	GetCountByPlan(plan Plan) (int64, error)
	Search(query string, offset, limit int) ([]*Tenant, int64, error)
	GetExpiringTrials(days int) ([]*Tenant, error)
	BulkUpdateStatus(ids []uuid.UUID, status Status) error
	GetTenantStats(tenantID uuid.UUID) (*TenantStatsResponse, error)
}

// ServiceInterface defines the contract for tenant business logic
type ServiceInterface interface {
	CreateTenant(req CreateTenantRequest) (*Tenant, error)
	GetTenant(id string) (*Tenant, error)
	GetTenantBySubdomain(subdomain string) (*Tenant, error)
	UpdateTenant(id string, req UpdateTenantRequest) (*Tenant, error)
	UpdatePlan(id string, req UpdatePlanRequest) (*Tenant, error)
	ListTenants(offset, limit int) ([]*Tenant, int64, error)
	DeactivateTenant(id string) error
	ActivateTenant(id string) error
	GetTenantStats(id string) (*TenantStatsResponse, error)
	CheckSubdomainAvailability(subdomain string) (bool, error)
	GetPlanUpgradeOptions(id string) ([]Plan, error)
	SuspendTenant(id string, reason string) error
	ValidateCustomDomain(id, domain string) error
}

// Module represents the tenant module
type Module struct {
	Repository RepositoryInterface
	Service    ServiceInterface
	Handler    *Handler
}

// NewModule creates a new tenant module with all dependencies
func NewModule(db *gorm.DB) *Module {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}

// RegisterRoutes registers all tenant routes
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	m.Handler.RegisterRoutes(r)
}
