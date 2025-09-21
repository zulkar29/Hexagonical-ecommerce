import React, { useEffect, useState } from 'react';
import { useAtom, useSetAtom } from 'jotai';
import { themeTemplates, sampleLayouts, componentTemplates } from '../../data/componentTemplates';
import { 
  DndContext, 
  closestCenter,
  useSensor,
  useSensors,
  PointerSensor,
  KeyboardSensor
} from '@dnd-kit/core';
import { 
  arrayMove, 
  SortableContext, 
  verticalListSortingStrategy,
  sortableKeyboardCoordinates
} from '@dnd-kit/sortable';
// import { restrictToVerticalAxis } from '@dnd-kit/modifiers';
import ComponentLibrary from './components/ComponentLibrary';
import CanvasArea from './components/CanvasArea';
import CustomizationPanel from './components/CustomizationPanel';
import PreviewPanel from './components/PreviewPanel';
import ThemeCustomizer from './components/ThemeCustomizer';
import SavedDesigns from './components/SavedDesigns';

import { applyThemeToDocument, resetTheme } from '../../utils/themeUtils';
import {
  componentsAtom,
  selectedComponentAtom,
  selectedThemeAtom,
  previewModeAtom,
  showThemeCustomizerAtom,
  addComponentAtom,
  insertComponentAtom,
  updateComponentAtom,
  deleteComponentAtom,
  reorderComponentsAtom,
  applyThemeAtom,
  undoAtom,
  redoAtom,
  saveDesignAtom,
  loadDesignAtom,
  getSavedDesignsAtom,
  saveStatusAtom,
  designNameAtom,
  builderSettingsAtom,
  historyAtom
} from './store/builderAtoms';

const ComponentBuilder = () => {
  const [components] = useAtom(componentsAtom);
  const [selectedComponent, setSelectedComponent] = useAtom(selectedComponentAtom);
  const [previewMode, setPreviewMode] = useAtom(previewModeAtom);
  const [selectedTheme] = useAtom(selectedThemeAtom);
  const [showThemeCustomizer] = useAtom(showThemeCustomizerAtom);
  
  const addComponent = useSetAtom(addComponentAtom);
  const insertComponent = useSetAtom(insertComponentAtom);
  const updateComponent = useSetAtom(updateComponentAtom);
  const deleteComponent = useSetAtom(deleteComponentAtom);
  const reorderComponents = useSetAtom(reorderComponentsAtom);
  const applyTheme = useSetAtom(applyThemeAtom);
  const undo = useSetAtom(undoAtom);
  const redo = useSetAtom(redoAtom);
  const saveDesign = useSetAtom(saveDesignAtom);
  const loadDesign = useSetAtom(loadDesignAtom);
  const [saveStatus] = useAtom(saveStatusAtom);
  const [designName, setDesignName] = useAtom(designNameAtom);
  const [builderSettings, setBuilderSettings] = useAtom(builderSettingsAtom);
  const [history] = useAtom(historyAtom);
  const [savedDesigns] = useAtom(getSavedDesignsAtom);
  const [showSavedDesigns, setShowSavedDesigns] = useState(false);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8, // 8px movement required to start drag
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const handleDragEnd = (event) => {
    const { active, over } = event;

    console.log('🔍 Drag ended:', {
      activeId: active?.id,
      overId: over?.id,
      activeType: active?.data?.current?.type,
      activeData: active?.data?.current
    });

    if (!over) {
      console.log('❌ No drop target found');
      return;
    }

    // If dragging theme to canvas
    if (active.data.current?.type === 'theme-item' && over.id === 'canvas') {
      const theme = active.data.current.theme;
      console.log('Applying theme via drag:', theme);
      handleThemeSelect(theme);
      return;
    }

    // If dragging from library to canvas (general drop)
    if (active.data.current?.type === 'library-item' && over.id === 'canvas') {
      const template = active.data.current?.template || componentTemplates.find(t => t.id === active.id);
      console.log('📦 Found template:', template);

      if (template) {
        const newComponent = {
          id: `${template.id}-${Date.now()}`,
          templateId: template.id,
          type: template.type,
          name: template.name,
          props: { ...template.defaultProps },
          styles: { ...template.defaultStyles }
        };
        console.log('✅ Adding component to end:', newComponent);
        addComponent(newComponent);
      } else {
        console.error('❌ Template not found for ID:', active.id);
        console.log('Available templates:', componentTemplates.map(t => t.id));
      }
      return;
    }

    // If dragging from library to a specific drop zone
    if (active.data.current?.type === 'library-item' && over.data.current?.type === 'drop-zone') {
      const template = active.data.current?.template || componentTemplates.find(t => t.id === active.id);
      const insertIndex = over.data.current?.index || 0;

      console.log('📦 Inserting template at index:', insertIndex, 'template:', template?.name);

      if (template) {
        const newComponent = {
          id: `${template.id}-${Date.now()}`,
          templateId: template.id,
          type: template.type,
          name: template.name,
          props: { ...template.defaultProps },
          styles: { ...template.defaultStyles }
        };
        console.log('✅ Inserting component at index:', insertIndex, newComponent);

        // Insert component at specific index using the new insertComponent function
        insertComponent({ component: newComponent, index: insertIndex });
      } else {
        console.error('❌ Template not found for ID:', active.id);
      }
      return;
    }

    // If reordering components in canvas
    if (active.data.current?.type === 'canvas-item' && over.data.current?.type === 'canvas-item') {
      const oldIndex = components.findIndex(c => c.id === active.id);
      const newIndex = components.findIndex(c => c.id === over.id);
      
      if (oldIndex !== newIndex) {
        reorderComponents({ oldIndex, newIndex });
      }
    }
  };

  const handleComponentSelect = (component) => {
    setSelectedComponent(component);
  };

  const handleComponentUpdate = (componentId, updates) => {
    updateComponent({ componentId, updates });
  };

  const handleComponentDelete = (componentId) => {
    deleteComponent(componentId);
  };

  const handleSave = async () => {
    try {
      await saveDesign({ name: designName, description: 'Store design' });
    } catch (error) {
      console.error('Failed to save design:', error);
    }
  };

  const handleLoad = (design) => {
    if (window.confirm('Loading a design will replace your current work. Continue?')) {
      loadDesign(design);
    }
  };

  const handleUndo = () => {
    if (history.past.length > 0) {
      undo();
    }
  };

  const handleRedo = () => {
    if (history.future.length > 0) {
      redo();
    }
  };

  const handleQuickStart = (template) => {
    console.log('🚀 Quick starting with template:', template.name);

    // Clear existing components
    components.forEach(comp => deleteComponent(comp.id));

    // Create components from template
    template.components.forEach((componentConfig, index) => {
      const newComponent = {
        id: `${componentConfig.type}-${Date.now()}-${index}`,
        templateId: `${componentConfig.type}-template`,
        type: componentConfig.type,
        name: componentConfig.title || componentConfig.companyName || `${componentConfig.type} Component`,
        props: componentConfig,
        styles: selectedTheme ? {
          ...selectedTheme.components?.[componentConfig.type],
          backgroundColor: selectedTheme.colors?.background || '#ffffff',
          textColor: selectedTheme.colors?.text || '#000000'
        } : {
          backgroundColor: '#ffffff',
          textColor: '#000000'
        }
      };
      addComponent(newComponent);
    });

    console.log('✅ Quick start template applied successfully');
  };


  const createComponentsFromTheme = (theme) => {
    const themeLayout = sampleLayouts[theme.id];
    if (!themeLayout) return [];

    return themeLayout.map((layoutComponent, index) => ({
      id: `${layoutComponent.type}-${Date.now()}-${index}`,
      templateId: `${layoutComponent.type}-${theme.id}`,
      type: layoutComponent.type,
      name: layoutComponent.content?.title || `${layoutComponent.type} Component`,
      props: {
        ...layoutComponent.content
      },
      styles: {
        ...theme.components[layoutComponent.type],
        backgroundColor: theme.colors?.background || '#ffffff',
        textColor: theme.colors?.text || '#000000'
      },
      theme: theme.id
    }));
  };

  const handleThemeSelect = (theme) => {
    try {
      console.log('🎨 Applying theme:', theme.name);

      // Apply theme using Jotai action
      applyTheme(theme);

      // Apply theme styles to document
      applyThemeToDocument(theme);

      // If canvas is empty, create starter components from theme
      if (components.length === 0) {
        console.log('📱 Creating starter layout for empty canvas');
        const themeComponents = createComponentsFromTheme(theme);
        console.log('Generated components:', themeComponents.length);

        if (themeComponents.length > 0) {
          themeComponents.forEach(component => {
            addComponent(component);
          });
        }
      } else {
        // Apply theme to existing components
        console.log('🔄 Applying theme to existing components');
        const updatedComponents = components.map(component => ({
          ...component,
          theme: theme.id,
          styles: {
            ...component.styles,
            ...theme.components?.[component.type],
            backgroundColor: theme.colors?.background || component.styles?.backgroundColor,
            textColor: theme.colors?.text || component.styles?.textColor
          }
        }));

        // Clear and re-add components with theme applied
        components.forEach(comp => deleteComponent(comp.id));
        updatedComponents.forEach(comp => addComponent(comp));
      }

      console.log('✅ Theme applied successfully:', theme.name);
    } catch (error) {
      console.error('❌ Error applying theme:', error);
    }
  };

  const handleThemeUpdate = (updatedTheme) => {
    // Apply updated theme using Jotai action
    applyTheme(updatedTheme);
    
    // Apply updated theme styles
    applyThemeToDocument(updatedTheme);
  };

  const handleResetTheme = () => {
    if (selectedTheme) {
      // Reset to original theme
      applyThemeToDocument(selectedTheme);
      applyTheme(selectedTheme);
    }
  };


  // Cleanup theme on unmount
  useEffect(() => {
    return () => {
      resetTheme();
    };
  }, []);

  const handleExport = () => {
    // TODO: Implement export functionality
    const exportData = {
      components,
      selectedTheme,
      timestamp: new Date().toISOString()
    };
    console.log('Exporting:', exportData);
    
    // Download as JSON for now
    const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'store-design.json';
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="h-screen flex flex-col bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Store Designer</h1>
            <p className="text-sm text-gray-600">Design your professional storefront with drag-and-drop simplicity</p>
          </div>
          <div className="flex items-center space-x-4">
            {/* Design Name Input */}
            <input
              type="text"
              value={designName}
              onChange={(e) => setDesignName(e.target.value)}
              className="px-3 py-1 border border-gray-300 rounded-md text-sm"
              placeholder="Design name"
            />

            {/* Undo/Redo */}
            <div className="flex items-center space-x-1">
              <button
                onClick={handleUndo}
                disabled={history.past.length === 0}
                className="p-2 text-gray-600 hover:text-gray-900 disabled:opacity-50 disabled:cursor-not-allowed"
                title="Undo"
              >
                ↶
              </button>
              <button
                onClick={handleRedo}
                disabled={history.future.length === 0}
                className="p-2 text-gray-600 hover:text-gray-900 disabled:opacity-50 disabled:cursor-not-allowed"
                title="Redo"
              >
                ↷
              </button>
            </div>

            {/* Device Preview Toggle */}
            <div className="flex items-center bg-gray-100 rounded-md p-1">
              {['desktop', 'tablet', 'mobile'].map((device) => (
                <button
                  key={device}
                  onClick={() => setBuilderSettings({ ...builderSettings, devicePreview: device })}
                  className={`px-3 py-1 text-xs font-medium rounded transition-colors ${
                    builderSettings.devicePreview === device
                      ? 'bg-white text-gray-900 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  {device.charAt(0).toUpperCase() + device.slice(1)}
                </button>
              ))}
            </div>

            <button
              onClick={() => setPreviewMode(!previewMode)}
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                previewMode
                  ? 'bg-blue-600 text-white hover:bg-blue-700'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              {previewMode ? 'Edit Mode' : 'Preview Mode'}
            </button>

            <button
              onClick={handleSave}
              disabled={saveStatus === 'saving'}
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                saveStatus === 'saving'
                  ? 'bg-gray-400 text-white cursor-not-allowed'
                  : saveStatus === 'saved'
                  ? 'bg-green-600 text-white hover:bg-green-700'
                  : 'bg-blue-600 text-white hover:bg-blue-700'
              }`}
            >
              {saveStatus === 'saving' ? 'Saving...' : saveStatus === 'saved' ? 'Saved!' : 'Save'}
            </button>

            <button
              onClick={() => setShowSavedDesigns(true)}
              className="px-4 py-2 bg-gray-600 text-white rounded-md text-sm font-medium hover:bg-gray-700 transition-colors"
            >
              Load Design
            </button>

            <button
              onClick={handleExport}
              className="px-4 py-2 bg-gray-600 text-white rounded-md text-sm font-medium hover:bg-gray-700 transition-colors"
            >
              Export
            </button>

            {process.env.NODE_ENV === 'development' && (
              <>
                <button
                  onClick={() => {
                    const testComponent = {
                      id: `test-${Date.now()}`,
                      templateId: 'header-simple',
                      type: 'header',
                      name: 'Test Header',
                      props: { title: 'Test Store', menuItems: ['Home', 'Products', 'About'] },
                      styles: { backgroundColor: '#ffffff', textColor: '#000000' }
                    };
                    addComponent(testComponent);
                  }}
                  className="px-3 py-2 bg-orange-600 text-white rounded-md text-xs font-medium hover:bg-orange-700 transition-colors"
                >
                  + Test
                </button>
                <button
                  onClick={() => {
                    console.log('📊 Debug Info:');
                    console.log('Components in state:', components.length);
                    console.log('Component templates available:', componentTemplates.length);
                    console.log('Canvas element:', document.querySelector('[data-id="canvas"]'));
                  }}
                  className="px-3 py-2 bg-purple-600 text-white rounded-md text-xs font-medium hover:bg-purple-700 transition-colors"
                >
                  Debug
                </button>
              </>
            )}

            <button
              onClick={() => {
                if (components.length === 0) {
                  alert('Please add some components to your store design before publishing.');
                  return;
                }
                if (window.confirm('Publish your store design? This will make it live for customers.')) {
                  // In a real app, this would deploy the design
                  alert('Store published successfully! 🎉\n\nYour customers can now see the new design.');
                }
              }}
              className="px-4 py-2 bg-green-600 text-white rounded-md text-sm font-medium hover:bg-green-700 transition-colors font-semibold"
            >
              Publish Store
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragEnd={handleDragEnd}
          onDragStart={(event) => {
            console.log('🚀 Drag started:', event.active.id, event.active.data.current);
          }}
          onDragOver={(event) => {
            console.log('🎯 Drag over:', event.over?.id);
          }}
        >
          {/* Component Library */}
          {!previewMode && (
            <div className="w-80 bg-white border-r border-gray-200 overflow-y-auto">
              <ComponentLibrary
                onThemeSelect={handleThemeSelect}
                selectedTheme={selectedTheme}
                onQuickStart={handleQuickStart}
              />
            </div>
          )}

          {/* Canvas Area */}
          <div className="flex-1 flex">
            <div className={`${previewMode ? 'w-full' : 'flex-1'} overflow-y-auto transition-all duration-300 ${
              selectedTheme
                ? 'bg-[var(--theme-background,#f9fafb)]'
                : 'bg-gray-50'
            } ${
              builderSettings.devicePreview === 'mobile' ? 'max-w-sm mx-auto' :
              builderSettings.devicePreview === 'tablet' ? 'max-w-4xl mx-auto' :
              'w-full'
            }`}>
              {previewMode ? (
                <PreviewPanel components={components} />
              ) : (
                <SortableContext items={components.map(c => c.id)} strategy={verticalListSortingStrategy}>
                  <CanvasArea
                    components={components}
                    selectedComponent={selectedComponent}
                    onComponentSelect={handleComponentSelect}
                    onComponentDelete={handleComponentDelete}
                    selectedTheme={selectedTheme}
                  />
                </SortableContext>
              )}

              {/* Debug Info */}
              {process.env.NODE_ENV === 'development' && (
                <div className="fixed bottom-4 right-4 bg-white p-2 rounded shadow text-xs">
                  <div>Components: {components.length}</div>
                  <div>Selected: {selectedComponent?.name || 'None'}</div>
                  <div>Theme: {selectedTheme?.name || 'None'}</div>
                </div>
              )}
            </div>

            {/* Customization Panel or Theme Customizer */}
            {!previewMode && (
              <div className="w-80 bg-white border-l border-gray-200 overflow-y-auto">
                {showThemeCustomizer && selectedTheme ? (
                  <ThemeCustomizer
                    selectedTheme={selectedTheme}
                    onThemeUpdate={handleThemeUpdate}
                    onResetTheme={handleResetTheme}
                  />
                ) : selectedComponent ? (
                  <CustomizationPanel
                    component={selectedComponent}
                    onUpdate={(updates) => handleComponentUpdate(selectedComponent.id, updates)}
                  />
                ) : (
                  <div className="p-4 text-center text-gray-500">
                    <p>Select a component to customize or choose a theme to get started</p>
                  </div>
                )}
              </div>
            )}
          </div>
        </DndContext>
      </div>

      {/* Saved Designs Modal */}
      <SavedDesigns
        isOpen={showSavedDesigns}
        onClose={() => setShowSavedDesigns(false)}
        onLoad={(design) => {
          console.log('Loaded design:', design.name);
        }}
      />
    </div>
  );
};

export default ComponentBuilder;