
const PreviewPanel = ({ components }) => {
  const renderComponent = (component) => {
    switch (component.type) {
      case 'header':
        return (
          <header 
            key={component.id}
            className="w-full p-4 border-b"
            style={{
              backgroundColor: component.styles?.backgroundColor || '#ffffff',
              color: component.styles?.textColor || '#000000'
            }}
          >
            <div className="flex items-center justify-between max-w-7xl mx-auto">
              <div className="flex items-center space-x-4">
                <h1 className="text-xl font-bold">{component.props?.title || 'Your Logo'}</h1>
              </div>
              <nav className="hidden md:flex space-x-6">
                {(component.props?.menuItems || ['Home', 'Products', 'About', 'Contact']).map((item, index) => (
                  <a key={index} href={typeof item === 'object' ? item.href || '#' : '#'} className="hover:text-blue-600 transition-colors">
                    {typeof item === 'object' ? item.label : item}
                  </a>
                ))}
              </nav>
              <div className="flex items-center space-x-4">
                <button className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors">
                  {component.props?.ctaText || 'Get Started'}
                </button>
              </div>
            </div>
          </header>
        );
      
      case 'footer':
        return (
          <footer 
            key={component.id}
            className="w-full p-6 border-t mt-auto"
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
                    {(component.props?.quickLinks || ['About', 'Services', 'Contact']).map((link, index) => (
                      <li key={index}>
                        <a href={typeof link === 'object' ? link.href || '#' : '#'} className="hover:text-blue-600 transition-colors">
                          {typeof link === 'object' ? link.label : link}
                        </a>
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
            key={component.id}
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
                  {component.props?.primaryCta || 'Get Started'}
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
          <div key={component.id} className="w-full p-8 bg-gray-100 border-2 border-dashed border-gray-300 text-center">
            <p className="text-gray-500">Unknown component type: {component.type}</p>
          </div>
        );
    }
  };

  return (
    <div className="h-full bg-white">
      {/* Preview Header */}
      <div className="p-4 border-b border-gray-200 bg-gray-50">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-gray-900">Preview</h2>
            <p className="text-sm text-gray-600 mt-1">
              This is how your layout will look to visitors
            </p>
          </div>
          <div className="flex items-center space-x-2">
            <div className="flex items-center space-x-1">
              <div className="w-3 h-3 bg-red-500 rounded-full"></div>
              <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
              <div className="w-3 h-3 bg-green-500 rounded-full"></div>
            </div>
            <div className="text-xs text-gray-500 ml-4">Browser Preview</div>
          </div>
        </div>
      </div>
      
      {/* Preview Content */}
      <div className="h-full overflow-y-auto">
        {components.length === 0 ? (
          <div className="flex items-center justify-center h-96">
            <div className="text-center">
              <div className="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center mx-auto mb-4">
                <span className="text-2xl text-gray-400">👁️</span>
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">No Components Yet</h3>
              <p className="text-gray-600 max-w-sm">
                Add some components to see the preview. Switch back to edit mode to start building.
              </p>
            </div>
          </div>
        ) : (
          <div className="min-h-full flex flex-col">
            {components.map(component => renderComponent(component))}
          </div>
        )}
      </div>
    </div>
  );
};

export default PreviewPanel;