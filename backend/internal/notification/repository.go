package notification

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Notification CRUD
	Create(notification *Notification) error
	GetByID(tenantID, notificationID uuid.UUID) (*Notification, error)
	Update(notification *Notification) error
	Delete(tenantID, notificationID uuid.UUID) error
	List(tenantID uuid.UUID, userID *uuid.UUID, offset, limit int) ([]*Notification, int64, error)
	ListByStatus(tenantID uuid.UUID, status string, offset, limit int) ([]*Notification, int64, error)

	// Template operations
	CreateTemplate(template *NotificationTemplate) error
	GetTemplate(tenantID, templateID uuid.UUID) (*NotificationTemplate, error)
	UpdateTemplate(template *NotificationTemplate) error
	DeleteTemplate(tenantID, templateID uuid.UUID) error
	ListTemplates(tenantID uuid.UUID, notificationType, channel string) ([]*NotificationTemplate, error)

	// Preference operations
	GetPreferences(tenantID, userID uuid.UUID, channel string) (*NotificationPreference, error)
	UpdatePreferences(preference *NotificationPreference) error
	CreatePreferences(preference *NotificationPreference) error

	// Log operations
	CreateLog(log *NotificationLog) error
	GetLogs(notificationID uuid.UUID) ([]*NotificationLog, error)

	// Stats operations
	GetNotificationStats(tenantID uuid.UUID) (map[string]int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(notification *Notification) error {
	return r.db.Create(notification).Error
}

func (r *repository) GetByID(tenantID, notificationID uuid.UUID) (*Notification, error) {
	var notification Notification
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, notificationID).First(&notification).Error
	return &notification, err
}

func (r *repository) Update(notification *Notification) error {
	return r.db.Save(notification).Error
}

func (r *repository) Delete(tenantID, notificationID uuid.UUID) error {
	return r.db.Where("tenant_id = ? AND id = ?", tenantID, notificationID).Delete(&Notification{}).Error
}

func (r *repository) List(tenantID uuid.UUID, userID *uuid.UUID, offset, limit int) ([]*Notification, int64, error) {
	var notifications []*Notification
	var total int64

	query := r.db.Model(&Notification{}).Where("tenant_id = ?", tenantID)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *repository) ListByStatus(tenantID uuid.UUID, status string, offset, limit int) ([]*Notification, int64, error) {
	var notifications []*Notification
	var total int64

	query := r.db.Model(&Notification{}).Where("tenant_id = ? AND status = ?", tenantID, status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// Template operations
func (r *repository) CreateTemplate(template *NotificationTemplate) error {
	return r.db.Create(template).Error
}

func (r *repository) GetTemplate(tenantID, templateID uuid.UUID) (*NotificationTemplate, error) {
	var template NotificationTemplate
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, templateID).First(&template).Error
	return &template, err
}

func (r *repository) UpdateTemplate(template *NotificationTemplate) error {
	return r.db.Save(template).Error
}

func (r *repository) DeleteTemplate(tenantID, templateID uuid.UUID) error {
	return r.db.Where("tenant_id = ? AND id = ?", tenantID, templateID).Delete(&NotificationTemplate{}).Error
}

func (r *repository) ListTemplates(tenantID uuid.UUID, notificationType, channel string) ([]*NotificationTemplate, error) {
	var templates []*NotificationTemplate

	query := r.db.Where("tenant_id = ? AND is_active = ?", tenantID, true)

	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	if channel != "" {
		query = query.Where("channel = ?", channel)
	}

	err := query.Order("name").Find(&templates).Error
	return templates, err
}

// Preference operations
func (r *repository) GetPreferences(tenantID, userID uuid.UUID, channel string) (*NotificationPreference, error) {
	var preference NotificationPreference
	err := r.db.Where("tenant_id = ? AND user_id = ? AND channel = ?", tenantID, userID, channel).First(&preference).Error
	return &preference, err
}

func (r *repository) UpdatePreferences(preference *NotificationPreference) error {
	return r.db.Save(preference).Error
}

func (r *repository) CreatePreferences(preference *NotificationPreference) error {
	return r.db.Create(preference).Error
}

// Log operations
func (r *repository) CreateLog(log *NotificationLog) error {
	return r.db.Create(log).Error
}

func (r *repository) GetLogs(notificationID uuid.UUID) ([]*NotificationLog, error) {
	var logs []*NotificationLog
	err := r.db.Where("notification_id = ?", notificationID).Order("created_at").Find(&logs).Error
	return logs, err
}

// Stats operations
func (r *repository) GetNotificationStats(tenantID uuid.UUID) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Count total notifications
	var total int64
	if err := r.db.Model(&Notification{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// Count by status
	var sent, delivered, failed int64
	if err := r.db.Model(&Notification{}).Where("tenant_id = ? AND status = ?", tenantID, StatusSent).Count(&sent).Error; err != nil {
		return nil, err
	}
	stats["sent"] = sent

	if err := r.db.Model(&Notification{}).Where("tenant_id = ? AND status = ?", tenantID, StatusDelivered).Count(&delivered).Error; err != nil {
		return nil, err
	}
	stats["delivered"] = delivered

	if err := r.db.Model(&Notification{}).Where("tenant_id = ? AND status = ?", tenantID, StatusFailed).Count(&failed).Error; err != nil {
		return nil, err
	}
	stats["failed"] = failed

	// By type
	var email, sms, push, inApp int64
	r.db.Model(&Notification{}).Where("tenant_id = ? AND type = ?", tenantID, TypeEmail).Count(&email)
	stats["email"] = email
	r.db.Model(&Notification{}).Where("tenant_id = ? AND type = ?", tenantID, TypeSMS).Count(&sms)
	stats["sms"] = sms
	r.db.Model(&Notification{}).Where("tenant_id = ? AND type = ?", tenantID, TypePush).Count(&push)
	stats["push"] = push
	r.db.Model(&Notification{}).Where("tenant_id = ? AND type = ?", tenantID, TypeInApp).Count(&inApp)
	stats["in_app"] = inApp

	return stats, nil
}
