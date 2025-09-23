-- Migration for Notification module
-- Run this SQL script to create the required notification tables
-- Note: Payment tables are already created in migration 006

-- ================================
-- NOTIFICATION MODULE TABLES
-- ================================

-- Note: notifications table is already created in migration 007

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
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Note: notification_templates table is already created in migration 007

-- Note: notification_preferences table is already created in migration 007

-- ================================
-- FOREIGN KEY CONSTRAINTS
-- ================================

-- Notification module foreign keys (main notification tables already have constraints from migration 007)
ALTER TABLE notification_deliveries ADD CONSTRAINT fk_notification_deliveries_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;

-- ================================
-- INDEXES FOR BETTER PERFORMANCE
-- ================================

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

-- Function to update updated_at timestamp (if not already exists)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for notification deliveries
CREATE TRIGGER update_notification_deliveries_updated_at BEFORE UPDATE ON notification_deliveries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Note: Triggers for notifications, notification_templates, and notification_preferences tables are already created in migration 007
