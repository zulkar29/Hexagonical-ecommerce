import React from 'react';
import { Zap, HelpCircle, Monitor, Info } from 'lucide-react';
import SettingsCard from '../SettingsCard';
import ToggleSetting from '../ToggleSetting';

const FeaturesTab = ({ settings, onSettingChange }) => (
  <div className="space-y-6">
    <SettingsCard
      title="Core Features"
      description="Enable or disable core platform features"
      icon={Zap}
    >
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <ToggleSetting
          label="Multi-Tenant Mode"
          description="Enable separate data isolation for each tenant"
          checked={settings.features.multiTenantMode}
          onChange={(checked) => onSettingChange('features', 'multiTenantMode', checked)}
        />

        <ToggleSetting
          label="Advanced Analytics"
          description="Enable detailed analytics and reporting"
          checked={settings.features.advancedAnalytics}
          onChange={(checked) => onSettingChange('features', 'advancedAnalytics', checked)}
        />

        <ToggleSetting
          label="Custom Branding"
          description="Allow tenants to customize their branding"
          checked={settings.features.customBranding}
          onChange={(checked) => onSettingChange('features', 'customBranding', checked)}
        />

        <ToggleSetting
          label="API Access"
          description="Provide REST API access to tenants"
          checked={settings.features.apiAccess}
          onChange={(checked) => onSettingChange('features', 'apiAccess', checked)}
        />

        <ToggleSetting
          label="Webhook Support"
          description="Enable webhook notifications for events"
          checked={settings.features.webhookSupport}
          onChange={(checked) => onSettingChange('features', 'webhookSupport', checked)}
        />

        <ToggleSetting
          label="Export Functionality"
          description="Allow data export in various formats"
          checked={settings.features.exportFunctionality}
          onChange={(checked) => onSettingChange('features', 'exportFunctionality', checked)}
        />
      </div>
    </SettingsCard>

    <SettingsCard
      title="Support Features"
      description="Configure customer support and help features"
      icon={HelpCircle}
    >
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <ToggleSetting
          label="Support Ticketing"
          description="Built-in support ticket system"
          checked={settings.features.supportTicketing}
          onChange={(checked) => onSettingChange('features', 'supportTicketing', checked)}
        />

        <ToggleSetting
          label="Knowledge Base"
          description="Self-service knowledge base articles"
          checked={settings.features.knowledgeBase}
          onChange={(checked) => onSettingChange('features', 'knowledgeBase', checked)}
        />

        <ToggleSetting
          label="Live Chat"
          description="Real-time chat support for users"
          checked={settings.features.liveChat}
          onChange={(checked) => onSettingChange('features', 'liveChat', checked)}
          disabled={true}
        />

        <ToggleSetting
          label="Backup System"
          description="Automated daily backups"
          checked={settings.features.backupSystem}
          onChange={(checked) => onSettingChange('features', 'backupSystem', checked)}
        />
      </div>
    </SettingsCard>

    <SettingsCard
      title="Mobile & Offline Features"
      description="Configure mobile app and offline capabilities"
      icon={Monitor}
      badge="Beta"
    >
      <ToggleSetting
        label="Mobile App"
        description="Enable mobile app access for tenants"
        checked={settings.features.mobileApp}
        onChange={(checked) => onSettingChange('features', 'mobileApp', checked)}
      />

      <ToggleSetting
        label="Offline Mode"
        description="Allow limited functionality when offline"
        checked={settings.features.offlineMode}
        onChange={(checked) => onSettingChange('features', 'offlineMode', checked)}
        disabled={true}
      />

      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <div className="flex items-start gap-3">
          <Info className="h-5 w-5 text-blue-600 mt-0.5" />
          <div>
            <h4 className="text-sm font-medium text-blue-800">Mobile App Status</h4>
            <p className="text-sm text-blue-700 mt-1">
              Mobile app is currently in beta testing. Offline mode will be available in the next release.
            </p>
          </div>
        </div>
      </div>
    </SettingsCard>
  </div>
);

export default FeaturesTab;