-- Migration: Create billing system tables
-- Description: Creates tables for multi-tenant billing with subscription management, usage tracking, invoicing, and dunning

-- Billing Plans Table
CREATE TABLE IF NOT EXISTS billing_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    base_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    billing_cycle VARCHAR(20) NOT NULL DEFAULT 'monthly' CHECK (billing_cycle IN ('monthly', 'quarterly', 'yearly')),
    limits JSONB DEFAULT '{}',
    features JSONB DEFAULT '[]',
    trial_period_days INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Usage Tiers Table (for usage-based pricing)
CREATE TABLE IF NOT EXISTS usage_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    billing_plan_id UUID NOT NULL,
    usage_type VARCHAR(50) NOT NULL CHECK (usage_type IN (
        'api_requests', 'storage_gb', 'bandwidth_gb', 'orders', 
        'products', 'users', 'emails_sent', 'transactions'
    )),
    min_units BIGINT NOT NULL DEFAULT 0,
    max_units BIGINT, -- NULL for unlimited
    price_per_unit DECIMAL(10,6) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    FOREIGN KEY (billing_plan_id) REFERENCES billing_plans(id) ON DELETE CASCADE
);

-- Tenant Subscriptions Table
CREATE TABLE IF NOT EXISTS tenant_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    plan_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN (
        'active', 'pending', 'suspended', 'canceled', 'expired'
    )),
    billing_cycle VARCHAR(20) NOT NULL CHECK (billing_cycle IN ('monthly', 'quarterly', 'yearly')),
    current_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    current_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    trial_end TIMESTAMP WITH TIME ZONE,
    base_amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    
    -- Plan change management
    pending_new_plan_id UUID,
    pending_effective_date TIMESTAMP WITH TIME ZONE,
    pending_proration_amount DECIMAL(10,2),
    pending_change_reason VARCHAR(255),
    
    -- Payment and billing
    payment_method_id VARCHAR(255),
    next_billing_date TIMESTAMP WITH TIME ZONE NOT NULL,
    canceled_at TIMESTAMP WITH TIME ZONE,
    cancellation_reason TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    FOREIGN KEY (plan_id) REFERENCES billing_plans(id),
    FOREIGN KEY (pending_new_plan_id) REFERENCES billing_plans(id),
    UNIQUE(tenant_id) -- One subscription per tenant
);

-- Usage Records Table
CREATE TABLE IF NOT EXISTS usage_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    usage_type VARCHAR(50) NOT NULL CHECK (usage_type IN (
        'api_requests', 'storage_gb', 'bandwidth_gb', 'orders', 
        'products', 'users', 'emails_sent', 'transactions'
    )),
    quantity BIGINT NOT NULL,
    units VARCHAR(50) NOT NULL,
    billing_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    billing_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    resource_id VARCHAR(255), -- ID of resource that generated usage
    metadata JSONB DEFAULT '{}',
    recorded_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Invoices Table
CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    subscription_id UUID NOT NULL,
    invoice_number VARCHAR(100) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN (
        'draft', 'pending', 'paid', 'overdue', 'voided', 'refunded'
    )),
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    subtotal_amount DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) DEFAULT 0.00,
    total_amount DECIMAL(10,2) NOT NULL,
    paid_amount DECIMAL(10,2) DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    paid_at TIMESTAMP WITH TIME ZONE,
    payment_method VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    FOREIGN KEY (subscription_id) REFERENCES tenant_subscriptions(id)
);

-- Invoice Line Items Table
CREATE TABLE IF NOT EXISTS invoice_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    description TEXT NOT NULL,
    quantity BIGINT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('subscription', 'usage', 'addon', 'discount')),
    usage_type VARCHAR(50), -- NULL for non-usage items
    period_start TIMESTAMP WITH TIME ZONE,
    period_end TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE
);

-- Payment Attempts Table
CREATE TABLE IF NOT EXISTS payment_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN (
        'pending', 'success', 'failed', 'retrying', 'abandoned'
    )),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    payment_method VARCHAR(255) NOT NULL,
    provider_id VARCHAR(100), -- e.g., 'stripe', 'bkash'
    provider_charge_id VARCHAR(255),
    provider_response JSONB,
    attempted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    failure_reason TEXT,
    retry_count INTEGER DEFAULT 0,
    next_retry_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    FOREIGN KEY (invoice_id) REFERENCES invoices(id)
);

-- Dunning Processes Table
CREATE TABLE IF NOT EXISTS dunning_processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    invoice_id UUID NOT NULL,
    subscription_id UUID NOT NULL,
    current_step INTEGER DEFAULT 1,
    total_steps INTEGER NOT NULL DEFAULT 5,
    is_completed BOOLEAN DEFAULT false,
    is_abandoned BOOLEAN DEFAULT false,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    next_action_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    emails_sent INTEGER DEFAULT 0,
    payment_attempts INTEGER DEFAULT 0,
    service_suspended BOOLEAN DEFAULT false,
    service_suspended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    FOREIGN KEY (invoice_id) REFERENCES invoices(id),
    FOREIGN KEY (subscription_id) REFERENCES tenant_subscriptions(id),
    UNIQUE(invoice_id) -- One dunning process per invoice
);

-- Dunning Actions Table
CREATE TABLE IF NOT EXISTS dunning_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dunning_process_id UUID NOT NULL,
    action_type VARCHAR(20) NOT NULL CHECK (action_type IN ('email', 'suspend', 'retry_payment', 'cancel')),
    step_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed')),
    result JSONB,
    error_message TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    executed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    FOREIGN KEY (dunning_process_id) REFERENCES dunning_processes(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_billing_plans_active ON billing_plans(is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_billing_plans_public ON billing_plans(is_public) WHERE is_public = true;

CREATE INDEX IF NOT EXISTS idx_usage_tiers_plan_type ON usage_tiers(billing_plan_id, usage_type);

CREATE INDEX IF NOT EXISTS idx_tenant_subscriptions_tenant_id ON tenant_subscriptions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_tenant_subscriptions_plan_id ON tenant_subscriptions(plan_id);
CREATE INDEX IF NOT EXISTS idx_tenant_subscriptions_status ON tenant_subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_tenant_subscriptions_next_billing ON tenant_subscriptions(next_billing_date) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_tenant_subscriptions_pending_changes ON tenant_subscriptions(pending_new_plan_id) WHERE pending_new_plan_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_usage_records_tenant_id ON usage_records(tenant_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_usage_type ON usage_records(usage_type);
CREATE INDEX IF NOT EXISTS idx_usage_records_period ON usage_records(billing_period_start, billing_period_end);
CREATE INDEX IF NOT EXISTS idx_usage_records_recorded_at ON usage_records(recorded_at);

CREATE INDEX IF NOT EXISTS idx_invoices_tenant_id ON invoices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_invoices_subscription_id ON invoices(subscription_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_due_date ON invoices(due_date);
CREATE INDEX IF NOT EXISTS idx_invoices_overdue ON invoices(status, due_date) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_invoices_created_at ON invoices(created_at);

CREATE INDEX IF NOT EXISTS idx_invoice_line_items_invoice_id ON invoice_line_items(invoice_id);
CREATE INDEX IF NOT EXISTS idx_invoice_line_items_type ON invoice_line_items(item_type);

CREATE INDEX IF NOT EXISTS idx_payment_attempts_invoice_id ON payment_attempts(invoice_id);
CREATE INDEX IF NOT EXISTS idx_payment_attempts_tenant_id ON payment_attempts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payment_attempts_status ON payment_attempts(status);
CREATE INDEX IF NOT EXISTS idx_payment_attempts_retry ON payment_attempts(status, next_retry_at) WHERE status = 'failed';

CREATE INDEX IF NOT EXISTS idx_dunning_processes_tenant_id ON dunning_processes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_dunning_processes_invoice_id ON dunning_processes(invoice_id);
CREATE INDEX IF NOT EXISTS idx_dunning_processes_active ON dunning_processes(is_completed, is_abandoned) WHERE is_completed = false AND is_abandoned = false;
CREATE INDEX IF NOT EXISTS idx_dunning_processes_next_action ON dunning_processes(next_action_at) WHERE is_completed = false AND is_abandoned = false;

CREATE INDEX IF NOT EXISTS idx_dunning_actions_process_id ON dunning_actions(dunning_process_id);
CREATE INDEX IF NOT EXISTS idx_dunning_actions_scheduled ON dunning_actions(status, scheduled_at) WHERE status = 'pending';

-- Triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_billing_plans_updated_at BEFORE UPDATE ON billing_plans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_usage_tiers_updated_at BEFORE UPDATE ON usage_tiers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tenant_subscriptions_updated_at BEFORE UPDATE ON tenant_subscriptions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_invoices_updated_at BEFORE UPDATE ON invoices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_invoice_line_items_updated_at BEFORE UPDATE ON invoice_line_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_payment_attempts_updated_at BEFORE UPDATE ON payment_attempts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_dunning_processes_updated_at BEFORE UPDATE ON dunning_processes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_dunning_actions_updated_at BEFORE UPDATE ON dunning_actions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
