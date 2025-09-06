import React, { useState } from 'react';
import { X, ExternalLink, Maximize2, Minimize2 } from 'lucide-react';

const ThemePreviewModal = ({ theme, isOpen, onClose, onApply, onCustomize }) => {
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [iframeKey, setIframeKey] = useState(0);

  if (!isOpen || !theme) return null;

  const previewUrl = `/theme-preview/${theme.id}`;

  const refreshPreview = () => {
    setIframeKey(prev => prev + 1);
  };

  const openInNewTab = () => {
    window.open(previewUrl, '_blank');
  };

  const toggleFullscreen = () => {
    setIsFullscreen(!isFullscreen);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div 
        className={`bg-white rounded-lg shadow-xl transition-all duration-300 ${
          isFullscreen 
            ? 'w-full h-full rounded-none' 
            : 'w-11/12 h-5/6 max-w-6xl'
        }`}
      >
        {/* Modal Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200">
          <div className="flex items-center space-x-4">
            <h2 className="text-xl font-semibold text-gray-900">
              {theme.name} Preview
            </h2>
            <span className="px-2 py-1 text-xs font-medium bg-gray-100 text-gray-600 rounded-full">
              {theme.category}
            </span>
          </div>
          
          <div className="flex items-center space-x-2">
            <button
              onClick={refreshPreview}
              className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
              title="Refresh Preview"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </button>
            
            <button
              onClick={openInNewTab}
              className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
              title="Open in New Tab"
            >
              <ExternalLink className="w-4 h-4" />
            </button>
            
            <button
              onClick={toggleFullscreen}
              className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
              title={isFullscreen ? 'Exit Fullscreen' : 'Fullscreen'}
            >
              {isFullscreen ? (
                <Minimize2 className="w-4 h-4" />
              ) : (
                <Maximize2 className="w-4 h-4" />
              )}
            </button>
            
            <button
              onClick={onClose}
              className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
              title="Close Preview"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        </div>

        {/* Preview Content */}
        <div className="flex-1 bg-gray-50" style={{ height: 'calc(100% - 73px)' }}>
          <iframe
            key={iframeKey}
            src={previewUrl}
            className="w-full h-full border-0"
            title={`${theme.name} Preview`}
            sandbox="allow-same-origin allow-scripts allow-forms"
          />
        </div>

        {/* Theme Info Footer */}
        <div className="px-4 py-3 bg-gray-50 border-t border-gray-200 rounded-b-lg">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="text-sm text-gray-600">
                <span className="font-medium">Colors:</span>
                <div className="inline-flex ml-2 space-x-1">
                  <div 
                    className="w-4 h-4 rounded-full border border-gray-300"
                    style={{ backgroundColor: theme.colors.primary }}
                    title="Primary"
                  />
                  <div 
                    className="w-4 h-4 rounded-full border border-gray-300"
                    style={{ backgroundColor: theme.colors.secondary }}
                    title="Secondary"
                  />
                  <div 
                    className="w-4 h-4 rounded-full border border-gray-300"
                    style={{ backgroundColor: theme.colors.accent }}
                    title="Accent"
                  />
                </div>
              </div>
              <div className="text-sm text-gray-600">
                <span className="font-medium">Font:</span>
                <span className="ml-1">{theme.typography.fontFamily}</span>
              </div>
            </div>
            
            <div className="flex items-center space-x-2">
              <button
                onClick={() => {
                  if (onApply) {
                    onApply(theme);
                    onClose();
                  }
                }}
                className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors"
              >
                Apply Theme
              </button>
              <button
                onClick={() => {
                  if (onCustomize) {
                    onCustomize(theme);
                    onClose();
                  }
                }}
                className="px-4 py-2 bg-gray-600 text-white text-sm font-medium rounded-lg hover:bg-gray-700 transition-colors"
              >
                Customize
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ThemePreviewModal;