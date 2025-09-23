-- Add affiliate marketing features to existing referral system
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Add new fields to existing referrals table for affiliate marketing functionality
ALTER TABLE referrals ADD COLUMN IF NOT EXISTS affiliate_type VARCHAR(20) DEFAULT 'customer' CHECK (affiliate_type IN ('customer', 'influencer', 'partner', 'employee', 'digital_marketer'));
ALTER TABLE referrals ADD COLUMN IF NOT EXISTS tracking_data JSONB DEFAULT '{}';
ALTER TABLE referrals ADD COLUMN IF NOT EXISTS payout_threshold DECIMAL(10,2) DEFAULT 50.0;

-- Update existing referral_commissions table to support order commissions
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS order_id UUID REFERENCES orders(id) ON DELETE SET NULL;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS order_amount DECIMAL(10,2) DEFAULT 0;
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS conversion_type VARCHAR(20) DEFAULT 'order' CHECK (conversion_type IN ('order', 'subscription', 'signup'));
ALTER TABLE referral_commissions ADD COLUMN IF NOT EXISTS click_data JSONB DEFAULT '{}';

-- Note: referral_commissions table already supports order commissions via order_id column

-- Create affiliate_clicks table for detailed click tracking
CREATE TABLE IF NOT EXISTS affiliate_clicks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    referral_id UUID NOT NULL REFERENCES referrals(id) ON DELETE CASCADE,
    referrer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    referrer TEXT,
    utm_source VARCHAR(100),
    utm_medium VARCHAR(100),
    utm_campaign VARCHAR(100),
    utm_term VARCHAR(100),
    utm_content VARCHAR(100),
    device_type VARCHAR(20),
    country VARCHAR(2),
    city VARCHAR(100),
    converted BOOLEAN DEFAULT false,
    converted_at TIMESTAMPTZ,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create affiliate_payout_batches table for batch payouts
CREATE TABLE IF NOT EXISTS affiliate_payout_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    batch_number VARCHAR(50) NOT NULL,
    total_amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    affiliate_count INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    processed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, batch_number)
);

-- Create indexes for affiliate_clicks table
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_tenant_id ON affiliate_clicks(tenant_id);
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_referral_id ON affiliate_clicks(referral_id);
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_referrer_id ON affiliate_clicks(referrer_id);
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_ip_address ON affiliate_clicks(ip_address);
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_converted ON affiliate_clicks(converted);
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_created_at ON affiliate_clicks(created_at);
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_utm_source ON affiliate_clicks(utm_source);
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_utm_campaign ON affiliate_clicks(utm_campaign);

-- Create indexes for affiliate_payout_batches table
CREATE INDEX IF NOT EXISTS idx_affiliate_payout_batches_tenant_id ON affiliate_payout_batches(tenant_id);
CREATE INDEX IF NOT EXISTS idx_affiliate_payout_batches_status ON affiliate_payout_batches(status);
CREATE INDEX IF NOT EXISTS idx_affiliate_payout_batches_created_at ON affiliate_payout_batches(created_at);

-- Create indexes for referrals table new fields
CREATE INDEX IF NOT EXISTS idx_referrals_affiliate_type ON referrals(affiliate_type);
CREATE INDEX IF NOT EXISTS idx_referrals_payout_threshold ON referrals(payout_threshold);

-- Create indexes for referral_commissions table new fields
CREATE INDEX IF NOT EXISTS idx_referral_commissions_conversion_type ON referral_commissions(conversion_type);

-- Create updated_at trigger for affiliate_payout_batches
CREATE TRIGGER update_affiliate_payout_batches_updated_at
    BEFORE UPDATE ON affiliate_payout_batches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Update the referral stats function to include affiliate click tracking
CREATE OR REPLACE FUNCTION update_affiliate_stats()
RETURNS TRIGGER AS $$
BEGIN
    -- Update stats when clicks are tracked
    IF TG_TABLE_NAME = 'affiliate_clicks' THEN
        INSERT INTO referral_stats (tenant_id, user_id, total_referrals, last_referral_at)
        VALUES (NEW.tenant_id, NEW.referrer_id, 0, NEW.created_at)
        ON CONFLICT (tenant_id, user_id) DO UPDATE SET
            last_referral_at = GREATEST(referral_stats.last_referral_at, NEW.created_at),
            updated_at = NOW();

        -- Update conversion if click converted
        IF NEW.converted = true AND (OLD.converted IS NULL OR OLD.converted = false) THEN
            UPDATE referral_stats SET
                successful_referrals = successful_referrals + 1,
                updated_at = NOW()
            WHERE tenant_id = NEW.tenant_id AND user_id = NEW.referrer_id;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for affiliate click stats
CREATE TRIGGER trigger_update_affiliate_stats_on_click
    AFTER INSERT OR UPDATE ON affiliate_clicks
    FOR EACH ROW
    EXECUTE FUNCTION update_affiliate_stats();

-- Update existing referrals to have proper affiliate type
UPDATE referrals
SET affiliate_type = 'customer',
    tracking_data = '{}',
    payout_threshold = 50.0,
    updated_at = NOW()
WHERE affiliate_type IS NULL;


-- Add comments for documentation
COMMENT ON TABLE affiliate_clicks IS 'Tracks clicks on affiliate links for detailed analytics';
COMMENT ON TABLE affiliate_payout_batches IS 'Manages batch payouts to affiliates';
COMMENT ON COLUMN referrals.affiliate_type IS 'Type of affiliate: customer, influencer, partner, employee, digital_marketer';
COMMENT ON COLUMN referrals.tracking_data IS 'JSON data for tracking affiliate-specific information';
COMMENT ON COLUMN referrals.payout_threshold IS 'Minimum amount required before payout';
COMMENT ON COLUMN referral_commissions.conversion_type IS 'Type of conversion: order, subscription, signup';
COMMENT ON COLUMN affiliate_clicks.utm_source IS 'UTM source parameter for tracking campaign source';
COMMENT ON COLUMN affiliate_clicks.utm_campaign IS 'UTM campaign parameter for tracking specific campaigns';