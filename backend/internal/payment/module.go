package payment

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the payment module
type Module struct {
	Repository Repository
	Service    Service
	Handler    *Handler
}

// NewModule creates a new payment module with all dependencies
func NewModule(db *gorm.DB) *Module {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}

// RegisterRoutes registers all payment routes
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	m.Handler.RegisterRoutes(r)
}
