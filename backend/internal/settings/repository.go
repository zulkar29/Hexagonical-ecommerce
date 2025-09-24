package settings

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new settings repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// GetSettings retrieves settings for a tenant, optionally filtered by section
func (r *repository) GetSettings(tenantID uuid.UUID, section string) ([]Setting, error) {
	var settings []Setting
	query := r.db.Where("tenant_id = ?", tenantID)

	if section != "" {
		query = query.Where("section = ?", section)
	}

	if err := query.Order("section, key").Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	return settings, nil
}

// GetSetting retrieves a specific setting
func (r *repository) GetSetting(tenantID uuid.UUID, section, key string) (*Setting, error) {
	var setting Setting
	err := r.db.Where("tenant_id = ? AND section = ? AND key = ?", tenantID, section, key).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("setting not found")
		}
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}

	return &setting, nil
}

// CreateSetting creates a new setting
func (r *repository) CreateSetting(setting *Setting) error {
	// Generate UUID for new setting if not set
	if setting.ID == uuid.Nil {
		setting.ID = uuid.New()
	}
	if err := r.db.Create(setting).Error; err != nil {
		return fmt.Errorf("failed to create setting: %w", err)
	}
	return nil
}

// UpdateSetting updates an existing setting
func (r *repository) UpdateSetting(setting *Setting) error {
	if err := r.db.Save(setting).Error; err != nil {
		return fmt.Errorf("failed to update setting: %w", err)
	}
	return nil
}

// DeleteSetting deletes a setting
func (r *repository) DeleteSetting(tenantID uuid.UUID, section, key string) error {
	result := r.db.Where("tenant_id = ? AND section = ? AND key = ?", tenantID, section, key).Delete(&Setting{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete setting: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("setting not found")
	}
	return nil
}

// GetPublicSettings retrieves all public settings for a tenant
func (r *repository) GetPublicSettings(tenantID uuid.UUID) ([]Setting, error) {
	var settings []Setting
	if err := r.db.Where("tenant_id = ? AND is_public = ?", tenantID, true).Order("section, key").Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("failed to get public settings: %w", err)
	}

	return settings, nil
}

// GetSettingsByKeys retrieves specific settings by their keys
func (r *repository) GetSettingsByKeys(tenantID uuid.UUID, keys []string) ([]Setting, error) {
	var settings []Setting
	if err := r.db.Where("tenant_id = ? AND key IN ?", tenantID, keys).Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("failed to get settings by keys: %w", err)
	}

	return settings, nil
}

// BulkCreateSettings creates multiple settings in a transaction
func (r *repository) BulkCreateSettings(settings []Setting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i := range settings {
			// Generate UUID for new setting if not set
			if settings[i].ID == uuid.Nil {
				settings[i].ID = uuid.New()
			}
			if err := tx.Create(&settings[i]).Error; err != nil {
				return fmt.Errorf("failed to create setting %s.%s: %w", settings[i].Section, settings[i].Key, err)
			}
		}
		return nil
	})
}

// BulkUpdateSettings updates multiple settings in a transaction
func (r *repository) BulkUpdateSettings(settings []Setting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, setting := range settings {
			if err := tx.Save(&setting).Error; err != nil {
				return fmt.Errorf("failed to update setting %s.%s: %w", setting.Section, setting.Key, err)
			}
		}
		return nil
	})
}

// InitializeDefaultSettings creates default settings for a new tenant
func (r *repository) InitializeDefaultSettings(ctx context.Context, tenantID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for section, sectionSettings := range DefaultSettings {
			for key, value := range sectionSettings {
				// Check if setting already exists
				var count int64
				tx.Model(&Setting{}).Where("tenant_id = ? AND section = ? AND key = ?", tenantID, section, key).Count(&count)
				if count > 0 {
					continue // Skip if already exists
				}

				// Determine if setting is public
				isPublic := false
				if publicKeys, exists := PublicSettingKeys[section]; exists {
					for _, publicKey := range publicKeys {
						if publicKey == key {
							isPublic = true
							break
						}
					}
				}

				// Determine value type
				valueType := "string"
				valueStr := fmt.Sprintf("%v", value)

				switch value.(type) {
				case bool:
					valueType = "boolean"
				case int, int32, int64, float32, float64:
					valueType = "number"
				}

				setting := Setting{
					ID:       uuid.New(),
					TenantID: tenantID,
					Section:  section,
					Key:      key,
					Value:    valueStr,
					Type:     valueType,
					IsPublic: isPublic,
				}

				if err := tx.Create(&setting).Error; err != nil {
					return fmt.Errorf("failed to create default setting %s.%s: %w", section, key, err)
				}
			}
		}
		return nil
	})
}

// GetSettingsSummary returns a summary of settings by section
func (r *repository) GetSettingsSummary(ctx context.Context, tenantID uuid.UUID) ([]*Setting, error) {
	var settings []*Setting
	if err := r.db.Where("tenant_id = ?", tenantID).Order("section, key").Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("failed to get settings summary: %w", err)
	}

	return settings, nil
}
