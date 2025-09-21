import { useState, useEffect } from 'react';
import { useForm, useFieldArray } from 'react-hook-form';
import { Palette, Type, Settings, ChevronDown, ChevronRight } from 'lucide-react';

const CustomizationPanel = ({ component, onUpdate }) => {
  const [expandedSections, setExpandedSections] = useState({
    content: true,
    styling: true,
    layout: true
  });

  // Initialize React Hook Form with component data
  const {
    register,
    control,
    watch,
    setValue,
    reset,
    handleSubmit
  } = useForm({
    defaultValues: {
      props: component?.props || {},
      styles: component?.styles || {}
    },
    mode: 'onChange'
  });

  // Field arrays for dynamic lists (menu items, quick links, etc.)
  const menuItemsArray = useFieldArray({
    control,
    name: 'props.menuItems'
  });

  const quickLinksArray = useFieldArray({
    control,
    name: 'props.quickLinks'
  });

  // Watch form changes and update component
  const formData = watch();

  useEffect(() => {
    if (component) {
      reset({
        props: component.props || {},
        styles: component.styles || {}
      });
    }
  }, [component?.id, reset]);

  // Debounced update on form changes
  useEffect(() => {
    const subscription = watch((value) => {
      const timeoutId = setTimeout(() => {
        if (value.props || value.styles) {
          const updates = {};
          if (value.props) updates.props = value.props;
          if (value.styles) updates.styles = value.styles;
          onUpdate(updates);
        }
      }, 300);

      return () => clearTimeout(timeoutId);
    });

    return () => subscription.unsubscribe();
  }, [watch, onUpdate]);

  const toggleSection = (section) => {
    setExpandedSections(prev => ({
      ...prev,
      [section]: !prev[section]
    }));
  };

  const renderContentFields = () => {
    switch (component.type) {
      case 'header':
        return (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Logo/Title</label>
              <input
                {...register('props.title')}
                type="text"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Your Logo"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">CTA Button Text</label>
              <input
                {...register('props.ctaText')}
                type="text"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Get Started"
              />
            </div>

            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-sm font-medium text-gray-700">Menu Items</label>
                <button
                  type="button"
                  onClick={() => menuItemsArray.append('New Item')}
                  className="text-sm text-blue-600 hover:text-blue-700"
                >
                  + Add Item
                </button>
              </div>
              {menuItemsArray.fields.map((field, index) => (
                <div key={field.id} className="flex items-center space-x-2 mb-2">
                  <input
                    {...register(`props.menuItems.${index}`)}
                    type="text"
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                  <button
                    type="button"
                    onClick={() => menuItemsArray.remove(index)}
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
                {...register('props.companyName')}
                type="text"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Company"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Description</label>
              <textarea
                {...register('props.description')}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                rows={3}
                placeholder="Your company description here."
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Email</label>
              <input
                {...register('props.email')}
                type="email"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="contact@company.com"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Phone</label>
              <input
                {...register('props.phone')}
                type="tel"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="+1 (555) 123-4567"
              />
            </div>
            
            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-sm font-medium text-gray-700">Quick Links</label>
                <button
                  type="button"
                  onClick={() => quickLinksArray.append('New Link')}
                  className="text-sm text-blue-600 hover:text-blue-700"
                >
                  + Add Link
                </button>
              </div>
              {quickLinksArray.fields.map((field, index) => (
                <div key={field.id} className="flex items-center space-x-2 mb-2">
                  <input
                    {...register(`props.quickLinks.${index}`)}
                    type="text"
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                  <button
                    type="button"
                    onClick={() => quickLinksArray.remove(index)}
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
                {...register('props.title')}
                type="text"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Welcome to Our Platform"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Subtitle</label>
              <textarea
                {...register('props.subtitle')}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                rows={2}
                placeholder="Build amazing experiences with our tools"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Primary CTA</label>
              <input
                {...register('props.primaryCta')}
                type="text"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Get Started"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Secondary CTA</label>
              <input
                {...register('props.secondaryCta')}
                type="text"
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
            {...register('styles.backgroundColor')}
            type="color"
            className="w-12 h-10 border border-gray-300 rounded cursor-pointer"
          />
          <input
            {...register('styles.backgroundColor')}
            type="text"
            className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="#ffffff"
          />
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Text Color</label>
        <div className="flex items-center space-x-2">
          <input
            {...register('styles.textColor')}
            type="color"
            className="w-12 h-10 border border-gray-300 rounded cursor-pointer"
          />
          <input
            {...register('styles.textColor')}
            type="text"
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