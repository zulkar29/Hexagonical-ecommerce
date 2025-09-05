import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  BarChart3,
  TrendingUp,
  Package,
  Users,
  ShoppingCart,
  DollarSign,
  Calendar,
  Download,
  FileText,
  Eye,
  ArrowRight,
  TrendingDown,
  AlertTriangle,
  CheckCircle,
  Clock
} from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ResponsiveContainer, AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, PieChart, Pie, Cell, BarChart, Bar } from 'recharts';

// Mock data for dashboard
const mockSalesData = [
  { date: '2025-01-15', sales: 12000, orders: 18 },
  { date: '2025-01-16', sales: 14500, orders: 22 },
  { date: '2025-01-17', sales: 9800, orders: 14 },
  { date: '2025-01-18', sales: 17200, orders: 25 },
  { date: '2025-01-19', sales: 15800, orders: 21 },
  { date: '2025-01-20', sales: 13400, orders: 19 },
  { date: '2025-01-21', sales: 16200, orders: 23 }
];

const mockCategoryData = [
  { name: 'Rice & Grains', value: 35, color: '#22c55e' },
  { name: 'Vegetables', value: 25, color: '#3b82f6' },
  { name: 'Dairy', value: 20, color: '#f59e0b' },
  { name: 'Cooking', value: 12, color: '#ef4444' },
  { name: 'Others', value: 8, color: '#8b5cf6' }
];

const mockTopProducts = [
  { name: 'Basmati Rice 5kg', sales: 450, revenue: 225000 },
  { name: 'Onion 1kg', sales: 320, revenue: 64000 },
  { name: 'Milk 1L', sales: 280, revenue: 84000 },
  { name: 'Cooking Oil 1L', sales: 200, revenue: 80000 },
  { name: 'Sugar 1kg', sales: 180, revenue: 54000 }
];

const reportCards = [
  {
    title: 'Sales Report',
    description: 'Daily, weekly, and monthly sales performance with trends and comparisons',
    icon: TrendingUp,
    path: '/reports/sales',
    color: 'text-green-600',
    bgColor: 'bg-green-50',
    metrics: { value: '৳98,400', change: '+12.5%', trend: 'up' }
  },
  {
    title: 'Inventory Report',
    description: 'Stock levels, low stock alerts, and inventory movement tracking',
    icon: Package,
    path: '/reports/inventory',
    color: 'text-blue-600',
    bgColor: 'bg-blue-50',
    metrics: { value: '156 Items', change: '3 Low Stock', trend: 'warning' }
  },
  {
    title: 'Customer Insights',
    description: 'Customer behavior, purchase history, and loyalty analytics',
    icon: Users,
    path: '/reports/customers',
    color: 'text-purple-600',
    bgColor: 'bg-purple-50',
    metrics: { value: '1,234 Customers', change: '+8.2%', trend: 'up' }
  }
];

const Reports = () => {
  const navigate = useNavigate();
  const [dateRange, setDateRange] = useState('7d');

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  
  const totalSales = mockSalesData.reduce((sum, d) => sum + d.sales, 0);
  const totalOrders = mockSalesData.reduce((sum, d) => sum + d.orders, 0);
  const avgOrderValue = totalOrders ? Math.round(totalSales / totalOrders) : 0;

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Analytics & Reports</h1>
          <p className="text-muted-foreground">
            Comprehensive business insights and performance analytics for your shop
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm">
            <Calendar className="h-4 w-4 mr-2" />
            Last 7 Days
          </Button>
          <Button variant="outline" size="sm">
            <Download className="h-4 w-4 mr-2" />
            Export All
          </Button>
        </div>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <DollarSign className="h-8 w-8 text-green-600" />
            <div>
              <p className="text-sm text-muted-foreground">Total Revenue</p>
              <p className="text-2xl font-bold">{formatCurrency(totalSales)}</p>
              <div className="flex items-center gap-1 text-sm text-green-600">
                <TrendingUp className="h-3 w-3" />
                <span>+12.5%</span>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <ShoppingCart className="h-8 w-8 text-blue-600" />
            <div>
              <p className="text-sm text-muted-foreground">Total Orders</p>
              <p className="text-2xl font-bold">{totalOrders}</p>
              <div className="flex items-center gap-1 text-sm text-blue-600">
                <TrendingUp className="h-3 w-3" />
                <span>+8.2%</span>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <BarChart3 className="h-8 w-8 text-purple-600" />
            <div>
              <p className="text-sm text-muted-foreground">Avg Order Value</p>
              <p className="text-2xl font-bold">{formatCurrency(avgOrderValue)}</p>
              <div className="flex items-center gap-1 text-sm text-purple-600">
                <TrendingUp className="h-3 w-3" />
                <span>+5.8%</span>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <Users className="h-8 w-8 text-orange-600" />
            <div>
              <p className="text-sm text-muted-foreground">Active Customers</p>
              <p className="text-2xl font-bold">1,234</p>
              <div className="flex items-center gap-1 text-sm text-orange-600">
                <TrendingUp className="h-3 w-3" />
                <span>+15.3%</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Report Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        {reportCards.map((report, index) => (
          <Card key={index} className="hover:shadow-lg transition-shadow cursor-pointer">
            <CardContent className="p-6">
              <div className="flex items-start justify-between mb-4">
                <div className={`p-3 rounded-lg ${report.bgColor}`}>
                  <report.icon className={`h-6 w-6 ${report.color}`} />
                </div>
                <Button 
                  variant="ghost" 
                  size="sm" 
                  onClick={() => navigate(report.path)}
                  className="h-8 w-8 p-0"
                >
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </div>
              <h3 className="text-lg font-semibold mb-2">{report.title}</h3>
              <p className="text-sm text-muted-foreground mb-4">{report.description}</p>
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-xl font-bold">{report.metrics.value}</p>
                  <div className={`flex items-center gap-1 text-sm ${
                    report.metrics.trend === 'up' ? 'text-green-600' : 
                    report.metrics.trend === 'down' ? 'text-red-600' : 
                    'text-orange-600'
                  }`}>
                    {report.metrics.trend === 'up' && <TrendingUp className="h-3 w-3" />}
                    {report.metrics.trend === 'down' && <TrendingDown className="h-3 w-3" />}
                    {report.metrics.trend === 'warning' && <AlertTriangle className="h-3 w-3" />}
                    <span>{report.metrics.change}</span>
                  </div>
                </div>
                <Button variant="outline" size="sm" onClick={() => navigate(report.path)}>
                  <Eye className="h-4 w-4 mr-2" />
                  View Report
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        {/* Sales Trend */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5 text-green-600" />
              Sales Trend (Last 7 Days)
            </CardTitle>
            <CardDescription>Daily sales performance and order volume</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={mockSalesData}>
                <defs>
                  <linearGradient id="salesGradient" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#22c55e" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#22c55e" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                <XAxis 
                  dataKey="date" 
                  tickFormatter={(value) => new Date(value).toLocaleDateString('en-US', { day: 'numeric', month: 'short' })}
                />
                <YAxis tickFormatter={(value) => `৳${(value/1000).toFixed(0)}k`} />
                <Tooltip 
                  formatter={(value, name) => [
                    name === 'sales' ? formatCurrency(value) : value,
                    name === 'sales' ? 'Sales' : 'Orders'
                  ]}
                  labelFormatter={(label) => new Date(label).toLocaleDateString()}
                />
                <Area 
                  type="monotone" 
                  dataKey="sales" 
                  stroke="#22c55e" 
                  fillOpacity={1} 
                  fill="url(#salesGradient)" 
                />
              </AreaChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Category Distribution */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Package className="h-5 w-5 text-blue-600" />
              Sales by Category
            </CardTitle>
            <CardDescription>Revenue distribution across product categories</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={mockCategoryData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={120}
                  paddingAngle={5}
                  dataKey="value"
                >
                  {mockCategoryData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip formatter={(value) => [`${value}%`, 'Share']} />
              </PieChart>
            </ResponsiveContainer>
            <div className="grid grid-cols-2 gap-2 mt-4">
              {mockCategoryData.map((category, index) => (
                <div key={index} className="flex items-center gap-2">
                  <div 
                    className="w-3 h-3 rounded-full" 
                    style={{ backgroundColor: category.color }}
                  />
                  <span className="text-sm">{category.name}</span>
                  <span className="text-sm font-medium ml-auto">{category.value}%</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Top Products */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <BarChart3 className="h-5 w-5 text-purple-600" />
                Top Performing Products
              </CardTitle>
              <CardDescription>Best selling products by quantity and revenue</CardDescription>
            </div>
            <Button variant="outline" size="sm" onClick={() => navigate('/inventory/products')}>
              View All Products
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {mockTopProducts.map((product, index) => (
              <div key={index} className="flex items-center justify-between p-3 rounded-lg border">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
                    <span className="text-sm font-bold text-primary">#{index + 1}</span>
                  </div>
                  <div>
                    <h4 className="font-medium">{product.name}</h4>
                    <p className="text-sm text-muted-foreground">{product.sales} units sold</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="font-bold">{formatCurrency(product.revenue)}</p>
                  <Badge variant="secondary" className="text-xs">
                    {formatCurrency(Math.round(product.revenue / product.sales))}/unit
                  </Badge>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Clock className="h-5 w-5 text-orange-600" />
            Quick Actions
          </CardTitle>
          <CardDescription>Common reporting and export actions</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            <Button variant="outline" className="h-auto p-4 flex flex-col gap-2">
              <FileText className="h-5 w-5" />
              <span className="text-sm">Daily Sales</span>
            </Button>
            <Button variant="outline" className="h-auto p-4 flex flex-col gap-2">
              <Package className="h-5 w-5" />
              <span className="text-sm">Stock Report</span>
            </Button>
            <Button variant="outline" className="h-auto p-4 flex flex-col gap-2">
              <Users className="h-5 w-5" />
              <span className="text-sm">Customer List</span>
            </Button>
            <Button variant="outline" className="h-auto p-4 flex flex-col gap-2">
              <Download className="h-5 w-5" />
              <span className="text-sm">Export Data</span>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default Reports;
