import React, { useState } from 'react';
import {
  Settings,
  Store,
  User,
  Bell,
  Shield,
  CreditCard,
  Printer,
  Smartphone,
  Globe,
  Database,
  Upload,
  Download,
  Trash2,
  Save,
  Eye,
  EyeOff,
  Camera,
  MapPin,
  Phone,
  Mail,
  Clock,
  DollarSign,
  Percent,
  Users,
  Package,
  BarChart3,
  RefreshCw,
  CheckCircle,
  AlertTriangle,
  Info,
  Plus,
  X,
  Edit,
  Copy,
  Key,
  Wifi,
  Monitor,
  Sun,
  Moon,
  Volume2,
  VolumeX,
  Languages,
  Calendar,
  Building,
  FileText,
  Lock,
  Unlock,
  HelpCircle
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import { Progress } from '@/components/ui/progress';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs';
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
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  Alert,
  AlertDescription,
  AlertTitle,
} from '@/components/ui/alert';

const SettingsPage = () => {
  const [showPasswordDialog, setShowPasswordDialog] = useState(false);
  const [showBackupDialog, setShowBackupDialog] = useState(false);
  const [settings, setSettings] = useState({
    // Business Info
    businessName: 'Rahman General Store',
    businessAddress: 'House 45, Road 12, Dhanmondi, Dhaka-1205',
    businessPhone: '01712345678',
    businessEmail: 'rahman.shop@gmail.com',
    businessLicense: 'TIN-123456789012',
    
    // User Profile
    ownerName: 'Mohammad Rahman',
    ownerEmail: 'm.rahman@gmail.com',
    ownerPhone: '01712345678',
    
    // System Settings
    currency: 'BDT',
    language: 'en',
    timezone: 'Asia/Dhaka',
    dateFormat: 'DD/MM/YYYY',
    theme: 'light',
    
    // Notifications
    emailNotifications: true,
    smsNotifications: true,
    lowStockAlerts: true,
    dailyReports: true,
    
    // POS Settings
    taxRate: 0,
    defaultPaymentMethod: 'cash',
    receiptTemplate: 'standard',
    autoPrint: false,
    
    // Security
    twoFactorAuth: false,
    sessionTimeout: 30,
    passwordExpiry: 90,
    
    // Backup
    autoBackup: true,
    backupFrequency: 'daily',
    lastBackup: '2025-01-23 10:30 AM'
  });

  const [unsavedChanges, setUnsavedChanges] = useState(false);

  const handleSettingChange = (key, value) => {
    setSettings(prev => ({ ...prev, [key]: value }));
    setUnsavedChanges(true);
  };

  const saveSettings = () => {
    // Save settings logic here
    console.log('Saving settings:', settings);
    setUnsavedChanges(false);
    // Show success toast
  };

  const resetSettings = () => {
    // Reset to defaults logic
    setUnsavedChanges(false);
  };

  const exportSettings = () => {
    const dataStr = JSON.stringify(settings, null, 2);
    const dataUri = 'data:application/json;charset=utf-8,'+ encodeURIComponent(dataStr);
    const exportFileDefaultName = 'shop-settings.json';
    const linkElement = document.createElement('a');
    linkElement.setAttribute('href', dataUri);
    linkElement.setAttribute('download', exportFileDefaultName);
    linkElement.click();
  };

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Settings</h1>
          <p className="text-muted-foreground mt-1">
            Configure your business settings and preferences
          </p>
        </div>
        <div className="flex space-x-3">
          {unsavedChanges && (
            <Alert className="w-auto inline-flex items-center p-2 border-orange-200 bg-orange-50">
              <AlertTriangle className="h-4 w-4 text-orange-600" />
              <span className="text-sm text-orange-700 ml-2">Unsaved changes</span>
            </Alert>
          )}
          <Button variant="outline" onClick={resetSettings}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Reset
          </Button>
          <Button onClick={saveSettings} disabled={!unsavedChanges}>
            <Save className="h-4 w-4 mr-2" />
            Save Changes
          </Button>
        </div>
      </div>

      {/* Settings Tabs */}
      <Card>
        <CardContent className="p-0">
          <Tabs defaultValue="business" className="w-full">
            <div className="border-b">
              <TabsList className="grid w-full grid-cols-6 bg-transparent">
                <TabsTrigger value="business" className="flex items-center space-x-2">
                  <Store className="h-4 w-4" />
                  <span className="hidden sm:inline">Business</span>
                </TabsTrigger>
                <TabsTrigger value="profile" className="flex items-center space-x-2">
                  <User className="h-4 w-4" />
                  <span className="hidden sm:inline">Profile</span>
                </TabsTrigger>
                <TabsTrigger value="system" className="flex items-center space-x-2">
                  <Settings className="h-4 w-4" />
                  <span className="hidden sm:inline">System</span>
                </TabsTrigger>
                <TabsTrigger value="pos" className="flex items-center space-x-2">
                  <CreditCard className="h-4 w-4" />
                  <span className="hidden sm:inline">POS</span>
                </TabsTrigger>
                <TabsTrigger value="security" className="flex items-center space-x-2">
                  <Shield className="h-4 w-4" />
                  <span className="hidden sm:inline">Security</span>
                </TabsTrigger>
                <TabsTrigger value="backup" className="flex items-center space-x-2">
                  <Database className="h-4 w-4" />
                  <span className="hidden sm:inline">Backup</span>
                </TabsTrigger>
              </TabsList>
            </div>

            {/* Business Settings */}
            <TabsContent value="business" className="p-6 space-y-6">
              <div>
                <h3 className="text-lg font-semibold mb-4">Business Information</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="businessName">Business Name</Label>
                      <Input
                        id="businessName"
                        value={settings.businessName}
                        onChange={(e) => handleSettingChange('businessName', e.target.value)}
                        placeholder="Your Business Name"
                      />
                    </div>
                    <div>
                      <Label htmlFor="businessPhone">Business Phone</Label>
                      <Input
                        id="businessPhone"
                        value={settings.businessPhone}
                        onChange={(e) => handleSettingChange('businessPhone', e.target.value)}
                        placeholder="01712345678"
                      />
                    </div>
                    <div>
                      <Label htmlFor="businessEmail">Business Email</Label>
                      <Input
                        id="businessEmail"
                        type="email"
                        value={settings.businessEmail}
                        onChange={(e) => handleSettingChange('businessEmail', e.target.value)}
                        placeholder="business@example.com"
                      />
                    </div>
                    <div>
                      <Label htmlFor="businessLicense">Business License/TIN</Label>
                      <Input
                        id="businessLicense"
                        value={settings.businessLicense}
                        onChange={(e) => handleSettingChange('businessLicense', e.target.value)}
                        placeholder="TIN-123456789012"
                      />
                    </div>
                  </div>
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="businessAddress">Business Address</Label>
                      <Textarea
                        id="businessAddress"
                        value={settings.businessAddress}
                        onChange={(e) => handleSettingChange('businessAddress', e.target.value)}
                        placeholder="Complete business address"
                        rows={4}
                      />
                    </div>
                    <div>
                      <Label>Business Logo</Label>
                      <div className="flex items-center space-x-4 mt-2">
                        <Avatar className="h-16 w-16">
                          <AvatarFallback className="text-lg bg-primary text-primary-foreground">
                            RS
                          </AvatarFallback>
                        </Avatar>
                        <div className="space-y-2">
                          <Button variant="outline" size="sm">
                            <Camera className="h-4 w-4 mr-2" />
                            Change Logo
                          </Button>
                          <p className="text-xs text-muted-foreground">
                            Recommended: 200x200px, PNG/JPG
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <Separator />

              <div>
                <h3 className="text-lg font-semibold mb-4">Operating Hours</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                  {['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'].map((day) => (
                    <Card key={day}>
                      <CardContent className="p-4">
                        <div className="space-y-3">
                          <div className="flex items-center justify-between">
                            <Label className="text-sm font-medium">{day}</Label>
                            <Switch defaultChecked={day !== 'Friday'} />
                          </div>
                          <div className="grid grid-cols-2 gap-2">
                            <div>
                              <Label className="text-xs">Open</Label>
                              <Input type="time" defaultValue="09:00" className="h-8 text-xs" />
                            </div>
                            <div>
                              <Label className="text-xs">Close</Label>
                              <Input type="time" defaultValue="21:00" className="h-8 text-xs" />
                            </div>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </div>
            </TabsContent>

            {/* Profile Settings */}
            <TabsContent value="profile" className="p-6 space-y-6">
              <div>
                <h3 className="text-lg font-semibold mb-4">Personal Information</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <div className="flex items-center space-x-4">
                      <Avatar className="h-20 w-20">
                        <AvatarFallback className="text-xl bg-primary text-primary-foreground">
                          MR
                        </AvatarFallback>
                      </Avatar>
                      <div className="space-y-2">
                        <Button variant="outline" size="sm">
                          <Camera className="h-4 w-4 mr-2" />
                          Change Photo
                        </Button>
                        <p className="text-xs text-muted-foreground">
                          JPG, PNG up to 2MB
                        </p>
                      </div>
                    </div>
                    <div>
                      <Label htmlFor="ownerName">Full Name</Label>
                      <Input
                        id="ownerName"
                        value={settings.ownerName}
                        onChange={(e) => handleSettingChange('ownerName', e.target.value)}
                      />
                    </div>
                    <div>
                      <Label htmlFor="ownerEmail">Email Address</Label>
                      <Input
                        id="ownerEmail"
                        type="email"
                        value={settings.ownerEmail}
                        onChange={(e) => handleSettingChange('ownerEmail', e.target.value)}
                      />
                    </div>
                    <div>
                      <Label htmlFor="ownerPhone">Phone Number</Label>
                      <Input
                        id="ownerPhone"
                        value={settings.ownerPhone}
                        onChange={(e) => handleSettingChange('ownerPhone', e.target.value)}
                      />
                    </div>
                  </div>
                  <div className="space-y-4">
                    <Card>
                      <CardHeader>
                        <CardTitle className="text-lg">Account Security</CardTitle>
                        <CardDescription>
                          Manage your account security settings
                        </CardDescription>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <Dialog open={showPasswordDialog} onOpenChange={setShowPasswordDialog}>
                          <DialogTrigger asChild>
                            <Button variant="outline" className="w-full justify-start">
                              <Key className="h-4 w-4 mr-2" />
                              Change Password
                            </Button>
                          </DialogTrigger>
                          <DialogContent>
                            <DialogHeader>
                              <DialogTitle>Change Password</DialogTitle>
                              <DialogDescription>
                                Update your account password for better security
                              </DialogDescription>
                            </DialogHeader>
                            <div className="space-y-4">
                              <div>
                                <Label htmlFor="currentPassword">Current Password</Label>
                                <Input id="currentPassword" type="password" />
                              </div>
                              <div>
                                <Label htmlFor="newPassword">New Password</Label>
                                <Input id="newPassword" type="password" />
                              </div>
                              <div>
                                <Label htmlFor="confirmPassword">Confirm Password</Label>
                                <Input id="confirmPassword" type="password" />
                              </div>
                            </div>
                            <DialogFooter>
                              <Button variant="outline" onClick={() => setShowPasswordDialog(false)}>
                                Cancel
                              </Button>
                              <Button onClick={() => setShowPasswordDialog(false)}>
                                Update Password
                              </Button>
                            </DialogFooter>
                          </DialogContent>
                        </Dialog>
                        <Button variant="outline" className="w-full justify-start">
                          <Smartphone className="h-4 w-4 mr-2" />
                          Two-Factor Authentication
                        </Button>
                        <Button variant="outline" className="w-full justify-start">
                          <Download className="h-4 w-4 mr-2" />
                          Download Account Data
                        </Button>
                      </CardContent>
                    </Card>
                  </div>
                </div>
              </div>

              <Separator />

              <div>
                <h3 className="text-lg font-semibold mb-4">Notification Preferences</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <Card>
                    <CardHeader>
                      <CardTitle className="text-base">Email Notifications</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">Daily Sales Report</Label>
                          <p className="text-xs text-muted-foreground">Receive daily sales summary</p>
                        </div>
                        <Switch
                          checked={settings.dailyReports}
                          onCheckedChange={(checked) => handleSettingChange('dailyReports', checked)}
                        />
                      </div>
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">Low Stock Alerts</Label>
                          <p className="text-xs text-muted-foreground">Alert when stock is low</p>
                        </div>
                        <Switch
                          checked={settings.lowStockAlerts}
                          onCheckedChange={(checked) => handleSettingChange('lowStockAlerts', checked)}
                        />
                      </div>
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">New Orders</Label>
                          <p className="text-xs text-muted-foreground">Notify for new customer orders</p>
                        </div>
                        <Switch defaultChecked />
                      </div>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardHeader>
                      <CardTitle className="text-base">SMS Notifications</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">Security Alerts</Label>
                          <p className="text-xs text-muted-foreground">Login attempts and security</p>
                        </div>
                        <Switch defaultChecked />
                      </div>
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">Payment Confirmations</Label>
                          <p className="text-xs text-muted-foreground">Payment received notifications</p>
                        </div>
                        <Switch
                          checked={settings.smsNotifications}
                          onCheckedChange={(checked) => handleSettingChange('smsNotifications', checked)}
                        />
                      </div>
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">System Updates</Label>
                          <p className="text-xs text-muted-foreground">Important system notifications</p>
                        </div>
                        <Switch defaultChecked />
                      </div>
                    </CardContent>
                  </Card>
                </div>
              </div>
            </TabsContent>

            {/* System Settings */}
            <TabsContent value="system" className="p-6 space-y-6">
              <div>
                <h3 className="text-lg font-semibold mb-4">System Preferences</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="currency">Currency</Label>
                      <Select value={settings.currency} onValueChange={(value) => handleSettingChange('currency', value)}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="BDT">Bangladeshi Taka (৳)</SelectItem>
                          <SelectItem value="USD">US Dollar ($)</SelectItem>
                          <SelectItem value="EUR">Euro (€)</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div>
                      <Label htmlFor="language">Language</Label>
                      <Select value={settings.language} onValueChange={(value) => handleSettingChange('language', value)}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="en">English</SelectItem>
                          <SelectItem value="bn">বাংলা (Bengali)</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div>
                      <Label htmlFor="timezone">Timezone</Label>
                      <Select value={settings.timezone} onValueChange={(value) => handleSettingChange('timezone', value)}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="Asia/Dhaka">Asia/Dhaka (GMT+6)</SelectItem>
                          <SelectItem value="UTC">UTC (GMT+0)</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div>
                      <Label htmlFor="dateFormat">Date Format</Label>
                      <Select value={settings.dateFormat} onValueChange={(value) => handleSettingChange('dateFormat', value)}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="DD/MM/YYYY">DD/MM/YYYY</SelectItem>
                          <SelectItem value="MM/DD/YYYY">MM/DD/YYYY</SelectItem>
                          <SelectItem value="YYYY-MM-DD">YYYY-MM-DD</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                  <div className="space-y-4">
                    <Card>
                      <CardHeader>
                        <CardTitle className="text-base">Appearance</CardTitle>
                        <CardDescription>Customize the look and feel</CardDescription>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Theme</Label>
                            <p className="text-xs text-muted-foreground">Light or dark mode</p>
                          </div>
                          <Select value={settings.theme} onValueChange={(value) => handleSettingChange('theme', value)}>
                            <SelectTrigger className="w-32">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="light">
                                <div className="flex items-center">
                                  <Sun className="h-4 w-4 mr-2" />
                                  Light
                                </div>
                              </SelectItem>
                              <SelectItem value="dark">
                                <div className="flex items-center">
                                  <Moon className="h-4 w-4 mr-2" />
                                  Dark
                                </div>
                              </SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Sound Effects</Label>
                            <p className="text-xs text-muted-foreground">Button clicks and alerts</p>
                          </div>
                          <Switch defaultChecked />
                        </div>
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Animations</Label>
                            <p className="text-xs text-muted-foreground">UI animations and transitions</p>
                          </div>
                          <Switch defaultChecked />
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                </div>
              </div>
            </TabsContent>

            {/* POS Settings */}
            <TabsContent value="pos" className="p-6 space-y-6">
              <div>
                <h3 className="text-lg font-semibold mb-4">Point of Sale Configuration</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="taxRate">Tax Rate (%)</Label>
                      <Input
                        id="taxRate"
                        type="number"
                        value={settings.taxRate}
                        onChange={(e) => handleSettingChange('taxRate', parseFloat(e.target.value) || 0)}
                        placeholder="0"
                        min="0"
                        max="100"
                        step="0.1"
                      />
                      <p className="text-xs text-muted-foreground mt-1">
                        VAT/Tax percentage applied to sales
                      </p>
                    </div>
                    <div>
                      <Label htmlFor="defaultPayment">Default Payment Method</Label>
                      <Select 
                        value={settings.defaultPaymentMethod} 
                        onValueChange={(value) => handleSettingChange('defaultPaymentMethod', value)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="cash">Cash</SelectItem>
                          <SelectItem value="bkash">bKash</SelectItem>
                          <SelectItem value="nagad">Nagad</SelectItem>
                          <SelectItem value="card">Card</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div>
                      <Label htmlFor="receiptTemplate">Receipt Template</Label>
                      <Select 
                        value={settings.receiptTemplate} 
                        onValueChange={(value) => handleSettingChange('receiptTemplate', value)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="standard">Standard</SelectItem>
                          <SelectItem value="compact">Compact</SelectItem>
                          <SelectItem value="detailed">Detailed</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                  <div className="space-y-4">
                    <Card>
                      <CardHeader>
                        <CardTitle className="text-base">POS Features</CardTitle>
                        <CardDescription>Configure POS system behavior</CardDescription>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Auto Print Receipts</Label>
                            <p className="text-xs text-muted-foreground">Print receipt after each sale</p>
                          </div>
                          <Switch
                            checked={settings.autoPrint}
                            onCheckedChange={(checked) => handleSettingChange('autoPrint', checked)}
                          />
                        </div>
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Barcode Scanner</Label>
                            <p className="text-xs text-muted-foreground">Enable barcode scanning</p>
                          </div>
                          <Switch defaultChecked />
                        </div>
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Customer Display</Label>
                            <p className="text-xs text-muted-foreground">Show prices to customer</p>
                          </div>
                          <Switch defaultChecked />
                        </div>
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Quick Keys</Label>
                            <p className="text-xs text-muted-foreground">Keyboard shortcuts</p>
                          </div>
                          <Switch defaultChecked />
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                </div>
              </div>

              <Separator />

              <div>
                <h3 className="text-lg font-semibold mb-4">Payment Methods</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                  {[
                    { name: 'Cash', icon: DollarSign, enabled: true, color: 'green' },
                    { name: 'bKash', icon: Smartphone, enabled: true, color: 'pink' },
                    { name: 'Nagad', icon: Smartphone, enabled: true, color: 'orange' },
                    { name: 'Card', icon: CreditCard, enabled: false, color: 'blue' }
                  ].map((method) => (
                    <Card key={method.name}>
                      <CardContent className="p-4">
                        <div className="flex items-center justify-between mb-3">
                          <div className="flex items-center space-x-2">
                            <method.icon className={`h-5 w-5 text-${method.color}-600`} />
                            <span className="font-medium">{method.name}</span>
                          </div>
                          <Switch defaultChecked={method.enabled} />
                        </div>
                        {method.name !== 'Cash' && (
                          <Button variant="outline" size="sm" className="w-full">
                            Configure
                          </Button>
                        )}
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </div>
            </TabsContent>

            {/* Security Settings */}
            <TabsContent value="security" className="p-6 space-y-6">
              <div>
                <h3 className="text-lg font-semibold mb-4">Security Configuration</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <Card>
                      <CardHeader>
                        <CardTitle className="text-base flex items-center space-x-2">
                          <Shield className="h-4 w-4" />
                          <span>Account Security</span>
                        </CardTitle>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Two-Factor Authentication</Label>
                            <p className="text-xs text-muted-foreground">Add extra layer of security</p>
                          </div>
                          <Switch
                            checked={settings.twoFactorAuth}
                            onCheckedChange={(checked) => handleSettingChange('twoFactorAuth', checked)}
                          />
                        </div>
                        <div>
                          <Label htmlFor="sessionTimeout">Session Timeout (minutes)</Label>
                          <Input
                            id="sessionTimeout"
                            type="number"
                            value={settings.sessionTimeout}
                            onChange={(e) => handleSettingChange('sessionTimeout', parseInt(e.target.value) || 30)}
                            min="5"
                            max="120"
                          />
                          <p className="text-xs text-muted-foreground mt-1">
                            Auto logout after inactivity
                          </p>
                        </div>
                        <div>
                          <Label htmlFor="passwordExpiry">Password Expiry (days)</Label>
                          <Input
                            id="passwordExpiry"
                            type="number"
                            value={settings.passwordExpiry}
                            onChange={(e) => handleSettingChange('passwordExpiry', parseInt(e.target.value) || 90)}
                            min="30"
                            max="365"
                          />
                          <p className="text-xs text-muted-foreground mt-1">
                            Force password change interval
                          </p>
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                  <div className="space-y-4">
                    <Card>
                      <CardHeader>
                        <CardTitle className="text-base flex items-center space-x-2">
                          <Lock className="h-4 w-4" />
                          <span>Access Control</span>
                        </CardTitle>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Employee Access</Label>
                            <p className="text-xs text-muted-foreground">Allow staff login access</p>
                          </div>
                          <Switch defaultChecked />
                        </div>
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Remote Access</Label>
                            <p className="text-xs text-muted-foreground">Access from outside network</p>
                          </div>
                          <Switch defaultChecked />
                        </div>
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Login Notifications</Label>
                            <p className="text-xs text-muted-foreground">Alert on new logins</p>
                          </div>
                          <Switch defaultChecked />
                        </div>
                        <Button variant="outline" className="w-full">
                          <Eye className="h-4 w-4 mr-2" />
                          View Login History
                        </Button>
                      </CardContent>
                    </Card>
                  </div>
                </div>
              </div>

              <Separator />

              <div>
                <h3 className="text-lg font-semibold mb-4">Data Protection</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <Card>
                    <CardHeader>
                      <CardTitle className="text-base">Privacy Settings</CardTitle>
                      <CardDescription>Control how your data is handled</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">Analytics Tracking</Label>
                          <p className="text-xs text-muted-foreground">Anonymous usage statistics</p>
                        </div>
                        <Switch defaultChecked />
                      </div>
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">Error Reporting</Label>
                          <p className="text-xs text-muted-foreground">Send crash reports</p>
                        </div>
                        <Switch defaultChecked />
                      </div>
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">Data Retention</Label>
                          <p className="text-xs text-muted-foreground">Keep data for compliance</p>
                        </div>
                        <Switch defaultChecked />
                      </div>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardHeader>
                      <CardTitle className="text-base">Encryption</CardTitle>
                      <CardDescription>Data security measures</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">Database Encryption</Label>
                          <p className="text-xs text-muted-foreground">Encrypt stored data</p>
                        </div>
                        <Badge variant="secondary" className="text-green-700 bg-green-100">
                          <CheckCircle className="h-3 w-3 mr-1" />
                          Enabled
                        </Badge>
                      </div>
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">SSL Certificate</Label>
                          <p className="text-xs text-muted-foreground">Secure connections</p>
                        </div>
                        <Badge variant="secondary" className="text-green-700 bg-green-100">
                          <CheckCircle className="h-3 w-3 mr-1" />
                          Active
                        </Badge>
                      </div>
                      <div className="flex items-center justify-between">
                        <div>
                          <Label className="text-sm font-medium">Backup Encryption</Label>
                          <p className="text-xs text-muted-foreground">Encrypt backup files</p>
                        </div>
                        <Badge variant="secondary" className="text-green-700 bg-green-100">
                          <CheckCircle className="h-3 w-3 mr-1" />
                          Enabled
                        </Badge>
                      </div>
                    </CardContent>
                  </Card>
                </div>
              </div>
            </TabsContent>

            {/* Backup Settings */}
            <TabsContent value="backup" className="p-6 space-y-6">
              <div>
                <h3 className="text-lg font-semibold mb-4">Backup Configuration</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <Card>
                      <CardHeader>
                        <CardTitle className="text-base">Automatic Backup</CardTitle>
                        <CardDescription>Schedule regular data backups</CardDescription>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="flex items-center justify-between">
                          <div>
                            <Label className="text-sm font-medium">Enable Auto Backup</Label>
                            <p className="text-xs text-muted-foreground">Automatically backup data</p>
                          </div>
                          <Switch
                            checked={settings.autoBackup}
                            onCheckedChange={(checked) => handleSettingChange('autoBackup', checked)}
                          />
                        </div>
                        <div>
                          <Label htmlFor="backupFrequency">Backup Frequency</Label>
                          <Select 
                            value={settings.backupFrequency} 
                            onValueChange={(value) => handleSettingChange('backupFrequency', value)}
                          >
                            <SelectTrigger>
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="hourly">Every Hour</SelectItem>
                              <SelectItem value="daily">Daily</SelectItem>
                              <SelectItem value="weekly">Weekly</SelectItem>
                              <SelectItem value="monthly">Monthly</SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                        <div className="flex items-center justify-between text-sm">
                          <span>Last Backup:</span>
                          <span className="font-medium">{settings.lastBackup}</span>
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                  <div className="space-y-4">
                    <Card>
                      <CardHeader>
                        <CardTitle className="text-base">Backup Storage</CardTitle>
                        <CardDescription>Configure backup storage options</CardDescription>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="space-y-3">
                          <div className="flex items-center justify-between">
                            <Label className="text-sm">Local Storage</Label>
                            <Switch defaultChecked />
                          </div>
                          <div className="flex items-center justify-between">
                            <Label className="text-sm">Google Drive</Label>
                            <Switch />
                          </div>
                          <div className="flex items-center justify-between">
                            <Label className="text-sm">Dropbox</Label>
                            <Switch />
                          </div>
                        </div>
                        <div className="space-y-2">
                          <Label className="text-sm">Storage Usage</Label>
                          <Progress value={65} className="h-2" />
                          <div className="flex justify-between text-xs text-muted-foreground">
                            <span>2.3 GB used</span>
                            <span>5 GB total</span>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                </div>
              </div>

              <Separator />

              <div>
                <h3 className="text-lg font-semibold mb-4">Backup Actions</h3>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <Card>
                    <CardContent className="p-6 text-center">
                      <Database className="h-12 w-12 mx-auto mb-4 text-blue-600" />
                      <h4 className="font-semibold mb-2">Create Backup</h4>
                      <p className="text-sm text-muted-foreground mb-4">
                        Create a manual backup of your data
                      </p>
                      <Button className="w-full">
                        <Download className="h-4 w-4 mr-2" />
                        Backup Now
                      </Button>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="p-6 text-center">
                      <Upload className="h-12 w-12 mx-auto mb-4 text-green-600" />
                      <h4 className="font-semibold mb-2">Restore Data</h4>
                      <p className="text-sm text-muted-foreground mb-4">
                        Restore from a previous backup
                      </p>
                      <Dialog open={showBackupDialog} onOpenChange={setShowBackupDialog}>
                        <DialogTrigger asChild>
                          <Button variant="outline" className="w-full">
                            <Upload className="h-4 w-4 mr-2" />
                            Restore
                          </Button>
                        </DialogTrigger>
                        <DialogContent>
                          <DialogHeader>
                            <DialogTitle>Restore from Backup</DialogTitle>
                            <DialogDescription>
                              Select a backup file to restore your data
                            </DialogDescription>
                          </DialogHeader>
                          <div className="space-y-4">
                            <div>
                              <Label>Select Backup File</Label>
                              <Input type="file" accept=".json,.sql,.zip" className="mt-1" />
                            </div>
                            <Alert>
                              <AlertTriangle className="h-4 w-4" />
                              <AlertTitle>Warning</AlertTitle>
                              <AlertDescription>
                                This will overwrite your current data. Make sure to create a backup first.
                              </AlertDescription>
                            </Alert>
                          </div>
                          <DialogFooter>
                            <Button variant="outline" onClick={() => setShowBackupDialog(false)}>
                              Cancel
                            </Button>
                            <Button onClick={() => setShowBackupDialog(false)}>
                              Restore Data
                            </Button>
                          </DialogFooter>
                        </DialogContent>
                      </Dialog>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="p-6 text-center">
                      <FileText className="h-12 w-12 mx-auto mb-4 text-purple-600" />
                      <h4 className="font-semibold mb-2">Export Settings</h4>
                      <p className="text-sm text-muted-foreground mb-4">
                        Export your configuration settings
                      </p>
                      <Button variant="outline" className="w-full" onClick={exportSettings}>
                        <Download className="h-4 w-4 mr-2" />
                        Export
                      </Button>
                    </CardContent>
                  </Card>
                </div>
              </div>

              <div>
                <h3 className="text-lg font-semibold mb-4">Recent Backups</h3>
                <Card>
                  <CardContent className="p-0">
                    <div className="space-y-0">
                      {[
                        { date: '2025-01-23 10:30 AM', size: '245 MB', type: 'Auto', status: 'success' },
                        { date: '2025-01-22 10:30 AM', size: '238 MB', type: 'Auto', status: 'success' },
                        { date: '2025-01-21 10:30 AM', size: '232 MB', type: 'Auto', status: 'success' },
                        { date: '2025-01-20 03:15 PM', size: '229 MB', type: 'Manual', status: 'success' },
                        { date: '2025-01-20 10:30 AM', size: '225 MB', type: 'Auto', status: 'failed' }
                      ].map((backup, index) => (
                        <div key={index} className="flex items-center justify-between p-4 border-b last:border-b-0">
                          <div className="flex items-center space-x-3">
                            <div className={cn(
                              "w-2 h-2 rounded-full",
                              backup.status === 'success' ? 'bg-green-500' : 'bg-red-500'
                            )}></div>
                            <div>
                              <p className="font-medium text-sm">{backup.date}</p>
                              <p className="text-xs text-muted-foreground">
                                {backup.size} • {backup.type} Backup
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center space-x-2">
                            <Badge variant={backup.status === 'success' ? 'default' : 'destructive'}>
                              {backup.status === 'success' ? 'Success' : 'Failed'}
                            </Badge>
                            <Button variant="ghost" size="sm">
                              <Download className="h-4 w-4" />
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>

      {/* Danger Zone */}
      <Card className="border-destructive/20">
        <CardHeader>
          <CardTitle className="text-destructive flex items-center space-x-2">
            <AlertTriangle className="h-5 w-5" />
            <span>Danger Zone</span>
          </CardTitle>
          <CardDescription>
            Irreversible and destructive actions
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between p-4 border border-destructive/20 rounded-lg">
            <div>
              <h4 className="font-medium">Reset All Settings</h4>
              <p className="text-sm text-muted-foreground">
                Reset all settings to default values
              </p>
            </div>
            <Button variant="destructive" size="sm">
              <RefreshCw className="h-4 w-4 mr-2" />
              Reset Settings
            </Button>
          </div>
          <div className="flex items-center justify-between p-4 border border-destructive/20 rounded-lg">
            <div>
              <h4 className="font-medium">Delete All Data</h4>
              <p className="text-sm text-muted-foreground">
                Permanently delete all business data
              </p>
            </div>
            <Button variant="destructive" size="sm">
              <Trash2 className="h-4 w-4 mr-2" />
              Delete Data
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default SettingsPage;