package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// Service defines the settings service interface
type Service interface {
	GetSettings(ctx context.Context, tenantID uint, req *GetSettingsRequest) (*GetSettingsResponse, error)
	UpdateSettings(ctx context.Context, tenantID uint, req *UpdateSettingsRequest) (*UpdateSettingsResponse, error)
	GetPublicSettings(ctx context.Context, tenantID uint) (*PublicSettingsResponse, error)
}

// Repository defines the settings repository interface
type Repository interface {
	GetSettings(tenantID uint, section string) ([]Setting, error)
	GetSetting(tenantID uint, section, key string) (*Setting, error)
	CreateSetting(setting *Setting) error
	UpdateSetting(setting *Setting) error
	DeleteSetting(tenantID uint, section, key string) error
	GetPublicSettings(tenantID uint) ([]Setting, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new settings service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetSettings retrieves settings for a tenant
func (s *service) GetSettings(ctx context.Context, tenantID uint, req *GetSettingsRequest) (*GetSettingsResponse, error) {
	log.Printf("Getting settings for tenant %d, section: %s", tenantID, req.Section)

	settings, err := s.repo.GetSettings(tenantID, req.Section)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	// Convert settings to map
	settingsMap := make(map[string]interface{})
	sections := make(map[string]bool)

	for _, setting := range settings {
		// Parse JSON values
		var value interface{}
		if setting.Type == "json" {
			if err := json.Unmarshal([]byte(setting.Value), &value); err != nil {
				value = setting.Value // fallback to string
			}
		} else {
			value = setting.Value
		}

		// Group by section
		if req.Section == "" || req.Section == setting.Section {
			if settingsMap[setting.Section] == nil {
				settingsMap[setting.Section] = make(map[string]interface{})
			}
			settingsMap[setting.Section].(map[string]interface{})[setting.Key] = value
			sections[setting.Section] = true
		}
	}

	// If no settings found, return defaults for the section
	if len(settingsMap) == 0 && req.Section != "" {
		if defaults, exists := DefaultSettings[req.Section]; exists {
			settingsMap[req.Section] = defaults
		}
	}

	// Convert sections map to slice
	sectionsList := make([]string, 0, len(sections))
	for section := range sections {
		sectionsList = append(sectionsList, section)
	}

	response := &GetSettingsResponse{
		Settings: settingsMap,
		Sections: sectionsList,
	}

	return response, nil
}

// UpdateSettings updates settings for a tenant
func (s *service) UpdateSettings(ctx context.Context, tenantID uint, req *UpdateSettingsRequest) (*UpdateSettingsResponse, error) {
	log.Printf("Updating settings for tenant %d", tenantID)

	if len(req.Settings) == 0 {
		return nil, fmt.Errorf("no settings provided")
	}

	updated := make([]string, 0)
	updatedSettings := make(map[string]interface{})

	// Process each section in the settings
	for section, sectionData := range req.Settings {
		sectionMap, ok := sectionData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid settings format for section %s", section)
		}

		// Process each key-value pair in the section
		for key, value := range sectionMap {
			// Get existing setting
			existingSetting, err := s.repo.GetSetting(tenantID, section, key)
			if err != nil && err.Error() != "setting not found" {
				return nil, fmt.Errorf("failed to get existing setting: %w", err)
			}

			// Determine value type and convert to string
			var valueStr string
			var valueType string

			switch v := value.(type) {
			case string:
				valueStr = v
				valueType = "string"
			case bool:
				valueStr = fmt.Sprintf("%t", v)
				valueType = "boolean"
			case float64:
				valueStr = fmt.Sprintf("%g", v)
				valueType = "number"
			default:
				// Convert complex types to JSON
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal value for %s.%s: %w", section, key, err)
				}
				valueStr = string(jsonBytes)
				valueType = "json"
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

			// Create or update setting
			if existingSetting != nil {
				// Update existing setting
				existingSetting.Value = valueStr
				existingSetting.Type = valueType
				existingSetting.IsPublic = isPublic
				if err := s.repo.UpdateSetting(existingSetting); err != nil {
					return nil, fmt.Errorf("failed to update setting %s.%s: %w", section, key, err)
				}
			} else {
				// Create new setting
				newSetting := &Setting{
					TenantID: tenantID,
					Section:  section,
					Key:      key,
					Value:    valueStr,
					Type:     valueType,
					IsPublic: isPublic,
				}
				if err := s.repo.CreateSetting(newSetting); err != nil {
					return nil, fmt.Errorf("failed to create setting %s.%s: %w", section, key, err)
				}
			}

			updated = append(updated, fmt.Sprintf("%s.%s", section, key))
			updatedSettings[fmt.Sprintf("%s.%s", section, key)] = value
		}
	}

	response := &UpdateSettingsResponse{
		Message:  fmt.Sprintf("Successfully updated %d settings", len(updated)),
		Settings: updatedSettings,
		Updated:  updated,
	}

	return response, nil
}

// GetPublicSettings retrieves public settings for a tenant
func (s *service) GetPublicSettings(ctx context.Context, tenantID uint) (*PublicSettingsResponse, error) {
	log.Printf("Getting public settings for tenant %d", tenantID)

	settings, err := s.repo.GetPublicSettings(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get public settings: %w", err)
	}

	// Initialize response with defaults
	response := &PublicSettingsResponse{
		StoreName:   "My Store",
		Logo:        "",
		Contact:     make(map[string]interface{}),
		Theme:       make(map[string]interface{}),
		SEO:         make(map[string]interface{}),
		SocialLinks: make(map[string]interface{}),
	}

	// Process settings
	for _, setting := range settings {
		var value interface{}
		if setting.Type == "json" {
			if err := json.Unmarshal([]byte(setting.Value), &value); err != nil {
				value = setting.Value // fallback to string
			}
		} else {
			value = setting.Value
		}

		switch setting.Section {
		case "general":
			switch setting.Key {
			case "store_name":
				response.StoreName = setting.Value
			case "store_email", "store_phone", "store_address":
				response.Contact[setting.Key] = value
			}
		case "appearance":
			switch setting.Key {
			case "logo":
				response.Logo = setting.Value
			default:
				response.Theme[setting.Key] = value
			}
		case "seo":
			response.SEO[setting.Key] = value
		}
	}

	return response, nil
}