package billing

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ecommerce-saas/internal/analytics"
	"ecommerce-saas/internal/contact"
	"ecommerce-saas/internal/payment"
)

// Module represents the billing module
type Module struct {
	repository Repository
	service    Service
	handler    Handler
}

// NewModule creates a new billing module instance
func NewModule(db *gorm.DB, paymentService payment.Service, contactService contact.Service, analyticsService analytics.Service) *Module {
	repo := NewBillingRepository(db)
	service := NewBillingService(repo, paymentService, contactService, analyticsService)
	handler := NewHandler(service)

	return &Module{
		repository: repo,
		service:    service,
		handler:    handler,
	}
}

// RegisterRoutes registers all billing routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the billing handler
func (m *Module) GetHandler() Handler {
	return m.handler
}

// GetService returns the billing service
func (m *Module) GetService() Service {
	return m.service
}

// GetRepository returns the billing repository
func (m *Module) GetRepository() Repository {
	return m.repository
}
