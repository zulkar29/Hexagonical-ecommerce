import React from 'react';
import { Settings, AlertTriangle, Search, Save, RefreshCw, MoreHorizontal, Download, Upload } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

const SettingsHeader = ({
  hasUnsavedChanges,
  isSaving,
  onSave,
  onExport,
  onImport,
  onReset
}) => (
  <div className="flex items-center gap-2">
    <Button
      variant={hasUnsavedChanges ? "default" : "outline"}
      size="sm"
      onClick={onSave}
      disabled={!hasUnsavedChanges || isSaving}
      className="shrink-0"
    >
      {isSaving ? (
        <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
      ) : (
        <Save className="h-4 w-4 mr-2" />
      )}
      <span className="hidden sm:inline">
        {isSaving ? 'Saving...' : 'Save Changes'}
      </span>
      <span className="sm:hidden">Save</span>
    </Button>

    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="shrink-0">
          <MoreHorizontal className="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuLabel>Platform Actions</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={onExport}>
          <Download className="h-4 w-4 mr-2" />
          Export Settings
        </DropdownMenuItem>
        <DropdownMenuItem onClick={onImport}>
          <Upload className="h-4 w-4 mr-2" />
          Import Settings
        </DropdownMenuItem>
        <DropdownMenuItem onClick={onReset}>
          <RefreshCw className="h-4 w-4 mr-2" />
          Reset to Defaults
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
);

export default SettingsHeader;