import React from 'react';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';

const ToggleSetting = ({ label, description, checked, onChange, disabled = false }) => (
  <div className="flex items-center justify-between p-3 border rounded-lg">
    <div className="flex-1">
      <div className="flex items-center gap-2">
        <Label className="font-medium">{label}</Label>
        {disabled && <Badge variant="outline" className="text-xs">Pro Only</Badge>}
      </div>
      {description && (
        <p className="text-sm text-muted-foreground mt-1">{description}</p>
      )}
    </div>
    <Switch
      checked={checked}
      onCheckedChange={onChange}
      disabled={disabled}
    />
  </div>
);

export default ToggleSetting;