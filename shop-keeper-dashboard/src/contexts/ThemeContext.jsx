import React, { createContext, useContext, useState, useEffect } from 'react';
import { atom, useAtom } from 'jotai';

// Theme atoms
const currentThemeAtom = atom(null);
const selectedThemeAtom = atom(null);
const themeCustomizationAtom = atom({});
const isPreviewModeAtom = atom(false);

const ThemeContext = createContext();

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
};

export const ThemeProvider = ({ children }) => {
  const [currentTheme, setCurrentTheme] = useAtom(currentThemeAtom);
  const [selectedTheme, setSelectedTheme] = useAtom(selectedThemeAtom);
  const [themeCustomization, setThemeCustomization] = useAtom(themeCustomizationAtom);
  const [isPreviewMode, setIsPreviewMode] = useAtom(isPreviewModeAtom);

  // Load theme from localStorage on mount
  useEffect(() => {
    const savedTheme = localStorage.getItem('currentTheme');
    const savedCustomization = localStorage.getItem('themeCustomization');
    
    if (savedTheme) {
      try {
        setCurrentTheme(JSON.parse(savedTheme));
      } catch (error) {
        console.error('Failed to parse saved theme:', error);
      }
    }
    
    if (savedCustomization) {
      try {
        setThemeCustomization(JSON.parse(savedCustomization));
      } catch (error) {
        console.error('Failed to parse saved customization:', error);
      }
    }
  }, [setCurrentTheme, setThemeCustomization]);

  // Save theme to localStorage when it changes
  useEffect(() => {
    if (currentTheme) {
      localStorage.setItem('currentTheme', JSON.stringify(currentTheme));
    }
  }, [currentTheme]);

  // Save customization to localStorage when it changes
  useEffect(() => {
    if (Object.keys(themeCustomization).length > 0) {
      localStorage.setItem('themeCustomization', JSON.stringify(themeCustomization));
    }
  }, [themeCustomization]);

  const applyTheme = (theme) => {
    setCurrentTheme(theme);
    setSelectedTheme(theme);
  };

  const updateCustomization = (customization) => {
    setThemeCustomization(prev => ({ ...prev, ...customization }));
  };

  const resetCustomization = () => {
    setThemeCustomization({});
    localStorage.removeItem('themeCustomization');
  };

  const value = {
    currentTheme,
    selectedTheme,
    themeCustomization,
    isPreviewMode,
    setCurrentTheme,
    setSelectedTheme,
    setIsPreviewMode,
    applyTheme,
    updateCustomization,
    resetCustomization
  };

  return (
    <ThemeContext.Provider value={value}>
      {children}
    </ThemeContext.Provider>
  );
};

export default ThemeContext;