package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
	
	"ecommerce-saas/internal/shared/utils"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer token format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Validate token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("token_id", claims.TokenID)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT token if present but doesn't require it
func OptionalAuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if claims, err := jwtManager.ValidateToken(token); err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("tenant_id", claims.TenantID)
				c.Set("user_email", claims.Email)
				c.Set("user_role", claims.Role)
				c.Set("token_id", claims.TokenID)
			}
		}

		c.Next()
	}
}



// TenantMiddleware ensures the request has valid tenant context
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.GetHeader("Host")
		if host == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Host header required"})
			c.Abort()
			return
		}

		// Parse subdomain from host
		parts := strings.Split(host, ".")
		if len(parts) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host format"})
			c.Abort()
			return
		}

		subdomain := parts[0]
		if subdomain == "" || subdomain == "www" || subdomain == "api" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant subdomain"})
			c.Abort()
			return
		}

		// Set tenant context from subdomain
		c.Set("tenant_subdomain", subdomain)
		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role format"})
			c.Abort()
			return
		}

		// Check if user has required role or higher privileges
		if !hasRequiredRole(role, requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasRequiredRole checks if user role has required permissions
func hasRequiredRole(userRole, requiredRole string) bool {
	// Define role hierarchy (higher roles include lower role permissions)
	roleHierarchy := map[string]int{
		"customer":    1,
		"staff":       2,
		"manager":     3,
		"admin":       4,
		"super_admin": 5,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}

// RateLimitMiddleware implements basic rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	// Simple in-memory rate limiting
	clients := make(map[string][]time.Time)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()
		
		// Clean old entries (older than 1 minute)
		if requests, exists := clients[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < time.Minute {
					validRequests = append(validRequests, reqTime)
				}
			}
			clients[clientIP] = validRequests
		}
		
		// Check rate limit (100 requests per minute)
		if len(clients[clientIP]) >= 100 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		
		// Add current request
		clients[clientIP] = append(clients[clientIP], now)
		
		c.Next()
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}
