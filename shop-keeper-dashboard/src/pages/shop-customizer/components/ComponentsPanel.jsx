import React, { useState } from 'react';
import { useDraggable } from '@dnd-kit/core';
import { componentTemplates } from '../../../data/shopTemplates';
import {
  Search,
  Grid,
  List,
  ChevronDown,
  ChevronRight,
  Star,
  ShoppingCart,
  User,
  Package,
  Layout,
  Type,
  Image as ImageIcon,
  Mail,
  MessageSquare,
  Zap,
  CreditCard,
  Heart,
  Filter,
  Tag,
  Users
} from 'lucide-react';

const ComponentPreview = ({ template }) => {
  const renderPreview = () => {
    switch (template.type) {
      case 'header':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="flex items-center justify-between text-xs">
              <div className="font-semibold text-gray-800">Store</div>
              <div className="flex items-center space-x-2">
                <Search size={10} className="text-gray-500" />
                <ShoppingCart size={10} className="text-gray-500" />
                <User size={10} className="text-gray-500" />
              </div>
            </div>
          </div>
        );
      case 'hero':
        return (
          <div className="bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-md p-2">
            <div className="text-center">
              <div className="text-xs font-bold mb-1">Hero Title</div>
              <div className="bg-white text-blue-600 px-2 py-0.5 rounded text-xs inline-block">
                Shop Now
              </div>
            </div>
          </div>
        );
      case 'product':
      case 'product-grid':
      case 'product-card':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="grid grid-cols-2 gap-1">
              {[1, 2, 3, 4].map(i => (
                <div key={i} className="border border-gray-100 rounded p-1">
                  <div className="bg-gray-100 h-6 rounded mb-1"></div>
                  <div className="text-xs text-gray-800">Product {i}</div>
                  <div className="text-xs text-blue-600 font-bold">$99</div>
                </div>
              ))}
            </div>
          </div>
        );
      case 'cart':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="flex items-center justify-between text-xs mb-1">
              <span className="font-semibold">Cart (3)</span>
              <ShoppingCart size={10} className="text-blue-600" />
            </div>
            <div className="space-y-1">
              {[1, 2].map(i => (
                <div key={i} className="flex items-center space-x-1">
                  <div className="w-3 h-3 bg-gray-100 rounded"></div>
                  <div className="flex-1 text-xs">Item {i}</div>
                  <div className="text-xs font-bold">$29</div>
                </div>
              ))}
            </div>
            <div className="border-t mt-1 pt-1 text-xs font-bold">Total: $58</div>
          </div>
        );
      case 'checkout':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="text-xs font-semibold mb-1">Checkout</div>
            <div className="space-y-1">
              <div className="flex space-x-1">
                <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                <div className="w-2 h-2 bg-gray-300 rounded-full"></div>
                <div className="w-2 h-2 bg-gray-300 rounded-full"></div>
              </div>
              <div className="text-xs text-gray-600">Step 1 of 3</div>
              <div className="bg-gray-100 h-4 rounded text-xs flex items-center px-1">Form</div>
            </div>
          </div>
        );
      case 'promotion':
        return (
          <div className="bg-red-500 text-white rounded-md p-2">
            <div className="text-center">
              <div className="text-xs font-bold">SALE 50% OFF</div>
              <div className="text-xs opacity-90">Limited Time</div>
            </div>
          </div>
        );
      case 'reviews':
      case 'testimonials':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="text-xs font-semibold mb-1">Reviews</div>
            <div className="space-y-1">
              {[1, 2].map(i => (
                <div key={i} className="flex items-center space-x-1">
                  <div className="flex">
                    {[1, 2, 3, 4, 5].map(star => (
                      <div key={star} className="w-1 h-1 bg-yellow-400 rounded-full"></div>
                    ))}
                  </div>
                  <div className="text-xs text-gray-600">Great!</div>
                </div>
              ))}
            </div>
          </div>
        );
      case 'shipping':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="flex items-center space-x-1 mb-1">
              <Package size={10} className="text-blue-600" />
              <div className="text-xs font-semibold">Shipping</div>
            </div>
            <div className="text-xs text-gray-600">ZIP: _____</div>
            <div className="text-xs text-green-600">Free shipping available</div>
          </div>
        );
      case 'search':
      case 'filter':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="flex items-center space-x-1 mb-1">
              <Search size={10} className="text-gray-500" />
              <div className="bg-gray-100 h-2 rounded flex-1"></div>
            </div>
            <div className="flex space-x-1">
              <div className="text-xs bg-blue-100 text-blue-800 px-1 rounded">Filter</div>
              <div className="text-xs bg-gray-100 text-gray-600 px-1 rounded">Sort</div>
            </div>
          </div>
        );
      case 'wishlist':
      case 'recent':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="flex items-center space-x-1 mb-1">
              <Heart size={10} className="text-red-500" />
              <div className="text-xs font-semibold">Saved Items</div>
            </div>
            <div className="grid grid-cols-3 gap-1">
              {[1, 2, 3].map(i => (
                <div key={i} className="bg-gray-100 h-4 rounded"></div>
              ))}
            </div>
          </div>
        );
      case 'comparison':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="text-xs font-semibold mb-1">Compare</div>
            <div className="grid grid-cols-2 gap-1">
              <div className="border border-gray-200 p-1">
                <div className="bg-gray-100 h-3 rounded mb-1"></div>
                <div className="text-xs">Product A</div>
              </div>
              <div className="border border-gray-200 p-1">
                <div className="bg-gray-100 h-3 rounded mb-1"></div>
                <div className="text-xs">Product B</div>
              </div>
            </div>
          </div>
        );
      case 'policy':
        return (
          <div className="bg-blue-50 border border-blue-200 rounded-md p-2">
            <div className="text-xs font-semibold text-blue-900 mb-1">30-Day Returns</div>
            <div className="text-xs text-blue-700">Free & easy returns</div>
          </div>
        );
      case 'features':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="grid grid-cols-2 gap-1">
              {[1, 2, 3, 4].map(i => (
                <div key={i} className="text-center">
                  <div className="w-3 h-3 bg-blue-500 rounded-full mx-auto mb-1"></div>
                  <div className="text-xs">Feature {i}</div>
                </div>
              ))}
            </div>
          </div>
        );
      case 'testimonial':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="flex items-start space-x-1">
              <div className="w-3 h-3 bg-gray-300 rounded-full"></div>
              <div>
                <div className="text-xs text-gray-600">"Great product!"</div>
                <div className="text-xs font-semibold">- John D.</div>
              </div>
            </div>
          </div>
        );
      case 'newsletter':
        return (
          <div className="bg-blue-50 border border-blue-200 rounded-md p-2">
            <div className="text-xs font-semibold text-blue-900 mb-1">Subscribe</div>
            <div className="flex space-x-1">
              <div className="bg-white h-2 rounded flex-1"></div>
              <div className="bg-blue-600 text-white text-xs px-1 rounded">Sign up</div>
            </div>
          </div>
        );
      case 'blog':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="space-y-1">
              <div className="bg-gray-100 h-3 rounded"></div>
              <div className="text-xs font-semibold">Latest Blog Posts</div>
              <div className="text-xs text-gray-600">Read our latest articles...</div>
            </div>
          </div>
        );
      case 'footer':
        return (
          <div className="bg-gray-800 text-white rounded-md p-2">
            <div className="grid grid-cols-3 gap-1 mb-1">
              <div className="text-xs">Links</div>
              <div className="text-xs">Contact</div>
              <div className="text-xs">Social</div>
            </div>
            <div className="text-xs opacity-75">© 2024 Store</div>
          </div>
        );
      case 'breadcrumb':
        return (
          <div className="bg-gray-50 border border-gray-200 rounded-md p-2">
            <div className="flex items-center space-x-1 text-xs">
              <span>Home</span>
              <span className="text-gray-400">></span>
              <span>Products</span>
              <span className="text-gray-400">></span>
              <span className="text-blue-600">Item</span>
            </div>
          </div>
        );
      case 'announcement':
      case 'banner':
        return (
          <div className="bg-yellow-100 border border-yellow-300 rounded-md p-2">
            <div className="text-center text-xs font-semibold text-yellow-800">
              🎉 Free shipping on orders over $50!
            </div>
          </div>
        );
      case 'category':
        return (
          <div className="bg-white border border-gray-200 rounded-md p-2">
            <div className="grid grid-cols-2 gap-1">
              {['Electronics', 'Fashion', 'Home', 'Sports'].map(cat => (
                <div key={cat} className="bg-gray-100 text-center py-1 rounded">
                  <div className="text-xs">{cat}</div>
                </div>
              ))}
            </div>
          </div>
        );
      case 'trust':
        return (
          <div className="bg-green-50 border border-green-200 rounded-md p-2">
            <div className="flex justify-center space-x-1">
              <div className="text-xs text-green-800">✓ Secure</div>
              <div className="text-xs text-green-800">✓ Fast</div>
              <div className="text-xs text-green-800">✓ Trusted</div>
            </div>
          </div>
        );
      default:
        return (
          <div className="bg-gray-50 border border-gray-200 rounded-md p-2 flex items-center justify-center">
            <div className="text-xs text-gray-500">{template.type}</div>
          </div>
        );
    }
  };

  return (
    <div className="w-full h-16 flex items-center justify-center bg-gray-50 rounded-lg border border-gray-200">
      {renderPreview()}
    </div>
  );
};

const DraggableComponent = ({ template, viewMode }) => {
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

  const iconMap = {
    header: Layout,
    hero: Zap,
    product: Package,
    'product-grid': Grid,
    'product-card': Package,
    cart: ShoppingCart,
    checkout: CreditCard,
    promotion: Tag,
    reviews: Star,
    testimonials: Users,
    shipping: Package,
    search: Search,
    filter: Filter,
    wishlist: Heart,
    recent: Heart,
    comparison: Grid,
    policy: Type,
    features: Star,
    testimonial: MessageSquare,
    newsletter: Mail,
    blog: Type,
    footer: Layout,
    breadcrumb: Type,
    announcement: Tag,
    banner: Tag,
    category: Grid,
    trust: Star,
    text: Type,
    image: ImageIcon
  };

  const Icon = iconMap[template.type] || Package;

  if (viewMode === 'list') {
    return (
      <div
        ref={setNodeRef}
        style={style}
        {...listeners}
        {...attributes}
        className="flex items-center space-x-3 p-3 bg-white border border-gray-200 rounded-lg hover:border-blue-300 hover:shadow-md transition-all cursor-grab active:cursor-grabbing"
      >
        <div className="w-10 h-10 bg-blue-50 rounded-lg flex items-center justify-center">
          <Icon className="w-5 h-5 text-blue-600" />
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="font-medium text-gray-900 text-sm truncate">{template.name}</h3>
          <p className="text-xs text-gray-500 truncate">{template.description}</p>
        </div>
        <div className="text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded-full">
          {template.type}
        </div>
      </div>
    );
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      className="bg-white border border-gray-200 rounded-lg p-3 hover:border-blue-300 hover:shadow-md transition-all cursor-grab active:cursor-grabbing"
    >
      <div className="flex items-center justify-between mb-2">
        <div className="w-6 h-6 bg-blue-50 rounded flex items-center justify-center">
          <Icon className="w-4 h-4 text-blue-600" />
        </div>
        <div className="text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full">
          {template.type}
        </div>
      </div>

      <ComponentPreview template={template} />

      <div className="mt-2">
        <h3 className="font-medium text-gray-900 text-sm truncate">{template.name}</h3>
        <p className="text-xs text-gray-500 mt-0.5 line-clamp-2">{template.description}</p>
      </div>

      <div className="mt-2 text-center">
        <div className="text-xs text-blue-600 font-medium opacity-0 group-hover:opacity-100 transition-opacity">
          Drag to Canvas
        </div>
      </div>
    </div>
  );
};

const ComponentCategory = ({ title, icon, components, isExpanded, onToggle, viewMode }) => {
  const Icon = icon;
  return (
  <div className="mb-4">
    <button
      onClick={onToggle}
      className="w-full flex items-center justify-between p-3 bg-gray-50 hover:bg-gray-100 rounded-lg transition-colors"
    >
      <div className="flex items-center space-x-3">
        <Icon className="w-4 h-4 text-gray-600" />
        <span className="font-medium text-gray-800">{title}</span>
        <span className="px-2 py-0.5 text-xs bg-blue-100 text-blue-700 rounded-full">
          {components.length}
        </span>
      </div>
      {isExpanded ? (
        <ChevronDown className="w-4 h-4 text-gray-500" />
      ) : (
        <ChevronRight className="w-4 h-4 text-gray-500" />
      )}
    </button>

    {isExpanded && (
      <div className={`mt-3 ${
        viewMode === 'grid'
          ? 'grid grid-cols-1 gap-3'
          : 'space-y-2'
      }`}>
        {components.map(template => (
          <DraggableComponent
            key={template.id}
            template={template}
            viewMode={viewMode}
          />
        ))}
      </div>
    )}
  </div>
  );
};

const ComponentsPanel = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [viewMode, setViewMode] = useState('grid'); // 'grid' or 'list'
  const [expandedCategories, setExpandedCategories] = useState({
    header: true,
    hero: true,
    product: true,
    cart: false,
    promotional: false,
    navigation: false,
    socialProof: false,
    userAccount: false,
    utility: false
  });

  const toggleCategory = (category) => {
    setExpandedCategories(prev => ({
      ...prev,
      [category]: !prev[category]
    }));
  };

  // Filter components based on search
  const filteredComponents = componentTemplates.filter(component =>
    component.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    component.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
    component.type.toLowerCase().includes(searchQuery.toLowerCase())
  );

  // Group components by category for e-commerce
  const componentsByCategory = {
    header: filteredComponents.filter(t => t.category === 'header'),
    hero: filteredComponents.filter(t => t.category === 'hero'),
    product: filteredComponents.filter(t => t.category === 'product' || t.type.includes('product')),
    cart: filteredComponents.filter(t => ['cart', 'checkout'].includes(t.category)),
    promotional: filteredComponents.filter(t => ['promotional', 'testimonial', 'newsletter', 'cta', 'banner'].includes(t.category)),
    navigation: filteredComponents.filter(t => ['navigation', 'breadcrumb'].includes(t.category)),
    socialProof: filteredComponents.filter(t => t.category === 'social-proof'),
    userAccount: filteredComponents.filter(t => t.category === 'user-account'),
    utility: filteredComponents.filter(t => ['utility', 'footer', 'text', 'image', 'video', 'gallery'].includes(t.category))
  };

  const categories = [
    {
      key: 'header',
      title: 'Headers & Navigation',
      icon: Layout,
      description: 'Store headers, navigation menus, and top bars'
    },
    {
      key: 'hero',
      title: 'Hero Sections',
      icon: Zap,
      description: 'Eye-catching banners and promotional sections'
    },
    {
      key: 'product',
      title: 'Product Components',
      icon: Package,
      description: 'Product cards, grids, comparison, and showcase components'
    },
    {
      key: 'cart',
      title: 'Cart & Checkout',
      icon: CreditCard,
      description: 'Shopping cart, checkout flows, and payment components'
    },
    {
      key: 'promotional',
      title: 'Marketing & Promotions',
      icon: Tag,
      description: 'Sales banners, coupons, discounts, and promotional elements'
    },
    {
      key: 'navigation',
      title: 'Search & Filters',
      icon: Filter,
      description: 'Product search, filters, and navigation tools'
    },
    {
      key: 'socialProof',
      title: 'Reviews & Social Proof',
      icon: Users,
      description: 'Customer reviews, testimonials, and trust signals'
    },
    {
      key: 'userAccount',
      title: 'User Account',
      icon: Heart,
      description: 'Wishlist, recently viewed, and account-related components'
    },
    {
      key: 'utility',
      title: 'Utility & Content',
      icon: Type,
      description: 'Shipping, returns, policies, and general content'
    }
  ];

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-lg font-semibold text-gray-900">Components</h2>
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
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            type="text"
            placeholder="Search components..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-4">
        {searchQuery ? (
          // Show search results
          <div className={`${
            viewMode === 'grid'
              ? 'grid grid-cols-1 gap-3'
              : 'space-y-2'
          }`}>
            {filteredComponents.map(template => (
              <DraggableComponent
                key={template.id}
                template={template}
                viewMode={viewMode}
              />
            ))}
          </div>
        ) : (
          // Show categories
          categories.map(category => (
            <ComponentCategory
              key={category.key}
              title={category.title}
              icon={category.icon}
              components={componentsByCategory[category.key] || []}
              isExpanded={expandedCategories[category.key]}
              onToggle={() => toggleCategory(category.key)}
              viewMode={viewMode}
            />
          ))
        )}

        {searchQuery && filteredComponents.length === 0 && (
          <div className="text-center py-8 text-gray-500">
            <Package className="w-12 h-12 mx-auto mb-3 text-gray-300" />
            <p>No components found</p>
            <p className="text-sm">Try adjusting your search terms</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default ComponentsPanel;