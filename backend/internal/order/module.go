package order

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Module represents the order module
type Module struct {
	Handler    *Handler
	Service    *Service
	Repository *Repository
}

// NewModule creates a new order module with all dependencies
func NewModule(db *gorm.DB) *Module {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	return &Module{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

// RegisterRoutes registers all order routes with the router
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.Handler.RegisterRoutes(router)
}

// Migrate runs database migrations for order module
func (m *Module) Migrate() error {
	return m.Repository.db.AutoMigrate(
		&Order{},
		&OrderItem{},
	)
}

// GetService returns the order service for integration with other modules
func (m *Module) GetService() *Service {
	return m.Service
}

// GetRepository returns the order repository for direct database access
func (m *Module) GetRepository() *Repository {
	return m.Repository
}
