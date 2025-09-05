package contact

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Module represents the contact module
type Module struct {
	repository Repository
	service    Service
	handler    *Handler
}

// NewModule creates a new contact module instance
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

// RegisterRoutes registers all contact routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the contact handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}

// GetService returns the contact service
func (m *Module) GetService() Service {
	return m.service
}

// GetRepository returns the contact repository
func (m *Module) GetRepository() Repository {
	return m.repository
}