import { useState } from 'react';
import { themeTemplates, themeCategories } from '../../../data/componentTemplates';
import { 
  Palette, 
  Check, 
  Eye, 
  Settings,
  Minimize2,
  ShoppingBag,
  Briefcase,
  Crown,
  ExternalLink
} from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import ThemePreviewModal from './ThemePreviewModal';

const ThemeCard = ({ theme, isSelected, onSelect, onPreview, onCustomize, onBrowserPreview, onApply }) => {
  const [isHovered, setIsHovered] = useState(false);

  const getCategoryIcon = (category) => {
    const icons = {
      minimalist: <Minimize2 className="w-4 h-4" />,
      ecommerce: <ShoppingBag className="w-4 h-4" />,
      creative: <Palette className="w-4 h-4" />,
      business: <Briefcase className="w-4 h-4" />,
      luxury: <Crown className="w-4 h-4" />
    };
    return icons[category] || <Palette className="w-4 h-4" />;
  };

  return (
    <div 
      className={`relative bg-white rounded-lg border-2 transition-all duration-200 cursor-pointer ${
        isSelected 
          ? 'border-blue-500 shadow-lg' 
          : 'border-gray-200 hover:border-gray-300 hover:shadow-md'
      }`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      onClick={() => onSelect(theme)}
    >
      {/* Selection Indicator */}
      {isSelected && (
        <div className="absolute top-2 right-2 w-6 h-6 bg-blue-500 rounded-full flex items-center justify-center z-10">
          <Check className="w-4 h-4 text-white" />
        </div>
      )}

      {/* Theme Preview */}
      <div className="relative overflow-hidden rounded-t-lg">
        <div 
          className="h-32 bg-gradient-to-br p-4"
          style={{
            background: `linear-gradient(135deg, ${theme.colors.primary}20 0%, ${theme.colors.secondary}20 100%)`
          }}
        >
          {/* Mock Layout Preview */}
          <div className="space-y-2">
            {/* Header */}
            <div 
              className="h-2 rounded"
              style={{ backgroundColor: theme.colors.primary }}
            />
            {/* Content */}
            <div className="flex space-x-2">
              <div 
                className="flex-1 h-8 rounded"
                style={{ backgroundColor: theme.colors.surface }}
              />
              <div 
                className="w-8 h-8 rounded"
                style={{ backgroundColor: theme.colors.accent }}
              />
            </div>
            {/* Footer */}
            <div className="flex space-x-1">
              <div 
                className="flex-1 h-1 rounded"
                style={{ backgroundColor: theme.colors.secondary }}
              />
              <div 
                className="flex-1 h-1 rounded"
                style={{ backgroundColor: theme.colors.secondary }}
              />
            </div>
          </div>
        </div>

        {/* Hover Actions */}
        {isHovered && (
          <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center space-x-2">
            <button
              onClick={(e) => {
                e.stopPropagation();
                onPreview(theme);
              }}
              className="p-2 bg-white rounded-full hover:bg-gray-100 transition-colors"
              title="Preview"
            >
              <Eye className="w-4 h-4 text-gray-700" />
            </button>
            <button
              onClick={(e) => {
                e.stopPropagation();
                onBrowserPreview(theme);
              }}
              className="p-2 bg-white rounded-full hover:bg-gray-100 transition-colors"
              title="Browser Preview"
            >
              <ExternalLink className="w-4 h-4 text-gray-700" />
            </button>
            <button
              onClick={(e) => {
                e.stopPropagation();
                onApply(theme);
              }}
              className="p-2 bg-white rounded-full hover:bg-gray-100 transition-colors"
              title="Apply Theme"
            >
              <Check className="w-4 h-4 text-gray-700" />
            </button>
            <button
              onClick={(e) => {
                e.stopPropagation();
                onCustomize(theme);
              }}
              className="p-2 bg-white rounded-full hover:bg-gray-100 transition-colors"
              title="Customize"
            >
              <Settings className="w-4 h-4 text-gray-700" />
            </button>
          </div>
        )}
      </div>

      {/* Theme Info */}
      <div className="p-4">
        <div className="flex items-center justify-between mb-2">
          <h3 className="font-semibold text-gray-900">{theme.name}</h3>
          <div className="flex items-center space-x-1 text-gray-500">
            {getCategoryIcon(theme.category)}
            <span className="text-xs">
              {themeCategories[theme.category]?.name || theme.category}
            </span>
          </div>
        </div>
        <p className="text-sm text-gray-600 mb-3">{theme.description}</p>
        
        {/* Color Palette */}
        <div className="flex items-center space-x-1">
          <span className="text-xs text-gray-500 mr-2">Colors:</span>
          <div 
            className="w-4 h-4 rounded-full border border-gray-200"
            style={{ backgroundColor: theme.colors.primary }}
            title="Primary"
          />
          <div 
            className="w-4 h-4 rounded-full border border-gray-200"
            style={{ backgroundColor: theme.colors.secondary }}
            title="Secondary"
          />
          <div 
            className="w-4 h-4 rounded-full border border-gray-200"
            style={{ backgroundColor: theme.colors.accent }}
            title="Accent"
          />
        </div>
      </div>
    </div>
  );
};

const ThemeTemplates = ({ onThemeSelect, selectedTheme }) => {
  const navigate = useNavigate();
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [previewTheme, setPreviewTheme] = useState(null);
  const [isPreviewModalOpen, setIsPreviewModalOpen] = useState(false);
  const [customizeTheme, setCustomizeTheme] = useState(null);

  const themes = Object.values(themeTemplates);
  const categories = Object.entries(themeCategories);

  const filteredThemes = selectedCategory === 'all' 
    ? themes 
    : themes.filter(theme => theme.category === selectedCategory);

  const handleThemeSelect = (theme) => {
    onThemeSelect(theme);
  };

  const handlePreview = (theme) => {
    // Open standalone preview in new tab instead of modal
    const previewUrl = `/standalone-preview/${theme.id}`;
    window.open(previewUrl, '_blank', 'noopener,noreferrer');
  };

  const closePreviewModal = () => {
    setIsPreviewModalOpen(false);
    setPreviewTheme(null);
  };

  const handleBrowserPreview = (theme) => {
    const previewUrl = `/standalone-preview/${theme.id}`;
    window.open(previewUrl, '_blank', 'noopener,noreferrer');
  };

  const handleApplyTheme = (theme) => {
    console.log('Applying theme:', theme);
    // Apply the theme to the current project
    if (onThemeSelect) {
      onThemeSelect(theme);
    }
  };

  const handleCustomize = (theme) => {
    console.log('Customizing theme:', theme);
    // Apply the theme first
    if (onThemeSelect) {
      onThemeSelect(theme);
    }
    // Switch to customization mode
    setCustomizeTheme(theme);
  };

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center space-x-2 mb-3">
          <div className="p-1.5 bg-purple-100 rounded-md text-purple-600">
            <Palette className="w-4 h-4" />
          </div>
          <h2 className="text-lg font-semibold text-gray-900">Theme Templates</h2>
        </div>
        <p className="text-sm text-gray-600">Choose a pre-designed theme to get started quickly</p>
      </div>

      {/* Category Filter */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex flex-wrap gap-2">
          <button
            onClick={() => setSelectedCategory('all')}
            className={`px-3 py-1.5 text-xs font-medium rounded-full transition-colors ${
              selectedCategory === 'all'
                ? 'bg-blue-100 text-blue-700'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            All Themes
          </button>
          {categories.map(([key, category]) => (
            <button
              key={key}
              onClick={() => setSelectedCategory(key)}
              className={`px-3 py-1.5 text-xs font-medium rounded-full transition-colors ${
                selectedCategory === key
                  ? 'bg-blue-100 text-blue-700'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {category.name}
            </button>
          ))}
        </div>
      </div>

      {/* Theme Grid */}
      <div className="flex-1 overflow-y-auto p-4">
        <div className="grid grid-cols-1 gap-4">
          {filteredThemes.map(theme => (
            <ThemeCard
              key={theme.id}
              theme={theme}
              isSelected={selectedTheme?.id === theme.id}
              onSelect={handleThemeSelect}
              onPreview={handlePreview}
              onCustomize={handleCustomize}
              onBrowserPreview={handleBrowserPreview}
              onApply={handleApplyTheme}
            />
          ))}
        </div>
        
        {filteredThemes.length === 0 && (
          <div className="text-center py-12">
            <Palette className="w-12 h-12 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-500">No themes found in this category</p>
          </div>
        )}
      </div>

      {/* Selected Theme Info */}
      {selectedTheme && (
        <div className="p-4 border-t border-gray-200 bg-blue-50">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="font-medium text-blue-900">{selectedTheme.name}</h3>
              <p className="text-sm text-blue-700">Theme selected</p>
            </div>
            <button
              onClick={() => handleCustomize(selectedTheme)}
              className="px-3 py-1.5 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 transition-colors"
            >
              Customize
            </button>
          </div>
        </div>
      )}
      
      {/* Theme Preview Modal */}
      <ThemePreviewModal 
        theme={previewTheme}
        isOpen={isPreviewModalOpen}
        onClose={closePreviewModal}
        onApply={handleApplyTheme}
        onCustomize={handleCustomize}
      />
    </div>
  );
};

export default ThemeTemplates;