// Theme utility functions for applying theme configurations

/**
 * Apply theme styles to the document root
 * @param {Object} theme - The selected theme object
 */
export const applyThemeToDocument = (theme) => {
  if (!theme) return;

  const root = document.documentElement;
  const { colors, typography, spacing, borderRadius } = theme;

  // Apply color variables
  if (colors) {
    root.style.setProperty('--theme-primary', colors.primary);
    root.style.setProperty('--theme-secondary', colors.secondary);
    root.style.setProperty('--theme-accent', colors.accent);
    root.style.setProperty('--theme-background', colors.background);
    root.style.setProperty('--theme-surface', colors.surface);
    root.style.setProperty('--theme-text-primary', colors.text.primary);
    root.style.setProperty('--theme-text-secondary', colors.text.secondary);
    root.style.setProperty('--theme-text-muted', colors.text.muted);
  }

  // Apply typography variables
  if (typography) {
    root.style.setProperty('--theme-font-primary', typography.primary);
    root.style.setProperty('--theme-font-secondary', typography.secondary);
    root.style.setProperty('--theme-font-size-xs', typography.sizes.xs);
    root.style.setProperty('--theme-font-size-sm', typography.sizes.sm);
    root.style.setProperty('--theme-font-size-base', typography.sizes.base);
    root.style.setProperty('--theme-font-size-lg', typography.sizes.lg);
    root.style.setProperty('--theme-font-size-xl', typography.sizes.xl);
    root.style.setProperty('--theme-font-size-2xl', typography.sizes['2xl']);
    root.style.setProperty('--theme-font-size-3xl', typography.sizes['3xl']);
  }

  // Apply spacing variables
  if (spacing) {
    Object.entries(spacing).forEach(([key, value]) => {
      root.style.setProperty(`--theme-spacing-${key}`, value);
    });
  }

  // Apply border radius
  if (borderRadius) {
    root.style.setProperty('--theme-border-radius-sm', borderRadius.sm);
    root.style.setProperty('--theme-border-radius-md', borderRadius.md);
    root.style.setProperty('--theme-border-radius-lg', borderRadius.lg);
    root.style.setProperty('--theme-border-radius-xl', borderRadius.xl);
  }

  // Apply responsive breakpoints
  root.style.setProperty('--theme-breakpoint-sm', '640px');
  root.style.setProperty('--theme-breakpoint-md', '768px');
  root.style.setProperty('--theme-breakpoint-lg', '1024px');
  root.style.setProperty('--theme-breakpoint-xl', '1280px');
};

/**
 * Generate CSS classes for theme-aware components with responsive utilities
 * @param {Object} theme - The selected theme object
 * @returns {String} CSS class definitions
 */
export const generateThemeClasses = (theme) => {
  if (!theme) return '';
  
  return `
    .theme-${theme.id} {
      --primary: ${theme.colors.primary};
      --secondary: ${theme.colors.secondary};
      --accent: ${theme.colors.accent};
      --background: ${theme.colors.background};
      --surface: ${theme.colors.surface};
      --text-primary: ${theme.colors.textPrimary};
      --text-secondary: ${theme.colors.textSecondary};
      --border: ${theme.colors.border};
      
      --font-primary: ${theme.typography.primary};
      --font-secondary: ${theme.typography.secondary};
      --font-size-base: ${theme.typography.fontSize};
      
      --spacing-xs: ${theme.spacing.xs};
      --spacing-sm: ${theme.spacing.sm};
      --spacing-md: ${theme.spacing.md};
      --spacing-lg: ${theme.spacing.lg};
      --spacing-xl: ${theme.spacing.xl};
      
      --border-radius: ${theme.borderRadius};
    }
    
    /* Responsive utilities */
    .theme-${theme.id} .responsive-grid {
      display: grid;
      gap: var(--spacing-md);
      grid-template-columns: 1fr;
    }
    
    @media (min-width: 640px) {
      .theme-${theme.id} .responsive-grid {
        grid-template-columns: repeat(2, 1fr);
      }
    }
    
    @media (min-width: 1024px) {
      .theme-${theme.id} .responsive-grid {
        grid-template-columns: repeat(3, 1fr);
      }
    }
    
    .theme-${theme.id} .responsive-text {
      font-size: 0.875rem;
    }
    
    @media (min-width: 768px) {
      .theme-${theme.id} .responsive-text {
        font-size: 1rem;
      }
    }
    
    @media (min-width: 1024px) {
      .theme-${theme.id} .responsive-text {
        font-size: 1.125rem;
      }
    }
  `;
};

/**
 * Get theme classes for components (legacy function for backward compatibility)
 * @param {Object} theme - The selected theme object
 * @returns {Object} CSS class mappings
 */
export const getThemeClasses = (theme) => {
  if (!theme) return {};

  return {
    container: `bg-[var(--theme-background)] text-[var(--theme-text-primary)] font-[var(--theme-font-primary)]`,
    surface: `bg-[var(--theme-surface)] rounded-[var(--theme-border-radius-md)]`,
    primary: `bg-[var(--theme-primary)] text-white`,
    secondary: `bg-[var(--theme-secondary)] text-white`,
    accent: `bg-[var(--theme-accent)] text-white`,
    text: {
      primary: `text-[var(--theme-text-primary)]`,
      secondary: `text-[var(--theme-text-secondary)]`,
      muted: `text-[var(--theme-text-muted)]`
    },
    spacing: {
      xs: `p-[var(--theme-spacing-xs)]`,
      sm: `p-[var(--theme-spacing-sm)]`,
      md: `p-[var(--theme-spacing-md)]`,
      lg: `p-[var(--theme-spacing-lg)]`,
      xl: `p-[var(--theme-spacing-xl)]`
    }
  };
};

/**
 * Apply theme to existing components
 * @param {Array} components - Array of component objects
 * @param {Object} theme - The selected theme object
 * @returns {Array} Updated components with theme styles
 */
export const applyThemeToComponents = (components, theme) => {
  if (!theme || !components) return components;

  const themeClasses = getThemeClasses(theme);

  return components.map(component => {
    const updatedComponent = { ...component };

    // Apply theme-based styling based on component type
    switch (component.type) {
      case 'header':
        updatedComponent.style = {
          ...updatedComponent.style,
          backgroundColor: theme.colors?.primary || updatedComponent.style?.backgroundColor,
          color: 'white',
          fontFamily: theme.typography?.primary || updatedComponent.style?.fontFamily
        };
        break;
      
      case 'button':
        updatedComponent.style = {
          ...updatedComponent.style,
          backgroundColor: theme.colors?.accent || updatedComponent.style?.backgroundColor,
          color: 'white',
          borderRadius: theme.borderRadius?.md || updatedComponent.style?.borderRadius,
          fontFamily: theme.typography?.primary || updatedComponent.style?.fontFamily
        };
        break;
      
      case 'card':
        updatedComponent.style = {
          ...updatedComponent.style,
          backgroundColor: theme.colors?.surface || updatedComponent.style?.backgroundColor,
          color: theme.colors?.text?.primary || updatedComponent.style?.color,
          borderRadius: theme.borderRadius?.lg || updatedComponent.style?.borderRadius,
          fontFamily: theme.typography?.primary || updatedComponent.style?.fontFamily
        };
        break;
      
      case 'text':
        updatedComponent.style = {
          ...updatedComponent.style,
          color: theme.colors?.text?.primary || updatedComponent.style?.color,
          fontFamily: theme.typography?.primary || updatedComponent.style?.fontFamily
        };
        break;
      
      default:
        // Apply general theme styling
        updatedComponent.style = {
          ...updatedComponent.style,
          fontFamily: theme.typography?.primary || updatedComponent.style?.fontFamily,
          color: theme.colors?.text?.primary || updatedComponent.style?.color
        };
    }

    return updatedComponent;
  });
};

/**
 * Reset theme styles from document
 */
export const resetTheme = () => {
  const root = document.documentElement;
  const themeProperties = [
    '--theme-primary', '--theme-secondary', '--theme-accent',
    '--theme-background', '--theme-surface',
    '--theme-text-primary', '--theme-text-secondary', '--theme-text-muted',
    '--theme-font-primary', '--theme-font-secondary',
    '--theme-font-size-xs', '--theme-font-size-sm', '--theme-font-size-base',
    '--theme-font-size-lg', '--theme-font-size-xl', '--theme-font-size-2xl', '--theme-font-size-3xl',
    '--theme-border-radius-sm', '--theme-border-radius-md', '--theme-border-radius-lg', '--theme-border-radius-xl'
  ];

  themeProperties.forEach(prop => {
    root.style.removeProperty(prop);
  });

  // Remove spacing properties
  ['xs', 'sm', 'md', 'lg', 'xl', '2xl', '3xl'].forEach(size => {
    root.style.removeProperty(`--theme-spacing-${size}`);
  });
}