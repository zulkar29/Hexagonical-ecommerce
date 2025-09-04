package payment

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: Implement payment entities
// This will handle:
// - Payment processing
// - Transaction management
// - Payment gateway integration
// - Refund handling

type Payment struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID          uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	OrderID           uuid.UUID  `json:"order_id" gorm:"type:uuid;not null;index"`
	UserID            uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	PaymentIntentID   string     `json:"payment_intent_id" gorm:"size:255"`
	PaymentMethodID   string     `json:"payment_method_id" gorm:"size:255"`
	Amount            float64    `json:"amount" gorm:"not null"`
	Currency          string     `json:"currency" gorm:"size:3;not null;default:'USD'"`
	Status            string     `json:"status" gorm:"size:50;not null;default:'pending'"` // pending, processing, succeeded, failed, cancelled, refunded
	Gateway           string     `json:"gateway" gorm:"size:50;not null"`                 // stripe, paypal, square, etc.
	GatewayResponse   string     `json:"gateway_response" gorm:"type:text"`
	FailureReason     string     `json:"failure_reason" gorm:"type:text"`
	RefundedAmount    float64    `json:"refunded_amount" gorm:"default:0"`
	RefundedAt        *time.Time `json:"refunded_at"`
	ProcessedAt       *time.Time `json:"processed_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type PaymentMethod struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID     uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Type         string    `json:"type" gorm:"size:50;not null"`         // card, bank_account, digital_wallet
	Provider     string    `json:"provider" gorm:"size:50;not null"`     // stripe, paypal, apple_pay, google_pay
	ProviderID   string    `json:"provider_id" gorm:"size:255;not null"` // External provider's payment method ID
	Last4        string    `json:"last4" gorm:"size:4"`                  // Last 4 digits for cards
	Brand        string    `json:"brand" gorm:"size:50"`                 // visa, mastercard, amex, etc.
	ExpiryMonth  int       `json:"expiry_month"`
	ExpiryYear   int       `json:"expiry_year"`
	IsDefault    bool      `json:"is_default" gorm:"default:false"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Refund struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID        uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	PaymentID       uuid.UUID `json:"payment_id" gorm:"type:uuid;not null;index"`
	OrderID         uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
	Amount          float64   `json:"amount" gorm:"not null"`
	Currency        string    `json:"currency" gorm:"size:3;not null;default:'USD'"`
	Reason          string    `json:"reason" gorm:"size:255"`
	Status          string    `json:"status" gorm:"size:50;not null;default:'pending'"` // pending, succeeded, failed
	RefundID        string    `json:"refund_id" gorm:"size:255"`                        // Gateway refund ID
	GatewayResponse string    `json:"gateway_response" gorm:"type:text"`
	ProcessedAt     *time.Time `json:"processed_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type PaymentHistory struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID  uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	PaymentID uuid.UUID `json:"payment_id" gorm:"type:uuid;not null;index"`
	Status    string    `json:"status" gorm:"size:50;not null"`
	Notes     string    `json:"notes" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}
