package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseRepository provides common repository functionality
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// PaginationFilter represents common pagination parameters
type PaginationFilter struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}

// PaginationResponse represents paginated response metadata
type PaginationResponse struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	ItemsPerPage int  `json:"items_per_page"`
	HasNextPage bool  `json:"has_next_page"`
	HasPrevPage bool  `json:"has_prev_page"`
}

// Paginate applies pagination to a GORM query
func (r *BaseRepository) Paginate(query *gorm.DB, filter PaginationFilter) (*gorm.DB, int, int) {
	// Set defaults
	page := filter.Page
	if page < 1 {
		page = 1
	}

	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Apply sorting
	orderBy := "created_at DESC"
	if filter.Sort != "" {
		order := "ASC"
		if filter.Order == "desc" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filter.Sort, order)
	}

	return query.Order(orderBy).Limit(limit).Offset(offset), limit, offset
}

// CreatePaginationResponse creates pagination metadata
func (r *BaseRepository) CreatePaginationResponse(page, limit int, totalItems int64) *PaginationResponse {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	return &PaginationResponse{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: limit,
		HasNextPage:  page < totalPages,
		HasPrevPage:  page > 1,
	}
}

// WithTransaction executes a function within a database transaction
func (r *BaseRepository) WithTransaction(fn func(*gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// ApplyTenantFilter applies tenant isolation to queries
func (r *BaseRepository) ApplyTenantFilter(query *gorm.DB, tenantID uuid.UUID) *gorm.DB {
	return query.Where("tenant_id = ?", tenantID)
}

// ExistsBy checks if a record exists with given conditions
func (r *BaseRepository) ExistsBy(model interface{}, condition string, args ...interface{}) (bool, error) {
	var count int64
	err := r.db.Model(model).Where(condition, args...).Count(&count).Error
	return count > 0, err
}

// SoftDelete performs soft delete with tenant isolation
func (r *BaseRepository) SoftDelete(model interface{}, tenantID uuid.UUID, id uuid.UUID) error {
	result := r.db.Where("tenant_id = ? AND id = ?", tenantID, id).Delete(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found or access denied")
	}
	return nil
}

// BulkUpdate performs bulk update with tenant isolation
func (r *BaseRepository) BulkUpdate(model interface{}, tenantID uuid.UUID, ids []uuid.UUID, updates map[string]interface{}) error {
	if len(ids) == 0 {
		return errors.New("no IDs provided")
	}

	result := r.db.Model(model).Where("tenant_id = ? AND id IN ?", tenantID, ids).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records updated or access denied")
	}
	return nil
}

// GetDB returns the underlying database connection
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}