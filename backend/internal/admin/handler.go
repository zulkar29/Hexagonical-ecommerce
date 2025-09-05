package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles admin-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new admin handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers admin routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	admin := router.Group("/admin")
	{
		// Dashboard endpoints
		admin.GET("/dashboard", h.GetDashboard)
		admin.GET("/quick-stats", h.GetQuickStats)
		
		// Staff management
		admin.GET("/staff", h.ListStaff)
		admin.PATCH("/staff/:id", h.ManageStaff)
		
		// Role management
		admin.GET("/roles", h.ListRoles)
		admin.PATCH("/roles/:id", h.ManageRoles)
		
		// Activity logs
		admin.GET("/activity-logs", h.GetActivityLogs)
		
		// System health
		admin.GET("/system-health", h.GetSystemHealth)
	}
}

// GetDashboard handles GET /admin/dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	// Extract query parameters
	period := c.DefaultQuery("period", "week")
	metricsParam := c.Query("metrics")
	
	var metrics []string
	if metricsParam != "" {
		// Parse comma-separated metrics
		// For simplicity, we'll handle this in service
		metrics = []string{metricsParam}
	}
	
	req := DashboardRequest{
		Period:  period,
		Metrics: metrics,
	}
	
	// TODO: Extract tenant ID from context
	tenantID := h.extractTenantID(c)
	
	stats, err := h.service.GetDashboardStats(c.Request.Context(), tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetQuickStats handles GET /admin/quick-stats
func (h *Handler) GetQuickStats(c *gin.Context) {
	// TODO: Extract tenant ID from context
	tenantID := h.extractTenantID(c)
	
	stats, err := h.service.GetQuickStats(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// ListStaff handles GET /admin/staff
func (h *Handler) ListStaff(c *gin.Context) {
	role := c.Query("role")
	status := c.Query("status")
	
	// TODO: Extract tenant ID from context
	tenantID := h.extractTenantID(c)
	
	staff, err := h.service.ListStaff(c.Request.Context(), tenantID, role, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": staff})
}

// ManageStaff handles PATCH /admin/staff/:id
func (h *Handler) ManageStaff(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid staff ID"})
		return
	}
	
	action := c.Query("action")
	if action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Action parameter is required"})
		return
	}
	
	// TODO: Extract tenant ID from context
	tenantID := h.extractTenantID(c)
	
	switch action {
	case "create":
		var req StaffRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		staff, err := h.service.CreateStaff(c.Request.Context(), tenantID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, gin.H{"data": staff})
		
	case "update":
		var req StaffRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		staff, err := h.service.UpdateStaff(c.Request.Context(), tenantID, id, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"data": staff})
		
	case "delete":
		err := h.service.DeleteStaff(c.Request.Context(), tenantID, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Staff deleted successfully"})
		
	case "assign_roles":
		var req struct {
			Roles []string `json:"roles"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		err := h.service.AssignRoles(c.Request.Context(), tenantID, id, req.Roles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Roles assigned successfully"})
		
	case "change_status":
		var req struct {
			Status string `json:"status"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		err := h.service.ChangeStaffStatus(c.Request.Context(), tenantID, id, req.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Status changed successfully"})
		
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
	}
}

// ListRoles handles GET /admin/roles
func (h *Handler) ListRoles(c *gin.Context) {
	includePermissions := c.Query("include_permissions") == "true"
	
	// TODO: Extract tenant ID from context
	tenantID := h.extractTenantID(c)
	
	roles, err := h.service.ListRoles(c.Request.Context(), tenantID, includePermissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// ManageRoles handles PATCH /admin/roles/:id
func (h *Handler) ManageRoles(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}
	
	action := c.Query("action")
	if action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Action parameter is required"})
		return
	}
	
	// TODO: Extract tenant ID from context
	tenantID := h.extractTenantID(c)
	
	switch action {
	case "create":
		var req RoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		role, err := h.service.CreateRole(c.Request.Context(), tenantID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, gin.H{"data": role})
		
	case "update":
		var req RoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		role, err := h.service.UpdateRole(c.Request.Context(), tenantID, id, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"data": role})
		
	case "delete":
		err := h.service.DeleteRole(c.Request.Context(), tenantID, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
		
	case "assign_permissions":
		var req struct {
			Permissions []string `json:"permissions"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		err := h.service.AssignPermissions(c.Request.Context(), tenantID, id, req.Permissions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Permissions assigned successfully"})
		
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
	}
}

// GetActivityLogs handles GET /admin/activity-logs
func (h *Handler) GetActivityLogs(c *gin.Context) {
	var filter ActivityLogFilter
	
	// Parse query parameters
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			filter.UserID = &userID
		}
	}
	
	filter.Action = c.Query("action")
	
	// Parse date parameters (simplified)
	// TODO: Implement proper date parsing
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}
	
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}
	
	// Set defaults
	if filter.Limit == 0 {
		filter.Limit = 50
	}
	
	// TODO: Extract tenant ID from context
	tenantID := h.extractTenantID(c)
	
	logs, err := h.service.GetActivityLogs(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// GetSystemHealth handles GET /admin/system-health
func (h *Handler) GetSystemHealth(c *gin.Context) {
	health, err := h.service.GetSystemHealth(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": health})
}

// Helper methods

// extractTenantID extracts tenant ID from context
// TODO: Implement proper tenant ID extraction from JWT or context
func (h *Handler) extractTenantID(c *gin.Context) *uuid.UUID {
	// This is a placeholder implementation
	// In a real implementation, you would extract this from JWT token or middleware
	return nil
}

// extractUserID extracts user ID from context
// TODO: Implement proper user ID extraction from JWT
func (h *Handler) extractUserID(c *gin.Context) uuid.UUID {
	// This is a placeholder implementation
	// In a real implementation, you would extract this from JWT token
	return uuid.New()
}