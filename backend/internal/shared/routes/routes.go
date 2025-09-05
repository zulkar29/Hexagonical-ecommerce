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
	"ecommerce-saas/internal/payment"
	"ecommerce-saas/internal/notification"
	"ecommerce-saas/internal/address"
	"ecommerce-saas/internal/analytics"
	// "ecommerce-saas/internal/billing"
	// "ecommerce-saas/internal/cart"
	// "ecommerce-saas/internal/category"
	"ecommerce-saas/internal/contact"
	"ecommerce-saas/internal/content"
	"ecommerce-saas/internal/discount"
	"ecommerce-saas/internal/finance"
	"ecommerce-saas/internal/marketing"
	// "ecommerce-saas/internal/observability"
	"ecommerce-saas/internal/returns"
	"ecommerce-saas/internal/reviews"
	"ecommerce-saas/internal/shipping"
	"ecommerce-saas/internal/support"
	// "ecommerce-saas/internal/webhook"
	// "ecommerce-saas/internal/wishlist"
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
		setupAnalyticsRoutes(protected, cfg)
		setupBillingRoutes(protected, cfg)
		setupCartRoutes(protected, cfg)
		setupCategoryRoutes(protected, cfg)
		setupContactRoutes(protected, cfg)
		setupContentRoutes(protected, cfg)
		setupDiscountRoutes(protected, cfg)
		setupMarketingRoutes(protected, cfg)
		setupObservabilityRoutes(protected, cfg)
		setupReviewsRoutes(protected, cfg)
		setupShippingRoutes(protected, cfg)
		setupSupportRoutes(protected, cfg)
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
		
		// TODO: Public order tracking (no auth required) - requires order module
		// public.GET("/orders/track/:number", orderModule.Handler.TrackOrder)
		// public.GET("/orders/number/:number", orderModule.Handler.GetOrderByNumber)
	}
}

func setupOrderRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// TODO: Initialize order module - requires service dependencies
	// orderModule := order.NewModule(cfg.DB, productService, discountService, paymentService, inventoryService, notificationService)
	
	// TODO: Register order routes
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
	// TODO: Initialize finance module - NewModule function not implemented yet
	// financeModule := finance.NewModule(cfg.DB)
	
	// Initialize finance handler directly
	financeRepo := finance.NewRepository(cfg.DB)
	financeService := finance.NewService(financeRepo)
	financeHandler := finance.NewHandler(financeService)
	
	// Register finance routes
	financeHandler.RegisterRoutes(v1)
}

func setupReturnsRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// TODO: Initialize returns module - NewModule function not implemented yet
	// returnsModule := returns.NewModule(cfg.DB)
	
	// Initialize returns handler directly
	returnsRepo := returns.NewRepository(cfg.DB)
	returnsService := returns.NewService(returnsRepo)
	returnsHandler := returns.NewHandler(returnsService)
	
	// Register returns routes
	returnsHandler.RegisterRoutes(v1)
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

// Setup billing routes
func setupBillingRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// billingRepo := billing.NewBillingRepository(cfg.DB)
	// TODO: Initialize payment provider, email service, and analytics service
	// billingService := billing.NewBillingService(billingRepo, paymentProvider, emailService, analyticsService)
	// billingHandler := billing.NewBillingHandler(billingService)
	// billingHandler.RegisterRoutes(v1)
}

// Setup cart routes
func setupCartRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// cartRepo := cart.NewRepository(cfg.DB)
	// TODO: Initialize product, discount, tax, and shipping services
	// cartService := cart.NewService(cartRepo, productService, discountService, taxService, shippingService)
	// cartHandler := cart.NewHandler(cartService)
	// cartHandler.RegisterRoutes(v1)
}

// Setup category routes
func setupCategoryRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// TODO: Implement category module
	// categoryRepo := category.NewRepository(cfg.DB)
	// categoryService := category.NewService(categoryRepo)
	// categoryHandler := category.NewHandler(categoryService)
	// categoryHandler.RegisterRoutes(v1)
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

// Setup observability routes
func setupObservabilityRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// TODO: Implement observability module
	// observabilityRepo := observability.NewRepository(cfg.DB)
	// observabilityService := observability.NewService(observabilityRepo)
	// observabilityHandler := observability.NewHandler(observabilityService)
	// observabilityHandler.RegisterRoutes(v1)
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

// Setup webhook routes
func setupWebhookRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// webhookRepo := webhook.NewRepository(cfg.DB)
	// TODO: Initialize signing key for webhook service
	// signingKey := []byte("your-webhook-signing-key")
	// webhookService := webhook.NewService(webhookRepo, signingKey)
	// webhookHandler := webhook.NewHandler(webhookService)
	// webhookHandler.RegisterRoutes(v1)
}

// Setup wishlist routes
func setupWishlistRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// TODO: Implement wishlist module
	// wishlistRepo := wishlist.NewRepository(cfg.DB)
	// wishlistService := wishlist.NewService(wishlistRepo)
	// wishlistHandler := wishlist.NewHandler(wishlistService)
	// wishlistHandler.RegisterRoutes(v1)
}
