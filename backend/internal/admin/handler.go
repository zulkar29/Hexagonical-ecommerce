package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"ecommerce-saas/internal/shared/handlers"
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

	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	stats, err := h.service.GetDashboardStats(c.Request.Context(), &tenantID, req)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"data": stats})
}

// GetQuickStats handles GET /admin/quick-stats
func (h *Handler) GetQuickStats(c *gin.Context) {
	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	stats, err := h.service.GetQuickStats(c.Request.Context(), &tenantID)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"data": stats})
}

// ListStaff handles GET /admin/staff
func (h *Handler) ListStaff(c *gin.Context) {
	role := c.Query("role")
	status := c.Query("status")

	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	staff, err := h.service.ListStaff(c.Request.Context(), &tenantID, role, status)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"data": staff})
}

// ManageStaff handles PATCH /admin/staff/:id
func (h *Handler) ManageStaff(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid staff ID"))
		return
	}

	action := c.Query("action")
	if action == "" {
		handlers.HandleError(c, fmt.Errorf("action parameter is required"))
		return
	}

	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	switch action {
	case "create":
		var req StaffRequest
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}

		staff, err := h.service.CreateStaff(c.Request.Context(), &tenantID, req)
		if err != nil {
			handlers.HandleError(c, err)
			return
		}

		handlers.RespondWithCreated(c, gin.H{"data": staff})

	case "update":
		var req StaffRequest
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}

		staff, err := h.service.UpdateStaff(c.Request.Context(), &tenantID, id, req)
		if err != nil {
			handlers.HandleError(c, err)
			return
		}

		handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"data": staff})

	case "delete":
		err := h.service.DeleteStaff(c.Request.Context(), &tenantID, id)
		if err != nil {
			handlers.HandleError(c, err)
			return
		}

		handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Staff deleted successfully"})

	case "assign_roles":
		var req struct {
			Roles []string `json:"roles"`
		}
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}

		err := h.service.AssignRoles(c.Request.Context(), &tenantID, id, req.Roles)
		if err != nil {
			handlers.HandleError(c, err)
			return
		}

		handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Roles assigned successfully"})

	case "change_status":
		var req struct {
			Status string `json:"status"`
		}
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}

		err := h.service.ChangeStaffStatus(c.Request.Context(), &tenantID, id, req.Status)
		if err != nil {
			handlers.HandleError(c, err)
			return
		}

		handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Status changed successfully"})

	default:
		handlers.HandleError(c, fmt.Errorf("invalid action"))
	}
}

// ListRoles handles GET /admin/roles
func (h *Handler) ListRoles(c *gin.Context) {
	includePermissions := c.Query("include_permissions") == "true"

	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	roles, err := h.service.ListRoles(c.Request.Context(), &tenantID, includePermissions)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"data": roles})
}

// ManageRoles handles PATCH /admin/roles/:id
func (h *Handler) ManageRoles(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handlers.HandleError(c, fmt.Errorf("invalid role ID"))
		return
	}

	action := c.Query("action")
	if action == "" {
		handlers.HandleError(c, fmt.Errorf("action parameter is required"))
		return
	}

	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	switch action {
	case "create":
		var req RoleRequest
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}

		role, err := h.service.CreateRole(c.Request.Context(), &tenantID, req)
		if err != nil {
			handlers.HandleError(c, err)
			return
		}

		handlers.RespondWithCreated(c, gin.H{"data": role})

	case "update":
		var req RoleRequest
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}

		role, err := h.service.UpdateRole(c.Request.Context(), &tenantID, id, req)
		if err != nil {
			handlers.HandleError(c, err)
			return
		}

		handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"data": role})

	case "delete":
		err := h.service.DeleteRole(c.Request.Context(), &tenantID, id)
		if err != nil {
			handlers.HandleError(c, err)
			return
		}

		handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Role deleted successfully"})

	case "assign_permissions":
		var req struct {
			Permissions []string `json:"permissions"`
		}
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}

		err := h.service.AssignPermissions(c.Request.Context(), &tenantID, id, req.Permissions)
		if err != nil {
			handlers.HandleError(c, err)
			return
		}

		handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Permissions assigned successfully"})

	default:
		handlers.HandleError(c, fmt.Errorf("invalid action"))
	}
}

// GetActivityLogs handles GET /admin/activity-logs
func (h *Handler) GetActivityLogs(c *gin.Context) {
	var filter ActivityLogFilter

	// Parse query parameters
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, parseErr := uuid.Parse(userIDStr); parseErr == nil {
			filter.UserID = &userID
		}
	}

	filter.Action = c.Query("action")

	// Parse date parameters (simplified)
	// TODO: Implement proper date parsing

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, limitErr := strconv.Atoi(limitStr); limitErr == nil {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, offsetErr := strconv.Atoi(offsetStr); offsetErr == nil {
			filter.Offset = offset
		}
	}

	// Set defaults
	if filter.Limit == 0 {
		filter.Limit = 50
	}

	tenantID, ok := handlers.RequireTenantID(c)
	if !ok {
		return
	}

	logs, err := h.service.GetActivityLogs(c.Request.Context(), &tenantID, filter)
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"data": logs})
}

// GetSystemHealth handles GET /admin/system-health
func (h *Handler) GetSystemHealth(c *gin.Context) {
	health, err := h.service.GetSystemHealth(c.Request.Context())
	if err != nil {
		handlers.HandleError(c, err)
		return
	}

	handlers.RespondWithSuccess(c, http.StatusOK, gin.H{"data": health})
}

// Helper methods
