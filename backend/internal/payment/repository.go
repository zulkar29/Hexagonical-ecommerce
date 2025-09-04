package payment

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(payment *Payment) error
	GetByID(tenantID, paymentID uuid.UUID) (*Payment, error)
	GetByOrderID(tenantID, orderID uuid.UUID) ([]*Payment, error)
	GetByTransactionID(transactionID string) (*Payment, error)
	Update(payment *Payment) error
	Delete(tenantID, paymentID uuid.UUID) error
	List(tenantID uuid.UUID, orderID *uuid.UUID, offset, limit int) ([]*Payment, int64, error)
	
	// Refund operations
	CreateRefund(refund *Refund) error
	GetRefund(tenantID, refundID uuid.UUID) (*Refund, error)
	ListRefunds(tenantID uuid.UUID, paymentID *uuid.UUID, offset, limit int) ([]*Refund, int64, error)
	
	// Payment method operations
	CreatePaymentMethod(method *PaymentMethod) error
	GetPaymentMethod(tenantID, methodID uuid.UUID) (*PaymentMethod, error)
	ListPaymentMethods(tenantID, userID uuid.UUID) ([]*PaymentMethod, error)
	UpdatePaymentMethod(method *PaymentMethod) error
	DeletePaymentMethod(tenantID, methodID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(payment *Payment) error {
	return r.db.Create(payment).Error
}

func (r *repository) GetByID(tenantID, paymentID uuid.UUID) (*Payment, error) {
	var payment Payment
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, paymentID).First(&payment).Error
	return &payment, err
}

func (r *repository) GetByOrderID(tenantID, orderID uuid.UUID) ([]*Payment, error) {
	var payments []*Payment
	err := r.db.Where("tenant_id = ? AND order_id = ?", tenantID, orderID).Find(&payments).Error
	return payments, err
}

func (r *repository) GetByTransactionID(transactionID string) (*Payment, error) {
	var payment Payment
	err := r.db.Where("payment_intent_id = ?", transactionID).First(&payment).Error
	return &payment, err
}

func (r *repository) Update(payment *Payment) error {
	return r.db.Save(payment).Error
}

func (r *repository) Delete(tenantID, paymentID uuid.UUID) error {
	return r.db.Where("tenant_id = ? AND id = ?", tenantID, paymentID).Delete(&Payment{}).Error
}

func (r *repository) List(tenantID uuid.UUID, orderID *uuid.UUID, offset, limit int) ([]*Payment, int64, error) {
	var payments []*Payment
	var total int64

	query := r.db.Model(&Payment{}).Where("tenant_id = ?", tenantID)
	
	if orderID != nil {
		query = query.Where("order_id = ?", *orderID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// Refund operations
func (r *repository) CreateRefund(refund *Refund) error {
	return r.db.Create(refund).Error
}

func (r *repository) GetRefund(tenantID, refundID uuid.UUID) (*Refund, error) {
	var refund Refund
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, refundID).First(&refund).Error
	return &refund, err
}

func (r *repository) ListRefunds(tenantID uuid.UUID, paymentID *uuid.UUID, offset, limit int) ([]*Refund, int64, error) {
	var refunds []*Refund
	var total int64

	query := r.db.Model(&Refund{}).Where("tenant_id = ?", tenantID)
	
	if paymentID != nil {
		query = query.Where("payment_id = ?", *paymentID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&refunds).Error; err != nil {
		return nil, 0, err
	}

	return refunds, total, nil
}

// Payment method operations
func (r *repository) CreatePaymentMethod(method *PaymentMethod) error {
	return r.db.Create(method).Error
}

func (r *repository) GetPaymentMethod(tenantID, methodID uuid.UUID) (*PaymentMethod, error) {
	var method PaymentMethod
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, methodID).First(&method).Error
	return &method, err
}

func (r *repository) ListPaymentMethods(tenantID, userID uuid.UUID) ([]*PaymentMethod, error) {
	var methods []*PaymentMethod
	err := r.db.Where("tenant_id = ? AND user_id = ? AND is_active = ?", tenantID, userID, true).
		Order("is_default DESC, created_at DESC").Find(&methods).Error
	return methods, err
}

func (r *repository) UpdatePaymentMethod(method *PaymentMethod) error {
	return r.db.Save(method).Error
}

func (r *repository) DeletePaymentMethod(tenantID, methodID uuid.UUID) error {
	return r.db.Where("tenant_id = ? AND id = ?", tenantID, methodID).Delete(&PaymentMethod{}).Error
}
