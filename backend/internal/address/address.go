package address

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Address represents a customer address
type Address struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	CustomerID  uuid.UUID  `json:"customer_id" gorm:"type:uuid;not null;index"`
	Type        string     `json:"type" gorm:"type:varchar(20);not null;default:'shipping'"` // shipping, billing, both
	Label       string     `json:"label" gorm:"type:varchar(100)"` // Home, Work, etc.
	FirstName   string     `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName    string     `json:"last_name" gorm:"type:varchar(100);not null"`
	Company     string     `json:"company" gorm:"type:varchar(200)"`
	Address1    string     `json:"address1" gorm:"type:varchar(255);not null"`
	Address2    string     `json:"address2" gorm:"type:varchar(255)"`
	City        string     `json:"city" gorm:"type:varchar(100);not null"`
	State       string     `json:"state" gorm:"type:varchar(100);not null"`
	PostalCode  string     `json:"postal_code" gorm:"type:varchar(20);not null"`
	Country     string     `json:"country" gorm:"type:varchar(2);not null"` // ISO 3166-1 alpha-2
	Phone       string     `json:"phone" gorm:"type:varchar(20)"`
	IsDefault   bool       `json:"is_default" gorm:"default:false;index"`
	IsValidated bool       `json:"is_validated" gorm:"default:false"`
	Latitude    *float64   `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude   *float64   `json:"longitude" gorm:"type:decimal(11,8)"`
	Instructions string    `json:"instructions" gorm:"type:text"`
	Metadata    string     `json:"metadata" gorm:"type:jsonb"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

// AddressValidation represents address validation result
type AddressValidation struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID         uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	AddressID        uuid.UUID `json:"address_id" gorm:"type:uuid;not null;index"`
	Provider         string    `json:"provider" gorm:"type:varchar(50);not null"` // google, usps, etc.
	IsValid          bool      `json:"is_valid" gorm:"not null"`
	ConfidenceScore  float64   `json:"confidence_score" gorm:"type:decimal(3,2)"`
	SuggestedAddress string    `json:"suggested_address" gorm:"type:jsonb"`
	ValidationErrors string    `json:"validation_errors" gorm:"type:jsonb"`
	ValidatedAt      time.Time `json:"validated_at" gorm:"autoCreateTime"`

	// Relationships
	Address Address `json:"address" gorm:"foreignKey:AddressID;constraint:OnDelete:CASCADE"`
}

// Constants for address types
const (
	AddressTypeShipping = "shipping"
	AddressTypeBilling  = "billing"
	AddressTypeBoth     = "both"
)

// Business logic methods for Address

// GetFullName returns the full name
func (a *Address) GetFullName() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", a.FirstName, a.LastName))
}

// GetFormattedAddress returns a formatted address string
func (a *Address) GetFormattedAddress() string {
	var parts []string
	
	if a.Company != "" {
		parts = append(parts, a.Company)
	}
	
	parts = append(parts, a.GetFullName())
	parts = append(parts, a.Address1)
	
	if a.Address2 != "" {
		parts = append(parts, a.Address2)
	}
	
	parts = append(parts, fmt.Sprintf("%s, %s %s", a.City, a.State, a.PostalCode))
	parts = append(parts, a.GetCountryName())
	
	return strings.Join(parts, "\n")
}

// GetCountryName returns the full country name from ISO code
func (a *Address) GetCountryName() string {
	// This would typically use a country code lookup table
	// For now, return the code itself
	return strings.ToUpper(a.Country)
}

// IsShippingAddress checks if address can be used for shipping
func (a *Address) IsShippingAddress() bool {
	return a.Type == AddressTypeShipping || a.Type == AddressTypeBoth
}

// IsBillingAddress checks if address can be used for billing
func (a *Address) IsBillingAddress() bool {
	return a.Type == AddressTypeBilling || a.Type == AddressTypeBoth
}

// IsComplete checks if address has all required fields
func (a *Address) IsComplete() bool {
	return a.FirstName != "" &&
		a.LastName != "" &&
		a.Address1 != "" &&
		a.City != "" &&
		a.State != "" &&
		a.PostalCode != "" &&
		a.Country != ""
}

// HasCoordinates checks if address has latitude and longitude
func (a *Address) HasCoordinates() bool {
	return a.Latitude != nil && a.Longitude != nil
}

// Validate validates the address data
func (a *Address) Validate() error {
	if a.TenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	
	if a.CustomerID == uuid.Nil {
		return ErrInvalidCustomerID
	}
	
	if !a.IsValidAddressType() {
		return ErrInvalidAddressType
	}
	
	if strings.TrimSpace(a.FirstName) == "" {
		return ErrFirstNameRequired
	}
	
	if strings.TrimSpace(a.LastName) == "" {
		return ErrLastNameRequired
	}
	
	if strings.TrimSpace(a.Address1) == "" {
		return ErrAddress1Required
	}
	
	if strings.TrimSpace(a.City) == "" {
		return ErrCityRequired
	}
	
	if strings.TrimSpace(a.State) == "" {
		return ErrStateRequired
	}
	
	if strings.TrimSpace(a.PostalCode) == "" {
		return ErrPostalCodeRequired
	}
	
	if !a.IsValidCountryCode() {
		return ErrInvalidCountryCode
	}
	
	if a.Phone != "" && !a.IsValidPhoneNumber() {
		return ErrInvalidPhoneNumber
	}
	
	if len(a.FirstName) > 100 {
		return ErrFirstNameTooLong
	}
	
	if len(a.LastName) > 100 {
		return ErrLastNameTooLong
	}
	
	if len(a.Company) > 200 {
		return ErrCompanyTooLong
	}
	
	if len(a.Address1) > 255 {
		return ErrAddress1TooLong
	}
	
	if len(a.Address2) > 255 {
		return ErrAddress2TooLong
	}
	
	if len(a.City) > 100 {
		return ErrCityTooLong
	}
	
	if len(a.State) > 100 {
		return ErrStateTooLong
	}
	
	if len(a.PostalCode) > 20 {
		return ErrPostalCodeTooLong
	}
	
	if len(a.Phone) > 20 {
		return ErrPhoneTooLong
	}
	
	if len(a.Label) > 100 {
		return ErrLabelTooLong
	}
	
	if len(a.Instructions) > 1000 {
		return ErrInstructionsTooLong
	}
	
	return nil
}

// IsValidAddressType checks if the address type is valid
func (a *Address) IsValidAddressType() bool {
	return a.Type == AddressTypeShipping || a.Type == AddressTypeBilling || a.Type == AddressTypeBoth
}

// IsValidCountryCode checks if the country code is valid (basic validation)
func (a *Address) IsValidCountryCode() bool {
	if len(a.Country) != 2 {
		return false
	}
	
	// Basic regex for ISO 3166-1 alpha-2 codes
	regex := regexp.MustCompile(`^[A-Z]{2}$`)
	return regex.MatchString(strings.ToUpper(a.Country))
}

// IsValidPhoneNumber checks if the phone number format is valid (basic validation)
func (a *Address) IsValidPhoneNumber() bool {
	if a.Phone == "" {
		return true
	}
	
	// Remove common phone number characters
	phone := strings.ReplaceAll(a.Phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, "+", "")
	
	// Check if remaining characters are digits
	regex := regexp.MustCompile(`^[0-9]{7,15}$`)
	return regex.MatchString(phone)
}

// NormalizeData normalizes address data
func (a *Address) NormalizeData() {
	a.FirstName = strings.TrimSpace(a.FirstName)
	a.LastName = strings.TrimSpace(a.LastName)
	a.Company = strings.TrimSpace(a.Company)
	a.Address1 = strings.TrimSpace(a.Address1)
	a.Address2 = strings.TrimSpace(a.Address2)
	a.City = strings.TrimSpace(a.City)
	a.State = strings.TrimSpace(a.State)
	a.PostalCode = strings.TrimSpace(a.PostalCode)
	a.Country = strings.ToUpper(strings.TrimSpace(a.Country))
	a.Phone = strings.TrimSpace(a.Phone)
	a.Label = strings.TrimSpace(a.Label)
	a.Instructions = strings.TrimSpace(a.Instructions)
	
	// Set default type if empty
	if a.Type == "" {
		a.Type = AddressTypeShipping
	}
}

// GORM hooks

// BeforeCreate is called before creating an address
func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	
	a.NormalizeData()
	return a.Validate()
}

// BeforeUpdate is called before updating an address
func (a *Address) BeforeUpdate(tx *gorm.DB) error {
	a.NormalizeData()
	return a.Validate()
}

// Business logic methods for AddressValidation

// IsHighConfidence checks if validation has high confidence
func (av *AddressValidation) IsHighConfidence() bool {
	return av.ConfidenceScore >= 0.8
}

// NeedsRevalidation checks if address needs revalidation
func (av *AddressValidation) NeedsRevalidation() bool {
	return time.Since(av.ValidatedAt) > 30*24*time.Hour // 30 days
}

// Request/Response structures

// CreateAddressRequest represents a request to create an address
type CreateAddressRequest struct {
	TenantID     uuid.UUID `json:"tenant_id"`
	CustomerID   uuid.UUID `json:"customer_id" validate:"required"`
	Type         string    `json:"type" validate:"required,oneof=shipping billing both"`
	Label        string    `json:"label"`
	FirstName    string    `json:"first_name" validate:"required,max=100"`
	LastName     string    `json:"last_name" validate:"required,max=100"`
	Company      string    `json:"company" validate:"max=200"`
	Address1     string    `json:"address1" validate:"required,max=255"`
	Address2     string    `json:"address2" validate:"max=255"`
	City         string    `json:"city" validate:"required,max=100"`
	State        string    `json:"state" validate:"required,max=100"`
	PostalCode   string    `json:"postal_code" validate:"required,max=20"`
	Country      string    `json:"country" validate:"required,len=2"`
	Phone        string    `json:"phone" validate:"max=20"`
	IsDefault    bool      `json:"is_default"`
	Latitude     *float64  `json:"latitude"`
	Longitude    *float64  `json:"longitude"`
	Instructions string    `json:"instructions" validate:"max=1000"`
}

// UpdateAddressRequest represents a request to update an address
type UpdateAddressRequest struct {
	Type         *string   `json:"type,omitempty" validate:"omitempty,oneof=shipping billing both"`
	Label        *string   `json:"label,omitempty"`
	FirstName    *string   `json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName     *string   `json:"last_name,omitempty" validate:"omitempty,max=100"`
	Company      *string   `json:"company,omitempty" validate:"omitempty,max=200"`
	Address1     *string   `json:"address1,omitempty" validate:"omitempty,max=255"`
	Address2     *string   `json:"address2,omitempty" validate:"omitempty,max=255"`
	City         *string   `json:"city,omitempty" validate:"omitempty,max=100"`
	State        *string   `json:"state,omitempty" validate:"omitempty,max=100"`
	PostalCode   *string   `json:"postal_code,omitempty" validate:"omitempty,max=20"`
	Country      *string   `json:"country,omitempty" validate:"omitempty,len=2"`
	Phone        *string   `json:"phone,omitempty" validate:"omitempty,max=20"`
	IsDefault    *bool     `json:"is_default,omitempty"`
	Latitude     *float64  `json:"latitude,omitempty"`
	Longitude    *float64  `json:"longitude,omitempty"`
	Instructions *string   `json:"instructions,omitempty" validate:"omitempty,max=1000"`
}

// AddressResponse represents an address response
type AddressResponse struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	CustomerID       uuid.UUID  `json:"customer_id"`
	Type             string     `json:"type"`
	Label            string     `json:"label"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	FullName         string     `json:"full_name"`
	Company          string     `json:"company"`
	Address1         string     `json:"address1"`
	Address2         string     `json:"address2"`
	City             string     `json:"city"`
	State            string     `json:"state"`
	PostalCode       string     `json:"postal_code"`
	Country          string     `json:"country"`
	CountryName      string     `json:"country_name"`
	Phone            string     `json:"phone"`
	IsDefault        bool       `json:"is_default"`
	IsValidated      bool       `json:"is_validated"`
	Latitude         *float64   `json:"latitude"`
	Longitude        *float64   `json:"longitude"`
	HasCoordinates   bool       `json:"has_coordinates"`
	Instructions     string     `json:"instructions"`
	FormattedAddress string     `json:"formatted_address"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// AddressFilter represents address filtering options
type AddressFilter struct {
	CustomerID  *uuid.UUID `json:"customer_id,omitempty"`
	Type        string     `json:"type,omitempty"`
	Label       string     `json:"label,omitempty"`
	Country     string     `json:"country,omitempty"`
	State       string     `json:"state,omitempty"`
	City        string     `json:"city,omitempty"`
	IsDefault   *bool      `json:"is_default,omitempty"`
	IsValidated *bool      `json:"is_validated,omitempty"`
	Search      string     `json:"search,omitempty"`
}

// ValidateAddressRequest represents a request to validate an address
type ValidateAddressRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google usps"`
}

// AddressValidationResponse represents address validation response
type AddressValidationResponse struct {
	ID               uuid.UUID `json:"id"`
	AddressID        uuid.UUID `json:"address_id"`
	Provider         string    `json:"provider"`
	IsValid          bool      `json:"is_valid"`
	ConfidenceScore  float64   `json:"confidence_score"`
	SuggestedAddress string    `json:"suggested_address,omitempty"`
	ValidationErrors []string  `json:"validation_errors,omitempty"`
	ValidatedAt      time.Time `json:"validated_at"`
}

// AddressStats represents address statistics
type AddressStats struct {
	TotalAddresses     int64            `json:"total_addresses"`
	ValidatedAddresses int64            `json:"validated_addresses"`
	DefaultAddresses   int64            `json:"default_addresses"`
	AddressesByType    map[string]int64 `json:"addresses_by_type"`
	AddressesByCountry map[string]int64 `json:"addresses_by_country"`
	RecentAddresses    int64            `json:"recent_addresses"`
}

// Business logic errors
var (
	ErrAddressNotFound         = errors.New("address not found")
	ErrValidationNotFound      = errors.New("address validation not found")
	ErrInvalidTenantID         = errors.New("invalid tenant ID")
	ErrInvalidCustomerID       = errors.New("invalid customer ID")
	ErrInvalidAddressID        = errors.New("invalid address ID")
	ErrInvalidAddressType      = errors.New("invalid address type")
	ErrInvalidCountryCode      = errors.New("invalid country code")
	ErrInvalidPhoneNumber      = errors.New("invalid phone number")
	ErrFirstNameRequired       = errors.New("first name is required")
	ErrLastNameRequired        = errors.New("last name is required")
	ErrAddress1Required        = errors.New("address line 1 is required")
	ErrCityRequired            = errors.New("city is required")
	ErrStateRequired           = errors.New("state is required")
	ErrPostalCodeRequired      = errors.New("postal code is required")
	ErrCountryRequired         = errors.New("country is required")
	ErrFirstNameTooLong        = errors.New("first name is too long")
	ErrLastNameTooLong         = errors.New("last name is too long")
	ErrCompanyTooLong          = errors.New("company name is too long")
	ErrAddress1TooLong         = errors.New("address line 1 is too long")
	ErrAddress2TooLong         = errors.New("address line 2 is too long")
	ErrCityTooLong             = errors.New("city name is too long")
	ErrStateTooLong            = errors.New("state name is too long")
	ErrPostalCodeTooLong       = errors.New("postal code is too long")
	ErrPhoneTooLong            = errors.New("phone number is too long")
	ErrLabelTooLong            = errors.New("label is too long")
	ErrInstructionsTooLong     = errors.New("instructions are too long")
	ErrCannotDeleteDefaultAddress = errors.New("cannot delete default address")
	ErrAddressLimitExceeded    = errors.New("address limit exceeded")
	ErrValidationProviderNotSupported = errors.New("validation provider not supported")
	ErrAddressValidationFailed = errors.New("address validation failed")
)