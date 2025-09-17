package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ecommerce-saas/internal/shared/repository"
	"ecommerce-saas/internal/shared/utils"
)

// CRUDService defines the interface for CRUD operations
type CRUDService[T any, CreateReq any, UpdateReq any] interface {
	Create(tenantID uuid.UUID, req CreateReq) (*T, error)
	GetByID(tenantID uuid.UUID, id uuid.UUID) (*T, error)
	Update(tenantID uuid.UUID, id uuid.UUID, req UpdateReq) (*T, error)
	Delete(tenantID uuid.UUID, id uuid.UUID) error
	List(tenantID uuid.UUID, filter repository.PaginationFilter) ([]*T, *repository.PaginationResponse, error)
}

// BaseCRUDHandler provides common CRUD operations for entities
type BaseCRUDHandler[T any, CreateReq any, UpdateReq any] struct {
	service     CRUDService[T, CreateReq, UpdateReq]
	entityName  string
	routePrefix string
}

// NewBaseCRUDHandler creates a new base CRUD handler
func NewBaseCRUDHandler[T any, CreateReq any, UpdateReq any](
	service CRUDService[T, CreateReq, UpdateReq],
	entityName string,
	routePrefix string,
) *BaseCRUDHandler[T, CreateReq, UpdateReq] {
	return &BaseCRUDHandler[T, CreateReq, UpdateReq]{
		service:     service,
		entityName:  entityName,
		routePrefix: routePrefix,
	}
}

// Create handles POST requests to create a new entity
func (h *BaseCRUDHandler[T, CreateReq, UpdateReq]) Create(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	var req CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}

	entity, err := h.service.Create(tenantID, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	RespondWithCreated(c, entity, h.entityName+" created successfully")
}

// GetByID handles GET requests to retrieve an entity by ID
func (h *BaseCRUDHandler[T, CreateReq, UpdateReq]) GetByID(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	id, err := utils.ParseID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	entity, err := h.service.GetByID(tenantID, id)
	if err != nil {
		HandleError(c, err)
		return
	}

	RespondWithSuccess(c, http.StatusOK, entity)
}

// Update handles PUT/PATCH requests to update an entity
func (h *BaseCRUDHandler[T, CreateReq, UpdateReq]) Update(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	id, err := utils.ParseID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}

	entity, err := h.service.Update(tenantID, id, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	RespondWithUpdated(c, entity, h.entityName+" updated successfully")
}

// Delete handles DELETE requests to remove an entity
func (h *BaseCRUDHandler[T, CreateReq, UpdateReq]) Delete(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	id, err := utils.ParseID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = h.service.Delete(tenantID, id)
	if err != nil {
		HandleError(c, err)
		return
	}

	RespondWithDeleted(c, h.entityName+" deleted successfully")
}

// List handles GET requests to list entities with pagination
func (h *BaseCRUDHandler[T, CreateReq, UpdateReq]) List(c *gin.Context) {
	tenantID, err := utils.GetTenantIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Parse pagination and sorting parameters
	page, limit := utils.ParsePaginationParams(c)
	sort, order := utils.ParseSortParams(c, "created_at")

	filter := repository.PaginationFilter{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
	}

	entities, pagination, err := h.service.List(tenantID, filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	RespondWithSuccessAndMeta(c, http.StatusOK, entities, pagination)
}

// RegisterRoutes registers all CRUD routes for the handler
func (h *BaseCRUDHandler[T, CreateReq, UpdateReq]) RegisterRoutes(router *gin.RouterGroup) {
	entityRoutes := router.Group(h.routePrefix)
	{
		entityRoutes.POST("/", h.Create)
		entityRoutes.GET("/", h.List)
		entityRoutes.GET("/:id", h.GetByID)
		entityRoutes.PUT("/:id", h.Update)
		entityRoutes.PATCH("/:id", h.Update)
		entityRoutes.DELETE("/:id", h.Delete)
	}
}

// CustomHandler allows adding custom routes beyond basic CRUD
type CustomHandler interface {
	RegisterCustomRoutes(router *gin.RouterGroup)
}

// RegisterRoutesWithCustom registers both CRUD and custom routes
func (h *BaseCRUDHandler[T, CreateReq, UpdateReq]) RegisterRoutesWithCustom(
	router *gin.RouterGroup,
	customHandler CustomHandler,
) {
	// Register standard CRUD routes
	h.RegisterRoutes(router)

	// Register custom routes if handler implements CustomHandler
	if customHandler != nil {
		customHandler.RegisterCustomRoutes(router)
	}
}