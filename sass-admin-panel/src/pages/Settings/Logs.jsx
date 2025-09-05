import React, { useState, useEffect } from 'react';
import {
  FileText,
  Search,
  Filter,
  Download,
  Eye,
  Calendar,
  Clock,
  User,
  Shield,
  AlertTriangle,
  CheckCircle,
  Info,
  X,
  ChevronDown,
  ChevronRight,
  MoreHorizontal,
  RefreshCw,
  Settings,
  Activity,
  Database,
  Globe,
  CreditCard,
  UserPlus,
  Store,
  Trash2,
  Edit,
  Key,
  ExternalLink,
  ArrowUpRight,
  ArrowDownRight,
  Minus,
  Bell,
  Copy,
  Users,
  Package
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

const AuditLogsPage = () => {
  const [selectedFilters, setSelectedFilters] = useState({
    dateRange: 'today',
    action: 'all',
    user: 'all',
    resource: 'all',
    severity: 'all'
  });
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedLog, setSelectedLog] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [expandedRows, setExpandedRows] = useState(new Set());
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);

  // Sample audit log data - more comprehensive
  const auditLogs = [
    {
      id: 'AL001',
      timestamp: '2025-01-24T10:30:00Z',
      user: { name: 'Admin Rahman', email: 'admin@shopowner.com', role: 'Super Admin', avatar: 'AR' },
      action: 'tenant_created',
      resource: 'tenant',
      resourceId: 'TNT-156',
      details: {
        tenantName: 'Rahman Electronics',
        plan: 'Business',
        initialUsers: 3,
        monthlyRevenue: 15000
      },
      severity: 'info',
      ipAddress: '192.168.1.100',
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
      location: 'Dhaka, Bangladesh',
      changes: {
        before: null,
        after: {
          name: 'Rahman Electronics',
          status: 'active',
          plan: 'business',
          created_at: '2025-01-24T10:30:00Z'
        }
      }
    },
    {
      id: 'AL002',
      timestamp: '2025-01-24T10:25:00Z',
      user: { name: 'Support Fatima', email: 'fatima@shopowner.com', role: 'Support', avatar: 'SF' },
      action: 'subscription_updated',
      resource: 'subscription',
      resourceId: 'SUB-789',
      details: {
        tenantName: 'Karim Store',
        oldPlan: 'Starter',
        newPlan: 'Business',
        priceChange: 3000
      },
      severity: 'warning',
      ipAddress: '192.168.1.101',
      userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)',
      location: 'Chittagong, Bangladesh',
      changes: {
        before: { plan: 'starter', price: 2000, max_products: 500 },
        after: { plan: 'business', price: 5000, max_products: 2000 }
      }
    },
    {
      id: 'AL003',
      timestamp: '2025-01-24T10:20:00Z',
      user: { name: 'System', email: null, role: 'System', avatar: 'SY' },
      action: 'payment_failed',
      resource: 'payment',
      resourceId: 'PAY-456',
      details: {
        tenantName: 'Ahmed Pharmacy',
        amount: 5000,
        paymentMethod: 'bKash',
        reason: 'Insufficient funds',
        attemptCount: 3
      },
      severity: 'error',
      ipAddress: 'system',
      userAgent: 'system',
      location: 'System',
      changes: {
        before: { status: 'pending', attempt: 2 },
        after: { status: 'failed', failureReason: 'insufficient_funds', attempt: 3 }
      }
    },
    {
      id: 'AL004',
      timestamp: '2025-01-24T10:15:00Z',
      user: { name: 'Admin Rahman', email: 'admin@shopowner.com', role: 'Super Admin', avatar: 'AR' },
      action: 'user_deleted',
      resource: 'user',
      resourceId: 'USR-123',
      details: {
        deletedUser: 'test@example.com',
        tenantName: 'Test Store',
        reason: 'Account closure request'
      },
      severity: 'critical',
      ipAddress: '192.168.1.100',
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',
      location: 'Dhaka, Bangladesh',
      changes: {
        before: { status: 'active', role: 'owner', last_login: '2025-01-20T08:00:00Z' },
        after: null
      }
    },
    {
      id: 'AL005',
      timestamp: '2025-01-24T10:10:00Z',
      user: { name: 'Admin Rahman', email: 'admin@shopowner.com', role: 'Super Admin', avatar: 'AR' },
      action: 'settings_updated',
      resource: 'platform_settings',
      resourceId: 'SET-001',
      details: {
        setting: 'max_products_per_plan',
        category: 'subscription_limits',
        changes: 'Updated Business plan limit from 1500 to 2000'
      },
      severity: 'info',
      ipAddress: '192.168.1.100',
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',
      location: 'Dhaka, Bangladesh',
      changes: {
        before: { business_max_products: 1500, enterprise_max_products: 'unlimited' },
        after: { business_max_products: 2000, enterprise_max_products: 'unlimited' }
      }
    },
    {
      id: 'AL006',
      timestamp: '2025-01-24T10:05:00Z',
      user: { name: 'Support Nasir', email: 'nasir@shopowner.com', role: 'Support', avatar: 'SN' },
      action: 'ticket_resolved',
      resource: 'support_ticket',
      resourceId: 'TIC-892',
      details: {
        ticketTitle: 'POS System Not Syncing',
        tenantName: 'Hasan General Store',
        resolutionTime: '2 hours 15 minutes',
        satisfaction: 5
      },
      severity: 'info',
      ipAddress: '192.168.1.102',
      userAgent: 'Mozilla/5.0 (X11; Linux x86_64)',
      location: 'Rajshahi, Bangladesh',
      changes: {
        before: { status: 'in_progress', assigned_to: 'nasir@shopowner.com' },
        after: { status: 'resolved', resolved_at: '2025-01-24T10:05:00Z', satisfaction_rating: 5 }
      }
    },
    {
      id: 'AL007',
      timestamp: '2025-01-24T09:55:00Z',
      user: { name: 'System', email: null, role: 'System', avatar: 'SY' },
      action: 'backup_completed',
      resource: 'database',
      resourceId: 'BCK-daily-240124',
      details: {
        backupType: 'Full Daily Backup',
        size: '2.4 GB',
        duration: '45 minutes',
        tables: 127
      },
      severity: 'info',
      ipAddress: 'system',
      userAgent: 'cron-job',
      location: 'System',
      changes: {
        before: { last_backup: '2025-01-23T09:55:00Z' },
        after: { last_backup: '2025-01-24T09:55:00Z', backup_size: '2.4GB' }
      }
    },
    {
      id: 'AL008',
      timestamp: '2025-01-24T09:45:00Z',
      user: { name: 'Merchant Rashid', email: 'rashid@rashhidstore.com', role: 'Tenant', avatar: 'MR' },
      action: 'login_suspicious',
      resource: 'authentication',
      resourceId: 'AUTH-991',
      details: {
        tenantName: 'Rashid Electronics',
        loginLocation: 'Sylhet, Bangladesh',
        deviceType: 'Mobile',
        flaggedReason: 'Login from new location'
      },
      severity: 'warning',
      ipAddress: '103.230.105.45',
      userAgent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)',
      location: 'Sylhet, Bangladesh',
      changes: {
        before: { last_location: 'Dhaka, Bangladesh', device_trusted: false },
        after: { last_location: 'Sylhet, Bangladesh', security_flag: 'location_change' }
      }
    }
  ];

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  
  const formatTimestamp = (timestamp) => {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMinutes = Math.floor((now - date) / (1000 * 60));
    
    if (diffMinutes < 1) return 'Just now';
    if (diffMinutes < 60) return `${diffMinutes}m ago`;
    if (diffMinutes < 1440) return `${Math.floor(diffMinutes / 60)}h ago`;
    
    return {
      date: date.toLocaleDateString('en-BD'),
      time: date.toLocaleTimeString('en-BD', { hour12: true }),
      relative: `${Math.floor(diffMinutes / 1440)}d ago`
    };
  };

  const getSeverityConfig = (severity) => {
    const configs = {
      info: { 
        color: 'bg-blue-100 text-blue-800 border-blue-200',
        icon: Info,
        label: 'Info',
        badgeVariant: 'secondary'
      },
      warning: { 
        color: 'bg-yellow-100 text-yellow-800 border-yellow-200',
        icon: AlertTriangle,
        label: 'Warning',
        badgeVariant: 'outline'
      },
      error: { 
        color: 'bg-red-100 text-red-800 border-red-200',
        icon: X,
        label: 'Error',
        badgeVariant: 'destructive'
      },
      critical: { 
        color: 'bg-red-100 text-red-900 border-red-300',
        icon: AlertTriangle,
        label: 'Critical',
        badgeVariant: 'destructive'
      }
    };
    return configs[severity] || configs.info;
  };

  const getActionConfig = (action) => {
    const configs = {
      tenant_created: { icon: Store, label: 'Tenant Created', color: 'text-green-600' },
      tenant_updated: { icon: Edit, label: 'Tenant Updated', color: 'text-blue-600' },
      tenant_deleted: { icon: Trash2, label: 'Tenant Deleted', color: 'text-red-600' },
      user_created: { icon: UserPlus, label: 'User Created', color: 'text-green-600' },
      user_updated: { icon: User, label: 'User Updated', color: 'text-blue-600' },
      user_deleted: { icon: Trash2, label: 'User Deleted', color: 'text-red-600' },
      subscription_updated: { icon: CreditCard, label: 'Subscription Updated', color: 'text-orange-600' },
      payment_failed: { icon: X, label: 'Payment Failed', color: 'text-red-600' },
      payment_success: { icon: CheckCircle, label: 'Payment Success', color: 'text-green-600' },
      settings_updated: { icon: Settings, label: 'Settings Updated', color: 'text-purple-600' },
      login: { icon: Key, label: 'Login', color: 'text-blue-600' },
      login_suspicious: { icon: Shield, label: 'Suspicious Login', color: 'text-red-600' },
      logout: { icon: Key, label: 'Logout', color: 'text-gray-600' },
      ticket_resolved: { icon: CheckCircle, label: 'Ticket Resolved', color: 'text-green-600' },
      backup_completed: { icon: Database, label: 'Backup Completed', color: 'text-blue-600' }
    };
    return configs[action] || { icon: Activity, label: action.replace('_', ' '), color: 'text-gray-600' };
  };

  const toggleRowExpansion = (logId) => {
    const newExpanded = new Set(expandedRows);
    if (newExpanded.has(logId)) {
      newExpanded.delete(logId);
    } else {
      newExpanded.add(logId);
    }
    setExpandedRows(newExpanded);
  };

  const handleExport = () => {
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
      // Create CSV content
      const csvContent = [
        'ID,Timestamp,User,Action,Resource,Severity,Details,IP Address',
        ...auditLogs.map(log => 
          `${log.id},${log.timestamp},${log.user.name},${log.action},${log.resource},${log.severity},"${Object.entries(log.details).map(([k,v]) => `${k}: ${v}`).join('; ')}",${log.ipAddress}`
        )
      ].join('\n');
      
      const blob = new Blob([csvContent], { type: 'text/csv' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `audit-logs-${new Date().toISOString().split('T')[0]}.csv`;
      a.click();
      window.URL.revokeObjectURL(url);
    }, 2000);
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
  };

  // Auto-refresh effect
  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(() => {
        // In real app, fetch new logs
        console.log('Auto-refreshing logs...');
      }, 30000);
      return () => clearInterval(interval);
    }
  }, [autoRefresh]);

  const summaryStats = [
    { 
      label: 'Total Entries Today', 
      value: '247', 
      change: '+12%', 
      icon: FileText, 
      color: 'text-blue-500',
      trend: 'up'
    },
    { 
      label: 'Critical Events', 
      value: '3', 
      change: '-50%', 
      icon: AlertTriangle, 
      color: 'text-red-500',
      trend: 'down'
    },
    { 
      label: 'Failed Actions', 
      value: '18', 
      change: '+25%', 
      icon: X, 
      color: 'text-orange-500',
      trend: 'up'
    },
    { 
      label: 'Active Users', 
      value: '24', 
      change: '+8%', 
      icon: Users, 
      color: 'text-green-500',
      trend: 'up'
    }
  ];

  const quickFilters = [
    { label: 'All Events', value: 'all', count: 247 },
    { label: 'User Actions', value: 'user_actions', count: 89 },
    { label: 'System Events', value: 'system', count: 156 },
    { label: 'Payments', value: 'payments', count: 34 },
    { label: 'Security', value: 'security', count: 12 }
  ];

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <FileText className="h-5 w-5 text-primary" />
              <h1 className="text-xl font-semibold">Audit Logs</h1>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="text-xs">
                Live Monitoring
              </Badge>
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <div className="h-2 w-2 bg-green-500 rounded-full animate-pulse" />
                Real-time
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
            
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="outline" size="sm" onClick={() => window.location.reload()}>
                    <RefreshCw className={`h-4 w-4 ${autoRefresh ? 'animate-spin' : ''}`} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Refresh logs</TooltipContent>
              </Tooltip>
            </TooltipProvider>
            
            <Button variant="outline" size="sm" onClick={handleExport} disabled={isLoading}>
              {isLoading ? (
                <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
              ) : (
                <Download className="h-4 w-4 mr-2" />
              )}
              Export CSV
            </Button>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Quick Actions</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <Settings className="h-4 w-4 mr-2" />
                  Log Settings
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Shield className="h-4 w-4 mr-2" />
                  Security Report
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Database className="h-4 w-4 mr-2" />
                  Archive Logs
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      {/* Summary Stats */}
      <div className="border-b p-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
          {summaryStats.map((stat, index) => (
            <Card key={index}>
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-muted-foreground">{stat.label}</p>
                    <p className="text-2xl font-bold">{stat.value}</p>
                    <div className="flex items-center gap-1 mt-1">
                      {stat.trend === 'up' ? (
                        <ArrowUpRight className="h-3 w-3 text-green-500" />
                      ) : (
                        <ArrowDownRight className="h-3 w-3 text-red-500" />
                      )}
                      <p className={`text-xs ${stat.change.startsWith('+') ? 'text-green-600' : 'text-red-600'}`}>
                        {stat.change} from yesterday
                      </p>
                    </div>
                  </div>
                  <stat.icon className={`h-8 w-8 ${stat.color}`} />
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Quick Filter Tabs */}
        <Tabs defaultValue="all" className="w-full">
          <TabsList className="grid w-full grid-cols-5">
            {quickFilters.map((filter) => (
              <TabsTrigger key={filter.value} value={filter.value} className="text-xs">
                {filter.label}
                <Badge variant="secondary" className="ml-2 text-xs">
                  {filter.count}
                </Badge>
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
      </div>

      {/* Filters */}
      <div className="border-b p-6">
        <div className="flex flex-col lg:flex-row gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search logs by user, action, resource, or details..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>
          
          <div className="flex flex-wrap gap-2">
            <Select value={selectedFilters.dateRange} onValueChange={(value) => 
              setSelectedFilters(prev => ({ ...prev, dateRange: value }))
            }>
              <SelectTrigger className="w-36">
                <Calendar className="h-4 w-4 mr-2" />
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="today">Today</SelectItem>
                <SelectItem value="yesterday">Yesterday</SelectItem>
                <SelectItem value="week">This Week</SelectItem>
                <SelectItem value="month">This Month</SelectItem>
                <SelectItem value="quarter">This Quarter</SelectItem>
                <SelectItem value="custom">Custom Range</SelectItem>
              </SelectContent>
            </Select>

            <Select value={selectedFilters.action} onValueChange={(value) => 
              setSelectedFilters(prev => ({ ...prev, action: value }))
            }>
              <SelectTrigger className="w-44">
                <Activity className="h-4 w-4 mr-2" />
                <SelectValue placeholder="All Actions" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Actions</SelectItem>
                <SelectItem value="tenant_created">Tenant Created</SelectItem>
                <SelectItem value="user_created">User Created</SelectItem>
                <SelectItem value="payment_failed">Payment Failed</SelectItem>
                <SelectItem value="settings_updated">Settings Updated</SelectItem>
                <SelectItem value="login_suspicious">Suspicious Login</SelectItem>
              </SelectContent>
            </Select>

            <Select value={selectedFilters.severity} onValueChange={(value) => 
              setSelectedFilters(prev => ({ ...prev, severity: value }))
            }>
              <SelectTrigger className="w-36">
                <Filter className="h-4 w-4 mr-2" />
                <SelectValue placeholder="All Levels" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Levels</SelectItem>
                <SelectItem value="info">Info</SelectItem>
                <SelectItem value="warning">Warning</SelectItem>
                <SelectItem value="error">Error</SelectItem>
                <SelectItem value="critical">Critical</SelectItem>
              </SelectContent>
            </Select>

            <Select value={selectedFilters.user} onValueChange={(value) => 
              setSelectedFilters(prev => ({ ...prev, user: value }))
            }>
              <SelectTrigger className="w-36">
                <User className="h-4 w-4 mr-2" />
                <SelectValue placeholder="All Users" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Users</SelectItem>
                <SelectItem value="admin">Admins Only</SelectItem>
                <SelectItem value="support">Support Only</SelectItem>
                <SelectItem value="system">System Only</SelectItem>
                <SelectItem value="tenant">Tenants Only</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </div>

      {/* Audit Logs Table */}
      <div className="flex-1 overflow-hidden">
        <div className="h-full overflow-auto">
          <Table>
            <TableHeader className="sticky top-0 bg-background border-b">
              <TableRow>
                <TableHead className="w-10"></TableHead>
                <TableHead className="w-32">Time</TableHead>
                <TableHead className="w-44">User</TableHead>
                <TableHead className="w-44">Action</TableHead>
                <TableHead className="w-36">Resource</TableHead>
                <TableHead className="w-24">Severity</TableHead>
                <TableHead>Details</TableHead>
                <TableHead className="w-28">Location</TableHead>
                <TableHead className="w-20"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {auditLogs.map((log) => {
                const isExpanded = expandedRows.has(log.id);
                const severityConfig = getSeverityConfig(log.severity);
                const actionConfig = getActionConfig(log.action);
                const timestamp = formatTimestamp(log.timestamp);
                
                return (
                  <React.Fragment key={log.id}>
                    <TableRow className="hover:bg-muted/50">
                      <TableCell>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-6 w-6 p-0"
                          onClick={() => toggleRowExpansion(log.id)}
                        >
                          {isExpanded ? (
                            <ChevronDown className="h-3 w-3" />
                          ) : (
                            <ChevronRight className="h-3 w-3" />
                          )}
                        </Button>
                      </TableCell>
                      
                      <TableCell>
                        <div className="text-sm">
                          <div className="font-medium">
                            {typeof timestamp === 'string' ? timestamp : timestamp.relative}
                          </div>
                          {typeof timestamp === 'object' && (
                            <div className="text-muted-foreground text-xs">
                              {timestamp.time}
                            </div>
                          )}
                        </div>
                      </TableCell>
                      
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <Avatar className="h-7 w-7">
                            <AvatarFallback className="text-xs bg-primary/10">
                              {log.user.avatar}
                            </AvatarFallback>
                          </Avatar>
                          <div className="text-sm min-w-0">
                            <div className="font-medium truncate">{log.user.name || 'System'}</div>
                            <div className="text-muted-foreground text-xs truncate">{log.user.role}</div>
                          </div>
                        </div>
                      </TableCell>
                      
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <actionConfig.icon className={`h-4 w-4 ${actionConfig.color} shrink-0`} />
                          <span className="text-sm font-medium truncate">{actionConfig.label}</span>
                        </div>
                      </TableCell>
                      
                      <TableCell>
                        <div className="text-sm">
                          <div className="font-medium capitalize">{log.resource.replace('_', ' ')}</div>
                          <div className="text-muted-foreground text-xs font-mono">{log.resourceId}</div>
                        </div>
                      </TableCell>
                      
                      <TableCell>
                        <Badge variant={severityConfig.badgeVariant} className="text-xs">
                          <severityConfig.icon className="h-3 w-3 mr-1" />
                          {severityConfig.label}
                        </Badge>
                      </TableCell>
                      
                      <TableCell>
                        <div className="text-sm text-muted-foreground max-w-80 truncate">
                          {Object.entries(log.details).slice(0, 2).map(([key, value]) => 
                            `${key.replace('_', ' ')}: ${value}`
                          ).join(', ')}
                          {Object.keys(log.details).length > 2 && '...'}
                        </div>
                      </TableCell>

                      <TableCell>
                        <div className="text-xs text-muted-foreground">
                          <div className="truncate">{log.location}</div>
                          <div className="font-mono">{log.ipAddress !== 'system' ? log.ipAddress.split('.').slice(0, 2).join('.') + '.***' : 'System'}</div>
                        </div>
                      </TableCell>
                      
                      <TableCell>
                        <div className="flex items-center gap-1">
                          <Dialog>
                            <DialogTrigger asChild>
                              <Button variant="ghost" size="sm" className="h-7 w-7 p-0">
                                <Eye className="h-3 w-3" />
                              </Button>
                            </DialogTrigger>
                            <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
                              <DialogHeader>
                                <DialogTitle className="flex items-center gap-2">
                                  <actionConfig.icon className={`h-5 w-5 ${actionConfig.color}`} />
                                  Audit Log Details - {log.id}
                                </DialogTitle>
                                <DialogDescription>
                                  Complete information about this audit log entry
                                </DialogDescription>
                              </DialogHeader>
                              <div className="space-y-6">
                                {/* Basic Information */}
                                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                                  <div>
                                    <label className="text-sm font-medium text-muted-foreground">Log ID</label>
                                    <div className="flex items-center gap-2">
                                      <p className="text-sm font-mono">{log.id}</p>
                                      <Button
                                        variant="ghost"
                                        size="sm"
                                        className="h-5 w-5 p-0"
                                        onClick={() => copyToClipboard(log.id)}
                                      >
                                        <Copy className="h-3 w-3" />
                                      </Button>
                                    </div>
                                  </div>
                                  <div>
                                    <label className="text-sm font-medium text-muted-foreground">Timestamp</label>
                                    <p className="text-sm">{new Date(log.timestamp).toLocaleString('en-BD')}</p>
                                  </div>
                                  <div>
                                    <label className="text-sm font-medium text-muted-foreground">IP Address</label>
                                    <p className="text-sm font-mono">{log.ipAddress}</p>
                                  </div>
                                  <div>
                                    <label className="text-sm font-medium text-muted-foreground">Severity</label>
                                    <Badge variant={severityConfig.badgeVariant} className="text-xs">
                                      <severityConfig.icon className="h-3 w-3 mr-1" />
                                      {severityConfig.label}
                                    </Badge>
                                  </div>
                                </div>
                                
                                <Separator />
                                
                                {/* User Information */}
                                <div>
                                  <h4 className="text-sm font-medium mb-3">User Information</h4>
                                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                    <div>
                                      <label className="text-sm font-medium text-muted-foreground">Name</label>
                                      <p className="text-sm">{log.user.name || 'System'}</p>
                                    </div>
                                    <div>
                                      <label className="text-sm font-medium text-muted-foreground">Email</label>
                                      <p className="text-sm">{log.user.email || 'N/A'}</p>
                                    </div>
                                    <div>
                                      <label className="text-sm font-medium text-muted-foreground">Role</label>
                                      <p className="text-sm">{log.user.role}</p>
                                    </div>
                                  </div>
                                </div>

                                <Separator />

                                {/* Action Details */}
                                <div>
                                  <h4 className="text-sm font-medium mb-3">Action Details</h4>
                                  <div className="bg-muted/50 rounded-lg p-4">
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                                      <div>
                                        <label className="text-sm font-medium text-muted-foreground">Action</label>
                                        <p className="text-sm font-medium">{actionConfig.label}</p>
                                      </div>
                                      <div>
                                        <label className="text-sm font-medium text-muted-foreground">Resource</label>
                                        <p className="text-sm">{log.resource} ({log.resourceId})</p>
                                      </div>
                                    </div>
                                    <div>
                                      <label className="text-sm font-medium text-muted-foreground">Details</label>
                                      <div className="mt-2 space-y-2">
                                        {Object.entries(log.details).map(([key, value]) => (
                                          <div key={key} className="flex justify-between items-center py-1 border-b border-border/50 last:border-0">
                                            <span className="text-sm text-muted-foreground capitalize">
                                              {key.replace('_', ' ')}
                                            </span>
                                            <span className="text-sm font-mono">
                                              {typeof value === 'number' && key.includes('amount') || key.includes('revenue') || key.includes('price') 
                                                ? formatCurrency(value) 
                                                : String(value)
                                              }
                                            </span>
                                          </div>
                                        ))}
                                      </div>
                                    </div>
                                  </div>
                                </div>
                                
                                <Separator />
                                
                                {/* Technical Details */}
                                <div>
                                  <h4 className="text-sm font-medium mb-3">Technical Information</h4>
                                  <div className="space-y-3">
                                    <div>
                                      <label className="text-sm font-medium text-muted-foreground">User Agent</label>
                                      <p className="text-sm font-mono text-muted-foreground break-all bg-muted/50 p-2 rounded">
                                        {log.userAgent}
                                      </p>
                                    </div>
                                    <div>
                                      <label className="text-sm font-medium text-muted-foreground">Location</label>
                                      <p className="text-sm">{log.location}</p>
                                    </div>
                                  </div>
                                </div>

                                <Separator />

                                {/* Changes */}
                                <div>
                                  <h4 className="text-sm font-medium mb-3">Changes Made</h4>
                                  <div className="space-y-3">
                                    {log.changes.before && (
                                      <div className="bg-red-50 border border-red-200 rounded-lg p-3">
                                        <div className="flex items-center gap-2 mb-2">
                                          <ArrowDownRight className="h-4 w-4 text-red-600" />
                                          <p className="text-sm font-medium text-red-800">Before</p>
                                        </div>
                                        <pre className="text-xs text-red-700 overflow-auto bg-red-100/50 p-2 rounded">
                                          {JSON.stringify(log.changes.before, null, 2)}
                                        </pre>
                                      </div>
                                    )}
                                    {log.changes.after && (
                                      <div className="bg-green-50 border border-green-200 rounded-lg p-3">
                                        <div className="flex items-center gap-2 mb-2">
                                          <ArrowUpRight className="h-4 w-4 text-green-600" />
                                          <p className="text-sm font-medium text-green-800">After</p>
                                        </div>
                                        <pre className="text-xs text-green-700 overflow-auto bg-green-100/50 p-2 rounded">
                                          {JSON.stringify(log.changes.after, null, 2)}
                                        </pre>
                                      </div>
                                    )}
                                    {!log.changes.before && !log.changes.after && (
                                      <div className="bg-muted/50 border rounded-lg p-3 text-center">
                                        <p className="text-sm text-muted-foreground">No changes recorded</p>
                                      </div>
                                    )}
                                  </div>
                                </div>
                              </div>
                            </DialogContent>
                          </Dialog>

                          <TooltipProvider>
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  className="h-7 w-7 p-0"
                                  onClick={() => copyToClipboard(log.id)}
                                >
                                  <Copy className="h-3 w-3" />
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>Copy Log ID</TooltipContent>
                            </Tooltip>
                          </TooltipProvider>
                        </div>
                      </TableCell>
                    </TableRow>
                    
                    {/* Expanded Row Details */}
                    {isExpanded && (
                      <TableRow>
                        <TableCell colSpan={9} className="bg-muted/20 p-0">
                          <div className="p-6 space-y-4">
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-sm">
                              <div className="space-y-2">
                                <label className="font-medium text-muted-foreground">Technical Info</label>
                                <div className="space-y-1">
                                  <div className="flex justify-between">
                                    <span className="text-muted-foreground">IP Address:</span>
                                    <span className="font-mono text-xs">{log.ipAddress}</span>
                                  </div>
                                  <div className="flex justify-between">
                                    <span className="text-muted-foreground">Location:</span>
                                    <span className="text-xs">{log.location}</span>
                                  </div>
                                  <div className="flex justify-between">
                                    <span className="text-muted-foreground">Resource ID:</span>
                                    <span className="font-mono text-xs">{log.resourceId}</span>
                                  </div>
                                </div>
                              </div>

                              <div className="space-y-2">
                                <label className="font-medium text-muted-foreground">User Details</label>
                                <div className="space-y-1">
                                  <div className="flex justify-between">
                                    <span className="text-muted-foreground">Email:</span>
                                    <span className="text-xs">{log.user.email || 'N/A'}</span>
                                  </div>
                                  <div className="flex justify-between">
                                    <span className="text-muted-foreground">Role:</span>
                                    <span className="text-xs">{log.user.role}</span>
                                  </div>
                                </div>
                              </div>

                              <div className="space-y-2">
                                <label className="font-medium text-muted-foreground">Action Context</label>
                                <div className="space-y-1">
                                  {Object.entries(log.details).slice(0, 3).map(([key, value]) => (
                                    <div key={key} className="flex justify-between">
                                      <span className="text-muted-foreground capitalize">
                                        {key.replace('_', ' ')}:
                                      </span>
                                      <span className="text-xs font-mono">
                                        {typeof value === 'number' && (key.includes('amount') || key.includes('revenue') || key.includes('price'))
                                          ? formatCurrency(value)
                                          : String(value).length > 20 
                                            ? String(value).substring(0, 20) + '...'
                                            : String(value)
                                        }
                                      </span>
                                    </div>
                                  ))}
                                </div>
                              </div>

                              <div className="space-y-2">
                                <label className="font-medium text-muted-foreground">Quick Actions</label>
                                <div className="flex flex-col gap-2">
                                  <Button variant="outline" size="sm" className="h-7 text-xs">
                                    <ExternalLink className="h-3 w-3 mr-1" />
                                    View Resource
                                  </Button>
                                  <Button variant="outline" size="sm" className="h-7 text-xs">
                                    <User className="h-3 w-3 mr-1" />
                                    View User
                                  </Button>
                                  {log.severity === 'error' || log.severity === 'critical' ? (
                                    <Button variant="outline" size="sm" className="h-7 text-xs text-red-600">
                                      <AlertTriangle className="h-3 w-3 mr-1" />
                                      Investigate
                                    </Button>
                                  ) : null}
                                </div>
                              </div>
                            </div>

                            {/* User Agent */}
                            <div className="space-y-2">
                              <label className="font-medium text-muted-foreground">User Agent</label>
                              <div className="bg-background rounded border p-3">
                                <p className="text-xs font-mono text-muted-foreground break-all">
                                  {log.userAgent}
                                </p>
                              </div>
                            </div>

                            {/* Changes Summary */}
                            <div className="space-y-2">
                              <label className="font-medium text-muted-foreground">Changes Summary</label>
                              <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                {log.changes.before && (
                                  <div className="bg-red-50 border border-red-200 rounded p-3">
                                    <p className="text-xs font-medium text-red-800 mb-1 flex items-center gap-1">
                                      <ArrowDownRight className="h-3 w-3" />
                                      Before
                                    </p>
                                    <div className="space-y-1">
                                      {Object.entries(log.changes.before).map(([key, value]) => (
                                        <div key={key} className="flex justify-between text-xs">
                                          <span className="text-red-700">{key}:</span>
                                          <span className="text-red-800 font-mono">
                                            {typeof value === 'number' && (key.includes('price') || key.includes('amount'))
                                              ? formatCurrency(value)
                                              : String(value)
                                            }
                                          </span>
                                        </div>
                                      ))}
                                    </div>
                                  </div>
                                )}
                                {log.changes.after && (
                                  <div className="bg-green-50 border border-green-200 rounded p-3">
                                    <p className="text-xs font-medium text-green-800 mb-1 flex items-center gap-1">
                                      <ArrowUpRight className="h-3 w-3" />
                                      After
                                    </p>
                                    <div className="space-y-1">
                                      {Object.entries(log.changes.after).map(([key, value]) => (
                                        <div key={key} className="flex justify-between text-xs">
                                          <span className="text-green-700">{key}:</span>
                                          <span className="text-green-800 font-mono">
                                            {typeof value === 'number' && (key.includes('price') || key.includes('amount'))
                                              ? formatCurrency(value)
                                              : String(value)
                                            }
                                          </span>
                                        </div>
                                      ))}
                                    </div>
                                  </div>
                                )}
                              </div>
                            </div>
                          </div>
                        </TableCell>
                      </TableRow>
                    )}
                  </React.Fragment>
                );
              })}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* Pagination Footer */}
      <div className="border-t p-4 bg-background">
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="text-sm text-muted-foreground">
            Showing <span className="font-medium">1-8</span> of <span className="font-medium">247</span> entries
            {searchQuery && (
              <span className="ml-2">
                • Filtered by "<span className="font-medium">{searchQuery}</span>"
              </span>
            )}
          </div>
          
          <div className="flex items-center gap-2">
            <div className="flex items-center gap-1 mr-4">
              <label className="text-sm text-muted-foreground">Rows per page:</label>
              <Select defaultValue="10">
                <SelectTrigger className="w-16 h-8">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="10">10</SelectItem>
                  <SelectItem value="25">25</SelectItem>
                  <SelectItem value="50">50</SelectItem>
                  <SelectItem value="100">100</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <Button variant="outline" size="sm" disabled>
              Previous
            </Button>
            <div className="flex items-center gap-1">
              {[1, 2, 3, '...', 25].map((page, index) => (
                <Button
                  key={index}
                  variant={page === 1 ? "default" : "ghost"}
                  size="sm"
                  className="w-8 h-8 p-0"
                  disabled={page === '...'}
                  onClick={() => page !== '...' && setCurrentPage(page)}
                >
                  {page}
                </Button>
              ))}
            </div>
            <Button variant="outline" size="sm">
              Next
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AuditLogsPage;