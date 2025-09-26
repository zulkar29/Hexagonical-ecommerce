package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"ecommerce-saas/internal/shared/config"
)

// EmailProvider represents different email service providers
type EmailProvider string

const (
	ProviderSendGrid  EmailProvider = "sendgrid"
	ProviderMailgun   EmailProvider = "mailgun"
	ProviderSMTP      EmailProvider = "smtp"
	ProviderConsole   EmailProvider = "console" // For development
)

// EmailRequest represents an email to be sent
type EmailRequest struct {
	To          []string               `json:"to" validate:"required,min=1"`
	Subject     string                 `json:"subject" validate:"required"`
	Content     string                 `json:"content" validate:"required"`
	ContentType string                 `json:"content_type,omitempty"` // "text/plain" or "text/html"
	Variables   map[string]interface{} `json:"variables,omitempty"`
	TemplateID  string                 `json:"template_id,omitempty"`
	From        string                 `json:"from,omitempty"`
	FromName    string                 `json:"from_name,omitempty"`
}

// EmailResponse represents the response from email service
type EmailResponse struct {
	MessageID   string `json:"message_id"`
	Status      string `json:"status"`
	Provider    string `json:"provider"`
	SentAt      time.Time `json:"sent_at"`
	Error       string `json:"error,omitempty"`
}

// EmailService interface for different email providers
type EmailService interface {
	SendEmail(req *EmailRequest) (*EmailResponse, error)
	SendVerificationEmail(email, token string) error
	SendPasswordResetEmail(email, token string) error
	SendOrderConfirmationEmail(email, orderNumber string, orderDetails interface{}) error
	SendPaymentSuccessEmail(email, orderNumber string, paymentDetails interface{}) error
	SendPaymentFailedEmail(email, orderNumber string, reason string) error
}

// Service implements EmailService
type Service struct {
	provider    EmailProvider
	config      *config.Config
	httpClient  *http.Client

	// SendGrid config
	sendGridAPIKey string

	// Mailgun config
	mailgunAPIKey string
	mailgunDomain string

	// SMTP config
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string

	// Default sender info
	defaultFrom     string
	defaultFromName string
}

// NewService creates a new email service
func NewService(cfg *config.Config) EmailService {
	provider := EmailProvider(getEnvOrDefault("EMAIL_PROVIDER", "console"))

	service := &Service{
		provider:   provider,
		config:     cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},

		// Load configuration based on provider
		sendGridAPIKey:  getEnvOrDefault("SENDGRID_API_KEY", ""),
		mailgunAPIKey:   getEnvOrDefault("MAILGUN_API_KEY", ""),
		mailgunDomain:   getEnvOrDefault("MAILGUN_DOMAIN", ""),
		smtpHost:        getEnvOrDefault("SMTP_HOST", ""),
		smtpPort:        getEnvOrDefault("SMTP_PORT", "587"),
		smtpUsername:    getEnvOrDefault("SMTP_USERNAME", ""),
		smtpPassword:    getEnvOrDefault("SMTP_PASSWORD", ""),
		defaultFrom:     getEnvOrDefault("DEFAULT_FROM_EMAIL", "noreply@example.com"),
		defaultFromName: getEnvOrDefault("DEFAULT_FROM_NAME", "E-commerce SaaS"),
	}

	return service
}

// SendEmail sends an email using the configured provider
func (s *Service) SendEmail(req *EmailRequest) (*EmailResponse, error) {
	// Set defaults
	if req.From == "" {
		req.From = s.defaultFrom
	}
	if req.FromName == "" {
		req.FromName = s.defaultFromName
	}
	if req.ContentType == "" {
		req.ContentType = "text/html"
	}

	// Replace template variables if provided
	if len(req.Variables) > 0 {
		req.Content = s.replaceVariables(req.Content, req.Variables)
		req.Subject = s.replaceVariables(req.Subject, req.Variables)
	}

	switch s.provider {
	case ProviderSendGrid:
		return s.sendWithSendGrid(req)
	case ProviderMailgun:
		return s.sendWithMailgun(req)
	case ProviderSMTP:
		return s.sendWithSMTP(req)
	case ProviderConsole:
		return s.sendWithConsole(req)
	default:
		return s.sendWithConsole(req) // Fallback to console
	}
}

// SendVerificationEmail sends account verification email
func (s *Service) SendVerificationEmail(email, token string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.config.App.BaseURL, token)

	req := &EmailRequest{
		To:      []string{email},
		Subject: "Verify Your Account - {{company_name}}",
		Content: s.getVerificationEmailTemplate(),
		Variables: map[string]interface{}{
			"company_name": s.defaultFromName,
			"verify_url":   verifyURL,
			"token":        token,
		},
	}

	_, err := s.SendEmail(req)
	return err
}

// SendPasswordResetEmail sends password reset email
func (s *Service) SendPasswordResetEmail(email, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.config.App.BaseURL, token)

	req := &EmailRequest{
		To:      []string{email},
		Subject: "Reset Your Password - {{company_name}}",
		Content: s.getPasswordResetEmailTemplate(),
		Variables: map[string]interface{}{
			"company_name": s.defaultFromName,
			"reset_url":    resetURL,
			"token":        token,
		},
	}

	_, err := s.SendEmail(req)
	return err
}

// SendOrderConfirmationEmail sends order confirmation email
func (s *Service) SendOrderConfirmationEmail(email, orderNumber string, orderDetails interface{}) error {
	req := &EmailRequest{
		To:      []string{email},
		Subject: "Order Confirmation #{{order_number}} - {{company_name}}",
		Content: s.getOrderConfirmationEmailTemplate(),
		Variables: map[string]interface{}{
			"company_name":   s.defaultFromName,
			"order_number":   orderNumber,
			"order_details":  orderDetails,
		},
	}

	_, err := s.SendEmail(req)
	return err
}

// SendPaymentSuccessEmail sends payment success email
func (s *Service) SendPaymentSuccessEmail(email, orderNumber string, paymentDetails interface{}) error {
	req := &EmailRequest{
		To:      []string{email},
		Subject: "Payment Successful #{{order_number}} - {{company_name}}",
		Content: s.getPaymentSuccessEmailTemplate(),
		Variables: map[string]interface{}{
			"company_name":    s.defaultFromName,
			"order_number":    orderNumber,
			"payment_details": paymentDetails,
		},
	}

	_, err := s.SendEmail(req)
	return err
}

// SendPaymentFailedEmail sends payment failed email
func (s *Service) SendPaymentFailedEmail(email, orderNumber, reason string) error {
	req := &EmailRequest{
		To:      []string{email},
		Subject: "Payment Failed #{{order_number}} - {{company_name}}",
		Content: s.getPaymentFailedEmailTemplate(),
		Variables: map[string]interface{}{
			"company_name": s.defaultFromName,
			"order_number": orderNumber,
			"reason":       reason,
		},
	}

	_, err := s.SendEmail(req)
	return err
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}