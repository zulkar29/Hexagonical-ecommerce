package finance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for finance business logic
type Service interface {
	// Account operations
	CreateAccount(ctx context.Context, account *Account) (*Account, error)
	GetAccount(ctx context.Context, tenantID, accountID uuid.UUID) (*Account, error)
	ListAccounts(ctx context.Context, tenantID uuid.UUID, filter AccountFilters) ([]*Account, int64, error)
	UpdateAccount(ctx context.Context, account *Account) (*Account, error)
	DeleteAccount(ctx context.Context, tenantID, accountID uuid.UUID) error

	// Transaction operations
	CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error)
	GetTransaction(ctx context.Context, tenantID, transactionID uuid.UUID) (*Transaction, error)
	ListTransactions(ctx context.Context, tenantID uuid.UUID, filter TransactionFilters) ([]*Transaction, int64, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error)
	DeleteTransaction(ctx context.Context, tenantID, transactionID uuid.UUID) error

	// Payout operations
	CreatePayout(ctx context.Context, payout *Payout) (*Payout, error)
	GetPayout(ctx context.Context, tenantID, payoutID uuid.UUID) (*Payout, error)
	ListPayouts(ctx context.Context, tenantID uuid.UUID, filter PayoutFilters) ([]*Payout, int64, error)
	ProcessPayout(ctx context.Context, tenantID, payoutID uuid.UUID, processedBy uuid.UUID) (*Payout, error)

	// Reconciliation operations
	CreateReconciliation(ctx context.Context, reconciliation *ReconciliationRecord) (*ReconciliationRecord, error)
	GetReconciliation(ctx context.Context, tenantID, reconciliationID uuid.UUID) (*ReconciliationRecord, error)
	ListReconciliations(ctx context.Context, tenantID uuid.UUID, filter ReconciliationFilters) ([]*ReconciliationRecord, int64, error)

	// Reporting operations
	GetTrialBalance(ctx context.Context, tenantID uuid.UUID, asOfDate time.Time) ([]*TrialBalanceEntry, error)
	GetProfitAndLoss(ctx context.Context, tenantID uuid.UUID, period ReportPeriod, startDate, endDate time.Time) (*ProfitAndLossReport, error)
	GetBalanceSheet(ctx context.Context, tenantID uuid.UUID, asOfDate time.Time) (*BalanceSheetReport, error)
	GetCashFlow(ctx context.Context, tenantID uuid.UUID, period ReportPeriod, startDate, endDate time.Time) (*CashFlowReport, error)
	GetRevenueReport(ctx context.Context, tenantID uuid.UUID, period ReportPeriod, startDate, endDate time.Time) (*RevenueReport, error)
	GetExpenseReport(ctx context.Context, tenantID uuid.UUID, period ReportPeriod, startDate, endDate time.Time) (*ExpenseReport, error)
	GetTaxReport(ctx context.Context, tenantID uuid.UUID, period ReportPeriod, startDate, endDate time.Time) (*TaxReport, error)
}

// service implements the Service interface
type service struct {
	repo Repository
	// Add external service dependencies here
	// orderService OrderService
	// paymentService PaymentService
	// notificationService NotificationService
}

// NewService creates a new finance service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// CreateAccount creates a new account
func (s *service) CreateAccount(ctx context.Context, account *Account) (*Account, error) {
	// Validate account
	if err := s.validateAccount(account); err != nil {
		return nil, err
	}
	
	// Set tenant ID from context
	account.TenantID = getTenantIDFromContext(ctx)
	account.ID = uuid.New()
	
	// Create account
	if err := s.repo.CreateAccount(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	
	return account, nil
}

// GetAccount retrieves an account by ID
func (s *service) GetAccount(ctx context.Context, tenantID, accountID uuid.UUID) (*Account, error) {
	return s.repo.GetAccount(ctx, tenantID, accountID)
}

// ListAccounts retrieves accounts with filtering and pagination
func (s *service) ListAccounts(ctx context.Context, tenantID uuid.UUID, filter AccountFilters) ([]*Account, int64, error) {
	return s.repo.ListAccounts(ctx, tenantID, filter)
}

// UpdateAccount updates an existing account
func (s *service) UpdateAccount(ctx context.Context, account *Account) (*Account, error) {
	// Validate account
	if err := s.validateAccount(account); err != nil {
		return nil, err
	}
	
	// Update account
	if err := s.repo.UpdateAccount(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}
	
	return account, nil
}

// DeleteAccount deletes an account
func (s *service) DeleteAccount(ctx context.Context, tenantID, accountID uuid.UUID) error {
	return s.repo.DeleteAccount(ctx, tenantID, accountID)
}

// CreateTransaction creates a new transaction
func (s *service) CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	// Validate transaction
	if err := s.validateTransaction(transaction); err != nil {
		return nil, err
	}
	
	// Set tenant ID from context
	transaction.TenantID = getTenantIDFromContext(ctx)
	transaction.ID = uuid.New()
	
	// Generate transaction number
	transaction.GenerateTransactionNumber()
	
	// Create transaction
	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	
	return transaction, nil
}

// GetTransaction retrieves a transaction by ID
func (s *service) GetTransaction(ctx context.Context, tenantID, transactionID uuid.UUID) (*Transaction, error) {
	return s.repo.GetTransaction(ctx, tenantID, transactionID)
}

// ListTransactions retrieves transactions with filtering and pagination
func (s *service) ListTransactions(ctx context.Context, tenantID uuid.UUID, filter TransactionFilters) ([]*Transaction, int64, error) {
	return s.repo.ListTransactions(ctx, tenantID, filter)
}

// UpdateTransaction updates an existing transaction
func (s *service) UpdateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	// Validate transaction
	if err := s.validateTransaction(transaction); err != nil {
		return nil, err
	}
	
	// Update transaction
	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}
	
	return transaction, nil
}

// DeleteTransaction deletes a transaction
func (s *service) DeleteTransaction(ctx context.Context, tenantID, transactionID uuid.UUID) error {
	return s.repo.DeleteTransaction(ctx, tenantID, transactionID)
}

// CreatePayout creates a new payout
func (s *service) CreatePayout(ctx context.Context, payout *Payout) (*Payout, error) {
	// Validate payout
	if err := s.validatePayout(payout); err != nil {
		return nil, err
	}
	
	// Set tenant ID from context
	payout.TenantID = getTenantIDFromContext(ctx)
	payout.ID = uuid.New()
	payout.Status = PayoutStatusPending
	
	// Create payout
	if err := s.repo.CreatePayout(ctx, payout); err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}
	
	return payout, nil
}

// GetPayout retrieves a payout by ID
func (s *service) GetPayout(ctx context.Context, tenantID, payoutID uuid.UUID) (*Payout, error) {
	return s.repo.GetPayout(ctx, tenantID, payoutID)
}

// ListPayouts retrieves payouts with filtering and pagination
func (s *service) ListPayouts(ctx context.Context, tenantID uuid.UUID, filter PayoutFilters) ([]*Payout, int64, error) {
	return s.repo.ListPayouts(ctx, tenantID, filter)
}

// ProcessPayout processes a payout
func (s *service) ProcessPayout(ctx context.Context, tenantID, payoutID uuid.UUID, processedBy uuid.UUID) (*Payout, error) {
	// Get payout
	payout, err := s.repo.GetPayout(ctx, tenantID, payoutID)
	if err != nil {
		return nil, err
	}
	
	// Check if payout can be processed
	if !payout.CanProcess() {
		return nil, errors.New("payout cannot be processed")
	}
	
	// Update payout status
	payout.Status = PayoutStatusProcessed
	payout.ProcessedAt = &time.Time{}
	*payout.ProcessedAt = time.Now()
	payout.ProcessedBy = &processedBy
	
	// Update payout
	if err := s.repo.UpdatePayout(ctx, payout); err != nil {
		return nil, fmt.Errorf("failed to process payout: %w", err)
	}
	
	return payout, nil
}

// CreateReconciliation creates a new reconciliation record
func (s *service) CreateReconciliation(ctx context.Context, reconciliation *ReconciliationRecord) (*ReconciliationRecord, error) {
	// Set tenant ID from context
	reconciliation.TenantID = getTenantIDFromContext(ctx)
	reconciliation.ID = uuid.New()
	
	// Calculate difference
	reconciliation.CalculateDifference()
	
	// Create reconciliation
	if err := s.repo.CreateReconciliation(ctx, reconciliation); err != nil {
		return nil, fmt.Errorf("failed to create reconciliation: %w", err)
	}
	
	return reconciliation, nil
}

// GetReconciliation retrieves a reconciliation record by ID
func (s *service) GetReconciliation(ctx context.Context, tenantID, reconciliationID uuid.UUID) (*ReconciliationRecord, error) {
	return s.repo.GetReconciliation(ctx, tenantID, reconciliationID)
}

// ListReconciliations retrieves reconciliation records with filtering and pagination
func (s *service) ListReconciliations(ctx context.Context, tenantID uuid.UUID, filter ReconciliationFilters) ([]*ReconciliationRecord, int64, error) {
	return s.repo.ListReconciliations(ctx, tenantID, filter)
}

// GetTrialBalance generates a trial balance report
func (s *service) GetTrialBalance(ctx context.Context, tenantID uuid.UUID, asOfDate time.Time) ([]*TrialBalanceEntry, error) {
	return s.repo.GetTrialBalance(ctx, tenantID, asOfDate)
}

// GetProfitAndLoss generates a profit and loss report
func (s *service) GetProfitAndLoss(ctx context.Context, tenantID uuid.UUID, period ReportPeriod, startDate, endDate time.Time) (*ProfitAndLossReport, error) {
	return s.repo.GetProfitAndLoss(ctx, tenantID, startDate, endDate)
}

// GetBalanceSheet generates a balance sheet report
func (s *service) GetBalanceSheet(ctx context.Context, tenantID uuid.UUID, asOfDate time.Time) (*BalanceSheetReport, error) {
	return s.repo.GetBalanceSheet(ctx, tenantID, asOfDate)
}

// GetCashFlow generates a cash flow report
func (s *service) GetCashFlow(ctx context.Context, tenantID uuid.UUID, period ReportPeriod, startDate, endDate time.Time) (*CashFlowReport, error) {
	return s.repo.GetCashFlow(ctx, tenantID, period, startDate, endDate)
}

// GetRevenueReport generates a revenue report
func (s *service) GetRevenueReport(ctx context.Context, tenantID uuid.UUID, period ReportPeriod, startDate, endDate time.Time) (*RevenueReport, error) {
	return s.repo.GetRevenueReport(ctx, tenantID, period, startDate, endDate)
}

// GetExpenseReport generates an expense report
func (s *service) GetExpenseReport(ctx context.Context, tenantID uuid.UUID, period ReportPeriod, startDate, endDate time.Time) (*ExpenseReport, error) {
	return s.repo.GetExpenseReport(ctx, tenantID, period, startDate, endDate)
}

// GetTaxReport generates a tax report
func (s *service) GetTaxReport(ctx context.Context, tenantID uuid.UUID, period ReportPeriod, startDate, endDate time.Time) (*TaxReport, error) {
	return s.repo.GetTaxReport(ctx, tenantID, period, startDate, endDate)
}

// Helper functions
func (s *service) validateAccount(account *Account) error {
	if account.Code == "" {
		return errors.New("account code is required")
	}
	if account.Name == "" {
		return errors.New("account name is required")
	}
	return nil
}

func (s *service) validateTransaction(transaction *Transaction) error {
	if transaction.Description == "" {
		return errors.New("transaction description is required")
	}
	if transaction.Amount <= 0 {
		return errors.New("transaction amount must be greater than zero")
	}
	if len(transaction.Entries) < 2 {
		return errors.New("transaction must have at least 2 entries")
	}
	return nil
}

func (s *service) validatePayout(payout *Payout) error {
	if payout.Amount <= 0 {
		return errors.New("payout amount must be greater than zero")
	}
	if payout.RecipientID == uuid.Nil {
		return errors.New("payout recipient is required")
	}
	return nil
}