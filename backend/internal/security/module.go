package security

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the security module
type Module struct {
	service *SecurityService
	handler *Handler
}

// NewModule creates a new security module
// Note: This is a simplified module. In practice, you'd inject UserRepository from user module
func NewModule(db *gorm.DB, userRepo UserRepository) *Module {
	service := NewSecurityService(db, userRepo)
	handler := NewHandler(service)

	return &Module{
		service: service,
		handler: handler,
	}
}

// RegisterRoutes registers the security module routes
func (m *Module) RegisterRoutes(router gin.IRouter) {
	securityGroup := router.Group("/security")
	m.handler.SetupRoutes(securityGroup)
}

// Note: Database schema is handled by raw SQL migrations in /migrations directory

// GetService returns the security service
func (m *Module) GetService() *SecurityService {
	return m.service
}

// GetHandler returns the security handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}
