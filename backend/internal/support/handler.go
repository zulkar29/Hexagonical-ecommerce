package support

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler defines the support HTTP handlers
type Handler struct {
	service Service
}

// NewHandler creates a new support handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all support routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	support := router.Group("/support")
	{
		// Ticket routes
		tickets := support.Group("/tickets")
		{
			tickets.POST("", h.createTicket)
			tickets.GET("", h.getTickets)
			tickets.GET("/:id", h.getTicket)
			tickets.PUT("/:id", h.updateTicket)
			tickets.DELETE("/:id", h.deleteTicket)
			tickets.POST("/:id/assign", h.assignTicket)
			tickets.POST("/:id/resolve", h.resolveTicket)
			tickets.GET("/:id/messages", h.getTicketMessages)
			tickets.POST("/:id/messages", h.addTicketMessage)
		}
		
		// FAQ routes
		faqs := support.Group("/faqs")
		{
			faqs.POST("", h.createFAQ)
			faqs.GET("", h.getFAQs)
			faqs.GET("/:id", h.getFAQ)
			faqs.PUT("/:id", h.updateFAQ)
			faqs.DELETE("/:id", h.deleteFAQ)
		}
		
		// Knowledge base routes
		kb := support.Group("/knowledge-base")
		{
			kb.POST("", h.createArticle)
			kb.GET("", h.getArticles)
			kb.GET("/:slug", h.getArticle)
			kb.PUT("/:id", h.updateArticle)
			kb.DELETE("/:id", h.deleteArticle)
		}
		
		// Settings routes
		support.GET("/settings", h.getSettings)
		support.PUT("/settings", h.updateSettings)
		
		// Analytics routes
		support.GET("/stats", h.getTicketStats)
	}
}

// Ticket handlers
func (h *Handler) createTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	// tenantID := c.GetString("tenant_id")
	
	ticket, err := h.service.CreateTicket(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": ticket})
}

func (h *Handler) getTickets(c *gin.Context) {
	// TODO: Get tenant ID from context
	// tenantID, _ := uuid.Parse(c.GetString("tenant_id"))
	tenantID := uuid.New() // Placeholder
	
	var filter TicketFilter
	
	// Parse query parameters
	if status := c.QueryArray("status"); len(status) > 0 {
		for _, s := range status {
			filter.Status = append(filter.Status, TicketStatus(s))
		}
	}
	
	if priority := c.QueryArray("priority"); len(priority) > 0 {
		for _, p := range priority {
			filter.Priority = append(filter.Priority, TicketPriority(p))
		}
	}
	
	if category := c.QueryArray("category"); len(category) > 0 {
		for _, cat := range category {
			filter.Category = append(filter.Category, TicketCategory(cat))
		}
	}
	
	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		if assignedTo == "unassigned" {
			nilUUID := uuid.Nil
			filter.AssignedToID = &nilUUID
		} else if id, err := uuid.Parse(assignedTo); err == nil {
			filter.AssignedToID = &id
		}
	}
	
	filter.Search = c.Query("search")
	
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	tickets, err := h.service.GetTickets(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func (h *Handler) getTicket(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	ticket, err := h.service.GetTicket(c.Request.Context(), tenantID, ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": ticket})
}

func (h *Handler) updateTicket(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	
	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	ticket, err := h.service.UpdateTicket(c.Request.Context(), tenantID, ticketID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": ticket})
}

func (h *Handler) deleteTicket(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.DeleteTicket(c.Request.Context(), tenantID, ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

func (h *Handler) assignTicket(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	
	var req struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.AssignTicket(c.Request.Context(), tenantID, ticketID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Ticket assigned successfully"})
}

func (h *Handler) resolveTicket(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.ResolveTicket(c.Request.Context(), tenantID, ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Ticket resolved successfully"})
}

func (h *Handler) getTicketMessages(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	messages, err := h.service.GetMessages(c.Request.Context(), tenantID, ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": messages})
}

func (h *Handler) addTicketMessage(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	
	var req AddMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	req.TicketID = ticketID
	
	message, err := h.service.AddMessage(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": message})
}

// FAQ handlers
func (h *Handler) createFAQ(c *gin.Context) {
	var req CreateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	faq, err := h.service.CreateFAQ(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": faq})
}

func (h *Handler) getFAQs(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	var filter FAQFilter
	filter.Category = c.Query("category")
	filter.Search = c.Query("search")
	
	if published := c.Query("published"); published != "" {
		if published == "true" {
			isPublished := true
			filter.IsPublished = &isPublished
		} else if published == "false" {
			isPublished := false
			filter.IsPublished = &isPublished
		}
	}
	
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	faqs, err := h.service.GetFAQs(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": faqs})
}

func (h *Handler) getFAQ(c *gin.Context) {
	faqID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid FAQ ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	faq, err := h.service.GetFAQ(c.Request.Context(), tenantID, faqID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": faq})
}

func (h *Handler) updateFAQ(c *gin.Context) {
	faqID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid FAQ ID"})
		return
	}
	
	var req UpdateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	faq, err := h.service.UpdateFAQ(c.Request.Context(), tenantID, faqID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": faq})
}

func (h *Handler) deleteFAQ(c *gin.Context) {
	faqID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid FAQ ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.DeleteFAQ(c.Request.Context(), tenantID, faqID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "FAQ deleted successfully"})
}

// Knowledge base handlers
func (h *Handler) createArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	article, err := h.service.CreateArticle(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": article})
}

func (h *Handler) getArticles(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	var filter ArticleFilter
	filter.Category = c.Query("category")
	filter.Search = c.Query("search")
	
	if published := c.Query("published"); published != "" {
		if published == "true" {
			isPublished := true
			filter.IsPublished = &isPublished
		} else if published == "false" {
			isPublished := false
			filter.IsPublished = &isPublished
		}
	}
	
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	articles, err := h.service.GetArticles(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": articles})
}

func (h *Handler) getArticle(c *gin.Context) {
	slug := c.Param("slug")
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	article, err := h.service.GetArticle(c.Request.Context(), tenantID, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": article})
}

func (h *Handler) updateArticle(c *gin.Context) {
	articleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}
	
	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	article, err := h.service.UpdateArticle(c.Request.Context(), tenantID, articleID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": article})
}

func (h *Handler) deleteArticle(c *gin.Context) {
	articleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.DeleteArticle(c.Request.Context(), tenantID, articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
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
func (h *Handler) getTicketStats(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	period := c.DefaultQuery("period", "30d")
	
	stats, err := h.service.GetTicketStats(c.Request.Context(), tenantID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": stats})
}