package platform

import (
	"time"

	"github.com/google/uuid"
)

// PlatformStats represents platform-wide statistics
type PlatformStats struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TotalTenants     int       `json:"total_tenants"`
	ActiveTenants    int       `json:"active_tenants"`
	TotalRevenue     float64   `json:"total_revenue"`
	MonthlyRevenue   float64   `json:"monthly_revenue"`
	TotalUsers       int       `json:"total_users"`
	ActiveUsers      int       `json:"active_users"`
	GrowthRate       float64   `json:"growth_rate"`
	SystemUptime     float64   `json:"system_uptime"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// PlatformAdmin represents platform administrators
type PlatformAdmin struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      *uuid.UUID `json:"user_id" gorm:"type:uuid"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Password    string    `json:"-" gorm:"not null"`
	Role        string    `json:"role" gorm:"not null"` // super_admin, platform_admin, support
	Permissions []string  `json:"permissions" gorm:"type:text[]"`
	Status      string    `json:"status" gorm:"default:'active'"` // active, inactive, suspended
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PlatformRole represents platform-level roles
type PlatformRole struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions" gorm:"type:text[]"`
	IsSystem    bool      `json:"is_system" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PlatformSettings represents platform-wide settings
type PlatformSettings struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Category    string                 `json:"category" gorm:"not null"` // billing, notifications, security, features
	Key         string                 `json:"key" gorm:"not null"`
	Value       map[string]interface{} `json:"value" gorm:"type:jsonb"`
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// PlatformAuditLog represents platform audit logs
type PlatformAuditLog struct {
	ID         uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     *uuid.UUID             `json:"user_id" gorm:"type:uuid"`
	TenantID   *uuid.UUID             `json:"tenant_id" gorm:"type:uuid"`
	Action     string                 `json:"action" gorm:"not null"`
	Resource   string                 `json:"resource"`
	ResourceID *uuid.UUID             `json:"resource_id" gorm:"type:uuid"`
	Details    map[string]interface{} `json:"details" gorm:"type:jsonb"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	CreatedAt  time.Time              `json:"created_at"`
}

// SystemStatus represents overall system health
type SystemStatus struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Status          string    `json:"status"` // healthy, degraded, down
	DatabaseStatus  string    `json:"database_status"`
	RedisStatus     string    `json:"redis_status"`
	APIResponseTime float64   `json:"api_response_time"`
	CPUUsage        float64   `json:"cpu_usage"`
	MemoryUsage     float64   `json:"memory_usage"`
	DiskUsage       float64   `json:"disk_usage"`
	ActiveConnections int     `json:"active_connections"`
	CreatedAt       time.Time `json:"created_at"`
}

// Request/Response DTOs
type DashboardRequest struct {
	Period  string   `json:"period"`  // day, week, month
	Metrics []string `json:"metrics"` // tenants, revenue, users, system
}

type DashboardResponse struct {
	Stats          *PlatformStats            `json:"stats"`
	TenantGrowth   []map[string]interface{}  `json:"tenant_growth"`
	RevenueMetrics []map[string]interface{}  `json:"revenue_metrics"`
	UserMetrics    []map[string]interface{}  `json:"user_metrics"`
	SystemMetrics  []map[string]interface{}  `json:"system_metrics"`
}

type AdminRequest struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	Role        string    `json:"role" binding:"required"`
	Permissions []string  `json:"permissions"`
	Status      string    `json:"status"`
}

type RoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type SettingsRequest struct {
	Category string                 `json:"category" binding:"required"`
	Key      string                 `json:"key" binding:"required"`
	Value    map[string]interface{} `json:"value" binding:"required"`
}

type TenantRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Domain      string                 `json:"domain" binding:"required"`
	Plan        string                 `json:"plan"`
	Settings    map[string]interface{} `json:"settings"`
	OwnerEmail  string                 `json:"owner_email" binding:"required,email"`
}

type AuditLogFilter struct {
	UserID    *uuid.UUID `json:"user_id"`
	TenantID  *uuid.UUID `json:"tenant_id"`
	Action    string     `json:"action"`
	Resource  string     `json:"resource"`
	DateFrom  *time.Time `json:"date_from"`
	DateTo    *time.Time `json:"date_to"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// Platform Admin Request Types
type CreatePlatformAdminRequest struct {
	Email       string   `json:"email" binding:"required,email"`
	Name        string   `json:"name" binding:"required"`
	Password    string   `json:"password" binding:"required,min=8"`
	Role        string   `json:"role" binding:"required"`
	Permissions []string `json:"permissions"`
}

type UpdatePlatformAdminRequest struct {
	Name        *string  `json:"name"`
	Email       *string  `json:"email"`
	Role        *string  `json:"role"`
	Status      *string  `json:"status"`
	Permissions []string `json:"permissions"`
	Password    *string  `json:"password"`
}

// Platform Role Request Types
type CreatePlatformRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type UpdatePlatformRoleRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Permissions []string `json:"permissions"`
}

// Platform Settings Request Types
type UpdatePlatformSettingsRequest struct {
	Category    string                 `json:"category" binding:"required"`
	Key         string                 `json:"key" binding:"required"`
	Value       map[string]interface{} `json:"value" binding:"required"`
	Description *string                `json:"description"`
}

// Tenant Management Request/Response Types
type ListTenantsRequest struct {
	Status   string   `json:"status"`
	Plan     string   `json:"plan"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	Search   string   `json:"search"`
	SortBy   string   `json:"sort_by"`
	SortDesc bool     `json:"sort_desc"`
	Include  []string `json:"include"`
}

type ListTenantsResponse struct {
	Tenants []map[string]interface{} `json:"tenants"`
	Total   int                      `json:"total"`
	Limit   int                      `json:"limit"`
	Offset  int                      `json:"offset"`
}

type UpdateTenantRequest struct {
	Name     *string                `json:"name"`
	Domain   *string                `json:"domain"`
	Plan     *string                `json:"plan"`
	Status   *string                `json:"status"`
	Settings map[string]interface{} `json:"settings"`
}

// Audit Log Request/Response Types
type GetAuditLogsRequest struct {
	UserID    *uuid.UUID `json:"user_id"`
	TenantID  *uuid.UUID `json:"tenant_id"`
	Action    string     `json:"action"`
	Resource  string     `json:"resource"`
	DateFrom  *time.Time `json:"date_from"`
	DateTo    *time.Time `json:"date_to"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

type GetAuditLogsResponse struct {
	Logs   []*PlatformAuditLog `json:"logs"`
	Total  int                 `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}