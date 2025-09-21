import { useDroppable } from '@dnd-kit/core';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Trash2, GripVertical, Edit3, Plus } from 'lucide-react';
import { componentTemplates } from '../../../data/componentTemplates';

// Drop zone component for inserting components between existing ones
const DropZone = ({ id, index, onDrop }) => {
  const { isOver, setNodeRef } = useDroppable({
    id: `drop-zone-${id}`,
    data: {
      type: 'drop-zone',
      index: index
    }
  });

  return (
    <div
      ref={setNodeRef}
      className={`transition-all duration-200 ${
        isOver
          ? 'h-16 bg-blue-50 border-2 border-dashed border-blue-300 rounded-lg'
          : 'h-2 hover:h-8 hover:bg-gray-50 hover:border hover:border-dashed hover:border-gray-300 hover:rounded-lg'
      }`}
    >
      {isOver && (
        <div className="h-full flex items-center justify-center">
          <div className="flex items-center space-x-2 text-blue-600 text-sm font-medium">
            <Plus className="w-4 h-4" />
            <span>Drop component here</span>
          </div>
        </div>
      )}
    </div>
  );
};

const SortableComponent = ({ component, isSelected, onSelect, onDelete }) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging
  } = useSortable({
    id: component.id,
    data: {
      type: 'canvas-item',
      component
    }
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1
  };

  const template = componentTemplates.find(t => t.id === component.templateId);
  
  const renderComponent = () => {
    switch (component.type) {
      case 'header':
        return (
          <header 
            className="w-full p-4 border-b"
            style={{
              backgroundColor: component.styles?.backgroundColor || '#ffffff',
              color: component.styles?.textColor || '#000000'
            }}
          >
            <div className="flex items-center justify-between max-w-7xl mx-auto">
              <div className="flex items-center space-x-4">
                <h1 className="text-xl font-bold">{component.props?.title || component.props?.companyName || 'Your Logo'}</h1>
              </div>
              <nav className="hidden md:flex space-x-6">
                {(component.props?.navigation || component.props?.menuItems || ['Home', 'Products', 'About', 'Contact']).map((item, index) => (
                  <a key={index} href="#" className="hover:text-blue-600 transition-colors">
                    {item}
                  </a>
                ))}
              </nav>
              <div className="flex items-center space-x-4">
                <button className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors">
                  {component.props?.ctaText || component.props?.buttonText || 'Get Started'}
                </button>
              </div>
            </div>
          </header>
        );
      
      case 'footer':
        return (
          <footer 
            className="w-full p-6 border-t"
            style={{
              backgroundColor: component.styles?.backgroundColor || '#f8f9fa',
              color: component.styles?.textColor || '#6b7280'
            }}
          >
            <div className="max-w-7xl mx-auto">
              <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
                <div>
                  <h3 className="font-semibold mb-4">{component.props?.companyName || 'Company'}</h3>
                  <p className="text-sm">{component.props?.description || 'Your company description here.'}</p>
                </div>
                <div>
                  <h4 className="font-medium mb-3">Quick Links</h4>
                  <ul className="space-y-2 text-sm">
                    {(component.props?.links || component.props?.quickLinks || ['About', 'Services', 'Contact']).map((link, index) => (
                      <li key={index}>
                        <a href="#" className="hover:text-blue-600 transition-colors">{link}</a>
                      </li>
                    ))}
                  </ul>
                </div>
                <div>
                  <h4 className="font-medium mb-3">Support</h4>
                  <ul className="space-y-2 text-sm">
                    <li><a href="#" className="hover:text-blue-600 transition-colors">Help Center</a></li>
                    <li><a href="#" className="hover:text-blue-600 transition-colors">Privacy Policy</a></li>
                    <li><a href="#" className="hover:text-blue-600 transition-colors">Terms of Service</a></li>
                  </ul>
                </div>
                <div>
                  <h4 className="font-medium mb-3">Contact</h4>
                  <p className="text-sm">{component.props?.email || 'contact@company.com'}</p>
                  <p className="text-sm">{component.props?.phone || '+1 (555) 123-4567'}</p>
                </div>
              </div>
              <div className="border-t mt-8 pt-6 text-center text-sm">
                <p>&copy; 2024 {component.props?.companyName || 'Company'}. All rights reserved.</p>
              </div>
            </div>
          </footer>
        );
      
      case 'hero':
        return (
          <section
            className="w-full py-20 px-4"
            style={{
              backgroundColor: component.styles?.backgroundColor || '#f8f9fa',
              color: component.styles?.textColor || '#000000'
            }}
          >
            <div className="max-w-4xl mx-auto text-center">
              <h1 className="text-4xl md:text-6xl font-bold mb-6">
                {component.props?.title || 'Welcome to Our Platform'}
              </h1>
              <p className="text-xl mb-8 text-gray-600">
                {component.props?.subtitle || 'Build amazing experiences with our tools'}
              </p>
              <div className="space-x-4">
                <button className="bg-blue-600 text-white px-8 py-3 rounded-lg text-lg hover:bg-blue-700 transition-colors">
                  {component.props?.buttonText || component.props?.primaryCta || 'Get Started'}
                </button>
                <button className="border border-gray-300 px-8 py-3 rounded-lg text-lg hover:bg-gray-50 transition-colors">
                  {component.props?.secondaryCta || 'Learn More'}
                </button>
              </div>
            </div>
          </section>
        );

      case 'product':
        return (
          <section
            className="w-full py-16 px-4"
            style={{
              backgroundColor: component.styles?.backgroundColor || '#ffffff',
              color: component.styles?.textColor || '#000000'
            }}
          >
            <div className="max-w-7xl mx-auto">
              <h2 className="text-3xl font-bold text-center mb-12">
                {component.props?.title || 'Featured Products'}
              </h2>
              <div className={`grid gap-6 ${
                component.props?.columns === 2 ? 'grid-cols-1 md:grid-cols-2' :
                component.props?.columns === 4 ? 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4' :
                'grid-cols-1 md:grid-cols-2 lg:grid-cols-3'
              }`}>
                {(component.props?.products || []).map((product, index) => (
                  <div key={index} className="bg-white border border-gray-200 rounded-lg overflow-hidden hover:shadow-lg transition-shadow">
                    <img
                      src={product.image}
                      alt={product.name}
                      className="w-full h-48 object-cover"
                    />
                    <div className="p-4">
                      <h3 className="font-semibold text-lg mb-2">{product.name}</h3>
                      <div className="flex items-center justify-between">
                        <span className="text-xl font-bold text-blue-600">{product.price}</span>
                        {component.props?.showRatings && product.rating && (
                          <span className="text-sm text-gray-600">⭐ {product.rating}</span>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </section>
        );

      case 'category':
        return (
          <section
            className="w-full py-16 px-4 relative"
            style={{
              backgroundColor: component.styles?.backgroundColor || '#000000',
              color: component.styles?.textColor || '#ffffff',
              backgroundImage: component.props?.backgroundImage ? `url(${component.props.backgroundImage})` : 'none',
              backgroundSize: 'cover',
              backgroundPosition: 'center'
            }}
          >
            {component.props?.overlay && (
              <div className="absolute inset-0 bg-black bg-opacity-50"></div>
            )}
            <div className="max-w-4xl mx-auto text-center relative z-10">
              <h2 className="text-4xl md:text-5xl font-bold mb-4">
                {component.props?.title || 'New Collection'}
              </h2>
              <p className="text-xl mb-8">
                {component.props?.subtitle || 'Discover our latest arrivals'}
              </p>
              <button className="bg-white text-black px-8 py-3 rounded-lg text-lg hover:bg-gray-100 transition-colors">
                {component.props?.buttonText || 'Shop Now'}
              </button>
            </div>
          </section>
        );
      
      default:
        return (
          <div className="w-full p-8 bg-gray-100 border-2 border-dashed border-gray-300 text-center rounded-lg">
            <p className="text-gray-500 mb-2">Unknown component type: <strong>{component.type}</strong></p>
            <p className="text-xs text-gray-400">Template ID: {component.templateId}</p>
            <p className="text-xs text-gray-400 mt-2">Available props: {Object.keys(component.props || {}).join(', ') || 'None'}</p>
          </div>
        );
    }
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`relative group border-2 transition-all duration-200 ${
        isSelected 
          ? 'border-blue-500 shadow-lg' 
          : 'border-transparent hover:border-gray-300'
      }`}
      onClick={() => onSelect(component)}
    >
      {/* Component Content */}
      <div className="pointer-events-none">
        {renderComponent()}
      </div>
      
      {/* Controls Overlay */}
      <div className="absolute inset-0 bg-blue-500 bg-opacity-10 opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none" />
      
      {/* Control Buttons */}
      <div className="absolute top-2 right-2 flex space-x-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
        <button
          onClick={(e) => {
            e.stopPropagation();
            onSelect(component);
          }}
          className="p-2 bg-white rounded-md shadow-md hover:bg-gray-50 transition-colors pointer-events-auto"
          title="Edit Component"
        >
          <Edit3 className="w-4 h-4 text-gray-600" />
        </button>
        <button
          onClick={(e) => {
            e.stopPropagation();
            onDelete(component.id);
          }}
          className="p-2 bg-white rounded-md shadow-md hover:bg-red-50 transition-colors pointer-events-auto"
          title="Delete Component"
        >
          <Trash2 className="w-4 h-4 text-red-600" />
        </button>
      </div>
      
      {/* Drag Handle */}
      <div 
        {...attributes}
        {...listeners}
        className="absolute top-2 left-2 p-2 bg-white rounded-md shadow-md opacity-0 group-hover:opacity-100 transition-opacity duration-200 cursor-grab active:cursor-grabbing pointer-events-auto"
        title="Drag to Reorder"
      >
        <GripVertical className="w-4 h-4 text-gray-600" />
      </div>
      
      {/* Component Label */}
      <div className="absolute bottom-2 left-2 px-2 py-1 bg-white rounded-md shadow-md opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none">
        <span className="text-xs font-medium text-gray-700">
          {template?.name || component.type}
        </span>
      </div>
    </div>
  );
};

const CanvasArea = ({ components, selectedComponent, onComponentSelect, onComponentDelete, selectedTheme }) => {
  const { isOver, setNodeRef, active } = useDroppable({
    id: 'canvas'
  });

  // Debug logging for development
  if (process.env.NODE_ENV === 'development' && (isOver || active)) {
    console.log('🎯 Canvas area - isOver:', isOver, 'active:', active?.id, 'activeType:', active?.data?.current?.type);
  }

  const isThemeDrag = active?.data?.current?.type === 'theme-item';
  const isComponentDrag = active?.data?.current?.type === 'library-item';

  const canvasStyle = selectedTheme ? {
    backgroundColor: selectedTheme.colors?.background || '#ffffff',
    fontFamily: selectedTheme.typography?.primary || 'inherit'
  } : {};
  
  const responsiveClasses = selectedTheme ? `theme-${selectedTheme.id}` : '';

  return (
    <div className="h-full bg-gray-50">
      <div className="p-4 border-b border-gray-200 bg-white">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-gray-900">Store Preview</h2>
            <p className="text-sm text-gray-600 mt-1">
              {components.length === 0
                ? 'Choose a theme or drag components to start building'
                : `${components.length} component${components.length !== 1 ? 's' : ''} in your store`
              }
            </p>
          </div>
          <div className="flex items-center space-x-2">
            {selectedTheme && (
              <div className="flex items-center space-x-2">
                <div className="flex space-x-1">
                  <div
                    className="w-3 h-3 rounded-full border border-gray-300"
                    style={{ backgroundColor: selectedTheme.colors.primary }}
                    title="Primary Color"
                  />
                  <div
                    className="w-3 h-3 rounded-full border border-gray-300"
                    style={{ backgroundColor: selectedTheme.colors.secondary }}
                    title="Secondary Color"
                  />
                  <div
                    className="w-3 h-3 rounded-full border border-gray-300"
                    style={{ backgroundColor: selectedTheme.colors.accent }}
                    title="Accent Color"
                  />
                </div>
                <span className="px-3 py-1 bg-purple-100 text-purple-800 text-xs font-medium rounded-full">
                  🎨 {selectedTheme.name}
                </span>
              </div>
            )}
          </div>
        </div>
      </div>
      
      <div
        ref={setNodeRef}
        data-id="canvas"
        className={`min-h-full transition-all duration-300 p-2 sm:p-4 ${responsiveClasses} ${
          isOver && isThemeDrag ? 'bg-purple-50/30 border-2 border-dashed border-purple-300 scale-[0.98]' :
          isOver && isComponentDrag ? 'bg-blue-50/30 border-2 border-dashed border-blue-300 scale-[0.98]' :
          isOver ? 'bg-gray-50/20' : ''
        }`}
        style={isOver ? {} : canvasStyle}
      >
        {/* Drop Feedback - Floating badge only */}
        {isOver && isThemeDrag && (
          <div className="absolute top-4 left-1/2 transform -translate-x-1/2 pointer-events-none z-20">
            <div className="bg-purple-600 text-white px-4 py-2 rounded-full shadow-lg border-2 border-purple-200 animate-pulse">
              <div className="flex items-center space-x-2">
                <span>🎨</span>
                <span className="text-sm font-medium">Drop to apply theme</span>
              </div>
            </div>
          </div>
        )}
        {isOver && isComponentDrag && (
          <div className="absolute top-4 left-1/2 transform -translate-x-1/2 pointer-events-none z-20">
            <div className="bg-blue-600 text-white px-4 py-2 rounded-full shadow-lg border-2 border-blue-200 animate-pulse">
              <div className="flex items-center space-x-2">
                <span>📦</span>
                <span className="text-sm font-medium">Drop to add component</span>
              </div>
            </div>
          </div>
        )}
        {components.length === 0 ? (
          <div className="flex items-center justify-center h-96">
            <div className="text-center max-w-lg">
              <div className="w-24 h-24 bg-gradient-to-br from-blue-100 to-purple-100 rounded-full flex items-center justify-center mx-auto mb-6">
                <span className="text-4xl">🏪</span>
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-3">Create Your Store Design</h3>
              <p className="text-gray-600 mb-8">
                Build a professional storefront that converts visitors into customers.
              </p>

              <div className="bg-blue-50 border border-blue-200 rounded-lg p-6 mb-6">
                <h4 className="font-semibold text-blue-900 mb-3">🎨 Quick Start Options:</h4>
                <div className="space-y-3 text-sm">
                  <div className="flex items-center space-x-3 text-blue-800">
                    <span className="w-6 h-6 bg-blue-600 text-white rounded-full flex items-center justify-center text-xs font-bold">1</span>
                    <span><strong>Choose a Theme:</strong> Click "Apply Theme" on any theme template</span>
                  </div>
                  <div className="flex items-center space-x-3 text-blue-800">
                    <span className="w-6 h-6 bg-blue-600 text-white rounded-full flex items-center justify-center text-xs font-bold">2</span>
                    <span><strong>Drag Components:</strong> Build from scratch using the component library</span>
                  </div>
                  <div className="flex items-center space-x-3 text-blue-800">
                    <span className="w-6 h-6 bg-blue-600 text-white rounded-full flex items-center justify-center text-xs font-bold">3</span>
                    <span><strong>Customize:</strong> Modify colors, text, and layout to match your brand</span>
                  </div>
                </div>
              </div>

              <div className="text-xs text-gray-500">
                💪 Professional templates optimized for e-commerce conversions
              </div>
            </div>
          </div>
        ) : (
          <div className="space-y-0">
            {/* Drop zone at the beginning */}
            <DropZone id="start" index={0} />

            {components.map((component, index) => (
              <div key={component.id}>
                <SortableComponent
                  component={component}
                  isSelected={selectedComponent?.id === component.id}
                  onSelect={onComponentSelect}
                  onDelete={onComponentDelete}
                />
                {/* Drop zone after each component */}
                <DropZone id={`after-${component.id}`} index={index + 1} />
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default CanvasArea;