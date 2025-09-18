import React from 'react';
import { Lock, Clock, Database, Shield, CheckCircle, Calendar } from 'lucide-react';
import SettingsCard from '../SettingsCard';
import ToggleSetting from '../ToggleSetting';
import InputSetting from '../InputSetting';
import DateTimeSetting from '../DateTimeSetting';

const SecurityTab = ({ settings, onSettingChange }) => (
  <div className="space-y-6">
    <SettingsCard
      title="Password Policy"
      description="Configure password requirements for all users"
      icon={Lock}
    >
      <InputSetting
        label="Minimum Password Length"
        description="Minimum number of characters required"
        value={settings.security.passwordMinLength}
        onChange={(value) => onSettingChange('security', 'passwordMinLength', parseInt(value) || 0)}
        type="number"
        suffix="characters"
      />

      <div className="space-y-3">
        <ToggleSetting
          label="Require Uppercase Letters"
          description="Password must contain at least one uppercase letter"
          checked={settings.security.requireUppercase}
          onChange={(checked) => onSettingChange('security', 'requireUppercase', checked)}
        />

        <ToggleSetting
          label="Require Numbers"
          description="Password must contain at least one number"
          checked={settings.security.requireNumbers}
          onChange={(checked) => onSettingChange('security', 'requireNumbers', checked)}
        />

        <ToggleSetting
          label="Require Special Characters"
          description="Password must contain at least one special character"
          checked={settings.security.requireSpecialChars}
          onChange={(checked) => onSettingChange('security', 'requireSpecialChars', checked)}
        />
      </div>

      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <div className="flex items-start gap-3">
          <Shield className="h-5 w-5 text-blue-600 mt-0.5" />
          <div>
            <h4 className="text-sm font-medium text-blue-800">Password Strength Preview</h4>
            <p className="text-sm text-blue-700 mt-1">
              Current policy requires: {settings.security.passwordMinLength}+ characters
              {settings.security.requireUppercase && ', uppercase letters'}
              {settings.security.requireNumbers && ', numbers'}
              {settings.security.requireSpecialChars && ', special characters'}
            </p>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard
      title="Session Management"
      description="Configure user session and login security"
      icon={Clock}
    >
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <InputSetting
          label="Session Timeout"
          description="Hours before automatic logout"
          value={settings.security.sessionTimeout}
          onChange={(value) => onSettingChange('security', 'sessionTimeout', parseInt(value) || 0)}
          type="number"
          suffix="hours"
        />

        <InputSetting
          label="Max Login Attempts"
          description="Failed attempts before account lockout"
          value={settings.security.maxLoginAttempts}
          onChange={(value) => onSettingChange('security', 'maxLoginAttempts', parseInt(value) || 0)}
          type="number"
          suffix="attempts"
        />

        <InputSetting
          label="Lockout Duration"
          description="Minutes to lock account after max attempts"
          value={settings.security.lockoutDuration}
          onChange={(value) => onSettingChange('security', 'lockoutDuration', parseInt(value) || 0)}
          type="number"
          suffix="minutes"
        />

        <InputSetting
          label="API Rate Limit"
          description="API requests per hour per user"
          value={settings.security.apiRateLimit}
          onChange={(value) => onSettingChange('security', 'apiRateLimit', parseInt(value) || 0)}
          type="number"
          suffix="req/hour"
        />
      </div>

      <ToggleSetting
        label="Two-Factor Authentication Required"
        description="Require 2FA for all admin users"
        checked={settings.security.twoFactorRequired}
        onChange={(checked) => onSettingChange('security', 'twoFactorRequired', checked)}
      />

      <ToggleSetting
        label="IP Whitelisting"
        description="Restrict admin access to specific IP addresses"
        checked={settings.security.ipWhitelisting}
        onChange={(checked) => onSettingChange('security', 'ipWhitelisting', checked)}
        disabled={true}
      />
    </SettingsCard>

    <SettingsCard
      title="Data Protection"
      description="Configure data retention and audit settings"
      icon={Database}
    >
      <InputSetting
        label="Data Retention Period"
        description="Days to retain user data after account deletion"
        value={settings.security.dataRetention}
        onChange={(value) => onSettingChange('security', 'dataRetention', parseInt(value) || 0)}
        type="number"
        suffix="days"
      />

      <ToggleSetting
        label="Audit Logging"
        description="Log all admin actions for compliance"
        checked={settings.security.auditLogging}
        onChange={(checked) => onSettingChange('security', 'auditLogging', checked)}
      />

      <div className="bg-green-50 border border-green-200 rounded-lg p-4">
        <div className="flex items-start gap-3">
          <CheckCircle className="h-5 w-5 text-green-600 mt-0.5" />
          <div>
            <h4 className="text-sm font-medium text-green-800">GDPR Compliance</h4>
            <p className="text-sm text-green-700 mt-1">
              Current settings are compliant with GDPR data protection requirements.
            </p>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard
      title="Backup & Recovery"
      description="Configure automated backup schedules and data recovery settings"
      icon={Calendar}
    >
      <ToggleSetting
        label="Automated Backups"
        description="Enable automatic daily backups of platform data"
        checked={settings.security?.automatedBackups || false}
        onChange={(checked) => onSettingChange('security', 'automatedBackups', checked)}
      />

      <DateTimeSetting
        label="Backup Schedule"
        description="Set the time for daily automated backups"
        value={settings.security?.backupTime || '02:00'}
        onChange={(value) => onSettingChange('security', 'backupTime', value)}
        type="time"
        disabled={!settings.security?.automatedBackups}
      />

      <InputSetting
        label="Backup Retention"
        description="Number of days to keep backup files"
        value={settings.security?.backupRetention || 30}
        onChange={(value) => onSettingChange('security', 'backupRetention', parseInt(value) || 0)}
        type="number"
        suffix="days"
      />

      <ToggleSetting
        label="Offsite Backup"
        description="Store backups in external cloud storage"
        checked={settings.security?.offsiteBackup || false}
        onChange={(checked) => onSettingChange('security', 'offsiteBackup', checked)}
      />
    </SettingsCard>
  </div>
);

export default SecurityTab;