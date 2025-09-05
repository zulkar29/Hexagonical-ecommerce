package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	CreatePayment(tenantID uuid.UUID, req *CreatePaymentRequest) (*CreatePaymentResponse, error)
	ProcessPayment(tenantID uuid.UUID, req *ProcessPaymentRequest) error
	GetPayment(tenantID uuid.UUID, paymentID string) (*Payment, error)
	RefundPayment(tenantID uuid.UUID, req *RefundPaymentRequest) error
	ListPayments(tenantID uuid.UUID, orderID *uuid.UUID, offset, limit int) ([]*Payment, int64, error)
	
	// SSLCommerz specific
	InitiateSSLCommerzPayment(tenantID uuid.UUID, req *CreatePaymentRequest) (*SSLCommerzPaymentResponse, error)
	ValidateSSLCommerzPayment(ipnData *SSLCommerzIPNResponse) error
}

type service struct {
	repository Repository
	validator  *validator.Validate
	
	// SSLCommerz configuration
	sslCommerzStoreID    string
	sslCommerzStorePass  string
	sslCommerzSandbox    bool
	sslCommerzBaseURL    string
}

func NewService(repository Repository) Service {
	return &service{
		repository:        repository,
		validator:         validator.New(),
		sslCommerzStoreID: "test_store", // TODO: Load from config
		sslCommerzStorePass: "test_pass", // TODO: Load from config
		sslCommerzSandbox: true,
		sslCommerzBaseURL: "https://sandbox.sslcommerz.com", // Sandbox URL
	}
}

func (s *service) CreatePayment(tenantID uuid.UUID, req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
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
		sslResponse, err := s.InitiateSSLCommerzPayment(tenantID, req)
		if err != nil {
			return nil, fmt.Errorf("failed to initiate SSLCommerz payment: %w", err)
		}
		response.PaymentURL = sslResponse.GatewayPageURL
		response.SessionKey = sslResponse.SessionKey
		response.GatewayPageURL = sslResponse.GatewayPageURL
		
		// Update payment with session key
		payment.PaymentIntentID = sslResponse.SessionKey
		s.repository.Update(payment)
		
	case GatewayBKash:
		// TODO: Implement bKash integration
		return nil, errors.New("bKash integration not implemented yet")
		
	case GatewayNagad:
		// TODO: Implement Nagad integration
		return nil, errors.New("Nagad integration not implemented yet")
		
	default:
		return nil, fmt.Errorf("unsupported payment gateway: %s", req.Gateway)
	}

	return response, nil
}

func (s *service) InitiateSSLCommerzPayment(tenantID uuid.UUID, req *CreatePaymentRequest) (*SSLCommerzPaymentResponse, error) {
	// Build SSLCommerz request
	sslReq := &SSLCommerzPaymentRequest{
		StoreID:           s.sslCommerzStoreID,
		StorePassword:     s.sslCommerzStorePass,
		TotalAmount:       req.Amount,
		Currency:          req.Currency,
		TransactionID:     req.OrderID,
		SuccessURL:        fmt.Sprintf("%s/payment/success", req.ReturnURL),
		FailURL:           fmt.Sprintf("%s/payment/fail", req.ReturnURL),
		CancelURL:         fmt.Sprintf("%s/payment/cancel", req.ReturnURL),
		IPNListenerURL:    fmt.Sprintf("%s/webhooks/sslcommerz", req.ReturnURL),
		CustomerEmail:     req.CustomerEmail,
		CustomerPhone:     req.CustomerPhone,
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

func (s *service) ProcessPayment(tenantID uuid.UUID, req *ProcessPaymentRequest) error {
	if err := s.validator.Struct(req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return fmt.Errorf("invalid payment ID: %w", err)
	}

	payment, err := s.repository.GetByID(tenantID, paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
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
		return fmt.Errorf("unsupported payment gateway: %s", req.Gateway)
	}

	// Store gateway response
	gatewayResponseJSON, _ := json.Marshal(req.GatewayResponse)
	payment.GatewayResponse = string(gatewayResponseJSON)

	return s.repository.Update(payment)
}

func (s *service) ValidateSSLCommerzPayment(ipnData *SSLCommerzIPNResponse) error {
	// This would involve validating the IPN data with SSLCommerz
	// For now, we'll do basic validation
	if ipnData.Status == "VALID" {
		return nil
	}
	return fmt.Errorf("invalid payment status: %s", ipnData.Status)
}

func (s *service) GetPayment(tenantID uuid.UUID, paymentID string) (*Payment, error) {
	id, err := uuid.Parse(paymentID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}

	return s.repository.GetByID(tenantID, id)
}

func (s *service) RefundPayment(tenantID uuid.UUID, req *RefundPaymentRequest) error {
	if err := s.validator.Struct(req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return fmt.Errorf("invalid payment ID: %w", err)
	}

	payment, err := s.repository.GetByID(tenantID, paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	if payment.Status != StatusSucceeded {
		return errors.New("can only refund succeeded payments")
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
		return fmt.Errorf("failed to create refund: %w", err)
	}

	// Update payment
	payment.RefundedAmount += req.Amount
	if payment.RefundedAmount >= payment.Amount {
		payment.Status = StatusRefunded
		now := time.Now()
		payment.RefundedAt = &now
	}

	return s.repository.Update(payment)
}

func (s *service) ListPayments(tenantID uuid.UUID, orderID *uuid.UUID, offset, limit int) ([]*Payment, int64, error) {
	return s.repository.List(tenantID, orderID, offset, limit)
}
