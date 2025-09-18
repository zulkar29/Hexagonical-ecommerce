import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
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
  Globe
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
  ResponsiveContainer
} from 'recharts';

// Mock analytics data
const mockAnalyticsData = {
  overview: {
    totalRevenue: 2847650,
    revenueGrowth: 12.5,
    totalTenants: 1247,
    tenantGrowth: 8.3,
    activeUsers: 15420,
    userGrowth: -2.1,
    systemUptime: 99.97,
    uptimeChange: 0.02
  },
  revenueData: [
    { month: 'Jan', revenue: 185000, tenants: 980, users: 12400 },
    { month: 'Feb', revenue: 198000, tenants: 1020, users: 13100 },
    { month: 'Mar', revenue: 215000, tenants: 1080, users: 13800 },
    { month: 'Apr', revenue: 235000, tenants: 1150, users: 14200 },
    { month: 'May', revenue: 248000, tenants: 1190, users: 14800 },
    { month: 'Jun', revenue: 285000, tenants: 1247, users: 15420 }
  ],
  tenantDistribution: [
    { name: 'Starter', value: 45, count: 561, color: '#8884d8' },
    { name: 'Professional', value: 35, count: 437, color: '#82ca9d' },
    { name: 'Enterprise', value: 15, count: 187, color: '#ffc658' },
    { name: 'Custom', value: 5, count: 62, color: '#ff7300' }
  ],
  trafficData: [
    { date: '2024-01-15', pageViews: 45200, uniqueVisitors: 12400, sessions: 18600 },
    { date: '2024-01-16', pageViews: 48100, uniqueVisitors: 13200, sessions: 19800 },
    { date: '2024-01-17', pageViews: 52300, uniqueVisitors: 14100, sessions: 21200 },
    { date: '2024-01-18', pageViews: 49800, uniqueVisitors: 13800, sessions: 20400 },
    { date: '2024-01-19', pageViews: 55600, uniqueVisitors: 15200, sessions: 22800 },
    { date: '2024-01-20', pageViews: 58900, uniqueVisitors: 16100, sessions: 24200 },
    { date: '2024-01-21', pageViews: 61200, uniqueVisitors: 16800, sessions: 25600 }
  ],
  topPerformingTenants: [
    { id: 1, name: 'TechCorp Solutions', revenue: 45200, growth: 18.5, plan: 'Enterprise' },
    { id: 2, name: 'Digital Dynamics', revenue: 38900, growth: 15.2, plan: 'Professional' },
    { id: 3, name: 'Innovation Labs', revenue: 32100, growth: 22.8, plan: 'Enterprise' },
    { id: 4, name: 'StartupHub', revenue: 28700, growth: 12.4, plan: 'Professional' },
    { id: 5, name: 'CloudFirst Inc', revenue: 25400, growth: 9.8, plan: 'Enterprise' }
  ],
  systemMetrics: {
    responseTime: 245,
    errorRate: 0.12,
    throughput: 1250,
    cpuUsage: 68,
    memoryUsage: 72,
    diskUsage: 45
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

const StatCard = ({ title, value, change, icon: Icon, format = 'number' }) => {
  const isPositive = change > 0;
  const formattedValue = format === 'currency' ? formatCurrency(value) : 
                       format === 'percentage' ? `${value}%` : 
                       formatNumber(value);
  
  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div>
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
          </div>
          <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center">
            <Icon className="h-6 w-6 text-primary" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default function AnalyticsHome() {
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
    console.log('Exporting analytics data...');
  };

  return (
    <div className="space-y-6 p-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-3xl font-bold tracking-tight">Analytics Dashboard</h2>
              <p className="text-muted-foreground">
                Platform performance metrics and business insights
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

          {/* Last Updated */}
          <div className="text-sm text-muted-foreground">
            Last updated: {lastUpdated.toLocaleString()}
          </div>

          {/* Overview Stats */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <StatCard
              title="Total Revenue"
              value={mockAnalyticsData.overview.totalRevenue}
              change={mockAnalyticsData.overview.revenueGrowth}
              icon={DollarSign}
              format="currency"
            />
            <StatCard
              title="Active Tenants"
              value={mockAnalyticsData.overview.totalTenants}
              change={mockAnalyticsData.overview.tenantGrowth}
              icon={Users}
            />
            <StatCard
              title="Platform Users"
              value={mockAnalyticsData.overview.activeUsers}
              change={mockAnalyticsData.overview.userGrowth}
              icon={Activity}
            />
            <StatCard
              title="System Uptime"
              value={mockAnalyticsData.overview.systemUptime}
              change={mockAnalyticsData.overview.uptimeChange}
              icon={Globe}
              format="percentage"
            />
          </div>

          {/* Analytics Tabs */}
          <Tabs defaultValue="revenue" className="space-y-6">
            <TabsList className="grid w-full grid-cols-4">
              <TabsTrigger value="revenue">Revenue Analytics</TabsTrigger>
              <TabsTrigger value="tenants">Tenant Insights</TabsTrigger>
              <TabsTrigger value="traffic">Traffic Analytics</TabsTrigger>
              <TabsTrigger value="performance">System Performance</TabsTrigger>
            </TabsList>

            {/* Revenue Analytics */}
            <TabsContent value="revenue" className="space-y-6">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Revenue Trend</CardTitle>
                    <CardDescription>Monthly revenue growth over time</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <ResponsiveContainer width="100%" height={300}>
                      <AreaChart data={mockAnalyticsData.revenueData}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="month" />
                        <YAxis />
                        <Tooltip formatter={(value) => formatCurrency(value)} />
                        <Area 
                          type="monotone" 
                          dataKey="revenue" 
                          stroke="#8884d8" 
                          fill="#8884d8" 
                          fillOpacity={0.3}
                        />
                      </AreaChart>
                    </ResponsiveContainer>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle>Top Performing Tenants</CardTitle>
                    <CardDescription>Highest revenue generating tenants</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      {mockAnalyticsData.topPerformingTenants.map((tenant, index) => (
                        <div key={tenant.id} className="flex items-center justify-between p-3 border rounded-lg">
                          <div className="flex items-center space-x-3">
                            <div className="h-8 w-8 bg-primary/10 rounded-full flex items-center justify-center text-sm font-medium">
                              {index + 1}
                            </div>
                            <div>
                              <p className="font-medium">{tenant.name}</p>
                              <div className="flex items-center space-x-2">
                                <Badge variant="outline">{tenant.plan}</Badge>
                                <span className="text-sm text-green-600">+{tenant.growth}%</span>
                              </div>
                            </div>
                          </div>
                          <div className="text-right">
                            <p className="font-semibold">{formatCurrency(tenant.revenue)}</p>
                            <p className="text-sm text-muted-foreground">this month</p>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>

            {/* Tenant Insights */}
            <TabsContent value="tenants" className="space-y-6">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Tenant Growth</CardTitle>
                    <CardDescription>New tenant acquisitions over time</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <ResponsiveContainer width="100%" height={300}>
                      <LineChart data={mockAnalyticsData.revenueData}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="month" />
                        <YAxis />
                        <Tooltip />
                        <Line 
                          type="monotone" 
                          dataKey="tenants" 
                          stroke="#82ca9d" 
                          strokeWidth={3}
                        />
                      </LineChart>
                    </ResponsiveContainer>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle>Plan Distribution</CardTitle>
                    <CardDescription>Breakdown of tenant subscription plans</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <ResponsiveContainer width="100%" height={300}>
                      <RechartsPieChart>
                        <Pie
                          data={mockAnalyticsData.tenantDistribution}
                          cx="50%"
                          cy="50%"
                          outerRadius={80}
                          dataKey="value"
                        >
                          {mockAnalyticsData.tenantDistribution.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={entry.color} />
                          ))}
                        </Pie>
                        <Tooltip formatter={(value) => `${value}%`} />
                        <Legend />
                      </RechartsPieChart>
                    </ResponsiveContainer>
                    <div className="grid grid-cols-2 gap-4 mt-4">
                      {mockAnalyticsData.tenantDistribution.map((item) => (
                        <div key={item.name} className="flex items-center space-x-2">
                          <div 
                            className="h-3 w-3 rounded-full" 
                            style={{ backgroundColor: item.color }}
                          ></div>
                          <span className="text-sm">{item.name}: {item.count}</span>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>

            {/* Traffic Analytics */}
            <TabsContent value="traffic" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Platform Traffic</CardTitle>
                  <CardDescription>Page views, visitors, and session data</CardDescription>
                </CardHeader>
                <CardContent>
                  <ResponsiveContainer width="100%" height={400}>
                    <AreaChart data={mockAnalyticsData.trafficData}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="date" />
                      <YAxis />
                      <Tooltip />
                      <Legend />
                      <Area 
                        type="monotone" 
                        dataKey="pageViews" 
                        stackId="1" 
                        stroke="#8884d8" 
                        fill="#8884d8" 
                        name="Page Views"
                      />
                      <Area 
                        type="monotone" 
                        dataKey="sessions" 
                        stackId="1" 
                        stroke="#82ca9d" 
                        fill="#82ca9d" 
                        name="Sessions"
                      />
                      <Area 
                        type="monotone" 
                        dataKey="uniqueVisitors" 
                        stackId="1" 
                        stroke="#ffc658" 
                        fill="#ffc658" 
                        name="Unique Visitors"
                      />
                    </AreaChart>
                  </ResponsiveContainer>
                </CardContent>
              </Card>
            </TabsContent>

            {/* System Performance */}
            <TabsContent value="performance" className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center space-x-2">
                      <Clock className="h-5 w-5" />
                      <span>Response Time</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{mockAnalyticsData.systemMetrics.responseTime}ms</div>
                    <p className="text-sm text-muted-foreground">Average response time</p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center space-x-2">
                      <Activity className="h-5 w-5" />
                      <span>Error Rate</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{mockAnalyticsData.systemMetrics.errorRate}%</div>
                    <p className="text-sm text-muted-foreground">System error rate</p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center space-x-2">
                      <BarChart3 className="h-5 w-5" />
                      <span>Throughput</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{formatNumber(mockAnalyticsData.systemMetrics.throughput)}</div>
                    <p className="text-sm text-muted-foreground">Requests per minute</p>
                  </CardContent>
                </Card>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle>CPU Usage</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2">
                      <div className="flex justify-between">
                        <span>Current</span>
                        <span>{mockAnalyticsData.systemMetrics.cpuUsage}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-blue-600 h-2 rounded-full" 
                          style={{ width: `${mockAnalyticsData.systemMetrics.cpuUsage}%` }}
                        ></div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle>Memory Usage</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2">
                      <div className="flex justify-between">
                        <span>Current</span>
                        <span>{mockAnalyticsData.systemMetrics.memoryUsage}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-green-600 h-2 rounded-full" 
                          style={{ width: `${mockAnalyticsData.systemMetrics.memoryUsage}%` }}
                        ></div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle>Disk Usage</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2">
                      <div className="flex justify-between">
                        <span>Current</span>
                        <span>{mockAnalyticsData.systemMetrics.diskUsage}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-orange-600 h-2 rounded-full" 
                          style={{ width: `${mockAnalyticsData.systemMetrics.diskUsage}%` }}
                        ></div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>
          </Tabs>

          {/* Quick Actions */}
          <Card>
            <CardHeader>
              <CardTitle>Quick Actions</CardTitle>
              <CardDescription>Access detailed analytics and reports</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <Link to="/analytics/revenue">
                  <Button variant="outline" className="w-full justify-start">
                    <DollarSign className="h-4 w-4 mr-2" />
                    Revenue Analytics
                  </Button>
                </Link>
                <Link to="/analytics/tenants">
                  <Button variant="outline" className="w-full justify-start">
                    <Users className="h-4 w-4 mr-2" />
                    Tenant Analytics
                  </Button>
                </Link>
                <Link to="/analytics/performance">
                  <Button variant="outline" className="w-full justify-start">
                    <Activity className="h-4 w-4 mr-2" />
                    Performance Metrics
                  </Button>
                </Link>
                <Link to="/analytics/custom">
                  <Button variant="outline" className="w-full justify-start">
                    <BarChart3 className="h-4 w-4 mr-2" />
                    Custom Reports
                  </Button>
                </Link>
              </div>
            </CardContent>
          </Card>
    </div>
  );
}
