package referral

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Request/Response types
type GenerateReferralRequest struct {
	UserID         uuid.UUID `json:"user_id" binding:"required"`
	CommissionRate float64   `json:"commission_rate" binding:"required,min=0,max=1"`
	ExpiresInDays  int       `json:"expires_in_days" binding:"required,min=1"`
}

type GenerateReferralResponse struct {
	ID             uuid.UUID  `json:"id"`
	ReferralCode   string     `json:"referral_code"`
	ReferrerID     uuid.UUID  `json:"referrer_id"`
	CommissionRate float64    `json:"commission_rate"`
	ExpiresAt      *time.Time `json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

type ApplyReferralRequest struct {
	ReferralCode string    `json:"referral_code" binding:"required"`
	RefereeID    uuid.UUID `json:"referee_id" binding:"required"`
}

type ApplyReferralResponse struct {
	Success        bool      `json:"success"`
	ReferralCode   string    `json:"referral_code"`
	ReferrerID     uuid.UUID `json:"referrer_id"`
	CommissionRate float64   `json:"commission_rate"`
	Message        string    `json:"message"`
}

type ReferralStatsResponse struct {
	TotalReferrals          int     `json:"total_referrals"`
	SuccessfulReferrals     int     `json:"successful_referrals"`
	PendingReferrals        int     `json:"pending_referrals"`
	TotalCommissionEarned   float64 `json:"total_commission_earned"`
	TotalCommissionPaid     float64 `json:"total_commission_paid"`
	TotalCommissionPending  float64 `json:"total_commission_pending"`
	LastReferralAt          *time.Time `json:"last_referral_at"`
	LastCommissionAt        *time.Time `json:"last_commission_at"`
}

type CommissionHistoryResponse struct {
	Commissions []ReferralCommission `json:"commissions"`
	Total       int                  `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
}

// Handler handles HTTP requests for referral operations
type Handler struct {
	service Service
}

// NewHandler creates a new referral handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers referral routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	referrals := router.Group("/referrals")
	{
		referrals.POST("/generate", h.GenerateReferralCode)
		referrals.POST("/apply", h.ApplyReferralCode)
		referrals.GET("/validate/:code", h.ValidateReferralCode)
		referrals.GET("/my", h.GetMyReferrals)
		referrals.DELETE("/:id", h.DeactivateReferral)
		referrals.GET("/stats", h.GetMyReferralStats)
		referrals.GET("/stats/tenant", h.GetTenantReferralStats) // Admin only
	}

	commissions := router.Group("/commissions")
	{
		commissions.GET("/my", h.GetMyCommissions)
		commissions.GET("/pending", h.GetPendingCommissions) // Admin only
		commissions.PUT("/:id/paid", h.MarkCommissionAsPaid)  // Admin only
	}
}

// GenerateReferralCodeRequest represents the request to generate a referral code
type GenerateReferralCodeRequest struct {
	CommissionRate float64    `json:"commission_rate" binding:"required,min=0,max=1"`
	ExpiresAt      *time.Time `json:"expires_at"`
}

// ApplyReferralCodeRequest represents the request to apply a referral code
type ApplyReferralCodeRequest struct {
	ReferralCode string    `json:"referral_code" binding:"required"`
	RefereeID    uuid.UUID `json:"referee_id" binding:"required"`
}

// GenerateReferralCode generates a new referral code
func (h *Handler) GenerateReferralCode(c *gin.Context) {
	var req GenerateReferralCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get tenant and user from context (assuming middleware sets these)
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	referral, err := h.service.GenerateReferralCode(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		userID.(uuid.UUID),
		req.CommissionRate,
		req.ExpiresAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := GenerateReferralResponse{
		ID:             referral.ID,
		ReferralCode:   referral.ReferralCode,
		ReferrerID:     referral.ReferrerID,
		CommissionRate: referral.CommissionRate,
		ExpiresAt:      referral.ExpiresAt,
		CreatedAt:      referral.CreatedAt,
	}
	c.JSON(http.StatusCreated, response)
}

// ApplyReferralCode applies a referral code
func (h *Handler) ApplyReferralCode(c *gin.Context) {
	var req ApplyReferralCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	referral, err := h.service.ApplyReferralCode(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		req.ReferralCode,
		req.RefereeID,
	)
	if err != nil {
		if err.Error() == "invalid referral code" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	response := ApplyReferralResponse{
		Success:        true,
		ReferralCode:   referral.ReferralCode,
		ReferrerID:     referral.ReferrerID,
		CommissionRate: referral.CommissionRate,
		Message:        "Referral code applied successfully",
	}
	c.JSON(http.StatusOK, response)
}

// ValidateReferralCode validates a referral code
func (h *Handler) ValidateReferralCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "referral code is required"})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	referral, err := h.service.ValidateReferralCode(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		code,
	)
	if err != nil {
		if err.Error() == "invalid referral code" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":           true,
		"referral_code":   referral.ReferralCode,
		"commission_rate": referral.CommissionRate,
		"expires_at":      referral.ExpiresAt,
	})
}

// GetMyReferrals gets the current user's referrals
func (h *Handler) GetMyReferrals(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	referrals, err := h.service.GetUserReferrals(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		userID.(uuid.UUID),
		limit,
		offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"referrals": referrals,
		"limit":     limit,
		"offset":    offset,
	})
}

// DeactivateReferral deactivates a referral
func (h *Handler) DeactivateReferral(c *gin.Context) {
	referralIDStr := c.Param("id")
	referralID, err := uuid.Parse(referralIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid referral ID"})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	err = h.service.DeactivateReferral(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		referralID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "referral deactivated successfully"})
}

// GetMyReferralStats gets the current user's referral statistics
func (h *Handler) GetMyReferralStats(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	stats, err := h.service.GetUserReferralStats(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		userID.(uuid.UUID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ReferralStatsResponse{
		TotalReferrals:          int(stats.TotalReferrals),
		SuccessfulReferrals:     int(stats.CompletedReferrals),
		PendingReferrals:        int(stats.ActiveReferrals),
		TotalCommissionEarned:   stats.TotalCommissions,
		TotalCommissionPaid:     stats.PaidCommissions,
		TotalCommissionPending:  stats.PendingCommissions,
		LastReferralAt:          nil,
		LastCommissionAt:        nil,
	}
	c.JSON(http.StatusOK, response)
}

// GetTenantReferralStats gets tenant-wide referral statistics (admin only)
func (h *Handler) GetTenantReferralStats(c *gin.Context) {
	// Check if user is admin (assuming middleware sets this)
	isAdmin, exists := c.Get("is_admin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	stats, err := h.service.GetTenantReferralStats(
		c.Request.Context(),
		tenantID.(uuid.UUID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetMyCommissions gets the current user's commissions
func (h *Handler) GetMyCommissions(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	commissions, err := h.service.GetUserCommissions(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		userID.(uuid.UUID),
		limit,
		offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to slice of values for response
	commissionValues := make([]ReferralCommission, len(commissions))
	for i, commission := range commissions {
		commissionValues[i] = *commission
	}

	response := CommissionHistoryResponse{
		Commissions: commissionValues,
		Total:       len(commissions),
		Page:        offset/limit + 1,
		Limit:       limit,
	}
	c.JSON(http.StatusOK, response)
}

// GetPendingCommissions gets pending commissions (admin only)
func (h *Handler) GetPendingCommissions(c *gin.Context) {
	// Check if user is admin
	isAdmin, exists := c.Get("is_admin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	commissions, err := h.service.ProcessPendingCommissions(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"commissions": commissions,
		"limit":       limit,
	})
}

// MarkCommissionAsPaid marks a commission as paid (admin only)
func (h *Handler) MarkCommissionAsPaid(c *gin.Context) {
	// Check if user is admin
	isAdmin, exists := c.Get("is_admin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	commissionIDStr := c.Param("id")
	commissionID, err := uuid.Parse(commissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commission ID"})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	err = h.service.MarkCommissionAsPaid(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		commissionID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "commission marked as paid successfully"})
}