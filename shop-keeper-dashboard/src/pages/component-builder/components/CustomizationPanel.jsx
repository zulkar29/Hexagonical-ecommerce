import { useState } from 'react';
import { Palette, Type, Settings, ChevronDown, ChevronRight } from 'lucide-react';

const CustomizationPanel = ({ component, onUpdate }) => {
  const [activeTab, setActiveTab] = useState('content');
  const [expandedSections, setExpandedSections] = useState({
    content: true,
    styling: true,
    layout: true
  });

  const toggleSection = (section) => {
    setExpandedSections(prev => ({
      ...prev,
      [section]: !prev[section]
    }));
  };

  const handlePropChange = (key, value) => {
    onUpdate({
      props: {
        ...component.props,
        [key]: value
      }
    });
  };

  const handleStyleChange = (key, value) => {
    onUpdate({
      styles: {
        ...component.styles,
        [key]: value
      }
    });
  };

  const handleArrayPropChange = (key, index, value) => {
    const currentArray = component.props?.[key] || [];
    const newArray = [...currentArray];
    newArray[index] = value;
    handlePropChange(key, newArray);
  };

  const addArrayItem = (key, defaultValue = '') => {
    const currentArray = component.props?.[key] || [];
    handlePropChange(key, [...currentArray, defaultValue]);
  };

  const removeArrayItem = (key, index) => {
    const currentArray = component.props?.[key] || [];
    const newArray = currentArray.filter((_, i) => i !== index);
    handlePropChange(key, newArray);
  };

  const renderContentFields = () => {
    switch (component.type) {
      case 'header':
        return (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Logo/Title</label>
              <input
                type="text"
                value={component.props?.title || ''}
                onChange={(e) => handlePropChange('title', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Your Logo"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">CTA Button Text</label>
              <input
                type="text"
                value={component.props?.ctaText || ''}
                onChange={(e) => handlePropChange('ctaText', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Get Started"
              />
            </div>
            
            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-sm font-medium text-gray-700">Menu Items</label>
                <button
                  onClick={() => addArrayItem('menuItems', 'New Item')}
                  className="text-sm text-blue-600 hover:text-blue-700"
                >
                  + Add Item
                </button>
              </div>
              {(component.props?.menuItems || ['Home', 'Products', 'About', 'Contact']).map((item, index) => (
                <div key={index} className="flex items-center space-x-2 mb-2">
                  <input
                    type="text"
                    value={item}
                    onChange={(e) => handleArrayPropChange('menuItems', index, e.target.value)}
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                  <button
                    onClick={() => removeArrayItem('menuItems', index)}
                    className="text-red-600 hover:text-red-700 text-sm"
                  >
                    Remove
                  </button>
                </div>
              ))}
            </div>
          </div>
        );
      
      case 'footer':
        return (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Company Name</label>
              <input
                type="text"
                value={component.props?.companyName || ''}
                onChange={(e) => handlePropChange('companyName', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Company"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Description</label>
              <textarea
                value={component.props?.description || ''}
                onChange={(e) => handlePropChange('description', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                rows={3}
                placeholder="Your company description here."
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Email</label>
              <input
                type="email"
                value={component.props?.email || ''}
                onChange={(e) => handlePropChange('email', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="contact@company.com"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Phone</label>
              <input
                type="tel"
                value={component.props?.phone || ''}
                onChange={(e) => handlePropChange('phone', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="+1 (555) 123-4567"
              />
            </div>
            
            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-sm font-medium text-gray-700">Quick Links</label>
                <button
                  onClick={() => addArrayItem('quickLinks', 'New Link')}
                  className="text-sm text-blue-600 hover:text-blue-700"
                >
                  + Add Link
                </button>
              </div>
              {(component.props?.quickLinks || ['About', 'Services', 'Contact']).map((link, index) => (
                <div key={index} className="flex items-center space-x-2 mb-2">
                  <input
                    type="text"
                    value={link}
                    onChange={(e) => handleArrayPropChange('quickLinks', index, e.target.value)}
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                  <button
                    onClick={() => removeArrayItem('quickLinks', index)}
                    className="text-red-600 hover:text-red-700 text-sm"
                  >
                    Remove
                  </button>
                </div>
              ))}
            </div>
          </div>
        );
      
      case 'hero':
        return (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Title</label>
              <input
                type="text"
                value={component.props?.title || ''}
                onChange={(e) => handlePropChange('title', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Welcome to Our Platform"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Subtitle</label>
              <textarea
                value={component.props?.subtitle || ''}
                onChange={(e) => handlePropChange('subtitle', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                rows={2}
                placeholder="Build amazing experiences with our tools"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Primary CTA</label>
              <input
                type="text"
                value={component.props?.primaryCta || ''}
                onChange={(e) => handlePropChange('primaryCta', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Get Started"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Secondary CTA</label>
              <input
                type="text"
                value={component.props?.secondaryCta || ''}
                onChange={(e) => handlePropChange('secondaryCta', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Learn More"
              />
            </div>
          </div>
        );
      
      default:
        return (
          <div className="text-center py-8 text-gray-500">
            No customization options available for this component type.
          </div>
        );
    }
  };

  const renderStyleFields = () => (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Background Color</label>
        <div className="flex items-center space-x-2">
          <input
            type="color"
            value={component.styles?.backgroundColor || '#ffffff'}
            onChange={(e) => handleStyleChange('backgroundColor', e.target.value)}
            className="w-12 h-10 border border-gray-300 rounded cursor-pointer"
          />
          <input
            type="text"
            value={component.styles?.backgroundColor || '#ffffff'}
            onChange={(e) => handleStyleChange('backgroundColor', e.target.value)}
            className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="#ffffff"
          />
        </div>
      </div>
      
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Text Color</label>
        <div className="flex items-center space-x-2">
          <input
            type="color"
            value={component.styles?.textColor || '#000000'}
            onChange={(e) => handleStyleChange('textColor', e.target.value)}
            className="w-12 h-10 border border-gray-300 rounded cursor-pointer"
          />
          <input
            type="text"
            value={component.styles?.textColor || '#000000'}
            onChange={(e) => handleStyleChange('textColor', e.target.value)}
            className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="#000000"
          />
        </div>
      </div>
    </div>
  );

  const Section = ({ title, icon: Icon, children, sectionKey }) => (
    <div className="border-b border-gray-200 last:border-b-0">
      <button
        onClick={() => toggleSection(sectionKey)}
        className="w-full flex items-center justify-between p-4 text-left hover:bg-gray-50 transition-colors"
      >
        <div className="flex items-center space-x-2">
          <Icon className="w-4 h-4 text-gray-600" />
          <span className="font-medium text-gray-900">{title}</span>
        </div>
        {expandedSections[sectionKey] ? (
          <ChevronDown className="w-4 h-4 text-gray-600" />
        ) : (
          <ChevronRight className="w-4 h-4 text-gray-600" />
        )}
      </button>
      {expandedSections[sectionKey] && (
        <div className="px-4 pb-4">
          {children}
        </div>
      )}
    </div>
  );

  if (!component) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="text-center">
          <Settings className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">No Component Selected</h3>
          <p className="text-gray-600 max-w-sm">
            Select a component from the canvas to customize its properties and styling.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-gray-200">
        <h2 className="text-lg font-semibold text-gray-900">Customize Component</h2>
        <p className="text-sm text-gray-600 mt-1">
          {component.name || component.type}
        </p>
      </div>
      
      {/* Content */}
      <div className="flex-1 overflow-y-auto">
        <Section title="Content" icon={Type} sectionKey="content">
          {renderContentFields()}
        </Section>
        
        <Section title="Styling" icon={Palette} sectionKey="styling">
          {renderStyleFields()}
        </Section>
      </div>
    </div>
  );
};

export default CustomizationPanel;