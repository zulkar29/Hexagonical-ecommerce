package category

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Module represents the category module
type Module struct {
	repository *Repository
	service    *Service
	handler    *Handler
}

// NewModule creates a new category module instance
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

// RegisterRoutes registers all category routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the category handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}

// GetService returns the category service
func (m *Module) GetService() *Service {
	return m.service
}

// GetRepository returns the category repository
func (m *Module) GetRepository() *Repository {
	return m.repository
}

// Migrate runs database migrations for category module
func (m *Module) Migrate() error {
	// Skip GORM auto-migration to avoid conflicts with SQL migrations
	// The database schema is managed by SQL migration files
	return nil
}