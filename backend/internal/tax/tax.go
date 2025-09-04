package tax

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tax calculation types
const (
	TaxTypePercentage = "percentage"
	TaxTypeFixed      = "fixed"
	TaxTypeCompound   = "compound"
)

// Tax rule types
const (
	RuleTypeProduct  = "product"
	RuleTypeCategory = "category"
	RuleTypeLocation = "location"
	RuleTypeCustomer = "customer"
	RuleTypeGlobal   = "global"
)

// Tax statuses
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusArchived = "archived"
)

// Tax calculation methods
const (
	MethodInclusive = "inclusive" // Tax included in price
	MethodExclusive = "exclusive" // Tax added to price
)

// Pagination constants
const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Tax represents a tax calculation result
type Tax struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	OrderID     *uuid.UUID `json:"order_id,omitempty" gorm:"type:uuid;index"`
	ProductID   *uuid.UUID `json:"product_id,omitempty" gorm:"type:uuid;index"`
	CustomerID  *uuid.UUID `json:"customer_id,omitempty" gorm:"type:uuid;index"`
	
	// Tax calculation details
	TaxableAmount float64 `json:"taxable_amount" gorm:"type:decimal(10,2);not null"`
	TaxAmount     float64 `json:"tax_amount" gorm:"type:decimal(10,2);not null"`
	TaxRate       float64 `json:"tax_rate" gorm:"type:decimal(5,4);not null"`
	TaxType       string  `json:"tax_type" gorm:"type:varchar(20);not null"`
	Method        string  `json:"method" gorm:"type:varchar(20);not null;default:'exclusive'"`
	
	// Location details
	Country     string `json:"country" gorm:"type:varchar(2);not null"`
	State       string `json:"state" gorm:"type:varchar(100)"`
	City        string `json:"city" gorm:"type:varchar(100)"`
	PostalCode  string `json:"postal_code" gorm:"type:varchar(20)"`
	
	// Applied rules
	AppliedRules []TaxRuleApplication `json:"applied_rules" gorm:"foreignKey:TaxID"`
	
	// Metadata
	CalculatedAt time.Time `json:"calculated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TaxRule represents a tax rule configuration
type TaxRule struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	
	// Rule identification
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text"`
	Code        string `json:"code" gorm:"type:varchar(50);not null;uniqueIndex:idx_tenant_code"`
	Type        string `json:"type" gorm:"type:varchar(20);not null"`
	Status      string `json:"status" gorm:"type:varchar(20);not null;default:'active'"`
	
	// Tax calculation
	TaxType string  `json:"tax_type" gorm:"type:varchar(20);not null"`
	Rate    float64 `json:"rate" gorm:"type:decimal(5,4);not null"`
	Method  string  `json:"method" gorm:"type:varchar(20);not null;default:'exclusive'"`
	
	// Priority and conditions
	Priority    int    `json:"priority" gorm:"not null;default:0"`
	Conditions  string `json:"conditions" gorm:"type:jsonb"`
	IsCompound  bool   `json:"is_compound" gorm:"not null;default:false"`
	IsInclusive bool   `json:"is_inclusive" gorm:"not null;default:false"`
	
	// Validity period
	ValidFrom *time.Time `json:"valid_from,omitempty"`
	ValidTo   *time.Time `json:"valid_to,omitempty"`
	
	// Location targeting
	Countries   []string `json:"countries" gorm:"type:jsonb"`
	States      []string `json:"states" gorm:"type:jsonb"`
	Cities      []string `json:"cities" gorm:"type:jsonb"`
	PostalCodes []string `json:"postal_codes" gorm:"type:jsonb"`
	
	// Product/Category targeting
	ProductIDs  []uuid.UUID `json:"product_ids" gorm:"type:jsonb"`
	CategoryIDs []uuid.UUID `json:"category_ids" gorm:"type:jsonb"`
	
	// Customer targeting
	CustomerIDs    []uuid.UUID `json:"customer_ids" gorm:"type:jsonb"`
	CustomerGroups []string    `json:"customer_groups" gorm:"type:jsonb"`
	
	// Thresholds
	MinAmount *float64 `json:"min_amount,omitempty" gorm:"type:decimal(10,2)"`
	MaxAmount *float64 `json:"max_amount,omitempty" gorm:"type:decimal(10,2)"`
	
	// Relations
	Rates        []TaxRate             `json:"rates" gorm:"foreignKey:RuleID"`
	Applications []TaxRuleApplication  `json:"applications" gorm:"foreignKey:RuleID"`
	
	// Timestamps
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TaxRate represents specific tax rates for different scenarios
type TaxRate struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	RuleID   uuid.UUID `json:"rule_id" gorm:"type:uuid;not null;index"`
	
	// Rate details
	Name        string  `json:"name" gorm:"type:varchar(255);not null"`
	Rate        float64 `json:"rate" gorm:"type:decimal(5,4);not null"`
	TaxType     string  `json:"tax_type" gorm:"type:varchar(20);not null"`
	Description string  `json:"description" gorm:"type:text"`
	
	// Location specificity
	Country    string `json:"country" gorm:"type:varchar(2);not null"`
	State      string `json:"state" gorm:"type:varchar(100)"`
	City       string `json:"city" gorm:"type:varchar(100)"`
	PostalCode string `json:"postal_code" gorm:"type:varchar(20)"`
	
	// Validity
	ValidFrom *time.Time `json:"valid_from,omitempty"`
	ValidTo   *time.Time `json:"valid_to,omitempty"`
	IsActive  bool       `json:"is_active" gorm:"not null;default:true"`
	
	// Relations
	Rule *TaxRule `json:"rule,omitempty" gorm:"foreignKey:RuleID"`
	
	// Timestamps
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TaxRuleApplication tracks which rules were applied to a tax calculation
type TaxRuleApplication struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	TaxID    uuid.UUID `json:"tax_id" gorm:"type:uuid;not null;index"`
	RuleID   uuid.UUID `json:"rule_id" gorm:"type:uuid;not null;index"`
	
	// Application details
	RuleName      string  `json:"rule_name" gorm:"type:varchar(255);not null"`
	RuleCode      string  `json:"rule_code" gorm:"type:varchar(50);not null"`
	AppliedRate   float64 `json:"applied_rate" gorm:"type:decimal(5,4);not null"`
	TaxableAmount float64 `json:"taxable_amount" gorm:"type:decimal(10,2);not null"`
	TaxAmount     float64 `json:"tax_amount" gorm:"type:decimal(10,2);not null"`
	Priority      int     `json:"priority" gorm:"not null"`
	
	// Relations
	Tax  *Tax     `json:"tax,omitempty" gorm:"foreignKey:TaxID"`
	Rule *TaxRule `json:"rule,omitempty" gorm:"foreignKey:RuleID"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// Business logic methods for Tax

// GetEffectiveRate returns the effective tax rate
func (t *Tax) GetEffectiveRate() float64 {
	if t.TaxableAmount == 0 {
		return 0
	}
	return t.TaxAmount / t.TaxableAmount
}

// GetTotalAmount returns the total amount including tax
func (t *Tax) GetTotalAmount() float64 {
	if t.Method == MethodInclusive {
		return t.TaxableAmount
	}
	return t.TaxableAmount + t.TaxAmount
}

// GetBaseAmount returns the base amount excluding tax
func (t *Tax) GetBaseAmount() float64 {
	if t.Method == MethodInclusive {
		return t.TaxableAmount - t.TaxAmount
	}
	return t.TaxableAmount
}

// IsValid validates the tax calculation
func (t *Tax) IsValid() bool {
	return t.TaxableAmount >= 0 && t.TaxAmount >= 0 && t.TaxRate >= 0
}

// GetLocationString returns formatted location string
func (t *Tax) GetLocationString() string {
	parts := []string{}
	if t.City != "" {
		parts = append(parts, t.City)
	}
	if t.State != "" {
		parts = append(parts, t.State)
	}
	if t.Country != "" {
		parts = append(parts, t.Country)
	}
	if t.PostalCode != "" {
		parts = append(parts, t.PostalCode)
	}
	return strings.Join(parts, ", ")
}

// Business logic methods for TaxRule

// IsActive checks if the rule is currently active
func (tr *TaxRule) IsActive() bool {
	return tr.Status == StatusActive
}

// IsValidForDate checks if the rule is valid for a specific date
func (tr *TaxRule) IsValidForDate(date time.Time) bool {
	if tr.ValidFrom != nil && date.Before(*tr.ValidFrom) {
		return false
	}
	if tr.ValidTo != nil && date.After(*tr.ValidTo) {
		return false
	}
	return true
}

// IsValidForLocation checks if the rule applies to a location
func (tr *TaxRule) IsValidForLocation(country, state, city, postalCode string) bool {
	// Check countries
	if len(tr.Countries) > 0 {
		found := false
		for _, c := range tr.Countries {
			if strings.EqualFold(c, country) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check states
	if len(tr.States) > 0 {
		found := false
		for _, s := range tr.States {
			if strings.EqualFold(s, state) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check cities
	if len(tr.Cities) > 0 {
		found := false
		for _, c := range tr.Cities {
			if strings.EqualFold(c, city) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check postal codes
	if len(tr.PostalCodes) > 0 {
		found := false
		for _, pc := range tr.PostalCodes {
			if strings.EqualFold(pc, postalCode) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	return true
}

// IsValidForProduct checks if the rule applies to a product
func (tr *TaxRule) IsValidForProduct(productID uuid.UUID, categoryIDs []uuid.UUID) bool {
	// Check specific products
	if len(tr.ProductIDs) > 0 {
		for _, pid := range tr.ProductIDs {
			if pid == productID {
				return true
			}
		}
		return false
	}
	
	// Check categories
	if len(tr.CategoryIDs) > 0 {
		for _, cid := range tr.CategoryIDs {
			for _, categoryID := range categoryIDs {
				if cid == categoryID {
					return true
				}
			}
		}
		return false
	}
	
	// If no specific targeting, applies to all
	return true
}

// IsValidForCustomer checks if the rule applies to a customer
func (tr *TaxRule) IsValidForCustomer(customerID uuid.UUID, customerGroups []string) bool {
	// Check specific customers
	if len(tr.CustomerIDs) > 0 {
		for _, cid := range tr.CustomerIDs {
			if cid == customerID {
				return true
			}
		}
		return false
	}
	
	// Check customer groups
	if len(tr.CustomerGroups) > 0 {
		for _, group := range tr.CustomerGroups {
			for _, customerGroup := range customerGroups {
				if strings.EqualFold(group, customerGroup) {
					return true
				}
			}
		}
		return false
	}
	
	// If no specific targeting, applies to all
	return true
}

// IsValidForAmount checks if the rule applies to an amount
func (tr *TaxRule) IsValidForAmount(amount float64) bool {
	if tr.MinAmount != nil && amount < *tr.MinAmount {
		return false
	}
	if tr.MaxAmount != nil && amount > *tr.MaxAmount {
		return false
	}
	return true
}

// CalculateTax calculates tax for a given amount
func (tr *TaxRule) CalculateTax(amount float64) float64 {
	switch tr.TaxType {
	case TaxTypePercentage:
		return amount * tr.Rate / 100
	case TaxTypeFixed:
		return tr.Rate
	default:
		return 0
	}
}

// GetDisplayName returns a display-friendly name
func (tr *TaxRule) GetDisplayName() string {
	if tr.Name != "" {
		return tr.Name
	}
	return tr.Code
}

// CanDelete checks if the rule can be deleted
func (tr *TaxRule) CanDelete() bool {
	return tr.Status != StatusActive
}

// Business logic methods for TaxRate

// IsValidForDate checks if the rate is valid for a specific date
func (tr *TaxRate) IsValidForDate(date time.Time) bool {
	if !tr.IsActive {
		return false
	}
	if tr.ValidFrom != nil && date.Before(*tr.ValidFrom) {
		return false
	}
	if tr.ValidTo != nil && date.After(*tr.ValidTo) {
		return false
	}
	return true
}

// CalculateTax calculates tax for a given amount
func (tr *TaxRate) CalculateTax(amount float64) float64 {
	switch tr.TaxType {
	case TaxTypePercentage:
		return amount * tr.Rate / 100
	case TaxTypeFixed:
		return tr.Rate
	default:
		return 0
	}
}

// GetLocationString returns formatted location string
func (tr *TaxRate) GetLocationString() string {
	parts := []string{}
	if tr.City != "" {
		parts = append(parts, tr.City)
	}
	if tr.State != "" {
		parts = append(parts, tr.State)
	}
	if tr.Country != "" {
		parts = append(parts, tr.Country)
	}
	if tr.PostalCode != "" {
		parts = append(parts, tr.PostalCode)
	}
	return strings.Join(parts, ", ")
}

// GORM hooks

// BeforeCreate sets up the tax rule before creation
func (tr *TaxRule) BeforeCreate(tx *gorm.DB) error {
	if tr.ID == uuid.Nil {
		tr.ID = uuid.New()
	}
	
	// Validate required fields
	if tr.Name == "" {
		return errors.New("name is required")
	}
	if tr.Code == "" {
		return errors.New("code is required")
	}
	if tr.Rate < 0 {
		return errors.New("rate cannot be negative")
	}
	
	// Set defaults
	if tr.Status == "" {
		tr.Status = StatusActive
	}
	if tr.Method == "" {
		tr.Method = MethodExclusive
	}
	if tr.TaxType == "" {
		tr.TaxType = TaxTypePercentage
	}
	
	return nil
}

// BeforeUpdate validates the tax rule before update
func (tr *TaxRule) BeforeUpdate(tx *gorm.DB) error {
	if tr.Rate < 0 {
		return errors.New("rate cannot be negative")
	}
	return nil
}

// BeforeCreate sets up the tax rate before creation
func (tr *TaxRate) BeforeCreate(tx *gorm.DB) error {
	if tr.ID == uuid.Nil {
		tr.ID = uuid.New()
	}
	
	// Validate required fields
	if tr.Name == "" {
		return errors.New("name is required")
	}
	if tr.Rate < 0 {
		return errors.New("rate cannot be negative")
	}
	if tr.Country == "" {
		return errors.New("country is required")
	}
	
	// Set defaults
	if tr.TaxType == "" {
		tr.TaxType = TaxTypePercentage
	}
	
	return nil
}

// BeforeUpdate validates the tax rate before update
func (tr *TaxRate) BeforeUpdate(tx *gorm.DB) error {
	if tr.Rate < 0 {
		return errors.New("rate cannot be negative")
	}
	return nil
}

// Request/Response structures

// CreateTaxRuleRequest represents a request to create a tax rule
type CreateTaxRuleRequest struct {
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description"`
	Code        string    `json:"code" validate:"required,min=1,max=50"`
	Type        string    `json:"type" validate:"required,oneof=product category location customer global"`
	TaxType     string    `json:"tax_type" validate:"required,oneof=percentage fixed compound"`
	Rate        float64   `json:"rate" validate:"required,min=0"`
	Method      string    `json:"method" validate:"oneof=inclusive exclusive"`
	Priority    int       `json:"priority"`
	Conditions  string    `json:"conditions"`
	IsCompound  bool      `json:"is_compound"`
	IsInclusive bool      `json:"is_inclusive"`
	ValidFrom   *time.Time `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
	
	// Location targeting
	Countries   []string `json:"countries"`
	States      []string `json:"states"`
	Cities      []string `json:"cities"`
	PostalCodes []string `json:"postal_codes"`
	
	// Product/Category targeting
	ProductIDs  []uuid.UUID `json:"product_ids"`
	CategoryIDs []uuid.UUID `json:"category_ids"`
	
	// Customer targeting
	CustomerIDs    []uuid.UUID `json:"customer_ids"`
	CustomerGroups []string    `json:"customer_groups"`
	
	// Thresholds
	MinAmount *float64 `json:"min_amount"`
	MaxAmount *float64 `json:"max_amount"`
}

// UpdateTaxRuleRequest represents a request to update a tax rule
type UpdateTaxRuleRequest struct {
	Name        *string    `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string    `json:"description"`
	Status      *string    `json:"status" validate:"omitempty,oneof=active inactive archived"`
	TaxType     *string    `json:"tax_type" validate:"omitempty,oneof=percentage fixed compound"`
	Rate        *float64   `json:"rate" validate:"omitempty,min=0"`
	Method      *string    `json:"method" validate:"omitempty,oneof=inclusive exclusive"`
	Priority    *int       `json:"priority"`
	Conditions  *string    `json:"conditions"`
	IsCompound  *bool      `json:"is_compound"`
	IsInclusive *bool      `json:"is_inclusive"`
	ValidFrom   *time.Time `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
	
	// Location targeting
	Countries   []string `json:"countries"`
	States      []string `json:"states"`
	Cities      []string `json:"cities"`
	PostalCodes []string `json:"postal_codes"`
	
	// Product/Category targeting
	ProductIDs  []uuid.UUID `json:"product_ids"`
	CategoryIDs []uuid.UUID `json:"category_ids"`
	
	// Customer targeting
	CustomerIDs    []uuid.UUID `json:"customer_ids"`
	CustomerGroups []string    `json:"customer_groups"`
	
	// Thresholds
	MinAmount *float64 `json:"min_amount"`
	MaxAmount *float64 `json:"max_amount"`
}

// CreateTaxRateRequest represents a request to create a tax rate
type CreateTaxRateRequest struct {
	RuleID      uuid.UUID  `json:"rule_id" validate:"required"`
	Name        string     `json:"name" validate:"required,min=1,max=255"`
	Rate        float64    `json:"rate" validate:"required,min=0"`
	TaxType     string     `json:"tax_type" validate:"required,oneof=percentage fixed"`
	Description string     `json:"description"`
	Country     string     `json:"country" validate:"required,len=2"`
	State       string     `json:"state"`
	City        string     `json:"city"`
	PostalCode  string     `json:"postal_code"`
	ValidFrom   *time.Time `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
}

// UpdateTaxRateRequest represents a request to update a tax rate
type UpdateTaxRateRequest struct {
	Name        *string    `json:"name" validate:"omitempty,min=1,max=255"`
	Rate        *float64   `json:"rate" validate:"omitempty,min=0"`
	TaxType     *string    `json:"tax_type" validate:"omitempty,oneof=percentage fixed"`
	Description *string    `json:"description"`
	State       *string    `json:"state"`
	City        *string    `json:"city"`
	PostalCode  *string    `json:"postal_code"`
	ValidFrom   *time.Time `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
	IsActive    *bool      `json:"is_active"`
}

// TaxCalculationRequest represents a request to calculate tax
type TaxCalculationRequest struct {
	Amount     float64     `json:"amount" validate:"required,min=0"`
	ProductID  *uuid.UUID  `json:"product_id"`
	CustomerID *uuid.UUID  `json:"customer_id"`
	Country    string      `json:"country" validate:"required,len=2"`
	State      string      `json:"state"`
	City       string      `json:"city"`
	PostalCode string      `json:"postal_code"`
	Method     string      `json:"method" validate:"oneof=inclusive exclusive"`
	Date       *time.Time  `json:"date"`
}

// TaxRuleResponse represents a tax rule response
type TaxRuleResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Code        string     `json:"code"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	TaxType     string     `json:"tax_type"`
	Rate        float64    `json:"rate"`
	Method      string     `json:"method"`
	Priority    int        `json:"priority"`
	IsCompound  bool       `json:"is_compound"`
	IsInclusive bool       `json:"is_inclusive"`
	ValidFrom   *time.Time `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TaxRateResponse represents a tax rate response
type TaxRateResponse struct {
	ID          uuid.UUID  `json:"id"`
	RuleID      uuid.UUID  `json:"rule_id"`
	Name        string     `json:"name"`
	Rate        float64    `json:"rate"`
	TaxType     string     `json:"tax_type"`
	Description string     `json:"description"`
	Country     string     `json:"country"`
	State       string     `json:"state"`
	City        string     `json:"city"`
	PostalCode  string     `json:"postal_code"`
	ValidFrom   *time.Time `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TaxCalculationResponse represents a tax calculation response
type TaxCalculationResponse struct {
	TaxableAmount float64                 `json:"taxable_amount"`
	TaxAmount     float64                 `json:"tax_amount"`
	TotalAmount   float64                 `json:"total_amount"`
	EffectiveRate float64                 `json:"effective_rate"`
	Method        string                  `json:"method"`
	Location      string                  `json:"location"`
	AppliedRules  []AppliedTaxRuleResponse `json:"applied_rules"`
	CalculatedAt  time.Time               `json:"calculated_at"`
}

// AppliedTaxRuleResponse represents an applied tax rule in calculation
type AppliedTaxRuleResponse struct {
	RuleID        uuid.UUID `json:"rule_id"`
	RuleName      string    `json:"rule_name"`
	RuleCode      string    `json:"rule_code"`
	AppliedRate   float64   `json:"applied_rate"`
	TaxableAmount float64   `json:"taxable_amount"`
	TaxAmount     float64   `json:"tax_amount"`
	Priority      int       `json:"priority"`
}

// Filter structures

// TaxRuleFilter represents filters for tax rules
type TaxRuleFilter struct {
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	TaxType     string     `json:"tax_type"`
	Country     string     `json:"country"`
	State       string     `json:"state"`
	City        string     `json:"city"`
	ProductID   *uuid.UUID `json:"product_id"`
	CategoryID  *uuid.UUID `json:"category_id"`
	CustomerID  *uuid.UUID `json:"customer_id"`
	ValidDate   *time.Time `json:"valid_date"`
	Search      string     `json:"search"`
}

// TaxRateFilter represents filters for tax rates
type TaxRateFilter struct {
	RuleID     *uuid.UUID `json:"rule_id"`
	TaxType    string     `json:"tax_type"`
	Country    string     `json:"country"`
	State      string     `json:"state"`
	City       string     `json:"city"`
	PostalCode string     `json:"postal_code"`
	IsActive   *bool      `json:"is_active"`
	ValidDate  *time.Time `json:"valid_date"`
	Search     string     `json:"search"`
}

// TaxFilter represents filters for tax calculations
type TaxFilter struct {
	OrderID    *uuid.UUID `json:"order_id"`
	ProductID  *uuid.UUID `json:"product_id"`
	CustomerID *uuid.UUID `json:"customer_id"`
	Country    string     `json:"country"`
	State      string     `json:"state"`
	City       string     `json:"city"`
	TaxType    string     `json:"tax_type"`
	Method     string     `json:"method"`
	DateFrom   *time.Time `json:"date_from"`
	DateTo     *time.Time `json:"date_to"`
}

// Statistics structures

// TaxStats represents tax statistics
type TaxStats struct {
	TotalRules       int64   `json:"total_rules"`
	ActiveRules      int64   `json:"active_rules"`
	InactiveRules    int64   `json:"inactive_rules"`
	TotalRates       int64   `json:"total_rates"`
	ActiveRates      int64   `json:"active_rates"`
	TotalCalculations int64  `json:"total_calculations"`
	TotalTaxAmount   float64 `json:"total_tax_amount"`
	AverageTaxRate   float64 `json:"average_tax_rate"`
}

// TaxByLocation represents tax statistics by location
type TaxByLocation struct {
	Country       string  `json:"country"`
	State         string  `json:"state"`
	City          string  `json:"city"`
	Calculations  int64   `json:"calculations"`
	TotalTax      float64 `json:"total_tax"`
	AverageRate   float64 `json:"average_rate"`
}

// TaxByType represents tax statistics by type
type TaxByType struct {
	TaxType       string  `json:"tax_type"`
	Calculations  int64   `json:"calculations"`
	TotalTax      float64 `json:"total_tax"`
	AverageRate   float64 `json:"average_rate"`
}

// Business logic errors
var (
	ErrTaxRuleNotFound     = errors.New("tax rule not found")
	ErrTaxRateNotFound     = errors.New("tax rate not found")
	ErrTaxNotFound         = errors.New("tax calculation not found")
	ErrInvalidTaxType      = errors.New("invalid tax type")
	ErrInvalidMethod       = errors.New("invalid calculation method")
	ErrInvalidRate         = errors.New("invalid tax rate")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInvalidLocation     = errors.New("invalid location")
	ErrRuleCodeExists      = errors.New("tax rule code already exists")
	ErrInvalidDateRange    = errors.New("invalid date range")
	ErrCannotDeleteRule    = errors.New("cannot delete active tax rule")
	ErrNoApplicableRules   = errors.New("no applicable tax rules found")
	ErrCircularDependency  = errors.New("circular dependency in compound tax rules")
	ErrBulkSizeExceeded    = errors.New("bulk operation size exceeded")
)

// Helper functions

// RoundToTwoDecimals rounds a float64 to 2 decimal places
func RoundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

// FormatRate formats a tax rate for display
func FormatRate(rate float64, taxType string) string {
	switch taxType {
	case TaxTypePercentage:
		return fmt.Sprintf("%.2f%%", rate)
	case TaxTypeFixed:
		return fmt.Sprintf("$%.2f", rate)
	default:
		return fmt.Sprintf("%.2f", rate)
	}
}

// ValidateTaxType validates a tax type
func ValidateTaxType(taxType string) bool {
	return taxType == TaxTypePercentage || taxType == TaxTypeFixed || taxType == TaxTypeCompound
}

// ValidateMethod validates a calculation method
func ValidateMethod(method string) bool {
	return method == MethodInclusive || method == MethodExclusive
}

// ValidateRuleType validates a rule type
func ValidateRuleType(ruleType string) bool {
	return ruleType == RuleTypeProduct || ruleType == RuleTypeCategory || 
		   ruleType == RuleTypeLocation || ruleType == RuleTypeCustomer || 
		   ruleType == RuleTypeGlobal
}

// ValidateStatus validates a status
func ValidateStatus(status string) bool {
	return status == StatusActive || status == StatusInactive || status == StatusArchived
}