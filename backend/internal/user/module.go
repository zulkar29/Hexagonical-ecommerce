package user

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"ecommerce-saas/internal/shared/utils"
)

// Module represents the user module
type Module struct {
	repository Repository
	service    Service
	handler    *Handler
}

// NewModule creates a new user module instance
func NewModule(db *gorm.DB, jwtManager *utils.JWTManager) *Module {
	repo := NewRepository(db)
	svc := NewService(repo, jwtManager)
	handler := NewHandler(svc)

	return &Module{
		repository: repo,
		service:    *svc,
		handler:    handler,
	}
}

// RegisterRoutes registers all user routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the user handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}

// GetService returns the user service
func (m *Module) GetService() *Service {
	return &m.service
}

// GetRepository returns the user repository
func (m *Module) GetRepository() Repository {
	return m.repository
}