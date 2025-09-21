import { atom } from 'jotai';

// Core shop customizer state atoms
export const componentsAtom = atom([]);
export const selectedComponentAtom = atom(null);
export const selectedThemeAtom = atom(null);
export const previewModeAtom = atom(false);
export const showThemeCustomizerAtom = atom(false);

// Derived atoms for computed values
export const componentsCountAtom = atom((get) => get(componentsAtom).length);

// Theme-related atoms
export const themeColorsAtom = atom((get) => {
  const theme = get(selectedThemeAtom);
  return theme?.colors || {};
});

export const themeTypographyAtom = atom((get) => {
  const theme = get(selectedThemeAtom);
  return theme?.typography || {};
});

// Customizer settings atoms
export const customizerSettingsAtom = atom({
  autoSave: true,
  showGrid: false,
  snapToGrid: true,
  gridSize: 8,
  devicePreview: 'desktop' // desktop, tablet, mobile
});

// History management atoms
export const historyAtom = atom({
  past: [],
  present: null,
  future: []
});

// Save/Load atoms
export const saveStatusAtom = atom('idle'); // idle, saving, saved, error
export const lastSavedAtom = atom(null);
export const designNameAtom = atom('Untitled Design');

// Actions atoms (write-only atoms for complex operations)
export const addComponentAtom = atom(
  null,
  (get, set, component) => {
    const currentComponents = get(componentsAtom);
    set(componentsAtom, [...currentComponents, component]);
  }
);

export const insertComponentAtom = atom(
  null,
  (get, set, { component, index }) => {
    const currentComponents = get(componentsAtom);
    const newComponents = [...currentComponents];
    newComponents.splice(index, 0, component);
    set(componentsAtom, newComponents);
  }
);

export const updateComponentAtom = atom(
  null,
  (get, set, { componentId, updates }) => {
    const currentComponents = get(componentsAtom);
    const updatedComponents = currentComponents.map(comp =>
      comp.id === componentId ? { ...comp, ...updates } : comp
    );
    set(componentsAtom, updatedComponents);

    // Update selected component if it's the one being updated
    const selectedComponent = get(selectedComponentAtom);
    if (selectedComponent?.id === componentId) {
      set(selectedComponentAtom, { ...selectedComponent, ...updates });
    }
  }
);

export const deleteComponentAtom = atom(
  null,
  (get, set, componentId) => {
    const currentComponents = get(componentsAtom);
    const filteredComponents = currentComponents.filter(comp => comp.id !== componentId);
    set(componentsAtom, filteredComponents);

    // Clear selected component if it was deleted
    const selectedComponent = get(selectedComponentAtom);
    if (selectedComponent?.id === componentId) {
      set(selectedComponentAtom, null);
    }
  }
);

export const reorderComponentsAtom = atom(
  null,
  (get, set, { oldIndex, newIndex }) => {
    const currentComponents = get(componentsAtom);
    const reorderedComponents = [...currentComponents];
    const [removed] = reorderedComponents.splice(oldIndex, 1);
    reorderedComponents.splice(newIndex, 0, removed);
    set(componentsAtom, reorderedComponents);
  }
);

export const applyThemeAtom = atom(
  null,
  (get, set, theme) => {
    set(selectedThemeAtom, theme);
    set(showThemeCustomizerAtom, true);

    // Apply theme to existing components
    const currentComponents = get(componentsAtom);
    if (currentComponents.length > 0) {
      const themedComponents = currentComponents.map(component => ({
        ...component,
        theme: theme.id,
        styles: {
          ...component.styles,
          ...theme.components?.[component.type],
          backgroundColor: theme.colors?.background || component.styles?.backgroundColor,
          textColor: theme.colors?.text || component.styles?.textColor
        }
      }));
      set(componentsAtom, themedComponents);
    }
  }
);

export const clearBuilderAtom = atom(
  null,
  (get, set) => {
    set(componentsAtom, []);
    set(selectedComponentAtom, null);
    set(selectedThemeAtom, null);
    set(showThemeCustomizerAtom, false);
    set(historyAtom, { past: [], present: null, future: [] });
  }
);

// Undo/Redo functionality
export const undoAtom = atom(
  null,
  (get, set) => {
    const history = get(historyAtom);
    if (history.past.length === 0) return;

    const previous = history.past[history.past.length - 1];
    const newPast = history.past.slice(0, history.past.length - 1);
    const newFuture = [get(componentsAtom), ...history.future];

    set(historyAtom, {
      past: newPast,
      present: previous,
      future: newFuture
    });
    set(componentsAtom, previous);
  }
);

export const redoAtom = atom(
  null,
  (get, set) => {
    const history = get(historyAtom);
    if (history.future.length === 0) return;

    const next = history.future[0];
    const newFuture = history.future.slice(1);
    const newPast = [...history.past, get(componentsAtom)];

    set(historyAtom, {
      past: newPast,
      present: next,
      future: newFuture
    });
    set(componentsAtom, next);
  }
);

// Save design atom
export const saveDesignAtom = atom(
  null,
  async (get, set, { name, description } = {}) => {
    set(saveStatusAtom, 'saving');

    try {
      const designData = {
        id: Date.now().toString(),
        name: name || get(designNameAtom),
        description: description || '',
        components: get(componentsAtom),
        theme: get(selectedThemeAtom),
        settings: get(builderSettingsAtom),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      };

      // Store in localStorage for now
      const savedDesigns = JSON.parse(localStorage.getItem('store-designs') || '[]');
      const existingIndex = savedDesigns.findIndex(d => d.name === designData.name);

      if (existingIndex >= 0) {
        savedDesigns[existingIndex] = { ...savedDesigns[existingIndex], ...designData, updatedAt: designData.updatedAt };
      } else {
        savedDesigns.push(designData);
      }

      localStorage.setItem('store-designs', JSON.stringify(savedDesigns));

      set(saveStatusAtom, 'saved');
      set(lastSavedAtom, new Date());
      set(designNameAtom, designData.name);

      // Auto-reset status after 3 seconds
      setTimeout(() => set(saveStatusAtom, 'idle'), 3000);

      return designData;
    } catch (error) {
      console.error('Save failed:', error);
      set(saveStatusAtom, 'error');
      setTimeout(() => set(saveStatusAtom, 'idle'), 3000);
      throw error;
    }
  }
);

// Load design atom
export const loadDesignAtom = atom(
  null,
  (get, set, designData) => {
    set(componentsAtom, designData.components || []);
    set(selectedThemeAtom, designData.theme || null);
    set(builderSettingsAtom, { ...get(builderSettingsAtom), ...designData.settings });
    set(designNameAtom, designData.name || 'Untitled Design');
    set(selectedComponentAtom, null);
    set(historyAtom, { past: [], present: null, future: [] });
  }
);

// Get saved designs atom
export const getSavedDesignsAtom = atom(
  () => {
    try {
      return JSON.parse(localStorage.getItem('store-designs') || '[]');
    } catch {
      return [];
    }
  }
);

// Export for Next.js storefront atom
export const exportForNextJSAtom = atom(
  (get) => {
    const components = get(componentsAtom);
    const theme = get(selectedThemeAtom);
    const settings = get(builderSettingsAtom);
    const designName = get(designNameAtom);

    // Transform components into Next.js-friendly structure
    const transformedComponents = components.map(component => {
      const baseComponent = {
        id: component.id,
        type: component.type,
        order: component.order || 0,
        visible: component.visible !== false,
        content: {},
        styles: {
          backgroundColor: component.styles?.backgroundColor || theme?.colors?.background || '#ffffff',
          textColor: component.styles?.textColor || theme?.colors?.text || '#000000',
          fontSize: component.styles?.fontSize || theme?.typography?.fontSize || '16px',
          fontFamily: component.styles?.fontFamily || theme?.typography?.fontFamily || 'Inter, sans-serif',
          padding: component.styles?.padding || '1rem',
          margin: component.styles?.margin || '0'
        }
      };

      // Extract content based on component type
      switch (component.type) {
        case 'header':
          const menuItems = component.props?.menuItems || component.content?.menuItems || ['Home', 'Products', 'About', 'Contact'];
          baseComponent.content = {
            logo: component.props?.title || component.content?.logo || 'Your Store',
            navigation: menuItems.map(item => 
              typeof item === 'object' ? item : { label: item, href: `/${item.toLowerCase()}` }
            ),
            ctaText: component.props?.ctaText || component.content?.ctaText || 'Get Started',
            showCart: component.content?.showCart !== false,
            showSearch: component.content?.showSearch !== false
          };
          break;

        case 'hero':
          baseComponent.content = {
            title: component.props?.title || component.content?.title || 'Welcome to Our Store',
            subtitle: component.props?.subtitle || component.content?.subtitle || 'Discover amazing products',
            backgroundImage: component.props?.backgroundImage || component.content?.backgroundImage || '',
            primaryCta: component.props?.primaryCta || component.content?.primaryCta || 'Shop Now',
            secondaryCta: component.props?.secondaryCta || component.content?.secondaryCta || 'Learn More',
            ctaButtons: [
              { text: component.props?.primaryCta || 'Shop Now', href: '/products', variant: 'primary' },
              { text: component.props?.secondaryCta || 'Learn More', href: '/about', variant: 'secondary' }
            ]
          };
          break;

        case 'footer':
          const quickLinks = component.props?.quickLinks || component.content?.quickLinks || ['About', 'Privacy', 'Terms'];
          baseComponent.content = {
            companyName: component.props?.companyName || component.content?.companyName || 'Your Store',
            description: component.props?.description || component.content?.description || 'Your trusted e-commerce partner',
            email: component.props?.email || component.content?.email || 'contact@yourstore.com',
            phone: component.props?.phone || component.content?.phone || '+1 (555) 123-4567',
            quickLinks: quickLinks.map(link => 
              typeof link === 'object' ? link : { label: link, href: `/${link.toLowerCase()}` }
            ),
            socialLinks: component.content?.socialLinks || [],
            showNewsletter: component.content?.showNewsletter !== false
          };
          break;

        case 'product-grid':
          baseComponent.content = {
            title: component.content?.title || 'Featured Products',
            productsPerRow: component.content?.productsPerRow || 4,
            showFilters: component.content?.showFilters !== false,
            showSorting: component.content?.showSorting !== false,
            category: component.content?.category || 'all'
          };
          break;

        case 'testimonials':
          baseComponent.content = {
            title: component.content?.title || 'What Our Customers Say',
            testimonials: component.content?.testimonials || [
              {
                name: 'John Doe',
                rating: 5,
                comment: 'Great products and excellent service!',
                avatar: ''
              }
            ]
          };
          break;

        default:
          baseComponent.content = component.content || {};
      }

      return baseComponent;
    });

    // Create the complete storefront configuration
    const storefrontConfig = {
      meta: {
        name: designName,
        version: '1.0.0',
        generatedAt: new Date().toISOString(),
        platform: 'next-js'
      },
      theme: {
        id: theme?.id || 'default',
        name: theme?.name || 'Default Theme',
        colors: {
          primary: theme?.colors?.primary || '#3b82f6',
          secondary: theme?.colors?.secondary || '#64748b',
          background: theme?.colors?.background || '#ffffff',
          text: theme?.colors?.text || '#1f2937',
          accent: theme?.colors?.accent || '#f59e0b',
          muted: theme?.colors?.muted || '#f8fafc'
        },
        typography: {
          fontFamily: theme?.typography?.fontFamily || 'Inter, sans-serif',
          fontSize: {
            xs: '0.75rem',
            sm: '0.875rem',
            base: theme?.typography?.fontSize || '1rem',
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
          xl: '3rem'
        },
        borderRadius: {
          sm: '0.25rem',
          md: '0.375rem',
          lg: '0.5rem',
          xl: '0.75rem'
        }
      },
      layout: {
        components: transformedComponents.sort((a, b) => (a.order || 0) - (b.order || 0)),
        settings: {
          maxWidth: settings?.maxWidth || '1200px',
          containerPadding: settings?.containerPadding || '1rem',
          sectionSpacing: settings?.sectionSpacing || '4rem'
        }
      },
      seo: {
        title: `${designName} - E-commerce Store`,
        description: `Welcome to ${designName}, your trusted online shopping destination.`,
        keywords: ['ecommerce', 'online store', 'shopping', designName.toLowerCase()],
        ogImage: '',
        favicon: '/favicon.ico'
      },
      features: {
        cart: true,
        wishlist: true,
        userAccounts: true,
        productReviews: true,
        newsletter: true,
        search: true,
        filters: true,
        multiCurrency: false,
        multiLanguage: false
      }
    };

    return storefrontConfig;
  }
);

// Export design for Next.js action atom
export const exportDesignForNextJSAtom = atom(
  null,
  (get, set) => {
    const config = get(exportForNextJSAtom);
    const fileName = `${config.meta.name.toLowerCase().replace(/\s+/g, '-')}-config.json`;
    
    // Create downloadable JSON file
    const dataStr = JSON.stringify(config, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    
    // Create download link
    const link = document.createElement('a');
    link.href = url;
    link.download = fileName;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
    
    return config;
  }
);