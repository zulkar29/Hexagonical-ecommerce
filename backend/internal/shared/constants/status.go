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

// Additional Status Types
type ComponentStatus string
type ReviewStatus string
type NotificationStatus string
type WebhookStatus string
type IntegrationStatus string
type SyncStatus string
type CartStatus string
type ReturnStatus string
type ReferralStatus string
type CommissionStatus string
type FinanceStatus string
type PayoutStatus string
type CampaignStatus string
type EmailStatus string
type SegmentStatus string
type SubscriptionStatus string
type InvoiceStatus string
type TrackingStatus string
type FulfillmentStatus string
type CategoryStatus string
type PageStatus string
type LoginAttemptStatus string
type DeviceStatus string
type AlertStatus string
type DunningActionStatus string

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
	PaymentStatusAbandoned  PaymentStatus = "abandoned"
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
		string(PaymentStatusAbandoned),
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

// Component Status Values
const (
	ComponentStatusDraft    ComponentStatus = "draft"
	ComponentStatusActive   ComponentStatus = "active"
	ComponentStatusInactive ComponentStatus = "inactive"
	ComponentStatusArchived ComponentStatus = "archived"
)

// Review Status Values
const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
)

// Notification Status Values
const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusDelivered NotificationStatus = "delivered"
	NotificationStatusFailed    NotificationStatus = "failed"
)

// Webhook Status Values
const (
	WebhookStatusPending   WebhookStatus = "pending"
	WebhookStatusSent      WebhookStatus = "sent"
	WebhookStatusDelivered WebhookStatus = "delivered"
	WebhookStatusFailed    WebhookStatus = "failed"
)

// Integration Status Values
const (
	IntegrationStatusPending IntegrationStatus = "pending"
	IntegrationStatusActive  IntegrationStatus = "active"
	IntegrationStatusFailed  IntegrationStatus = "failed"
)

// Sync Status Values
const (
	SyncStatusPending   SyncStatus = "pending"
	SyncStatusRunning   SyncStatus = "running"
	SyncStatusCompleted SyncStatus = "completed"
	SyncStatusFailed    SyncStatus = "failed"
)

// Cart Status Values
const (
	CartStatusActive    CartStatus = "active"
	CartStatusAbandoned CartStatus = "abandoned"
	CartStatusExpired   CartStatus = "expired"
)

// Return Status Values
const (
	ReturnStatusPending    ReturnStatus = "pending"
	ReturnStatusApproved   ReturnStatus = "approved"
	ReturnStatusProcessing ReturnStatus = "processing"
	ReturnStatusCompleted  ReturnStatus = "completed"
	ReturnStatusCancelled  ReturnStatus = "cancelled"
)

// Referral Status Values
const (
	ReferralStatusPending   ReferralStatus = "pending"
	ReferralStatusActive    ReferralStatus = "active"
	ReferralStatusCompleted ReferralStatus = "completed"
	ReferralStatusExpired   ReferralStatus = "expired"
	ReferralStatusCancelled ReferralStatus = "cancelled"
)

// Commission Status Values
const (
	CommissionStatusPending   CommissionStatus = "pending"
	CommissionStatusApproved  CommissionStatus = "approved"
	CommissionStatusPaid      CommissionStatus = "paid"
	CommissionStatusFailed    CommissionStatus = "failed"
)

// Finance Status Values
const (
	FinanceStatusPending   FinanceStatus = "pending"
	FinanceStatusProcessing FinanceStatus = "processing"
	FinanceStatusCompleted FinanceStatus = "completed"
	FinanceStatusFailed    FinanceStatus = "failed"
)

// Payout Status Values
const (
	PayoutStatusPending   PayoutStatus = "pending"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted PayoutStatus = "completed"
	PayoutStatusFailed    PayoutStatus = "failed"
	PayoutStatusCancelled PayoutStatus = "cancelled"
)

// Campaign Status Values
const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusScheduled CampaignStatus = "scheduled"
	CampaignStatusRunning   CampaignStatus = "running"
	CampaignStatusCompleted CampaignStatus = "completed"
	CampaignStatusCancelled CampaignStatus = "cancelled"
)

// Email Status Values
const (
	EmailStatusPending   EmailStatus = "pending"
	EmailStatusSent      EmailStatus = "sent"
	EmailStatusDelivered EmailStatus = "delivered"
	EmailStatusOpened    EmailStatus = "opened"
	EmailStatusFailed    EmailStatus = "failed"
)

// Segment Status Values
const (
	SegmentStatusActive   SegmentStatus = "active"
	SegmentStatusInactive SegmentStatus = "inactive"
)

// Subscription Status Values
const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusPending   SubscriptionStatus = "pending"
	SubscriptionStatusTrialing  SubscriptionStatus = "trialing"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusSuspended SubscriptionStatus = "suspended"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
)

// Invoice Status Values
const (
	InvoiceStatusDraft    InvoiceStatus = "draft"
	InvoiceStatusPending  InvoiceStatus = "pending"
	InvoiceStatusPaid     InvoiceStatus = "paid"
	InvoiceStatusOverdue  InvoiceStatus = "overdue"
	InvoiceStatusVoid     InvoiceStatus = "void"
	InvoiceStatusRefunded InvoiceStatus = "refunded"
)

// Tracking Status Values
const (
	TrackingStatusPending   TrackingStatus = "pending"
	TrackingStatusPicked    TrackingStatus = "picked"
	TrackingStatusInTransit TrackingStatus = "in_transit"
	TrackingStatusDelivered TrackingStatus = "delivered"
	TrackingStatusFailed    TrackingStatus = "failed"
	TrackingStatusReturned  TrackingStatus = "returned"
)

// Fulfillment Status Values
const (
	FulfillmentStatusPending   FulfillmentStatus = "pending"
	FulfillmentStatusPicked    FulfillmentStatus = "picked"
	FulfillmentStatusPacked    FulfillmentStatus = "packed"
	FulfillmentStatusShipped   FulfillmentStatus = "shipped"
	FulfillmentStatusDelivered FulfillmentStatus = "delivered"
)

// Category Status Values
const (
	CategoryStatusActive   CategoryStatus = "active"
	CategoryStatusInactive CategoryStatus = "inactive"
	CategoryStatusArchived CategoryStatus = "archived"
)

// Page Status Values
const (
	PageStatusDraft     PageStatus = "draft"
	PageStatusPublished PageStatus = "published"
	PageStatusArchived  PageStatus = "archived"
	PageStatusScheduled PageStatus = "scheduled"
)

// Login Attempt Status Values
const (
	LoginAttemptStatusSuccess LoginAttemptStatus = "success"
	LoginAttemptStatusFailed  LoginAttemptStatus = "failed"
	LoginAttemptStatusBlocked LoginAttemptStatus = "blocked"
)

// Device Status Values
const (
	DeviceStatusTrusted DeviceStatus = "trusted"
	DeviceStatusBlocked DeviceStatus = "blocked"
)

// Alert Status Values
const (
	AlertStatusActive   AlertStatus = "active"
	AlertStatusResolved AlertStatus = "resolved"
	AlertStatusIgnored  AlertStatus = "ignored"
)

// Dunning Action Status Values
const (
	DunningActionStatusPending   DunningActionStatus = "pending"
	DunningActionStatusCompleted DunningActionStatus = "completed"
	DunningActionStatusFailed    DunningActionStatus = "failed"
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

// Validation functions for new status types
func IsValidComponentStatus(status string) bool {
	return status == string(ComponentStatusDraft) ||
		status == string(ComponentStatusActive) ||
		status == string(ComponentStatusInactive) ||
		status == string(ComponentStatusArchived)
}

func IsValidReviewStatus(status string) bool {
	return status == string(ReviewStatusPending) ||
		status == string(ReviewStatusApproved) ||
		status == string(ReviewStatusRejected)
}

func IsValidNotificationStatus(status string) bool {
	return status == string(NotificationStatusPending) ||
		status == string(NotificationStatusSent) ||
		status == string(NotificationStatusDelivered) ||
		status == string(NotificationStatusFailed)
}

func IsValidWebhookStatus(status string) bool {
	return status == string(WebhookStatusPending) ||
		status == string(WebhookStatusSent) ||
		status == string(WebhookStatusDelivered) ||
		status == string(WebhookStatusFailed)
}

func IsValidIntegrationStatus(status string) bool {
	return status == string(IntegrationStatusPending) ||
		status == string(IntegrationStatusActive) ||
		status == string(IntegrationStatusFailed)
}

func IsValidSyncStatus(status string) bool {
	return status == string(SyncStatusPending) ||
		status == string(SyncStatusRunning) ||
		status == string(SyncStatusCompleted) ||
		status == string(SyncStatusFailed)
}

func IsValidCartStatus(status string) bool {
	return status == string(CartStatusActive) ||
		status == string(CartStatusAbandoned) ||
		status == string(CartStatusExpired)
}

func IsValidReturnStatus(status string) bool {
	return status == string(ReturnStatusPending) ||
		status == string(ReturnStatusApproved) ||
		status == string(ReturnStatusProcessing) ||
		status == string(ReturnStatusCompleted) ||
		status == string(ReturnStatusCancelled)
}

func IsValidReferralStatus(status string) bool {
	return status == string(ReferralStatusPending) ||
		status == string(ReferralStatusActive) ||
		status == string(ReferralStatusCompleted) ||
		status == string(ReferralStatusExpired) ||
		status == string(ReferralStatusCancelled)
}

func IsValidCommissionStatus(status string) bool {
	return status == string(CommissionStatusPending) ||
		status == string(CommissionStatusApproved) ||
		status == string(CommissionStatusPaid) ||
		status == string(CommissionStatusFailed)
}

func IsValidFinanceStatus(status string) bool {
	return status == string(FinanceStatusPending) ||
		status == string(FinanceStatusProcessing) ||
		status == string(FinanceStatusCompleted) ||
		status == string(FinanceStatusFailed)
}

func IsValidPayoutStatus(status string) bool {
	return status == string(PayoutStatusPending) ||
		status == string(PayoutStatusProcessing) ||
		status == string(PayoutStatusCompleted) ||
		status == string(PayoutStatusFailed) ||
		status == string(PayoutStatusCancelled)
}

func IsValidCampaignStatus(status string) bool {
	return status == string(CampaignStatusDraft) ||
		status == string(CampaignStatusScheduled) ||
		status == string(CampaignStatusRunning) ||
		status == string(CampaignStatusCompleted) ||
		status == string(CampaignStatusCancelled)
}

func IsValidEmailStatus(status string) bool {
	return status == string(EmailStatusPending) ||
		status == string(EmailStatusSent) ||
		status == string(EmailStatusDelivered) ||
		status == string(EmailStatusOpened) ||
		status == string(EmailStatusFailed)
}

func IsValidSegmentStatus(status string) bool {
	return status == string(SegmentStatusActive) ||
		status == string(SegmentStatusInactive)
}

func IsValidSubscriptionStatus(status string) bool {
	return status == string(SubscriptionStatusActive) ||
		status == string(SubscriptionStatusPending) ||
		status == string(SubscriptionStatusTrialing) ||
		status == string(SubscriptionStatusCancelled) ||
		status == string(SubscriptionStatusSuspended) ||
		status == string(SubscriptionStatusExpired)
}

func IsValidInvoiceStatus(status string) bool {
	return status == string(InvoiceStatusDraft) ||
		status == string(InvoiceStatusPending) ||
		status == string(InvoiceStatusPaid) ||
		status == string(InvoiceStatusOverdue) ||
		status == string(InvoiceStatusVoid) ||
		status == string(InvoiceStatusRefunded)
}

func IsValidTrackingStatus(status string) bool {
	return status == string(TrackingStatusPending) ||
		status == string(TrackingStatusPicked) ||
		status == string(TrackingStatusInTransit) ||
		status == string(TrackingStatusDelivered) ||
		status == string(TrackingStatusFailed) ||
		status == string(TrackingStatusReturned)
}

func IsValidFulfillmentStatus(status string) bool {
	return status == string(FulfillmentStatusPending) ||
		status == string(FulfillmentStatusPicked) ||
		status == string(FulfillmentStatusPacked) ||
		status == string(FulfillmentStatusShipped) ||
		status == string(FulfillmentStatusDelivered)
}

func IsValidCategoryStatus(status string) bool {
	return status == string(CategoryStatusActive) ||
		status == string(CategoryStatusInactive) ||
		status == string(CategoryStatusArchived)
}

func IsValidPageStatus(status string) bool {
	return status == string(PageStatusDraft) ||
		status == string(PageStatusPublished) ||
		status == string(PageStatusArchived) ||
		status == string(PageStatusScheduled)
}

func IsValidLoginAttemptStatus(status string) bool {
	return status == string(LoginAttemptStatusSuccess) ||
		status == string(LoginAttemptStatusFailed) ||
		status == string(LoginAttemptStatusBlocked)
}

func IsValidDeviceStatus(status string) bool {
	return status == string(DeviceStatusTrusted) ||
		status == string(DeviceStatusBlocked)
}

func IsValidAlertStatus(status string) bool {
	return status == string(AlertStatusActive) ||
		status == string(AlertStatusResolved) ||
		status == string(AlertStatusIgnored)
}

func IsValidDunningActionStatus(status string) bool {
	return status == string(DunningActionStatusPending) ||
		status == string(DunningActionStatusCompleted) ||
		status == string(DunningActionStatusFailed)
}