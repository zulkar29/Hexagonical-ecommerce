import React, { useState, useEffect } from 'react';
import {
  Globe,
  Key,
  Activity,
  Shield,
  BarChart3,
  Clock,
  Users,
  Zap,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Eye,
  EyeOff,
  Copy,
  Plus,
  Edit,
  Trash2,
  RefreshCw,
  Download,
  Upload,
  Search,
  Filter,
  MoreHorizontal,
  ExternalLink,
  Code,
  FileText,
  Settings,
  Lock,
  Unlock,
  Target,
  Gauge,
  Timer,
  Server,
  Database,
  Network,
  Bug,
  Info,
  AlertCircle,
  TrendingUp,
  TrendingDown,
  Minus,
  Hash,
  Link2,
  Webhook,
  Terminal,
  BookOpen,
  HelpCircle,
  PlayCircle,
  PauseCircle,
  StopCircle,
  RotateCcw,
  Calendar,
  MapPin,
  Smartphone,
  Monitor,  
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

// Custom Table Components
const Table = ({ children, ...props }) => (
  <div className="w-full overflow-auto">
    <table className="w-full caption-bottom text-sm" {...props}>
      {children}
    </table>
  </div>
);

const TableHeader = ({ children, ...props }) => (
  <thead className="[&_tr]:border-b" {...props}>
    {children}
  </thead>
);

const TableBody = ({ children, ...props }) => (
  <tbody className="[&_tr:last-child]:border-0" {...props}>
    {children}
  </tbody>
);

const TableRow = ({ children, className = "", ...props }) => (
  <tr className={`border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted ${className}`} {...props}>
    {children}
  </tr>
);

const TableHead = ({ children, className = "", ...props }) => (
  <th className={`h-12 px-4 text-left align-middle font-medium text-muted-foreground [&:has([role=checkbox])]:pr-0 ${className}`} {...props}>
    {children}
  </th>
);

const TableCell = ({ children, className = "", ...props }) => (
  <td className={`p-4 align-middle [&:has([role=checkbox])]:pr-0 ${className}`} {...props}>
    {children}
  </td>
);

const APIManagementPage = () => {
  const [activeTab, setActiveTab] = useState('overview');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedTimeframe, setSelectedTimeframe] = useState('24h');
  const [showApiKey, setShowApiKey] = useState({});
  const [isCreatingKey, setIsCreatingKey] = useState(false);
  const [newKeyData, setNewKeyData] = useState({
    name: '',
    description: '',
    permissions: [],
    rateLimit: 1000,
    expiresIn: '1year'
  });

  // API Keys Data
  const apiKeys = [
    {
      id: 'ak_001',
      name: 'Production API Key',
      description: 'Main production API key for tenant operations',
      key: 'sk_prod_1234567890abcdef1234567890abcdef',
      status: 'active',
      permissions: ['read', 'write', 'admin'],
      rateLimit: 5000,
      usage: 2847,
      usagePercentage: 57,
      lastUsed: '2025-01-24T10:30:00Z',
      createdAt: '2025-01-01T00:00:00Z',
      expiresAt: '2026-01-01T00:00:00Z',
      tenant: 'System',
      ipWhitelist: ['192.168.1.0/24'],
      webhookUrl: 'https://api.shopowner.com/webhooks'
    },
    {
      id: 'ak_002',
      name: 'Mobile App API',
      description: 'API key for mobile application integration',
      key: 'sk_mobile_abcdef1234567890abcdef1234567890',
      status: 'active',
      permissions: ['read', 'write'],
      rateLimit: 2000,
      usage: 1234,
      usagePercentage: 62,
      lastUsed: '2025-01-24T10:25:00Z',
      createdAt: '2025-01-15T00:00:00Z',
      expiresAt: '2025-07-15T00:00:00Z',
      tenant: 'Mobile Team',
      ipWhitelist: [],
      webhookUrl: null
    },
    {
      id: 'ak_003',
      name: 'Analytics Integration',
      description: 'Read-only access for analytics and reporting',
      key: 'sk_analytics_9876543210fedcba9876543210fedcba',
      status: 'active',
      permissions: ['read'],
      rateLimit: 10000,
      usage: 8234,
      usagePercentage: 82,
      lastUsed: '2025-01-24T10:20:00Z',
      createdAt: '2025-01-10T00:00:00Z',
      expiresAt: '2025-12-31T00:00:00Z',
      tenant: 'Analytics Team',
      ipWhitelist: ['203.190.37.0/24'],
      webhookUrl: 'https://analytics.shopowner.com/webhook'
    },
    {
      id: 'ak_004',
      name: 'Backup Service',
      description: 'Automated backup service API access',
      key: 'sk_backup_fedcba0987654321fedcba0987654321',
      status: 'suspended',
      permissions: ['read'],
      rateLimit: 500,
      usage: 456,
      usagePercentage: 91,
      lastUsed: '2025-01-23T15:45:00Z',
      createdAt: '2025-01-05T00:00:00Z',
      expiresAt: '2025-06-05T00:00:00Z',
      tenant: 'System',
      ipWhitelist: ['10.0.0.0/8'],
      webhookUrl: null
    }
  ];

  // API Endpoints Data
  const apiEndpoints = [
    {
      id: 'ep_001',
      method: 'GET',
      path: '/api/v1/tenants',
      description: 'List all tenants',
      status: 'healthy',
      avgResponseTime: 120,
      requests24h: 15420,
      successRate: 99.8,
      lastError: null,
      rateLimit: 1000,
      authRequired: true,
      deprecated: false
    },
    {
      id: 'ep_002',
      method: 'POST',
      path: '/api/v1/tenants',
      description: 'Create new tenant',
      status: 'healthy',
      avgResponseTime: 340,
      requests24h: 245,
      successRate: 98.2,
      lastError: null,
      rateLimit: 100,
      authRequired: true,
      deprecated: false
    },
    {
      id: 'ep_003',
      method: 'GET',
      path: '/api/v1/subscriptions',
      description: 'List subscriptions',
      status: 'warning',
      avgResponseTime: 890,
      requests24h: 8750,
      successRate: 97.5,
      lastError: '2025-01-24T09:15:00Z',
      rateLimit: 2000,
      authRequired: true,
      deprecated: false
    },
    {
      id: 'ep_004',
      method: 'POST',
      path: '/api/v1/payments',
      description: 'Process payment',
      status: 'critical',
      avgResponseTime: 2340,
      requests24h: 1240,
      successRate: 89.2,
      lastError: '2025-01-24T10:20:00Z',
      rateLimit: 500,
      authRequired: true,
      deprecated: false
    },
    {
      id: 'ep_005',
      method: 'GET',
      path: '/api/v1/users/profile',
      description: 'Get user profile (deprecated)',
      status: 'healthy',
      avgResponseTime: 85,
      requests24h: 890,
      successRate: 99.9,
      lastError: null,
      rateLimit: 5000,
      authRequired: true,
      deprecated: true
    }
  ];

  // API Usage Analytics
  const usageStats = {
    totalRequests24h: 127840,
    totalRequests7d: 892150,
    avgResponseTime: 245,
    successRate: 98.7,
    activeKeys: 3,
    totalKeys: 4,
    rateLimitHits: 234,
    errorRate: 1.3,
    topEndpoints: [
      { endpoint: '/api/v1/tenants', requests: 15420, percentage: 35.2 },
      { endpoint: '/api/v1/products', requests: 12340, percentage: 28.7 },
      { endpoint: '/api/v1/orders', requests: 8950, percentage: 20.8 },
      { endpoint: '/api/v1/payments', requests: 1240, percentage: 2.9 }
    ]
  };

  const formatNumber = (num) => {
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
    return num.toString();
  };

  const getStatusConfig = (status) => {
    const configs = {
      active: { color: 'bg-green-100 text-green-800', icon: CheckCircle, label: 'Active' },
      suspended: { color: 'bg-red-100 text-red-800', icon: XCircle, label: 'Suspended' },
      expired: { color: 'bg-gray-100 text-gray-800', icon: Clock, label: 'Expired' },
      healthy: { color: 'bg-green-100 text-green-800', icon: CheckCircle, label: 'Healthy' },
      warning: { color: 'bg-yellow-100 text-yellow-800', icon: AlertTriangle, label: 'Warning' },
      critical: { color: 'bg-red-100 text-red-800', icon: XCircle, label: 'Critical' }
    };
    return configs[status] || configs.active;
  };

  const getMethodColor = (method) => {
    const colors = {
      GET: 'text-blue-600 bg-blue-50 border-blue-200',
      POST: 'text-green-600 bg-green-50 border-green-200',
      PUT: 'text-orange-600 bg-orange-50 border-orange-200',
      DELETE: 'text-red-600 bg-red-50 border-red-200',
      PATCH: 'text-purple-600 bg-purple-50 border-purple-200'
    };
    return colors[method] || 'text-gray-600 bg-gray-50 border-gray-200';
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
  };

  const toggleKeyVisibility = (keyId) => {
    setShowApiKey(prev => ({
      ...prev,
      [keyId]: !prev[keyId]
    }));
  };

  const handleCreateKey = () => {
    setIsCreatingKey(true);
    // Simulate API call
    setTimeout(() => {
      setIsCreatingKey(false);
      setNewKeyData({
        name: '',
        description: '',
        permissions: [],
        rateLimit: 1000,
        expiresIn: '1year'
      });
    }, 2000);
  };

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <Globe className="h-5 w-5 text-primary" />
              <h1 className="text-xl font-semibold">API Management</h1>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="text-xs">
                v1.0
              </Badge>
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <div className="h-2 w-2 bg-green-500 rounded-full animate-pulse" />
                All endpoints operational
              </div>
            </div>
          </div>
          
          <div className="flex items-center gap-2">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search APIs, keys, endpoints..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 w-64"
              />
            </div>

            <Select value={selectedTimeframe} onValueChange={setSelectedTimeframe}>
              <SelectTrigger className="w-24">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="1h">1H</SelectItem>
                <SelectItem value="24h">24H</SelectItem>
                <SelectItem value="7d">7D</SelectItem>
                <SelectItem value="30d">30D</SelectItem>
              </SelectContent>
            </Select>
            
            <Button size="sm">
              <Plus className="h-4 w-4 mr-2" />
              New API Key
            </Button>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>API Actions</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <FileText className="h-4 w-4 mr-2" />
                  API Documentation
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Download className="h-4 w-4 mr-2" />
                  Export Usage Report
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Settings className="h-4 w-4 mr-2" />
                  Global Settings
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      {/* API Stats Overview */}
      <div className="border-b p-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Requests (24h)</p>
                  <p className="text-2xl font-bold">{formatNumber(usageStats.totalRequests24h)}</p>
                  <div className="flex items-center gap-1 mt-1">
                    <TrendingUp className="h-3 w-3 text-green-500" />
                    <p className="text-xs text-green-600">+12% from yesterday</p>
                  </div>
                </div>
                <BarChart3 className="h-8 w-8 text-blue-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Success Rate</p>
                  <p className="text-2xl font-bold">{usageStats.successRate}%</p>
                  <div className="flex items-center gap-1 mt-1">
                    <TrendingUp className="h-3 w-3 text-green-500" />
                    <p className="text-xs text-green-600">+0.3% improvement</p>
                  </div>
                </div>
                <CheckCircle className="h-8 w-8 text-green-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Avg Response Time</p>
                  <p className="text-2xl font-bold">{usageStats.avgResponseTime}ms</p>
                  <div className="flex items-center gap-1 mt-1">
                    <TrendingDown className="h-3 w-3 text-green-500" />
                    <p className="text-xs text-green-600">-15ms faster</p>
                  </div>
                </div>
                <Timer className="h-8 w-8 text-orange-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Active API Keys</p>
                  <p className="text-2xl font-bold">{usageStats.activeKeys}/{usageStats.totalKeys}</p>
                  <div className="flex items-center gap-1 mt-1">
                    <Minus className="h-3 w-3 text-gray-500" />
                    <p className="text-xs text-muted-foreground">No change</p>
                  </div>
                </div>
                <Key className="h-8 w-8 text-purple-500" />
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Main Content Tabs */}
      <div className="flex-1 overflow-hidden">
        <Tabs value={activeTab} onValueChange={setActiveTab} className="h-full flex flex-col">
          <div className="border-b">
            <TabsList className="w-full justify-start h-12 p-1 ml-6">
              <TabsTrigger value="overview" className="flex items-center gap-2">
                <Activity className="h-4 w-4" />
                Overview
              </TabsTrigger>
              <TabsTrigger value="keys" className="flex items-center gap-2">
                <Key className="h-4 w-4" />
                API Keys
              </TabsTrigger>
              <TabsTrigger value="endpoints" className="flex items-center gap-2">
                <Globe className="h-4 w-4" />
                Endpoints
              </TabsTrigger>
              <TabsTrigger value="analytics" className="flex items-center gap-2">
                <BarChart3 className="h-4 w-4" />
                Analytics
              </TabsTrigger>
              <TabsTrigger value="documentation" className="flex items-center gap-2">
                <BookOpen className="h-4 w-4" />
                Documentation
              </TabsTrigger>
            </TabsList>
          </div>

          <div className="flex-1 overflow-auto">
            {/* Overview Tab */}
            <TabsContent value="overview" className="h-full p-6 space-y-6">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Top Endpoints */}
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Target className="h-5 w-5" />
                      Top Endpoints (24h)
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      {usageStats.topEndpoints.map((endpoint, index) => (
                        <div key={index} className="flex items-center justify-between">
                          <div className="flex-1">
                            <div className="flex items-center gap-2 mb-1">
                              <span className="text-sm font-mono">{endpoint.endpoint}</span>
                              <Badge variant="outline" className="text-xs">
                                {formatNumber(endpoint.requests)}
                              </Badge>
                            </div>
                            <Progress value={endpoint.percentage} className="h-2" />
                          </div>
                          <span className="text-sm text-muted-foreground ml-4">
                            {endpoint.percentage}%
                          </span>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>

                {/* Recent Activity */}
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Clock className="h-5 w-5" />
                      Recent Activity
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div className="flex items-start gap-3">
                        <div className="w-2 h-2 bg-green-500 rounded-full mt-2" />
                        <div className="flex-1">
                          <p className="text-sm font-medium">New API key created</p>
                          <p className="text-xs text-muted-foreground">Mobile App API - 2 minutes ago</p>
                        </div>
                      </div>
                      <div className="flex items-start gap-3">
                        <div className="w-2 h-2 bg-yellow-500 rounded-full mt-2" />
                        <div className="flex-1">
                          <p className="text-sm font-medium">Rate limit exceeded</p>
                          <p className="text-xs text-muted-foreground">Analytics Integration - 15 minutes ago</p>
                        </div>
                      </div>
                      <div className="flex items-start gap-3">
                        <div className="w-2 h-2 bg-red-500 rounded-full mt-2" />
                        <div className="flex-1">
                          <p className="text-sm font-medium">Endpoint error spike</p>
                          <p className="text-xs text-muted-foreground">/api/v1/payments - 1 hour ago</p>
                        </div>
                      </div>
                      <div className="flex items-start gap-3">
                        <div className="w-2 h-2 bg-blue-500 rounded-full mt-2" />
                        <div className="flex-1">
                          <p className="text-sm font-medium">Documentation updated</p>
                          <p className="text-xs text-muted-foreground">API v1.2 release notes - 3 hours ago</p>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Quick Actions */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Zap className="h-5 w-5" />
                    Quick Actions
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                    <Button variant="outline" className="h-20 flex-col gap-2">
                      <Plus className="h-6 w-6" />
                      <span>Create API Key</span>
                    </Button>
                    <Button variant="outline" className="h-20 flex-col gap-2">
                      <FileText className="h-6 w-6" />
                      <span>View Docs</span>
                    </Button>
                    <Button variant="outline" className="h-20 flex-col gap-2">
                      <BarChart3 className="h-6 w-6" />
                      <span>Usage Report</span>
                    </Button>
                    <Button variant="outline" className="h-20 flex-col gap-2">
                      <Settings className="h-6 w-6" />
                      <span>API Settings</span>
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            {/* API Keys Tab */}
            <TabsContent value="keys" className="h-full p-6 space-y-6">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-medium">API Keys Management</h3>
                <Dialog>
                  <DialogTrigger asChild>
                    <Button>
                      <Plus className="h-4 w-4 mr-2" />
                      Create New Key
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="max-w-2xl">
                    <DialogHeader>
                      <DialogTitle>Create New API Key</DialogTitle>
                      <DialogDescription>
                        Generate a new API key with specific permissions and rate limits
                      </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                      <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                          <Label>Key Name</Label>
                          <Input
                            placeholder="e.g., Mobile App API"
                            value={newKeyData.name}
                            onChange={(e) => setNewKeyData(prev => ({...prev, name: e.target.value}))}
                          />
                        </div>
                        <div className="space-y-2">
                          <Label>Rate Limit (requests/hour)</Label>
                          <Input
                            type="number"
                            placeholder="1000"
                            value={newKeyData.rateLimit}
                            onChange={(e) => setNewKeyData(prev => ({...prev, rateLimit: parseInt(e.target.value)}))}
                          />
                        </div>
                      </div>
                      
                      <div className="space-y-2">
                        <Label>Description</Label>
                        <Textarea
                          placeholder="Describe the purpose of this API key..."
                          value={newKeyData.description}
                          onChange={(e) => setNewKeyData(prev => ({...prev, description: e.target.value}))}
                        />
                      </div>

                      <div className="space-y-2">
                        <Label>Permissions</Label>
                        <div className="flex gap-2">
                          {['read', 'write', 'admin'].map(permission => (
                            <div key={permission} className="flex items-center space-x-2">
                              <input
                                type="checkbox"
                                id={permission}
                                checked={newKeyData.permissions.includes(permission)}
                                onChange={(e) => {
                                  if (e.target.checked) {
                                    setNewKeyData(prev => ({
                                      ...prev,
                                      permissions: [...prev.permissions, permission]
                                    }));
                                  } else {
                                    setNewKeyData(prev => ({
                                      ...prev,
                                      permissions: prev.permissions.filter(p => p !== permission)
                                    }));
                                  }
                                }}
                              />
                              <Label htmlFor={permission} className="capitalize">{permission}</Label>
                            </div>
                          ))}
                        </div>
                      </div>

                      <div className="flex justify-end gap-2">
                        <Button variant="outline">Cancel</Button>
                        <Button onClick={handleCreateKey} disabled={isCreatingKey}>
                          {isCreatingKey ? (
                            <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                          ) : (
                            <Plus className="h-4 w-4 mr-2" />
                          )}
                          {isCreatingKey ? 'Creating...' : 'Create Key'}
                        </Button>
                      </div>
                    </div>
                  </DialogContent>
                </Dialog>
              </div>

              <div className="space-y-4">
                {apiKeys.map((key) => {
                  const statusConfig = getStatusConfig(key.status);
                  const isKeyVisible = showApiKey[key.id];
                  
                  return (
                    <Card key={key.id}>
                      <CardContent className="p-6">
                        <div className="flex items-start justify-between mb-4">
                          <div className="flex-1">
                            <div className="flex items-center gap-3 mb-2">
                              <h4 className="font-medium text-lg">{key.name}</h4>
                              <Badge className={statusConfig.color}>
                                <statusConfig.icon className="h-3 w-3 mr-1" />
                                {statusConfig.label}
                              </Badge>
                              <div className="flex items-center gap-1">
                                {key.permissions.map(permission => (
                                  <Badge key={permission} variant="outline" className="text-xs capitalize">
                                    {permission}
                                  </Badge>
                                ))}
                              </div>
                            </div>
                            <p className="text-sm text-muted-foreground mb-3">{key.description}</p>
                            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                              <div>
                                <p className="text-muted-foreground">Usage</p>
                                <p className="font-medium">{key.usage.toLocaleString()}/{key.rateLimit.toLocaleString()}</p>
                                <Progress value={key.usagePercentage} className="h-1 mt-1" />
                              </div>
                              <div>
                                <p className="text-muted-foreground">Last Used</p>
                                <p className="font-medium">{new Date(key.lastUsed).toLocaleDateString('en-BD')}</p>
                              </div>
                              <div>
                                <p className="text-muted-foreground">Created</p>
                                <p className="font-medium">{new Date(key.createdAt).toLocaleDateString('en-BD')}</p>
                              </div>
                              <div>
                                <p className="text-muted-foreground">Expires</p>
                                <p className="font-medium">{new Date(key.expiresAt).toLocaleDateString('en-BD')}</p>
                              </div>
                            </div>
                          </div>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="sm">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem>
                                <Edit className="h-4 w-4 mr-2" />
                                Edit Key
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                <RefreshCw className="h-4 w-4 mr-2" />
                                Regenerate
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                {key.status === 'active' ? (
                                  <>
                                    <PauseCircle className="h-4 w-4 mr-2" />
                                    Suspend
                                  </>
                                ) : (
                                  <>
                                    <PlayCircle className="h-4 w-4 mr-2" />
                                    Activate
                                  </>
                                )}
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem className="text-red-600">
                                <Trash2 className="h-4 w-4 mr-2" />
                                Delete Key
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </div>

                        <div className="space-y-3">
                          <div className="bg-muted/50 rounded-lg p-3">
                            <div className="flex items-center justify-between">
                              <div className="flex items-center gap-2">
                                <Key className="h-4 w-4 text-muted-foreground" />
                                <span className="text-sm font-medium">API Key</span>
                              </div>
                              <div className="flex items-center gap-2">
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => toggleKeyVisibility(key.id)}
                                >
                                  {isKeyVisible ? (
                                    <EyeOff className="h-4 w-4" />
                                  ) : (
                                    <Eye className="h-4 w-4" />
                                  )}
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => copyToClipboard(key.key)}
                                >
                                  <Copy className="h-4 w-4" />
                                </Button>
                              </div>
                            </div>
                            <code className="text-xs font-mono bg-background p-2 rounded mt-2 block">
                              {isKeyVisible ? key.key : '••••••••••••••••••••••••••••••••'}
                            </code>
                          </div>

                          {(key.ipWhitelist.length > 0 || key.webhookUrl) && (
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                              {key.ipWhitelist.length > 0 && (
                                <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                                  <div className="flex items-center gap-2 mb-2">
                                    <Shield className="h-4 w-4 text-blue-600" />
                                    <span className="text-sm font-medium text-blue-800">IP Whitelist</span>
                                  </div>
                                  <div className="space-y-1">
                                    {key.ipWhitelist.map((ip, index) => (
                                      <code key={index} className="text-xs text-blue-700 bg-blue-100 px-2 py-1 rounded">
                                        {ip}
                                      </code>
                                    ))}
                                  </div>
                                </div>
                              )}
                              
                              {key.webhookUrl && (
                                <div className="bg-green-50 border border-green-200 rounded-lg p-3">
                                  <div className="flex items-center gap-2 mb-2">
                                    <Webhook className="h-4 w-4 text-green-600" />
                                    <span className="text-sm font-medium text-green-800">Webhook URL</span>
                                  </div>
                                  <code className="text-xs text-green-700 bg-green-100 px-2 py-1 rounded break-all">
                                    {key.webhookUrl}
                                  </code>
                                </div>
                              )}
                            </div>
                          )}
                        </div>
                      </CardContent>
                    </Card>
                  );
                })}
              </div>
            </TabsContent>

            {/* Endpoints Tab */}
            <TabsContent value="endpoints" className="h-full p-6 space-y-6">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-medium">API Endpoints</h3>
                <div className="flex items-center gap-2">
                  <Select defaultValue="all">
                    <SelectTrigger className="w-32">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Status</SelectItem>
                      <SelectItem value="healthy">Healthy</SelectItem>
                      <SelectItem value="warning">Warning</SelectItem>
                      <SelectItem value="critical">Critical</SelectItem>
                    </SelectContent>
                  </Select>
                  <Button variant="outline" size="sm">
                    <RefreshCw className="h-4 w-4 mr-2" />
                    Refresh
                  </Button>
                </div>
              </div>

              <div className="border rounded-lg">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Endpoint</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Response Time</TableHead>
                      <TableHead>Requests (24h)</TableHead>
                      <TableHead>Success Rate</TableHead>
                      <TableHead>Rate Limit</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {apiEndpoints.map((endpoint) => {
                      const statusConfig = getStatusConfig(endpoint.status);
                      const methodColor = getMethodColor(endpoint.method);
                      
                      return (
                        <TableRow key={endpoint.id}>
                          <TableCell>
                            <div>
                              <div className="flex items-center gap-2 mb-1">
                                <Badge className={`text-xs border ${methodColor}`}>
                                  {endpoint.method}
                                </Badge>
                                <code className="text-sm font-mono">{endpoint.path}</code>
                                {endpoint.deprecated && (
                                  <Badge variant="outline" className="text-xs text-orange-600">
                                    Deprecated
                                  </Badge>
                                )}
                              </div>
                              <p className="text-xs text-muted-foreground">{endpoint.description}</p>
                            </div>
                          </TableCell>
                          <TableCell>
                            <Badge className={statusConfig.color}>
                              <statusConfig.icon className="h-3 w-3 mr-1" />
                              {statusConfig.label}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            <span className={`font-medium ${
                              endpoint.avgResponseTime > 1000 ? 'text-red-600' :
                              endpoint.avgResponseTime > 500 ? 'text-yellow-600' :
                              'text-green-600'
                            }`}>
                              {endpoint.avgResponseTime}ms
                            </span>
                          </TableCell>
                          <TableCell>{formatNumber(endpoint.requests24h)}</TableCell>
                          <TableCell>
                            <span className={`font-medium ${
                              endpoint.successRate > 99 ? 'text-green-600' :
                              endpoint.successRate > 95 ? 'text-yellow-600' :
                              'text-red-600'
                            }`}>
                              {endpoint.successRate}%
                            </span>
                          </TableCell>
                          <TableCell>
                            <span className="text-sm">{endpoint.rateLimit}/hr</span>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-1">
                              <TooltipProvider>
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button variant="ghost" size="sm" className="h-7 w-7 p-0">
                                      <BarChart3 className="h-3 w-3" />
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>View Analytics</TooltipContent>
                                </Tooltip>
                              </TooltipProvider>
                              
                              <TooltipProvider>
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button variant="ghost" size="sm" className="h-7 w-7 p-0">
                                      <FileText className="h-3 w-3" />
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>View Documentation</TooltipContent>
                                </Tooltip>
                              </TooltipProvider>

                              <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                  <Button variant="ghost" size="sm" className="h-7 w-7 p-0">
                                    <MoreHorizontal className="h-3 w-3" />
                                  </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                  <DropdownMenuItem>
                                    <PlayCircle className="h-4 w-4 mr-2" />
                                    Test Endpoint
                                  </DropdownMenuItem>
                                  <DropdownMenuItem>
                                    <Settings className="h-4 w-4 mr-2" />
                                    Configure
                                  </DropdownMenuItem>
                                  {endpoint.deprecated && (
                                    <DropdownMenuItem className="text-red-600">
                                      <XCircle className="h-4 w-4 mr-2" />
                                      Disable
                                    </DropdownMenuItem>
                                  )}
                                </DropdownMenuContent>
                              </DropdownMenu>
                            </div>
                          </TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              </div>

              {/* Endpoint Health Summary */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <Card>
                  <CardContent className="p-4 text-center">
                    <div className="text-2xl font-bold text-green-600">
                      {apiEndpoints.filter(e => e.status === 'healthy').length}
                    </div>
                    <div className="text-sm text-muted-foreground">Healthy Endpoints</div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardContent className="p-4 text-center">
                    <div className="text-2xl font-bold text-yellow-600">
                      {apiEndpoints.filter(e => e.status === 'warning').length}
                    </div>
                    <div className="text-sm text-muted-foreground">Warning Endpoints</div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardContent className="p-4 text-center">
                    <div className="text-2xl font-bold text-red-600">
                      {apiEndpoints.filter(e => e.status === 'critical').length}
                    </div>
                    <div className="text-sm text-muted-foreground">Critical Endpoints</div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>

            {/* Analytics Tab */}
            <TabsContent value="analytics" className="h-full p-6 space-y-6">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Request Volume Chart */}
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <BarChart3 className="h-5 w-5" />
                      Request Volume ({selectedTimeframe})
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="h-64 flex items-center justify-center bg-muted/20 rounded">
                      <div className="text-center">
                        <BarChart3 className="h-12 w-12 text-muted-foreground mx-auto mb-2" />
                        <p className="text-sm text-muted-foreground">Request Volume Chart</p>
                        <div className="mt-4 grid grid-cols-2 gap-4 text-xs">
                          <div className="text-center">
                            <p className="text-2xl font-bold text-blue-600">{formatNumber(usageStats.totalRequests24h)}</p>
                            <p className="text-muted-foreground">Total Requests</p>
                          </div>
                          <div className="text-center">
                            <p className="text-2xl font-bold text-green-600">{usageStats.successRate}%</p>
                            <p className="text-muted-foreground">Success Rate</p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                {/* Response Time Trends */}
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Timer className="h-5 w-5" />
                      Response Time Trends
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="h-64 flex items-center justify-center bg-muted/20 rounded">
                      <div className="text-center">
                        <Timer className="h-12 w-12 text-muted-foreground mx-auto mb-2" />
                        <p className="text-sm text-muted-foreground">Response Time Chart</p>
                        <div className="mt-4 space-y-2">
                          <div className="flex items-center justify-between text-xs">
                            <span className="flex items-center gap-2">
                              <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
                              Average
                            </span>
                            <span>{usageStats.avgResponseTime}ms</span>
                          </div>
                          <div className="flex items-center justify-between text-xs">
                            <span className="flex items-center gap-2">
                              <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                              P95
                            </span>
                            <span>340ms</span>
                          </div>
                          <div className="flex items-center justify-between text-xs">
                            <span className="flex items-center gap-2">
                              <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                              P99
                            </span>
                            <span>1.2s</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Error Analysis */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Bug className="h-5 w-5" />
                    Error Analysis
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <div className="text-center">
                      <div className="text-2xl font-bold text-red-600">{usageStats.errorRate}%</div>
                      <div className="text-sm text-muted-foreground">Error Rate</div>
                      <div className="text-xs text-green-600 mt-1">-0.2% from yesterday</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-orange-600">{usageStats.rateLimitHits}</div>
                      <div className="text-sm text-muted-foreground">Rate Limit Hits</div>
                      <div className="text-xs text-red-600 mt-1">+15% from yesterday</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-blue-600">2.1s</div>
                      <div className="text-sm text-muted-foreground">Avg Error Response</div>
                      <div className="text-xs text-green-600 mt-1">-0.3s faster</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-purple-600">12</div>
                      <div className="text-sm text-muted-foreground">Unique Errors</div>
                      <div className="text-xs text-muted-foreground mt-1">No change</div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Usage by Client */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Users className="h-5 w-5" />
                    Usage by API Key
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {apiKeys.map((key, index) => (
                      <div key={key.id} className="flex items-center justify-between p-3 border rounded">
                        <div className="flex items-center gap-3">
                          <Avatar className="h-8 w-8">
                            <AvatarFallback className="text-xs">
                              {key.name.split(' ').map(n => n[0]).join('')}
                            </AvatarFallback>
                          </Avatar>
                          <div>
                            <p className="font-medium text-sm">{key.name}</p>
                            <p className="text-xs text-muted-foreground">{key.tenant}</p>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="font-medium">{key.usage.toLocaleString()}</p>
                          <p className="text-xs text-muted-foreground">requests</p>
                        </div>
                        <div className="w-20">
                          <Progress value={key.usagePercentage} className="h-2" />
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            {/* Documentation Tab */}
            <TabsContent value="documentation" className="h-full p-6 space-y-6">
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* API Reference */}
                <div className="lg:col-span-2 space-y-6">
                  <Card>
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2">
                        <BookOpen className="h-5 w-5" />
                        API Reference
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="bg-muted/50 rounded-lg p-4">
                        <h4 className="font-medium mb-2">Base URL</h4>
                        <code className="text-sm bg-background p-2 rounded border block">
                          https://api.shopowner.com/v1
                        </code>
                      </div>

                      <div className="bg-muted/50 rounded-lg p-4">
                        <h4 className="font-medium mb-2">Authentication</h4>
                        <p className="text-sm text-muted-foreground mb-3">
                          Include your API key in the Authorization header:
                        </p>
                        <code className="text-sm bg-background p-2 rounded border block">
                          Authorization: Bearer your_api_key_here
                        </code>
                      </div>

                      <div className="space-y-3">
                        <h4 className="font-medium">Available Endpoints</h4>
                        {['Tenants', 'Users', 'Subscriptions', 'Payments', 'Products', 'Orders'].map((category) => (
                          <div key={category} className="border rounded-lg p-3">
                            <div className="flex items-center justify-between">
                              <span className="font-medium">{category}</span>
                              <div className="flex items-center gap-2">
                                <Badge variant="outline" className="text-xs">v1.0</Badge>
                                <Button variant="ghost" size="sm">
                                  <ExternalLink className="h-3 w-3" />
                                </Button>
                              </div>
                            </div>
                            <p className="text-sm text-muted-foreground mt-1">
                              Manage {category.toLowerCase()} data and operations
                            </p>
                          </div>
                        ))}
                      </div>
                    </CardContent>
                  </Card>
                </div>

                {/* Quick Links */}
                <div className="space-y-6">
                  <Card>
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2">
                        <Link2 className="h-5 w-5" />
                        Quick Links
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-3">
                      {[
                        { title: 'Getting Started', icon: PlayCircle, url: '#' },
                        { title: 'Authentication Guide', icon: Key, url: '#' },
                        { title: 'Rate Limiting', icon: Gauge, url: '#' },
                        { title: 'Webhooks', icon: Webhook, url: '#' },
                        { title: 'SDKs & Libraries', icon: Code, url: '#' },
                        { title: 'Postman Collection', icon: Download, url: '#' }
                      ].map((link, index) => (
                        <Button
                          key={index}
                          variant="ghost"
                          className="w-full justify-start gap-2"
                        >
                          <link.icon className="h-4 w-4" />
                          {link.title}
                          <ExternalLink className="h-3 w-3 ml-auto" />
                        </Button>
                      ))}
                    </CardContent>
                  </Card>

                  <Card>
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2">
                        <Terminal className="h-5 w-5" />
                        Code Examples
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-3">
                      {['JavaScript', 'Python', 'PHP', 'cURL'].map((lang, index) => (
                        <Button
                          key={index}
                          variant="outline"
                          className="w-full justify-start gap-2"
                        >
                          <Code className="h-4 w-4" />
                          {lang} Examples
                        </Button>
                      ))}
                    </CardContent>
                  </Card>

                  <Card>
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2">
                        <HelpCircle className="h-5 w-5" />
                        Support Resources
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-3">
                      <Button variant="ghost" className="w-full justify-start gap-2">
                        <FileText className="h-4 w-4" />
                        FAQ
                      </Button>
                      <Button variant="ghost" className="w-full justify-start gap-2">
                        <Users className="h-4 w-4" />
                        Developer Community
                      </Button>
                      <Button variant="ghost" className="w-full justify-start gap-2">
                        <Bug className="h-4 w-4" />
                        Report Issues
                      </Button>
                    </CardContent>
                  </Card>
                </div>
              </div>
            </TabsContent>
          </div>
        </Tabs>
      </div>
    </div>
  );
};

export default APIManagementPage;