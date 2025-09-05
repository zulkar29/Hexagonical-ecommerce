import React, { useState } from 'react';
import {
  Download,
  Upload,
  RefreshCw,
  Save,
  Info,
  CheckCircle,
  AlertTriangle,
  Calendar,
  Clock,
  Trash2,
  Database,
  Folder,
  Eye,
  EyeOff,
  Copy
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { Textarea } from '@/components/ui/textarea';

const initialHistory = [
  { id: 1, date: '2024-07-25 02:00', size: '120MB', status: 'Success', location: '/backups/db/backup-2024-07-25.sql' },
  { id: 2, date: '2024-07-24 02:00', size: '119MB', status: 'Success', location: '/backups/db/backup-2024-07-24.sql' },
  { id: 3, date: '2024-07-23 02:00', size: '118MB', status: 'Success', location: '/backups/db/backup-2024-07-23.sql' },
  { id: 4, date: '2024-07-22 02:00', size: '117MB', status: 'Failed', location: '/backups/db/backup-2024-07-22.sql' },
];

export default function BackupsSettings() {
  const [settings, setSettings] = useState({
    autoBackup: true,
    backupTime: '02:00',
    backupLocation: '/backups/db',
    retentionDays: 30,
    notifyOnFailure: true,
    notes: ''
  });
  const [isSaving, setIsSaving] = useState(false);
  const [history, setHistory] = useState(initialHistory);
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

  const handleManualBackup = () => {
    setHistory([
      { id: Date.now(), date: new Date().toISOString().slice(0, 16).replace('T', ' '), size: '121MB', status: 'Success', location: '/backups/db/backup-' + new Date().toISOString().slice(0, 10) + '.sql' },
      ...history
    ]);
  };

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center gap-4 px-6">
          <Folder className="h-5 w-5 text-primary" />
          <h1 className="text-xl font-semibold">Backups</h1>
          <Badge variant="outline" className="ml-2">Database</Badge>
          <div className="ml-auto flex gap-2">
            <Button size="sm" variant="outline" onClick={handleSave} disabled={isSaving}>
              {isSaving ? <RefreshCw className="h-4 w-4 mr-2 animate-spin" /> : <Save className="h-4 w-4 mr-2" />}
              {isSaving ? 'Saving...' : 'Save Settings'}
            </Button>
          </div>
        </div>
      </div>
      <div className="flex-1 overflow-auto p-6 space-y-6">
        {/* Backup Overview */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Database className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-base">Backup Overview</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <Label>Backup Location</Label>
                <Input value={settings.backupLocation} onChange={e => handleChange('backupLocation', e.target.value)} />
              </div>
              <div>
                <Label>Retention Period (days)</Label>
                <Input type="number" value={settings.retentionDays} onChange={e => handleChange('retentionDays', e.target.value)} />
              </div>
              <div className="flex items-center gap-2 mt-2">
                <Switch checked={settings.autoBackup} onCheckedChange={v => handleChange('autoBackup', v)} id="autoBackup" />
                <Label htmlFor="autoBackup">Enable Automatic Backups</Label>
              </div>
              <div className="flex items-center gap-2 mt-2">
                <Switch checked={settings.notifyOnFailure} onCheckedChange={v => handleChange('notifyOnFailure', v)} id="notifyOnFailure" />
                <Label htmlFor="notifyOnFailure">Notify on Failure</Label>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Manual Backup */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Download className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-base">Manual Backup</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2">
              <Button variant="outline" size="sm" onClick={handleManualBackup}>
                <Download className="h-4 w-4 mr-2" />
                Run Manual Backup
              </Button>
              <Button variant="outline" size="sm">
                <Upload className="h-4 w-4 mr-2" />
                Upload Backup
              </Button>
            </div>
            <div className="flex items-center gap-2 mt-2">
              <CheckCircle className="h-4 w-4 text-green-600" />
              <span className="text-sm text-green-700">Last backup was successful.</span>
            </div>
            <div className="flex items-center gap-2 mt-2">
              <AlertTriangle className="h-4 w-4 text-yellow-600" />
              <span className="text-sm text-yellow-700">Remember to download backups regularly.</span>
            </div>
          </CardContent>
        </Card>

        {/* Scheduled Backups */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-base">Scheduled Backups</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <Label>Backup Time (24h)</Label>
                <Input type="time" value={settings.backupTime} onChange={e => handleChange('backupTime', e.target.value)} />
              </div>
              <div>
                <Label>Notes</Label>
                <Textarea value={settings.notes} onChange={e => handleChange('notes', e.target.value)} placeholder="Add any notes about backup schedule..." />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Restore Options */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Upload className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-base">Restore Options</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2">
              <Button variant="outline" size="sm">
                <Upload className="h-4 w-4 mr-2" />
                Restore from File
              </Button>
              <Button variant="outline" size="sm">
                <Database className="h-4 w-4 mr-2" />
                Restore from Cloud
              </Button>
            </div>
            <div className="flex items-center gap-2 mt-2">
              <Info className="h-4 w-4 text-blue-600" />
              <span className="text-sm text-blue-700">Restoring will overwrite current data. Proceed with caution.</span>
            </div>
          </CardContent>
        </Card>

        {/* Backup History */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <HistoryIcon />
              <CardTitle className="text-base">Backup History</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="overflow-x-auto">
              <table className="min-w-full text-sm border">
                <thead>
                  <tr className="bg-muted/50">
                    <th className="p-2 text-left font-medium">Date</th>
                    <th className="p-2 text-left font-medium">Size</th>
                    <th className="p-2 text-left font-medium">Status</th>
                    <th className="p-2 text-left font-medium">Location</th>
                    <th className="p-2 text-left font-medium">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {history.map((item) => (
                    <tr key={item.id} className="border-t">
                      <td className="p-2">{item.date}</td>
                      <td className="p-2">{item.size}</td>
                      <td className="p-2">
                        {item.status === 'Success' ? (
                          <Badge variant="success">Success</Badge>
                        ) : (
                          <Badge variant="destructive">Failed</Badge>
                        )}
                      </td>
                      <td className="p-2 flex items-center gap-2">
                        <span>{item.location}</span>
                        <Button size="icon" variant="outline" onClick={() => handleCopy(item.location)}>
                          <Copy className="h-4 w-4" />
                        </Button>
                        {copied && <Badge variant="success">Copied!</Badge>}
                      </td>
                      <td className="p-2 flex gap-2">
                        <Button size="icon" variant="outline">
                          <Download className="h-4 w-4" />
                        </Button>
                        <Button size="icon" variant="outline">
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

// Custom icon for history
function HistoryIcon(props) {
  return <Clock {...props} />;
} 