package referral

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the referral module
type Module struct {
	Repository Repository
	Service    Service
	Handler    *Handler
}

// NewModule creates a new referral module with all dependencies wired up
func NewModule(db *gorm.DB) *Module {
	// Create repository
	repo := NewGormRepository(db)

	// Create service
	service := NewService(repo)

	// Create handler
	handler := NewHandler(service)

	return &Module{
		Repository: repo,
		Service:    service,
		Handler:    handler,
	}
}

// RegisterRoutes registers all referral routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.Handler.RegisterRoutes(router)
}

// GetService returns the referral service
func (m *Module) GetService() Service {
	return m.Service
}

// GetRepository returns the referral repository
func (m *Module) GetRepository() Repository {
	return m.Repository
}

// GetHandler returns the referral handler
func (m *Module) GetHandler() *Handler {
	return m.Handler
}