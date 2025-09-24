package contact

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"ecommerce-saas/internal/shared/utils"
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
	// Validate required fields
	if contact.Name == "" {
		return errors.New("contact name is required")
	}
	if contact.Email == "" {
		return errors.New("contact email is required")
	}
	if contact.Subject == "" {
		return errors.New("contact subject is required")
	}
	if contact.Message == "" {
		return errors.New("contact message is required")
	}
	if contact.TenantID == uuid.Nil {
		return errors.New("tenant ID is required")
	}

	// Normalize and validate email
	contact.Email = strings.ToLower(strings.TrimSpace(contact.Email))
	if !strings.Contains(contact.Email, "@") {
		return errors.New("invalid email format")
	}

	// Set defaults
	if contact.Status == "" {
		contact.Status = StatusNew
	}
	if contact.Priority == "" {
		contact.Priority = PriorityMedium
	}
	if contact.Type == "" {
		contact.Type = TypeGeneral
	}

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
	// Validate required fields
	if contact.ID == uuid.Nil {
		return errors.New("contact ID is required")
	}
	if contact.TenantID == uuid.Nil {
		return errors.New("tenant ID is required")
	}
	if contact.Name == "" {
		return errors.New("contact name is required")
	}
	if contact.Email == "" {
		return errors.New("contact email is required")
	}
	if contact.Subject == "" {
		return errors.New("contact subject is required")
	}
	if contact.Message == "" {
		return errors.New("contact message is required")
	}

	// Normalize email
	contact.Email = strings.ToLower(strings.TrimSpace(contact.Email))
	if !strings.Contains(contact.Email, "@") {
		return errors.New("invalid email format")
	}

	// Ensure tenant isolation by adding WHERE clause
	return r.db.WithContext(ctx).
		Where("tenant_id = ?", contact.TenantID).
		Save(contact).Error
}

func (r *repository) DeleteContact(ctx context.Context, tenantID, contactID uuid.UUID) error {
	// Validate inputs
	if contactID == uuid.Nil {
		return errors.New("contact ID is required")
	}
	if tenantID == uuid.Nil {
		return errors.New("tenant ID is required")
	}

	// Start transaction for cascade delete
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Soft delete related contact replies first
	if err := tx.Where("contact_id = ? AND tenant_id = ?", contactID, tenantID).Delete(&ContactReply{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Soft delete the contact
	if err := tx.Where("id = ? AND tenant_id = ?", contactID, tenantID).Delete(&Contact{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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
	if len(filter.Tags) > 0 {
		// Use PostgreSQL JSONB array containment operator for efficient tag filtering
		for _, tag := range filter.Tags {
			query = query.Where("tags ? ?", tag)
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
	if reply.ContactID == uuid.Nil {
		return errors.New("contact ID is required")
	}

	if strings.TrimSpace(reply.Content) == "" {
		return errors.New("content is required")
	}

	reply.ID = uuid.New()
	reply.CreatedAt = time.Now()
	reply.UpdatedAt = time.Now()

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

func (r *repository) UpdateContactReplyStatus(ctx context.Context, replyID uuid.UUID, status ContactStatus) error {
	// ContactReply doesn't have a status field, this method should update the parent contact status
	var reply ContactReply
	if err := r.db.WithContext(ctx).Where("id = ?", replyID).First(&reply).Error; err != nil {
		return err
	}
	
	return r.db.WithContext(ctx).
		Model(&Contact{}).
		Where("id = ?", reply.ContactID).
		Update("status", status).
		Error
}

func (r *repository) DeleteContactReply(ctx context.Context, tenantID, replyID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", replyID).
		Delete(&ContactReply{}).Error
}

func (r *repository) GetContactRepliesByType(ctx context.Context, tenantID uuid.UUID, replyType ContactType, limit, offset int) ([]*ContactReply, error) {
	var replies []*ContactReply
	err := r.db.WithContext(ctx).
		Joins("JOIN contacts ON contact_replies.contact_id = contacts.id").
		Where("contacts.tenant_id = ? AND contacts.type = ?", tenantID, replyType).
		Order("contact_replies.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&replies).Error
	return replies, err
}

func (r *repository) GetContactRepliesByStatus(ctx context.Context, tenantID uuid.UUID, status ContactStatus, limit, offset int) ([]*ContactReply, error) {
	var replies []*ContactReply
	err := r.db.WithContext(ctx).
		Joins("JOIN contacts ON contact_replies.contact_id = contacts.id").
		Where("contacts.tenant_id = ? AND contacts.status = ?", tenantID, status).
		Order("contact_replies.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&replies).Error
	return replies, err
}

// Contact forms
func (r *repository) CreateContactForm(ctx context.Context, form *ContactForm) error {
	// Validate required fields
	if form.Name == "" {
		return errors.New("form name is required")
	}
	if form.TenantID == uuid.Nil {
		return errors.New("tenant ID is required")
	}

	// Validate default type
	if form.DefaultType == "" {
		form.DefaultType = TypeGeneral
	}

	// Validate form type
	validTypes := []ContactType{TypeGeneral, TypeSupport, TypeSales, TypeTechnical}
	validType := false
	for _, validT := range validTypes {
		if form.DefaultType == validT {
			validType = true
			break
		}
	}
	if !validType {
		return errors.New("invalid default type")
	}

	// Set default active status
	if !form.IsActive {
		form.IsActive = true
	}

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
	// Validate required fields
	if form.ID == uuid.Nil {
		return errors.New("form ID is required")
	}
	if form.TenantID == uuid.Nil {
		return errors.New("tenant ID is required")
	}
	if form.Name == "" {
		return errors.New("form name is required")
	}

	// Validate default type
	if form.DefaultType == "" {
		form.DefaultType = TypeGeneral
	}

	// Validate form type
	validTypes := []ContactType{TypeGeneral, TypeSupport, TypeSales, TypeTechnical}
	validType := false
	for _, validT := range validTypes {
		if form.DefaultType == validT {
			validType = true
			break
		}
	}
	if !validType {
		return errors.New("invalid default type")
	}

	// Ensure tenant isolation
	return r.db.WithContext(ctx).
		Where("tenant_id = ?", form.TenantID).
		Save(form).Error
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
	// Validate required fields
	if template.Name == "" {
		return errors.New("template name is required")
	}
	if template.TenantID == uuid.Nil {
		return errors.New("tenant ID is required")
	}
	if template.Type == "" {
		return errors.New("template type is required")
	}
	if template.Subject == "" {
		return errors.New("template subject is required")
	}
	if template.Content == "" {
		return errors.New("template content is required")
	}

	// Validate template type
	validTypes := []ContactTemplateType{TemplateTypeAutoReply, TemplateTypeFollowUp, TemplateTypeResolution, TemplateTypeEscalation}
	validType := false
	for _, validT := range validTypes {
		if template.Type == validT {
			validType = true
			break
		}
	}
	if !validType {
		return errors.New("invalid template type")
	}

	// Set default active status
	if !template.IsActive {
		template.IsActive = true
	}

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
	// Validate required fields
	if template.ID == uuid.Nil {
		return errors.New("template ID is required")
	}
	if template.TenantID == uuid.Nil {
		return errors.New("tenant ID is required")
	}
	if template.Name == "" {
		return errors.New("template name is required")
	}
	if template.Type == "" {
		return errors.New("template type is required")
	}
	if template.Subject == "" {
		return errors.New("template subject is required")
	}
	if template.Content == "" {
		return errors.New("template content is required")
	}

	// Validate template type
	validTypes := []ContactTemplateType{TemplateTypeAutoReply, TemplateTypeFollowUp, TemplateTypeResolution, TemplateTypeEscalation}
	validType := false
	for _, validT := range validTypes {
		if template.Type == validT {
			validType = true
			break
		}
	}
	if !validType {
		return errors.New("invalid template type")
	}

	// Ensure tenant isolation
	return r.db.WithContext(ctx).
		Where("tenant_id = ?", template.TenantID).
		Save(template).Error
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
			// Return default settings
			return &ContactSettings{
				TenantID:              tenantID,
				ContactEmail:          "contact@example.com",
				SupportEmail:          "support@example.com",
				BusinessHours:         "9 AM - 5 PM",
				Timezone:              "Asia/Dhaka",
				AutoReplyEnabled:      true,
				SLAResponseTime:       24,
				SLAResolutionTime:     72,
				EmailNotifications:    true,
				NotifyOnNewContact:    true,
				AllowAnonymousContact: true,
				EnableSpamFilter:      true,
				MaxDailySubmissions:   10,
			}, nil
		}
		return nil, err
	}
	return &settings, nil
}

func (r *repository) UpdateContactSettings(ctx context.Context, settings *ContactSettings) error {
	// Validate required fields
	if settings.TenantID == uuid.Nil {
		return errors.New("tenant ID is required")
	}
	if settings.ContactEmail == "" {
		return errors.New("contact email is required")
	}
	if settings.SupportEmail == "" {
		return errors.New("support email is required")
	}

	// Validate email formats
	if !utils.IsValidEmail(settings.ContactEmail) {
		return errors.New("invalid contact email format")
	}
	if !utils.IsValidEmail(settings.SupportEmail) {
		return errors.New("invalid support email format")
	}

	// Validate SLA times (must be positive)
	if settings.SLAResponseTime <= 0 {
		return errors.New("SLA response time must be positive")
	}
	if settings.SLAResolutionTime <= 0 {
		return errors.New("SLA resolution time must be positive")
	}

	// Validate max daily submissions (must be positive)
	if settings.MaxDailySubmissions <= 0 {
		return errors.New("max daily submissions must be positive")
	}

	// Ensure tenant isolation
	return r.db.WithContext(ctx).
		Where("tenant_id = ?", settings.TenantID).
		Save(settings).Error
}

// Analytics
func (r *repository) GetContactAnalytics(ctx context.Context, tenantID uuid.UUID, period AnalyticsPeriod) (*ContactAnalytics, error) {
	now := time.Now()

	// TODO: Implement comprehensive analytics queries
	// This is a complex query that would need to aggregate data from multiple periods
	// For now, return basic structure
	analytics := ContactAnalytics{
		TenantID: tenantID,
		Date:     now,
	}

	return &analytics, nil
}

func (r *repository) GetContactMetrics(ctx context.Context, tenantID uuid.UUID, from, to time.Time) (*ContactMetrics, error) {
	// TODO: Implement comprehensive metrics calculation
	// This would involve complex aggregation queries for:
	// - Total contacts, new contacts, resolved contacts
	// - Average response time, resolution time
	// - Customer satisfaction scores
	// - Agent performance metrics

	metrics := ContactMetrics{
		TenantID:  tenantID,
		StartDate: from,
		EndDate:   to,
	}

	return &metrics, nil
}