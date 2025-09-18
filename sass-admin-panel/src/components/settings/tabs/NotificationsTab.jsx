import React from 'react';
import { Bell, Mail, MessageSquare, Users, Info, Calendar } from 'lucide-react';
import SettingsCard from '../SettingsCard';
import ToggleSetting from '../ToggleSetting';
import InputSetting from '../InputSetting';
import DateTimeSetting from '../DateTimeSetting';

const NotificationsTab = ({ settings, onSettingChange }) => (
  <div className="space-y-6">
    <SettingsCard
      title="Email Notifications"
      description="Configure platform-wide email notification settings"
      icon={Mail}
    >
      <ToggleSetting
        label="Email Notifications"
        description="Enable email notifications system-wide"
        checked={settings.notifications?.emailNotifications || true}
        onChange={(checked) => onSettingChange('notifications', 'emailNotifications', checked)}
      />

      <ToggleSetting
        label="Welcome Emails"
        description="Send welcome emails to new tenant signups"
        checked={settings.notifications?.welcomeEmails || true}
        onChange={(checked) => onSettingChange('notifications', 'welcomeEmails', checked)}
      />

      <ToggleSetting
        label="Subscription Notifications"
        description="Send emails for subscription changes and renewals"
        checked={settings.notifications?.subscriptionEmails || true}
        onChange={(checked) => onSettingChange('notifications', 'subscriptionEmails', checked)}
      />

      <ToggleSetting
        label="Payment Reminders"
        description="Send payment reminder emails for overdue accounts"
        checked={settings.notifications?.paymentReminders || true}
        onChange={(checked) => onSettingChange('notifications', 'paymentReminders', checked)}
      />

      <InputSetting
        label="Trial Expiry Notice"
        description="Days before trial expires to send notification"
        value={settings.notifications?.trialExpiryNotice || 7}
        onChange={(value) => onSettingChange('notifications', 'trialExpiryNotice', parseInt(value) || 0)}
        type="number"
        suffix="days before"
      />
    </SettingsCard>

    <SettingsCard
      title="SMS & Push Notifications"
      description="Configure mobile and SMS notification settings"
      icon={MessageSquare}
    >
      <ToggleSetting
        label="SMS Notifications"
        description="Enable SMS notifications for critical updates"
        checked={settings.notifications?.smsNotifications || true}
        onChange={(checked) => onSettingChange('notifications', 'smsNotifications', checked)}
      />

      <ToggleSetting
        label="Push Notifications"
        description="Enable push notifications for mobile apps"
        checked={settings.notifications?.pushNotifications || true}
        onChange={(checked) => onSettingChange('notifications', 'pushNotifications', checked)}
      />

      <ToggleSetting
        label="System Alerts"
        description="Send alerts for system maintenance and updates"
        checked={settings.notifications?.systemAlerts || true}
        onChange={(checked) => onSettingChange('notifications', 'systemAlerts', checked)}
      />

      <ToggleSetting
        label="Security Alerts"
        description="Send notifications for security-related events"
        checked={settings.notifications?.securityAlerts || true}
        onChange={(checked) => onSettingChange('notifications', 'securityAlerts', checked)}
      />
    </SettingsCard>

    <SettingsCard
      title="Admin Notifications"
      description="Configure notifications for platform administrators"
      icon={Users}
    >
      <ToggleSetting
        label="New Tenant Alerts"
        description="Notify admins when new tenants sign up"
        checked={settings.notifications?.newTenantAlerts || true}
        onChange={(checked) => onSettingChange('notifications', 'newTenantAlerts', checked)}
      />

      <ToggleSetting
        label="Payment Failure Alerts"
        description="Notify admins of payment failures and issues"
        checked={settings.notifications?.paymentFailureAlerts || true}
        onChange={(checked) => onSettingChange('notifications', 'paymentFailureAlerts', checked)}
      />

      <ToggleSetting
        label="Support Ticket Alerts"
        description="Notify admins of new support tickets"
        checked={settings.notifications?.supportTicketAlerts || true}
        onChange={(checked) => onSettingChange('notifications', 'supportTicketAlerts', checked)}
      />

      <ToggleSetting
        label="System Performance Alerts"
        description="Notify admins of system performance issues"
        checked={settings.notifications?.performanceAlerts || true}
        onChange={(checked) => onSettingChange('notifications', 'performanceAlerts', checked)}
      />

      <DateTimeSetting
        label="Daily Summary Reports"
        description="Send daily summary emails to admins at specified time"
        value={settings.notifications?.dailySummaryTime || '09:00'}
        onChange={(value) => onSettingChange('notifications', 'dailySummaryTime', value)}
        type="time"
      />
    </SettingsCard>

    <SettingsCard
      title="Notification Delivery"
      description="Configure notification delivery settings and rate limits"
      icon={Bell}
    >
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <InputSetting
          label="Email Rate Limit"
          description="Maximum emails per hour per tenant"
          value={settings.notifications?.emailRateLimit || 100}
          onChange={(value) => onSettingChange('notifications', 'emailRateLimit', parseInt(value) || 0)}
          type="number"
          suffix="emails/hour"
        />

        <InputSetting
          label="SMS Rate Limit"
          description="Maximum SMS per hour per tenant"
          value={settings.notifications?.smsRateLimit || 50}
          onChange={(value) => onSettingChange('notifications', 'smsRateLimit', parseInt(value) || 0)}
          type="number"
          suffix="SMS/hour"
        />
      </div>

      <ToggleSetting
        label="Batch Delivery"
        description="Group similar notifications together to reduce volume"
        checked={settings.notifications?.batchDelivery || true}
        onChange={(checked) => onSettingChange('notifications', 'batchDelivery', checked)}
      />

      <ToggleSetting
        label="Delivery Retry"
        description="Retry failed notification deliveries"
        checked={settings.notifications?.deliveryRetry || true}
        onChange={(checked) => onSettingChange('notifications', 'deliveryRetry', checked)}
      />

      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <div className="flex items-start gap-3">
          <Info className="h-5 w-5 text-blue-600 mt-0.5" />
          <div>
            <h4 className="text-sm font-medium text-blue-800">Delivery Status</h4>
            <p className="text-sm text-blue-700 mt-1">
              All notification services are operational. Email delivery: 99.9% success rate. SMS delivery: 98.5% success rate.
            </p>
          </div>
        </div>
      </div>
    </SettingsCard>
  </div>
);

export default NotificationsTab;