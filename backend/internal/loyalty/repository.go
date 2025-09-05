package loyalty

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the loyalty repository interface
type Repository interface {
	// Loyalty Programs
	CreateProgram(program *LoyaltyProgram) (*LoyaltyProgram, error)
	GetProgram(tenantID, programID uuid.UUID) (*LoyaltyProgram, error)
	UpdateProgram(program *LoyaltyProgram) (*LoyaltyProgram, error)
	DeleteProgram(tenantID, programID uuid.UUID) error
	ListPrograms(tenantID uuid.UUID) ([]*LoyaltyProgram, int64, error)

	// Loyalty Accounts
	CreateAccount(account *LoyaltyAccount) (*LoyaltyAccount, error)
	GetAccount(tenantID, accountID uuid.UUID) (*LoyaltyAccount, error)
	GetAccountByUser(tenantID, userID uuid.UUID) (*LoyaltyAccount, error)
	UpdateAccount(account *LoyaltyAccount) (*LoyaltyAccount, error)
	ListAccounts(tenantID uuid.UUID, userID *uuid.UUID, status *string, limit, offset int) ([]*LoyaltyAccount, int64, error)

	// Loyalty Transactions
	CreateTransaction(transaction *LoyaltyTransaction) (*LoyaltyTransaction, error)
	UpdateTransaction(transaction *LoyaltyTransaction) (*LoyaltyTransaction, error)
	GetTransaction(tenantID, transactionID uuid.UUID) (*LoyaltyTransaction, error)
	ListTransactions(tenantID uuid.UUID, accountID *uuid.UUID, limit, offset int) ([]*LoyaltyTransaction, int64, error)

	// Loyalty Rewards
	CreateReward(reward *LoyaltyReward) (*LoyaltyReward, error)
	GetReward(tenantID, rewardID uuid.UUID) (*LoyaltyReward, error)
	UpdateReward(reward *LoyaltyReward) (*LoyaltyReward, error)
	DeleteReward(tenantID, rewardID uuid.UUID) error
	ListRewards(tenantID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]*LoyaltyReward, int64, error)

	// Analytics
	GetStats(tenantID uuid.UUID) (*LoyaltyStats, error)
}

// GormRepository implements Repository using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM repository
func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// Loyalty Programs
func (r *GormRepository) CreateProgram(program *LoyaltyProgram) (*LoyaltyProgram, error) {
	err := r.db.Create(program).Error
	return program, err
}

func (r *GormRepository) GetProgram(tenantID, programID uuid.UUID) (*LoyaltyProgram, error) {
	var program LoyaltyProgram
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, programID).First(&program).Error
	if err != nil {
		return nil, err
	}
	return &program, nil
}

func (r *GormRepository) UpdateProgram(program *LoyaltyProgram) (*LoyaltyProgram, error) {
	err := r.db.Save(program).Error
	return program, err
}

func (r *GormRepository) DeleteProgram(tenantID, programID uuid.UUID) error {
	return r.db.Where("tenant_id = ? AND id = ?", tenantID, programID).Delete(&LoyaltyProgram{}).Error
}

func (r *GormRepository) ListPrograms(tenantID uuid.UUID) ([]*LoyaltyProgram, int64, error) {
	var programs []*LoyaltyProgram
	var total int64
	
	query := r.db.Model(&LoyaltyProgram{}).Where("tenant_id = ?", tenantID)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = query.Order("created_at DESC").Find(&programs).Error
	return programs, total, err
}

// Loyalty Accounts
func (r *GormRepository) CreateAccount(account *LoyaltyAccount) (*LoyaltyAccount, error) {
	err := r.db.Create(account).Error
	return account, err
}

func (r *GormRepository) GetAccount(tenantID, accountID uuid.UUID) (*LoyaltyAccount, error) {
	var account LoyaltyAccount
	err := r.db.Preload("Program").Where("tenant_id = ? AND id = ?", tenantID, accountID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *GormRepository) GetAccountByUser(tenantID, userID uuid.UUID) (*LoyaltyAccount, error) {
	var account LoyaltyAccount
	err := r.db.Preload("Program").Where("tenant_id = ? AND user_id = ?", tenantID, userID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *GormRepository) UpdateAccount(account *LoyaltyAccount) (*LoyaltyAccount, error) {
	err := r.db.Save(account).Error
	return account, err
}

func (r *GormRepository) ListAccounts(tenantID uuid.UUID, userID *uuid.UUID, status *string, limit, offset int) ([]*LoyaltyAccount, int64, error) {
	var accounts []*LoyaltyAccount
	var total int64
	
	query := r.db.Model(&LoyaltyAccount{}).Preload("Program").Where("tenant_id = ?", tenantID)
	
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&accounts).Error
	return accounts, total, err
}

// Loyalty Transactions
func (r *GormRepository) CreateTransaction(transaction *LoyaltyTransaction) (*LoyaltyTransaction, error) {
	err := r.db.Create(transaction).Error
	return transaction, err
}

func (r *GormRepository) GetTransaction(tenantID, transactionID uuid.UUID) (*LoyaltyTransaction, error) {
	var transaction LoyaltyTransaction
	err := r.db.Preload("Account").Preload("Account.Program").Where("tenant_id = ? AND id = ?", tenantID, transactionID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *GormRepository) UpdateTransaction(transaction *LoyaltyTransaction) (*LoyaltyTransaction, error) {
	err := r.db.Save(transaction).Error
	return transaction, err
}

func (r *GormRepository) ListTransactions(tenantID uuid.UUID, accountID *uuid.UUID, limit, offset int) ([]*LoyaltyTransaction, int64, error) {
	var transactions []*LoyaltyTransaction
	var total int64
	
	query := r.db.Model(&LoyaltyTransaction{}).Preload("Account").Preload("Account.Program").Where("tenant_id = ?", tenantID)
	
	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}
	
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error
	return transactions, total, err
}

// Loyalty Rewards
func (r *GormRepository) CreateReward(reward *LoyaltyReward) (*LoyaltyReward, error) {
	err := r.db.Create(reward).Error
	return reward, err
}

func (r *GormRepository) GetReward(tenantID, rewardID uuid.UUID) (*LoyaltyReward, error) {
	var reward LoyaltyReward
	err := r.db.Preload("Program").Where("tenant_id = ? AND id = ?", tenantID, rewardID).First(&reward).Error
	if err != nil {
		return nil, err
	}
	return &reward, nil
}

func (r *GormRepository) UpdateReward(reward *LoyaltyReward) (*LoyaltyReward, error) {
	err := r.db.Save(reward).Error
	return reward, err
}

func (r *GormRepository) DeleteReward(tenantID, rewardID uuid.UUID) error {
	return r.db.Where("tenant_id = ? AND id = ?", tenantID, rewardID).Delete(&LoyaltyReward{}).Error
}

func (r *GormRepository) ListRewards(tenantID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]*LoyaltyReward, int64, error) {
	var rewards []*LoyaltyReward
	var total int64

	query := r.db.Model(&LoyaltyReward{}).Preload("Program").Where("tenant_id = ?", tenantID)

	// Apply filters
	for key, value := range filters {
		switch key {
		case "status":
			query = query.Where("status = ?", value)
		case "type":
			query = query.Where("type = ?", value)
		case "program_id":
			query = query.Where("program_id = ?", value)
		case "search":
			query = query.Where("name ILIKE ? OR description ILIKE ?", fmt.Sprintf("%%%s%%", value), fmt.Sprintf("%%%s%%", value))
		}
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&rewards).Error
	return rewards, total, err
}

// Analytics
func (r *GormRepository) GetStats(tenantID uuid.UUID) (*LoyaltyStats, error) {
	stats := &LoyaltyStats{}

	// Total programs
	r.db.Model(&LoyaltyProgram{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalPrograms)

	// Active programs
	r.db.Model(&LoyaltyProgram{}).Where("tenant_id = ? AND status = 'active'", tenantID).Count(&stats.ActivePrograms)

	// Total accounts
	r.db.Model(&LoyaltyAccount{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalAccounts)

	// Active accounts
	r.db.Model(&LoyaltyAccount{}).Where("tenant_id = ? AND status = 'active'", tenantID).Count(&stats.ActiveAccounts)

	// Total points
	var totalPoints sql.NullInt64
	r.db.Model(&LoyaltyAccount{}).Where("tenant_id = ?", tenantID).Select("COALESCE(SUM(points), 0)").Scan(&totalPoints)
	stats.TotalPoints = totalPoints.Int64

	// Total redemptions
	r.db.Model(&LoyaltyTransaction{}).Where("tenant_id = ? AND type = 'redeemed'", tenantID).Count(&stats.TotalRedemptions)

	return stats, nil
}