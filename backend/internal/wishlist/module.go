package wishlist

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Module represents the wishlist module
type Module struct {
	repository Repository
	service    Service
	handler    *Handler
}

// NewModule creates a new wishlist module instance
func NewModule(db *gorm.DB) *Module {
	repo := NewGormRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	return &Module{
		repository: repo,
		service:    svc,
		handler:    handler,
	}
}

// RegisterRoutes registers all wishlist routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the wishlist handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}

// GetService returns the wishlist service
func (m *Module) GetService() Service {
	return m.service
}

// GetRepository returns the wishlist repository
func (m *Module) GetRepository() Repository {
	return m.repository
}