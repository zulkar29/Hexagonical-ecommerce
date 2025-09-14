package webhook

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the webhook module
type Module struct {
	handler    *Handler
	service    *Service
	repository *Repository
}

// NewModule creates a new webhook module
func NewModule(db *gorm.DB) *Module {
	// Initialize repository
	repo := NewRepository(db)

	// Initialize service with signing key from environment or default
	signingKey := getSigningKeyFromConfig()
	service := NewService(repo, signingKey)

	// Initialize handler
	handler := NewHandler(service)

	return &Module{
		handler:    handler,
		service:    service,
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

// getSigningKeyFromConfig retrieves the webhook signing key from environment variables
// Falls back to a default key if not configured (should be changed in production)
func getSigningKeyFromConfig() []byte {
	// Try to get from environment variable
	if key := os.Getenv("WEBHOOK_SIGNING_KEY"); key != "" {
		return []byte(key)
	}

	// Try alternative environment variable names
	if key := os.Getenv("WEBHOOK_SECRET"); key != "" {
		return []byte(key)
	}

	if key := os.Getenv("WEBHOOK_SECRET_KEY"); key != "" {
		return []byte(key)
	}

	// Default fallback (should be changed in production)
	return []byte("webhook-signing-key-change-in-production")
}
