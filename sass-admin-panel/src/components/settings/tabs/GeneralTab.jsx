import React from 'react';
import { Globe, Settings, Clock, Info } from 'lucide-react';
import SettingsCard from '../SettingsCard';
import ToggleSetting from '../ToggleSetting';
import InputSetting from '../InputSetting';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

const GeneralTab = ({ settings, onSettingChange }) => (
  <div className="space-y-6">
    <SettingsCard
      title="Platform Information"
      description="Basic platform configuration and branding"
      icon={Settings}
    >
      <InputSetting
        label="Platform Name"
        description="The name of your platform displayed to users"
        value={settings.general.platformName}
        onChange={(value) => onSettingChange('general', 'platformName', value)}
        placeholder="Enter platform name"
      />

      <InputSetting
        label="Platform Description"
        description="Brief description of your platform"
        value={settings.general.platformDescription}
        onChange={(value) => onSettingChange('general', 'platformDescription', value)}
        placeholder="Enter platform description"
      />

      <div className="space-y-2">
        <label className="text-sm font-medium">Default Language</label>
        <Select
          value={settings.general.defaultLanguage}
          onValueChange={(value) => onSettingChange('general', 'defaultLanguage', value)}
        >
          <SelectTrigger>
            <SelectValue placeholder="Select language" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="bn">Bengali (বাংলা)</SelectItem>
            <SelectItem value="en">English</SelectItem>
          </SelectContent>
        </Select>
        <p className="text-xs text-muted-foreground">Primary language for the platform interface</p>
      </div>

      <div className="space-y-2">
        <label className="text-sm font-medium">Currency</label>
        <Select
          value={settings.general.currency}
          onValueChange={(value) => onSettingChange('general', 'currency', value)}
        >
          <SelectTrigger>
            <SelectValue placeholder="Select currency" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="BDT">Bangladeshi Taka (৳)</SelectItem>
            <SelectItem value="USD">US Dollar ($)</SelectItem>
            <SelectItem value="EUR">Euro (€)</SelectItem>
          </SelectContent>
        </Select>
        <p className="text-xs text-muted-foreground">Default currency for pricing and billing</p>
      </div>
    </SettingsCard>

    <SettingsCard
      title="Regional Settings"
      description="Configure timezone and date formats"
      icon={Clock}
    >
      <div className="space-y-2">
        <label className="text-sm font-medium">Timezone</label>
        <Select
          value={settings.general.timezone}
          onValueChange={(value) => onSettingChange('general', 'timezone', value)}
        >
          <SelectTrigger>
            <SelectValue placeholder="Select timezone" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="Asia/Dhaka">Asia/Dhaka (GMT+6)</SelectItem>
            <SelectItem value="UTC">UTC (GMT+0)</SelectItem>
            <SelectItem value="America/New_York">America/New_York (EST)</SelectItem>
          </SelectContent>
        </Select>
        <p className="text-xs text-muted-foreground">Default timezone for the platform</p>
      </div>

      <div className="space-y-2">
        <label className="text-sm font-medium">Date Format</label>
        <Select
          value={settings.general.dateFormat}
          onValueChange={(value) => onSettingChange('general', 'dateFormat', value)}
        >
          <SelectTrigger>
            <SelectValue placeholder="Select date format" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="DD/MM/YYYY">DD/MM/YYYY</SelectItem>
            <SelectItem value="MM/DD/YYYY">MM/DD/YYYY</SelectItem>
            <SelectItem value="YYYY-MM-DD">YYYY-MM-DD</SelectItem>
          </SelectContent>
        </Select>
        <p className="text-xs text-muted-foreground">How dates are displayed throughout the platform</p>
      </div>
    </SettingsCard>

    <SettingsCard
      title="System Settings"
      description="Platform-wide system configuration"
      icon={Globe}
    >
      <ToggleSetting
        label="Maintenance Mode"
        description="Put the platform in maintenance mode to prevent user access"
        checked={settings.general.maintenanceMode}
        onChange={(checked) => onSettingChange('general', 'maintenanceMode', checked)}
      />

      <ToggleSetting
        label="Debug Mode"
        description="Enable debug mode for development and troubleshooting"
        checked={settings.general.debugMode}
        onChange={(checked) => onSettingChange('general', 'debugMode', checked)}
      />

      <ToggleSetting
        label="Analytics Tracking"
        description="Enable analytics and usage tracking"
        checked={settings.general.analyticsEnabled}
        onChange={(checked) => onSettingChange('general', 'analyticsEnabled', checked)}
      />

      {settings.general.debugMode && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <div className="flex items-start gap-3">
            <Info className="h-5 w-5 text-yellow-600 mt-0.5" />
            <div>
              <h4 className="text-sm font-medium text-yellow-800">Debug Mode Active</h4>
              <p className="text-sm text-yellow-700 mt-1">
                Debug mode is currently enabled. This may expose sensitive information and should only be used in development environments.
              </p>
            </div>
          </div>
        </div>
      )}
    </SettingsCard>
  </div>
);

export default GeneralTab;