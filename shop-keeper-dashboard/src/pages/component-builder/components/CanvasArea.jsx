import { useDroppable } from '@dnd-kit/core';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Trash2, GripVertical, Edit3 } from 'lucide-react';
import { componentTemplates } from '../../../data/componentTemplates';

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
      
      default:
        return (
          <div className="w-full p-8 bg-gray-100 border-2 border-dashed border-gray-300 text-center">
            <p className="text-gray-500">Unknown component type: {component.type}</p>
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
  const { isOver, setNodeRef } = useDroppable({
    id: 'canvas'
  });

  const canvasStyle = selectedTheme ? {
    backgroundColor: selectedTheme.colors?.background || '#ffffff',
    fontFamily: selectedTheme.typography?.primary || 'inherit'
  } : {};
  
  const responsiveClasses = selectedTheme ? `theme-${selectedTheme.id}` : '';

  return (
    <div className="h-full bg-gray-50">
      <div className="p-4 border-b border-gray-200 bg-white">
        <h2 className="text-lg font-semibold text-gray-900">Canvas</h2>
        <p className="text-sm text-gray-600 mt-1">
          {components.length === 0 
            ? 'Drop components here to start building' 
            : `${components.length} component${components.length !== 1 ? 's' : ''} added`
          }
          {selectedTheme && (
            <span className="ml-2 px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">
              Theme: {selectedTheme.name}
            </span>
          )}
        </p>
      </div>
      
      <div 
        ref={setNodeRef}
        className={`min-h-full transition-colors duration-200 p-2 sm:p-4 ${responsiveClasses} ${
          isOver ? 'bg-blue-50' : ''
        }`}
        style={isOver ? {} : canvasStyle}
      >
        {components.length === 0 ? (
          <div className="flex items-center justify-center h-96">
            <div className="text-center">
              <div className="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center mx-auto mb-4">
                <span className="text-2xl text-gray-400">📱</span>
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">Start Building</h3>
              <p className="text-gray-600 max-w-sm">
                Drag components from the library to create your layout. 
                You can reorder them by dragging the grip handle.
              </p>
            </div>
          </div>
        ) : (
          <div className="space-y-2 sm:space-y-4">
            {components.map(component => (
              <SortableComponent
                key={component.id}
                component={component}
                isSelected={selectedComponent?.id === component.id}
                onSelect={onComponentSelect}
                onDelete={onComponentDelete}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default CanvasArea;