import React, { useState } from 'react';
import {
  Database,
  Download,
  Upload,
  RefreshCw,
  Shield,
  Trash2,
  Save,
  Info,
  CheckCircle,
  AlertTriangle,
  Eye,
  Copy,
  EyeOff,
  Wrench,
  Activity
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';

export default function DatabaseSettings() {
  const [settings, setSettings] = useState({
    host: 'db.shopowner.com',
    port: 5432,
    username: 'admin',
    password: '••••••••••',
    dbName: 'shopowner_prod',
    ssl: true,
    autoBackup: true,
    backupLocation: '/backups/db',
    lastBackup: '2024-07-25 02:00',
    maintenanceMode: false,
    allowRemote: false,
    auditLogging: true,
    notes: ''
  });
  const [isSaving, setIsSaving] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [copied, setCopied] = useState(false);

  const handleChange = (key, value) => {
    setSettings((prev) => ({ ...prev, [key]: value }));
  };

  const handleSave = () => {
    setIsSaving(true);
    setTimeout(() => setIsSaving(false), 1200);
  };

  const handleCopy = (value) => {
    navigator.clipboard.writeText(value);
    setCopied(true);
    setTimeout(() => setCopied(false), 1000);
  };

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center gap-4 px-6">
          <Database className="h-5 w-5 text-primary" />
          <h1 className="text-xl font-semibold">Database Settings</h1>
          <Badge variant="outline" className="ml-2">Production</Badge>
          <div className="ml-auto flex gap-2">
            <Button size="sm" variant="outline" onClick={handleSave} disabled={isSaving}>
              {isSaving ? <RefreshCw className="h-4 w-4 mr-2 animate-spin" /> : <Save className="h-4 w-4 mr-2" />}
              {isSaving ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </div>
      </div>
      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Database Connection Info */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Database className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-base">Connection Info</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <Label>Host</Label>
                <Input value={settings.host} onChange={e => handleChange('host', e.target.value)} />
              </div>
              <div>
                <Label>Port</Label>
                <Input type="number" value={settings.port} onChange={e => handleChange('port', e.target.value)} />
              </div>
              <div>
                <Label>Username</Label>
                <Input value={settings.username} onChange={e => handleChange('username', e.target.value)} />
              </div>
              <div>
                <Label>Password</Label>
                <div className="flex gap-2">
                  <Input type={showPassword ? 'text' : 'password'} value={settings.password} onChange={e => handleChange('password', e.target.value)} />
                  <Button type="button" size="icon" variant="outline" onClick={() => setShowPassword(v => !v)}>
                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                  </Button>
                </div>
              </div>
              <div>
                <Label>Database Name</Label>
                <Input value={settings.dbName} onChange={e => handleChange('dbName', e.target.value)} />
              </div>
              <div className="flex items-center gap-2 mt-2">
                <Switch checked={settings.ssl} onCheckedChange={v => handleChange('ssl', v)} id="ssl" />
                <Label htmlFor="ssl">Require SSL</Label>
              </div>
              <div className="flex items-center gap-2 mt-2">
                <Switch checked={settings.allowRemote} onCheckedChange={v => handleChange('allowRemote', v)} id="allowRemote" />
                <Label htmlFor="allowRemote">Allow Remote Connections</Label>
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Info className="h-4 w-4 text-muted-foreground cursor-pointer" />
                    </TooltipTrigger>
                    <TooltipContent side="top">Enabling this may reduce security.</TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Backup & Restore */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Download className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-base">Backup & Restore</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex flex-col md:flex-row md:items-center gap-4">
              <div className="flex-1">
                <Label>Backup Location</Label>
                <Input value={settings.backupLocation} onChange={e => handleChange('backupLocation', e.target.value)} />
              </div>
              <div className="flex-1">
                <Label>Last Backup</Label>
                <div className="flex items-center gap-2">
                  <Input value={settings.lastBackup} readOnly />
                  <Button size="icon" variant="outline" onClick={() => handleCopy(settings.lastBackup)}>
                    <Copy className="h-4 w-4" />
                  </Button>
                  {copied && <Badge variant="success">Copied!</Badge>}
                </div>
              </div>
            </div>
            <div className="flex gap-2">
              <Button variant="outline" size="sm">
                <Download className="h-4 w-4 mr-2" />
                Download Backup
              </Button>
              <Button variant="outline" size="sm">
                <Upload className="h-4 w-4 mr-2" />
                Restore Backup
              </Button>
              <Button variant="outline" size="sm">
                <RefreshCw className="h-4 w-4 mr-2" />
                Run Backup Now
              </Button>
            </div>
            <div className="flex items-center gap-2 mt-2">
              <Switch checked={settings.autoBackup} onCheckedChange={v => handleChange('autoBackup', v)} id="autoBackup" />
              <Label htmlFor="autoBackup">Enable Automatic Backups</Label>
            </div>
          </CardContent>
        </Card>

        {/* Maintenance Tools */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Wrench className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-base">Maintenance Tools</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2">
              <Button variant="outline" size="sm">
                <Activity className="h-4 w-4 mr-2" />
                Optimize Database
              </Button>
              <Button variant="outline" size="sm">
                <Trash2 className="h-4 w-4 mr-2" />
                Clear Logs
              </Button>
              <Button variant="outline" size="sm">
                <RefreshCw className="h-4 w-4 mr-2" />
                Reindex Tables
              </Button>
            </div>
            <div className="flex items-center gap-2 mt-2">
              <Switch checked={settings.maintenanceMode} onCheckedChange={v => handleChange('maintenanceMode', v)} id="maintenanceMode" />
              <Label htmlFor="maintenanceMode">Enable Maintenance Mode</Label>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Info className="h-4 w-4 text-muted-foreground cursor-pointer" />
                  </TooltipTrigger>
                  <TooltipContent side="top">Only admins can access the platform when enabled.</TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>
          </CardContent>
        </Card>

        {/* Security */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Shield className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-base">Security</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-2 mt-2">
              <Switch checked={settings.auditLogging} onCheckedChange={v => handleChange('auditLogging', v)} id="auditLogging" />
              <Label htmlFor="auditLogging">Enable Audit Logging</Label>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Info className="h-4 w-4 text-muted-foreground cursor-pointer" />
                  </TooltipTrigger>
                  <TooltipContent side="top">Track all database changes for compliance.</TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>
            <div>
              <Label>Notes</Label>
              <Textarea value={settings.notes} onChange={e => handleChange('notes', e.target.value)} placeholder="Add any notes or comments about database security..." />
            </div>
            <div className="flex items-center gap-2 mt-2">
              <CheckCircle className="h-4 w-4 text-green-600" />
              <span className="text-sm text-green-700">Database is secure and compliant.</span>
            </div>
            <div className="flex items-center gap-2 mt-2">
              <AlertTriangle className="h-4 w-4 text-yellow-600" />
              <span className="text-sm text-yellow-700">Remember to regularly update your backup and security settings.</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
} 