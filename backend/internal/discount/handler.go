package discount

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler defines the discount HTTP handlers
type Handler struct {
	service Service
}

// NewHandler creates a new discount handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all discount routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Discount/Coupon routes
	discounts := router.Group("/discounts")
	{
		discounts.POST("", h.createDiscount)
		discounts.GET("", h.getDiscounts)
		discounts.GET("/:id", h.getDiscount)
		discounts.PUT("/:id", h.updateDiscount)
		discounts.DELETE("/:id", h.deleteDiscount)
		discounts.GET("/:id/usage", h.getDiscountUsage)
		discounts.GET("/stats", h.getDiscountStats)
		discounts.GET("/performance", h.getTopDiscounts)
		discounts.GET("/revenue-impact", h.getDiscountRevenue)
	}
	
	// Public discount validation (for cart/checkout)
	router.POST("/validate-discount", h.validateDiscountCode)
	router.POST("/apply-discount", h.applyDiscount)
	router.DELETE("/remove-discount/:orderId", h.removeDiscount)
	
	// Gift card routes
	giftCards := router.Group("/gift-cards")
	{
		giftCards.POST("", h.createGiftCard)
		giftCards.GET("", h.getGiftCards)
		giftCards.GET("/:code", h.getGiftCard)
		giftCards.PUT("/:id", h.updateGiftCard)
		giftCards.DELETE("/:id", h.deleteGiftCard)
		giftCards.GET("/:id/transactions", h.getGiftCardTransactions)
		giftCards.POST("/:id/refill", h.refillGiftCard)
	}
	
	// Public gift card validation
	router.POST("/validate-gift-card", h.validateGiftCard)
	router.POST("/use-gift-card", h.useGiftCard)
	
	// Store credit routes
	storeCredit := router.Group("/store-credit")
	{
		storeCredit.GET("/:customerId", h.getStoreCredit)
		storeCredit.POST("/:customerId/add", h.addStoreCredit)
		storeCredit.POST("/:customerId/use", h.useStoreCredit)
		storeCredit.GET("/:customerId/transactions", h.getStoreCreditTransactions)
	}
}

// Discount handlers
func (h *Handler) createDiscount(c *gin.Context) {
	var req CreateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	discount, err := h.service.CreateDiscount(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": discount})
}

func (h *Handler) getDiscounts(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	filter := h.parseDiscountFilter(c)
	
	discounts, err := h.service.GetDiscounts(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": discounts})
}

func (h *Handler) getDiscount(c *gin.Context) {
	discountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	discount, err := h.service.GetDiscount(c.Request.Context(), tenantID, discountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": discount})
}

func (h *Handler) updateDiscount(c *gin.Context) {
	discountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount ID"})
		return
	}
	
	var req UpdateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	discount, err := h.service.UpdateDiscount(c.Request.Context(), tenantID, discountID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": discount})
}

func (h *Handler) deleteDiscount(c *gin.Context) {
	discountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.DeleteDiscount(c.Request.Context(), tenantID, discountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Discount deleted successfully"})
}

func (h *Handler) getDiscountUsage(c *gin.Context) {
	discountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	filter := h.parseUsageFilter(c)
	
	usage, err := h.service.GetDiscountUsage(c.Request.Context(), tenantID, discountID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": usage})
}

func (h *Handler) validateDiscountCode(c *gin.Context) {
	var req ValidateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	validation, err := h.service.ValidateDiscountCode(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": validation})
}

func (h *Handler) applyDiscount(c *gin.Context) {
	var req ApplyDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context and add client info
	req.TenantID = uuid.New() // Placeholder
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")
	
	application, err := h.service.ApplyDiscount(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": application})
}

func (h *Handler) removeDiscount(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.RemoveDiscount(c.Request.Context(), tenantID, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Discount removed successfully"})
}

// Analytics handlers
func (h *Handler) getDiscountStats(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	period := c.DefaultQuery("period", "30d")
	
	stats, err := h.service.GetDiscountStats(c.Request.Context(), tenantID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handler) getTopDiscounts(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	discounts, err := h.service.GetTopDiscounts(c.Request.Context(), tenantID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": discounts})
}

func (h *Handler) getDiscountRevenue(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	period := c.DefaultQuery("period", "30d")
	
	revenue, err := h.service.GetDiscountRevenue(c.Request.Context(), tenantID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": revenue})
}

// Gift card handlers
func (h *Handler) createGiftCard(c *gin.Context) {
	var req CreateGiftCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	giftCard, err := h.service.CreateGiftCard(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": giftCard})
}

func (h *Handler) getGiftCards(c *gin.Context) {
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	filter := h.parseGiftCardFilter(c)
	
	giftCards, err := h.service.GetGiftCards(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": giftCards})
}

func (h *Handler) getGiftCard(c *gin.Context) {
	code := c.Param("code")
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	giftCard, err := h.service.GetGiftCard(c.Request.Context(), tenantID, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": giftCard})
}

func (h *Handler) updateGiftCard(c *gin.Context) {
	giftCardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gift card ID"})
		return
	}
	
	var req UpdateGiftCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	giftCard, err := h.service.UpdateGiftCard(c.Request.Context(), tenantID, giftCardID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": giftCard})
}

func (h *Handler) deleteGiftCard(c *gin.Context) {
	giftCardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gift card ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	err = h.service.DeleteGiftCard(c.Request.Context(), tenantID, giftCardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Gift card deleted successfully"})
}

func (h *Handler) getGiftCardTransactions(c *gin.Context) {
	giftCardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gift card ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	transactions, err := h.service.GetGiftCardTransactions(c.Request.Context(), tenantID, giftCardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

func (h *Handler) refillGiftCard(c *gin.Context) {
	giftCardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gift card ID"})
		return
	}
	
	var req RefillGiftCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID and user ID from context
	req.TenantID = uuid.New()   // Placeholder
	req.GiftCardID = giftCardID
	req.ProcessedBy = uuid.New() // Placeholder
	
	transaction, err := h.service.RefillGiftCard(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

func (h *Handler) validateGiftCard(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	validation, err := h.service.ValidateGiftCard(c.Request.Context(), tenantID, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": validation})
}

func (h *Handler) useGiftCard(c *gin.Context) {
	var req UseGiftCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	req.TenantID = uuid.New() // Placeholder
	
	transaction, err := h.service.UseGiftCard(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// Store credit handlers
func (h *Handler) getStoreCredit(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	storeCredit, err := h.service.GetStoreCredit(c.Request.Context(), tenantID, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": storeCredit})
}

func (h *Handler) addStoreCredit(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}
	
	var req AddStoreCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID and user ID from context
	req.TenantID = uuid.New()   // Placeholder
	req.CustomerID = customerID
	req.ProcessedBy = &uuid.UUID{} // Placeholder
	
	transaction, err := h.service.AddStoreCredit(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": transaction})
}

func (h *Handler) useStoreCredit(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}
	
	var req UseStoreCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Get tenant ID from context
	req.TenantID = uuid.New() // Placeholder
	req.CustomerID = customerID
	
	transaction, err := h.service.UseStoreCredit(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

func (h *Handler) getStoreCreditTransactions(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}
	
	// TODO: Get tenant ID from context
	tenantID := uuid.New() // Placeholder
	
	filter := h.parseStoreCreditFilter(c)
	
	transactions, err := h.service.GetStoreCreditTransactions(c.Request.Context(), tenantID, customerID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

// Helper functions
func (h *Handler) parseDiscountFilter(c *gin.Context) DiscountFilter {
	filter := DiscountFilter{}
	
	// Parse status array
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		for _, s := range statuses {
			filter.Status = append(filter.Status, DiscountStatus(s))
		}
	}
	
	// Parse type array
	if types := c.QueryArray("type"); len(types) > 0 {
		for _, t := range types {
			filter.Type = append(filter.Type, DiscountType(t))
		}
	}
	
	// Parse target array
	if targets := c.QueryArray("target"); len(targets) > 0 {
		for _, t := range targets {
			filter.Target = append(filter.Target, DiscountTarget(t))
		}
	}
	
	filter.Search = c.Query("search")
	
	// Parse boolean flags
	if expired := c.Query("expired"); expired != "" {
		isExpired := expired == "true"
		filter.IsExpired = &isExpired
	}
	
	if active := c.Query("active"); active != "" {
		isActive := active == "true"
		filter.IsActive = &isActive
	}
	
	// Parse creator ID
	if createdBy := c.Query("created_by"); createdBy != "" {
		if id, err := uuid.Parse(createdBy); err == nil {
			filter.CreatedBy = &id
		}
	}
	
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

func (h *Handler) parseUsageFilter(c *gin.Context) UsageFilter {
	filter := UsageFilter{}
	
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
	
	// Parse pagination
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	return filter
}

func (h *Handler) parseGiftCardFilter(c *gin.Context) GiftCardFilter {
	filter := GiftCardFilter{}
	
	// Parse status array
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		for _, s := range statuses {
			filter.Status = append(filter.Status, DiscountStatus(s))
		}
	}
	
	filter.Search = c.Query("search")
	
	// Parse boolean flags
	if expired := c.Query("expired"); expired != "" {
		isExpired := expired == "true"
		filter.IsExpired = &isExpired
	}
	
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
	
	// Parse pagination
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	return filter
}

func (h *Handler) parseStoreCreditFilter(c *gin.Context) StoreCreditFilter {
	filter := StoreCreditFilter{}
	
	// Parse type array
	filter.Type = c.QueryArray("type")
	
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
	
	// Parse pagination
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}
	
	return filter
}