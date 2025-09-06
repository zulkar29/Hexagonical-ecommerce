package components

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module represents the components module
type Module struct {
	db      *gorm.DB
	handler *Handler
	service Service
	repo    Repository
}

// NewModule creates a new components module
func NewModule(db *gorm.DB) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)
	
	return &Module{
		db:      db,
		handler: handler,
		service: service,
		repo:    repo,
	}
}

// RegisterRoutes registers the component routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	// Component routes
	components := router.Group("/components")
	{
		components.POST("", m.handler.CreateComponent)
		components.GET("", m.handler.ListComponents)
		components.GET("/:id", m.handler.GetComponent)
		components.GET("/slug/:slug", m.handler.GetComponentBySlug)
		components.PUT("/:id", m.handler.UpdateComponent)
		components.DELETE("/:id", m.handler.DeleteComponent)
		components.POST("/:id/duplicate", m.handler.DuplicateComponent)
	}
	
	// Component instance routes
	instances := router.Group("/instances")
	{
		instances.POST("", m.handler.CreateInstance)
		instances.GET("", m.handler.ListInstances)
		instances.GET("/:id", m.handler.GetInstance)
	}
	
	// Theme routes
	themes := router.Group("/themes")
	{
		themes.GET("", m.handler.ListThemes)
		themes.GET("/active", m.handler.GetActiveTheme)
	}
	
	// Template routes
	templates := router.Group("/templates")
	{
		templates.GET("", m.handler.ListTemplates)
		templates.GET("/:id", m.handler.GetTemplate)
	}
	
	// Statistics routes
	stats := router.Group("/stats")
	{
		stats.GET("", m.handler.GetStats)
	}
}

// Migrate runs the database migrations for components
func (m *Module) Migrate() error {
	return m.db.AutoMigrate(
		&Component{},
		&ComponentTemplate{},
		&ComponentInstance{},
		&Theme{},
	)
}

// GetService returns the components service
func (m *Module) GetService() Service {
	return m.service
}

// GetRepository returns the components repository
func (m *Module) GetRepository() Repository {
	return m.repo
}

// GetHandler returns the components handler
func (m *Module) GetHandler() *Handler {
	return m.handler
}