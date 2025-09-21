import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { themeTemplates, sampleLayouts, sampleProducts } from '../../../data/shopTemplates';

const StandaloneThemePreview = () => {
  const { themeId } = useParams();
  const [theme, setTheme] = useState(null);
  const [layout, setLayout] = useState([]);
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (themeId && themeTemplates[themeId]) {
      const selectedTheme = themeTemplates[themeId];
      setTheme(selectedTheme);
      setLayout(sampleLayouts[themeId] || []);
      setProducts(sampleProducts[themeId] || []);
    }
    setLoading(false);
  }, [themeId]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading theme preview...</p>
        </div>
      </div>
    );
  }

  if (!theme) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Theme Not Found</h1>
          <p className="text-gray-600">The requested theme "{themeId}" could not be loaded.</p>
        </div>
      </div>
    );
  }

  // Apply theme styles to document
  useEffect(() => {
    if (theme) {
      const root = document.documentElement;
      Object.entries(theme.colors).forEach(([key, value]) => {
        root.style.setProperty(`--theme-${key}`, value);
      });

      // Apply typography
      if (theme.typography) {
        root.style.setProperty('--theme-font-family', theme.typography.fontFamily);
        root.style.setProperty('--theme-heading-font', theme.typography.headingFont);
      }
    }

    return () => {
      // Cleanup theme styles
      const root = document.documentElement;
      Object.keys(theme?.colors || {}).forEach(key => {
        root.style.removeProperty(`--theme-${key}`);
      });
      if (theme?.typography) {
        root.style.removeProperty('--theme-font-family');
        root.style.removeProperty('--theme-heading-font');
      }
    };
  }, [theme]);

  const renderComponent = (component) => {
    switch (component.type) {
      case 'header':
        return (
          <header key={component.id} className="bg-[var(--theme-background)] border-b border-[var(--theme-border)] p-4">
            <div className="container mx-auto flex items-center justify-between">
              <h1 className="text-2xl font-bold text-[var(--theme-text)]">{component.content?.title || 'Store Name'}</h1>
              <nav className="space-x-6">
                {['Home', 'Shop', 'About', 'Contact'].map(item => (
                  <a key={item} href="#" className="text-[var(--theme-text)] hover:text-[var(--theme-primary)]">{item}</a>
                ))}
              </nav>
            </div>
          </header>
        );

      case 'hero':
        return (
          <section key={component.id} className="bg-[var(--theme-surface)] py-16">
            <div className="container mx-auto text-center">
              <h2 className="text-4xl font-bold text-[var(--theme-text)] mb-4">
                {component.content?.title || 'Welcome to Our Store'}
              </h2>
              <p className="text-lg text-[var(--theme-textSecondary)] mb-8">
                {component.content?.subtitle || 'Discover amazing products'}
              </p>
              <button className="bg-[var(--theme-primary)] text-white px-8 py-3 rounded-lg hover:opacity-90">
                {component.content?.buttonText || 'Shop Now'}
              </button>
            </div>
          </section>
        );

      case 'product-grid':
        return (
          <section key={component.id} className="py-16 bg-[var(--theme-background)]">
            <div className="container mx-auto">
              <h2 className="text-3xl font-bold text-[var(--theme-text)] text-center mb-8">
                {component.content?.title || 'Featured Products'}
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {products.slice(0, 6).map((product, index) => (
                  <div key={index} className="bg-[var(--theme-surface)] rounded-lg p-4 border border-[var(--theme-border)]">
                    <div className="h-48 bg-gray-200 rounded mb-4"></div>
                    <h3 className="font-semibold text-[var(--theme-text)]">{product.name}</h3>
                    <p className="text-[var(--theme-primary)] font-bold">${product.price}</p>
                  </div>
                ))}
              </div>
            </div>
          </section>
        );

      case 'footer':
        return (
          <footer key={component.id} className="bg-[var(--theme-text)] text-[var(--theme-background)] py-8">
            <div className="container mx-auto text-center">
              <p>&copy; 2024 {component.content?.title || 'Your Store'}. All rights reserved.</p>
            </div>
          </footer>
        );

      default:
        return (
          <div key={component.id} className="py-8 bg-[var(--theme-surface)]">
            <div className="container mx-auto text-center">
              <p className="text-[var(--theme-text)]">Component: {component.type}</p>
            </div>
          </div>
        );
    }
  };

  return (
    <div className="min-h-screen" style={{ fontFamily: 'var(--theme-font-family, Inter, sans-serif)' }}>
      {layout.map(renderComponent)}
    </div>
  );
};

export default StandaloneThemePreview;