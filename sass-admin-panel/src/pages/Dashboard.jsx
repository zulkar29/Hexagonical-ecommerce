import React, { useState, useEffect } from 'react';
import {
  Users,
  Store,
  DollarSign,
  TrendingUp,
  TrendingDown,
  AlertTriangle,
  CheckCircle,
  Clock,
  CreditCard,
  X,
  Calendar,
  BarChart3,
  PieChart,
  Settings,
  Bell,
  Search,
  Filter,
  Download,
  RefreshCw,
  Eye,
  Edit,
  Trash2,
  Plus,
  MoreHorizontal,
  Phone,
  Mail,
  MapPin,
  Package,
  UserPlus,
  Shield,
  AlertCircle
} from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
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
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { LineChart, Line, AreaChart, Area, BarChart, Bar, PieChart as RechartsPieChart, Cell, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend, Pie } from 'recharts';

const SaaSAdminDashboard = () => {
  const [currentTime, setCurrentTime] = useState(new Date());
  const [selectedPeriod, setSelectedPeriod] = useState('30d');
  const [searchQuery, setSearchQuery] = useState('');
  const [showAddTenantDialog, setShowAddTenantDialog] = useState(false);
  const [showMessageDialog, setShowMessageDialog] = useState(false);
  const [selectedTenant, setSelectedTenant] = useState(null);

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  // Mock data for dashboard
  const stats = [
    {
      id: 1,
      title: 'Total Revenue',
      value: '৳2,45,000',
      change: '+12.5%',
      trend: 'up',
      icon: DollarSign,
      description: 'Monthly recurring revenue',
      color: 'text-green-600'
    },
    {
      id: 2,
      title: 'Active Tenants',
      value: '156',
      change: '+8.2%',
      trend: 'up',
      icon: Store,
      description: 'Currently subscribed shops',
      color: 'text-blue-600'
    },
    {
      id: 3,
      title: 'New Signups',
      value: '24',
      change: '+15.3%',
      trend: 'up',
      icon: UserPlus,
      description: 'This month',
      color: 'text-purple-600'
    },
    {
      id: 4,
      title: 'Churn Rate',
      value: '2.8%',
      change: '-0.5%',
      trend: 'up',
      icon: TrendingDown,
      description: 'Monthly churn rate',
      color: 'text-orange-600'
    }
  ];

  const quickActions = [
    {
      id: 1,
      title: 'Add New Tenant',
      description: 'Register a new shop',
      icon: Plus,
      color: 'bg-green-50 text-green-600 hover:bg-green-100'
    },
    {
      id: 2,
      title: 'View Reports',
      description: 'Generate analytics',
      icon: BarChart3,
      color: 'bg-blue-50 text-blue-600 hover:bg-blue-100'
    },
    {
      id: 3,
      title: 'Support Queue',
      description: '5 pending tickets',
      icon: AlertCircle,
      color: 'bg-red-50 text-red-600 hover:bg-red-100'
    },
    {
      id: 4,
      title: 'System Health',
      description: 'All systems operational',
      icon: Shield,
      color: 'bg-green-50 text-green-600 hover:bg-green-100'
    }
  ];

  const revenueData = [
    { month: 'Jan', revenue: 180000, tenants: 120, newSignups: 15 },
    { month: 'Feb', revenue: 195000, tenants: 128, newSignups: 18 },
    { month: 'Mar', revenue: 210000, tenants: 135, newSignups: 22 },
    { month: 'Apr', revenue: 225000, tenants: 142, newSignups: 20 },
    { month: 'May', revenue: 240000, tenants: 148, newSignups: 25 },
    { month: 'Jun', revenue: 245000, tenants: 156, newSignups: 24 }
  ];

  const planDistribution = [
    { name: 'Starter', value: 65, count: 89, color: '#22c55e', revenue: 178000 },
    { name: 'Business', value: 25, count: 45, color: '#3b82f6', revenue: 225000 },
    { name: 'Enterprise', value: 10, count: 22, color: '#f59e0b', revenue: 220000 }
  ];

  const topPerformingTenants = [
    { name: 'Rahman Electronics', revenue: 15000, growth: '+25%' },
    { name: 'Modern Pharmacy', revenue: 12000, growth: '+18%' },
    { name: 'Tech Solutions', revenue: 10000, growth: '+22%' },
    { name: 'Green Grocers', revenue: 8500, growth: '+15%' },
    { name: 'Fashion Hub', revenue: 7200, growth: '+12%' }
  ];

  const recentTenants = [
    {
      id: 1,
      name: 'Rahman Electronics',
      owner: 'Abdul Rahman',
      email: 'rahman@email.com',
      phone: '01712345678',
      plan: 'Business',
      status: 'active',
      joinDate: '2024-07-20',
      revenue: 5000,
      location: 'Dhanmondi, Dhaka',
      lastLogin: '2024-07-24',
      productsCount: 245
    },
    {
      id: 2,
      name: 'Fatima Fashion',
      owner: 'Fatima Khatun',
      email: 'fatima@email.com',
      phone: '01798765432',
      plan: 'Starter',
      status: 'trial',
      joinDate: '2024-07-22',
      revenue: 2000,
      location: 'Gulshan, Dhaka',
      lastLogin: '2024-07-23',
      productsCount: 89
    },
    {
      id: 3,
      name: 'Modern Pharmacy',
      owner: 'Dr. Ahmed Ali',
      email: 'ahmed@email.com',
      phone: '01656789012',
      plan: 'Enterprise',
      status: 'active',
      joinDate: '2024-07-18',
      revenue: 10000,
      location: 'Uttara, Dhaka',
      lastLogin: '2024-07-24',
      productsCount: 567
    },
    {
      id: 4,
      name: 'Green Grocers',
      owner: 'Mohammad Hasan',
      email: 'hasan@email.com',
      phone: '01534567890',
      plan: 'Business',
      status: 'suspended',
      joinDate: '2024-07-15',
      revenue: 5000,
      location: 'Mirpur, Dhaka',
      lastLogin: '2024-07-20',
      productsCount: 156
    },
    {
      id: 5,
      name: 'Tech Solutions',
      owner: 'Rashida Begum',
      email: 'rashida@email.com',
      phone: '01445678901',
      plan: 'Enterprise',
      status: 'active',
      joinDate: '2024-07-21',
      revenue: 10000,
      location: 'Banani, Dhaka',
      lastLogin: '2024-07-24',
      productsCount: 423
    }
  ];

  const recentPayments = [
    {
      id: 1,
      tenant: 'Rahman Electronics',
      amount: 5000,
      method: 'bKash',
      status: 'completed',
      date: '2024-07-24',
      plan: 'Business',
      transactionId: 'BKS789123456',
      invoiceNumber: 'INV-2024-001'
    },
    {
      id: 2,
      tenant: 'Modern Pharmacy',
      amount: 10000,
      method: 'Bank Transfer',
      status: 'completed',
      date: '2024-07-23',
      plan: 'Enterprise',
      transactionId: 'BT456789123',
      invoiceNumber: 'INV-2024-002'
    },
    {
      id: 3,
      tenant: 'Fatima Fashion',
      amount: 2000,
      method: 'Nagad',
      status: 'pending',
      date: '2024-07-22',
      plan: 'Starter',
      transactionId: 'NGD123456789',
      invoiceNumber: 'INV-2024-003'
    },
    {
      id: 4,
      tenant: 'Green Grocers',
      amount: 5000,
      method: 'bKash',
      status: 'failed',
      date: '2024-07-21',
      plan: 'Business',
      transactionId: 'BKS987654321',
      invoiceNumber: 'INV-2024-004'
    },
    {
      id: 5,
      tenant: 'Tech Solutions',
      amount: 10000,
      method: 'Rocket',
      status: 'completed',
      date: '2024-07-20',
      plan: 'Enterprise',
      transactionId: 'RKT555444333',
      invoiceNumber: 'INV-2024-005'
    }
  ];

  const supportTickets = [
    {
      id: 1,
      tenant: 'Rahman Electronics',
      subject: 'POS System Not Working',
      priority: 'high',
      status: 'open',
      assignee: 'Support Team A',
      created: '2024-07-24',
      category: 'Technical',
      description: 'Customer unable to process sales, POS system freezing'
    },
    {
      id: 2,
      tenant: 'Fatima Fashion',
      subject: 'Inventory Sync Issue',
      priority: 'medium',
      status: 'in-progress',
      assignee: 'Support Team B',
      created: '2024-07-23',
      category: 'Bug Report',
      description: 'Product quantities not updating correctly'
    },
    {
      id: 3,
      tenant: 'Modern Pharmacy',
      subject: 'Report Generation Problem',
      priority: 'low',
      status: 'resolved',
      assignee: 'Support Team A',
      created: '2024-07-22',
      category: 'Feature Request',
      description: 'Need custom report format for pharmacy regulations'
    },
    {
      id: 4,
      tenant: 'Tech Solutions',
      subject: 'Payment Gateway Integration',
      priority: 'high',
      status: 'open',
      assignee: 'Support Team C',
      created: '2024-07-21',
      category: 'Integration',
      description: 'Need to integrate additional payment methods'
    },
    {
      id: 5,
      tenant: 'Green Grocers',
      subject: 'Mobile App Access Issue',
      priority: 'medium',
      status: 'in-progress',
      assignee: 'Support Team B',
      created: '2024-07-20',
      category: 'Access',
      description: 'Unable to access mobile version of the dashboard'
    }
  ];

  const systemAlerts = [
    {
      id: 1,
      type: 'warning',
      title: 'Server Load High',
      message: 'Database server at 85% capacity',
      time: '2 hours ago'
    },
    {
      id: 2,
      type: 'info',
      title: 'Scheduled Maintenance',
      message: 'System update planned for tomorrow 2 AM',
      time: '5 hours ago'
    },
    {
      id: 3,
      type: 'success',
      title: 'Backup Completed',
      message: 'Daily backup completed successfully',
      time: '8 hours ago'
    }
  ];

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  const getStatusColor = (status) => {
    switch (status) {
      case 'active': return 'bg-green-100 text-green-700';
      case 'trial': return 'bg-blue-100 text-blue-700';
      case 'suspended': return 'bg-red-100 text-red-700';
      case 'completed': return 'bg-green-100 text-green-700';
      case 'pending': return 'bg-yellow-100 text-yellow-700';
      case 'failed': return 'bg-red-100 text-red-700';
      case 'open': return 'bg-red-100 text-red-700';
      case 'in-progress': return 'bg-blue-100 text-blue-700';
      case 'resolved': return 'bg-green-100 text-green-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getPriorityColor = (priority) => {
    switch (priority) {
      case 'high': return 'bg-red-100 text-red-700';
      case 'medium': return 'bg-yellow-100 text-yellow-700';
      case 'low': return 'bg-green-100 text-green-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getAlertColor = (type) => {
    switch (type) {
      case 'warning': return 'border-l-yellow-500 bg-yellow-50';
      case 'info': return 'border-l-blue-500 bg-blue-50';
      case 'success': return 'border-l-green-500 bg-green-50';
      case 'error': return 'border-l-red-500 bg-red-50';
      default: return 'border-l-gray-500 bg-gray-50';
    }
  };

  const filteredTenants = recentTenants.filter(tenant =>
    tenant.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    tenant.owner.toLowerCase().includes(searchQuery.toLowerCase()) ||
    tenant.email.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="min-h-screen bg-background">
      {/* Main Content */}
      <div className="p-6 space-y-6">
        {/* Welcome Section */}
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Dashboard Overview</h2>
            <p className="text-muted-foreground">
              Welcome back! Here's what's happening with your SaaS platform today.
            </p>
          </div>
          <div className="flex items-center space-x-2">
            <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="7d">Last 7 days</SelectItem>
                <SelectItem value="30d">Last 30 days</SelectItem>
                <SelectItem value="90d">Last 90 days</SelectItem>
                <SelectItem value="1y">Last year</SelectItem>
              </SelectContent>
            </Select>
            <Button variant="outline" size="sm">
              <RefreshCw className="h-4 w-4 mr-2" />
              Refresh
            </Button>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {stats.map((stat) => (
            <Card key={stat.id} className="hover:shadow-md transition-shadow">
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div className="space-y-1">
                    <p className="text-sm font-medium text-muted-foreground">
                      {stat.title}
                    </p>
                    <p className="text-2xl font-bold">{stat.value}</p>
                    <div className="flex items-center space-x-1">
                      {stat.trend === 'up' ? (
                        <TrendingUp className="h-3 w-3 text-green-500" />
                      ) : (
                        <TrendingDown className="h-3 w-3 text-red-500" />
                      )}
                      <span className={`text-xs font-medium ${
                        stat.trend === 'up' ? 'text-green-500' : 'text-red-500'
                      }`}>
                        {stat.change}
                      </span>
                      <span className="text-xs text-muted-foreground">from last month</span>
                    </div>
                    <p className="text-xs text-muted-foreground">
                      {stat.description}
                    </p>
                  </div>
                  <div className={`p-3 rounded-full bg-opacity-10 ${stat.color}`}>
                    <stat.icon className={`h-6 w-6 ${stat.color}`} />
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
            <CardDescription>Common administrative tasks</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {quickActions.map((action) => (
                <Button
                  key={action.id}
                  variant="outline"
                  className={`h-auto p-4 justify-start ${action.color}`}
                  onClick={() => {
                    if (action.id === 1) setShowAddTenantDialog(true);
                  }}
                >
                  <div className="flex items-center space-x-3">
                    <action.icon className="h-5 w-5" />
                    <div className="text-left">
                      <p className="font-medium">{action.title}</p>
                      <p className="text-xs opacity-70">{action.description}</p>
                    </div>
                  </div>
                </Button>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Charts Section */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Revenue Chart */}
          <Card className="lg:col-span-2">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Revenue & Growth Metrics</CardTitle>
                  <CardDescription>Monthly revenue, tenant count, and new signups</CardDescription>
                </div>
                <div className="flex items-center space-x-2">
                  <Badge variant="outline" className="text-green-600">
                    +12.5% Revenue Growth
                  </Badge>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={350}>
                <AreaChart data={revenueData}>
                  <defs>
                    <linearGradient id="revenueGradient" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#22c55e" stopOpacity={0.3}/>
                      <stop offset="95%" stopColor="#22c55e" stopOpacity={0}/>
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                  <XAxis dataKey="month" />
                  <YAxis />
                  <Tooltip 
                    formatter={(value, name) => {
                      if (name === 'revenue') return [formatCurrency(value), 'Revenue'];
                      if (name === 'tenants') return [value, 'Active Tenants'];
                      if (name === 'newSignups') return [value, 'New Signups'];
                      return [value, name];
                    }}
                    labelStyle={{ color: '#000' }}
                  />
                  <Legend />
                  <Area 
                    type="monotone" 
                    dataKey="revenue" 
                    stroke="#22c55e" 
                    fill="url(#revenueGradient)"
                    strokeWidth={2}
                    name="revenue"
                  />
                  <Line 
                    type="monotone" 
                    dataKey="tenants" 
                    stroke="#3b82f6" 
                    strokeWidth={2}
                    name="tenants"
                  />
                  <Line 
                    type="monotone" 
                    dataKey="newSignups" 
                    stroke="#f59e0b" 
                    strokeWidth={2}
                    name="newSignups"
                  />
                </AreaChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>

          {/* Plan Distribution */}
          <Card>
            <CardHeader>
              <CardTitle>Subscription Plans</CardTitle>
              <CardDescription>Distribution of active subscriptions</CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={200}>
                <RechartsPieChart>
                  <Pie
                    data={planDistribution}
                    cx="50%"
                    cy="50%"
                    innerRadius={40}
                    outerRadius={80}
                    paddingAngle={5}
                    dataKey="value"
                  >
                    {planDistribution.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip formatter={(value) => `${value}%`} />
                </RechartsPieChart>
              </ResponsiveContainer>
              <div className="space-y-3 mt-4">
                {planDistribution.map((plan) => (
                  <div key={plan.name} className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <div 
                        className="w-3 h-3 rounded-full" 
                        style={{ backgroundColor: plan.color }}
                      />
                      <span className="text-sm font-medium">{plan.name}</span>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-semibold">{plan.count} tenants</p>
                      <p className="text-xs text-muted-foreground">{formatCurrency(plan.revenue)}/mo</p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Top Performers & System Alerts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Top Performing Tenants */}
          <Card>
            <CardHeader>
              <CardTitle>Top Performing Tenants</CardTitle>
              <CardDescription>Highest revenue generators this month</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {topPerformingTenants.map((tenant, index) => (
                  <div key={index} className="flex items-center justify-between p-3 bg-muted/30 rounded-lg">
                    <div className="flex items-center space-x-3">
                      <div className="flex items-center justify-center w-8 h-8 bg-primary/10 text-primary rounded-full text-sm font-semibold">
                        {index + 1}
                      </div>
                      <div>
                        <p className="font-medium">{tenant.name}</p>
                        <p className="text-xs text-muted-foreground">Monthly Revenue</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold">{formatCurrency(tenant.revenue)}</p>
                      <p className="text-xs text-green-600">{tenant.growth}</p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* System Alerts */}
          <Card>
            <CardHeader>
              <CardTitle>System Alerts</CardTitle>
              <CardDescription>Recent system notifications and alerts</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {systemAlerts.map((alert) => (
                  <div key={alert.id} className={`p-4 border-l-4 rounded-r-lg ${getAlertColor(alert.type)}`}>
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-medium text-sm">{alert.title}</p>
                        <p className="text-xs text-muted-foreground mt-1">{alert.message}</p>
                      </div>
                      <span className="text-xs text-muted-foreground">{alert.time}</span>
                    </div>
                  </div>
                ))}
                <Button variant="outline" size="sm" className="w-full mt-3">
                  View All Alerts
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Data Tables Section */}
        <Tabs defaultValue="tenants" className="space-y-4">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="tenants" className="flex items-center space-x-2">
              <Users className="h-4 w-4" />
              <span>Tenants ({recentTenants.length})</span>
            </TabsTrigger>
            <TabsTrigger value="payments" className="flex items-center space-x-2">
              <CreditCard className="h-4 w-4" />
              <span>Payments ({recentPayments.length})</span>
            </TabsTrigger>
            <TabsTrigger value="support" className="flex items-center space-x-2">
              <AlertTriangle className="h-4 w-4" />
              <span>Support ({supportTickets.length})</span>
            </TabsTrigger>
          </TabsList>

          {/* Recent Tenants Tab */}
          <TabsContent value="tenants">
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>Tenant Management</CardTitle>
                    <CardDescription>Manage shop registrations and subscriptions</CardDescription>
                  </div>
                  <div className="flex items-center space-x-2">
                    <div className="relative">
                      <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                      <Input
                        placeholder="Search tenants..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-10 w-64"
                      />
                    </div>
                    <Button variant="outline" size="sm">
                      <Filter className="h-4 w-4 mr-2" />
                      Filter
                    </Button>
                    <Button size="sm" onClick={() => setShowAddTenantDialog(true)}>
                      <Plus className="h-4 w-4 mr-2" />
                      Add Tenant
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Shop Details</TableHead>
                      <TableHead>Owner Info</TableHead>
                      <TableHead>Plan & Status</TableHead>
                      <TableHead>Performance</TableHead>
                      <TableHead>Last Activity</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filteredTenants.map((tenant) => (
                      <TableRow key={tenant.id}>
                        <TableCell>
                          <div className="space-y-1">
                            <p className="font-medium">{tenant.name}</p>
                            <div className="flex items-center text-sm text-muted-foreground">
                              <MapPin className="h-3 w-3 mr-1" />
                              <span>{tenant.location}</span>
                            </div>
                            <div className="flex items-center text-sm text-muted-foreground">
                              <Package className="h-3 w-3 mr-1" />
                              <span>{tenant.productsCount} products</span>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            <div className="flex items-center space-x-2">
                              <Avatar className="h-6 w-6">
                                <AvatarFallback className="text-xs">
                                  {tenant.owner.split(' ').map(n => n[0]).join('')}
                                </AvatarFallback>
                              </Avatar>
                              <span className="font-medium text-sm">{tenant.owner}</span>
                            </div>
                            <div className="flex items-center text-xs text-muted-foreground">
                              <Phone className="h-3 w-3 mr-1" />
                              <span>{tenant.phone}</span>
                            </div>
                            <div className="flex items-center text-xs text-muted-foreground">
                              <Mail className="h-3 w-3 mr-1" />
                              <span>{tenant.email}</span>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="space-y-2">
                            <Badge variant="outline" className="text-xs">
                              {tenant.plan}
                            </Badge>
                            <Badge className={`${getStatusColor(tenant.status)} text-xs block w-fit`}>
                              {tenant.status}
                            </Badge>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            <p className="font-semibold text-sm">
                              {formatCurrency(tenant.revenue)}<span className="text-xs text-muted-foreground">/mo</span>
                            </p>
                            <p className="text-xs text-muted-foreground">
                              Joined {new Date(tenant.joinDate).toLocaleDateString()}
                            </p>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            <p className="text-sm text-muted-foreground">
                              {new Date(tenant.lastLogin).toLocaleDateString()}
                            </p>
                            <div className="flex items-center text-xs">
                              <div className={`w-2 h-2 rounded-full mr-1 ${
                                tenant.status === 'active' ? 'bg-green-500' : 
                                tenant.status === 'trial' ? 'bg-blue-500' : 'bg-red-500'
                              }`} />
                              <span className="text-muted-foreground">
                                {tenant.status === 'active' ? 'Online' : 
                                 tenant.status === 'trial' ? 'Trial' : 'Offline'}
                              </span>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="sm">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem>
                                <Eye className="h-4 w-4 mr-2" />
                                View Dashboard
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                <Edit className="h-4 w-4 mr-2" />
                                Edit Details
                              </DropdownMenuItem>
                              <DropdownMenuItem onClick={() => {
                                setSelectedTenant(tenant);
                                setShowMessageDialog(true);
                              }}>
                                <Mail className="h-4 w-4 mr-2" />
                                Send Message
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                <BarChart3 className="h-4 w-4 mr-2" />
                                View Analytics
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem>
                                <CreditCard className="h-4 w-4 mr-2" />
                                Billing History
                              </DropdownMenuItem>
                              <DropdownMenuItem className="text-orange-600">
                                <Clock className="h-4 w-4 mr-2" />
                                Suspend Account
                              </DropdownMenuItem>
                              <DropdownMenuItem className="text-red-600">
                                <Trash2 className="h-4 w-4 mr-2" />
                                Delete Account
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Recent Payments Tab */}
          <TabsContent value="payments">
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>Payment Management</CardTitle>
                    <CardDescription>Track subscription payments and transactions</CardDescription>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Badge variant="outline" className="text-green-600">
                      ৳2,45,000 This Month
                    </Badge>
                    <Button variant="outline" size="sm">
                      <Download className="h-4 w-4 mr-2" />
                      Export
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Invoice Details</TableHead>
                      <TableHead>Tenant</TableHead>
                      <TableHead>Amount & Method</TableHead>
                      <TableHead>Plan & Period</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Transaction</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {recentPayments.map((payment) => (
                      <TableRow key={payment.id}>
                        <TableCell>
                          <div className="space-y-1">
                            <p className="font-mono text-sm font-medium">
                              {payment.invoiceNumber}
                            </p>
                            <p className="text-xs text-muted-foreground">
                              {new Date(payment.date).toLocaleDateString('en-US', {
                                year: 'numeric',
                                month: 'short',
                                day: 'numeric'
                              })}
                            </p>
                          </div>
                        </TableCell>
                        <TableCell>
                          <p className="font-medium">{payment.tenant}</p>
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            <p className="font-semibold">
                              {formatCurrency(payment.amount)}
                            </p>
                            <div className="flex items-center text-sm text-muted-foreground">
                              <CreditCard className="h-3 w-3 mr-1" />
                              <span>{payment.method}</span>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            <Badge variant="outline" className="text-xs">
                              {payment.plan}
                            </Badge>
                            <p className="text-xs text-muted-foreground">Monthly</p>
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge className={`${getStatusColor(payment.status)} text-xs`}>
                            {payment.status}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <p className="font-mono text-xs text-muted-foreground">
                            {payment.transactionId}
                          </p>
                        </TableCell>
                        <TableCell>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="sm">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem>
                                <Eye className="h-4 w-4 mr-2" />
                                View Receipt
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                <Download className="h-4 w-4 mr-2" />
                                Download Invoice
                              </DropdownMenuItem>
                              {payment.status === 'failed' && (
                                <DropdownMenuItem>
                                  <RefreshCw className="h-4 w-4 mr-2" />
                                  Retry Payment
                                </DropdownMenuItem>
                              )}
                              <DropdownMenuItem>
                                <Mail className="h-4 w-4 mr-2" />
                                Send Receipt
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem>
                                <Edit className="h-4 w-4 mr-2" />
                                Update Status
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Support Tickets Tab */}
          <TabsContent value="support">
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>Support Ticket Management</CardTitle>
                    <CardDescription>Customer support requests and issue tracking</CardDescription>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Badge variant="destructive" className="flex items-center space-x-1">
                      <AlertTriangle className="h-3 w-3" />
                      <span>2 High Priority</span>
                    </Badge>
                    <Badge variant="secondary">
                      3 Open Tickets
                    </Badge>
                    <Button size="sm">
                      <Plus className="h-4 w-4 mr-2" />
                      New Ticket
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Ticket Info</TableHead>
                      <TableHead>Tenant</TableHead>
                      <TableHead>Subject & Category</TableHead>
                      <TableHead>Priority</TableHead>
                      <TableHead>Status & Assignee</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {supportTickets.map((ticket) => (
                      <TableRow key={ticket.id}>
                        <TableCell>
                          <div className="space-y-1">
                            <p className="font-mono text-sm font-medium">
                              #{ticket.id.toString().padStart(4, '0')}
                            </p>
                            <Badge variant="outline" className="text-xs">
                              {ticket.category}
                            </Badge>
                          </div>
                        </TableCell>
                        <TableCell>
                          <p className="font-medium">{ticket.tenant}</p>
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            <p className="font-medium text-sm">{ticket.subject}</p>
                            <p className="text-xs text-muted-foreground line-clamp-2">
                              {ticket.description}
                            </p>
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge className={`${getPriorityColor(ticket.priority)} text-xs`}>
                            {ticket.priority}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            <Badge className={`${getStatusColor(ticket.status)} text-xs`}>
                              {ticket.status}
                            </Badge>
                            <p className="text-xs text-muted-foreground">
                              {ticket.assignee}
                            </p>
                          </div>
                        </TableCell>
                        <TableCell>
                          <p className="text-sm text-muted-foreground">
                            {new Date(ticket.created).toLocaleDateString('en-US', {
                              month: 'short',
                              day: 'numeric',
                              hour: 'numeric',
                              minute: '2-digit'
                            })}
                          </p>
                        </TableCell>
                        <TableCell>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="sm">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem>
                                <Eye className="h-4 w-4 mr-2" />
                                View Details
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                <Edit className="h-4 w-4 mr-2" />
                                Update Status
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                <Users className="h-4 w-4 mr-2" />
                                Reassign
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                <Mail className="h-4 w-4 mr-2" />
                                Contact Customer
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem>
                                <CheckCircle className="h-4 w-4 mr-2" />
                                Mark Resolved
                              </DropdownMenuItem>
                              <DropdownMenuItem className="text-red-600">
                                <X className="h-4 w-4 mr-2" />
                                Close Ticket
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Add Tenant Dialog */}
        <Dialog open={showAddTenantDialog} onOpenChange={setShowAddTenantDialog}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Add New Tenant</DialogTitle>
              <DialogDescription>
                Register a new shop to the platform
              </DialogDescription>
            </DialogHeader>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="shopName">Shop Name</Label>
                <Input id="shopName" placeholder="Enter shop name" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="ownerName">Owner Name</Label>
                <Input id="ownerName" placeholder="Enter owner name" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="email">Email Address</Label>
                <Input id="email" type="email" placeholder="Enter email" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="phone">Phone Number</Label>
                <Input id="phone" placeholder="Enter phone number" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="location">Location</Label>
                <Input id="location" placeholder="Enter location" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="plan">Subscription Plan</Label>
                <Select>
                  <SelectTrigger>
                    <SelectValue placeholder="Select plan" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="starter">Starter - ৳2,000/month</SelectItem>
                    <SelectItem value="business">Business - ৳5,000/month</SelectItem>
                    <SelectItem value="enterprise">Enterprise - ৳10,000/month</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="md:col-span-2 space-y-2">
                <Label htmlFor="notes">Notes (Optional)</Label>
                <Textarea id="notes" placeholder="Additional notes about the tenant" />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setShowAddTenantDialog(false)}>
                Cancel
              </Button>
              <Button onClick={() => setShowAddTenantDialog(false)}>
                Create Tenant
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Send Message Dialog */}
        <Dialog open={showMessageDialog} onOpenChange={setShowMessageDialog}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Send Message</DialogTitle>
              <DialogDescription>
                Send a message to {selectedTenant?.name}
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="subject">Subject</Label>
                <Input id="subject" placeholder="Enter message subject" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="message">Message</Label>
                <Textarea 
                  id="message" 
                  placeholder="Type your message here..." 
                  rows={5}
                />
              </div>
              <div className="flex items-center space-x-2">
                <input type="checkbox" id="sms" className="rounded" />
                <Label htmlFor="sms" className="text-sm">Also send as SMS</Label>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setShowMessageDialog(false)}>
                Cancel
              </Button>
              <Button onClick={() => setShowMessageDialog(false)}>
                <Mail className="h-4 w-4 mr-2" />
                Send Message
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
};

export default SaaSAdminDashboard;