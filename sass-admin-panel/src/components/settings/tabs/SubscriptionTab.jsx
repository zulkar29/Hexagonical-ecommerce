import React from 'react';
import { Package, DollarSign, Clock, CheckCircle } from 'lucide-react';
import SettingsCard from '../SettingsCard';
import ToggleSetting from '../ToggleSetting';
import InputSetting from '../InputSetting';

const SubscriptionTab = ({ settings, onSettingChange }) => {
  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  return (
    <div className="space-y-6">
      <SettingsCard
        title="Subscription Plans"
        description="Configure pricing and features for different subscription tiers"
        icon={Package}
      >
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <InputSetting
            label="Starter Plan Price"
            description="Monthly price for Starter plan"
            value={settings.subscription?.starterPrice || 1990}
            onChange={(value) => onSettingChange('subscription', 'starterPrice', parseInt(value) || 0)}
            type="number"
            prefix="৳"
            suffix="/month"
          />

          <InputSetting
            label="Professional Plan Price"
            description="Monthly price for Professional plan"
            value={settings.subscription?.professionalPrice || 4990}
            onChange={(value) => onSettingChange('subscription', 'professionalPrice', parseInt(value) || 0)}
            type="number"
            prefix="৳"
            suffix="/month"
          />

          <InputSetting
            label="Enterprise Plan Price"
            description="Monthly price for Enterprise plan"
            value={settings.subscription?.enterprisePrice || 12990}
            onChange={(value) => onSettingChange('subscription', 'enterprisePrice', parseInt(value) || 0)}
            type="number"
            prefix="৳"
            suffix="/month"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <InputSetting
            label="Trial Duration"
            description="Free trial period for new users"
            value={settings.subscription?.trialDuration || 14}
            onChange={(value) => onSettingChange('subscription', 'trialDuration', parseInt(value) || 0)}
            type="number"
            suffix="days"
          />

          <InputSetting
            label="Grace Period"
            description="Grace period after payment failure"
            value={settings.subscription?.gracePeriod || 7}
            onChange={(value) => onSettingChange('subscription', 'gracePeriod', parseInt(value) || 0)}
            type="number"
            suffix="days"
          />
        </div>
      </SettingsCard>

      <SettingsCard
        title="Plan Limits & Features"
        description="Configure resource limits for each subscription tier"
        icon={CheckCircle}
      >
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <InputSetting
            label="Starter - Max Products"
            description="Maximum products for Starter plan"
            value={settings.subscription?.maxProductsStarter || 500}
            onChange={(value) => onSettingChange('subscription', 'maxProductsStarter', parseInt(value) || 0)}
            type="number"
            suffix="products"
          />

          <InputSetting
            label="Professional - Max Products"
            description="Maximum products for Professional plan"
            value={settings.subscription?.maxProductsProfessional || 2000}
            onChange={(value) => onSettingChange('subscription', 'maxProductsProfessional', parseInt(value) || 0)}
            type="number"
            suffix="products"
          />

          <div className="space-y-2">
            <label className="text-sm font-medium">Enterprise - Products</label>
            <div className="flex items-center p-3 border rounded-lg bg-green-50 border-green-200">
              <CheckCircle className="h-4 w-4 text-green-600 mr-2" />
              <span className="text-sm text-green-800">Unlimited</span>
            </div>
          </div>
        </div>
      </SettingsCard>

      <SettingsCard
        title="Billing & Payment Settings"
        description="Configure billing behavior and payment policies"
        icon={DollarSign}
      >
        <div className="space-y-4">
          <ToggleSetting
            label="Allow Plan Downgrades"
            description="Allow customers to downgrade their subscription plan"
            checked={settings.subscription?.allowDowngrades || true}
            onChange={(checked) => onSettingChange('subscription', 'allowDowngrades', checked)}
          />

          <ToggleSetting
            label="Proration Enabled"
            description="Prorate charges when customers upgrade/downgrade mid-cycle"
            checked={settings.subscription?.prorationEnabled || true}
            onChange={(checked) => onSettingChange('subscription', 'prorationEnabled', checked)}
          />

          <InputSetting
            label="Refund Policy Period"
            description="Days within which customers can request refunds"
            value={settings.subscription?.refundPolicy || 30}
            onChange={(value) => onSettingChange('subscription', 'refundPolicy', parseInt(value) || 0)}
            type="number"
            suffix="days"
          />
        </div>

        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div className="flex items-start gap-3">
            <DollarSign className="h-5 w-5 text-blue-600 mt-0.5" />
            <div>
              <h4 className="text-sm font-medium text-blue-800">Current Pricing Summary</h4>
              <div className="text-sm text-blue-700 mt-1 space-y-1">
                <p>Starter: {formatCurrency(settings.subscription?.starterPrice || 1990)}/month</p>
                <p>Professional: {formatCurrency(settings.subscription?.professionalPrice || 4990)}/month</p>
                <p>Enterprise: {formatCurrency(settings.subscription?.enterprisePrice || 12990)}/month</p>
              </div>
            </div>
          </div>
        </div>
      </SettingsCard>
    </div>
  );
};

export default SubscriptionTab;