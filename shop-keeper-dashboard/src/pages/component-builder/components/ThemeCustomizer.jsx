import React, { useState } from 'react';
import { Palette, Type, Layout, RotateCcw } from 'lucide-react';

const ThemeCustomizer = ({ selectedTheme, onThemeUpdate, onResetTheme }) => {
  const [activeTab, setActiveTab] = useState('colors');
  const [customTheme, setCustomTheme] = useState(selectedTheme || {});

  if (!selectedTheme) {
    return (
      <div className="p-4 text-center text-gray-500">
        <Palette className="w-8 h-8 mx-auto mb-2 text-gray-400" />
        <p>Select a theme to customize</p>
      </div>
    );
  }

  const handleColorChange = (colorKey, value) => {
    const updatedTheme = {
      ...customTheme,
      colors: {
        ...customTheme.colors,
        [colorKey]: value
      }
    };
    setCustomTheme(updatedTheme);
    onThemeUpdate(updatedTheme);
  };

  const handleTextColorChange = (textKey, value) => {
    const updatedTheme = {
      ...customTheme,
      colors: {
        ...customTheme.colors,
        text: {
          ...customTheme.colors.text,
          [textKey]: value
        }
      }
    };
    setCustomTheme(updatedTheme);
    onThemeUpdate(updatedTheme);
  };

  const handleTypographyChange = (fontKey, value) => {
    const updatedTheme = {
      ...customTheme,
      typography: {
        ...customTheme.typography,
        [fontKey]: value
      }
    };
    setCustomTheme(updatedTheme);
    onThemeUpdate(updatedTheme);
  };

  const handleSpacingChange = (spacingKey, value) => {
    const updatedTheme = {
      ...customTheme,
      spacing: {
        ...customTheme.spacing,
        [spacingKey]: value
      }
    };
    setCustomTheme(updatedTheme);
    onThemeUpdate(updatedTheme);
  };

  const handleBorderRadiusChange = (radiusKey, value) => {
    const updatedTheme = {
      ...customTheme,
      borderRadius: {
        ...customTheme.borderRadius,
        [radiusKey]: value
      }
    };
    setCustomTheme(updatedTheme);
    onThemeUpdate(updatedTheme);
  };

  const resetToOriginal = () => {
    setCustomTheme(selectedTheme);
    onResetTheme();
  };

  const tabs = [
    { id: 'colors', label: 'Colors', icon: Palette },
    { id: 'typography', label: 'Typography', icon: Type },
    { id: 'layout', label: 'Layout', icon: Layout }
  ];

  const fontOptions = [
    'Inter, sans-serif',
    'Roboto, sans-serif',
    'Open Sans, sans-serif',
    'Lato, sans-serif',
    'Montserrat, sans-serif',
    'Poppins, sans-serif',
    'Playfair Display, serif',
    'Merriweather, serif',
    'Georgia, serif'
  ];

  return (
    <div className="bg-white border-l border-gray-200 w-80 flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-gray-900">Theme Customizer</h3>
          <button
            onClick={resetToOriginal}
            className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-md transition-colors"
            title="Reset to original theme"
          >
            <RotateCcw className="w-4 h-4" />
          </button>
        </div>
        <p className="text-sm text-gray-600 mt-1">Customizing: {selectedTheme.name}</p>
      </div>

      {/* Tabs */}
      <div className="flex border-b border-gray-200">
        {tabs.map(tab => {
          const Icon = tab.icon;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`flex-1 flex items-center justify-center px-3 py-2 text-sm font-medium transition-colors ${
                activeTab === tab.id
                  ? 'text-blue-600 border-b-2 border-blue-600 bg-blue-50'
                  : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
              }`}
            >
              <Icon className="w-4 h-4 mr-1" />
              {tab.label}
            </button>
          );
        })}
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-4">
        {activeTab === 'colors' && (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Primary Color</label>
              <div className="flex items-center space-x-2">
                <input
                  type="color"
                  value={customTheme.colors?.primary || '#3b82f6'}
                  onChange={(e) => handleColorChange('primary', e.target.value)}
                  className="w-10 h-10 rounded border border-gray-300"
                />
                <input
                  type="text"
                  value={customTheme.colors?.primary || '#3b82f6'}
                  onChange={(e) => handleColorChange('primary', e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Secondary Color</label>
              <div className="flex items-center space-x-2">
                <input
                  type="color"
                  value={customTheme.colors?.secondary || '#6b7280'}
                  onChange={(e) => handleColorChange('secondary', e.target.value)}
                  className="w-10 h-10 rounded border border-gray-300"
                />
                <input
                  type="text"
                  value={customTheme.colors?.secondary || '#6b7280'}
                  onChange={(e) => handleColorChange('secondary', e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Accent Color</label>
              <div className="flex items-center space-x-2">
                <input
                  type="color"
                  value={customTheme.colors?.accent || '#10b981'}
                  onChange={(e) => handleColorChange('accent', e.target.value)}
                  className="w-10 h-10 rounded border border-gray-300"
                />
                <input
                  type="text"
                  value={customTheme.colors?.accent || '#10b981'}
                  onChange={(e) => handleColorChange('accent', e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Background Color</label>
              <div className="flex items-center space-x-2">
                <input
                  type="color"
                  value={customTheme.colors?.background || '#ffffff'}
                  onChange={(e) => handleColorChange('background', e.target.value)}
                  className="w-10 h-10 rounded border border-gray-300"
                />
                <input
                  type="text"
                  value={customTheme.colors?.background || '#ffffff'}
                  onChange={(e) => handleColorChange('background', e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Surface Color</label>
              <div className="flex items-center space-x-2">
                <input
                  type="color"
                  value={customTheme.colors?.surface || '#f9fafb'}
                  onChange={(e) => handleColorChange('surface', e.target.value)}
                  className="w-10 h-10 rounded border border-gray-300"
                />
                <input
                  type="text"
                  value={customTheme.colors?.surface || '#f9fafb'}
                  onChange={(e) => handleColorChange('surface', e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm"
                />
              </div>
            </div>

            <div className="border-t pt-4">
              <h4 className="text-sm font-medium text-gray-700 mb-3">Text Colors</h4>
              
              <div className="space-y-3">
                <div>
                  <label className="block text-xs text-gray-600 mb-1">Primary Text</label>
                  <div className="flex items-center space-x-2">
                    <input
                      type="color"
                      value={customTheme.colors?.text?.primary || '#111827'}
                      onChange={(e) => handleTextColorChange('primary', e.target.value)}
                      className="w-8 h-8 rounded border border-gray-300"
                    />
                    <input
                      type="text"
                      value={customTheme.colors?.text?.primary || '#111827'}
                      onChange={(e) => handleTextColorChange('primary', e.target.value)}
                      className="flex-1 px-2 py-1 border border-gray-300 rounded text-xs"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-xs text-gray-600 mb-1">Secondary Text</label>
                  <div className="flex items-center space-x-2">
                    <input
                      type="color"
                      value={customTheme.colors?.text?.secondary || '#6b7280'}
                      onChange={(e) => handleTextColorChange('secondary', e.target.value)}
                      className="w-8 h-8 rounded border border-gray-300"
                    />
                    <input
                      type="text"
                      value={customTheme.colors?.text?.secondary || '#6b7280'}
                      onChange={(e) => handleTextColorChange('secondary', e.target.value)}
                      className="flex-1 px-2 py-1 border border-gray-300 rounded text-xs"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-xs text-gray-600 mb-1">Muted Text</label>
                  <div className="flex items-center space-x-2">
                    <input
                      type="color"
                      value={customTheme.colors?.text?.muted || '#9ca3af'}
                      onChange={(e) => handleTextColorChange('muted', e.target.value)}
                      className="w-8 h-8 rounded border border-gray-300"
                    />
                    <input
                      type="text"
                      value={customTheme.colors?.text?.muted || '#9ca3af'}
                      onChange={(e) => handleTextColorChange('muted', e.target.value)}
                      className="flex-1 px-2 py-1 border border-gray-300 rounded text-xs"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'typography' && (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Primary Font</label>
              <select
                value={customTheme.typography?.primary || 'Inter, sans-serif'}
                onChange={(e) => handleTypographyChange('primary', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
              >
                {fontOptions.map(font => (
                  <option key={font} value={font} style={{ fontFamily: font }}>
                    {font.split(',')[0]}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Secondary Font</label>
              <select
                value={customTheme.typography?.secondary || 'Inter, sans-serif'}
                onChange={(e) => handleTypographyChange('secondary', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
              >
                {fontOptions.map(font => (
                  <option key={font} value={font} style={{ fontFamily: font }}>
                    {font.split(',')[0]}
                  </option>
                ))}
              </select>
            </div>

            <div className="border-t pt-4">
              <h4 className="text-sm font-medium text-gray-700 mb-3">Font Sizes</h4>
              <div className="space-y-3">
                {Object.entries(customTheme.typography?.sizes || {}).map(([size, value]) => (
                  <div key={size}>
                    <label className="block text-xs text-gray-600 mb-1 capitalize">{size}</label>
                    <input
                      type="text"
                      value={value}
                      onChange={(e) => handleTypographyChange(`sizes.${size}`, e.target.value)}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-xs"
                      placeholder="e.g., 1rem, 16px"
                    />
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {activeTab === 'layout' && (
          <div className="space-y-4">
            <div>
              <h4 className="text-sm font-medium text-gray-700 mb-3">Spacing</h4>
              <div className="space-y-3">
                {Object.entries(customTheme.spacing || {}).map(([size, value]) => (
                  <div key={size}>
                    <label className="block text-xs text-gray-600 mb-1 capitalize">{size}</label>
                    <input
                      type="text"
                      value={value}
                      onChange={(e) => handleSpacingChange(size, e.target.value)}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-xs"
                      placeholder="e.g., 0.5rem, 8px"
                    />
                  </div>
                ))}
              </div>
            </div>

            <div className="border-t pt-4">
              <h4 className="text-sm font-medium text-gray-700 mb-3">Border Radius</h4>
              <div className="space-y-3">
                {Object.entries(customTheme.borderRadius || {}).map(([size, value]) => (
                  <div key={size}>
                    <label className="block text-xs text-gray-600 mb-1 capitalize">{size}</label>
                    <input
                      type="text"
                      value={value}
                      onChange={(e) => handleBorderRadiusChange(size, e.target.value)}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-xs"
                      placeholder="e.g., 0.375rem, 6px"
                    />
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ThemeCustomizer;