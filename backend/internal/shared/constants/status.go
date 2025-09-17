package constants

// Common Status Types
type Status string
type OrderStatus string
type PaymentStatus string
type UserStatus string
type TenantStatus string
type ProductStatus string
type DiscountStatus string
type ShipmentStatus string

// Common Status Values
const (
	StatusPending   Status = "pending"
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusCancelled Status = "cancelled"
	StatusCompleted Status = "completed"
	StatusDraft     Status = "draft"
	StatusArchived  Status = "archived"
)

// Order Status Values
const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusReturned   OrderStatus = "returned"
	OrderStatusRefunded   OrderStatus = "refunded"
	OrderStatusFailed     OrderStatus = "failed"
	OrderStatusOnHold     OrderStatus = "on_hold"
)

// Payment Status Values
const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusPartial    PaymentStatus = "partial"
	PaymentStatusExpired    PaymentStatus = "expired"
)

// User Status Values
const (
	UserStatusActive      UserStatus = "active"
	UserStatusInactive    UserStatus = "inactive"
	UserStatusSuspended   UserStatus = "suspended"
	UserStatusPending     UserStatus = "pending"
	UserStatusBlocked     UserStatus = "blocked"
	UserStatusDeactivated UserStatus = "deactivated"
)

// Tenant Status Values
const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusInactive  TenantStatus = "inactive"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusTrial     TenantStatus = "trial"
	TenantStatusExpired   TenantStatus = "expired"
	TenantStatusCancelled TenantStatus = "cancelled"
)

// Product Status Values
const (
	ProductStatusActive      ProductStatus = "active"
	ProductStatusInactive    ProductStatus = "inactive"
	ProductStatusDraft       ProductStatus = "draft"
	ProductStatusArchived    ProductStatus = "archived"
	ProductStatusOutOfStock  ProductStatus = "out_of_stock"
	ProductStatusDiscontinued ProductStatus = "discontinued"
)

// Discount Status Values
const (
	DiscountStatusActive    DiscountStatus = "active"
	DiscountStatusInactive  DiscountStatus = "inactive"
	DiscountStatusScheduled DiscountStatus = "scheduled"
	DiscountStatusExpired   DiscountStatus = "expired"
	DiscountStatusUsed      DiscountStatus = "used"
	DiscountStatusCancelled DiscountStatus = "cancelled"
)

// Shipment Status Values
const (
	ShipmentStatusPending    ShipmentStatus = "pending"
	ShipmentStatusPacked     ShipmentStatus = "packed"
	ShipmentStatusShipped    ShipmentStatus = "shipped"
	ShipmentStatusInTransit  ShipmentStatus = "in_transit"
	ShipmentStatusDelivered  ShipmentStatus = "delivered"
	ShipmentStatusFailed     ShipmentStatus = "failed"
	ShipmentStatusReturned   ShipmentStatus = "returned"
	ShipmentStatusCancelled  ShipmentStatus = "cancelled"
)

// Status Groups for validation
var (
	ValidOrderStatuses = []string{
		string(OrderStatusPending),
		string(OrderStatusConfirmed),
		string(OrderStatusProcessing),
		string(OrderStatusShipped),
		string(OrderStatusDelivered),
		string(OrderStatusCancelled),
		string(OrderStatusReturned),
		string(OrderStatusRefunded),
		string(OrderStatusFailed),
		string(OrderStatusOnHold),
	}

	ValidPaymentStatuses = []string{
		string(PaymentStatusPending),
		string(PaymentStatusProcessing),
		string(PaymentStatusCompleted),
		string(PaymentStatusFailed),
		string(PaymentStatusCancelled),
		string(PaymentStatusRefunded),
		string(PaymentStatusPartial),
		string(PaymentStatusExpired),
	}

	ValidUserStatuses = []string{
		string(UserStatusActive),
		string(UserStatusInactive),
		string(UserStatusSuspended),
		string(UserStatusPending),
		string(UserStatusBlocked),
		string(UserStatusDeactivated),
	}

	ValidTenantStatuses = []string{
		string(TenantStatusActive),
		string(TenantStatusInactive),
		string(TenantStatusSuspended),
		string(TenantStatusTrial),
		string(TenantStatusExpired),
		string(TenantStatusCancelled),
	}

	ValidProductStatuses = []string{
		string(ProductStatusActive),
		string(ProductStatusInactive),
		string(ProductStatusDraft),
		string(ProductStatusArchived),
		string(ProductStatusOutOfStock),
		string(ProductStatusDiscontinued),
	}
)

// Helper functions to check valid statuses
func IsValidOrderStatus(status string) bool {
	for _, validStatus := range ValidOrderStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func IsValidPaymentStatus(status string) bool {
	for _, validStatus := range ValidPaymentStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func IsValidUserStatus(status string) bool {
	for _, validStatus := range ValidUserStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func IsValidTenantStatus(status string) bool {
	for _, validStatus := range ValidTenantStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func IsValidProductStatus(status string) bool {
	for _, validStatus := range ValidProductStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}