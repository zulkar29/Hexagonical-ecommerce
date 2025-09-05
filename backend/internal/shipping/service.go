package shipping

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// Request/Response Models

type CreateShippingZoneRequest struct {
	Name        string                        `json:"name" binding:"required"`
	Description string                        `json:"description"`
	IsActive    bool                          `json:"is_active"`
	IsDefault   bool                          `json:"is_default"`
	Countries   []CreateShippingZoneCountryRequest `json:"countries"`
}

type CreateShippingZoneCountryRequest struct {
	Country     string   `json:"country" binding:"required"`
	States      []string `json:"states"`
	Cities      []string `json:"cities"`
	PostalCodes []string `json:"postal_codes"`
}

type UpdateShippingZoneRequest struct {
	Name        *string                       `json:"name"`
	Description *string                       `json:"description"`
	IsActive    *bool                         `json:"is_active"`
	IsDefault   *bool                         `json:"is_default"`
	Countries   []CreateShippingZoneCountryRequest `json:"countries"`
}

type CreateShippingRateRequest struct {
	ZoneID          string           `json:"zone_id" binding:"required"`
	Provider        ShippingProvider `json:"provider" binding:"required"`
	Method          ShippingMethod   `json:"method" binding:"required"`
	Name            string           `json:"name" binding:"required"`
	Description     string           `json:"description"`
	BaseRate        float64          `json:"base_rate" binding:"required,min=0"`
	WeightRate      float64          `json:"weight_rate"`
	VolumeRate      float64          `json:"volume_rate"`
	MinWeight       float64          `json:"min_weight"`
	MaxWeight       float64          `json:"max_weight"`
	FreeShippingMin float64          `json:"free_shipping_min"`
	EstimatedDays   int              `json:"estimated_days" binding:"min=1"`
	IsActive        bool             `json:"is_active"`
	Priority        int              `json:"priority"`
}

type ShippingRateRequest struct {
	DestinationCountry string  `json:"destination_country" binding:"required"`
	DestinationState   string  `json:"destination_state"`
	DestinationCity    string  `json:"destination_city"`
	PostalCode         string  `json:"postal_code"`
	Weight             float64 `json:"weight" binding:"required,min=0"`
	Length             float64 `json:"length"`
	Width              float64 `json:"width"`
	Height             float64 `json:"height"`
	OrderValue         float64 `json:"order_value" binding:"required,min=0"`
}

type ShippingRateResponse struct {
	RateID        uuid.UUID        `json:"rate_id"`
	Provider      ShippingProvider `json:"provider"`
	Method        ShippingMethod   `json:"method"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Cost          float64          `json:"cost"`
	Currency      string           `json:"currency"`
	EstimatedDays int              `json:"estimated_days"`
	IsFree        bool             `json:"is_free"`
}

type CreateShippingLabelRequest struct {
	OrderID         string           `json:"order_id" binding:"required"`
	Provider        ShippingProvider `json:"provider" binding:"required"`
	RateID          string           `json:"rate_id" binding:"required"`
	SenderAddress   Address          `json:"sender_address" binding:"required"`
	ReceiverAddress Address          `json:"receiver_address" binding:"required"`
	PackageDetails  PackageDetails   `json:"package_details" binding:"required"`
}

type Address struct {
	Name       string `json:"name" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Email      string `json:"email"`
	Street     string `json:"street" binding:"required"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state"`
	Country    string `json:"country" binding:"required"`
	PostalCode string `json:"postal_code"`
}

type PackageDetails struct {
	Weight      float64 `json:"weight" binding:"required,min=0"`
	Length      float64 `json:"length" binding:"required,min=0"`
	Width       float64 `json:"width" binding:"required,min=0"`
	Height      float64 `json:"height" binding:"required,min=0"`
	Value       float64 `json:"value" binding:"required,min=0"`
	Description string  `json:"description" binding:"required"`
	Items       []PackageItem `json:"items"`
}

type PackageItem struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	Weight      float64 `json:"weight"`
	Value       float64 `json:"value"`
	Description string  `json:"description"`
}

type AddressValidationRequest struct {
	Address Address `json:"address" binding:"required"`
}

type AddressValidationResponse struct {
	IsValid     bool     `json:"is_valid"`
	Confidence  float64  `json:"confidence"`
	Suggestions []Address `json:"suggestions"`
	Errors      []string `json:"errors"`
}

type DeliveryEstimateRequest struct {
	FromAddress Address        `json:"from_address" binding:"required"`
	ToAddress   Address        `json:"to_address" binding:"required"`
	Weight      float64        `json:"weight" binding:"required,min=0"`
	Method      ShippingMethod `json:"method"`
}

type DeliveryEstimateResponse struct {
	Provider      ShippingProvider `json:"provider"`
	Method        ShippingMethod   `json:"method"`
	EstimatedDays int              `json:"estimated_days"`
	MinDays       int              `json:"min_days"`
	MaxDays       int              `json:"max_days"`
	DeliveryDate  time.Time        `json:"delivery_date"`
}

type ProviderConfigRequest struct {
	APIKey      string                 `json:"api_key" binding:"required"`
	APISecret   string                 `json:"api_secret"`
	SandboxMode bool                   `json:"sandbox_mode"`
	Settings    map[string]interface{} `json:"settings"`
}

type ShippingStats struct {
	TotalLabels     int                           `json:"total_labels"`
	ActiveLabels    int                           `json:"active_labels"`
	DeliveredLabels int                           `json:"delivered_labels"`
	CancelledLabels int                           `json:"cancelled_labels"`
	TotalCost       float64                       `json:"total_cost"`
	AverageDeliveryTime float64                   `json:"average_delivery_time"`
	ProviderStats   []ProviderStats               `json:"provider_stats"`
	MonthlyStats    []MonthlyShippingStats        `json:"monthly_stats"`
}

type ProviderStats struct {
	Provider     ShippingProvider `json:"provider"`
	TotalLabels  int              `json:"total_labels"`
	DeliveredRate float64         `json:"delivered_rate"`
	AverageTime  float64          `json:"average_time"`
	TotalCost    float64          `json:"total_cost"`
}

type MonthlyShippingStats struct {
	Month       string  `json:"month"`
	TotalLabels int     `json:"total_labels"`
	TotalCost   float64 `json:"total_cost"`
	DeliveredRate float64 `json:"delivered_rate"`
}

// Shipping Zone Services

func (s *Service) CreateShippingZone(tenantID uuid.UUID, req CreateShippingZoneRequest) (*ShippingZone, error) {
	zone := &ShippingZone{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		IsDefault:   req.IsDefault,
	}

	// Create countries
	for _, countryReq := range req.Countries {
		country := ShippingZoneCountry{
			Country:     countryReq.Country,
			States:      countryReq.States,
			Cities:      countryReq.Cities,
			PostalCodes: countryReq.PostalCodes,
		}
		zone.Countries = append(zone.Countries, country)
	}

	return s.repository.CreateShippingZone(zone)
}

func (s *Service) GetShippingZones(tenantID uuid.UUID) ([]ShippingZone, error) {
	return s.repository.GetShippingZones(tenantID)
}

func (s *Service) GetShippingZone(tenantID uuid.UUID, zoneID string) (*ShippingZone, error) {
	id, err := uuid.Parse(zoneID)
	if err != nil {
		return nil, errors.New("invalid zone ID")
	}
	return s.repository.GetShippingZone(tenantID, id)
}

func (s *Service) UpdateShippingZone(tenantID uuid.UUID, zoneID string, req UpdateShippingZoneRequest) (*ShippingZone, error) {
	id, err := uuid.Parse(zoneID)
	if err != nil {
		return nil, errors.New("invalid zone ID")
	}

	zone, err := s.repository.GetShippingZone(tenantID, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		zone.Name = *req.Name
	}
	if req.Description != nil {
		zone.Description = *req.Description
	}
	if req.IsActive != nil {
		zone.IsActive = *req.IsActive
	}
	if req.IsDefault != nil {
		zone.IsDefault = *req.IsDefault
	}

	// Update countries if provided
	if req.Countries != nil {
		zone.Countries = []ShippingZoneCountry{}
		for _, countryReq := range req.Countries {
			country := ShippingZoneCountry{
				ZoneID:      zone.ID,
				Country:     countryReq.Country,
				States:      countryReq.States,
				Cities:      countryReq.Cities,
				PostalCodes: countryReq.PostalCodes,
			}
			zone.Countries = append(zone.Countries, country)
		}
	}

	return s.repository.UpdateShippingZone(zone)
}

func (s *Service) DeleteShippingZone(tenantID uuid.UUID, zoneID string) error {
	id, err := uuid.Parse(zoneID)
	if err != nil {
		return errors.New("invalid zone ID")
	}
	return s.repository.DeleteShippingZone(tenantID, id)
}

// Shipping Rate Services

func (s *Service) CreateShippingRate(tenantID uuid.UUID, req CreateShippingRateRequest) (*ShippingRate, error) {
	zoneID, err := uuid.Parse(req.ZoneID)
	if err != nil {
		return nil, errors.New("invalid zone ID")
	}

	rate := &ShippingRate{
		TenantID:        tenantID,
		ZoneID:          zoneID,
		Provider:        req.Provider,
		Method:          req.Method,
		Name:            req.Name,
		Description:     req.Description,
		BaseRate:        req.BaseRate,
		WeightRate:      req.WeightRate,
		VolumeRate:      req.VolumeRate,
		MinWeight:       req.MinWeight,
		MaxWeight:       req.MaxWeight,
		FreeShippingMin: req.FreeShippingMin,
		EstimatedDays:   req.EstimatedDays,
		IsActive:        req.IsActive,
		Priority:        req.Priority,
	}

	return s.repository.CreateShippingRate(rate)
}

func (s *Service) CalculateShippingRates(tenantID uuid.UUID, req ShippingRateRequest) ([]ShippingRateResponse, error) {
	// Find applicable zones for destination
	zones, err := s.repository.GetShippingZonesForDestination(tenantID, req.DestinationCountry, req.DestinationState, req.DestinationCity)
	if err != nil {
		return nil, err
	}

	var responses []ShippingRateResponse

	for _, zone := range zones {
		rates, err := s.repository.GetShippingRatesForZone(zone.ID)
		if err != nil {
			continue
		}

		for _, rate := range rates {
			if rate.IsEligible(req.Weight) {
				cost := rate.CalculateRate(req.Weight, req.Length, req.Width, req.Height, req.OrderValue)
				
				response := ShippingRateResponse{
					RateID:        rate.ID,
					Provider:      rate.Provider,
					Method:        rate.Method,
					Name:          rate.Name,
					Description:   rate.Description,
					Cost:          cost,
					Currency:      "BDT",
					EstimatedDays: rate.EstimatedDays,
					IsFree:        cost == 0,
				}
				responses = append(responses, response)
			}
		}
	}

	return responses, nil
}

// Shipping Label Services

func (s *Service) CreateShippingLabel(tenantID uuid.UUID, req CreateShippingLabelRequest) (*ShippingLabel, error) {
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	rateID, err := uuid.Parse(req.RateID)
	if err != nil {
		return nil, errors.New("invalid rate ID")
	}

	// Get the shipping rate to calculate cost
	rate, err := s.repository.GetShippingRate(rateID)
	if err != nil {
		return nil, errors.New("shipping rate not found")
	}

	// Generate tracking number
	trackingNumber := s.generateTrackingNumber(req.Provider)

	// Calculate cost
	cost := rate.CalculateRate(
		req.PackageDetails.Weight,
		req.PackageDetails.Length,
		req.PackageDetails.Width,
		req.PackageDetails.Height,
		req.PackageDetails.Value,
	)

	label := &ShippingLabel{
		TenantID:       tenantID,
		OrderID:        orderID,
		Provider:       req.Provider,
		TrackingNumber: trackingNumber,
		Cost:           cost,
		Currency:       "BDT",
		Status:         "created",
		EstimatedDelivery: func() *time.Time {
			t := time.Now().AddDate(0, 0, rate.EstimatedDays)
			return &t
		}(),
	}

	// Call provider API to create shipping label
	err = s.createProviderLabel(label, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create label with provider: %v", err)
	}

	return s.repository.CreateShippingLabel(label)
}

func (s *Service) GetShippingLabel(tenantID uuid.UUID, labelID string) (*ShippingLabel, error) {
	id, err := uuid.Parse(labelID)
	if err != nil {
		return nil, errors.New("invalid label ID")
	}
	return s.repository.GetShippingLabel(tenantID, id)
}

func (s *Service) GetShippingLabels(tenantID uuid.UUID, offset, limit int) ([]ShippingLabel, int64, error) {
	return s.repository.GetShippingLabels(tenantID, offset, limit)
}

func (s *Service) CancelShipment(tenantID uuid.UUID, labelID string) error {
	id, err := uuid.Parse(labelID)
	if err != nil {
		return errors.New("invalid label ID")
	}

	label, err := s.repository.GetShippingLabel(tenantID, id)
	if err != nil {
		return err
	}

	if label.Status == "delivered" {
		return errors.New("cannot cancel delivered shipment")
	}

	// Call provider API to cancel
	err = s.cancelProviderLabel(label)
	if err != nil {
		return err
	}

	label.Status = "cancelled"
	_, err = s.repository.UpdateShippingLabel(label)
	return err
}

// Package Tracking Services

func (s *Service) TrackPackage(trackingNumber string) (*ShippingLabel, error) {
	label, err := s.repository.GetShippingLabelByTrackingNumber(trackingNumber)
	if err != nil {
		return nil, err
	}

	// Update tracking from provider
	err = s.updateTrackingFromProvider(label)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to update tracking from provider: %v\n", err)
	}

	return label, nil
}

func (s *Service) GetTrackingHistory(trackingNumber string) ([]ShippingTracking, error) {
	label, err := s.repository.GetShippingLabelByTrackingNumber(trackingNumber)
	if err != nil {
		return nil, err
	}

	return s.repository.GetTrackingHistory(label.ID)
}

// Address Validation Services

func (s *Service) ValidateAddress(req AddressValidationRequest) (*AddressValidationResponse, error) {
	// Basic validation logic
	response := &AddressValidationResponse{
		IsValid:    true,
		Confidence: 0.8,
	}

	var errors []string

	if req.Address.Name == "" {
		errors = append(errors, "Name is required")
	}
	if req.Address.Street == "" {
		errors = append(errors, "Street address is required")
	}
	if req.Address.City == "" {
		errors = append(errors, "City is required")
	}
	if req.Address.Country == "" {
		errors = append(errors, "Country is required")
	}

	if len(errors) > 0 {
		response.IsValid = false
		response.Confidence = 0.0
		response.Errors = errors
	}

	return response, nil
}

// Delivery Estimate Services

func (s *Service) GetDeliveryEstimate(req DeliveryEstimateRequest) (*DeliveryEstimateResponse, error) {
	// Basic estimation logic
	estimatedDays := 3 // Default

	// Adjust based on method
	switch req.Method {
	case MethodExpress:
		estimatedDays = 1
	case MethodSameDay:
		estimatedDays = 1
	case MethodStandard:
		estimatedDays = 3
	}

	// Adjust based on distance (simplified)
	if req.FromAddress.Country != req.ToAddress.Country {
		estimatedDays += 5
	} else if req.FromAddress.City != req.ToAddress.City {
		estimatedDays += 1
	}

	response := &DeliveryEstimateResponse{
		Provider:      ProviderPathao, // Default
		Method:        req.Method,
		EstimatedDays: estimatedDays,
		MinDays:       estimatedDays - 1,
		MaxDays:       estimatedDays + 2,
		DeliveryDate:  time.Now().AddDate(0, 0, estimatedDays),
	}

	return response, nil
}

// Provider Management Services

func (s *Service) GetShippingProviders(tenantID uuid.UUID) ([]ShippingProviderConfig, error) {
	return s.repository.GetShippingProviders(tenantID)
}

func (s *Service) ConfigureProvider(tenantID uuid.UUID, providerName string, req ProviderConfigRequest) error {
	provider := &ShippingProviderConfig{
		TenantID:    tenantID,
		Provider:    ShippingProvider(providerName),
		Name:        providerName,
		IsActive:    true,
		APIKey:      req.APIKey,
		APISecret:   req.APISecret,
		SandboxMode: req.SandboxMode,
		Settings:    req.Settings,
	}

	return s.repository.CreateOrUpdateProvider(provider)
}

// Statistics Services

func (s *Service) GetShippingStats(tenantID uuid.UUID) (*ShippingStats, error) {
	return s.repository.GetShippingStats(tenantID)
}

func (s *Service) GetShippingHistory(tenantID uuid.UUID, offset, limit int) ([]ShippingLabel, int64, error) {
	return s.repository.GetShippingLabels(tenantID, offset, limit)
}

// Helper Methods

func (s *Service) generateTrackingNumber(provider ShippingProvider) string {
	// Generate unique tracking number based on provider format
	timestamp := time.Now().Unix()
	switch provider {
	case ProviderPathao:
		return fmt.Sprintf("PAT%d", timestamp)
	case ProviderRedX:
		return fmt.Sprintf("RDX%d", timestamp)
	case ProviderPaperfly:
		return fmt.Sprintf("PPF%d", timestamp)
	case ProviderDHL:
		return fmt.Sprintf("DHL%d", timestamp)
	case ProviderFedEx:
		return fmt.Sprintf("FDX%d", timestamp)
	default:
		return fmt.Sprintf("SHP%d", timestamp)
	}
}

// Provider Integration Methods (stubs for now)

func (s *Service) createProviderLabel(label *ShippingLabel, req CreateShippingLabelRequest) error {
	// TODO: Implement provider-specific API calls
	switch label.Provider {
	case ProviderPathao:
		return s.createPathaoLabel(label, req)
	case ProviderRedX:
		return s.createRedXLabel(label, req)
	case ProviderPaperfly:
		return s.createPaperflyLabel(label, req)
	case ProviderDHL:
		return s.createDHLLabel(label, req)
	case ProviderFedEx:
		return s.createFedExLabel(label, req)
	default:
		return errors.New("unsupported provider")
	}
}

func (s *Service) createPathaoLabel(label *ShippingLabel, req CreateShippingLabelRequest) error {
	// TODO: Implement Pathao API integration
	label.LabelURL = "https://pathao.com/label/" + label.TrackingNumber
	label.ProviderOrderID = "PATHAO-" + label.TrackingNumber
	return nil
}

func (s *Service) createRedXLabel(label *ShippingLabel, req CreateShippingLabelRequest) error {
	// TODO: Implement RedX API integration
	label.LabelURL = "https://redx.com.bd/label/" + label.TrackingNumber
	label.ProviderOrderID = "REDX-" + label.TrackingNumber
	return nil
}

func (s *Service) createPaperflyLabel(label *ShippingLabel, req CreateShippingLabelRequest) error {
	// TODO: Implement Paperfly API integration
	label.LabelURL = "https://paperfly.com.bd/label/" + label.TrackingNumber
	label.ProviderOrderID = "PPF-" + label.TrackingNumber
	return nil
}

func (s *Service) createDHLLabel(label *ShippingLabel, req CreateShippingLabelRequest) error {
	// TODO: Implement DHL API integration
	label.LabelURL = "https://dhl.com/label/" + label.TrackingNumber
	label.ProviderOrderID = "DHL-" + label.TrackingNumber
	return nil
}

func (s *Service) createFedExLabel(label *ShippingLabel, req CreateShippingLabelRequest) error {
	// TODO: Implement FedEx API integration
	label.LabelURL = "https://fedex.com/label/" + label.TrackingNumber
	label.ProviderOrderID = "FDX-" + label.TrackingNumber
	return nil
}

func (s *Service) cancelProviderLabel(label *ShippingLabel) error {
	// TODO: Implement provider-specific cancellation
	return nil
}

func (s *Service) updateTrackingFromProvider(label *ShippingLabel) error {
	// TODO: Implement provider-specific tracking updates
	return nil
}

// Webhook Processing Methods

func (s *Service) ProcessPathaoWebhook(payload map[string]interface{}) error {
	// Extract tracking number and status from Pathao webhook
	trackingNumber, ok := payload["tracking_number"].(string)
	if !ok {
		return errors.New("missing tracking_number in webhook payload")
	}

	status, ok := payload["status"].(string)
	if !ok {
		return errors.New("missing status in webhook payload")
	}

	// Find shipping label by tracking number
	label, err := s.repository.GetShippingLabelByTrackingNumber(trackingNumber)
	if err != nil {
		return fmt.Errorf("failed to find shipping label: %w", err)
	}

	// Update tracking status based on Pathao status
	var trackingStatus TrackingStatus
	switch status {
	case "picked_up":
		trackingStatus = StatusPickedUp
	case "in_transit":
		trackingStatus = StatusInTransit
	case "out_for_delivery":
		trackingStatus = StatusOutForDelivery
	case "delivered":
		trackingStatus = StatusDelivered
	case "failed":
		trackingStatus = StatusFailed
	case "returned":
		trackingStatus = StatusReturned
	default:
		trackingStatus = StatusInTransit
	}

	// Create tracking update
	tracking := &ShippingTracking{
		ID:          uuid.New(),
		LabelID:     label.ID,
		Status:      string(trackingStatus),
		Location:    getStringFromPayload(payload, "location"),
		Description: getStringFromPayload(payload, "description"),
		Timestamp:   time.Now(),
		CreatedAt:   time.Now(),
	}

	_, err = s.repository.CreateShippingTracking(tracking)
	return err
}

func (s *Service) ProcessRedXWebhook(payload map[string]interface{}) error {
	// Extract tracking number and status from RedX webhook
	trackingNumber, ok := payload["tracking_id"].(string)
	if !ok {
		return errors.New("missing tracking_id in webhook payload")
	}

	status, ok := payload["delivery_status"].(string)
	if !ok {
		return errors.New("missing delivery_status in webhook payload")
	}

	// Find shipping label by tracking number
	label, err := s.repository.GetShippingLabelByTrackingNumber(trackingNumber)
	if err != nil {
		return fmt.Errorf("failed to find shipping label: %w", err)
	}

	// Update tracking status based on RedX status
	var trackingStatus TrackingStatus
	switch status {
	case "PICKED_UP":
		trackingStatus = StatusPickedUp
	case "IN_TRANSIT":
		trackingStatus = StatusInTransit
	case "OUT_FOR_DELIVERY":
		trackingStatus = StatusOutForDelivery
	case "DELIVERED":
		trackingStatus = StatusDelivered
	case "DELIVERY_FAILED":
		trackingStatus = StatusFailed
	case "RETURNED":
		trackingStatus = StatusReturned
	default:
		trackingStatus = StatusInTransit
	}

	// Create tracking update
	tracking := &ShippingTracking{
		ID:          uuid.New(),
		LabelID:     label.ID,
		Status:      string(trackingStatus),
		Location:    getStringFromPayload(payload, "current_location"),
		Description: getStringFromPayload(payload, "remarks"),
		Timestamp:   time.Now(),
		CreatedAt:   time.Now(),
	}

	_, err = s.repository.CreateShippingTracking(tracking)
	return err
}

func (s *Service) ProcessPaperflyWebhook(payload map[string]interface{}) error {
	// Extract tracking number and status from Paperfly webhook
	trackingNumber, ok := payload["consignment_id"].(string)
	if !ok {
		return errors.New("missing consignment_id in webhook payload")
	}

	status, ok := payload["status"].(string)
	if !ok {
		return errors.New("missing status in webhook payload")
	}

	// Find shipping label by tracking number
	label, err := s.repository.GetShippingLabelByTrackingNumber(trackingNumber)
	if err != nil {
		return fmt.Errorf("failed to find shipping label: %w", err)
	}

	// Update tracking status based on Paperfly status
	var trackingStatus TrackingStatus
	switch status {
	case "picked":
		trackingStatus = StatusPickedUp
	case "in_transit":
		trackingStatus = StatusInTransit
	case "out_for_delivery":
		trackingStatus = StatusOutForDelivery
	case "delivered":
		trackingStatus = StatusDelivered
	case "failed":
		trackingStatus = StatusFailed
	case "returned":
		trackingStatus = StatusReturned
	default:
		trackingStatus = StatusInTransit
	}

	// Create tracking update
	tracking := &ShippingTracking{
		ID:          uuid.New(),
		LabelID:     label.ID,
		Status:      string(trackingStatus),
		Location:    getStringFromPayload(payload, "location"),
		Description: getStringFromPayload(payload, "note"),
		Timestamp:   time.Now(),
		CreatedAt:   time.Now(),
	}

	_, err = s.repository.CreateShippingTracking(tracking)
	return err
}

func (s *Service) ProcessDHLWebhook(payload map[string]interface{}) error {
	// Extract tracking number and status from DHL webhook
	trackingNumber, ok := payload["trackingNumber"].(string)
	if !ok {
		return errors.New("missing trackingNumber in webhook payload")
	}

	// DHL sends events array
	events, ok := payload["events"].([]interface{})
	if !ok || len(events) == 0 {
		return errors.New("missing events in webhook payload")
	}

	// Get the latest event
	latestEvent, ok := events[0].(map[string]interface{})
	if !ok {
		return errors.New("invalid event format in webhook payload")
	}

	status, ok := latestEvent["status"].(string)
	if !ok {
		return errors.New("missing status in event")
	}

	// Find shipping label by tracking number
	label, err := s.repository.GetShippingLabelByTrackingNumber(trackingNumber)
	if err != nil {
		return fmt.Errorf("failed to find shipping label: %w", err)
	}

	// Update tracking status based on DHL status
	var trackingStatus TrackingStatus
	switch status {
	case "transit":
		trackingStatus = StatusInTransit
	case "delivered":
		trackingStatus = StatusDelivered
	case "failure":
		trackingStatus = StatusFailed
	case "unknown":
		trackingStatus = StatusPending
	default:
		trackingStatus = StatusInTransit
	}

	// Create tracking update
	tracking := &ShippingTracking{
		ID:          uuid.New(),
		LabelID:     label.ID,
		Status:      string(trackingStatus),
		Location:    getStringFromPayload(latestEvent, "location"),
		Description: getStringFromPayload(latestEvent, "description"),
		Timestamp:   time.Now(),
		CreatedAt:   time.Now(),
	}

	_, err = s.repository.CreateShippingTracking(tracking)
	return err
}

func (s *Service) ProcessFedExWebhook(payload map[string]interface{}) error {
	// Extract tracking number and status from FedEx webhook
	trackingNumber, ok := payload["trackingNumber"].(string)
	if !ok {
		return errors.New("missing trackingNumber in webhook payload")
	}

	// FedEx sends scanEvents array
	scanEvents, ok := payload["scanEvents"].([]interface{})
	if !ok || len(scanEvents) == 0 {
		return errors.New("missing scanEvents in webhook payload")
	}

	// Get the latest scan event
	latestScan, ok := scanEvents[0].(map[string]interface{})
	if !ok {
		return errors.New("invalid scanEvent format in webhook payload")
	}

	status, ok := latestScan["eventType"].(string)
	if !ok {
		return errors.New("missing eventType in scanEvent")
	}

	// Find shipping label by tracking number
	label, err := s.repository.GetShippingLabelByTrackingNumber(trackingNumber)
	if err != nil {
		return fmt.Errorf("failed to find shipping label: %w", err)
	}

	// Update tracking status based on FedEx status
	var trackingStatus TrackingStatus
	switch status {
	case "PU":
		trackingStatus = StatusPickedUp
	case "IT":
		trackingStatus = StatusInTransit
	case "OD":
		trackingStatus = StatusOutForDelivery
	case "DL":
		trackingStatus = StatusDelivered
	case "DE":
		trackingStatus = StatusFailed
	default:
		trackingStatus = StatusInTransit
	}

	// Create tracking update
	tracking := &ShippingTracking{
		ID:          uuid.New(),
		LabelID:     label.ID,
		Status:      string(trackingStatus),
		Location:    getStringFromPayload(latestScan, "scanLocation"),
		Description: getStringFromPayload(latestScan, "eventDescription"),
		Timestamp:   time.Now(),
		CreatedAt:   time.Now(),
	}

	_, err = s.repository.CreateShippingTracking(tracking)
	return err
}

// Helper function to safely extract string from payload
func getStringFromPayload(payload map[string]interface{}, key string) string {
	if value, ok := payload[key].(string); ok {
		return value
	}
	return ""
}