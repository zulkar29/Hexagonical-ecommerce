package security

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the security module
type Module struct {
	repository SecurityRepository
	service    SecurityService
	handler    *Handler
}

// NewModule creates a new security module
func NewModule(db *gorm.DB) *Module {
	repository := NewSecurityRepository(db)
	service := NewSecurityService(repository)
	handler := NewHandler(service)

	return &Module{
		repository: repository,
		service:    service,
		handler:    handler,
	}
}

// RegisterRoutes registers the security module routes
func (m *Module) RegisterRoutes(router gin.IRouter) {
	securityGroup := router.Group("/security")
	m.handler.RegisterRoutes(securityGroup)
}

// Migrate runs the security module migrations
func (m *Module) Migrate() error {
	return m.repository.Migrate()
}

// GetRepository returns the security repository
func (m *Module) GetRepository() SecurityRepository {
	return m.repository
}

// GetService returns the security service
func (m *Module) GetService() SecurityService {
	return m.service
}

// GetHandler returns the security handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}