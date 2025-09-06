package components

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for components
type Handler struct {
	service Service
}

// NewHandler creates a new component handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Component handlers

// CreateComponent creates a new component
// @Summary Create a new component
// @Description Create a new customizable component
// @Tags components
// @Accept json
// @Produce json
// @Param component body CreateComponentRequest true "Component data"
// @Success 201 {object} Component
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /components [post]
func (h *Handler) CreateComponent(c *gin.Context) {
	var req CreateComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body", Details: err.Error()})
		return
	}
	
	tenantID := getTenantID(c)
	component, err := h.service.CreateComponent(c.Request.Context(), tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create component", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, component)
}

// GetComponent retrieves a component by ID
// @Summary Get component by ID
// @Description Get a component by its ID
// @Tags components
// @Produce json
// @Param id path string true "Component ID"
// @Success 200 {object} Component
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /components/{id} [get]
func (h *Handler) GetComponent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid component ID"})
		return
	}
	
	tenantID := getTenantID(c)
	component, err := h.service.GetComponent(c.Request.Context(), tenantID, id)
	if err != nil {
		if err.Error() == "component not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Component not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get component", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, component)
}

// GetComponentBySlug retrieves a component by slug
// @Summary Get component by slug
// @Description Get a component by its slug
// @Tags components
// @Produce json
// @Param slug path string true "Component slug"
// @Success 200 {object} Component
// @Failure 404 {object} ErrorResponse
// @Router /components/slug/{slug} [get]
func (h *Handler) GetComponentBySlug(c *gin.Context) {
	slug := c.Param("slug")
	tenantID := getTenantID(c)
	
	component, err := h.service.GetComponentBySlug(c.Request.Context(), tenantID, slug)
	if err != nil {
		if err.Error() == "component not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Component not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get component", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, component)
}

// UpdateComponent updates a component
// @Summary Update component
// @Description Update an existing component
// @Tags components
// @Accept json
// @Produce json
// @Param id path string true "Component ID"
// @Param component body UpdateComponentRequest true "Component update data"
// @Success 200 {object} Component
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /components/{id} [put]
func (h *Handler) UpdateComponent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid component ID"})
		return
	}
	
	var req UpdateComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body", Details: err.Error()})
		return
	}
	
	tenantID := getTenantID(c)
	component, err := h.service.UpdateComponent(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if err.Error() == "component not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Component not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update component", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, component)
}

// DeleteComponent deletes a component
// @Summary Delete component
// @Description Delete a component by ID
// @Tags components
// @Param id path string true "Component ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /components/{id} [delete]
func (h *Handler) DeleteComponent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid component ID"})
		return
	}
	
	tenantID := getTenantID(c)
	err = h.service.DeleteComponent(c.Request.Context(), tenantID, id)
	if err != nil {
		if err.Error() == "component not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Component not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete component", Details: err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ListComponents lists components with filtering
// @Summary List components
// @Description List components with optional filtering
// @Tags components
// @Produce json
// @Param type query string false "Component type"
// @Param status query string false "Component status"
// @Param category query string false "Component category"
// @Param search query string false "Search term"
// @Param featured query bool false "Filter by featured"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort field" default("created_at")
// @Param sort_order query string false "Sort order (ASC/DESC)" default("DESC")
// @Success 200 {object} ComponentListResponse
// @Router /components [get]
func (h *Handler) ListComponents(c *gin.Context) {
	filters := ComponentFilters{
		Type:      ComponentType(c.Query("type")),
		Status:    ComponentStatus(c.Query("status")),
		Category:  ComponentCategory(c.Query("category")),
		Search:    c.Query("search"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "DESC"),
	}
	
	// Parse featured filter
	if featuredStr := c.Query("featured"); featuredStr != "" {
		if featured, err := strconv.ParseBool(featuredStr); err == nil {
			filters.Featured = &featured
		}
	}
	
	// Parse pagination
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		} else {
			filters.Page = 1
		}
	}
	
	if limitStr := c.DefaultQuery("limit", "20"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		} else {
			filters.Limit = 20
		}
	}
	
	tenantID := getTenantID(c)
	response, err := h.service.ListComponents(c.Request.Context(), tenantID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list components", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// DuplicateComponent duplicates a component
// @Summary Duplicate component
// @Description Create a copy of an existing component
// @Tags components
// @Accept json
// @Produce json
// @Param id path string true "Component ID"
// @Param request body DuplicateRequest true "Duplicate request"
// @Success 201 {object} Component
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /components/{id}/duplicate [post]
func (h *Handler) DuplicateComponent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid component ID"})
		return
	}
	
	var req DuplicateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body", Details: err.Error()})
		return
	}
	
	tenantID := getTenantID(c)
	component, err := h.service.DuplicateComponent(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if err.Error() == "component not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Component not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to duplicate component", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, component)
}

// Component Instance handlers

// CreateInstance creates a new component instance
// @Summary Create component instance
// @Description Create a new instance of a component in a theme
// @Tags instances
// @Accept json
// @Produce json
// @Param instance body CreateInstanceRequest true "Instance data"
// @Success 201 {object} ComponentInstance
// @Failure 400 {object} ErrorResponse
// @Router /instances [post]
func (h *Handler) CreateInstance(c *gin.Context) {
	var req CreateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body", Details: err.Error()})
		return
	}
	
	tenantID := getTenantID(c)
	instance, err := h.service.CreateInstance(c.Request.Context(), tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create instance", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, instance)
}

// GetInstance retrieves an instance by ID
// @Summary Get instance by ID
// @Description Get a component instance by its ID
// @Tags instances
// @Produce json
// @Param id path string true "Instance ID"
// @Success 200 {object} ComponentInstance
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /instances/{id} [get]
func (h *Handler) GetInstance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid instance ID"})
		return
	}
	
	tenantID := getTenantID(c)
	instance, err := h.service.GetInstance(c.Request.Context(), tenantID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get instance", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, instance)
}

// ListInstances lists instances for a theme
// @Summary List instances
// @Description List component instances for a specific theme
// @Tags instances
// @Produce json
// @Param theme_id query string true "Theme ID"
// @Success 200 {array} ComponentInstance
// @Failure 400 {object} ErrorResponse
// @Router /instances [get]
func (h *Handler) ListInstances(c *gin.Context) {
	themeIDStr := c.Query("theme_id")
	if themeIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "theme_id is required"})
		return
	}
	
	themeID, err := uuid.Parse(themeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid theme ID"})
		return
	}
	
	tenantID := getTenantID(c)
	instances, err := h.service.ListInstances(c.Request.Context(), tenantID, themeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list instances", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, instances)
}

// Theme handlers

// ListThemes lists themes with filtering
// @Summary List themes
// @Description List themes with optional filtering
// @Tags themes
// @Produce json
// @Param status query string false "Theme status"
// @Param search query string false "Search term"
// @Param active query bool false "Filter by active status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} ThemeListResponse
// @Router /themes [get]
func (h *Handler) ListThemes(c *gin.Context) {
	filters := ThemeFilters{
		Status:    ComponentStatus(c.Query("status")),
		Search:    c.Query("search"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "DESC"),
	}
	
	// Parse active filter
	if activeStr := c.Query("active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			filters.Active = &active
		}
	}
	
	// Parse pagination
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		} else {
			filters.Page = 1
		}
	}
	
	if limitStr := c.DefaultQuery("limit", "20"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		} else {
			filters.Limit = 20
		}
	}
	
	tenantID := getTenantID(c)
	response, err := h.service.ListThemes(c.Request.Context(), tenantID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list themes", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// GetActiveTheme gets the currently active theme
// @Summary Get active theme
// @Description Get the currently active theme for the tenant
// @Tags themes
// @Produce json
// @Success 200 {object} Theme
// @Failure 404 {object} ErrorResponse
// @Router /themes/active [get]
func (h *Handler) GetActiveTheme(c *gin.Context) {
	tenantID := getTenantID(c)
	theme, err := h.service.GetActiveTheme(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "No active theme found"})
		return
	}
	
	c.JSON(http.StatusOK, theme)
}

// Template handlers

// ListTemplates lists component templates
// @Summary List templates
// @Description List available component templates
// @Tags templates
// @Produce json
// @Param type query string false "Template type"
// @Param category query string false "Template category"
// @Param search query string false "Search term"
// @Param free query bool false "Filter by free templates"
// @Param featured query bool false "Filter by featured templates"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} TemplateListResponse
// @Router /templates [get]
func (h *Handler) ListTemplates(c *gin.Context) {
	filters := TemplateFilters{
		Type:      ComponentType(c.Query("type")),
		Category:  ComponentCategory(c.Query("category")),
		Search:    c.Query("search"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "DESC"),
	}
	
	// Parse boolean filters
	if freeStr := c.Query("free"); freeStr != "" {
		if free, err := strconv.ParseBool(freeStr); err == nil {
			filters.Free = &free
		}
	}
	
	if featuredStr := c.Query("featured"); featuredStr != "" {
		if featured, err := strconv.ParseBool(featuredStr); err == nil {
			filters.Featured = &featured
		}
	}
	
	// Parse pagination
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		} else {
			filters.Page = 1
		}
	}
	
	if limitStr := c.DefaultQuery("limit", "20"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		} else {
			filters.Limit = 20
		}
	}
	
	response, err := h.service.ListTemplates(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list templates", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// GetTemplate retrieves a template by ID
// @Summary Get template by ID
// @Description Get a component template by its ID
// @Tags templates
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} ComponentTemplate
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /templates/{id} [get]
func (h *Handler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid template ID"})
		return
	}
	
	template, err := h.service.GetTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Template not found"})
		return
	}
	
	c.JSON(http.StatusOK, template)
}

// Theme Template handlers

// CreateThemeTemplate creates a new theme template
// @Summary Create theme template
// @Description Create a new theme template
// @Tags theme-templates
// @Accept json
// @Produce json
// @Param request body CreateThemeTemplateRequest true "Theme template data"
// @Success 201 {object} ThemeTemplate
// @Failure 400 {object} ErrorResponse
// @Router /theme-templates [post]
func (h *Handler) CreateThemeTemplate(c *gin.Context) {
	var req CreateThemeTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body", Details: err.Error()})
		return
	}
	
	themeTemplate, err := h.service.CreateThemeTemplate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create theme template", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, themeTemplate)
}

// GetThemeTemplate retrieves a theme template by ID
// @Summary Get theme template by ID
// @Description Get a theme template by its ID
// @Tags theme-templates
// @Produce json
// @Param id path string true "Theme template ID"
// @Success 200 {object} ThemeTemplate
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /theme-templates/{id} [get]
func (h *Handler) GetThemeTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid theme template ID"})
		return
	}
	
	themeTemplate, err := h.service.GetThemeTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Theme template not found"})
		return
	}
	
	c.JSON(http.StatusOK, themeTemplate)
}

// UpdateThemeTemplate updates a theme template
// @Summary Update theme template
// @Description Update an existing theme template
// @Tags theme-templates
// @Accept json
// @Produce json
// @Param id path string true "Theme template ID"
// @Param request body UpdateThemeTemplateRequest true "Updated theme template data"
// @Success 200 {object} ThemeTemplate
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /theme-templates/{id} [put]
func (h *Handler) UpdateThemeTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid theme template ID"})
		return
	}
	
	var req UpdateThemeTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body", Details: err.Error()})
		return
	}
	
	themeTemplate, err := h.service.UpdateThemeTemplate(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update theme template", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, themeTemplate)
}

// DeleteThemeTemplate deletes a theme template
// @Summary Delete theme template
// @Description Delete a theme template by ID
// @Tags theme-templates
// @Param id path string true "Theme template ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /theme-templates/{id} [delete]
func (h *Handler) DeleteThemeTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid theme template ID"})
		return
	}
	
	err = h.service.DeleteThemeTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete theme template", Details: err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ListThemeTemplates lists theme templates
// @Summary List theme templates
// @Description List available theme templates
// @Tags theme-templates
// @Produce json
// @Param category query string false "Template category"
// @Param search query string false "Search term"
// @Param free query bool false "Filter by free templates"
// @Param featured query bool false "Filter by featured templates"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} ThemeTemplateListResponse
// @Router /theme-templates [get]
func (h *Handler) ListThemeTemplates(c *gin.Context) {
	filters := ThemeTemplateFilter{
		Category:  c.Query("category"),
		Search:    c.Query("search"),
	}
	
	// Parse boolean filters
	if freeStr := c.Query("free"); freeStr != "" {
		if free, err := strconv.ParseBool(freeStr); err == nil {
			filters.IsFree = &free
		}
	}
	
	if featuredStr := c.Query("featured"); featuredStr != "" {
		if featured, err := strconv.ParseBool(featuredStr); err == nil {
			filters.IsFeatured = &featured
		}
	}
	
	// Parse pagination
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Offset = (page - 1) * filters.Limit
		}
	}
	
	if limitStr := c.DefaultQuery("limit", "20"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		} else {
			filters.Limit = 20
		}
	}
	
	response, err := h.service.ListThemeTemplates(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list theme templates", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// Statistics handler

// GetStats retrieves component statistics
// @Summary Get statistics
// @Description Get component usage statistics for the tenant
// @Tags statistics
// @Produce json
// @Success 200 {object} ComponentStats
// @Router /stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	tenantID := getTenantID(c)
	stats, err := h.service.GetStats(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get statistics", Details: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// Helper types and functions

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// getTenantID extracts tenant ID from context
// This should be set by middleware in a real application
func getTenantID(c *gin.Context) uuid.UUID {
	// In a real application, this would be extracted from JWT token or session
	// For now, we'll use a default tenant ID or extract from header
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		// Return a default tenant ID for development
		return uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}
	
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		// Return a default tenant ID if parsing fails
		return uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}
	
	return tenantID
}