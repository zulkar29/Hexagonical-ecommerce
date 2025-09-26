package notification

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/shared/email"
)

type Module struct {
	repository Repository
	service    Service
	handler    *Handler
}

func NewModule(db *gorm.DB, cfg *config.Config) *Module {
	repository := NewRepository(db)
	emailService := email.NewService(cfg)
	service := NewService(repository, emailService)
	handler := NewHandler(service)

	log.Println("✅ Notification module initialized successfully")

	return &Module{
		repository: repository,
		service:    service,
		handler:    handler,
	}
}

func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	m.handler.RegisterRoutes(r)
	log.Println("✅ Notification routes registered at /api/notifications")
}

func (m *Module) GetService() Service {
	return m.service
}

func (m *Module) GetRepository() Repository {
	return m.repository
}

func (m *Module) GetHandler() *Handler {
	return m.handler
}
