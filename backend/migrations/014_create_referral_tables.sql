-- Create referral tables
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create referrals table
CREATE TABLE IF NOT EXISTS referrals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    referrer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    referral_code VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'expired', 'cancelled')),
    commission_rate DECIMAL(5,4) DEFAULT 0.10, -- 10% default commission
    expires_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create referral_commissions table
CREATE TABLE IF NOT EXISTS referral_commissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    referral_id UUID NOT NULL REFERENCES referrals(id) ON DELETE CASCADE,
    referrer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    commission_amount DECIMAL(10,2) NOT NULL,
    commission_rate DECIMAL(5,4) NOT NULL,
    order_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'cancelled')),
    paid_at TIMESTAMPTZ,
    payment_method VARCHAR(50),
    payment_reference VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create referral_stats table for caching statistics
CREATE TABLE IF NOT EXISTS referral_stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_referrals INTEGER DEFAULT 0,
    successful_referrals INTEGER DEFAULT 0,
    pending_referrals INTEGER DEFAULT 0,
    total_commission_earned DECIMAL(10,2) DEFAULT 0.00,
    total_commission_paid DECIMAL(10,2) DEFAULT 0.00,
    total_commission_pending DECIMAL(10,2) DEFAULT 0.00,
    last_referral_at TIMESTAMPTZ,
    last_commission_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, user_id)
);

-- Add referral_code to tenants table if not exists
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS referral_code VARCHAR(50) UNIQUE;

-- Create indexes for referrals table
CREATE INDEX IF NOT EXISTS idx_referrals_tenant_id ON referrals(tenant_id);
CREATE INDEX IF NOT EXISTS idx_referrals_referrer_id ON referrals(referrer_id);
CREATE INDEX IF NOT EXISTS idx_referrals_referee_id ON referrals(referee_id);
CREATE INDEX IF NOT EXISTS idx_referrals_code ON referrals(referral_code);
CREATE INDEX IF NOT EXISTS idx_referrals_status ON referrals(status);
CREATE INDEX IF NOT EXISTS idx_referrals_expires_at ON referrals(expires_at);
CREATE INDEX IF NOT EXISTS idx_referrals_created_at ON referrals(created_at);

-- Create indexes for referral_commissions table
CREATE INDEX IF NOT EXISTS idx_referral_commissions_tenant_id ON referral_commissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_referral_commissions_referral_id ON referral_commissions(referral_id);
CREATE INDEX IF NOT EXISTS idx_referral_commissions_referrer_id ON referral_commissions(referrer_id);
CREATE INDEX IF NOT EXISTS idx_referral_commissions_referee_id ON referral_commissions(referee_id);
CREATE INDEX IF NOT EXISTS idx_referral_commissions_order_id ON referral_commissions(order_id);
CREATE INDEX IF NOT EXISTS idx_referral_commissions_status ON referral_commissions(status);
CREATE INDEX IF NOT EXISTS idx_referral_commissions_created_at ON referral_commissions(created_at);

-- Create indexes for referral_stats table
CREATE INDEX IF NOT EXISTS idx_referral_stats_tenant_id ON referral_stats(tenant_id);
CREATE INDEX IF NOT EXISTS idx_referral_stats_user_id ON referral_stats(user_id);
CREATE INDEX IF NOT EXISTS idx_referral_stats_tenant_user ON referral_stats(tenant_id, user_id);

-- Create updated_at triggers
CREATE TRIGGER update_referrals_updated_at 
    BEFORE UPDATE ON referrals 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_referral_commissions_updated_at 
    BEFORE UPDATE ON referral_commissions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_referral_stats_updated_at 
    BEFORE UPDATE ON referral_stats 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to update referral stats
CREATE OR REPLACE FUNCTION update_referral_stats()
RETURNS TRIGGER AS $$
BEGIN
    -- Update stats when referral status changes
    IF TG_TABLE_NAME = 'referrals' THEN
        INSERT INTO referral_stats (tenant_id, user_id, total_referrals, successful_referrals, pending_referrals, last_referral_at)
        VALUES (NEW.tenant_id, NEW.referrer_id, 1, 
                CASE WHEN NEW.status = 'completed' THEN 1 ELSE 0 END,
                CASE WHEN NEW.status = 'pending' THEN 1 ELSE 0 END,
                NEW.created_at)
        ON CONFLICT (tenant_id, user_id) DO UPDATE SET
            total_referrals = referral_stats.total_referrals + 
                CASE WHEN TG_OP = 'INSERT' THEN 1 ELSE 0 END,
            successful_referrals = referral_stats.successful_referrals + 
                CASE WHEN NEW.status = 'completed' AND (OLD.status IS NULL OR OLD.status != 'completed') THEN 1
                     WHEN OLD.status = 'completed' AND NEW.status != 'completed' THEN -1
                     ELSE 0 END,
            pending_referrals = referral_stats.pending_referrals + 
                CASE WHEN NEW.status = 'pending' AND (OLD.status IS NULL OR OLD.status != 'pending') THEN 1
                     WHEN OLD.status = 'pending' AND NEW.status != 'pending' THEN -1
                     ELSE 0 END,
            last_referral_at = GREATEST(referral_stats.last_referral_at, NEW.created_at),
            updated_at = NOW();
    END IF;
    
    -- Update stats when commission changes
    IF TG_TABLE_NAME = 'referral_commissions' THEN
        INSERT INTO referral_stats (tenant_id, user_id, total_commission_earned, total_commission_paid, total_commission_pending, last_commission_at)
        VALUES (NEW.tenant_id, NEW.referrer_id, NEW.commission_amount,
                CASE WHEN NEW.status = 'paid' THEN NEW.commission_amount ELSE 0 END,
                CASE WHEN NEW.status = 'pending' THEN NEW.commission_amount ELSE 0 END,
                NEW.created_at)
        ON CONFLICT (tenant_id, user_id) DO UPDATE SET
            total_commission_earned = referral_stats.total_commission_earned + 
                CASE WHEN TG_OP = 'INSERT' THEN NEW.commission_amount ELSE NEW.commission_amount - OLD.commission_amount END,
            total_commission_paid = referral_stats.total_commission_paid + 
                CASE WHEN NEW.status = 'paid' AND (OLD.status IS NULL OR OLD.status != 'paid') THEN NEW.commission_amount
                     WHEN OLD.status = 'paid' AND NEW.status != 'paid' THEN -OLD.commission_amount
                     WHEN NEW.status = 'paid' AND OLD.status = 'paid' THEN NEW.commission_amount - OLD.commission_amount
                     ELSE 0 END,
            total_commission_pending = referral_stats.total_commission_pending + 
                CASE WHEN NEW.status = 'pending' AND (OLD.status IS NULL OR OLD.status != 'pending') THEN NEW.commission_amount
                     WHEN OLD.status = 'pending' AND NEW.status != 'pending' THEN -OLD.commission_amount
                     WHEN NEW.status = 'pending' AND OLD.status = 'pending' THEN NEW.commission_amount - OLD.commission_amount
                     ELSE 0 END,
            last_commission_at = GREATEST(referral_stats.last_commission_at, NEW.created_at),
            updated_at = NOW();
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for automatic stats updates
CREATE TRIGGER trigger_update_referral_stats_on_referral
    AFTER INSERT OR UPDATE ON referrals
    FOR EACH ROW
    EXECUTE FUNCTION update_referral_stats();

CREATE TRIGGER trigger_update_referral_stats_on_commission
    AFTER INSERT OR UPDATE ON referral_commissions
    FOR EACH ROW
    EXECUTE FUNCTION update_referral_stats();