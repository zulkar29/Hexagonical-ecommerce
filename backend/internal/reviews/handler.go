package reviews

import (
	"net/http"
	"strconv"
	"time"

	"ecommerce-saas/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler defines the reviews HTTP handlers
type Handler struct {
	service Service
}

// NewHandler creates a new reviews handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all review routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// 📍 CORE REVIEW ENDPOINTS (5)
	reviews := router.Group("/reviews")
	{
		reviews.POST("", h.createReview) // CreateReview
		reviews.GET("", h.getReviews)    // GetReviews (with filtering, stats, trends, recent)
		// 📍 SETTINGS (2) - Must come before /:id routes
		reviews.GET("/settings", h.getSettings)    // GetSettings
		reviews.PUT("/settings", h.updateSettings) // UpdateSettings
		reviews.GET("/:id", h.getReview)           // GetReview
		reviews.PUT("/:id", h.updateReview)        // UpdateReview (handles moderation actions)
		reviews.DELETE("/:id", h.deleteReview)     // DeleteReview
	}

	// 📍 REVIEW OPERATIONS (4)
	reviews.POST("/bulk", h.bulkModerateReviews)       // BulkModerateReviews
	reviews.GET("/:id/replies", h.getReplies)          // GetReplies
	reviews.POST("/:id/replies", h.addReply)           // AddReply
	reviews.PUT("/replies/:replyId", h.updateReply)    // UpdateReply
	reviews.DELETE("/replies/:replyId", h.deleteReply) // DeleteReply
	reviews.POST("/:id/react", h.reactToReview)        // ReactToReview
	reviews.DELETE("/:id/react", h.removeReaction)     // RemoveReaction

	// 📍 PRODUCT REVIEWS (3) - Use reviews prefix to avoid conflicts
	reviews.GET("/product/:productId", h.getProductReviews)                            // GetProductReviews
	reviews.GET("/product/:productId/summary", h.getProductReviewSummary)              // GetProductReviewSummary
	reviews.POST("/product/:productId/summary/refresh", h.refreshProductReviewSummary) // RefreshProductReviewSummary

	// 📍 REVIEW INVITATIONS (4)
	invitations := router.Group("/review-invitations")
	{
		invitations.POST("", h.createReviewInvitation)       // CreateReviewInvitation
		invitations.GET("", h.getReviewInvitations)          // GetReviewInvitations (with pending filter)
		invitations.PUT("/:id", h.updateReviewInvitation)    // UpdateReviewInvitation (handles send/remind)
		invitations.DELETE("/:id", h.deleteReviewInvitation) // DeleteReviewInvitation
	}

	// 📍 PUBLIC ENDPOINTS (1)
	router.GET("/review-invite/:token", h.processInvitationClick) // ProcessInvitationClick

}

// Review CRUD handlers
func (h *Handler) createReview(c *gin.Context) {
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	review, err := h.service.CreateReview(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": review})
}

func (h *Handler) getReviews(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	// Check for consolidated operations via query parameters
	operationType := c.Query("type")

	switch operationType {
	case "stats":
		h.handleGetStats(c, tenantID)
		return
	case "trends":
		h.handleGetTrends(c, tenantID)
		return
	case "top-products":
		h.handleGetTopProducts(c, tenantID)
		return
	case "recent":
		h.handleGetRecent(c, tenantID)
		return
	case "pending":
		h.handleGetPending(c, tenantID)
		return
	default:
		// Standard review listing with filtering
		filter := h.parseReviewFilter(c)

		reviews, err := h.service.GetReviews(c.Request.Context(), tenantID, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": reviews})
	}
}

func (h *Handler) getReview(c *gin.Context) {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	review, err := h.service.GetReview(c.Request.Context(), tenantID, reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": review})
}

func (h *Handler) updateReview(c *gin.Context) {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	// Check for moderation actions via query parameters
	action := c.Query("action")

	switch action {
	case "approve":
		h.handleApproveReview(c, tenantID, reviewID)
		return
	case "reject":
		h.handleRejectReview(c, tenantID, reviewID)
		return
	case "spam":
		h.handleMarkAsSpam(c, tenantID, reviewID)
		return
	default:
		// Standard review update
		var req UpdateReviewRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		review, err := h.service.UpdateReview(c.Request.Context(), tenantID, reviewID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": review})
	}
}

func (h *Handler) deleteReview(c *gin.Context) {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	err = h.service.DeleteReview(c.Request.Context(), tenantID, reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

// Moderation handlers

func (h *Handler) bulkModerateReviews(c *gin.Context) {
	var req BulkModerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	moderatorID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	req.ModeratorID = moderatorID // Set moderator ID from context

	err = h.service.BulkModerateReviews(c.Request.Context(), tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bulk moderation completed successfully"})
}

// Reply handlers
func (h *Handler) addReply(c *gin.Context) {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req AddReplyRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	req.ReviewID = reviewID

	reply, err := h.service.AddReply(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": reply})
}

func (h *Handler) getReplies(c *gin.Context) {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	replies, err := h.service.GetReplies(c.Request.Context(), tenantID, reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": replies})
}

func (h *Handler) updateReply(c *gin.Context) {
	replyID, err := uuid.Parse(c.Param("replyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reply ID"})
		return
	}

	var req UpdateReplyRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	reply, err := h.service.UpdateReply(c.Request.Context(), tenantID, replyID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reply})
}

func (h *Handler) deleteReply(c *gin.Context) {
	replyID, err := uuid.Parse(c.Param("replyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reply ID"})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	err = h.service.DeleteReply(c.Request.Context(), tenantID, replyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reply deleted successfully"})
}

// Reaction handlers
func (h *Handler) reactToReview(c *gin.Context) {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req ReviewReactionRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	req.ReviewID = reviewID
	req.IPAddress = c.ClientIP()

	err = h.service.ReactToReview(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaction added successfully"})
}

func (h *Handler) removeReaction(c *gin.Context) {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	customerEmail := c.Query("email")
	if customerEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer email is required"})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	err = h.service.RemoveReaction(c.Request.Context(), tenantID, reviewID, customerEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaction removed successfully"})
}

// Product review handlers
func (h *Handler) getProductReviews(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	filter := h.parseProductReviewFilter(c)

	reviews, err := h.service.GetProductReviews(c.Request.Context(), tenantID, productID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reviews})
}

func (h *Handler) getProductReviewSummary(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	summary, err := h.service.GetReviewSummary(c.Request.Context(), tenantID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func (h *Handler) refreshProductReviewSummary(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	summary, err := h.service.RefreshReviewSummary(c.Request.Context(), tenantID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}

// Invitation handlers
func (h *Handler) createReviewInvitation(c *gin.Context) {
	var req CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invitation, err := h.service.CreateReviewInvitation(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": invitation})
}

func (h *Handler) getReviewInvitations(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	status := c.DefaultQuery("status", "")

	var invitations []ReviewInvitation

	if status != "" {
		invitations, err = h.service.GetPendingInvitations(c.Request.Context(), tenantID)
	} else {
		invitations, err = h.service.GetPendingInvitations(c.Request.Context(), tenantID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": invitations})
}

func (h *Handler) updateReviewInvitation(c *gin.Context) {
	invitationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation ID"})
		return
	}

	// Check for action query parameter
	action := c.Query("action")
	if action != "" {
		// Handle action-based updates (send/remind)
		req := UpdateInvitationRequest{Action: action}

		tenantID, err := utils.GetTenantIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
			return
		}

		invitation, updateErr := h.service.UpdateReviewInvitation(c.Request.Context(), tenantID, invitationID, req)
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": invitation, "message": "Action completed successfully"})
		return
	}

	// Handle field updates
	var req UpdateInvitationRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	invitation, err := h.service.UpdateReviewInvitation(c.Request.Context(), tenantID, invitationID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": invitation})
}

func (h *Handler) deleteReviewInvitation(c *gin.Context) {
	invitationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation ID"})
		return
	}

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	err = h.service.DeleteReviewInvitation(c.Request.Context(), tenantID, invitationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation deleted successfully"})
}

func (h *Handler) processInvitationClick(c *gin.Context) {
	token := c.Param("token")

	invitation, err := h.service.ProcessInvitationClick(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Redirect to review form or return invitation data
	c.JSON(http.StatusOK, gin.H{"data": invitation})
}

// Analytics handlers

// Settings handlers
func (h *Handler) getSettings(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

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

	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid tenant context"})
		return
	}

	settings, err := h.service.UpdateSettings(c.Request.Context(), tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": settings})
}

// Consolidated operation helper functions
func (h *Handler) handleGetStats(c *gin.Context, tenantID uuid.UUID) {
	period := c.DefaultQuery("period", "30d")

	stats, err := h.service.GetReviewStats(c.Request.Context(), tenantID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handler) handleGetTrends(c *gin.Context, tenantID uuid.UUID) {
	period := c.DefaultQuery("period", "30d")

	trends, err := h.service.GetReviewTrends(c.Request.Context(), tenantID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trends})
}

func (h *Handler) handleGetTopProducts(c *gin.Context, tenantID uuid.UUID) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.service.GetTopRatedProducts(c.Request.Context(), tenantID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *Handler) handleGetRecent(c *gin.Context, tenantID uuid.UUID) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, err := h.service.GetRecentReviews(c.Request.Context(), tenantID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reviews})
}

func (h *Handler) handleGetPending(c *gin.Context, tenantID uuid.UUID) {
	filter := ReviewFilter{
		Status:    []ReviewStatus{StatusPending},
		SortBy:    "created_at",
		SortOrder: "asc",
		Page:      1,
		Limit:     50,
	}

	// Parse pagination parameters
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}

	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "50")); err == nil {
		filter.Limit = limit
	}

	reviews, err := h.service.GetReviews(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reviews})
}

func (h *Handler) handleApproveReview(c *gin.Context, tenantID, reviewID uuid.UUID) {
	moderatorID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	err = h.service.ApproveReview(c.Request.Context(), tenantID, reviewID, moderatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review approved successfully"})
}

func (h *Handler) handleRejectReview(c *gin.Context, tenantID, reviewID uuid.UUID) {
	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	moderatorID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	err = h.service.RejectReview(c.Request.Context(), tenantID, reviewID, moderatorID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review rejected successfully"})
}

func (h *Handler) handleMarkAsSpam(c *gin.Context, tenantID, reviewID uuid.UUID) {
	moderatorID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	err = h.service.MarkAsSpam(c.Request.Context(), tenantID, reviewID, moderatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review marked as spam successfully"})
}

// Helper functions
func (h *Handler) parseReviewFilter(c *gin.Context) ReviewFilter {
	filter := ReviewFilter{}

	// Parse product ID
	if productID := c.Query("product_id"); productID != "" {
		if id, err := uuid.Parse(productID); err == nil {
			filter.ProductID = &id
		}
	}

	// Parse order ID
	if orderID := c.Query("order_id"); orderID != "" {
		if id, err := uuid.Parse(orderID); err == nil {
			filter.OrderID = &id
		}
	}

	// Parse user ID
	if userID := c.Query("user_id"); userID != "" {
		if id, err := uuid.Parse(userID); err == nil {
			filter.UserID = &id
		}
	}

	// Parse status array
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		for _, s := range statuses {
			filter.Status = append(filter.Status, ReviewStatus(s))
		}
	}

	// Parse ratings
	if ratings := c.QueryArray("rating"); len(ratings) > 0 {
		for _, r := range ratings {
			if rating, err := strconv.Atoi(r); err == nil && rating >= 1 && rating <= 5 {
				filter.Rating = append(filter.Rating, rating)
			}
		}
	}

	// Parse boolean flags
	if verified := c.Query("verified"); verified != "" {
		isVerified := verified == "true"
		filter.IsVerified = &isVerified
	}

	if hasImages := c.Query("has_images"); hasImages != "" {
		hasImagesFlag := hasImages == "true"
		filter.HasImages = &hasImagesFlag
	}

	if hasVideos := c.Query("has_videos"); hasVideos != "" {
		hasVideosFlag := hasVideos == "true"
		filter.HasVideos = &hasVideosFlag
	}

	// Parse search
	filter.Search = c.Query("search")

	// Parse dates
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

	// Parse sorting
	filter.SortBy = c.DefaultQuery("sort_by", "created_at")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	// Parse pagination
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}

	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}

	return filter
}

func (h *Handler) parseProductReviewFilter(c *gin.Context) ProductReviewFilter {
	filter := ProductReviewFilter{}

	// Parse status array
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		for _, s := range statuses {
			filter.Status = append(filter.Status, ReviewStatus(s))
		}
	} else {
		// Default to approved reviews for public API
		filter.Status = []ReviewStatus{StatusApproved}
	}

	// Parse ratings
	if ratings := c.QueryArray("rating"); len(ratings) > 0 {
		for _, r := range ratings {
			if rating, err := strconv.Atoi(r); err == nil && rating >= 1 && rating <= 5 {
				filter.Rating = append(filter.Rating, rating)
			}
		}
	}

	// Parse boolean flags
	if verified := c.Query("verified"); verified != "" {
		isVerified := verified == "true"
		filter.IsVerified = &isVerified
	}

	if hasImages := c.Query("has_images"); hasImages != "" {
		hasImagesFlag := hasImages == "true"
		filter.HasImages = &hasImagesFlag
	}

	if hasVideos := c.Query("has_videos"); hasVideos != "" {
		hasVideosFlag := hasVideos == "true"
		filter.HasVideos = &hasVideosFlag
	}

	// Parse search
	filter.Search = c.Query("search")

	// Parse sorting
	filter.SortBy = c.DefaultQuery("sort_by", "created_at")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	// Parse pagination
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}

	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}

	return filter
}
