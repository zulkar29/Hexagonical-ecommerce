package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ecommerce-saas/internal/shared/config"
	"ecommerce-saas/internal/shared/database"
	"ecommerce-saas/internal/shared/middleware"
	"ecommerce-saas/internal/tenant"
	"ecommerce-saas/internal/product"
	"ecommerce-saas/internal/order"
	"ecommerce-saas/internal/user"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	config *config.Config
	db     *gorm.DB
}

// New creates a new server instance
func New(cfg *config.Config, db *gorm.DB) *Server {
	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.New()

	server := &Server{
		router: router,
		config: cfg,
		db:     db,
	}

	// Setup middleware
	server.setupMiddleware()

	// Setup routes
	server.setupRoutes()

	return server
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.router.Use(gin.Recovery())

	// Logger middleware
	if s.config.App.Debug {
		s.router.Use(middleware.LoggingMiddleware())
	}

	// CORS middleware
	s.router.Use(middleware.CORSMiddleware(
		s.config.App.CORS.AllowedOrigins,
		s.config.App.CORS.AllowedMethods,
		s.config.App.CORS.AllowedHeaders,
		s.config.App.CORS.AllowCredentials,
	))

	// Rate limiting middleware
	s.router.Use(middleware.RateLimitMiddleware())

	// Request ID middleware
	s.router.Use(middleware.RequestIDMiddleware())

	// TODO: Add more middleware as needed
	// - Authentication middleware
	// - Tenant context middleware
	// - Metrics middleware
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthHandler)
	s.router.GET("/ready", s.readyHandler)

	// API routes
	api := s.router.Group("/api/v1")
	
	// Setup module routes
	s.setupTenantRoutes(api)
	s.setupProductRoutes(api)
	s.setupOrderRoutes(api)
	s.setupUserRoutes(api)
	
	// TODO: Add more route groups
	// - Analytics routes
	// - Payment routes
	// - Notification routes
}

// setupTenantRoutes configures tenant routes
func (s *Server) setupTenantRoutes(api *gin.RouterGroup) {
	// Initialize tenant module
	tenantRepo := tenant.NewRepository(s.db)
	tenantService := tenant.NewService(tenantRepo)
	tenantHandler := tenant.NewHandler(tenantService)

	// Tenant routes
	tenants := api.Group("/tenants")
	{
		tenants.POST("", tenantHandler.CreateTenant)
		tenants.GET("", tenantHandler.ListTenants)
		tenants.GET("/:id", tenantHandler.GetTenant)
		tenants.PUT("/:id", tenantHandler.UpdateTenant)
		tenants.PUT("/:id/plan", tenantHandler.UpdatePlan)
		tenants.POST("/:id/activate", tenantHandler.ActivateTenant)
		tenants.POST("/:id/deactivate", tenantHandler.DeactivateTenant)
		tenants.GET("/:id/stats", tenantHandler.GetTenantStats)
		tenants.GET("/subdomain/:subdomain", tenantHandler.GetTenantBySubdomain)
		tenants.GET("/check-subdomain/:subdomain", tenantHandler.CheckSubdomainAvailability)
	}
}

// setupProductRoutes configures product routes
func (s *Server) setupProductRoutes(api *gin.RouterGroup) {
	// Initialize product module
	productRepo := product.NewRepository(s.db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	// Product routes (tenant-scoped)
	products := api.Group("/products")
	products.Use(s.tenantMiddleware()) // Extract tenant from subdomain/header
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("", productHandler.ListProducts)
		products.GET("/stats", productHandler.GetProductStats)
		products.PATCH("/bulk", productHandler.BulkUpdateProducts)
		products.GET("/:id", productHandler.GetProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
		products.PATCH("/:id/inventory", productHandler.UpdateInventory)
		products.GET("/slug/:slug", productHandler.GetProductBySlug)
	}

	// Category routes (tenant-scoped)
	categories := api.Group("/categories")
	categories.Use(s.tenantMiddleware())
	{
		categories.POST("", productHandler.CreateCategory)
		categories.GET("", productHandler.ListCategories)
		categories.GET("/:id", productHandler.GetCategory)
	}
}

// setupOrderRoutes configures order routes
func (s *Server) setupOrderRoutes(api *gin.RouterGroup) {
	// TODO: Initialize order module
	// orderRepo := order.NewRepository(s.db)
	// orderService := order.NewService(orderRepo)
	// orderHandler := order.NewHandler(orderService)

	// Orders routes (tenant-scoped)
	orders := api.Group("/orders")
	orders.Use(s.tenantMiddleware())
	{
		// TODO: Add order routes
		orders.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Orders endpoint - TODO"})
		})
	}
}

// setupUserRoutes configures user routes
func (s *Server) setupUserRoutes(api *gin.RouterGroup) {
	// TODO: Initialize user module
	// userRepo := user.NewRepository(s.db)
	// userService := user.NewService(userRepo)
	// userHandler := user.NewHandler(userService)

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		// TODO: Add auth routes
		auth.POST("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - TODO"})
		})
		auth.POST("/register", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Register endpoint - TODO"})
		})
	}

	// User routes (authenticated)
	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware()) // TODO: Implement auth middleware
	{
		// TODO: Add user routes
		users.GET("/profile", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Profile endpoint - TODO"})
		})
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	
	server := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Close database connection
	database.Close()

	log.Println("Server exited")
	return nil
}

// GetRouter returns the gin router (for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// Health handlers
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().UTC(),
	})
}

func (s *Server) readyHandler(c *gin.Context) {
	// Check database connection
	if err := database.Health(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"error":  "database not available",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"time":   time.Now().UTC(),
	})
}

// TODO: Additional middleware to implement:
// - Security headers middleware
// - Metrics collection middleware
// - Enhanced logging middleware
