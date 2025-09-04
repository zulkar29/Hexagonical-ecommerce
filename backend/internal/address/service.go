package address

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for address business logic
type Service interface {
	// Address CRUD operations
	CreateAddress(ctx context.Context, tenantID uuid.UUID, req CreateAddressRequest) (*AddressResponse, error)
	GetAddress(ctx context.Context, tenantID, addressID uuid.UUID) (*AddressResponse, error)
	UpdateAddress(ctx context.Context, tenantID, addressID uuid.UUID, req UpdateAddressRequest) (*AddressResponse, error)
	DeleteAddress(ctx context.Context, tenantID, addressID uuid.UUID) error
	ListAddresses(ctx context.Context, tenantID uuid.UUID, filter AddressFilter, limit, offset int) (*AddressListResponse, error)
	
	// Customer address operations
	GetCustomerAddresses(ctx context.Context, tenantID, customerID uuid.UUID) ([]*AddressResponse, error)
	GetDefaultAddress(ctx context.Context, tenantID, customerID uuid.UUID, addressType string) (*AddressResponse, error)
	SetDefaultAddress(ctx context.Context, tenantID, customerID, addressID uuid.UUID) error
	UnsetDefaultAddresses(ctx context.Context, tenantID, customerID uuid.UUID, addressType string) error
	
	// Address validation operations
	ValidateAddress(ctx context.Context, tenantID, addressID uuid.UUID, req ValidateAddressRequest) (*AddressValidationResponse, error)
	GetAddressValidation(ctx context.Context, tenantID, addressID uuid.UUID) (*AddressValidationResponse, error)
	ListAddressValidations(ctx context.Context, tenantID uuid.UUID, limit, offset int) (*AddressValidationListResponse, error)
	
	// Bulk operations
	BulkCreateAddresses(ctx context.Context, tenantID uuid.UUID, requests []CreateAddressRequest) ([]*AddressResponse, error)
	BulkUpdateAddresses(ctx context.Context, tenantID uuid.UUID, updates []BulkUpdateAddressRequest) ([]*AddressResponse, error)
	BulkDeleteAddresses(ctx context.Context, tenantID uuid.UUID, addressIDs []uuid.UUID) error
	
	// Statistics and analytics
	GetAddressStats(ctx context.Context, tenantID uuid.UUID) (*AddressStats, error)
	GetAddressesByCountry(ctx context.Context, tenantID uuid.UUID) (map[string]int64, error)
	GetAddressesByType(ctx context.Context, tenantID uuid.UUID) (map[string]int64, error)
	GetRecentAddresses(ctx context.Context, tenantID uuid.UUID, days int) ([]*AddressResponse, error)
	
	// Maintenance operations
	CleanupUnvalidatedAddresses(ctx context.Context, tenantID uuid.UUID, days int) (int64, error)
	CleanupOrphanedValidations(ctx context.Context, tenantID uuid.UUID) (int64, error)
	
	// Utility operations
	NormalizeAddress(ctx context.Context, req NormalizeAddressRequest) (*NormalizeAddressResponse, error)
	SuggestAddresses(ctx context.Context, req AddressSuggestionRequest) (*AddressSuggestionResponse, error)
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new address service
func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

// Address CRUD operations

// CreateAddress creates a new address
func (s *ServiceImpl) CreateAddress(ctx context.Context, tenantID uuid.UUID, req CreateAddressRequest) (*AddressResponse, error) {
	// Validate request
	if err := s.validateCreateAddressRequest(req); err != nil {
		return nil, err
	}
	
	// Check if customer exists (this would typically call user service)
	// For now, we'll assume the customer exists
	
	// Check address limit per customer
	count, err := s.repo.CountCustomerAddresses(ctx, tenantID, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to count customer addresses: %w", err)
	}
	
	if count >= MaxAddressesPerCustomer {
		return nil, ErrTooManyAddresses
	}
	
	// Create address entity
	address := &Address{
		ID:         uuid.New(),
		TenantID:   tenantID,
		CustomerID: req.CustomerID,
		Type:       req.Type,
		Label:      req.Label,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Company:    req.Company,
		Address1:   req.Address1,
		Address2:   req.Address2,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    strings.ToUpper(req.Country),
		Phone:      req.Phone,
		IsDefault:  req.IsDefault,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	// If this is set as default, unset other defaults
	if req.IsDefault {
		if err := s.repo.UnsetDefaultAddresses(ctx, tenantID, req.CustomerID, req.Type); err != nil {
			return nil, fmt.Errorf("failed to unset default addresses: %w", err)
		}
	}
	
	// Create the address
	if err := s.repo.Create(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}
	
	return s.buildAddressResponse(address), nil
}

// GetAddress retrieves an address by ID
func (s *ServiceImpl) GetAddress(ctx context.Context, tenantID, addressID uuid.UUID) (*AddressResponse, error) {
	address, err := s.repo.GetByID(ctx, tenantID, addressID)
	if err != nil {
		return nil, err
	}
	
	return s.buildAddressResponse(address), nil
}

// UpdateAddress updates an existing address
func (s *ServiceImpl) UpdateAddress(ctx context.Context, tenantID, addressID uuid.UUID, req UpdateAddressRequest) (*AddressResponse, error) {
	// Get existing address
	address, err := s.repo.GetByID(ctx, tenantID, addressID)
	if err != nil {
		return nil, err
	}
	
	// Validate request
	if err := s.validateUpdateAddressRequest(req); err != nil {
		return nil, err
	}
	
	// Update fields
	if req.Type != nil {
		address.Type = *req.Type
	}
	if req.Label != nil {
		address.Label = *req.Label
	}
	if req.FirstName != nil {
		address.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		address.LastName = *req.LastName
	}
	if req.Company != nil {
		address.Company = req.Company
	}
	if req.Address1 != nil {
		address.Address1 = *req.Address1
	}
	if req.Address2 != nil {
		address.Address2 = req.Address2
	}
	if req.City != nil {
		address.City = *req.City
	}
	if req.State != nil {
		address.State = *req.State
	}
	if req.PostalCode != nil {
		address.PostalCode = *req.PostalCode
	}
	if req.Country != nil {
		address.Country = strings.ToUpper(*req.Country)
	}
	if req.Phone != nil {
		address.Phone = req.Phone
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
		
		// If setting as default, unset other defaults
		if *req.IsDefault {
			if err := s.repo.UnsetDefaultAddresses(ctx, tenantID, address.CustomerID, address.Type); err != nil {
				return nil, fmt.Errorf("failed to unset default addresses: %w", err)
			}
		}
	}
	
	address.UpdatedAt = time.Now()
	
	// Reset validation if address details changed
	if s.addressDetailsChanged(req) {
		address.IsValidated = false
		address.ValidationScore = nil
	}
	
	// Update the address
	if err := s.repo.Update(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}
	
	return s.buildAddressResponse(address), nil
}

// DeleteAddress soft deletes an address
func (s *ServiceImpl) DeleteAddress(ctx context.Context, tenantID, addressID uuid.UUID) error {
	// Check if address exists
	exists, err := s.repo.ExistsByID(ctx, tenantID, addressID)
	if err != nil {
		return fmt.Errorf("failed to check address existence: %w", err)
	}
	
	if !exists {
		return ErrAddressNotFound
	}
	
	// Delete the address
	if err := s.repo.Delete(ctx, tenantID, addressID); err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}
	
	return nil
}

// ListAddresses retrieves addresses with filtering and pagination
func (s *ServiceImpl) ListAddresses(ctx context.Context, tenantID uuid.UUID, filter AddressFilter, limit, offset int) (*AddressListResponse, error) {
	// Validate pagination
	if limit <= 0 || limit > MaxPageSize {
		limit = DefaultPageSize
	}
	if offset < 0 {
		offset = 0
	}
	
	addresses, total, err := s.repo.List(ctx, tenantID, filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}
	
	responses := make([]*AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = s.buildAddressResponse(address)
	}
	
	return &AddressListResponse{
		Addresses: responses,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	}, nil
}

// Customer address operations

// GetCustomerAddresses retrieves all addresses for a customer
func (s *ServiceImpl) GetCustomerAddresses(ctx context.Context, tenantID, customerID uuid.UUID) ([]*AddressResponse, error) {
	addresses, err := s.repo.GetCustomerAddresses(ctx, tenantID, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer addresses: %w", err)
	}
	
	responses := make([]*AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = s.buildAddressResponse(address)
	}
	
	return responses, nil
}

// GetDefaultAddress retrieves the default address for a customer and type
func (s *ServiceImpl) GetDefaultAddress(ctx context.Context, tenantID, customerID uuid.UUID, addressType string) (*AddressResponse, error) {
	address, err := s.repo.GetDefaultAddress(ctx, tenantID, customerID, addressType)
	if err != nil {
		return nil, err
	}
	
	return s.buildAddressResponse(address), nil
}

// SetDefaultAddress sets an address as default for a customer
func (s *ServiceImpl) SetDefaultAddress(ctx context.Context, tenantID, customerID, addressID uuid.UUID) error {
	// Verify the address belongs to the customer
	address, err := s.repo.GetByID(ctx, tenantID, addressID)
	if err != nil {
		return err
	}
	
	if address.CustomerID != customerID {
		return ErrAddressNotFound
	}
	
	return s.repo.SetDefaultAddress(ctx, tenantID, customerID, addressID)
}

// UnsetDefaultAddresses unsets default addresses for a customer and type
func (s *ServiceImpl) UnsetDefaultAddresses(ctx context.Context, tenantID, customerID uuid.UUID, addressType string) error {
	return s.repo.UnsetDefaultAddresses(ctx, tenantID, customerID, addressType)
}

// Address validation operations

// ValidateAddress validates an address using external service
func (s *ServiceImpl) ValidateAddress(ctx context.Context, tenantID, addressID uuid.UUID, req ValidateAddressRequest) (*AddressValidationResponse, error) {
	// Get the address
	address, err := s.repo.GetByID(ctx, tenantID, addressID)
	if err != nil {
		return nil, err
	}
	
	// Create validation record
	validation := &AddressValidation{
		ID:               uuid.New(),
		TenantID:         tenantID,
		AddressID:        addressID,
		Provider:         req.Provider,
		IsValid:          req.IsValid,
		Score:            req.Score,
		NormalizedData:   req.NormalizedData,
		Suggestions:      req.Suggestions,
		ValidationErrors: req.ValidationErrors,
		ValidatedAt:      time.Now(),
	}
	
	// Save validation
	if err := s.repo.CreateValidation(ctx, validation); err != nil {
		return nil, fmt.Errorf("failed to create validation: %w", err)
	}
	
	// Update address validation status
	address.IsValidated = req.IsValid
	address.ValidationScore = &req.Score
	address.UpdatedAt = time.Now()
	
	if err := s.repo.Update(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to update address validation: %w", err)
	}
	
	return s.buildAddressValidationResponse(validation), nil
}

// GetAddressValidation retrieves the latest validation for an address
func (s *ServiceImpl) GetAddressValidation(ctx context.Context, tenantID, addressID uuid.UUID) (*AddressValidationResponse, error) {
	validation, err := s.repo.GetValidation(ctx, tenantID, addressID)
	if err != nil {
		return nil, err
	}
	
	return s.buildAddressValidationResponse(validation), nil
}

// ListAddressValidations retrieves validations with pagination
func (s *ServiceImpl) ListAddressValidations(ctx context.Context, tenantID uuid.UUID, limit, offset int) (*AddressValidationListResponse, error) {
	// Validate pagination
	if limit <= 0 || limit > MaxPageSize {
		limit = DefaultPageSize
	}
	if offset < 0 {
		offset = 0
	}
	
	validations, total, err := s.repo.ListValidations(ctx, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list validations: %w", err)
	}
	
	responses := make([]*AddressValidationResponse, len(validations))
	for i, validation := range validations {
		responses[i] = s.buildAddressValidationResponse(validation)
	}
	
	return &AddressValidationListResponse{
		Validations: responses,
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	}, nil
}

// Bulk operations

// BulkCreateAddresses creates multiple addresses
func (s *ServiceImpl) BulkCreateAddresses(ctx context.Context, tenantID uuid.UUID, requests []CreateAddressRequest) ([]*AddressResponse, error) {
	if len(requests) == 0 {
		return []*AddressResponse{}, nil
	}
	
	if len(requests) > MaxBulkSize {
		return nil, ErrBulkSizeExceeded
	}
	
	addresses := make([]*Address, len(requests))
	for i, req := range requests {
		// Validate each request
		if err := s.validateCreateAddressRequest(req); err != nil {
			return nil, fmt.Errorf("validation failed for request %d: %w", i, err)
		}
		
		// Create address entity
		addresses[i] = &Address{
			ID:         uuid.New(),
			TenantID:   tenantID,
			CustomerID: req.CustomerID,
			Type:       req.Type,
			Label:      req.Label,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Company:    req.Company,
			Address1:   req.Address1,
			Address2:   req.Address2,
			City:       req.City,
			State:      req.State,
			PostalCode: req.PostalCode,
			Country:    strings.ToUpper(req.Country),
			Phone:      req.Phone,
			IsDefault:  req.IsDefault,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}
	
	// Create all addresses
	if err := s.repo.BulkCreate(ctx, addresses); err != nil {
		return nil, fmt.Errorf("failed to bulk create addresses: %w", err)
	}
	
	// Build responses
	responses := make([]*AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = s.buildAddressResponse(address)
	}
	
	return responses, nil
}

// BulkUpdateAddresses updates multiple addresses
func (s *ServiceImpl) BulkUpdateAddresses(ctx context.Context, tenantID uuid.UUID, updates []BulkUpdateAddressRequest) ([]*AddressResponse, error) {
	if len(updates) == 0 {
		return []*AddressResponse{}, nil
	}
	
	if len(updates) > MaxBulkSize {
		return nil, ErrBulkSizeExceeded
	}
	
	addresses := make([]*Address, len(updates))
	for i, update := range updates {
		// Get existing address
		address, err := s.repo.GetByID(ctx, tenantID, update.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get address %s: %w", update.ID, err)
		}
		
		// Apply updates
		if update.Label != nil {
			address.Label = *update.Label
		}
		if update.IsDefault != nil {
			address.IsDefault = *update.IsDefault
		}
		
		address.UpdatedAt = time.Now()
		addresses[i] = address
	}
	
	// Update all addresses
	if err := s.repo.BulkUpdate(ctx, addresses); err != nil {
		return nil, fmt.Errorf("failed to bulk update addresses: %w", err)
	}
	
	// Build responses
	responses := make([]*AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = s.buildAddressResponse(address)
	}
	
	return responses, nil
}

// BulkDeleteAddresses soft deletes multiple addresses
func (s *ServiceImpl) BulkDeleteAddresses(ctx context.Context, tenantID uuid.UUID, addressIDs []uuid.UUID) error {
	if len(addressIDs) == 0 {
		return nil
	}
	
	if len(addressIDs) > MaxBulkSize {
		return ErrBulkSizeExceeded
	}
	
	return s.repo.BulkDelete(ctx, tenantID, addressIDs)
}

// Statistics and analytics

// GetAddressStats retrieves address statistics
func (s *ServiceImpl) GetAddressStats(ctx context.Context, tenantID uuid.UUID) (*AddressStats, error) {
	return s.repo.GetStats(ctx, tenantID)
}

// GetAddressesByCountry retrieves address count by country
func (s *ServiceImpl) GetAddressesByCountry(ctx context.Context, tenantID uuid.UUID) (map[string]int64, error) {
	return s.repo.GetAddressesByCountry(ctx, tenantID)
}

// GetAddressesByType retrieves address count by type
func (s *ServiceImpl) GetAddressesByType(ctx context.Context, tenantID uuid.UUID) (map[string]int64, error) {
	return s.repo.GetAddressesByType(ctx, tenantID)
}

// GetRecentAddresses retrieves recent addresses
func (s *ServiceImpl) GetRecentAddresses(ctx context.Context, tenantID uuid.UUID, days int) ([]*AddressResponse, error) {
	if days <= 0 {
		days = 30 // Default to 30 days
	}
	
	addresses, err := s.repo.GetRecentAddresses(ctx, tenantID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent addresses: %w", err)
	}
	
	responses := make([]*AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = s.buildAddressResponse(address)
	}
	
	return responses, nil
}

// Maintenance operations

// CleanupUnvalidatedAddresses removes old unvalidated addresses
func (s *ServiceImpl) CleanupUnvalidatedAddresses(ctx context.Context, tenantID uuid.UUID, days int) (int64, error) {
	if days <= 0 {
		days = 90 // Default to 90 days
	}
	
	return s.repo.CleanupUnvalidatedAddresses(ctx, tenantID, days)
}

// CleanupOrphanedValidations removes validations for non-existent addresses
func (s *ServiceImpl) CleanupOrphanedValidations(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return s.repo.CleanupOrphanedValidations(ctx, tenantID)
}

// Utility operations

// NormalizeAddress normalizes address data
func (s *ServiceImpl) NormalizeAddress(ctx context.Context, req NormalizeAddressRequest) (*NormalizeAddressResponse, error) {
	// This would typically call an external address normalization service
	// For now, we'll implement basic normalization
	
	normalized := NormalizeAddressResponse{
		Address1:   strings.TrimSpace(strings.Title(strings.ToLower(req.Address1))),
		Address2:   strings.TrimSpace(strings.Title(strings.ToLower(req.Address2))),
		City:       strings.TrimSpace(strings.Title(strings.ToLower(req.City))),
		State:      strings.TrimSpace(strings.ToUpper(req.State)),
		PostalCode: strings.TrimSpace(strings.ToUpper(req.PostalCode)),
		Country:    strings.TrimSpace(strings.ToUpper(req.Country)),
	}
	
	return &normalized, nil
}

// SuggestAddresses provides address suggestions
func (s *ServiceImpl) SuggestAddresses(ctx context.Context, req AddressSuggestionRequest) (*AddressSuggestionResponse, error) {
	// This would typically call an external address suggestion service
	// For now, we'll return an empty response
	
	return &AddressSuggestionResponse{
		Suggestions: []AddressSuggestion{},
	}, nil
}

// Helper methods

// validateCreateAddressRequest validates create address request
func (s *ServiceImpl) validateCreateAddressRequest(req CreateAddressRequest) error {
	if req.CustomerID == uuid.Nil {
		return ErrInvalidCustomerID
	}
	
	if req.Type == "" {
		req.Type = AddressTypeBoth
	}
	
	if req.Type != AddressTypeShipping && req.Type != AddressTypeBilling && req.Type != AddressTypeBoth {
		return ErrInvalidAddressType
	}
	
	if strings.TrimSpace(req.FirstName) == "" {
		return ErrInvalidFirstName
	}
	
	if strings.TrimSpace(req.LastName) == "" {
		return ErrInvalidLastName
	}
	
	if strings.TrimSpace(req.Address1) == "" {
		return ErrInvalidAddress
	}
	
	if strings.TrimSpace(req.City) == "" {
		return ErrInvalidCity
	}
	
	if strings.TrimSpace(req.State) == "" {
		return ErrInvalidState
	}
	
	if strings.TrimSpace(req.PostalCode) == "" {
		return ErrInvalidPostalCode
	}
	
	if strings.TrimSpace(req.Country) == "" {
		return ErrInvalidCountry
	}
	
	return nil
}

// validateUpdateAddressRequest validates update address request
func (s *ServiceImpl) validateUpdateAddressRequest(req UpdateAddressRequest) error {
	if req.Type != nil {
		if *req.Type != AddressTypeShipping && *req.Type != AddressTypeBilling && *req.Type != AddressTypeBoth {
			return ErrInvalidAddressType
		}
	}
	
	if req.FirstName != nil && strings.TrimSpace(*req.FirstName) == "" {
		return ErrInvalidFirstName
	}
	
	if req.LastName != nil && strings.TrimSpace(*req.LastName) == "" {
		return ErrInvalidLastName
	}
	
	if req.Address1 != nil && strings.TrimSpace(*req.Address1) == "" {
		return ErrInvalidAddress
	}
	
	if req.City != nil && strings.TrimSpace(*req.City) == "" {
		return ErrInvalidCity
	}
	
	if req.State != nil && strings.TrimSpace(*req.State) == "" {
		return ErrInvalidState
	}
	
	if req.PostalCode != nil && strings.TrimSpace(*req.PostalCode) == "" {
		return ErrInvalidPostalCode
	}
	
	if req.Country != nil && strings.TrimSpace(*req.Country) == "" {
		return ErrInvalidCountry
	}
	
	return nil
}

// addressDetailsChanged checks if address details changed
func (s *ServiceImpl) addressDetailsChanged(req UpdateAddressRequest) bool {
	return req.Address1 != nil || req.Address2 != nil || req.City != nil ||
		req.State != nil || req.PostalCode != nil || req.Country != nil
}

// buildAddressResponse builds address response
func (s *ServiceImpl) buildAddressResponse(address *Address) *AddressResponse {
	return &AddressResponse{
		ID:              address.ID,
		CustomerID:      address.CustomerID,
		Type:            address.Type,
		Label:           address.Label,
		FirstName:       address.FirstName,
		LastName:        address.LastName,
		Company:         address.Company,
		Address1:        address.Address1,
		Address2:        address.Address2,
		City:            address.City,
		State:           address.State,
		PostalCode:      address.PostalCode,
		Country:         address.Country,
		Phone:           address.Phone,
		IsDefault:       address.IsDefault,
		IsValidated:     address.IsValidated,
		ValidationScore: address.ValidationScore,
		CreatedAt:       address.CreatedAt,
		UpdatedAt:       address.UpdatedAt,
	}
}

// buildAddressValidationResponse builds address validation response
func (s *ServiceImpl) buildAddressValidationResponse(validation *AddressValidation) *AddressValidationResponse {
	return &AddressValidationResponse{
		ID:               validation.ID,
		AddressID:        validation.AddressID,
		Provider:         validation.Provider,
		IsValid:          validation.IsValid,
		Score:            validation.Score,
		NormalizedData:   validation.NormalizedData,
		Suggestions:      validation.Suggestions,
		ValidationErrors: validation.ValidationErrors,
		ValidatedAt:      validation.ValidatedAt,
	}
}