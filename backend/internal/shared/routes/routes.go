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
	"ecommerce-saas/internal/payment"
	"ecommerce-saas/internal/notification"
	"ecommerce-saas/internal/finance"
	"ecommerce-saas/internal/returns"
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
	protected.Use(middleware.TenantMiddleware()) // Add tenant resolution middleware
	{
		// Setup tenant routes
		setupTenantRoutes(protected, cfg)
		
		// Setup product routes
		setupProductRoutes(protected, cfg)
		
		// Setup order routes
		setupOrderRoutes(protected, cfg)
		
		// Setup payment routes
		setupPaymentRoutes(protected, cfg)
		
		// Setup notification routes
		setupNotificationRoutes(protected, cfg)
		
		// Setup finance routes
		setupFinanceRoutes(protected, cfg)
		
		// Setup returns routes
		setupReturnsRoutes(protected, cfg)
		
		// TODO: Add other protected routes here
		// setupAnalyticsRoutes(protected, cfg)
		// setupBillingRoutes(protected, cfg)
		// setupMarketingRoutes(protected, cfg)
		// setupSupportRoutes(protected, cfg)
		// setupObservabilityRoutes(protected, cfg)
	}

	// Public routes (for storefront)
	storefront := v1.Group("/public")
	{
		// Public product routes (no auth needed for browsing)
		setupPublicProductRoutes(storefront, cfg)
	}

}

// Setup tenant routes
func setupTenantRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize tenant module
	tenantModule := tenant.NewModule(cfg.DB)
	
	// Register tenant routes
	tenantModule.Handler.RegisterRoutes(v1)
}

func setupProductRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize product module
	productModule := product.NewModule(cfg.DB)
	
	// Register product routes
	productModule.RegisterRoutes(v1)
}

func setupPublicProductRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize product module for public access
	productModule := product.NewModule(cfg.DB)
	
	// Initialize order module for public tracking
	orderModule := order.NewModule(cfg.DB)
	
	// Public product routes (read-only, no auth required)
	public := v1.Group("")
	public.Use(middleware.TenantMiddleware()) // Still need tenant resolution
	{
		// Public product browsing endpoints
		public.GET("/products", productModule.Handler.ListProducts)
		public.GET("/products/search", productModule.Handler.SearchProducts)
		public.GET("/products/slug/:slug", productModule.Handler.GetProductBySlug)
		public.GET("/products/:id", productModule.Handler.GetProduct)
		public.GET("/products/:id/variants", productModule.Handler.GetProductVariants)
		
		// Public category browsing
		public.GET("/categories", productModule.Handler.ListCategories)
		public.GET("/categories/root", productModule.Handler.GetRootCategories)
		public.GET("/categories/:id", productModule.Handler.GetCategory)
		public.GET("/categories/:id/children", productModule.Handler.GetCategoryChildren)
		
		// Public order tracking (no auth required)
		public.GET("/orders/track/:number", orderModule.Handler.TrackOrder)
		public.GET("/orders/number/:number", orderModule.Handler.GetOrderByNumber)
	}
}

func setupOrderRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize order module
	orderModule := order.NewModule(cfg.DB)
	
	// Register order routes
	orderModule.RegisterRoutes(v1)
}

func setupPaymentRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Convert GORM DB to sql.DB for payment module
	sqlDB, err := cfg.DB.DB()
	if err != nil {
		panic("Failed to get sql.DB from GORM: " + err.Error())
	}
	
	// Initialize payment module
	paymentModule := payment.NewModule(sqlDB)
	
	// Register payment routes
	paymentModule.RegisterRoutes(v1)
}

func setupNotificationRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Convert GORM DB to sql.DB for notification module
	sqlDB, err := cfg.DB.DB()
	if err != nil {
		panic("Failed to get sql.DB from GORM: " + err.Error())
	}
	
	// Initialize notification module
	notificationModule := notification.NewModule(sqlDB)
	
	// Register notification routes
	notificationModule.RegisterRoutes(v1)
}

func setupFinanceRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize finance module
	financeModule := finance.NewModule(cfg.DB)
	
	// Register finance routes
	financeModule.RegisterRoutes(v1)
}

func setupReturnsRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize returns module
	returnsModule := returns.NewModule(cfg.DB)
	
	// Register returns routes
	returnsModule.RegisterRoutes(v1)
}
