import React, { useState } from 'react';
import { themeTemplates } from '../../../data/shopTemplates';
import {
  Palette,
  Check,
  Eye,
  Download,
  Star,
  Zap,
  Crown,
  Sparkles,
  Search,
  Filter,
  Grid,
  List
} from 'lucide-react';

const ThemePreview = ({ theme }) => {
  return (
    <div className="relative bg-white border border-gray-200 rounded-lg overflow-hidden">
      {/* Theme Preview Canvas */}
      <div className="h-32 relative">
        {/* Mock browser frame */}
        <div className="absolute top-2 left-2 right-2 bg-gray-100 rounded-t border-b">
          <div className="flex items-center space-x-1 p-1">
            <div className="w-2 h-2 bg-red-400 rounded-full"></div>
            <div className="w-2 h-2 bg-yellow-400 rounded-full"></div>
            <div className="w-2 h-2 bg-green-400 rounded-full"></div>
          </div>
        </div>

        {/* Theme preview content */}
        <div className="absolute top-6 left-2 right-2 bottom-2 rounded-b" style={{
          background: `linear-gradient(135deg, ${theme.colors.primary}15, ${theme.colors.secondary}15)`
        }}>
          {/* Header */}
          <div className="h-6 flex items-center justify-between px-2 border-b border-gray-200/50" style={{
            backgroundColor: theme.colors.background,
            borderColor: theme.colors.border
          }}>
            <div className="text-xs font-semibold" style={{ color: theme.colors.text }}>
              Store
            </div>
            <div className="flex space-x-1">
              <div className="w-2 h-2 rounded" style={{ backgroundColor: theme.colors.primary }}></div>
              <div className="w-2 h-2 rounded" style={{ backgroundColor: theme.colors.secondary }}></div>
            </div>
          </div>

          {/* Hero section */}
          <div className="h-8 flex items-center justify-center" style={{
            backgroundColor: theme.colors.primary + '20'
          }}>
            <div className="text-xs font-bold" style={{ color: theme.colors.text }}>
              {theme.name}
            </div>
          </div>

          {/* Content grid */}
          <div className="p-2 grid grid-cols-2 gap-1">
            {[1, 2, 3, 4].map(i => (
              <div key={i} className="h-4 rounded" style={{
                backgroundColor: theme.colors.surface,
                border: `1px solid ${theme.colors.border}`
              }}></div>
            ))}
          </div>
        </div>
      </div>

      {/* Color palette */}
      <div className="p-3 border-t border-gray-200">
        <div className="flex items-center justify-between mb-2">
          <span className="text-sm font-medium text-gray-900">{theme.name}</span>
          <div className="flex space-x-1">
            <div className="w-3 h-3 rounded-full border border-gray-300" style={{
              backgroundColor: theme.colors.primary
            }} title="Primary"></div>
            <div className="w-3 h-3 rounded-full border border-gray-300" style={{
              backgroundColor: theme.colors.secondary
            }} title="Secondary"></div>
            <div className="w-3 h-3 rounded-full border border-gray-300" style={{
              backgroundColor: theme.colors.accent
            }} title="Accent"></div>
          </div>
        </div>
        <p className="text-xs text-gray-500 line-clamp-2">{theme.description}</p>
      </div>
    </div>
  );
};

const ThemeCard = ({ theme, isSelected, onSelect, onPreview }) => {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <div
      className={`relative group cursor-pointer transition-all duration-200 ${
        isSelected ? 'ring-2 ring-blue-500 ring-offset-2' : ''
      }`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      onClick={onSelect}
    >
      <ThemePreview theme={theme} />

      {/* Selected indicator */}
      {isSelected && (
        <div className="absolute top-2 right-2 w-6 h-6 bg-blue-500 rounded-full flex items-center justify-center">
          <Check className="w-4 h-4 text-white" />
        </div>
      )}

      {/* Hover overlay */}
      {isHovered && !isSelected && (
        <div className="absolute inset-0 bg-black bg-opacity-10 rounded-lg flex items-center justify-center">
          <div className="bg-white rounded-lg shadow-lg p-2 flex space-x-2">
            <button
              onClick={(e) => {
                e.stopPropagation();
                onPreview?.(theme);
              }}
              className="flex items-center space-x-1 px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded"
            >
              <Eye className="w-3 h-3" />
              <span>Preview</span>
            </button>
            <button
              onClick={onSelect}
              className="flex items-center space-x-1 px-3 py-1 text-sm bg-blue-600 text-white hover:bg-blue-700 rounded"
            >
              <Download className="w-3 h-3" />
              <span>Apply</span>
            </button>
          </div>
        </div>
      )}

      {/* Premium badge */}
      {theme.premium && (
        <div className="absolute top-2 left-2">
          <div className="bg-gradient-to-r from-yellow-400 to-orange-500 text-white text-xs px-2 py-1 rounded-full flex items-center space-x-1">
            <Crown className="w-3 h-3" />
            <span>Pro</span>
          </div>
        </div>
      )}
    </div>
  );
};

const ThemesPanel = ({ onThemeSelect, selectedTheme }) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [viewMode, setViewMode] = useState('grid');

  const categories = [
    { id: 'all', name: 'All Themes', icon: Palette },
    { id: 'popular', name: 'Popular', icon: Star },
    { id: 'modern', name: 'Modern', icon: Zap },
    { id: 'classic', name: 'Classic', icon: Crown },
    { id: 'creative', name: 'Creative', icon: Sparkles }
  ];

  // Convert themeTemplates object to array
  const themesArray = Object.values(themeTemplates || {});

  // Filter themes
  const filteredThemes = themesArray.filter(theme => {
    const matchesSearch = theme.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         theme.description.toLowerCase().includes(searchQuery.toLowerCase());

    if (selectedCategory === 'all') return matchesSearch;

    // Add category filtering logic here based on theme properties
    return matchesSearch;
  });

  const handleThemeSelect = (theme) => {
    onThemeSelect?.(theme);
  };

  const handleThemePreview = (theme) => {
    // Open theme preview in new tab/modal
    const previewUrl = `/standalone-preview/${theme.id}`;
    window.open(previewUrl, '_blank', 'width=1200,height=800');
  };

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-lg font-semibold text-gray-900">Themes</h2>
          <div className="flex items-center space-x-1">
            <button
              onClick={() => setViewMode('grid')}
              className={`p-2 rounded-md transition-colors ${
                viewMode === 'grid'
                  ? 'bg-blue-100 text-blue-600'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <Grid className="w-4 h-4" />
            </button>
            <button
              onClick={() => setViewMode('list')}
              className={`p-2 rounded-md transition-colors ${
                viewMode === 'list'
                  ? 'bg-blue-100 text-blue-600'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <List className="w-4 h-4" />
            </button>
          </div>
        </div>

        {/* Search */}
        <div className="relative mb-3">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            type="text"
            placeholder="Search themes..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* Categories */}
        <div className="flex flex-wrap gap-1">
          {categories.map(category => {
            const Icon = category.icon;
            return (
              <button
                key={category.id}
                onClick={() => setSelectedCategory(category.id)}
                className={`flex items-center space-x-1 px-3 py-1 rounded-full text-sm transition-colors ${
                  selectedCategory === category.id
                    ? 'bg-blue-100 text-blue-700'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                <Icon className="w-3 h-3" />
                <span>{category.name}</span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Themes Grid */}
      <div className="flex-1 overflow-y-auto p-4">
        {selectedTheme && (
          <div className="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                <Palette className="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <h3 className="font-medium text-blue-900">Current Theme</h3>
                <p className="text-sm text-blue-700">{selectedTheme.name}</p>
              </div>
            </div>
          </div>
        )}

        <div className={`${
          viewMode === 'grid'
            ? 'grid grid-cols-1 gap-4'
            : 'space-y-3'
        }`}>
          {filteredThemes.map(theme => (
            <ThemeCard
              key={theme.id}
              theme={theme}
              isSelected={selectedTheme?.id === theme.id}
              onSelect={() => handleThemeSelect(theme)}
              onPreview={() => handleThemePreview(theme)}
            />
          ))}
        </div>

        {filteredThemes.length === 0 && (
          <div className="text-center py-8 text-gray-500">
            <Palette className="w-12 h-12 mx-auto mb-3 text-gray-300" />
            <p>No themes found</p>
            <p className="text-sm">Try adjusting your search or filters</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default ThemesPanel;