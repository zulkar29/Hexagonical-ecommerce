import React, { useState } from 'react';
import {
  Store,
  Package,
  Truck,
  CreditCard,
  Shield,
  Zap,
  Globe,
  Settings,
  CheckCircle,
  ExternalLink
} from 'lucide-react';

const VendorFeatures = ({ onFeatureToggle }) => {
  const [enabledFeatures, setEnabledFeatures] = useState({
    inventory: true,
    shipping: true,
    payments: true,
    security: true,
    seo: false,
    analytics: true,
    multilanguage: false,
    customDomain: false
  });

  const features = [
    {
      id: 'inventory',
      name: 'Inventory Management',
      description: 'Track stock levels, low stock alerts, and automated reordering',
      icon: Package,
      category: 'core',
      premium: false
    },
    {
      id: 'shipping',
      name: 'Shipping Integration',
      description: 'Connect with shipping providers for real-time rates and tracking',
      icon: Truck,
      category: 'core',
      premium: false
    },
    {
      id: 'payments',
      name: 'Payment Gateway',
      description: 'Accept credit cards, PayPal, and digital wallets securely',
      icon: CreditCard,
      category: 'core',
      premium: false
    },
    {
      id: 'security',
      name: 'SSL & Security',
      description: 'Enterprise-grade security with fraud protection',
      icon: Shield,
      category: 'core',
      premium: false
    },
    {
      id: 'seo',
      name: 'SEO Optimization',
      description: 'Built-in SEO tools, meta tags, and search engine optimization',
      icon: Globe,
      category: 'growth',
      premium: true
    },
    {
      id: 'analytics',
      name: 'Advanced Analytics',
      description: 'Detailed sales reports, customer insights, and conversion tracking',
      icon: Zap,
      category: 'growth',
      premium: false
    },
    {
      id: 'multilanguage',
      name: 'Multi-language Support',
      description: 'Serve customers in multiple languages and currencies',
      icon: Globe,
      category: 'premium',
      premium: true
    },
    {
      id: 'customDomain',
      name: 'Custom Domain',
      description: 'Use your own domain name for professional branding',
      icon: ExternalLink,
      category: 'premium',
      premium: true
    }
  ];

  const handleFeatureToggle = (featureId) => {
    const newState = {
      ...enabledFeatures,
      [featureId]: !enabledFeatures[featureId]
    };
    setEnabledFeatures(newState);
    onFeatureToggle?.(featureId, newState[featureId]);
  };

  const FeatureCard = ({ feature }) => {
    const Icon = feature.icon;
    const isEnabled = enabledFeatures[feature.id];

    return (
      <div className={`p-4 border rounded-lg transition-all ${
        isEnabled
          ? 'border-blue-200 bg-blue-50'
          : 'border-gray-200 bg-white hover:border-gray-300'
      }`}>
        <div className="flex items-start justify-between mb-3">
          <div className="flex items-center space-x-3">
            <div className={`p-2 rounded-lg ${
              isEnabled ? 'bg-blue-100 text-blue-600' : 'bg-gray-100 text-gray-600'
            }`}>
              <Icon className="w-5 h-5" />
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <h4 className="font-medium text-gray-900">{feature.name}</h4>
                {feature.premium && (
                  <span className="px-2 py-1 text-xs bg-gradient-to-r from-purple-100 to-pink-100 text-purple-700 rounded-full">
                    Premium
                  </span>
                )}
              </div>
              <p className="text-sm text-gray-600 mt-1">{feature.description}</p>
            </div>
          </div>

          <button
            onClick={() => handleFeatureToggle(feature.id)}
            disabled={feature.premium && !isEnabled}
            className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
              isEnabled
                ? 'bg-blue-600'
                : feature.premium && !isEnabled
                ? 'bg-gray-200 cursor-not-allowed'
                : 'bg-gray-200 hover:bg-gray-300'
            }`}
          >
            <span
              className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                isEnabled ? 'translate-x-6' : 'translate-x-1'
              }`}
            />
          </button>
        </div>

        {isEnabled && (
          <div className="flex items-center text-sm text-blue-600">
            <CheckCircle className="w-4 h-4 mr-1" />
            Active in your store
          </div>
        )}
      </div>
    );
  };

  const categoryGroups = {
    core: features.filter(f => f.category === 'core'),
    growth: features.filter(f => f.category === 'growth'),
    premium: features.filter(f => f.category === 'premium')
  };

  return (
    <div className="p-6 space-y-6">
      <div className="text-center">
        <Store className="w-12 h-12 text-blue-600 mx-auto mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 mb-2">Store Features</h2>
        <p className="text-gray-600">
          Configure your store's capabilities and integrations
        </p>
      </div>

      {/* Core Features */}
      <div>
        <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
          <div className="w-2 h-2 bg-green-500 rounded-full mr-2"></div>
          Core Features
        </h3>
        <div className="space-y-4">
          {categoryGroups.core.map(feature => (
            <FeatureCard key={feature.id} feature={feature} />
          ))}
        </div>
      </div>

      {/* Growth Features */}
      <div>
        <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
          <div className="w-2 h-2 bg-blue-500 rounded-full mr-2"></div>
          Growth Features
        </h3>
        <div className="space-y-4">
          {categoryGroups.growth.map(feature => (
            <FeatureCard key={feature.id} feature={feature} />
          ))}
        </div>
      </div>

      {/* Premium Features */}
      <div>
        <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
          <div className="w-2 h-2 bg-purple-500 rounded-full mr-2"></div>
          Premium Features
        </h3>
        <div className="space-y-4">
          {categoryGroups.premium.map(feature => (
            <FeatureCard key={feature.id} feature={feature} />
          ))}
        </div>

        <div className="mt-4 p-4 bg-gradient-to-r from-purple-50 to-pink-50 rounded-lg border border-purple-200">
          <div className="flex items-center justify-between">
            <div>
              <h4 className="font-medium text-purple-900">Upgrade to Premium</h4>
              <p className="text-sm text-purple-700 mt-1">
                Unlock advanced features to grow your business
              </p>
            </div>
            <button className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors">
              Upgrade
            </button>
          </div>
        </div>
      </div>

      {/* Feature Summary */}
      <div className="bg-gray-50 p-4 rounded-lg">
        <h4 className="font-medium text-gray-900 mb-2">Active Features Summary</h4>
        <div className="text-sm text-gray-600">
          {Object.values(enabledFeatures).filter(Boolean).length} of {features.length} features enabled
        </div>
        <div className="mt-2 text-xs text-gray-500">
          Changes will be applied to your live store design
        </div>
      </div>
    </div>
  );
};

export default VendorFeatures;