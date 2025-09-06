import React, { useEffect } from 'react';
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
import { restrictToVerticalAxis } from '@dnd-kit/modifiers';
import ComponentLibrary from './components/ComponentLibrary';
import CanvasArea from './components/CanvasArea';
import CustomizationPanel from './components/CustomizationPanel';
import PreviewPanel from './components/PreviewPanel';
import ThemeCustomizer from './components/ThemeCustomizer';

import { applyThemeToDocument, resetTheme } from '../../utils/themeUtils';
import {
  componentsAtom,
  selectedComponentAtom,
  selectedThemeAtom,
  previewModeAtom,
  showThemeCustomizerAtom,
  addComponentAtom,
  updateComponentAtom,
  deleteComponentAtom,
  reorderComponentsAtom,
  applyThemeAtom
} from './store/builderAtoms';

const ComponentBuilder = () => {
  const [components] = useAtom(componentsAtom);
  const [selectedComponent, setSelectedComponent] = useAtom(selectedComponentAtom);
  const [previewMode, setPreviewMode] = useAtom(previewModeAtom);
  const [selectedTheme] = useAtom(selectedThemeAtom);
  const [showThemeCustomizer] = useAtom(showThemeCustomizerAtom);
  
  const addComponent = useSetAtom(addComponentAtom);
  const updateComponent = useSetAtom(updateComponentAtom);
  const deleteComponent = useSetAtom(deleteComponentAtom);
  const reorderComponents = useSetAtom(reorderComponentsAtom);
  const applyTheme = useSetAtom(applyThemeAtom);

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const handleDragEnd = (event) => {
    const { active, over } = event;

    if (!over) return;

    // If dragging from library to canvas
    if (active.data.current?.type === 'library-item' && over.id === 'canvas') {
      const template = componentTemplates.find(t => t.id === active.id);
      if (template) {
        const newComponent = {
          id: `${template.id}-${Date.now()}`,
          templateId: template.id,
          type: template.type,
          name: template.name,
          props: { ...template.defaultProps },
          styles: { ...template.defaultStyles }
        };
        addComponent(newComponent);
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

  const handleSave = () => {
    // TODO: Implement save functionality
    console.log('Saving components:', components);
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
    // Apply theme using Jotai action
    applyTheme(theme);
    
    // Apply theme styles to document
    applyThemeToDocument(theme);
    
    // Create components from theme layout if canvas is empty
    if (components.length === 0) {
      const themeComponents = createComponentsFromTheme(theme);
      themeComponents.forEach(component => addComponent(component));
    }
    
    console.log('Theme applied:', theme.name);
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
            <p className="text-sm text-gray-600">Drag and drop components to build your layout</p>
          </div>
          <div className="flex items-center space-x-4">
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
              className="px-4 py-2 bg-green-600 text-white rounded-md text-sm font-medium hover:bg-green-700 transition-colors"
            >
              Save
            </button>
            <button
              onClick={handleExport}
              className="px-4 py-2 bg-purple-600 text-white rounded-md text-sm font-medium hover:bg-purple-700 transition-colors"
            >
              Export
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
          modifiers={[restrictToVerticalAxis]}
        >
          {/* Component Library */}
          {!previewMode && (
            <div className="w-80 bg-white border-r border-gray-200 overflow-y-auto">
              <ComponentLibrary 
                onThemeSelect={handleThemeSelect}
                selectedTheme={selectedTheme}
              />
            </div>
          )}

          {/* Canvas Area */}
          <div className="flex-1 flex">
            <div className={`${previewMode ? 'w-full' : 'flex-1'} overflow-y-auto ${
              selectedTheme 
                ? 'bg-[var(--theme-background,#f9fafb)]' 
                : 'bg-gray-50'
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
    </div>
  );
};

export default ComponentBuilder;