package product

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Module represents the product module
type Module struct {
	Handler *Handler
	Service *Service
	Repository Repository
	InventoryService *InventoryService
}

// NewModule creates a new product module with all dependencies
func NewModule(db *gorm.DB) *Module {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)
	inventoryService := NewInventoryService(repository)

	return &Module{
		Handler:    handler,
		Service:    service,
		Repository: repository,
		InventoryService: inventoryService,
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
