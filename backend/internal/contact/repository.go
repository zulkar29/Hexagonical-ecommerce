package contact

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Contact management
	CreateContact(ctx context.Context, contact *Contact) error
	GetContactByID(ctx context.Context, tenantID, contactID uuid.UUID) (*Contact, error)
	UpdateContact(ctx context.Context, contact *Contact) error
	DeleteContact(ctx context.Context, tenantID, contactID uuid.UUID) error
	ListContacts(ctx context.Context, tenantID uuid.UUID, filter ContactFilter) ([]*Contact, int64, error)

	// Contact replies
	CreateContactReply(ctx context.Context, reply *ContactReply) error
	ListContactReplies(ctx context.Context, tenantID, contactID uuid.UUID) ([]*ContactReply, error)
	DeleteContactReply(ctx context.Context, tenantID, replyID uuid.UUID) error

	// Contact forms
	CreateContactForm(ctx context.Context, form *ContactForm) error
	GetContactFormByID(ctx context.Context, tenantID, formID uuid.UUID) (*ContactForm, error)
	UpdateContactForm(ctx context.Context, form *ContactForm) error
	DeleteContactForm(ctx context.Context, tenantID, formID uuid.UUID) error
	ListContactForms(ctx context.Context, tenantID uuid.UUID, filter ContactFormFilter) ([]*ContactForm, int64, error)
	GetActiveContactForm(ctx context.Context, tenantID uuid.UUID, formType ContactFormType) (*ContactForm, error)

	// Templates
	CreateContactTemplate(ctx context.Context, template *ContactTemplate) error
	GetContactTemplateByID(ctx context.Context, tenantID, templateID uuid.UUID) (*ContactTemplate, error)
	UpdateContactTemplate(ctx context.Context, template *ContactTemplate) error
	DeleteContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error
	ListContactTemplates(ctx context.Context, tenantID uuid.UUID, filter ContactTemplateFilter) ([]*ContactTemplate, int64, error)
	GetContactTemplateByType(ctx context.Context, tenantID uuid.UUID, templateType ContactTemplateType) (*ContactTemplate, error)

	// Settings
	GetContactSettings(ctx context.Context, tenantID uuid.UUID) (*ContactSettings, error)
	UpdateContactSettings(ctx context.Context, settings *ContactSettings) error

	// Analytics
	GetContactAnalytics(ctx context.Context, tenantID uuid.UUID, period AnalyticsPeriod) (*ContactAnalytics, error)
	GetContactMetrics(ctx context.Context, tenantID uuid.UUID, from, to time.Time) (*ContactMetrics, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Contact management
func (r *repository) CreateContact(ctx context.Context, contact *Contact) error {
	// TODO: Implement contact creation with proper validation and tenant isolation
	return r.db.WithContext(ctx).Create(contact).Error
}

func (r *repository) GetContactByID(ctx context.Context, tenantID, contactID uuid.UUID) (*Contact, error) {
	var contact Contact
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, contactID).
		Preload("Replies").
		Preload("Form").
		First(&contact).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *repository) UpdateContact(ctx context.Context, contact *Contact) error {
	// TODO: Implement contact update with proper validation and audit logging
	return r.db.WithContext(ctx).Save(contact).Error
}

func (r *repository) DeleteContact(ctx context.Context, tenantID, contactID uuid.UUID) error {
	// TODO: Implement soft delete with cascade to replies
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, contactID).
		Delete(&Contact{}).Error
}

func (r *repository) ListContacts(ctx context.Context, tenantID uuid.UUID, filter ContactFilter) ([]*Contact, int64, error) {
	var contacts []*Contact
	var total int64

	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	// Apply filters
	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}
	if len(filter.Priority) > 0 {
		query = query.Where("priority IN ?", filter.Priority)
	}
	if len(filter.Type) > 0 {
		query = query.Where("type IN ?", filter.Type)
	}
	if filter.AssignedToID != nil {
		query = query.Where("assigned_to_id = ?", *filter.AssignedToID)
	}
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	if len(filter.Source) > 0 {
		query = query.Where("source IN ?", filter.Source)
	}
	if filter.Search != "" {
		query = query.Where("(name ILIKE ? OR email ILIKE ? OR subject ILIKE ? OR message ILIKE ?)", 
			"%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	if filter.Tags != nil && len(filter.Tags) > 0 {
		// TODO: Implement proper JSONB array containment query for PostgreSQL
		for _, tag := range filter.Tags {
			query = query.Where("tags::text ILIKE ?", "%"+tag+"%")
		}
	}

	// Count total
	if err := query.Model(&Contact{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", filter.SortBy, direction))
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	err := query.Preload("Replies").Find(&contacts).Error
	return contacts, total, err
}

// Contact replies
func (r *repository) CreateContactReply(ctx context.Context, reply *ContactReply) error {
	// TODO: Implement reply creation with notification triggers
	return r.db.WithContext(ctx).Create(reply).Error
}

func (r *repository) ListContactReplies(ctx context.Context, tenantID, contactID uuid.UUID) ([]*ContactReply, error) {
	var replies []*ContactReply
	err := r.db.WithContext(ctx).
		Where("contact_id = ?", contactID).
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}

func (r *repository) DeleteContactReply(ctx context.Context, tenantID, replyID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", replyID).
		Delete(&ContactReply{}).Error
}

// Contact forms
func (r *repository) CreateContactForm(ctx context.Context, form *ContactForm) error {
	// TODO: Implement form creation with validation
	return r.db.WithContext(ctx).Create(form).Error
}

func (r *repository) GetContactFormByID(ctx context.Context, tenantID, formID uuid.UUID) (*ContactForm, error) {
	var form ContactForm
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, formID).
		First(&form).Error
	if err != nil {
		return nil, err
	}
	return &form, nil
}

func (r *repository) UpdateContactForm(ctx context.Context, form *ContactForm) error {
	// TODO: Implement form update with validation
	return r.db.WithContext(ctx).Save(form).Error
}

func (r *repository) DeleteContactForm(ctx context.Context, tenantID, formID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, formID).
		Delete(&ContactForm{}).Error
}

func (r *repository) ListContactForms(ctx context.Context, tenantID uuid.UUID, filter ContactFormFilter) ([]*ContactForm, int64, error) {
	var forms []*ContactForm
	var total int64

	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	// Apply filters
	if len(filter.Type) > 0 {
		query = query.Where("default_type IN ?", filter.Type)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	// Count total
	if err := query.Model(&ContactForm{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortDesc {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", filter.SortBy, direction))
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	err := query.Find(&forms).Error
	return forms, total, err
}

func (r *repository) GetActiveContactForm(ctx context.Context, tenantID uuid.UUID, formType ContactFormType) (*ContactForm, error) {
	var form ContactForm
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND default_type = ? AND is_active = ?", tenantID, formType, true).
		First(&form).Error
	if err != nil {
		return nil, err
	}
	return &form, nil
}

// Templates
func (r *repository) CreateContactTemplate(ctx context.Context, template *ContactTemplate) error {
	// TODO: Implement template creation with validation
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *repository) GetContactTemplateByID(ctx context.Context, tenantID, templateID uuid.UUID) (*ContactTemplate, error) {
	var template ContactTemplate
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, templateID).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *repository) UpdateContactTemplate(ctx context.Context, template *ContactTemplate) error {
	// TODO: Implement template update with validation
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *repository) DeleteContactTemplate(ctx context.Context, tenantID, templateID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, templateID).
		Delete(&ContactTemplate{}).Error
}

func (r *repository) ListContactTemplates(ctx context.Context, tenantID uuid.UUID, filter ContactTemplateFilter) ([]*ContactTemplate, int64, error) {
	var templates []*ContactTemplate
	var total int64

	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	// Apply filters
	if len(filter.Type) > 0 {
		query = query.Where("type IN ?", filter.Type)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	// Count total
	if err := query.Model(&ContactTemplate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortDesc {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", filter.SortBy, direction))
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	err := query.Find(&templates).Error
	return templates, total, err
}

func (r *repository) GetContactTemplateByType(ctx context.Context, tenantID uuid.UUID, templateType ContactTemplateType) (*ContactTemplate, error) {
	var template ContactTemplate
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND type = ? AND is_active = ?", tenantID, templateType, true).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Settings
func (r *repository) GetContactSettings(ctx context.Context, tenantID uuid.UUID) (*ContactSettings, error) {
	var settings ContactSettings
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default settings if not found
			return &ContactSettings{
				TenantID:             tenantID,
				ContactEmail:         "contact@example.com",
				AutoAssignEnabled:    false,
				AutoReplyEnabled:     true,
				EmailNotifications:   true,
				SLAResponseTime:      24, // 24 hours
				SLAResolutionTime:    72, // 72 hours
				AllowAnonymousContact: true,
				EnableSpamFilter:     true,
				DataRetentionDays:    365,
				ConsentRequired:      true,
			}, nil
		}
		return nil, err
	}
	return &settings, nil
}

func (r *repository) UpdateContactSettings(ctx context.Context, settings *ContactSettings) error {
	// TODO: Implement settings update with validation
	return r.db.WithContext(ctx).Save(settings).Error
}

// Analytics
func (r *repository) GetContactAnalytics(ctx context.Context, tenantID uuid.UUID, period AnalyticsPeriod) (*ContactAnalytics, error) {
	now := time.Now()
	
	var analytics ContactAnalytics
	
	// TODO: Implement comprehensive analytics queries
	// This is a complex query that would need to aggregate data from multiple periods
	// For now, return basic structure
	analytics = ContactAnalytics{
		TenantID: tenantID,
		Date:     now,
	}
	
	return &analytics, nil
}

func (r *repository) GetContactMetrics(ctx context.Context, tenantID uuid.UUID, from, to time.Time) (*ContactMetrics, error) {
	var metrics ContactMetrics

	// TODO: Implement comprehensive metrics calculation
	// This would involve complex aggregation queries for:
	// - Total contacts, new contacts, resolved contacts
	// - Average response time, resolution time
	// - Customer satisfaction scores
	// - Agent performance metrics

	metrics = ContactMetrics{
		TenantID:  tenantID,
		StartDate: from,
		EndDate:   to,
	}

	return &metrics, nil
}