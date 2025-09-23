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

// RegisterRoutes registers referral/affiliate routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Referral/Affiliate routes
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

	// Affiliate-specific routes
	affiliates := router.Group("/affiliates")
	{
		affiliates.POST("/register", h.CreateAffiliateAccount)
		affiliates.PUT("/:id/settings", h.UpdateAffiliateSettings)
		affiliates.GET("/performance", h.GetAffiliatePerformance)
		affiliates.GET("/top", h.GetTopPerformingAffiliates) // Admin only
	}

	// Click tracking routes
	tracking := router.Group("/tracking")
	{
		tracking.POST("/click", h.TrackAffiliateClick)
		tracking.GET("/clicks", h.GetAffiliateClicks)
		tracking.GET("/clicks/range", h.GetClicksByDateRange)
	}

	// Commission routes
	commissions := router.Group("/commissions")
	{
		commissions.GET("/my", h.GetMyCommissions)
		commissions.GET("/pending", h.GetPendingCommissions) // Admin only
		commissions.PUT("/:id/paid", h.MarkCommissionAsPaid)  // Admin only
		commissions.POST("/order", h.CreateOrderCommission)   // Internal use
	}

	// Payout routes
	payouts := router.Group("/payouts")
	{
		payouts.POST("/batch", h.CreatePayoutBatch)     // Admin only
		payouts.PUT("/batch/:id/process", h.ProcessPayoutBatch) // Admin only
		payouts.GET("/batches", h.GetPayoutBatches)     // Admin only
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

// Affiliate marketing handlers

// CreateAffiliateAccountRequest represents the request to create an affiliate account
type CreateAffiliateAccountRequest struct {
	AffiliateType   string  `json:"affiliate_type" binding:"required"`
	CommissionRate  float64 `json:"commission_rate" binding:"required,min=0,max=1"`
	PayoutThreshold float64 `json:"payout_threshold" binding:"required,min=0"`
}

// TrackAffiliateClickRequest represents the request to track an affiliate click
type TrackAffiliateClickRequest struct {
	ReferralCode string                 `json:"referral_code" binding:"required"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
	Referrer     string                 `json:"referrer"`
	UTMSource    string                 `json:"utm_source"`
	UTMMedium    string                 `json:"utm_medium"`
	UTMCampaign  string                 `json:"utm_campaign"`
	UTMTerm      string                 `json:"utm_term"`
	UTMContent   string                 `json:"utm_content"`
	DeviceType   string                 `json:"device_type"`
	Country      string                 `json:"country"`
	City         string                 `json:"city"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// CreateOrderCommissionRequest represents the request to create an order commission
type CreateOrderCommissionRequest struct {
	OrderID     string  `json:"order_id" binding:"required"`
	OrderAmount float64 `json:"order_amount" binding:"required,min=0"`
}

// CreatePayoutBatchRequest represents the request to create a payout batch
type CreatePayoutBatchRequest struct {
	AffiliateIDs []string `json:"affiliate_ids" binding:"required,min=1"`
}

// CreateAffiliateAccount creates a new affiliate account
func (h *Handler) CreateAffiliateAccount(c *gin.Context) {
	var req CreateAffiliateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	affiliateType := AffiliateType(req.AffiliateType)
	affiliate, err := h.service.CreateAffiliateAccount(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		userID.(uuid.UUID),
		affiliateType,
		req.CommissionRate,
		req.PayoutThreshold,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, affiliate)
}

// UpdateAffiliateSettings updates affiliate settings
func (h *Handler) UpdateAffiliateSettings(c *gin.Context) {
	affiliateIDStr := c.Param("id")
	affiliateID, err := uuid.Parse(affiliateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid affiliate ID"})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateAffiliateSettings(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		affiliateID,
		settings,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "affiliate settings updated successfully"})
}

// TrackAffiliateClick tracks an affiliate click
func (h *Handler) TrackAffiliateClick(c *gin.Context) {
	var req TrackAffiliateClickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	// Create click data
	clickData := &AffiliateClick{
		IPAddress:   req.IPAddress,
		UserAgent:   req.UserAgent,
		Referrer:    req.Referrer,
		UTMSource:   req.UTMSource,
		UTMMedium:   req.UTMMedium,
		UTMCampaign: req.UTMCampaign,
		UTMTerm:     req.UTMTerm,
		UTMContent:  req.UTMContent,
		DeviceType:  req.DeviceType,
		Country:     req.Country,
		City:        req.City,
		Metadata:    req.Metadata,
	}

	click, err := h.service.TrackAffiliateClick(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		req.ReferralCode,
		clickData,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, click)
}

// GetAffiliateClicks gets affiliate clicks
func (h *Handler) GetAffiliateClicks(c *gin.Context) {
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

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	clicks, err := h.service.GetAffiliateClicks(
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
		"clicks": clicks,
		"limit":  limit,
		"offset": offset,
	})
}

// GetClicksByDateRange gets clicks by date range
func (h *Handler) GetClicksByDateRange(c *gin.Context) {
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

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
		return
	}

	clicks, err := h.service.GetClicksByDateRange(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		userID.(uuid.UUID),
		startDate,
		endDate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clicks": clicks})
}

// GetAffiliatePerformance gets affiliate performance metrics
func (h *Handler) GetAffiliatePerformance(c *gin.Context) {
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

	performance, err := h.service.GetAffiliatePerformance(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		userID.(uuid.UUID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, performance)
}

// GetTopPerformingAffiliates gets top performing affiliates (admin only)
func (h *Handler) GetTopPerformingAffiliates(c *gin.Context) {
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

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	affiliates, err := h.service.GetTopPerformingAffiliates(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"affiliates": affiliates})
}

// CreateOrderCommission creates a commission for an order
func (h *Handler) CreateOrderCommission(c *gin.Context) {
	var req CreateOrderCommissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	// Parse referral ID from header or query param
	referralIDStr := c.GetHeader("X-Referral-ID")
	if referralIDStr == "" {
		referralIDStr = c.Query("referral_id")
	}

	if referralIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "referral_id is required"})
		return
	}

	referralID, err := uuid.Parse(referralIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid referral_id"})
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
		return
	}

	commission, err := h.service.CreateOrderCommission(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		referralID,
		orderID,
		req.OrderAmount,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, commission)
}

// CreatePayoutBatch creates a payout batch (admin only)
func (h *Handler) CreatePayoutBatch(c *gin.Context) {
	isAdmin, exists := c.Get("is_admin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	var req CreatePayoutBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	// Convert string IDs to UUIDs
	var affiliateIDs []uuid.UUID
	for _, idStr := range req.AffiliateIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid affiliate ID: " + idStr})
			return
		}
		affiliateIDs = append(affiliateIDs, id)
	}

	batch, err := h.service.CreatePayoutBatch(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		affiliateIDs,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, batch)
}

// ProcessPayoutBatch processes a payout batch (admin only)
func (h *Handler) ProcessPayoutBatch(c *gin.Context) {
	isAdmin, exists := c.Get("is_admin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	batchIDStr := c.Param("id")
	batchID, err := uuid.Parse(batchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch ID"})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	err = h.service.ProcessPayoutBatch(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		batchID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payout batch processed successfully"})
}

// GetPayoutBatches gets payout batches (admin only)
func (h *Handler) GetPayoutBatches(c *gin.Context) {
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

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	batches, err := h.service.GetPayoutBatches(
		c.Request.Context(),
		tenantID.(uuid.UUID),
		limit,
		offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"batches": batches,
		"limit":   limit,
		"offset":  offset,
	})
}