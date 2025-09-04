package content

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// Request/Response DTOs

type CreatePageRequest struct {
	Type            ContentType `json:"type" binding:"required"`
	Title           string      `json:"title" binding:"required"`
	Slug            string      `json:"slug"`
	Content         string      `json:"content"`
	Excerpt         string      `json:"excerpt"`
	Status          PageStatus  `json:"status"`
	MetaTitle       string      `json:"meta_title"`
	MetaDescription string      `json:"meta_description"`
	MetaKeywords    string      `json:"meta_keywords"`
	OGTitle         string      `json:"og_title"`
	OGDescription   string      `json:"og_description"`
	OGImage         string      `json:"og_image"`
	FeaturedImageID *uuid.UUID  `json:"featured_image_id"`
	ParentID        *uuid.UUID  `json:"parent_id"`
	Template        string      `json:"template"`
	LayoutData      string      `json:"layout_data"`
	TagIDs          []uuid.UUID `json:"tag_ids"`
	CategoryIDs     []uuid.UUID `json:"category_ids"`
	PublishedAt     *time.Time  `json:"published_at"`
	ScheduledAt     *time.Time  `json:"scheduled_at"`
}

type UpdatePageRequest struct {
	Title           *string     `json:"title"`
	Slug            *string     `json:"slug"`
	Content         *string     `json:"content"`
	Excerpt         *string     `json:"excerpt"`
	Status          *PageStatus `json:"status"`
	MetaTitle       *string     `json:"meta_title"`
	MetaDescription *string     `json:"meta_description"`
	MetaKeywords    *string     `json:"meta_keywords"`
	OGTitle         *string     `json:"og_title"`
	OGDescription   *string     `json:"og_description"`
	OGImage         *string     `json:"og_image"`
	FeaturedImageID *uuid.UUID  `json:"featured_image_id"`
	ParentID        *uuid.UUID  `json:"parent_id"`
	Template        *string     `json:"template"`
	LayoutData      *string     `json:"layout_data"`
	TagIDs          []uuid.UUID `json:"tag_ids"`
	CategoryIDs     []uuid.UUID `json:"category_ids"`
	PublishedAt     *time.Time  `json:"published_at"`
	ScheduledAt     *time.Time  `json:"scheduled_at"`
}

type DuplicatePageRequest struct {
	Title  string     `json:"title" binding:"required"`
	Slug   string     `json:"slug"`
	Status PageStatus `json:"status"`
}

type PageListFilter struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Search string `json:"search"`
}

type UploadMediaRequest struct {
	File        multipart.File   `json:"-"`
	Header      *multipart.FileHeader `json:"-"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	AltText     string           `json:"alt_text"`
}

type UpdateMediaRequest struct {
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	AltText         *string `json:"alt_text"`
	MetaTitle       *string `json:"meta_title"`
	MetaDescription *string `json:"meta_description"`
}

type MediaLibraryFilter struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Type   string `json:"type"`
	Search string `json:"search"`
}

type CreateMenuRequest struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type UpdateMenuRequest struct {
	Name     *string `json:"name"`
	Location *string `json:"location"`
	IsActive *bool   `json:"is_active"`
}

type CreateMenuItemRequest struct {
	ParentID  *uuid.UUID `json:"parent_id"`
	Title     string     `json:"title" binding:"required"`
	URL       string     `json:"url"`
	PageID    *uuid.UUID `json:"page_id"`
	Target    string     `json:"target"`
	CSSClass  string     `json:"css_class"`
	IconClass string     `json:"icon_class"`
	IsActive  bool       `json:"is_active"`
	SortOrder int        `json:"sort_order"`
}

type UpdateMenuItemRequest struct {
	ParentID  *uuid.UUID `json:"parent_id"`
	Title     *string    `json:"title"`
	URL       *string    `json:"url"`
	PageID    *uuid.UUID `json:"page_id"`
	Target    *string    `json:"target"`
	CSSClass  *string    `json:"css_class"`
	IconClass *string    `json:"icon_class"`
	IsActive  *bool      `json:"is_active"`
	SortOrder *int       `json:"sort_order"`
}

type ReorderMenuItemsRequest struct {
	ItemOrders []MenuItemOrder `json:"item_orders" binding:"required"`
}

type MenuItemOrder struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	SortOrder int       `json:"sort_order"`
}

type CreateTagRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateTagRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}

type CreateCategoryRequest struct {
	Name        string     `json:"name" binding:"required"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsActive    bool       `json:"is_active"`
	SortOrder   int        `json:"sort_order"`
}

type UpdateCategoryRequest struct {
	Name        *string    `json:"name"`
	Slug        *string    `json:"slug"`
	Description *string    `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsActive    *bool      `json:"is_active"`
	SortOrder   *int       `json:"sort_order"`
}

type UpdateSEOSettingsRequest struct {
	SiteTitle            *string `json:"site_title"`
	SiteDescription      *string `json:"site_description"`
	DefaultMetaTitle     *string `json:"default_meta_title"`
	DefaultMetaDesc      *string `json:"default_meta_description"`
	DefaultOGImage       *string `json:"default_og_image"`
	GoogleAnalyticsID    *string `json:"google_analytics_id"`
	GoogleTagManagerID   *string `json:"google_tag_manager_id"`
	FacebookPixelID      *string `json:"facebook_pixel_id"`
	GoogleVerification   *string `json:"google_verification"`
	BingVerification     *string `json:"bing_verification"`
	RobotsTxt            *string `json:"robots_txt"`
	EnableSitemap        *bool   `json:"enable_sitemap"`
	SitemapFrequency     *string `json:"sitemap_frequency"`
}

type ContentAnalytics struct {
	TotalPages      int                   `json:"total_pages"`
	PublishedPages  int                   `json:"published_pages"`
	DraftPages      int                   `json:"draft_pages"`
	TotalMedia      int                   `json:"total_media"`
	TotalViews      int64                 `json:"total_views"`
	PopularPages    []PageAnalytics       `json:"popular_pages"`
	RecentPages     []Page                `json:"recent_pages"`
	ContentByType   []ContentTypeStats    `json:"content_by_type"`
	ViewsByMonth    []MonthlyViewStats    `json:"views_by_month"`
}

type PageAnalytics struct {
	Page      Page  `json:"page"`
	Views     int64 `json:"views"`
	ViewsRate float64 `json:"views_rate"`
}

type ContentTypeStats struct {
	Type  ContentType `json:"type"`
	Count int         `json:"count"`
}

type MonthlyViewStats struct {
	Month string `json:"month"`
	Views int64  `json:"views"`
}

type ContentSearchFilter struct {
	Query  string `json:"query"`
	Type   string `json:"type"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type ContentSearchResult struct {
	ID          uuid.UUID   `json:"id"`
	Type        ContentType `json:"type"`
	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	Excerpt     string      `json:"excerpt"`
	Status      PageStatus  `json:"status"`
	PublishedAt *time.Time  `json:"published_at"`
	Relevance   float64     `json:"relevance"`
}

// Page Management Service Methods

func (s *Service) CreatePage(tenantID, authorID uuid.UUID, req CreatePageRequest) (*Page, error) {
	page := &Page{
		TenantID:        tenantID,
		Type:            req.Type,
		Title:           req.Title,
		Slug:            req.Slug,
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		Status:          req.Status,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		MetaKeywords:    req.MetaKeywords,
		OGTitle:         req.OGTitle,
		OGDescription:   req.OGDescription,
		OGImage:         req.OGImage,
		FeaturedImageID: req.FeaturedImageID,
		AuthorID:        authorID,
		ParentID:        req.ParentID,
		Template:        req.Template,
		LayoutData:      req.LayoutData,
		PublishedAt:     req.PublishedAt,
		ScheduledAt:     req.ScheduledAt,
	}

	// Generate slug if not provided
	if page.Slug == "" {
		page.Slug = s.generateSlug(page.Title)
	}

	// Set default status
	if page.Status == "" {
		page.Status = StatusDraft
	}

	// Set default template
	if page.Template == "" {
		page.Template = "default"
	}

	// Validate slug uniqueness
	if err := s.validateSlugUniqueness(tenantID, page.Slug, uuid.Nil); err != nil {
		return nil, err
	}

	// Create the page
	createdPage, err := s.repository.CreatePage(page)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	// Associate tags and categories
	if len(req.TagIDs) > 0 {
		if err := s.repository.AssociatePageTags(createdPage.ID, req.TagIDs); err != nil {
			return nil, fmt.Errorf("failed to associate tags: %w", err)
		}
	}

	if len(req.CategoryIDs) > 0 {
		if err := s.repository.AssociatePageCategories(createdPage.ID, req.CategoryIDs); err != nil {
			return nil, fmt.Errorf("failed to associate categories: %w", err)
		}
	}

	return s.repository.GetPage(tenantID, createdPage.ID)
}

func (s *Service) UpdatePage(tenantID, pageID uuid.UUID, req UpdatePageRequest) (*Page, error) {
	page, err := s.repository.GetPage(tenantID, pageID)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	// Update fields
	if req.Title != nil {
		page.Title = *req.Title
	}
	if req.Slug != nil {
		if err := s.validateSlugUniqueness(tenantID, *req.Slug, pageID); err != nil {
			return nil, err
		}
		page.Slug = *req.Slug
	}
	if req.Content != nil {
		page.Content = *req.Content
	}
	if req.Excerpt != nil {
		page.Excerpt = *req.Excerpt
	}
	if req.Status != nil {
		page.Status = *req.Status
	}
	if req.MetaTitle != nil {
		page.MetaTitle = *req.MetaTitle
	}
	if req.MetaDescription != nil {
		page.MetaDescription = *req.MetaDescription
	}
	if req.MetaKeywords != nil {
		page.MetaKeywords = *req.MetaKeywords
	}
	if req.OGTitle != nil {
		page.OGTitle = *req.OGTitle
	}
	if req.OGDescription != nil {
		page.OGDescription = *req.OGDescription
	}
	if req.OGImage != nil {
		page.OGImage = *req.OGImage
	}
	if req.FeaturedImageID != nil {
		page.FeaturedImageID = req.FeaturedImageID
	}
	if req.ParentID != nil {
		page.ParentID = req.ParentID
	}
	if req.Template != nil {
		page.Template = *req.Template
	}
	if req.LayoutData != nil {
		page.LayoutData = *req.LayoutData
	}
	if req.PublishedAt != nil {
		page.PublishedAt = req.PublishedAt
	}
	if req.ScheduledAt != nil {
		page.ScheduledAt = req.ScheduledAt
	}

	// Update tags and categories
	if req.TagIDs != nil {
		if err := s.repository.UpdatePageTags(pageID, req.TagIDs); err != nil {
			return nil, fmt.Errorf("failed to update tags: %w", err)
		}
	}

	if req.CategoryIDs != nil {
		if err := s.repository.UpdatePageCategories(pageID, req.CategoryIDs); err != nil {
			return nil, fmt.Errorf("failed to update categories: %w", err)
		}
	}

	return s.repository.UpdatePage(page)
}

func (s *Service) DeletePage(tenantID, pageID uuid.UUID) error {
	return s.repository.DeletePage(tenantID, pageID)
}

func (s *Service) GetPage(tenantID, pageID uuid.UUID) (*Page, error) {
	return s.repository.GetPage(tenantID, pageID)
}

func (s *Service) GetPageBySlug(tenantID uuid.UUID, slug string) (*Page, error) {
	return s.repository.GetPageBySlug(tenantID, slug)
}

func (s *Service) GetPublishedPageBySlug(tenantID uuid.UUID, slug string) (*Page, error) {
	page, err := s.repository.GetPageBySlug(tenantID, slug)
	if err != nil {
		return nil, err
	}

	if !page.IsPublished() {
		return nil, errors.New("page not published")
	}

	return page, nil
}

func (s *Service) GetPages(tenantID uuid.UUID, filter PageListFilter) ([]Page, int64, error) {
	return s.repository.GetPages(tenantID, filter)
}

func (s *Service) PublishPage(tenantID, pageID uuid.UUID) (*Page, error) {
	page, err := s.repository.GetPage(tenantID, pageID)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	page.Status = StatusPublished
	now := time.Now()
	page.PublishedAt = &now

	return s.repository.UpdatePage(page)
}

func (s *Service) UnpublishPage(tenantID, pageID uuid.UUID) (*Page, error) {
	page, err := s.repository.GetPage(tenantID, pageID)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	page.Status = StatusDraft
	page.PublishedAt = nil

	return s.repository.UpdatePage(page)
}

func (s *Service) DuplicatePage(tenantID, authorID, pageID uuid.UUID, req DuplicatePageRequest) (*Page, error) {
	originalPage, err := s.repository.GetPage(tenantID, pageID)
	if err != nil {
		return nil, fmt.Errorf("original page not found: %w", err)
	}

	duplicatePage := &Page{
		TenantID:        tenantID,
		Type:            originalPage.Type,
		Title:           req.Title,
		Slug:            req.Slug,
		Content:         originalPage.Content,
		Excerpt:         originalPage.Excerpt,
		Status:          req.Status,
		MetaTitle:       originalPage.MetaTitle,
		MetaDescription: originalPage.MetaDescription,
		MetaKeywords:    originalPage.MetaKeywords,
		OGTitle:         originalPage.OGTitle,
		OGDescription:   originalPage.OGDescription,
		OGImage:         originalPage.OGImage,
		FeaturedImageID: originalPage.FeaturedImageID,
		AuthorID:        authorID,
		ParentID:        originalPage.ParentID,
		Template:        originalPage.Template,
		LayoutData:      originalPage.LayoutData,
	}

	// Generate slug if not provided
	if duplicatePage.Slug == "" {
		duplicatePage.Slug = s.generateSlug(duplicatePage.Title)
	}

	// Set default status
	if duplicatePage.Status == "" {
		duplicatePage.Status = StatusDraft
	}

	// Validate slug uniqueness
	if err := s.validateSlugUniqueness(tenantID, duplicatePage.Slug, uuid.Nil); err != nil {
		return nil, err
	}

	// Create the duplicate page
	createdPage, err := s.repository.CreatePage(duplicatePage)
	if err != nil {
		return nil, fmt.Errorf("failed to create duplicate page: %w", err)
	}

	// Copy tags and categories from original
	originalPageWithRelations, err := s.repository.GetPage(tenantID, pageID)
	if err == nil {
		tagIDs := make([]uuid.UUID, len(originalPageWithRelations.Tags))
		for i, tag := range originalPageWithRelations.Tags {
			tagIDs[i] = tag.ID
		}
		if len(tagIDs) > 0 {
			s.repository.AssociatePageTags(createdPage.ID, tagIDs)
		}

		categoryIDs := make([]uuid.UUID, len(originalPageWithRelations.Categories))
		for i, category := range originalPageWithRelations.Categories {
			categoryIDs[i] = category.ID
		}
		if len(categoryIDs) > 0 {
			s.repository.AssociatePageCategories(createdPage.ID, categoryIDs)
		}
	}

	return s.repository.GetPage(tenantID, createdPage.ID)
}

func (s *Service) IncrementPageViews(tenantID, pageID uuid.UUID) error {
	return s.repository.IncrementPageViews(tenantID, pageID)
}

// Media Management Service Methods

func (s *Service) UploadMedia(tenantID, userID uuid.UUID, req UploadMediaRequest) (*Media, error) {
	// Validate file type
	mediaType, err := s.getMediaType(req.Header.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("unsupported file type: %w", err)
	}

	// Generate file path
	fileName := s.generateFileName(req.Header.Filename)
	filePath := s.generateFilePath(tenantID, fileName)

	// Save file to disk
	if err := s.saveFile(req.File, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Get file size
	fileSize, _ := req.File.Seek(0, 2)
	req.File.Seek(0, 0)

	// Create media record
	media := &Media{
		TenantID:    tenantID,
		Type:        mediaType,
		Title:       req.Title,
		Description: req.Description,
		FileName:    fileName,
		FilePath:    filePath,
		FileSize:    fileSize,
		MimeType:    req.Header.Header.Get("Content-Type"),
		AltText:     req.AltText,
		UploadedBy:  userID,
	}

	// Set default title if empty
	if media.Title == "" {
		media.Title = req.Header.Filename
	}

	// Get image dimensions if it's an image
	if media.IsImage() {
		width, height, err := s.getImageDimensions(filePath)
		if err == nil {
			media.Width = width
			media.Height = height
		}
	}

	return s.repository.CreateMedia(media)
}

func (s *Service) UpdateMedia(tenantID, mediaID uuid.UUID, req UpdateMediaRequest) (*Media, error) {
	media, err := s.repository.GetMedia(tenantID, mediaID)
	if err != nil {
		return nil, fmt.Errorf("media not found: %w", err)
	}

	if req.Title != nil {
		media.Title = *req.Title
	}
	if req.Description != nil {
		media.Description = *req.Description
	}
	if req.AltText != nil {
		media.AltText = *req.AltText
	}
	if req.MetaTitle != nil {
		media.MetaTitle = *req.MetaTitle
	}
	if req.MetaDescription != nil {
		media.MetaDescription = *req.MetaDescription
	}

	return s.repository.UpdateMedia(media)
}

func (s *Service) DeleteMedia(tenantID, mediaID uuid.UUID) error {
	media, err := s.repository.GetMedia(tenantID, mediaID)
	if err != nil {
		return fmt.Errorf("media not found: %w", err)
	}

	// Delete file from disk
	if err := os.Remove(media.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return s.repository.DeleteMedia(tenantID, mediaID)
}

func (s *Service) GetMedia(tenantID, mediaID uuid.UUID) (*Media, error) {
	return s.repository.GetMedia(tenantID, mediaID)
}

func (s *Service) GetMediaLibrary(tenantID uuid.UUID, filter MediaLibraryFilter) ([]Media, int64, error) {
	return s.repository.GetMediaLibrary(tenantID, filter)
}

// Menu Management Service Methods

func (s *Service) CreateMenu(tenantID uuid.UUID, req CreateMenuRequest) (*Menu, error) {
	menu := &Menu{
		TenantID: tenantID,
		Name:     req.Name,
		Location: req.Location,
		IsActive: req.IsActive,
	}

	return s.repository.CreateMenu(menu)
}

func (s *Service) UpdateMenu(tenantID, menuID uuid.UUID, req UpdateMenuRequest) (*Menu, error) {
	menu, err := s.repository.GetMenu(tenantID, menuID)
	if err != nil {
		return nil, fmt.Errorf("menu not found: %w", err)
	}

	if req.Name != nil {
		menu.Name = *req.Name
	}
	if req.Location != nil {
		menu.Location = *req.Location
	}
	if req.IsActive != nil {
		menu.IsActive = *req.IsActive
	}

	return s.repository.UpdateMenu(menu)
}

func (s *Service) DeleteMenu(tenantID, menuID uuid.UUID) error {
	return s.repository.DeleteMenu(tenantID, menuID)
}

func (s *Service) GetMenu(tenantID, menuID uuid.UUID) (*Menu, error) {
	return s.repository.GetMenu(tenantID, menuID)
}

func (s *Service) GetMenus(tenantID uuid.UUID, location string, activeOnly bool) ([]Menu, error) {
	return s.repository.GetMenus(tenantID, location, activeOnly)
}

func (s *Service) GetMenuByLocation(tenantID uuid.UUID, location string) (*Menu, error) {
	return s.repository.GetMenuByLocation(tenantID, location)
}

func (s *Service) GetPublicMenuByLocation(tenantID uuid.UUID, location string) (*Menu, error) {
	menu, err := s.repository.GetMenuByLocation(tenantID, location)
	if err != nil {
		return nil, err
	}

	if !menu.IsActive {
		return nil, errors.New("menu not active")
	}

	return menu, nil
}

// Menu Item Management Service Methods

func (s *Service) CreateMenuItem(tenantID, menuID uuid.UUID, req CreateMenuItemRequest) (*MenuItem, error) {
	menuItem := &MenuItem{
		MenuID:    menuID,
		ParentID:  req.ParentID,
		Title:     req.Title,
		URL:       req.URL,
		PageID:    req.PageID,
		Target:    req.Target,
		CSSClass:  req.CSSClass,
		IconClass: req.IconClass,
		IsActive:  req.IsActive,
		SortOrder: req.SortOrder,
	}

	// Set default target
	if menuItem.Target == "" {
		menuItem.Target = "_self"
	}

	return s.repository.CreateMenuItem(menuItem)
}

func (s *Service) UpdateMenuItem(tenantID, itemID uuid.UUID, req UpdateMenuItemRequest) (*MenuItem, error) {
	menuItem, err := s.repository.GetMenuItem(tenantID, itemID)
	if err != nil {
		return nil, fmt.Errorf("menu item not found: %w", err)
	}

	if req.ParentID != nil {
		menuItem.ParentID = req.ParentID
	}
	if req.Title != nil {
		menuItem.Title = *req.Title
	}
	if req.URL != nil {
		menuItem.URL = *req.URL
	}
	if req.PageID != nil {
		menuItem.PageID = req.PageID
	}
	if req.Target != nil {
		menuItem.Target = *req.Target
	}
	if req.CSSClass != nil {
		menuItem.CSSClass = *req.CSSClass
	}
	if req.IconClass != nil {
		menuItem.IconClass = *req.IconClass
	}
	if req.IsActive != nil {
		menuItem.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		menuItem.SortOrder = *req.SortOrder
	}

	return s.repository.UpdateMenuItem(menuItem)
}

func (s *Service) DeleteMenuItem(tenantID, itemID uuid.UUID) error {
	return s.repository.DeleteMenuItem(tenantID, itemID)
}

func (s *Service) ReorderMenuItems(tenantID, menuID uuid.UUID, req ReorderMenuItemsRequest) error {
	return s.repository.ReorderMenuItems(tenantID, menuID, req.ItemOrders)
}

// Tag Management Service Methods

func (s *Service) CreateTag(tenantID uuid.UUID, req CreateTagRequest) (*Tag, error) {
	tag := &Tag{
		TenantID:    tenantID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
	}

	// Generate slug if not provided
	if tag.Slug == "" {
		tag.Slug = s.generateSlug(tag.Name)
	}

	return s.repository.CreateTag(tag)
}

func (s *Service) UpdateTag(tenantID, tagID uuid.UUID, req UpdateTagRequest) (*Tag, error) {
	tag, err := s.repository.GetTag(tenantID, tagID)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	if req.Name != nil {
		tag.Name = *req.Name
	}
	if req.Slug != nil {
		tag.Slug = *req.Slug
	}
	if req.Description != nil {
		tag.Description = *req.Description
	}
	if req.Color != nil {
		tag.Color = *req.Color
	}

	return s.repository.UpdateTag(tag)
}

func (s *Service) DeleteTag(tenantID, tagID uuid.UUID) error {
	return s.repository.DeleteTag(tenantID, tagID)
}

func (s *Service) GetTags(tenantID uuid.UUID, search string) ([]Tag, error) {
	return s.repository.GetTags(tenantID, search)
}

// Category Management Service Methods

func (s *Service) CreateCategory(tenantID uuid.UUID, req CreateCategoryRequest) (*Category, error) {
	category := &Category{
		TenantID:    tenantID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
	}

	// Generate slug if not provided
	if category.Slug == "" {
		category.Slug = s.generateSlug(category.Name)
	}

	return s.repository.CreateCategory(category)
}

func (s *Service) UpdateCategory(tenantID, categoryID uuid.UUID, req UpdateCategoryRequest) (*Category, error) {
	category, err := s.repository.GetCategory(tenantID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	return s.repository.UpdateCategory(category)
}

func (s *Service) DeleteCategory(tenantID, categoryID uuid.UUID) error {
	return s.repository.DeleteCategory(tenantID, categoryID)
}

func (s *Service) GetCategories(tenantID uuid.UUID, activeOnly bool) ([]Category, error) {
	return s.repository.GetCategories(tenantID, activeOnly)
}

// SEO Management Service Methods

func (s *Service) UpdateSEOSettings(tenantID uuid.UUID, req UpdateSEOSettingsRequest) (*SEOSettings, error) {
	settings, err := s.repository.GetSEOSettings(tenantID)
	if err != nil {
		// Create new settings if not exists
		settings = &SEOSettings{
			TenantID: tenantID,
		}
	}

	if req.SiteTitle != nil {
		settings.SiteTitle = *req.SiteTitle
	}
	if req.SiteDescription != nil {
		settings.SiteDescription = *req.SiteDescription
	}
	if req.DefaultMetaTitle != nil {
		settings.DefaultMetaTitle = *req.DefaultMetaTitle
	}
	if req.DefaultMetaDesc != nil {
		settings.DefaultMetaDesc = *req.DefaultMetaDesc
	}
	if req.DefaultOGImage != nil {
		settings.DefaultOGImage = *req.DefaultOGImage
	}
	if req.GoogleAnalyticsID != nil {
		settings.GoogleAnalyticsID = *req.GoogleAnalyticsID
	}
	if req.GoogleTagManagerID != nil {
		settings.GoogleTagManagerID = *req.GoogleTagManagerID
	}
	if req.FacebookPixelID != nil {
		settings.FacebookPixelID = *req.FacebookPixelID
	}
	if req.GoogleVerification != nil {
		settings.GoogleVerification = *req.GoogleVerification
	}
	if req.BingVerification != nil {
		settings.BingVerification = *req.BingVerification
	}
	if req.RobotsTxt != nil {
		settings.RobotsTxt = *req.RobotsTxt
	}
	if req.EnableSitemap != nil {
		settings.EnableSitemap = *req.EnableSitemap
	}
	if req.SitemapFrequency != nil {
		settings.SitemapFrequency = *req.SitemapFrequency
	}

	if settings.ID == uuid.Nil {
		return s.repository.CreateSEOSettings(settings)
	}

	return s.repository.UpdateSEOSettings(settings)
}

func (s *Service) GetSEOSettings(tenantID uuid.UUID) (*SEOSettings, error) {
	return s.repository.GetSEOSettings(tenantID)
}

func (s *Service) GenerateSitemap(tenantID uuid.UUID) (string, error) {
	pages, err := s.repository.GetPublishedPages(tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to get published pages: %w", err)
	}

	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`

	for _, page := range pages {
		priority := "0.8"
		if page.Type == TypePage {
			priority = "1.0"
		} else if page.Type == TypeBlogPost {
			priority = "0.6"
		}

		changefreq := "weekly"
		if page.Type == TypeBlogPost {
			changefreq = "monthly"
		}

		lastmod := page.UpdatedAt.Format("2006-01-02")
		if page.PublishedAt != nil {
			lastmod = page.PublishedAt.Format("2006-01-02")
		}

		sitemap += fmt.Sprintf(`  <url>
    <loc>%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>%s</changefreq>
    <priority>%s</priority>
  </url>
`, page.GetURL(), lastmod, changefreq, priority)
	}

	sitemap += "</urlset>"

	return sitemap, nil
}

func (s *Service) GetRobotsTxt(tenantID uuid.UUID) (string, error) {
	settings, err := s.repository.GetSEOSettings(tenantID)
	if err != nil || settings.RobotsTxt == "" {
		// Return default robots.txt
		return `User-agent: *
Allow: /

Sitemap: /sitemap.xml`, nil
	}

	return settings.RobotsTxt, nil
}

// Analytics Service Methods

func (s *Service) GetContentAnalytics(tenantID uuid.UUID) (*ContentAnalytics, error) {
	analytics := &ContentAnalytics{}

	// Get basic stats
	stats, err := s.repository.GetContentStats(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content stats: %w", err)
	}

	analytics.TotalPages = stats.TotalPages
	analytics.PublishedPages = stats.PublishedPages
	analytics.DraftPages = stats.DraftPages
	analytics.TotalMedia = stats.TotalMedia
	analytics.TotalViews = stats.TotalViews

	// Get popular pages
	popularPages, err := s.repository.GetPopularPages(tenantID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular pages: %w", err)
	}

	analytics.PopularPages = make([]PageAnalytics, len(popularPages))
	for i, page := range popularPages {
		analytics.PopularPages[i] = PageAnalytics{
			Page:      page,
			Views:     int64(page.ViewCount),
			ViewsRate: float64(page.ViewCount) / float64(analytics.TotalViews) * 100,
		}
	}

	// Get recent pages
	recentPages, err := s.repository.GetRecentPages(tenantID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent pages: %w", err)
	}
	analytics.RecentPages = recentPages

	return analytics, nil
}

func (s *Service) GetPopularContent(tenantID uuid.UUID, limit int, contentType string) ([]Page, error) {
	return s.repository.GetPopularPages(tenantID, limit)
}

// Search Service Methods

func (s *Service) SearchContent(tenantID uuid.UUID, filter ContentSearchFilter) ([]ContentSearchResult, int64, error) {
	return s.repository.SearchContent(tenantID, filter)
}

// Helper methods

func (s *Service) generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)
	
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// Remove special characters (keep only alphanumeric and hyphens)
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	slug = reg.ReplaceAllString(slug, "")
	
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")
	
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	
	return slug
}

func (s *Service) validateSlugUniqueness(tenantID uuid.UUID, slug string, excludeID uuid.UUID) error {
	exists, err := s.repository.SlugExists(tenantID, slug, excludeID)
	if err != nil {
		return fmt.Errorf("failed to check slug uniqueness: %w", err)
	}

	if exists {
		return errors.New("slug already exists")
	}

	return nil
}

func (s *Service) getMediaType(mimeType string) (MediaType, error) {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return MediaImage, nil
	case strings.HasPrefix(mimeType, "video/"):
		return MediaVideo, nil
	case strings.HasPrefix(mimeType, "audio/"):
		return MediaAudio, nil
	case mimeType == "application/pdf" || strings.HasPrefix(mimeType, "text/"):
		return MediaDocument, nil
	default:
		return "", errors.New("unsupported media type")
	}
}

func (s *Service) generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)
	
	// Clean the filename
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	cleanName := reg.ReplaceAllString(nameWithoutExt, "_")
	
	// Add timestamp to make it unique
	timestamp := time.Now().Unix()
	
	return fmt.Sprintf("%s_%d%s", cleanName, timestamp, ext)
}

func (s *Service) generateFilePath(tenantID uuid.UUID, fileName string) string {
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	
	return fmt.Sprintf("uploads/%s/%s/%s/%s", tenantID.String(), year, month, fileName)
}

func (s *Service) saveFile(src multipart.File, filePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func (s *Service) getImageDimensions(filePath string) (int, int, error) {
	// This is a placeholder - in a real implementation, you would use
	// an image library like "image" package to get actual dimensions
	return 0, 0, errors.New("not implemented")
}