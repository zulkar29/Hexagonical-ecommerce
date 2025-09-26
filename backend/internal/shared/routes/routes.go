package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	"ecommerce-saas/internal/referral"
	"ecommerce-saas/internal/returns"
	"ecommerce-saas/internal/reviews"
	"ecommerce-saas/internal/search"

	"ecommerce-saas/internal/settings"
	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/shared/email"
	"ecommerce-saas/internal/shared/middleware"
	"ecommerce-saas/internal/shared/utils"
	"ecommerce-saas/internal/shipping"
	"ecommerce-saas/internal/support"
	"ecommerce-saas/internal/tenant"
	"ecommerce-saas/internal/user"
	"ecommerce-saas/internal/webhook"
	"ecommerce-saas/internal/wishlist"
)

// RouteConfig holds dependencies for route setup
type RouteConfig struct {
	DB         *gorm.DB
	Config     *config.Config
	JWTManager *utils.JWTManager
}

// ServiceContainer holds shared service instances to avoid duplications
type ServiceContainer struct {
	ProductRepo         product.Repository
	DiscountRepo        discount.Repository
	ContactRepo         contact.Repository
	AnalyticsRepo       analytics.Repository
	ShippingRepo        *shipping.Repository
	PaymentRepo         payment.Repository
	ReferralRepo        referral.Repository
	SettingsRepo        settings.Repository
	OrderRepo           order.Repository
	NotificationRepo    notification.Repository
	ProductService      *product.Service
	DiscountService     discount.Service
	ContactService      contact.Service
	AnalyticsService    analytics.Service
	ShippingService     *shipping.Service
	PaymentService      payment.Service
	ReferralService     referral.Service
	OrderService        *order.Service
	NotificationService notification.Service
}

// OrderServiceAdapter adapts order.Service to cart.OrderService interface
type OrderServiceAdapter struct {
	orderService *order.Service
}

func (a *OrderServiceAdapter) CreateOrderFromCart(ctx context.Context, tenantID, cartID uuid.UUID) (*cart.OrderFromCartResult, error) {
	result, err := a.orderService.CreateOrderFromCart(ctx, tenantID, cartID)
	if err != nil {
		return nil, err
	}
	
	// Convert order.OrderFromCartResult to cart.OrderFromCartResult
	return &cart.OrderFromCartResult{
		OrderID:     result.OrderID,
		OrderNumber: result.OrderNumber,
		Total:       result.Total,
		ItemCount:   result.ItemCount,
	}, nil
}

// SettingsRepositoryWrapper wraps settings.Repository to provide adapter functionality
type SettingsRepositoryWrapper struct {
	repo settings.Repository
	adapter *utils.SettingsRepositoryAdapter
}

func NewSettingsRepositoryWrapper(repo settings.Repository) *SettingsRepositoryWrapper {
	// Create an adapter that converts settings.Repository to the expected interface
	adapterRepo := &settingsRepoAdapter{repo: repo}
	adapter := utils.NewSettingsRepositoryAdapter(adapterRepo)
	return &SettingsRepositoryWrapper{
		repo: repo,
		adapter: adapter,
	}
}

// settingsRepoAdapter adapts settings.Repository to match the interface expected by utils.NewSettingsRepositoryAdapter
type settingsRepoAdapter struct {
	repo settings.Repository
}

func (s *settingsRepoAdapter) GetSetting(tenantID uuid.UUID, section, key string) (interface{}, error) {
	return s.repo.GetSetting(tenantID, section, key)
}

// Implement settings.Repository interface by delegating to the wrapped repo
func (w *SettingsRepositoryWrapper) GetSettings(tenantID uuid.UUID, section string) ([]settings.Setting, error) {
	return w.repo.GetSettings(tenantID, section)
}

func (w *SettingsRepositoryWrapper) GetSetting(tenantID uuid.UUID, section, key string) (*settings.Setting, error) {
	return w.repo.GetSetting(tenantID, section, key)
}

func (w *SettingsRepositoryWrapper) CreateSetting(setting *settings.Setting) error {
	return w.repo.CreateSetting(setting)
}

func (w *SettingsRepositoryWrapper) UpdateSetting(setting *settings.Setting) error {
	return w.repo.UpdateSetting(setting)
}

func (w *SettingsRepositoryWrapper) DeleteSetting(tenantID uuid.UUID, section, key string) error {
	return w.repo.DeleteSetting(tenantID, section, key)
}

func (w *SettingsRepositoryWrapper) GetPublicSettings(tenantID uuid.UUID) ([]settings.Setting, error) {
	return w.repo.GetPublicSettings(tenantID)
}

// GetAdapter returns the utils adapter for use with discount service
func (w *SettingsRepositoryWrapper) GetAdapter() *utils.SettingsRepositoryAdapter {
	return w.adapter
}

// NewServiceContainer creates and initializes shared services
func NewServiceContainer(cfg *RouteConfig) *ServiceContainer {
	// Initialize repositories
	productRepo := product.NewRepository(cfg.DB)
	settingsRepo := settings.NewRepository(cfg.DB)
	// Create wrapper for settings repository that provides both interfaces
	settingsWrapper := NewSettingsRepositoryWrapper(settingsRepo)
	discountRepo := discount.NewRepository(cfg.DB, settingsWrapper.GetAdapter())
	contactRepo := contact.NewRepository(cfg.DB)
	analyticsRepo := analytics.NewRepository(cfg.DB)
	shippingRepo := shipping.NewRepository(cfg.DB)
	paymentRepo := payment.NewRepository(cfg.DB)
	referralRepo := referral.NewGormRepository(cfg.DB)
	orderRepo := order.NewRepository(cfg.DB)
	notificationRepo := notification.NewRepository(cfg.DB)

	// Initialize services
	productService := product.NewService(productRepo, cfg.DB)
	discountService := discount.NewService(discountRepo)
	contactService := contact.NewService(contactRepo)
	analyticsService := analytics.NewService(analyticsRepo)
	shippingService := shipping.NewService(shippingRepo)
	paymentService := payment.NewService(paymentRepo, cfg.Config)
	referralService := referral.NewService(referralRepo)
	emailService := email.NewService(cfg.Config)
	notificationService := notification.NewService(notificationRepo, emailService)
	orderService := order.NewService(orderRepo, cfg.DB, productService, discountService, paymentService, notificationService)

	return &ServiceContainer{
		ProductRepo:      productRepo,
		DiscountRepo:     discountRepo,
		ContactRepo:      contactRepo,
		AnalyticsRepo:    analyticsRepo,
		ShippingRepo:     shippingRepo,
		PaymentRepo:      paymentRepo,
		ReferralRepo:     referralRepo,
		SettingsRepo:        settingsWrapper,
		OrderRepo:           orderRepo,
		NotificationRepo:    notificationRepo,
		ProductService:      productService,
		DiscountService:  discountService,
		ContactService:   contactService,
		AnalyticsService: analyticsService,
		ShippingService:  shippingService,
		PaymentService:   paymentService,
		ReferralService:     referralService,
		OrderService:        orderService,
		NotificationService: notificationService,
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
	tenantIsolation := middleware.NewTenantIsolationMiddleware(cfg.DB, cfg.Config.App.Domain)
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
		setupReferralRoutes(protected, cfg, services)
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
	publicTenantIsolation := middleware.NewTenantIsolationMiddleware(cfg.DB, cfg.Config.App.Domain)
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
	tenantModule := tenant.NewModule(cfg.DB, cfg.Config)

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

		// Public category browsing - already handled in storefront section above
		// setupPublicCategoryRoutes(public, cfg) // REMOVED: duplicate route registration

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
	notificationModule := notification.NewModule(cfg.DB, cfg.Config)
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
	notificationModule := notification.NewModule(cfg.DB, cfg.Config)

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
	billingModule := billing.NewModule(cfg.DB, services.PaymentService, services.ContactService, services.AnalyticsService, services.ReferralService)

	// Register billing routes
	billingModule.RegisterRoutes(v1)
}

// Setup referral routes
func setupReferralRoutes(v1 *gin.RouterGroup, cfg *RouteConfig, services *ServiceContainer) {
	// Initialize referral module using shared services
	referralModule := referral.NewModule(cfg.DB)

	// Register referral routes
	referralModule.RegisterRoutes(v1)
}

// Setup cart routes
func setupCartRoutes(v1 *gin.RouterGroup, cfg *RouteConfig, services *ServiceContainer) {
	// Create adapter for order service to match cart.OrderService interface
	orderServiceAdapter := &OrderServiceAdapter{orderService: services.OrderService}
	
	// Initialize cart module using shared services
	cartModule := cart.NewModule(cfg.DB, services.ProductService, services.DiscountService, services.ShippingService, orderServiceAdapter)

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
	componentsModule.RegisterRoutes(v1)
}

// Setup security routes
func setupSecurityRoutes(v1 *gin.RouterGroup, cfg *RouteConfig) {
	// For now, we'll skip standalone security routes since they require
	// complex adapter patterns. Security functionality is accessible
	// through the user module endpoints.
	//
	// If you need standalone security endpoints, you can enable them by:
	// 1. Creating a user repository adapter
	// 2. Initializing: securityModule := security.NewModule(cfg.DB, userAdapter)
	// 3. Registering: securityModule.RegisterRoutes(v1)
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
