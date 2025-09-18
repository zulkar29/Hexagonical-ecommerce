import React from 'react';
import { DatePicker } from '@/components/ui/date-picker';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { cn } from '@/lib/utils';

const DateTimeSetting = ({
  label,
  description,
  value,
  onChange,
  type = 'date', // 'date', 'time', 'datetime'
  disabled = false,
  className,
  ...props
}) => {
  const handleDateChange = (date) => {
    onChange(date);
  };

  const handleTimeChange = (e) => {
    onChange(e.target.value);
  };

  const renderInput = () => {
    switch (type) {
      case 'date':
        return (
          <DatePicker
            date={value}
            onDateChange={handleDateChange}
            disabled={disabled}
            className={className}
            {...props}
          />
        );
      case 'time':
        return (
          <Input
            type="time"
            value={value || ''}
            onChange={handleTimeChange}
            disabled={disabled}
            className={className}
            {...props}
          />
        );
      case 'datetime':
        return (
          <div className="space-y-2">
            <DatePicker
              date={value?.date}
              onDateChange={(date) => onChange({ ...value, date })}
              disabled={disabled}
              className={className}
            />
            <Input
              type="time"
              value={value?.time || ''}
              onChange={(e) => onChange({ ...value, time: e.target.value })}
              disabled={disabled}
              className={className}
            />
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <div className="space-y-2">
      <div className="space-y-1">
        <Label className="text-sm font-medium text-gray-900">
          {label}
        </Label>
        {description && (
          <p className="text-sm text-gray-500">
            {description}
          </p>
        )}
      </div>
      {renderInput()}
    </div>
  );
};

export default DateTimeSetting;