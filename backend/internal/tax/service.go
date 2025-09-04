package tax

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for tax business logic
type Service interface {
	// Tax Rule operations
	CreateTaxRule(ctx context.Context, tenantID uuid.UUID, req CreateTaxRuleRequest) (*TaxRuleResponse, error)
	GetTaxRule(ctx context.Context, tenantID, ruleID uuid.UUID) (*TaxRuleResponse, error)
	GetTaxRuleByCode(ctx context.Context, tenantID uuid.UUID, code string) (*TaxRuleResponse, error)
	UpdateTaxRule(ctx context.Context, tenantID, ruleID uuid.UUID, req UpdateTaxRuleRequest) (*TaxRuleResponse, error)
	DeleteTaxRule(ctx context.Context, tenantID, ruleID uuid.UUID) error
	ListTaxRules(ctx context.Context, tenantID uuid.UUID, filter TaxRuleFilter, page, pageSize int) ([]*TaxRuleResponse, int64, error)
	GetActiveTaxRules(ctx context.Context, tenantID uuid.UUID, date time.Time) ([]*TaxRuleResponse, error)
	
	// Tax Rate operations
	CreateTaxRate(ctx context.Context, tenantID uuid.UUID, req CreateTaxRateRequest) (*TaxRateResponse, error)
	GetTaxRate(ctx context.Context, tenantID, rateID uuid.UUID) (*TaxRateResponse, error)
	UpdateTaxRate(ctx context.Context, tenantID, rateID uuid.UUID, req UpdateTaxRateRequest) (*TaxRateResponse, error)
	DeleteTaxRate(ctx context.Context, tenantID, rateID uuid.UUID) error
	ListTaxRates(ctx context.Context, tenantID uuid.UUID, filter TaxRateFilter, page, pageSize int) ([]*TaxRateResponse, int64, error)
	GetTaxRatesByRule(ctx context.Context, tenantID, ruleID uuid.UUID) ([]*TaxRateResponse, error)
	
	// Tax Calculation operations
	CalculateTax(ctx context.Context, tenantID uuid.UUID, req TaxCalculationRequest) (*TaxCalculationResponse, error)
	GetTaxCalculation(ctx context.Context, tenantID, taxID uuid.UUID) (*TaxCalculationResponse, error)
	ListTaxCalculations(ctx context.Context, tenantID uuid.UUID, filter TaxFilter, page, pageSize int) ([]*TaxCalculationResponse, int64, error)
	GetTaxCalculationsByOrder(ctx context.Context, tenantID, orderID uuid.UUID) ([]*TaxCalculationResponse, error)
	GetTaxCalculationsByProduct(ctx context.Context, tenantID, productID uuid.UUID) ([]*TaxCalculationResponse, error)
	GetTaxCalculationsByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*TaxCalculationResponse, error)
	
	// Bulk operations
	BulkCreateTaxRules(ctx context.Context, tenantID uuid.UUID, requests []CreateTaxRuleRequest) ([]*TaxRuleResponse, error)
	BulkUpdateTaxRuleStatus(ctx context.Context, tenantID uuid.UUID, ruleIDs []uuid.UUID, status string) error
	BulkDeleteTaxRules(ctx context.Context, tenantID uuid.UUID, ruleIDs []uuid.UUID) error
	BulkCreateTaxRates(ctx context.Context, tenantID uuid.UUID, requests []CreateTaxRateRequest) ([]*TaxRateResponse, error)
	BulkUpdateTaxRateStatus(ctx context.Context, tenantID uuid.UUID, rateIDs []uuid.UUID, isActive bool) error
	
	// Validation operations
	ValidateTaxRule(ctx context.Context, tenantID uuid.UUID, req CreateTaxRuleRequest) error
	ValidateTaxRate(ctx context.Context, tenantID uuid.UUID, req CreateTaxRateRequest) error
	ValidateTaxCalculation(ctx context.Context, tenantID uuid.UUID, req TaxCalculationRequest) error
	
	// Statistics and analytics
	GetTaxStats(ctx context.Context, tenantID uuid.UUID) (*TaxStats, error)
	GetTaxStatsByLocation(ctx context.Context, tenantID uuid.UUID, limit int) ([]*TaxByLocation, error)
	GetTaxStatsByType(ctx context.Context, tenantID uuid.UUID) ([]*TaxByType, error)
	GetTaxTrends(ctx context.Context, tenantID uuid.UUID, days int) (map[string]interface{}, error)
	
	// Maintenance operations
	CleanupExpiredRules(ctx context.Context, tenantID uuid.UUID) (int64, error)
	CleanupExpiredRates(ctx context.Context, tenantID uuid.UUID) (int64, error)
	ArchiveOldCalculations(ctx context.Context, tenantID uuid.UUID, olderThan time.Time) (int64, error)
	
	// Utility operations
	GetApplicableTaxRules(ctx context.Context, tenantID uuid.UUID, req TaxCalculationRequest) ([]*TaxRuleResponse, error)
	PreviewTaxCalculation(ctx context.Context, tenantID uuid.UUID, req TaxCalculationRequest) (*TaxCalculationResponse, error)
	ValidateLocation(ctx context.Context, country, state, city, postalCode string) error
	GetSupportedLocations(ctx context.Context, tenantID uuid.UUID) (map[string]interface{}, error)
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new tax service
func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

// Tax Rule operations

// CreateTaxRule creates a new tax rule
func (s *ServiceImpl) CreateTaxRule(ctx context.Context, tenantID uuid.UUID, req CreateTaxRuleRequest) (*TaxRuleResponse, error) {
	// Validate request
	if err := s.ValidateTaxRule(ctx, tenantID, req); err != nil {
		return nil, err
	}
	
	// Check if code already exists
	exists, err := s.repo.CheckTaxRuleCodeExists(ctx, tenantID, req.Code, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check code existence: %w", err)
	}
	if exists {
		return nil, ErrRuleCodeExists
	}
	
	// Create tax rule
	rule := &TaxRule{
		TenantID:       tenantID,
		Name:           req.Name,
		Description:    req.Description,
		Code:           req.Code,
		Type:           req.Type,
		Status:         StatusActive,
		TaxType:        req.TaxType,
		Rate:           req.Rate,
		Method:         req.Method,
		Priority:       req.Priority,
		Conditions:     req.Conditions,
		IsCompound:     req.IsCompound,
		IsInclusive:    req.IsInclusive,
		ValidFrom:      req.ValidFrom,
		ValidTo:        req.ValidTo,
		Countries:      req.Countries,
		States:         req.States,
		Cities:         req.Cities,
		PostalCodes:    req.PostalCodes,
		ProductIDs:     req.ProductIDs,
		CategoryIDs:    req.CategoryIDs,
		CustomerIDs:    req.CustomerIDs,
		CustomerGroups: req.CustomerGroups,
		MinAmount:      req.MinAmount,
		MaxAmount:      req.MaxAmount,
	}
	
	// Set defaults
	if rule.Method == "" {
		rule.Method = MethodExclusive
	}
	
	if err := s.repo.CreateTaxRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to create tax rule: %w", err)
	}
	
	return s.buildTaxRuleResponse(rule), nil
}

// GetTaxRule retrieves a tax rule by ID
func (s *ServiceImpl) GetTaxRule(ctx context.Context, tenantID, ruleID uuid.UUID) (*TaxRuleResponse, error) {
	rule, err := s.repo.GetTaxRule(ctx, tenantID, ruleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rule: %w", err)
	}
	return s.buildTaxRuleResponse(rule), nil
}

// GetTaxRuleByCode retrieves a tax rule by code
func (s *ServiceImpl) GetTaxRuleByCode(ctx context.Context, tenantID uuid.UUID, code string) (*TaxRuleResponse, error) {
	rule, err := s.repo.GetTaxRuleByCode(ctx, tenantID, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rule by code: %w", err)
	}
	return s.buildTaxRuleResponse(rule), nil
}

// UpdateTaxRule updates a tax rule
func (s *ServiceImpl) UpdateTaxRule(ctx context.Context, tenantID, ruleID uuid.UUID, req UpdateTaxRuleRequest) (*TaxRuleResponse, error) {
	// Get existing rule
	rule, err := s.repo.GetTaxRule(ctx, tenantID, ruleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rule: %w", err)
	}
	
	// Update fields
	if req.Name != nil {
		rule.Name = *req.Name
	}
	if req.Description != nil {
		rule.Description = *req.Description
	}
	if req.Status != nil {
		rule.Status = *req.Status
	}
	if req.TaxType != nil {
		rule.TaxType = *req.TaxType
	}
	if req.Rate != nil {
		rule.Rate = *req.Rate
	}
	if req.Method != nil {
		rule.Method = *req.Method
	}
	if req.Priority != nil {
		rule.Priority = *req.Priority
	}
	if req.Conditions != nil {
		rule.Conditions = *req.Conditions
	}
	if req.IsCompound != nil {
		rule.IsCompound = *req.IsCompound
	}
	if req.IsInclusive != nil {
		rule.IsInclusive = *req.IsInclusive
	}
	if req.ValidFrom != nil {
		rule.ValidFrom = req.ValidFrom
	}
	if req.ValidTo != nil {
		rule.ValidTo = req.ValidTo
	}
	if req.Countries != nil {
		rule.Countries = req.Countries
	}
	if req.States != nil {
		rule.States = req.States
	}
	if req.Cities != nil {
		rule.Cities = req.Cities
	}
	if req.PostalCodes != nil {
		rule.PostalCodes = req.PostalCodes
	}
	if req.ProductIDs != nil {
		rule.ProductIDs = req.ProductIDs
	}
	if req.CategoryIDs != nil {
		rule.CategoryIDs = req.CategoryIDs
	}
	if req.CustomerIDs != nil {
		rule.CustomerIDs = req.CustomerIDs
	}
	if req.CustomerGroups != nil {
		rule.CustomerGroups = req.CustomerGroups
	}
	if req.MinAmount != nil {
		rule.MinAmount = req.MinAmount
	}
	if req.MaxAmount != nil {
		rule.MaxAmount = req.MaxAmount
	}
	
	if err := s.repo.UpdateTaxRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to update tax rule: %w", err)
	}
	
	return s.buildTaxRuleResponse(rule), nil
}

// DeleteTaxRule deletes a tax rule
func (s *ServiceImpl) DeleteTaxRule(ctx context.Context, tenantID, ruleID uuid.UUID) error {
	// Get existing rule to check if it can be deleted
	rule, err := s.repo.GetTaxRule(ctx, tenantID, ruleID)
	if err != nil {
		return fmt.Errorf("failed to get tax rule: %w", err)
	}
	
	if !rule.CanDelete() {
		return ErrCannotDeleteRule
	}
	
	return s.repo.DeleteTaxRule(ctx, tenantID, ruleID)
}

// ListTaxRules retrieves tax rules with filtering and pagination
func (s *ServiceImpl) ListTaxRules(ctx context.Context, tenantID uuid.UUID, filter TaxRuleFilter, page, pageSize int) ([]*TaxRuleResponse, int64, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}
	
	rules, total, err := s.repo.ListTaxRules(ctx, tenantID, filter, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tax rules: %w", err)
	}
	
	responses := make([]*TaxRuleResponse, len(rules))
	for i, rule := range rules {
		responses[i] = s.buildTaxRuleResponse(rule)
	}
	
	return responses, total, nil
}

// GetActiveTaxRules retrieves active tax rules for a specific date
func (s *ServiceImpl) GetActiveTaxRules(ctx context.Context, tenantID uuid.UUID, date time.Time) ([]*TaxRuleResponse, error) {
	rules, err := s.repo.GetActiveTaxRules(ctx, tenantID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get active tax rules: %w", err)
	}
	
	responses := make([]*TaxRuleResponse, len(rules))
	for i, rule := range rules {
		responses[i] = s.buildTaxRuleResponse(rule)
	}
	
	return responses, nil
}

// Tax Rate operations

// CreateTaxRate creates a new tax rate
func (s *ServiceImpl) CreateTaxRate(ctx context.Context, tenantID uuid.UUID, req CreateTaxRateRequest) (*TaxRateResponse, error) {
	// Validate request
	if err := s.ValidateTaxRate(ctx, tenantID, req); err != nil {
		return nil, err
	}
	
	// Verify rule exists
	_, err := s.repo.GetTaxRule(ctx, tenantID, req.RuleID)
	if err != nil {
		return nil, fmt.Errorf("tax rule not found: %w", err)
	}
	
	// Create tax rate
	rate := &TaxRate{
		TenantID:    tenantID,
		RuleID:      req.RuleID,
		Name:        req.Name,
		Rate:        req.Rate,
		TaxType:     req.TaxType,
		Description: req.Description,
		Country:     req.Country,
		State:       req.State,
		City:        req.City,
		PostalCode:  req.PostalCode,
		ValidFrom:   req.ValidFrom,
		ValidTo:     req.ValidTo,
		IsActive:    true,
	}
	
	if err := s.repo.CreateTaxRate(ctx, rate); err != nil {
		return nil, fmt.Errorf("failed to create tax rate: %w", err)
	}
	
	return s.buildTaxRateResponse(rate), nil
}

// GetTaxRate retrieves a tax rate by ID
func (s *ServiceImpl) GetTaxRate(ctx context.Context, tenantID, rateID uuid.UUID) (*TaxRateResponse, error) {
	rate, err := s.repo.GetTaxRate(ctx, tenantID, rateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rate: %w", err)
	}
	return s.buildTaxRateResponse(rate), nil
}

// UpdateTaxRate updates a tax rate
func (s *ServiceImpl) UpdateTaxRate(ctx context.Context, tenantID, rateID uuid.UUID, req UpdateTaxRateRequest) (*TaxRateResponse, error) {
	// Get existing rate
	rate, err := s.repo.GetTaxRate(ctx, tenantID, rateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rate: %w", err)
	}
	
	// Update fields
	if req.Name != nil {
		rate.Name = *req.Name
	}
	if req.Rate != nil {
		rate.Rate = *req.Rate
	}
	if req.TaxType != nil {
		rate.TaxType = *req.TaxType
	}
	if req.Description != nil {
		rate.Description = *req.Description
	}
	if req.State != nil {
		rate.State = *req.State
	}
	if req.City != nil {
		rate.City = *req.City
	}
	if req.PostalCode != nil {
		rate.PostalCode = *req.PostalCode
	}
	if req.ValidFrom != nil {
		rate.ValidFrom = req.ValidFrom
	}
	if req.ValidTo != nil {
		rate.ValidTo = req.ValidTo
	}
	if req.IsActive != nil {
		rate.IsActive = *req.IsActive
	}
	
	if err := s.repo.UpdateTaxRate(ctx, rate); err != nil {
		return nil, fmt.Errorf("failed to update tax rate: %w", err)
	}
	
	return s.buildTaxRateResponse(rate), nil
}

// DeleteTaxRate deletes a tax rate
func (s *ServiceImpl) DeleteTaxRate(ctx context.Context, tenantID, rateID uuid.UUID) error {
	return s.repo.DeleteTaxRate(ctx, tenantID, rateID)
}

// ListTaxRates retrieves tax rates with filtering and pagination
func (s *ServiceImpl) ListTaxRates(ctx context.Context, tenantID uuid.UUID, filter TaxRateFilter, page, pageSize int) ([]*TaxRateResponse, int64, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}
	
	rates, total, err := s.repo.ListTaxRates(ctx, tenantID, filter, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tax rates: %w", err)
	}
	
	responses := make([]*TaxRateResponse, len(rates))
	for i, rate := range rates {
		responses[i] = s.buildTaxRateResponse(rate)
	}
	
	return responses, total, nil
}

// GetTaxRatesByRule retrieves tax rates for a specific rule
func (s *ServiceImpl) GetTaxRatesByRule(ctx context.Context, tenantID, ruleID uuid.UUID) ([]*TaxRateResponse, error) {
	rates, err := s.repo.GetTaxRatesByRule(ctx, tenantID, ruleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rates by rule: %w", err)
	}
	
	responses := make([]*TaxRateResponse, len(rates))
	for i, rate := range rates {
		responses[i] = s.buildTaxRateResponse(rate)
	}
	
	return responses, nil
}

// Tax Calculation operations

// CalculateTax calculates tax for a given request
func (s *ServiceImpl) CalculateTax(ctx context.Context, tenantID uuid.UUID, req TaxCalculationRequest) (*TaxCalculationResponse, error) {
	// Validate request
	if err := s.ValidateTaxCalculation(ctx, tenantID, req); err != nil {
		return nil, err
	}
	
	// Get applicable tax rules
	rules, err := s.repo.GetApplicableTaxRules(ctx, tenantID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get applicable tax rules: %w", err)
	}
	
	if len(rules) == 0 {
		return nil, ErrNoApplicableRules
	}
	
	// Sort rules by priority (highest first)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})
	
	// Calculate tax
	calcResult, err := s.performTaxCalculation(req, rules)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate tax: %w", err)
	}
	
	// Store tax calculation
	tax := &Tax{
		TenantID:      tenantID,
		OrderID:       nil, // Set when called from order context
		ProductID:     req.ProductID,
		CustomerID:    req.CustomerID,
		TaxableAmount: calcResult.TaxableAmount,
		TaxAmount:     calcResult.TaxAmount,
		TaxRate:       calcResult.EffectiveRate,
		TaxType:       TaxTypePercentage, // Default, could be determined from rules
		Method:        req.Method,
		Country:       req.Country,
		State:         req.State,
		City:          req.City,
		PostalCode:    req.PostalCode,
		CalculatedAt:  time.Now(),
	}
	
	// Set default method if not specified
	if tax.Method == "" {
		tax.Method = MethodExclusive
	}
	
	if err := s.repo.CreateTax(ctx, tax); err != nil {
		return nil, fmt.Errorf("failed to store tax calculation: %w", err)
	}
	
	// Store applied rules
	for _, appliedRule := range calcResult.AppliedRules {
		application := &TaxRuleApplication{
			TenantID:      tenantID,
			TaxID:         tax.ID,
			RuleID:        appliedRule.RuleID,
			RuleName:      appliedRule.RuleName,
			RuleCode:      appliedRule.RuleCode,
			AppliedRate:   appliedRule.AppliedRate,
			TaxableAmount: appliedRule.TaxableAmount,
			TaxAmount:     appliedRule.TaxAmount,
			Priority:      appliedRule.Priority,
		}
		if err := s.repo.CreateTaxRuleApplication(ctx, application); err != nil {
			return nil, fmt.Errorf("failed to store tax rule application: %w", err)
		}
	}
	
	return calcResult, nil
}

// GetTaxCalculation retrieves a tax calculation by ID
func (s *ServiceImpl) GetTaxCalculation(ctx context.Context, tenantID, taxID uuid.UUID) (*TaxCalculationResponse, error) {
	tax, err := s.repo.GetTax(ctx, tenantID, taxID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax calculation: %w", err)
	}
	return s.buildTaxCalculationResponse(tax), nil
}

// ListTaxCalculations retrieves tax calculations with filtering and pagination
func (s *ServiceImpl) ListTaxCalculations(ctx context.Context, tenantID uuid.UUID, filter TaxFilter, page, pageSize int) ([]*TaxCalculationResponse, int64, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}
	
	taxes, total, err := s.repo.ListTaxes(ctx, tenantID, filter, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tax calculations: %w", err)
	}
	
	responses := make([]*TaxCalculationResponse, len(taxes))
	for i, tax := range taxes {
		responses[i] = s.buildTaxCalculationResponse(tax)
	}
	
	return responses, total, nil
}

// GetTaxCalculationsByOrder retrieves tax calculations for a specific order
func (s *ServiceImpl) GetTaxCalculationsByOrder(ctx context.Context, tenantID, orderID uuid.UUID) ([]*TaxCalculationResponse, error) {
	taxes, err := s.repo.GetTaxesByOrder(ctx, tenantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax calculations by order: %w", err)
	}
	
	responses := make([]*TaxCalculationResponse, len(taxes))
	for i, tax := range taxes {
		responses[i] = s.buildTaxCalculationResponse(tax)
	}
	
	return responses, nil
}

// GetTaxCalculationsByProduct retrieves tax calculations for a specific product
func (s *ServiceImpl) GetTaxCalculationsByProduct(ctx context.Context, tenantID, productID uuid.UUID) ([]*TaxCalculationResponse, error) {
	taxes, err := s.repo.GetTaxesByProduct(ctx, tenantID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax calculations by product: %w", err)
	}
	
	responses := make([]*TaxCalculationResponse, len(taxes))
	for i, tax := range taxes {
		responses[i] = s.buildTaxCalculationResponse(tax)
	}
	
	return responses, nil
}

// GetTaxCalculationsByCustomer retrieves tax calculations for a specific customer
func (s *ServiceImpl) GetTaxCalculationsByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*TaxCalculationResponse, error) {
	taxes, err := s.repo.GetTaxesByCustomer(ctx, tenantID, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax calculations by customer: %w", err)
	}
	
	responses := make([]*TaxCalculationResponse, len(taxes))
	for i, tax := range taxes {
		responses[i] = s.buildTaxCalculationResponse(tax)
	}
	
	return responses, nil
}

// Bulk operations

// BulkCreateTaxRules creates multiple tax rules
func (s *ServiceImpl) BulkCreateTaxRules(ctx context.Context, tenantID uuid.UUID, requests []CreateTaxRuleRequest) ([]*TaxRuleResponse, error) {
	if len(requests) == 0 {
		return []*TaxRuleResponse{}, nil
	}
	if len(requests) > 100 {
		return nil, ErrBulkSizeExceeded
	}
	
	rules := make([]*TaxRule, len(requests))
	for i, req := range requests {
		// Validate each request
		if err := s.ValidateTaxRule(ctx, tenantID, req); err != nil {
			return nil, fmt.Errorf("validation failed for rule %d: %w", i, err)
		}
		
		// Check if code already exists
		exists, err := s.repo.CheckTaxRuleCodeExists(ctx, tenantID, req.Code, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check code existence for rule %d: %w", i, err)
		}
		if exists {
			return nil, fmt.Errorf("rule code '%s' already exists", req.Code)
		}
		
		rules[i] = &TaxRule{
			TenantID:       tenantID,
			Name:           req.Name,
			Description:    req.Description,
			Code:           req.Code,
			Type:           req.Type,
			Status:         StatusActive,
			TaxType:        req.TaxType,
			Rate:           req.Rate,
			Method:         req.Method,
			Priority:       req.Priority,
			Conditions:     req.Conditions,
			IsCompound:     req.IsCompound,
			IsInclusive:    req.IsInclusive,
			ValidFrom:      req.ValidFrom,
			ValidTo:        req.ValidTo,
			Countries:      req.Countries,
			States:         req.States,
			Cities:         req.Cities,
			PostalCodes:    req.PostalCodes,
			ProductIDs:     req.ProductIDs,
			CategoryIDs:    req.CategoryIDs,
			CustomerIDs:    req.CustomerIDs,
			CustomerGroups: req.CustomerGroups,
			MinAmount:      req.MinAmount,
			MaxAmount:      req.MaxAmount,
		}
		
		// Set defaults
		if rules[i].Method == "" {
			rules[i].Method = MethodExclusive
		}
	}
	
	if err := s.repo.BulkCreateTaxRules(ctx, rules); err != nil {
		return nil, fmt.Errorf("failed to bulk create tax rules: %w", err)
	}
	
	responses := make([]*TaxRuleResponse, len(rules))
	for i, rule := range rules {
		responses[i] = s.buildTaxRuleResponse(rule)
	}
	
	return responses, nil
}

// BulkUpdateTaxRuleStatus updates status for multiple tax rules
func (s *ServiceImpl) BulkUpdateTaxRuleStatus(ctx context.Context, tenantID uuid.UUID, ruleIDs []uuid.UUID, status string) error {
	if !ValidateStatus(status) {
		return fmt.Errorf("invalid status: %s", status)
	}
	return s.repo.BulkUpdateTaxRuleStatus(ctx, tenantID, ruleIDs, status)
}

// BulkDeleteTaxRules deletes multiple tax rules
func (s *ServiceImpl) BulkDeleteTaxRules(ctx context.Context, tenantID uuid.UUID, ruleIDs []uuid.UUID) error {
	return s.repo.BulkDeleteTaxRules(ctx, tenantID, ruleIDs)
}

// BulkCreateTaxRates creates multiple tax rates
func (s *ServiceImpl) BulkCreateTaxRates(ctx context.Context, tenantID uuid.UUID, requests []CreateTaxRateRequest) ([]*TaxRateResponse, error) {
	if len(requests) == 0 {
		return []*TaxRateResponse{}, nil
	}
	if len(requests) > 100 {
		return nil, ErrBulkSizeExceeded
	}
	
	rates := make([]*TaxRate, len(requests))
	for i, req := range requests {
		// Validate each request
		if err := s.ValidateTaxRate(ctx, tenantID, req); err != nil {
			return nil, fmt.Errorf("validation failed for rate %d: %w", i, err)
		}
		
		// Verify rule exists
		_, err := s.repo.GetTaxRule(ctx, tenantID, req.RuleID)
		if err != nil {
			return nil, fmt.Errorf("tax rule not found for rate %d: %w", i, err)
		}
		
		rates[i] = &TaxRate{
			TenantID:    tenantID,
			RuleID:      req.RuleID,
			Name:        req.Name,
			Rate:        req.Rate,
			TaxType:     req.TaxType,
			Description: req.Description,
			Country:     req.Country,
			State:       req.State,
			City:        req.City,
			PostalCode:  req.PostalCode,
			ValidFrom:   req.ValidFrom,
			ValidTo:     req.ValidTo,
			IsActive:    true,
		}
	}
	
	if err := s.repo.BulkCreateTaxRates(ctx, rates); err != nil {
		return nil, fmt.Errorf("failed to bulk create tax rates: %w", err)
	}
	
	responses := make([]*TaxRateResponse, len(rates))
	for i, rate := range rates {
		responses[i] = s.buildTaxRateResponse(rate)
	}
	
	return responses, nil
}

// BulkUpdateTaxRateStatus updates status for multiple tax rates
func (s *ServiceImpl) BulkUpdateTaxRateStatus(ctx context.Context, tenantID uuid.UUID, rateIDs []uuid.UUID, isActive bool) error {
	return s.repo.BulkUpdateTaxRateStatus(ctx, tenantID, rateIDs, isActive)
}

// Validation operations

// ValidateTaxRule validates a tax rule request
func (s *ServiceImpl) ValidateTaxRule(ctx context.Context, tenantID uuid.UUID, req CreateTaxRuleRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return fmt.Errorf("code is required")
	}
	if !ValidateRuleType(req.Type) {
		return fmt.Errorf("invalid rule type: %s", req.Type)
	}
	if !ValidateTaxType(req.TaxType) {
		return fmt.Errorf("invalid tax type: %s", req.TaxType)
	}
	if req.Rate < 0 {
		return ErrInvalidRate
	}
	if req.Method != "" && !ValidateMethod(req.Method) {
		return fmt.Errorf("invalid method: %s", req.Method)
	}
	if req.ValidFrom != nil && req.ValidTo != nil && req.ValidFrom.After(*req.ValidTo) {
		return ErrInvalidDateRange
	}
	if req.MinAmount != nil && req.MaxAmount != nil && *req.MinAmount > *req.MaxAmount {
		return fmt.Errorf("min amount cannot be greater than max amount")
	}
	return nil
}

// ValidateTaxRate validates a tax rate request
func (s *ServiceImpl) ValidateTaxRate(ctx context.Context, tenantID uuid.UUID, req CreateTaxRateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if req.Rate < 0 {
		return ErrInvalidRate
	}
	if !ValidateTaxType(req.TaxType) {
		return fmt.Errorf("invalid tax type: %s", req.TaxType)
	}
	if len(req.Country) != 2 {
		return fmt.Errorf("country must be a 2-letter code")
	}
	if req.ValidFrom != nil && req.ValidTo != nil && req.ValidFrom.After(*req.ValidTo) {
		return ErrInvalidDateRange
	}
	return nil
}

// ValidateTaxCalculation validates a tax calculation request
func (s *ServiceImpl) ValidateTaxCalculation(ctx context.Context, tenantID uuid.UUID, req TaxCalculationRequest) error {
	if req.Amount < 0 {
		return ErrInvalidAmount
	}
	if len(req.Country) != 2 {
		return ErrInvalidLocation
	}
	if req.Method != "" && !ValidateMethod(req.Method) {
		return fmt.Errorf("invalid method: %s", req.Method)
	}
	return nil
}

// Statistics and analytics

// GetTaxStats retrieves tax statistics
func (s *ServiceImpl) GetTaxStats(ctx context.Context, tenantID uuid.UUID) (*TaxStats, error) {
	return s.repo.GetTaxStats(ctx, tenantID)
}

// GetTaxStatsByLocation retrieves tax statistics by location
func (s *ServiceImpl) GetTaxStatsByLocation(ctx context.Context, tenantID uuid.UUID, limit int) ([]*TaxByLocation, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	return s.repo.GetTaxStatsByLocation(ctx, tenantID, limit)
}

// GetTaxStatsByType retrieves tax statistics by type
func (s *ServiceImpl) GetTaxStatsByType(ctx context.Context, tenantID uuid.UUID) ([]*TaxByType, error) {
	return s.repo.GetTaxStatsByType(ctx, tenantID)
}

// GetTaxTrends retrieves tax trends over time
func (s *ServiceImpl) GetTaxTrends(ctx context.Context, tenantID uuid.UUID, days int) (map[string]interface{}, error) {
	if days <= 0 || days > 365 {
		days = 30
	}
	return s.repo.GetTaxTrends(ctx, tenantID, days)
}

// Maintenance operations

// CleanupExpiredRules cleans up expired tax rules
func (s *ServiceImpl) CleanupExpiredRules(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return s.repo.CleanupExpiredRules(ctx, tenantID)
}

// CleanupExpiredRates cleans up expired tax rates
func (s *ServiceImpl) CleanupExpiredRates(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return s.repo.CleanupExpiredRates(ctx, tenantID)
}

// ArchiveOldCalculations archives old tax calculations
func (s *ServiceImpl) ArchiveOldCalculations(ctx context.Context, tenantID uuid.UUID, olderThan time.Time) (int64, error) {
	return s.repo.ArchiveOldTaxCalculations(ctx, tenantID, olderThan)
}

// Utility operations

// GetApplicableTaxRules retrieves applicable tax rules for a calculation request
func (s *ServiceImpl) GetApplicableTaxRules(ctx context.Context, tenantID uuid.UUID, req TaxCalculationRequest) ([]*TaxRuleResponse, error) {
	rules, err := s.repo.GetApplicableTaxRules(ctx, tenantID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get applicable tax rules: %w", err)
	}
	
	responses := make([]*TaxRuleResponse, len(rules))
	for i, rule := range rules {
		responses[i] = s.buildTaxRuleResponse(rule)
	}
	
	return responses, nil
}

// PreviewTaxCalculation previews tax calculation without storing it
func (s *ServiceImpl) PreviewTaxCalculation(ctx context.Context, tenantID uuid.UUID, req TaxCalculationRequest) (*TaxCalculationResponse, error) {
	// Validate request
	if err := s.ValidateTaxCalculation(ctx, tenantID, req); err != nil {
		return nil, err
	}
	
	// Get applicable tax rules
	rules, err := s.repo.GetApplicableTaxRules(ctx, tenantID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get applicable tax rules: %w", err)
	}
	
	if len(rules) == 0 {
		return nil, ErrNoApplicableRules
	}
	
	// Sort rules by priority (highest first)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})
	
	// Calculate tax without storing
	return s.performTaxCalculation(req, rules)
}

// ValidateLocation validates a location
func (s *ServiceImpl) ValidateLocation(ctx context.Context, country, state, city, postalCode string) error {
	if len(country) != 2 {
		return fmt.Errorf("country must be a 2-letter code")
	}
	// Additional location validation logic could be added here
	return nil
}

// GetSupportedLocations retrieves supported locations
func (s *ServiceImpl) GetSupportedLocations(ctx context.Context, tenantID uuid.UUID) (map[string]interface{}, error) {
	// This would typically return a list of supported countries, states, etc.
	// For now, return a basic structure
	return map[string]interface{}{
		"countries": []string{"US", "CA", "GB", "AU", "DE", "FR"},
		"message":   "Location validation is basic. Extend as needed.",
	}, nil
}

// Helper methods

// performTaxCalculation performs the actual tax calculation
func (s *ServiceImpl) performTaxCalculation(req TaxCalculationRequest, rules []*TaxRule) (*TaxCalculationResponse, error) {
	method := req.Method
	if method == "" {
		method = MethodExclusive
	}
	
	taxableAmount := req.Amount
	totalTaxAmount := 0.0
	appliedRules := []AppliedTaxRuleResponse{}
	
	// Handle inclusive vs exclusive calculation
	if method == MethodInclusive {
		// For inclusive, we need to extract tax from the total amount
		totalRate := 0.0
		for _, rule := range rules {
			if rule.IsValidForAmount(req.Amount) {
				totalRate += rule.Rate
			}
		}
		
		if totalRate > 0 {
			taxableAmount = req.Amount / (1 + totalRate/100)
		}
	}
	
	// Apply each rule
	for _, rule := range rules {
		if !rule.IsValidForAmount(req.Amount) {
			continue
		}
		
		// Calculate tax for this rule
		ruleTaxAmount := rule.CalculateTax(taxableAmount)
		
		// Handle compound tax (tax on tax)
		if rule.IsCompound && totalTaxAmount > 0 {
			ruleTaxAmount += rule.CalculateTax(totalTaxAmount)
		}
		
		ruleTaxAmount = RoundToTwoDecimals(ruleTaxAmount)
		totalTaxAmount += ruleTaxAmount
		
		appliedRules = append(appliedRules, AppliedTaxRuleResponse{
			RuleID:        rule.ID,
			RuleName:      rule.Name,
			RuleCode:      rule.Code,
			AppliedRate:   rule.Rate,
			TaxableAmount: taxableAmount,
			TaxAmount:     ruleTaxAmount,
			Priority:      rule.Priority,
		})
	}
	
	totalTaxAmount = RoundToTwoDecimals(totalTaxAmount)
	totalAmount := taxableAmount + totalTaxAmount
	if method == MethodInclusive {
		totalAmount = req.Amount
	}
	
	effectiveRate := 0.0
	if taxableAmount > 0 {
		effectiveRate = (totalTaxAmount / taxableAmount) * 100
	}
	effectiveRate = RoundToTwoDecimals(effectiveRate)
	
	location := fmt.Sprintf("%s", req.Country)
	if req.State != "" {
		location += ", " + req.State
	}
	if req.City != "" {
		location += ", " + req.City
	}
	if req.PostalCode != "" {
		location += " " + req.PostalCode
	}
	
	return &TaxCalculationResponse{
		TaxableAmount: RoundToTwoDecimals(taxableAmount),
		TaxAmount:     totalTaxAmount,
		TotalAmount:   RoundToTwoDecimals(totalAmount),
		EffectiveRate: effectiveRate,
		Method:        method,
		Location:      location,
		AppliedRules:  appliedRules,
		CalculatedAt:  time.Now(),
	}, nil
}

// buildTaxRuleResponse builds a tax rule response
func (s *ServiceImpl) buildTaxRuleResponse(rule *TaxRule) *TaxRuleResponse {
	return &TaxRuleResponse{
		ID:          rule.ID,
		Name:        rule.Name,
		Description: rule.Description,
		Code:        rule.Code,
		Type:        rule.Type,
		Status:      rule.Status,
		TaxType:     rule.TaxType,
		Rate:        rule.Rate,
		Method:      rule.Method,
		Priority:    rule.Priority,
		IsCompound:  rule.IsCompound,
		IsInclusive: rule.IsInclusive,
		ValidFrom:   rule.ValidFrom,
		ValidTo:     rule.ValidTo,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}
}

// buildTaxRateResponse builds a tax rate response
func (s *ServiceImpl) buildTaxRateResponse(rate *TaxRate) *TaxRateResponse {
	return &TaxRateResponse{
		ID:          rate.ID,
		RuleID:      rate.RuleID,
		Name:        rate.Name,
		Rate:        rate.Rate,
		TaxType:     rate.TaxType,
		Description: rate.Description,
		Country:     rate.Country,
		State:       rate.State,
		City:        rate.City,
		PostalCode:  rate.PostalCode,
		ValidFrom:   rate.ValidFrom,
		ValidTo:     rate.ValidTo,
		IsActive:    rate.IsActive,
		CreatedAt:   rate.CreatedAt,
		UpdatedAt:   rate.UpdatedAt,
	}
}

// buildTaxCalculationResponse builds a tax calculation response
func (s *ServiceImpl) buildTaxCalculationResponse(tax *Tax) *TaxCalculationResponse {
	appliedRules := make([]AppliedTaxRuleResponse, len(tax.AppliedRules))
	for i, app := range tax.AppliedRules {
		appliedRules[i] = AppliedTaxRuleResponse{
			RuleID:        app.RuleID,
			RuleName:      app.RuleName,
			RuleCode:      app.RuleCode,
			AppliedRate:   app.AppliedRate,
			TaxableAmount: app.TaxableAmount,
			TaxAmount:     app.TaxAmount,
			Priority:      app.Priority,
		}
	}
	
	return &TaxCalculationResponse{
		TaxableAmount: tax.TaxableAmount,
		TaxAmount:     tax.TaxAmount,
		TotalAmount:   tax.GetTotalAmount(),
		EffectiveRate: tax.GetEffectiveRate() * 100, // Convert to percentage
		Method:        tax.Method,
		Location:      tax.GetLocationString(),
		AppliedRules:  appliedRules,
		CalculatedAt:  tax.CalculatedAt,
	}
}