package settings

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the settings module
type Module struct {
	repository Repository
	service    Service
	handler    *Handler
}

// NewModule creates a new settings module
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

// RegisterRoutes registers all settings routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the settings handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}

// GetService returns the settings service
func (m *Module) GetService() Service {
	return m.service
}

// GetRepository returns the settings repository
func (m *Module) GetRepository() Repository {
	return m.repository
}