package payment

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// InitiateSSLCommerzPayment initiates a payment with SSLCommerz
func (s *service) InitiateSSLCommerzPayment(ctx context.Context, payment *Payment) (*SSLCommerzPaymentResponse, error) {
	// Build payment request
	req := &SSLCommerzPaymentRequest{
		StoreID:         s.sslCommerzStoreID,
		StorePassword:   s.sslCommerzStorePass,
		TotalAmount:     payment.Amount,
		Currency:        payment.Currency,
		TransactionID:   payment.ID.String(),
		SuccessURL:      s.config.App.BaseURL + "/api/v1/payments/sslcommerz/success",
		FailURL:         s.config.App.BaseURL + "/api/v1/payments/sslcommerz/fail",
		CancelURL:       s.config.App.BaseURL + "/api/v1/payments/sslcommerz/cancel",
		IPNListenerURL:  s.config.App.BaseURL + "/api/v1/payments/sslcommerz/ipn",
		ProductName:     "Order Payment",
		ProductCategory: "ecommerce",
		ProductProfile:  "general",
	}

	// Add customer information if available
	// In a real implementation, you'd get this from the order/customer
	req.CustomerName = "Customer"
	req.CustomerEmail = "customer@example.com"
	req.CustomerPhone = "01700000000"
	req.CustomerAddress1 = "Dhaka"
	req.CustomerCity = "Dhaka"
	req.CustomerState = "Dhaka"
	req.CustomerPostcode = "1000"
	req.CustomerCountry = "Bangladesh"

	// Convert to form data
	formData := url.Values{}
	formData.Set("store_id", req.StoreID)
	formData.Set("store_passwd", req.StorePassword)
	formData.Set("total_amount", fmt.Sprintf("%.2f", req.TotalAmount))
	formData.Set("currency", req.Currency)
	formData.Set("tran_id", req.TransactionID)
	formData.Set("success_url", req.SuccessURL)
	formData.Set("fail_url", req.FailURL)
	formData.Set("cancel_url", req.CancelURL)
	formData.Set("ipn_url", req.IPNListenerURL)
	formData.Set("cus_name", req.CustomerName)
	formData.Set("cus_email", req.CustomerEmail)
	formData.Set("cus_phone", req.CustomerPhone)
	formData.Set("cus_add1", req.CustomerAddress1)
	formData.Set("cus_city", req.CustomerCity)
	formData.Set("cus_state", req.CustomerState)
	formData.Set("cus_postcode", req.CustomerPostcode)
	formData.Set("cus_country", req.CustomerCountry)
	formData.Set("shipping_method", req.ShippingMethodName)
	formData.Set("product_name", req.ProductName)
	formData.Set("product_category", req.ProductCategory)
	formData.Set("product_profile", req.ProductProfile)

	// Make request to SSLCommerz
	initURL := s.sslCommerzBaseURL + "/gwprocess/v3/api.php"
	resp, err := s.httpClient.PostForm(initURL, formData)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate SSLCommerz payment: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSLCommerz response: %w", err)
	}

	var sslResponse SSLCommerzPaymentResponse
	if err := json.Unmarshal(body, &sslResponse); err != nil {
		return nil, fmt.Errorf("failed to parse SSLCommerz response: %w", err)
	}

	// Check if initialization was successful
	if sslResponse.Status != "SUCCESS" {
		return nil, fmt.Errorf("SSLCommerz payment initialization failed: %s", sslResponse.FailedReason)
	}

	return &sslResponse, nil
}

// ValidateSSLCommerzPayment validates payment data from SSLCommerz IPN/callback
func (s *service) ValidateSSLCommerzPayment(ctx context.Context, ipnData *SSLCommerzIPNResponse) error {
	// Step 1: Basic validation
	if ipnData.ValID == "" || ipnData.TransactionID == "" {
		return fmt.Errorf("missing required fields in IPN data")
	}

	// Step 2: Verify payment status
	if ipnData.Status != "VALID" && ipnData.Status != "VALIDATED" {
		return fmt.Errorf("invalid payment status: %s", ipnData.Status)
	}

	// Step 3: Validate transaction with SSLCommerz server
	if err := s.validateWithSSLCommerzServer(ipnData.ValID); err != nil {
		return fmt.Errorf("server validation failed: %w", err)
	}

	// Step 4: Find and update the payment record
	tenantID, err := s.getTenantIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tenant ID: %w", err)
	}

	paymentID, err := uuid.Parse(ipnData.TransactionID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID format: %w", err)
	}

	// Get payment from database
	payment, err := s.repository.GetPaymentByID(tenantID, paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// Step 5: Verify payment amount matches
	if payment.Amount != ipnData.Amount {
		return fmt.Errorf("amount mismatch: expected %.2f, got %.2f", payment.Amount, ipnData.Amount)
	}

	// Step 6: Verify payment currency matches
	if payment.Currency != ipnData.Currency {
		return fmt.Errorf("currency mismatch: expected %s, got %s", payment.Currency, ipnData.Currency)
	}

	// Step 7: Check if payment is already processed
	if payment.Status == StatusSucceeded {
		return nil // Already processed, ignore
	}

	// Step 8: Update payment status
	updates := map[string]interface{}{
		"status":            StatusSucceeded,
		"gateway_response":  s.formatGatewayResponse(ipnData),
		"payment_intent_id": ipnData.ValID,
		"processed_at":      time.Now(),
		"updated_at":        time.Now(),
	}

	if _, err := s.repository.UpdatePayment(tenantID, paymentID, updates); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// validateWithSSLCommerzServer validates payment with SSLCommerz server
func (s *service) validateWithSSLCommerzServer(valID string) error {
	// Validation API endpoint
	validateURL := s.sslCommerzBaseURL + "/validator/api/validationserverAPI.php"

	// Build validation request
	formData := url.Values{}
	formData.Set("val_id", valID)
	formData.Set("store_id", s.sslCommerzStoreID)
	formData.Set("store_passwd", s.sslCommerzStorePass)
	formData.Set("v", "1")
	formData.Set("format", "json")

	// Make validation request
	resp, err := s.httpClient.PostForm(validateURL, formData)
	if err != nil {
		return fmt.Errorf("failed to validate with SSLCommerz server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SSLCommerz validation returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read validation response: %w", err)
	}

	// Parse validation response
	var validationResp map[string]interface{}
	if err := json.Unmarshal(body, &validationResp); err != nil {
		return fmt.Errorf("failed to parse validation response: %w", err)
	}

	// Check validation status
	status, exists := validationResp["status"]
	if !exists {
		return fmt.Errorf("validation response missing status")
	}

	if status != "VALID" && status != "VALIDATED" {
		return fmt.Errorf("payment validation failed with status: %v", status)
	}

	return nil
}

// formatGatewayResponse formats IPN data for storage
func (s *service) formatGatewayResponse(ipnData *SSLCommerzIPNResponse) string {
	response := map[string]interface{}{
		"status":              ipnData.Status,
		"val_id":              ipnData.ValID,
		"tran_id":             ipnData.TransactionID,
		"amount":              ipnData.Amount,
		"currency":            ipnData.Currency,
		"store_amount":        ipnData.StoreAmount,
		"card_type":           ipnData.CardType,
		"card_no":             ipnData.CardNo,
		"bank_tran_id":        ipnData.BankTransactionID,
		"tran_date":           ipnData.TransactionDate,
		"currency_type":       ipnData.CurrencyType,
		"currency_amount":     ipnData.CurrencyAmount,
		"currency_rate":       ipnData.CurrencyRate,
		"base_amount":         ipnData.BaseAmount,
		"verify_sign":         ipnData.VerifySign,
	}

	jsonData, _ := json.Marshal(response)
	return string(jsonData)
}

// generateSSLCommerzHash generates verification hash for SSLCommerz
func (s *service) generateSSLCommerzHash(data map[string]string) string {
	// Sort keys and concatenate values
	var values []string
	keys := []string{"store_id", "store_passwd", "total_amount", "currency", "tran_id"}

	for _, key := range keys {
		if value, exists := data[key]; exists {
			values = append(values, value)
		}
	}

	// Create hash string
	hashString := strings.Join(values, "")
	hash := md5.Sum([]byte(hashString))
	return fmt.Sprintf("%x", hash)
}

// HandleSSLCommerzSuccess handles successful payment callback
func (s *service) HandleSSLCommerzSuccess(ctx context.Context, formData map[string]string) error {
	// Extract IPN data from form
	ipnData := &SSLCommerzIPNResponse{
		Status:           formData["status"],
		TransactionID:    formData["tran_id"],
		ValID:            formData["val_id"],
		CardType:         formData["card_type"],
		CardNo:           formData["card_no"],
		BankTransactionID: formData["bank_tran_id"],
		TransactionDate:  formData["tran_date"],
		CurrencyType:     formData["currency_type"],
		VerifySign:       formData["verify_sign"],
	}

	// Parse numeric fields
	if amount, err := strconv.ParseFloat(formData["amount"], 64); err == nil {
		ipnData.Amount = amount
	}
	if storeAmount, err := strconv.ParseFloat(formData["store_amount"], 64); err == nil {
		ipnData.StoreAmount = storeAmount
	}
	if currencyAmount, err := strconv.ParseFloat(formData["currency_amount"], 64); err == nil {
		ipnData.CurrencyAmount = currencyAmount
	}
	if currencyRate, err := strconv.ParseFloat(formData["currency_rate"], 64); err == nil {
		ipnData.CurrencyRate = currencyRate
	}
	if baseAmount, err := strconv.ParseFloat(formData["base_amount"], 64); err == nil {
		ipnData.BaseAmount = baseAmount
	}

	ipnData.Currency = formData["currency"]

	// Validate the payment
	return s.ValidateSSLCommerzPayment(ctx, ipnData)
}

// HandleSSLCommerzFail handles failed payment callback
func (s *service) HandleSSLCommerzFail(ctx context.Context, formData map[string]string) error {
	transactionID := formData["tran_id"]
	if transactionID == "" {
		return fmt.Errorf("missing transaction ID in failure callback")
	}

	tenantID, err := s.getTenantIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tenant ID: %w", err)
	}

	paymentID, err := uuid.Parse(transactionID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID format: %w", err)
	}

	// Update payment status to failed
	updates := map[string]interface{}{
		"status":         StatusFailed,
		"failure_reason": formData["error"],
		"updated_at":     time.Now(),
	}

	_, err = s.repository.UpdatePayment(tenantID, paymentID, updates)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// HandleSSLCommerzCancel handles cancelled payment callback
func (s *service) HandleSSLCommerzCancel(ctx context.Context, formData map[string]string) error {
	transactionID := formData["tran_id"]
	if transactionID == "" {
		return fmt.Errorf("missing transaction ID in cancel callback")
	}

	tenantID, err := s.getTenantIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tenant ID: %w", err)
	}

	paymentID, err := uuid.Parse(transactionID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID format: %w", err)
	}

	// Update payment status to cancelled
	updates := map[string]interface{}{
		"status":     StatusCancelled,
		"updated_at": time.Now(),
	}

	_, err = s.repository.UpdatePayment(tenantID, paymentID, updates)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}