package notification

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-saas/internal/shared/testhelpers"
)

// Integration tests for notification service - critical for e-commerce operations

func TestNotificationService_NotificationLifecycle(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&Notification{}, &NotificationTemplate{}, &NotificationPreference{}, &NotificationLog{})
	require.NoError(t, err)

	// Setup services
	notificationRepo := NewRepository(testDB.DB)
	notificationService := NewService(notificationRepo)

	tenantID := uuid.New()
	userID := uuid.New()

	t.Run("Send notification with template", func(t *testing.T) {
		// Step 1: Create email template
		templateReq := &CreateTemplateRequest{
			Name:    "Order Confirmation",
			Type:    TypeEmail,
			Channel: ChannelOrderConfirmation,
			Subject: "Order Confirmed - #{{order_number}}",
			Content: "Hello {{customer_name}}, your order #{{order_number}} has been confirmed. Total: {{total_amount}}",
			Variables: map[string]interface{}{
				"customer_name": "string",
				"order_number":  "string",
				"total_amount":  "string",
			},
		}

		template, err := notificationService.CreateTemplate(tenantID, templateReq)
		require.NoError(t, err)
		assert.Equal(t, templateReq.Name, template.Name)
		assert.Equal(t, templateReq.Type, template.Type)
		assert.Equal(t, templateReq.Channel, template.Channel)
		assert.True(t, template.IsActive)

		// Step 2: Send notification using template
		sendReq := &SendNotificationRequest{
			Type:       TypeEmail,
			Channel:    ChannelOrderConfirmation,
			Recipients: []string{"customer@example.com"},
			TemplateID: template.ID.String(),
			Variables: map[string]interface{}{
				"customer_name": "John Doe",
				"order_number":  "ORD-12345",
				"total_amount":  "৳2,500",
			},
			Priority: PriorityHigh,
			UserID:   userID.String(),
		}

		response, err := notificationService.SendNotification(tenantID, sendReq)
		require.NoError(t, err)
		assert.NotEmpty(t, response.NotificationIDs)
		assert.Equal(t, "queued", response.Status)

		// Wait for async notification processing to complete
		time.Sleep(100 * time.Millisecond)

		// Step 3: Verify notification was created
		notificationID, err := uuid.Parse(response.NotificationIDs[0])
		require.NoError(t, err)

		notification, err := notificationService.GetNotification(tenantID, notificationID.String())
		require.NoError(t, err)
		assert.Equal(t, tenantID, notification.TenantID)
		assert.Equal(t, &userID, notification.UserID)
		assert.Equal(t, TypeEmail, notification.Type)
		assert.Equal(t, ChannelOrderConfirmation, notification.Channel)
		assert.Equal(t, "customer@example.com", notification.Recipient)
		assert.Equal(t, PriorityHigh, notification.Priority)
		assert.Contains(t, notification.Subject, "ORD-12345")
		assert.Contains(t, notification.Content, "John Doe")
	})

	t.Run("SMS notification without template", func(t *testing.T) {
		// Send direct SMS notification
		smsReq := &SendSMSRequest{
			To:      []string{"+8801712345678"},
			Message: "Your order has been shipped. Track: bit.ly/track123",
			UserID:  userID.String(),
		}

		err := notificationService.SendSMS(tenantID, smsReq)
		require.NoError(t, err)

		// Wait for async notification processing to complete
		time.Sleep(100 * time.Millisecond)

		// Verify notification was created
		notifications, _, err := notificationService.ListNotifications(tenantID, &userID, 0, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(notifications), 1)

		// Find the SMS notification
		var smsNotification *Notification
		for _, n := range notifications {
			if n.Type == TypeSMS {
				smsNotification = n
				break
			}
		}
		require.NotNil(t, smsNotification)
		assert.Equal(t, "+8801712345678", smsNotification.Recipient)
		assert.Contains(t, smsNotification.Content, "shipped")
	})
}

func TestNotificationService_TemplateManagement(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&NotificationTemplate{})
	require.NoError(t, err)

	// Setup services
	notificationRepo := NewRepository(testDB.DB)
	notificationService := NewService(notificationRepo)

	tenantID := uuid.New()

	t.Run("Template CRUD operations", func(t *testing.T) {
		// Create password reset template
		createReq := &CreateTemplateRequest{
			Name:    "Password Reset",
			Type:    TypeEmail,
			Channel: ChannelPasswordReset,
			Subject: "Reset Your Password",
			Content: "Click here to reset: {{reset_link}}. Link expires in {{expiry_minutes}} minutes.",
			Variables: map[string]interface{}{
				"reset_link":     "string",
				"expiry_minutes": "number",
			},
		}

		template, err := notificationService.CreateTemplate(tenantID, createReq)
		require.NoError(t, err)
		assert.Equal(t, createReq.Name, template.Name)
		assert.Equal(t, createReq.Type, template.Type)
		assert.Equal(t, createReq.Channel, template.Channel)

		// Get template
		retrievedTemplate, err := notificationService.GetTemplate(tenantID, template.ID.String())
		require.NoError(t, err)
		assert.Equal(t, template.ID, retrievedTemplate.ID)
		assert.Equal(t, template.Content, retrievedTemplate.Content)

		// Update template
		updateReq := &UpdateTemplateRequest{
			Content: "New reset link: {{reset_link}}. Expires in {{expiry_minutes}} minutes. Contact support if needed.",
		}

		err = notificationService.UpdateTemplate(tenantID, template.ID.String(), updateReq)
		require.NoError(t, err)

		// Verify update
		updatedTemplate, err := notificationService.GetTemplate(tenantID, template.ID.String())
		require.NoError(t, err)
		assert.Equal(t, updateReq.Content, updatedTemplate.Content)

		// List templates
		templates, err := notificationService.ListTemplates(tenantID, TypeEmail, "")
		require.NoError(t, err)
		assert.Len(t, templates, 1)
		assert.Equal(t, template.ID, templates[0].ID)

		// List by channel
		channelTemplates, err := notificationService.ListTemplates(tenantID, "", ChannelPasswordReset)
		require.NoError(t, err)
		assert.Len(t, channelTemplates, 1)
		assert.Equal(t, ChannelPasswordReset, channelTemplates[0].Channel)
	})

	t.Run("Multiple template types", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create email template
		emailTemplate := &CreateTemplateRequest{
			Name:    "Welcome Email",
			Type:    TypeEmail,
			Channel: ChannelWelcome,
			Subject: "Welcome to {{store_name}}!",
			Content: "Hi {{customer_name}}, welcome to our store!",
		}

		_, err := notificationService.CreateTemplate(tenantID, emailTemplate)
		require.NoError(t, err)

		// Create SMS template
		smsTemplate := &CreateTemplateRequest{
			Name:    "Welcome SMS",
			Type:    TypeSMS,
			Channel: ChannelWelcome,
			Content: "Welcome to {{store_name}}! Use code WELCOME10 for 10% off.",
		}

		_, err = notificationService.CreateTemplate(tenantID, smsTemplate)
		require.NoError(t, err)

		// List all templates
		allTemplates, err := notificationService.ListTemplates(tenantID, "", "")
		require.NoError(t, err)
		assert.Len(t, allTemplates, 2)

		// Filter by type
		emailTemplates, err := notificationService.ListTemplates(tenantID, TypeEmail, "")
		require.NoError(t, err)
		assert.Len(t, emailTemplates, 1)
		assert.Equal(t, TypeEmail, emailTemplates[0].Type)

		smsTemplates, err := notificationService.ListTemplates(tenantID, TypeSMS, "")
		require.NoError(t, err)
		assert.Len(t, smsTemplates, 1)
		assert.Equal(t, TypeSMS, smsTemplates[0].Type)
	})
}

func TestNotificationService_UserPreferences(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&NotificationPreference{})
	require.NoError(t, err)

	// Setup services
	notificationRepo := NewRepository(testDB.DB)
	notificationService := NewService(notificationRepo)

	tenantID := uuid.New()
	userID := uuid.New()

	t.Run("Preference management", func(t *testing.T) {
		// Get default preferences for marketing channel (should return defaults)
		prefs, err := notificationService.GetPreferences(tenantID, userID, ChannelMarketing)
		require.NoError(t, err)
		// Default values should be set
		assert.True(t, prefs.EmailEnabled) // Default
		assert.False(t, prefs.SMSEnabled)  // Default

		// Update marketing preferences
		marketingReq := &NotificationPreferenceRequest{
			Channel:      ChannelMarketing,
			EmailEnabled: boolPtr(false),
			SMSEnabled:   boolPtr(true),
		}

		err = notificationService.UpdatePreferences(tenantID, userID, marketingReq)
		require.NoError(t, err)

		// Verify marketing channel update
		updatedPrefs, err := notificationService.GetPreferences(tenantID, userID, ChannelMarketing)
		require.NoError(t, err)
		t.Logf("Marketing prefs: EmailEnabled=%v, SMSEnabled=%v", updatedPrefs.EmailEnabled, updatedPrefs.SMSEnabled)
		assert.False(t, updatedPrefs.EmailEnabled)
		assert.True(t, updatedPrefs.SMSEnabled)

		// Update order notification preferences
		orderReq := &NotificationPreferenceRequest{
			Channel:      ChannelOrderConfirmation,
			EmailEnabled: boolPtr(true),
			SMSEnabled:   boolPtr(true),
			PushEnabled:  boolPtr(false),
		}

		err = notificationService.UpdatePreferences(tenantID, userID, orderReq)
		require.NoError(t, err)

		// Verify order confirmation channel preferences
		orderPrefs, err := notificationService.GetPreferences(tenantID, userID, ChannelOrderConfirmation)
		require.NoError(t, err)
		t.Logf("Order prefs: EmailEnabled=%v, SMSEnabled=%v, PushEnabled=%v", orderPrefs.EmailEnabled, orderPrefs.SMSEnabled, orderPrefs.PushEnabled)
		assert.True(t, orderPrefs.EmailEnabled)
		assert.True(t, orderPrefs.SMSEnabled)
		assert.False(t, orderPrefs.PushEnabled)
	})
}

func TestNotificationService_Statistics(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&Notification{}, &NotificationLog{})
	require.NoError(t, err)

	// Setup services
	notificationRepo := NewRepository(testDB.DB)
	notificationService := NewService(notificationRepo)

	tenantID := uuid.New()

	t.Run("Notification statistics", func(t *testing.T) {
		// Create sample notifications with different statuses
		notifications := []*Notification{
			{
				ID:        uuid.New(),
				TenantID:  tenantID,
				Type:      TypeEmail,
				Channel:   ChannelOrderConfirmation,
				Recipient: "test1@example.com",
				Status:    StatusSent,
				Content:   "Order confirmed",
				Metadata:  "{}",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				TenantID:  tenantID,
				Type:      TypeEmail,
				Channel:   ChannelOrderConfirmation,
				Recipient: "test2@example.com",
				Status:    StatusDelivered,
				Content:   "Order confirmed",
				Metadata:  "{}",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				TenantID:  tenantID,
				Type:      TypeSMS,
				Channel:   ChannelShippingUpdate,
				Recipient: "+8801712345678",
				Status:    StatusFailed,
				Content:   "Shipped",
				Metadata:  "{}",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		for _, notification := range notifications {
			err := notificationRepo.Create(notification)
			require.NoError(t, err)
		}

		// Get statistics
		stats, err := notificationService.GetStats(tenantID)
		require.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(2), stats.TotalSent)     // sent + delivered
		assert.Equal(t, int64(1), stats.TotalDelivered) // only delivered
		assert.Equal(t, int64(1), stats.TotalFailed)   // failed
		assert.Greater(t, stats.DeliveryRate, float64(0))
		assert.Greater(t, stats.FailureRate, float64(0))
	})
}

func TestNotificationService_MultiTenantIsolation(t *testing.T) {
	// Setup test database
	testDB := testhelpers.SetupSimpleTestDatabase(t)
	defer testDB.TeardownTestDatabase(t)

	// Migrate schemas
	err := testDB.DB.AutoMigrate(&Notification{}, &NotificationTemplate{})
	require.NoError(t, err)

	// Setup services
	notificationRepo := NewRepository(testDB.DB)
	notificationService := NewService(notificationRepo)

	tenant1ID := uuid.New()
	tenant2ID := uuid.New()

	t.Run("Notification isolation between tenants", func(t *testing.T) {
		// Send notification for tenant 1
		tenant1Req := &SendNotificationRequest{
			Type:       TypeEmail,
			Channel:    ChannelOrderConfirmation,
			Recipients: []string{"tenant1@example.com"},
			Subject:    "Tenant 1 Order",
			Content:    "Your order from Tenant 1 store",
			Priority:   PriorityNormal,
		}

		tenant1Response, err := notificationService.SendNotification(tenant1ID, tenant1Req)
		require.NoError(t, err)
		assert.NotEmpty(t, tenant1Response.NotificationIDs)

		// Send notification for tenant 2
		tenant2Req := &SendNotificationRequest{
			Type:       TypeEmail,
			Channel:    ChannelOrderConfirmation,
			Recipients: []string{"tenant2@example.com"},
			Subject:    "Tenant 2 Order",
			Content:    "Your order from Tenant 2 store",
			Priority:   PriorityHigh,
		}

		tenant2Response, err := notificationService.SendNotification(tenant2ID, tenant2Req)
		require.NoError(t, err)
		assert.NotEmpty(t, tenant2Response.NotificationIDs)

		// Verify tenant 1 can only see their notifications
		tenant1Notifications, _, err := notificationService.ListNotifications(tenant1ID, nil, 0, 10)
		require.NoError(t, err)
		assert.Len(t, tenant1Notifications, 1)
		assert.Equal(t, tenant1ID, tenant1Notifications[0].TenantID)
		assert.Contains(t, tenant1Notifications[0].Subject, "Tenant 1")

		// Verify tenant 2 can only see their notifications
		tenant2Notifications, _, err := notificationService.ListNotifications(tenant2ID, nil, 0, 10)
		require.NoError(t, err)
		assert.Len(t, tenant2Notifications, 1)
		assert.Equal(t, tenant2ID, tenant2Notifications[0].TenantID)
		assert.Contains(t, tenant2Notifications[0].Subject, "Tenant 2")

		// Verify cross-tenant access is blocked
		tenant1NotificationID, _ := uuid.Parse(tenant1Response.NotificationIDs[0])
		_, err = notificationService.GetNotification(tenant2ID, tenant1NotificationID.String())
		assert.Error(t, err) // Should not be able to access other tenant's notification
	})

	t.Run("Template isolation between tenants", func(t *testing.T) {
		// Create template for tenant 1
		tenant1Template := &CreateTemplateRequest{
			Name:    "Tenant 1 Welcome",
			Type:    TypeEmail,
			Channel: ChannelWelcome,
			Subject: "Welcome to Tenant 1 Store",
			Content: "Welcome to our exclusive Tenant 1 store!",
		}

		template1, err := notificationService.CreateTemplate(tenant1ID, tenant1Template)
		require.NoError(t, err)

		// Create template for tenant 2
		tenant2Template := &CreateTemplateRequest{
			Name:    "Tenant 2 Welcome",
			Type:    TypeEmail,
			Channel: ChannelWelcome,
			Subject: "Welcome to Tenant 2 Store",
			Content: "Welcome to our amazing Tenant 2 store!",
		}

		template2, err := notificationService.CreateTemplate(tenant2ID, tenant2Template)
		require.NoError(t, err)

		// Verify tenant 1 can only see their templates
		tenant1Templates, err := notificationService.ListTemplates(tenant1ID, "", "")
		require.NoError(t, err)
		assert.Len(t, tenant1Templates, 1)
		assert.Equal(t, template1.ID, tenant1Templates[0].ID)
		assert.Contains(t, tenant1Templates[0].Subject, "Tenant 1")

		// Verify tenant 2 can only see their templates
		tenant2Templates, err := notificationService.ListTemplates(tenant2ID, "", "")
		require.NoError(t, err)
		assert.Len(t, tenant2Templates, 1)
		assert.Equal(t, template2.ID, tenant2Templates[0].ID)
		assert.Contains(t, tenant2Templates[0].Subject, "Tenant 2")

		// Verify cross-tenant template access is blocked
		_, err = notificationService.GetTemplate(tenant2ID, template1.ID.String())
		assert.Error(t, err) // Should not be able to access other tenant's template
	})
}

// Helper function
func boolPtr(b bool) *bool {
	return &b
}