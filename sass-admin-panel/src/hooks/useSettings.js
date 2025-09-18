import { useState } from 'react';

const defaultSettings = {
  general: {
    platformName: 'Hexagonal Ecommerce Platform',
    platformDescription: 'Complete multi-tenant ecommerce solution for Bangladesh businesses',
    defaultLanguage: 'bn',
    supportedLanguages: ['bn', 'en'],
    timezone: 'Asia/Dhaka',
    currency: 'BDT',
    dateFormat: 'DD/MM/YYYY',
    maintenanceMode: false,
    debugMode: false,
    analyticsEnabled: true
  },
  subscription: {
    starterPrice: 1990,
    professionalPrice: 4990,
    enterprisePrice: 12990,
    trialDuration: 14,
    gracePeriod: 7,
    maxProductsStarter: 500,
    maxProductsProfessional: 2000,
    maxProductsEnterprise: -1, // unlimited
    allowDowngrades: true,
    prorationEnabled: true,
    refundPolicy: 30
  },
  payments: {
    sslcommerzEnabled: true,
    sslcommerzStoreId: '',
    sslcommerzPassword: '',
    stripeEnabled: false,
    stripePublishableKey: '',
    stripeSecretKey: '',
    paymentTimeout: 30,
    maxRetryAttempts: 3,
    autoRetryEnabled: true,
    testMode: false,
    webhookUrl: '',
    webhookSecret: ''
  },
  notifications: {
    emailNotifications: true,
    smsNotifications: true,
    pushNotifications: true,
    systemAlerts: true,
    securityAlerts: true,
    paymentReminders: true,
    subscriptionEmails: true,
    welcomeEmails: true,
    newTenantAlerts: true,
    paymentFailureAlerts: true,
    supportTicketAlerts: true,
    performanceAlerts: true,
    trialExpiryNotice: 7,
    dailySummaryTime: '09:00',
    emailRateLimit: 100,
    smsRateLimit: 50,
    batchDelivery: true,
    deliveryRetry: true
  },
  security: {
    passwordMinLength: 8,
    requireUppercase: true,
    requireNumbers: true,
    requireSpecialChars: true,
    sessionTimeout: 24,
    maxLoginAttempts: 5,
    lockoutDuration: 30,
    twoFactorRequired: false,
    ipWhitelisting: false,
    apiRateLimit: 1000,
    dataRetention: 365,
    auditLogging: true,
    automatedBackups: true,
    backupTime: '02:00',
    backupRetention: 30,
    offsiteBackup: false
  },
  features: {
    multiTenantMode: true,
    advancedAnalytics: true,
    customBranding: true,
    apiAccess: true,
    webhookSupport: true,
    exportFunctionality: true,
    backupSystem: true,
    supportTicketing: true,
    knowledgeBase: true,
    liveChat: false,
    mobileApp: true,
    offlineMode: false
  },
  integrations: {
    googleAnalytics: {
      enabled: true,
      trackingId: 'GA-XXXX-XXXX',
      ecommerceTracking: true
    },
    facebook: {
      enabled: false,
      appId: '',
      pixelId: ''
    },
    whatsapp: {
      enabled: true,
      businessNumber: '+8801XXXXXXXXX',
      apiToken: '••••••••••••••••'
    },
    sms: {
      provider: 'ssl_wireless',
      apiKey: '••••••••••••••••',
      senderId: 'SHOPOWNER'
    }
  }
};

export const useSettings = () => {
  const [settings, setSettings] = useState(defaultSettings);
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  const handleSettingChange = (category, key, value) => {
    setSettings(prev => ({
      ...prev,
      [category]: {
        ...prev[category],
        [key]: key.includes('.') ?
          { ...prev[category][key.split('.')[0]], [key.split('.')[1]]: value } :
          value
      }
    }));
    setHasUnsavedChanges(true);
  };

  const handleSave = async () => {
    setIsSaving(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000));
      setHasUnsavedChanges(false);
      // Show success message
    } catch (error) {
      console.error('Failed to save settings:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const resetChanges = () => {
    setSettings(defaultSettings);
    setHasUnsavedChanges(false);
  };

  return {
    settings,
    hasUnsavedChanges,
    isSaving,
    handleSettingChange,
    handleSave,
    resetChanges
  };
};