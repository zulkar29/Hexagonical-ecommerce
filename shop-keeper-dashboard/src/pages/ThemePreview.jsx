import React from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import { themeTemplates, sampleLayouts, sampleProducts } from '../data/componentTemplates';
import { useEffect, useState } from 'react';

const ThemePreview = () => {
  const { themeId } = useParams();
  const [searchParams] = useSearchParams();
  const [theme, setTheme] = useState(null);
  const [layout, setLayout] = useState([]);
  const [products, setProducts] = useState([]);

  useEffect(() => {
    if (themeId && themeTemplates[themeId]) {
      const selectedTheme = themeTemplates[themeId];
      setTheme(selectedTheme);
      setLayout(sampleLayouts[themeId] || []);
      setProducts(sampleProducts[themeId] || []);
    }
  }, [themeId]);

  if (!theme) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Theme Not Found</h1>
          <p className="text-gray-600">The requested theme could not be loaded.</p>
        </div>
      </div>
    );
  }

  const renderComponent = (component) => {
    const { type, content } = component;
    
    switch (type) {
      case 'header':
        return (
          <header 
            className="w-full py-4 px-6 border-b"
            style={{
              backgroundColor: theme.colors.background,
              borderColor: theme.colors.border,
              fontFamily: theme.typography.fontFamily
            }}
          >
            <div className="max-w-7xl mx-auto flex items-center justify-between">
              <h1 
                className="text-2xl font-bold"
                style={{
                  color: theme.colors.text,
                  fontFamily: theme.typography.headingFont
                }}
              >
                {content.title}
              </h1>
              <nav className="hidden md:flex space-x-6">
                {content.navigation?.map((item, index) => (
                  <a 
                    key={index}
                    href="#" 
                    className="hover:opacity-75 transition-opacity"
                    style={{ color: theme.colors.textSecondary }}
                  >
                    {item}
                  </a>
                ))}
              </nav>
            </div>
          </header>
        );

      case 'hero':
        return (
          <section 
            className="relative py-20 px-6"
            style={{
              backgroundImage: content.backgroundImage ? `url(${content.backgroundImage})` : 'none',
              backgroundColor: theme.colors.surface,
              backgroundSize: 'cover',
              backgroundPosition: 'center'
            }}
          >
            <div className="absolute inset-0 bg-black bg-opacity-30"></div>
            <div className="relative max-w-4xl mx-auto text-center">
              <h1 
                className="text-5xl font-bold mb-6"
                style={{
                  color: content.backgroundImage ? '#ffffff' : theme.colors.text,
                  fontFamily: theme.typography.headingFont
                }}
              >
                {content.title}
              </h1>
              <p 
                className="text-xl mb-8"
                style={{
                  color: content.backgroundImage ? '#f3f4f6' : theme.colors.textSecondary,
                  fontFamily: theme.typography.fontFamily
                }}
              >
                {content.subtitle}
              </p>
              <button 
                className="px-8 py-3 font-semibold transition-all duration-200 hover:opacity-90"
                style={{
                  backgroundColor: theme.colors.primary,
                  color: '#ffffff',
                  borderRadius: theme.borderRadius.md,
                  fontFamily: theme.typography.fontFamily
                }}
              >
                {content.buttonText}
              </button>
            </div>
          </section>
        );

      case 'products':
        return (
          <section 
            className="py-16 px-6"
            style={{
              backgroundColor: theme.colors.background,
              fontFamily: theme.typography.fontFamily
            }}
          >
            <div className="max-w-6xl mx-auto">
              <h2 
                className="text-3xl font-bold text-center mb-12"
                style={{
                  color: theme.colors.text,
                  fontFamily: theme.typography.headingFont
                }}
              >
                Featured Products
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                {products.slice(0, 6).map((product, index) => (
                  <div 
                    key={product.id || index}
                    className="group cursor-pointer"
                    style={{
                      backgroundColor: theme.colors.surface,
                      borderRadius: theme.borderRadius.lg,
                      padding: theme.spacing.md
                    }}
                  >
                    <div className="aspect-square mb-4 overflow-hidden" style={{ borderRadius: theme.borderRadius.md }}>
                      <img 
                        src={product.image} 
                        alt={product.name}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-200"
                      />
                    </div>
                    <h3 
                      className="font-semibold mb-2"
                      style={{ color: theme.colors.text }}
                    >
                      {product.name}
                    </h3>
                    <p 
                      className="text-lg font-bold"
                      style={{ color: theme.colors.primary }}
                    >
                      {product.price}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          </section>
        );

      case 'footer':
        return (
          <footer 
            className="py-12 px-6 border-t"
            style={{
              backgroundColor: theme.colors.surface,
              borderColor: theme.colors.border,
              fontFamily: theme.typography.fontFamily
            }}
          >
            <div className="max-w-6xl mx-auto">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                <div>
                  <h3 
                    className="text-xl font-bold mb-4"
                    style={{
                      color: theme.colors.text,
                      fontFamily: theme.typography.headingFont
                    }}
                  >
                    {content.companyName}
                  </h3>
                  <p style={{ color: theme.colors.textSecondary }}>
                    Creating exceptional experiences for our customers.
                  </p>
                </div>
                <div>
                  <h4 
                    className="font-semibold mb-4"
                    style={{ color: theme.colors.text }}
                  >
                    Quick Links
                  </h4>
                  <ul className="space-y-2">
                    {content.links?.map((link, index) => (
                      <li key={index}>
                        <a 
                          href="#" 
                          className="hover:opacity-75 transition-opacity"
                          style={{ color: theme.colors.textSecondary }}
                        >
                          {link}
                        </a>
                      </li>
                    ))}
                  </ul>
                </div>
                <div>
                  <h4 
                    className="font-semibold mb-4"
                    style={{ color: theme.colors.text }}
                  >
                    Follow Us
                  </h4>
                  <div className="flex space-x-4">
                    {content.social?.map((platform, index) => (
                      <a 
                        key={index}
                        href="#" 
                        className="hover:opacity-75 transition-opacity"
                        style={{ color: theme.colors.primary }}
                      >
                        {platform}
                      </a>
                    ))}
                  </div>
                </div>
              </div>
              <div 
                className="mt-8 pt-8 border-t text-center"
                style={{ borderColor: theme.colors.border }}
              >
                <p style={{ color: theme.colors.textSecondary }}>
                  © 2024 {content.companyName}. All rights reserved.
                </p>
              </div>
            </div>
          </footer>
        );

      default:
        return null;
    }
  };

  return (
    <div className="min-h-screen" style={{
      backgroundColor: theme.colors.background,
      color: theme.colors.text,
      fontFamily: theme.typography.fontFamily,
      fontSize: theme.typography.fontSize,
      lineHeight: theme.typography.lineHeight
    }}>
      {/* Theme Preview Content */}
      <div className="theme-preview" style={{
        '--primary-color': theme.colors.primary,
        '--secondary-color': theme.colors.secondary,
        '--accent-color': theme.colors.accent,
        '--background-color': theme.colors.background,
        '--text-color': theme.colors.text,
        '--border-radius': theme.borderRadius
      }}>
      {/* Theme Preview Header */}
      <div 
        className="sticky top-0 z-50 px-6 py-3 border-b"
        style={{
          backgroundColor: theme.colors.surface,
          borderColor: theme.colors.border
        }}
      >
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <button 
              onClick={() => window.close()}
              className="px-4 py-2 text-sm font-medium rounded-md transition-colors"
              style={{
                backgroundColor: theme.colors.border,
                color: theme.colors.text
              }}
            >
              Close Preview
            </button>
            <div>
              <h1 
                className="font-semibold"
                style={{ color: theme.colors.text }}
              >
                {theme.name} Preview
              </h1>
              <p 
                className="text-sm"
                style={{ color: theme.colors.textSecondary }}
              >
                {theme.description}
              </p>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <span 
              className="text-sm"
              style={{ color: theme.colors.textSecondary }}
            >
              Theme Colors:
            </span>
            <div className="flex space-x-1">
              <div 
                className="w-4 h-4 rounded-full border"
                style={{ backgroundColor: theme.colors.primary, borderColor: theme.colors.border }}
                title="Primary"
              />
              <div 
                className="w-4 h-4 rounded-full border"
                style={{ backgroundColor: theme.colors.secondary, borderColor: theme.colors.border }}
                title="Secondary"
              />
              <div 
                className="w-4 h-4 rounded-full border"
                style={{ backgroundColor: theme.colors.accent, borderColor: theme.colors.border }}
                title="Accent"
              />
            </div>
          </div>
        </div>
      </div>

      {/* Theme Content */}
      <div className="min-h-screen">
        {layout.map((component, index) => (
          <div key={component.id || index}>
            {renderComponent(component)}
          </div>
        ))}
        
        {/* Always show products section if available */}
        {products.length > 0 && (
          <div>
            {renderComponent({ type: 'products', content: {} })}
          </div>
        )}
      </div>
      </div>
    </div>
  );
};

export default ThemePreview;