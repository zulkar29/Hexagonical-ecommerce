package admin

import (
	"time"

	"github.com/google/uuid"
)

// AdminDashboardStats represents dashboard overview statistics
type AdminDashboardStats struct {
	Period        string                 `json:"period"`
	SalesStats    *SalesStats           `json:"sales_stats,omitempty"`
	OrderStats    *OrderStats           `json:"order_stats,omitempty"`
	CustomerStats *CustomerStats        `json:"customer_stats,omitempty"`
	ProductStats  *ProductStats         `json:"product_stats,omitempty"`
	Metrics       map[string]interface{} `json:"metrics,omitempty"`
}

// SalesStats represents sales statistics
type SalesStats struct {
	TotalRevenue    float64 `json:"total_revenue"`
	RevenueGrowth   float64 `json:"revenue_growth"`
	AverageOrder    float64 `json:"average_order"`
	TotalSales      int64   `json:"total_sales"`
	SalesGrowth     float64 `json:"sales_growth"`
	TopProducts     []TopProduct `json:"top_products,omitempty"`
}

// OrderStats represents order statistics
type OrderStats struct {
	TotalOrders     int64   `json:"total_orders"`
	OrderGrowth     float64 `json:"order_growth"`
	PendingOrders   int64   `json:"pending_orders"`
	ProcessingOrders int64  `json:"processing_orders"`
	ShippedOrders   int64   `json:"shipped_orders"`
	CompletedOrders int64   `json:"completed_orders"`
	CancelledOrders int64   `json:"cancelled_orders"`
}

// CustomerStats represents customer statistics
type CustomerStats struct {
	TotalCustomers    int64   `json:"total_customers"`
	CustomerGrowth    float64 `json:"customer_growth"`
	ActiveCustomers   int64   `json:"active_customers"`
	NewCustomers      int64   `json:"new_customers"`
	ReturningCustomers int64  `json:"returning_customers"`
}

// ProductStats represents product statistics
type ProductStats struct {
	TotalProducts   int64 `json:"total_products"`
	ActiveProducts  int64 `json:"active_products"`
	DraftProducts   int64 `json:"draft_products"`
	LowStockProducts int64 `json:"low_stock_products"`
	OutOfStockProducts int64 `json:"out_of_stock_products"`
}

// TopProduct represents top selling product
type TopProduct struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Sales    int64     `json:"sales"`
	Revenue  float64   `json:"revenue"`
}

// QuickStat represents a single quick statistic
type QuickStat struct {
	Value  interface{} `json:"value"`
	Change float64     `json:"change"`
	Trend  string      `json:"trend"` // up, down, stable
}

// AdminStaff represents admin staff member
type AdminStaff struct {
	ID          uuid.UUID  `json:"id" gorm:"primarykey"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty" gorm:"index"`
	UserID      uuid.UUID  `json:"user_id" gorm:"not null;index"`
	Role        string     `json:"role" gorm:"not null"`
	Status      string     `json:"status" gorm:"default:active"`
	Permissions []string   `json:"permissions" gorm:"type:jsonb"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	
	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// User represents user for staff relations
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Avatar    string    `json:"avatar,omitempty"`
}

// AdminRole represents admin role with permissions
type AdminRole struct {
	ID          uuid.UUID  `json:"id" gorm:"primarykey"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty" gorm:"index"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Permissions []string   `json:"permissions" gorm:"type:jsonb"`
	IsSystem    bool       `json:"is_system" gorm:"default:false"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}





// HealthCheck represents a health check result
type HealthCheck struct {
	Status      string  `json:"status"`
	Latency     string  `json:"latency,omitempty"`
	Message     string  `json:"message,omitempty"`
	LastChecked string  `json:"last_checked"`
}

// MemoryMetrics represents memory usage metrics
type MemoryMetrics struct {
	Used      int64   `json:"used"`
	Total     int64   `json:"total"`
	Usage     float64 `json:"usage"`
	Available int64   `json:"available"`
}

// CPUMetrics represents CPU usage metrics
type CPUMetrics struct {
	Usage   float64 `json:"usage"`
	Cores   int     `json:"cores"`
	Load1   float64 `json:"load_1"`
	Load5   float64 `json:"load_5"`
	Load15  float64 `json:"load_15"`
}

// DiskMetrics represents disk usage metrics
type DiskMetrics struct {
	Used      int64   `json:"used"`
	Total     int64   `json:"total"`
	Usage     float64 `json:"usage"`
	Available int64   `json:"available"`
}