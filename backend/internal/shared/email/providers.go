package email

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

// SendGrid implementation
func (s *Service) sendWithSendGrid(req *EmailRequest) (*EmailResponse, error) {
	if s.sendGridAPIKey == "" {
		return nil, fmt.Errorf("SendGrid API key not configured")
	}

	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": func() []map[string]string {
					recipients := make([]map[string]string, len(req.To))
					for i, email := range req.To {
						recipients[i] = map[string]string{"email": email}
					}
					return recipients
				}(),
				"subject": req.Subject,
			},
		},
		"from": map[string]string{
			"email": req.From,
			"name":  req.FromName,
		},
		"content": []map[string]string{
			{
				"type":  req.ContentType,
				"value": req.Content,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SendGrid payload: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create SendGrid request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.sendGridAPIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send SendGrid request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &EmailResponse{
			Provider: "sendgrid",
			Status:   "failed",
			Error:    fmt.Sprintf("SendGrid API error: %d", resp.StatusCode),
			SentAt:   time.Now(),
		}, nil
	}

	return &EmailResponse{
		MessageID: resp.Header.Get("X-Message-Id"),
		Status:    "sent",
		Provider:  "sendgrid",
		SentAt:    time.Now(),
	}, nil
}

// Mailgun implementation
func (s *Service) sendWithMailgun(req *EmailRequest) (*EmailResponse, error) {
	if s.mailgunAPIKey == "" || s.mailgunDomain == "" {
		return nil, fmt.Errorf("Mailgun configuration not complete")
	}

	url := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", s.mailgunDomain)

	payload := fmt.Sprintf(
		"from=%s <%s>&to=%s&subject=%s&html=%s",
		req.FromName,
		req.From,
		strings.Join(req.To, ","),
		req.Subject,
		req.Content,
	)

	httpReq, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create Mailgun request: %w", err)
	}

	httpReq.SetBasicAuth("api", s.mailgunAPIKey)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send Mailgun request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &EmailResponse{
			Provider: "mailgun",
			Status:   "failed",
			Error:    fmt.Sprintf("Mailgun API error: %d", resp.StatusCode),
			SentAt:   time.Now(),
		}, nil
	}

	var result struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		return &EmailResponse{
			MessageID: result.ID,
			Status:    "sent",
			Provider:  "mailgun",
			SentAt:    time.Now(),
		}, nil
	}

	return &EmailResponse{
		Status:   "sent",
		Provider: "mailgun",
		SentAt:   time.Now(),
	}, nil
}

// SMTP implementation
func (s *Service) sendWithSMTP(req *EmailRequest) (*EmailResponse, error) {
	if s.smtpHost == "" || s.smtpUsername == "" || s.smtpPassword == "" {
		return nil, fmt.Errorf("SMTP configuration not complete")
	}

	// Create message
	message := fmt.Sprintf(
		"From: %s <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: %s; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		req.FromName,
		req.From,
		strings.Join(req.To, ", "),
		req.Subject,
		req.ContentType,
		req.Content,
	)

	// Setup TLS config
	tlsConfig := &tls.Config{
		ServerName: s.smtpHost,
	}

	// Connect to server
	conn, err := tls.Dial("tcp", s.smtpHost+":"+s.smtpPort, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.smtpHost)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	if err := client.Auth(auth); err != nil {
		return nil, fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Send email
	if err := client.Mail(req.From); err != nil {
		return nil, fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range req.To {
		if err := client.Rcpt(recipient); err != nil {
			return nil, fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to get data writer: %w", err)
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		writer.Close()
		return nil, fmt.Errorf("failed to write message: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close message writer: %w", err)
	}

	return &EmailResponse{
		Status:   "sent",
		Provider: "smtp",
		SentAt:   time.Now(),
	}, nil
}

// Console implementation (for development)
func (s *Service) sendWithConsole(req *EmailRequest) (*EmailResponse, error) {
	log.Printf("=== EMAIL (Console Provider) ===")
	log.Printf("From: %s <%s>", req.FromName, req.From)
	log.Printf("To: %s", strings.Join(req.To, ", "))
	log.Printf("Subject: %s", req.Subject)
	log.Printf("Content-Type: %s", req.ContentType)
	log.Printf("--- Content ---")
	log.Printf("%s", req.Content)
	log.Printf("===============")

	return &EmailResponse{
		MessageID: fmt.Sprintf("console_%d", time.Now().Unix()),
		Status:    "sent",
		Provider:  "console",
		SentAt:    time.Now(),
	}, nil
}

// replaceVariables replaces template variables in content
func (s *Service) replaceVariables(content string, variables map[string]interface{}) string {
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		replacement := fmt.Sprintf("%v", value)
		content = strings.ReplaceAll(content, placeholder, replacement)
	}
	return content
}