package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	
	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/shared/middleware"
	"ecommerce-saas/internal/shared/utils"
	"ecommerce-saas/internal/user"
	"ecommerce-saas/internal/tenant"
	"ecommerce-saas/internal/product"
	"ecommerce-saas/internal/order"
)

// RouteConfig holds dependencies for route setup
type RouteConfig struct {
	DB         *gorm.DB
	Config     *config.Config
	JWTManager *utils.JWTManager
}

// SetupRoutes configures all application routes
func SetupRoutes(r *gin.Engine, cfg *RouteConfig) {
	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware(
		[]string{"*"}, // Allow all origins in development
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		[]string{"Authorization", "Content-Type", "X-Request-ID"},
		true,
	))
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": "2025-09-04T12:00:00Z",
			"version":   "1.0.0",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	
	// Initialize repositories and services
	userRepo := user.NewRepository(cfg.DB)
	userService := user.NewService(userRepo, cfg.JWTManager)
	userHandler := user.NewHandler(userService)

	// Public routes (no authentication required)
	public := v1.Group("")
	{
		// User authentication routes
		userHandler.RegisterRoutes(public)
	}

	// Protected routes (authentication required)
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTManager))
	{
		// TODO: Add protected routes here
		// tenant routes, product management, order management, etc.
	}

	// TODO: Implement other module routes when ready
	// setupTenantRoutes(v1, cfg)
	// setupProductRoutes(v1, cfg)
	// setupOrderRoutes(v1, cfg)
	// setupAnalyticsRoutes(v1, cfg)
}

// TODO: Implement other route setup functions
func setupTenantRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// tenantRepo := tenant.NewRepository(cfg.DB)
	// tenantService := tenant.NewService(tenantRepo)
	// tenantHandler := tenant.NewHandler(tenantService)
	
	// tenantGroup := v1.Group("/tenants")
	// tenantGroup.Use(middleware.AuthMiddleware(cfg.JWTManager))
	// tenantHandler.RegisterRoutes(tenantGroup)
}

func setupProductRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// productRepo := product.NewRepository(cfg.DB)
	// productService := product.NewService(productRepo)
	// productHandler := product.NewHandler(productService)
	
	// productGroup := v1.Group("/products")
	// productGroup.Use(middleware.TenantMiddleware())
	// productHandler.RegisterRoutes(productGroup)
}

func setupOrderRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// orderRepo := order.NewRepository(cfg.DB)
	// orderService := order.NewService(orderRepo)
	// orderHandler := order.NewHandler(orderService)
	
	// orderGroup := v1.Group("/orders")
	// orderGroup.Use(middleware.TenantMiddleware())
	// orderGroup.Use(middleware.AuthMiddleware(cfg.JWTManager))
	// orderHandler.RegisterRoutes(orderGroup)
}
		// analyticsGroup.Use(middleware.TenantMiddleware())
		// analyticsGroup.Use(middleware.AuthMiddleware())
		// analyticsHandler.RegisterRoutes(analyticsGroup)

		// TODO: Register payment routes
		// paymentGroup := v1.Group("/payments")
		// paymentGroup.Use(middleware.TenantMiddleware())
		// paymentGroup.Use(middleware.AuthMiddleware())
		// paymentHandler.RegisterRoutes(paymentGroup)

		// TODO: Register notification routes
		// notificationGroup := v1.Group("/notifications")
		// notificationGroup.Use(middleware.TenantMiddleware())
		// notificationGroup.Use(middleware.AuthMiddleware())
		// notificationHandler.RegisterRoutes(notificationGroup)
	}

	// Public routes (no auth required)
	public := r.Group("/api/public")
	{
		// TODO: Add public routes like product catalog, registration, etc.
		public.GET("/products", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Public product catalog"})
		})
	}

	// Webhook routes
	webhooks := r.Group("/webhooks")
	{
		// TODO: Add webhook handlers for payments, notifications, etc.
		webhooks.POST("/stripe", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Stripe webhook"})
		})
		webhooks.POST("/paypal", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "PayPal webhook"})
		})
	}
}
