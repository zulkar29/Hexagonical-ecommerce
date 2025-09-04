package payment

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Payment statuses
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusSucceeded  = "succeeded"
	StatusFailed     = "failed"
	StatusCancelled  = "cancelled"
	StatusRefunded   = "refunded"
)

// Payment gateways
const (
	GatewaySSLCommerz = "sslcommerz"
	GatewayBKash      = "bkash"
	GatewayNagad      = "nagad"
	GatewayStripe     = "stripe"
	GatewayPayPal     = "paypal"
)

type Payment struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID          uuid.UUID  `json:"tenant_id" gorm:"type:uuid;not null;index"`
	OrderID           uuid.UUID  `json:"order_id" gorm:"type:uuid;not null;index"`
	UserID            uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	PaymentIntentID   string     `json:"payment_intent_id" gorm:"size:255"`
	PaymentMethodID   string     `json:"payment_method_id" gorm:"size:255"`
	Amount            float64    `json:"amount" gorm:"not null"`
	Currency          string     `json:"currency" gorm:"size:3;not null;default:'BDT'"`
	Status            string     `json:"status" gorm:"size:50;not null;default:'pending'"`
	Gateway           string     `json:"gateway" gorm:"size:50;not null"`
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
	Provider     string    `json:"provider" gorm:"size:50;not null"`     // sslcommerz, bkash, nagad, stripe, paypal
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
	Currency        string    `json:"currency" gorm:"size:3;not null;default:'BDT'"`
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

// SSLCommerz Specific Structures
type SSLCommerzPaymentRequest struct {
	StoreID              string  `json:"store_id"`
	StorePassword        string  `json:"store_passwd"`
	TotalAmount          float64 `json:"total_amount"`
	Currency             string  `json:"currency"`
	TransactionID        string  `json:"tran_id"`
	SuccessURL           string  `json:"success_url"`
	FailURL              string  `json:"fail_url"`
	CancelURL            string  `json:"cancel_url"`
	IPNListenerURL       string  `json:"ipn_url"`
	CustomerName         string  `json:"cus_name"`
	CustomerEmail        string  `json:"cus_email"`
	CustomerPhone        string  `json:"cus_phone"`
	CustomerAddress1     string  `json:"cus_add1"`
	CustomerCity         string  `json:"cus_city"`
	CustomerState        string  `json:"cus_state"`
	CustomerPostcode     string  `json:"cus_postcode"`
	CustomerCountry      string  `json:"cus_country"`
	ShippingMethodName   string  `json:"shipping_method"`
	ProductName          string  `json:"product_name"`
	ProductCategory      string  `json:"product_category"`
	ProductProfile       string  `json:"product_profile"`
}

type SSLCommerzPaymentResponse struct {
	Status          string `json:"status"`
	FailedReason    string `json:"failedreason"`
	SessionKey      string `json:"sessionkey"`
	GatewayPageURL  string `json:"GatewayPageURL"`
	StoreBanner     string `json:"storeBanner"`
	StoreLogo       string `json:"storeLogo"`
	IsMobileAPIFail string `json:"is_mobile_api_fail"`
	Description     string `json:"desc"`
}

type SSLCommerzIPNResponse struct {
	Status           string  `json:"status"`
	TransactionID    string  `json:"tran_id"`
	ValID            string  `json:"val_id"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	StoreAmount      float64 `json:"store_amount"`
	CardType         string  `json:"card_type"`
	CardNo           string  `json:"card_no"`
	BankTransactionID string `json:"bank_tran_id"`
	TransactionDate  string  `json:"tran_date"`
	Error            string  `json:"error"`
	CurrencyType     string  `json:"currency_type"`
	CurrencyAmount   float64 `json:"currency_amount"`
	CurrencyRate     float64 `json:"currency_rate"`
	BaseAmount       float64 `json:"base_amount"`
	ValueA           string  `json:"value_a"`
	ValueB           string  `json:"value_b"`
	ValueC           string  `json:"value_c"`
	ValueD           string  `json:"value_d"`
	VerifySign       string  `json:"verify_sign"`
	VerifyKey        string  `json:"verify_key"`
	RiskLevel        string  `json:"risk_level"`
	RiskTitle        string  `json:"risk_title"`
}

// Request/Response Structures
type CreatePaymentRequest struct {
	OrderID         string  `json:"order_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,min=0.01"`
	Currency        string  `json:"currency" validate:"required,len=3"`
	Gateway         string  `json:"gateway" validate:"required"`
	PaymentMethodID string  `json:"payment_method_id,omitempty"`
	CustomerEmail   string  `json:"customer_email" validate:"required,email"`
	CustomerPhone   string  `json:"customer_phone,omitempty"`
	ReturnURL       string  `json:"return_url,omitempty"`
}

type CreatePaymentResponse struct {
	PaymentID      string `json:"payment_id"`
	Status         string `json:"status"`
	PaymentURL     string `json:"payment_url,omitempty"`
	SessionKey     string `json:"session_key,omitempty"`
	GatewayPageURL string `json:"gateway_page_url,omitempty"`
}

type ProcessPaymentRequest struct {
	PaymentID       string                 `json:"payment_id" validate:"required"`
	Gateway         string                 `json:"gateway" validate:"required"`
	GatewayResponse map[string]interface{} `json:"gateway_response"`
}

type RefundPaymentRequest struct {
	PaymentID string  `json:"payment_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,min=0.01"`
	Reason    string  `json:"reason,omitempty"`
}
