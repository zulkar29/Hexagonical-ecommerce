package finance

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for finance operations
type Handler struct {
	service Service
}

// NewHandler creates a new finance handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers all finance routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Account management endpoints
	accounts := router.Group("/accounts")
	{
		accounts.POST("", h.CreateAccount)                    // POST /api/v1/accounts
		accounts.GET("", h.ListAccounts)                     // GET /api/v1/accounts
		accounts.GET("/:id", h.GetAccount)                   // GET /api/v1/accounts/:id
		accounts.PATCH("/:id", h.UpdateAccount)              // PATCH /api/v1/accounts/:id
		accounts.DELETE("/:id", h.DeleteAccount)             // DELETE /api/v1/accounts/:id
	}

	// Transaction management endpoints
	transactions := router.Group("/transactions")
	{
		transactions.POST("", h.CreateTransaction)                    // POST /api/v1/transactions
		transactions.GET("", h.ListTransactions)                     // GET /api/v1/transactions
		transactions.GET("/:id", h.GetTransaction)                   // GET /api/v1/transactions/:id
		transactions.PATCH("/:id", h.UpdateTransaction)              // PATCH /api/v1/transactions/:id
		transactions.DELETE("/:id", h.DeleteTransaction)             // DELETE /api/v1/transactions/:id
	}

	// Payout management endpoints
	payouts := router.Group("/payouts")
	{
		payouts.POST("", h.CreatePayout)                    // POST /api/v1/payouts
		payouts.GET("", h.ListPayouts)                     // GET /api/v1/payouts
		payouts.GET("/:id", h.GetPayout)                   // GET /api/v1/payouts/:id
		payouts.POST("/:id/process", h.ProcessPayout)      // POST /api/v1/payouts/:id/process
	}

	// Reconciliation management endpoints
	reconciliations := router.Group("/reconciliations")
	{
		reconciliations.POST("", h.CreateReconciliation)                    // POST /api/v1/reconciliations
		reconciliations.GET("", h.ListReconciliations)                     // GET /api/v1/reconciliations
		reconciliations.GET("/:id", h.GetReconciliation)                   // GET /api/v1/reconciliations/:id
	}

	// Financial reports endpoints
	reports := router.Group("/reports")
	{
		reports.GET("/trial-balance", h.GetTrialBalance)           // GET /api/v1/reports/trial-balance
		reports.GET("/profit-loss", h.GetProfitAndLoss)           // GET /api/v1/reports/profit-loss
		reports.GET("/balance-sheet", h.GetBalanceSheet)          // GET /api/v1/reports/balance-sheet
		reports.GET("/cash-flow", h.GetCashFlow)                  // GET /api/v1/reports/cash-flow
		reports.GET("/revenue", h.GetRevenueReport)               // GET /api/v1/reports/revenue
		reports.GET("/expense", h.GetExpenseReport)               // GET /api/v1/reports/expense
		reports.GET("/tax", h.GetTaxReport)                       // GET /api/v1/reports/tax
	}
}

// Account Handlers

// CreateAccount creates a new account
// @Summary Create a new account
// @Description Create a new financial account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body Account true "Account creation request"
// @Success 201 {object} Account
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	var account Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	createdAccount, err := h.service.CreateAccount(c.Request.Context(), &account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdAccount)
}

// ListAccounts retrieves accounts with filtering and pagination
// @Summary List accounts
// @Description Get a list of accounts with optional filtering and pagination
// @Tags accounts
// @Accept json
// @Produce json
// @Param type query []string false "Filter by account type" collectionFormat(multi)
// @Param parent_id query string false "Filter by parent account ID"
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search in account code, name, or description"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort field" default(code)
// @Param sort_order query string false "Sort order (asc/desc)" default(asc)
// @Success 200 {object} PaginatedAccountsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts [get]
func (h *Handler) ListAccounts(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse query parameters
	filter := AccountFilters{
		Page:      getIntQueryParam(c, "page", 1),
		Limit:     getIntQueryParam(c, "limit", 20),
		SortBy:    c.DefaultQuery("sort_by", "code"),
		SortOrder: c.DefaultQuery("sort_order", "asc"),
		Search:    c.Query("search"),
	}

	// Parse type filter
	if typeParams := c.QueryArray("type"); len(typeParams) > 0 {
		for _, accountType := range typeParams {
			filter.Type = append(filter.Type, AccountType(accountType))
		}
	}

	// Parse parent_id filter
	if parentID := c.Query("parent_id"); parentID != "" {
		if id, err := uuid.Parse(parentID); err == nil {
			filter.ParentID = &id
		}
	}

	// Parse is_active filter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &isActive
		}
	}

	accounts, total, err := h.service.ListAccounts(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list accounts", "details": err.Error()})
		return
	}

	response := PaginatedAccountsResponse{
		Data:       accounts,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: (total + int64(filter.Limit) - 1) / int64(filter.Limit),
	}

	c.JSON(http.StatusOK, response)
}

// GetAccount retrieves a specific account
// @Summary Get account by ID
// @Description Get a specific account by its ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} Account
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/{id} [get]
func (h *Handler) GetAccount(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := h.service.GetAccount(c.Request.Context(), tenantID, accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// UpdateAccount updates an existing account
// @Summary Update account
// @Description Update an existing account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param account body Account true "Account update request"
// @Success 200 {object} Account
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/{id} [patch]
func (h *Handler) UpdateAccount(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var account Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	account.ID = accountID
	updatedAccount, err := h.service.UpdateAccount(c.Request.Context(), &account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAccount)
}

// DeleteAccount deletes an account
// @Summary Delete account
// @Description Delete an account by ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/{id} [delete]
func (h *Handler) DeleteAccount(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	err = h.service.DeleteAccount(c.Request.Context(), tenantID, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account", "details": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Transaction Handlers

// CreateTransaction creates a new transaction
func (h *Handler) CreateTransaction(c *gin.Context) {
	var transaction Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	createdTransaction, err := h.service.CreateTransaction(c.Request.Context(), &transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTransaction)
}

// ListTransactions retrieves transactions with filtering and pagination
func (h *Handler) ListTransactions(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse query parameters
	filter := TransactionFilters{
		Page:      getIntQueryParam(c, "page", 1),
		Limit:     getIntQueryParam(c, "limit", 20),
		SortBy:    c.DefaultQuery("sort_by", "transaction_date"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
		Search:    c.Query("search"),
	}

	// Parse type filter
	if typeParams := c.QueryArray("type"); len(typeParams) > 0 {
		for _, transactionType := range typeParams {
			filter.Type = append(filter.Type, TransactionType(transactionType))
		}
	}

	// Parse account_id filter
	if accountID := c.Query("account_id"); accountID != "" {
		if id, err := uuid.Parse(accountID); err == nil {
			filter.AccountID = &id
		}
	}

	// Parse date filters
	if dateAfter := c.Query("date_after"); dateAfter != "" {
		if date, err := time.Parse(time.RFC3339, dateAfter); err == nil {
			filter.DateAfter = &date
		}
	}

	if dateBefore := c.Query("date_before"); dateBefore != "" {
		if date, err := time.Parse(time.RFC3339, dateBefore); err == nil {
			filter.DateBefore = &date
		}
	}

	// Parse amount filters
	if minAmount := c.Query("min_amount"); minAmount != "" {
		if amount, err := strconv.ParseFloat(minAmount, 64); err == nil {
			filter.MinAmount = &amount
		}
	}

	if maxAmount := c.Query("max_amount"); maxAmount != "" {
		if amount, err := strconv.ParseFloat(maxAmount, 64); err == nil {
			filter.MaxAmount = &amount
		}
	}

	transactions, total, err := h.service.ListTransactions(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list transactions", "details": err.Error()})
		return
	}

	response := PaginatedTransactionsResponse{
		Data:       transactions,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: (total + int64(filter.Limit) - 1) / int64(filter.Limit),
	}

	c.JSON(http.StatusOK, response)
}

// GetTransaction retrieves a specific transaction
func (h *Handler) GetTransaction(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.service.GetTransaction(c.Request.Context(), tenantID, transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction updates an existing transaction
func (h *Handler) UpdateTransaction(c *gin.Context) {
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var transaction Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	transaction.ID = transactionID
	updatedTransaction, err := h.service.UpdateTransaction(c.Request.Context(), &transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTransaction)
}

// DeleteTransaction deletes a transaction
func (h *Handler) DeleteTransaction(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	err = h.service.DeleteTransaction(c.Request.Context(), tenantID, transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction", "details": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Payout Handlers

// CreatePayout creates a new payout
func (h *Handler) CreatePayout(c *gin.Context) {
	var payout Payout
	if err := c.ShouldBindJSON(&payout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	createdPayout, err := h.service.CreatePayout(c.Request.Context(), &payout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payout", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPayout)
}

// ListPayouts retrieves payouts with filtering and pagination
func (h *Handler) ListPayouts(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse query parameters
	filter := PayoutFilters{
		Page:      getIntQueryParam(c, "page", 1),
		Limit:     getIntQueryParam(c, "limit", 20),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
		Search:    c.Query("search"),
	}

	// Parse status filter
	if statusParams := c.QueryArray("status"); len(statusParams) > 0 {
		for _, status := range statusParams {
			filter.Status = append(filter.Status, PayoutStatus(status))
		}
	}

	// Parse recipient_id filter
	if recipientID := c.Query("recipient_id"); recipientID != "" {
		if id, err := uuid.Parse(recipientID); err == nil {
			filter.RecipientID = &id
		}
	}

	payouts, total, err := h.service.ListPayouts(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list payouts", "details": err.Error()})
		return
	}

	response := PaginatedPayoutsResponse{
		Data:       payouts,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: (total + int64(filter.Limit) - 1) / int64(filter.Limit),
	}

	c.JSON(http.StatusOK, response)
}

// GetPayout retrieves a specific payout
func (h *Handler) GetPayout(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)
	payoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payout ID"})
		return
	}

	payout, err := h.service.GetPayout(c.Request.Context(), tenantID, payoutID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payout not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payout)
}

// ProcessPayout processes a payout
func (h *Handler) ProcessPayout(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)
	payoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payout ID"})
		return
	}

	// Get processed_by from request body or context
	var req struct {
		ProcessedBy uuid.UUID `json:"processed_by" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	processedPayout, err := h.service.ProcessPayout(c.Request.Context(), tenantID, payoutID, req.ProcessedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payout", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, processedPayout)
}

// Reconciliation Handlers

// CreateReconciliation creates a new reconciliation record
func (h *Handler) CreateReconciliation(c *gin.Context) {
	var reconciliation ReconciliationRecord
	if err := c.ShouldBindJSON(&reconciliation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	createdReconciliation, err := h.service.CreateReconciliation(c.Request.Context(), &reconciliation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reconciliation", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdReconciliation)
}

// ListReconciliations retrieves reconciliation records with filtering and pagination
func (h *Handler) ListReconciliations(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse query parameters
	filter := ReconciliationFilters{
		Page:      getIntQueryParam(c, "page", 1),
		Limit:     getIntQueryParam(c, "limit", 20),
		SortBy:    c.DefaultQuery("sort_by", "reconciliation_date"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}

	// Parse account_id filter
	if accountID := c.Query("account_id"); accountID != "" {
		if id, err := uuid.Parse(accountID); err == nil {
			filter.AccountID = &id
		}
	}

	// Parse is_reconciled filter
	if isReconciledStr := c.Query("is_reconciled"); isReconciledStr != "" {
		if isReconciled, err := strconv.ParseBool(isReconciledStr); err == nil {
			filter.IsReconciled = &isReconciled
		}
	}

	reconciliations, total, err := h.service.ListReconciliations(c.Request.Context(), tenantID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list reconciliations", "details": err.Error()})
		return
	}

	response := PaginatedReconciliationsResponse{
		Data:       reconciliations,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: (total + int64(filter.Limit) - 1) / int64(filter.Limit),
	}

	c.JSON(http.StatusOK, response)
}

// GetReconciliation retrieves a specific reconciliation record
func (h *Handler) GetReconciliation(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)
	reconciliationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reconciliation ID"})
		return
	}

	reconciliation, err := h.service.GetReconciliation(c.Request.Context(), tenantID, reconciliationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reconciliation not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reconciliation)
}

// Report Handlers

// GetTrialBalance generates a trial balance report
func (h *Handler) GetTrialBalance(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse as_of_date parameter
	asOfDateStr := c.Query("as_of_date")
	if asOfDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "as_of_date parameter is required"})
		return
	}

	asOfDate, err := time.Parse(time.RFC3339, asOfDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid as_of_date format, use RFC3339"})
		return
	}

	trialBalance, err := h.service.GetTrialBalance(c.Request.Context(), tenantID, asOfDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate trial balance", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trialBalance)
}

// GetProfitAndLoss generates a profit and loss report
func (h *Handler) GetProfitAndLoss(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse parameters
	period := ReportPeriod(c.DefaultQuery("period", "monthly"))
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use RFC3339"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use RFC3339"})
		return
	}

	report, err := h.service.GetProfitAndLoss(c.Request.Context(), tenantID, period, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate profit and loss report", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetBalanceSheet generates a balance sheet report
func (h *Handler) GetBalanceSheet(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse as_of_date parameter
	asOfDateStr := c.Query("as_of_date")
	if asOfDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "as_of_date parameter is required"})
		return
	}

	asOfDate, err := time.Parse(time.RFC3339, asOfDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid as_of_date format, use RFC3339"})
		return
	}

	balanceSheet, err := h.service.GetBalanceSheet(c.Request.Context(), tenantID, asOfDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate balance sheet", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balanceSheet)
}

// GetCashFlow generates a cash flow report
func (h *Handler) GetCashFlow(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse parameters
	period := ReportPeriod(c.DefaultQuery("period", "monthly"))
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use RFC3339"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use RFC3339"})
		return
	}

	report, err := h.service.GetCashFlow(c.Request.Context(), tenantID, period, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate cash flow report", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetRevenueReport generates a revenue report
func (h *Handler) GetRevenueReport(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse parameters
	period := ReportPeriod(c.DefaultQuery("period", "monthly"))
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use RFC3339"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use RFC3339"})
		return
	}

	report, err := h.service.GetRevenueReport(c.Request.Context(), tenantID, period, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate revenue report", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetExpenseReport generates an expense report
func (h *Handler) GetExpenseReport(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse parameters
	period := ReportPeriod(c.DefaultQuery("period", "monthly"))
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use RFC3339"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use RFC3339"})
		return
	}

	report, err := h.service.GetExpenseReport(c.Request.Context(), tenantID, period, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate expense report", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetTaxReport generates a tax report
func (h *Handler) GetTaxReport(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)

	// Parse parameters
	period := ReportPeriod(c.DefaultQuery("period", "monthly"))
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use RFC3339"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use RFC3339"})
		return
	}

	report, err := h.service.GetTaxReport(c.Request.Context(), tenantID, period, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tax report", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// Helper functions

func getTenantIDFromContext(c *gin.Context) uuid.UUID {
	// TODO: Extract tenant ID from context/JWT token
	// For now, return a placeholder
	if tenantID, exists := c.Get("tenant_id"); exists {
		if id, ok := tenantID.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.New() // Placeholder
}

func getIntQueryParam(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Response types for pagination

type PaginatedAccountsResponse struct {
	Data       []*Account `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int64      `json:"total_pages"`
}

type PaginatedTransactionsResponse struct {
	Data       []*Transaction `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int64          `json:"total_pages"`
}

type PaginatedPayoutsResponse struct {
	Data       []*Payout `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int64     `json:"total_pages"`
}

type PaginatedReconciliationsResponse struct {
	Data       []*ReconciliationRecord `json:"data"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int64                   `json:"total_pages"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}