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
	
	// TODO: Implement tenant statistics
	// This would include:
	// - Product count
	// - Order count
	// - Revenue
	// - Storage usage
	// - Bandwidth usage
	
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"tenant_id":      tenantID,
			"product_count":  0,
			"order_count":    0,
			"revenue":        0,
			"storage_used":   0,
			"bandwidth_used": 0,
		},
	})
}

// CheckSubdomainAvailability handles GET /api/tenants/check-subdomain/:subdomain
func (h *Handler) CheckSubdomainAvailability(c *gin.Context) {
	subdomain := c.Param("subdomain")
	
	_, err := h.service.GetTenantBySubdomain(subdomain)
	available := err != nil // If error (not found), then available
	
	c.JSON(http.StatusOK, gin.H{
		"subdomain": subdomain,
		"available": available,
	})
}

// TODO: Add more handlers
// - UploadLogo(c *gin.Context) - Handle logo upload
// - GetBillingInfo(c *gin.Context) - Get billing information
// - UpdateDomainSettings(c *gin.Context) - Update custom domain
// - GetUsageMetrics(c *gin.Context) - Get usage metrics
// - ExportTenantData(c *gin.Context) - Export tenant data
