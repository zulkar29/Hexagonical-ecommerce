import { useState } from 'react';
import { useDraggable } from '@dnd-kit/core';
import { componentTemplates, getCategories, getCategoryDisplayName, getTemplatesByCategory } from '../../../data/componentTemplates';
import ThemeTemplates from './ThemeTemplates';
import { 
  Palette, 
  Layout, 
  Type, 
  ShoppingBag, 
  Image, 
  Settings, 
  Layers,
  Grid3X3,
  Sparkles
} from 'lucide-react';

const DraggableComponent = ({ template }) => {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: template.id,
    data: {
      type: 'library-item',
      template
    }
  });

  const style = transform ? {
    transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
    opacity: isDragging ? 0.5 : 1
  } : undefined;

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      className="p-4 bg-white border border-gray-200 rounded-lg cursor-grab hover:border-blue-300 hover:shadow-md transition-all duration-200 active:cursor-grabbing"
    >
      <div className="flex items-center space-x-3">
        <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
          <span className="text-white text-lg">{template.icon}</span>
        </div>
        <div>
          <h3 className="font-medium text-gray-900">{template.name}</h3>
          <p className="text-sm text-gray-500">{template.description}</p>
        </div>
      </div>
      
      {/* Preview */}
      <div className="mt-3 p-2 bg-gray-50 rounded border">
        <div className="text-xs text-gray-600 mb-1">Preview:</div>
        <div className="bg-white rounded border p-2 text-xs">
          {template.preview}
        </div>
      </div>
    </div>
  );
};

const ComponentLibrary = ({ onThemeSelect, selectedTheme }) => {
  const [activeTab, setActiveTab] = useState('themes');

  // Organized sections with proper categorization
  const sections = [
    {
      id: 'layout-components',
      title: 'Layout Components',
      icon: <Layout className="w-4 h-4" />,
      description: 'Headers, footers, and structure',
      components: [
        ...getTemplatesByCategory('header'),
        ...getTemplatesByCategory('footer')
      ]
    },
    {
      id: 'content-blocks',
      title: 'Content Blocks',
      icon: <Type className="w-4 h-4" />,
      description: 'Text, features, and content sections',
      components: getTemplatesByCategory('content')
    },
    {
      id: 'product-widgets',
      title: 'Product Widgets',
      icon: <ShoppingBag className="w-4 h-4" />,
      description: 'Product cards, grids, and commerce',
      components: getTemplatesByCategory('product')
    },
    {
      id: 'media-gallery',
      title: 'Media Gallery',
      icon: <Image className="w-4 h-4" />,
      description: 'Sliders, carousels, and media',
      components: getTemplatesByCategory('slider')
    }
  ];

  const ComponentSection = ({ section }) => (
    <div className="mb-8">
      <div className="px-4 mb-4">
        <div className="flex items-center space-x-2 mb-2">
          <div className="p-1.5 bg-blue-100 rounded-md text-blue-600">
            {section.icon}
          </div>
          <h3 className="text-sm font-semibold text-gray-900">
            {section.title}
          </h3>
        </div>
        <p className="text-xs text-gray-500 ml-8">{section.description}</p>
      </div>
      
      <div className="space-y-3 px-4">
        {section.components.map(template => (
          <DraggableComponent key={template.id} template={template} />
        ))}
      </div>
    </div>
  );

  return (
    <div className="h-full flex flex-col">
      {/* Header with Tabs */}
      <div className="p-4 border-b border-gray-200 bg-gradient-to-r from-blue-50 to-purple-50">
        <div className="flex items-center space-x-2 mb-3">
          <Palette className="w-5 h-5 text-blue-600" />
          <h2 className="text-lg font-semibold text-gray-900">Design Library</h2>
        </div>
        
        {/* Tab Navigation */}
        <div className="flex space-x-1 bg-white rounded-lg p-1">
          <button
            onClick={() => setActiveTab('themes')}
            className={`flex-1 px-3 py-2 text-sm font-medium rounded-md transition-colors ${
              activeTab === 'themes'
                ? 'bg-blue-100 text-blue-700'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            <Sparkles className="w-4 h-4 inline mr-1" />
            Themes
          </button>
          <button
            onClick={() => setActiveTab('components')}
            className={`flex-1 px-3 py-2 text-sm font-medium rounded-md transition-colors ${
              activeTab === 'components'
                ? 'bg-blue-100 text-blue-700'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            <Layout className="w-4 h-4 inline mr-1" />
            Components
          </button>
        </div>
      </div>
      
      {/* Tab Content */}
      <div className="flex-1 overflow-hidden">
        {activeTab === 'themes' ? (
          <ThemeTemplates 
            onThemeSelect={onThemeSelect}
            selectedTheme={selectedTheme}
          />
        ) : (
          <div className="h-full overflow-y-auto py-6">
            {sections.map(section => (
              <ComponentSection key={section.id} section={section} />
            ))}
            
            {/* Settings Section */}
            <div className="px-4 mt-8 pt-6 border-t border-gray-200">
              <div className="flex items-center space-x-2 mb-4">
                <div className="p-1.5 bg-gray-100 rounded-md text-gray-600">
                  <Settings className="w-4 h-4" />
                </div>
                <h3 className="text-sm font-semibold text-gray-900">Settings</h3>
              </div>
              <div className="space-y-2">
                <button className="w-full text-left px-3 py-2 text-sm text-gray-600 hover:bg-gray-50 rounded-md transition-colors">
                  Global Styles
                </button>
                <button className="w-full text-left px-3 py-2 text-sm text-gray-600 hover:bg-gray-50 rounded-md transition-colors">
                  Theme Settings
                </button>
                <button className="w-full text-left px-3 py-2 text-sm text-gray-600 hover:bg-gray-50 rounded-md transition-colors">
                  Export/Import
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ComponentLibrary;