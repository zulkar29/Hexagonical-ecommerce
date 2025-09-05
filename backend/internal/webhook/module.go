package webhook

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the webhook module
type Module struct {
	handler *Handler
	service *Service
	repository *Repository
}

// NewModule creates a new webhook module
func NewModule(db *gorm.DB) *Module {
	// Initialize repository
	repo := NewRepository(db)
	
	// Initialize service with signing key
	// TODO: Get signing key from config
	signingKey := []byte("webhook-signing-key-change-in-production")
	service := NewService(repo, signingKey)
	
	// Initialize handler
	handler := NewHandler(service)
	
	return &Module{
		handler: handler,
		service: service,
		repository: repo,
	}
}

// RegisterRoutes registers all webhook routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the webhook handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}

// GetService returns the webhook service
func (m *Module) GetService() *Service {
	return m.service
}

// GetRepository returns the webhook repository
func (m *Module) GetRepository() *Repository {
	return m.repository
}