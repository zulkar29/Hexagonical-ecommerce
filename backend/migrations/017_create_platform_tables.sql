-- Migration: Create platform management tables
-- Created: 2024-09-23

-- Create platform_stats table
CREATE TABLE IF NOT EXISTS platform_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    total_tenants INTEGER NOT NULL DEFAULT 0,
    active_tenants INTEGER NOT NULL DEFAULT 0,
    total_revenue DECIMAL(15,2) NOT NULL DEFAULT 0,
    monthly_revenue DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_users INTEGER NOT NULL DEFAULT 0,
    active_users INTEGER NOT NULL DEFAULT 0,
    growth_rate DECIMAL(5,2) NOT NULL DEFAULT 0,
    system_uptime DECIMAL(5,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create platform_admins table
CREATE TABLE IF NOT EXISTS platform_admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('super_admin', 'platform_admin', 'support')),
    permissions JSONB NOT NULL DEFAULT '[]',
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create platform_roles table
CREATE TABLE IF NOT EXISTS platform_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    permissions JSONB NOT NULL DEFAULT '[]',
    is_system BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create platform_settings table
CREATE TABLE IF NOT EXISTS platform_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(100) NOT NULL,
    key VARCHAR(255) NOT NULL,
    value JSONB NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(category, key)
);

-- Create platform_audit_logs table
CREATE TABLE IF NOT EXISTS platform_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    tenant_id UUID,
    action VARCHAR(255) NOT NULL,
    resource VARCHAR(255),
    resource_id UUID,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create system_status table
CREATE TABLE IF NOT EXISTS system_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status VARCHAR(20) NOT NULL CHECK (status IN ('healthy', 'degraded', 'down')),
    database_status VARCHAR(20) NOT NULL,
    redis_status VARCHAR(20) NOT NULL,
    api_response_time DECIMAL(10,2) NOT NULL DEFAULT 0,
    cpu_usage DECIMAL(5,2) NOT NULL DEFAULT 0,
    memory_usage DECIMAL(5,2) NOT NULL DEFAULT 0,
    disk_usage DECIMAL(5,2) NOT NULL DEFAULT 0,
    active_connections INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_platform_admins_email ON platform_admins(email);
CREATE INDEX IF NOT EXISTS idx_platform_admins_role ON platform_admins(role);
CREATE INDEX IF NOT EXISTS idx_platform_admins_status ON platform_admins(status);

CREATE INDEX IF NOT EXISTS idx_platform_roles_name ON platform_roles(name);
CREATE INDEX IF NOT EXISTS idx_platform_roles_is_system ON platform_roles(is_system);

CREATE INDEX IF NOT EXISTS idx_platform_settings_category ON platform_settings(category);
CREATE INDEX IF NOT EXISTS idx_platform_settings_key ON platform_settings(key);

CREATE INDEX IF NOT EXISTS idx_platform_audit_logs_user_id ON platform_audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_platform_audit_logs_tenant_id ON platform_audit_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_platform_audit_logs_action ON platform_audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_platform_audit_logs_created_at ON platform_audit_logs(created_at);

CREATE INDEX IF NOT EXISTS idx_system_status_created_at ON system_status(created_at);

-- Add updated_at trigger for platform tables
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers
DROP TRIGGER IF EXISTS update_platform_stats_updated_at ON platform_stats;
CREATE TRIGGER update_platform_stats_updated_at
    BEFORE UPDATE ON platform_stats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_platform_admins_updated_at ON platform_admins;
CREATE TRIGGER update_platform_admins_updated_at
    BEFORE UPDATE ON platform_admins
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_platform_roles_updated_at ON platform_roles;
CREATE TRIGGER update_platform_roles_updated_at
    BEFORE UPDATE ON platform_roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_platform_settings_updated_at ON platform_settings;
CREATE TRIGGER update_platform_settings_updated_at
    BEFORE UPDATE ON platform_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();