package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	tenantpkg "ecommerce-saas/internal/tenant"
)

// TenantContextKey is the key used to store tenant in context
type TenantContextKey string

const (
	TenantKey   TenantContextKey = "tenant"
	TenantIDKey TenantContextKey = "tenant_id"
)

// TenantMiddleware handles tenant resolution based on domain/subdomain
type TenantMiddleware struct {
	tenantRepo tenantpkg.Repository
	baseDomain string // e.g., "esass.com"
}

// NewTenantMiddleware creates a new tenant middleware
func NewTenantMiddleware(tenantRepo tenantpkg.Repository, baseDomain string) *TenantMiddleware {
	return &TenantMiddleware{
		tenantRepo: tenantRepo,
		baseDomain: baseDomain,
	}
}

// ResolveTenant middleware resolves tenant from request domain
func (tm *TenantMiddleware) ResolveTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host

		// Remove port if present
		if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
			host = host[:colonIndex]
		}

		tenant, err := tm.resolveTenantFromHost(host)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Tenant not found",
				"message": "The domain you're trying to access is not associated with any active tenant",
			})
			c.Abort()
			return
		}

		if !tenant.IsActive() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Tenant inactive",
				"message": "This tenant is currently inactive",
			})
			c.Abort()
			return
		}

		// Set tenant in context
		c.Set(string(TenantKey), tenant)
		c.Set(string(TenantIDKey), tenant.ID)

		// Set tenant in request context for downstream services
		ctx := context.WithValue(c.Request.Context(), TenantKey, tenant)
		ctx = context.WithValue(ctx, TenantIDKey, tenant.ID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// resolveTenantFromHost resolves tenant from the host header
func (tm *TenantMiddleware) resolveTenantFromHost(host string) (*tenantpkg.Tenant, error) {
	// First, try to find by custom domain
	if tenantByDomain, err := tm.tenantRepo.FindByCustomDomain(host); err == nil {
		return tenantByDomain, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// If not found by custom domain, check if it's a subdomain
	if strings.HasSuffix(host, "."+tm.baseDomain) {
		// Extract subdomain
		subdomain := strings.TrimSuffix(host, "."+tm.baseDomain)

		// Handle nested subdomains (take the first part)
		if dotIndex := strings.Index(subdomain, "."); dotIndex != -1 {
			subdomain = subdomain[:dotIndex]
		}

		// Find tenant by subdomain
		return tm.tenantRepo.FindBySubdomain(subdomain)
	}

	// If neither custom domain nor subdomain, return error
	return nil, gorm.ErrRecordNotFound
}

// RequireTenant middleware ensures a tenant is present in context
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant, exists := c.Get(string(TenantKey))
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Tenant required",
				"message": "This endpoint requires tenant context",
			})
			c.Abort()
			return
		}

		// Validate tenant type
		if _, ok := tenant.(*tenantpkg.Tenant); !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid tenant context",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetTenantFromContext extracts tenant from gin context
func GetTenantFromContext(c *gin.Context) (*tenantpkg.Tenant, bool) {
	tenant, exists := c.Get(string(TenantKey))
	if !exists {
		return nil, false
	}

	tenantObj, ok := tenant.(*tenantpkg.Tenant)
	return tenantObj, ok
}

// GetTenantIDFromContext extracts tenant ID from gin context
func GetTenantIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	tenantID, exists := c.Get(string(TenantIDKey))
	if !exists {
		return uuid.Nil, false
	}

	tenantUUID, ok := tenantID.(uuid.UUID)
	return tenantUUID, ok
}

// GetTenantFromRequestContext extracts tenant from request context
func GetTenantFromRequestContext(ctx context.Context) (*tenantpkg.Tenant, bool) {
	tenant := ctx.Value(TenantKey)
	if tenant == nil {
		return nil, false
	}

	tenantObj, ok := tenant.(*tenantpkg.Tenant)
	return tenantObj, ok
}

// GetTenantIDFromRequestContext extracts tenant ID from request context
func GetTenantIDFromRequestContext(ctx context.Context) (uuid.UUID, bool) {
	tenantID := ctx.Value(TenantIDKey)
	if tenantID == nil {
		return uuid.Nil, false
	}

	tenantUUID, ok := tenantID.(uuid.UUID)
	return tenantUUID, ok
}
