package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"ecommerce-saas/internal/shared/email"
)

type Service interface {
	SendNotification(tenantID uuid.UUID, req *SendNotificationRequest) (*SendNotificationResponse, error)
	SendEmail(tenantID uuid.UUID, req *SendEmailRequest) error
	SendSMS(tenantID uuid.UUID, req *SendSMSRequest) error
	GetNotification(tenantID uuid.UUID, notificationID string) (*Notification, error)
	ListNotifications(tenantID uuid.UUID, userID *uuid.UUID, offset, limit int) ([]*Notification, int64, error)
	MarkAsRead(tenantID uuid.UUID, notificationID string) error
	
	// Template management
	CreateTemplate(tenantID uuid.UUID, req *CreateTemplateRequest) (*NotificationTemplate, error)
	UpdateTemplate(tenantID uuid.UUID, templateID string, req *UpdateTemplateRequest) error
	GetTemplate(tenantID uuid.UUID, templateID string) (*NotificationTemplate, error)
	ListTemplates(tenantID uuid.UUID, notificationType, channel string) ([]*NotificationTemplate, error)
	
	// Preferences
	GetPreferences(tenantID, userID uuid.UUID, channel string) (*NotificationPreference, error)
	UpdatePreferences(tenantID, userID uuid.UUID, req *NotificationPreferenceRequest) error
	
	// Stats
	GetStats(tenantID uuid.UUID) (*NotificationStatsResponse, error)

	// Convenience methods for specific notification types
	SendOrderConfirmationEmail(email, orderNumber string, orderDetails interface{}) error
	SendPaymentSuccessEmail(email, orderNumber string, paymentDetails interface{}) error
	SendPaymentFailedEmail(email, orderNumber, reason string) error
	SendEmailVerificationEmail(email, token string) error
	SendPasswordResetEmail(email, token string) error
}

type service struct {
	repository   Repository
	validator    *validator.Validate
	emailService email.EmailService

	// SMS configuration - keeping legacy SMS provider for now
	smsProvider SMSProvider
}

func NewService(repository Repository, emailService email.EmailService) Service {
	return &service{
		repository:   repository,
		validator:    validator.New(),
		emailService: emailService,
		smsProvider: SMSProvider{
			Name:      getEnvOrDefault("SMS_PROVIDER", "local_bd"),
			APIKey:    getEnvOrDefault("SMS_API_KEY", ""),
			APISecret: getEnvOrDefault("SMS_API_SECRET", ""),
		},
	}
}

func (s *service) SendNotification(tenantID uuid.UUID, req *SendNotificationRequest) (*SendNotificationResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	var notificationIDs []string
	var content string
	var subject string

	// Process template if provided
	if req.TemplateID != "" {
		templateID, err := uuid.Parse(req.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("invalid template ID: %w", err)
		}
		
		template, err := s.repository.GetTemplate(tenantID, templateID)
		if err != nil {
			return nil, fmt.Errorf("template not found: %w", err)
		}
		
		content = s.processTemplate(template.Content, req.Variables)
		subject = s.processTemplate(template.Subject, req.Variables)
	} else {
		content = req.Content
		subject = req.Subject
	}

	// Create notifications for each recipient
	for _, recipient := range req.Recipients {
		notification := &Notification{
			TenantID:    tenantID,
			Type:        req.Type,
			Channel:     req.Channel,
			Subject:     subject,
			Content:     content,
			Recipient:   recipient,
			Status:      StatusPending,
			Priority:    req.Priority,
			ScheduledAt: req.ScheduledAt,
			Metadata:    "{}", // Initialize as empty JSON object
		}

		if req.UserID != "" {
			userID, err := uuid.Parse(req.UserID)
			if err == nil {
				notification.UserID = &userID
			}
		}

		if err := s.repository.Create(notification); err != nil {
			log.Printf("Failed to create notification: %v", err)
			continue // Log error but continue with other notifications
		}

		notificationIDs = append(notificationIDs, notification.ID.String())

		// Send immediately if not scheduled
		if req.ScheduledAt == nil {
			go s.sendNotificationAsync(notification)
		}
	}

	return &SendNotificationResponse{
		NotificationIDs: notificationIDs,
		Status:          "queued",
		Message:         fmt.Sprintf("Successfully queued %d notifications", len(notificationIDs)),
	}, nil
}

func (s *service) sendNotificationAsync(notification *Notification) {
	switch notification.Type {
	case TypeEmail:
		emailErr := s.sendEmailNotification(notification)
		s.updateNotificationStatus(notification, emailErr)
	case TypeSMS:
		smsErr := s.sendSMSNotification(notification)
		s.updateNotificationStatus(notification, smsErr)
	case TypePush:
		pushErr := s.sendPushNotification(notification)
		s.updateNotificationStatus(notification, pushErr)
	case TypeInApp:
		// In-app notifications are just stored in database
		notification.Status = StatusDelivered
		now := time.Now()
		notification.DeliveredAt = &now
		if err := s.repository.Update(notification); err != nil {
			log.Printf("Failed to update notification status: %v", err)
		}
	}
}

func (s *service) sendEmailNotification(notification *Notification) error {
	// Use our integrated email service
	emailRequest := &email.EmailRequest{
		To:      []string{notification.Recipient},
		Subject: notification.Subject,
		Content: notification.Content,
	}

	// Send email through our email service
	_, err := s.emailService.SendEmail(emailRequest)
	if err != nil {
		log.Printf("Failed to send email notification to %s: %v", notification.Recipient, err)
		return err
	}

	log.Printf("Email notification sent successfully to %s", notification.Recipient)
	return nil
}

func (s *service) sendEmailViaSendGrid(notification *Notification) error {
	// SendGrid API integration
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{"email": notification.Recipient},
				},
				"subject": notification.Subject,
			},
		},
		"from": map[string]string{
			"email": "noreply@example.com", // TODO: Get from config
			"name":  "E-commerce SaaS",     // TODO: Get from config
		},
		"content": []map[string]string{
			{
				"type":  "text/html",
				"value": notification.Content,
			},
		},
	}

	jsonData, _ := json.Marshal(payload)
	
	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+getEnvOrDefault("SENDGRID_API_KEY", ""))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("sendgrid API error: status %d", resp.StatusCode)
	}

	return nil
}

func (s *service) sendSMSNotification(notification *Notification) error {
	switch s.smsProvider.Name {
	case "local_bd":
		return s.sendSMSViaLocalBD(notification)
	case "twilio":
		return s.sendSMSViaTwilio(notification)
	default:
		return fmt.Errorf("unsupported SMS provider: %s", s.smsProvider.Name)
	}
}

func (s *service) sendSMSViaLocalBD(notification *Notification) error {
	// Bangladesh local SMS gateway integration
	payload := BDSMSGatewayRequest{
		Username: s.smsProvider.APIKey,
		Password: s.smsProvider.APISecret,
		Number:   notification.Recipient,
		Message:  notification.Content,
		Type:     "text",
	}

	jsonData, _ := json.Marshal(payload)
	
	// This would be the actual SMS gateway URL
	req, err := http.NewRequest("POST", "https://api.local-sms-bd.com/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var smsResp BDSMSGatewayResponse
	if err := json.NewDecoder(resp.Body).Decode(&smsResp); err != nil {
		return err
	}

	if smsResp.Status != "success" {
		return fmt.Errorf("SMS gateway error: %s", smsResp.Message)
	}

	return nil
}

func (s *service) sendSMSViaTwilio(notification *Notification) error {
	// Twilio SMS integration - placeholder
	return fmt.Errorf("twilio SMS not implemented yet")
}

func (s *service) sendPushNotification(notification *Notification) error {
	// Firebase Cloud Messaging (FCM) integration
	// This is a basic implementation - in production, you'd want to use the official FCM SDK
	pushProvider := getEnvOrDefault("PUSH_PROVIDER", "fcm")
	
	switch pushProvider {
	case "fcm":
		return s.sendPushViaFCM(notification)
	default:
		return fmt.Errorf("unsupported push provider: %s", pushProvider)
	}
}

func (s *service) sendPushViaFCM(notification *Notification) error {
	serverKey := getEnvOrDefault("FCM_SERVER_KEY", "")
	if serverKey == "" {
		return fmt.Errorf("FCM server key not configured")
	}

	payload := map[string]interface{}{
		"to": notification.Recipient, // This should be the FCM token
		"notification": map[string]string{
			"title": notification.Subject,
			"body":  notification.Content,
		},
		"data": map[string]string{
			"tenant_id": notification.TenantID.String(),
			"type":      notification.Type,
		},
	}

	jsonData, _ := json.Marshal(payload)
	
	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "key="+serverKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FCM API error: status %d", resp.StatusCode)
	}

	return nil
}

func (s *service) updateNotificationStatus(notification *Notification, err error) {
	now := time.Now()
	
	if err != nil {
		notification.Status = StatusFailed
		notification.FailedAt = &now
		notification.FailureReason = err.Error()
		notification.RetryCount++
	} else {
		notification.Status = StatusSent
		notification.SentAt = &now
	}

	if err := s.repository.Update(notification); err != nil {
		log.Printf("Failed to update notification status: %v", err)
	}
}

func (s *service) processTemplate(template string, variables map[string]interface{}) string {
	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

func (s *service) SendEmail(tenantID uuid.UUID, req *SendEmailRequest) error {
	sendReq := &SendNotificationRequest{
		Type:       TypeEmail,
		Channel:    ChannelMarketing, // Default channel
		Recipients: req.To,
		Subject:    req.Subject,
		Content:    req.Content,
		Variables:  req.Variables,
		TemplateID: req.TemplateID,
	}

	_, err := s.SendNotification(tenantID, sendReq)
	return err
}

func (s *service) SendSMS(tenantID uuid.UUID, req *SendSMSRequest) error {
	sendReq := &SendNotificationRequest{
		Type:       TypeSMS,
		Channel:    ChannelMarketing, // Default channel
		Recipients: req.To,
		Content:    req.Message,
		Variables:  req.Variables,
		TemplateID: req.TemplateID,
		UserID:     req.UserID,
	}

	_, err := s.SendNotification(tenantID, sendReq)
	return err
}

func (s *service) GetNotification(tenantID uuid.UUID, notificationID string) (*Notification, error) {
	id, err := uuid.Parse(notificationID)
	if err != nil {
		return nil, fmt.Errorf("invalid notification ID: %w", err)
	}

	return s.repository.GetByID(tenantID, id)
}

func (s *service) ListNotifications(tenantID uuid.UUID, userID *uuid.UUID, offset, limit int) ([]*Notification, int64, error) {
	return s.repository.List(tenantID, userID, offset, limit)
}

func (s *service) MarkAsRead(tenantID uuid.UUID, notificationID string) error {
	notification, err := s.GetNotification(tenantID, notificationID)
	if err != nil {
		return err
	}

	now := time.Now()
	notification.ReadAt = &now
	return s.repository.Update(notification)
}

func (s *service) CreateTemplate(tenantID uuid.UUID, req *CreateTemplateRequest) (*NotificationTemplate, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	template := &NotificationTemplate{
		TenantID: tenantID,
		Name:     req.Name,
		Type:     req.Type,
		Channel:  req.Channel,
		Subject:  req.Subject,
		Content:  req.Content,
		IsActive: true,
	}

	// Handle Variables field - ensure it's never null
	if req.Variables == nil {
		template.Variables = "{}"
	} else {
		variablesJSON, _ := json.Marshal(req.Variables)
		template.Variables = string(variablesJSON)
	}

	if err := s.repository.CreateTemplate(template); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return template, nil
}

func (s *service) UpdateTemplate(tenantID uuid.UUID, templateID string, req *UpdateTemplateRequest) error {
	id, err := uuid.Parse(templateID)
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	template, err := s.repository.GetTemplate(tenantID, id)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Subject != "" {
		template.Subject = req.Subject
	}
	if req.Content != "" {
		template.Content = req.Content
	}
	if req.Variables != nil {
		variablesJSON, _ := json.Marshal(req.Variables)
		template.Variables = string(variablesJSON)
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}

	return s.repository.UpdateTemplate(template)
}

func (s *service) GetTemplate(tenantID uuid.UUID, templateID string) (*NotificationTemplate, error) {
	id, err := uuid.Parse(templateID)
	if err != nil {
		return nil, fmt.Errorf("invalid template ID: %w", err)
	}

	return s.repository.GetTemplate(tenantID, id)
}

func (s *service) ListTemplates(tenantID uuid.UUID, notificationType, channel string) ([]*NotificationTemplate, error) {
	return s.repository.ListTemplates(tenantID, notificationType, channel)
}

func (s *service) GetPreferences(tenantID, userID uuid.UUID, channel string) (*NotificationPreference, error) {
	preference, err := s.repository.GetPreferences(tenantID, userID, channel)
	if err != nil {
		// Return default preferences if not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &NotificationPreference{
				TenantID:     tenantID,
				UserID:       userID,
				Channel:      channel,
				EmailEnabled: true,  // Default to enabled
				SMSEnabled:   false, // Default to disabled
				PushEnabled:  true,  // Default to enabled
				InAppEnabled: true,  // Default to enabled
			}, nil
		}
		return nil, err
	}
	return preference, nil
}

func (s *service) UpdatePreferences(tenantID, userID uuid.UUID, req *NotificationPreferenceRequest) error {
	// Check if preference exists for this channel
	preference, err := s.repository.GetPreferences(tenantID, userID, req.Channel)
	isNewRecord := false
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new preference for this channel
			isNewRecord = true
			preference = &NotificationPreference{
				ID:           uuid.New(),
				TenantID:     tenantID,
				UserID:       userID,
				Channel:      req.Channel,
				EmailEnabled: true,  // Default values
				SMSEnabled:   false,
				PushEnabled:  true,
				InAppEnabled: true,
				CreatedAt:    time.Now(),
			}
		} else {
			return err
		}
	}

	// Update preferences based on request
	if req.EmailEnabled != nil {
		preference.EmailEnabled = *req.EmailEnabled
	}
	if req.SMSEnabled != nil {
		preference.SMSEnabled = *req.SMSEnabled
	}
	if req.PushEnabled != nil {
		preference.PushEnabled = *req.PushEnabled
	}
	if req.InAppEnabled != nil {
		preference.InAppEnabled = *req.InAppEnabled
	}

	preference.UpdatedAt = time.Now()

	// Save preference
	if isNewRecord {
		return s.repository.CreatePreferences(preference)
	} else {
		return s.repository.UpdatePreferences(preference)
	}
}

func (s *service) GetStats(tenantID uuid.UUID) (*NotificationStatsResponse, error) {
	statsMap, err := s.repository.GetNotificationStats(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification stats: %w", err)
	}

	// Convert map to structured response
	stats := &NotificationStatsResponse{
		TotalSent:      statsMap["sent"] + statsMap["delivered"], // sent + delivered
		TotalDelivered: statsMap["delivered"],
		TotalFailed:    statsMap["failed"],
	}

	// Calculate rates
	total := statsMap["total"]
	if total > 0 {
		stats.DeliveryRate = float64(statsMap["delivered"]) / float64(total) * 100
		stats.FailureRate = float64(statsMap["failed"]) / float64(total) * 100
	}

	// Email stats
	stats.EmailStats.Sent = statsMap["email"]
	stats.EmailStats.Delivered = statsMap["email"] // Simplified for now

	// SMS stats
	stats.SMSStats.Sent = statsMap["sms"]
	stats.SMSStats.Delivered = statsMap["sms"] // Simplified for now

	return stats, nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Convenience methods for specific email types

// SendOrderConfirmationEmail sends order confirmation email
func (s *service) SendOrderConfirmationEmail(email, orderNumber string, orderDetails interface{}) error {
	return s.emailService.SendOrderConfirmationEmail(email, orderNumber, orderDetails)
}

// SendPaymentSuccessEmail sends payment success email
func (s *service) SendPaymentSuccessEmail(email, orderNumber string, paymentDetails interface{}) error {
	return s.emailService.SendPaymentSuccessEmail(email, orderNumber, paymentDetails)
}

// SendPaymentFailedEmail sends payment failed email
func (s *service) SendPaymentFailedEmail(email, orderNumber, reason string) error {
	return s.emailService.SendPaymentFailedEmail(email, orderNumber, reason)
}

// SendEmailVerificationEmail sends email verification email
func (s *service) SendEmailVerificationEmail(email, token string) error {
	return s.emailService.SendVerificationEmail(email, token)
}

// SendPasswordResetEmail sends password reset email
func (s *service) SendPasswordResetEmail(email, token string) error {
	return s.emailService.SendPasswordResetEmail(email, token)
}
