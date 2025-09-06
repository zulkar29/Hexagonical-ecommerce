package components

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for component data operations
type Repository interface {
	// Component operations
	CreateComponent(ctx context.Context, tenantID uuid.UUID, component *Component) error
	GetComponent(ctx context.Context, tenantID, id uuid.UUID) (*Component, error)
	GetComponentBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Component, error)
	ListComponents(ctx context.Context, tenantID uuid.UUID, filter ComponentFilters) ([]Component, int64, error)
	UpdateComponent(ctx context.Context, tenantID uuid.UUID, component *Component) error
	DeleteComponent(ctx context.Context, tenantID, id uuid.UUID) error
	
	// Component instance operations
	CreateComponentInstance(ctx context.Context, tenantID uuid.UUID, instance *ComponentInstance) error
	GetComponentInstance(ctx context.Context, tenantID, id uuid.UUID) (*ComponentInstance, error)
	ListComponentInstances(ctx context.Context, tenantID uuid.UUID, filter ComponentInstanceFilter) ([]*ComponentInstance, error)
	UpdateComponentInstance(ctx context.Context, tenantID uuid.UUID, instance *ComponentInstance) error
	DeleteComponentInstance(ctx context.Context, tenantID, id uuid.UUID) error
	
	// Instance operations (aliases for compatibility)
	CreateInstance(ctx context.Context, tenantID uuid.UUID, instance *ComponentInstance) error
	GetInstance(ctx context.Context, tenantID, id uuid.UUID) (*ComponentInstance, error)
	ListInstances(ctx context.Context, tenantID uuid.UUID, filter ComponentInstanceFilter) ([]*ComponentInstance, error)
	DeleteInstance(ctx context.Context, tenantID, id uuid.UUID) error
	
	// Theme operations
	CreateTheme(ctx context.Context, tenantID uuid.UUID, theme *Theme) error
	GetTheme(ctx context.Context, tenantID, id uuid.UUID) (*Theme, error)
	GetThemeBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Theme, error)
	ListThemes(ctx context.Context, tenantID uuid.UUID, filters ThemeFilters) ([]Theme, int64, error)
	UpdateTheme(ctx context.Context, theme *Theme) error
	DeleteTheme(ctx context.Context, tenantID, id uuid.UUID) error
	GetActiveTheme(ctx context.Context, tenantID uuid.UUID) (*Theme, error)
	
	// Component template operations
	CreateComponentTemplate(ctx context.Context, template *ComponentTemplate) error
	GetComponentTemplate(ctx context.Context, id uuid.UUID) (*ComponentTemplate, error)
	ListComponentTemplates(ctx context.Context, filter ComponentTemplateFilter) ([]*ComponentTemplate, error)
	UpdateComponentTemplate(ctx context.Context, template *ComponentTemplate) error
	DeleteComponentTemplate(ctx context.Context, id uuid.UUID) error
	
	// Theme template operations
	CreateThemeTemplate(ctx context.Context, template *ThemeTemplate) error
	GetThemeTemplate(ctx context.Context, id uuid.UUID) (*ThemeTemplate, error)
	ListThemeTemplates(ctx context.Context, filter ThemeTemplateFilter) ([]*ThemeTemplate, error)
	UpdateThemeTemplate(ctx context.Context, template *ThemeTemplate) error
	DeleteThemeTemplate(ctx context.Context, id uuid.UUID) error
	
	// Template operations (general)
	ListTemplates(ctx context.Context, filters TemplateFilters) ([]ComponentTemplate, int64, error)
	GetTemplate(ctx context.Context, id uuid.UUID) (*ComponentTemplate, error)
	
	// Stats operations
	GetComponentStats(ctx context.Context, tenantID uuid.UUID) (*ComponentStats, error)
	GetStats(ctx context.Context, tenantID uuid.UUID) (*ComponentStats, error)
}



// gormRepository implements Repository using GORM
type gormRepository struct {
	db *gorm.DB
}

// NewRepository creates a new component repository
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// Component operations

func (r *gormRepository) CreateComponent(ctx context.Context, tenantID uuid.UUID, component *Component) error {
	if component.ID == uuid.Nil {
		component.ID = uuid.New()
	}
	component.TenantID = tenantID
	return r.db.WithContext(ctx).Create(component).Error
}

func (r *gormRepository) GetComponent(ctx context.Context, tenantID, id uuid.UUID) (*Component, error) {
	var component Component
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&component).Error
	if err != nil {
		return nil, err
	}
	return &component, nil
}

func (r *gormRepository) GetComponentBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Component, error) {
	var component Component
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND slug = ?", tenantID, slug).
		First(&component).Error
	if err != nil {
		return nil, err
	}
	return &component, nil
}

func (r *gormRepository) UpdateComponent(ctx context.Context, tenantID uuid.UUID, component *Component) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, component.ID).
		Updates(component).Error
}

func (r *gormRepository) DeleteComponent(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&Component{}).Error
}

func (r *gormRepository) ListComponents(ctx context.Context, tenantID uuid.UUID, filters ComponentFilters) ([]Component, int64, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", 
			fmt.Sprintf("%%%s%%", filters.Search),
			fmt.Sprintf("%%%s%%", filters.Search))
	}
	if filters.Featured != nil {
		query = query.Where("is_featured = ?", *filters.Featured)
	}
	
	// Count total
	var total int64
	if err := query.Model(&Component{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination and sorting
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}
	
	offset := (filters.Page - 1) * filters.Limit
	query = query.Offset(offset).Limit(filters.Limit)
	
	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder == "ASC" {
		sortOrder = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	
	var components []Component
	err := query.Find(&components).Error
	return components, total, err
}

// Component Instance operations

func (r *gormRepository) CreateComponentInstance(ctx context.Context, tenantID uuid.UUID, instance *ComponentInstance) error {
	if instance.ID == uuid.Nil {
		instance.ID = uuid.New()
	}
	instance.TenantID = tenantID
	return r.db.WithContext(ctx).Create(instance).Error
}

func (r *gormRepository) GetComponentInstance(ctx context.Context, tenantID, id uuid.UUID) (*ComponentInstance, error) {
	var instance ComponentInstance
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Preload("Component").
		Preload("Theme").
		First(&instance).Error
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (r *gormRepository) UpdateComponentInstance(ctx context.Context, tenantID uuid.UUID, instance *ComponentInstance) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", instance.TenantID, instance.ID).
		Save(instance).Error
}

func (r *gormRepository) DeleteComponentInstance(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&ComponentInstance{}).Error
}

func (r *gormRepository) ListComponentInstances(ctx context.Context, tenantID uuid.UUID, filter ComponentInstanceFilter) ([]*ComponentInstance, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filter.ComponentID != nil {
		query = query.Where("component_id = ?", *filter.ComponentID)
	}
	if filter.ThemeID != nil {
		query = query.Where("theme_id = ?", *filter.ThemeID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Zone != "" {
		query = query.Where("zone = ?", filter.Zone)
	}
	if filter.Visible != nil {
		query = query.Where("visible = ?", *filter.Visible)
	}
	
	// Apply pagination and sorting
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)
	
	// Apply sorting
	sortBy := "position"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "ASC"
	if filter.SortOrder == "DESC" {
		sortOrder = "DESC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	
	var instances []ComponentInstance
	err := query.Preload("Component").Find(&instances).Error
	if err != nil {
		return nil, err
	}
	
	// Convert to pointer slice
	result := make([]*ComponentInstance, len(instances))
	for i := range instances {
		result[i] = &instances[i]
	}
	
	return result, nil
}

// Instance operations (aliases for compatibility)
func (r *gormRepository) CreateInstance(ctx context.Context, tenantID uuid.UUID, instance *ComponentInstance) error {
	return r.CreateComponentInstance(ctx, tenantID, instance)
}

func (r *gormRepository) GetInstance(ctx context.Context, tenantID, id uuid.UUID) (*ComponentInstance, error) {
	return r.GetComponentInstance(ctx, tenantID, id)
}

func (r *gormRepository) ListInstances(ctx context.Context, tenantID uuid.UUID, filter ComponentInstanceFilter) ([]*ComponentInstance, error) {
	return r.ListComponentInstances(ctx, tenantID, filter)
}

func (r *gormRepository) DeleteInstance(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.DeleteComponentInstance(ctx, tenantID, id)
}

// Theme operations

func (r *gormRepository) CreateTheme(ctx context.Context, tenantID uuid.UUID, theme *Theme) error {
	if theme.ID == uuid.Nil {
		theme.ID = uuid.New()
	}
	theme.TenantID = tenantID
	return r.db.WithContext(ctx).Create(theme).Error
}

func (r *gormRepository) GetTheme(ctx context.Context, tenantID, id uuid.UUID) (*Theme, error) {
	var theme Theme
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Preload("Instances").
		Preload("Instances.Component").
		First(&theme).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func (r *gormRepository) GetThemeBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Theme, error) {
	var theme Theme
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND slug = ?", tenantID, slug).
		Preload("Instances").
		Preload("Instances.Component").
		First(&theme).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func (r *gormRepository) UpdateTheme(ctx context.Context, theme *Theme) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", theme.TenantID, theme.ID).
		Save(theme).Error
}

func (r *gormRepository) DeleteTheme(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&Theme{}).Error
}

func (r *gormRepository) ListThemes(ctx context.Context, tenantID uuid.UUID, filters ThemeFilters) ([]Theme, int64, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	// Apply filters
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", 
			fmt.Sprintf("%%%s%%", filters.Search),
			fmt.Sprintf("%%%s%%", filters.Search))
	}
	if filters.Active != nil {
		query = query.Where("is_active = ?", *filters.Active)
	}
	
	// Count total
	var total int64
	if err := query.Model(&Theme{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination and sorting
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}
	
	offset := (filters.Page - 1) * filters.Limit
	query = query.Offset(offset).Limit(filters.Limit)
	
	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder == "ASC" {
		sortOrder = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	
	var themes []Theme
	err := query.Preload("Instances").Find(&themes).Error
	return themes, total, err
}

func (r *gormRepository) GetActiveTheme(ctx context.Context, tenantID uuid.UUID) (*Theme, error) {
	var theme Theme
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Preload("Instances").
		Preload("Instances.Component").
		First(&theme).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

// Template operations

func (r *gormRepository) ListTemplates(ctx context.Context, filters TemplateFilters) ([]ComponentTemplate, int64, error) {
	query := r.db.WithContext(ctx)
	
	// Apply filters
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", 
			fmt.Sprintf("%%%s%%", filters.Search),
			fmt.Sprintf("%%%s%%", filters.Search))
	}
	if filters.Free != nil {
		query = query.Where("is_free = ?", *filters.Free)
	}
	if filters.Featured != nil {
		query = query.Where("is_featured = ?", *filters.Featured)
	}
	
	// Count total
	var total int64
	if err := query.Model(&ComponentTemplate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination and sorting
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}
	
	offset := (filters.Page - 1) * filters.Limit
	query = query.Offset(offset).Limit(filters.Limit)
	
	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder == "ASC" {
		sortOrder = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	
	var templates []ComponentTemplate
	err := query.Find(&templates).Error
	return templates, total, err
}



// Statistics

func (r *gormRepository) GetComponentStats(ctx context.Context, tenantID uuid.UUID) (*ComponentStats, error) {
	return r.GetStats(ctx, tenantID)
}

func (r *gormRepository) GetStats(ctx context.Context, tenantID uuid.UUID) (*ComponentStats, error) {
	stats := &ComponentStats{}
	
	// Component stats
	r.db.WithContext(ctx).Model(&Component{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalComponents)
	r.db.WithContext(ctx).Model(&Component{}).Where("tenant_id = ? AND status = ?", tenantID, StatusActive).Count(&stats.ActiveComponents)
	r.db.WithContext(ctx).Model(&Component{}).Where("tenant_id = ? AND status = ?", tenantID, StatusDraft).Count(&stats.DraftComponents)
	
	// Theme stats
	r.db.WithContext(ctx).Model(&Theme{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalThemes)
	r.db.WithContext(ctx).Model(&Theme{}).Where("tenant_id = ? AND is_active = ?", tenantID, true).Count(&stats.ActiveThemes)
	
	// Instance stats
	r.db.WithContext(ctx).Model(&ComponentInstance{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalInstances)
	
	// Most used components
	r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("usage_count DESC").
		Limit(5).
		Find(&stats.MostUsedComponents)
	
	return stats, nil
}

// Theme Template operations

func (r *gormRepository) CreateThemeTemplate(ctx context.Context, themeTemplate *ThemeTemplate) error {
	if themeTemplate.ID == uuid.Nil {
		themeTemplate.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(themeTemplate).Error
}

func (r *gormRepository) GetThemeTemplate(ctx context.Context, id uuid.UUID) (*ThemeTemplate, error) {
	var themeTemplate ThemeTemplate
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Components").
		First(&themeTemplate).Error
	if err != nil {
		return nil, err
	}
	return &themeTemplate, nil
}

func (r *gormRepository) UpdateThemeTemplate(ctx context.Context, themeTemplate *ThemeTemplate) error {
	return r.db.WithContext(ctx).
		Where("id = ?", themeTemplate.ID).
		Save(themeTemplate).Error
}

func (r *gormRepository) DeleteThemeTemplate(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&ThemeTemplate{}).Error
}

func (r *gormRepository) ListThemeTemplates(ctx context.Context, filters ThemeTemplateFilter) ([]*ThemeTemplate, error) {
	themeTemplates, _, err := r.listThemeTemplatesWithCount(ctx, filters)
	return themeTemplates, err
}

func (r *gormRepository) listThemeTemplatesWithCount(ctx context.Context, filters ThemeTemplateFilter) ([]*ThemeTemplate, int64, error) {
	query := r.db.WithContext(ctx)
	
	// Apply filters
	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", 
			fmt.Sprintf("%%%s%%", filters.Search),
			fmt.Sprintf("%%%s%%", filters.Search))
	}
	if filters.IsFree != nil {
		query = query.Where("is_free = ?", *filters.IsFree)
	}
	if filters.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filters.IsFeatured)
	}
	if filters.IsDefault != nil {
		query = query.Where("is_default = ?", *filters.IsDefault)
	}
	if len(filters.Tags) > 0 {
		query = query.Where("tags && ?", pq.Array(filters.Tags))
	}
	
	// Count total
	var total int64
	if err := query.Model(&ThemeTemplate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}
	
	offset := (filters.Page - 1) * filters.Limit
	query = query.Offset(offset).Limit(filters.Limit)
	
	// Apply sorting
	query = query.Order("created_at DESC")
	
	var themeTemplates []ThemeTemplate
	err := query.Preload("Components").Find(&themeTemplates).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Convert to pointer slice
	result := make([]*ThemeTemplate, len(themeTemplates))
	for i := range themeTemplates {
		result[i] = &themeTemplates[i]
	}
	
	return result, total, nil
}

// Component Template operations

func (r *gormRepository) CreateComponentTemplate(ctx context.Context, componentTemplate *ComponentTemplate) error {
	if componentTemplate.ID == uuid.Nil {
		componentTemplate.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(componentTemplate).Error
}

func (r *gormRepository) GetComponentTemplate(ctx context.Context, id uuid.UUID) (*ComponentTemplate, error) {
	var componentTemplate ComponentTemplate
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&componentTemplate).Error
	if err != nil {
		return nil, err
	}
	return &componentTemplate, nil
}

func (r *gormRepository) UpdateComponentTemplate(ctx context.Context, componentTemplate *ComponentTemplate) error {
	return r.db.WithContext(ctx).
		Where("id = ?", componentTemplate.ID).
		Save(componentTemplate).Error
}

func (r *gormRepository) DeleteComponentTemplate(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&ComponentTemplate{}).Error
}

func (r *gormRepository) ListComponentTemplates(ctx context.Context, filters ComponentTemplateFilter) ([]*ComponentTemplate, error) {
	templates, _, err := r.listComponentTemplatesWithCount(ctx, filters)
	return templates, err
}

func (r *gormRepository) listComponentTemplatesWithCount(ctx context.Context, filters ComponentTemplateFilter) ([]*ComponentTemplate, int64, error) {
	query := r.db.WithContext(ctx)
	
	// Apply filters
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", 
			fmt.Sprintf("%%%s%%", filters.Search),
			fmt.Sprintf("%%%s%%", filters.Search))
	}
	if filters.IsFree != nil {
		query = query.Where("is_free = ?", *filters.IsFree)
	}
	if filters.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filters.IsFeatured)
	}
	if len(filters.Tags) > 0 {
		query = query.Where("tags && ?", pq.Array(filters.Tags))
	}
	
	// Count total
	var total int64
	if err := query.Model(&ComponentTemplate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}
	
	offset := (filters.Page - 1) * filters.Limit
	query = query.Offset(offset).Limit(filters.Limit)
	
	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder == "ASC" {
		sortOrder = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	
	var componentTemplates []ComponentTemplate
	err := query.Find(&componentTemplates).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Convert to pointer slice
	result := make([]*ComponentTemplate, len(componentTemplates))
	for i := range componentTemplates {
		result[i] = &componentTemplates[i]
	}
	
	return result, total, err
}

// Template operations (general aliases)
func (r *gormRepository) GetTemplate(ctx context.Context, id uuid.UUID) (*ComponentTemplate, error) {
	return r.GetComponentTemplate(ctx, id)
}