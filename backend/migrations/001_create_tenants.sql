-- Create tenants table
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(100) UNIQUE NOT NULL,
    custom_domain VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    plan VARCHAR(20) DEFAULT 'starter' CHECK (plan IN ('starter', 'professional', 'premium', 'enterprise')),
    
    -- Business Information
    description TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    address TEXT,
    logo VARCHAR(500),
    
    -- Settings
    currency VARCHAR(3) DEFAULT 'BDT',
    language VARCHAR(5) DEFAULT 'bn',
    timezone VARCHAR(50) DEFAULT 'Asia/Dhaka',
    
    -- Limits
    product_limit INTEGER DEFAULT 100,
    storage_limit INTEGER DEFAULT 1024, -- in MB
    bandwidth_limit INTEGER DEFAULT 10240, -- in MB
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_tenants_subdomain ON tenants(subdomain);
CREATE INDEX IF NOT EXISTS idx_tenants_custom_domain ON tenants(custom_domain);
CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);
CREATE INDEX IF NOT EXISTS idx_tenants_plan ON tenants(plan);

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tenants_updated_at 
    BEFORE UPDATE ON tenants 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
