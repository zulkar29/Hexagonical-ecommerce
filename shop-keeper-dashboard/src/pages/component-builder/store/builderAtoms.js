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