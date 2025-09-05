package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ecommerce-saas/internal/address"
	"ecommerce-saas/internal/admin"
	"ecommerce-saas/internal/analytics"
	"ecommerce-saas/internal/billing"
	// "ecommerce-saas/internal/cart" // Temporarily disabled due to interface compatibility issues
	"ecommerce-saas/internal/contact"
	"ecommerce-saas/internal/content"
	"ecommerce-saas/internal/discount"
	"ecommerce-saas/internal/finance"
	"ecommerce-saas/internal/loyalty"
	"ecommerce-saas/internal/marketing"
	"ecommerce-saas/internal/notification"
	"ecommerce-saas/internal/observability"
	"ecommerce-saas/internal/payment"
	"ecommerce-saas/internal/product"
	"ecommerce-saas/internal/returns"
	"ecommerce-saas/internal/reviews"
	"ecommerce-saas/internal/search"
	"ecommerce-saas/internal/settings"
	"ecommerce-saas/internal/shipping"
	"ecommerce-saas/internal/support"
	"ecommerce-saas/internal/tax"
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

// SetupRoutes configures all application routes
func SetupRoutes(r *gin.Engine, cfg *RouteConfig) {
	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
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
	protected.Use(middleware.TenantMiddleware(cfg.DB)) // Add tenant resolution middleware
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
		
		// Setup other protected routes
		setupAddressRoutes(protected, cfg)
		setupAdminRoutes(protected, cfg)
		setupAnalyticsRoutes(protected, cfg)
		setupBillingRoutes(protected, cfg)
		setupCartRoutes(protected, cfg)
		setupContactRoutes(protected, cfg)
		setupContentRoutes(protected, cfg)
		setupDiscountRoutes(protected, cfg)
		setupLoyaltyRoutes(protected, cfg)
		setupMarketingRoutes(protected, cfg)
		setupObservabilityRoutes(protected, cfg)
		setupReviewsRoutes(protected, cfg)
		setupSearchRoutes(protected, cfg)
		setupSettingsRoutes(protected, cfg)
		setupShippingRoutes(protected, cfg)
		setupSupportRoutes(protected, cfg)
		setupTaxRoutes(protected, cfg)
		setupWebhookRoutes(protected, cfg)
		setupWishlistRoutes(protected, cfg)
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
	
	// TODO: Initialize order module for public tracking - requires service dependencies
	// orderModule := order.NewModule(cfg.DB, productService, discountService, paymentService, inventoryService, notificationService)
	
	// Public product routes (read-only, no auth required)
	public := v1.Group("")
	public.Use(middleware.TenantMiddleware(cfg.DB)) // Still need tenant resolution
	{
		// Public product browsing endpoints
		public.GET("/products", productModule.Handler.GetPublicProducts)
		public.GET("/products/search", productModule.Handler.SearchProducts)
		public.GET("/products/slug/:slug", productModule.Handler.GetProductBySlug)
		public.GET("/products/:id", productModule.Handler.GetPublicProduct)
		public.GET("/products/:id/variants", productModule.Handler.GetProductVariants)
		
		// Public category browsing
		public.GET("/categories", productModule.Handler.GetPublicCategories)
		public.GET("/categories/root", productModule.Handler.GetRootCategories)
		public.GET("/categories/:id", productModule.Handler.GetCategory)
		public.GET("/categories/:id/children", productModule.Handler.GetCategoryChildren)
		
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

func setupOrderRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// TODO: Implement proper service dependencies
	// For now, comment out order routes to avoid compilation errors
	// orderModule := order.NewModule(cfg.DB, productService, discountService, paymentService, inventoryService, notificationService)
	// orderModule.RegisterRoutes(v1)
}

func setupPaymentRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize payment module
	paymentModule := payment.NewModule(cfg.DB)
	
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
func setupAnalyticsRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	analyticsRepo := analytics.NewRepository(cfg.DB)
	analyticsService := analytics.NewService(analyticsRepo)
	analyticsHandler := analytics.NewHandler(analyticsService)
	
	analyticsHandler.RegisterRoutes(v1)
}



// Setup contact routes
func setupContactRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	contactRepo := contact.NewRepository(cfg.DB)
	contactService := contact.NewService(contactRepo)
	contactHandler := contact.NewHandler(contactService)
	
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
func setupDiscountRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	discountRepo := discount.NewRepository(cfg.DB)
	discountService := discount.NewService(discountRepo)
	discountHandler := discount.NewHandler(discountService)
	
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
func setupShippingRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	shippingRepo := shipping.NewRepository(cfg.DB)
	shippingService := shipping.NewService(shippingRepo)
	shippingHandler := shipping.NewHandler(shippingService)
	
	shippingHandler.RegisterRoutes(v1)
}

// Setup support routes
func setupSupportRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	supportRepo := support.NewRepository(cfg.DB)
	supportService := support.NewService(supportRepo)
	supportHandler := support.NewHandler(supportService)
	
	supportHandler.RegisterRoutes(v1)
}



func setupTaxRoutes(protected *gin.RouterGroup, cfg *RouteConfig) {
	taxRepo := tax.NewGormRepository(cfg.DB)
	taxService := tax.NewService(taxRepo)
	taxHandler := tax.NewHandler(taxService)
	taxHandler.RegisterRoutes(protected)
}

// Setup loyalty routes
func setupLoyaltyRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// Initialize loyalty module
	loyaltyModule := loyalty.NewModule(cfg.DB)
	
	// Register loyalty routes
	loyaltyModule.RegisterRoutes(v1)
}

// Setup admin routes
func setupAdminRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	adminRepo := admin.NewRepository(cfg.DB)
	adminService := admin.NewService(adminRepo)
	adminHandler := admin.NewHandler(adminService)
	
	adminHandler.RegisterRoutes(v1)
}

// Setup billing routes
func setupBillingRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	billingRepo := billing.NewRepository(cfg.DB)
	billingService := billing.NewService(billingRepo)
	billingHandler := billing.NewHandler(billingService)
	
	billingHandler.RegisterRoutes(v1)
}

// Setup cart routes
func setupCartRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// TODO: Fix service interface compatibility issues
	// The cart module expects different interfaces than what the services provide
	// Temporarily commented out to allow compilation
	
	// Initialize required service dependencies
	// productModule := product.NewModule(cfg.DB)
	// discountModule := discount.NewModule(cfg.DB)
	// taxModule := tax.NewModule(cfg.DB)
	// shippingModule := shipping.NewModule(cfg.DB)
	
	// Initialize cart module with dependencies
	// cartModule := cart.NewModule(cfg.DB, productModule.Service, discountModule.GetService(), taxModule.GetService(), shippingModule.GetService())
	
	// Register cart routes
	// cartModule.RegisterRoutes(v1)
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
