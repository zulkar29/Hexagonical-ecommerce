package tax

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for tax operations
type Handler struct {
	service Service
}

// NewHandler creates a new tax handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Tax Rule endpoints

// CreateTaxRule creates a new tax rule
// @Summary Create tax rule
// @Description Create a new tax rule with validation
// @Tags tax-rules
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param request body CreateTaxRuleRequest true "Tax rule data"
// @Success 201 {object} TaxRuleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tenants/{tenant_id}/tax/rules [post]
func (h *Handler) CreateTaxRule(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	var req CreateTaxRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	rule, err := h.service.CreateTaxRule(c.Request.Context(), tenantID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == ErrRuleCodeExists || err == ErrInvalidRate || err == ErrInvalidDateRange {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, rule)
}

// GetTaxRule retrieves a tax rule by ID
// @Summary Get tax rule
// @Description Get a tax rule by ID
// @Tags tax-rules
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param rule_id path string true "Tax Rule ID"
// @Success 200 {object} TaxRuleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/tenants/{tenant_id}/tax/rules/{rule_id} [get]
func (h *Handler) GetTaxRule(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	ruleID, err := uuid.Parse(c.Param("rule_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}
	
	rule, err := h.service.GetTaxRule(c.Request.Context(), tenantID, ruleID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "tax rule not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, rule)
}

// UpdateTaxRule updates a tax rule
// @Summary Update tax rule
// @Description Update a tax rule
// @Tags tax-rules
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param rule_id path string true "Tax Rule ID"
// @Param request body UpdateTaxRuleRequest true "Tax rule update data"
// @Success 200 {object} TaxRuleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/tenants/{tenant_id}/tax/rules/{rule_id} [put]
func (h *Handler) UpdateTaxRule(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	ruleID, err := uuid.Parse(c.Param("rule_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}
	
	var req UpdateTaxRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	rule, err := h.service.UpdateTaxRule(c.Request.Context(), tenantID, ruleID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "tax rule not found" {
			statusCode = http.StatusNotFound
		} else if err == ErrInvalidRate || err == ErrInvalidDateRange {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, rule)
}

// DeleteTaxRule deletes a tax rule
// @Summary Delete tax rule
// @Description Delete a tax rule
// @Tags tax-rules
// @Param tenant_id path string true "Tenant ID"
// @Param rule_id path string true "Tax Rule ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/tenants/{tenant_id}/tax/rules/{rule_id} [delete]
func (h *Handler) DeleteTaxRule(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	ruleID, err := uuid.Parse(c.Param("rule_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}
	
	err = h.service.DeleteTaxRule(c.Request.Context(), tenantID, ruleID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "tax rule not found" {
			statusCode = http.StatusNotFound
		} else if err == ErrCannotDeleteRule {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}

// ListTaxRules lists tax rules with filtering and pagination
// @Summary List tax rules
// @Description List tax rules with filtering and pagination
// @Tags tax-rules
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param status query string false "Filter by status"
// @Param type query string false "Filter by type"
// @Param tax_type query string false "Filter by tax type"
// @Param country query string false "Filter by country"
// @Param state query string false "Filter by state"
// @Param search query string false "Search in name, description, or code"
// @Success 200 {object} PaginatedTaxRulesResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/tenants/{tenant_id}/tax/rules [get]
func (h *Handler) ListTaxRules(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	// Parse filters
	filter := TaxRuleFilter{
		Status:  c.Query("status"),
		Type:    c.Query("type"),
		TaxType: c.Query("tax_type"),
		Country: c.Query("country"),
		State:   c.Query("state"),
		Search:  c.Query("search"),
	}
	
	rules, total, err := h.service.ListTaxRules(c.Request.Context(), tenantID, filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	response := gin.H{
		"data":       rules,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}
	
	c.JSON(http.StatusOK, response)
}

// CalculateTax calculates tax for a given request
// @Summary Calculate tax
// @Description Calculate tax for a given amount and location
// @Tags tax-calculations
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param request body TaxCalculationRequest true "Tax calculation data"
// @Success 200 {object} TaxCalculationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tenants/{tenant_id}/tax/calculate [post]
func (h *Handler) CalculateTax(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	var req TaxCalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	result, err := h.service.CalculateTax(c.Request.Context(), tenantID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == ErrInvalidAmount || err == ErrInvalidLocation || err == ErrNoApplicableRules {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

// GetTaxStats retrieves tax statistics
// @Summary Get tax statistics
// @Description Get comprehensive tax statistics
// @Tags tax-analytics
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {object} TaxStats
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/tenants/{tenant_id}/tax/stats [get]
func (h *Handler) GetTaxStats(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	stats, err := h.service.GetTaxStats(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// CleanupExpiredRules cleans up expired tax rules
// @Summary Cleanup expired rules
// @Description Clean up expired tax rules
// @Tags tax-maintenance
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {object} CleanupResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/tenants/{tenant_id}/tax/cleanup/rules [post]
func (h *Handler) CleanupExpiredRules(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	count, err := h.service.CleanupExpiredRules(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Expired rules cleaned up successfully",
		"count":   count,
	})
}

// GetApplicableTaxRules retrieves applicable tax rules for a calculation request
// @Summary Get applicable tax rules
// @Description Get applicable tax rules for a calculation request
// @Tags tax-utilities
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param request body TaxCalculationRequest true "Tax calculation data"
// @Success 200 {array} TaxRuleResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/tenants/{tenant_id}/tax/rules/applicable [post]
func (h *Handler) GetApplicableTaxRules(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}
	
	var req TaxCalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	rules, err := h.service.GetApplicableTaxRules(c.Request.Context(), tenantID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == ErrInvalidLocation || err == ErrNoApplicableRules {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

// ValidateLocation validates a location for tax calculation
// @Summary Validate location
// @Description Validate a location for tax calculation
// @Tags tax-utilities
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param request body LocationValidationRequest true "Location data"
// @Success 200 {object} LocationValidationResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/tenants/{tenant_id}/tax/validate/location [post]
func (h *Handler) ValidateLocation(c *gin.Context) {
	var req struct {
		Country string `json:"country" binding:"required"`
		State   string `json:"state"`
		City    string `json:"city"`
		ZipCode string `json:"zip_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.service.ValidateLocation(c.Request.Context(), req.Country, req.State, req.City, req.ZipCode)
	isValid := err == nil
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"valid":   isValid,
		"country": req.Country,
		"state":   req.State,
		"city":    req.City,
		"zip_code": req.ZipCode,
	})
}