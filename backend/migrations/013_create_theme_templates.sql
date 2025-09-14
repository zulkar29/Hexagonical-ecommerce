-- Migration: Create theme template tables
-- Created: 2024-09-04

CREATE TABLE IF NOT EXISTS theme_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Theme template content
    global_css TEXT,
    global_js TEXT,
    settings JSONB DEFAULT '{}',
    
    -- Template metadata
    thumbnail VARCHAR(500),
    preview_url VARCHAR(500),
    tags JSONB DEFAULT '[]',
    category VARCHAR(100),
    is_free BOOLEAN DEFAULT true,
    is_featured BOOLEAN DEFAULT false,
    is_default BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS theme_template_components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    theme_template_id UUID NOT NULL REFERENCES theme_templates(id) ON DELETE CASCADE,
    component_template_id UUID NOT NULL REFERENCES component_templates(id) ON DELETE CASCADE,
    
    -- Component positioning and configuration
    zone VARCHAR(50) DEFAULT 'main', -- header, footer, main, sidebar
    position INTEGER DEFAULT 0,
    settings JSONB DEFAULT '{}',
    custom_css TEXT,
    custom_js TEXT,
    is_visible BOOLEAN DEFAULT true,
    
    -- Responsive settings
    breakpoints JSONB DEFAULT '{}',
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Foreign key constraints
ALTER TABLE theme_templates ADD CONSTRAINT fk_theme_templates_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
ALTER TABLE theme_template_components ADD CONSTRAINT fk_theme_template_components_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;

-- Unique constraints
ALTER TABLE theme_templates ADD CONSTRAINT unique_theme_templates_tenant_slug UNIQUE (tenant_id, slug);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_theme_templates_tenant_id ON theme_templates(tenant_id);
CREATE INDEX IF NOT EXISTS idx_theme_templates_slug ON theme_templates(slug);
CREATE INDEX IF NOT EXISTS idx_theme_templates_category ON theme_templates(category);
CREATE INDEX IF NOT EXISTS idx_theme_templates_featured ON theme_templates(is_featured);
CREATE INDEX IF NOT EXISTS idx_theme_templates_default ON theme_templates(is_default);
CREATE INDEX IF NOT EXISTS idx_theme_templates_created_at ON theme_templates(created_at);

CREATE INDEX IF NOT EXISTS idx_theme_template_components_tenant_id ON theme_template_components(tenant_id);
CREATE INDEX IF NOT EXISTS idx_theme_template_components_theme_id ON theme_template_components(theme_template_id);
CREATE INDEX IF NOT EXISTS idx_theme_template_components_template_id ON theme_template_components(component_template_id);
CREATE INDEX IF NOT EXISTS idx_theme_template_components_zone ON theme_template_components(zone);
CREATE INDEX IF NOT EXISTS idx_theme_template_components_position ON theme_template_components(position);

-- Note: Default theme templates should be inserted per tenant during tenant creation
-- This ensures proper tenant isolation

-- Note: Role-based permissions should be configured at the application level
-- or through database-specific role management