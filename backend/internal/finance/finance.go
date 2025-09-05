package finance

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionType represents the type of financial transaction
type TransactionType string

const (
	TransactionTypeDebit  TransactionType = "debit"
	TransactionTypeCredit TransactionType = "credit"
)

// AccountType represents the type of financial account
type AccountType string

const (
	AccountTypeAsset     AccountType = "asset"
	AccountTypeLiability AccountType = "liability"
	AccountTypeEquity    AccountType = "equity"
	AccountTypeRevenue   AccountType = "revenue"
	AccountTypeExpense   AccountType = "expense"
)

// PayoutStatus represents the status of a payout
type PayoutStatus string

const (
	PayoutStatusPending   PayoutStatus = "pending"
	PayoutStatusProcessed PayoutStatus = "processed"
	PayoutStatusCompleted PayoutStatus = "completed"
	PayoutStatusFailed    PayoutStatus = "failed"
	PayoutStatusCancelled PayoutStatus = "cancelled"
)

// ReportPeriod represents the period for financial reports
type ReportPeriod string

const (
	ReportPeriodDaily     ReportPeriod = "daily"
	ReportPeriodWeekly    ReportPeriod = "weekly"
	ReportPeriodMonthly   ReportPeriod = "monthly"
	ReportPeriodQuarterly ReportPeriod = "quarterly"
	ReportPeriodYearly    ReportPeriod = "yearly"
)

// Account represents a financial account in the chart of accounts
type Account struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID   `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Code        string      `json:"code" gorm:"size:20;not null;uniqueIndex:idx_tenant_account_code"`
	Name        string      `json:"name" gorm:"size:255;not null"`
	Description string      `json:"description" gorm:"type:text"`
	Type        AccountType `json:"type" gorm:"size:20;not null;index"`
	ParentID    *uuid.UUID  `json:"parent_id,omitempty" gorm:"type:uuid;index"`
	IsActive    bool        `json:"is_active" gorm:"default:true;index"`
	Balance     float64     `json:"balance" gorm:"type:decimal(15,2);default:0"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Parent      *Account      `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []*Account    `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Transactions []*Transaction `json:"transactions,omitempty" gorm:"many2many:transaction_accounts;"`
}

// Transaction represents a financial transaction
type Transaction struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID        uuid.UUID       `json:"tenant_id" gorm:"type:uuid;not null;index"`
	TransactionNumber string        `json:"transaction_number" gorm:"size:50;not null;uniqueIndex:idx_tenant_transaction_number"`
	Description     string          `json:"description" gorm:"size:500;not null"`
	Reference       string          `json:"reference" gorm:"size:100"`
	Amount          float64         `json:"amount" gorm:"type:decimal(15,2);not null"`
	Type            TransactionType `json:"type" gorm:"size:10;not null;index"`
	TransactionDate time.Time       `json:"transaction_date" gorm:"not null;index"`
	OrderID         *uuid.UUID      `json:"order_id,omitempty" gorm:"type:uuid;index"`
	PaymentID       *uuid.UUID      `json:"payment_id,omitempty" gorm:"type:uuid;index"`
	RefundID        *uuid.UUID      `json:"refund_id,omitempty" gorm:"type:uuid;index"`
	Metadata        map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `json:"-" gorm:"index"`

	// Relations
	Accounts []*Account `json:"accounts,omitempty" gorm:"many2many:transaction_accounts;"`
	Entries  []*TransactionEntry `json:"entries,omitempty" gorm:"foreignKey:TransactionID"`
}

// TransactionEntry represents individual entries in a transaction (double-entry bookkeeping)
type TransactionEntry struct {
	ID            uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TransactionID uuid.UUID       `json:"transaction_id" gorm:"type:uuid;not null;index"`
	AccountID     uuid.UUID       `json:"account_id" gorm:"type:uuid;not null;index"`
	Type          TransactionType `json:"type" gorm:"size:10;not null"`
	Amount        float64         `json:"amount" gorm:"type:decimal(15,2);not null"`
	Description   string          `json:"description" gorm:"size:255"`
	CreatedAt     time.Time       `json:"created_at"`

	// Relations
	Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
	Account     *Account     `json:"account,omitempty" gorm:"foreignKey:AccountID"`
}

// Payout represents a payout to vendors or other parties
type Payout struct {
	ID              uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID        uuid.UUID    `json:"tenant_id" gorm:"type:uuid;not null;index"`
	PayoutNumber    string       `json:"payout_number" gorm:"size:50;not null;uniqueIndex:idx_tenant_payout_number"`
	RecipientID     uuid.UUID    `json:"recipient_id" gorm:"type:uuid;not null;index"`
	RecipientType   string       `json:"recipient_type" gorm:"size:20;not null"` // vendor, affiliate, etc.
	Amount          float64      `json:"amount" gorm:"type:decimal(15,2);not null"`
	Currency        string       `json:"currency" gorm:"size:3;not null;default:'USD'"`
	Status          PayoutStatus `json:"status" gorm:"size:20;not null;default:'pending';index"`
	Description     string       `json:"description" gorm:"size:500"`
	PaymentMethod   string       `json:"payment_method" gorm:"size:50"`
	PaymentDetails  map[string]interface{} `json:"payment_details,omitempty" gorm:"type:jsonb"`
	ScheduledDate   *time.Time   `json:"scheduled_date,omitempty" gorm:"index"`
	ProcessedDate   *time.Time   `json:"processed_date,omitempty" gorm:"index"`
	CompletedDate   *time.Time   `json:"completed_date,omitempty" gorm:"index"`
	FailureReason   string       `json:"failure_reason,omitempty" gorm:"size:500"`
	TransactionID   *uuid.UUID   `json:"transaction_id,omitempty" gorm:"type:uuid;index"`
	Metadata        map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
}

// ReconciliationRecord represents account reconciliation data
type ReconciliationRecord struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID          uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	AccountID         uuid.UUID  `json:"account_id" gorm:"type:uuid;not null;index"`
	ReconciliationDate time.Time `json:"reconciliation_date" gorm:"not null;index"`
	BookBalance       float64    `json:"book_balance" gorm:"type:decimal(15,2);not null"`
	BankBalance       float64    `json:"bank_balance" gorm:"type:decimal(15,2);not null"`
	Difference        float64    `json:"difference" gorm:"type:decimal(15,2);not null"`
	IsReconciled      bool       `json:"is_reconciled" gorm:"default:false;index"`
	Notes             string     `json:"notes" gorm:"type:text"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relations
	Account *Account `json:"account,omitempty" gorm:"foreignKey:AccountID"`
}

// Business Logic Methods

// Account methods
func (a *Account) IsDebitAccount() bool {
	return a.Type == AccountTypeAsset || a.Type == AccountTypeExpense
}

func (a *Account) IsCreditAccount() bool {
	return a.Type == AccountTypeLiability || a.Type == AccountTypeEquity || a.Type == AccountTypeRevenue
}

func (a *Account) UpdateBalance(amount float64, transactionType TransactionType) {
	if a.IsDebitAccount() {
		if transactionType == TransactionTypeDebit {
			a.Balance += amount
		} else {
			a.Balance -= amount
		}
	} else {
		if transactionType == TransactionTypeCredit {
			a.Balance += amount
		} else {
			a.Balance -= amount
		}
	}
}

func (a *Account) GetFullCode() string {
	if a.Parent != nil {
		return a.Parent.GetFullCode() + "." + a.Code
	}
	return a.Code
}

// Transaction methods
func (t *Transaction) GenerateTransactionNumber() string {
	return "TXN-" + time.Now().Format("20060102") + "-" + t.ID.String()[:8]
}

func (t *Transaction) IsBalanced() bool {
	var debitTotal, creditTotal float64
	for _, entry := range t.Entries {
		if entry.Type == TransactionTypeDebit {
			debitTotal += entry.Amount
		} else {
			creditTotal += entry.Amount
		}
	}
	return debitTotal == creditTotal
}

func (t *Transaction) GetTotalAmount() float64 {
	var total float64
	for _, entry := range t.Entries {
		if entry.Type == TransactionTypeDebit {
			total += entry.Amount
		}
	}
	return total
}

// Payout methods
func (p *Payout) GeneratePayoutNumber() string {
	return "PAY-" + time.Now().Format("20060102") + "-" + p.ID.String()[:8]
}

func (p *Payout) CanProcess() bool {
	return p.Status == PayoutStatusPending
}

func (p *Payout) CanCancel() bool {
	return p.Status == PayoutStatusPending || p.Status == PayoutStatusProcessed
}

func (p *Payout) IsCompleted() bool {
	return p.Status == PayoutStatusCompleted
}

func (p *Payout) IsFailed() bool {
	return p.Status == PayoutStatusFailed
}

func (p *Payout) MarkAsProcessed() {
	p.Status = PayoutStatusProcessed
	now := time.Now()
	p.ProcessedDate = &now
}

func (p *Payout) MarkAsCompleted() {
	p.Status = PayoutStatusCompleted
	now := time.Now()
	p.CompletedDate = &now
}

func (p *Payout) MarkAsFailed(reason string) {
	p.Status = PayoutStatusFailed
	p.FailureReason = reason
}

// ReconciliationRecord methods
func (r *ReconciliationRecord) CalculateDifference() {
	r.Difference = r.BankBalance - r.BookBalance
}

func (r *ReconciliationRecord) IsInBalance() bool {
	return r.Difference == 0
}

func (r *ReconciliationRecord) HasDiscrepancy() bool {
	return r.Difference != 0
}

// GORM Hooks
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.TransactionNumber == "" {
		t.TransactionNumber = t.GenerateTransactionNumber()
	}
	return nil
}

func (p *Payout) BeforeCreate(tx *gorm.DB) error {
	if p.PayoutNumber == "" {
		p.PayoutNumber = p.GeneratePayoutNumber()
	}
	return nil
}

func (r *ReconciliationRecord) BeforeSave(tx *gorm.DB) error {
	r.CalculateDifference()
	r.IsReconciled = r.IsInBalance()
	return nil
}

// Filter structs for API endpoints

// AccountFilters represents filters for account listing
type AccountFilters struct {
	Page      int           `json:"page"`
	Limit     int           `json:"limit"`
	SortBy    string        `json:"sort_by"`
	SortOrder string        `json:"sort_order"`
	Search    string        `json:"search"`
	Type      []AccountType `json:"type"`
	ParentID  *uuid.UUID    `json:"parent_id"`
	IsActive  *bool         `json:"is_active"`
}

// TransactionFilters represents filters for transaction listing
type TransactionFilters struct {
	Page        int               `json:"page"`
	Limit       int               `json:"limit"`
	SortBy      string            `json:"sort_by"`
	SortOrder   string            `json:"sort_order"`
	Search      string            `json:"search"`
	Type        []TransactionType `json:"type"`
	AccountID   *uuid.UUID        `json:"account_id"`
	DateAfter   *time.Time        `json:"date_after"`
	DateBefore  *time.Time        `json:"date_before"`
	MinAmount   *float64          `json:"min_amount"`
	MaxAmount   *float64          `json:"max_amount"`
}

// PayoutFilters represents filters for payout listing
type PayoutFilters struct {
	Page          int            `json:"page"`
	Limit         int            `json:"limit"`
	SortBy        string         `json:"sort_by"`
	SortOrder     string         `json:"sort_order"`
	Search        string         `json:"search"`
	Status        []PayoutStatus `json:"status"`
	RecipientID   *uuid.UUID     `json:"recipient_id"`
	RecipientType string         `json:"recipient_type"`
	StartDate     *time.Time     `json:"start_date"`
	EndDate       *time.Time     `json:"end_date"`
	MinAmount     *float64       `json:"min_amount"`
	MaxAmount     *float64       `json:"max_amount"`
}

// ReconciliationFilters represents filters for reconciliation listing
type ReconciliationFilters struct {
	Page           int        `json:"page"`
	Limit          int        `json:"limit"`
	SortBy         string     `json:"sort_by"`
	SortOrder      string     `json:"sort_order"`
	AccountID      *uuid.UUID `json:"account_id"`
	IsReconciled   *bool      `json:"is_reconciled"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	HasDiscrepancy *bool      `json:"has_discrepancy"`
}