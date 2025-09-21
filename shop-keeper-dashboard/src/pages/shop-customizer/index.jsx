import React, { useEffect, useState } from 'react';
import { useAtom, useSetAtom } from 'jotai';
import { sampleLayouts, componentTemplates } from '../../data/shopTemplates';
import {
  DndContext,
  closestCenter,
  useSensor,
  useSensors,
  PointerSensor,
  KeyboardSensor
} from '@dnd-kit/core';
import {
  SortableContext,
  verticalListSortingStrategy,
  sortableKeyboardCoordinates
} from '@dnd-kit/sortable';

// Modern Tool Components
import ToolNavigation from './components/ToolNavigation';
import ComponentsPanel from './components/ComponentsPanel';
import ThemesPanel from './components/ThemesPanel';
import SettingsPanel from './components/SettingsPanel';
import CanvasArea from './components/CanvasArea';
import PreviewPanel from './components/PreviewPanel';
import PropertiesPanel from './components/PropertiesPanel';

import { applyThemeToDocument, resetTheme } from './utils/themeUtils';
import {
  componentsAtom,
  selectedComponentAtom,
  selectedThemeAtom,
  previewModeAtom,
  addComponentAtom,
  insertComponentAtom,
  updateComponentAtom,
  deleteComponentAtom,
  reorderComponentsAtom,
  applyThemeAtom,
  undoAtom,
  redoAtom,
  saveDesignAtom,
  designNameAtom,
  customizerSettingsAtom,
  historyAtom
} from './store/customizerAtoms';

const StoreDesigner = () => {
  const [activePanel, setActivePanel] = useState('components'); // components, themes, settings, properties
  const [panelCollapsed, setPanelCollapsed] = useState(false);

  // Set document title
  useEffect(() => {
    document.title = 'Customize My Shop - Design Your Storefront';
    return () => {
      document.title = 'Shop Keeper Dashboard';
    };
  }, []);

  const [components] = useAtom(componentsAtom);
  const [selectedComponent, setSelectedComponent] = useAtom(selectedComponentAtom);
  const [previewMode, setPreviewMode] = useAtom(previewModeAtom);
  const [selectedTheme] = useAtom(selectedThemeAtom);

  const addComponent = useSetAtom(addComponentAtom);
  const insertComponent = useSetAtom(insertComponentAtom);
  const updateComponent = useSetAtom(updateComponentAtom);
  const deleteComponent = useSetAtom(deleteComponentAtom);
  const reorderComponents = useSetAtom(reorderComponentsAtom);
  const applyTheme = useSetAtom(applyThemeAtom);
  const undo = useSetAtom(undoAtom);
  const redo = useSetAtom(redoAtom);
  const saveDesign = useSetAtom(saveDesignAtom);
  const [designName, setDesignName] = useAtom(designNameAtom);
  const [customizerSettings, setCustomizerSettings] = useAtom(customizerSettingsAtom);
  const [history] = useAtom(historyAtom);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const handleDragEnd = (event) => {
    const { active, over } = event;

    if (!over) return;

    // If dragging theme to canvas
    if (active.data.current?.type === 'theme-item' && over.id === 'canvas') {
      const theme = active.data.current.theme;
      handleThemeSelect(theme);
      return;
    }

    // If dragging from library to canvas
    if (active.data.current?.type === 'library-item' && over.id === 'canvas') {
      const template = active.data.current?.template || componentTemplates.find(t => t.id === active.id);

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

    // If dragging from library to a specific drop zone
    if (active.data.current?.type === 'library-item' && over.data.current?.type === 'drop-zone') {
      const template = active.data.current?.template || componentTemplates.find(t => t.id === active.id);
      const insertIndex = over.data.current?.index || 0;

      if (template) {
        const newComponent = {
          id: `${template.id}-${Date.now()}`,
          templateId: template.id,
          type: template.type,
          name: template.name,
          props: { ...template.defaultProps },
          styles: { ...template.defaultStyles }
        };
        insertComponent({ component: newComponent, index: insertIndex });
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
    setActivePanel('properties');
  };

  const handleComponentUpdate = (componentId, updates) => {
    updateComponent({ componentId, updates });
  };

  const handleComponentDelete = (componentId) => {
    deleteComponent(componentId);
    if (selectedComponent?.id === componentId) {
      setSelectedComponent(null);
    }
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
      applyTheme(theme);
      applyThemeToDocument(theme);

      if (components.length === 0) {
        const themeComponents = createComponentsFromTheme(theme);
        if (themeComponents.length > 0) {
          themeComponents.forEach(component => {
            addComponent(component);
          });
        }
      } else {
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

        components.forEach(comp => deleteComponent(comp.id));
        updatedComponents.forEach(comp => addComponent(comp));
      }
    } catch (error) {
      console.error('Error applying theme:', error);
    }
  };

  const handleThemeUpdate = (updatedTheme) => {
    applyTheme(updatedTheme);
    applyThemeToDocument(updatedTheme);
  };

  const handleQuickStart = (template) => {
    components.forEach(comp => deleteComponent(comp.id));

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
  };

  // Cleanup theme on unmount
  useEffect(() => {
    return () => {
      resetTheme();
    };
  }, []);

  const renderPanel = () => {
    if (panelCollapsed) return null;

    switch (activePanel) {
      case 'components':
        return (
          <ComponentsPanel
            onThemeSelect={handleThemeSelect}
            selectedTheme={selectedTheme}
            onQuickStart={handleQuickStart}
          />
        );
      case 'themes':
        return (
          <ThemesPanel
            onThemeSelect={handleThemeSelect}
            selectedTheme={selectedTheme}
            onThemeUpdate={handleThemeUpdate}
          />
        );
      case 'settings':
        return (
          <SettingsPanel
            designName={designName}
            setDesignName={setDesignName}
            customizerSettings={customizerSettings}
            setCustomizerSettings={setCustomizerSettings}
            onSave={() => saveDesign(designName)}
            onUndo={undo}
            onRedo={redo}
            history={history}
            components={components}
            selectedTheme={selectedTheme}
          />
        );
      case 'properties':
        return selectedComponent ? (
          <PropertiesPanel
            component={selectedComponent}
            onUpdate={(updates) => handleComponentUpdate(selectedComponent.id, updates)}
            onClose={() => setSelectedComponent(null)}
          />
        ) : (
          <div className="p-6 text-center text-gray-500">
            <p>Select a component to view its properties</p>
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <div className="h-screen flex flex-col bg-gray-50">
      {/* Top Navigation */}
      <ToolNavigation
        activePanel={activePanel}
        setActivePanel={setActivePanel}
        panelCollapsed={panelCollapsed}
        setPanelCollapsed={setPanelCollapsed}
        previewMode={previewMode}
        setPreviewMode={setPreviewMode}
        selectedTheme={selectedTheme}
        components={components}
        designName={designName}
        customizerSettings={customizerSettings}
        setCustomizerSettings={setCustomizerSettings}
      />

      {/* Main Content Area */}
      <div className="flex-1 flex overflow-hidden">
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragEnd={handleDragEnd}
        >
          {/* Left Panel */}
          {!previewMode && (
            <div className={`bg-white border-r border-gray-200 transition-all duration-300 ${
              panelCollapsed ? 'w-0' : 'w-80'
            } overflow-hidden`}>
              {renderPanel()}
            </div>
          )}

          {/* Canvas Area */}
          <div className="flex-1 overflow-hidden">
            <div className={`h-full overflow-y-auto ${
              selectedTheme
                ? 'bg-[var(--theme-background,#f9fafb)]'
                : 'bg-gray-50'
            } ${
              customizerSettings.devicePreview === 'mobile' ? 'flex justify-center' :
              customizerSettings.devicePreview === 'tablet' ? 'flex justify-center' :
              ''
            }`}>
              <div className={`${
                customizerSettings.devicePreview === 'mobile' ? 'w-80 min-h-full' :
                customizerSettings.devicePreview === 'tablet' ? 'w-[768px] min-h-full' :
                'w-full min-h-full'
              } transition-all duration-300`}>
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
            </div>
          </div>
        </DndContext>
      </div>
    </div>
  );
};

export default StoreDesigner;