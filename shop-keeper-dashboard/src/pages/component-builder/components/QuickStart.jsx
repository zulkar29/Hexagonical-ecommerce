import React, { useState } from 'react';
import {
  Zap,
  ShoppingBag,
  Palette,
  Rocket,
  ArrowRight,
  Briefcase,
  Crown,
  Sparkles
} from 'lucide-react';

const QuickStart = ({ onSelectTemplate, onSelectTheme, selectedTheme }) => {
  const [activeTab, setActiveTab] = useState('templates');

  // Theme options with working quick start templates
  const themeTemplates = [
    {
      id: 'modern-clean',
      name: 'Modern Clean',
      description: 'Professional and minimalist design perfect for any business',
      icon: <Zap className="w-5 h-5" />,
      colors: { primary: '#2563eb', secondary: '#64748b', accent: '#0ea5e9', background: '#ffffff', text: '#1e293b' },
      recommended: true,
      preview: '🎨 Blue accents, clean typography, professional look'
    },
    {
      id: 'classic-commerce',
      name: 'Classic Commerce',
      description: 'Traditional e-commerce design with proven conversion patterns',
      icon: <ShoppingBag className="w-5 h-5" />,
      colors: { primary: '#dc2626', secondary: '#374151', accent: '#f59e0b', background: '#ffffff', text: '#111827' },
      recommended: true,
      preview: '🛒 Red call-to-actions, warm colors, trusted feel'
    },
    {
      id: 'luxury-brand',
      name: 'Luxury Brand',
      description: 'Elegant and premium design for high-end products',
      icon: <Crown className="w-5 h-5" />,
      colors: { primary: '#000000', secondary: '#6b7280', accent: '#d4af37', background: '#ffffff', text: '#000000' },
      preview: '✨ Black & gold, sophisticated, premium feel'
    },
    {
      id: 'creative-bold',
      name: 'Creative Bold',
      description: 'Vibrant and eye-catching for creative businesses',
      icon: <Palette className="w-5 h-5" />,
      colors: { primary: '#7c3aed', secondary: '#ec4899', accent: '#06b6d4', background: '#fefefe', text: '#1f2937' },
      preview: '🌈 Purple & pink, creative, energetic vibe'
    },
    {
      id: 'professional-corporate',
      name: 'Professional Corporate',
      description: 'Sophisticated business design for B2B companies',
      icon: <Briefcase className="w-5 h-5" />,
      colors: { primary: '#1e40af', secondary: '#475569', accent: '#059669', background: '#ffffff', text: '#0f172a' },
      preview: '💼 Navy blue, corporate, trustworthy appearance'
    }
  ];

  const quickStartTemplates = [
    {
      id: 'minimal-store',
      name: 'Minimal Store',
      description: 'Clean, professional layout perfect for any product',
      icon: <Zap className="w-5 h-5" />,
      preview: '🏪 Header + Hero + Products + Footer',
      recommended: true,
      components: [
        { type: 'header', title: 'Your Store' },
        { type: 'hero', title: 'Welcome to Our Store', subtitle: 'Discover amazing products' },
        { type: 'product', title: 'Featured Products' },
        { type: 'footer', companyName: 'Your Store' }
      ]
    },
    {
      id: 'product-showcase',
      name: 'Product Showcase',
      description: 'Highlight your best products with style',
      icon: <ShoppingBag className="w-5 h-5" />,
      preview: '🎯 Hero + Product Grid + Category Banner',
      components: [
        { type: 'hero', title: 'Premium Products', subtitle: 'Quality you can trust' },
        { type: 'product', title: 'Best Sellers' },
        { type: 'category', title: 'New Collection' }
      ]
    },
    {
      id: 'brand-experience',
      name: 'Brand Experience',
      description: 'Tell your brand story while selling products',
      icon: <Palette className="w-5 h-5" />,
      preview: '🎨 Header + Brand Story + Products + About',
      components: [
        { type: 'header', title: 'Brand Name' },
        { type: 'hero', title: 'Our Story', subtitle: 'Crafted with passion since 2020' },
        { type: 'product', title: 'Our Products' },
        { type: 'footer', companyName: 'Brand Name' }
      ]
    }
  ];

  const handleTemplateSelect = (template) => {
    onSelectTemplate(template);
  };

  const handleThemeSelect = (theme) => {
    if (onSelectTheme) {
      onSelectTheme(theme);
    }
  };

  return (
    <div className="p-6 space-y-6">
      <div className="text-center">
        <Rocket className="w-12 h-12 text-blue-600 mx-auto mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 mb-2">Quick Start Store Builder</h2>
        <p className="text-gray-600 mb-6">
          Choose a theme and template to get your store design up and running in seconds
        </p>
      </div>

      {/* Tab Navigation */}
      <div className="flex space-x-1 bg-gray-100 p-1 rounded-lg">
        <button
          onClick={() => setActiveTab('templates')}
          className={`flex-1 px-4 py-2 text-sm font-medium rounded-md transition-colors ${
            activeTab === 'templates'
              ? 'bg-white text-blue-600 shadow-sm'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          <Rocket className="w-4 h-4 inline mr-2" />
          Store Templates
        </button>
        <button
          onClick={() => setActiveTab('themes')}
          className={`flex-1 px-4 py-2 text-sm font-medium rounded-md transition-colors ${
            activeTab === 'themes'
              ? 'bg-white text-blue-600 shadow-sm'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          <Sparkles className="w-4 h-4 inline mr-2" />
          Color Themes
        </button>
      </div>

      {/* Content Area */}
      {activeTab === 'templates' ? (
        <div className="space-y-4">
          {quickStartTemplates.map((template) => (
            <div
              key={template.id}
              className="border border-gray-200 rounded-lg p-4 hover:border-blue-300 hover:shadow-md transition-all cursor-pointer"
              onClick={() => handleTemplateSelect(template)}
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-blue-100 rounded-lg text-blue-600">
                    {template.icon}
                  </div>
                  <div>
                    <div className="flex items-center space-x-2">
                      <h3 className="font-medium text-gray-900">{template.name}</h3>
                      {template.recommended && (
                        <span className="px-2 py-1 text-xs bg-green-100 text-green-700 rounded-full">
                          ⭐ Recommended
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-gray-600 mt-1">{template.description}</p>
                  </div>
                </div>
                <ArrowRight className="w-5 h-5 text-gray-400" />
              </div>

              <div className="bg-gray-50 rounded-lg p-3 mb-3">
                <div className="text-xs text-gray-600 mb-1">Template Preview:</div>
                <div className="text-sm font-mono text-gray-800">{template.preview}</div>
              </div>

              <div className="flex items-center justify-between">
                <div className="text-xs text-gray-500">
                  {template.components.length} components included
                </div>
                <button className="px-3 py-1.5 bg-blue-600 text-white text-sm rounded-md hover:bg-blue-700 transition-colors">
                  Use Template
                </button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="space-y-4">
          {themeTemplates.map((theme) => (
            <div
              key={theme.id}
              className="border border-gray-200 rounded-lg p-4 hover:border-blue-300 hover:shadow-md transition-all cursor-pointer"
              onClick={() => handleThemeSelect(theme)}
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-gray-100 rounded-lg text-gray-600">
                    {theme.icon}
                  </div>
                  <div>
                    <div className="flex items-center space-x-2">
                      <h3 className="font-medium text-gray-900">{theme.name}</h3>
                      {theme.recommended && (
                        <span className="px-2 py-1 text-xs bg-green-100 text-green-700 rounded-full">
                          ⭐ Recommended
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-gray-600 mt-1">{theme.description}</p>
                  </div>
                </div>
                <ArrowRight className="w-5 h-5 text-gray-400" />
              </div>

              <div className="bg-gray-50 rounded-lg p-3 mb-3">
                <div className="text-xs text-gray-600 mb-1">Color Preview:</div>
                <div className="flex space-x-2 items-center">
                  <div
                    className="w-4 h-4 rounded-full border"
                    style={{ backgroundColor: theme.colors.primary }}
                    title="Primary"
                  />
                  <div
                    className="w-4 h-4 rounded-full border"
                    style={{ backgroundColor: theme.colors.secondary }}
                    title="Secondary"
                  />
                  <div
                    className="w-4 h-4 rounded-full border"
                    style={{ backgroundColor: theme.colors.accent }}
                    title="Accent"
                  />
                  <span className="text-sm text-gray-600 ml-2">{theme.preview}</span>
                </div>
              </div>

              <div className="flex items-center justify-between">
                <div className="text-xs text-gray-500">
                  Apply to existing components
                </div>
                <button className="px-3 py-1.5 bg-purple-600 text-white text-sm rounded-md hover:bg-purple-700 transition-colors">
                  Apply Theme
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {selectedTheme && (
        <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
          <div className="flex items-center space-x-2 mb-2">
            <Palette className="w-4 h-4 text-purple-600" />
            <span className="text-sm font-medium text-purple-900">
              Active Theme: {selectedTheme.name}
            </span>
          </div>
          <p className="text-xs text-purple-700">
            Quick start templates will automatically use your selected theme colors and styling.
          </p>
        </div>
      )}

      <div className="text-center pt-4 border-t border-gray-200">
        <p className="text-xs text-gray-500">
          💡 All templates are fully customizable after creation
        </p>
      </div>
    </div>
  );
};

export default QuickStart;