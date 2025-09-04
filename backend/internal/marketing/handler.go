package marketing

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler defines the marketing HTTP handlers
type Handler struct {
	service Service
}

// NewHandler creates a new marketing handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all marketing routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	marketing := router.Group("/marketing")
	{
		// Campaign routes
		campaigns := marketing.Group("/campaigns")
		{
			campaigns.POST("", h.createCampaign)
			campaigns.GET("", h.getCampaigns)
			campaigns.GET("/:id", h.getCampaign)
			campaigns.PUT("/:id", h.updateCampaign)
			campaigns.DELETE("/:id", h.deleteCampaign)
			campaigns.POST("/:id/schedule", h.scheduleCampaign)
			campaigns.POST("/:id/start", h.startCampaign)
			campaigns.POST("/:id/pause", h.pauseCampaign)
			campaigns.POST("/:id/stop", h.stopCampaign)
			campaigns.GET("/:id/emails", h.getCampaignEmails)
			campaigns.GET("/:id/stats", h.getCampaignStats)
		}
		
		// Template routes
		templates := marketing.Group("/templates")
		{
			templates.POST("", h.createTemplate)
			templates.GET("", h.getTemplates)
			templates.GET("/:id", h.getTemplate)
			templates.PUT("/:id", h.updateTemplate)
			templates.DELETE("/:id", h.deleteTemplate)
		}
		
		// Segment routes
		segments := marketing.Group("/segments")
		{
			segments.POST("", h.createSegment)
			segments.GET("", h.getSegments)
			segments.GET("/:id", h.getSegment)
			segments.PUT("/:id", h.updateSegment)
			segments.DELETE("/:id", h.deleteSegment)
			segments.POST("/:id/refresh", h.refreshSegment)
		}
		
		// Newsletter routes
		newsletter := marketing.Group("/newsletter")
		{
			newsletter.POST("/subscribe", h.subscribe)
			newsletter.POST("/unsubscribe", h.unsubscribe)
			newsletter.GET("/subscribers", h.getSubscribers)
			newsletter.GET("/subscribers/:email", h.getSubscriber)
		}
		
		// Abandoned cart routes
		abandonedCarts := marketing.Group("/abandoned-carts")
		{
			abandonedCarts.POST("", h.createAbandonedCart)
			abandonedCarts.GET("", h.getAbandonedCarts)
			abandonedCarts.POST("/:id/recover", h.markCartRecovered)
			abandonedCarts.POST("/:id/send-email", h.sendAbandonedCartEmail)
		}
		
		// Settings routes
		marketing.GET("/settings", h.getSettings)
		marketing.PUT("/settings", h.updateSettings)
		
		// Analytics routes
		marketing.GET("/overview", h.getMarketingOverview)
		
		// Tracking routes (for email opens/clicks)
		tracking := marketing.Group("/track")
		{
			tracking.GET("/open/:emailId", h.trackEmailOpen)
			tracking.GET("/click/:emailId", h.trackEmailClick)
		}
	}
}

// Campaign handlers
func (h *Handler) createCampaign(c *gin.Context) {
	var req CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	campaign, err := h.service.CreateCampaign(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": campaign})
}

func (h *Handler) getCampaigns(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	var filter CampaignFilter
	
	// Parse query parameters
	if types := c.QueryArray("type"); len(types) > 0 {
		for _, t := range types {
			filter.Type = append(filter.Type, CampaignType(t))
		}
	}
	
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		for _, s := range statuses {
			filter.Status = append(filter.Status, CampaignStatus(s))
		}
	}
	
	filter.Search = c.Query("search")
	
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}
	
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}
	
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	campaigns, err := h.service.GetCampaigns(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": campaigns})
}

func (h *Handler) getCampaign(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	campaign, err := h.service.GetCampaign(c.Request.Context(), tenantID, campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": campaign})
}

func (h *Handler) updateCampaign(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	
	var req UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	campaign, err := h.service.UpdateCampaign(c.Request.Context(), tenantID, campaignID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": campaign})
}

func (h *Handler) deleteCampaign(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.DeleteCampaign(c.Request.Context(), tenantID, campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}

func (h *Handler) scheduleCampaign(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	
	var req struct {
		ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.ScheduleCampaign(c.Request.Context(), tenantID, campaignID, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Campaign scheduled successfully"})
}

func (h *Handler) startCampaign(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.StartCampaign(c.Request.Context(), tenantID, campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Campaign started successfully"})
}

func (h *Handler) pauseCampaign(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.PauseCampaign(c.Request.Context(), tenantID, campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Campaign paused successfully"})
}

func (h *Handler) stopCampaign(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.StopCampaign(c.Request.Context(), tenantID, campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Campaign stopped successfully"})
}

func (h *Handler) getCampaignEmails(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	var filter EmailFilter
	
	// Parse query parameters
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		for _, s := range statuses {
			filter.Status = append(filter.Status, EmailStatus(s))
		}
	}
	
	filter.Search = c.Query("search")
	
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}
	
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}
	
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	emails, err := h.service.GetCampaignEmails(c.Request.Context(), tenantID, campaignID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": emails})
}

func (h *Handler) getCampaignStats(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	stats, err := h.service.GetCampaignStats(c.Request.Context(), tenantID, campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// Template handlers
func (h *Handler) createTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	template, err := h.service.CreateTemplate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": template})
}

func (h *Handler) getTemplates(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	var filter TemplateFilter
	filter.Category = c.Query("category")
	filter.Search = c.Query("search")
	
	if types := c.QueryArray("type"); len(types) > 0 {
		for _, t := range types {
			filter.Type = append(filter.Type, CampaignType(t))
		}
	}
	
	if active := c.Query("active"); active != "" {
		if active == "true" {
			isActive := true
			filter.IsActive = &isActive
		} else if active == "false" {
			isActive := false
			filter.IsActive = &isActive
		}
	}
	
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	templates, err := h.service.GetTemplates(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": templates})
}

func (h *Handler) getTemplate(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	template, err := h.service.GetTemplate(c.Request.Context(), tenantID, templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": template})
}

func (h *Handler) updateTemplate(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}
	
	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	template, err := h.service.UpdateTemplate(c.Request.Context(), tenantID, templateID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": template})
}

func (h *Handler) deleteTemplate(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.DeleteTemplate(c.Request.Context(), tenantID, templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// Segment handlers (similar pattern for brevity)
func (h *Handler) createSegment(c *gin.Context) {
	var req CreateSegmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	segment, err := h.service.CreateSegment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": segment})
}

func (h *Handler) getSegments(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	segments, err := h.service.GetSegments(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": segments})
}

func (h *Handler) getSegment(c *gin.Context) {
	segmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid segment ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	segment, err := h.service.GetSegment(c.Request.Context(), tenantID, segmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": segment})
}

func (h *Handler) updateSegment(c *gin.Context) {
	segmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid segment ID"})
		return
	}
	
	var req UpdateSegmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	segment, err := h.service.UpdateSegment(c.Request.Context(), tenantID, segmentID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": segment})
}

func (h *Handler) deleteSegment(c *gin.Context) {
	segmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid segment ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.DeleteSegment(c.Request.Context(), tenantID, segmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Segment deleted successfully"})
}

func (h *Handler) refreshSegment(c *gin.Context) {
	segmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid segment ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.RefreshSegment(c.Request.Context(), tenantID, segmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Segment refreshed successfully"})
}

// Newsletter handlers
func (h *Handler) subscribe(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	subscriber, err := h.service.Subscribe(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": subscriber})
}

func (h *Handler) unsubscribe(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err := h.service.Unsubscribe(c.Request.Context(), tenantID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

func (h *Handler) getSubscribers(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	var filter SubscriberFilter
	filter.Status = c.QueryArray("status")
	filter.Tags = c.QueryArray("tags")
	filter.Search = c.Query("search")
	
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}
	
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}
	
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	subscribers, err := h.service.GetSubscribers(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": subscribers})
}

func (h *Handler) getSubscriber(c *gin.Context) {
	email := c.Param("email")
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	subscriber, err := h.service.GetSubscriber(c.Request.Context(), tenantID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": subscriber})
}

// Abandoned cart handlers
func (h *Handler) createAbandonedCart(c *gin.Context) {
	var req CreateAbandonedCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	cart, err := h.service.CreateAbandonedCart(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": cart})
}

func (h *Handler) getAbandonedCarts(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	var filter AbandonedCartFilter
	
	if recovered := c.Query("recovered"); recovered != "" {
		if recovered == "true" {
			isRecovered := true
			filter.IsRecovered = &isRecovered
		} else if recovered == "false" {
			isRecovered := false
			filter.IsRecovered = &isRecovered
		}
	}
	
	if minValue := c.Query("min_value"); minValue != "" {
		if val, err := strconv.ParseFloat(minValue, 64); err == nil {
			filter.MinValue = &val
		}
	}
	
	if maxValue := c.Query("max_value"); maxValue != "" {
		if val, err := strconv.ParseFloat(maxValue, 64); err == nil {
			filter.MaxValue = &val
		}
	}
	
	filter.Search = c.Query("search")
	
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}
	
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}
	
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	carts, err := h.service.GetAbandonedCarts(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": carts})
}

func (h *Handler) markCartRecovered(c *gin.Context) {
	cartID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}
	
	var req struct {
		RecoveredValue float64 `json:"recovered_value" binding:"required,min=0"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.MarkCartRecovered(c.Request.Context(), tenantID, cartID, req.RecoveredValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Cart marked as recovered"})
}

func (h *Handler) sendAbandonedCartEmail(c *gin.Context) {
	cartID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.SendAbandonedCartEmail(c.Request.Context(), tenantID, cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Abandoned cart email sent"})
}

// Settings handlers
func (h *Handler) getSettings(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	settings, err := h.service.GetSettings(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": settings})
}

func (h *Handler) updateSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	settings, err := h.service.UpdateSettings(c.Request.Context(), tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": settings})
}

// Analytics handlers
func (h *Handler) getMarketingOverview(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	period := c.DefaultQuery("period", "30d")
	
	overview, err := h.service.GetMarketingOverview(c.Request.Context(), tenantID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": overview})
}

// Tracking handlers
func (h *Handler) trackEmailOpen(c *gin.Context) {
	emailID, err := uuid.Parse(c.Param("emailId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID"})
		return
	}
	
	err = h.service.TrackEmailOpen(c.Request.Context(), emailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Return tracking pixel (1x1 transparent PNG)
	pixel := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
		0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
	
	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Data(http.StatusOK, "image/png", pixel)
}

func (h *Handler) trackEmailClick(c *gin.Context) {
	emailID, err := uuid.Parse(c.Param("emailId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID"})
		return
	}
	
	redirectURL := c.Query("url")
	if redirectURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing redirect URL"})
		return
	}
	
	err = h.service.TrackEmailClick(c.Request.Context(), emailID)
	if err != nil {
		// Log error but don't fail the redirect
		// TODO: Add proper logging
	}
	
	c.Redirect(http.StatusFound, redirectURL)
}