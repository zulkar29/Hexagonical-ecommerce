package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ecommerce-saas/internal/tenant"
)

// TenantIsolationMiddleware provides comprehensive tenant isolation and validation
type TenantIsolationMiddleware struct {
	db         *gorm.DB
	tenantRepo tenant.Repository
	baseDomain string
}

// NewTenantIsolationMiddleware creates a new tenant isolation middleware
func NewTenantIsolationMiddleware(db *gorm.DB, baseDomain string) *TenantIsolationMiddleware {
	tenantRepo := tenant.NewRepository(db)
	return &TenantIsolationMiddleware{
		db:         db,
		tenantRepo: tenantRepo,
		baseDomain: baseDomain,
	}
}

// ResolveTenantWithIsolation resolves tenant and ensures proper data isolation
func (tim *TenantIsolationMiddleware) ResolveTenantWithIsolation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get host from request
		host := c.Request.Host
		if host == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Host header required",
				"code":  "MISSING_HOST",
			})
			c.Abort()
			return
		}

		// Remove port if present
		if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
			host = host[:colonIndex]
		}

		// Resolve tenant from host
		tenantEntity, err := tim.resolveTenantFromHost(host)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"error":   "Tenant not found",
					"message": "The domain you're trying to access is not associated with any active tenant",
					"code":    "TENANT_NOT_FOUND",
					"host":    host,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to resolve tenant",
					"message": "An error occurred while resolving tenant information",
					"code":    "TENANT_RESOLUTION_ERROR",
				})
			}
			c.Abort()
			return
		}

		// Validate tenant status
		if !tenantEntity.IsActive() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Tenant inactive",
				"message": "This tenant is currently inactive or suspended",
				"code":    "TENANT_INACTIVE",
			})
			c.Abort()
			return
		}

		// Set tenant context for data isolation
		tim.setTenantContext(c, tenantEntity, host)

		// Configure database session with tenant isolation
		tim.configureTenantIsolatedDB(c, tenantEntity.ID)

		c.Next()
	}
}

// ValidateTenantAccess validates that user has access to the resolved tenant
func (tim *TenantIsolationMiddleware) ValidateTenantAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant from context (should be set by ResolveTenantWithIsolation)
		tenantID, exists := c.Get("tenant_id")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Tenant context missing",
				"code":  "MISSING_TENANT_CONTEXT",
			})
			c.Abort()
			return
		}

		// Get user's tenant from JWT claims (if authenticated)
		userTenantID, userExists := c.Get("user_tenant_id")
		if userExists {
			// Validate that user belongs to the resolved tenant
			if userTenantID != tenantID {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "Tenant access denied",
					"message": "User does not have access to this tenant",
					"code":    "TENANT_ACCESS_DENIED",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// resolveTenantFromHost resolves tenant from the host header
func (tim *TenantIsolationMiddleware) resolveTenantFromHost(host string) (*tenant.Tenant, error) {
	// First, try to find by custom domain
	if tenantByDomain, err := tim.tenantRepo.FindByCustomDomain(host); err == nil {
		return tenantByDomain, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// If not found by custom domain, check if it's a subdomain
	if strings.HasSuffix(host, "."+tim.baseDomain) {
		// Extract subdomain
		subdomain := strings.TrimSuffix(host, "."+tim.baseDomain)

		// Handle nested subdomains (take the first part)
		if dotIndex := strings.Index(subdomain, "."); dotIndex != -1 {
			subdomain = subdomain[:dotIndex]
		}

		// Skip reserved subdomains
		if tim.isReservedSubdomain(subdomain) {
			return nil, gorm.ErrRecordNotFound
		}

		// Find tenant by subdomain
		return tim.tenantRepo.FindBySubdomain(subdomain)
	}

	// If neither custom domain nor subdomain, return error
	return nil, gorm.ErrRecordNotFound
}

// isReservedSubdomain checks if subdomain is reserved
func (tim *TenantIsolationMiddleware) isReservedSubdomain(subdomain string) bool {
	reserved := []string{"www", "api", "admin", "app", "mail", "ftp", "blog", "support", "help", "docs"}
	for _, r := range reserved {
		if subdomain == r {
			return true
		}
	}
	return false
}

// setTenantContext sets tenant information in gin context
func (tim *TenantIsolationMiddleware) setTenantContext(c *gin.Context, tenantEntity *tenant.Tenant, host string) {
	// Set tenant in gin context
	c.Set("tenant", tenantEntity)
	c.Set("tenant_id", tenantEntity.ID)
	c.Set("tenant_subdomain", tenantEntity.Subdomain)
	c.Set("tenant_domain", host)
	c.Set("tenant_name", tenantEntity.Name)
	c.Set("tenant_status", tenantEntity.Status)

	// Set tenant in request context for downstream services
	ctx := context.WithValue(c.Request.Context(), "tenant", tenantEntity)
	ctx = context.WithValue(ctx, "tenant_id", tenantEntity.ID)
	c.Request = c.Request.WithContext(ctx)
}

// configureTenantIsolatedDB configures database session with tenant isolation
func (tim *TenantIsolationMiddleware) configureTenantIsolatedDB(c *gin.Context, tenantID uuid.UUID) {
	// Create a new DB session with tenant isolation
	tenantDB := tim.db.Session(&gorm.Session{})

	// Add global scope for tenant isolation
	tenantDB = tenantDB.Scopes(func(db *gorm.DB) *gorm.DB {
		// Only apply tenant filtering to tables that have tenant_id column
		if db.Statement.Schema != nil {
			for _, field := range db.Statement.Schema.Fields {
				if field.DBName == "tenant_id" {
					return db.Where("tenant_id = ?", tenantID)
				}
			}
		}
		return db
	})

	// Set the tenant-isolated DB in context
	c.Set("tenant_db", tenantDB)
}

// GetTenantFromContext extracts tenant from gin context
func GetTenantFromContext(c *gin.Context) (*tenant.Tenant, bool) {
	tenantValue, exists := c.Get("tenant")
	if !exists {
		return nil, false
	}

	tenantObj, ok := tenantValue.(*tenant.Tenant)
	return tenantObj, ok
}

// GetTenantIDFromContext extracts tenant ID from gin context
func GetTenantIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		return uuid.Nil, false
	}

	tenantUUID, ok := tenantID.(uuid.UUID)
	return tenantUUID, ok
}

// GetTenantDBFromContext extracts tenant-isolated DB from gin context
func GetTenantDBFromContext(c *gin.Context) (*gorm.DB, bool) {
	tenantDB, exists := c.Get("tenant_db")
	if !exists {
		return nil, false
	}

	db, ok := tenantDB.(*gorm.DB)
	return db, ok
}

// RequireTenant middleware ensures a tenant is present in context
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := GetTenantFromContext(c); !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Tenant required",
				"message": "This endpoint requires tenant context",
				"code":    "TENANT_REQUIRED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
