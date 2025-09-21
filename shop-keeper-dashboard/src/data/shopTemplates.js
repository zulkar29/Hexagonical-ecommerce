// Shop Templates and Themes for Store Customization

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
        background: '#ffffff',
        padding: '1rem 0',
        borderBottom: '1px solid #e2e8f0'
      },
      button: {
        borderRadius: '0.5rem',
        padding: '0.75rem 1.5rem',
        fontWeight: '500'
      },
      card: {
        borderRadius: '0.75rem',
        padding: '1.5rem',
        shadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)'
      }
    }
  },

  'dark-modern': {
    id: 'dark-modern',
    name: 'Dark Modern',
    description: 'Sophisticated dark theme for premium brands',
    category: 'dark',
    premium: true,
    preview: '/themes/dark-modern.jpg',
    colors: {
      primary: '#3b82f6',
      secondary: '#6366f1',
      accent: '#10b981',
      background: '#0f172a',
      surface: '#1e293b',
      text: '#f1f5f9',
      textSecondary: '#94a3b8',
      border: '#334155'
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
      xs: '0.5rem',
      sm: '1rem',
      md: '1.5rem',
      lg: '2rem',
      xl: '3rem',
      '2xl': '4rem'
    },
    borderRadius: {
      sm: '0.5rem',
      md: '0.75rem',
      lg: '1rem',
      xl: '1.5rem'
    },
    components: {
      header: {
        background: 'rgba(15, 23, 42, 0.8)',
        padding: '1rem 0',
        borderBottom: '1px solid #334155',
        backdropFilter: 'blur(10px)'
      },
      button: {
        borderRadius: '0.75rem',
        padding: '0.75rem 1.5rem',
        fontWeight: '600'
      },
      card: {
        borderRadius: '1rem',
        padding: '1.5rem',
        shadow: '0 10px 15px -3px rgb(0 0 0 / 0.3)'
      }
    }
  },

  'vibrant-creative': {
    id: 'vibrant-creative',
    name: 'Vibrant Creative',
    description: 'Bold, colorful design for creative businesses',
    category: 'creative',
    preview: '/themes/vibrant-creative.jpg',
    colors: {
      primary: '#f59e0b',
      secondary: '#ec4899',
      accent: '#10b981',
      background: '#fefce8',
      surface: '#ffffff',
      text: '#1f2937',
      textSecondary: '#6b7280',
      border: '#f3f4f6'
    },
    typography: {
      fontFamily: 'Nunito, system-ui, sans-serif',
      headingFont: 'Nunito, system-ui, sans-serif',
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
        background: 'linear-gradient(135deg, #f59e0b 0%, #ec4899 100%)',
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
    description: 'Conservative, trustworthy design for corporate use',
    category: 'corporate',
    preview: '/themes/professional-corporate.jpg',
    colors: {
      primary: '#1e40af',
      secondary: '#374151',
      accent: '#059669',
      background: '#f9fafb',
      surface: '#ffffff',
      text: '#111827',
      textSecondary: '#6b7280',
      border: '#d1d5db'
    },
    typography: {
      fontFamily: 'Roboto, system-ui, sans-serif',
      headingFont: 'Roboto, system-ui, sans-serif',
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
      md: '0.375rem',
      lg: '0.5rem',
      xl: '0.75rem'
    },
    components: {
      header: {
        background: '#ffffff',
        padding: '1rem 0',
        borderBottom: '2px solid #e5e7eb'
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

  'luxe-elegant': {
    id: 'luxe-elegant',
    name: 'Luxe Elegant',
    description: 'Sophisticated luxury design with gold accents',
    category: 'luxury',
    premium: true,
    preview: '/themes/luxe-elegant.jpg',
    colors: {
      primary: '#d97706',
      secondary: '#78716c',
      accent: '#fbbf24',
      background: '#fafaf9',
      surface: '#ffffff',
      text: '#1c1917',
      textSecondary: '#57534e',
      border: '#e7e5e4'
    },
    typography: {
      fontFamily: 'Playfair Display, serif',
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
      '2xl': '6rem'
    },
    borderRadius: {
      sm: '0.125rem',
      md: '0.25rem',
      lg: '0.375rem',
      xl: '0.5rem'
    },
    components: {
      header: {
        background: 'linear-gradient(to right, #fafaf9, #f5f5f4)',
        padding: '1.5rem 0',
        borderBottom: '1px solid #d97706'
      },
      button: {
        borderRadius: '0.25rem',
        padding: '1rem 2rem',
        fontWeight: '600'
      },
      card: {
        borderRadius: '0.375rem',
        padding: '2rem',
        shadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)'
      }
    }
  },

  'nature-organic': {
    id: 'nature-organic',
    name: 'Nature Organic',
    description: 'Earthy, sustainable design for eco-friendly brands',
    category: 'nature',
    preview: '/themes/nature-organic.jpg',
    colors: {
      primary: '#16a34a',
      secondary: '#84cc16',
      accent: '#eab308',
      background: '#f7fee7',
      surface: '#ffffff',
      text: '#365314',
      textSecondary: '#65a30d',
      border: '#d9f99d'
    },
    typography: {
      fontFamily: 'Lora, serif',
      headingFont: 'Lora, serif',
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
      sm: '1rem',
      md: '1.5rem',
      lg: '2rem',
      xl: '3rem'
    },
    components: {
      header: {
        background: '#ffffff',
        padding: '1.25rem 0',
        borderBottom: '2px solid #16a34a'
      },
      button: {
        borderRadius: '1.5rem',
        padding: '0.75rem 2rem',
        fontWeight: '500'
      },
      card: {
        borderRadius: '2rem',
        padding: '1.5rem',
        shadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)'
      }
    }
  }
};

export const componentTemplates = [
  // ===== HEADER COMPONENTS =====
  {
    id: 'modern-header',
    name: 'Modern Header',
    type: 'header',
    category: 'header',
    description: 'Clean modern header with logo and navigation',
    icon: '🎯',
    preview: 'Logo | Nav Menu | Actions',
    defaultProps: {
      storeName: 'Modern Store',
      menuItems: ['Shop', 'About', 'Contact'],
      showSearch: true,
      showCart: true,
      showAccount: true,
      cartCount: 0,
      logoSize: 'medium'
    },
    defaultStyles: {
      backgroundColor: '#ffffff',
      textColor: '#1f2937',
      borderBottom: '1px solid #e5e7eb',
      padding: '1rem 0'
    }
  },

  {
    id: 'mega-menu-header',
    name: 'Mega Menu Header',
    type: 'header',
    category: 'header',
    description: 'Header with dropdown mega menu for categories',
    icon: '📋',
    preview: 'Logo | Mega Menu | Search | Cart',
    defaultProps: {
      storeName: 'Mega Store',
      menuItems: ['Electronics', 'Fashion', 'Home', 'Sports'],
      showMegaMenu: true,
      showSearch: true,
      showCart: true,
      showWishlist: true,
      cartCount: 3
    },
    defaultStyles: {
      backgroundColor: '#ffffff',
      textColor: '#111827',
      borderBottom: '2px solid #f3f4f6',
      padding: '1.25rem 0'
    }
  },

  {
    id: 'sticky-header',
    name: 'Sticky Header',
    type: 'header',
    category: 'header',
    description: 'Header that stays fixed at top while scrolling',
    icon: '📌',
    preview: 'Fixed | Logo | Menu | Cart',
    defaultProps: {
      storeName: 'Sticky Store',
      menuItems: ['Products', 'Deals', 'Support'],
      isSticky: true,
      showSearch: true,
      showCart: true,
      cartCount: 1
    },
    defaultStyles: {
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      textColor: '#1f2937',
      backdropFilter: 'blur(10px)',
      borderBottom: '1px solid #e5e7eb'
    }
  },

  // ===== HERO COMPONENTS =====
  {
    id: 'hero-video-bg',
    name: 'Video Background Hero',
    type: 'hero',
    category: 'hero',
    description: 'Hero section with video background',
    icon: '🎬',
    preview: 'Video BG | Title | CTA',
    defaultProps: {
      title: 'Transform Your Style',
      subtitle: 'Discover our latest collection of premium products',
      primaryCta: 'Shop Collection',
      secondaryCta: 'Watch Story',
      videoUrl: '/video/hero-bg.mp4',
      overlayOpacity: 0.4
    },
    defaultStyles: {
      minHeight: '80vh',
      position: 'relative',
      textAlign: 'center',
      color: '#ffffff'
    }
  },

  {
    id: 'hero-split-content',
    name: 'Split Content Hero',
    type: 'hero',
    category: 'hero',
    description: 'Hero with content on left, image on right',
    icon: '↔️',
    preview: 'Text Left | Image Right',
    defaultProps: {
      title: 'Premium Quality Products',
      subtitle: 'Crafted with attention to detail and built to last',
      primaryCta: 'Explore Products',
      secondaryCta: 'Learn More',
      heroImage: '/images/hero-product.jpg',
      imagePosition: 'right'
    },
    defaultStyles: {
      padding: '4rem 0',
      backgroundColor: '#f9fafb'
    }
  },

  {
    id: 'hero-carousel',
    name: 'Carousel Hero',
    type: 'hero',
    category: 'hero',
    description: 'Rotating hero slides with multiple offers',
    icon: '🎠',
    preview: 'Slide 1 | Slide 2 | Slide 3',
    defaultProps: {
      slides: [
        {
          title: 'Summer Collection',
          subtitle: 'New arrivals for the season',
          cta: 'Shop Summer',
          background: '/images/summer.jpg'
        },
        {
          title: 'Special Offers',
          subtitle: 'Up to 50% off selected items',
          cta: 'View Deals',
          background: '/images/deals.jpg'
        },
        {
          title: 'Free Shipping',
          subtitle: 'On orders over $100',
          cta: 'Shop Now',
          background: '/images/shipping.jpg'
        }
      ],
      autoplay: true,
      showDots: true,
      showArrows: true
    },
    defaultStyles: {
      height: '70vh',
      position: 'relative'
    }
  },

  // ===== PRODUCT COMPONENTS =====
  {
    id: 'product-grid-modern',
    name: 'Modern Product Grid',
    type: 'product-grid',
    category: 'product',
    description: 'Modern grid layout with hover effects',
    icon: '🛍️',
    preview: '3x2 Grid | Hover Effects',
    defaultProps: {
      title: 'Featured Products',
      columns: 3,
      rows: 2,
      showPrices: true,
      showRatings: true,
      showQuickView: true,
      showWishlist: true,
      hoverEffect: 'zoom'
    },
    defaultStyles: {
      padding: '3rem 0',
      backgroundColor: '#ffffff'
    }
  },

  {
    id: 'product-carousel',
    name: 'Product Carousel',
    type: 'product',
    category: 'product',
    description: 'Horizontal scrolling product showcase',
    icon: '🎡',
    preview: 'Scrollable | Multiple Items',
    defaultProps: {
      title: 'Trending Now',
      itemsToShow: 4,
      showArrows: true,
      showDots: false,
      autoplay: false,
      showPrices: true,
      showRatings: true
    },
    defaultStyles: {
      padding: '2rem 0',
      backgroundColor: '#f8fafc'
    }
  },

  {
    id: 'product-showcase',
    name: 'Product Showcase',
    type: 'product',
    category: 'product',
    description: 'Featured product with detailed presentation',
    icon: '⭐',
    preview: 'Large Image | Details | Actions',
    defaultProps: {
      productName: 'Premium Product',
      description: 'High-quality product with excellent features',
      price: '$299.99',
      originalPrice: '$399.99',
      images: ['/product1.jpg', '/product2.jpg'],
      showGallery: true,
      showReviews: true,
      showSpecs: true
    },
    defaultStyles: {
      padding: '3rem 0',
      backgroundColor: '#ffffff'
    }
  },

  // ===== CONTENT COMPONENTS =====
  {
    id: 'features-grid',
    name: 'Features Grid',
    type: 'features',
    category: 'content',
    description: 'Grid of features or benefits',
    icon: '✨',
    preview: 'Icon | Title | Description',
    defaultProps: {
      title: 'Why Choose Us',
      features: [
        {
          icon: '🚚',
          title: 'Free Shipping',
          description: 'Free delivery on orders over $50'
        },
        {
          icon: '💯',
          title: 'Quality Guarantee',
          description: '100% satisfaction guaranteed'
        },
        {
          icon: '🔒',
          title: 'Secure Payment',
          description: 'Your payment is safe and secure'
        },
        {
          icon: '🎧',
          title: '24/7 Support',
          description: 'Round-the-clock customer service'
        }
      ],
      columns: 4
    },
    defaultStyles: {
      padding: '3rem 0',
      backgroundColor: '#f9fafb'
    }
  },

  {
    id: 'testimonials-carousel',
    name: 'Customer Testimonials',
    type: 'testimonial',
    category: 'marketing',
    description: 'Customer reviews and testimonials',
    icon: '💬',
    preview: 'Quote | Author | Rating',
    defaultProps: {
      title: 'What Our Customers Say',
      testimonials: [
        {
          quote: 'Amazing quality and fast shipping!',
          author: 'Sarah Johnson',
          rating: 5,
          avatar: '/avatars/sarah.jpg'
        },
        {
          quote: 'Best customer service I have experienced.',
          author: 'Mike Chen',
          rating: 5,
          avatar: '/avatars/mike.jpg'
        },
        {
          quote: 'Products exceeded my expectations.',
          author: 'Emma Davis',
          rating: 5,
          avatar: '/avatars/emma.jpg'
        }
      ],
      showRatings: true,
      autoplay: true
    },
    defaultStyles: {
      padding: '3rem 0',
      backgroundColor: '#ffffff'
    }
  },

  {
    id: 'newsletter-signup',
    name: 'Newsletter Signup',
    type: 'newsletter',
    category: 'marketing',
    description: 'Email subscription form',
    icon: '📧',
    preview: 'Email Input | Subscribe Button',
    defaultProps: {
      title: 'Stay Updated',
      subtitle: 'Get the latest news and exclusive offers',
      placeholder: 'Enter your email address',
      buttonText: 'Subscribe',
      showPrivacyNote: true,
      privacyNote: 'We respect your privacy. Unsubscribe anytime.'
    },
    defaultStyles: {
      padding: '2rem 0',
      backgroundColor: '#f3f4f6',
      textAlign: 'center'
    }
  },

  {
    id: 'blog-preview',
    name: 'Blog Preview',
    type: 'blog',
    category: 'content',
    description: 'Latest blog posts preview',
    icon: '📝',
    preview: 'Article Cards | Read More',
    defaultProps: {
      title: 'Latest Articles',
      postsToShow: 3,
      showExcerpt: true,
      showDate: true,
      showAuthor: true,
      showReadTime: true
    },
    defaultStyles: {
      padding: '3rem 0',
      backgroundColor: '#ffffff'
    }
  },

  // ===== LAYOUT COMPONENTS =====
  {
    id: 'modern-footer',
    name: 'Modern Footer',
    type: 'footer',
    category: 'layout',
    description: 'Comprehensive footer with links and info',
    icon: '🦶',
    preview: 'Links | Social | Contact',
    defaultProps: {
      companyName: 'Your Store',
      description: 'Your trusted partner for quality products',
      sections: [
        {
          title: 'Shop',
          links: ['All Products', 'New Arrivals', 'Best Sellers', 'Sale']
        },
        {
          title: 'Support',
          links: ['Contact Us', 'FAQ', 'Shipping', 'Returns']
        },
        {
          title: 'Company',
          links: ['About', 'Careers', 'Press', 'Blog']
        }
      ],
      social: [
        { platform: 'Facebook', url: '#' },
        { platform: 'Twitter', url: '#' },
        { platform: 'Instagram', url: '#' }
      ],
      showNewsletter: true,
      showPaymentMethods: true
    },
    defaultStyles: {
      backgroundColor: '#1f2937',
      textColor: '#f9fafb',
      padding: '3rem 0 1rem'
    }
  },

  {
    id: 'breadcrumb-nav',
    name: 'Breadcrumb Navigation',
    type: 'breadcrumb',
    category: 'layout',
    description: 'Navigation breadcrumb trail',
    icon: '🧭',
    preview: 'Home > Category > Product',
    defaultProps: {
      items: [
        { label: 'Home', url: '/' },
        { label: 'Category', url: '/category' },
        { label: 'Current Page', url: '#' }
      ],
      separator: '>',
      showHome: true
    },
    defaultStyles: {
      padding: '1rem 0',
      backgroundColor: '#f9fafb',
      fontSize: '0.875rem'
    }
  },

  // ===== PROMOTIONAL COMPONENTS =====
  {
    id: 'announcement-bar',
    name: 'Announcement Bar',
    type: 'announcement',
    category: 'marketing',
    description: 'Top announcement or promotional bar',
    icon: '📢',
    preview: 'Free Shipping on Orders $50+',
    defaultProps: {
      message: 'Free shipping on orders over $50! Limited time offer.',
      showCloseButton: true,
      link: '/shipping',
      linkText: 'Learn More',
      isScrolling: false
    },
    defaultStyles: {
      backgroundColor: '#3b82f6',
      textColor: '#ffffff',
      padding: '0.5rem 0',
      textAlign: 'center'
    }
  },

  {
    id: 'countdown-banner',
    name: 'Countdown Banner',
    type: 'banner',
    category: 'marketing',
    description: 'Sale countdown timer banner',
    icon: '⏰',
    preview: 'Sale Ends In: 2d 4h 32m',
    defaultProps: {
      title: 'Flash Sale',
      subtitle: 'Up to 50% off selected items',
      endDate: '2024-12-31T23:59:59',
      showTimer: true,
      ctaText: 'Shop Sale',
      ctaLink: '/sale'
    },
    defaultStyles: {
      backgroundColor: '#dc2626',
      textColor: '#ffffff',
      padding: '2rem 0',
      textAlign: 'center'
    }
  },

  {
    id: 'category-showcase',
    name: 'Category Showcase',
    type: 'category',
    category: 'product',
    description: 'Featured product categories',
    icon: '🏷️',
    preview: 'Category Grid | Images | Links',
    defaultProps: {
      title: 'Shop by Category',
      categories: [
        {
          name: 'Electronics',
          image: '/categories/electronics.jpg',
          url: '/category/electronics'
        },
        {
          name: 'Fashion',
          image: '/categories/fashion.jpg',
          url: '/category/fashion'
        },
        {
          name: 'Home & Garden',
          image: '/categories/home.jpg',
          url: '/category/home'
        },
        {
          name: 'Sports',
          image: '/categories/sports.jpg',
          url: '/category/sports'
        }
      ],
      layout: 'grid',
      columns: 4
    },
    defaultStyles: {
      padding: '3rem 0',
      backgroundColor: '#ffffff'
    }
  },

  {
    id: 'trust-badges',
    name: 'Trust Badges',
    type: 'trust',
    category: 'marketing',
    description: 'Security and trust indicators',
    icon: '🛡️',
    preview: 'SSL | Guarantee | Reviews',
    defaultProps: {
      badges: [
        {
          icon: '🔒',
          title: 'SSL Secured',
          description: '256-bit encryption'
        },
        {
          icon: '✅',
          title: 'Money Back Guarantee',
          description: '30-day returns'
        },
        {
          icon: '⭐',
          title: '4.9/5 Rating',
          description: '10,000+ reviews'
        },
        {
          icon: '🚚',
          title: 'Fast Shipping',
          description: '2-day delivery'
        }
      ],
      layout: 'horizontal'
    },
    defaultStyles: {
      padding: '2rem 0',
      backgroundColor: '#f9fafb',
      textAlign: 'center'
    }
  },

  // === E-COMMERCE PRODUCT COMPONENTS ===
  {
    id: 'product-card-basic',
    name: 'Basic Product Card',
    type: 'product-card',
    category: 'product',
    description: 'Simple product card with image, title, price, and add to cart',
    icon: '🛍️',
    preview: 'Image | Title | $Price | Add to Cart',
    defaultProps: {
      productName: 'Premium T-Shirt',
      price: 29.99,
      originalPrice: 39.99,
      rating: 4.5,
      reviewCount: 128,
      productImage: '/images/product-1.jpg',
      badges: ['Sale', 'Popular'],
      inStock: true,
      quickView: true
    },
    defaultStyles: {
      border: '1px solid #e2e8f0',
      borderRadius: '0.5rem',
      padding: '1rem',
      backgroundColor: '#ffffff'
    }
  },

  {
    id: 'product-card-hover',
    name: 'Hover Effect Product Card',
    type: 'product-card',
    category: 'product',
    description: 'Product card with hover animations and secondary image',
    icon: '✨',
    preview: 'Hover Effects | Secondary Image | Quick Actions',
    defaultProps: {
      productName: 'Designer Sneakers',
      price: 149.99,
      originalPrice: null,
      rating: 4.8,
      reviewCount: 89,
      primaryImage: '/images/product-primary.jpg',
      secondaryImage: '/images/product-secondary.jpg',
      colorVariants: ['Black', 'White', 'Gray'],
      sizeVariants: ['38', '39', '40', '41', '42'],
      quickActions: ['wishlist', 'compare', 'quickView']
    },
    defaultStyles: {
      border: '1px solid #e2e8f0',
      borderRadius: '0.75rem',
      overflow: 'hidden',
      transition: 'all 0.3s ease'
    }
  },

  {
    id: 'product-grid-featured',
    name: 'Featured Products Grid',
    type: 'product-grid',
    category: 'product',
    description: 'Grid layout for featured/bestselling products',
    icon: '🌟',
    preview: 'Featured Grid | 4 Columns | Filters',
    defaultProps: {
      title: 'Featured Products',
      subtitle: 'Handpicked items just for you',
      gridColumns: 4,
      showFilters: true,
      showSortBy: true,
      filterOptions: ['Category', 'Price', 'Brand', 'Rating'],
      sortOptions: ['Popularity', 'Price: Low to High', 'Price: High to Low', 'Newest'],
      productsPerPage: 12,
      showPagination: true
    },
    defaultStyles: {
      padding: '3rem 0',
      backgroundColor: '#ffffff'
    }
  },

  {
    id: 'product-comparison',
    name: 'Product Comparison Table',
    type: 'comparison',
    category: 'product',
    description: 'Side-by-side product comparison',
    icon: '⚖️',
    preview: 'Compare Features | Specs | Prices',
    defaultProps: {
      title: 'Compare Products',
      maxProducts: 3,
      comparisonFeatures: ['Price', 'Rating', 'Material', 'Warranty', 'Colors Available'],
      showDifferencesOnly: false,
      enablePrint: true
    },
    defaultStyles: {
      padding: '2rem 0',
      backgroundColor: '#f9fafb'
    }
  },

  // === E-COMMERCE CART & CHECKOUT ===
  {
    id: 'mini-cart',
    name: 'Mini Cart Dropdown',
    type: 'cart',
    category: 'cart',
    description: 'Compact cart preview with quick actions',
    icon: '🛒',
    preview: 'Cart Items | Subtotal | Checkout',
    defaultProps: {
      showProductImages: true,
      showQuantityControls: true,
      showRemoveButton: true,
      showSubtotal: true,
      showTax: true,
      showShipping: true,
      quickCheckout: true,
      freeShippingThreshold: 50,
      maxItems: 5
    },
    defaultStyles: {
      width: '400px',
      maxHeight: '500px',
      backgroundColor: '#ffffff',
      border: '1px solid #e2e8f0',
      borderRadius: '0.5rem',
      boxShadow: '0 10px 25px rgba(0,0,0,0.15)'
    }
  },

  {
    id: 'cart-page',
    name: 'Full Cart Page',
    type: 'cart',
    category: 'cart',
    description: 'Complete cart page with all functionality',
    icon: '🛍️',
    preview: 'Full Cart | Quantity | Coupons | Totals',
    defaultProps: {
      showProductDetails: true,
      showQuantityControls: true,
      showRemoveButton: true,
      showSaveForLater: true,
      showRecommendations: true,
      showCouponField: true,
      showEstimatedDelivery: true,
      showSecurityBadges: true,
      enableGuestCheckout: true
    },
    defaultStyles: {
      padding: '2rem 0',
      backgroundColor: '#ffffff'
    }
  },

  {
    id: 'checkout-steps',
    name: 'Multi-Step Checkout',
    type: 'checkout',
    category: 'checkout',
    description: 'Step-by-step checkout process',
    icon: '📋',
    preview: 'Step 1-4 | Progress | Validation',
    defaultProps: {
      steps: ['Cart', 'Shipping', 'Payment', 'Confirmation'],
      showProgressBar: true,
      showStepNumbers: true,
      allowStepSkipping: false,
      showOrderSummary: true,
      showSecurityInfo: true,
      enableGuestCheckout: true,
      saveProgress: true
    },
    defaultStyles: {
      padding: '2rem 0',
      backgroundColor: '#f9fafb'
    }
  },

  {
    id: 'one-page-checkout',
    name: 'One Page Checkout',
    type: 'checkout',
    category: 'checkout',
    description: 'Single page checkout with all fields',
    icon: '⚡',
    preview: 'All-in-One | Quick Checkout | Express Pay',
    defaultProps: {
      showOrderSummary: true,
      showExpressPayment: true,
      expressPaymentMethods: ['PayPal', 'Apple Pay', 'Google Pay'],
      showGuestOption: true,
      showCreateAccount: true,
      showNewsletterSignup: true,
      autoFillAddress: true,
      showTrustSignals: true
    },
    defaultStyles: {
      padding: '2rem 0',
      backgroundColor: '#ffffff'
    }
  },

  // === E-COMMERCE PROMOTIONAL ===
  {
    id: 'sale-banner',
    name: 'Sale Banner',
    type: 'promotion',
    category: 'promotional',
    description: 'Eye-catching sale/discount banner',
    icon: '🏷️',
    preview: 'SALE 50% OFF | Limited Time | Shop Now',
    defaultProps: {
      bannerText: 'SUMMER SALE - UP TO 50% OFF',
      discountPercent: 50,
      validUntil: '2024-08-31',
      showCountdown: true,
      ctaText: 'Shop Now',
      backgroundColor: '#ef4444',
      textColor: '#ffffff',
      position: 'top',
      dismissible: true
    },
    defaultStyles: {
      padding: '1rem',
      textAlign: 'center',
      fontWeight: 'bold'
    }
  },

  {
    id: 'flash-sale',
    name: 'Flash Sale Timer',
    type: 'promotion',
    category: 'promotional',
    description: 'Countdown timer for flash sales',
    icon: '⏰',
    preview: 'Flash Sale | 2h 30m 45s | Limited Stock',
    defaultProps: {
      title: 'Flash Sale',
      description: 'Limited time offer - Don\'t miss out!',
      endTime: '2024-12-31T23:59:59',
      showDays: true,
      showHours: true,
      showMinutes: true,
      showSeconds: true,
      urgencyMessages: ['Only 3 items left!', 'Sale ends soon!'],
      productLimit: 100,
      showProgressBar: true
    },
    defaultStyles: {
      padding: '2rem',
      backgroundColor: '#fef2f2',
      border: '2px solid #ef4444',
      borderRadius: '0.5rem',
      textAlign: 'center'
    }
  },

  {
    id: 'coupon-popup',
    name: 'Discount Coupon Popup',
    type: 'promotion',
    category: 'promotional',
    description: 'Exit-intent or timed discount popup',
    icon: '🎟️',
    preview: 'Get 10% OFF | Email Signup | Exclusive Deal',
    defaultProps: {
      title: 'Get 10% OFF Your First Order',
      description: 'Join our newsletter for exclusive deals and updates',
      discountCode: 'WELCOME10',
      discountPercent: 10,
      minOrderValue: 50,
      showEmailField: true,
      showPhoneField: false,
      triggerType: 'exit-intent',
      triggerDelay: 30000,
      showOnce: true
    },
    defaultStyles: {
      maxWidth: '500px',
      backgroundColor: '#ffffff',
      border: '2px solid #10b981',
      borderRadius: '1rem',
      boxShadow: '0 25px 50px rgba(0,0,0,0.25)'
    }
  },

  // === E-COMMERCE REVIEW & TESTIMONIALS ===
  {
    id: 'product-reviews',
    name: 'Product Reviews Section',
    type: 'reviews',
    category: 'social-proof',
    description: 'Customer reviews with ratings and photos',
    icon: '⭐',
    preview: 'Reviews | Ratings | Photos | Helpful',
    defaultProps: {
      showOverallRating: true,
      showRatingBreakdown: true,
      showReviewPhotos: true,
      showVerifiedPurchase: true,
      showHelpfulVotes: true,
      allowSorting: true,
      allowFiltering: true,
      sortOptions: ['Most Recent', 'Most Helpful', 'Highest Rating', 'Lowest Rating'],
      filterOptions: ['5 Stars', '4 Stars', '3 Stars', '2 Stars', '1 Star', 'With Photos'],
      reviewsPerPage: 10
    },
    defaultStyles: {
      padding: '2rem 0',
      backgroundColor: '#ffffff'
    }
  },

  {
    id: 'review-carousel',
    name: 'Customer Review Carousel',
    type: 'testimonials',
    category: 'social-proof',
    description: 'Rotating customer testimonials with photos',
    icon: '💬',
    preview: 'Customer Photos | Testimonials | Auto-rotate',
    defaultProps: {
      showCustomerPhotos: true,
      showCustomerNames: true,
      showCustomerLocation: true,
      showRatings: true,
      showProductPurchased: true,
      autoRotate: true,
      rotationInterval: 5000,
      showNavigation: true,
      showDots: true,
      itemsPerSlide: 3
    },
    defaultStyles: {
      padding: '3rem 0',
      backgroundColor: '#f9fafb'
    }
  },

  // === E-COMMERCE SHIPPING & RETURNS ===
  {
    id: 'shipping-calculator',
    name: 'Shipping Calculator',
    type: 'shipping',
    category: 'utility',
    description: 'Calculate shipping costs by location',
    icon: '📦',
    preview: 'Enter ZIP | Calculate | Delivery Options',
    defaultProps: {
      showDeliveryOptions: true,
      showEstimatedDays: true,
      showTrackingInfo: true,
      freeShippingThreshold: 50,
      expeditedOptions: ['Standard', 'Express', 'Overnight'],
      allowPickup: true,
      showInsurance: true,
      showSignatureRequired: false
    },
    defaultStyles: {
      padding: '1.5rem',
      backgroundColor: '#f8fafc',
      border: '1px solid #e2e8f0',
      borderRadius: '0.5rem'
    }
  },

  {
    id: 'return-policy',
    name: 'Return Policy Widget',
    type: 'policy',
    category: 'utility',
    description: 'Clear return policy information',
    icon: '↩️',
    preview: '30-Day Returns | Free Returns | Easy Process',
    defaultProps: {
      returnPeriod: 30,
      freeReturns: true,
      returnConditions: ['Unused', 'Original packaging', 'With tags'],
      showReturnProcess: true,
      showReturnForm: true,
      showFAQ: true,
      contactInfo: true,
      showGuarantees: true
    },
    defaultStyles: {
      padding: '2rem',
      backgroundColor: '#ffffff',
      border: '1px solid #e2e8f0',
      borderRadius: '0.5rem'
    }
  },

  // === E-COMMERCE SEARCH & FILTERS ===
  {
    id: 'advanced-search',
    name: 'Advanced Product Search',
    type: 'search',
    category: 'navigation',
    description: 'Search with filters and suggestions',
    icon: '🔍',
    preview: 'Search Bar | Filters | Autocomplete | Results',
    defaultProps: {
      showAutocomplete: true,
      showSearchSuggestions: true,
      showRecentSearches: true,
      showPopularSearches: true,
      showCategoryFilter: true,
      showPriceFilter: true,
      showBrandFilter: true,
      showRatingFilter: true,
      showAvailabilityFilter: true,
      voiceSearch: true
    },
    defaultStyles: {
      padding: '1rem 0',
      backgroundColor: '#ffffff'
    }
  },

  {
    id: 'product-filter-sidebar',
    name: 'Product Filter Sidebar',
    type: 'filter',
    category: 'navigation',
    description: 'Comprehensive product filtering options',
    icon: '🎛️',
    preview: 'Categories | Price | Brand | Rating | More',
    defaultProps: {
      showCategories: true,
      showPriceRange: true,
      showBrands: true,
      showRatings: true,
      showAvailability: true,
      showColors: true,
      showSizes: true,
      showMaterials: true,
      showDiscounts: true,
      collapsibleSections: true,
      showProductCount: true,
      showClearFilters: true
    },
    defaultStyles: {
      width: '280px',
      padding: '1.5rem',
      backgroundColor: '#f9fafb',
      borderRight: '1px solid #e2e8f0'
    }
  },

  // === E-COMMERCE WISHLIST & FAVORITES ===
  {
    id: 'wishlist-grid',
    name: 'Wishlist Grid',
    type: 'wishlist',
    category: 'user-account',
    description: 'Grid layout for saved products',
    icon: '❤️',
    preview: 'Saved Items | Move to Cart | Share | Remove',
    defaultProps: {
      showProductImages: true,
      showProductDetails: true,
      showPriceChanges: true,
      showStockStatus: true,
      showMoveToCart: true,
      showRemoveButton: true,
      showShareWishlist: true,
      showWishlistName: true,
      allowMultipleWishlists: true,
      showDateAdded: true
    },
    defaultStyles: {
      padding: '2rem 0',
      backgroundColor: '#ffffff'
    }
  },

  {
    id: 'recently-viewed',
    name: 'Recently Viewed Products',
    type: 'recent',
    category: 'user-account',
    description: 'Horizontal scroll of recently viewed items',
    icon: '👁️',
    preview: 'Recently Viewed | Horizontal Scroll | Quick Add',
    defaultProps: {
      maxItems: 10,
      showScrollArrows: true,
      showProductNames: true,
      showPrices: true,
      showRatings: true,
      showQuickAdd: true,
      autoScroll: false,
      showClearHistory: true,
      itemsPerView: 5
    },
    defaultStyles: {
      padding: '2rem 0',
      backgroundColor: '#f9fafb'
    }
  }
];

// Sample layouts for different themes
export const sampleLayouts = {
  'modern-minimalist': [
    {
      id: 'header-1',
      type: 'header',
      content: {
        title: 'Modern Store',
        navigation: ['Shop', 'About', 'Contact']
      }
    },
    {
      id: 'hero-1',
      type: 'hero',
      content: {
        title: 'Premium Quality Products',
        subtitle: 'Discover our curated collection of modern essentials',
        buttonText: 'Shop Now'
      }
    },
    {
      id: 'features-1',
      type: 'features',
      content: {
        title: 'Why Choose Us'
      }
    },
    {
      id: 'footer-1',
      type: 'footer',
      content: {
        companyName: 'Modern Store',
        links: ['Support', 'Privacy', 'Terms'],
        social: ['Twitter', 'Instagram', 'Facebook']
      }
    }
  ],

  'dark-modern': [
    {
      id: 'header-2',
      type: 'header',
      content: {
        title: 'Dark Store',
        navigation: ['Products', 'Collections', 'About']
      }
    },
    {
      id: 'hero-2',
      type: 'hero',
      content: {
        title: 'Experience the Dark Side',
        subtitle: 'Premium products for the sophisticated customer',
        buttonText: 'Explore Collection'
      }
    },
    {
      id: 'products-2',
      type: 'products',
      content: {
        title: 'Featured Collection'
      }
    },
    {
      id: 'footer-2',
      type: 'footer',
      content: {
        companyName: 'Dark Store'
      }
    }
  ],

  'vibrant-creative': [
    {
      id: 'header-3',
      type: 'header',
      content: {
        title: 'Creative Hub',
        navigation: ['Gallery', 'Shop', 'Artists', 'About']
      }
    },
    {
      id: 'hero-3',
      type: 'hero',
      content: {
        title: 'Unleash Your Creativity',
        subtitle: 'Unique products from independent artists',
        buttonText: 'Discover Art'
      }
    },
    {
      id: 'category-3',
      type: 'category',
      content: {
        title: 'Browse Collections'
      }
    },
    {
      id: 'testimonials-3',
      type: 'testimonial',
      content: {
        title: 'Artist Spotlights'
      }
    },
    {
      id: 'footer-3',
      type: 'footer',
      content: {
        companyName: 'Creative Hub'
      }
    }
  ]
};

// Sample products for different themes
export const sampleProducts = {
  'modern-minimalist': [
    {
      id: 1,
      name: 'Minimalist Watch',
      price: '$199',
      image: '/products/watch.jpg'
    },
    {
      id: 2,
      name: 'Clean Desk Lamp',
      price: '$89',
      image: '/products/lamp.jpg'
    },
    {
      id: 3,
      name: 'Simple Notebook',
      price: '$29',
      image: '/products/notebook.jpg'
    }
  ],

  'dark-modern': [
    {
      id: 1,
      name: 'Dark Edition Headphones',
      price: '$299',
      image: '/products/headphones.jpg'
    },
    {
      id: 2,
      name: 'Black Steel Bottle',
      price: '$45',
      image: '/products/bottle.jpg'
    },
    {
      id: 3,
      name: 'Night Mode Keyboard',
      price: '$149',
      image: '/products/keyboard.jpg'
    }
  ],

  'vibrant-creative': [
    {
      id: 1,
      name: 'Colorful Art Print',
      price: '$79',
      image: '/products/print.jpg'
    },
    {
      id: 2,
      name: 'Rainbow Mug Set',
      price: '$39',
      image: '/products/mugs.jpg'
    },
    {
      id: 3,
      name: 'Artistic Phone Case',
      price: '$25',
      image: '/products/case.jpg'
    }
  ]
};

// Quick start templates for different business types
export const quickStartTemplates = [
  {
    id: 'fashion-store',
    name: 'Fashion Store',
    description: 'Complete setup for fashion and apparel',
    icon: '👗',
    category: 'fashion',
    theme: 'vibrant-creative',
    components: [
      {
        type: 'header',
        title: 'Fashion Forward',
        menuItems: ['Women', 'Men', 'Accessories', 'Sale']
      },
      {
        type: 'hero',
        title: 'New Collection',
        subtitle: 'Discover the latest trends'
      },
      {
        type: 'category',
        title: 'Shop by Category'
      },
      {
        type: 'product-grid',
        title: 'Featured Products'
      },
      {
        type: 'newsletter',
        title: 'Style Updates'
      },
      {
        type: 'footer',
        companyName: 'Fashion Forward'
      }
    ]
  },

  {
    id: 'tech-store',
    name: 'Tech Store',
    description: 'Electronics and gadgets store',
    icon: '💻',
    category: 'technology',
    theme: 'dark-modern',
    components: [
      {
        type: 'header',
        title: 'Tech Hub',
        menuItems: ['Computers', 'Phones', 'Accessories', 'Support']
      },
      {
        type: 'hero',
        title: 'Latest Technology',
        subtitle: 'Cutting-edge gadgets and devices'
      },
      {
        type: 'features',
        title: 'Why Choose Tech Hub'
      },
      {
        type: 'product-grid',
        title: 'Best Sellers'
      },
      {
        type: 'testimonials',
        title: 'Customer Reviews'
      },
      {
        type: 'footer',
        companyName: 'Tech Hub'
      }
    ]
  },

  {
    id: 'home-decor',
    name: 'Home Decor',
    description: 'Furniture and home accessories',
    icon: '🏠',
    category: 'home',
    theme: 'nature-organic',
    components: [
      {
        type: 'header',
        title: 'Home Haven',
        menuItems: ['Furniture', 'Decor', 'Lighting', 'Inspiration']
      },
      {
        type: 'hero',
        title: 'Beautiful Homes',
        subtitle: 'Transform your space with our curated collection'
      },
      {
        type: 'category',
        title: 'Room by Room'
      },
      {
        type: 'product-showcase',
        title: 'Featured Furniture'
      },
      {
        type: 'trust-badges',
        title: 'Our Promise'
      },
      {
        type: 'footer',
        companyName: 'Home Haven'
      }
    ]
  }
];

// Helper function to get templates by category
export const getTemplatesByCategory = (category) => {
  return componentTemplates.filter(template => template.category === category);
};

// Helper function to get theme by id
export const getThemeById = (themeId) => {
  return themeTemplates[themeId];
};

// Helper function to get quick start template
export const getQuickStartTemplate = (templateId) => {
  return quickStartTemplates.find(template => template.id === templateId);
};