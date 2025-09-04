package shipping

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: Implement shipping entities
// This will handle:
// - Bangladesh shipping providers (Pathao, RedX, Paperfly, SA Paribahan)
// - International shipping (DHL, FedEx, UPS)
// - Shipping zone management
// - Rate calculation and delivery tracking

type ShippingProvider string
type ShippingMethod string
type DeliveryType string

const (
	// Bangladesh Local Providers
	ProviderPathao     ShippingProvider = "pathao"
	ProviderRedX       ShippingProvider = "redx"
	ProviderPaperfly   ShippingProvider = "paperfly"
	ProviderSAParibahan ShippingProvider = "sa_paribahan"
	ProviderSteadfast  ShippingProvider = "steadfast"
	ProviderEcourier   ShippingProvider = "ecourier"
	
	// International Providers
	ProviderDHL        ShippingProvider = "dhl"
	ProviderFedEx      ShippingProvider = "fedex"
	ProviderUPS        ShippingProvider = "ups"
	ProviderBanglapost ShippingProvider = "banglapost"
)

const (
	MethodStandard ShippingMethod = "standard"
	MethodExpress  ShippingMethod = "express"
	MethodSameDay  ShippingMethod = "same_day"
	MethodPickup   ShippingMethod = "pickup"
)

const (
	DeliveryHome   DeliveryType = "home"
	DeliveryOffice DeliveryType = "office"
	DeliveryHub    DeliveryType = "hub"
	DeliveryPickup DeliveryType = "pickup_point"
)

type ShippingZone struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	// Zone coverage
	Countries []ShippingZoneCountry `json:"countries" gorm:"foreignKey:ZoneID"`
	Rates     []ShippingRate       `json:"rates" gorm:"foreignKey:ZoneID"`
}

type ShippingZoneCountry struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ZoneID    uuid.UUID `json:"zone_id" gorm:"type:uuid;not null;index"`
	Country   string    `json:"country" gorm:"size:2;not null"` // ISO country code
	States    []string  `json:"states" gorm:"serializer:json"`  // States/divisions within country
	Cities    []string  `json:"cities" gorm:"serializer:json"`  // Specific cities
	PostalCodes []string `json:"postal_codes" gorm:"serializer:json"` // Postal code ranges
	CreatedAt time.Time `json:"created_at"`
}

type ShippingRate struct {
	ID               uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID         uuid.UUID        `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ZoneID           uuid.UUID        `json:"zone_id" gorm:"type:uuid;not null;index"`
	Provider         ShippingProvider `json:"provider" gorm:"size:50;not null"`
	Method           ShippingMethod   `json:"method" gorm:"size:50;not null"`
	Name             string           `json:"name" gorm:"size:100;not null"`
	Description      string           `json:"description" gorm:"type:text"`
	BaseRate         float64          `json:"base_rate" gorm:"not null"`
	WeightRate       float64          `json:"weight_rate" gorm:"default:0"` // Per kg
	VolumeRate       float64          `json:"volume_rate" gorm:"default:0"` // Per cubic cm
	MinWeight        float64          `json:"min_weight" gorm:"default:0"`
	MaxWeight        float64          `json:"max_weight" gorm:"default:0"` // 0 = unlimited
	FreeShippingMin  float64          `json:"free_shipping_min" gorm:"default:0"` // Free shipping threshold
	EstimatedDays    int              `json:"estimated_days" gorm:"default:1"`
	IsActive         bool             `json:"is_active" gorm:"default:true"`
	Priority         int              `json:"priority" gorm:"default:0"` // Display order
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	DeletedAt        gorm.DeletedAt   `json:"deleted_at" gorm:"index"`
}

type ShippingLabel struct {
	ID           uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID     uuid.UUID        `json:"tenant_id" gorm:"type:uuid;not null;index"`
	OrderID      uuid.UUID        `json:"order_id" gorm:"type:uuid;not null;index"`
	Provider     ShippingProvider `json:"provider" gorm:"size:50;not null"`
	TrackingNumber string         `json:"tracking_number" gorm:"size:100;not null;index"`
	LabelURL     string           `json:"label_url" gorm:"size:500"`
	Cost         float64          `json:"cost" gorm:"not null"`
	Currency     string           `json:"currency" gorm:"size:3;not null;default:'BDT'"`
	Status       string           `json:"status" gorm:"size:50;not null;default:'created'"` // created, printed, shipped, delivered, failed
	
	// Provider specific data
	ProviderOrderID   string `json:"provider_order_id" gorm:"size:100"`
	ProviderResponse  string `json:"provider_response" gorm:"type:json"`
	
	// Delivery details
	EstimatedDelivery *time.Time `json:"estimated_delivery"`
	ActualDelivery    *time.Time `json:"actual_delivery"`
	DeliveredTo       string     `json:"delivered_to" gorm:"size:255"`
	DeliveryNotes     string     `json:"delivery_notes" gorm:"type:text"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type ShippingTracking struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LabelID      uuid.UUID `json:"label_id" gorm:"type:uuid;not null;index"`
	Status       string    `json:"status" gorm:"size:50;not null"`
	Description  string    `json:"description" gorm:"type:text"`
	Location     string    `json:"location" gorm:"size:255"`
	Timestamp    time.Time `json:"timestamp" gorm:"not null"`
	IsDelivered  bool      `json:"is_delivered" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
}

type ShippingProvider struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID       uuid.UUID `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Provider       ShippingProvider `json:"provider" gorm:"size:50;not null"`
	Name           string    `json:"name" gorm:"size:100;not null"`
	IsActive       bool      `json:"is_active" gorm:"default:true"`
	APIKey         string    `json:"api_key" gorm:"size:255"`
	APISecret      string    `json:"api_secret" gorm:"size:255"`
	SandboxMode    bool      `json:"sandbox_mode" gorm:"default:true"`
	
	// Provider specific settings
	Settings map[string]interface{} `json:"settings" gorm:"serializer:json"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Business Logic Methods

// IsInternational checks if shipping zone covers international destinations
func (sz *ShippingZone) IsInternational() bool {
	for _, country := range sz.Countries {
		if country.Country != "BD" { // Bangladesh
			return true
		}
	}
	return false
}

// IsDomestic checks if shipping zone is domestic only
func (sz *ShippingZone) IsDomestic() bool {
	return !sz.IsInternational()
}

// CalculateRate calculates shipping rate based on weight and dimensions
func (sr *ShippingRate) CalculateRate(weight, length, width, height, orderValue float64) float64 {
	// Check free shipping threshold
	if sr.FreeShippingMin > 0 && orderValue >= sr.FreeShippingMin {
		return 0
	}
	
	rate := sr.BaseRate
	
	// Add weight-based charges
	if sr.WeightRate > 0 && weight > 0 {
		rate += weight * sr.WeightRate
	}
	
	// Add volume-based charges
	if sr.VolumeRate > 0 && length > 0 && width > 0 && height > 0 {
		volume := length * width * height
		rate += volume * sr.VolumeRate
	}
	
	return rate
}

// IsEligible checks if shipping rate is eligible for given weight
func (sr *ShippingRate) IsEligible(weight float64) bool {
	if !sr.IsActive {
		return false
	}
	
	if sr.MinWeight > 0 && weight < sr.MinWeight {
		return false
	}
	
	if sr.MaxWeight > 0 && weight > sr.MaxWeight {
		return false
	}
	
	return true
}

// IsDelivered checks if package has been delivered
func (sl *ShippingLabel) IsDelivered() bool {
	return sl.Status == "delivered" && sl.ActualDelivery != nil
}

// GetEstimatedDeliveryDate calculates estimated delivery date
func (sr *ShippingRate) GetEstimatedDeliveryDate() time.Time {
	return time.Now().AddDate(0, 0, sr.EstimatedDays)
}

// TODO: Add provider-specific API integration methods
// - PathaoCreateOrder(orderData) (*ShippingLabel, error)
// - RedXCreateOrder(orderData) (*ShippingLabel, error)
// - PaperflyCreateOrder(orderData) (*ShippingLabel, error)
// - DHLCreateOrder(orderData) (*ShippingLabel, error)
// - TrackPackage(trackingNumber, provider) (*ShippingTracking, error)
