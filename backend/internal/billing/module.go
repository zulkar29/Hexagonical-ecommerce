package billing

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Module represents the billing module
type Module struct {
	repository Repository
	service    Service
	handler    Handler
}

// NewModule creates a new billing module instance
func NewModule(db *gorm.DB) *Module {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	return &Module{
		repository: repo,
		service:    svc,
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