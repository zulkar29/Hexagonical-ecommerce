package support

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Module represents the support module
type Module struct {
	repository Repository
	service    Service
	handler    *Handler
}

// NewModule creates a new support module instance
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

// RegisterRoutes registers all support routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the support handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}

// GetService returns the support service
func (m *Module) GetService() Service {
	return m.service
}

// GetRepository returns the support repository
func (m *Module) GetRepository() Repository {
	return m.repository
}