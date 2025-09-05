package loyalty

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for loyalty operations
type Handler struct {
	service Service
}

// NewHandler creates a new loyalty handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Loyalty Programs

// CreateProgram creates a new loyalty program
func (h *Handler) CreateProgram(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req CreateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.TenantID = tenantID.(uuid.UUID)

	program, err := h.service.CreateProgram(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, program)
}

// GetProgram retrieves a loyalty program by ID
func (h *Handler) GetProgram(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	program, err := h.service.GetProgram(c.Request.Context(), tenantID.(uuid.UUID), programID)
	if err != nil {
		if err.Error() == "loyalty program not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, program)
}

// UpdateProgram updates a loyalty program
func (h *Handler) UpdateProgram(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	var req UpdateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	program, err := h.service.UpdateProgram(c.Request.Context(), tenantID.(uuid.UUID), programID, &req)
	if err != nil {
		if err.Error() == "loyalty program not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, program)
}

// DeleteProgram deletes a loyalty program
func (h *Handler) DeleteProgram(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	err = h.service.DeleteProgram(c.Request.Context(), tenantID.(uuid.UUID), programID)
	if err != nil {
		if err.Error() == "loyalty program not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListPrograms lists loyalty programs with optional filtering
func (h *Handler) ListPrograms(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	// Check if this is a stats request
	if c.Query("type") == "stats" {
		h.GetStats(c)
		return
	}

	var req ListProgramsRequest

	// Parse query parameters
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}
	if programType := c.Query("type"); programType != "" {
		req.Type = &programType
	}
	if search := c.Query("search"); search != "" {
		req.Search = &search
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = &limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = &offset
		}
	}

	response, err := h.service.ListPrograms(c.Request.Context(), tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Loyalty Accounts

// CreateAccount creates a new loyalty account
func (h *Handler) CreateAccount(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.TenantID = tenantID.(uuid.UUID)

	account, err := h.service.CreateAccount(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccount retrieves a loyalty account by ID
func (h *Handler) GetAccount(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := h.service.GetAccount(c.Request.Context(), tenantID.(uuid.UUID), accountID)
	if err != nil {
		if err.Error() == "loyalty account not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// UpdateAccount updates a loyalty account
func (h *Handler) UpdateAccount(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.service.UpdateAccount(c.Request.Context(), tenantID.(uuid.UUID), accountID, &req)
	if err != nil {
		if err.Error() == "loyalty account not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// ListAccounts lists loyalty accounts with optional filtering
func (h *Handler) ListAccounts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req ListAccountsRequest

	// Parse query parameters
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}
	if tier := c.Query("tier"); tier != "" {
		req.Tier = &tier
	}
	if programIDStr := c.Query("program_id"); programIDStr != "" {
		if programID, err := uuid.Parse(programIDStr); err == nil {
			req.ProgramID = &programID
		}
	}
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			req.UserID = &userID
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	response, err := h.service.ListAccounts(c.Request.Context(), tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Points Operations

// EarnPoints adds points to a loyalty account
func (h *Handler) EarnPoints(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var reqData struct {
		Points  int64      `json:"points" binding:"required,min=1"`
		Reason  string     `json:"reason" binding:"required"`
		OrderID *uuid.UUID `json:"order_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &EarnPointsRequest{
		TenantID:  tenantID.(uuid.UUID),
		AccountID: accountID,
		Points:    reqData.Points,
		Reason:    reqData.Reason,
		OrderID:   reqData.OrderID,
	}

	transaction, err := h.service.EarnPoints(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// RedeemPoints deducts points from a loyalty account
func (h *Handler) RedeemPoints(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var reqData struct {
		Points  int64      `json:"points" binding:"required,min=1"`
		Reason  string     `json:"reason" binding:"required"`
		OrderID *uuid.UUID `json:"order_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &RedeemPointsRequest{
		TenantID:  tenantID.(uuid.UUID),
		AccountID: accountID,
		Points:    reqData.Points,
		Reason:    reqData.Reason,
		OrderID:   reqData.OrderID,
	}

	transaction, err := h.service.RedeemPoints(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// RedeemReward redeems a specific reward
func (h *Handler) RedeemReward(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var req struct {
		RewardID uuid.UUID `json:"reward_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.service.RedeemReward(c.Request.Context(), tenantID.(uuid.UUID), accountID, req.RewardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// Loyalty Transactions

// GetTransaction retrieves a loyalty transaction by ID
func (h *Handler) GetTransaction(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.service.GetTransaction(c.Request.Context(), tenantID.(uuid.UUID), transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// ListTransactions lists loyalty transactions with optional filtering
func (h *Handler) ListTransactions(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req ListTransactionsRequest

	// Parse query parameters
	if transactionType := c.Query("type"); transactionType != "" {
		req.Type = &transactionType
	}
	if accountIDStr := c.Query("account_id"); accountIDStr != "" {
		if accountID, err := uuid.Parse(accountIDStr); err == nil {
			req.AccountID = &accountID
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = &limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = &offset
		}
	}

	response, err := h.service.ListTransactions(c.Request.Context(), tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListAccountTransactions lists transactions for a specific account
func (h *Handler) ListAccountTransactions(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var req ListAccountTransactionsRequest

	// Parse query parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = &limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = &offset
		}
	}

	response, err := h.service.ListAccountTransactions(c.Request.Context(), tenantID.(uuid.UUID), accountID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Loyalty Rewards

// CreateReward creates a new loyalty reward
func (h *Handler) CreateReward(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.TenantID = tenantID.(uuid.UUID)

	reward, err := h.service.CreateReward(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reward)
}

// GetReward retrieves a loyalty reward by ID
func (h *Handler) GetReward(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	rewardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reward ID"})
		return
	}

	reward, err := h.service.GetReward(c.Request.Context(), tenantID.(uuid.UUID), rewardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reward)
}

// UpdateReward updates a loyalty reward
func (h *Handler) UpdateReward(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	rewardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reward ID"})
		return
	}

	var req UpdateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reward, err := h.service.UpdateReward(c.Request.Context(), tenantID.(uuid.UUID), rewardID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reward)
}

// DeleteReward deletes a loyalty reward
func (h *Handler) DeleteReward(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	rewardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reward ID"})
		return
	}

	err = h.service.DeleteReward(c.Request.Context(), tenantID.(uuid.UUID), rewardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListRewards lists loyalty rewards with optional filtering
func (h *Handler) ListRewards(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req ListRewardsRequest

	// Parse query parameters
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}
	if rewardType := c.Query("type"); rewardType != "" {
		req.Type = &rewardType
	}
	if programIDStr := c.Query("program_id"); programIDStr != "" {
		if programID, err := uuid.Parse(programIDStr); err == nil {
			req.ProgramID = &programID
		}
	}
	if search := c.Query("search"); search != "" {
		req.Search = &search
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	response, err := h.service.ListRewards(c.Request.Context(), tenantID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListProgramRewards lists rewards for a specific program
func (h *Handler) ListProgramRewards(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	var req ListProgramRewardsRequest

	// Parse query parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	response, err := h.service.ListProgramRewards(c.Request.Context(), tenantID.(uuid.UUID), programID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Analytics

// GetStats retrieves loyalty analytics and statistics
func (h *Handler) GetStats(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	stats, err := h.service.GetStats(c.Request.Context(), tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// RegisterRoutes registers all loyalty routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	loyalty := router.Group("/loyalty")
	{
		// Loyalty Programs
		programs := loyalty.Group("/programs")
		{
			programs.POST("", h.CreateProgram)
			programs.GET("", h.ListPrograms) // Supports ?type=stats
			programs.GET("/:id", h.GetProgram)
			programs.PUT("/:id", h.UpdateProgram)
			programs.DELETE("/:id", h.DeleteProgram)
			programs.GET("/:id/rewards", h.ListProgramRewards)
		}

		// Loyalty Accounts
		accounts := loyalty.Group("/accounts")
		{
			accounts.POST("", h.CreateAccount)
			accounts.GET("", h.ListAccounts)
			accounts.GET("/:id", h.GetAccount)
			accounts.PUT("/:id", h.UpdateAccount)
			accounts.POST("/:id/earn", h.EarnPoints)
			accounts.POST("/:id/redeem", h.RedeemPoints)
			accounts.POST("/:id/redeem-reward", h.RedeemReward)
			accounts.GET("/:id/transactions", h.ListAccountTransactions)
		}

		// Loyalty Transactions
		transactions := loyalty.Group("/transactions")
		{
			transactions.GET("", h.ListTransactions)
			transactions.GET("/:id", h.GetTransaction)
		}

		// Loyalty Rewards
		rewards := loyalty.Group("/rewards")
		{
			rewards.POST("", h.CreateReward)
			rewards.GET("", h.ListRewards)
			rewards.GET("/:id", h.GetReward)
			rewards.PUT("/:id", h.UpdateReward)
			rewards.DELETE("/:id", h.DeleteReward)
		}
	}
}