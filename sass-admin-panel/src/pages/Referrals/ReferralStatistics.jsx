import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  LineChart,
  Line,
  PieChart,
  Pie,
  Cell,
  Area,
  AreaChart
} from 'recharts';
import {
  TrendingUp,
  TrendingDown,
  Users,
  DollarSign,
  Share2,
  Target,
  Calendar,
  Download,
  RefreshCw,
  Eye,
  MousePointer,
  UserCheck
} from 'lucide-react';
import { toast } from 'sonner';
import referralService from '@/services/referralService';

const ReferralStatistics = () => {
  const [stats, setStats] = useState(null);
  const [chartData, setChartData] = useState({
    conversions: [],
    commissions: [],
    topCodes: [],
    trends: []
  });
  const [timeRange, setTimeRange] = useState('30d');
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const timeRangeOptions = [
    { value: '7d', label: 'Last 7 days' },
    { value: '30d', label: 'Last 30 days' },
    { value: '90d', label: 'Last 3 months' },
    { value: '1y', label: 'Last year' }
  ];

  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8'];

  useEffect(() => {
    fetchStatistics();
  }, [timeRange]);

  const fetchStatistics = async () => {
    try {
      setLoading(true);
      const [statsResponse, analyticsResponse] = await Promise.all([
        referralService.getReferralStatistics(timeRange),
        referralService.getReferralAnalytics(timeRange)
      ]);

      setStats(statsResponse.data);
      setChartData({
        conversions: analyticsResponse.data.conversions || [],
        commissions: analyticsResponse.data.commissions || [],
        topCodes: analyticsResponse.data.topCodes || [],
        trends: analyticsResponse.data.trends || []
      });
    } catch (error) {
      toast.error('Failed to fetch referral statistics');
      console.error('Statistics fetch error:', error);
    } finally {
      setLoading(false);
    }
  };

  const refreshData = async () => {
    setRefreshing(true);
    await fetchStatistics();
    setRefreshing(false);
    toast.success('Statistics refreshed');
  };

  const exportData = async () => {
    try {
      const response = await referralService.exportReferralData({
        timeRange,
        format: 'csv',
        includeCharts: true
      });
      
      // Create download link
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `referral-statistics-${timeRange}.csv`);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      toast.success('Statistics exported successfully');
    } catch (error) {
      toast.error('Failed to export statistics');
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-BD', {
      style: 'currency',
      currency: 'BDT',
      minimumFractionDigits: 0
    }).format(amount);
  };

  const formatPercentage = (value) => {
    return `${(value * 100).toFixed(1)}%`;
  };

  const StatCard = ({ title, value, change, icon: Icon, color = 'blue', format = 'number' }) => {
    const isPositive = change >= 0;
    const formattedValue = format === 'currency' ? formatCurrency(value) : 
                          format === 'percentage' ? formatPercentage(value) : 
                          value.toLocaleString();

    return (
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-muted-foreground">{title}</p>
              <p className="text-2xl font-bold">{formattedValue}</p>
              <div className="flex items-center mt-1">
                {isPositive ? (
                  <TrendingUp className="h-4 w-4 text-green-500 mr-1" />
                ) : (
                  <TrendingDown className="h-4 w-4 text-red-500 mr-1" />
                )}
                <span className={`text-sm ${
                  isPositive ? 'text-green-500' : 'text-red-500'
                }`}>
                  {Math.abs(change).toFixed(1)}%
                </span>
                <span className="text-sm text-muted-foreground ml-1">
                  vs last period
                </span>
              </div>
            </div>
            <div className={`p-3 rounded-full bg-${color}-100`}>
              <Icon className={`h-6 w-6 text-${color}-600`} />
            </div>
          </div>
        </CardContent>
      </Card>
    );
  };

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h2 className="text-2xl font-bold">Referral Statistics</h2>
          <div className="flex gap-2">
            <div className="w-32 h-10 bg-muted animate-pulse rounded" />
            <div className="w-24 h-10 bg-muted animate-pulse rounded" />
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="h-32 bg-muted animate-pulse rounded-lg" />
          ))}
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="h-80 bg-muted animate-pulse rounded-lg" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-2xl font-bold">Referral Statistics</h2>
          <p className="text-muted-foreground">
            Track your referral performance and earnings
          </p>
        </div>
        <div className="flex gap-2">
          <Select value={timeRange} onValueChange={setTimeRange}>
            <SelectTrigger className="w-40">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {timeRangeOptions.map(option => (
                <SelectItem key={option.value} value={option.value}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button variant="outline" onClick={refreshData} disabled={refreshing}>
            <RefreshCw className={`h-4 w-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          <Button variant="outline" onClick={exportData}>
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
        </div>
      </div>

      {/* Key Metrics */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <StatCard
            title="Total Referrals"
            value={stats.totalReferrals}
            change={stats.referralGrowth}
            icon={Users}
            color="blue"
          />
          <StatCard
            title="Total Commissions"
            value={stats.totalCommissions}
            change={stats.commissionGrowth}
            icon={DollarSign}
            color="green"
            format="currency"
          />
          <StatCard
            title="Conversion Rate"
            value={stats.conversionRate}
            change={stats.conversionGrowth}
            icon={Target}
            color="purple"
            format="percentage"
          />
          <StatCard
            title="Active Codes"
            value={stats.activeCodes}
            change={stats.codeGrowth}
            icon={Share2}
            color="orange"
          />
        </div>
      )}

      {/* Charts Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Conversion Funnel */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Target className="h-5 w-5" />
              Conversion Funnel
            </CardTitle>
            <CardDescription>
              Track user journey from clicks to conversions
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={chartData.conversions}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="stage" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="count" fill="#8884d8" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Commission Trends */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <DollarSign className="h-5 w-5" />
              Commission Trends
            </CardTitle>
            <CardDescription>
              Daily commission earnings over time
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={chartData.commissions}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip formatter={(value) => [formatCurrency(value), 'Commission']} />
                <Area 
                  type="monotone" 
                  dataKey="amount" 
                  stroke="#82ca9d" 
                  fill="#82ca9d" 
                  fillOpacity={0.6}
                />
              </AreaChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Top Performing Codes */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Share2 className="h-5 w-5" />
              Top Performing Codes
            </CardTitle>
            <CardDescription>
              Best performing referral codes by conversions
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={chartData.topCodes}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="conversions"
                >
                  {chartData.topCodes.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Traffic Trends */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Traffic Trends
            </CardTitle>
            <CardDescription>
              Referral traffic patterns over time
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={chartData.trends}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip />
                <Line 
                  type="monotone" 
                  dataKey="clicks" 
                  stroke="#8884d8" 
                  name="Clicks"
                />
                <Line 
                  type="monotone" 
                  dataKey="conversions" 
                  stroke="#82ca9d" 
                  name="Conversions"
                />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Detailed Metrics */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Eye className="h-5 w-5" />
                Traffic Metrics
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Total Clicks</span>
                <span className="font-medium">{stats.totalClicks?.toLocaleString()}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Unique Visitors</span>
                <span className="font-medium">{stats.uniqueVisitors?.toLocaleString()}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Bounce Rate</span>
                <span className="font-medium">{formatPercentage(stats.bounceRate || 0)}</span>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <UserCheck className="h-5 w-5" />
                Conversion Metrics
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Sign-ups</span>
                <span className="font-medium">{stats.signups?.toLocaleString()}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Purchases</span>
                <span className="font-medium">{stats.purchases?.toLocaleString()}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Avg. Order Value</span>
                <span className="font-medium">{formatCurrency(stats.avgOrderValue || 0)}</span>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <DollarSign className="h-5 w-5" />
                Revenue Metrics
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Total Revenue</span>
                <span className="font-medium">{formatCurrency(stats.totalRevenue || 0)}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Pending Commissions</span>
                <span className="font-medium">{formatCurrency(stats.pendingCommissions || 0)}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Paid Commissions</span>
                <span className="font-medium">{formatCurrency(stats.paidCommissions || 0)}</span>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
};

export default ReferralStatistics;