package platform

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles platform admin HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new platform handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers platform admin routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Platform dashboard and stats
	router.GET("/dashboard", h.GetPlatformDashboard)
	router.GET("/stats", h.GetPlatformStats)
	router.GET("/system/status", h.GetSystemStatus)

	// Platform admin management
	admins := router.Group("/admins")
	{
		admins.GET("", h.ListPlatformAdmins)
		admins.POST("", h.CreatePlatformAdmin)
		admins.GET("/:id", h.GetPlatformAdmin)
		admins.PATCH("/:id", h.UpdatePlatformAdmin)
		admins.DELETE("/:id", h.DeletePlatformAdmin)
	}

	// Platform role management
	roles := router.Group("/roles")
	{
		roles.GET("", h.ListPlatformRoles)
		roles.POST("", h.CreatePlatformRole)
		roles.GET("/:id", h.GetPlatformRole)
		roles.PATCH("/:id", h.UpdatePlatformRole)
		roles.DELETE("/:id", h.DeletePlatformRole)
	}

	// Platform settings
	settings := router.Group("/settings")
	{
		settings.GET("", h.GetPlatformSettings)
		settings.PATCH("", h.UpdatePlatformSettings)
	}

	// Platform tenant management
	tenants := router.Group("/tenants")
	{
		tenants.GET("", h.ListAllTenants)
		tenants.GET("/:id", h.GetTenantDetails)
		tenants.PATCH("/:id", h.UpdateTenant)
		tenants.DELETE("/:id", h.DeleteTenant)
	}

	// Platform audit logs
	router.GET("/audit-logs", h.GetPlatformAuditLogs)
}

// GetPlatformDashboard retrieves platform dashboard data
func (h *Handler) GetPlatformDashboard(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	// Get platform stats
	stats, err := h.service.GetPlatformStats(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get platform stats",
			"details": err.Error(),
		})
		return
	}

	// Get system status
	systemStatus, err := h.service.GetSystemStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get system status",
			"details": err.Error(),
		})
		return
	}

	dashboard := gin.H{
		"stats":         stats,
		"system_status": systemStatus,
		"period":        period,
		"generated_at":  time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dashboard,
	})
}

// GetPlatformStats retrieves platform statistics
func (h *Handler) GetPlatformStats(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	stats, err := h.service.GetPlatformStats(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get platform stats",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetSystemStatus retrieves system health status
func (h *Handler) GetSystemStatus(c *gin.Context) {
	status, err := h.service.GetSystemStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get system status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}

// ListPlatformAdmins retrieves platform administrators
func (h *Handler) ListPlatformAdmins(c *gin.Context) {
	role := c.Query("role")

	admins, err := h.service.ListPlatformAdmins(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list platform admins",
			"details": err.Error(),
		})
		return
	}

	// Remove password from response
	for _, admin := range admins {
		admin.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    admins,
		"total":   len(admins),
	})
}

// GetPlatformAdmin retrieves a platform administrator by ID
func (h *Handler) GetPlatformAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid admin ID",
		})
		return
	}

	admin, err := h.service.GetPlatformAdmin(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Platform admin not found",
			"details": err.Error(),
		})
		return
	}

	// Remove password from response
	admin.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    admin,
	})
}

// CreatePlatformAdmin creates a new platform administrator
func (h *Handler) CreatePlatformAdmin(c *gin.Context) {
	var req CreatePlatformAdminRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": bindErr.Error(),
		})
		return
	}

	admin, err := h.service.CreatePlatformAdmin(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create platform admin",
			"details": err.Error(),
		})
		return
	}

	// Remove password from response
	admin.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    admin,
		"message": "Platform admin created successfully",
	})
}

// UpdatePlatformAdmin updates a platform administrator
func (h *Handler) UpdatePlatformAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid admin ID",
		})
		return
	}

	var req UpdatePlatformAdminRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": bindErr.Error(),
		})
		return
	}

	admin, err := h.service.UpdatePlatformAdmin(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update platform admin",
			"details": err.Error(),
		})
		return
	}

	// Remove password from response
	admin.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    admin,
		"message": "Platform admin updated successfully",
	})
}

// DeletePlatformAdmin deletes a platform administrator
func (h *Handler) DeletePlatformAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid admin ID",
		})
		return
	}

	err = h.service.DeletePlatformAdmin(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete platform admin",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Platform admin deleted successfully",
	})
}

// ListPlatformRoles retrieves platform roles
func (h *Handler) ListPlatformRoles(c *gin.Context) {
	roles, err := h.service.ListPlatformRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list platform roles",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roles,
		"total":   len(roles),
	})
}

// GetPlatformRole retrieves a platform role by ID
func (h *Handler) GetPlatformRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role ID",
		})
		return
	}

	role, err := h.service.GetPlatformRole(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Platform role not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    role,
	})
}

// CreatePlatformRole creates a new platform role
func (h *Handler) CreatePlatformRole(c *gin.Context) {
	var req CreatePlatformRoleRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": bindErr.Error(),
		})
		return
	}

	role, err := h.service.CreatePlatformRole(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create platform role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    role,
		"message": "Platform role created successfully",
	})
}

// UpdatePlatformRole updates a platform role
func (h *Handler) UpdatePlatformRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role ID",
		})
		return
	}

	var req UpdatePlatformRoleRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": bindErr.Error(),
		})
		return
	}

	role, err := h.service.UpdatePlatformRole(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update platform role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    role,
		"message": "Platform role updated successfully",
	})
}

// DeletePlatformRole deletes a platform role
func (h *Handler) DeletePlatformRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role ID",
		})
		return
	}

	err = h.service.DeletePlatformRole(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete platform role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Platform role deleted successfully",
	})
}

// GetPlatformSettings retrieves platform settings
func (h *Handler) GetPlatformSettings(c *gin.Context) {
	category := c.Query("category")

	settings, err := h.service.GetPlatformSettings(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get platform settings",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
		"total":   len(settings),
	})
}

// UpdatePlatformSettings updates platform settings
func (h *Handler) UpdatePlatformSettings(c *gin.Context) {
	var req UpdatePlatformSettingsRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": bindErr.Error(),
		})
		return
	}

	settings, err := h.service.UpdatePlatformSettings(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update platform settings",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
		"message": "Platform settings updated successfully",
	})
}

// ListAllTenants retrieves all tenants for platform admin
func (h *Handler) ListAllTenants(c *gin.Context) {
	req := &ListTenantsRequest{
		Status: c.Query("status"),
	}

	// Parse include parameter
	if includeStr := c.Query("include"); includeStr != "" {
		req.Include = strings.Split(includeStr, ",")
	}

	// Parse pagination
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, parseErr := strconv.Atoi(limitStr); parseErr == nil {
			req.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, parseErr := strconv.Atoi(offsetStr); parseErr == nil {
			req.Offset = offset
		}
	}

	response, err := h.service.ListAllTenants(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list tenants",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response.Tenants,
		"total":   response.Total,
		"limit":   response.Limit,
		"offset":  response.Offset,
	})
}

// GetTenantDetails retrieves detailed tenant information
func (h *Handler) GetTenantDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant ID",
		})
		return
	}

	// Parse include parameter
	var include []string
	if includeStr := c.Query("include"); includeStr != "" {
		include = strings.Split(includeStr, ",")
	}

	tenant, err := h.service.GetTenantDetails(c.Request.Context(), id, include)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tenant not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tenant,
	})
}

// UpdateTenant updates tenant information
func (h *Handler) UpdateTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant ID",
		})
		return
	}

	var req UpdateTenantRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": bindErr.Error(),
		})
		return
	}

	err = h.service.UpdateTenant(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update tenant",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tenant updated successfully",
	})
}

// DeleteTenant deletes a tenant
func (h *Handler) DeleteTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant ID",
		})
		return
	}

	err = h.service.DeleteTenant(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete tenant",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tenant deleted successfully",
	})
}

// GetPlatformAuditLogs retrieves platform audit logs
func (h *Handler) GetPlatformAuditLogs(c *gin.Context) {
	req := &GetAuditLogsRequest{}

	// Parse query parameters
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, parseErr := uuid.Parse(userIDStr); parseErr == nil {
			req.UserID = &userID
		}
	}

	if tenantIDStr := c.Query("tenant_id"); tenantIDStr != "" {
		if tenantID, parseErr := uuid.Parse(tenantIDStr); parseErr == nil {
			req.TenantID = &tenantID
		}
	}

	req.Action = c.Query("action")
	req.Resource = c.Query("resource")

	// Parse date filters
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if dateFrom, parseErr := time.Parse(time.RFC3339, dateFromStr); parseErr == nil {
			req.DateFrom = &dateFrom
		}
	}

	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if dateTo, parseErr := time.Parse(time.RFC3339, dateToStr); parseErr == nil {
			req.DateTo = &dateTo
		}
	}

	// Parse pagination
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, parseErr := strconv.Atoi(limitStr); parseErr == nil {
			req.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, parseErr := strconv.Atoi(offsetStr); parseErr == nil {
			req.Offset = offset
		}
	}

	response, err := h.service.GetPlatformAuditLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get audit logs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response.Logs,
		"total":   response.Total,
		"limit":   response.Limit,
		"offset":  response.Offset,
	})
}