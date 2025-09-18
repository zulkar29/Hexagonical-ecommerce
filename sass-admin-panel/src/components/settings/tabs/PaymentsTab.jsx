import React from 'react';
import { CreditCard, Shield, Clock, Webhook } from 'lucide-react';
import SettingsCard from '../SettingsCard';
import ToggleSetting from '../ToggleSetting';
import InputSetting from '../InputSetting';

const PaymentsTab = ({ settings, onSettingChange }) => (
  <div className="space-y-6">
    <SettingsCard
      title="Payment Gateways"
      description="Configure available payment methods for tenants"
      icon={CreditCard}
    >
      <div className="space-y-4">
        <ToggleSetting
          label="SSLCommerz Gateway"
          description="Enable SSLCommerz payment gateway for Bangladesh market"
          checked={settings.payments?.sslcommerzEnabled || true}
          onChange={(checked) => onSettingChange('payments', 'sslcommerzEnabled', checked)}
        />

        {settings.payments?.sslcommerzEnabled && (
          <div className="space-y-4 pl-4 border-l-2 border-muted">
            <InputSetting
              label="SSLCommerz Store ID"
              value={settings.payments?.sslcommerzStoreId || ''}
              onChange={(value) => onSettingChange('payments', 'sslcommerzStoreId', value)}
              placeholder="Your SSLCommerz Store ID"
            />

            <InputSetting
              label="SSLCommerz Store Password"
              value={settings.payments?.sslcommerzPassword || ''}
              onChange={(value) => onSettingChange('payments', 'sslcommerzPassword', value)}
              type="password"
              placeholder="Your SSLCommerz Store Password"
            />
          </div>
        )}

        <ToggleSetting
          label="Stripe Gateway"
          description="Enable Stripe for international payments"
          checked={settings.payments?.stripeEnabled || false}
          onChange={(checked) => onSettingChange('payments', 'stripeEnabled', checked)}
        />

        {settings.payments?.stripeEnabled && (
          <div className="space-y-4 pl-4 border-l-2 border-muted">
            <InputSetting
              label="Stripe Publishable Key"
              value={settings.payments?.stripePublishableKey || ''}
              onChange={(value) => onSettingChange('payments', 'stripePublishableKey', value)}
              placeholder="pk_test_..."
            />

            <InputSetting
              label="Stripe Secret Key"
              value={settings.payments?.stripeSecretKey || ''}
              onChange={(value) => onSettingChange('payments', 'stripeSecretKey', value)}
              type="password"
              placeholder="sk_test_..."
            />
          </div>
        )}
      </div>
    </SettingsCard>

    <SettingsCard
      title="Payment Processing"
      description="Configure payment processing behavior and retry logic"
      icon={Clock}
    >
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <InputSetting
          label="Payment Timeout"
          description="Maximum time to wait for payment confirmation"
          value={settings.payments?.paymentTimeout || 30}
          onChange={(value) => onSettingChange('payments', 'paymentTimeout', parseInt(value) || 0)}
          type="number"
          suffix="seconds"
        />

        <InputSetting
          label="Max Retry Attempts"
          description="Maximum payment retry attempts for failed transactions"
          value={settings.payments?.maxRetryAttempts || 3}
          onChange={(value) => onSettingChange('payments', 'maxRetryAttempts', parseInt(value) || 0)}
          type="number"
          suffix="attempts"
        />
      </div>

      <ToggleSetting
        label="Auto-Retry Failed Payments"
        description="Automatically retry failed subscription payments"
        checked={settings.payments?.autoRetryEnabled || true}
        onChange={(checked) => onSettingChange('payments', 'autoRetryEnabled', checked)}
      />

      <ToggleSetting
        label="Test Mode"
        description="Enable test mode for payment gateways (sandbox environment)"
        checked={settings.payments?.testMode || false}
        onChange={(checked) => onSettingChange('payments', 'testMode', checked)}
      />
    </SettingsCard>

    <SettingsCard
      title="Webhooks & Security"
      description="Configure payment webhooks and security settings"
      icon={Webhook}
    >
      <InputSetting
        label="Webhook Endpoint URL"
        description="URL to receive payment notifications"
        value={settings.payments?.webhookUrl || ''}
        onChange={(value) => onSettingChange('payments', 'webhookUrl', value)}
        placeholder="https://your-domain.com/webhooks/payments"
      />

      <InputSetting
        label="Webhook Secret Key"
        description="Secret key for webhook signature verification"
        value={settings.payments?.webhookSecret || ''}
        onChange={(value) => onSettingChange('payments', 'webhookSecret', value)}
        type="password"
        placeholder="Your webhook secret key"
      />

      <div className="bg-amber-50 border border-amber-200 rounded-lg p-4">
        <div className="flex items-start gap-3">
          <Shield className="h-5 w-5 text-amber-600 mt-0.5" />
          <div>
            <h4 className="text-sm font-medium text-amber-800">Security Recommendations</h4>
            <ul className="text-sm text-amber-700 mt-1 space-y-1">
              <li>• Always use HTTPS for webhook endpoints</li>
              <li>• Verify webhook signatures to ensure authenticity</li>
              <li>• Keep API keys secure and rotate them regularly</li>
              <li>• Use test mode during development and testing</li>
            </ul>
          </div>
        </div>
      </div>
    </SettingsCard>
  </div>
);

export default PaymentsTab;