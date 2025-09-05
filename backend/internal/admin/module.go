package admin

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Module represents the admin module
type Module struct {
	repository Repository
	service    Service
	handler    *Handler
}

// NewModule creates a new admin module instance
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

// RegisterRoutes registers all admin routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the admin handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}

// GetService returns the admin service
func (m *Module) GetService() Service {
	return m.service
}

// GetRepository returns the admin repository
func (m *Module) GetRepository() Repository {
	return m.repository
}