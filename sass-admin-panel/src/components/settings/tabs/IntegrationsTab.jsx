import React from 'react';
import { BarChart3, MessageSquare, Phone, Globe, Eye } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import SettingsCard from '../SettingsCard';
import ToggleSetting from '../ToggleSetting';
import InputSetting from '../InputSetting';

const IntegrationsTab = ({ settings, onSettingChange }) => (
  <div className="space-y-6">
    <SettingsCard
      title="Google Analytics"
      description="Configure Google Analytics tracking"
      icon={BarChart3}
      badge={settings.integrations.googleAnalytics.enabled ? "Connected" : "Disconnected"}
    >
      <ToggleSetting
        label="Enable Google Analytics"
        description="Track user behavior and platform usage"
        checked={settings.integrations.googleAnalytics.enabled}
        onChange={(checked) => onSettingChange('integrations', 'googleAnalytics.enabled', checked)}
      />

      {settings.integrations.googleAnalytics.enabled && (
        <div className="space-y-4 pl-4 border-l-2 border-muted">
          <InputSetting
            label="Tracking ID"
            value={settings.integrations.googleAnalytics.trackingId}
            onChange={(value) => onSettingChange('integrations', 'googleAnalytics.trackingId', value)}
            placeholder="GA-XXXX-XXXX"
          />

          <ToggleSetting
            label="Enhanced Ecommerce Tracking"
            description="Track subscription and payment events"
            checked={settings.integrations.googleAnalytics.ecommerceTracking}
            onChange={(checked) => onSettingChange('integrations', 'googleAnalytics.ecommerceTracking', checked)}
          />
        </div>
      )}
    </SettingsCard>

    <SettingsCard
      title="WhatsApp Business"
      description="Configure WhatsApp Business API integration"
      icon={MessageSquare}
      badge={settings.integrations.whatsapp.enabled ? "Connected" : "Disconnected"}
    >
      <ToggleSetting
        label="Enable WhatsApp Integration"
        description="Send notifications via WhatsApp Business API"
        checked={settings.integrations.whatsapp.enabled}
        onChange={(checked) => onSettingChange('integrations', 'whatsapp.enabled', checked)}
      />

      {settings.integrations.whatsapp.enabled && (
        <div className="space-y-4 pl-4 border-l-2 border-muted">
          <InputSetting
            label="Business Phone Number"
            value={settings.integrations.whatsapp.businessNumber}
            onChange={(value) => onSettingChange('integrations', 'whatsapp.businessNumber', value)}
            placeholder="+8801XXXXXXXXX"
          />

          <div className="space-y-2">
            <Label className="font-medium">API Token</Label>
            <div className="flex gap-2">
              <Input
                type="password"
                value={settings.integrations.whatsapp.apiToken}
                onChange={(e) => onSettingChange('integrations', 'whatsapp.apiToken', e.target.value)}
                placeholder="Enter WhatsApp API token"
              />
              <Button variant="outline" size="sm">
                <Eye className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      )}
    </SettingsCard>

    <SettingsCard
      title="SMS Provider"
      description="Configure SMS service provider settings"
      icon={Phone}
    >
      <div className="space-y-4">
        <div className="space-y-2">
          <Label className="font-medium">SMS Provider</Label>
          <Select
            value={settings.integrations.sms.provider}
            onValueChange={(value) => onSettingChange('integrations', 'sms.provider', value)}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ssl_wireless">SSL Wireless</SelectItem>
              <SelectItem value="robi_axiata">Robi Axiata</SelectItem>
              <SelectItem value="banglalink">Banglalink</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <InputSetting
          label="Sender ID"
          description="SMS sender ID (approved by provider)"
          value={settings.integrations.sms.senderId}
          onChange={(value) => onSettingChange('integrations', 'sms.senderId', value)}
          placeholder="SHOPOWNER"
        />

        <div className="space-y-2">
          <Label className="font-medium">API Key</Label>
          <div className="flex gap-2">
            <Input
              type="password"
              value={settings.integrations.sms.apiKey}
              onChange={(e) => onSettingChange('integrations', 'sms.apiKey', e.target.value)}
              placeholder="Enter SMS provider API key"
            />
            <Button variant="outline" size="sm">
              <Eye className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard
      title="Facebook Integration"
      description="Configure Facebook app and pixel integration"
      icon={Globe}
      badge={settings.integrations.facebook.enabled ? "Connected" : "Disconnected"}
    >
      <ToggleSetting
        label="Enable Facebook Integration"
        description="Connect with Facebook for social login and marketing"
        checked={settings.integrations.facebook.enabled}
        onChange={(checked) => onSettingChange('integrations', 'facebook.enabled', checked)}
      />

      {settings.integrations.facebook.enabled && (
        <div className="space-y-4 pl-4 border-l-2 border-muted">
          <InputSetting
            label="Facebook App ID"
            value={settings.integrations.facebook.appId}
            onChange={(value) => onSettingChange('integrations', 'facebook.appId', value)}
            placeholder="Enter Facebook App ID"
          />

          <InputSetting
            label="Facebook Pixel ID"
            value={settings.integrations.facebook.pixelId}
            onChange={(value) => onSettingChange('integrations', 'facebook.pixelId', value)}
            placeholder="Enter Facebook Pixel ID"
          />
        </div>
      )}
    </SettingsCard>
  </div>
);

export default IntegrationsTab;