import React from 'react';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

const InputSetting = ({
  label,
  description,
  value,
  onChange,
  type = "text",
  placeholder,
  suffix,
  prefix
}) => (
  <div className="space-y-2">
    <Label className="font-medium">{label}</Label>
    {description && (
      <p className="text-sm text-muted-foreground">{description}</p>
    )}
    <div className="flex">
      {prefix && (
        <span className="inline-flex items-center px-3 rounded-l-md border border-r-0 border-input bg-muted text-muted-foreground text-sm">
          {prefix}
        </span>
      )}
      <Input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className={prefix ? "rounded-l-none" : suffix ? "rounded-r-none" : ""}
      />
      {suffix && (
        <span className="inline-flex items-center px-3 rounded-r-md border border-l-0 border-input bg-muted text-muted-foreground text-sm">
          {suffix}
        </span>
      )}
    </div>
  </div>
);

export default InputSetting;