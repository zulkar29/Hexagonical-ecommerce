package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
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
	GetPreferences(tenantID, userID uuid.UUID) (*NotificationPreference, error)
	UpdatePreferences(tenantID, userID uuid.UUID, req *NotificationPreferenceRequest) error
	
	// Stats
	GetStats(tenantID uuid.UUID) (*NotificationStatsResponse, error)
}

type service struct {
	repository Repository
	validator  *validator.Validate
	
	// Email configuration
	emailProvider EmailProvider
	
	// SMS configuration
	smsProvider SMSProvider
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
		validator:  validator.New(),
		emailProvider: EmailProvider{
			Name:      "sendgrid",
			APIKey:    "your_sendgrid_key", // TODO: Load from config
			FromEmail: "noreply@yourdomain.com",
			FromName:  "Your Platform",
		},
		smsProvider: SMSProvider{
			Name:      "local_bd",
			APIKey:    "your_sms_key", // TODO: Load from config
			APISecret: "your_sms_secret",
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
		}

		if req.UserID != "" {
			userID, err := uuid.Parse(req.UserID)
			if err == nil {
				notification.UserID = &userID
			}
		}

		if err := s.repository.Create(notification); err != nil {
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
		err := s.sendEmailNotification(notification)
		s.updateNotificationStatus(notification, err)
	case TypeSMS:
		err := s.sendSMSNotification(notification)
		s.updateNotificationStatus(notification, err)
	case TypePush:
		// TODO: Implement push notification
		s.updateNotificationStatus(notification, fmt.Errorf("push notifications not implemented"))
	case TypeInApp:
		// In-app notifications are just stored in database
		notification.Status = StatusDelivered
		now := time.Now()
		notification.DeliveredAt = &now
		s.repository.Update(notification)
	}
}

func (s *service) sendEmailNotification(notification *Notification) error {
	switch s.emailProvider.Name {
	case "sendgrid":
		return s.sendEmailViaSendGrid(notification)
	default:
		return fmt.Errorf("unsupported email provider: %s", s.emailProvider.Name)
	}
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
			"email": s.emailProvider.FromEmail,
			"name":  s.emailProvider.FromName,
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

	req.Header.Set("Authorization", "Bearer "+s.emailProvider.APIKey)
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

	s.repository.Update(notification)
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

	variablesJSON, _ := json.Marshal(req.Variables)
	template.Variables = string(variablesJSON)

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

func (s *service) GetPreferences(tenantID, userID uuid.UUID) (*NotificationPreference, error) {
	return s.repository.GetPreferences(tenantID, userID)
}

func (s *service) UpdatePreferences(tenantID, userID uuid.UUID, req *NotificationPreferenceRequest) error {
	preference, err := s.repository.GetPreferences(tenantID, userID)
	if err != nil {
		// Create new preference if not exists
		preference = &NotificationPreference{
			TenantID: tenantID,
			UserID:   userID,
			Channel:  req.Channel,
		}
	}

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

	return s.repository.UpdatePreferences(preference)
}

func (s *service) GetStats(tenantID uuid.UUID) (*NotificationStatsResponse, error) {
	// TODO: Implement comprehensive stats
	return &NotificationStatsResponse{
		TotalSent:      0,
		TotalDelivered: 0,
		TotalFailed:    0,
		DeliveryRate:   0,
		FailureRate:    0,
	}, nil
}
