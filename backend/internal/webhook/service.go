package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	repo       *Repository
	httpClient *http.Client
	signingKey []byte
}

type TestResult struct {
	Success      bool   `json:"success"`
	ResponseCode int    `json:"response_code"`
	ResponseTime int    `json:"response_time"`
	Error        string `json:"error,omitempty"`
}

func NewService(repo *Repository, signingKey []byte) *Service {
	return &Service{
		repo:       repo,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		signingKey: signingKey,
	}
}

// Webhook Endpoint Management Methods

func (s *Service) CreateEndpoint(tenantID uuid.UUID, endpoint *WebhookEndpoint) (*WebhookEndpoint, error) {
	endpoint.ID = uuid.New()
	endpoint.TenantID = tenantID
	endpoint.CreatedAt = time.Now()
	endpoint.UpdatedAt = time.Now()
	
	// Generate a secret for the endpoint
	if endpoint.Secret == "" {
		endpoint.Secret = s.generateSecret()
	}
	
	return s.repo.CreateEndpoint(endpoint)
}

func (s *Service) UpdateEndpoint(tenantID uuid.UUID, endpointID uuid.UUID, updates *WebhookEndpoint) (*WebhookEndpoint, error) {
	existing, err := s.repo.GetEndpointByID(tenantID, endpointID)
	if err != nil {
		return nil, err
	}
	
	// Update allowed fields
	if updates.URL != "" {
		existing.URL = updates.URL
	}
	if updates.Description != "" {
		existing.Description = updates.Description
	}
	if len(updates.Events) > 0 {
		existing.Events = updates.Events
	}
	if updates.Secret != "" {
		existing.Secret = updates.Secret
	}
	existing.IsActive = updates.IsActive
	existing.UpdatedAt = time.Now()
	
	return s.repo.UpdateEndpoint(existing)
}

func (s *Service) DeleteEndpoint(tenantID uuid.UUID, endpointID uuid.UUID) error {
	return s.repo.DeleteEndpoint(tenantID, endpointID)
}

func (s *Service) GetEndpoints(tenantID uuid.UUID) ([]*WebhookEndpoint, error) {
	return s.repo.GetEndpoints(tenantID)
}

func (s *Service) GetEndpoint(tenantID uuid.UUID, endpointID uuid.UUID) (*WebhookEndpoint, error) {
	return s.repo.GetEndpointByID(tenantID, endpointID)
}

func (s *Service) generateSecret() string {
	return hex.EncodeToString([]byte(uuid.New().String()))
}

// Event Dispatch Methods

func (s *Service) DispatchEvent(tenantID uuid.UUID, event WebhookEvent, eventID uuid.UUID, payload interface{}) error {
	// Get active endpoints that support this event
	endpoints, err := s.repo.GetEndpointsByEvent(tenantID, event)
	if err != nil {
		return err
	}
	
	// Create deliveries for each endpoint
	for _, endpoint := range endpoints {
		if !endpoint.IsActive {
			continue
		}
		
		// Check rate limiting
		if s.IsRateLimited(tenantID, endpoint.ID) {
			continue
		}
		
		// Create delivery record
		delivery := &WebhookDelivery{
			ID:         uuid.New(),
			TenantID:   tenantID,
			EndpointID: endpoint.ID,
			Event:      event,
			EventID:    eventID,
			Payload:    payload,
			Status:     StatusPending,
			Attempts:   0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		
		// Save delivery record
		if _, err := s.repo.CreateDelivery(delivery); err != nil {
			continue // Log error but continue with other endpoints
		}
		
		// Attempt immediate delivery
		go s.DeliverWebhook(delivery)
	}
	
	return nil
}

func (s *Service) DeliverWebhook(delivery *WebhookDelivery) error {
	// Get endpoint details
	endpoint, err := s.repo.GetEndpointByID(delivery.TenantID, delivery.EndpointID)
	if err != nil {
		return err
	}
	
	// Prepare payload
	payloadBytes, err := json.Marshal(delivery.Payload)
	if err != nil {
		return err
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", endpoint.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Webhook-Service/1.0")
	req.Header.Set("X-Webhook-Event", string(delivery.Event))
	req.Header.Set("X-Webhook-ID", delivery.ID.String())
	req.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	
	// Generate and set signature
	signature := s.GenerateWebhookSignature(payloadBytes, endpoint.Secret)
	req.Header.Set("X-Webhook-Signature", signature)
	
	// Update delivery attempt
	delivery.Attempts++
	delivery.LastAttemptAt = &time.Time{}
	*delivery.LastAttemptAt = time.Now()
	delivery.UpdatedAt = time.Now()
	
	// Make HTTP request
	start := time.Now()
	resp, err := s.httpClient.Do(req)
	responseTime := int(time.Since(start).Milliseconds())
	
	delivery.ResponseTime = responseTime
	
	if err != nil {
		delivery.Status = StatusFailed
		delivery.ErrorMessage = err.Error()
	} else {
		delivery.ResponseStatus = resp.StatusCode
		
		// Read response body
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		delivery.ResponseBody = string(body)
		
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			delivery.Status = StatusDelivered
			delivery.DeliveredAt = delivery.LastAttemptAt
		} else {
			delivery.Status = StatusFailed
			delivery.ErrorMessage = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body))
		}
	}
	
	// Update delivery record
	s.repo.UpdateDelivery(delivery)
	
	// Schedule retry if failed and retries remaining
	if delivery.Status == StatusFailed && delivery.ShouldRetry() {
		s.ScheduleRetry(delivery)
	}
	
	// Increment rate limit
	s.IncrementRateLimit(delivery.TenantID, delivery.EndpointID)
	
	return nil
}

func (s *Service) RetryFailedDeliveries() error {
	// Get pending retries
	deliveries, err := s.repo.GetPendingRetries()
	if err != nil {
		return err
	}
	
	for _, delivery := range deliveries {
		if delivery.ShouldRetry() {
			go s.DeliverWebhook(delivery)
		}
	}
	
	return nil
}

func (s *Service) ProcessRetryQueue() error {
	return s.RetryFailedDeliveries()
}

// Event-specific dispatch methods

func (s *Service) DispatchOrderCreated(tenantID uuid.UUID, orderID uuid.UUID, order interface{}) error {
	return s.DispatchEvent(tenantID, EventOrderCreated, orderID, order)
}

func (s *Service) DispatchOrderUpdated(tenantID uuid.UUID, orderID uuid.UUID, order interface{}) error {
	return s.DispatchEvent(tenantID, EventOrderUpdated, orderID, order)
}

func (s *Service) DispatchOrderCancelled(tenantID uuid.UUID, orderID uuid.UUID, order interface{}) error {
	return s.DispatchEvent(tenantID, EventOrderCancelled, orderID, order)
}

func (s *Service) DispatchPaymentSucceeded(tenantID uuid.UUID, paymentID uuid.UUID, payment interface{}) error {
	return s.DispatchEvent(tenantID, EventPaymentSucceeded, paymentID, payment)
}

func (s *Service) DispatchPaymentFailed(tenantID uuid.UUID, paymentID uuid.UUID, payment interface{}) error {
	return s.DispatchEvent(tenantID, EventPaymentFailed, paymentID, payment)
}

func (s *Service) DispatchPaymentRefunded(tenantID uuid.UUID, paymentID uuid.UUID, payment interface{}) error {
	return s.DispatchEvent(tenantID, EventPaymentRefunded, paymentID, payment)
}

func (s *Service) DispatchProductCreated(tenantID uuid.UUID, productID uuid.UUID, product interface{}) error {
	return s.DispatchEvent(tenantID, EventProductCreated, productID, product)
}

func (s *Service) DispatchProductUpdated(tenantID uuid.UUID, productID uuid.UUID, product interface{}) error {
	return s.DispatchEvent(tenantID, EventProductUpdated, productID, product)
}

func (s *Service) DispatchInventoryLow(tenantID uuid.UUID, productID uuid.UUID, inventoryData interface{}) error {
	return s.DispatchEvent(tenantID, EventInventoryLow, productID, inventoryData)
}

func (s *Service) DispatchShipmentCreated(tenantID uuid.UUID, shipmentID uuid.UUID, shipment interface{}) error {
	return s.DispatchEvent(tenantID, EventShipmentCreated, shipmentID, shipment)
}

func (s *Service) DispatchShipmentDelivered(tenantID uuid.UUID, shipmentID uuid.UUID, shipment interface{}) error {
	return s.DispatchEvent(tenantID, EventShipmentDelivered, shipmentID, shipment)
}

func (s *Service) DispatchCustomerCreated(tenantID uuid.UUID, customerID uuid.UUID, customer interface{}) error {
	return s.DispatchEvent(tenantID, EventCustomerCreated, customerID, customer)
}

func (s *Service) DispatchCustomerUpdated(tenantID uuid.UUID, customerID uuid.UUID, customer interface{}) error {
	return s.DispatchEvent(tenantID, EventCustomerUpdated, customerID, customer)
}

// Incoming Webhook Processing

func (s *Service) ProcessStripeWebhook(tenantID uuid.UUID, signature string, body []byte) error {
	// Validate Stripe webhook signature
	if !s.validateStripeSignature(signature, body) {
		return fmt.Errorf("invalid stripe webhook signature")
	}
	
	// Parse Stripe event
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	
	// Create incoming webhook record
	incoming := &WebhookIncoming{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Provider:  ProviderStripe,
		EventType: fmt.Sprintf("%v", event["type"]),
		Payload:   event,
		Signature: signature,
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return s.processIncomingWebhook(incoming)
}

func (s *Service) ProcessPayPalWebhook(tenantID uuid.UUID, signature string, body []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	
	incoming := &WebhookIncoming{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Provider:  ProviderPayPal,
		EventType: fmt.Sprintf("%v", event["event_type"]),
		Payload:   event,
		Signature: signature,
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return s.processIncomingWebhook(incoming)
}

func (s *Service) ProcessBkashWebhook(tenantID uuid.UUID, signature string, body []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	
	incoming := &WebhookIncoming{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Provider:  ProviderBkash,
		EventType: fmt.Sprintf("%v", event["event_type"]),
		Payload:   event,
		Signature: signature,
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return s.processIncomingWebhook(incoming)
}

func (s *Service) ProcessNagadWebhook(tenantID uuid.UUID, signature string, body []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	
	incoming := &WebhookIncoming{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Provider:  ProviderNagad,
		EventType: fmt.Sprintf("%v", event["event_type"]),
		Payload:   event,
		Signature: signature,
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return s.processIncomingWebhook(incoming)
}

func (s *Service) ProcessPathaoWebhook(tenantID uuid.UUID, signature string, body []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	
	incoming := &WebhookIncoming{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Provider:  ProviderPathao,
		EventType: fmt.Sprintf("%v", event["event_type"]),
		Payload:   event,
		Signature: signature,
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return s.processIncomingWebhook(incoming)
}

func (s *Service) ProcessRedXWebhook(tenantID uuid.UUID, signature string, body []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	
	incoming := &WebhookIncoming{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Provider:  ProviderRedX,
		EventType: fmt.Sprintf("%v", event["event_type"]),
		Payload:   event,
		Signature: signature,
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return s.processIncomingWebhook(incoming)
}

func (s *Service) ProcessPaperflyWebhook(tenantID uuid.UUID, signature string, body []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	
	incoming := &WebhookIncoming{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Provider:  ProviderPaperfly,
		EventType: fmt.Sprintf("%v", event["event_type"]),
		Payload:   event,
		Signature: signature,
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return s.processIncomingWebhook(incoming)
}

func (s *Service) processIncomingWebhook(incoming *WebhookIncoming) error {
	// Save incoming webhook
	if _, err := s.repo.CreateIncomingWebhook(incoming); err != nil {
		return err
	}
	
	// Process webhook asynchronously
	go s.handleIncomingWebhook(incoming)
	
	return nil
}

func (s *Service) handleIncomingWebhook(incoming *WebhookIncoming) {
	// Mark as processed
	defer func() {
		incoming.Processed = true
		incoming.ProcessedAt = &time.Time{}
		*incoming.ProcessedAt = time.Now()
		incoming.UpdatedAt = time.Now()
		s.repo.UpdateIncomingWebhook(incoming)
	}()
	
	// Handle based on provider and event type
	switch incoming.Provider {
	case ProviderStripe:
		s.handleStripeEvent(incoming)
	case ProviderPayPal:
		s.handlePayPalEvent(incoming)
	case ProviderBkash, ProviderNagad:
		s.handleMobilePaymentEvent(incoming)
	case ProviderPathao, ProviderRedX, ProviderPaperfly:
		s.handleShippingEvent(incoming)
	}
}

func (s *Service) handleStripeEvent(incoming *WebhookIncoming) {
	// Handle Stripe-specific events
	// This would integrate with payment processing logic
}

func (s *Service) handlePayPalEvent(incoming *WebhookIncoming) {
	// Handle PayPal-specific events
}

func (s *Service) handleMobilePaymentEvent(incoming *WebhookIncoming) {
	// Handle mobile payment events (bKash, Nagad)
}

func (s *Service) handleShippingEvent(incoming *WebhookIncoming) {
	// Handle shipping provider events
}

func (s *Service) validateStripeSignature(signature string, body []byte) bool {
	// Implement Stripe signature validation
	// This is a simplified version - real implementation would use Stripe's validation
	return signature != ""
}

// Security and Validation

func (s *Service) ValidateWebhookSignature(payload []byte, signature string, secret string) bool {
	// Generate expected signature using HMAC-SHA256
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	
	// Compare signatures (constant time comparison)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (s *Service) GenerateWebhookSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func (s *Service) IsRateLimited(tenantID uuid.UUID, endpointID uuid.UUID) bool {
	// Get current rate limit for endpoint
	rateLimit, err := s.repo.GetRateLimit(tenantID, endpointID)
	if err != nil {
		// No rate limit exists, not limited
		return false
	}
	
	// Check if we're in a new time window (1 hour)
	if time.Since(rateLimit.WindowStart) > time.Hour {
		// Reset the window
		rateLimit.Count = 0
		rateLimit.WindowStart = time.Now()
		rateLimit.UpdatedAt = time.Now()
		s.repo.UpdateRateLimit(rateLimit)
		return false
	}
	
	// Check if rate limit exceeded (default: 1000 requests per hour)
	maxRequests := 1000
	return rateLimit.Count >= maxRequests
}

func (s *Service) IncrementRateLimit(tenantID uuid.UUID, endpointID uuid.UUID) error {
	// Get current rate limit
	rateLimit, err := s.repo.GetRateLimit(tenantID, endpointID)
	if err != nil {
		// No rate limit exists, create one
		rateLimit = &WebhookRateLimit{
			ID:         uuid.New(),
			TenantID:   tenantID,
			EndpointID: endpointID,
			Count:      1,
			WindowStart: time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_, err := s.repo.CreateRateLimit(rateLimit)
		return err
	}
	
	// Check if we're in a new time window (1 hour)
	if time.Since(rateLimit.WindowStart) > time.Hour {
		// Reset the window
		rateLimit.Count = 1
		rateLimit.WindowStart = time.Now()
	} else {
		// Increment count
		rateLimit.Count++
	}
	
	rateLimit.UpdatedAt = time.Now()
	return s.repo.UpdateRateLimit(rateLimit)
}

// Monitoring and Analytics

func (s *Service) GetDeliveryStats(tenantID uuid.UUID, startDate, endDate time.Time) (*WebhookStats, error) {
	return s.repo.GetDeliveryStats(tenantID, startDate, endDate)
}

func (s *Service) GetFailedDeliveries(tenantID uuid.UUID, limit int) ([]*WebhookDelivery, error) {
	return s.repo.GetFailedDeliveries(tenantID, limit)
}

func (s *Service) GetEndpointHealth(tenantID uuid.UUID, endpointID uuid.UUID) (*EndpointHealth, error) {
	return s.repo.GetEndpointHealth(tenantID, endpointID)
}

func (s *Service) TestEndpoint(tenantID uuid.UUID, endpointID uuid.UUID) (*TestResult, error) {
	// Get endpoint
	endpoint, err := s.repo.GetEndpointByID(tenantID, endpointID)
	if err != nil {
		return nil, err
	}
	
	// Create test payload
	testPayload := map[string]interface{}{
		"event": "test",
		"timestamp": time.Now().Unix(),
		"data": map[string]string{
			"message": "This is a test webhook",
		},
	}
	
	payloadBytes, err := json.Marshal(testPayload)
	if err != nil {
		return &TestResult{
			Success: false,
			Error:   "Failed to marshal test payload",
		}, nil
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", endpoint.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return &TestResult{
			Success: false,
			Error:   "Failed to create HTTP request",
		}, nil
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Webhook-Service/1.0")
	req.Header.Set("X-Webhook-Event", "test")
	req.Header.Set("X-Webhook-Test", "true")
	
	// Generate and set signature
	signature := s.GenerateWebhookSignature(payloadBytes, endpoint.Secret)
	req.Header.Set("X-Webhook-Signature", signature)
	
	// Make HTTP request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	
	start := time.Now()
	resp, err := s.httpClient.Do(req)
	responseTime := int(time.Since(start).Milliseconds())
	
	if err != nil {
		return &TestResult{
			Success:      false,
			ResponseTime: responseTime,
			Error:        err.Error(),
		}, nil
	}
	defer resp.Body.Close()
	
	return &TestResult{
		Success:      resp.StatusCode >= 200 && resp.StatusCode < 300,
		ResponseCode: resp.StatusCode,
		ResponseTime: responseTime,
	}, nil
}

// Background Processing

func (s *Service) ScheduleRetry(delivery *WebhookDelivery) error {
	// Calculate next retry time using exponential backoff
	backoffDuration := time.Duration(delivery.Attempts*delivery.Attempts) * time.Minute
	if backoffDuration > 24*time.Hour {
		backoffDuration = 24 * time.Hour
	}
	
	nextRetry := time.Now().Add(backoffDuration)
	delivery.NextRetryAt = &nextRetry
	delivery.UpdatedAt = time.Now()
	
	return s.repo.UpdateDelivery(delivery)
}

func (s *Service) CleanupOldDeliveries(olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)
	return s.repo.DeleteOldDeliveries(cutoffTime)
}

func (s *Service) DisableFailingEndpoints() error {
	// Get endpoints with high failure rates
	failingEndpoints, err := s.repo.GetFailingEndpoints(0.8, 100) // 80% failure rate, min 100 attempts
	if err != nil {
		return err
	}
	
	for _, endpoint := range failingEndpoints {
		endpoint.IsActive = false
		endpoint.UpdatedAt = time.Now()
		if err := s.repo.UpdateEndpoint(endpoint); err != nil {
			// Log error but continue with other endpoints
			continue
		}
	}
	
	return nil
}

func (s *Service) ProcessBackgroundTasks() error {
	// Process retry queue
	if err := s.ProcessRetryQueue(); err != nil {
		return err
	}
	
	// Cleanup old deliveries (older than 30 days)
	if err := s.CleanupOldDeliveries(30 * 24 * time.Hour); err != nil {
		return err
	}
	
	// Disable failing endpoints
	if err := s.DisableFailingEndpoints(); err != nil {
		return err
	}
	
	return nil
}
