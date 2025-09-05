package observability

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// observabilityRepository handles persistence of observability data
type observabilityRepository struct {
	db *gorm.DB
}

// NewObservabilityRepository creates a new observability repository
func NewObservabilityRepository(db *gorm.DB) ObservabilityRepository {
	return &observabilityRepository{
		db: db,
	}
}

// Database models for persistence

// LogEntryModel represents a log entry in the database
type LogEntryModel struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	Level     string    `gorm:"not null;index"`
	Message   string    `gorm:"not null"`
	Service   string    `gorm:"not null;index"`
	Version   string    `gorm:"not null"`
	RequestID string    `gorm:"index"`
	UserID    string    `gorm:"index"`
	TenantID  string    `gorm:"index"`
	TraceID   string    `gorm:"index"`
	SpanID    string    `gorm:"index"`
	Fields    string    `gorm:"type:jsonb"` // PostgreSQL JSONB
	Error     string    `gorm:"type:jsonb"` // Error details as JSON
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}

// MetricModel represents a metric in the database
type MetricModel struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	Name      string    `gorm:"not null;index"`
	Type      string    `gorm:"not null"`
	Value     float64   `gorm:"not null"`
	Service   string    `gorm:"not null;index"`
	Version   string    `gorm:"not null"`
	Tags      string    `gorm:"type:jsonb"` // Tags as JSON
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}

// TraceModel represents a trace in the database
type TraceModel struct {
	ID            string    `gorm:"primaryKey;type:uuid"`
	OperationName string    `gorm:"not null"`
	Service       string    `gorm:"not null;index"`
	Version       string    `gorm:"not null"`
	StartTime     time.Time `gorm:"not null;index"`
	EndTime       *time.Time
	Duration      int64     // Duration in milliseconds
	SpanCount     int       `gorm:"default:0"`
	Tags          string    `gorm:"type:jsonb"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// SpanModel represents a span in the database
type SpanModel struct {
	ID            string    `gorm:"primaryKey;type:uuid"`
	TraceID       string    `gorm:"not null;index"`
	ParentSpanID  string    `gorm:"index"`
	OperationName string    `gorm:"not null"`
	Service       string    `gorm:"not null;index"`
	StartTime     time.Time `gorm:"not null"`
	EndTime       *time.Time
	Duration      int64     // Duration in milliseconds
	Tags          string    `gorm:"type:jsonb"`
	Logs          string    `gorm:"type:jsonb"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// AlertModel represents an alert in the database
type AlertModel struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Severity    string    `gorm:"not null;index"`
	Service     string    `gorm:"not null;index"`
	Condition   string    `gorm:"not null"`
	Status      string    `gorm:"not null;default:'active';index"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	ResolvedAt  *time.Time
}

// Migrate creates or updates the database tables
func (r *observabilityRepository) Migrate() error {
	return r.db.AutoMigrate(
		&LogEntryModel{},
		&MetricModel{},
		&TraceModel{},
		&SpanModel{},
		&AlertModel{},
	)
}

// Log storage methods

// SaveLogEntry saves a log entry to the database
func (r *observabilityRepository) SaveLogEntry(ctx context.Context, entry *LogEntry) error {
	fieldsJSON, _ := json.Marshal(entry.Fields)
	errorJSON := ""
	if entry.Error != nil {
		errorData, _ := json.Marshal(entry.Error)
		errorJSON = string(errorData)
	}

	model := &LogEntryModel{
		ID:        entry.ID.String(),
		Level:     string(entry.Level),
		Message:   entry.Message,
		Service:   entry.Service,
		Version:   entry.Version,
		RequestID: entry.RequestID,
		UserID:    entry.UserID,
		TenantID:  entry.TenantID,
		TraceID:   entry.TraceID,
		SpanID:    entry.SpanID,
		Fields:    string(fieldsJSON),
		Error:     errorJSON,
		CreatedAt: entry.Timestamp,
	}

	return r.db.WithContext(ctx).Create(model).Error
}

// GetLogEntries retrieves log entries with filters
func (r *observabilityRepository) GetLogEntries(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*LogEntry, error) {
	var models []LogEntryModel
	
	query := r.db.WithContext(ctx).Model(&LogEntryModel{})
	
	// Apply filters
	for key, value := range filters {
		switch key {
		case "level":
			query = query.Where("level = ?", value)
		case "service":
			query = query.Where("service = ?", value)
		case "tenant_id":
			query = query.Where("tenant_id = ?", value)
		case "user_id":
			query = query.Where("user_id = ?", value)
		case "trace_id":
			query = query.Where("trace_id = ?", value)
		case "from":
			query = query.Where("created_at >= ?", value)
		case "to":
			query = query.Where("created_at <= ?", value)
		}
	}
	
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&models).Error
	if err != nil {
		return nil, err
	}

	entries := make([]*LogEntry, len(models))
	for i, model := range models {
		entry, err := r.logModelToEntry(&model)
		if err != nil {
			continue
		}
		entries[i] = entry
	}

	return entries, nil
}

// Metric storage methods

// SaveMetric saves a metric to the database
func (r *observabilityRepository) SaveMetric(ctx context.Context, metric *Metric) error {
	tagsJSON, _ := json.Marshal(metric.Tags)

	model := &MetricModel{
		ID:        metric.ID.String(),
		Name:      metric.Name,
		Type:      string(metric.Type),
		Value:     metric.Value,
		Service:   metric.Service,
		Version:   metric.Version,
		Tags:      string(tagsJSON),
		CreatedAt: metric.Timestamp,
	}

	return r.db.WithContext(ctx).Create(model).Error
}

// GetMetrics retrieves metrics with filters
func (r *observabilityRepository) GetMetrics(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*Metric, error) {
	var models []MetricModel
	
	query := r.db.WithContext(ctx).Model(&MetricModel{})
	
	// Apply filters
	for key, value := range filters {
		switch key {
		case "name":
			query = query.Where("name = ?", value)
		case "type":
			query = query.Where("type = ?", value)
		case "service":
			query = query.Where("service = ?", value)
		case "from":
			query = query.Where("created_at >= ?", value)
		case "to":
			query = query.Where("created_at <= ?", value)
		}
	}
	
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&models).Error
	if err != nil {
		return nil, err
	}

	metrics := make([]*Metric, len(models))
	for i, model := range models {
		metric, err := r.metricModelToMetric(&model)
		if err != nil {
			continue
		}
		metrics[i] = metric
	}

	return metrics, nil
}

// Trace storage methods

// SaveTrace saves a trace to the database
func (r *observabilityRepository) SaveTrace(ctx context.Context, trace *Trace) error {
	tagsJSON, _ := json.Marshal(trace.Tags)

	model := &TraceModel{
		ID:            trace.ID,
		OperationName: trace.OperationName,
		Service:       trace.Service,
		Version:       trace.Version,
		StartTime:     trace.StartTime,
		SpanCount:     trace.SpanCount,
		Tags:          string(tagsJSON),
	}

	if !trace.EndTime.IsZero() {
		model.EndTime = &trace.EndTime
		model.Duration = trace.Duration.Milliseconds()
	}

	return r.db.WithContext(ctx).Create(model).Error
}

// SaveSpan saves a span to the database
func (r *observabilityRepository) SaveSpan(ctx context.Context, span *Span) error {
	tagsJSON, _ := json.Marshal(span.Tags)
	logsJSON, _ := json.Marshal(span.Logs)

	model := &SpanModel{
		ID:            span.ID,
		TraceID:       span.TraceID,
		ParentSpanID:  span.ParentSpanID,
		OperationName: span.OperationName,
		Service:       span.Service,
		StartTime:     span.StartTime,
		Tags:          string(tagsJSON),
		Logs:          string(logsJSON),
	}

	if !span.EndTime.IsZero() {
		model.EndTime = &span.EndTime
		model.Duration = span.Duration.Milliseconds()
	}

	return r.db.WithContext(ctx).Create(model).Error
}

// GetTrace retrieves a trace by ID
func (r *observabilityRepository) GetTrace(ctx context.Context, traceID string) (*Trace, error) {
	var model TraceModel
	err := r.db.WithContext(ctx).Where("id = ?", traceID).First(&model).Error
	if err != nil {
		return nil, err
	}

	return r.traceModelToTrace(&model)
}

// GetSpansByTraceID retrieves spans for a trace
func (r *observabilityRepository) GetSpansByTraceID(ctx context.Context, traceID string) ([]*Span, error) {
	var models []SpanModel
	err := r.db.WithContext(ctx).Where("trace_id = ?", traceID).Order("start_time").Find(&models).Error
	if err != nil {
		return nil, err
	}

	spans := make([]*Span, len(models))
	for i, model := range models {
		span, err := r.spanModelToSpan(&model)
		if err != nil {
			continue
		}
		spans[i] = span
	}

	return spans, nil
}

// Alert storage methods

// SaveAlert saves an alert to the database
func (r *observabilityRepository) SaveAlert(ctx context.Context, alert *Alert) error {
	model := &AlertModel{
		ID:          alert.ID,
		Title:       alert.Title,
		Description: alert.Description,
		Severity:    string(alert.Severity),
		Service:     alert.Service,
		Condition:   alert.Condition,
		Status:      string(alert.Status),
		CreatedAt:   alert.CreatedAt,
		UpdatedAt:   alert.UpdatedAt,
	}

	if !alert.ResolvedAt.IsZero() {
		model.ResolvedAt = &alert.ResolvedAt
	}

	return r.db.WithContext(ctx).Create(model).Error
}

// GetAlerts retrieves alerts with filters
func (r *observabilityRepository) GetAlerts(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*Alert, error) {
	var models []AlertModel
	
	query := r.db.WithContext(ctx).Model(&AlertModel{})
	
	// Apply filters
	for key, value := range filters {
		switch key {
		case "severity":
			query = query.Where("severity = ?", value)
		case "service":
			query = query.Where("service = ?", value)
		case "status":
			query = query.Where("status = ?", value)
		}
	}
	
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&models).Error
	if err != nil {
		return nil, err
	}

	alerts := make([]*Alert, len(models))
	for i, model := range models {
		alert := r.alertModelToAlert(&model)
		alerts[i] = alert
	}

	return alerts, nil
}

// Helper methods for model conversion

func (r *observabilityRepository) logModelToEntry(model *LogEntryModel) (*LogEntry, error) {
	var fields map[string]interface{}
	if model.Fields != "" {
		json.Unmarshal([]byte(model.Fields), &fields)
	}

	var errorDetails *ErrorDetails
	if model.Error != "" {
		json.Unmarshal([]byte(model.Error), &errorDetails)
	}

	return &LogEntry{
		Level:     LogLevel(model.Level),
		Message:   model.Message,
		Timestamp: model.CreatedAt,
		Service:   model.Service,
		Version:   model.Version,
		RequestID: model.RequestID,
		UserID:    model.UserID,
		TenantID:  model.TenantID,
		TraceID:   model.TraceID,
		SpanID:    model.SpanID,
		Fields:    fields,
		Error:     errorDetails,
	}, nil
}

func (r *observabilityRepository) metricModelToMetric(model *MetricModel) (*Metric, error) {
	var tags map[string]string
	if model.Tags != "" {
		json.Unmarshal([]byte(model.Tags), &tags)
	}

	return &Metric{
		Name:      model.Name,
		Type:      MetricType(model.Type),
		Value:     model.Value,
		Tags:      tags,
		Timestamp: model.CreatedAt,
		Service:   model.Service,
		Version:   model.Version,
	}, nil
}

func (r *observabilityRepository) traceModelToTrace(model *TraceModel) (*Trace, error) {
	var tags map[string]interface{}
	if model.Tags != "" {
		json.Unmarshal([]byte(model.Tags), &tags)
	}

	trace := &Trace{
		ID:            model.ID,
		OperationName: model.OperationName,
		StartTime:     model.StartTime,
		Service:       model.Service,
		Version:       model.Version,
		Tags:          tags,
		SpanCount:     model.SpanCount,
	}

	if model.EndTime != nil {
		trace.EndTime = *model.EndTime
		trace.Duration = time.Duration(model.Duration) * time.Millisecond
	}

	return trace, nil
}

func (r *observabilityRepository) spanModelToSpan(model *SpanModel) (*Span, error) {
	var tags map[string]interface{}
	if model.Tags != "" {
		json.Unmarshal([]byte(model.Tags), &tags)
	}

	var logs []SpanLog
	if model.Logs != "" {
		json.Unmarshal([]byte(model.Logs), &logs)
	}

	span := &Span{
		ID:            model.ID,
		TraceID:       model.TraceID,
		ParentSpanID:  model.ParentSpanID,
		OperationName: model.OperationName,
		StartTime:     model.StartTime,
		Service:       model.Service,
		Tags:          tags,
		Logs:          logs,
	}

	if model.EndTime != nil {
		span.EndTime = *model.EndTime
		span.Duration = time.Duration(model.Duration) * time.Millisecond
	}

	return span, nil
}

func (r *observabilityRepository) alertModelToAlert(model *AlertModel) *Alert {
	alert := &Alert{
		ID:          model.ID,
		Title:       model.Title,
		Description: model.Description,
		Severity:    AlertSeverity(model.Severity),
		Service:     model.Service,
		Condition:   model.Condition,
		Status:      AlertStatus(model.Status),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}

	if model.ResolvedAt != nil {
		alert.ResolvedAt = *model.ResolvedAt
	}

	return alert
}
