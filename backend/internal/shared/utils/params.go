package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrInvalidUUID     = errors.New("invalid UUID format")
	ErrInvalidTenantID = errors.New("invalid tenant ID")
	ErrInvalidUserID   = errors.New("invalid user ID")
	ErrInvalidID       = errors.New("invalid ID")
)

// ParseUUIDParam extracts and validates a UUID parameter from the request
func ParseUUIDParam(c *gin.Context, paramName string) (uuid.UUID, error) {
	paramValue := c.Param(paramName)
	if paramValue == "" {
		return uuid.Nil, errors.New(paramName + " is required")
	}

	parsedUUID, err := uuid.Parse(paramValue)
	if err != nil {
		return uuid.Nil, ErrInvalidUUID
	}

	return parsedUUID, nil
}

// ParseTenantID extracts and validates tenant ID from request parameters
func ParseTenantID(c *gin.Context) (uuid.UUID, error) {
	tenantID, err := ParseUUIDParam(c, "tenant_id")
	if err != nil {
		return uuid.Nil, ErrInvalidTenantID
	}
	return tenantID, nil
}

// ParseUserID extracts and validates user ID from request parameters
func ParseUserID(c *gin.Context) (uuid.UUID, error) {
	userID, err := ParseUUIDParam(c, "user_id")
	if err != nil {
		return uuid.Nil, ErrInvalidUserID
	}
	return userID, nil
}

// ParseID extracts and validates a generic ID parameter
func ParseID(c *gin.Context) (uuid.UUID, error) {
	id, err := ParseUUIDParam(c, "id")
	if err != nil {
		return uuid.Nil, ErrInvalidID
	}
	return id, nil
}

// Note: GetTenantIDFromContext and GetUserIDFromContext already exist in utils.go

// ParseIntParam extracts and validates an integer parameter from the request
func ParseIntParam(c *gin.Context, paramName string, defaultValue int) (int, error) {
	paramValue := c.Param(paramName)
	if paramValue == "" {
		return defaultValue, nil
	}

	intValue, err := strconv.Atoi(paramValue)
	if err != nil {
		return defaultValue, errors.New("invalid " + paramName + " format")
	}

	return intValue, nil
}

// ParseIntQuery extracts and validates an integer query parameter
func ParseIntQuery(c *gin.Context, queryName string, defaultValue int) int {
	queryValue := c.Query(queryName)
	if queryValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(queryValue)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// ParsePaginationParams extracts common pagination parameters
func ParsePaginationParams(c *gin.Context) (page int, limit int) {
	page = ParseIntQuery(c, "page", 1)
	limit = ParseIntQuery(c, "limit", 20)

	// Ensure minimum values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	// Set maximum limit
	if limit > 100 {
		limit = 100
	}

	return page, limit
}

// ParseSortParams extracts sorting parameters
func ParseSortParams(c *gin.Context, defaultSort string) (sort string, order string) {
	sort = c.Query("sort")
	if sort == "" {
		sort = defaultSort
	}

	order = c.Query("order")
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return sort, order
}

// ParseUUIDFromString safely parses a UUID from string
func ParseUUIDFromString(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, errors.New("empty UUID string")
	}
	return uuid.Parse(s)
}

// ParseUUIDSlice parses a slice of UUID strings
func ParseUUIDSlice(strings []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, 0, len(strings))
	for _, s := range strings {
		if parsed, err := uuid.Parse(s); err == nil {
			uuids = append(uuids, parsed)
		} else {
			return nil, err
		}
	}
	return uuids, nil
}