package platform


import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the platform module
type Module struct {
	handler *Handler
	service Service
	repo    Repository
}

// NewModule creates a new platform module
func NewModule(db *gorm.DB) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{
		handler: handler,
		service: service,
		repo:    repo,
	}
}

// RegisterRoutes registers platform routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetService returns the platform service
func (m *Module) GetService() Service {
	return m.service
}

// GetRepository returns the platform repository
func (m *Module) GetRepository() Repository {
	return m.repo
}

// GetHandler returns the platform handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}