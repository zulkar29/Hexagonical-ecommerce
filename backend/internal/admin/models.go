package admin

import (
	"time"

	"github.com/google/uuid"
)

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalUsers    int64     `json:"total_users"`
	TotalOrders   int64     `json:"total_orders"`
	TotalRevenue  float64   `json:"total_revenue"`
	TotalProducts int64     `json:"total_products"`
	Period        string    `json:"period"`
	GeneratedAt   time.Time `json:"generated_at"`
}

// QuickStats represents quick dashboard statistics
type QuickStats struct {
	ActiveUsers   int64     `json:"active_users"`
	PendingOrders int64     `json:"pending_orders"`
	LowStockItems int64     `json:"low_stock_items"`
	TodayRevenue  float64   `json:"today_revenue"`
	GeneratedAt   time.Time `json:"generated_at"`
}

// Staff represents a staff member
type Staff struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Email     string     `json:"email" gorm:"uniqueIndex"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	TenantID  *uuid.UUID `json:"tenant_id" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Role represents a user role
type Role struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	TenantID    *uuid.UUID `json:"tenant_id" gorm:"type:uuid"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ActivityLog represents an activity log entry
type ActivityLog struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	Action    string     `json:"action"`
	Resource  string     `json:"resource"`
	Details   string     `json:"details"`
	IPAddress string     `json:"ip_address"`
	UserAgent string     `json:"user_agent"`
	TenantID  *uuid.UUID `json:"tenant_id" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at"`
}

// SystemHealth represents system health status
type SystemHealth struct {
	Status      string                 `json:"status"`
	Services    map[string]interface{} `json:"services"`
	GeneratedAt time.Time              `json:"generated_at"`
}

// Request/Response types

// DashboardRequest represents a dashboard stats request
type DashboardRequest struct {
	Period  string   `json:"period" form:"period"`
	Metrics []string `json:"metrics" form:"metrics"`
}

// StaffRequest represents a staff creation/update request
type StaffRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

// RoleRequest represents a role creation/update request
type RoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// ActivityLogFilter represents filters for activity logs
type ActivityLogFilter struct {
	UserID    *uuid.UUID `json:"user_id" form:"user_id"`
	Action    string     `json:"action" form:"action"`
	Resource  string     `json:"resource" form:"resource"`
	FromDate  *time.Time `json:"from_date" form:"from_date"`
	ToDate    *time.Time `json:"to_date" form:"to_date"`
	StartDate *time.Time `json:"start_date" form:"start_date"`
	EndDate   *time.Time `json:"end_date" form:"end_date"`
	Limit     int        `json:"limit" form:"limit"`
	Offset    int        `json:"offset" form:"offset"`
}