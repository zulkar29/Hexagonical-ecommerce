package loyalty

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Module represents the loyalty module
type Module struct {
	Repository Repository
	Service    Service
	Handler    *Handler
}

// NewModule creates a new loyalty module with all dependencies
func NewModule(db *gorm.DB) *Module {
	// Initialize repository
	repo := NewGormRepository(db)

	// Initialize service
	service := NewService(repo)

	// Initialize handler
	handler := NewHandler(service)

	return &Module{
		Repository: repo,
		Service:    service,
		Handler:    handler,
	}
}

// RegisterRoutes registers all loyalty routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.Handler.RegisterRoutes(router)
}

// Migrate runs database migrations for loyalty tables
func (m *Module) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&LoyaltyProgram{},
		&LoyaltyAccount{},
		&LoyaltyTransaction{},
		&LoyaltyReward{},
	)
}