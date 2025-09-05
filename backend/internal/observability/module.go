package observability

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the observability module
type Module struct {
	handler *ObservabilityHandler
	service *ObservabilityService
	repository ObservabilityRepository
}

// NewModule creates a new observability module
func NewModule(db *gorm.DB) *Module {
	// Initialize repository
	repo := NewObservabilityRepository(db)
	
	// Initialize service (needs logger, metrics, tracer)
	// Initialize logger, metrics, and tracer
	logger := NewLogger("ecommerce-saas", "1.0.0", LogLevelInfo)
	metrics := NewMetricsCollector(MetricsConfig{})
	tracer := NewTracer(TracingConfig{})
	service := NewObservabilityService(logger, metrics, tracer, repo)
	
	// Initialize handler
	handler := NewObservabilityHandler(service)
	
	return &Module{
		handler: handler,
		service: service,
		repository: repo,
	}
}

// RegisterRoutes registers all observability routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}

// GetHandler returns the observability handler
func (m *Module) GetHandler() *ObservabilityHandler {
	return m.handler
}

// GetService returns the observability service
func (m *Module) GetService() *ObservabilityService {
	return m.service
}

// GetRepository returns the observability repository
func (m *Module) GetRepository() ObservabilityRepository {
	return m.repository
}