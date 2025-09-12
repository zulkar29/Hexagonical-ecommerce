package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
	"ecommerce-saas/internal/shared/config"
)

type Service interface {
	CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error)
	ProcessPayment(ctx context.Context, req *ProcessPaymentRequest) (*Payment, error)
	GetPayment(ctx context.Context, id string) (*Payment, error)
	ListPayments(ctx context.Context, req *ListPaymentsRequest) (*ListPaymentsResponse, error)
	UpdatePayment(ctx context.Context, id string, updates map[string]interface{}) (*Payment, error)
	RefundPayment(ctx context.Context, req *RefundPaymentRequest) (*Payment, error)

	// Payment Methods Management
	GetPaymentMethods(ctx context.Context, userID string) ([]*PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, id string, req *UpdatePaymentMethodRequest) (*PaymentMethod, error)

	// SSLCommerz specific methods
	InitiateSSLCommerzPayment(ctx context.Context, payment *Payment) (*SSLCommerzPaymentResponse, error)
	ValidateSSLCommerzPayment(ctx context.Context, ipnData *SSLCommerzIPNResponse) error
}

type service struct {
	repository Repository
	validator  *validator.Validate
	config     *config.Config
	
	// SSLCommerz configuration
	sslCommerzStoreID    string
	sslCommerzStorePass  string
	sslCommerzSandbox    bool
	sslCommerzBaseURL    string
}

func NewService(repository Repository, cfg *config.Config) Service {
	s := &service{
		repository: repository,
		validator:  validator.New(),
		config:     cfg,
	}
	
	// Load SSLCommerz configuration from config
	s.loadSSLCommerzConfig()
	
	return s
}

// loadSSLCommerzConfig loads SSLCommerz configuration from config
func (s *service) loadSSLCommerzConfig() {
	// Load from environment variables or use defaults
	s.sslCommerzStoreID = s.getEnvOrDefault("SSLCOMMERZ_STORE_ID", "test_store")
	s.sslCommerzStorePass = s.getEnvOrDefault("SSLCOMMERZ_STORE_PASSWORD", "test_pass")
	s.sslCommerzSandbox = s.config.App.Environment != "production"
	
	if s.sslCommerzSandbox {
		s.sslCommerzBaseURL = "https://sandbox.sslcommerz.com"
	} else {
		s.sslCommerzBaseURL = "https://securepay.sslcommerz.com"
	}
}

// getEnvOrDefault gets environment variable or returns default
func (s *service) getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getTenantIDFromContext extracts tenant ID from context
func (s *service) getTenantIDFromContext(ctx context.Context) (uuid.UUID, error) {
	if tenantID := ctx.Value("tenant_id"); tenantID != nil {
		if tenantIDStr, ok := tenantID.(string); ok {
			return uuid.Parse(tenantIDStr)
		}
		if tenantIDUUID, ok := tenantID.(uuid.UUID); ok {
			return tenantIDUUID, nil
		}
	}
	return uuid.Nil, errors.New("tenant ID not found in context")
}

func (s *service) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	// Extract tenant ID from context
	tenantID, err := s.getTenantIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant ID: %w", err)
	}
	
	// Create payment record
	payment := &Payment{
		TenantID: tenantID,
		OrderID:  orderID,
		Amount:   req.Amount,
		Currency: req.Currency,
		Status:   StatusPending,
		Gateway:  req.Gateway,
	}

	if err := s.repository.Create(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	response := &CreatePaymentResponse{
		PaymentID: payment.ID.String(),
		Status:    payment.Status,
	}

	// Initialize gateway-specific payment
	switch req.Gateway {
	case GatewaySSLCommerz:
		sslResponse, err := s.InitiateSSLCommerzPayment(ctx, payment)
		if err != nil {
			return nil, fmt.Errorf("failed to initiate SSLCommerz payment: %w", err)
		}
		response.PaymentURL = sslResponse.GatewayPageURL
		response.SessionKey = sslResponse.SessionKey
		response.GatewayPageURL = sslResponse.GatewayPageURL
		
		// Update payment with session key
		payment.PaymentIntentID = sslResponse.SessionKey
		if err := s.repository.Update(payment); err != nil {
			return nil, fmt.Errorf("failed to update payment: %w", err)
		}
		
	case GatewayBKash:
		// TODO: Implement bKash integration
		return nil, errors.New("bkash integration not implemented yet")
		
	case GatewayNagad:
		// TODO: Implement Nagad integration
		return nil, errors.New("nagad integration not implemented yet")
		
	default:
		return nil, fmt.Errorf("unsupported payment gateway: %s", req.Gateway)
	}

	return response, nil
}

func (s *service) InitiateSSLCommerzPayment(ctx context.Context, payment *Payment) (*SSLCommerzPaymentResponse, error) {
	// Build SSLCommerz request
	sslReq := &SSLCommerzPaymentRequest{
		StoreID:           s.sslCommerzStoreID,
		StorePassword:     s.sslCommerzStorePass,
		TotalAmount:       payment.Amount,
		Currency:          payment.Currency,
		TransactionID:     payment.OrderID.String(),
		SuccessURL:        "https://example.com/payment/success", // TODO: Get from config
		FailURL:           "https://example.com/payment/fail",
		CancelURL:         "https://example.com/payment/cancel",
		IPNListenerURL:    "https://example.com/webhooks/sslcommerz",
		CustomerEmail:     "customer@example.com", // TODO: Get from payment data
		CustomerPhone:     "+8801234567890",
		CustomerName:      "Customer", // TODO: Get from customer data
		CustomerAddress1:  "Dhaka",    // TODO: Get from customer data
		CustomerCity:      "Dhaka",
		CustomerState:     "Dhaka",
		CustomerPostcode:  "1000",
		CustomerCountry:   "Bangladesh",
		ShippingMethodName: "NO",
		ProductName:       "E-commerce Purchase",
		ProductCategory:   "general",
		ProductProfile:    "general",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(sslReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SSLCommerz request: %w", err)
	}

	// Make API call to SSLCommerz
	url := fmt.Sprintf("%s/gwprocess/v3/api.php", s.sslCommerzBaseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to call SSLCommerz API: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var sslResp SSLCommerzPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&sslResp); err != nil {
		return nil, fmt.Errorf("failed to decode SSLCommerz response: %w", err)
	}

	if sslResp.Status != "SUCCESS" {
		return nil, fmt.Errorf("SSLCommerz payment initiation failed: %s", sslResp.FailedReason)
	}

	return &sslResp, nil
}

func (s *service) ProcessPayment(ctx context.Context, req *ProcessPaymentRequest) (*Payment, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}

	// Extract tenant ID from context
	tenantID, err := s.getTenantIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant ID: %w", err)
	}
	
	payment, err := s.repository.GetByID(tenantID, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	// Process based on gateway
	switch req.Gateway {
	case GatewaySSLCommerz:
		// Validate with SSLCommerz
		// This would typically involve validating the transaction with SSLCommerz API
		payment.Status = StatusSucceeded
		now := time.Now()
		payment.ProcessedAt = &now
		
	default:
		return nil, fmt.Errorf("unsupported payment gateway: %s", req.Gateway)
	}

	// Store gateway response
	gatewayResponseJSON, _ := json.Marshal(req.GatewayResponse)
	payment.GatewayResponse = string(gatewayResponseJSON)

	if err := s.repository.Update(payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *service) ValidateSSLCommerzPayment(ctx context.Context, ipnData *SSLCommerzIPNResponse) error {
	// This would involve validating the IPN data with SSLCommerz
	// For now, we'll do basic validation
	if ipnData.Status == "VALID" {
		return nil
	}
	return fmt.Errorf("invalid payment status: %s", ipnData.Status)
}

func (s *service) GetPayment(ctx context.Context, paymentID string) (*Payment, error) {
	id, err := uuid.Parse(paymentID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}

	// Extract tenant ID from context
	tenantID, err := s.getTenantIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant ID: %w", err)
	}

	return s.repository.GetByID(tenantID, id)
}

func (s *service) RefundPayment(ctx context.Context, req *RefundPaymentRequest) (*Payment, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}

	// Extract tenant ID from context
	tenantID, err := s.getTenantIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant ID: %w", err)
	}

	payment, err := s.repository.GetByID(tenantID, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	if payment.Status != StatusSucceeded {
		return nil, errors.New("can only refund succeeded payments")
	}

	// Create refund record
	refund := &Refund{
		TenantID:  tenantID,
		PaymentID: paymentID,
		OrderID:   payment.OrderID,
		Amount:    req.Amount,
		Currency:  payment.Currency,
		Reason:    req.Reason,
		Status:    StatusPending,
	}

	if err := s.repository.CreateRefund(refund); err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	// Update payment
	payment.RefundedAmount += req.Amount
	if payment.RefundedAmount >= payment.Amount {
		payment.Status = StatusRefunded
		now := time.Now()
		payment.RefundedAt = &now
	}

	if err := s.repository.Update(payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *service) ListPayments(ctx context.Context, req *ListPaymentsRequest) (*ListPaymentsResponse, error) {
	// Set default pagination
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Extract tenant ID from context
	tenantID, err := s.getTenantIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant ID: %w", err)
	}
	
	// List payments with basic filtering
	payments, total, err := s.repository.List(tenantID, nil, req.Offset, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	response := &ListPaymentsResponse{
		Payments: payments,
		Total:    total,
		Offset:   req.Offset,
		Limit:    req.Limit,
	}

	// Add stats if requested
	if req.View == "stats" {
		// TODO: Implement payment stats calculation
		response.Stats = nil
	}

	return response, nil
}

func (s *service) UpdatePayment(ctx context.Context, id string, updates map[string]interface{}) (*Payment, error) {
	// Validate payment exists and user has access
	payment, err := s.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates to payment object
	for key, value := range updates {
		switch key {
		case "status":
			if status, ok := value.(string); ok {
				payment.Status = status
			}
		case "gateway_response":
			if response, ok := value.(string); ok {
				payment.GatewayResponse = response
			}
		case "payment_intent_id":
			if intentID, ok := value.(string); ok {
				payment.PaymentIntentID = intentID
			}
		}
	}

	// Update payment in repository
	if err := s.repository.Update(payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	return payment, nil
}

func (s *service) GetPaymentMethods(ctx context.Context, userID string) ([]*PaymentMethod, error) {
	// Extract tenant ID from context
	tenantID, err := s.getTenantIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	
	methods, err := s.repository.ListPaymentMethods(tenantID, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	return methods, nil
}

func (s *service) UpdatePaymentMethod(ctx context.Context, id string, req *UpdatePaymentMethodRequest) (*PaymentMethod, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(string)

	// Extract tenant ID from context
	tenantID, err := s.getTenantIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant ID: %w", err)
	}
	methodID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid method ID: %w", err)
	}
	
	// Validate payment method exists and belongs to user
	method, err := s.repository.GetPaymentMethod(tenantID, methodID)
	if err != nil {
		return nil, fmt.Errorf("payment method not found: %w", err)
	}

	if method.UserID.String() != userID {
		return nil, fmt.Errorf("unauthorized access to payment method")
	}

	// TODO: Implement default method management and updates
	// For now, just return the existing method
	return method, nil
}
