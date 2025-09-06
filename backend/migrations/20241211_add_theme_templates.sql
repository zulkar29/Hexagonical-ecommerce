-- Add theme template tables
CREATE TABLE IF NOT EXISTS theme_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
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
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS theme_template_components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_theme_templates_slug ON theme_templates(slug);
CREATE INDEX IF NOT EXISTS idx_theme_templates_category ON theme_templates(category);
CREATE INDEX IF NOT EXISTS idx_theme_templates_featured ON theme_templates(is_featured);
CREATE INDEX IF NOT EXISTS idx_theme_templates_default ON theme_templates(is_default);

CREATE INDEX IF NOT EXISTS idx_theme_template_components_theme_id ON theme_template_components(theme_template_id);
CREATE INDEX IF NOT EXISTS idx_theme_template_components_template_id ON theme_template_components(component_template_id);
CREATE INDEX IF NOT EXISTS idx_theme_template_components_zone ON theme_template_components(zone);
CREATE INDEX IF NOT EXISTS idx_theme_template_components_position ON theme_template_components(position);

-- Insert some default theme templates
INSERT INTO theme_templates (name, slug, description, category, is_default, is_featured) VALUES
('Modern E-commerce', 'modern-ecommerce', 'A clean and modern theme perfect for online stores', 'ecommerce', true, true),
('Minimalist Blog', 'minimalist-blog', 'Simple and elegant theme for blogs and content sites', 'blog', false, true),
('Corporate Business', 'corporate-business', 'Professional theme for business and corporate websites', 'business', false, false),
('Creative Portfolio', 'creative-portfolio', 'Showcase your work with this creative portfolio theme', 'portfolio', false, true);

-- Grant permissions to roles
GRANT SELECT, INSERT, UPDATE, DELETE ON theme_templates TO authenticated;
GRANT SELECT, INSERT, UPDATE, DELETE ON theme_template_components TO authenticated;
GRANT SELECT ON theme_templates TO anon;
GRANT SELECT ON theme_template_components TO anon;