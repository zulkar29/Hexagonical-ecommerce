-- Migration: Create component system tables
-- Created: 2024-09-04

CREATE TABLE IF NOT EXISTS component_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(150) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('header', 'footer', 'banner', 'sidebar', 'content', 'navigation', 'hero', 'testimonial', 'cta', 'gallery')),
    category VARCHAR(50),
    
    -- Template content
    html_template TEXT NOT NULL,
    css_template TEXT,
    js_template TEXT,
    
    -- Configuration
    config_schema JSONB, -- JSON schema for component configuration
    default_config JSONB, -- Default configuration values
    
    -- Metadata
    preview_image VARCHAR(500),
    tags JSONB,
    is_free BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,
    price DECIMAL(10,2) DEFAULT 0.00,
    
    -- SEO
    meta_title VARCHAR(255),
    meta_description VARCHAR(500),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create themes table
CREATE TABLE IF NOT EXISTS themes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(150) NOT NULL,
    description TEXT,
    
    -- Theme configuration
    global_styles JSONB, -- Global CSS variables, fonts, colors
    layout_config JSONB, -- Layout configuration
    
    -- Status
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'inactive', 'archived')),
    is_active BOOLEAN DEFAULT FALSE, -- Only one theme can be active per tenant
    
    -- Metadata
    preview_image VARCHAR(500),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    UNIQUE(tenant_id, slug)
);

-- Create components table
CREATE TABLE IF NOT EXISTS components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    template_id UUID REFERENCES component_templates(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(150) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('header', 'footer', 'banner', 'sidebar', 'content', 'navigation', 'hero', 'testimonial', 'cta', 'gallery')),
    category VARCHAR(50),
    
    -- Component content
    html_content TEXT NOT NULL,
    css_content TEXT,
    js_content TEXT,
    
    -- Configuration
    config JSONB, -- Component-specific configuration
    
    -- Status and visibility
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'inactive', 'archived')),
    is_featured BOOLEAN DEFAULT FALSE,
    
    -- Metadata
    preview_image VARCHAR(500),
    tags JSONB,
    
    -- SEO
    meta_title VARCHAR(255),
    meta_description VARCHAR(500),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    UNIQUE(tenant_id, slug)
);

-- Create component_instances table
CREATE TABLE IF NOT EXISTS component_instances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    theme_id UUID NOT NULL REFERENCES themes(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES components(id) ON DELETE CASCADE,
    
    -- Instance configuration
    instance_config JSONB, -- Instance-specific overrides
    custom_css TEXT, -- Instance-specific CSS
    custom_js TEXT, -- Instance-specific JS
    
    -- Position and layout
    position VARCHAR(20) NOT NULL CHECK (position IN ('header', 'footer', 'sidebar_left', 'sidebar_right', 'content_top', 'content_bottom', 'hero', 'banner')),
    sort_order INTEGER DEFAULT 0,
    
    -- Visibility
    is_visible BOOLEAN DEFAULT TRUE,
    visibility_rules JSONB, -- Rules for when to show/hide (e.g., specific pages, user roles)
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    UNIQUE(theme_id, component_id, position) -- One component per position per theme
);

-- Add foreign key constraints
ALTER TABLE themes ADD CONSTRAINT fk_themes_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;

-- Add unique constraints
ALTER TABLE component_templates ADD CONSTRAINT unique_component_template_slug UNIQUE (tenant_id, slug);

-- Create indexes for performance

-- Component templates indexes
CREATE INDEX IF NOT EXISTS idx_component_templates_tenant_id ON component_templates(tenant_id);
CREATE INDEX IF NOT EXISTS idx_component_templates_type ON component_templates(type);
CREATE INDEX IF NOT EXISTS idx_component_templates_category ON component_templates(category);
CREATE INDEX IF NOT EXISTS idx_component_templates_slug ON component_templates(tenant_id, slug);
CREATE INDEX IF NOT EXISTS idx_component_templates_featured ON component_templates(is_featured);
CREATE INDEX IF NOT EXISTS idx_component_templates_free ON component_templates(is_free);

-- Themes indexes
CREATE INDEX IF NOT EXISTS idx_themes_tenant_id ON themes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_themes_slug ON themes(tenant_id, slug);
CREATE INDEX IF NOT EXISTS idx_themes_status ON themes(status);
CREATE INDEX IF NOT EXISTS idx_themes_active ON themes(tenant_id, is_active);

-- Components indexes
CREATE INDEX IF NOT EXISTS idx_components_tenant_id ON components(tenant_id);
CREATE INDEX IF NOT EXISTS idx_components_template_id ON components(template_id);
CREATE INDEX IF NOT EXISTS idx_components_slug ON components(tenant_id, slug);
CREATE INDEX IF NOT EXISTS idx_components_type ON components(type);
CREATE INDEX IF NOT EXISTS idx_components_category ON components(category);
CREATE INDEX IF NOT EXISTS idx_components_status ON components(status);
CREATE INDEX IF NOT EXISTS idx_components_featured ON components(is_featured);

-- Component instances indexes
CREATE INDEX IF NOT EXISTS idx_component_instances_tenant_id ON component_instances(tenant_id);
CREATE INDEX IF NOT EXISTS idx_component_instances_theme_id ON component_instances(theme_id);
CREATE INDEX IF NOT EXISTS idx_component_instances_component_id ON component_instances(component_id);
CREATE INDEX IF NOT EXISTS idx_component_instances_position ON component_instances(position);
CREATE INDEX IF NOT EXISTS idx_component_instances_sort_order ON component_instances(sort_order);
CREATE INDEX IF NOT EXISTS idx_component_instances_visible ON component_instances(is_visible);

-- Create GIN indexes for JSONB columns
CREATE INDEX IF NOT EXISTS idx_component_templates_config_schema ON component_templates USING GIN (config_schema);
CREATE INDEX IF NOT EXISTS idx_component_templates_default_config ON component_templates USING GIN (default_config);
CREATE INDEX IF NOT EXISTS idx_component_templates_tags ON component_templates USING GIN (tags);

CREATE INDEX IF NOT EXISTS idx_themes_global_styles ON themes USING GIN (global_styles);
CREATE INDEX IF NOT EXISTS idx_themes_layout_config ON themes USING GIN (layout_config);

CREATE INDEX IF NOT EXISTS idx_components_config ON components USING GIN (config);
CREATE INDEX IF NOT EXISTS idx_components_tags ON components USING GIN (tags);

CREATE INDEX IF NOT EXISTS idx_component_instances_config ON component_instances USING GIN (instance_config);
CREATE INDEX IF NOT EXISTS idx_component_instances_visibility_rules ON component_instances USING GIN (visibility_rules);

-- Create triggers for updated_at columns
CREATE TRIGGER update_component_templates_updated_at 
    BEFORE UPDATE ON component_templates 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_themes_updated_at 
    BEFORE UPDATE ON themes 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_components_updated_at 
    BEFORE UPDATE ON components 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_component_instances_updated_at 
    BEFORE UPDATE ON component_instances 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to ensure only one active theme per tenant
CREATE OR REPLACE FUNCTION ensure_single_active_theme()
RETURNS TRIGGER AS $$
BEGIN
    -- If setting a theme to active, deactivate all other themes for this tenant
    IF NEW.is_active = TRUE THEN
        UPDATE themes 
        SET is_active = FALSE 
        WHERE tenant_id = NEW.tenant_id 
        AND id != NEW.id 
        AND is_active = TRUE;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to ensure only one active theme per tenant
CREATE TRIGGER ensure_single_active_theme_trigger
    BEFORE INSERT OR UPDATE ON themes
    FOR EACH ROW
    EXECUTE FUNCTION ensure_single_active_theme();

-- Note: Default component templates should be created through the application API
-- rather than hardcoded in migrations to avoid foreign key constraint violations.
-- This ensures proper tenant validation and data integrity.