package database

import (
	"github.com/esass/internal/analytics"
	"github.com/esass/internal/content"
	"github.com/esass/internal/notification"
	"github.com/esass/internal/order"
	"github.com/esass/internal/payment"
	"github.com/esass/internal/product"
	"github.com/esass/internal/shipping"
	"github.com/esass/internal/tenant"
	"github.com/esass/internal/user"
	"github.com/esass/internal/webhook"
	"gorm.io/gorm"
)

// TODO: Implement database migrations
// This will handle:
// - Table creation and updates
// - Data migrations
// - Schema versioning

func AutoMigrate(db *gorm.DB) error {
	// Migrate all entities
	return db.AutoMigrate(
		// Tenant entities
		&tenant.Tenant{},
		
		// Product entities
		&product.Product{},
		&product.Category{},
		&product.ProductVariant{},
		
		// Order entities
		&order.Order{},
		&order.OrderItem{},
		&order.ShippingAddress{},
		
		// User entities
		&user.User{},
		&user.UserSession{},
		&user.UserProfile{},
		
		// Analytics entities
		&analytics.AnalyticsEvent{},
		&analytics.PageView{},
		&analytics.ProductView{},
		&analytics.Purchase{},
		
		// Payment entities
		&payment.Payment{},
		&payment.PaymentMethod{},
		&payment.Refund{},
		&payment.PaymentHistory{},
		
		// Notification entities
		&notification.Notification{},
		&notification.NotificationTemplate{},
		&notification.NotificationPreference{},
		&notification.NotificationLog{},
		
		// Shipping entities
		&shipping.ShippingZone{},
		&shipping.ShippingZoneCountry{},
		&shipping.ShippingRate{},
		&shipping.ShippingLabel{},
		&shipping.ShippingTracking{},
		&shipping.ShippingProvider{},
		
		// Webhook entities
		&webhook.WebhookEndpoint{},
		&webhook.WebhookDelivery{},
		&webhook.WebhookIncoming{},
		&webhook.WebhookRateLimit{},
		
		// Content entities
		&content.Page{},
		&content.Media{},
		&content.Menu{},
		&content.MenuItem{},
		&content.Tag{},
		&content.Category{},
		&content.SEOSettings{},
	)
}
