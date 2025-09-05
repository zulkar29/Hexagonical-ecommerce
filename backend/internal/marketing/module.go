package marketing

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Module represents the marketing module
type Module struct {
	repository Repository
	service    Service
	handler    *Handler
}

// NewModule creates a new marketing module instance
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

// RegisterRoutes registers all marketing routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the marketing handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}

// GetService returns the marketing service
func (m *Module) GetService() Service {
	return m.service
}

// GetRepository returns the marketing repository
func (m *Module) GetRepository() Repository {
	return m.repository
}