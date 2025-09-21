import { useState } from 'react';
import { useDraggable } from '@dnd-kit/core';
import { componentTemplates, getTemplatesByCategory } from '../../../data/componentTemplates';
import VendorFeatures from './VendorFeatures';
import QuickStart from './QuickStart';
import {
  Palette,
  Layout,
  Type,
  ShoppingBag,
  Settings,
  Rocket
} from 'lucide-react';

const DraggableComponent = ({ template }) => {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: template.id,
    data: {
      type: 'library-item',
      template: template
    }
  });

  const style = transform ? {
    transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
    opacity: isDragging ? 0.5 : 1
  } : undefined;

  const getComponentTypeColor = (type) => {
    switch (type) {
      case 'header':
        return 'from-blue-500 to-blue-600';
      case 'footer':
        return 'from-gray-500 to-gray-600';
      case 'hero':
        return 'from-purple-500 to-purple-600';
      case 'product':
        return 'from-green-500 to-green-600';
      case 'category':
        return 'from-orange-500 to-orange-600';
      default:
        return 'from-blue-500 to-purple-600';
    }
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      className="group p-4 bg-white border border-gray-200 rounded-xl hover:border-blue-300 hover:shadow-lg transition-all duration-200 hover:scale-[1.02] select-none cursor-grab active:cursor-grabbing"
    >
      <div className="flex items-start space-x-3">
        <div className={`w-12 h-12 bg-gradient-to-br ${getComponentTypeColor(template.type)} rounded-xl flex items-center justify-center shadow-sm`}>
          <span className="text-white text-xl">{template.icon}</span>
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center space-x-2 mb-1">
            <h3 className="font-semibold text-gray-900 truncate">{template.name}</h3>
            <span className="px-2 py-0.5 text-xs font-medium bg-gray-100 text-gray-700 rounded-full">
              {template.type}
            </span>
          </div>
          <p className="text-sm text-gray-600 leading-snug">{template.description}</p>
        </div>
      </div>

      {/* Enhanced Preview */}
      <div className="mt-4 p-3 bg-gradient-to-br from-gray-50 to-gray-100 rounded-lg border">
        <div className="flex items-center justify-between mb-2">
          <div className="text-xs font-medium text-gray-700">Component Preview</div>
          <div className="text-xs text-blue-600 font-medium">⚡ Drag to add</div>
        </div>
        <div className="bg-white rounded-md border p-3 text-xs font-mono text-gray-800 whitespace-pre-line leading-relaxed shadow-sm">
          {template.preview || 'No preview available'}
        </div>
      </div>

      {/* Actions */}
      <div className="mt-3 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
        <div className="flex items-center justify-between text-xs text-gray-500">
          <span>Ready to use</span>
          <div className="flex items-center space-x-2">
            <span>✨ Customizable</span>
            <div className="text-blue-600 font-medium">
              🚀 Drag to add
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

const ComponentLibrary = ({ onThemeSelect, selectedTheme, onQuickStart }) => {
  const [activeTab, setActiveTab] = useState('quickstart');

  // Organized sections with proper categorization for e-commerce
  const sections = [
    {
      id: 'store-basics',
      title: 'Store Basics',
      icon: <Layout className="w-4 h-4" />,
      description: 'Essential components every store needs',
      badge: 'Required',
      components: [
        ...getTemplatesByCategory('header'),
        ...getTemplatesByCategory('footer')
      ]
    },
    {
      id: 'products',
      title: 'Product Display',
      icon: <ShoppingBag className="w-4 h-4" />,
      description: 'Showcase your products beautifully',
      badge: 'Popular',
      components: getTemplatesByCategory('product')
    },
    {
      id: 'content',
      title: 'Content Sections',
      icon: <Type className="w-4 h-4" />,
      description: 'Add content and promotional areas',
      badge: 'Flexible',
      components: [
        ...getTemplatesByCategory('content'),
        ...getTemplatesByCategory('slider')
      ]
    }
  ];

  const ComponentSection = ({ section }) => (
    <div className="mb-8">
      {/* Enhanced Section Header */}
      <div className="px-4 mb-6">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-gradient-to-br from-blue-100 to-blue-200 rounded-lg text-blue-600 shadow-sm">
              {section.icon}
            </div>
            <div>
              <h3 className="text-lg font-bold text-gray-900">
                {section.title}
              </h3>
              <p className="text-sm text-gray-600">{section.description}</p>
            </div>
          </div>
          {section.badge && (
            <div className="flex flex-col items-end">
              <span className={`px-3 py-1 text-sm font-semibold rounded-full shadow-sm ${
                section.badge === 'Required' ? 'bg-red-100 text-red-700 border border-red-200' :
                section.badge === 'Popular' ? 'bg-blue-100 text-blue-700 border border-blue-200' :
                section.badge === 'Flexible' ? 'bg-green-100 text-green-700 border border-green-200' :
                'bg-gray-100 text-gray-700 border border-gray-200'
              }`}>
                {section.badge}
              </span>
              <span className="text-xs text-gray-500 mt-1">
                {section.components.length} components
              </span>
            </div>
          )}
        </div>
        <div className="h-px bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200"></div>
      </div>

      {/* Component Grid */}
      <div className="px-4">
        {section.components.length > 0 ? (
          <div className="grid gap-4">
            {section.components.map(template => (
              <DraggableComponent key={template.id} template={template} />
            ))}
          </div>
        ) : (
          <div className="text-center py-12 text-gray-500 bg-gray-50 rounded-xl border-2 border-dashed border-gray-200">
            <div className="text-lg mb-2">📦</div>
            <div className="text-base font-medium mb-1">No components available</div>
            <div className="text-sm">More components coming soon!</div>
          </div>
        )}
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
        <div className="grid grid-cols-3 gap-2 bg-white rounded-lg p-1">
          <button
            onClick={() => setActiveTab('quickstart')}
            className={`px-4 py-3 text-sm font-medium rounded-md transition-colors ${
              activeTab === 'quickstart'
                ? 'bg-blue-100 text-blue-700 shadow-sm'
                : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
            }`}
          >
            <Rocket className="w-4 h-4 inline mr-2" />
            Quick Start
          </button>
          <button
            onClick={() => setActiveTab('components')}
            className={`px-4 py-3 text-sm font-medium rounded-md transition-colors ${
              activeTab === 'components'
                ? 'bg-blue-100 text-blue-700 shadow-sm'
                : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
            }`}
          >
            <Layout className="w-4 h-4 inline mr-2" />
            Components
          </button>
          <button
            onClick={() => setActiveTab('tools')}
            className={`px-4 py-3 text-sm font-medium rounded-md transition-colors ${
              activeTab === 'tools'
                ? 'bg-blue-100 text-blue-700 shadow-sm'
                : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
            }`}
          >
            <Settings className="w-4 h-4 inline mr-2" />
            Store Tools
          </button>
        </div>
      </div>
      
      {/* Tab Content */}
      <div className="flex-1 overflow-hidden">
        {activeTab === 'quickstart' ? (
          <div className="h-full overflow-y-auto">
            <QuickStart
              onSelectTemplate={onQuickStart}
              onSelectTheme={onThemeSelect}
              selectedTheme={selectedTheme}
            />
          </div>
        ) : activeTab === 'tools' ? (
          <div className="h-full overflow-y-auto">
            <VendorFeatures onFeatureToggle={(feature, enabled) => {
              console.log(`Feature ${feature} ${enabled ? 'enabled' : 'disabled'}`);
            }} />
          </div>
        ) : (
          <div className="h-full overflow-y-auto">
            {/* Components Overview Header */}
            <div className="px-4 py-6 bg-gradient-to-br from-blue-50 to-purple-50 border-b border-gray-200">
              <div className="text-center mb-4">
                <Layout className="w-8 h-8 text-blue-600 mx-auto mb-2" />
                <h2 className="text-xl font-bold text-gray-900 mb-1">Component Library</h2>
                <p className="text-gray-600 text-sm">
                  Drag and drop components to build your custom store layout
                </p>
              </div>

              {/* Quick Stats */}
              <div className="grid grid-cols-3 gap-4 max-w-md mx-auto">
                {sections.map(section => (
                  <div key={section.id} className="text-center">
                    <div className="text-lg font-bold text-gray-900">
                      {section.components.length}
                    </div>
                    <div className="text-xs text-gray-600">
                      {section.title.split(' ')[0]}
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="py-6">
              {sections.map(section => (
                <ComponentSection key={section.id} section={section} />
              ))}

              {/* Tips Section */}
              <div className="px-4 mt-8 pt-6 border-t border-gray-200">
                <div className="bg-gradient-to-br from-blue-50 to-indigo-50 border border-blue-200 rounded-xl p-6">
                  <div className="flex items-center space-x-3 mb-4">
                    <div className="p-2 bg-blue-100 rounded-lg text-blue-600">
                      <Rocket className="w-5 h-5" />
                    </div>
                    <h3 className="text-lg font-bold text-blue-900">Building Tips</h3>
                  </div>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm text-blue-800">
                    <div className="flex items-start space-x-2">
                      <span>💡</span>
                      <span>Start with Quick Start templates for fastest setup</span>
                    </div>
                    <div className="flex items-start space-x-2">
                      <span>🎨</span>
                      <span>Apply themes before adding custom components</span>
                    </div>
                    <div className="flex items-start space-x-2">
                      <span>📱</span>
                      <span>Preview on mobile before publishing</span>
                    </div>
                    <div className="flex items-start space-x-2">
                      <span>⚡</span>
                      <span>Use drag & drop to arrange components</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ComponentLibrary;