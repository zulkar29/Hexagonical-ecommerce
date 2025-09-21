import React, { useState } from 'react';
import {
  Type,
  Palette,
  Layout,
  Settings,
  X,
  ChevronDown,
  ChevronRight,
  Eye,
  EyeOff
} from 'lucide-react';

const PropertiesPanel = ({ component, onUpdate, onClose }) => {
  const [expandedSections, setExpandedSections] = useState({
    content: true,
    styling: true,
    layout: false,
    advanced: false
  });

  const toggleSection = (section) => {
    setExpandedSections(prev => ({
      ...prev,
      [section]: !prev[section]
    }));
  };

  const handlePropertyChange = (path, value) => {
    if (path.startsWith('props.')) {
      const propPath = path.substring(6);
      onUpdate({
        props: {
          ...component.props,
          [propPath]: value
        }
      });
    } else if (path.startsWith('styles.')) {
      const stylePath = path.substring(7);
      onUpdate({
        styles: {
          ...component.styles,
          [stylePath]: value
        }
      });
    }
  };

  const Section = ({ title, icon, children, sectionKey }) => {
    const Icon = icon;
    return (
    <div className="border-b border-gray-200 last:border-b-0">
      <button
        onClick={() => toggleSection(sectionKey)}
        className="w-full flex items-center justify-between p-3 text-left hover:bg-gray-50 transition-colors"
      >
        <div className="flex items-center space-x-2">
          <Icon className="w-4 h-4 text-gray-600" />
          <span className="font-medium text-gray-900">{title}</span>
        </div>
        {expandedSections[sectionKey] ? (
          <ChevronDown className="w-4 h-4 text-gray-500" />
        ) : (
          <ChevronRight className="w-4 h-4 text-gray-500" />
        )}
      </button>
      {expandedSections[sectionKey] && (
        <div className="px-3 pb-3 space-y-3">
          {children}
        </div>
      )}
    </div>
  );
  };

  const InputField = ({ label, value, onChange, type = 'text', placeholder, options }) => (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-1">
        {label}
      </label>
      {type === 'select' ? (
        <select
          value={value || ''}
          onChange={(e) => onChange(e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        >
          {options?.map(option => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      ) : type === 'textarea' ? (
        <textarea
          value={value || ''}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          rows={3}
          className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        />
      ) : type === 'color' ? (
        <div className="flex space-x-2">
          <input
            type="color"
            value={value || '#000000'}
            onChange={(e) => onChange(e.target.value)}
            className="w-12 h-8 border border-gray-300 rounded cursor-pointer"
          />
          <input
            type="text"
            value={value || ''}
            onChange={(e) => onChange(e.target.value)}
            placeholder="#000000"
            className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
      ) : (
        <input
          type={type}
          value={value || ''}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        />
      )}
    </div>
  );

  const renderContentFields = () => {
    const props = component.props || {};

    switch (component.type) {
      case 'header':
        return (
          <>
            <InputField
              label="Store Name"
              value={props.storeName}
              onChange={(value) => handlePropertyChange('props.storeName', value)}
              placeholder="Your Store"
            />
            <InputField
              label="Show Search"
              value={props.showSearch}
              onChange={(value) => handlePropertyChange('props.showSearch', value)}
              type="select"
              options={[
                { value: true, label: 'Show' },
                { value: false, label: 'Hide' }
              ]}
            />
            <InputField
              label="Show Cart"
              value={props.showCart}
              onChange={(value) => handlePropertyChange('props.showCart', value)}
              type="select"
              options={[
                { value: true, label: 'Show' },
                { value: false, label: 'Hide' }
              ]}
            />
          </>
        );

      case 'hero':
        return (
          <>
            <InputField
              label="Title"
              value={props.title}
              onChange={(value) => handlePropertyChange('props.title', value)}
              placeholder="Welcome to Our Store"
            />
            <InputField
              label="Subtitle"
              value={props.subtitle}
              onChange={(value) => handlePropertyChange('props.subtitle', value)}
              type="textarea"
              placeholder="Discover amazing products..."
            />
            <InputField
              label="Button Text"
              value={props.primaryCta}
              onChange={(value) => handlePropertyChange('props.primaryCta', value)}
              placeholder="Shop Now"
            />
            <InputField
              label="Button Link"
              value={props.ctaLink}
              onChange={(value) => handlePropertyChange('props.ctaLink', value)}
              placeholder="/products"
            />
          </>
        );

      case 'product':
      case 'product-grid':
        return (
          <>
            <InputField
              label="Section Title"
              value={props.title}
              onChange={(value) => handlePropertyChange('props.title', value)}
              placeholder="Featured Products"
            />
            <InputField
              label="Products to Show"
              value={props.limit}
              onChange={(value) => handlePropertyChange('props.limit', parseInt(value))}
              type="number"
              placeholder="6"
            />
            <InputField
              label="Show Prices"
              value={props.showPrices}
              onChange={(value) => handlePropertyChange('props.showPrices', value)}
              type="select"
              options={[
                { value: true, label: 'Show' },
                { value: false, label: 'Hide' }
              ]}
            />
          </>
        );

      default:
        return (
          <div className="text-sm text-gray-500">
            No content options available for this component type.
          </div>
        );
    }
  };

  const renderStyleFields = () => {
    const styles = component.styles || {};

    return (
      <>
        <InputField
          label="Background Color"
          value={styles.backgroundColor}
          onChange={(value) => handlePropertyChange('styles.backgroundColor', value)}
          type="color"
        />
        <InputField
          label="Text Color"
          value={styles.textColor}
          onChange={(value) => handlePropertyChange('styles.textColor', value)}
          type="color"
        />
        <InputField
          label="Padding"
          value={styles.padding}
          onChange={(value) => handlePropertyChange('styles.padding', value)}
          placeholder="16px"
        />
        <InputField
          label="Margin"
          value={styles.margin}
          onChange={(value) => handlePropertyChange('styles.margin', value)}
          placeholder="0px"
        />
        <InputField
          label="Border Radius"
          value={styles.borderRadius}
          onChange={(value) => handlePropertyChange('styles.borderRadius', value)}
          placeholder="8px"
        />
      </>
    );
  };

  const renderLayoutFields = () => {
    const styles = component.styles || {};

    return (
      <>
        <InputField
          label="Width"
          value={styles.width}
          onChange={(value) => handlePropertyChange('styles.width', value)}
          placeholder="100%"
        />
        <InputField
          label="Height"
          value={styles.height}
          onChange={(value) => handlePropertyChange('styles.height', value)}
          placeholder="auto"
        />
        <InputField
          label="Display"
          value={styles.display}
          onChange={(value) => handlePropertyChange('styles.display', value)}
          type="select"
          options={[
            { value: 'block', label: 'Block' },
            { value: 'flex', label: 'Flex' },
            { value: 'grid', label: 'Grid' },
            { value: 'inline-block', label: 'Inline Block' }
          ]}
        />
        <InputField
          label="Text Align"
          value={styles.textAlign}
          onChange={(value) => handlePropertyChange('styles.textAlign', value)}
          type="select"
          options={[
            { value: 'left', label: 'Left' },
            { value: 'center', label: 'Center' },
            { value: 'right', label: 'Right' }
          ]}
        />
      </>
    );
  };

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-gray-900">Properties</h2>
            <p className="text-sm text-gray-600">{component.name || component.type}</p>
          </div>
          <button
            onClick={onClose}
            className="p-1 text-gray-500 hover:text-gray-700 rounded"
          >
            <X className="w-4 h-4" />
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto">
        <Section title="Content" icon={Type} sectionKey="content">
          {renderContentFields()}
        </Section>

        <Section title="Styling" icon={Palette} sectionKey="styling">
          {renderStyleFields()}
        </Section>

        <Section title="Layout" icon={Layout} sectionKey="layout">
          {renderLayoutFields()}
        </Section>

        <Section title="Advanced" icon={Settings} sectionKey="advanced">
          <InputField
            label="Custom CSS Class"
            value={component.styles?.customClass}
            onChange={(value) => handlePropertyChange('styles.customClass', value)}
            placeholder="my-custom-class"
          />
          <InputField
            label="Component ID"
            value={component.id}
            onChange={(value) => onUpdate({ id: value })}
            placeholder="unique-id"
          />
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium text-gray-700">Visible</span>
            <button
              onClick={() => handlePropertyChange('styles.display',
                component.styles?.display === 'none' ? 'block' : 'none'
              )}
              className={`p-2 rounded-lg transition-colors ${
                component.styles?.display === 'none'
                  ? 'bg-red-100 text-red-600'
                  : 'bg-green-100 text-green-600'
              }`}
            >
              {component.styles?.display === 'none' ? (
                <EyeOff className="w-4 h-4" />
              ) : (
                <Eye className="w-4 h-4" />
              )}
            </button>
          </div>
        </Section>
      </div>
    </div>
  );
};

export default PropertiesPanel;