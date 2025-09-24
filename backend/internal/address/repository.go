package address

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sharedErrors "ecommerce-saas/internal/shared/errors"
)

// Repository defines the interface for address data operations
type Repository interface {
	// Address CRUD operations
	Create(ctx context.Context, address *Address) error
	GetByID(ctx context.Context, tenantID, addressID uuid.UUID) (*Address, error)
	Update(ctx context.Context, address *Address) error
	Delete(ctx context.Context, tenantID, addressID uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, filter AddressFilter, limit, offset int) ([]*Address, int64, error)

	// Customer address operations
	GetCustomerAddresses(ctx context.Context, tenantID, customerID uuid.UUID) ([]*Address, error)
	GetDefaultAddress(ctx context.Context, tenantID, customerID uuid.UUID, addressType string) (*Address, error)
	SetDefaultAddress(ctx context.Context, tenantID, customerID, addressID uuid.UUID) error
	UnsetDefaultAddresses(ctx context.Context, tenantID, customerID uuid.UUID, addressType string) error

	// Address validation operations
	CreateValidation(ctx context.Context, validation *AddressValidation) error
	GetValidation(ctx context.Context, tenantID, addressID uuid.UUID) (*AddressValidation, error)
	GetValidationByID(ctx context.Context, tenantID, validationID uuid.UUID) (*AddressValidation, error)
	UpdateValidation(ctx context.Context, validation *AddressValidation) error
	ListValidations(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*AddressValidation, int64, error)

	// Bulk operations
	BulkCreate(ctx context.Context, addresses []*Address) error
	BulkUpdate(ctx context.Context, addresses []*Address) error
	BulkDelete(ctx context.Context, tenantID uuid.UUID, addressIDs []uuid.UUID) error

	// Validation and checks
	ExistsByID(ctx context.Context, tenantID, addressID uuid.UUID) (bool, error)
	CountCustomerAddresses(ctx context.Context, tenantID, customerID uuid.UUID) (int64, error)
	HasDefaultAddress(ctx context.Context, tenantID, customerID uuid.UUID, addressType string) (bool, error)

	// Statistics and analytics
	GetStats(ctx context.Context, tenantID uuid.UUID) (*AddressStats, error)
	GetAddressesByCountry(ctx context.Context, tenantID uuid.UUID) (map[string]int64, error)
	GetAddressesByType(ctx context.Context, tenantID uuid.UUID) (map[string]int64, error)
	GetRecentAddresses(ctx context.Context, tenantID uuid.UUID, days int) ([]*Address, error)

	// Maintenance operations
	CleanupUnvalidatedAddresses(ctx context.Context, tenantID uuid.UUID, days int) (int64, error)
	CleanupOrphanedValidations(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// GormRepository implements Repository using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM repository
func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// NewRepository creates a new repository (alias for NewGormRepository)
func NewRepository(db *gorm.DB) Repository {
	return NewGormRepository(db)
}

// Address CRUD operations

// Create creates a new address
func (r *GormRepository) Create(ctx context.Context, address *Address) error {
	if err := r.db.WithContext(ctx).Create(address).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create address", err)
	}
	return nil
}

// GetByID retrieves an address by ID
func (r *GormRepository) GetByID(ctx context.Context, tenantID, addressID uuid.UUID) (*Address, error) {
	var address Address
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, addressID).
		First(&address).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Address not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get address", err)
	}

	return &address, nil
}

// Update updates an existing address
func (r *GormRepository) Update(ctx context.Context, address *Address) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", address.TenantID, address.ID).
		Updates(address)

	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to update address", result.Error)
	}

	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Address not found")
	}

	return nil
}

// Delete soft deletes an address
func (r *GormRepository) Delete(ctx context.Context, tenantID, addressID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, addressID).
		Delete(&Address{})

	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to delete address", result.Error)
	}

	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Address not found")
	}

	return nil
}

// List retrieves addresses with filtering and pagination
func (r *GormRepository) List(ctx context.Context, tenantID uuid.UUID, filter AddressFilter, limit, offset int) ([]*Address, int64, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	// Apply filters
	query = r.applyAddressFilters(query, filter)

	// Count total records
	var total int64
	if err := query.Model(&Address{}).Count(&total).Error; err != nil {
		return nil, 0, sharedErrors.NewInternalError("Failed to count addresses", err)
	}

	// Get paginated results
	var addresses []*Address
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&addresses).Error

	if err != nil {
		return nil, 0, sharedErrors.NewInternalError("Failed to list addresses", err)
	}

	return addresses, total, nil
}

// Customer address operations

// GetCustomerAddresses retrieves all addresses for a customer
func (r *GormRepository) GetCustomerAddresses(ctx context.Context, tenantID, customerID uuid.UUID) ([]*Address, error) {
	var addresses []*Address
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error

	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get customer addresses", err)
	}

	return addresses, nil
}

// GetDefaultAddress retrieves the default address for a customer and type
func (r *GormRepository) GetDefaultAddress(ctx context.Context, tenantID, customerID uuid.UUID, addressType string) (*Address, error) {
	var address Address
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND customer_id = ? AND is_default = ?", tenantID, customerID, true)

	if addressType != "" {
		query = query.Where("type = ? OR type = ?", addressType, AddressTypeBoth)
	}

	err := query.First(&address).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Default address not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get default address", err)
	}

	return &address, nil
}

// SetDefaultAddress sets an address as default for a customer
func (r *GormRepository) SetDefaultAddress(ctx context.Context, tenantID, customerID, addressID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First, unset all default addresses for this customer
		if err := tx.Model(&Address{}).
			Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
			Update("is_default", false).Error; err != nil {
			return sharedErrors.NewInternalError("Failed to unset default addresses", err)
		}

		// Then set the specified address as default
		result := tx.Model(&Address{}).
			Where("tenant_id = ? AND customer_id = ? AND id = ?", tenantID, customerID, addressID).
			Update("is_default", true)

		if result.Error != nil {
			return sharedErrors.NewInternalError("Failed to set default address", result.Error)
		}

		if result.RowsAffected == 0 {
			return sharedErrors.NewNotFoundError("Address not found")
		}

		return nil
	})
}

// UnsetDefaultAddresses unsets default addresses for a customer and type
func (r *GormRepository) UnsetDefaultAddresses(ctx context.Context, tenantID, customerID uuid.UUID, addressType string) error {
	query := r.db.WithContext(ctx).Model(&Address{}).
		Where("tenant_id = ? AND customer_id = ? AND is_default = ?", tenantID, customerID, true)

	if addressType != "" {
		query = query.Where("type = ? OR type = ?", addressType, AddressTypeBoth)
	}

	if err := query.Update("is_default", false).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to unset default addresses", err)
	}
	return nil
}

// Address validation operations

// CreateValidation creates a new address validation
func (r *GormRepository) CreateValidation(ctx context.Context, validation *AddressValidation) error {
	if err := r.db.WithContext(ctx).Create(validation).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to create address validation", err)
	}
	return nil
}

// GetValidation retrieves the latest validation for an address
func (r *GormRepository) GetValidation(ctx context.Context, tenantID, addressID uuid.UUID) (*AddressValidation, error) {
	var validation AddressValidation
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND address_id = ?", tenantID, addressID).
		Order("validated_at DESC").
		First(&validation).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Address validation not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get address validation", err)
	}

	return &validation, nil
}

// GetValidationByID retrieves a validation by ID
func (r *GormRepository) GetValidationByID(ctx context.Context, tenantID, validationID uuid.UUID) (*AddressValidation, error) {
	var validation AddressValidation
	err := r.db.WithContext(ctx).
		Preload("Address").
		Where("tenant_id = ? AND id = ?", tenantID, validationID).
		First(&validation).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("Address validation not found")
		}
		return nil, sharedErrors.NewInternalError("Failed to get address validation", err)
	}

	return &validation, nil
}

// UpdateValidation updates an existing validation
func (r *GormRepository) UpdateValidation(ctx context.Context, validation *AddressValidation) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", validation.TenantID, validation.ID).
		Updates(validation)

	if result.Error != nil {
		return sharedErrors.NewInternalError("Failed to update address validation", result.Error)
	}

	if result.RowsAffected == 0 {
		return sharedErrors.NewNotFoundError("Address validation not found")
	}

	return nil
}

// ListValidations retrieves validations with pagination
func (r *GormRepository) ListValidations(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*AddressValidation, int64, error) {
	query := r.db.WithContext(ctx).
		Preload("Address").
		Where("address_validations.tenant_id = ?", tenantID)

	// Count total records
	var total int64
	if err := query.Model(&AddressValidation{}).Count(&total).Error; err != nil {
		return nil, 0, sharedErrors.NewInternalError("Failed to count address validations", err)
	}

	// Get paginated results
	var validations []*AddressValidation
	err := query.Order("validated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&validations).Error

	if err != nil {
		return nil, 0, sharedErrors.NewInternalError("Failed to list address validations", err)
	}

	return validations, total, nil
}

// Bulk operations

// BulkCreate creates multiple addresses
func (r *GormRepository) BulkCreate(ctx context.Context, addresses []*Address) error {
	if len(addresses) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(addresses, 100).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to bulk create addresses", err)
	}
	return nil
}

// BulkUpdate updates multiple addresses
func (r *GormRepository) BulkUpdate(ctx context.Context, addresses []*Address) error {
	if len(addresses) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, address := range addresses {
			if err := tx.Where("tenant_id = ? AND id = ?", address.TenantID, address.ID).
				Updates(address).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return sharedErrors.NewInternalError("Failed to bulk update addresses", err)
	}
	return nil
}

// BulkDelete soft deletes multiple addresses
func (r *GormRepository) BulkDelete(ctx context.Context, tenantID uuid.UUID, addressIDs []uuid.UUID) error {
	if len(addressIDs) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id IN ?", tenantID, addressIDs).
		Delete(&Address{}).Error; err != nil {
		return sharedErrors.NewInternalError("Failed to bulk delete addresses", err)
	}
	return nil
}

// Validation and checks

// ExistsByID checks if an address exists
func (r *GormRepository) ExistsByID(ctx context.Context, tenantID, addressID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Address{}).
		Where("tenant_id = ? AND id = ?", tenantID, addressID).
		Count(&count).Error

	if err != nil {
		return false, sharedErrors.NewInternalError("Failed to check address existence", err)
	}

	return count > 0, nil
}

// CountCustomerAddresses counts addresses for a customer
func (r *GormRepository) CountCustomerAddresses(ctx context.Context, tenantID, customerID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Address{}).
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Count(&count).Error

	if err != nil {
		return 0, sharedErrors.NewInternalError("Failed to count customer addresses", err)
	}

	return count, nil
}

// HasDefaultAddress checks if customer has a default address of specified type
func (r *GormRepository) HasDefaultAddress(ctx context.Context, tenantID, customerID uuid.UUID, addressType string) (bool, error) {
	query := r.db.WithContext(ctx).
		Model(&Address{}).
		Where("tenant_id = ? AND customer_id = ? AND is_default = ?", tenantID, customerID, true)

	if addressType != "" {
		query = query.Where("type = ? OR type = ?", addressType, AddressTypeBoth)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, sharedErrors.NewInternalError("Failed to check default address", err)
	}
	return count > 0, nil
}

// Statistics and analytics

// GetStats retrieves address statistics
func (r *GormRepository) GetStats(ctx context.Context, tenantID uuid.UUID) (*AddressStats, error) {
	stats := &AddressStats{
		AddressesByType:    make(map[string]int64),
		AddressesByCountry: make(map[string]int64),
	}

	// Total addresses
	if err := r.db.WithContext(ctx).
		Model(&Address{}).
		Where("tenant_id = ?", tenantID).
		Count(&stats.TotalAddresses).Error; err != nil {
		return nil, sharedErrors.NewInternalError("Failed to count total addresses", err)
	}

	// Validated addresses
	if err := r.db.WithContext(ctx).
		Model(&Address{}).
		Where("tenant_id = ? AND is_validated = ?", tenantID, true).
		Count(&stats.ValidatedAddresses).Error; err != nil {
		return nil, sharedErrors.NewInternalError("Failed to count validated addresses", err)
	}

	// Default addresses
	if err := r.db.WithContext(ctx).
		Model(&Address{}).
		Where("tenant_id = ? AND is_default = ?", tenantID, true).
		Count(&stats.DefaultAddresses).Error; err != nil {
		return nil, sharedErrors.NewInternalError("Failed to count default addresses", err)
	}

	// Recent addresses (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if err := r.db.WithContext(ctx).
		Model(&Address{}).
		Where("tenant_id = ? AND created_at >= ?", tenantID, thirtyDaysAgo).
		Count(&stats.RecentAddresses).Error; err != nil {
		return nil, sharedErrors.NewInternalError("Failed to count recent addresses", err)
	}

	// Addresses by type
	typeStats, err := r.GetAddressesByType(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	stats.AddressesByType = typeStats

	// Addresses by country
	countryStats, err := r.GetAddressesByCountry(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	stats.AddressesByCountry = countryStats

	return stats, nil
}

// GetAddressesByCountry retrieves address count by country
func (r *GormRepository) GetAddressesByCountry(ctx context.Context, tenantID uuid.UUID) (map[string]int64, error) {
	var results []struct {
		Country string `json:"country"`
		Count   int64  `json:"count"`
	}

	err := r.db.WithContext(ctx).
		Model(&Address{}).
		Select("country, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("country").
		Order("count DESC").
		Find(&results).Error

	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get addresses by country", err)
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.Country] = result.Count
	}

	return stats, nil
}

// GetAddressesByType retrieves address count by type
func (r *GormRepository) GetAddressesByType(ctx context.Context, tenantID uuid.UUID) (map[string]int64, error) {
	var results []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}

	err := r.db.WithContext(ctx).
		Model(&Address{}).
		Select("type, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("type").
		Order("count DESC").
		Find(&results).Error

	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get addresses by type", err)
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.Type] = result.Count
	}

	return stats, nil
}

// GetRecentAddresses retrieves recent addresses
func (r *GormRepository) GetRecentAddresses(ctx context.Context, tenantID uuid.UUID, days int) ([]*Address, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	var addresses []*Address
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND created_at >= ?", tenantID, cutoffDate).
		Order("created_at DESC").
		Find(&addresses).Error

	if err != nil {
		return nil, sharedErrors.NewInternalError("Failed to get recent addresses", err)
	}

	return addresses, nil
}

// Maintenance operations

// CleanupUnvalidatedAddresses removes old unvalidated addresses
func (r *GormRepository) CleanupUnvalidatedAddresses(ctx context.Context, tenantID uuid.UUID, days int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_validated = ? AND created_at < ?", tenantID, false, cutoffDate).
		Delete(&Address{})

	if result.Error != nil {
		return 0, sharedErrors.NewInternalError("Failed to cleanup unvalidated addresses", result.Error)
	}

	return result.RowsAffected, nil
}

// CleanupOrphanedValidations removes validations for non-existent addresses
func (r *GormRepository) CleanupOrphanedValidations(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND address_id NOT IN (SELECT id FROM addresses WHERE tenant_id = ?)", tenantID, tenantID).
		Delete(&AddressValidation{})

	if result.Error != nil {
		return 0, sharedErrors.NewInternalError("Failed to cleanup orphaned validations", result.Error)
	}

	return result.RowsAffected, nil
}

// Helper methods

// applyAddressFilters applies filters to the query
func (r *GormRepository) applyAddressFilters(query *gorm.DB, filter AddressFilter) *gorm.DB {
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.Label != "" {
		query = query.Where("label ILIKE ?", "%"+filter.Label+"%")
	}

	if filter.Country != "" {
		query = query.Where("country = ?", strings.ToUpper(filter.Country))
	}

	if filter.State != "" {
		query = query.Where("state ILIKE ?", "%"+filter.State+"%")
	}

	if filter.City != "" {
		query = query.Where("city ILIKE ?", "%"+filter.City+"%")
	}

	if filter.IsDefault != nil {
		query = query.Where("is_default = ?", *filter.IsDefault)
	}

	if filter.IsValidated != nil {
		query = query.Where("is_validated = ?", *filter.IsValidated)
	}

	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where(
			"LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(company) LIKE ? OR LOWER(address1) LIKE ? OR LOWER(city) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	return query
}
