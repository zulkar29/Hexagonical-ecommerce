package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ObservabilityService handles observability operations
type ObservabilityService struct {
	logger  Logger
	metrics MetricsCollector
	tracer  Tracer
	repo    ObservabilityRepository
}

// NewObservabilityService creates a new observability service
func NewObservabilityService(logger Logger, metrics MetricsCollector, tracer Tracer, repo ObservabilityRepository) *ObservabilityService {
	return &ObservabilityService{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
		repo:    repo,
	}
}

// LogEvent logs an event with context
func (s *ObservabilityService) LogEvent(ctx context.Context, level LogLevel, message string, fields map[string]interface{}, err error) {
	switch level {
	case LogLevelDebug:
		s.logger.Debug(ctx, message, fields)
	case LogLevelInfo:
		s.logger.Info(ctx, message, fields)
	case LogLevelWarn:
		s.logger.Warn(ctx, message, fields)
	case LogLevelError:
		s.logger.Error(ctx, message, err, fields)
	case LogLevelFatal:
		s.logger.Fatal(ctx, message, err, fields)
	}
}

// RecordMetric records a metric
func (s *ObservabilityService) RecordMetric(ctx context.Context, name string, metricType MetricType, value float64, tags map[string]string) {
	switch metricType {
	case MetricTypeCounter:
		s.metrics.Counter(name, value, tags)
	case MetricTypeGauge:
		s.metrics.Gauge(name, value, tags)
	case MetricTypeHistogram:
		s.metrics.Histogram(name, value, tags)
	case MetricTypeTimer:
		s.metrics.Timer(name, time.Duration(value)*time.Millisecond, tags)
	case MetricTypeSet:
		s.metrics.Set(name, fmt.Sprintf("%.0f", value), tags)
	}
}

// TraceOperation traces an operation
func (s *ObservabilityService) TraceOperation(ctx context.Context, operationName string, tags map[string]interface{}, fn func(ctx context.Context) error) error {
	return s.tracer.TraceOperation(ctx, operationName, tags, fn)
}

// GetHealthStatus returns the health status of the system
func (s *ObservabilityService) GetHealthStatus(ctx context.Context) (*HealthStatus, error) {
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Services:  make(map[string]ServiceHealth),
		Version:   "1.0.0", // This should come from config
	}

	// Check database health
	dbHealth := s.checkDatabaseHealth(ctx)
	status.Services["database"] = dbHealth

	// Check cache health (if implemented)
	cacheHealth := s.checkCacheHealth(ctx)
	status.Services["cache"] = cacheHealth

	// Determine overall status
	for _, service := range status.Services {
		if service.Status != "healthy" {
			status.Status = "unhealthy"
			break
		}
	}

	return status, nil
}

// GetMetrics returns collected metrics
func (s *ObservabilityService) GetMetrics(ctx context.Context, filters map[string]string) (map[string]*Metric, error) {
	metrics := s.metrics.GetMetrics()
	
	// Apply filters if provided
	if len(filters) > 0 {
		filtered := make(map[string]*Metric)
		for key, metric := range metrics {
			include := true
			for filterKey, filterValue := range filters {
				if metricValue, exists := metric.Tags[filterKey]; !exists || metricValue != filterValue {
					include = false
					break
				}
			}
			if include {
				filtered[key] = metric
			}
		}
		return filtered, nil
	}

	return metrics, nil
}

// GetTraces returns traces
func (s *ObservabilityService) GetTraces(ctx context.Context, limit int) (map[string]*Trace, error) {
	traces := s.tracer.GetTraces()
	
	// If limit is specified, return only the most recent traces
	if limit > 0 && len(traces) > limit {
		// This is a simple implementation - in production you'd want proper pagination
		limited := make(map[string]*Trace)
		count := 0
		for k, v := range traces {
			if count >= limit {
				break
			}
			limited[k] = v
			count++
		}
		return limited, nil
	}

	return traces, nil
}

// GetAlerts returns current alerts
func (s *ObservabilityService) GetAlerts(ctx context.Context) ([]*Alert, error) {
	// This would typically come from an alert store
	// For now, return empty array
	return []*Alert{}, nil
}

// CreateAlert creates a new alert
func (s *ObservabilityService) CreateAlert(ctx context.Context, alert *Alert) error {
	alert.ID = time.Now().Unix() // Simple ID generation
	alert.CreatedAt = time.Now().UTC()
	alert.UpdatedAt = time.Now().UTC()

	// Log the alert
	s.logger.Warn(ctx, fmt.Sprintf("Alert created: %s", alert.Title), map[string]interface{}{
		"alert_id":   alert.ID,
		"severity":   alert.Severity,
		"service":    alert.Service,
		"condition":  alert.Condition,
	})

	// Record alert metric
	s.metrics.Increment("alerts.created", map[string]string{
		"severity": string(alert.Severity),
		"service":  alert.Service,
	})

	return nil
}

// Helper methods for health checks

func (s *ObservabilityService) checkDatabaseHealth(ctx context.Context) ServiceHealth {
	health := ServiceHealth{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Details:   make(map[string]interface{}),
	}

	// This would check actual database connectivity
	// For now, return healthy
	health.Details["connected"] = true
	health.Details["response_time"] = "5ms"

	return health
}

func (s *ObservabilityService) checkCacheHealth(ctx context.Context) ServiceHealth {
	health := ServiceHealth{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Details:   make(map[string]interface{}),
	}

	// This would check actual cache connectivity
	// For now, return healthy
	health.Details["connected"] = true
	health.Details["hit_rate"] = "95%"

	return health
}

// Business-specific observability methods

// TrackUserAction tracks user actions for analytics
func (s *ObservabilityService) TrackUserAction(ctx context.Context, userID, tenantID, action string, properties map[string]interface{}) {
	// Log the action
	s.logger.Info(ctx, fmt.Sprintf("User action: %s", action), map[string]interface{}{
		"user_id":    userID,
		"tenant_id":  tenantID,
		"action":     action,
		"properties": properties,
	})

	// Record metrics
	s.metrics.Increment("user.actions", map[string]string{
		"action":     action,
		"tenant_id":  tenantID,
	})
}

// TrackBusinessEvent tracks business events
func (s *ObservabilityService) TrackBusinessEvent(ctx context.Context, event string, tenantID string, value float64, properties map[string]interface{}) {
	// Log the event
	s.logger.Info(ctx, fmt.Sprintf("Business event: %s", event), map[string]interface{}{
		"tenant_id":  tenantID,
		"event":      event,
		"value":      value,
		"properties": properties,
	})

	// Record business metric
	s.metrics.RecordBusinessMetric(BusinessMetric{
		Name:     event,
		Category: "business",
		Type:     "event",
		Value:    value,
		Tags: map[string]string{
			"tenant_id": tenantID,
		},
	})
}

// TrackAPIPerformance tracks API endpoint performance
func (s *ObservabilityService) TrackAPIPerformance(ctx context.Context, endpoint, method string, statusCode int, duration time.Duration) {
	status := fmt.Sprintf("%d", statusCode)
	
	// Record performance metric
	s.metrics.RecordAPICall(endpoint, method, status, duration)

	// Log slow requests
	if duration > 1*time.Second {
		s.logger.Warn(ctx, "Slow API request", map[string]interface{}{
			"endpoint":     endpoint,
			"method":       method,
			"status_code":  statusCode,
			"duration_ms":  duration.Milliseconds(),
		})
	}
}

// Middleware functions

// LoggingMiddleware adds request logging
func (s *ObservabilityService) LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format that integrates with our structured logging
		fields := map[string]interface{}{
			"method":      param.Method,
			"path":        param.Path,
			"status_code": param.StatusCode,
			"latency":     param.Latency.String(),
			"client_ip":   param.ClientIP,
			"user_agent":  param.Request.UserAgent(),
		}

		level := LogLevelInfo
		if param.StatusCode >= 400 {
			level = LogLevelWarn
		}
		if param.StatusCode >= 500 {
			level = LogLevelError
		}

		s.LogEvent(param.Request.Context(), level, "HTTP Request", fields, nil)
		return "" // We handle logging ourselves
	})
}

// MetricsMiddleware adds request metrics
func (s *ObservabilityService) MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		
		// Track API performance
		s.TrackAPIPerformance(
			c.Request.Context(),
			c.FullPath(),
			c.Request.Method,
			c.Writer.Status(),
			duration,
		)
	}
}

// TracingMiddleware adds distributed tracing
func (s *ObservabilityService) TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start a new trace for the request
		ctx, trace := s.tracer.StartTrace(c.Request.Context(), 
			fmt.Sprintf("%s %s", c.Request.Method, c.FullPath()),
			map[string]interface{}{
				"http.method": c.Request.Method,
				"http.path":   c.FullPath(),
				"http.url":    c.Request.URL.String(),
			})

		// Add trace context to gin context
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// Finish the trace
		s.tracer.FinishTrace(trace)
	}
}
