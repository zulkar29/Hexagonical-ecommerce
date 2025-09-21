// Theme Templates for Store Designer

export const themeTemplates = {
  'modern-minimalist': {
    id: 'modern-minimalist',
    name: 'Modern Minimalist',
    description: 'Clean, simple design perfect for professional businesses',
    category: 'minimalist',
    recommended: true,
    preview: '/themes/modern-minimalist.jpg',
    colors: {
      primary: '#2563eb',
      secondary: '#64748b',
      accent: '#0ea5e9',
      background: '#ffffff',
      surface: '#f8fafc',
      text: '#1e293b',
      textSecondary: '#64748b',
      border: '#e2e8f0'
    },
    typography: {
      fontFamily: 'Inter, system-ui, sans-serif',
      headingFont: 'Inter, system-ui, sans-serif',
      fontSize: {
        xs: '0.75rem',
        sm: '0.875rem',
        base: '1rem',
        lg: '1.125rem',
        xl: '1.25rem',
        '2xl': '1.5rem',
        '3xl': '1.875rem',
        '4xl': '2.25rem'
      },
      fontWeight: {
        normal: '400',
        medium: '500',
        semibold: '600',
        bold: '700'
      }
    },
    spacing: {
      xs: '0.5rem',
      sm: '1rem',
      md: '1.5rem',
      lg: '2rem',
      xl: '3rem',
      '2xl': '4rem'
    },
    borderRadius: {
      sm: '0.25rem',
      md: '0.5rem',
      lg: '0.75rem',
      xl: '1rem'
    },
    components: {
      header: {
        backgroundColor: '#ffffff',
        textColor: '#1e293b',
        padding: '1rem 0',
        borderBottom: '1px solid #e2e8f0',
        fontFamily: 'Inter, system-ui, sans-serif'
      },
      hero: {
        backgroundColor: '#f8fafc',
        textColor: '#1e293b',
        padding: '4rem 0',
        fontFamily: 'Inter, system-ui, sans-serif'
      },
      footer: {
        backgroundColor: '#1f2937',
        textColor: '#ffffff',
        padding: '2rem 0',
        fontFamily: 'Inter, system-ui, sans-serif'
      },
      button: {
        backgroundColor: '#2563eb',
        textColor: '#ffffff',
        borderRadius: '0.5rem',
        padding: '0.75rem 1.5rem',
        fontWeight: '500'
      },
      card: {
        backgroundColor: '#ffffff',
        textColor: '#1e293b',
        borderRadius: '0.75rem',
        padding: '1.5rem',
        shadow: '0 1px 3px 0 rgb(0 0 0 / 0.1)'
      }
    }
  },

  'classic-ecommerce': {
    id: 'classic-ecommerce',
    name: 'Classic E-commerce',
    description: 'Proven layout optimized for online sales and conversions',
    category: 'ecommerce',
    recommended: true,
    preview: '/themes/classic-ecommerce.jpg',
    colors: {
      primary: '#dc2626',
      secondary: '#374151',
      accent: '#f59e0b',
      background: '#ffffff',
      surface: '#f9fafb',
      text: '#111827',
      textSecondary: '#6b7280',
      border: '#d1d5db'
    },
    typography: {
      fontFamily: 'Roboto, Arial, sans-serif',
      headingFont: 'Roboto, Arial, sans-serif',
      fontSize: {
        xs: '0.75rem',
        sm: '0.875rem',
        base: '1rem',
        lg: '1.125rem',
        xl: '1.25rem',
        '2xl': '1.5rem',
        '3xl': '1.875rem',
        '4xl': '2.25rem'
      },
      fontWeight: {
        normal: '400',
        medium: '500',
        semibold: '600',
        bold: '700'
      }
    },
    spacing: {
      xs: '0.5rem',
      sm: '1rem',
      md: '1.5rem',
      lg: '2rem',
      xl: '2.5rem',
      '2xl': '3rem'
    },
    borderRadius: {
      sm: '0.125rem',
      md: '0.25rem',
      lg: '0.5rem',
      xl: '0.75rem'
    },
    components: {
      header: {
        backgroundColor: '#ffffff',
        textColor: '#111827',
        padding: '1rem 0',
        borderBottom: '2px solid #e5e7eb',
        fontFamily: 'Roboto, Arial, sans-serif'
      },
      hero: {
        backgroundColor: '#f9fafb',
        textColor: '#111827',
        padding: '4rem 0',
        fontFamily: 'Roboto, Arial, sans-serif'
      },
      footer: {
        backgroundColor: '#374151',
        textColor: '#ffffff',
        padding: '2rem 0',
        fontFamily: 'Roboto, Arial, sans-serif'
      },
      button: {
        backgroundColor: '#dc2626',
        textColor: '#ffffff',
        borderRadius: '0.25rem',
        padding: '0.75rem 1.25rem',
        fontWeight: '600'
      },
      card: {
        backgroundColor: '#ffffff',
        textColor: '#111827',
        borderRadius: '0.5rem',
        padding: '1.25rem',
        shadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)'
      }
    }
  },

  'bold-colorful': {
    id: 'bold-colorful',
    name: 'Bold & Colorful',
    description: 'Eye-catching design for creative and lifestyle brands',
    category: 'creative',
    preview: '/themes/bold-colorful.jpg',
    colors: {
      primary: '#7c3aed',
      secondary: '#ec4899',
      accent: '#06b6d4',
      background: '#fefefe',
      surface: '#f0f9ff',
      text: '#1f2937',
      textSecondary: '#4b5563',
      border: '#e5e7eb'
    },
    typography: {
      fontFamily: 'Poppins, system-ui, sans-serif',
      headingFont: 'Poppins, system-ui, sans-serif',
      fontSize: {
        xs: '0.75rem',
        sm: '0.875rem',
        base: '1rem',
        lg: '1.125rem',
        xl: '1.25rem',
        '2xl': '1.5rem',
        '3xl': '1.875rem',
        '4xl': '2.25rem'
      },
      fontWeight: {
        normal: '400',
        medium: '500',
        semibold: '600',
        bold: '700'
      }
    },
    spacing: {
      xs: '0.75rem',
      sm: '1.25rem',
      md: '2rem',
      lg: '2.5rem',
      xl: '3.5rem',
      '2xl': '4.5rem'
    },
    borderRadius: {
      sm: '0.5rem',
      md: '0.75rem',
      lg: '1rem',
      xl: '1.5rem'
    },
    components: {
      header: {
        background: 'linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)',
        padding: '1.5rem 0',
        borderBottom: 'none'
      },
      button: {
        borderRadius: '0.75rem',
        padding: '1rem 2rem',
        fontWeight: '600'
      },
      card: {
        borderRadius: '1rem',
        padding: '2rem',
        shadow: '0 10px 25px -3px rgb(0 0 0 / 0.1)'
      }
    }
  },

  'professional-corporate': {
    id: 'professional-corporate',
    name: 'Professional Corporate',
    description: 'Enterprise-grade design for B2B and corporate sales',
    category: 'business',
    recommended: true,
    preview: '/themes/professional-corporate.jpg',
    colors: {
      primary: '#1e40af',
      secondary: '#475569',
      accent: '#059669',
      background: '#ffffff',
      surface: '#f8fafc',
      text: '#0f172a',
      textSecondary: '#64748b',
      border: '#cbd5e1'
    },
    typography: {
      fontFamily: 'Source Sans Pro, Arial, sans-serif',
      headingFont: 'Source Sans Pro, Arial, sans-serif',
      fontSize: {
        xs: '0.75rem',
        sm: '0.875rem',
        base: '1rem',
        lg: '1.125rem',
        xl: '1.25rem',
        '2xl': '1.5rem',
        '3xl': '1.875rem',
        '4xl': '2.25rem'
      },
      fontWeight: {
        normal: '400',
        medium: '500',
        semibold: '600',
        bold: '700'
      }
    },
    spacing: {
      xs: '0.5rem',
      sm: '1rem',
      md: '1.5rem',
      lg: '2rem',
      xl: '3rem',
      '2xl': '4rem'
    },
    borderRadius: {
      sm: '0.125rem',
      md: '0.375rem',
      lg: '0.5rem',
      xl: '0.75rem'
    },
    components: {
      header: {
        background: '#ffffff',
        padding: '1.25rem 0',
        borderBottom: '1px solid #e2e8f0'
      },
      button: {
        borderRadius: '0.375rem',
        padding: '0.75rem 1.5rem',
        fontWeight: '600'
      },
      card: {
        borderRadius: '0.5rem',
        padding: '1.5rem',
        shadow: '0 1px 3px 0 rgb(0 0 0 / 0.1)'
      }
    }
  },

  'creative-portfolio': {
    id: 'creative-portfolio',
    name: 'Creative Portfolio',
    description: 'Artistic and expressive design for creative professionals',
    category: 'creative',
    preview: '/themes/creative-portfolio.jpg',
    colors: {
      primary: '#f59e0b',
      secondary: '#8b5cf6',
      accent: '#ef4444',
      background: '#fafafa',
      surface: '#ffffff',
      text: '#18181b',
      textSecondary: '#71717a',
      border: '#e4e4e7'
    },
    typography: {
      fontFamily: 'Montserrat, system-ui, sans-serif',
      headingFont: 'Playfair Display, serif',
      fontSize: {
        xs: '0.75rem',
        sm: '0.875rem',
        base: '1rem',
        lg: '1.125rem',
        xl: '1.25rem',
        '2xl': '1.5rem',
        '3xl': '1.875rem',
        '4xl': '2.25rem'
      },
      fontWeight: {
        normal: '400',
        medium: '500',
        semibold: '600',
        bold: '700'
      }
    },
    spacing: {
      xs: '0.75rem',
      sm: '1.25rem',
      md: '2rem',
      lg: '3rem',
      xl: '4rem',
      '2xl': '5rem'
    },
    borderRadius: {
      sm: '0.25rem',
      md: '0.5rem',
      lg: '1rem',
      xl: '1.5rem'
    },
    components: {
      header: {
        background: 'transparent',
        padding: '2rem 0',
        borderBottom: 'none'
      },
      button: {
        borderRadius: '2rem',
        padding: '0.875rem 2rem',
        fontWeight: '500'
      },
      card: {
        borderRadius: '1rem',
        padding: '2rem',
        shadow: '0 8px 25px -8px rgb(0 0 0 / 0.1)'
      }
    }
  },

  'luxury-brand': {
    id: 'luxury-brand',
    name: 'Luxury Brand',
    description: 'Elegant and premium design for high-end brands',
    category: 'luxury',
    preview: '/themes/luxury-brand.jpg',
    colors: {
      primary: '#000000',
      secondary: '#6b7280',
      accent: '#d4af37',
      background: '#ffffff',
      surface: '#fafafa',
      text: '#000000',
      textSecondary: '#6b7280',
      border: '#e5e7eb'
    },
    typography: {
      fontFamily: 'Crimson Text, serif',
      headingFont: 'Crimson Text, serif',
      fontSize: {
        xs: '0.75rem',
        sm: '0.875rem',
        base: '1rem',
        lg: '1.125rem',
        xl: '1.25rem',
        '2xl': '1.5rem',
        '3xl': '1.875rem',
        '4xl': '2.25rem'
      },
      fontWeight: {
        normal: '400',
        medium: '500',
        semibold: '600',
        bold: '700'
      }
    },
    spacing: {
      xs: '0.5rem',
      sm: '1rem',
      md: '1.5rem',
      lg: '2.5rem',
      xl: '4rem',
      '2xl': '6rem'
    },
    borderRadius: {
      sm: '0rem',
      md: '0.125rem',
      lg: '0.25rem',
      xl: '0.375rem'
    },
    components: {
      header: {
        background: '#ffffff',
        padding: '2rem 0',
        borderBottom: '1px solid #000000'
      },
      button: {
        borderRadius: '0rem',
        padding: '1rem 2.5rem',
        fontWeight: '400'
      },
      card: {
        borderRadius: '0.125rem',
        padding: '2.5rem',
        shadow: '0 2px 4px 0 rgb(0 0 0 / 0.05)'
      }
    }
  }
};

export const themeCategories = {
  minimalist: {
    name: 'Minimalist',
    description: 'Clean and simple designs',
    icon: 'Minimize2'
  },
  ecommerce: {
    name: 'E-commerce',
    description: 'Online store layouts',
    icon: 'ShoppingBag'
  },
  creative: {
    name: 'Creative',
    description: 'Artistic and expressive',
    icon: 'Palette'
  },
  business: {
    name: 'Business',
    description: 'Professional and corporate',
    icon: 'Briefcase'
  },
  luxury: {
    name: 'Luxury',
    description: 'Premium and elegant',
    icon: 'Crown'
  }
};

export const sampleLayouts = {
  'modern-minimalist': [
    { 
      type: 'header', 
      id: 'header-1',
      content: {
        title: 'MinimalStore',
        navigation: ['Home', 'Products', 'About', 'Contact'],
        style: 'clean'
      }
    },
    { 
      type: 'hero', 
      id: 'hero-1',
      content: {
        title: 'Simple. Beautiful. Functional.',
        subtitle: 'Discover our curated collection of minimalist products',
        buttonText: 'Shop Now',
        backgroundImage: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=minimalist%20white%20background%20with%20subtle%20geometric%20shapes&image_size=landscape_16_9'
      }
    },
    { 
      type: 'footer', 
      id: 'footer-1',
      content: {
        companyName: 'MinimalStore',
        links: ['Privacy', 'Terms', 'Support'],
        social: ['Instagram', 'Twitter']
      }
    }
  ],
  'classic-ecommerce': [
    { 
      type: 'header', 
      id: 'header-2',
      content: {
        title: 'ClassicShop',
        navigation: ['Home', 'Categories', 'Deals', 'Account'],
        style: 'traditional'
      }
    },
    { 
      type: 'hero', 
      id: 'hero-2',
      content: {
        title: 'Welcome to ClassicShop',
        subtitle: 'Your trusted online marketplace since 1995',
        buttonText: 'Browse Products',
        backgroundImage: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=classic%20ecommerce%20storefront%20with%20warm%20lighting&image_size=landscape_16_9'
      }
    },
    { 
      type: 'footer', 
      id: 'footer-2',
      content: {
        companyName: 'ClassicShop',
        links: ['About Us', 'Customer Service', 'Returns'],
        social: ['Facebook', 'Twitter', 'Instagram']
      }
    }
  ],
  'bold-colorful': [
    { 
      type: 'header', 
      id: 'header-3',
      content: {
        title: 'VibrantMarket',
        navigation: ['Explore', 'Trending', 'Create', 'Community'],
        style: 'vibrant'
      }
    },
    { 
      type: 'hero', 
      id: 'hero-3',
      content: {
        title: 'Express Yourself!',
        subtitle: 'Bold products for bold personalities',
        buttonText: 'Get Creative',
        backgroundImage: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=vibrant%20colorful%20abstract%20background%20with%20bold%20geometric%20patterns&image_size=landscape_16_9'
      }
    },
    { 
      type: 'footer', 
      id: 'footer-3',
      content: {
        companyName: 'VibrantMarket',
        links: ['Community', 'Creators', 'Blog'],
        social: ['TikTok', 'Instagram', 'YouTube']
      }
    }
  ],
  'professional-corporate': [
    { 
      type: 'header', 
      id: 'header-4',
      content: {
        title: 'ProBusiness',
        navigation: ['Solutions', 'Services', 'Resources', 'Contact'],
        style: 'corporate'
      }
    },
    { 
      type: 'hero', 
      id: 'hero-4',
      content: {
        title: 'Professional Solutions',
        subtitle: 'Empowering businesses with reliable products and services',
        buttonText: 'Learn More',
        backgroundImage: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=professional%20corporate%20office%20environment%20with%20clean%20lines&image_size=landscape_16_9'
      }
    },
    { 
      type: 'footer', 
      id: 'footer-4',
      content: {
        companyName: 'ProBusiness Corp',
        links: ['Legal', 'Compliance', 'Careers'],
        social: ['LinkedIn', 'Twitter']
      }
    }
  ],
  'creative-portfolio': [
    { 
      type: 'header', 
      id: 'header-5',
      content: {
        title: 'ArtisanStudio',
        navigation: ['Portfolio', 'Process', 'Stories', 'Connect'],
        style: 'artistic'
      }
    },
    { 
      type: 'hero', 
      id: 'hero-5',
      content: {
        title: 'Handcrafted with Love',
        subtitle: 'Unique artisan creations that tell a story',
        buttonText: 'View Collection',
        backgroundImage: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=artistic%20creative%20workspace%20with%20handmade%20crafts%20and%20warm%20lighting&image_size=landscape_16_9'
      }
    },
    { 
      type: 'footer', 
      id: 'footer-5',
      content: {
        companyName: 'ArtisanStudio',
        links: ['Artist Bio', 'Process', 'Custom Orders'],
        social: ['Instagram', 'Pinterest', 'Etsy']
      }
    }
  ],
  'luxury-brand': [
    { 
      type: 'header', 
      id: 'header-6',
      content: {
        title: 'LUXE',
        navigation: ['Collections', 'Heritage', 'Atelier', 'Concierge'],
        style: 'luxury'
      }
    },
    { 
      type: 'hero', 
      id: 'hero-6',
      content: {
        title: 'Timeless Elegance',
        subtitle: 'Exquisite craftsmanship meets modern sophistication',
        buttonText: 'Explore Collection',
        backgroundImage: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=luxury%20elegant%20interior%20with%20gold%20accents%20and%20marble%20textures&image_size=landscape_16_9'
      }
    },
    { 
      type: 'footer', 
      id: 'footer-6',
      content: {
        companyName: 'LUXE Maison',
        links: ['Heritage', 'Craftsmanship', 'Boutiques'],
        social: ['Instagram', 'Facebook']
      }
    }
  ]
};

export const sampleProducts = {
  'modern-minimalist': [
    {
      id: 'min-1',
      name: 'Essential Tote',
      price: '$89',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=minimalist%20white%20leather%20tote%20bag%20on%20clean%20background&image_size=square'
    },
    {
      id: 'min-2',
      name: 'Pure Ceramic Mug',
      price: '$24',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=simple%20white%20ceramic%20coffee%20mug%20minimalist%20design&image_size=square'
    }
  ],
  'classic-ecommerce': [
    {
      id: 'classic-1',
      name: 'Heritage Watch',
      price: '$299',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=classic%20analog%20watch%20with%20leather%20strap%20traditional%20design&image_size=square'
    },
    {
      id: 'classic-2',
      name: 'Wool Sweater',
      price: '$79',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=classic%20wool%20sweater%20in%20navy%20blue%20traditional%20style&image_size=square'
    }
  ],
  'bold-colorful': [
    {
      id: 'bold-1',
      name: 'Neon Sneakers',
      price: '$149',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=bright%20neon%20colored%20sneakers%20with%20bold%20design%20elements&image_size=square'
    },
    {
      id: 'bold-2',
      name: 'Rainbow Backpack',
      price: '$89',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=colorful%20rainbow%20backpack%20with%20vibrant%20patterns&image_size=square'
    }
  ],
  'professional-corporate': [
    {
      id: 'corp-1',
      name: 'Executive Briefcase',
      price: '$399',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=professional%20black%20leather%20briefcase%20corporate%20style&image_size=square'
    },
    {
      id: 'corp-2',
      name: 'Business Planner',
      price: '$49',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=professional%20business%20planner%20with%20clean%20design&image_size=square'
    }
  ],
  'creative-portfolio': [
    {
      id: 'art-1',
      name: 'Handwoven Scarf',
      price: '$125',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=handwoven%20artistic%20scarf%20with%20unique%20patterns%20and%20textures&image_size=square'
    },
    {
      id: 'art-2',
      name: 'Ceramic Vase',
      price: '$89',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=handmade%20ceramic%20vase%20with%20artistic%20glaze%20and%20unique%20shape&image_size=square'
    }
  ],
  'luxury-brand': [
    {
      id: 'lux-1',
      name: 'Diamond Necklace',
      price: '$2,999',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=elegant%20diamond%20necklace%20on%20luxury%20velvet%20background&image_size=square'
    },
    {
      id: 'lux-2',
      name: 'Silk Evening Gown',
      price: '$1,299',
      image: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=luxury%20silk%20evening%20gown%20in%20elegant%20setting&image_size=square'
    }
  ]
};

export const componentTemplates = [
  // Header Components
  {
    id: 'header-simple',
    name: 'Simple Header',
    type: 'header',
    category: 'header',
    description: 'Clean header with logo and navigation',
    icon: '🏠',
    preview: 'Logo | Home Products About Contact | [CTA]',
    defaultProps: {
      title: 'Your Logo',
      menuItems: ['Home', 'Products', 'About', 'Contact'],
      ctaText: 'Get Started'
    },
    defaultStyles: {
      backgroundColor: '#ffffff',
      textColor: '#000000'
    }
  },
  {
    id: 'header-modern',
    name: 'Modern Header',
    type: 'header',
    category: 'header',
    description: 'Modern header with gradient background',
    icon: '✨',
    preview: 'Logo | Navigation | [CTA Button]',
    defaultProps: {
      title: 'ModernCorp',
      menuItems: ['Home', 'Services', 'Portfolio', 'Contact'],
      ctaText: 'Start Free Trial'
    },
    defaultStyles: {
      backgroundColor: '#1f2937',
      textColor: '#ffffff'
    }
  },
  {
    id: 'header-minimal',
    name: 'Minimal Header',
    type: 'header',
    category: 'header',
    description: 'Minimalist header design',
    icon: '⚡',
    preview: 'Brand | Menu | Action',
    defaultProps: {
      title: 'Brand',
      menuItems: ['Work', 'About', 'Contact'],
      ctaText: 'Hire Us'
    },
    defaultStyles: {
      backgroundColor: '#f8f9fa',
      textColor: '#374151'
    }
  },

  // Content Components
  {
    id: 'hero-centered',
    name: 'Centered Hero',
    type: 'hero',
    category: 'content',
    description: 'Hero section with centered content',
    icon: '🎯',
    preview: 'Large Title\nSubtitle text\n[Primary] [Secondary]',
    defaultProps: {
      title: 'Welcome to Our Platform',
      subtitle: 'Build amazing experiences with our powerful tools and intuitive interface.',
      primaryCta: 'Get Started',
      secondaryCta: 'Learn More'
    },
    defaultStyles: {
      backgroundColor: '#f8f9fa',
      textColor: '#000000'
    }
  },
  {
    id: 'hero-gradient',
    name: 'Gradient Hero',
    type: 'hero',
    category: 'content',
    description: 'Hero with gradient background',
    icon: '🌈',
    preview: 'Bold Statement\nDescription\n[Action Button]',
    defaultProps: {
      title: 'Transform Your Business',
      subtitle: 'Unlock new possibilities with our innovative solutions designed for modern enterprises.',
      primaryCta: 'Start Today',
      secondaryCta: 'Watch Demo'
    },
    defaultStyles: {
      backgroundColor: '#3b82f6',
      textColor: '#ffffff'
    }
  },
  {
    id: 'hero-minimal',
    name: 'Minimal Hero',
    type: 'hero',
    category: 'content',
    description: 'Clean and minimal hero section',
    icon: '🎨',
    preview: 'Simple Title\nClean subtitle\n[CTA]',
    defaultProps: {
      title: 'Simple. Powerful. Effective.',
      subtitle: 'Everything you need, nothing you don\'t.',
      primaryCta: 'Try It Free',
      secondaryCta: 'See Features'
    },
    defaultStyles: {
      backgroundColor: '#ffffff',
      textColor: '#1f2937'
    }
  },

  // Product Components
  {
    id: 'product-grid',
    name: 'Product Grid',
    type: 'product',
    category: 'product',
    description: 'Grid layout for showcasing products',
    icon: '🛍️',
    preview: '[Product 1] [Product 2] [Product 3]\n[Product 4] [Product 5] [Product 6]',
    defaultProps: {
      title: 'Featured Products',
      columns: 3,
      showPrices: true,
      showRatings: true,
      products: [
        { id: 1, name: 'Premium Wireless Headphones', price: '$199.99', rating: 4.5, image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=300&h=300&fit=crop' },
        { id: 2, name: 'Smart Fitness Watch', price: '$299.99', rating: 4.8, image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=300&h=300&fit=crop' },
        { id: 3, name: 'Laptop Stand', price: '$79.99', rating: 4.2, image: 'https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=300&h=300&fit=crop' }
      ]
    },
    defaultStyles: {
      backgroundColor: '#ffffff',
      textColor: '#000000',
      padding: '2rem 0'
    }
  },
  {
    id: 'product-showcase',
    name: 'Product Showcase',
    type: 'product',
    category: 'product',
    description: 'Single product spotlight with detailed info',
    icon: '⭐',
    preview: '[Large Product Image] | Product Details\nDescription & Reviews',
    defaultProps: {
      productName: 'Professional Wireless Speaker',
      price: '$299.99',
      originalPrice: '$399.99',
      rating: 4.9,
      reviews: 127,
      description: 'Premium sound quality with 360-degree audio and smart connectivity features.',
      features: ['360° Sound Technology', 'Waterproof Design', 'Smart Voice Assistant', '20-hour Battery Life'],
      images: ['https://images.unsplash.com/photo-1608043152269-423dbba4e7e1?w=500&h=500&fit=crop']
    },
    defaultStyles: {
      backgroundColor: '#f8f9fa',
      textColor: '#000000',
      padding: '3rem 0'
    }
  },
  {
    id: 'category-banner',
    name: 'Category Banner',
    type: 'category',
    category: 'product',
    description: 'Promotional banner for product categories',
    icon: '🏷️',
    preview: '[Background Image] Category Title\nShop Now Button',
    defaultProps: {
      title: 'New Collection',
      subtitle: 'Discover our latest arrivals',
      buttonText: 'Shop Now',
      backgroundImage: 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=1200&h=400&fit=crop',
      overlay: true
    },
    defaultStyles: {
      backgroundColor: '#000000',
      textColor: '#ffffff',
      padding: '4rem 0'
    }
  },

  // Footer Components
  {
    id: 'footer-comprehensive',
    name: 'Comprehensive Footer',
    type: 'footer',
    category: 'footer',
    description: 'Full-featured footer with multiple sections',
    icon: '🦶',
    preview: 'Company | Links | Contact\nSocial Links\nCopyright',
    defaultProps: {
      companyName: 'Your Store',
      description: 'Your trusted online shopping destination',
      sections: [
        {
          title: 'Shop',
          links: ['New Arrivals', 'Best Sellers', 'Sale', 'Gift Cards']
        },
        {
          title: 'Support',
          links: ['Contact Us', 'FAQ', 'Shipping', 'Returns']
        },
        {
          title: 'Company',
          links: ['About Us', 'Careers', 'Press', 'Blog']
        }
      ],
      socialLinks: [
        { platform: 'Facebook', url: '#' },
        { platform: 'Twitter', url: '#' },
        { platform: 'Instagram', url: '#' }
      ]
    },
    defaultStyles: {
      backgroundColor: '#1f2937',
      textColor: '#ffffff',
      padding: '3rem 0'
    }
  },
  {
    id: 'footer-simple',
    name: 'Simple Footer',
    type: 'footer',
    category: 'footer',
    description: 'Clean and simple footer',
    icon: '📄',
    preview: 'Company Info | Links | Contact\nCopyright',
    defaultProps: {
      companyName: 'StartupCo',
      description: 'Making technology accessible to everyone.',
      email: 'hello@startup.co',
      phone: '+1 (555) 987-6543',
      quickLinks: ['Home', 'About', 'Contact']
    },
    defaultStyles: {
      backgroundColor: '#f8f9fa',
      textColor: '#6b7280'
    }
  },
  {
    id: 'footer-minimal',
    name: 'Minimal Footer',
    type: 'footer',
    category: 'footer',
    description: 'Minimalist footer design',
    icon: '📝',
    preview: 'Brand | Essential Links\n© Year Company',
    defaultProps: {
      companyName: 'MinimalCorp',
      description: 'Less is more.',
      email: 'info@minimal.com',
      phone: '+1 (555) 111-2222',
      quickLinks: ['Privacy', 'Terms']
    },
    defaultStyles: {
      backgroundColor: '#ffffff',
      textColor: '#9ca3af'
    }
  }
];

// Helper function to get template by ID
export const getTemplateById = (id) => {
  return componentTemplates.find(template => template.id === id);
};

// Helper function to get templates by category
export const getTemplatesByCategory = (category) => {
  return componentTemplates.filter(template => template.category === category);
};

// Helper function to get templates by type
export const getTemplatesByType = (type) => {
  return componentTemplates.filter(template => template.type === type);
};

// Get all available categories
export const getCategories = () => {
  const categories = [...new Set(componentTemplates.map(template => template.category))];
  return categories.sort();
};

// Get category display names
export const getCategoryDisplayName = (category) => {
  const displayNames = {
    'header': 'Headers',
    'footer': 'Footers',
    'content': 'Content Blocks',
    'product': 'Product Components',
    'slider': 'Sliders & Carousels'
  };
  return displayNames[category] || category;
};

export default { themeTemplates, themeCategories, sampleLayouts, componentTemplates };