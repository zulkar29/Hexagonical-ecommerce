import React, { useState, useEffect } from 'react';
import {
  Settings,
  Shield,
  Globe,
  CreditCard,
  Mail,
  Bell,
  Database,
  Key,
  Users,
  Package,
  Server,
  Lock,
  Eye,
  EyeOff,
  Save,
  RefreshCw,
  Download,
  Upload,
  Copy,
  Check,
  AlertTriangle,
  Info,
  Zap,
  Clock,
  Target,
  Gauge,
  Monitor,
  Cloud,
  FileText,
  Search,
  Filter,
  MoreHorizontal,
  Edit,
  Trash2,
  Plus,
  ExternalLink,
  CheckCircle,
  XCircle,
  AlertCircle,
  HelpCircle,
  Wrench,
  Cog,
  ToggleLeft,
  ToggleRight,
  Sliders,
  Activity,
  BarChart3,
  PieChart,
  TrendingUp,
  Calendar,
  MapPin,
  Phone,
  Building,
  Briefcase,
  DollarSign,
  Percent,
  Hash,
  Type,
  Image,
  Link,
  Code,
  Palette,
  MessageSquare
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { Separator } from '@/components/ui/separator';
import { Progress } from '@/components/ui/progress';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

const PlatformSettingsPage = () => {
  const [activeTab, setActiveTab] = useState('general');
  const [isSaving, setIsSaving] = useState(false);
  const [showApiKey, setShowApiKey] = useState(false);
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  // Platform Settings State
  const [settings, setSettings] = useState({
    general: {
      platformName: 'Shop Owner Dashboard',
      platformDescription: 'Complete business management solution for Bangladesh small businesses',
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
      starterPrice: 2000,
      businessPrice: 5000,
      enterprisePrice: 10000,
      trialDuration: 14,
      gracePeriod: 7,
      maxProductsStarter: 500,
      maxProductsBusiness: 2000,
      maxProductsEnterprise: -1, // unlimited
      allowDowngrades: true,
      prorationEnabled: true,
      refundPolicy: 30
    },
    payments: {
      bkashEnabled: true,
      bkashMerchantId: 'BK123456789',
      bkashApiKey: '••••••••••••••••',
      nagadEnabled: true,
      nagadMerchantId: 'NG987654321',
      nagadApiKey: '••••••••••••••••',
      bankTransferEnabled: true,
      paymentRetryAttempts: 3,
      paymentTimeout: 30,
      webhookUrl: 'https://api.shopowner.com/webhooks/payments',
      testMode: false
    },
    notifications: {
      emailNotifications: true,
      smsNotifications: true,
      pushNotifications: true,
      marketingEmails: false,
      systemAlerts: true,
      paymentReminders: true,
      trialExpiryNotice: 7,
      invoiceGeneration: true,
      welcomeEmails: true,
      supportTicketUpdates: true
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
      auditLogging: true
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
  });

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

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
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 2000));
    setIsSaving(false);
    setHasUnsavedChanges(false);
    // Show success message
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
  };

  const settingsTabs = [
    { id: 'general', label: 'General', icon: Settings },
    { id: 'subscription', label: 'Subscription Plans', icon: Package },
    { id: 'payments', label: 'Payment Gateways', icon: CreditCard },
    { id: 'notifications', label: 'Notifications', icon: Bell },
    { id: 'security', label: 'Security', icon: Shield },
    { id: 'features', label: 'Feature Flags', icon: Zap },
    { id: 'integrations', label: 'Integrations', icon: Globe },
  ];

  const SettingCard = ({ title, description, children, icon: Icon, badge }) => (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            {Icon && <Icon className="h-4 w-4 text-muted-foreground" />}
            <CardTitle className="text-base">{title}</CardTitle>
            {badge && <Badge variant="secondary" className="text-xs">{badge}</Badge>}
          </div>
        </div>
        {description && (
          <p className="text-sm text-muted-foreground">{description}</p>
        )}
      </CardHeader>
      <CardContent className="space-y-4">
        {children}
      </CardContent>
    </Card>
  );

  const ToggleSetting = ({ label, description, checked, onChange, disabled = false }) => (
    <div className="flex items-center justify-between p-3 border rounded-lg">
      <div className="flex-1">
        <div className="flex items-center gap-2">
          <Label className="font-medium">{label}</Label>
          {disabled && <Badge variant="outline" className="text-xs">Pro Only</Badge>}
        </div>
        {description && (
          <p className="text-sm text-muted-foreground mt-1">{description}</p>
        )}
      </div>
      <Switch 
        checked={checked} 
        onCheckedChange={onChange}
        disabled={disabled}
      />
    </div>
  );

  const InputSetting = ({ label, description, value, onChange, type = "text", placeholder, suffix, prefix }) => (
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

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <Settings className="h-5 w-5 text-primary" />
              <h1 className="text-xl font-semibold">Platform Settings</h1>
            </div>
            {hasUnsavedChanges && (
              <Badge variant="outline" className="text-xs text-orange-600">
                <AlertTriangle className="h-3 w-3 mr-1" />
                Unsaved Changes
              </Badge>
            )}
          </div>
          
          <div className="flex items-center gap-2">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search settings..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 w-64"
              />
            </div>
            
            <Button 
              variant="outline" 
              size="sm"
              onClick={handleSave}
              disabled={!hasUnsavedChanges || isSaving}
            >
              {isSaving ? (
                <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
              ) : (
                <Save className="h-4 w-4 mr-2" />
              )}
              {isSaving ? 'Saving...' : 'Save Changes'}
            </Button>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Platform Actions</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <Download className="h-4 w-4 mr-2" />
                  Export Settings
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Upload className="h-4 w-4 mr-2" />
                  Import Settings
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Reset to Defaults
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      {/* Settings Content */}
      <div className="flex-1 overflow-hidden">
        <Tabs value={activeTab} onValueChange={setActiveTab} className="h-full flex">
          {/* Sidebar Navigation */}
          <div className="w-64 border-r bg-muted/30">
            <TabsList className="flex flex-col h-full w-full bg-transparent p-2 space-y-1">
              {settingsTabs.map((tab) => (
                <TabsTrigger
                  key={tab.id}
                  value={tab.id}
                  className="w-full justify-start gap-2 data-[state=active]:bg-background data-[state=active]:shadow-sm"
                >
                  <tab.icon className="h-4 w-4" />
                  {tab.label}
                </TabsTrigger>
              ))}
            </TabsList>
          </div>

          {/* Settings Content */}
          <div className="flex-1 overflow-auto">
            {/* General Settings */}
            <TabsContent value="general" className="h-full p-6 space-y-6">
              <div className="space-y-6">
                <SettingCard 
                  title="Email Notifications" 
                  description="Configure email notification settings"
                  icon={Mail}
                >
                  <ToggleSetting
                    label="Email Notifications"
                    description="Enable email notifications system-wide"
                    checked={settings.notifications.emailNotifications}
                    onChange={(checked) => handleSettingChange('notifications', 'emailNotifications', checked)}
                  />

                  <ToggleSetting
                    label="Welcome Emails"
                    description="Send welcome emails to new users"
                    checked={settings.notifications.welcomeEmails}
                    onChange={(checked) => handleSettingChange('notifications', 'welcomeEmails', checked)}
                  />

                  <ToggleSetting
                    label="Invoice Generation"
                    description="Automatically send invoices via email"
                    checked={settings.notifications.invoiceGeneration}
                    onChange={(checked) => handleSettingChange('notifications', 'invoiceGeneration', checked)}
                  />

                  <ToggleSetting
                    label="Payment Reminders"
                    description="Send payment reminder emails for overdue accounts"
                    checked={settings.notifications.paymentReminders}
                    onChange={(checked) => handleSettingChange('notifications', 'paymentReminders', checked)}
                  />

                  <InputSetting
                    label="Trial Expiry Notice"
                    description="Days before trial expires to send notification"
                    value={settings.notifications.trialExpiryNotice}
                    onChange={(value) => handleSettingChange('notifications', 'trialExpiryNotice', parseInt(value) || 0)}
                    type="number"
                    suffix="days before"
                  />
                </SettingCard>

                <SettingCard 
                  title="SMS & Push Notifications" 
                  description="Configure mobile and SMS notification settings"
                  icon={Phone}
                >
                  <ToggleSetting
                    label="SMS Notifications"
                    description="Enable SMS notifications for critical updates"
                    checked={settings.notifications.smsNotifications}
                    onChange={(checked) => handleSettingChange('notifications', 'smsNotifications', checked)}
                  />

                  <ToggleSetting
                    label="Push Notifications"
                    description="Enable mobile app push notifications"
                    checked={settings.notifications.pushNotifications}
                    onChange={(checked) => handleSettingChange('notifications', 'pushNotifications', checked)}
                  />

                  <ToggleSetting
                    label="System Alerts"
                    description="Send alerts for system maintenance and updates"
                    checked={settings.notifications.systemAlerts}
                    onChange={(checked) => handleSettingChange('notifications', 'systemAlerts', checked)}
                  />

                  <ToggleSetting
                    label="Support Ticket Updates"
                    description="Notify users about support ticket status changes"
                    checked={settings.notifications.supportTicketUpdates}
                    onChange={(checked) => handleSettingChange('notifications', 'supportTicketUpdates', checked)}
                  />
                </SettingCard>

                <SettingCard 
                  title="Marketing Communications" 
                  description="Configure marketing and promotional notifications"
                  icon={Target}
                >
                  <ToggleSetting
                    label="Marketing Emails"
                    description="Send promotional and feature update emails"
                    checked={settings.notifications.marketingEmails}
                    onChange={(checked) => handleSettingChange('notifications', 'marketingEmails', checked)}
                  />

                  <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                    <div className="flex items-start gap-3">
                      <Info className="h-5 w-5 text-yellow-600 mt-0.5" />
                      <div>
                        <h4 className="text-sm font-medium text-yellow-800">Marketing Compliance</h4>
                        <p className="text-sm text-yellow-700 mt-1">
                          Users must explicitly opt-in to marketing communications. This setting only affects users who have already consented.
                        </p>
                      </div>
                    </div>
                  </div>
                </SettingCard>
              </div>
            </TabsContent>

            {/* Security Settings */}
            <TabsContent value="security" className="h-full p-6 space-y-6">
              <div className="space-y-6">
                <SettingCard 
                  title="Password Policy" 
                  description="Configure password requirements for all users"
                  icon={Lock}
                >
                  <InputSetting
                    label="Minimum Password Length"
                    description="Minimum number of characters required"
                    value={settings.security.passwordMinLength}
                    onChange={(value) => handleSettingChange('security', 'passwordMinLength', parseInt(value) || 0)}
                    type="number"
                    suffix="characters"
                  />

                  <div className="space-y-3">
                    <ToggleSetting
                      label="Require Uppercase Letters"
                      description="Password must contain at least one uppercase letter"
                      checked={settings.security.requireUppercase}
                      onChange={(checked) => handleSettingChange('security', 'requireUppercase', checked)}
                    />

                    <ToggleSetting
                      label="Require Numbers"
                      description="Password must contain at least one number"
                      checked={settings.security.requireNumbers}
                      onChange={(checked) => handleSettingChange('security', 'requireNumbers', checked)}
                    />

                    <ToggleSetting
                      label="Require Special Characters"
                      description="Password must contain at least one special character"
                      checked={settings.security.requireSpecialChars}
                      onChange={(checked) => handleSettingChange('security', 'requireSpecialChars', checked)}
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
                </SettingCard>

                <SettingCard 
                  title="Session Management" 
                  description="Configure user session and login security"
                  icon={Clock}
                >
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <InputSetting
                      label="Session Timeout"
                      description="Hours before automatic logout"
                      value={settings.security.sessionTimeout}
                      onChange={(value) => handleSettingChange('security', 'sessionTimeout', parseInt(value) || 0)}
                      type="number"
                      suffix="hours"
                    />

                    <InputSetting
                      label="Max Login Attempts"
                      description="Failed attempts before account lockout"
                      value={settings.security.maxLoginAttempts}
                      onChange={(value) => handleSettingChange('security', 'maxLoginAttempts', parseInt(value) || 0)}
                      type="number"
                      suffix="attempts"
                    />

                    <InputSetting
                      label="Lockout Duration"
                      description="Minutes to lock account after max attempts"
                      value={settings.security.lockoutDuration}
                      onChange={(value) => handleSettingChange('security', 'lockoutDuration', parseInt(value) || 0)}
                      type="number"
                      suffix="minutes"
                    />

                    <InputSetting
                      label="API Rate Limit"
                      description="API requests per hour per user"
                      value={settings.security.apiRateLimit}
                      onChange={(value) => handleSettingChange('security', 'apiRateLimit', parseInt(value) || 0)}
                      type="number"
                      suffix="req/hour"
                    />
                  </div>

                  <ToggleSetting
                    label="Two-Factor Authentication Required"
                    description="Require 2FA for all admin users"
                    checked={settings.security.twoFactorRequired}
                    onChange={(checked) => handleSettingChange('security', 'twoFactorRequired', checked)}
                  />

                  <ToggleSetting
                    label="IP Whitelisting"
                    description="Restrict admin access to specific IP addresses"
                    checked={settings.security.ipWhitelisting}
                    onChange={(checked) => handleSettingChange('security', 'ipWhitelisting', checked)}
                    disabled={true}
                  />
                </SettingCard>

                <SettingCard 
                  title="Data Protection" 
                  description="Configure data retention and audit settings"
                  icon={Database}
                >
                  <InputSetting
                    label="Data Retention Period"
                    description="Days to retain user data after account deletion"
                    value={settings.security.dataRetention}
                    onChange={(value) => handleSettingChange('security', 'dataRetention', parseInt(value) || 0)}
                    type="number"
                    suffix="days"
                  />

                  <ToggleSetting
                    label="Audit Logging"
                    description="Log all admin actions for compliance"
                    checked={settings.security.auditLogging}
                    onChange={(checked) => handleSettingChange('security', 'auditLogging', checked)}
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
                </SettingCard>
              </div>
            </TabsContent>

            {/* Feature Flags */}
            <TabsContent value="features" className="h-full p-6 space-y-6">
              <div className="space-y-6">
                <SettingCard 
                  title="Core Features" 
                  description="Enable or disable core platform features"
                  icon={Zap}
                >
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <ToggleSetting
                      label="Multi-Tenant Mode"
                      description="Enable separate data isolation for each tenant"
                      checked={settings.features.multiTenantMode}
                      onChange={(checked) => handleSettingChange('features', 'multiTenantMode', checked)}
                    />

                    <ToggleSetting
                      label="Advanced Analytics"
                      description="Enable detailed analytics and reporting"
                      checked={settings.features.advancedAnalytics}
                      onChange={(checked) => handleSettingChange('features', 'advancedAnalytics', checked)}
                    />

                    <ToggleSetting
                      label="Custom Branding"
                      description="Allow tenants to customize their branding"
                      checked={settings.features.customBranding}
                      onChange={(checked) => handleSettingChange('features', 'customBranding', checked)}
                    />

                    <ToggleSetting
                      label="API Access"
                      description="Provide REST API access to tenants"
                      checked={settings.features.apiAccess}
                      onChange={(checked) => handleSettingChange('features', 'apiAccess', checked)}
                    />

                    <ToggleSetting
                      label="Webhook Support"
                      description="Enable webhook notifications for events"
                      checked={settings.features.webhookSupport}
                      onChange={(checked) => handleSettingChange('features', 'webhookSupport', checked)}
                    />

                    <ToggleSetting
                      label="Export Functionality"
                      description="Allow data export in various formats"
                      checked={settings.features.exportFunctionality}
                      onChange={(checked) => handleSettingChange('features', 'exportFunctionality', checked)}
                    />
                  </div>
                </SettingCard>

                <SettingCard 
                  title="Support Features" 
                  description="Configure customer support and help features"
                  icon={HelpCircle}
                >
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <ToggleSetting
                      label="Support Ticketing"
                      description="Built-in support ticket system"
                      checked={settings.features.supportTicketing}
                      onChange={(checked) => handleSettingChange('features', 'supportTicketing', checked)}
                    />

                    <ToggleSetting
                      label="Knowledge Base"
                      description="Self-service knowledge base articles"
                      checked={settings.features.knowledgeBase}
                      onChange={(checked) => handleSettingChange('features', 'knowledgeBase', checked)}
                    />

                    <ToggleSetting
                      label="Live Chat"
                      description="Real-time chat support for users"
                      checked={settings.features.liveChat}
                      onChange={(checked) => handleSettingChange('features', 'liveChat', checked)}
                      disabled={true}
                    />

                    <ToggleSetting
                      label="Backup System"
                      description="Automated daily backups"
                      checked={settings.features.backupSystem}
                      onChange={(checked) => handleSettingChange('features', 'backupSystem', checked)}
                    />
                  </div>
                </SettingCard>

                <SettingCard 
                  title="Mobile & Offline Features" 
                  description="Configure mobile app and offline capabilities"
                  icon={Monitor}
                  badge="Beta"
                >
                  <ToggleSetting
                    label="Mobile App"
                    description="Enable mobile app access for tenants"
                    checked={settings.features.mobileApp}
                    onChange={(checked) => handleSettingChange('features', 'mobileApp', checked)}
                  />

                  <ToggleSetting
                    label="Offline Mode"
                    description="Allow limited functionality when offline"
                    checked={settings.features.offlineMode}
                    onChange={(checked) => handleSettingChange('features', 'offlineMode', checked)}
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
                </SettingCard>
              </div>
            </TabsContent>

            {/* Integrations */}
            <TabsContent value="integrations" className="h-full p-6 space-y-6">
              <div className="space-y-6">
                <SettingCard 
                  title="Google Analytics" 
                  description="Configure Google Analytics tracking"
                  icon={BarChart3}
                  badge={settings.integrations.googleAnalytics.enabled ? "Connected" : "Disconnected"}
                >
                  <ToggleSetting
                    label="Enable Google Analytics"
                    description="Track user behavior and platform usage"
                    checked={settings.integrations.googleAnalytics.enabled}
                    onChange={(checked) => handleSettingChange('integrations', 'googleAnalytics.enabled', checked)}
                  />

                  {settings.integrations.googleAnalytics.enabled && (
                    <div className="space-y-4 pl-4 border-l-2 border-muted">
                      <InputSetting
                        label="Tracking ID"
                        value={settings.integrations.googleAnalytics.trackingId}
                        onChange={(value) => handleSettingChange('integrations', 'googleAnalytics.trackingId', value)}
                        placeholder="GA-XXXX-XXXX"
                      />

                      <ToggleSetting
                        label="Enhanced Ecommerce Tracking"
                        description="Track subscription and payment events"
                        checked={settings.integrations.googleAnalytics.ecommerceTracking}
                        onChange={(checked) => handleSettingChange('integrations', 'googleAnalytics.ecommerceTracking', checked)}
                      />
                    </div>
                  )}
                </SettingCard>

                <SettingCard 
                  title="WhatsApp Business" 
                  description="Configure WhatsApp Business API integration"
                  icon={MessageSquare}
                  badge={settings.integrations.whatsapp.enabled ? "Connected" : "Disconnected"}
                >
                  <ToggleSetting
                    label="Enable WhatsApp Integration"
                    description="Send notifications via WhatsApp Business API"
                    checked={settings.integrations.whatsapp.enabled}
                    onChange={(checked) => handleSettingChange('integrations', 'whatsapp.enabled', checked)}
                  />

                  {settings.integrations.whatsapp.enabled && (
                    <div className="space-y-4 pl-4 border-l-2 border-muted">
                      <InputSetting
                        label="Business Phone Number"
                        value={settings.integrations.whatsapp.businessNumber}
                        onChange={(value) => handleSettingChange('integrations', 'whatsapp.businessNumber', value)}
                        placeholder="+8801XXXXXXXXX"
                      />

                      <div className="space-y-2">
                        <Label className="font-medium">API Token</Label>
                        <div className="flex gap-2">
                          <Input
                            type="password"
                            value={settings.integrations.whatsapp.apiToken}
                            onChange={(e) => handleSettingChange('integrations', 'whatsapp.apiToken', e.target.value)}
                            placeholder="Enter WhatsApp API token"
                          />
                          <Button variant="outline" size="sm">
                            <Eye className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    </div>
                  )}
                </SettingCard>

                <SettingCard 
                  title="SMS Provider" 
                  description="Configure SMS service provider settings"
                  icon={Phone}
                >
                  <div className="space-y-4">
                    <div className="space-y-2">
                      <Label className="font-medium">SMS Provider</Label>
                      <Select value={settings.integrations.sms.provider} onValueChange={(value) => 
                        handleSettingChange('integrations', 'sms.provider', value)
                      }>
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
                      onChange={(value) => handleSettingChange('integrations', 'sms.senderId', value)}
                      placeholder="SHOPOWNER"
                    />

                    <div className="space-y-2">
                      <Label className="font-medium">API Key</Label>
                      <div className="flex gap-2">
                        <Input
                          type="password"
                          value={settings.integrations.sms.apiKey}
                          onChange={(e) => handleSettingChange('integrations', 'sms.apiKey', e.target.value)}
                          placeholder="Enter SMS provider API key"
                        />
                        <Button variant="outline" size="sm">
                          <Eye className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  </div>
                </SettingCard>

                <SettingCard 
                  title="Facebook Integration" 
                  description="Configure Facebook app and pixel integration"
                  icon={Globe}
                  badge={settings.integrations.facebook.enabled ? "Connected" : "Disconnected"}
                >
                  <ToggleSetting
                    label="Enable Facebook Integration"
                    description="Connect with Facebook for social login and marketing"
                    checked={settings.integrations.facebook.enabled}
                    onChange={(checked) => handleSettingChange('integrations', 'facebook.enabled', checked)}
                  />

                  {settings.integrations.facebook.enabled && (
                    <div className="space-y-4 pl-4 border-l-2 border-muted">
                      <InputSetting
                        label="Facebook App ID"
                        value={settings.integrations.facebook.appId}
                        onChange={(value) => handleSettingChange('integrations', 'facebook.appId', value)}
                        placeholder="Enter Facebook App ID"
                      />

                      <InputSetting
                        label="Facebook Pixel ID"
                        value={settings.integrations.facebook.pixelId}
                        onChange={(value) => handleSettingChange('integrations', 'facebook.pixelId', value)}
                        placeholder="Enter Facebook Pixel ID"
                      />
                    </div>
                  )}
                </SettingCard>
              </div>
            </TabsContent>
          </div>
        </Tabs>
      </div>

      {/* Save Confirmation */}
      {hasUnsavedChanges && (
        <div className="border-t bg-yellow-50 border-yellow-200 p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <AlertTriangle className="h-5 w-5 text-yellow-600" />
              <div>
                <p className="text-sm font-medium text-yellow-800">You have unsaved changes</p>
                <p className="text-xs text-yellow-700">Make sure to save your changes before leaving this page.</p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Button 
                variant="outline" 
                size="sm"
                onClick={() => {
                  setHasUnsavedChanges(false);
                  // Reset to original values
                }}
              >
                Discard Changes
              </Button>
              <Button 
                size="sm"
                onClick={handleSave}
                disabled={isSaving}
              >
                {isSaving ? (
                  <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                ) : (
                  <Save className="h-4 w-4 mr-2" />
                )}
                {isSaving ? 'Saving...' : 'Save Changes'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PlatformSettingsPage;
          