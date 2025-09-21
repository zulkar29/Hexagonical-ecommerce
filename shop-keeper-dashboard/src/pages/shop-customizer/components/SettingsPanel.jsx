import React, { useState } from 'react';
import {
  Save,
  Download,
  Upload,
  Trash2,
  Undo,
  Redo,
  Settings,
  Globe,
  Code,
  Monitor,
  Smartphone,
  Tablet,
  History,
  FileText,
  Eye,
  ExternalLink
} from 'lucide-react';

const SettingsPanel = ({
  designName,
  setDesignName,
  customizerSettings,
  setCustomizerSettings,
  onSave,
  onUndo,
  onRedo,
  history,
  components,
  selectedTheme
}) => {
  const [activeTab, setActiveTab] = useState('general');

  // Export functions for database storage
  const generateStoreConfigJSON = () => {
    const storeConfig = {
      version: '1.0.0',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      designName: designName || 'Untitled Store',
      theme: selectedTheme ? {
        id: selectedTheme.id,
        name: selectedTheme.name,
        colors: selectedTheme.colors,
        typography: selectedTheme.typography,
        spacing: selectedTheme.spacing,
        borderRadius: selectedTheme.borderRadius,
        components: selectedTheme.components
      } : null,
      layout: {
        devicePreview: customizerSettings.devicePreview,
        autoSave: customizerSettings.autoSave,
        realTimePreview: customizerSettings.realTimePreview
      },
      components: components.map((component, index) => ({
        id: component.id,
        templateId: component.templateId,
        type: component.type,
        name: component.name,
        category: component.category || 'content',
        order: index,
        props: component.props || {},
        styles: component.styles || {},
        theme: component.theme || null,
        visible: component.visible !== false,
        responsive: component.responsive || {
          desktop: true,
          tablet: true,
          mobile: true
        }
      })),
      metadata: {
        totalComponents: components.length,
        categories: [...new Set(components.map(c => c.category || 'content'))],
        hasTheme: !!selectedTheme,
        lastModified: new Date().toISOString()
      }
    };
    return storeConfig;
  };

  const exportAsJSON = () => {
    try {
      const config = generateStoreConfigJSON();
      const jsonString = JSON.stringify(config, null, 2);
      const blob = new Blob([jsonString], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${designName || 'store-design'}-${Date.now()}.json`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Error exporting JSON:', error);
      alert('Failed to export design. Please try again.');
    }
  };

  const saveToDatabase = async () => {
    try {
      const config = generateStoreConfigJSON();
      // This would typically call your API to save to database
      console.log('Store configuration for database:', config);

      // Simulate API call
      const response = await fetch('/api/store-designs', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(config)
      });

      if (response.ok) {
        alert('Store design saved successfully!');
      } else {
        throw new Error('Failed to save to database');
      }
    } catch (error) {
      console.error('Error saving to database:', error);
      alert('Failed to save design to database. Please try again.');
    }
  };

  const generatePreviewLink = () => {
    const config = generateStoreConfigJSON();
    const encodedConfig = encodeURIComponent(JSON.stringify(config));
    const baseUrl = window.location.origin;
    const previewUrl = `${baseUrl}/preview?config=${encodedConfig}`;

    navigator.clipboard.writeText(previewUrl).then(() => {
      alert('Preview link copied to clipboard!');
    }).catch(() => {
      prompt('Copy this preview link:', previewUrl);
    });
  };

  const tabs = [
    { id: 'general', name: 'General', icon: Settings },
    { id: 'export', name: 'Export', icon: Download },
    { id: 'history', name: 'History', icon: History }
  ];

  const renderGeneralSettings = () => (
    <div className="space-y-6">
      {/* Design Info */}
      <div>
        <h3 className="text-sm font-semibold text-gray-900 mb-3">Design Information</h3>
        <div className="space-y-3">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Design Name
            </label>
            <input
              type="text"
              value={designName}
              onChange={(e) => setDesignName(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="My Store Design"
            />
          </div>

          <div className="grid grid-cols-2 gap-3 text-sm">
            <div className="bg-gray-50 rounded-lg p-3">
              <div className="text-gray-600">Components</div>
              <div className="font-semibold">{components.length}</div>
            </div>
            <div className="bg-gray-50 rounded-lg p-3">
              <div className="text-gray-600">Last Saved</div>
              <div className="font-semibold">Never</div>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div>
        <h3 className="text-sm font-semibold text-gray-900 mb-3">Quick Actions</h3>
        <div className="grid grid-cols-2 gap-2">
          <button
            onClick={onUndo}
            disabled={history.past.length === 0}
            className="flex items-center justify-center space-x-2 px-3 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Undo className="w-4 h-4" />
            <span>Undo</span>
          </button>
          <button
            onClick={onRedo}
            disabled={history.future.length === 0}
            className="flex items-center justify-center space-x-2 px-3 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Redo className="w-4 h-4" />
            <span>Redo</span>
          </button>
        </div>
      </div>

      {/* Preview Settings */}
      <div>
        <h3 className="text-sm font-semibold text-gray-900 mb-3">Preview Settings</h3>
        <div className="space-y-3">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Device Preview
            </label>
            <div className="grid grid-cols-3 gap-2">
              {[
                { key: 'desktop', icon: Monitor, label: 'Desktop' },
                { key: 'tablet', icon: Tablet, label: 'Tablet' },
                { key: 'mobile', icon: Smartphone, label: 'Mobile' }
              ].map((device) => {
                const Icon = device.icon;
                const isActive = customizerSettings.devicePreview === device.key;
                return (
                  <button
                    key={device.key}
                    onClick={() => setCustomizerSettings({
                      ...customizerSettings,
                      devicePreview: device.key
                    })}
                    className={`flex flex-col items-center space-y-1 p-3 rounded-lg transition-colors ${
                      isActive
                        ? 'bg-blue-100 text-blue-700 border border-blue-200'
                        : 'bg-gray-50 text-gray-600 hover:bg-gray-100'
                    }`}
                  >
                    <Icon className="w-5 h-5" />
                    <span className="text-xs font-medium">{device.label}</span>
                  </button>
                );
              })}
            </div>
          </div>
        </div>
      </div>

      {/* Performance */}
      <div>
        <h3 className="text-sm font-semibold text-gray-900 mb-3">Performance</h3>
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-700">Auto-save</span>
            <button
              onClick={() => setCustomizerSettings({
                ...customizerSettings,
                autoSave: !customizerSettings.autoSave
              })}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                customizerSettings.autoSave ? 'bg-blue-600' : 'bg-gray-200'
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  customizerSettings.autoSave ? 'translate-x-6' : 'translate-x-1'
                }`}
              />
            </button>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-700">Real-time preview</span>
            <button
              onClick={() => setCustomizerSettings({
                ...customizerSettings,
                realTimePreview: !customizerSettings.realTimePreview
              })}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                customizerSettings.realTimePreview ? 'bg-blue-600' : 'bg-gray-200'
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  customizerSettings.realTimePreview ? 'translate-x-6' : 'translate-x-1'
                }`}
              />
            </button>
          </div>
        </div>
      </div>
    </div>
  );

  const renderExportSettings = () => (
    <div className="space-y-6">
      {/* Save Design */}
      <div>
        <h3 className="text-sm font-semibold text-gray-900 mb-3">Save Design</h3>
        <div className="space-y-3">
          <button
            onClick={onSave}
            disabled={!designName.trim()}
            className="w-full flex items-center justify-center space-x-2 px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Save className="w-4 h-4" />
            <span>Save Design</span>
          </button>
          <p className="text-xs text-gray-500">
            Save your current design to continue working on it later.
          </p>
        </div>
      </div>

      {/* Export Options */}
      <div>
        <h3 className="text-sm font-semibold text-gray-900 mb-3">Export Options</h3>
        <div className="space-y-2">
          <button
            onClick={saveToDatabase}
            className="w-full flex items-center justify-between px-4 py-3 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors text-left"
          >
            <div className="flex items-center space-x-3">
              <Code className="w-4 h-4 text-blue-600" />
              <div>
                <div className="font-medium text-blue-900">Save to Database</div>
                <div className="text-xs text-blue-700">Store config for website generation</div>
              </div>
            </div>
            <ExternalLink className="w-4 h-4 text-blue-400" />
          </button>

          <button
            onClick={exportAsJSON}
            className="w-full flex items-center justify-between px-4 py-3 bg-gray-50 hover:bg-gray-100 rounded-lg transition-colors text-left"
          >
            <div className="flex items-center space-x-3">
              <FileText className="w-4 h-4 text-gray-600" />
              <div>
                <div className="font-medium text-gray-900">Export as JSON</div>
                <div className="text-xs text-gray-500">Download design configuration</div>
              </div>
            </div>
            <Download className="w-4 h-4 text-gray-400" />
          </button>

          <button
            onClick={generatePreviewLink}
            className="w-full flex items-center justify-between px-4 py-3 bg-gray-50 hover:bg-gray-100 rounded-lg transition-colors text-left"
          >
            <div className="flex items-center space-x-3">
              <Eye className="w-4 h-4 text-gray-600" />
              <div>
                <div className="font-medium text-gray-900">Preview Link</div>
                <div className="text-xs text-gray-500">Share preview with others</div>
              </div>
            </div>
            <ExternalLink className="w-4 h-4 text-gray-400" />
          </button>
        </div>
      </div>

      {/* Publish */}
      <div>
        <h3 className="text-sm font-semibold text-gray-900 mb-3">Publish Store</h3>
        <div className="space-y-3">
          <button
            disabled={components.length === 0}
            className="w-full flex items-center justify-center space-x-2 px-4 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Globe className="w-4 h-4" />
            <span>Publish Live</span>
          </button>
          <p className="text-xs text-gray-500">
            Make your store live for customers to see.
          </p>
        </div>
      </div>

      {/* Import */}
      <div>
        <h3 className="text-sm font-semibold text-gray-900 mb-3">Import Design</h3>
        <div className="space-y-2">
          <button className="w-full flex items-center justify-center space-x-2 px-4 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors">
            <Upload className="w-4 h-4" />
            <span>Import from File</span>
          </button>
        </div>
      </div>
    </div>
  );

  const renderHistorySettings = () => (
    <div className="space-y-6">
      <div>
        <h3 className="text-sm font-semibold text-gray-900 mb-3">Design History</h3>
        <div className="space-y-2">
          {history.past.length === 0 && history.future.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <History className="w-8 h-8 mx-auto mb-2 text-gray-300" />
              <p className="text-sm">No history available</p>
              <p className="text-xs">Actions will appear here as you work</p>
            </div>
          ) : (
            <>
              {/* Current state */}
              <div className="flex items-center space-x-3 p-3 bg-blue-50 border border-blue-200 rounded-lg">
                <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                <div className="flex-1">
                  <div className="font-medium text-blue-900">Current State</div>
                  <div className="text-xs text-blue-700">{components.length} components</div>
                </div>
              </div>

              {/* Future states */}
              {history.future.slice(0, 5).map((state, index) => (
                <div key={`future-${index}`} className="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg opacity-60">
                  <div className="w-2 h-2 bg-gray-400 rounded-full"></div>
                  <div className="flex-1">
                    <div className="font-medium text-gray-700">Future State {index + 1}</div>
                    <div className="text-xs text-gray-500">{state.components?.length || 0} components</div>
                  </div>
                </div>
              ))}

              {/* Past states */}
              {history.past.slice(-5).reverse().map((state, index) => (
                <div key={`past-${index}`} className="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg">
                  <div className="w-2 h-2 bg-gray-400 rounded-full"></div>
                  <div className="flex-1">
                    <div className="font-medium text-gray-700">Previous State {index + 1}</div>
                    <div className="text-xs text-gray-500">{state.components?.length || 0} components</div>
                  </div>
                </div>
              ))}
            </>
          )}
        </div>
      </div>

      {/* Clear History */}
      <div>
        <h3 className="text-sm font-semibold text-gray-900 mb-3">Manage History</h3>
        <button
          disabled={history.past.length === 0 && history.future.length === 0}
          className="w-full flex items-center justify-center space-x-2 px-4 py-3 bg-red-50 text-red-700 rounded-lg hover:bg-red-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <Trash2 className="w-4 h-4" />
          <span>Clear History</span>
        </button>
      </div>
    </div>
  );

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-gray-200">
        <h2 className="text-lg font-semibold text-gray-900">Settings</h2>

        {/* Tabs */}
        <div className="mt-3 flex space-x-1 bg-gray-100 rounded-lg p-1">
          {tabs.map(tab => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center space-x-2 px-3 py-2 rounded-md text-sm font-medium transition-colors flex-1 justify-center ${
                  activeTab === tab.id
                    ? 'bg-white text-blue-600 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                <Icon className="w-4 h-4" />
                <span>{tab.name}</span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-4">
        {activeTab === 'general' && renderGeneralSettings()}
        {activeTab === 'export' && renderExportSettings()}
        {activeTab === 'history' && renderHistorySettings()}
      </div>
    </div>
  );
};

export default SettingsPanel;