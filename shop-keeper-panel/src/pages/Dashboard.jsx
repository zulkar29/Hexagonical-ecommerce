import React, { useState, useEffect } from 'react';
import {
  BarChart3,
  TrendingUp,
  TrendingDown,
  Package,
  Users,
  ShoppingCart,
  AlertTriangle,
  Eye,
  Plus,
  ArrowRight,
  Calendar,
  Clock,
  Star,
  Zap,
  DollarSign,
  Target,
  Activity,
  ShoppingBag,
  CreditCard,
  Truck
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Progress } from '@/components/ui/progress';
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend
} from 'recharts';
import { useNavigate } from 'react-router-dom';

const Dashboard = () => {
  const [currentTime, setCurrentTime] = useState(new Date());
  const navigate = useNavigate();

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  // Mock data for charts
  const salesData = [
    { name: 'Jan', sales: 45000, orders: 120, customers: 89 },
    { name: 'Feb', sales: 52000, orders: 145, customers: 102 },
    { name: 'Mar', sales: 48000, orders: 135, customers: 95 },
    { name: 'Apr', sales: 61000, orders: 168, customers: 123 },
    { name: 'May', sales: 55000, orders: 152, customers: 108 },
    { name: 'Jun', sales: 67000, orders: 189, customers: 134 },
    { name: 'Jul', sales: 71000, orders: 201, customers: 145 }
  ];

  const weeklyData = [
    { day: 'Mon', sales: 12000, target: 15000 },
    { day: 'Tue', sales: 18000, target: 15000 },
    { day: 'Wed', sales: 14000, target: 15000 },
    { day: 'Thu', sales: 22000, target: 15000 },
    { day: 'Fri', sales: 25000, target: 15000 },
    { day: 'Sat', sales: 28000, target: 15000 },
    { day: 'Sun', sales: 16000, target: 15000 }
  ];

  const categoryData = [
    { name: 'Rice & Grains', value: 35, color: '#22c55e' },
    { name: 'Vegetables', value: 25, color: '#3b82f6' },
    { name: 'Dairy Products', value: 20, color: '#f59e0b' },
    { name: 'Spices', value: 12, color: '#ef4444' },
    { name: 'Others', value: 8, color: '#8b5cf6' }
  ];

  const topProducts = [
    { name: 'Basmati Rice 5kg', sold: 145, revenue: 72500, trend: 'up' },
    { name: 'Onion 1kg', sold: 234, revenue: 23400, trend: 'up' },
    { name: 'Milk 1L', sold: 189, revenue: 18900, trend: 'down' },
    { name: 'Potato 1kg', sold: 156, revenue: 12480, trend: 'up' },
    { name: 'Cooking Oil 1L', sold: 98, revenue: 14700, trend: 'up' }
  ];

  const recentOrders = [
    { id: '#1234', customer: 'Fatima Khan', amount: 2450, status: 'completed', time: '10 min ago' },
    { id: '#1235', customer: 'Ahmed Ali', amount: 1890, status: 'processing', time: '15 min ago' },
    { id: '#1236', customer: 'Rashida Begum', amount: 3200, status: 'completed', time: '23 min ago' },
    { id: '#1237', customer: 'Mohammad Hasan', amount: 1650, status: 'pending', time: '35 min ago' },
    { id: '#1238', customer: 'Nasir Ahmed', amount: 2780, status: 'completed', time: '1 hour ago' }
  ];

  const lowStockItems = [
    { name: 'Sugar 1kg', current: 5, minimum: 20, status: 'critical' },
    { name: 'Tea Leaves 500g', current: 8, minimum: 15, status: 'low' },
    { name: 'Salt 1kg', current: 12, minimum: 25, status: 'low' },
    { name: 'Garlic 1kg', current: 3, minimum: 10, status: 'critical' }
  ];

  // Stats cards data
  const stats = [
    {
      title: 'Today\'s Sales',
      value: '৳15,420',
      change: '+12.5%',
      trend: 'up',
      icon: DollarSign,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
      description: 'vs yesterday'
    },
    {
      title: 'Total Orders',
      value: '47',
      change: '+8.2%',
      trend: 'up',
      icon: ShoppingCart,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
      description: 'vs yesterday'
    },
    {
      title: 'New Customers',
      value: '12',
      change: '+15.3%',
      trend: 'up',
      icon: Users,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
      description: 'vs yesterday'
    },
    {
      title: 'Low Stock Items',
      value: '8',
      change: '+2',
      trend: 'down',
      icon: AlertTriangle,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
      description: 'need attention'
    }
  ];

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      {/* Welcome Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">
            Good {currentTime.getHours() < 12 ? 'Morning' : currentTime.getHours() < 17 ? 'Afternoon' : 'Evening'}, Mohammad! 👋
          </h1>
          <p className="text-muted-foreground mt-1">
            Here's what's happening in your store today, {currentTime.toLocaleDateString('en-US', { 
              weekday: 'long', 
              year: 'numeric', 
              month: 'long', 
              day: 'numeric' 
            })}
          </p>
        </div>
        <div className="flex space-x-3">
          <Button onClick={() => navigate('/pos')}>
            <Plus className="h-4 w-4 mr-2" />
            Quick Sale
          </Button>
          <Button variant="outline" onClick={() => navigate('/reports')}>
            <BarChart3 className="h-4 w-4 mr-2" />
            View Reports
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat, index) => (
          <Card key={index} className="relative overflow-hidden">
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div className="space-y-2">
                  <p className="text-sm font-medium text-muted-foreground">
                    {stat.title}
                  </p>
                  <p className="text-2xl font-bold text-foreground">
                    {stat.value}
                  </p>
                  <div className="flex items-center space-x-2">
                    <div className={cn(
                      "flex items-center space-x-1 px-2 py-1 rounded-full text-xs font-medium",
                      stat.trend === 'up' 
                        ? "bg-green-100 text-green-700" 
                        : "bg-red-100 text-red-700"
                    )}>
                      {stat.trend === 'up' ? (
                        <TrendingUp className="h-3 w-3" />
                      ) : (
                        <TrendingDown className="h-3 w-3" />
                      )}
                      <span>{stat.change}</span>
                    </div>
                    <span className="text-xs text-muted-foreground">
                      {stat.description}
                    </span>
                  </div>
                </div>
                <div className={cn("p-3 rounded-lg", stat.bgColor)}>
                  <stat.icon className={cn("h-6 w-6", stat.color)} />
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Sales Trend */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Sales Overview</CardTitle>
                <CardDescription>Monthly sales performance</CardDescription>
              </div>
              <Badge variant="secondary">Last 7 months</Badge>
            </div>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={salesData}>
                <defs>
                  <linearGradient id="salesGradient" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#22c55e" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#22c55e" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip 
                  formatter={(value, name) => [formatCurrency(value), name]}
                  labelClassName="text-foreground"
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

        {/* Sales by Category */}
        <Card>
          <CardHeader>
            <CardTitle>Sales by Category</CardTitle>
            <CardDescription>Product category breakdown</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={categoryData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={100}
                  paddingAngle={5}
                  dataKey="value"
                >
                  {categoryData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip formatter={(value) => [`${value}%`, 'Percentage']} />
              </PieChart>
            </ResponsiveContainer>
            <div className="space-y-2 mt-4">
              {categoryData.map((category, index) => (
                <div key={index} className="flex items-center justify-between text-sm">
                  <div className="flex items-center space-x-2">
                    <div 
                      className="w-3 h-3 rounded-full" 
                      style={{ backgroundColor: category.color }}
                    ></div>
                    <span className="text-muted-foreground">{category.name}</span>
                  </div>
                  <span className="font-medium">{category.value}%</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Weekly Performance */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Weekly Performance</CardTitle>
              <CardDescription>Daily sales vs targets this week</CardDescription>
            </div>
            <div className="flex items-center space-x-4 text-sm">
              <div className="flex items-center space-x-2">
                <div className="w-3 h-3 bg-primary rounded-full"></div>
                <span>Actual Sales</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-3 h-3 bg-muted rounded-full"></div>
                <span>Target</span>
              </div>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={weeklyData}>
              <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
              <XAxis dataKey="day" />
              <YAxis />
              <Tooltip formatter={(value) => [formatCurrency(value), '']} />
              <Bar dataKey="target" fill="#e5e7eb" radius={4} />
              <Bar dataKey="sales" fill="#22c55e" radius={4} />
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Bottom Section */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Top Products */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Top Products</CardTitle>
                <CardDescription>Best performing items</CardDescription>
              </div>
              <Button variant="ghost" size="sm">
                <Eye className="h-4 w-4" />
              </Button>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            {topProducts.map((product, index) => (
              <div key={index} className="flex items-center justify-between p-3 bg-muted/30 rounded-lg">
                <div className="space-y-1">
                  <p className="font-medium text-sm">{product.name}</p>
                  <div className="flex items-center space-x-4 text-xs text-muted-foreground">
                    <span>Sold: {product.sold}</span>
                    <span>{formatCurrency(product.revenue)}</span>
                  </div>
                </div>
                <div className={cn(
                  "p-1 rounded-full",
                  product.trend === 'up' 
                    ? "bg-green-100 text-green-600" 
                    : "bg-red-100 text-red-600"
                )}>
                  {product.trend === 'up' ? (
                    <TrendingUp className="h-3 w-3" />
                  ) : (
                    <TrendingDown className="h-3 w-3" />
                  )}
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        {/* Recent Orders */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Recent Orders</CardTitle>
                <CardDescription>Latest customer orders</CardDescription>
              </div>
              <Button variant="ghost" size="sm">
                View All
                <ArrowRight className="h-4 w-4 ml-2" />
              </Button>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            {recentOrders.map((order, index) => (
              <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                <div className="flex items-center space-x-3">
                  <Avatar className="h-8 w-8">
                    <AvatarFallback className="text-xs">
                      {order.customer.split(' ').map(n => n[0]).join('')}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-medium text-sm">{order.customer}</p>
                    <p className="text-xs text-muted-foreground">{order.id} • {order.time}</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="font-medium text-sm">{formatCurrency(order.amount)}</p>
                  <Badge 
                    variant={
                      order.status === 'completed' ? 'default' : 
                      order.status === 'processing' ? 'secondary' : 'outline'
                    }
                    className="text-xs"
                  >
                    {order.status}
                  </Badge>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        {/* Low Stock Alert */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center space-x-2">
                  <AlertTriangle className="h-5 w-5 text-orange-500" />
                  <span>Stock Alerts</span>
                </CardTitle>
                <CardDescription>Items running low</CardDescription>
              </div>
              <Badge variant="destructive">{lowStockItems.length}</Badge>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            {lowStockItems.map((item, index) => (
              <div key={index} className="space-y-2 p-3 border rounded-lg">
                <div className="flex items-center justify-between">
                  <p className="font-medium text-sm">{item.name}</p>
                  <Badge 
                    variant={item.status === 'critical' ? 'destructive' : 'secondary'}
                    className="text-xs"
                  >
                    {item.status}
                  </Badge>
                </div>
                <div className="space-y-1">
                  <div className="flex justify-between text-xs">
                    <span>Current: {item.current}</span>
                    <span>Min: {item.minimum}</span>
                  </div>
                  <Progress 
                    value={(item.current / item.minimum) * 100} 
                    className="h-2"
                  />
                </div>
              </div>
            ))}
            <Button className="w-full" size="sm">
              <Package className="h-4 w-4 mr-2" />
              Restock Items
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions Footer */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="font-semibold text-lg">Quick Actions</h3>
              <p className="text-sm text-muted-foreground">Common tasks to boost your productivity</p>
            </div>
            <div className="flex space-x-3">
              <Button variant="outline" onClick={() => navigate('/inventory/products')}>
                <ShoppingBag className="h-4 w-4 mr-2" />
                Add Product
              </Button>
              <Button variant="outline" onClick={() => navigate('/customers')}>
                <Users className="h-4 w-4 mr-2" />
                Add Customer
              </Button>
              <Button variant="outline" onClick={() => navigate('/sales/invoices')}>
                <CreditCard className="h-4 w-4 mr-2" />
                Create Invoice
              </Button>
              <Button onClick={() => navigate('/reports/sales')}>
                <Activity className="h-4 w-4 mr-2" />
                View Analytics
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default Dashboard;