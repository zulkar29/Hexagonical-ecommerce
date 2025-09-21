import React, { useState } from 'react';
import { useAtom, useSetAtom } from 'jotai';
import { getSavedDesignsAtom, loadDesignAtom, designNameAtom } from '../store/builderAtoms';
import {
  Save,
  Download,
  Trash2,
  Clock,
  Eye,
  Copy,
  X
} from 'lucide-react';

const SavedDesigns = ({ isOpen, onClose, onLoad }) => {
  const [savedDesigns] = useAtom(getSavedDesignsAtom);
  const loadDesign = useSetAtom(loadDesignAtom);
  const [, setDesignName] = useAtom(designNameAtom);
  const [searchTerm, setSearchTerm] = useState('');

  const filteredDesigns = savedDesigns.filter(design =>
    design.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    design.description?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleLoadDesign = (design) => {
    if (window.confirm('Loading this design will replace your current work. Continue?')) {
      loadDesign(design);
      setDesignName(design.name);
      onLoad?.(design);
      onClose();
    }
  };

  const handleDeleteDesign = (designToDelete) => {
    if (window.confirm(`Are you sure you want to delete "${designToDelete.name}"?`)) {
      const designs = JSON.parse(localStorage.getItem('store-designs') || '[]');
      const filteredDesigns = designs.filter(d => d.id !== designToDelete.id);
      localStorage.setItem('store-designs', JSON.stringify(filteredDesigns));
      // Force re-render by updating the atom
      window.location.reload();
    }
  };

  const handleDuplicateDesign = (design) => {
    const newDesign = {
      ...design,
      id: Date.now().toString(),
      name: `${design.name} (Copy)`,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };

    const designs = JSON.parse(localStorage.getItem('store-designs') || '[]');
    designs.push(newDesign);
    localStorage.setItem('store-designs', JSON.stringify(designs));
    window.location.reload();
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-4xl h-[80vh] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <div>
            <h2 className="text-xl font-semibold text-gray-900">Saved Designs</h2>
            <p className="text-sm text-gray-600 mt-1">
              {savedDesigns.length} design{savedDesigns.length !== 1 ? 's' : ''} saved
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-md transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Search */}
        <div className="p-6 border-b">
          <input
            type="text"
            placeholder="Search designs..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Designs List */}
        <div className="flex-1 overflow-y-auto p-6">
          {filteredDesigns.length === 0 ? (
            <div className="text-center py-12">
              <Save className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">
                {savedDesigns.length === 0 ? 'No saved designs' : 'No designs found'}
              </h3>
              <p className="text-gray-600">
                {savedDesigns.length === 0
                  ? 'Create and save your first design to see it here.'
                  : 'Try adjusting your search terms.'
                }
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredDesigns.map((design) => (
                <div key={design.id} className="border border-gray-200 rounded-lg overflow-hidden hover:shadow-md transition-shadow">
                  {/* Design Preview */}
                  <div className="h-32 bg-gradient-to-br from-blue-50 to-purple-50 flex items-center justify-center">
                    <div className="text-center">
                      <div className="w-8 h-8 bg-blue-600 rounded mx-auto mb-2"></div>
                      <div className="text-xs text-gray-600">
                        {design.components?.length || 0} components
                      </div>
                    </div>
                  </div>

                  {/* Design Info */}
                  <div className="p-4">
                    <h3 className="font-medium text-gray-900 mb-1 truncate">
                      {design.name}
                    </h3>

                    {design.description && (
                      <p className="text-sm text-gray-600 mb-2 line-clamp-2">
                        {design.description}
                      </p>
                    )}

                    <div className="flex items-center text-xs text-gray-500 mb-3">
                      <Clock className="w-3 h-3 mr-1" />
                      {formatDate(design.updatedAt || design.createdAt)}
                    </div>

                    {design.theme && (
                      <div className="mb-3">
                        <span className="inline-block px-2 py-1 bg-gray-100 text-xs text-gray-700 rounded">
                          {design.theme.name}
                        </span>
                      </div>
                    )}

                    {/* Actions */}
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={() => handleLoadDesign(design)}
                        className="flex-1 px-3 py-2 bg-blue-600 text-white text-sm rounded hover:bg-blue-700 transition-colors flex items-center justify-center"
                      >
                        <Eye className="w-4 h-4 mr-1" />
                        Load
                      </button>

                      <button
                        onClick={() => handleDuplicateDesign(design)}
                        className="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded transition-colors"
                        title="Duplicate"
                      >
                        <Copy className="w-4 h-4" />
                      </button>

                      <button
                        onClick={() => handleDeleteDesign(design)}
                        className="p-2 text-red-600 hover:text-red-700 hover:bg-red-50 rounded transition-colors"
                        title="Delete"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="p-6 border-t bg-gray-50">
          <div className="flex items-center justify-between">
            <div className="text-sm text-gray-600">
              Designs are saved locally in your browser
            </div>
            <button
              onClick={onClose}
              className="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 transition-colors"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SavedDesigns;