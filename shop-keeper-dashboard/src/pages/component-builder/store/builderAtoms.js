import { atom } from 'jotai';

// Core builder state atoms
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

// Builder settings atoms
export const builderSettingsAtom = atom({
  autoSave: true,
  showGrid: false,
  snapToGrid: true,
  gridSize: 8
});

// Actions atoms (write-only atoms for complex operations)
export const addComponentAtom = atom(
  null,
  (get, set, component) => {
    const currentComponents = get(componentsAtom);
    set(componentsAtom, [...currentComponents, component]);
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
  }
);