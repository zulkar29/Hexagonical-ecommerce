package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// MetricType represents the type of metric
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
	MetricTypeTimer     MetricType = "timer"
	MetricTypeSet       MetricType = "set"
)

// AlertSeverity represents the severity of an alert
type AlertSeverity string

const (
	AlertSeverityLow      AlertSeverity = "low"
	AlertSeverityMedium   AlertSeverity = "medium"
	AlertSeverityHigh     AlertSeverity = "high"
	AlertSeverityCritical AlertSeverity = "critical"
)

// AlertStatus represents the status of an alert
type AlertStatus string

const (
	AlertStatusActive   AlertStatus = "active"
	AlertStatusResolved AlertStatus = "resolved"
	AlertStatusMuted    AlertStatus = "muted"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	ID        uuid.UUID              `json:"id"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	TenantID  string                 `json:"tenant_id,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Error     *ErrorDetails          `json:"error,omitempty"`
}

// ErrorDetails represents error information in logs
type ErrorDetails struct {
	Type       string                 `json:"type"`
	Message    string                 `json:"message"`
	StackTrace string                 `json:"stack_trace,omitempty"`
	Code       string                 `json:"code,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// Metric represents a metric entry
type Metric struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Type      MetricType             `json:"type"`
	Value     float64                `json:"value"`
	Tags      map[string]string      `json:"tags,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
}

// BusinessMetric represents business-specific metrics
type BusinessMetric struct {
	Name     string            `json:"name"`
	Category string            `json:"category"` // user, sales, inventory, etc.
	Type     string            `json:"type"`     // activity, revenue, count, etc.
	Value    float64           `json:"value"`
	Tags     map[string]string `json:"tags,omitempty"`
}

// PerformanceMetric represents performance metrics
type PerformanceMetric struct {
	Operation    string        `json:"operation"`
	Resource     string        `json:"resource,omitempty"`
	Status       string        `json:"status"`
	ResponseTime time.Duration `json:"response_time"`
	Throughput   float64       `json:"throughput,omitempty"`
	ErrorRate    float64       `json:"error_rate,omitempty"`
}

// Trace represents a distributed trace
type Trace struct {
	ID            string                 `json:"id"`
	OperationName string                 `json:"operation_name"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	Duration      time.Duration          `json:"duration"`
	Service       string                 `json:"service"`
	Version       string                 `json:"version"`
	Tags          map[string]interface{} `json:"tags,omitempty"`
	Spans         []*Span                `json:"spans,omitempty"`
	SpanCount     int                    `json:"span_count"`
}

// Span represents a span within a trace
type Span struct {
	ID            string                 `json:"id"`
	TraceID       string                 `json:"trace_id"`
	ParentSpanID  string                 `json:"parent_span_id,omitempty"`
	OperationName string                 `json:"operation_name"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	Duration      time.Duration          `json:"duration"`
	Service       string                 `json:"service"`
	Tags          map[string]interface{} `json:"tags,omitempty"`
	Logs          []SpanLog              `json:"logs,omitempty"`
}

// SpanLog represents a log entry within a span
type SpanLog struct {
	Timestamp time.Time              `json:"timestamp"`
	Fields    map[string]interface{} `json:"fields"`
}

// Alert represents a system alert
type Alert struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Severity    AlertSeverity `json:"severity"`
	Service     string        `json:"service"`
	Condition   string        `json:"condition"`
	Status      AlertStatus   `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	ResolvedAt  time.Time     `json:"resolved_at,omitempty"`
}

// HealthStatus represents the health status of the system
type HealthStatus struct {
	Status    string                    `json:"status"`
	Timestamp time.Time                 `json:"timestamp"`
	Version   string                    `json:"version"`
	Services  map[string]ServiceHealth  `json:"services"`
}

// ServiceHealth represents the health of a specific service
type ServiceHealth struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// Logger interface for structured logging
type Logger interface {
	Debug(ctx context.Context, message string, fields ...map[string]interface{})
	Info(ctx context.Context, message string, fields ...map[string]interface{})
	Warn(ctx context.Context, message string, fields ...map[string]interface{})
	Error(ctx context.Context, message string, err error, fields ...map[string]interface{})
	Fatal(ctx context.Context, message string, err error, fields ...map[string]interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithRequestID(requestID string) Logger
	WithUserID(userID string) Logger
	WithTenantID(tenantID string) Logger
	WithTrace(traceID, spanID string) Logger
}

// MetricsCollector interface for collecting metrics
type MetricsCollector interface {
	Counter(name string, value float64, tags map[string]string)
	Gauge(name string, value float64, tags map[string]string)
	Histogram(name string, value float64, tags map[string]string)
	Timer(name string, duration time.Duration, tags map[string]string)
	Set(name string, value string, tags map[string]string)
	Increment(name string, tags map[string]string)
	Decrement(name string, tags map[string]string)
	Timing(name string, tags map[string]string, fn func())
	RecordBusinessMetric(metric BusinessMetric)
	RecordPerformanceMetric(metric PerformanceMetric)
	GetMetrics() map[string]*Metric
	Reset()
}

// Tracer interface for distributed tracing
type Tracer interface {
	StartTrace(ctx context.Context, operationName string, tags map[string]interface{}) (context.Context, *Trace)
	StartSpan(ctx context.Context, operationName string, tags map[string]interface{}) (context.Context, *Span)
	FinishSpan(span *Span)
	FinishTrace(trace *Trace)
	AddSpanTag(span *Span, key string, value interface{})
	AddSpanLog(span *Span, fields map[string]interface{})
	SetSpanError(span *Span, err error)
	GetTrace(traceID string) *Trace
	GetTraces() map[string]*Trace
	TraceOperation(ctx context.Context, operationName string, tags map[string]interface{}, fn func(ctx context.Context) error) error
}

// ObservabilityRepository interface for persistence
type ObservabilityRepository interface {
	// Log methods
	SaveLogEntry(ctx context.Context, entry *LogEntry) error
	GetLogEntries(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*LogEntry, error)
	
	// Metric methods
	SaveMetric(ctx context.Context, metric *Metric) error
	GetMetrics(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*Metric, error)
	
	// Trace methods
	SaveTrace(ctx context.Context, trace *Trace) error
	SaveSpan(ctx context.Context, span *Span) error
	GetTrace(ctx context.Context, traceID string) (*Trace, error)
	GetSpansByTraceID(ctx context.Context, traceID string) ([]*Span, error)
	
	// Alert methods
	SaveAlert(ctx context.Context, alert *Alert) error
	GetAlerts(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*Alert, error)
	
	// Migration
	Migrate() error
}

// logger implements the Logger interface
type logger struct {
	service   string
	version   string
	level     LogLevel
	fields    map[string]interface{}
	requestID string
	userID    string
	tenantID  string
	traceID   string
	spanID    string
}

// NewLogger creates a new structured logger
func NewLogger(service, version string, level LogLevel) Logger {
	return &logger{
		service: service,
		version: version,
		level:   level,
		fields:  make(map[string]interface{}),
	}
}

// Debug logs debug level messages
func (l *logger) Debug(ctx context.Context, message string, fields ...map[string]interface{}) {
	if l.shouldLog(LogLevelDebug) {
		l.log(LogLevelDebug, message, nil, fields...)
	}
}

// Info logs info level messages
func (l *logger) Info(ctx context.Context, message string, fields ...map[string]interface{}) {
	if l.shouldLog(LogLevelInfo) {
		l.log(LogLevelInfo, message, nil, fields...)
	}
}

// Warn logs warning level messages
func (l *logger) Warn(ctx context.Context, message string, fields ...map[string]interface{}) {
	if l.shouldLog(LogLevelWarn) {
		l.log(LogLevelWarn, message, nil, fields...)
	}
}

// Error logs error level messages
func (l *logger) Error(ctx context.Context, message string, err error, fields ...map[string]interface{}) {
	if l.shouldLog(LogLevelError) {
		l.log(LogLevelError, message, err, fields...)
	}
}

// Fatal logs fatal level messages and exits
func (l *logger) Fatal(ctx context.Context, message string, err error, fields ...map[string]interface{}) {
	l.log(LogLevelFatal, message, err, fields...)
	os.Exit(1)
}

// WithField adds a field to the logger context
func (l *logger) WithField(key string, value interface{}) Logger {
	newLogger := l.copy()
	newLogger.fields[key] = value
	return newLogger
}

// WithFields adds multiple fields to the logger context
func (l *logger) WithFields(fields map[string]interface{}) Logger {
	newLogger := l.copy()
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

// WithRequestID adds request ID to the logger context
func (l *logger) WithRequestID(requestID string) Logger {
	newLogger := l.copy()
	newLogger.requestID = requestID
	return newLogger
}

// WithUserID adds user ID to the logger context
func (l *logger) WithUserID(userID string) Logger {
	newLogger := l.copy()
	newLogger.userID = userID
	return newLogger
}

// WithTenantID adds tenant ID to the logger context
func (l *logger) WithTenantID(tenantID string) Logger {
	newLogger := l.copy()
	newLogger.tenantID = tenantID
	return newLogger
}

// WithTrace adds trace information to the logger context
func (l *logger) WithTrace(traceID, spanID string) Logger {
	newLogger := l.copy()
	newLogger.traceID = traceID
	newLogger.spanID = spanID
	return newLogger
}

// copy creates a copy of the logger with the same context
func (l *logger) copy() *logger {
	fields := make(map[string]interface{})
	for k, v := range l.fields {
		fields[k] = v
	}
	
	return &logger{
		service:   l.service,
		version:   l.version,
		level:     l.level,
		fields:    fields,
		requestID: l.requestID,
		userID:    l.userID,
		tenantID:  l.tenantID,
		traceID:   l.traceID,
		spanID:    l.spanID,
	}
}

// shouldLog determines if the message should be logged based on level
func (l *logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		LogLevelDebug: 0,
		LogLevelInfo:  1,
		LogLevelWarn:  2,
		LogLevelError: 3,
		LogLevelFatal: 4,
	}
	
	return levels[level] >= levels[l.level]
}

// log creates and outputs a structured log entry
func (l *logger) log(level LogLevel, message string, err error, extraFields ...map[string]interface{}) {
	entry := LogEntry{
		ID:        uuid.New(),
		Level:     level,
		Message:   message,
		Timestamp: time.Now().UTC(),
		Service:   l.service,
		Version:   l.version,
		RequestID: l.requestID,
		UserID:    l.userID,
		TenantID:  l.tenantID,
		TraceID:   l.traceID,
		SpanID:    l.spanID,
		Fields:    make(map[string]interface{}),
	}

	// Add base fields
	for k, v := range l.fields {
		entry.Fields[k] = v
	}

	// Add extra fields
	for _, fields := range extraFields {
		for k, v := range fields {
			entry.Fields[k] = v
		}
	}

	// Add error details if present
	if err != nil {
		entry.Error = &ErrorDetails{
			Type:       fmt.Sprintf("%T", err),
			Message:    err.Error(),
			StackTrace: getStackTrace(),
		}
	}

	// Output the log entry
	l.output(entry)
}

// output writes the log entry to the configured output
func (l *logger) output(entry LogEntry) {
	// Convert to JSON for structured logging
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple logging if JSON marshaling fails
		log.Printf("[%s] %s: %s", entry.Level, entry.Service, entry.Message)
		return
	}

	// For development, format for readability
	if l.isDevelopment() {
		l.outputDevelopment(entry)
	} else {
		// Production: output as JSON
		fmt.Println(string(jsonData))
	}
}

// outputDevelopment formats logs for development readability
func (l *logger) outputDevelopment(entry LogEntry) {
	// Color codes for different log levels
	colors := map[LogLevel]string{
		LogLevelDebug: "\033[36m", // Cyan
		LogLevelInfo:  "\033[32m", // Green
		LogLevelWarn:  "\033[33m", // Yellow
		LogLevelError: "\033[31m", // Red
		LogLevelFatal: "\033[35m", // Magenta
	}
	reset := "\033[0m"

	color := colors[entry.Level]
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")
	
	// Basic log line
	logLine := fmt.Sprintf("%s[%s]%s %s %s[%s]%s %s",
		color, strings.ToUpper(string(entry.Level)), reset,
		timestamp,
		color, entry.Service, reset,
		entry.Message,
	)

	// Add context information
	var context []string
	if entry.RequestID != "" {
		context = append(context, fmt.Sprintf("req_id=%s", entry.RequestID))
	}
	if entry.UserID != "" {
		context = append(context, fmt.Sprintf("user_id=%s", entry.UserID))
	}
	if entry.TenantID != "" {
		context = append(context, fmt.Sprintf("tenant_id=%s", entry.TenantID))
	}
	if entry.TraceID != "" {
		context = append(context, fmt.Sprintf("trace_id=%s", entry.TraceID[:8])) // Shortened
	}

	if len(context) > 0 {
		logLine += fmt.Sprintf(" [%s]", strings.Join(context, ", "))
	}

	// Add fields
	if len(entry.Fields) > 0 {
		fieldsJSON, _ := json.Marshal(entry.Fields)
		logLine += fmt.Sprintf(" fields=%s", string(fieldsJSON))
	}

	// Add error information
	if entry.Error != nil {
		logLine += fmt.Sprintf(" error=%s", entry.Error.Message)
	}

	fmt.Println(logLine)

	// Print stack trace for errors in development
	if entry.Level == LogLevelError || entry.Level == LogLevelFatal {
		if entry.Error != nil && entry.Error.StackTrace != "" {
			fmt.Printf("Stack trace:\n%s\n", entry.Error.StackTrace)
		}
	}
}

// isDevelopment checks if running in development mode
func (l *logger) isDevelopment() bool {
	env := os.Getenv("ENVIRONMENT")
	return env == "development" || env == ""
}

// getStackTrace captures the current stack trace
func getStackTrace() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}

// metricsCollector implements the MetricsCollector interface
type metricsCollector struct {
	mu      sync.RWMutex
	metrics map[string]*Metric
	config  MetricsConfig
}

// MetricsConfig holds configuration for metrics collection
type MetricsConfig struct {
	Service       string
	Version       string
	Environment   string
	BatchSize     int
	FlushInterval time.Duration
	Tags          map[string]string
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(config MetricsConfig) MetricsCollector {
	if config.BatchSize == 0 {
		config.BatchSize = 100
	}
	if config.FlushInterval == 0 {
		config.FlushInterval = 30 * time.Second
	}
	if config.Tags == nil {
		config.Tags = make(map[string]string)
	}

	return &metricsCollector{
		metrics: make(map[string]*Metric),
		config:  config,
	}
}

// Counter increments a counter metric
func (m *metricsCollector) Counter(name string, value float64, tags map[string]string) {
	m.recordMetric(name, MetricTypeCounter, value, tags)
}

// Gauge sets a gauge metric value
func (m *metricsCollector) Gauge(name string, value float64, tags map[string]string) {
	m.recordMetric(name, MetricTypeGauge, value, tags)
}

// Histogram records a histogram metric
func (m *metricsCollector) Histogram(name string, value float64, tags map[string]string) {
	m.recordMetric(name, MetricTypeHistogram, value, tags)
}

// Timer records execution time
func (m *metricsCollector) Timer(name string, duration time.Duration, tags map[string]string) {
	m.recordMetric(name, MetricTypeTimer, float64(duration.Milliseconds()), tags)
}

// Set records a set metric (unique values)
func (m *metricsCollector) Set(name string, value string, tags map[string]string) {
	m.recordMetric(name, MetricTypeSet, 1, mergeTags(tags, map[string]string{"value": value}))
}

// Increment increments a counter by 1
func (m *metricsCollector) Increment(name string, tags map[string]string) {
	m.Counter(name, 1, tags)
}

// Decrement decrements a counter by 1
func (m *metricsCollector) Decrement(name string, tags map[string]string) {
	m.Counter(name, -1, tags)
}

// Timing records a timing metric with function execution
func (m *metricsCollector) Timing(name string, tags map[string]string, fn func()) {
	start := time.Now()
	fn()
	duration := time.Since(start)
	m.Timer(name, duration, tags)
}

// RecordBusinessMetric records a business-specific metric
func (m *metricsCollector) RecordBusinessMetric(metric BusinessMetric) {
	tags := map[string]string{
		"category": metric.Category,
		"type":     metric.Type,
	}
	
	for k, v := range metric.Tags {
		tags[k] = v
	}

	m.recordMetric(metric.Name, MetricTypeGauge, metric.Value, tags)
}

// RecordPerformanceMetric records a performance metric
func (m *metricsCollector) RecordPerformanceMetric(metric PerformanceMetric) {
	tags := map[string]string{
		"operation": metric.Operation,
		"status":    metric.Status,
	}
	
	if metric.Resource != "" {
		tags["resource"] = metric.Resource
	}

	if metric.ResponseTime > 0 {
		m.Timer(fmt.Sprintf("%s.response_time", metric.Operation), metric.ResponseTime, tags)
	}

	if metric.Throughput > 0 {
		m.Gauge(fmt.Sprintf("%s.throughput", metric.Operation), metric.Throughput, tags)
	}

	if metric.ErrorRate >= 0 {
		m.Gauge(fmt.Sprintf("%s.error_rate", metric.Operation), metric.ErrorRate, tags)
	}
}

// GetMetrics returns all collected metrics
func (m *metricsCollector) GetMetrics() map[string]*Metric {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*Metric)
	for k, v := range m.metrics {
		result[k] = v
	}
	return result
}

// Reset clears all metrics
func (m *metricsCollector) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.metrics = make(map[string]*Metric)
}

// recordMetric records a metric with the given parameters
func (m *metricsCollector) recordMetric(name string, metricType MetricType, value float64, tags map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	allTags := mergeTags(m.config.Tags, tags)
	key := m.createMetricKey(name, allTags)

	metric := &Metric{
		ID:        uuid.New(),
		Name:      name,
		Type:      metricType,
		Value:     value,
		Tags:      allTags,
		Timestamp: time.Now().UTC(),
		Service:   m.config.Service,
		Version:   m.config.Version,
	}

	m.metrics[key] = metric
}

// createMetricKey creates a unique key for a metric based on name and tags
func (m *metricsCollector) createMetricKey(name string, tags map[string]string) string {
	key := name
	for k, v := range tags {
		key += fmt.Sprintf(":%s=%s", k, v)
	}
	return key
}

// mergeTags merges multiple tag maps
func mergeTags(tagMaps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, tags := range tagMaps {
		for k, v := range tags {
			result[k] = v
		}
	}
	return result
}

// tracer implements the Tracer interface
type tracer struct {
	mu      sync.RWMutex
	traces  map[string]*Trace
	config  TracingConfig
}

// TracingConfig holds configuration for distributed tracing
type TracingConfig struct {
	Service     string
	Version     string
	Environment string
	SampleRate  float64
	MaxSpans    int
}

// NewTracer creates a new distributed tracer
func NewTracer(config TracingConfig) Tracer {
	if config.SampleRate == 0 {
		config.SampleRate = 1.0
	}
	if config.MaxSpans == 0 {
		config.MaxSpans = 1000
	}

	return &tracer{
		traces: make(map[string]*Trace),
		config: config,
	}
}

// StartTrace starts a new trace
func (t *tracer) StartTrace(ctx context.Context, operationName string, tags map[string]interface{}) (context.Context, *Trace) {
	traceID := uuid.New().String()
	
	trace := &Trace{
		ID:            traceID,
		OperationName: operationName,
		StartTime:     time.Now().UTC(),
		Service:       t.config.Service,
		Version:       t.config.Version,
		Tags:          tags,
		Spans:         make([]*Span, 0),
	}

	t.mu.Lock()
	t.traces[traceID] = trace
	t.mu.Unlock()

	ctx = context.WithValue(ctx, "trace_id", traceID)
	ctx = context.WithValue(ctx, "trace", trace)

	return ctx, trace
}

// StartSpan starts a new span within a trace
func (t *tracer) StartSpan(ctx context.Context, operationName string, tags map[string]interface{}) (context.Context, *Span) {
	var trace *Trace
	var parentSpanID string

	if traceValue := ctx.Value("trace"); traceValue != nil {
		trace = traceValue.(*Trace)
	}

	if parentValue := ctx.Value("span_id"); parentValue != nil {
		parentSpanID = parentValue.(string)
	}

	spanID := uuid.New().String()
	span := &Span{
		ID:            spanID,
		TraceID:       func() string {
			if trace != nil {
				return trace.ID
			}
			return ""
		}(),
		ParentSpanID:  parentSpanID,
		OperationName: operationName,
		StartTime:     time.Now().UTC(),
		Service:       t.config.Service,
		Tags:          tags,
		Logs:          make([]SpanLog, 0),
	}

	if trace != nil {
		t.mu.Lock()
		trace.Spans = append(trace.Spans, span)
		t.mu.Unlock()
	}

	ctx = context.WithValue(ctx, "span_id", spanID)
	ctx = context.WithValue(ctx, "span", span)

	return ctx, span
}

// FinishSpan finishes a span
func (t *tracer) FinishSpan(span *Span) {
	if span == nil {
		return
	}

	span.EndTime = time.Now().UTC()
	span.Duration = span.EndTime.Sub(span.StartTime)
}

// FinishTrace finishes a trace
func (t *tracer) FinishTrace(trace *Trace) {
	if trace == nil {
		return
	}

	trace.EndTime = time.Now().UTC()
	trace.Duration = trace.EndTime.Sub(trace.StartTime)
	trace.SpanCount = len(trace.Spans)
}

// AddSpanTag adds a tag to a span
func (t *tracer) AddSpanTag(span *Span, key string, value interface{}) {
	if span == nil {
		return
	}

	if span.Tags == nil {
		span.Tags = make(map[string]interface{})
	}
	span.Tags[key] = value
}

// AddSpanLog adds a log entry to a span
func (t *tracer) AddSpanLog(span *Span, fields map[string]interface{}) {
	if span == nil {
		return
	}

	logEntry := SpanLog{
		Timestamp: time.Now().UTC(),
		Fields:    fields,
	}

	span.Logs = append(span.Logs, logEntry)
}

// SetSpanError marks a span as having an error
func (t *tracer) SetSpanError(span *Span, err error) {
	if span == nil || err == nil {
		return
	}

	t.AddSpanTag(span, "error", true)
	t.AddSpanTag(span, "error.message", err.Error())
	t.AddSpanTag(span, "error.type", fmt.Sprintf("%T", err))

	t.AddSpanLog(span, map[string]interface{}{
		"event":   "error",
		"message": err.Error(),
	})
}

// GetTrace retrieves a trace by ID
func (t *tracer) GetTrace(traceID string) *Trace {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	return t.traces[traceID]
}

// GetTraces returns all traces
func (t *tracer) GetTraces() map[string]*Trace {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make(map[string]*Trace)
	for k, v := range t.traces {
		result[k] = v
	}
	return result
}

// TraceOperation is a helper function to trace an operation with automatic error handling
func (t *tracer) TraceOperation(ctx context.Context, operationName string, tags map[string]interface{}, fn func(ctx context.Context) error) error {
	ctx, span := t.StartSpan(ctx, operationName, tags)
	defer t.FinishSpan(span)

	err := fn(ctx)
	if err != nil {
		t.SetSpanError(span, err)
	}

	return err
}

// Global instances
var (
	globalLogger  Logger
	globalMetrics MetricsCollector
	globalTracer  Tracer
)

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(service, version string, level LogLevel) {
	globalLogger = NewLogger(service, version, level)
}

// InitGlobalMetrics initializes the global metrics collector
func InitGlobalMetrics(config MetricsConfig) {
	globalMetrics = NewMetricsCollector(config)
}

// InitGlobalTracer initializes the global tracer
func InitGlobalTracer(config TracingConfig) {
	globalTracer = NewTracer(config)
}

// GetGlobalLogger returns the global logger
func GetGlobalLogger() Logger {
	return globalLogger
}

// GetGlobalMetrics returns the global metrics collector
func GetGlobalMetrics() MetricsCollector {
	return globalMetrics
}

// GetGlobalTracer returns the global tracer
func GetGlobalTracer() Tracer {
	return globalTracer
}

// Convenience functions for global logger
func Debug(ctx context.Context, message string, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(ctx, message, fields...)
	}
}

func Info(ctx context.Context, message string, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Info(ctx, message, fields...)
	}
}

func Warn(ctx context.Context, message string, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(ctx, message, fields...)
	}
}

func Error(ctx context.Context, message string, err error, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Error(ctx, message, err, fields...)
	}
}

func Fatal(ctx context.Context, message string, err error, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Fatal(ctx, message, err, fields...)
	}
}
