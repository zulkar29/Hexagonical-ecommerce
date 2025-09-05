package loyalty

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LoyaltyProgram represents a loyalty program
type LoyaltyProgram struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type" gorm:"not null"` // points, tiers, cashback
	Status      string    `json:"status" gorm:"default:'active'"` // active, inactive, draft
	Settings    string    `json:"settings" gorm:"type:jsonb"` // JSON configuration
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// LoyaltyAccount represents a customer's loyalty account
type LoyaltyAccount struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID  uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	ProgramID uuid.UUID `json:"program_id" gorm:"type:uuid;not null"`
	Points    int       `json:"points" gorm:"default:0"`
	Tier      string    `json:"tier" gorm:"default:'bronze'"`
	Status    string    `json:"status" gorm:"default:'active'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Program *LoyaltyProgram `json:"program,omitempty" gorm:"foreignKey:ProgramID"`
}

// LoyaltyTransaction represents a loyalty points transaction
type LoyaltyTransaction struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	AccountID   uuid.UUID `json:"account_id" gorm:"type:uuid;not null"`
	Type        string    `json:"type" gorm:"not null"` // earned, redeemed, expired, adjusted
	Points      int       `json:"points" gorm:"not null"`
	Description string    `json:"description"`
	OrderID     *uuid.UUID `json:"order_id,omitempty" gorm:"type:uuid"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Account *LoyaltyAccount `json:"account,omitempty" gorm:"foreignKey:AccountID"`
}

// LoyaltyReward represents available rewards
type LoyaltyReward struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ProgramID   uuid.UUID `json:"program_id" gorm:"type:uuid;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type" gorm:"not null"` // discount, product, cashback
	PointsCost  int       `json:"points_cost" gorm:"not null"`
	Value       string    `json:"value" gorm:"type:jsonb"` // JSON value configuration
	Status      string    `json:"status" gorm:"default:'active'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Program *LoyaltyProgram `json:"program,omitempty" gorm:"foreignKey:ProgramID"`
}

// Request/Response DTOs

// Program requests
type CreateProgramRequest struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Type        string    `json:"type" binding:"required,oneof=points tiers cashback"`
	Settings    string    `json:"settings"`
}

type UpdateProgramRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	Settings    *string `json:"settings"`
	Status      *string `json:"status"`
}

type ListProgramsRequest struct {
	Status *string `json:"status"`
	Type   *string `json:"type"`
	Search *string `json:"search"`
	Limit  *int    `json:"limit"`
	Offset *int    `json:"offset"`
}

type ListProgramsResponse struct {
	Programs []*LoyaltyProgram `json:"programs"`
	Total    int64             `json:"total"`
}

// Account requests
type CreateAccountRequest struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	ProgramID uuid.UUID `json:"program_id" binding:"required"`
}

type UpdateAccountRequest struct {
	TenantID uuid.UUID `json:"tenant_id"`
	Status   *string   `json:"status,omitempty"`
	Tier     *string   `json:"tier,omitempty"`
}

type ListAccountsRequest struct {
	TenantID  uuid.UUID  `json:"tenant_id"`
	ProgramID *uuid.UUID `json:"program_id,omitempty"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Status    *string    `json:"status,omitempty"`
	Tier      *string    `json:"tier,omitempty"`
	Search    *string    `json:"search,omitempty"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

type ListAccountsResponse struct {
	Accounts []LoyaltyAccount `json:"accounts"`
	Total    int64            `json:"total"`
}

// Transaction requests
type CreateTransactionRequest struct {
	TenantID    uuid.UUID  `json:"tenant_id"`
	AccountID   uuid.UUID  `json:"account_id" binding:"required"`
	Type        string     `json:"type" binding:"required"`
	Points      int        `json:"points" binding:"required"`
	Description string     `json:"description"`
	OrderID     *uuid.UUID `json:"order_id"`
}

type ListTransactionsRequest struct {
	AccountID *uuid.UUID `json:"account_id"`
	Type      *string    `json:"type"`
	Limit     *int       `json:"limit"`
	Offset    *int       `json:"offset"`
}

type ListTransactionsResponse struct {
	Transactions []LoyaltyTransaction `json:"transactions"`
	Total        int64                `json:"total"`
}

type ListAccountTransactionsRequest struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	AccountID uuid.UUID `json:"account_id"`
	Limit     *int      `json:"limit,omitempty"`
	Offset    *int      `json:"offset,omitempty"`
}

type ListAccountTransactionsResponse struct {
	Transactions []LoyaltyTransaction `json:"transactions"`
	Total        int64                `json:"total"`
}

// Reward requests
type CreateRewardRequest struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	ProgramID   uuid.UUID `json:"program_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Type        string    `json:"type" binding:"required,oneof=discount product cashback"`
	PointsCost  int       `json:"points_cost" binding:"required,min=1"`
	Value       string    `json:"value" binding:"required"`
}

type UpdateRewardRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	PointsCost  *int    `json:"points_cost"`
	Value       *string `json:"value"`
	Status      *string `json:"status"`
}

type ListRewardsRequest struct {
	TenantID  uuid.UUID  `json:"tenant_id"`
	ProgramID *uuid.UUID `json:"program_id,omitempty"`
	Type      *string    `json:"type,omitempty"`
	Status    *string    `json:"status,omitempty"`
	Search    *string    `json:"search,omitempty"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

type ListRewardsResponse struct {
	Rewards []*LoyaltyReward `json:"rewards"`
	Total   int64            `json:"total"`
}

// Points operation requests
type EarnPointsRequest struct {
	TenantID  uuid.UUID  `json:"tenant_id"`
	AccountID uuid.UUID  `json:"account_id"`
	Points    int64      `json:"points"`
	Reason    string     `json:"reason"`
	OrderID   *uuid.UUID `json:"order_id,omitempty"`
}

type RedeemPointsRequest struct {
	TenantID  uuid.UUID  `json:"tenant_id"`
	AccountID uuid.UUID  `json:"account_id"`
	Points    int64      `json:"points"`
	Reason    string     `json:"reason"`
	OrderID   *uuid.UUID `json:"order_id,omitempty"`
}

type AdjustPointsRequest struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	AccountID uuid.UUID `json:"account_id"`
	Points    int64     `json:"points"`
	Reason    string    `json:"reason"`
}

type ListProgramRewardsRequest struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	ProgramID uuid.UUID `json:"program_id"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

type ListProgramRewardsResponse struct {
	Rewards []*LoyaltyReward `json:"rewards"`
	Total   int64            `json:"total"`
}

type LoyaltyAnalytics struct {
	TotalPrograms     int64 `json:"total_programs"`
	TotalAccounts     int64 `json:"total_accounts"`
	TotalTransactions int64 `json:"total_transactions"`
	TotalRewards      int64 `json:"total_rewards"`
	TotalPointsEarned int64 `json:"total_points_earned"`
	TotalPointsRedeemed int64 `json:"total_points_redeemed"`
}

type RedeemRewardRequest struct {
	RewardID uuid.UUID `json:"reward_id" binding:"required"`
}

// Legacy types for backward compatibility
type CreateLoyaltyProgramRequest = CreateProgramRequest
type UpdateLoyaltyProgramRequest = UpdateProgramRequest
type CreateLoyaltyRewardRequest = CreateRewardRequest
type UpdateLoyaltyRewardRequest = UpdateRewardRequest

type LoyaltyStats struct {
	TotalPrograms    int64 `json:"total_programs"`
	ActivePrograms   int64 `json:"active_programs"`
	TotalAccounts    int64 `json:"total_accounts"`
	ActiveAccounts   int64 `json:"total_active_accounts"`
	TotalPoints      int64 `json:"total_points"`
	TotalRedemptions int64 `json:"total_redemptions"`
}

// Table names
func (LoyaltyProgram) TableName() string {
	return "loyalty_programs"
}

func (LoyaltyAccount) TableName() string {
	return "loyalty_accounts"
}

func (LoyaltyTransaction) TableName() string {
	return "loyalty_transactions"
}

func (LoyaltyReward) TableName() string {
	return "loyalty_rewards"
}