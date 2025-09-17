package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ecommerce-saas/internal/shared/utils"
)

// Module represents the user module
type Module struct {
	Repository RepositoryInterface
	Service    *Service
	Handler    *Handler
}

// NewModule creates a new user module instance
func NewModule(db *gorm.DB, jwtManager *utils.JWTManager) *Module {
	repo := NewRepository(db)
	securityService := NewSecurityService(db)
	svc := NewService(repo, jwtManager, securityService)
	handler := NewHandler(svc)

	return &Module{
		Repository: repo,
		Service:    svc,
		Handler:    handler,
	}
}

// RegisterRoutes registers all user routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.Handler.RegisterRoutes(router)
}

// GetHandler returns the user handler
func (m *Module) GetHandler() *Handler {
	return m.Handler
}

// GetService returns the user service
func (m *Module) GetService() *Service {
	return m.Service
}

// GetRepository returns the user repository
func (m *Module) GetRepository() RepositoryInterface {
	return m.Repository
}