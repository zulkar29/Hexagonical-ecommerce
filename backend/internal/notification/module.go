package notification

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

type Module struct {
	repository Repository
	service    Service
	handler    *Handler
}

func NewModule(db *sql.DB) *Module {
	repository := NewRepository(db)
	service := NewService(repository)
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
