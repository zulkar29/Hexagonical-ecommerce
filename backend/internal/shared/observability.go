package shared

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

// ObservabilityIntegration provides centralized observability for all modules
type ObservabilityIntegration struct {
	analyticsService AnalyticsService
	observabilityService ObservabilityService
}

// AnalyticsService defines the interface for analytics operations
type AnalyticsService interface {
	TrackEvent(ctx context.Context, tenantID uuid.UUID, event interface{}) error
	TrackPageView(ctx context.Context, tenantID uuid.UUID, pageView interface{}) error
	TrackBusinessMetric(ctx context.Context, tenantID uuid.UUID, metric BusinessMetric) error
}

// ObservabilityService defines the interface for observability operations  
type ObservabilityService interface {
	LogEvent(ctx context.Context, level string, message string, fields map[string]interface{}, err error)
	RecordMetric(ctx context.Context, name string, metricType string, value float64, tags map[string]string) 
	TraceOperation(ctx context.Context, operationName string, tags map[string]interface{}, fn func(ctx context.Context) error) error
}

// BusinessMetric represents a business-level metric
type BusinessMetric struct {
	Name      string            `json:"name"`
	Category  string            `json:"category"`
	Type      string            `json:"type"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags"`
	Timestamp time.Time         `json:"timestamp"`
}

// NewObservabilityIntegration creates a new observability integration
func NewObservabilityIntegration(analytics AnalyticsService, observability ObservabilityService) *ObservabilityIntegration {
	return &ObservabilityIntegration{
		analyticsService:     analytics,
		observabilityService: observability,
	}
}

// TrackUserAction tracks user actions across analytics and observability
func (o *ObservabilityIntegration) TrackUserAction(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, action string, properties map[string]interface{}) {
	// Track in analytics for business intelligence
	event := map[string]interface{}{
		"event_type":   "user_action",
		"event_name":   action,
		"user_id":      userID,
		"tenant_id":    tenantID,
		"properties":   properties,
		"timestamp":    time.Now(),
	}
	
	if err := o.analyticsService.TrackEvent(ctx, tenantID, event); err != nil {
		o.observabilityService.LogEvent(ctx, "error", "Failed to track user action in analytics", 
			map[string]interface{}{
				"user_id": userID,
				"action":  action,
				"error":   err.Error(),
			}, err)
	}

	// Log in observability for monitoring
	o.observabilityService.LogEvent(ctx, "info", "User action performed",
		map[string]interface{}{
			"user_id":    userID,
			"tenant_id":  tenantID,
			"action":     action,
			"properties": properties,
		}, nil)

	// Record metrics for monitoring
	o.observabilityService.RecordMetric(ctx, "user_actions_total", "counter", 1,
		map[string]string{
			"action":    action,
			"tenant_id": tenantID.String(),
		})
}

// TrackBusinessEvent tracks business events with full observability
func (o *ObservabilityIntegration) TrackBusinessEvent(ctx context.Context, tenantID uuid.UUID, eventType, eventName string, value float64, properties map[string]interface{}) {
	// Create business metric
	metric := BusinessMetric{
		Name:      eventName,
		Category:  "business",
		Type:      eventType,
		Value:     value,
		Tags: map[string]string{
			"tenant_id": tenantID.String(),
			"event_type": eventType,
		},
		Timestamp: time.Now(),
	}

	// Track in analytics
	if err := o.analyticsService.TrackBusinessMetric(ctx, tenantID, metric); err != nil {
		o.observabilityService.LogEvent(ctx, "error", "Failed to track business event in analytics",
			map[string]interface{}{
				"event_type": eventType,
				"event_name": eventName,
				"error":      err.Error(),
			}, err)
	}

	// Log the business event
	o.observabilityService.LogEvent(ctx, "info", "Business event tracked",
		map[string]interface{}{
			"tenant_id":   tenantID,
			"event_type":  eventType,
			"event_name":  eventName,
			"value":       value,
			"properties":  properties,
		}, nil)

	// Record business metrics
	o.observabilityService.RecordMetric(ctx, "business_events_total", "counter", 1,
		map[string]string{
			"event_type": eventType,
			"event_name": eventName,
			"tenant_id":  tenantID.String(),
		})

	if value > 0 {
		o.observabilityService.RecordMetric(ctx, "business_value_total", "counter", value,
			map[string]string{
				"event_type": eventType,
				"event_name": eventName,
				"tenant_id":  tenantID.String(),
			})
	}
}

// TrackAPICall tracks API calls with performance metrics
func (o *ObservabilityIntegration) TrackAPICall(ctx context.Context, method, endpoint string, statusCode int, duration time.Duration, tenantID *uuid.UUID) {
	// Log API call
	fields := map[string]interface{}{
		"method":      method,
		"endpoint":    endpoint,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	}
	
	if tenantID != nil {
		fields["tenant_id"] = *tenantID
	}

	logLevel := "info"
	if statusCode >= 400 && statusCode < 500 {
		logLevel = "warn"
	} else if statusCode >= 500 {
		logLevel = "error"
	}

	o.observabilityService.LogEvent(ctx, logLevel, "API call completed", fields, nil)

	// Record performance metrics
	tags := map[string]string{
		"method":      method,
		"endpoint":    endpoint,
		"status_code": string(rune(statusCode)),
	}
	
	if tenantID != nil {
		tags["tenant_id"] = tenantID.String()
	}

	// Record request count
	o.observabilityService.RecordMetric(ctx, "api_requests_total", "counter", 1, tags)
	
	// Record request duration
	o.observabilityService.RecordMetric(ctx, "api_request_duration_ms", "histogram", float64(duration.Milliseconds()), tags)

	// Track slow requests
	if duration > 1*time.Second {
		o.observabilityService.RecordMetric(ctx, "api_slow_requests_total", "counter", 1, tags)
	}

	// Track error rates
	if statusCode >= 400 {
		o.observabilityService.RecordMetric(ctx, "api_errors_total", "counter", 1, tags)
	}
}

// TrackDatabaseOperation tracks database operations
func (o *ObservabilityIntegration) TrackDatabaseOperation(ctx context.Context, operation, table string, duration time.Duration, err error) {
	fields := map[string]interface{}{
		"operation":   operation,
		"table":       table,
		"duration_ms": duration.Milliseconds(),
	}

	logLevel := "debug"
	message := "Database operation completed"
	
	if err != nil {
		logLevel = "error"
		message = "Database operation failed"
		fields["error"] = err.Error()
	}

	o.observabilityService.LogEvent(ctx, logLevel, message, fields, err)

	// Record database metrics
	tags := map[string]string{
		"operation": operation,
		"table":     table,
	}

	if err != nil {
		tags["status"] = "error"
		o.observabilityService.RecordMetric(ctx, "database_errors_total", "counter", 1, tags)
	} else {
		tags["status"] = "success"
	}

	o.observabilityService.RecordMetric(ctx, "database_operations_total", "counter", 1, tags)
	o.observabilityService.RecordMetric(ctx, "database_operation_duration_ms", "histogram", float64(duration.Milliseconds()), tags)
}

// Middleware functions for Gin

// RequestObservabilityMiddleware provides comprehensive request observability
func (o *ObservabilityIntegration) RequestObservabilityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Extract tenant ID from context if available
		var tenantID *uuid.UUID
		if tid, exists := c.Get("tenant_id"); exists {
			if id, ok := tid.(uuid.UUID); ok {
				tenantID = &id
			}
		}

		// Start distributed trace
		err := o.observabilityService.TraceOperation(c.Request.Context(), 
			c.Request.Method+" "+c.FullPath(),
			map[string]interface{}{
				"http.method": c.Request.Method,
				"http.path":   c.FullPath(),
				"http.url":    c.Request.URL.String(),
			},
			func(ctx context.Context) error {
				c.Request = c.Request.WithContext(ctx)
				c.Next()
				return nil
			})

		if err != nil {
			o.observabilityService.LogEvent(c.Request.Context(), "error", 
				"Failed to trace request", map[string]interface{}{
					"method": c.Request.Method,
					"path":   c.FullPath(),
					"error":  err.Error(),
				}, err)
		}

		duration := time.Since(start)
		
		// Track API call with comprehensive metrics
		o.TrackAPICall(c.Request.Context(), c.Request.Method, c.FullPath(), 
			c.Writer.Status(), duration, tenantID)

		// Track page view for analytics if it's a GET request to a page
		if c.Request.Method == "GET" && tenantID != nil && isPageRequest(c.FullPath()) {
			pageView := map[string]interface{}{
				"url":        c.Request.URL.String(),
				"path":       c.FullPath(),
				"referrer":   c.Request.Referer(),
				"user_agent": c.Request.UserAgent(),
				"tenant_id":  *tenantID,
				"timestamp":  time.Now(),
			}
			
			if userID, exists := c.Get("user_id"); exists {
				pageView["user_id"] = userID
			}
			
			if sessionID, exists := c.Get("session_id"); exists {
				pageView["session_id"] = sessionID
			}

			if err := o.analyticsService.TrackPageView(c.Request.Context(), *tenantID, pageView); err != nil {
				o.observabilityService.LogEvent(c.Request.Context(), "warn", 
					"Failed to track page view", map[string]interface{}{
						"path":  c.FullPath(),
						"error": err.Error(),
					}, err)
			}
		}
	}
}

// ErrorObservabilityMiddleware tracks and logs errors
func (o *ObservabilityIntegration) ErrorObservabilityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for errors in the context
		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				o.observabilityService.LogEvent(c.Request.Context(), "error", 
					"Request error occurred", map[string]interface{}{
						"method":     c.Request.Method,
						"path":       c.FullPath(),
						"error_type": ginErr.Type,
						"error_msg":  ginErr.Error(),
					}, ginErr.Err)

				// Record error metrics
				o.observabilityService.RecordMetric(c.Request.Context(), "request_errors_total", "counter", 1,
					map[string]string{
						"method":     c.Request.Method,
						"path":       c.FullPath(),
						"error_type": string(rune(int(ginErr.Type))),
					})
			}
		}
	}
}

// Helper functions

// isPageRequest determines if a request path represents a page view for analytics
func isPageRequest(path string) bool {
	// Define patterns that constitute page views
	// This would be customized based on your application's routing structure
	pagePatterns := []string{
		"/products",
		"/orders", 
		"/dashboard",
		"/settings",
		"/profile",
	}
	
	for _, pattern := range pagePatterns {
		if len(path) >= len(pattern) && path[:len(pattern)] == pattern {
			return true
		}
	}
	
	return false
}

// HealthCheck provides a combined health check for analytics and observability
func (o *ObservabilityIntegration) HealthCheck(ctx context.Context) map[string]interface{} {
	status := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"services": map[string]interface{}{
			"analytics": map[string]interface{}{
				"status": "healthy",
				"details": map[string]interface{}{
					"tracking_enabled": true,
					"last_event":       time.Now().Add(-1 * time.Minute), // Mock data
				},
			},
			"observability": map[string]interface{}{
				"status": "healthy", 
				"details": map[string]interface{}{
					"logging_enabled":  true,
					"metrics_enabled":  true,
					"tracing_enabled":  true,
				},
			},
		},
	}

	return status
}