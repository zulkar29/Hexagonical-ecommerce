-- Migration for Payment and Notification modules
-- Run this SQL script to create the required tables

-- ================================
-- PAYMENT MODULE TABLES
-- ================================

-- Payments table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    order_id UUID,
    user_id UUID,
    
    -- Payment identification
    payment_number VARCHAR(100) UNIQUE NOT NULL,
    gateway VARCHAR(50) NOT NULL, -- 'sslcommerz', 'bkash', 'nagad', etc.
    gateway_payment_id VARCHAR(255),
    
    -- Amount and currency
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    exchange_rate DECIMAL(10,4) DEFAULT 1.0000,
    
    -- Status and metadata
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed, cancelled, refunded
    payment_method VARCHAR(50), -- card, mobile_banking, net_banking, etc.
    
    -- Payment details
    description TEXT,
    metadata JSONB,
    
    -- Gateway specific data
    gateway_response JSONB,
    gateway_fee DECIMAL(10,2) DEFAULT 0.00,
    
    -- Refund information
    refunded_amount DECIMAL(10,2) DEFAULT 0.00,
    refund_reason TEXT,
    
    -- Timestamps
    processed_at TIMESTAMP,
    failed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    refunded_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT payments_amount_positive CHECK (amount > 0),
    CONSTRAINT payments_refunded_amount_positive CHECK (refunded_amount >= 0),
    CONSTRAINT payments_refunded_amount_not_exceed CHECK (refunded_amount <= amount)
);

-- Payment transactions table (for tracking all payment events)
CREATE TABLE IF NOT EXISTS payment_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    
    -- Transaction details
    transaction_type VARCHAR(30) NOT NULL, -- payment, refund, partial_refund, void
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    
    -- Gateway response
    gateway_transaction_id VARCHAR(255),
    gateway_response JSONB,
    status VARCHAR(20) NOT NULL, -- success, failed, pending
    
    -- Additional info
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Refunds table
CREATE TABLE IF NOT EXISTS refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    
    -- Refund details
    refund_number VARCHAR(100) UNIQUE NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    reason TEXT,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed
    
    -- Gateway info
    gateway_refund_id VARCHAR(255),
    gateway_response JSONB,
    
    -- Timestamps
    processed_at TIMESTAMP,
    failed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT refunds_amount_positive CHECK (amount > 0)
);

-- ================================
-- NOTIFICATION MODULE TABLES
-- ================================

-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID,
    
    -- Notification content
    type VARCHAR(50) NOT NULL, -- order_confirmation, payment_success, shipping_update, etc.
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    
    -- Delivery channels
    channels TEXT[] NOT NULL, -- [email, sms, push, in_app]
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, sent, delivered, failed, read
    priority INTEGER NOT NULL DEFAULT 5, -- 1-10, 1 highest priority
    
    -- Scheduling
    scheduled_at TIMESTAMP,
    sent_at TIMESTAMP,
    delivered_at TIMESTAMP,
    read_at TIMESTAMP,
    
    -- Additional data
    metadata JSONB,
    template_id UUID,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT notifications_priority_range CHECK (priority >= 1 AND priority <= 10)
);

-- Notification deliveries table (tracks delivery per channel)
CREATE TABLE IF NOT EXISTS notification_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    
    -- Delivery details
    channel VARCHAR(20) NOT NULL, -- email, sms, push, in_app
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, sent, delivered, failed, bounced
    
    -- Delivery info
    recipient VARCHAR(255) NOT NULL, -- email address, phone number, device token, etc.
    provider VARCHAR(50), -- sendgrid, mailgun, twilio, firebase, etc.
    provider_id VARCHAR(255), -- external provider message ID
    
    -- Error handling
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    
    -- Timestamps
    sent_at TIMESTAMP,
    delivered_at TIMESTAMP,
    failed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Notification templates table
CREATE TABLE IF NOT EXISTS notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    
    -- Template identification
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- order_confirmation, payment_success, etc.
    channel VARCHAR(20) NOT NULL, -- email, sms, push, in_app
    
    -- Template content
    subject VARCHAR(255), -- for email templates
    body TEXT NOT NULL,
    variables TEXT[], -- Available template variables
    
    -- Template settings
    is_active BOOLEAN NOT NULL DEFAULT true,
    language VARCHAR(5) NOT NULL DEFAULT 'en', -- en, bn
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    UNIQUE(tenant_id, name, type, channel)
);

-- Notification preferences table (user preferences for receiving notifications)
CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    
    -- Preference settings
    notification_type VARCHAR(50) NOT NULL,
    channels TEXT[] NOT NULL DEFAULT '{}', -- enabled channels for this notification type
    enabled BOOLEAN NOT NULL DEFAULT true,
    
    -- Timing preferences
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    timezone VARCHAR(50) DEFAULT 'Asia/Dhaka',
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    UNIQUE(tenant_id, user_id, notification_type)
);

-- ================================
-- INDEXES FOR BETTER PERFORMANCE
-- ================================

-- Payment indexes
CREATE INDEX IF NOT EXISTS idx_payments_tenant_id ON payments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at);
CREATE INDEX IF NOT EXISTS idx_payments_gateway ON payments(gateway);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_payment_id ON payment_transactions(payment_id);
CREATE INDEX IF NOT EXISTS idx_refunds_payment_id ON refunds(payment_id);

-- Notification indexes
CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id ON notifications(tenant_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_notification_deliveries_notification_id ON notification_deliveries(notification_id);
CREATE INDEX IF NOT EXISTS idx_notification_deliveries_channel ON notification_deliveries(channel);
CREATE INDEX IF NOT EXISTS idx_notification_deliveries_status ON notification_deliveries(status);
CREATE INDEX IF NOT EXISTS idx_notification_templates_tenant_type ON notification_templates(tenant_id, type);
CREATE INDEX IF NOT EXISTS idx_notification_preferences_tenant_user ON notification_preferences(tenant_id, user_id);

-- ================================
-- UPDATE TRIGGERS
-- ================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for payments
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Triggers for refunds  
CREATE TRIGGER update_refunds_updated_at BEFORE UPDATE ON refunds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Triggers for notifications
CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Triggers for notification_deliveries
CREATE TRIGGER update_notification_deliveries_updated_at BEFORE UPDATE ON notification_deliveries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Triggers for notification_templates
CREATE TRIGGER update_notification_templates_updated_at BEFORE UPDATE ON notification_templates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Triggers for notification_preferences
CREATE TRIGGER update_notification_preferences_updated_at BEFORE UPDATE ON notification_preferences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================
-- SEED DATA
-- ================================

-- Insert default notification templates for Bangladesh e-commerce
INSERT INTO notification_templates (tenant_id, name, type, channel, subject, body, variables, language) 
VALUES 
-- Order confirmation templates (English)
(gen_random_uuid(), 'Order Confirmation Email', 'order_confirmation', 'email', 
 'Order Confirmation - {{order_number}}', 
 'Dear {{customer_name}},\n\nYour order {{order_number}} has been confirmed.\nTotal: {{total_amount}} BDT\n\nItems:\n{{order_items}}\n\nThank you for shopping with us!',
 ARRAY['customer_name', 'order_number', 'total_amount', 'order_items'],
 'en'),

(gen_random_uuid(), 'Order Confirmation SMS', 'order_confirmation', 'sms',
 NULL,
 'Order {{order_number}} confirmed. Total: {{total_amount}} BDT. Track: {{tracking_url}}',
 ARRAY['order_number', 'total_amount', 'tracking_url'],
 'en'),

-- Payment success templates (English)
(gen_random_uuid(), 'Payment Success Email', 'payment_success', 'email',
 'Payment Successful - {{order_number}}',
 'Dear {{customer_name}},\n\nYour payment of {{amount}} BDT for order {{order_number}} has been processed successfully.\n\nPayment Method: {{payment_method}}\nTransaction ID: {{transaction_id}}\n\nYour order will be processed shortly.',
 ARRAY['customer_name', 'order_number', 'amount', 'payment_method', 'transaction_id'],
 'en'),

(gen_random_uuid(), 'Payment Success SMS', 'payment_success', 'sms',
 NULL,
 'Payment {{amount}} BDT successful for order {{order_number}}. Transaction: {{transaction_id}}',
 ARRAY['order_number', 'amount', 'transaction_id'],
 'en'),

-- Bengali templates
(gen_random_uuid(), 'Order Confirmation Email Bengali', 'order_confirmation', 'email',
 'অর্ডার নিশ্চিতকরণ - {{order_number}}',
 'প্রিয় {{customer_name}},\n\nআপনার অর্ডার {{order_number}} নিশ্চিত হয়েছে।\nমোট: {{total_amount}} টাকা\n\nআইটেমসমূহ:\n{{order_items}}\n\nআমাদের সাথে কেনাকাটার জন্য ধন্যবাদ!',
 ARRAY['customer_name', 'order_number', 'total_amount', 'order_items'],
 'bn'),

(gen_random_uuid(), 'Payment Success SMS Bengali', 'payment_success', 'sms',
 NULL,
 'অর্ডার {{order_number}} এর {{amount}} টাকা পেমেন্ট সফল। ট্রানজেকশন: {{transaction_id}}',
 ARRAY['order_number', 'amount', 'transaction_id'],
 'bn')

ON CONFLICT (tenant_id, name, type, channel) DO NOTHING;

COMMIT;
