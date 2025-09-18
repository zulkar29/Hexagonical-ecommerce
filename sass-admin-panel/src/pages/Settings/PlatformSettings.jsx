import React, { useState, useMemo } from 'react';
import {
  Settings,
  Package,
  CreditCard,
  Bell,
  Shield,
  Zap,
  Globe,
  Search,
  X,
  Save,
  Download,
  Upload,
  RotateCcw
} from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';

// Components
import SettingsHeader from '@/components/settings/SettingsHeader';
import UnsavedChangesBar from '@/components/settings/UnsavedChangesBar';
import GeneralTab from '@/components/settings/tabs/GeneralTab';
import SecurityTab from '@/components/settings/tabs/SecurityTab';
import FeaturesTab from '@/components/settings/tabs/FeaturesTab';
import IntegrationsTab from '@/components/settings/tabs/IntegrationsTab';
import SubscriptionTab from '@/components/settings/tabs/SubscriptionTab';
import PaymentsTab from '@/components/settings/tabs/PaymentsTab';
import NotificationsTab from '@/components/settings/tabs/NotificationsTab';

// Hooks
import { useSettings } from '@/hooks/useSettings';

const PlatformSettingsPage = () => {
  const [activeTab, setActiveTab] = useState('general');
  const [searchQuery, setSearchQuery] = useState('');

  const {
    settings,
    hasUnsavedChanges,
    isSaving,
    handleSettingChange,
    handleSave,
    resetChanges
  } = useSettings();

  const settingsTabs = [
    {
      id: 'general',
      label: 'General',
      icon: Settings,
      component: GeneralTab,
      description: 'Platform name, language, and basic configuration'
    },
    {
      id: 'subscription',
      label: 'Subscriptions',
      icon: Package,
      component: SubscriptionTab,
      description: 'Pricing plans, limits, and billing settings'
    },
    {
      id: 'payments',
      label: 'Payments',
      icon: CreditCard,
      component: PaymentsTab,
      description: 'Payment gateways and processing settings'
    },
    {
      id: 'notifications',
      label: 'Notifications',
      icon: Bell,
      component: NotificationsTab,
      description: 'Email, SMS, and alert configurations'
    },
    {
      id: 'security',
      label: 'Security',
      icon: Shield,
      component: SecurityTab,
      description: 'Password policies, sessions, and data protection'
    },
    {
      id: 'features',
      label: 'Features',
      icon: Zap,
      component: FeaturesTab,
      description: 'Feature flags and platform capabilities'
    },
    {
      id: 'integrations',
      label: 'Integrations',
      icon: Globe,
      component: IntegrationsTab,
      description: 'Third-party services and API connections'
    },
  ];

  // Filter tabs based on search query
  const filteredTabs = useMemo(() => {
    if (!searchQuery.trim()) return settingsTabs;

    return settingsTabs.filter(tab =>
      tab.label.toLowerCase().includes(searchQuery.toLowerCase()) ||
      tab.description.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [searchQuery, settingsTabs]);

  const handleExport = () => {
    const dataStr = JSON.stringify(settings, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = 'platform-settings.json';
    link.click();
  };

  const handleImport = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json';
    input.onchange = (e) => {
      const file = e.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          try {
            const importedSettings = JSON.parse(e.target.result);
            // TODO: Validate and merge settings
            console.log('Import settings:', importedSettings);
          } catch (error) {
            console.error('Invalid JSON file');
          }
        };
        reader.readAsText(file);
      }
    };
    input.click();
  };

  const handleReset = () => {
    if (confirm('Are you sure you want to reset all settings to defaults?')) {
      resetChanges();
    }
  };

  const renderTabContent = () => {
    const activeTabData = settingsTabs.find(tab => tab.id === activeTab);
    if (!activeTabData) return null;
    
    const TabComponent = activeTabData.component;
    return (
      <div className="space-y-6">
        <div>
          <h2 className="text-2xl font-bold text-foreground">{activeTabData.label}</h2>
          <p className="text-muted-foreground mt-1">{activeTabData.description}</p>
        </div>
        <TabComponent
          settings={settings}
          onSettingChange={handleSettingChange}
        />
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="border-b bg-card/50 backdrop-blur-sm sticky top-0 z-10">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Settings className="h-6 w-6 text-primary" />
              <div>
                <h1 className="text-2xl font-bold text-foreground">Platform Settings</h1>
                <p className="text-sm text-muted-foreground">
                  Configure your platform's behavior and integrations
                </p>
              </div>
              {hasUnsavedChanges && (
                <Badge variant="secondary" className="ml-2 animate-pulse">
                  Unsaved Changes
                </Badge>
              )}
            </div>
            
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={handleExport}
                className="hidden sm:flex"
              >
                <Download className="h-4 w-4 mr-2" />
                Export
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={handleImport}
                className="hidden sm:flex"
              >
                <Upload className="h-4 w-4 mr-2" />
                Import
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={handleReset}
                className="hidden sm:flex"
              >
                <RotateCcw className="h-4 w-4 mr-2" />
                Reset
              </Button>
              <Button
                onClick={handleSave}
                disabled={!hasUnsavedChanges}
                size="sm"
              >
                <Save className="h-4 w-4 mr-2" />
                Save Changes
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Sidebar Navigation */}
          <div className="lg:col-span-1">
            <Card>
              <CardHeader className="pb-3">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    placeholder="Search settings..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </CardHeader>
              <CardContent className="p-0">
                <nav className="space-y-1">
                  {filteredTabs.map((tab, index) => (
                    <div key={tab.id}>
                      <button
                        onClick={() => setActiveTab(tab.id)}
                        className={`
                          w-full flex items-center gap-3 px-4 py-3 text-left transition-colors
                          hover:bg-muted/50
                          ${activeTab === tab.id 
                            ? 'bg-primary/10 text-primary border-r-2 border-primary' 
                            : 'text-muted-foreground hover:text-foreground'
                          }
                        `}
                      >
                        <tab.icon className="h-4 w-4 flex-shrink-0" />
                        <div className="flex-1">
                          <div className="font-medium text-sm">{tab.label}</div>
                          <div className="text-xs text-muted-foreground line-clamp-1">
                            {tab.description}
                          </div>
                        </div>
                      </button>
                      {index < filteredTabs.length - 1 && <Separator />}
                    </div>
                  ))}
                </nav>
                
                <div className="p-4 border-t bg-muted/20">
                  <div className="flex justify-between items-center text-xs text-muted-foreground">
                    <span>Total Settings</span>
                    <Badge variant="outline" className="text-xs">
                      {settingsTabs.length}
                    </Badge>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Content Area */}
          <div className="lg:col-span-3">
            <Card>
              <CardContent className="p-6">
                {filteredTabs.length > 0 ? (
                  renderTabContent()
                ) : (
                  <div className="flex flex-col items-center justify-center py-12 text-center">
                    <Search className="h-12 w-12 text-muted-foreground mb-4" />
                    <h3 className="text-lg font-medium text-foreground mb-2">No settings found</h3>
                    <p className="text-muted-foreground">
                      Try adjusting your search query or browse all settings.
                    </p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      {/* Unsaved Changes Bar */}
      {hasUnsavedChanges && (
        <div className="fixed bottom-0 left-0 right-0 bg-card border-t p-4 z-20">
          <div className="container mx-auto flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Badge variant="secondary">Unsaved Changes</Badge>
              <span className="text-sm text-muted-foreground">
                You have unsaved changes that will be lost if you leave.
              </span>
            </div>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={resetChanges}>
                Discard
              </Button>
              <Button size="sm" onClick={handleSave}>
                Save Changes
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PlatformSettingsPage;