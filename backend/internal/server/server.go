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
	"ecommerce-saas/internal/shared/routes"
	"ecommerce-saas/internal/shared/utils"
)

// Server represents the HTTP server
type Server struct {
	router     *gin.Engine
	config     *config.Config
	db         *gorm.DB
	jwtManager *utils.JWTManager
}

// New creates a new server instance
func New(cfg *config.Config, db *gorm.DB, jwtManager *utils.JWTManager) *Server {
	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.New()

	server := &Server{
		router:     router,
		config:     cfg,
		db:         db,
		jwtManager: jwtManager,
	}

	// Setup routes with dependencies
	routeConfig := &routes.RouteConfig{
		DB:         db,
		Config:     cfg,
		JWTManager: jwtManager,
	}
	routes.SetupRoutes(router, routeConfig)

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	
	server := &http.Server{
		Addr:           address,
		Handler:        s.router,
		ReadTimeout:    s.config.Server.ReadTimeout,
		WriteTimeout:   s.config.Server.WriteTimeout,
		IdleTimeout:    s.config.Server.IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on %s", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Server shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("✅ Server exited gracefully")
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
