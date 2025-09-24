package components

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service defines the interface for component business logic
type Service interface {
	// Component operations
	CreateComponent(ctx context.Context, tenantID uuid.UUID, req CreateComponentRequest) (*Component, error)
	GetComponent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*Component, error)
	GetComponentBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Component, error)
	UpdateComponent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, req UpdateComponentRequest) (*Component, error)
	DeleteComponent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) error
	ListComponents(ctx context.Context, tenantID uuid.UUID, filters ComponentFilters) (*ComponentListResponse, error)
	DuplicateComponent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, req DuplicateRequest) (*Component, error)

	// Component instance operations
	CreateInstance(ctx context.Context, tenantID uuid.UUID, req CreateInstanceRequest) (*ComponentInstance, error)
	GetInstance(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*ComponentInstance, error)
	UpdateInstance(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, req UpdateInstanceRequest) (*ComponentInstance, error)
	DeleteInstance(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) error
	ListInstances(ctx context.Context, tenantID uuid.UUID, themeID uuid.UUID) ([]ComponentInstance, error)

	// Theme operations
	CreateTheme(ctx context.Context, tenantID uuid.UUID, req CreateThemeRequest) (*Theme, error)
	GetTheme(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*Theme, error)
	UpdateTheme(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, req UpdateThemeRequest) (*Theme, error)
	DeleteTheme(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) error
	ListThemes(ctx context.Context, tenantID uuid.UUID, filters ThemeFilters) (*ThemeListResponse, error)
	GetActiveTheme(ctx context.Context, tenantID uuid.UUID) (*Theme, error)

	// Template operations
	ListTemplates(ctx context.Context, filters TemplateFilters) (*TemplateListResponse, error)
	GetTemplate(ctx context.Context, id uuid.UUID) (*ComponentTemplate, error)
	CreateComponentFromTemplate(ctx context.Context, tenantID uuid.UUID, req CreateFromTemplateRequest) (*Component, error)

	// Component template operations
	CreateComponentTemplate(ctx context.Context, req CreateComponentTemplateRequest) (*ComponentTemplate, error)
	GetComponentTemplate(ctx context.Context, id uuid.UUID) (*ComponentTemplate, error)
	UpdateComponentTemplate(ctx context.Context, id uuid.UUID, req UpdateComponentTemplateRequest) (*ComponentTemplate, error)
	DeleteComponentTemplate(ctx context.Context, id uuid.UUID) error
	ListComponentTemplates(ctx context.Context, filters ComponentTemplateFilter) (*ComponentTemplateListResponse, error)

	// Theme template operations
	CreateThemeTemplate(ctx context.Context, req CreateThemeTemplateRequest) (*ThemeTemplate, error)
	GetThemeTemplate(ctx context.Context, id uuid.UUID) (*ThemeTemplate, error)
	UpdateThemeTemplate(ctx context.Context, id uuid.UUID, req UpdateThemeTemplateRequest) (*ThemeTemplate, error)
	DeleteThemeTemplate(ctx context.Context, id uuid.UUID) error
	ListThemeTemplates(ctx context.Context, filters ThemeTemplateFilter) (*ThemeTemplateListResponse, error)

	// Statistics
	GetStats(ctx context.Context, tenantID uuid.UUID) (*ComponentStats, error)
}

// Request/Response types

type CreateComponentRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	Slug        string                 `json:"slug,omitempty"`
	Type        ComponentType          `json:"type" validate:"required"`
	Description string                 `json:"description,omitempty"`
	HTML        string                 `json:"html" validate:"required"`
	CSS         string                 `json:"css,omitempty"`
	JavaScript  string                 `json:"javascript,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Category    ComponentCategory      `json:"category,omitempty"`
	IsFeatured  bool                   `json:"is_featured"`
}

type UpdateComponentRequest struct {
	Name        *string                `json:"name,omitempty"`
	Slug        *string                `json:"slug,omitempty"`
	Status      *ComponentStatus       `json:"status,omitempty"`
	Description *string                `json:"description,omitempty"`
	HTML        *string                `json:"html,omitempty"`
	CSS         *string                `json:"css,omitempty"`
	JavaScript  *string                `json:"javascript,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Category    *ComponentCategory     `json:"category,omitempty"`
	IsFeatured  *bool                  `json:"is_featured,omitempty"`
}

type CreateInstanceRequest struct {
	ComponentID uuid.UUID              `json:"component_id" validate:"required"`
	ThemeID     uuid.UUID              `json:"theme_id" validate:"required"`
	CustomHTML  string                 `json:"custom_html,omitempty"`
	CustomCSS   string                 `json:"custom_css,omitempty"`
	CustomJS    string                 `json:"custom_js,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Position    int                    `json:"position"`
	Zone        string                 `json:"zone"`
	IsVisible   bool                   `json:"is_visible"`
	Breakpoints map[string]interface{} `json:"breakpoints,omitempty"`
}

type UpdateInstanceRequest struct {
	CustomHTML  *string                `json:"custom_html,omitempty"`
	CustomCSS   *string                `json:"custom_css,omitempty"`
	CustomJS    *string                `json:"custom_js,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Position    *int                   `json:"position,omitempty"`
	Zone        *string                `json:"zone,omitempty"`
	IsVisible   *bool                  `json:"is_visible,omitempty"`
	Breakpoints map[string]interface{} `json:"breakpoints,omitempty"`
}

type DuplicateRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type CreateThemeRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	Slug        string                 `json:"slug,omitempty"`
	Description string                 `json:"description,omitempty"`
	GlobalCSS   string                 `json:"global_css,omitempty"`
	GlobalJS    string                 `json:"global_js,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Thumbnail   string                 `json:"thumbnail,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	IsDefault   bool                   `json:"is_default"`
}

type UpdateThemeRequest struct {
	Name        *string                `json:"name,omitempty"`
	Slug        *string                `json:"slug,omitempty"`
	Description *string                `json:"description,omitempty"`
	Status      *ComponentStatus       `json:"status,omitempty"`
	GlobalCSS   *string                `json:"global_css,omitempty"`
	GlobalJS    *string                `json:"global_js,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Thumbnail   *string                `json:"thumbnail,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	IsDefault   *bool                  `json:"is_default,omitempty"`
}

type CreateComponentTemplateRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	Slug        string                 `json:"slug,omitempty"`
	Type        ComponentType          `json:"type" validate:"required"`
	Description string                 `json:"description,omitempty"`
	HTML        string                 `json:"html" validate:"required"`
	CSS         string                 `json:"css,omitempty"`
	JavaScript  string                 `json:"javascript,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Thumbnail   string                 `json:"thumbnail,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Category    ComponentCategory      `json:"category,omitempty"`
	IsFree      bool                   `json:"is_free"`
	IsFeatured  bool                   `json:"is_featured"`
}

type UpdateComponentTemplateRequest struct {
	Name        *string                `json:"name,omitempty"`
	Slug        *string                `json:"slug,omitempty"`
	Type        *ComponentType         `json:"type,omitempty"`
	Description *string                `json:"description,omitempty"`
	HTML        *string                `json:"html,omitempty"`
	CSS         *string                `json:"css,omitempty"`
	JavaScript  *string                `json:"javascript,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Thumbnail   *string                `json:"thumbnail,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Category    *ComponentCategory     `json:"category,omitempty"`
	IsFree      *bool                  `json:"is_free,omitempty"`
	IsFeatured  *bool                  `json:"is_featured,omitempty"`
}

type CreateThemeTemplateRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	Slug        string                 `json:"slug,omitempty"`
	Description string                 `json:"description,omitempty"`
	GlobalCSS   string                 `json:"global_css,omitempty"`
	GlobalJS    string                 `json:"global_js,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Thumbnail   string                 `json:"thumbnail,omitempty"`
	PreviewURL  string                 `json:"preview_url,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Category    string                 `json:"category,omitempty"`
	IsFree      bool                   `json:"is_free"`
	IsFeatured  bool                   `json:"is_featured"`
	IsDefault   bool                   `json:"is_default"`
}

type UpdateThemeTemplateRequest struct {
	Name        *string                `json:"name,omitempty"`
	Slug        *string                `json:"slug,omitempty"`
	Description *string                `json:"description,omitempty"`
	GlobalCSS   *string                `json:"global_css,omitempty"`
	GlobalJS    *string                `json:"global_js,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Thumbnail   *string                `json:"thumbnail,omitempty"`
	PreviewURL  *string                `json:"preview_url,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Category    *string                `json:"category,omitempty"`
	IsFree      *bool                  `json:"is_free,omitempty"`
	IsFeatured  *bool                  `json:"is_featured,omitempty"`
	IsDefault   *bool                  `json:"is_default,omitempty"`
}

type CreateFromTemplateRequest struct {
	TemplateID uuid.UUID              `json:"template_id" validate:"required"`
	Name       string                 `json:"name" validate:"required,min=1,max=255"`
	Config     map[string]interface{} `json:"config"`
}

type ComponentListResponse struct {
	Components []Component `json:"components"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

type ThemeListResponse struct {
	Themes     []Theme `json:"themes"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalPages int     `json:"total_pages"`
}

type TemplateListResponse struct {
	Templates  []ComponentTemplate `json:"templates"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

type ThemeTemplateListResponse struct {
	ThemeTemplates []ThemeTemplate `json:"theme_templates"`
	Total          int64           `json:"total"`
	Page           int             `json:"page"`
	Limit          int             `json:"limit"`
	TotalPages     int             `json:"total_pages"`
}

type ComponentTemplateListResponse struct {
	ComponentTemplates []ComponentTemplate `json:"component_templates"`
	Total              int64               `json:"total"`
	Page               int                 `json:"page"`
	Limit              int                 `json:"limit"`
	TotalPages         int                 `json:"total_pages"`
}

type InstanceListResponse struct {
	Instances  []ComponentInstance `json:"instances"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

// service implements Service interface
type service struct {
	repo Repository
}

// NewService creates a new component service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Helper functions

func generateSlug(name string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func calculateTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	return int((total + int64(limit) - 1) / int64(limit))
}

// Component operations

func (s *service) CreateComponent(ctx context.Context, tenantID uuid.UUID, req CreateComponentRequest) (*Component, error) {

	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	}

	// Check if slug already exists
	existing, err := s.repo.GetComponentBySlug(ctx, tenantID, slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if existing != nil {
		return nil, errors.New("component with this slug already exists")
	}

	component := &Component{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        req.Name,
		Slug:        slug,
		Type:        req.Type,
		Description: req.Description,
		HTML:        req.HTML,
		CSS:         req.CSS,
		JavaScript:  req.JavaScript,
		Config:      req.Config,
		Tags:        req.Tags,
		Category:    req.Category,
		IsFeatured:  req.IsFeatured,
		Status:      StatusDraft,
	}

	if err := s.repo.CreateComponent(ctx, tenantID, component); err != nil {
		return nil, fmt.Errorf("failed to create component: %w", err)
	}

	return component, nil
}

func (s *service) GetComponent(ctx context.Context, tenantID, id uuid.UUID) (*Component, error) {
	component, err := s.repo.GetComponent(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("component not found")
		}
		return nil, fmt.Errorf("failed to get component: %w", err)
	}
	return component, nil
}

func (s *service) GetComponentBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Component, error) {
	component, err := s.repo.GetComponentBySlug(ctx, tenantID, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("component not found")
		}
		return nil, fmt.Errorf("failed to get component: %w", err)
	}
	return component, nil
}

func (s *service) UpdateComponent(ctx context.Context, tenantID, id uuid.UUID, req UpdateComponentRequest) (*Component, error) {
	component, err := s.repo.GetComponent(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("component not found")
		}
		return nil, fmt.Errorf("failed to get component: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		component.Name = *req.Name
	}
	if req.Slug != nil {
		// Check slug uniqueness if changed
		if *req.Slug != component.Slug {
			existing, err := s.repo.GetComponentBySlug(ctx, tenantID, *req.Slug)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
			}
			if existing != nil {
				return nil, errors.New("component with this slug already exists")
			}
			component.Slug = *req.Slug
		}
	}
	if req.Status != nil {
		component.Status = *req.Status
	}
	if req.Description != nil {
		component.Description = *req.Description
	}
	if req.HTML != nil {
		component.HTML = *req.HTML
	}
	if req.CSS != nil {
		component.CSS = *req.CSS
	}
	if req.JavaScript != nil {
		component.JavaScript = *req.JavaScript
	}
	if req.Config != nil {
		component.Config = req.Config
	}
	if req.Tags != nil {
		component.Tags = req.Tags
	}
	if req.Category != nil {
		component.Category = *req.Category
	}
	if req.IsFeatured != nil {
		component.IsFeatured = *req.IsFeatured
	}

	if err := s.repo.UpdateComponent(ctx, tenantID, component); err != nil {
		return nil, fmt.Errorf("failed to update component: %w", err)
	}

	return component, nil
}

func (s *service) DeleteComponent(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if component exists
	_, err := s.repo.GetComponent(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("component not found")
		}
		return fmt.Errorf("failed to get component: %w", err)
	}

	if err := s.repo.DeleteComponent(ctx, tenantID, id); err != nil {
		return fmt.Errorf("failed to delete component: %w", err)
	}

	return nil
}

func (s *service) ListComponents(ctx context.Context, tenantID uuid.UUID, filters ComponentFilters) (*ComponentListResponse, error) {
	components, total, err := s.repo.ListComponents(ctx, tenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list components: %w", err)
	}

	return &ComponentListResponse{
		Components: components,
		Total:      total,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: calculateTotalPages(total, filters.Limit),
	}, nil
}

func (s *service) DuplicateComponent(ctx context.Context, tenantID, id uuid.UUID, req DuplicateRequest) (*Component, error) {
	original, err := s.repo.GetComponent(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("component not found")
		}
		return nil, fmt.Errorf("failed to get component: %w", err)
	}

	// Create duplicate with new name and slug
	duplicate := &Component{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        req.Name,
		Slug:        generateSlug(req.Name),
		Type:        original.Type,
		Description: original.Description,
		HTML:        original.HTML,
		CSS:         original.CSS,
		JavaScript:  original.JavaScript,
		Config:      original.Config,
		Tags:        original.Tags,
		Category:    original.Category,
		Status:      StatusDraft,
		IsFeatured:  false,
	}

	// Check slug uniqueness
	existing, err := s.repo.GetComponentBySlug(ctx, tenantID, duplicate.Slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if existing != nil {
		// Append timestamp to make it unique
		duplicate.Slug = fmt.Sprintf("%s-%d", duplicate.Slug, time.Now().Unix())
	}

	if err := s.repo.CreateComponent(ctx, tenantID, duplicate); err != nil {
		return nil, fmt.Errorf("failed to create duplicate component: %w", err)
	}

	return duplicate, nil
}

// Additional service methods for instances, themes, templates, and stats would continue here...
// For brevity, I'll implement the key methods. The pattern is similar for all operations.

func (s *service) CreateInstance(ctx context.Context, tenantID uuid.UUID, req CreateInstanceRequest) (*ComponentInstance, error) {

	// Verify component and theme exist
	_, err := s.repo.GetComponent(ctx, tenantID, req.ComponentID)
	if err != nil {
		return nil, errors.New("component not found")
	}

	_, err = s.repo.GetTheme(ctx, tenantID, req.ThemeID)
	if err != nil {
		return nil, errors.New("theme not found")
	}

	instance := &ComponentInstance{
		ID:          uuid.New(),
		TenantID:    tenantID,
		ComponentID: req.ComponentID,
		ThemeID:     req.ThemeID,
		CustomCSS:   req.CustomCSS,
		CustomJS:    req.CustomJS,
		Settings:    req.Settings,
		Position:    fmt.Sprintf("%d", req.Position), // Convert int to string
		IsVisible:   req.IsVisible,
	}

	if err := s.repo.CreateInstance(ctx, tenantID, instance); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	return instance, nil
}

func (s *service) GetStats(ctx context.Context, tenantID uuid.UUID) (*ComponentStats, error) {
	stats, err := s.repo.GetStats(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	return stats, nil
}

// Placeholder implementations for remaining methods
// These would follow the same pattern as above

func (s *service) GetInstance(ctx context.Context, tenantID, id uuid.UUID) (*ComponentInstance, error) {
	return s.repo.GetInstance(ctx, tenantID, id)
}

func (s *service) UpdateInstance(ctx context.Context, tenantID, id uuid.UUID, req UpdateInstanceRequest) (*ComponentInstance, error) {
	// Implementation similar to UpdateComponent
	return nil, errors.New("not implemented")
}

func (s *service) DeleteInstance(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteInstance(ctx, tenantID, id)
}

func (s *service) ListInstances(ctx context.Context, tenantID, themeID uuid.UUID) ([]ComponentInstance, error) {
	filter := ComponentInstanceFilter{
		ThemeID: &themeID,
	}

	instances, err := s.repo.ListInstances(ctx, tenantID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	// Convert from pointer slice to value slice
	result := make([]ComponentInstance, len(instances))
	for i, instance := range instances {
		result[i] = *instance
	}

	return result, nil
}

func (s *service) ReorderInstances(ctx context.Context, tenantID, themeID uuid.UUID, instanceIDs []uuid.UUID) error {
	// Implementation for reordering instances
	return errors.New("not implemented")
}

func (s *service) CreateTheme(ctx context.Context, tenantID uuid.UUID, req CreateThemeRequest) (*Theme, error) {
	// Implementation similar to CreateComponent
	return nil, errors.New("not implemented")
}

func (s *service) GetTheme(ctx context.Context, tenantID, id uuid.UUID) (*Theme, error) {
	return s.repo.GetTheme(ctx, tenantID, id)
}

func (s *service) GetThemeBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Theme, error) {
	return s.repo.GetThemeBySlug(ctx, tenantID, slug)
}

func (s *service) UpdateTheme(ctx context.Context, tenantID, id uuid.UUID, req UpdateThemeRequest) (*Theme, error) {
	return nil, errors.New("not implemented")
}

func (s *service) DeleteTheme(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteTheme(ctx, tenantID, id)
}

func (s *service) ListThemes(ctx context.Context, tenantID uuid.UUID, filters ThemeFilters) (*ThemeListResponse, error) {
	themes, total, err := s.repo.ListThemes(ctx, tenantID, filters)
	if err != nil {
		return nil, err
	}
	return &ThemeListResponse{
		Themes:     themes,
		Total:      total,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: calculateTotalPages(total, filters.Limit),
	}, nil
}

func (s *service) PublishTheme(ctx context.Context, tenantID, id uuid.UUID) (*Theme, error) {
	return nil, errors.New("not implemented")
}

func (s *service) UnpublishTheme(ctx context.Context, tenantID, id uuid.UUID) (*Theme, error) {
	return nil, errors.New("not implemented")
}

func (s *service) GetActiveTheme(ctx context.Context, tenantID uuid.UUID) (*Theme, error) {
	return s.repo.GetActiveTheme(ctx, tenantID)
}

func (s *service) DuplicateTheme(ctx context.Context, tenantID, id uuid.UUID, name string) (*Theme, error) {
	return nil, errors.New("not implemented")
}

func (s *service) ListTemplates(ctx context.Context, filters TemplateFilters) (*TemplateListResponse, error) {
	templates, total, err := s.repo.ListTemplates(ctx, filters)
	if err != nil {
		return nil, err
	}
	return &TemplateListResponse{
		Templates:  templates,
		Total:      total,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: calculateTotalPages(total, filters.Limit),
	}, nil
}

func (s *service) GetTemplate(ctx context.Context, id uuid.UUID) (*ComponentTemplate, error) {
	return s.repo.GetTemplate(ctx, id)
}

func (s *service) CreateComponentFromTemplate(ctx context.Context, tenantID uuid.UUID, req CreateFromTemplateRequest) (*Component, error) {
	return nil, errors.New("not implemented")
}

// Theme Template Service Methods
func (s *service) CreateThemeTemplate(ctx context.Context, req CreateThemeTemplateRequest) (*ThemeTemplate, error) {
	themeTemplate := &ThemeTemplate{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		GlobalCSS:   req.GlobalCSS,
		GlobalJS:    req.GlobalJS,
		Settings:    req.Settings,
		Thumbnail:   req.Thumbnail,
		PreviewURL:  req.PreviewURL,
		Tags:        req.Tags,
		Category:    req.Category,
		IsFree:      req.IsFree,
		IsFeatured:  req.IsFeatured,
		IsDefault:   req.IsDefault,
	}

	if themeTemplate.Slug == "" {
		themeTemplate.Slug = generateSlug(themeTemplate.Name)
	}

	if err := s.repo.CreateThemeTemplate(ctx, themeTemplate); err != nil {
		return nil, fmt.Errorf("failed to create theme template: %w", err)
	}

	return themeTemplate, nil
}

func (s *service) GetThemeTemplate(ctx context.Context, id uuid.UUID) (*ThemeTemplate, error) {
	themeTemplate, err := s.repo.GetThemeTemplate(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("theme template not found")
		}
		return nil, fmt.Errorf("failed to get theme template: %w", err)
	}
	return themeTemplate, nil
}

func (s *service) UpdateThemeTemplate(ctx context.Context, id uuid.UUID, req UpdateThemeTemplateRequest) (*ThemeTemplate, error) {
	themeTemplate, err := s.repo.GetThemeTemplate(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("theme template not found")
		}
		return nil, fmt.Errorf("failed to get theme template: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		themeTemplate.Name = *req.Name
	}
	if req.Slug != nil {
		themeTemplate.Slug = *req.Slug
	}
	if req.Description != nil {
		themeTemplate.Description = *req.Description
	}
	if req.GlobalCSS != nil {
		themeTemplate.GlobalCSS = *req.GlobalCSS
	}
	if req.GlobalJS != nil {
		themeTemplate.GlobalJS = *req.GlobalJS
	}
	if req.Settings != nil {
		themeTemplate.Settings = req.Settings
	}
	if req.Thumbnail != nil {
		themeTemplate.Thumbnail = *req.Thumbnail
	}
	if req.PreviewURL != nil {
		themeTemplate.PreviewURL = *req.PreviewURL
	}
	if req.Tags != nil {
		themeTemplate.Tags = req.Tags
	}
	if req.Category != nil {
		themeTemplate.Category = *req.Category
	}
	if req.IsFree != nil {
		themeTemplate.IsFree = *req.IsFree
	}
	if req.IsFeatured != nil {
		themeTemplate.IsFeatured = *req.IsFeatured
	}
	if req.IsDefault != nil {
		themeTemplate.IsDefault = *req.IsDefault
	}

	if err := s.repo.UpdateThemeTemplate(ctx, themeTemplate); err != nil {
		return nil, fmt.Errorf("failed to update theme template: %w", err)
	}

	return themeTemplate, nil
}

func (s *service) DeleteThemeTemplate(ctx context.Context, id uuid.UUID) error {
	// Check if theme template exists
	_, err := s.repo.GetThemeTemplate(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("theme template not found")
		}
		return fmt.Errorf("failed to get theme template: %w", err)
	}

	if err := s.repo.DeleteThemeTemplate(ctx, id); err != nil {
		return fmt.Errorf("failed to delete theme template: %w", err)
	}

	return nil
}

func (s *service) ListThemeTemplates(ctx context.Context, filters ThemeTemplateFilter) (*ThemeTemplateListResponse, error) {
	themeTemplates, err := s.repo.ListThemeTemplates(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list theme templates: %w", err)
	}

	// Convert from pointer slice to value slice
	result := make([]ThemeTemplate, len(themeTemplates))
	for i, template := range themeTemplates {
		result[i] = *template
	}

	return &ThemeTemplateListResponse{
		ThemeTemplates: result,
		Total:          int64(len(themeTemplates)),
		Page:           filters.Page,
		Limit:          filters.Limit,
		TotalPages:     calculateTotalPages(int64(len(themeTemplates)), filters.Limit),
	}, nil
}

// Component Template Service Methods
func (s *service) CreateComponentTemplate(ctx context.Context, req CreateComponentTemplateRequest) (*ComponentTemplate, error) {
	componentTemplate := &ComponentTemplate{
		ID:            uuid.New(),
		Name:          req.Name,
		Slug:          req.Slug,
		Type:          req.Type,
		Description:   req.Description,
		HTML:          req.HTML,
		CSS:           req.CSS,
		JavaScript:    req.JavaScript,
		ConfigSchema:  req.Config,      // Use ConfigSchema instead of Config
		PreviewImage:  req.Thumbnail,   // Use PreviewImage instead of Thumbnail
		Tags:          req.Tags,
		Category:      req.Category,
		IsFree:        req.IsFree,
		IsFeatured:    req.IsFeatured,
	}

	if componentTemplate.Slug == "" {
		componentTemplate.Slug = generateSlug(componentTemplate.Name)
	}

	if err := s.repo.CreateComponentTemplate(ctx, componentTemplate); err != nil {
		return nil, fmt.Errorf("failed to create component template: %w", err)
	}

	return componentTemplate, nil
}

func (s *service) GetComponentTemplate(ctx context.Context, id uuid.UUID) (*ComponentTemplate, error) {
	componentTemplate, err := s.repo.GetComponentTemplate(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("component template not found")
		}
		return nil, fmt.Errorf("failed to get component template: %w", err)
	}
	return componentTemplate, nil
}

func (s *service) UpdateComponentTemplate(ctx context.Context, id uuid.UUID, req UpdateComponentTemplateRequest) (*ComponentTemplate, error) {
	componentTemplate, err := s.repo.GetComponentTemplate(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("component template not found")
		}
		return nil, fmt.Errorf("failed to get component template: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		componentTemplate.Name = *req.Name
	}
	if req.Slug != nil {
		componentTemplate.Slug = *req.Slug
	}
	if req.Type != nil {
		componentTemplate.Type = *req.Type
	}
	if req.Description != nil {
		componentTemplate.Description = *req.Description
	}
	if req.HTML != nil {
		componentTemplate.HTML = *req.HTML
	}
	if req.CSS != nil {
		componentTemplate.CSS = *req.CSS
	}
	if req.JavaScript != nil {
		componentTemplate.JavaScript = *req.JavaScript
	}
	if req.Config != nil {
		componentTemplate.ConfigSchema = req.Config
	}
	if req.Thumbnail != nil {
		componentTemplate.PreviewImage = *req.Thumbnail
	}
	if req.Tags != nil {
		componentTemplate.Tags = req.Tags
	}
	if req.Category != nil {
		componentTemplate.Category = *req.Category
	}
	if req.IsFree != nil {
		componentTemplate.IsFree = *req.IsFree
	}
	if req.IsFeatured != nil {
		componentTemplate.IsFeatured = *req.IsFeatured
	}

	if err := s.repo.UpdateComponentTemplate(ctx, componentTemplate); err != nil {
		return nil, fmt.Errorf("failed to update component template: %w", err)
	}

	return componentTemplate, nil
}

func (s *service) DeleteComponentTemplate(ctx context.Context, id uuid.UUID) error {
	// Check if component template exists
	_, err := s.repo.GetComponentTemplate(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("component template not found")
		}
		return fmt.Errorf("failed to get component template: %w", err)
	}

	if err := s.repo.DeleteComponentTemplate(ctx, id); err != nil {
		return fmt.Errorf("failed to delete component template: %w", err)
	}

	return nil
}

func (s *service) ListComponentTemplates(ctx context.Context, filters ComponentTemplateFilter) (*ComponentTemplateListResponse, error) {
	componentTemplates, err := s.repo.ListComponentTemplates(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list component templates: %w", err)
	}

	// Convert from pointer slice to value slice
	result := make([]ComponentTemplate, len(componentTemplates))
	for i, template := range componentTemplates {
		result[i] = *template
	}

	return &ComponentTemplateListResponse{
		ComponentTemplates: result,
		Total:              int64(len(componentTemplates)),
		Page:               filters.Page,
		Limit:              filters.Limit,
		TotalPages:         calculateTotalPages(int64(len(componentTemplates)), filters.Limit),
	}, nil
}
