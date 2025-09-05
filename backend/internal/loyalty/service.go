package loyalty

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service defines the loyalty service interface
type Service interface {
	// Program management
	CreateProgram(ctx context.Context, req *CreateProgramRequest) (*LoyaltyProgram, error)
	GetProgram(ctx context.Context, tenantID, programID uuid.UUID) (*LoyaltyProgram, error)
	UpdateProgram(ctx context.Context, tenantID, programID uuid.UUID, req *UpdateProgramRequest) (*LoyaltyProgram, error)
	DeleteProgram(ctx context.Context, tenantID, programID uuid.UUID) error
	ListPrograms(ctx context.Context, tenantID uuid.UUID, req *ListProgramsRequest) (*ListProgramsResponse, error)

	// Account management
	CreateAccount(ctx context.Context, req *CreateAccountRequest) (*LoyaltyAccount, error)
	GetAccount(ctx context.Context, tenantID, accountID uuid.UUID) (*LoyaltyAccount, error)
	UpdateAccount(ctx context.Context, tenantID, accountID uuid.UUID, req *UpdateAccountRequest) (*LoyaltyAccount, error)
	GetAccountByUser(ctx context.Context, tenantID, userID uuid.UUID) (*LoyaltyAccount, error)
	ListAccounts(ctx context.Context, tenantID uuid.UUID, req *ListAccountsRequest) (*ListAccountsResponse, error)

	// Points operations
	EarnPoints(ctx context.Context, req *EarnPointsRequest) (*LoyaltyTransaction, error)
	RedeemPoints(ctx context.Context, req *RedeemPointsRequest) (*LoyaltyTransaction, error)
	AdjustPoints(ctx context.Context, req *AdjustPointsRequest) (*LoyaltyTransaction, error)
	RedeemReward(ctx context.Context, tenantID, accountID, rewardID uuid.UUID) (*LoyaltyTransaction, error)

	// Transaction management
	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*LoyaltyTransaction, error)
	GetTransaction(ctx context.Context, tenantID, transactionID uuid.UUID) (*LoyaltyTransaction, error)
	ListTransactions(ctx context.Context, tenantID uuid.UUID, req *ListTransactionsRequest) (*ListTransactionsResponse, error)
	ListAccountTransactions(ctx context.Context, tenantID, accountID uuid.UUID, req *ListAccountTransactionsRequest) (*ListAccountTransactionsResponse, error)

	// Reward management
	CreateReward(ctx context.Context, req *CreateRewardRequest) (*LoyaltyReward, error)
	GetReward(ctx context.Context, tenantID, rewardID uuid.UUID) (*LoyaltyReward, error)
	UpdateReward(ctx context.Context, tenantID, rewardID uuid.UUID, req *UpdateRewardRequest) (*LoyaltyReward, error)
	DeleteReward(ctx context.Context, tenantID, rewardID uuid.UUID) error
	ListRewards(ctx context.Context, tenantID uuid.UUID, req *ListRewardsRequest) (*ListRewardsResponse, error)
	ListProgramRewards(ctx context.Context, tenantID, programID uuid.UUID, req *ListProgramRewardsRequest) (*ListProgramRewardsResponse, error)

	// Analytics
	GetStats(ctx context.Context, tenantID uuid.UUID) (*LoyaltyStats, error)
	GetAnalytics(ctx context.Context, tenantID uuid.UUID) (*LoyaltyAnalytics, error)
}

// ServiceImpl implements the loyalty service
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new loyalty service
func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

// Program management
func (s *ServiceImpl) CreateProgram(ctx context.Context, req *CreateProgramRequest) (*LoyaltyProgram, error) {
	program := &LoyaltyProgram{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Settings:    req.Settings,
		Status:      "active",
	}
	return s.repo.CreateProgram(program)
}

func (s *ServiceImpl) GetProgram(ctx context.Context, tenantID, programID uuid.UUID) (*LoyaltyProgram, error) {
	return s.repo.GetProgram(tenantID, programID)
}

func (s *ServiceImpl) UpdateProgram(ctx context.Context, tenantID, programID uuid.UUID, req *UpdateProgramRequest) (*LoyaltyProgram, error) {
	program, err := s.repo.GetProgram(tenantID, programID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		program.Name = *req.Name
	}
	if req.Description != nil {
		program.Description = *req.Description
	}
	if req.Type != nil {
		program.Type = *req.Type
	}
	if req.Settings != nil {
		program.Settings = *req.Settings
	}
	if req.Status != nil {
		program.Status = *req.Status
	}

	return s.repo.UpdateProgram(program)
}

func (s *ServiceImpl) DeleteProgram(ctx context.Context, tenantID, programID uuid.UUID) error {
	return s.repo.DeleteProgram(tenantID, programID)
}

func (s *ServiceImpl) ListPrograms(ctx context.Context, tenantID uuid.UUID, req *ListProgramsRequest) (*ListProgramsResponse, error) {
	programs, total, err := s.repo.ListPrograms(tenantID)
	if err != nil {
		return nil, err
	}
	return &ListProgramsResponse{
		Programs: programs,
		Total:    total,
	}, nil
}

// Account management
func (s *ServiceImpl) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*LoyaltyAccount, error) {
	account := &LoyaltyAccount{
		TenantID:  req.TenantID,
		UserID:    req.UserID,
		ProgramID: req.ProgramID,
		Points:    0,
		Tier:      "bronze",
		Status:    "active",
	}
	return s.repo.CreateAccount(account)
}

func (s *ServiceImpl) GetAccount(ctx context.Context, tenantID, accountID uuid.UUID) (*LoyaltyAccount, error) {
	return s.repo.GetAccount(tenantID, accountID)
}

func (s *ServiceImpl) GetAccountByUser(ctx context.Context, tenantID, userID uuid.UUID) (*LoyaltyAccount, error) {
	return s.repo.GetAccountByUser(tenantID, userID)
}

func (s *ServiceImpl) UpdateAccount(ctx context.Context, tenantID, accountID uuid.UUID, req *UpdateAccountRequest) (*LoyaltyAccount, error) {
	// Get existing account
	account, err := s.repo.GetAccount(tenantID, accountID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Status != nil {
		account.Status = *req.Status
	}
	if req.Tier != nil {
		account.Tier = *req.Tier
	}

	// Save updated account
	updatedAccount, err := s.repo.UpdateAccount(account)
	if err != nil {
		return nil, err
	}

	return updatedAccount, nil
}

func (s *ServiceImpl) ListAccounts(ctx context.Context, tenantID uuid.UUID, req *ListAccountsRequest) (*ListAccountsResponse, error) {
	accounts, total, err := s.repo.ListAccounts(tenantID, req.UserID, req.Status, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Convert []*LoyaltyAccount to []LoyaltyAccount
	accountList := make([]LoyaltyAccount, len(accounts))
	for i, account := range accounts {
		accountList[i] = *account
	}

	return &ListAccountsResponse{
		Accounts: accountList,
		Total:    total,
	}, nil
}

// Points operations
func (s *ServiceImpl) EarnPoints(ctx context.Context, req *EarnPointsRequest) (*LoyaltyTransaction, error) {
	// Get account
	account, err := s.repo.GetAccount(req.TenantID, req.AccountID)
	if err != nil {
		return nil, err
	}

	// Create transaction
	transaction := &LoyaltyTransaction{
		TenantID:    req.TenantID,
		AccountID:   req.AccountID,
		Type:        "earned",
		Points:      int(req.Points),
		Description: req.Reason,
		OrderID:     req.OrderID,
	}

	transaction, err = s.repo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// Update account points
	account.Points += int(req.Points)
	_, err = s.repo.UpdateAccount(account)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *ServiceImpl) RedeemPoints(ctx context.Context, req *RedeemPointsRequest) (*LoyaltyTransaction, error) {
	// Get account
	account, err := s.repo.GetAccount(req.TenantID, req.AccountID)
	if err != nil {
		return nil, err
	}

	// Check if account has enough points
	if account.Points < int(req.Points) {
		return nil, fmt.Errorf("insufficient points")
	}

	// Create transaction
	transaction := &LoyaltyTransaction{
		TenantID:    req.TenantID,
		AccountID:   req.AccountID,
		Type:        "redeemed",
		Points:      -int(req.Points),
		Description: req.Reason,
		OrderID:     req.OrderID,
	}

	transaction, err = s.repo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// Update account points
	account.Points -= int(req.Points)
	_, err = s.repo.UpdateAccount(account)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *ServiceImpl) RedeemReward(ctx context.Context, tenantID, accountID, rewardID uuid.UUID) (*LoyaltyTransaction, error) {
	// Get reward
	reward, err := s.repo.GetReward(tenantID, rewardID)
	if err != nil {
		return nil, err
	}

	if reward.Status != "active" {
		return nil, errors.New("loyalty reward is not active")
	}

	// Create redeem points request
	req := &RedeemPointsRequest{
		TenantID:  tenantID,
		AccountID: accountID,
		Points:      int64(reward.PointsCost),
		Reason:    fmt.Sprintf("Redeemed reward: %s", reward.Name),
	}

	// Redeem points for the reward
	transaction, err := s.RedeemPoints(ctx, req)
	if err != nil {
		return nil, err
	}

	// Note: RewardID field doesn't exist in LoyaltyTransaction, using Description instead
	transaction.Description = fmt.Sprintf("Redeemed reward: %s", reward.Name)

	return transaction, nil
}

// Transaction management
func (s *ServiceImpl) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*LoyaltyTransaction, error) {
	transaction := &LoyaltyTransaction{
		ID:          uuid.New(),
		TenantID:    req.TenantID,
		AccountID:   req.AccountID,
		Type:        req.Type,
		Points:      req.Points,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreateTransaction(transaction)
}

func (s *ServiceImpl) GetTransaction(ctx context.Context, tenantID, transactionID uuid.UUID) (*LoyaltyTransaction, error) {
	return s.repo.GetTransaction(tenantID, transactionID)
}

func (s *ServiceImpl) ListTransactions(ctx context.Context, tenantID uuid.UUID, req *ListTransactionsRequest) (*ListTransactionsResponse, error) {
	limit := 50
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}
	transactions, total, err := s.repo.ListTransactions(tenantID, req.AccountID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert []*LoyaltyTransaction to []LoyaltyTransaction
	transactionList := make([]LoyaltyTransaction, len(transactions))
	for i, t := range transactions {
		transactionList[i] = *t
	}

	return &ListTransactionsResponse{
		Transactions: transactionList,
		Total:        total,
	}, nil
}

func (s *ServiceImpl) ListAccountTransactions(ctx context.Context, tenantID, accountID uuid.UUID, req *ListAccountTransactionsRequest) (*ListAccountTransactionsResponse, error) {
	limit := 50
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}
	transactions, total, err := s.repo.ListTransactions(tenantID, &accountID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert []*LoyaltyTransaction to []LoyaltyTransaction
	transactionList := make([]LoyaltyTransaction, len(transactions))
	for i, t := range transactions {
		transactionList[i] = *t
	}

	return &ListAccountTransactionsResponse{
		Transactions: transactionList,
		Total:        total,
	}, nil
}

// Reward management
func (s *ServiceImpl) CreateReward(ctx context.Context, req *CreateRewardRequest) (*LoyaltyReward, error) {
	reward := &LoyaltyReward{
		TenantID:    req.TenantID,
		ProgramID:   req.ProgramID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		PointsCost:  req.PointsCost,
		Value:       req.Value,
		Status:      "active",
	}
	return s.repo.CreateReward(reward)
}

func (s *ServiceImpl) GetReward(ctx context.Context, tenantID, rewardID uuid.UUID) (*LoyaltyReward, error) {
	return s.repo.GetReward(tenantID, rewardID)
}

func (s *ServiceImpl) UpdateReward(ctx context.Context, tenantID, rewardID uuid.UUID, req *UpdateRewardRequest) (*LoyaltyReward, error) {
	reward, err := s.repo.GetReward(tenantID, rewardID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		reward.Name = *req.Name
	}
	if req.Description != nil {
		reward.Description = *req.Description
	}
	if req.Type != nil {
		reward.Type = *req.Type
	}
	if req.PointsCost != nil {
		reward.PointsCost = *req.PointsCost
	}
	if req.Value != nil {
		reward.Value = *req.Value
	}
	if req.Status != nil {
		reward.Status = *req.Status
	}

	return s.repo.UpdateReward(reward)
}

func (s *ServiceImpl) DeleteReward(ctx context.Context, tenantID, rewardID uuid.UUID) error {
	return s.repo.DeleteReward(tenantID, rewardID)
}

func (s *ServiceImpl) ListRewards(ctx context.Context, tenantID uuid.UUID, req *ListRewardsRequest) (*ListRewardsResponse, error) {
	limit := 50
	offset := 0
	if req.Limit > 0 {
		limit = req.Limit
	}
	if req.Offset > 0 {
		offset = req.Offset
	}
	
	// Build filters map
	filters := make(map[string]interface{})
	if req.ProgramID != nil {
		filters["program_id"] = *req.ProgramID
	}
	if req.Type != nil {
		filters["type"] = *req.Type
	}
	if req.Status != nil {
		filters["status"] = *req.Status
	}
	if req.Search != nil {
		filters["search"] = *req.Search
	}
	
	rewards, total, err := s.repo.ListRewards(tenantID, filters, limit, offset)
	if err != nil {
		return nil, err
	}
	return &ListRewardsResponse{
		Rewards: rewards,
		Total:   total,
	}, nil
}

// Points operations
func (s *ServiceImpl) AdjustPoints(ctx context.Context, req *AdjustPointsRequest) (*LoyaltyTransaction, error) {
	// Create transaction
	transaction := &LoyaltyTransaction{
		TenantID:    req.TenantID,
		AccountID:   req.AccountID,
		Type:        "adjusted",
		Points:      int(req.Points),
		Description: req.Reason,
	}

	transaction, err := s.repo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// Update account points
	account, err := s.repo.GetAccount(req.TenantID, req.AccountID)
	if err != nil {
		return nil, err
	}

	account.Points += int(req.Points)
	_, err = s.repo.UpdateAccount(account)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Analytics
func (s *ServiceImpl) GetStats(ctx context.Context, tenantID uuid.UUID) (*LoyaltyStats, error) {
	return s.repo.GetStats(tenantID)
}

func (s *ServiceImpl) GetAnalytics(ctx context.Context, tenantID uuid.UUID) (*LoyaltyAnalytics, error) {
	// TODO: Implement analytics calculation
	// For now, return empty analytics
	analytics := &LoyaltyAnalytics{
		TotalPrograms:       0,
		TotalAccounts:       0,
		TotalTransactions:   0,
		TotalRewards:        0,
		TotalPointsEarned:   0,
		TotalPointsRedeemed: 0,
	}

	return analytics, nil
}

func (s *ServiceImpl) ListProgramRewards(ctx context.Context, tenantID, programID uuid.UUID, req *ListProgramRewardsRequest) (*ListProgramRewardsResponse, error) {
	limit := 50
	offset := 0
	if req.Limit > 0 {
		limit = req.Limit
	}
	if req.Offset > 0 {
		offset = req.Offset
	}
	filters := map[string]interface{}{
		"program_id": programID,
	}
	rewards, total, err := s.repo.ListRewards(tenantID, filters, limit, offset)
	if err != nil {
		return nil, err
	}

	return &ListProgramRewardsResponse{
		Rewards: rewards,
		Total:   total,
	}, nil
}