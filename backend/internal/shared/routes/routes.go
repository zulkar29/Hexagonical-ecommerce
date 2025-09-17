package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ecommerce-saas/internal/address"
	"ecommerce-saas/internal/admin"
	"ecommerce-saas/internal/analytics"
	"ecommerce-saas/internal/billing"
	"ecommerce-saas/internal/cart"
	"ecommerce-saas/internal/category"
	"ecommerce-saas/internal/components"
	"ecommerce-saas/internal/contact"
	"ecommerce-saas/internal/content"
	"ecommerce-saas/internal/discount"
	"ecommerce-saas/internal/finance"
	"ecommerce-saas/internal/marketing"
	"ecommerce-saas/internal/notification"
	"ecommerce-saas/internal/observability"
	"ecommerce-saas/internal/order"
	"ecommerce-saas/internal/payment"
	"ecommerce-saas/internal/platform"
	"ecommerce-saas/internal/product"
	"ecommerce-saas/internal/returns"
	"ecommerce-saas/internal/reviews"
	"ecommerce-saas/internal/search"
	"ecommerce-saas/internal/security"

	"ecommerce-saas/internal/settings"
	"ecommerce-saas/internal/shipping"
	"ecommerce-saas/internal/support"
	"ecommerce-saas/internal/tenant"
	"ecommerce-saas/internal/user"
	"ecommerce-saas/internal/webhook"
	"ecommerce-saas/internal/wishlist"
	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/shared/middleware"
	"ecommerce-saas/internal/shared/utils"
)

// RouteConfig holds dependencies for route setup
type RouteConfig struct {
	DB         *gorm.DB
	Config     *config.Config
	JWTManager *utils.JWTManager
}

// ServiceContainer holds shared service instances to avoid duplications
type ServiceContainer struct {
	// Repositories
	ProductRepo   product.Repository
	DiscountRepo  discount.Repository
	ContactRepo   contact.Repository
	AnalyticsRepo analytics.Repository
	ShippingRepo  *shipping.Repository
	PaymentRepo   payment.Repository
	
	// Services
	ProductService   *product.Service
	DiscountService  discount.Service
	ContactService   contact.Service
	AnalyticsService analytics.Service
	ShippingService  *shipping.Service
	PaymentService   payment.Service
}

// NewServiceContainer creates and initializes shared services
func NewServiceContainer(cfg *RouteConfig) *ServiceContainer {
	// Initialize repositories
	productRepo := product.NewRepository(cfg.DB)
	discountRepo := discount.NewRepository(cfg.DB)
	contactRepo := contact.NewRepository(cfg.DB)
	analyticsRepo := analytics.NewRepository(cfg.DB)
	shippingRepo := shipping.NewRepository(cfg.DB)
	paymentRepo := payment.NewRepository(cfg.DB)
	
	// Initialize services
	productService := product.NewService(productRepo)
	discountService := discount.NewService(discountRepo)
	contactService := contact.NewService(contactRepo)
	analyticsService := analytics.NewService(analyticsRepo)
	shippingService := shipping.NewService(shippingRepo)
	paymentService := payment.NewService(paymentRepo, cfg.Config)
	
	return &ServiceContainer{
		ProductRepo:      productRepo,
		DiscountRepo:     discountRepo,
		ContactRepo:      contactRepo,
		AnalyticsRepo:    analyticsRepo,
		ShippingRepo:     shippingRepo,
		PaymentRepo:      paymentRepo,
		ProductService:   productService,
		DiscountService:  discountService,
		ContactService:   contactService,
		AnalyticsService: analyticsService,
		ShippingService:  shippingService,
		PaymentService:   paymentService,
	}
}

// SetupRoutes configures all application routes
func SetupRoutes(r *gin.Engine, cfg *RouteConfig) {
	// Initialize shared service container to avoid duplications
	services := NewServiceContainer(cfg)
	
	// Initialize observability module for middleware
	observabilityModule := observability.NewModule(cfg.DB)
	observabilityService := observabilityModule.GetService()
	
	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
	r.Use(observabilityService.ObservabilityLoggingMiddleware())
	r.Use(observabilityService.MetricsMiddleware())
	r.Use(observabilityService.TracingMiddleware())
	r.Use(middleware.CORSMiddleware(
		[]string{"*"}, // Allow all origins in development
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
	
	// Initialize user module
	userModule := user.NewModule(cfg.DB, cfg.JWTManager)

	// Public routes (no authentication required)
	public := v1.Group("")
	{
		// User authentication routes
		userModule.RegisterRoutes(public)
	}

	// Protected routes (authentication required)
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTManager))
	
	// Create tenant isolation middleware
	tenantIsolation := middleware.NewTenantIsolationMiddleware(cfg.DB, "esass.com") // TODO: Make configurable
	protected.Use(tenantIsolation.ResolveTenantWithIsolation())
	protected.Use(tenantIsolation.ValidateTenantAccess())
	{
		// Core tenant management
		setupTenantRoutes(protected, cfg)
		
		// Product catalog management
		setupProductRoutes(protected, cfg)
		setupCategoryRoutes(protected, cfg)
		setupReviewsRoutes(protected, cfg)
		
		// Order and commerce flow
		setupOrderRoutes(protected, cfg, services)
		setupCartRoutes(protected, cfg, services)
		setupPaymentRoutes(protected, cfg)
		setupShippingRoutes(protected, cfg, services)
		setupReturnsRoutes(protected, cfg)
		
		// Customer engagement
		setupDiscountRoutes(protected, cfg, services)
		setupMarketingRoutes(protected, cfg)
		setupNotificationRoutes(protected, cfg)
		setupWishlistRoutes(protected, cfg)
		
		// Business operations
		setupFinanceRoutes(protected, cfg)
		setupBillingRoutes(protected, cfg, services)
		setupAnalyticsRoutes(protected, cfg, services)
		
		// Customer service
		setupSupportRoutes(protected, cfg)
		setupContactRoutes(protected, cfg, services)
		
		// System management
		setupAdminRoutes(protected, cfg)
		setupSettingsRoutes(protected, cfg)
		setupSecurityRoutes(protected, cfg)
		setupObservabilityRoutes(protected, cfg)
		
		// Content and platform
		setupContentRoutes(protected, cfg)
		setupComponentsRoutes(protected, cfg)
		setupPlatformRoutes(protected, cfg)
		setupSearchRoutes(protected, cfg)
		setupWebhookRoutes(protected, cfg)
		
		// Utility
		setupAddressRoutes(protected, cfg)
	}

	// Public routes (for storefront)
	storefront := v1.Group("/public")
	// Create tenant isolation middleware for public routes
	publicTenantIsolation := middleware.NewTenantIsolationMiddleware(cfg.DB, "esass.com") // TODO: Make configurable
	storefront.Use(publicTenantIsolation.ResolveTenantWithIsolation())
	{
		// Public product routes (no auth needed for browsing)
		setupPublicProductRoutes(storefront, cfg)
		
		// Public category routes (no auth needed for browsing)
		setupPublicCategoryRoutes(storefront, cfg)
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
	
	// TODO: Initialize order module for public tracking - requires service dependencies
	
	// Public product routes (read-only, no auth required)
	// Note: Tenant isolation is already applied at the storefront group level
	public := v1.Group("")
	{
		// Public product browsing endpoints
		public.GET("/products", productModule.Handler.GetPublicProducts)
		public.GET("/products/search", productModule.Handler.SearchProducts)
		public.GET("/products/slug/:slug", productModule.Handler.GetProductBySlug)
		public.GET("/products/:id", productModule.Handler.GetPublicProduct)
		public.GET("/products/:id/variants", productModule.Handler.GetProductVariants)
		
		// Public category browsing - handled by category module
	setupPublicCategoryRoutes(public, cfg)
		
		// TODO: Public order tracking (no auth required) - requires order module
		// public.GET("/orders/track/:number", orderModule.Handler.TrackOrder)
		// public.GET("/orders/number/:number", orderModule.Handler.GetOrderByNumber)
		
		// Public settings (no auth required)
		settingsModule := settings.NewModule(cfg.DB)
		public.GET("/settings", settingsModule.GetHandler().GetPublicSettings)
		
		// Public search (no auth required)
		searchModule := search.NewModule(cfg.DB)
		public.GET("/search", searchModule.GetHandler().Search)
		public.GET("/search/products", searchModule.GetHandler().SearchProducts)
		public.GET("/search/suggestions", searchModule.GetHandler().GetSuggestions)
	}
}

func setupOrderRoutes(v1 *gin.RouterGroup, cfg *RouteConfig, services *ServiceContainer) {
	// Initialize notification module (not in shared container as it's less commonly used)
	notificationModule := notification.NewModule(cfg.DB)
	notificationService := notificationModule.GetService()
	
	// Initialize order module using shared services
	orderModule := order.NewModule(cfg.DB, services.ProductService, services.DiscountService, services.PaymentService, notificationService)
	
	// Register order routes
	orderModule.RegisterRoutes(v1)
}

func setupPaymentRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize payment module
	paymentModule := payment.NewModule(cfg.DB, cfg.Config)
	
	// Register payment routes
	paymentModule.RegisterRoutes(v1)
}

func setupNotificationRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize notification module
	notificationModule := notification.NewModule(cfg.DB)
	
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

// Setup address routes
func setupAddressRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	addressRepo := address.NewGormRepository(cfg.DB)
	addressService := address.NewService(addressRepo)
	addressHandler := address.NewHandler(addressService)
	
	addressHandler.RegisterRoutes(v1)
}

// Setup analytics routes
func setupAnalyticsRoutes(v1 *gin.RouterGroup, cfg *RouteConfig, services *ServiceContainer) {
	analyticsHandler := analytics.NewHandler(services.AnalyticsService)
	
	analyticsHandler.RegisterRoutes(v1)
}



// Setup contact routes
func setupContactRoutes(v1 *gin.RouterGroup, cfg *RouteConfig, services *ServiceContainer) {
	contactHandler := contact.NewHandler(services.ContactService)
	
	contactHandler.RegisterRoutes(v1)
}

// Setup content routes
func setupContentRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	contentRepo := content.NewRepository(cfg.DB)
	contentService := content.NewService(contentRepo)
	contentHandler := content.NewHandler(contentService)
	
	contentHandler.RegisterRoutes(v1)
}

// Setup discount routes
func setupDiscountRoutes(v1 *gin.RouterGroup, cfg *RouteConfig, services *ServiceContainer) {
	discountHandler := discount.NewHandler(services.DiscountService)
	
	discountHandler.RegisterRoutes(v1)
}

// Setup marketing routes
func setupMarketingRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	marketingRepo := marketing.NewRepository(cfg.DB)
	marketingService := marketing.NewService(marketingRepo)
	marketingHandler := marketing.NewHandler(marketingService)
	
	marketingHandler.RegisterRoutes(v1)
}



// Setup reviews routes
func setupReviewsRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	reviewsRepo := reviews.NewRepository(cfg.DB)
	reviewsService := reviews.NewService(reviewsRepo)
	reviewsHandler := reviews.NewHandler(reviewsService)
	
	reviewsHandler.RegisterRoutes(v1)
}

// Setup shipping routes
func setupShippingRoutes(v1 *gin.RouterGroup, cfg *RouteConfig, services *ServiceContainer) {
	shippingHandler := shipping.NewHandler(services.ShippingService)
	
	shippingHandler.RegisterRoutes(v1)
}

// Setup support routes
func setupSupportRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	supportRepo := support.NewRepository(cfg.DB)
	supportService := support.NewService(supportRepo)
	supportHandler := support.NewHandler(supportService)
	
	supportHandler.RegisterRoutes(v1)
}




// Setup admin routes
func setupAdminRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	adminRepo := admin.NewRepository(cfg.DB)
	adminService := admin.NewService(adminRepo)
	adminHandler := admin.NewHandler(adminService)
	
	adminHandler.RegisterRoutes(v1)
}

// Setup billing routes
func setupBillingRoutes(v1 *gin.RouterGroup, cfg *RouteConfig, services *ServiceContainer) {
	// Initialize billing module using shared services
	billingModule := billing.NewModule(cfg.DB, services.PaymentService, services.ContactService, services.AnalyticsService)
	
	// Register billing routes
	billingModule.RegisterRoutes(v1)
}

// Setup cart routes
func setupCartRoutes(v1 *gin.RouterGroup, cfg *RouteConfig, services *ServiceContainer) {
	// Initialize cart module using shared services
	cartModule := cart.NewModule(cfg.DB, services.ProductService, services.DiscountService, services.ShippingService)
	
	// Register cart routes
	cartModule.RegisterRoutes(v1)
}

// Setup observability routes
func setupObservabilityRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize observability module
	observabilityModule := observability.NewModule(cfg.DB)
	
	// Register observability routes
	observabilityModule.RegisterRoutes(v1)
}

// Setup search routes
func setupSearchRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize search module
	searchModule := search.NewModule(cfg.DB)
	
	// Register search routes
	searchModule.RegisterRoutes(v1)
}

// Setup settings routes
func setupSettingsRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize settings module
	settingsModule := settings.NewModule(cfg.DB)
	
	// Register settings routes
	settingsModule.RegisterRoutes(v1)
}

// Setup webhook routes
func setupWebhookRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	webhookRepo := webhook.NewRepository(cfg.DB)
	// Use a default signing key for webhook validation
	signingKey := []byte("default-webhook-signing-key")
	webhookService := webhook.NewService(webhookRepo, signingKey)
	webhookHandler := webhook.NewHandler(webhookService)
	
	webhookHandler.RegisterRoutes(v1)
}

// Setup wishlist routes
func setupWishlistRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	wishlistRepo := wishlist.NewGormRepository(cfg.DB)
	wishlistService := wishlist.NewService(wishlistRepo)
	wishlistHandler := wishlist.NewHandler(wishlistService)
	
	wishlistHandler.RegisterRoutes(v1)
}

// Setup components routes
func setupComponentsRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize components module
	componentsModule := components.NewModule(cfg.DB)
	
	// Register components routes
	componentsModule.RegisterRoutes(v1)
}

// Setup security routes
func setupSecurityRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize security module
	securityModule := security.NewModule(cfg.DB)
	
	// Register security routes
	securityModule.RegisterRoutes(v1)
}

// Setup platform routes
func setupPlatformRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize platform module
	platformModule := platform.NewModule(cfg.DB)

	// Create platform route group with proper prefix
	platformGroup := v1.Group("/platform")

	// Register platform routes
	platformModule.RegisterRoutes(platformGroup)
}

func setupCategoryRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	categoryModule := category.NewModule(cfg.DB)
	categoryModule.RegisterRoutes(v1)
}

func setupPublicCategoryRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	categoryModule := category.NewModule(cfg.DB)
	handler := categoryModule.GetHandler()
	
	// Public category browsing routes
	v1.GET("/categories", handler.ListCategories)
	v1.GET("/categories/tree", handler.GetCategoryTree)
	v1.GET("/categories/featured", handler.GetFeaturedCategories)
	v1.GET("/categories/popular", handler.GetPopularCategories)
	v1.GET("/categories/slug/:slug", handler.GetCategoryBySlug)
	v1.GET("/categories/:id", handler.GetCategory)
	v1.GET("/categories/:id/products", handler.GetCategoryProducts)
}
