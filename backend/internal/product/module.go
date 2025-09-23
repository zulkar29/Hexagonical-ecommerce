package product

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the product module
type Module struct {
	Repository Repository
	Service    *Service
	Handler    *Handler
}

// NewModule creates a new product module with all dependencies
func NewModule(db *gorm.DB) *Module {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}

// RegisterRoutes registers all product routes with the router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.Handler.RegisterRoutes(router)
}

// Migrate runs database migrations for product module
func (m *Module) Migrate() error {
	// Skip GORM auto-migration to avoid conflicts with SQL migrations
	// The database schema is managed by SQL migration files
	return nil
}
