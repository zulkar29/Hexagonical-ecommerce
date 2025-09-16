package cart

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the cart module
type Module struct {
	repository Repository
	service    *CartService
	handler    *Handler
}

// NewModule creates a new cart module instance
func NewModule(db *gorm.DB, productSvc ProductService, discountSvc DiscountService, shippingSvc ShippingService) *Module {
	repo := NewRepository(db)
	svc := NewCartService(repo, productSvc, discountSvc, shippingSvc)
	handler := NewHandler(svc)

	return &Module{
		repository: repo,
		service:    svc,
		handler:    handler,
	}
}

// RegisterRoutes registers all cart routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the cart handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}

// GetService returns the cart service
func (m *Module) GetService() ServiceInterface {
	return m.service
}

// GetRepository returns the cart repository
func (m *Module) GetRepository() Repository {
	return m.repository
}
