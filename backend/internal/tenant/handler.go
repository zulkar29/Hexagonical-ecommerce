package tenant

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for tenant operations
type Handler struct {
	service *Service
}

// NewHandler creates a new tenant handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers all tenant routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	tenants := router.Group("/tenants")
	{
		// Basic CRUD operations
		tenants.POST("", h.CreateTenant)
		tenants.GET("", h.ListTenants)
		tenants.GET("/:id", h.GetTenant)
		tenants.PUT("/:id", h.UpdateTenant)
		
		// Plan management
		tenants.PUT("/:id/plan", h.UpdatePlan)
		tenants.GET("/:id/upgrade-options", h.GetUpgradeOptions)
		
		// Status management
		tenants.POST("/:id/activate", h.ActivateTenant)
		tenants.POST("/:id/deactivate", h.DeactivateTenant)
		tenants.POST("/:id/suspend", h.SuspendTenant)
		
		// Information and statistics
		tenants.GET("/:id/stats", h.GetTenantStats)
		tenants.GET("/subdomain/:subdomain", h.GetTenantBySubdomain)
		
		// Utility endpoints
		tenants.GET("/check-subdomain/:subdomain", h.CheckSubdomainAvailability)
		tenants.POST("/:id/validate-domain", h.ValidateCustomDomain)
	}
}

// CreateTenant handles POST /api/tenants
func (h *Handler) CreateTenant(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	tenant, err := h.service.CreateTenant(req)
	if err != nil {
		// TODO: Implement proper error handling with custom error types
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tenant created successfully",
		"data":    tenant,
	})
}

// GetTenant handles GET /api/tenants/:id
func (h *Handler) GetTenant(c *gin.Context) {
	tenantID := c.Param("id")
	
	tenant, err := h.service.GetTenant(tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tenant,
	})
}

// GetTenantBySubdomain handles GET /api/tenants/subdomain/:subdomain
func (h *Handler) GetTenantBySubdomain(c *gin.Context) {
	subdomain := c.Param("subdomain")
	
	tenant, err := h.service.GetTenantBySubdomain(subdomain)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tenant,
	})
}

// UpdateTenant handles PUT /api/tenants/:id
func (h *Handler) UpdateTenant(c *gin.Context) {
	tenantID := c.Param("id")
	
	var req UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	tenant, err := h.service.UpdateTenant(tenantID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant updated successfully",
		"data":    tenant,
	})
}

// UpdatePlan handles PUT /api/tenants/:id/plan
func (h *Handler) UpdatePlan(c *gin.Context) {
	tenantID := c.Param("id")
	
	var req UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	tenant, err := h.service.UpdatePlan(tenantID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plan updated successfully",
		"data":    tenant,
	})
}

// ListTenants handles GET /api/tenants
func (h *Handler) ListTenants(c *gin.Context) {
	// Parse pagination parameters
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	
	// Validate limit
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}

	tenants, total, err := h.service.ListTenants(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"tenants": tenants,
			"total":   total,
			"offset":  offset,
			"limit":   limit,
		},
	})
}

// DeactivateTenant handles POST /api/tenants/:id/deactivate
func (h *Handler) DeactivateTenant(c *gin.Context) {
	tenantID := c.Param("id")
	
	err := h.service.DeactivateTenant(tenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant deactivated successfully",
	})
}

// ActivateTenant handles POST /api/tenants/:id/activate
func (h *Handler) ActivateTenant(c *gin.Context) {
	tenantID := c.Param("id")
	
	err := h.service.ActivateTenant(tenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant activated successfully",
	})
}

// GetTenantStats handles GET /api/tenants/:id/stats
func (h *Handler) GetTenantStats(c *gin.Context) {
	tenantID := c.Param("id")
	
	stats, err := h.service.GetTenantStats(tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// CheckSubdomainAvailability handles GET /api/tenants/check-subdomain/:subdomain
func (h *Handler) CheckSubdomainAvailability(c *gin.Context) {
	subdomain := c.Param("subdomain")
	
	available, err := h.service.CheckSubdomainAvailability(subdomain)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"subdomain": subdomain,
		"available": available,
	})
}

// GetUpgradeOptions handles GET /api/tenants/:id/upgrade-options
func (h *Handler) GetUpgradeOptions(c *gin.Context) {
	tenantID := c.Param("id")
	
	options, err := h.service.GetPlanUpgradeOptions(tenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"tenant_id": tenantID,
			"upgrade_options": options,
		},
	})
}

// SuspendTenant handles POST /api/tenants/:id/suspend
func (h *Handler) SuspendTenant(c *gin.Context) {
	tenantID := c.Param("id")
	
	var req struct {
		Reason string `json:"reason,omitempty"`
	}
	c.ShouldBindJSON(&req)
	
	err := h.service.SuspendTenant(tenantID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant suspended successfully",
	})
}

// ValidateCustomDomain handles POST /api/tenants/:id/validate-domain
func (h *Handler) ValidateCustomDomain(c *gin.Context) {
	tenantID := c.Param("id")
	
	var req struct {
		Domain string `json:"domain" validate:"required,fqdn"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.ValidateCustomDomain(tenantID, req.Domain)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Domain validated and updated successfully",
	})
}

// TODO: Add more handlers
// - UploadLogo(c *gin.Context) - Handle logo upload
// - GetBillingInfo(c *gin.Context) - Get billing information
// - UpdateDomainSettings(c *gin.Context) - Update custom domain
// - GetUsageMetrics(c *gin.Context) - Get usage metrics
// - ExportTenantData(c *gin.Context) - Export tenant data
