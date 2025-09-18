import React from 'react';
import { AlertTriangle, Save, RefreshCw } from 'lucide-react';
import { Button } from '@/components/ui/button';

const UnsavedChangesBar = ({ hasUnsavedChanges, isSaving, onSave, onDiscard }) => {
  if (!hasUnsavedChanges) return null;

  return (
    <div className="border-t bg-yellow-50 border-yellow-200 p-3 md:p-4">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <AlertTriangle className="h-5 w-5 text-yellow-600 shrink-0" />
          <div>
            <p className="text-sm font-medium text-yellow-800">You have unsaved changes</p>
            <p className="text-xs text-yellow-700">Make sure to save your changes before leaving this page.</p>
          </div>
        </div>
        <div className="flex items-center gap-2 shrink-0">
          <Button
            variant="outline"
            size="sm"
            onClick={onDiscard}
            className="shrink-0"
          >
            Discard
          </Button>
          <Button
            size="sm"
            onClick={onSave}
            disabled={isSaving}
            className="shrink-0"
          >
            {isSaving ? (
              <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
            ) : (
              <Save className="h-4 w-4 mr-2" />
            )}
            {isSaving ? 'Saving...' : 'Save Changes'}
          </Button>
        </div>
      </div>
    </div>
  );
};

export default UnsavedChangesBar;