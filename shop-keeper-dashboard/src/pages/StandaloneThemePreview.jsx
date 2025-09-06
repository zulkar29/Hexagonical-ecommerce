import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

const StandaloneThemePreview = () => {
  const { themeId } = useParams();
  const [theme, setTheme] = useState(null);
  const [layout, setLayout] = useState([]);
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);

  // Mock themes data - in real app this would come from API
  const mockThemes = {
    'modern-minimal': {
      id: 'modern-minimal',
      name: 'Modern Minimal',
      description: 'Clean and contemporary design with focus on simplicity',
      colors: {
        primary: '#2563eb',
        secondary: '#64748b',
        accent: '#f59e0b',
        background: '#ffffff',
        surface: '#f8fafc',
        text: '#1e293b',
        textSecondary: '#64748b',
        border: '#e2e8f0'
      },
      typography: {
        fontFamily: 'Inter, system-ui, sans-serif',
        headingFont: 'Inter, system-ui, sans-serif',
        fontSize: '16px',
        lineHeight: '1.6'
      },
      borderRadius: '8px',
      spacing: {
        xs: '4px',
        sm: '8px',
        md: '16px',
        lg: '24px',
        xl: '32px'
      }
    },
    'elegant-dark': {
      id: 'elegant-dark',
      name: 'Elegant Dark',
      description: 'Sophisticated dark theme with premium feel',
      colors: {
        primary: '#8b5cf6',
        secondary: '#a78bfa',
        accent: '#fbbf24',
        background: '#0f172a',
        surface: '#1e293b',
        text: '#f1f5f9',
        textSecondary: '#94a3b8',
        border: '#334155'
      },
      typography: {
        fontFamily: 'Playfair Display, serif',
        headingFont: 'Playfair Display, serif',
        fontSize: '16px',
        lineHeight: '1.7'
      },
      borderRadius: '12px',
      spacing: {
        xs: '6px',
        sm: '12px',
        md: '20px',
        lg: '28px',
        xl: '36px'
      }
    },
    'vibrant-colorful': {
      id: 'vibrant-colorful',
      name: 'Vibrant & Colorful',
      description: 'Bold and energetic design with vibrant colors',
      colors: {
        primary: '#ec4899',
        secondary: '#06b6d4',
        accent: '#84cc16',
        background: '#fefefe',
        surface: '#fef7ff',
        text: '#1f2937',
        textSecondary: '#6b7280',
        border: '#e5e7eb'
      },
      typography: {
        fontFamily: 'Poppins, sans-serif',
        headingFont: 'Poppins, sans-serif',
        fontSize: '15px',
        lineHeight: '1.6'
      },
      borderRadius: '16px',
      spacing: {
        xs: '4px',
        sm: '8px',
        md: '16px',
        lg: '24px',
        xl: '32px'
      }
    }
  };

  // Mock content data
  const content = {
    companyName: 'Your Store',
    tagline: 'Discover Amazing Products',
    heroTitle: 'Welcome to Our Store',
    heroSubtitle: 'Find the perfect products for your needs',
    links: ['Home', 'Products', 'About', 'Contact'],
    social: ['Facebook', 'Twitter', 'Instagram']
  };

  // Mock products data
  const mockProducts = [
    { id: 1, name: 'Premium Product 1', price: '$99.99', image: 'https://via.placeholder.com/300x200' },
    { id: 2, name: 'Featured Item 2', price: '$149.99', image: 'https://via.placeholder.com/300x200' },
    { id: 3, name: 'Best Seller 3', price: '$79.99', image: 'https://via.placeholder.com/300x200' },
    { id: 4, name: 'New Arrival 4', price: '$199.99', image: 'https://via.placeholder.com/300x200' }
  ];

  useEffect(() => {
    // Simulate loading theme data
    const loadTheme = async () => {
      setLoading(true);
      
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 500));
      
      const selectedTheme = mockThemes[themeId] || mockThemes['modern-minimal'];
      setTheme(selectedTheme);
      
      // Set default layout
      setLayout([
        { id: 'header', type: 'header', content: {} },
        { id: 'hero', type: 'hero', content: {} },
        { id: 'products', type: 'products', content: {} },
        { id: 'footer', type: 'footer', content: {} }
      ]);
      
      setProducts(mockProducts);
      setLoading(false);
    };

    loadTheme();
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
          <p className="text-gray-600">The requested theme could not be loaded.</p>
        </div>
      </div>
    );
  }

  const renderComponent = (component) => {
    switch (component.type) {
      case 'header':
        return (
          <header 
            className="sticky top-0 z-50 px-6 py-4 border-b backdrop-blur-sm"
            style={{
              backgroundColor: `${theme.colors.surface}f0`,
              borderColor: theme.colors.border,
              fontFamily: theme.typography.fontFamily
            }}
          >
            <div className="max-w-7xl mx-auto flex items-center justify-between">
              <div className="flex items-center space-x-8">
                <h1 
                  className="text-2xl font-bold"
                  style={{
                    color: theme.colors.primary,
                    fontFamily: theme.typography.headingFont
                  }}
                >
                  {content.companyName}
                </h1>
                <nav className="hidden md:flex space-x-6">
                  {content.links?.map((link, index) => (
                    <a 
                      key={index}
                      href="#" 
                      className="hover:opacity-75 transition-opacity"
                      style={{ color: theme.colors.text }}
                    >
                      {link}
                    </a>
                  ))}
                </nav>
              </div>
              <button 
                className="px-6 py-2 rounded-md font-medium transition-colors"
                style={{
                  backgroundColor: theme.colors.primary,
                  color: theme.colors.background,
                  borderRadius: theme.borderRadius
                }}
              >
                Shop Now
              </button>
            </div>
          </header>
        );

      case 'hero':
        return (
          <section 
            className="py-20 px-6"
            style={{
              backgroundColor: theme.colors.background,
              fontFamily: theme.typography.fontFamily
            }}
          >
            <div className="max-w-4xl mx-auto text-center">
              <h1 
                className="text-5xl md:text-6xl font-bold mb-6"
                style={{
                  color: theme.colors.text,
                  fontFamily: theme.typography.headingFont,
                  lineHeight: theme.typography.lineHeight
                }}
              >
                {content.heroTitle}
              </h1>
              <p 
                className="text-xl mb-8"
                style={{ color: theme.colors.textSecondary }}
              >
                {content.heroSubtitle}
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <button 
                  className="px-8 py-3 rounded-md font-semibold transition-colors"
                  style={{
                    backgroundColor: theme.colors.primary,
                    color: theme.colors.background,
                    borderRadius: theme.borderRadius
                  }}
                >
                  Explore Products
                </button>
                <button 
                  className="px-8 py-3 rounded-md font-semibold border transition-colors"
                  style={{
                    borderColor: theme.colors.border,
                    color: theme.colors.text,
                    borderRadius: theme.borderRadius
                  }}
                >
                  Learn More
                </button>
              </div>
            </div>
          </section>
        );

      case 'products':
        return (
          <section 
            className="py-16 px-6"
            style={{
              backgroundColor: theme.colors.surface,
              fontFamily: theme.typography.fontFamily
            }}
          >
            <div className="max-w-6xl mx-auto">
              <div className="text-center mb-12">
                <h2 
                  className="text-3xl md:text-4xl font-bold mb-4"
                  style={{
                    color: theme.colors.text,
                    fontFamily: theme.typography.headingFont
                  }}
                >
                  Featured Products
                </h2>
                <p style={{ color: theme.colors.textSecondary }}>
                  Discover our most popular items
                </p>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {products.map((product) => (
                  <div 
                    key={product.id}
                    className="rounded-lg overflow-hidden border transition-transform hover:scale-105"
                    style={{
                      backgroundColor: theme.colors.background,
                      borderColor: theme.colors.border,
                      borderRadius: theme.borderRadius
                    }}
                  >
                    <img 
                      src={product.image} 
                      alt={product.name}
                      className="w-full h-48 object-cover"
                    />
                    <div className="p-4">
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
      {/* Theme CSS Variables */}
      <style>{`
        :root {
          --primary-color: ${theme.colors.primary};
          --secondary-color: ${theme.colors.secondary};
          --accent-color: ${theme.colors.accent};
          --background-color: ${theme.colors.background};
          --text-color: ${theme.colors.text};
          --border-radius: ${theme.borderRadius};
        }
      `}</style>
      
      {/* Theme Content */}
      <div className="theme-preview">
        {layout.map((component, index) => (
          <div key={component.id || index}>
            {renderComponent(component)}
          </div>
        ))}
      </div>
    </div>
  );
};

export default StandaloneThemePreview;