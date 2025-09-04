package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configures CORS settings
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	if len(allowedOrigins) == 0 {
		// Default allowed origins for development
		allowedOrigins = []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"https://*.vercel.app",
		}
	}

	return cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
			"X-Tenant-ID",
		},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// LoggingMiddleware logs all requests with structured format
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Get request ID from context
		requestID := ""
		if param.Keys != nil {
			if id, exists := param.Keys["request_id"]; exists {
				requestID = fmt.Sprintf("%v", id)
			}
		}

		// Structured log format
		return fmt.Sprintf(
			"[%s] %s %s %s %d %s %s %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			requestID,
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.ErrorMessage,
		)
	})
}

// CustomLoggerMiddleware provides more detailed logging
func CustomLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request ID
		requestID, _ := c.Get("request_id")

		// Log request details
		log.Printf(
			"[REQUEST] ID=%v | %s %s%s | Status=%d | Latency=%v | IP=%s | UserAgent=%s",
			requestID,
			c.Request.Method,
			path,
			raw,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}

// SecurityMiddleware adds comprehensive security headers
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:")
		
		// Strict Transport Security (HTTPS only)
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		
		// Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		// Remove server information
		c.Header("Server", "")
		
		c.Next()
	}
}

// ErrorMiddleware handles errors globally
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()
		
		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			requestID, _ := c.Get("request_id")
			
			// Log the error
			log.Printf(
				"[ERROR] RequestID=%v | Path=%s | Error=%v | Type=%v",
				requestID,
				c.Request.URL.Path,
				err.Error(),
				err.Type,
			)
			
			// Return appropriate error response based on error type
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{
					"error":      "Invalid request format",
					"message":    "Please check your request data",
					"request_id": requestID,
				})
			case gin.ErrorTypePublic:
				c.JSON(http.StatusBadRequest, gin.H{
					"error":      err.Error(),
					"request_id": requestID,
				})
			default:
				// Internal server error - don't expose internal details
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":      "Internal server error",
					"message":    "Something went wrong. Please try again later.",
					"request_id": requestID,
				})
			}
			
			// Abort to prevent further processing
			c.Abort()
		}
	}
}

// RecoveryMiddleware handles panics and converts them to 500 errors
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID, _ := c.Get("request_id")
		
		// Log the panic
		log.Printf(
			"[PANIC] RequestID=%v | Path=%s | Panic=%v",
			requestID,
			c.Request.URL.Path,
			recovered,
		)
		
		// Return 500 error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Internal server error",
			"message":    "An unexpected error occurred. Please try again later.",
			"request_id": requestID,
		})
	})
}
