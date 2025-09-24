package content

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sharedErrors "ecommerce-saas/internal/shared/errors"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Page Repository Methods

func (r *Repository) CreatePage(page *Page) (*Page, error) {
	if err := r.db.Create(page).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to create page", err)
	}
	return page, nil
}

func (r *Repository) UpdatePage(page *Page) (*Page, error) {
	if err := r.db.Save(page).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to update page", err)
	}
	return page, nil
}

func (r *Repository) DeletePage(tenantID, pageID uuid.UUID) error {
	if err := r.db.Where("tenant_id = ? AND id = ?", tenantID, pageID).Delete(&Page{}).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete page", err)
	}
	return nil
}

func (r *Repository) GetPage(tenantID, pageID uuid.UUID) (*Page, error) {
	var page Page
	err := r.db.Preload("FeaturedImage").Preload("Tags").Preload("Categories").Preload("Parent").Preload("Children").
		Where("tenant_id = ? AND id = ?", tenantID, pageID).
		First(&page).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("page not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get page", err)
	}
	return &page, nil
}

func (r *Repository) GetPageBySlug(tenantID uuid.UUID, slug string) (*Page, error) {
	var page Page
	err := r.db.Preload("FeaturedImage").Preload("Tags").Preload("Categories").Preload("Parent").Preload("Children").
		Where("tenant_id = ? AND slug = ?", tenantID, slug).
		First(&page).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("page not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get page by slug", err)
	}
	return &page, nil
}

func (r *Repository) GetPages(tenantID uuid.UUID, filter PageListFilter) ([]Page, int64, error) {
	var pages []Page
	var total int64

	query := r.db.Model(&Page{}).Where("tenant_id = ?", tenantID)

	// Apply filters
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ? OR excerpt ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, sharedErrors.NewInternalError("failed to count pages", err)
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset)
	}

	// Get results with preloads
	err := query.Preload("FeaturedImage").Preload("Tags").Preload("Categories").
		Order("created_at DESC").
		Find(&pages).Error
	if err != nil {
		return nil, 0, sharedErrors.NewInternalError("failed to get pages", err)
	}

	return pages, total, nil
}

func (r *Repository) GetPublishedPages(tenantID uuid.UUID) ([]Page, error) {
	var pages []Page
	err := r.db.Where("tenant_id = ? AND status = ? AND (published_at IS NULL OR published_at <= NOW())",
		tenantID, StatusPublished).
		Order("published_at DESC").
		Find(&pages).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get published pages", err)
	}
	return pages, nil
}

func (r *Repository) GetPopularPages(tenantID uuid.UUID, limit int) ([]Page, error) {
	var pages []Page
	err := r.db.Where("tenant_id = ? AND status = ?", tenantID, StatusPublished).
		Order("view_count DESC").
		Limit(limit).
		Find(&pages).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get popular pages", err)
	}
	return pages, nil
}

func (r *Repository) GetRecentPages(tenantID uuid.UUID, limit int) ([]Page, error) {
	var pages []Page
	err := r.db.Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Find(&pages).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get recent pages", err)
	}
	return pages, nil
}

func (r *Repository) IncrementPageViews(tenantID, pageID uuid.UUID) error {
	if err := r.db.Model(&Page{}).
		Where("tenant_id = ? AND id = ?", tenantID, pageID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
		return sharedErrors.NewInternalError("failed to increment page views", err)
	}
	return nil
}

func (r *Repository) SlugExists(tenantID uuid.UUID, slug string, excludeID uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&Page{}).Where("tenant_id = ? AND slug = ?", tenantID, slug)

	if excludeID != uuid.Nil {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, sharedErrors.NewInternalError("failed to check slug existence", err)
	}
	return count > 0, nil
}

func (r *Repository) AssociatePageTags(pageID uuid.UUID, tagIDs []uuid.UUID) error {
	// First, clear existing associations
	if err := r.db.Exec("DELETE FROM page_tags WHERE page_id = ? AND EXISTS (SELECT 1 FROM pages WHERE pages.id = page_tags.page_id)", pageID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to clear page tag associations", err)
	}

	// Add new associations
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			if err := r.db.Exec("INSERT INTO page_tags (page_id, tag_id) VALUES (?, ?)", pageID, tagID).Error; err != nil {
				return sharedErrors.NewInternalError("failed to associate page tag", err)
			}
		}
	}
	return nil
}

func (r *Repository) UpdatePageTags(pageID uuid.UUID, tagIDs []uuid.UUID) error {
	return r.AssociatePageTags(pageID, tagIDs)
}

func (r *Repository) AssociatePageCategories(pageID uuid.UUID, categoryIDs []uuid.UUID) error {
	// First, clear existing associations
	if err := r.db.Exec("DELETE FROM page_categories WHERE page_id = ? AND EXISTS (SELECT 1 FROM pages WHERE pages.id = page_categories.page_id)", pageID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to clear page category associations", err)
	}

	// Add new associations
	if len(categoryIDs) > 0 {
		for _, categoryID := range categoryIDs {
			if err := r.db.Exec("INSERT INTO page_categories (page_id, category_id) VALUES (?, ?)", pageID, categoryID).Error; err != nil {
				return sharedErrors.NewInternalError("failed to associate page category", err)
			}
		}
	}
	return nil
}

func (r *Repository) UpdatePageCategories(pageID uuid.UUID, categoryIDs []uuid.UUID) error {
	return r.AssociatePageCategories(pageID, categoryIDs)
}

// Media Repository Methods

func (r *Repository) CreateMedia(media *Media) (*Media, error) {
	if err := r.db.Create(media).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to create media", err)
	}
	return media, nil
}

func (r *Repository) UpdateMedia(media *Media) (*Media, error) {
	if err := r.db.Save(media).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to update media", err)
	}
	return media, nil
}

func (r *Repository) DeleteMedia(tenantID, mediaID uuid.UUID) error {
	if err := r.db.Where("tenant_id = ? AND id = ?", tenantID, mediaID).Delete(&Media{}).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete media", err)
	}
	return nil
}

func (r *Repository) GetMedia(tenantID, mediaID uuid.UUID) (*Media, error) {
	var media Media
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, mediaID).First(&media).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("media not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get media", err)
	}
	return &media, nil
}

func (r *Repository) GetMediaLibrary(tenantID uuid.UUID, filter MediaLibraryFilter) ([]Media, int64, error) {
	var media []Media
	var total int64

	query := r.db.Model(&Media{}).Where("tenant_id = ?", tenantID)

	// Apply filters
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ? OR file_name ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, sharedErrors.NewInternalError("failed to count media", err)
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset)
	}

	// Get results
	err := query.Order("created_at DESC").Find(&media).Error
	if err != nil {
		return nil, 0, sharedErrors.NewInternalError("failed to get media library", err)
	}

	return media, total, nil
}

// Menu Repository Methods

func (r *Repository) CreateMenu(menu *Menu) (*Menu, error) {
	if err := r.db.Create(menu).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to create menu", err)
	}
	return menu, nil
}

func (r *Repository) UpdateMenu(menu *Menu) (*Menu, error) {
	if err := r.db.Save(menu).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to update menu", err)
	}
	return menu, nil
}

func (r *Repository) DeleteMenu(tenantID, menuID uuid.UUID) error {
	// Delete menu items first
	if err := r.db.Where("menu_id = ?", menuID).Delete(&MenuItem{}).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete menu items", err)
	}

	// Delete menu
	if err := r.db.Where("tenant_id = ? AND id = ?", tenantID, menuID).Delete(&Menu{}).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete menu", err)
	}
	return nil
}

func (r *Repository) GetMenu(tenantID, menuID uuid.UUID) (*Menu, error) {
	var menu Menu
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = ?", true).Order("sort_order ASC")
	}).Preload("Items.Page").Preload("Items.Children", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = ?", true).Order("sort_order ASC")
	}).
		Where("tenant_id = ? AND id = ?", tenantID, menuID).
		First(&menu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("menu not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get menu", err)
	}
	return &menu, nil
}

func (r *Repository) GetMenus(tenantID uuid.UUID, location string, activeOnly bool) ([]Menu, error) {
	var menus []Menu
	query := r.db.Where("tenant_id = ?", tenantID)

	if location != "" {
		query = query.Where("location = ?", location)
	}
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("name ASC").Find(&menus).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get menus", err)
	}
	return menus, nil
}

func (r *Repository) GetMenuByLocation(tenantID uuid.UUID, location string) (*Menu, error) {
	var menu Menu
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = ? AND parent_id IS NULL", true).Order("sort_order ASC")
	}).Preload("Items.Page").Preload("Items.Children", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = ?", true).Order("sort_order ASC")
	}).
		Where("tenant_id = ? AND location = ? AND is_active = ?", tenantID, location, true).
		First(&menu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("menu not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get menu by location", err)
	}
	return &menu, nil
}

// Menu Item Repository Methods

func (r *Repository) CreateMenuItem(menuItem *MenuItem) (*MenuItem, error) {
	if err := r.db.Create(menuItem).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to create menu item", err)
	}
	return menuItem, nil
}

func (r *Repository) UpdateMenuItem(menuItem *MenuItem) (*MenuItem, error) {
	if err := r.db.Save(menuItem).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to update menu item", err)
	}
	return menuItem, nil
}

func (r *Repository) DeleteMenuItem(tenantID, itemID uuid.UUID) error {
	// Delete children first
	if err := r.db.Where("parent_id = ?", itemID).Delete(&MenuItem{}).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete menu item children", err)
	}

	// Delete the item
	if err := r.db.Where("id = ?", itemID).Delete(&MenuItem{}).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete menu item", err)
	}
	return nil
}

func (r *Repository) GetMenuItem(tenantID, itemID uuid.UUID) (*MenuItem, error) {
	var menuItem MenuItem
	// Join with menu to check tenant
	err := r.db.Joins("JOIN menus ON menu_items.menu_id = menus.id").
		Where("menus.tenant_id = ? AND menu_items.id = ?", tenantID, itemID).
		First(&menuItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("menu item not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get menu item", err)
	}
	return &menuItem, nil
}

func (r *Repository) ReorderMenuItems(tenantID, menuID uuid.UUID, itemOrders []MenuItemOrder) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, order := range itemOrders {
			if err := tx.Model(&MenuItem{}).
				Joins("JOIN menus ON menu_items.menu_id = menus.id").
				Where("menus.tenant_id = ? AND menu_items.menu_id = ? AND menu_items.id = ?", tenantID, menuID, order.ID).
				Update("sort_order", order.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return sharedErrors.NewInternalError("failed to reorder menu items", err)
	}
	return nil
}

// Tag Repository Methods

func (r *Repository) CreateTag(tag *Tag) (*Tag, error) {
	if err := r.db.Create(tag).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to create tag", err)
	}
	return tag, nil
}

func (r *Repository) UpdateTag(tag *Tag) (*Tag, error) {
	if err := r.db.Save(tag).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to update tag", err)
	}
	return tag, nil
}

func (r *Repository) DeleteTag(tenantID, tagID uuid.UUID) error {
	// Remove associations first
	if err := r.db.Exec("DELETE FROM page_tags WHERE tag_id = ? AND EXISTS (SELECT 1 FROM tags WHERE tags.id = page_tags.tag_id AND tags.tenant_id = ?)", tagID, tenantID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to remove tag associations", err)
	}

	// Delete tag
	if err := r.db.Where("tenant_id = ? AND id = ?", tenantID, tagID).Delete(&Tag{}).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete tag", err)
	}
	return nil
}

func (r *Repository) GetTag(tenantID, tagID uuid.UUID) (*Tag, error) {
	var tag Tag
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, tagID).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("tag not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get tag", err)
	}
	return &tag, nil
}

func (r *Repository) GetTags(tenantID uuid.UUID, search string) ([]Tag, error) {
	var tags []Tag
	query := r.db.Where("tenant_id = ?", tenantID)

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	err := query.Order("name ASC").Find(&tags).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get tags", err)
	}
	return tags, nil
}

// Category Repository Methods

func (r *Repository) CreateCategory(category *Category) (*Category, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to create category", err)
	}
	return category, nil
}

func (r *Repository) UpdateCategory(category *Category) (*Category, error) {
	if err := r.db.Save(category).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to update category", err)
	}
	return category, nil
}

func (r *Repository) DeleteCategory(tenantID, categoryID uuid.UUID) error {
	// Remove associations first
	if err := r.db.Exec("DELETE FROM page_categories WHERE category_id = ? AND EXISTS (SELECT 1 FROM categories WHERE categories.id = page_categories.category_id AND categories.tenant_id = ?)", categoryID, tenantID).Error; err != nil {
		return sharedErrors.NewInternalError("failed to remove category associations", err)
	}

	// Update children to remove parent reference
	if err := r.db.Model(&Category{}).
		Where("tenant_id = ? AND parent_id = ?", tenantID, categoryID).
		Update("parent_id", nil).Error; err != nil {
		return sharedErrors.NewInternalError("failed to update child categories", err)
	}

	// Delete category
	if err := r.db.Where("tenant_id = ? AND id = ?", tenantID, categoryID).Delete(&Category{}).Error; err != nil {
		return sharedErrors.NewInternalError("failed to delete category", err)
	}
	return nil
}

func (r *Repository) GetCategory(tenantID, categoryID uuid.UUID) (*Category, error) {
	var category Category
	err := r.db.Preload("Parent").Preload("Children").
		Where("tenant_id = ? AND id = ?", tenantID, categoryID).
		First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("category not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get category", err)
	}
	return &category, nil
}

func (r *Repository) GetCategories(tenantID uuid.UUID, activeOnly bool) ([]Category, error) {
	var categories []Category
	query := r.db.Preload("Parent").Preload("Children").Where("tenant_id = ?", tenantID)

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("sort_order ASC, name ASC").Find(&categories).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get categories", err)
	}
	return categories, nil
}

// SEO Repository Methods

func (r *Repository) CreateSEOSettings(settings *SEOSettings) (*SEOSettings, error) {
	if err := r.db.Create(settings).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to create SEO settings", err)
	}
	return settings, nil
}

func (r *Repository) UpdateSEOSettings(settings *SEOSettings) (*SEOSettings, error) {
	if err := r.db.Save(settings).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to update SEO settings", err)
	}
	return settings, nil
}

func (r *Repository) GetSEOSettings(tenantID uuid.UUID) (*SEOSettings, error) {
	var settings SEOSettings
	err := r.db.Where("tenant_id = ?", tenantID).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedErrors.NewNotFoundError("SEO settings not found")
		}
		return nil, sharedErrors.NewInternalError("failed to get SEO settings", err)
	}
	return &settings, nil
}

// Analytics Repository Methods

type ContentStats struct {
	TotalPages     int   `json:"total_pages"`
	PublishedPages int   `json:"published_pages"`
	DraftPages     int   `json:"draft_pages"`
	TotalMedia     int   `json:"total_media"`
	TotalViews     int64 `json:"total_views"`
}

func (r *Repository) GetContentStats(tenantID uuid.UUID) (*ContentStats, error) {
	stats := &ContentStats{}

	// Get page counts
	var totalPages, publishedPages, draftPages, totalMedia int64
	if err := r.db.Model(&Page{}).Where("tenant_id = ?", tenantID).Count(&totalPages).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to count pages", err)
	}
	if err := r.db.Model(&Page{}).Where("tenant_id = ? AND status = ?", tenantID, StatusPublished).Count(&publishedPages).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to count published pages", err)
	}
	if err := r.db.Model(&Page{}).Where("tenant_id = ? AND status = ?", tenantID, StatusDraft).Count(&draftPages).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to count draft pages", err)
	}

	// Get media count
	if err := r.db.Model(&Media{}).Where("tenant_id = ?", tenantID).Count(&totalMedia).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to count media", err)
	}

	// Get total views
	var totalViewsResult struct {
		TotalViews int64 `json:"total_views"`
	}
	if err := r.db.Model(&Page{}).
		Where("tenant_id = ?", tenantID).
		Select("SUM(view_count) as total_views").
		Scan(&totalViewsResult).Error; err != nil {
		return nil, sharedErrors.NewInternalError("failed to calculate total views", err)
	}

	stats.TotalPages = int(totalPages)
	stats.PublishedPages = int(publishedPages)
	stats.DraftPages = int(draftPages)
	stats.TotalMedia = int(totalMedia)
	stats.TotalViews = totalViewsResult.TotalViews

	return stats, nil
}

// Search Repository Methods

func (r *Repository) SearchContent(tenantID uuid.UUID, filter ContentSearchFilter) ([]ContentSearchResult, int64, error) {
	var results []ContentSearchResult
	var total int64

	query := r.db.Model(&Page{}).Where("tenant_id = ?", tenantID)

	// Apply search
	if filter.Query != "" {
		searchTerm := "%" + strings.ToLower(filter.Query) + "%"
		query = query.Where(
			"LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(excerpt) LIKE ? OR LOWER(meta_description) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// Apply type filter
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	// Get total count
	query.Count(&total)

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset)
	}

	// Get results
	var pages []Page
	err := query.Select("id, type, title, slug, excerpt, status, published_at").
		Order("view_count DESC, created_at DESC").
		Find(&pages).Error
	if err != nil {
		return nil, 0, sharedErrors.NewInternalError("failed to search content", err)
	}

	// Convert to search results
	results = make([]ContentSearchResult, len(pages))
	for i, page := range pages {
		relevance := 1.0
		if filter.Query != "" {
			// Simple relevance scoring based on title match
			if strings.Contains(strings.ToLower(page.Title), strings.ToLower(filter.Query)) {
				relevance = 1.0
			} else if strings.Contains(strings.ToLower(page.Excerpt), strings.ToLower(filter.Query)) {
				relevance = 0.8
			} else {
				relevance = 0.6
			}
		}

		results[i] = ContentSearchResult{
			ID:          page.ID,
			Type:        page.Type,
			Title:       page.Title,
			Slug:        page.Slug,
			Excerpt:     page.Excerpt,
			Status:      page.Status,
			PublishedAt: page.PublishedAt,
			Relevance:   relevance,
		}
	}

	return results, total, nil
}

// Additional utility methods

func (r *Repository) GetPagesByTag(tenantID uuid.UUID, tagID uuid.UUID, limit int) ([]Page, error) {
	var pages []Page
	err := r.db.Joins("JOIN page_tags ON pages.id = page_tags.page_id").
		Where("pages.tenant_id = ? AND page_tags.tag_id = ? AND pages.status = ?", tenantID, tagID, StatusPublished).
		Limit(limit).
		Order("pages.published_at DESC").
		Find(&pages).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get pages by tag", err)
	}
	return pages, nil
}

func (r *Repository) GetPagesByCategory(tenantID uuid.UUID, categoryID uuid.UUID, limit int) ([]Page, error) {
	var pages []Page
	err := r.db.Joins("JOIN page_categories ON pages.id = page_categories.page_id").
		Where("pages.tenant_id = ? AND page_categories.category_id = ? AND pages.status = ?", tenantID, categoryID, StatusPublished).
		Limit(limit).
		Order("pages.published_at DESC").
		Find(&pages).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get pages by category", err)
	}
	return pages, nil
}

func (r *Repository) GetRelatedPages(tenantID uuid.UUID, pageID uuid.UUID, limit int) ([]Page, error) {
	var pages []Page

	// Get pages that share tags or categories with the current page
	err := r.db.Raw(`
		SELECT DISTINCT p.* FROM pages p
		WHERE p.tenant_id = ? AND p.id != ? AND p.status = 'published'
		AND (
			p.id IN (
				SELECT DISTINCT pt2.page_id FROM page_tags pt1
				JOIN page_tags pt2 ON pt1.tag_id = pt2.tag_id
				WHERE pt1.page_id = ? AND pt2.page_id != ?
			)
			OR p.id IN (
				SELECT DISTINCT pc2.page_id FROM page_categories pc1
				JOIN page_categories pc2 ON pc1.category_id = pc2.category_id
				WHERE pc1.page_id = ? AND pc2.page_id != ?
			)
		)
		ORDER BY p.view_count DESC, p.published_at DESC
		LIMIT ?
	`, tenantID, pageID, pageID, pageID, pageID, pageID, limit).Scan(&pages).Error

	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get related pages", err)
	}
	return pages, nil
}

func (r *Repository) GetMediaUsage(tenantID uuid.UUID, mediaID uuid.UUID) ([]Page, error) {
	var pages []Page
	err := r.db.Where("tenant_id = ? AND featured_image_id = ?", tenantID, mediaID).
		Or("tenant_id = ? AND content LIKE ?", tenantID, fmt.Sprintf("%%%s%%", mediaID.String())).
		Find(&pages).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get media usage", err)
	}
	return pages, nil
}

func (r *Repository) GetMenuItemsByPage(tenantID uuid.UUID, pageID uuid.UUID) ([]MenuItem, error) {
	var menuItems []MenuItem
	err := r.db.Joins("JOIN menus ON menu_items.menu_id = menus.id").
		Where("menus.tenant_id = ? AND menu_items.page_id = ?", tenantID, pageID).
		Find(&menuItems).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get menu items by page", err)
	}
	return menuItems, nil
}

func (r *Repository) GetOrphanedMedia(tenantID uuid.UUID) ([]Media, error) {
	var media []Media
	err := r.db.Where(`tenant_id = ? AND id NOT IN (
		SELECT DISTINCT featured_image_id FROM pages 
		WHERE tenant_id = ? AND featured_image_id IS NOT NULL
	)`, tenantID, tenantID).Find(&media).Error
	if err != nil {
		return nil, sharedErrors.NewInternalError("failed to get orphaned media", err)
	}
	return media, nil
}

func (r *Repository) BulkUpdatePageStatus(tenantID uuid.UUID, pageIDs []uuid.UUID, status PageStatus) error {
	err := r.db.Model(&Page{}).
		Where("tenant_id = ? AND id IN ?", tenantID, pageIDs).
		Update("status", status).Error
	if err != nil {
		return sharedErrors.NewInternalError("failed to bulk update page status", err)
	}
	return nil
}

func (r *Repository) BulkDeletePages(tenantID uuid.UUID, pageIDs []uuid.UUID) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Delete tag associations
		if err := tx.Exec("DELETE FROM page_tags WHERE page_id IN ? AND EXISTS (SELECT 1 FROM pages WHERE pages.id = page_tags.page_id AND pages.tenant_id = ?)", pageIDs, tenantID).Error; err != nil {
			return err
		}

		// Delete category associations
		if err := tx.Exec("DELETE FROM page_categories WHERE page_id IN ? AND EXISTS (SELECT 1 FROM pages WHERE pages.id = page_categories.page_id AND pages.tenant_id = ?)", pageIDs, tenantID).Error; err != nil {
			return err
		}

		// Delete pages
		if err := tx.Where("tenant_id = ? AND id IN ?", tenantID, pageIDs).Delete(&Page{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return sharedErrors.NewInternalError("failed to bulk delete pages", err)
	}
	return nil
}

func (r *Repository) BulkDeleteMedia(tenantID uuid.UUID, mediaIDs []uuid.UUID) error {
	err := r.db.Where("tenant_id = ? AND id IN ?", tenantID, mediaIDs).Delete(&Media{}).Error
	if err != nil {
		return sharedErrors.NewInternalError("failed to bulk delete media", err)
	}
	return nil
}
