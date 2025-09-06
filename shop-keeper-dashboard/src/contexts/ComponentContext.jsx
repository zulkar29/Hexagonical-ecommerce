import React, { createContext, useContext, useState, useEffect } from 'react';
import { atom, useAtom } from 'jotai';
import { componentTemplates } from '../data/componentTemplates';

// General app component atoms (separate from builder-specific atoms)
const appComponentsAtom = atom([]);
const appSelectedComponentAtom = atom(null);
const appComponentHistoryAtom = atom([]);
const appIsEditingAtom = atom(false);

const ComponentContext = createContext();

export const useComponent = () => {
  const context = useContext(ComponentContext);
  if (!context) {
    throw new Error('useComponent must be used within a ComponentProvider');
  }
  return context;
};

export const ComponentProvider = ({ children }) => {
  const [components, setComponents] = useAtom(appComponentsAtom);
  const [selectedComponent, setSelectedComponent] = useAtom(appSelectedComponentAtom);
  const [componentHistory, setComponentHistory] = useAtom(appComponentHistoryAtom);
  const [isEditing, setIsEditing] = useAtom(appIsEditingAtom);

  // Load components from localStorage on mount
  useEffect(() => {
    const savedComponents = localStorage.getItem('components');
    if (savedComponents) {
      try {
        setComponents(JSON.parse(savedComponents));
      } catch (error) {
        console.error('Failed to parse saved components:', error);
        // Fallback to default templates
        setComponents(componentTemplates);
      }
    } else {
      // Initialize with default templates
      setComponents(componentTemplates);
    }
  }, [setComponents]);

  // Save components to localStorage when they change
  useEffect(() => {
    if (components.length > 0) {
      localStorage.setItem('components', JSON.stringify(components));
    }
  }, [components]);

  const addComponent = (component) => {
    const newComponent = {
      ...component,
      id: Date.now().toString(),
      createdAt: new Date().toISOString()
    };
    setComponents(prev => [...prev, newComponent]);
    return newComponent;
  };

  const updateComponent = (id, updates) => {
    setComponents(prev => 
      prev.map(comp => 
        comp.id === id ? { ...comp, ...updates, updatedAt: new Date().toISOString() } : comp
      )
    );
  };

  const deleteComponent = (id) => {
    setComponents(prev => prev.filter(comp => comp.id !== id));
    if (selectedComponent?.id === id) {
      setSelectedComponent(null);
    }
  };

  const duplicateComponent = (id) => {
    const component = components.find(comp => comp.id === id);
    if (component) {
      const duplicated = {
        ...component,
        id: Date.now().toString(),
        name: `${component.name} (Copy)`,
        createdAt: new Date().toISOString()
      };
      setComponents(prev => [...prev, duplicated]);
      return duplicated;
    }
  };

  const addToHistory = (action) => {
    setComponentHistory(prev => [
      ...prev.slice(-9), // Keep last 9 actions
      {
        ...action,
        timestamp: new Date().toISOString(),
        id: Date.now().toString()
      }
    ]);
  };

  const clearHistory = () => {
    setComponentHistory([]);
  };

  const value = {
    components,
    selectedComponent,
    componentHistory,
    isEditing,
    setComponents,
    setSelectedComponent,
    setIsEditing,
    addComponent,
    updateComponent,
    deleteComponent,
    duplicateComponent,
    addToHistory,
    clearHistory
  };

  return (
    <ComponentContext.Provider value={value}>
      {children}
    </ComponentContext.Provider>
  );
};

export default ComponentContext;