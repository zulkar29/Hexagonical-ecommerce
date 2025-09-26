package social

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the social commerce module
type Module struct {
	Handler    *Handler
	Service    *Service
	Repository Repository
}

// NewModule creates a new social commerce module
func NewModule(db *gorm.DB) *Module {
	// Create repository
	repo := NewRepository(db)

	// Create service
	service := NewService(repo)

	// Create handler
	handler := NewHandler(service)

	return &Module{
		Handler:    handler,
		Service:    service,
		Repository: repo,
	}
}

// RegisterRoutes registers all social commerce routes
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	m.Handler.RegisterRoutes(r)
}

// Migrate runs database migrations for social commerce tables
func (m *Module) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&SocialIntegration{},
		&SocialProduct{},
		&SocialAnalytics{},
	)
}