package search

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the search module
type Module struct {
	handler *Handler
}

// NewModule creates a new search module
func NewModule(db *gorm.DB) *Module {
	// Initialize repository
	repo := NewRepository(db)
	
	// Initialize service
	service := NewService(repo)
	
	// Initialize handler
	handler := NewHandler(service)
	
	return &Module{
		handler: handler,
	}
}

// RegisterRoutes registers all search routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the search handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}