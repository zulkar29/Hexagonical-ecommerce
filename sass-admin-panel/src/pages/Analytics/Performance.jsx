import React, { useState, useEffect } from 'react';
import {
  TrendingUp,
  TrendingDown,
  Users,
  DollarSign,
  Activity,
  BarChart3,
  PieChart,
  Calendar,
  Download,
  Filter,
  RefreshCw,
  ArrowUpRight,
  ArrowDownRight,
  Eye,
  MousePointer,
  Clock,
  Globe,
  Target,
  ShoppingCart,
  CreditCard,
  Zap,
  Server,
  Database,
  Cpu,
  HardDrive,
  Wifi,
  AlertTriangle,
  CheckCircle,
  XCircle
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Progress } from '@/components/ui/progress';
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  BarChart,
  Bar,
  PieChart as RechartsPieChart,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  ComposedChart
} from 'recharts';

// Mock performance data
const performanceData = {
  revenue: {
    current: 2847650,
    previous: 2534200,
    growth: 12.4,
    target: 3000000,
    monthlyData: [
      { month: 'Jan', revenue: 185000, target: 200000, conversion: 3.2 },
      { month: 'Feb', revenue: 198000, target: 210000, conversion: 3.5 },
      { month: 'Mar', revenue: 215000, target: 220000, conversion: 3.8 },
      { month: 'Apr', revenue: 235000, target: 240000, conversion: 4.1 },
      { month: 'May', revenue: 248000, target: 250000, conversion: 4.3 },
      { month: 'Jun', revenue: 285000, target: 260000, conversion: 4.6 }
    ]
  },
  userEngagement: {
    dailyActiveUsers: 15420,
    monthlyActiveUsers: 45680,
    sessionDuration: 8.5,
    bounceRate: 32.1,
    pageViews: 234567,
    engagementData: [
      { date: '2024-01-15', dau: 12400, sessions: 18600, duration: 7.2 },
      { date: '2024-01-16', dau: 13200, sessions: 19800, duration: 7.8 },
      { date: '2024-01-17', dau: 14100, sessions: 21200, duration: 8.1 },
      { date: '2024-01-18', dau: 13800, sessions: 20400, duration: 7.9 },
      { date: '2024-01-19', dau: 15200, sessions: 22800, duration: 8.5 },
      { date: '2024-01-20', dau: 16100, sessions: 24200, duration: 8.8 },
      { date: '2024-01-21', dau: 16800, sessions: 25600, duration: 9.1 }
    ]
  },
  conversionRates: {
    overall: 4.2,
    byChannel: [
      { channel: 'Organic Search', rate: 5.8, visitors: 12400, conversions: 719 },
      { channel: 'Direct', rate: 6.2, visitors: 8900, conversions: 552 },
      { channel: 'Social Media', rate: 3.1, visitors: 15600, conversions: 484 },
      { channel: 'Email', rate: 8.9, visitors: 4200, conversions: 374 },
      { channel: 'Paid Ads', rate: 4.7, visitors: 6800, conversions: 320 }
    ]
  },
  salesPerformance: {
    totalOrders: 8945,
    averageOrderValue: 156.78,
    salesGrowth: 18.5,
    topProducts: [
      { name: 'Premium Plan', sales: 2340, revenue: 234000, growth: 22.1 },
      { name: 'Professional Plan', sales: 1890, revenue: 189000, growth: 15.8 },
      { name: 'Basic Plan', sales: 3420, revenue: 102600, growth: 8.9 },
      { name: 'Enterprise Plan', sales: 890, revenue: 267000, growth: 35.2 }
    ]
  },
  systemPerformance: {
    uptime: 99.97,
    responseTime: 245,
    errorRate: 0.12,
    throughput: 1250,
    metrics: [
      { time: '00:00', cpu: 45, memory: 62, disk: 38, network: 78 },
      { time: '04:00', cpu: 38, memory: 58, disk: 42, network: 65 },
      { time: '08:00', cpu: 72, memory: 75, disk: 45, network: 89 },
      { time: '12:00', cpu: 68, memory: 72, disk: 48, network: 92 },
      { time: '16:00', cpu: 85, memory: 82, disk: 52, network: 95 },
      { time: '20:00', cpu: 58, memory: 65, disk: 46, network: 82 }
    ]
  }
};

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(amount);
};

const formatNumber = (num) => {
  return new Intl.NumberFormat('en-US').format(num);
};

const formatPercentage = (num) => {
  return `${num > 0 ? '+' : ''}${num.toFixed(1)}%`;
};

const MetricCard = ({ title, value, change, icon: Icon, format = 'number', target }) => {
  const isPositive = change > 0;
  const formattedValue = format === 'currency' ? formatCurrency(value) : 
                       format === 'percentage' ? `${value}%` : 
                       format === 'duration' ? `${value} min` :
                       formatNumber(value);
  
  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div className="flex-1">
            <p className="text-sm font-medium text-muted-foreground">{title}</p>
            <p className="text-2xl font-bold">{formattedValue}</p>
            <div className="flex items-center mt-1">
              {isPositive ? (
                <ArrowUpRight className="h-4 w-4 text-green-500 mr-1" />
              ) : (
                <ArrowDownRight className="h-4 w-4 text-red-500 mr-1" />
              )}
              <span className={`text-sm font-medium ${
                isPositive ? 'text-green-600' : 'text-red-600'
              }`}>
                {formatPercentage(Math.abs(change))}
              </span>
              <span className="text-sm text-muted-foreground ml-1">vs last month</span>
            </div>
            {target && (
              <div className="mt-2">
                <div className="flex justify-between text-xs text-muted-foreground mb-1">
                  <span>Progress to target</span>
                  <span>{Math.round((value / target) * 100)}%</span>
                </div>
                <Progress value={(value / target) * 100} className="h-2" />
              </div>
            )}
          </div>
          <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center ml-4">
            <Icon className="h-6 w-6 text-primary" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

const SystemMetricCard = ({ title, value, status, icon: Icon, unit = '%' }) => {
  const getStatusColor = (status) => {
    switch (status) {
      case 'excellent': return 'text-green-600';
      case 'good': return 'text-blue-600';
      case 'warning': return 'text-yellow-600';
      case 'critical': return 'text-red-600';
      default: return 'text-gray-600';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'excellent': return <CheckCircle className="h-4 w-4 text-green-500" />;
      case 'good': return <CheckCircle className="h-4 w-4 text-blue-500" />;
      case 'warning': return <AlertTriangle className="h-4 w-4 text-yellow-500" />;
      case 'critical': return <XCircle className="h-4 w-4 text-red-500" />;
      default: return null;
    }
  };

  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <Icon className="h-5 w-5 text-muted-foreground" />
            <div>
              <p className="text-sm font-medium">{title}</p>
              <p className="text-lg font-bold">{value}{unit}</p>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            {getStatusIcon(status)}
            <Badge variant="outline" className={getStatusColor(status)}>
              {status}
            </Badge>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default function PerformanceAnalytics() {
  const [timeRange, setTimeRange] = useState('30d');
  const [isLoading, setIsLoading] = useState(false);
  const [lastUpdated, setLastUpdated] = useState(new Date());

  const handleRefresh = () => {
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
      setLastUpdated(new Date());
    }, 1500);
  };

  const handleExport = () => {
    console.log('Exporting performance data...');
  };

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Performance Analytics</h2>
          <p className="text-muted-foreground">
            Comprehensive performance metrics and business intelligence
          </p>
        </div>
        <div className="flex items-center space-x-3">
          <Select value={timeRange} onValueChange={setTimeRange}>
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
          <Button variant="outline" onClick={handleExport}>
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <Button variant="outline" onClick={handleRefresh} disabled={isLoading}>
            <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      <Tabs defaultValue="overview" className="space-y-6">
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="revenue">Revenue</TabsTrigger>
          <TabsTrigger value="engagement">User Engagement</TabsTrigger>
          <TabsTrigger value="conversion">Conversion</TabsTrigger>
          <TabsTrigger value="system">System Performance</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Key Performance Indicators */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <MetricCard
              title="Total Revenue"
              value={performanceData.revenue.current}
              change={performanceData.revenue.growth}
              icon={DollarSign}
              format="currency"
              target={performanceData.revenue.target}
            />
            <MetricCard
              title="Daily Active Users"
              value={performanceData.userEngagement.dailyActiveUsers}
              change={8.3}
              icon={Users}
            />
            <MetricCard
              title="Conversion Rate"
              value={performanceData.conversionRates.overall}
              change={0.8}
              icon={Target}
              format="percentage"
            />
            <MetricCard
              title="System Uptime"
              value={performanceData.systemPerformance.uptime}
              change={0.02}
              icon={Server}
              format="percentage"
            />
          </div>

          {/* Revenue vs Target Chart */}
          <Card>
            <CardHeader>
              <CardTitle>Revenue Performance vs Target</CardTitle>
              <CardDescription>
                Monthly revenue performance compared to targets
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <ComposedChart data={performanceData.revenue.monthlyData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="month" />
                  <YAxis />
                  <Tooltip formatter={(value, name) => [formatCurrency(value), name]} />
                  <Legend />
                  <Bar dataKey="target" fill="#e2e8f0" name="Target" />
                  <Bar dataKey="revenue" fill="#3b82f6" name="Actual Revenue" />
                  <Line type="monotone" dataKey="conversion" stroke="#10b981" strokeWidth={2} name="Conversion Rate (%)" />
                </ComposedChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>

          {/* Quick Stats Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Sales Performance</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Total Orders</span>
                  <span className="font-semibold">{formatNumber(performanceData.salesPerformance.totalOrders)}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Avg Order Value</span>
                  <span className="font-semibold">{formatCurrency(performanceData.salesPerformance.averageOrderValue)}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Sales Growth</span>
                  <Badge variant="secondary" className="text-green-600">
                    {formatPercentage(performanceData.salesPerformance.salesGrowth)}
                  </Badge>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="text-lg">User Engagement</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Session Duration</span>
                  <span className="font-semibold">{performanceData.userEngagement.sessionDuration} min</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Bounce Rate</span>
                  <span className="font-semibold">{performanceData.userEngagement.bounceRate}%</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Page Views</span>
                  <span className="font-semibold">{formatNumber(performanceData.userEngagement.pageViews)}</span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="text-lg">System Health</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Response Time</span>
                  <span className="font-semibold">{performanceData.systemPerformance.responseTime}ms</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Error Rate</span>
                  <span className="font-semibold">{performanceData.systemPerformance.errorRate}%</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Throughput</span>
                  <span className="font-semibold">{formatNumber(performanceData.systemPerformance.throughput)} req/min</span>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="revenue" className="space-y-6">
          {/* Revenue Analytics */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <MetricCard
              title="Monthly Revenue"
              value={performanceData.revenue.monthlyData[performanceData.revenue.monthlyData.length - 1].revenue}
              change={15.2}
              icon={DollarSign}
              format="currency"
            />
            <MetricCard
              title="Revenue Growth"
              value={performanceData.revenue.growth}
              change={2.1}
              icon={TrendingUp}
              format="percentage"
            />
            <MetricCard
              title="Target Achievement"
              value={Math.round((performanceData.revenue.current / performanceData.revenue.target) * 100)}
              change={5.8}
              icon={Target}
              format="percentage"
            />
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Revenue Trend Analysis</CardTitle>
              <CardDescription>Monthly revenue performance over time</CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={400}>
                <AreaChart data={performanceData.revenue.monthlyData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="month" />
                  <YAxis />
                  <Tooltip formatter={(value) => [formatCurrency(value), 'Revenue']} />
                  <Area type="monotone" dataKey="revenue" stroke="#3b82f6" fill="#3b82f6" fillOpacity={0.3} />
                </AreaChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Top Performing Products</CardTitle>
              <CardDescription>Revenue breakdown by product category</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {performanceData.salesPerformance.topProducts.map((product, index) => (
                  <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex items-center space-x-4">
                      <div className="w-2 h-8 bg-primary rounded"></div>
                      <div>
                        <p className="font-medium">{product.name}</p>
                        <p className="text-sm text-muted-foreground">{product.sales} sales</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold">{formatCurrency(product.revenue)}</p>
                      <Badge variant="secondary" className="text-green-600">
                        {formatPercentage(product.growth)}
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="engagement" className="space-y-6">
          {/* User Engagement Analytics */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            <MetricCard
              title="Daily Active Users"
              value={performanceData.userEngagement.dailyActiveUsers}
              change={8.3}
              icon={Users}
            />
            <MetricCard
              title="Session Duration"
              value={performanceData.userEngagement.sessionDuration}
              change={12.5}
              icon={Clock}
              format="duration"
            />
            <MetricCard
              title="Bounce Rate"
              value={performanceData.userEngagement.bounceRate}
              change={-3.2}
              icon={MousePointer}
              format="percentage"
            />
            <MetricCard
              title="Page Views"
              value={performanceData.userEngagement.pageViews}
              change={18.7}
              icon={Eye}
            />
          </div>

          <Card>
            <CardHeader>
              <CardTitle>User Engagement Trends</CardTitle>
              <CardDescription>Daily active users and session metrics</CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={400}>
                <ComposedChart data={performanceData.userEngagement.engagementData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis yAxisId="left" />
                  <YAxis yAxisId="right" orientation="right" />
                  <Tooltip />
                  <Legend />
                  <Bar yAxisId="left" dataKey="dau" fill="#3b82f6" name="Daily Active Users" />
                  <Line yAxisId="right" type="monotone" dataKey="duration" stroke="#10b981" strokeWidth={2} name="Avg Session Duration (min)" />
                </ComposedChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="conversion" className="space-y-6">
          {/* Conversion Analytics */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <MetricCard
              title="Overall Conversion Rate"
              value={performanceData.conversionRates.overall}
              change={0.8}
              icon={Target}
              format="percentage"
            />
            <MetricCard
              title="Best Performing Channel"
              value={Math.max(...performanceData.conversionRates.byChannel.map(c => c.rate))}
              change={1.2}
              icon={TrendingUp}
              format="percentage"
            />
            <MetricCard
              title="Total Conversions"
              value={performanceData.conversionRates.byChannel.reduce((sum, c) => sum + c.conversions, 0)}
              change={15.3}
              icon={ShoppingCart}
            />
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Conversion Rate by Channel</CardTitle>
              <CardDescription>Performance comparison across different traffic sources</CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={400}>
                <BarChart data={performanceData.conversionRates.byChannel}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="channel" />
                  <YAxis />
                  <Tooltip formatter={(value, name) => [name === 'rate' ? `${value}%` : formatNumber(value), name]} />
                  <Legend />
                  <Bar dataKey="rate" fill="#3b82f6" name="Conversion Rate (%)" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Channel Performance Details</CardTitle>
              <CardDescription>Detailed breakdown of traffic and conversions by channel</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {performanceData.conversionRates.byChannel.map((channel, index) => (
                  <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex items-center space-x-4">
                      <div className="w-3 h-3 bg-primary rounded-full"></div>
                      <div>
                        <p className="font-medium">{channel.channel}</p>
                        <p className="text-sm text-muted-foreground">{formatNumber(channel.visitors)} visitors</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold">{channel.rate}% conversion</p>
                      <p className="text-sm text-muted-foreground">{formatNumber(channel.conversions)} conversions</p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="system" className="space-y-6">
          {/* System Performance */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <SystemMetricCard
              title="CPU Usage"
              value={68}
              status="good"
              icon={Cpu}
            />
            <SystemMetricCard
              title="Memory Usage"
              value={72}
              status="good"
              icon={HardDrive}
            />
            <SystemMetricCard
              title="Disk Usage"
              value={45}
              status="excellent"
              icon={Database}
            />
            <SystemMetricCard
              title="Network I/O"
              value={82}
              status="warning"
              icon={Wifi}
            />
          </div>

          <Card>
            <CardHeader>
              <CardTitle>System Performance Metrics</CardTitle>
              <CardDescription>Real-time system resource utilization</CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={400}>
                <LineChart data={performanceData.systemPerformance.metrics}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="time" />
                  <YAxis />
                  <Tooltip formatter={(value) => [`${value}%`, 'Usage']} />
                  <Legend />
                  <Line type="monotone" dataKey="cpu" stroke="#3b82f6" strokeWidth={2} name="CPU" />
                  <Line type="monotone" dataKey="memory" stroke="#10b981" strokeWidth={2} name="Memory" />
                  <Line type="monotone" dataKey="disk" stroke="#f59e0b" strokeWidth={2} name="Disk" />
                  <Line type="monotone" dataKey="network" stroke="#ef4444" strokeWidth={2} name="Network" />
                </LineChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Performance Metrics</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Response Time</span>
                  <span className="font-semibold">{performanceData.systemPerformance.responseTime}ms</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Throughput</span>
                  <span className="font-semibold">{formatNumber(performanceData.systemPerformance.throughput)} req/min</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Error Rate</span>
                  <span className="font-semibold">{performanceData.systemPerformance.errorRate}%</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Uptime</span>
                  <Badge variant="secondary" className="text-green-600">
                    {performanceData.systemPerformance.uptime}%
                  </Badge>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>System Health Status</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm">Overall Health</span>
                  <Badge variant="secondary" className="text-green-600">Excellent</Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm">Last Incident</span>
                  <span className="text-sm text-muted-foreground">7 days ago</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm">Monitoring Status</span>
                  <div className="flex items-center space-x-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    <span className="text-sm">Active</span>
                  </div>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm">Next Maintenance</span>
                  <span className="text-sm text-muted-foreground">In 14 days</span>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}