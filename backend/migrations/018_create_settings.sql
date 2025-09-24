-- Create settings table for tenant-specific configuration
CREATE TABLE IF NOT EXISTS settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    section VARCHAR(100) NOT NULL,
    key VARCHAR(100) NOT NULL,
    value TEXT,
    type VARCHAR(50) DEFAULT 'string',
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_settings_tenant_id ON settings(tenant_id);
CREATE INDEX IF NOT EXISTS idx_settings_section ON settings(section);
CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(key);
CREATE INDEX IF NOT EXISTS idx_settings_tenant_section ON settings(tenant_id, section);
CREATE INDEX IF NOT EXISTS idx_settings_tenant_section_key ON settings(tenant_id, section, key);

-- Create unique constraint to prevent duplicate settings
CREATE UNIQUE INDEX IF NOT EXISTS idx_settings_unique ON settings(tenant_id, section, key);

-- Add foreign key constraint to tenants table
ALTER TABLE settings ADD CONSTRAINT fk_settings_tenant_id 
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_settings_updated_at
    BEFORE UPDATE ON settings
    FOR EACH ROW
    EXECUTE FUNCTION update_settings_updated_at();