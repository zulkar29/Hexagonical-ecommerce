package content

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Page Management Endpoints

func (h *Handler) CreatePage(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	page, err := h.service.CreatePage(tenantID.(uuid.UUID), userID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": page})
}

func (h *Handler) UpdatePage(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var req UpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	page, err := h.service.UpdatePage(tenantID.(uuid.UUID), pageID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": page})
}

func (h *Handler) DeletePage(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	if err := h.service.DeletePage(tenantID.(uuid.UUID), pageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}

func (h *Handler) GetPage(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	page, err := h.service.GetPage(tenantID.(uuid.UUID), pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": page})
}

func (h *Handler) GetPageBySlug(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug is required"})
		return
	}

	page, err := h.service.GetPageBySlug(tenantID.(uuid.UUID), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Increment view count for published pages
	if page.IsPublished() {
		go h.service.IncrementPageViews(tenantID.(uuid.UUID), page.ID)
	}

	c.JSON(http.StatusOK, gin.H{"data": page})
}

func (h *Handler) GetPages(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	contentType := c.Query("type")
	search := c.Query("search")

	filter := PageListFilter{
		Page:   page,
		Limit:  limit,
		Status: status,
		Type:   contentType,
		Search: search,
	}

	pages, total, err := h.service.GetPages(tenantID.(uuid.UUID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": pages,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *Handler) PublishPage(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	page, err := h.service.PublishPage(tenantID.(uuid.UUID), pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": page})
}

func (h *Handler) UnpublishPage(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	page, err := h.service.UnpublishPage(tenantID.(uuid.UUID), pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": page})
}

func (h *Handler) DuplicatePage(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var req DuplicatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	page, err := h.service.DuplicatePage(tenantID.(uuid.UUID), userID.(uuid.UUID), pageID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": page})
}

// Media Management Endpoints

func (h *Handler) UploadMedia(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Optional metadata
	title := c.PostForm("title")
	description := c.PostForm("description")
	altText := c.PostForm("alt_text")

	req := UploadMediaRequest{
		File:        file,
		Header:      header,
		Title:       title,
		Description: description,
		AltText:     altText,
	}

	media, err := h.service.UploadMedia(tenantID.(uuid.UUID), userID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": media})
}

func (h *Handler) UpdateMedia(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	mediaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	var req UpdateMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	media, err := h.service.UpdateMedia(tenantID.(uuid.UUID), mediaID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": media})
}

func (h *Handler) DeleteMedia(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	mediaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	if err := h.service.DeleteMedia(tenantID.(uuid.UUID), mediaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

func (h *Handler) GetMedia(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	mediaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	media, err := h.service.GetMedia(tenantID.(uuid.UUID), mediaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": media})
}

func (h *Handler) GetMediaLibrary(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	mediaType := c.Query("type")
	search := c.Query("search")

	filter := MediaLibraryFilter{
		Page:   page,
		Limit:  limit,
		Type:   mediaType,
		Search: search,
	}

	media, total, err := h.service.GetMediaLibrary(tenantID.(uuid.UUID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": media,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// Menu Management Endpoints

func (h *Handler) CreateMenu(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	menu, err := h.service.CreateMenu(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": menu})
}

func (h *Handler) UpdateMenu(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	menuID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	var req UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	menu, err := h.service.UpdateMenu(tenantID.(uuid.UUID), menuID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

func (h *Handler) DeleteMenu(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	menuID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	if err := h.service.DeleteMenu(tenantID.(uuid.UUID), menuID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
}

func (h *Handler) GetMenu(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	menuID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	menu, err := h.service.GetMenu(tenantID.(uuid.UUID), menuID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

func (h *Handler) GetMenus(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	location := c.Query("location")
	activeOnly := c.DefaultQuery("active_only", "false") == "true"

	menus, err := h.service.GetMenus(tenantID.(uuid.UUID), location, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menus})
}

func (h *Handler) GetMenuByLocation(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	location := c.Param("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}

	menu, err := h.service.GetMenuByLocation(tenantID.(uuid.UUID), location)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// Menu Item Management Endpoints

func (h *Handler) CreateMenuItem(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	menuID, err := uuid.Parse(c.Param("menu_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	var req CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	menuItem, err := h.service.CreateMenuItem(tenantID.(uuid.UUID), menuID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": menuItem})
}

func (h *Handler) UpdateMenuItem(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	var req UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	menuItem, err := h.service.UpdateMenuItem(tenantID.(uuid.UUID), itemID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menuItem})
}

func (h *Handler) DeleteMenuItem(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	if err := h.service.DeleteMenuItem(tenantID.(uuid.UUID), itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
}

func (h *Handler) ReorderMenuItems(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	menuID, err := uuid.Parse(c.Param("menu_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	var req ReorderMenuItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.service.ReorderMenuItems(tenantID.(uuid.UUID), menuID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu items reordered successfully"})
}

// Tag Management Endpoints

func (h *Handler) CreateTag(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	tag, err := h.service.CreateTag(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": tag})
}

func (h *Handler) UpdateTag(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	tag, err := h.service.UpdateTag(tenantID.(uuid.UUID), tagID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tag})
}

func (h *Handler) DeleteTag(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := h.service.DeleteTag(tenantID.(uuid.UUID), tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

func (h *Handler) GetTags(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	search := c.Query("search")

	tags, err := h.service.GetTags(tenantID.(uuid.UUID), search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// Category Management Endpoints

func (h *Handler) CreateCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	category, err := h.service.CreateCategory(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	category, err := h.service.UpdateCategory(tenantID.(uuid.UUID), categoryID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.service.DeleteCategory(tenantID.(uuid.UUID), categoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func (h *Handler) GetCategories(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	activeOnly := c.DefaultQuery("active_only", "false") == "true"

	categories, err := h.service.GetCategories(tenantID.(uuid.UUID), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// SEO Management Endpoints

func (h *Handler) UpdateSEOSettings(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var req UpdateSEOSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	seoSettings, err := h.service.UpdateSEOSettings(tenantID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": seoSettings})
}

func (h *Handler) GetSEOSettings(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	seoSettings, err := h.service.GetSEOSettings(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": seoSettings})
}

func (h *Handler) GenerateSitemap(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	sitemap, err := h.service.GenerateSitemap(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, sitemap)
}

func (h *Handler) GetRobotsTxt(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	robots, err := h.service.GetRobotsTxt(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, robots)
}

// Content Analytics Endpoints

func (h *Handler) GetContentAnalytics(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	analytics, err := h.service.GetContentAnalytics(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": analytics})
}

func (h *Handler) GetPopularContent(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	contentType := c.Query("type")

	content, err := h.service.GetPopularContent(tenantID.(uuid.UUID), limit, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": content})
}

// Content Search Endpoints

func (h *Handler) SearchContent(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	contentType := c.Query("type")

	searchFilter := ContentSearchFilter{
		Query:  query,
		Type:   contentType,
		Page:   page,
		Limit:  limit,
	}

	results, total, err := h.service.SearchContent(tenantID.(uuid.UUID), searchFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"query": query,
		},
	})
}

// Public API Endpoints (no auth required)

func (h *Handler) GetPublicPage(c *gin.Context) {
	// Extract tenant from subdomain or domain
	tenantID, err := h.extractTenantFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant"})
		return
	}

	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug is required"})
		return
	}

	page, err := h.service.GetPublishedPageBySlug(tenantID, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	// Increment view count
	go h.service.IncrementPageViews(tenantID, page.ID)

	c.JSON(http.StatusOK, gin.H{"data": page})
}

func (h *Handler) GetPublicMenu(c *gin.Context) {
	// Extract tenant from subdomain or domain
	tenantID, err := h.extractTenantFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant"})
		return
	}

	location := c.Param("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}

	menu, err := h.service.GetPublicMenuByLocation(tenantID, location)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

func (h *Handler) GetPublicSitemap(c *gin.Context) {
	// Extract tenant from subdomain or domain
	tenantID, err := h.extractTenantFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant"})
		return
	}

	sitemap, err := h.service.GenerateSitemap(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate sitemap"})
		return
	}

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, sitemap)
}

func (h *Handler) GetPublicRobotsTxt(c *gin.Context) {
	// Extract tenant from subdomain or domain
	tenantID, err := h.extractTenantFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant"})
		return
	}

	robots, err := h.service.GetRobotsTxt(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get robots.txt"})
		return
	}

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, robots)
}

// Helper methods

func (h *Handler) extractTenantFromRequest(c *gin.Context) (uuid.UUID, error) {
	// Try to get tenant from header first
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	
	// If not in header, try to extract from subdomain
	if tenantIDStr == "" {
		host := c.GetHeader("Host")
		if host != "" {
			parts := strings.Split(host, ".")
			if len(parts) > 2 {
				// Assume first part is the tenant slug
				tenantSlug := parts[0]
				// You would need a method to convert slug to UUID
				// For now, we'll expect UUID in header
				return uuid.Nil, nil
			}
		}
	}
	
	if tenantIDStr == "" {
		return uuid.Nil, nil
	}
	
	return uuid.Parse(tenantIDStr)
}

// Route registration
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	// Protected routes (require authentication)
	v1 := router.Group("/api/v1")
	{
		// Page management
		pages := v1.Group("/pages")
		{
			pages.POST("", h.CreatePage)
			pages.GET("", h.GetPages)
			pages.GET("/:id", h.GetPage)
			pages.PUT("/:id", h.UpdatePage)
			pages.DELETE("/:id", h.DeletePage)
			pages.GET("/slug/:slug", h.GetPageBySlug)
			pages.POST("/:id/publish", h.PublishPage)
			pages.POST("/:id/unpublish", h.UnpublishPage)
			pages.POST("/:id/duplicate", h.DuplicatePage)
		}

		// Media management
		media := v1.Group("/media")
		{
			media.POST("/upload", h.UploadMedia)
			media.GET("", h.GetMediaLibrary)
			media.GET("/:id", h.GetMedia)
			media.PUT("/:id", h.UpdateMedia)
			media.DELETE("/:id", h.DeleteMedia)
		}

		// Menu management
		menus := v1.Group("/menus")
		{
			menus.POST("", h.CreateMenu)
			menus.GET("", h.GetMenus)
			menus.GET("/:id", h.GetMenu)
			menus.PUT("/:id", h.UpdateMenu)
			menus.DELETE("/:id", h.DeleteMenu)
			menus.GET("/location/:location", h.GetMenuByLocation)

			// Menu items
			menus.POST("/:menu_id/items", h.CreateMenuItem)
			menus.PUT("/items/:id", h.UpdateMenuItem)
			menus.DELETE("/items/:id", h.DeleteMenuItem)
			menus.POST("/:menu_id/reorder", h.ReorderMenuItems)
		}

		// Tag management
		tags := v1.Group("/tags")
		{
			tags.POST("", h.CreateTag)
			tags.GET("", h.GetTags)
			tags.PUT("/:id", h.UpdateTag)
			tags.DELETE("/:id", h.DeleteTag)
		}

		// Category management
		categories := v1.Group("/categories")
		{
			categories.POST("", h.CreateCategory)
			categories.GET("", h.GetCategories)
			categories.PUT("/:id", h.UpdateCategory)
			categories.DELETE("/:id", h.DeleteCategory)
		}

		// SEO management
		seo := v1.Group("/seo")
		{
			seo.GET("/settings", h.GetSEOSettings)
			seo.PUT("/settings", h.UpdateSEOSettings)
			seo.GET("/sitemap.xml", h.GenerateSitemap)
			seo.GET("/robots.txt", h.GetRobotsTxt)
		}

		// Analytics
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/content", h.GetContentAnalytics)
			analytics.GET("/popular", h.GetPopularContent)
		}

		// Search
		v1.GET("/search/content", h.SearchContent)
	}

	// Public routes (no authentication required)
	public := router.Group("/public")
	{
		public.GET("/page/:slug", h.GetPublicPage)
		public.GET("/menu/:location", h.GetPublicMenu)
		public.GET("/sitemap.xml", h.GetPublicSitemap)
		public.GET("/robots.txt", h.GetPublicRobotsTxt)
	}
}