import React from 'react';
import {
  Layers,
  Palette,
  Settings,
  Sliders,
  Eye,
  Smartphone,
  Tablet,
  Monitor,
  Undo,
  Redo,
  Save,
  Menu,
  X,
  Rocket
} from 'lucide-react';

const ToolNavigation = ({
  activePanel,
  setActivePanel,
  panelCollapsed,
  setPanelCollapsed,
  previewMode,
  setPreviewMode,
  selectedTheme,
  components,
  designName,
  customizerSettings,
  setCustomizerSettings
}) => {
  const tools = [
    {
      id: 'components',
      name: 'Components',
      icon: Layers,
      description: 'Add elements to your store'
    },
    {
      id: 'themes',
      name: 'Themes',
      icon: Palette,
      description: 'Choose a design theme'
    },
    {
      id: 'properties',
      name: 'Properties',
      icon: Sliders,
      description: 'Edit selected element'
    },
    {
      id: 'settings',
      name: 'Settings',
      icon: Settings,
      description: 'Store settings & export'
    }
  ];

  const deviceOptions = [
    { key: 'desktop', icon: Monitor, label: 'Desktop' },
    { key: 'tablet', icon: Tablet, label: 'Tablet' },
    { key: 'mobile', icon: Smartphone, label: 'Mobile' }
  ];

  return (
    <div className="bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
      {/* Left Section - Logo & Tools */}
      <div className="flex items-center space-x-4">
        {/* Logo */}
        <div className="flex items-center space-x-2">
          <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
            <Rocket className="w-4 h-4 text-white" />
          </div>
          <span className="font-semibold text-gray-900 hidden sm:inline">Shop Designer</span>
        </div>

        <div className="h-6 w-px bg-gray-300" />

        {/* Panel Toggle */}
        <button
          onClick={() => setPanelCollapsed(!panelCollapsed)}
          className={`p-2 rounded-lg transition-colors ${
            panelCollapsed
              ? 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              : 'bg-blue-100 text-blue-600 hover:bg-blue-200'
          }`}
          title={panelCollapsed ? 'Show Panel' : 'Hide Panel'}
        >
          {panelCollapsed ? <Menu className="w-4 h-4" /> : <X className="w-4 h-4" />}
        </button>

        {/* Tool Buttons */}
        <div className="flex items-center space-x-1">
          {tools.map((tool) => {
            const Icon = tool.icon;
            const isActive = activePanel === tool.id;
            return (
              <button
                key={tool.id}
                onClick={() => {
                  setActivePanel(tool.id);
                  if (panelCollapsed) setPanelCollapsed(false);
                }}
                className={`flex items-center space-x-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isActive
                    ? 'bg-blue-100 text-blue-700 border border-blue-200'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                }`}
                title={tool.description}
              >
                <Icon className="w-4 h-4" />
                <span className="hidden lg:inline">{tool.name}</span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Center Section - Status */}
      <div className="hidden md:flex items-center space-x-4 text-sm text-gray-600">
        <div className="flex items-center space-x-2">
          <div className="w-2 h-2 bg-green-500 rounded-full"></div>
          <span>{components.length} components</span>
        </div>
        {selectedTheme && (
          <div className="flex items-center space-x-2">
            <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
            <span>{selectedTheme.name}</span>
          </div>
        )}
        {designName && (
          <div className="flex items-center space-x-2">
            <div className="w-2 h-2 bg-purple-500 rounded-full"></div>
            <span>{designName}</span>
          </div>
        )}
      </div>

      {/* Right Section - Actions */}
      <div className="flex items-center space-x-3">
        {/* Device Preview */}
        <div className="flex items-center bg-gray-100 rounded-lg p-1">
          {deviceOptions.map((device) => {
            const Icon = device.icon;
            const isActive = customizerSettings.devicePreview === device.key;
            return (
              <button
                key={device.key}
                onClick={() => setCustomizerSettings({
                  ...customizerSettings,
                  devicePreview: device.key
                })}
                className={`p-2 rounded-md transition-colors ${
                  isActive
                    ? 'bg-white text-gray-900 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900'
                }`}
                title={device.label}
              >
                <Icon className="w-4 h-4" />
              </button>
            );
          })}
        </div>

        {/* Undo/Redo - Only on larger screens */}
        <div className="hidden sm:flex items-center bg-gray-100 rounded-lg p-1">
          <button
            className="p-2 text-gray-600 hover:text-gray-900 hover:bg-white rounded-md transition-colors"
            title="Undo"
          >
            <Undo className="w-4 h-4" />
          </button>
          <button
            className="p-2 text-gray-600 hover:text-gray-900 hover:bg-white rounded-md transition-colors"
            title="Redo"
          >
            <Redo className="w-4 h-4" />
          </button>
        </div>

        {/* Preview & Save */}
        <div className="flex items-center space-x-2">
          <button
            onClick={() => setPreviewMode(!previewMode)}
            className={`flex items-center space-x-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
              previewMode
                ? 'bg-blue-600 text-white hover:bg-blue-700'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            <Eye className="w-4 h-4" />
            <span className="hidden sm:inline">{previewMode ? 'Edit' : 'Preview'}</span>
          </button>

          <button
            className="flex items-center space-x-2 px-3 py-2 bg-green-600 text-white rounded-lg text-sm font-medium hover:bg-green-700 transition-colors"
            title="Save Design"
          >
            <Save className="w-4 h-4" />
            <span className="hidden sm:inline">Save</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default ToolNavigation;