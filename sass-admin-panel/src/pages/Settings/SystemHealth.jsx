import React, { useState, useEffect } from 'react';
import {
  Shield,
  Activity,
  Server,
  Database,
  Globe,
  Cpu,
  HardDrive,
  MemoryStick,
  Wifi,
  Zap,
  AlertTriangle,
  CheckCircle,
  Clock,
  TrendingUp,
  TrendingDown,
  RefreshCw,
  Settings,
  Download,
  Bell,
  Eye,
  BarChart3,
  PieChart,
  LineChart,
  Monitor,
  Cloud,
  Lock,
  Key,
  FileText,
  Users,
  Package,
  CreditCard,
  MessageSquare,
  Mail,
  Phone,
  Calendar,
  Archive,
  Search,
  Filter,
  MoreHorizontal,
  ExternalLink,
  AlertCircle,
  Info,
  XCircle,
  Thermometer,
  Signal,
  Timer,
  Target,
  Gauge,
  NetworkIcon as Network
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

const SystemHealthPage = () => {
  const [selectedTimeRange, setSelectedTimeRange] = useState('1h');
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [selectedMetric, setSelectedMetric] = useState('overview');
  const [isLoading, setIsLoading] = useState(false);
  const [lastUpdated, setLastUpdated] = useState(new Date());

  // System Status Data
  const systemStatus = {
    overall: 'operational', // operational, degraded, maintenance, outage
    uptime: '99.97%',
    lastIncident: '2025-01-20T14:30:00Z',
    totalServices: 12,
    healthyServices: 11,
    warningServices: 1,
    criticalServices: 0
  };

  // Core Services Status
  const coreServices = [
    {
      id: 'api',
      name: 'API Gateway',
      status: 'operational',
      uptime: '99.99%',
      responseTime: 145,
      requests: '2.4M',
      errors: 0.01,
      lastCheck: '2025-01-24T10:30:00Z',
      description: 'Main API endpoints for tenant operations'
    },
    {
      id: 'database',
      name: 'Primary Database',
      status: 'operational',
      uptime: '99.98%',
      responseTime: 23,
      connections: 85,
      errors: 0,
      lastCheck: '2025-01-24T10:30:00Z',
      description: 'MongoDB cluster for application data'
    },
    {
      id: 'auth',
      name: 'Authentication Service',
      status: 'operational',
      uptime: '99.96%',
      responseTime: 89,
      requests: '450K',
      errors: 0.02,
      lastCheck: '2025-01-24T10:30:00Z',
      description: 'JWT-based authentication system'
    },
    {
      id: 'payments',
      name: 'Payment Processing',
      status: 'warning',
      uptime: '99.45%',
      responseTime: 2340,
      requests: '125K',
      errors: 0.55,
      lastCheck: '2025-01-24T10:30:00Z',
      description: 'bKash, Nagad, and bank transfer integration',
      issues: ['High response time from bKash API', 'Intermittent timeout errors']
    },
    {
      id: 'storage',
      name: 'File Storage',
      status: 'operational',
      uptime: '99.99%',
      responseTime: 178,
      requests: '890K',
      errors: 0,
      lastCheck: '2025-01-24T10:30:00Z',
      description: 'AWS S3 for file uploads and backups'
    },
    {
      id: 'email',
      name: 'Email Service',
      status: 'operational',
      uptime: '99.94%',
      responseTime: 1200,
      requests: '89K',
      errors: 0.06,
      lastCheck: '2025-01-24T10:30:00Z',
      description: 'Transactional emails and notifications'
    },
    {
      id: 'search',
      name: 'Search Engine',
      status: 'operational',
      uptime: '99.87%',
      responseTime: 67,
      requests: '1.2M',
      errors: 0.13,
      lastCheck: '2025-01-24T10:30:00Z',
      description: 'Elasticsearch for product and tenant search'
    },
    {
      id: 'backup',
      name: 'Backup System',
      status: 'operational',
      uptime: '99.91%',
      responseTime: 0,
      requests: '24',
      errors: 0,
      lastCheck: '2025-01-24T10:30:00Z',
      description: 'Automated daily backups and disaster recovery'
    }
  ];

  // Infrastructure Metrics
  const infrastructureMetrics = [
    {
      name: 'CPU Usage',
      value: 68,
      threshold: 80,
      status: 'healthy',
      icon: Cpu,
      unit: '%',
      trend: 'up',
      change: '+5%'
    },
    {
      name: 'Memory Usage',
      value: 74,
      threshold: 85,
      status: 'healthy',
      icon: MemoryStick,
      unit: '%',
      trend: 'up',
      change: '+8%'
    },
    {
      name: 'Disk Usage',
      value: 45,
      threshold: 90,
      status: 'healthy',
      icon: HardDrive,
      unit: '%',
      trend: 'up',
      change: '+2%'
    },
    {
      name: 'Network I/O',
      value: 3.4,
      threshold: 10,
      status: 'healthy',
      icon: Network,
      unit: 'GB/s',
      trend: 'down',
      change: '-12%'
    },
    {
      name: 'Database Connections',
      value: 85,
      threshold: 200,
      status: 'healthy',
      icon: Database,
      unit: 'conn',
      trend: 'up',
      change: '+15%'
    },
    {
      name: 'Active Sessions',
      value: 1247,
      threshold: 5000,
      status: 'healthy',
      icon: Users,
      unit: 'sessions',
      trend: 'up',
      change: '+23%'
    }
  ];

  // Recent Incidents
  const recentIncidents = [
    {
      id: 'INC-001',
      title: 'Payment Gateway Intermittent Timeouts',
      status: 'investigating',
      severity: 'minor',
      startTime: '2025-01-24T09:15:00Z',
      duration: '1h 15m',
      affectedServices: ['payments'],
      description: 'Users experiencing timeout errors during bKash payments',
      updates: [
        {
          time: '2025-01-24T10:30:00Z',
          message: 'Engineering team investigating root cause',
          author: 'System'
        },
        {
          time: '2025-01-24T09:45:00Z',
          message: 'Issue confirmed with third-party payment provider',
          author: 'Tech Team'
        }
      ]
    },
    {
      id: 'INC-002',
      title: 'Database Maintenance Completed',
      status: 'resolved',
      severity: 'maintenance',
      startTime: '2025-01-23T02:00:00Z',
      duration: '2h 30m',
      affectedServices: ['database', 'api'],
      description: 'Scheduled database optimization and index rebuilding',
      updates: [
        {
          time: '2025-01-23T04:30:00Z',
          message: 'Maintenance completed successfully. All services restored.',
          author: 'DevOps Team'
        }
      ]
    }
  ];

  // Performance Metrics (mock data for charts)
  const performanceData = {
    responseTime: [
      { time: '09:00', api: 120, database: 15, payments: 1800 },
      { time: '09:30', api: 145, database: 23, payments: 2340 },
      { time: '10:00', api: 132, database: 19, payments: 2100 },
      { time: '10:30', api: 145, database: 23, payments: 2340 }
    ],
    throughput: [
      { time: '09:00', requests: 2100, errors: 5 },
      { time: '09:30', requests: 2400, errors: 8 },
      { time: '10:00', requests: 2200, errors: 3 },
      { time: '10:30', requests: 2400, errors: 12 }
    ]
  };

  const getStatusConfig = (status) => {
    const configs = {
      operational: {
        color: 'bg-green-100 text-green-800 border-green-200',
        icon: CheckCircle,
        label: 'Operational',
        badgeVariant: 'secondary'
      },
      warning: {
        color: 'bg-yellow-100 text-yellow-800 border-yellow-200',
        icon: AlertTriangle,
        label: 'Warning',
        badgeVariant: 'outline'
      },
      critical: {
        color: 'bg-red-100 text-red-800 border-red-200',
        icon: XCircle,
        label: 'Critical',
        badgeVariant: 'destructive'
      },
      maintenance: {
        color: 'bg-blue-100 text-blue-800 border-blue-200',
        icon: Settings,
        label: 'Maintenance',
        badgeVariant: 'secondary'
      },
      degraded: {
        color: 'bg-orange-100 text-orange-800 border-orange-200',
        icon: AlertCircle,
        label: 'Degraded',
        badgeVariant: 'outline'
      }
    };
    return configs[status] || configs.operational;
  };

  const getSeverityConfig = (severity) => {
    const configs = {
      minor: { color: 'text-yellow-600', label: 'Minor' },
      major: { color: 'text-orange-600', label: 'Major' },
      critical: { color: 'text-red-600', label: 'Critical' },
      maintenance: { color: 'text-blue-600', label: 'Maintenance' }
    };
    return configs[severity] || configs.minor;
  };

  const formatDuration = (startTime) => {
    const start = new Date(startTime);
    const now = new Date();
    const diff = now - start;
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
    return `${hours}h ${minutes}m`;
  };

  const handleRefresh = () => {
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
      setLastUpdated(new Date());
    }, 1500);
  };

  // Auto-refresh effect
  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(() => {
        setLastUpdated(new Date());
      }, 30000);
      return () => clearInterval(interval);
    }
  }, [autoRefresh]);

  const overallStatusConfig = getStatusConfig(systemStatus.overall);

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <Shield className="h-5 w-5 text-primary" />
              <h1 className="text-xl font-semibold">System Health</h1>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant={overallStatusConfig.badgeVariant} className="text-xs">
                <overallStatusConfig.icon className="h-3 w-3 mr-1" />
                {overallStatusConfig.label}
              </Badge>
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <Clock className="h-3 w-3" />
                Uptime: {systemStatus.uptime}
              </div>
            </div>
          </div>
          
          <div className="flex items-center gap-2">
            <div className="flex items-center gap-2 mr-4">
              <Label htmlFor="auto-refresh" className="text-xs">Auto-refresh</Label>
              <Switch
                id="auto-refresh"
                checked={autoRefresh}
                onCheckedChange={setAutoRefresh}
                size="sm"
              />
            </div>

            <div className="text-xs text-muted-foreground mr-4">
              Last updated: {lastUpdated.toLocaleTimeString('en-BD')}
            </div>
            
            <Button variant="outline" size="sm" onClick={handleRefresh} disabled={isLoading}>
              <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
            </Button>
            
            <Select value={selectedTimeRange} onValueChange={setSelectedTimeRange}>
              <SelectTrigger className="w-24">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="1h">1H</SelectItem>
                <SelectItem value="6h">6H</SelectItem>
                <SelectItem value="24h">24H</SelectItem>
                <SelectItem value="7d">7D</SelectItem>
                <SelectItem value="30d">30D</SelectItem>
              </SelectContent>
            </Select>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>System Actions</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <Download className="h-4 w-4 mr-2" />
                  Export Health Report
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Bell className="h-4 w-4 mr-2" />
                  Configure Alerts
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Settings className="h-4 w-4 mr-2" />
                  Health Settings
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      {/* System Status Overview */}
      <div className="border-b p-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Overall Status</p>
                  <p className="text-2xl font-bold capitalize">{systemStatus.overall}</p>
                  <div className="flex items-center gap-1 mt-1">
                    <div className="h-2 w-2 bg-green-500 rounded-full animate-pulse" />
                    <p className="text-xs text-green-600">All systems operational</p>
                  </div>
                </div>
                <overallStatusConfig.icon className="h-8 w-8 text-green-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">System Uptime</p>
                  <p className="text-2xl font-bold">{systemStatus.uptime}</p>
                  <p className="text-xs text-muted-foreground">Last 30 days</p>
                </div>
                <TrendingUp className="h-8 w-8 text-green-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Healthy Services</p>
                  <p className="text-2xl font-bold">
                    {systemStatus.healthyServices}/{systemStatus.totalServices}
                  </p>
                  <p className="text-xs text-green-600">
                    {Math.round((systemStatus.healthyServices / systemStatus.totalServices) * 100)}% operational
                  </p>
                </div>
                <CheckCircle className="h-8 w-8 text-green-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Active Incidents</p>
                  <p className="text-2xl font-bold">{recentIncidents.filter(i => i.status !== 'resolved').length}</p>
                  <p className="text-xs text-yellow-600">1 under investigation</p>
                </div>
                <AlertTriangle className="h-8 w-8 text-yellow-500" />
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Active Incidents Alert */}
        {recentIncidents.some(incident => incident.status !== 'resolved') && (
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4">
            <div className="flex items-start gap-3">
              <AlertTriangle className="h-5 w-5 text-yellow-600 mt-0.5" />
              <div className="flex-1">
                <h4 className="text-sm font-medium text-yellow-800 mb-1">Active Incident</h4>
                <p className="text-sm text-yellow-700">
                  {recentIncidents.find(i => i.status !== 'resolved')?.title} - 
                  <span className="ml-1 font-medium">
                    {formatDuration(recentIncidents.find(i => i.status !== 'resolved')?.startTime)}
                  </span>
                </p>
              </div>
              <Button variant="outline" size="sm" className="border-yellow-300">
                View Details
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Main Content Tabs */}
      <div className="flex-1 overflow-hidden">
        <Tabs value={selectedMetric} onValueChange={setSelectedMetric} className="h-full flex flex-col">
          <div className="border-b">
            <TabsList className="w-full justify-start h-12 p-1 ml-6">
              <TabsTrigger value="overview" className="flex items-center gap-2">
                <Activity className="h-4 w-4" />
                Services Overview
              </TabsTrigger>
              <TabsTrigger value="infrastructure" className="flex items-center gap-2">
                <Server className="h-4 w-4" />
                Infrastructure
              </TabsTrigger>
              <TabsTrigger value="performance" className="flex items-center gap-2">
                <BarChart3 className="h-4 w-4" />
                Performance
              </TabsTrigger>
              <TabsTrigger value="incidents" className="flex items-center gap-2">
                <AlertTriangle className="h-4 w-4" />
                Incidents
              </TabsTrigger>
            </TabsList>
          </div>

          <div className="flex-1 overflow-auto">
            {/* Services Overview Tab */}
            <TabsContent value="overview" className="h-full p-6 space-y-6">
              <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
                {coreServices.map((service) => {
                  const statusConfig = getStatusConfig(service.status);
                  return (
                    <Card key={service.id} className="relative">
                      <CardHeader className="pb-3">
                        <div className="flex items-center justify-between">
                          <CardTitle className="text-base font-medium">
                            {service.name}
                          </CardTitle>
                          <Badge variant={statusConfig.badgeVariant} className="text-xs">
                            <statusConfig.icon className="h-3 w-3 mr-1" />
                            {statusConfig.label}
                          </Badge>
                        </div>
                        <p className="text-xs text-muted-foreground">
                          {service.description}
                        </p>
                      </CardHeader>
                      <CardContent className="space-y-3">
                        <div className="grid grid-cols-2 gap-3 text-sm">
                          <div>
                            <p className="text-muted-foreground">Uptime</p>
                            <p className="font-medium">{service.uptime}</p>
                          </div>
                          <div>
                            <p className="text-muted-foreground">Response Time</p>
                            <p className="font-medium">{service.responseTime}ms</p>
                          </div>
                          <div>
                            <p className="text-muted-foreground">Requests</p>
                            <p className="font-medium">{service.requests}</p>
                          </div>
                          <div>
                            <p className="text-muted-foreground">Error Rate</p>
                            <p className="font-medium">{service.errors}%</p>
                          </div>
                        </div>

                        {service.issues && (
                          <div className="bg-yellow-50 border border-yellow-200 rounded p-2">
                            <p className="text-xs font-medium text-yellow-800 mb-1">Known Issues:</p>
                            <ul className="text-xs text-yellow-700 space-y-1">
                              {service.issues.map((issue, index) => (
                                <li key={index} className="flex items-start gap-1">
                                  <span className="w-1 h-1 bg-yellow-600 rounded-full mt-1.5 shrink-0" />
                                  {issue}
                                </li>
                              ))}
                            </ul>
                          </div>
                        )}

                        <div className="flex items-center justify-between pt-2 border-t">
                          <p className="text-xs text-muted-foreground">
                            Last check: {new Date(service.lastCheck).toLocaleTimeString('en-BD')}
                          </p>
                          <Button variant="ghost" size="sm" className="h-6 px-2 text-xs">
                            <Eye className="h-3 w-3 mr-1" />
                            Details
                          </Button>
                        </div>
                      </CardContent>
                    </Card>
                  );
                })}
              </div>
            </TabsContent>

            {/* Infrastructure Tab */}
            <TabsContent value="infrastructure" className="h-full p-6 space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {infrastructureMetrics.map((metric, index) => {
                  const isWarning = metric.value > metric.threshold * 0.8;
                  const isCritical = metric.value > metric.threshold;
                  
                  return (
                    <Card key={index}>
                      <CardContent className="p-4">
                        <div className="flex items-center justify-between mb-3">
                          <div className="flex items-center gap-2">
                            <metric.icon className="h-4 w-4 text-muted-foreground" />
                            <span className="text-sm font-medium">{metric.name}</span>
                          </div>
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            {metric.trend === 'up' ? (
                              <TrendingUp className="h-3 w-3 text-red-500" />
                            ) : (
                              <TrendingDown className="h-3 w-3 text-green-500" />
                            )}
                            {metric.change}
                          </div>
                        </div>

                        <div className="space-y-2">
                          <div className="flex items-center justify-between">
                            <span className="text-2xl font-bold">
                              {metric.value}{metric.unit}
                            </span>
                            <span className="text-xs text-muted-foreground">
                              / {metric.threshold}{metric.unit}
                            </span>
                          </div>
                          
                          <Progress 
                            value={(metric.value / metric.threshold) * 100} 
                            className="h-2"
                          />
                          
                          <div className="flex justify-between text-xs">
                            <span className={`${
                              isCritical ? 'text-red-600' : 
                              isWarning ? 'text-yellow-600' : 
                              'text-green-600'
                            }`}>
                              {isCritical ? 'Critical' : isWarning ? 'Warning' : 'Healthy'}
                            </span>
                            <span className="text-muted-foreground">
                              {Math.round((metric.value / metric.threshold) * 100)}% of limit
                            </span>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  );
                })}
              </div>

              {/* Server Details */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Server className="h-5 w-5" />
                    Server Details
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    <div className="space-y-3">
                      <h4 className="text-sm font-medium text-muted-foreground">Web Servers</h4>
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span>web-01.dhaka</span>
                          <Badge variant="secondary" className="text-xs">
                            <CheckCircle className="h-3 w-3 mr-1" />
                            Online
                          </Badge>
                        </div>
                        <div className="flex justify-between text-sm">
                          <span>web-02.dhaka</span>
                          <Badge variant="secondary" className="text-xs">
                            <CheckCircle className="h-3 w-3 mr-1" />
                            Online
                          </Badge>
                        </div>
                        <div className="flex justify-between text-sm">
                          <span>web-03.chittagong</span>
                          <Badge variant="secondary" className="text-xs">
                            <CheckCircle className="h-3 w-3 mr-1" />
                            Online
                          </Badge>
                        </div>
                      </div>
                    </div>

                    <div className="space-y-3">
                      <h4 className="text-sm font-medium text-muted-foreground">Database Servers</h4>
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span>db-primary</span>
                          <Badge variant="secondary" className="text-xs">
                            <CheckCircle className="h-3 w-3 mr-1" />
                            Primary
                          </Badge>
                        </div>
                        <div className="flex justify-between text-sm">
                          <span>db-replica-01</span>
                          <Badge variant="secondary" className="text-xs">
                            <CheckCircle className="h-3 w-3 mr-1" />
                            Replica
                          </Badge>
                        </div>
                      </div>
                    </div>

                    <div className="space-y-3">
                      <h4 className="text-sm font-medium text-muted-foreground">Cache Servers</h4>
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span>redis-01</span>
                          <Badge variant="secondary" className="text-xs">
                            <CheckCircle className="h-3 w-3 mr-1" />
                            Master
                          </Badge>
                        </div>
                        <div className="flex justify-between text-sm">
                          <span>redis-02</span>
                          <Badge variant="secondary" className="text-xs">
                            <CheckCircle className="h-3 w-3 mr-1" />
                            Slave
                          </Badge>
                        </div>
                      </div>
                    </div>

                    <div className="space-y-3">
                      <h4 className="text-sm font-medium text-muted-foreground">Load Balancers</h4>
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span>lb-primary</span>
                          <Badge variant="secondary" className="text-xs">
                            <CheckCircle className="h-3 w-3 mr-1" />
                            Active
                          </Badge>
                        </div>
                        <div className="flex justify-between text-sm">
                          <span>lb-backup</span>
                          <Badge variant="outline" className="text-xs">
                            <Clock className="h-3 w-3 mr-1" />
                            Standby
                          </Badge>
                        </div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            {/* Performance Tab */}
            <TabsContent value="performance" className="h-full p-6 space-y-6">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Response Time Chart */}
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Timer className="h-5 w-5" />
                      Response Times (Last {selectedTimeRange})
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="h-64 flex items-center justify-center bg-muted/20 rounded">
                      <div className="text-center">
                        <LineChart className="h-12 w-12 text-muted-foreground mx-auto mb-2" />
                        <p className="text-sm text-muted-foreground">
                          Response Time Chart Placeholder
                        </p>
                        <div className="mt-4 space-y-1">
                          <div className="flex items-center justify-between text-xs">
                            <span className="flex items-center gap-2">
                              <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
                              API Gateway
                            </span>
                            <span>145ms avg</span>
                          </div>
                          <div className="flex items-center justify-between text-xs">
                            <span className="flex items-center gap-2">
                              <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                              Database
                            </span>
                            <span>23ms avg</span>
                          </div>
                          <div className="flex items-center justify-between text-xs">
                            <span className="flex items-center gap-2">
                              <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                              Payments
                            </span>
                            <span>2.34s avg</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                {/* Throughput Chart */}
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <BarChart3 className="h-5 w-5" />
                      Request Throughput
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="h-64 flex items-center justify-center bg-muted/20 rounded">
                      <div className="text-center">
                        <BarChart3 className="h-12 w-12 text-muted-foreground mx-auto mb-2" />
                        <p className="text-sm text-muted-foreground">
                          Throughput Chart Placeholder
                        </p>
                        <div className="mt-4 grid grid-cols-2 gap-4 text-xs">
                          <div className="text-center">
                            <p className="text-2xl font-bold text-green-600">2.4K</p>
                            <p className="text-muted-foreground">Requests/min</p>
                          </div>
                          <div className="text-center">
                            <p className="text-2xl font-bold text-red-600">12</p>
                            <p className="text-muted-foreground">Errors/min</p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Performance Metrics Table */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Gauge className="h-5 w-5" />
                    Detailed Performance Metrics
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="overflow-x-auto">
                    <table className="w-full text-sm">
                      <thead>
                        <tr className="border-b">
                          <th className="text-left p-2">Service</th>
                          <th className="text-left p-2">Response Time</th>
                          <th className="text-left p-2">Throughput</th>
                          <th className="text-left p-2">Error Rate</th>
                          <th className="text-left p-2">Success Rate</th>
                          <th className="text-left p-2">Status</th>
                        </tr>
                      </thead>
                      <tbody>
                        {coreServices.map((service) => (
                          <tr key={service.id} className="border-b hover:bg-muted/50">
                            <td className="p-2 font-medium">{service.name}</td>
                            <td className="p-2">
                              <span className={`${
                                service.responseTime > 1000 ? 'text-red-600' :
                                service.responseTime > 500 ? 'text-yellow-600' :
                                'text-green-600'
                              }`}>
                                {service.responseTime}ms
                              </span>
                            </td>
                            <td className="p-2">{service.requests}</td>
                            <td className="p-2">
                              <span className={`${
                                service.errors > 0.5 ? 'text-red-600' :
                                service.errors > 0.1 ? 'text-yellow-600' :
                                'text-green-600'
                              }`}>
                                {service.errors}%
                              </span>
                            </td>
                            <td className="p-2 text-green-600">
                              {(100 - service.errors).toFixed(2)}%
                            </td>
                            <td className="p-2">
                              <Badge variant={getStatusConfig(service.status).badgeVariant} className="text-xs">
                                {getStatusConfig(service.status).label}
                              </Badge>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </CardContent>
              </Card>

              {/* SLA Metrics */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <Card>
                  <CardContent className="p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm text-muted-foreground">Monthly SLA</p>
                        <p className="text-2xl font-bold text-green-600">99.97%</p>
                        <p className="text-xs text-muted-foreground">Target: 99.90%</p>
                      </div>
                      <Target className="h-8 w-8 text-green-500" />
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardContent className="p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm text-muted-foreground">Avg Response Time</p>
                        <p className="text-2xl font-bold">187ms</p>
                        <p className="text-xs text-green-600">-12ms from last week</p>
                      </div>
                      <Timer className="h-8 w-8 text-blue-500" />
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardContent className="p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm text-muted-foreground">Error Budget</p>
                        <p className="text-2xl font-bold text-green-600">78%</p>
                        <p className="text-xs text-muted-foreground">Remaining this month</p>
                      </div>
                      <Gauge className="h-8 w-8 text-green-500" />
                    </div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>

            {/* Incidents Tab */}
            <TabsContent value="incidents" className="h-full p-6 space-y-6">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-medium">Recent Incidents</h3>
                <div className="flex items-center gap-2">
                  <Select defaultValue="all">
                    <SelectTrigger className="w-32">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Status</SelectItem>
                      <SelectItem value="investigating">Investigating</SelectItem>
                      <SelectItem value="resolved">Resolved</SelectItem>
                      <SelectItem value="monitoring">Monitoring</SelectItem>
                    </SelectContent>
                  </Select>
                  
                  <Select defaultValue="all">
                    <SelectTrigger className="w-32">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Severity</SelectItem>
                      <SelectItem value="critical">Critical</SelectItem>
                      <SelectItem value="major">Major</SelectItem>
                      <SelectItem value="minor">Minor</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="space-y-4">
                {recentIncidents.map((incident) => {
                  const severityConfig = getSeverityConfig(incident.severity);
                  const statusConfig = getStatusConfig(incident.status);
                  
                  return (
                    <Card key={incident.id}>
                      <CardContent className="p-6">
                        <div className="flex items-start justify-between mb-4">
                          <div className="flex-1">
                            <div className="flex items-center gap-3 mb-2">
                              <h4 className="font-medium">{incident.title}</h4>
                              <Badge variant={statusConfig.badgeVariant} className="text-xs">
                                {statusConfig.label}
                              </Badge>
                              <Badge variant="outline" className={`text-xs ${severityConfig.color}`}>
                                {severityConfig.label}
                              </Badge>
                            </div>
                            <p className="text-sm text-muted-foreground mb-3">
                              {incident.description}
                            </p>
                            <div className="flex items-center gap-4 text-xs text-muted-foreground">
                              <span>ID: {incident.id}</span>
                              <span>Started: {new Date(incident.startTime).toLocaleString('en-BD')}</span>
                              <span>Duration: {formatDuration(incident.startTime)}</span>
                              <span>Affected: {incident.affectedServices.join(', ')}</span>
                            </div>
                          </div>
                          <Button variant="outline" size="sm">
                            <ExternalLink className="h-4 w-4 mr-2" />
                            View Details
                          </Button>
                        </div>

                        {/* Incident Updates */}
                        <div className="border-t pt-4">
                          <h5 className="text-sm font-medium mb-3">Recent Updates</h5>
                          <div className="space-y-3">
                            {incident.updates.map((update, index) => (
                              <div key={index} className="flex gap-3">
                                <div className="w-2 h-2 bg-blue-500 rounded-full mt-2 shrink-0" />
                                <div className="flex-1">
                                  <div className="flex items-center gap-2 mb-1">
                                    <span className="text-xs text-muted-foreground">
                                      {new Date(update.time).toLocaleString('en-BD')}
                                    </span>
                                    <span className="text-xs text-muted-foreground">
                                      by {update.author}
                                    </span>
                                  </div>
                                  <p className="text-sm">{update.message}</p>
                                </div>
                              </div>
                            ))}
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  );
                })}
              </div>

              {/* Incident Stats */}
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mt-8">
                <Card>
                  <CardContent className="p-4 text-center">
                    <div className="text-2xl font-bold text-blue-600">24</div>
                    <div className="text-sm text-muted-foreground">Total Incidents</div>
                    <div className="text-xs text-green-600">-12% from last month</div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardContent className="p-4 text-center">
                    <div className="text-2xl font-bold text-green-600">2.4h</div>
                    <div className="text-sm text-muted-foreground">Avg Resolution Time</div>
                    <div className="text-xs text-green-600">-30min improvement</div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardContent className="p-4 text-center">
                    <div className="text-2xl font-bold text-yellow-600">1</div>
                    <div className="text-sm text-muted-foreground">Active Incidents</div>
                    <div className="text-xs text-muted-foreground">Under investigation</div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardContent className="p-4 text-center">
                    <div className="text-2xl font-bold text-red-600">0</div>
                    <div className="text-sm text-muted-foreground">Critical Incidents</div>
                    <div className="text-xs text-green-600">0 this month</div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>
          </div>
        </Tabs>
      </div>

      {/* Status Bar Footer */}
      <div className="border-t p-4 bg-background">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-6 text-sm">
            <div className="flex items-center gap-2">
              <div className="h-2 w-2 bg-green-500 rounded-full animate-pulse" />
              <span className="text-green-600">All Systems Operational</span>
            </div>
            <div className="flex items-center gap-2">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <span className="text-muted-foreground">
                Last incident: {new Date(systemStatus.lastIncident).toLocaleDateString('en-BD')}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <Activity className="h-4 w-4 text-muted-foreground" />
              <span className="text-muted-foreground">
                Monitoring {systemStatus.totalServices} services
              </span>
            </div>
          </div>
          
          <div className="flex items-center gap-4 text-xs text-muted-foreground">
            <span>BD Time: {new Date().toLocaleString('en-BD')}</span>
            <div className="flex items-center gap-1">
              <Signal className="h-3 w-3" />
              <span>Connected</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SystemHealthPage;