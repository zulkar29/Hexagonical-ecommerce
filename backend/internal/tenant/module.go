package tenant

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the tenant module
type Module struct {
	Repository Repository
	Service    Service
	Handler    *Handler
}

// NewModule creates a new tenant module with all dependencies
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

// RegisterRoutes registers all tenant routes
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	m.Handler.RegisterRoutes(r)
}
