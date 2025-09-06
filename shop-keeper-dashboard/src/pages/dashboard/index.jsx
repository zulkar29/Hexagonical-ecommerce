import { DollarSign, ShoppingCart, ArrowRight, TrendingUp, TrendingDown, Package, Download, Boxes, AlertTriangle, Star, Calendar, Truck, HelpCircle, Loader2, Eye, Edit, MoreVertical, Mail } from "lucide-react";
import SupportAnalytics from "@/components/support/SupportAnalytics";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { 
  Card, 
  CardContent, 
  CardDescription, 
  CardHeader, 
  CardTitle
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";
import { 
  useAdminDashboard, 
  useRecentOrders, 
  useProducts, 
  useRevenueAnalytics,
  useRecentActivity,
  useLowStockProducts,
  useQuickActionsData
} from "@/hooks/useApi";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

// Default fallback data for charts
const getDefaultRevenueData = () => [
  { name: "Jan", revenue: 0 },
  { name: "Feb", revenue: 0 },
  { name: "Mar", revenue: 0 },
  { name: "Apr", revenue: 0 },
  { name: "May", revenue: 0 },
  { name: "Jun", revenue: 0 }
];


const StatCard = ({ title, value, change, trend, icon, onClick, path }) => (
  <Card 
    className="cursor-pointer hover:shadow-md transition-all duration-200 hover:scale-105"
    onClick={onClick || (() => {})}
  >
    <CardContent className="pt-4 sm:pt-6">
      <div className="flex justify-between items-start gap-3">
        <div className="flex flex-col gap-1 min-w-0 flex-1">
          <span className="text-xs sm:text-sm font-medium text-muted-foreground truncate">{title}</span>
          <div className="flex flex-col sm:flex-row sm:items-baseline gap-1 sm:gap-2">
            <span className="text-lg sm:text-2xl font-bold tracking-tight">{value}</span>
            <div className="flex items-center gap-1 text-xs sm:text-sm">
              {trend === "up" ? (
                <TrendingUp className="h-3 w-3 sm:h-4 sm:w-4 text-emerald-500" />
              ) : (
                <TrendingDown className="h-3 w-3 sm:h-4 sm:w-4 text-red-500" />
              )}
              <span className={trend === "up" ? "text-emerald-500" : "text-red-500"}>
                {change}
              </span>
            </div>
          </div>
        </div>
        <div className="rounded-full p-2 sm:p-3 bg-primary/10 flex-shrink-0">
          {icon}
        </div>
      </div>
    </CardContent>
  </Card>
);

const OrderTable = ({ orders, navigate }) => {
  if (!orders || !Array.isArray(orders) || orders.length === 0) {
    return (
      <div className="text-center py-12 text-muted-foreground">
        <div className="w-16 h-16 mx-auto mb-4 bg-muted rounded-full flex items-center justify-center">
          <ShoppingCart className="w-8 h-8" />
        </div>
        <h3 className="text-lg font-medium mb-2">No Recent Orders</h3>
        <p className="text-sm text-muted-foreground mb-4">Orders will appear here once customers start purchasing</p>
        <Button 
          variant="outline" 
          size="sm"
          className="cursor-pointer"
          onClick={() => navigate('/orders')}
        >
          View All Orders
        </Button>
      </div>
    );
  }

  const getStatusConfig = (status) => {
    const statusMap = {
      'delivered': { variant: 'default', color: 'bg-emerald-100 text-emerald-800 border-emerald-200', icon: '✓' },
      'completed': { variant: 'default', color: 'bg-emerald-100 text-emerald-800 border-emerald-200', icon: '✓' },
      'processing': { variant: 'default', color: 'bg-blue-100 text-blue-800 border-blue-200', icon: '⟳' },
      'pending': { variant: 'secondary', color: 'bg-amber-100 text-amber-800 border-amber-200', icon: '⏱' },
      'cancelled': { variant: 'destructive', color: 'bg-red-100 text-red-800 border-red-200', icon: '✕' },
      'shipped': { variant: 'default', color: 'bg-purple-100 text-purple-800 border-purple-200', icon: '🚚' }
    };
    return statusMap[status?.toLowerCase()] || statusMap['pending'];
  };

  return (
    <div className="space-y-1">
      {orders.map((order) => {
        const statusConfig = getStatusConfig(order.status);
        return (
          <div 
            key={order.id} 
            className="group p-4 rounded-lg border bg-card hover:bg-muted/50 transition-colors cursor-pointer"
            onClick={() => navigate(`/orders/details/${order.id}`)}
          >
            <div className="flex items-center justify-between">
              {/* Left Section - Order Info */}
              <div className="flex items-center gap-4 flex-1">
                <div className="flex flex-col">
                  <div className="flex items-center gap-2">
                    <span className="font-semibold text-sm">
                      #{order.id?.toString().padStart(4, '0')}
                    </span>
                    <Badge className={`text-xs px-2 py-1 ${statusConfig.color} hover:${statusConfig.color}`}>
                      {statusConfig.icon} {order.status?.charAt(0).toUpperCase() + order.status?.slice(1) || 'Unknown'}
                    </Badge>
                  </div>
                  <span className="text-xs text-muted-foreground mt-1">
                    {new Date(order.created_at || order.date).toLocaleDateString('en-US', {
                      month: 'short',
                      day: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit'
                    })}
                  </span>
                </div>

                {/* Customer Info */}
                <div className="flex items-center gap-3 flex-1 min-w-0">
                  <Avatar className="h-10 w-10 ring-2 ring-background">
                    <AvatarFallback className="text-sm font-medium bg-primary/10 text-primary">
                      {order.user?.name?.[0] || order.customer?.[0] || 'U'}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex flex-col min-w-0">
                    <span className="font-medium text-sm truncate">
                      {order.user?.name || order.customer || 'Unknown Customer'}
                    </span>
                    <div className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Mail className="w-3 h-3" />
                      <span className="truncate max-w-[150px]">
                        {order.user?.email || 'No email'}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Right Section - Amount & Actions */}
              <div className="flex items-center gap-4">
                <div className="text-right">
                  <div className="font-bold text-lg">
                    ${Number(order.final_price || order.amount || 0).toFixed(2)}
                  </div>
                  <div className="text-xs text-muted-foreground">
                    {order.items?.length || 0} items
                  </div>
                </div>
                
                <div className="flex items-center gap-1">
                  <Button 
                    variant="ghost" 
                    size="sm" 
                    className="h-8 w-8 p-0"
                    onClick={(e) => {
                      e.stopPropagation();
                      navigate(`/orders/details/${order.id}`);
                    }}
                  >
                    <Eye className="h-4 w-4" />
                  </Button>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                        <MoreVertical className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-48">
                      <DropdownMenuItem 
                        className="flex items-center gap-2"
                        onClick={() => navigate(`/orders/details/${order.id}`)}
                      >
                        <Eye className="h-4 w-4" />
                        View Details
                      </DropdownMenuItem>
                      <DropdownMenuItem 
                        className="flex items-center gap-2"
                        onClick={() => navigate(`/orders/edit/${order.id}`)}
                      >
                        <Edit className="h-4 w-4" />
                        Update Status
                      </DropdownMenuItem>
                      <DropdownMenuItem className="flex items-center gap-2">
                        <Download className="h-4 w-4" />
                        Download Invoice
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};


const Dashboard = () => {
  const [selectedPeriod, setSelectedPeriod] = useState('monthly');
  const navigate = useNavigate();
  
  // API calls
  const { 
    data: dashboardData, 
    isLoading: isDashboardLoading, 
    isError: isDashboardError 
  } = useAdminDashboard();
  
  const { 
    data: recentOrdersData, 
    isLoading: isOrdersLoading 
  } = useRecentOrders();
  
  const { 
    data: productsData, 
    isLoading: isProductsLoading 
  } = useProducts({ per_page: 4, sort_by: 'created_at', sort_direction: 'desc' });
  
  const { 
    data: revenueData, 
    isLoading: isRevenueLoading 
  } = useRevenueAnalytics(selectedPeriod);
  
  const { 
    data: recentActivityData, 
    isLoading: isActivityLoading 
  } = useRecentActivity();
  
  const { 
    data: lowStockData, 
    isLoading: isLowStockLoading 
  } = useLowStockProducts();
  
  const { 
    data: quickActionsData, 
    isLoading: isQuickActionsLoading 
  } = useQuickActionsData();

  // Loading state
  if (isDashboardLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  // Error state
  if (isDashboardError) {
    return (
      <Card>
        <CardContent className="flex items-center justify-center h-[400px] text-destructive">
          Error loading dashboard data. Please try again later.
        </CardContent>
      </Card>
    );
  }

  const dashboard = dashboardData?.data || {};
  const recentOrders = Array.isArray(recentOrdersData?.data) ? recentOrdersData.data : [];
  const topProducts = Array.isArray(productsData?.data) ? productsData.data : [];
  const chartData = Array.isArray(revenueData?.data) ? revenueData.data : getDefaultRevenueData();
  const recentActivity = Array.isArray(recentActivityData?.data) ? recentActivityData.data : [];
  const lowStockProducts = Array.isArray(lowStockData?.data) ? lowStockData.data : [];
  const quickActions = quickActionsData?.data || {};

  // Calculate percentage change display
  const formatPercentage = (current, previous) => {
    if (!previous) return "+0%";
    const change = ((current - previous) / previous) * 100;
    return change >= 0 ? `+${change.toFixed(1)}%` : `${change.toFixed(1)}%`;
  };

  const stats = [
    {
      title: "Today's Revenue",
      value: dashboard.today_revenue ? `$${Number(dashboard.today_revenue).toFixed(2)}` : "$0.00",
      change: dashboard.revenue_change ? formatPercentage(dashboard.monthly_revenue, dashboard.last_month_revenue) : "+0%",
      trend: (dashboard.revenue_change || 0) >= 0 ? "up" : "down",
      icon: <DollarSign className="w-5 h-5 text-primary" />,
      path: "/reports",
      onClick: () => navigate('/reports')
    },
    {
      title: "Pending Orders",
      value: dashboard.pending_orders?.toString() || "0",
      change: `${dashboard.processing_orders || 0} processing`,
      trend: "up",
      icon: <ShoppingCart className="w-5 h-5 text-primary" />,
      path: "/orders?status=pending",
      onClick: () => navigate('/orders?status=pending')
    },
    {
      title: "Total Products",
      value: dashboard.total_products?.toString() || "0",
      change: `${dashboard.new_users_today || 0} new today`,
      trend: "up",
      icon: <Boxes className="w-5 h-5 text-primary" />,
      path: "/products",
      onClick: () => navigate('/products')
    },
    {
      title: "Low Stock Items",
      value: dashboard.low_stock_products?.toString() || "0",
      change: `${dashboard.out_of_stock_products || 0} out of stock`,
      trend: (dashboard.low_stock_products || 0) > 0 ? "up" : "down",
      icon: <AlertTriangle className="w-5 h-5 text-primary" />,
      path: "/products?filter=low_stock",
      onClick: () => navigate('/products?filter=low_stock')
    },
  ];

  const supportStats = [
    {
      title: "New Customers Today",
      value: dashboard.new_users_today?.toString() || "0",
      change: `${dashboard.total_users || 0} total`,
      trend: "up",
      icon: <HelpCircle className="w-5 h-5 text-primary" />,
      path: "/customers",
      onClick: () => navigate('/customers')
    },
  ];

  return (
    <div className="p-3 sm:p-4 lg:p-6 space-y-4 sm:space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:justify-between sm:items-center">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-sm sm:text-base text-muted-foreground">Welcome back to your store overview.</p>
        </div>
        <div className="flex flex-col gap-2 sm:flex-row sm:w-auto">
          <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
            <SelectTrigger className="w-full sm:w-[180px]">
              <SelectValue placeholder="Select period" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="daily">Daily</SelectItem>
              <SelectItem value="weekly">Weekly</SelectItem>
              <SelectItem value="monthly">Monthly</SelectItem>
              <SelectItem value="yearly">Yearly</SelectItem>
            </SelectContent>
          </Select>
          <Button className="w-full sm:w-auto cursor-pointer">
            <Download className="mr-2 h-4 w-4" />
            <span className="hidden xs:inline">Download Report</span>
            <span className="xs:hidden">Report</span>
          </Button>
        </div>
      </div>

      {/* Priority Action Alerts */}
      {(dashboard.pending_orders > 0 || dashboard.low_stock_products > 0) && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {dashboard.pending_orders > 0 && (
            <Card className="border-amber-200 bg-amber-50">
              <CardContent className="pt-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-amber-100 rounded-lg flex items-center justify-center">
                      <Truck className="w-5 h-5 text-amber-600" />
                    </div>
                    <div>
                      <h3 className="font-semibold text-amber-800">Pending Orders</h3>
                      <p className="text-sm text-amber-600">{dashboard.pending_orders} orders need processing</p>
                    </div>
                  </div>
                  <Button 
                    size="sm" 
                    className="bg-amber-600 hover:bg-amber-700 cursor-pointer"
                    onClick={() => navigate('/orders?status=pending')}
                  >
                    Process Now
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}
          {dashboard.low_stock_products > 0 && (
            <Card className="border-red-200 bg-red-50">
              <CardContent className="pt-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-red-100 rounded-lg flex items-center justify-center">
                      <AlertTriangle className="w-5 h-5 text-red-600" />
                    </div>
                    <div>
                      <h3 className="font-semibold text-red-800">Low Stock Alert</h3>
                      <p className="text-sm text-red-600">{dashboard.low_stock_products} items running low</p>
                    </div>
                  </div>
                  <Button 
                    size="sm" 
                    variant="destructive"
                    className="cursor-pointer"
                    onClick={() => navigate('/products?filter=low_stock')}
                  >
                    Manage Stock
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      )}

      {/* Main Grid Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-3 xl:grid-cols-4 gap-4 lg:gap-6">
        {/* Left Column - 3 cols */}
        <div className="lg:col-span-2 xl:col-span-3 space-y-4 lg:space-y-6 order-2 lg:order-1">
          {/* Stats Grid */}
          <div className="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-2 sm:gap-4">
            {stats.map((stat) => (
              <StatCard key={stat.title} {...stat} />
            ))}
            {supportStats.map((stat) => (
              <StatCard key={stat.title} {...stat} />
            ))}
          </div>

          {/* Revenue Chart */}
          <Card className="shadow-sm">
            <CardHeader className="pb-2">
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Revenue Overview</CardTitle>
                  <CardDescription>Daily revenue statistics</CardDescription>
                </div>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="outline" size="sm" className="cursor-pointer">
                      <Download className="mr-2 h-4 w-4" />
                      Export
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem className="cursor-pointer">Export as PDF</DropdownMenuItem>
                    <DropdownMenuItem className="cursor-pointer">Export as CSV</DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </CardHeader>
            <CardContent>
              {isRevenueLoading ? (
                <div className="flex items-center justify-center h-[300px]">
                  <Loader2 className="h-6 w-6 animate-spin" />
                </div>
              ) : (
                <ResponsiveContainer width="100%" height={300}>
                   <LineChart data={chartData}>
                     <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                     <XAxis dataKey="name" stroke="#888888" />
                     <YAxis stroke="#888888" />
                     <Tooltip contentStyle={{ background: "#fff", border: "1px solid #f0f0f0", borderRadius: "8px" }} />
                     <Line type="monotone" dataKey="revenue" stroke="#6366f1" strokeWidth={3} dot={{ fill: "#6366f1", r: 4 }} activeDot={{ r: 6 }} />
                   </LineChart>
                 </ResponsiveContainer>
               )}
             </CardContent>
           </Card>

          {/* Support Analytics */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Support Overview</CardTitle>
                  <CardDescription>Customer support metrics and performance</CardDescription>
                </div>
                <Button 
                  variant="outline" 
                  size="sm"
                  className="cursor-pointer"
                  onClick={() => navigate('/support')}
                >
                  View Details
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <SupportAnalytics />
            </CardContent>
          </Card>

          {/* Recent Orders */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Recent Orders</CardTitle>
                  <CardDescription>Latest customer orders and details</CardDescription>
                </div>
                <Button 
                  variant="outline" 
                  size="sm"
                  className="cursor-pointer"
                  onClick={() => navigate('/orders')}
                >
                  View All
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {isOrdersLoading ? (
                <div className="flex items-center justify-center h-[200px]">
                  <Loader2 className="h-6 w-6 animate-spin" />
                </div>
              ) : (
                <OrderTable orders={recentOrders.slice(0, 5)} navigate={navigate} />
              )}
            </CardContent>
          </Card>
        </div>

        {/* Right Column - 1 col */}
        <div className="space-y-4 lg:space-y-6 order-1 lg:order-2">
          {/* Quick Actions */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-lg flex items-center gap-2">
                <div className="w-8 h-8 bg-primary/10 rounded-lg flex items-center justify-center">
                  <Star className="w-4 h-4 text-primary" />
                </div>
                Quick Actions
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {isQuickActionsLoading ? (
                <div className="flex items-center justify-center h-[200px]">
                  <Loader2 className="h-6 w-6 animate-spin" />
                </div>
              ) : (
                <>
                  <Button 
                    className="w-full justify-start h-11 text-sm font-medium cursor-pointer" 
                    variant="outline"
                    onClick={() => navigate('/products/create')}
                  >
                    <Package className="mr-3 h-4 w-4 text-emerald-600" />
                    <span className="flex-1 text-left">Add New Product</span>
                    <ArrowRight className="h-4 w-4 text-muted-foreground" />
                  </Button>
                  <Button 
                    className="w-full justify-start h-11 text-sm font-medium cursor-pointer" 
                    variant="outline"
                    onClick={() => navigate('/orders')}
                  >
                    <Truck className="mr-3 h-4 w-4 text-blue-600" />
                    <div className="flex-1 text-left">
                      <div>Process Orders</div>
                      <div className="text-xs text-muted-foreground">
                        {quickActions.pending_orders || dashboard.pending_orders || 0} pending
                      </div>
                    </div>
                    <ArrowRight className="h-4 w-4 text-muted-foreground" />
                  </Button>
                  <Button 
                    className="w-full justify-start h-11 text-sm font-medium cursor-pointer" 
                    variant="outline"
                    onClick={() => navigate('/products?filter=low_stock')}
                  >
                    <AlertTriangle className="mr-3 h-4 w-4 text-amber-600" />
                    <div className="flex-1 text-left">
                      <div>Low Stock Items</div>
                      <div className="text-xs text-muted-foreground">
                        {quickActions.low_stock_count || dashboard.low_stock_products || 0} items
                      </div>
                    </div>
                    <ArrowRight className="h-4 w-4 text-muted-foreground" />
                  </Button>
                  <Button 
                    className="w-full justify-start h-11 text-sm font-medium cursor-pointer" 
                    variant="outline"
                    onClick={() => navigate('/discounts/create')}
                  >
                    <Calendar className="mr-3 h-4 w-4 text-purple-600" />
                    <span className="flex-1 text-left">Schedule Promotion</span>
                    <ArrowRight className="h-4 w-4 text-muted-foreground" />
                  </Button>
                </>
              )}
            </CardContent>
          </Card>


                {/* Low Stock Alert */}
          <Card className="border-red-200 bg-red-50">
            <CardHeader>
              <CardTitle className="text-red-700 flex items-center gap-2">
                <AlertTriangle className="h-5 w-5" />
                Low Stock Alert
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {isLowStockLoading ? (
                <div className="flex items-center justify-center h-[150px]">
                  <Loader2 className="h-6 w-6 animate-spin" />
                </div>
              ) : (
                lowStockProducts.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    <Package className="w-8 h-8 mx-auto mb-2" />
                    <p className="text-sm">No low stock items</p>
                  </div>
                ) : (
                  (lowStockProducts || []).map(product => (
                    <div key={product.id} className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center">
                        {product.featured_image || (product.images && product.images[0]) ? (
                          <img 
                            src={product.featured_image || product.images[0]} 
                            alt={product.name}
                            className="w-full h-full object-cover rounded-lg"
                          />
                        ) : (
                          <Package className="w-6 h-6 text-red-600" />
                        )}
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium truncate">{product.name}</p>
                        <p className="text-xs text-red-600">
                          {product.stock === 0 ? 'Out of stock' : `Only ${product.stock} left`}
                        </p>
                      </div>
                      <Button size="sm" variant="destructive">
                        Restock
                      </Button>
                    </div>
                  ))
                )
              )}
            </CardContent>
          </Card>

          

          {/* Top Products */}
          <Card>
            <CardHeader>
              <CardTitle>Top Products</CardTitle>
              <CardDescription>Best selling items this month</CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              {isProductsLoading ? (
                <div className="flex items-center justify-center h-[200px]">
                  <Loader2 className="h-6 w-6 animate-spin" />
                </div>
              ) : topProducts.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <div className="w-16 h-16 mx-auto mb-4 bg-muted rounded-full flex items-center justify-center">
                    <Package className="w-8 h-8" />
                  </div>
                  <h3 className="font-medium mb-2">No Products Found</h3>
                  <p className="text-sm mb-4">Add some products to see them here</p>
                  <Button 
                    variant="outline" 
                    size="sm"
                    className="cursor-pointer"
                    onClick={() => navigate('/products/create')}
                  >
                    Add Product
                  </Button>
                </div>
              ) : (
                (topProducts || []).slice(0, 4).map(product => (
                  <div key={product.id} className="group p-3 rounded-lg border bg-card hover:bg-muted/50 transition-all duration-200 hover:shadow-sm">
                    <div className="flex items-start gap-3">
                      <div className="relative w-14 h-14 bg-muted rounded-lg flex items-center justify-center overflow-hidden">
                        {product.featured_image || (product.images && product.images[0]) ? (
                          <img 
                            src={product.featured_image || product.images[0]} 
                            alt={product.name}
                            className="w-full h-full object-cover transition-transform group-hover:scale-105"
                          />
                        ) : (
                          <Package className="w-6 h-6 text-muted-foreground" />
                        )}
                        <div className="absolute top-1 right-1">
                          <div className={`w-2 h-2 rounded-full ${
                            product.stock > (product.min_stock_level || 10) ? 'bg-emerald-500' : 
                            product.stock > 0 ? 'bg-amber-500' : 'bg-red-500'
                          }`} />
                        </div>
                      </div>
                      
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between">
                          <div className="flex-1 min-w-0">
                            <h4 className="text-sm font-semibold truncate group-hover:text-primary transition-colors">
                              {product.name}
                            </h4>
                            <p className="text-xs text-muted-foreground mt-1">
                              {product.category?.name || 'No Category'}
                            </p>
                            
                            <div className="flex items-center gap-2 mt-2">
                              <div className="flex items-center text-amber-500">
                                <Star className="w-3 h-3 fill-current" />
                                <span className="text-xs ml-1 font-medium">
                                  {product.reviews?.length > 0 
                                    ? (product.reviews.reduce((acc, review) => acc + review.rating, 0) / product.reviews.length).toFixed(1)
                                    : '0.0'
                                  }
                                </span>
                              </div>
                              <span className="text-xs text-muted-foreground">
                                •
                              </span>
                              <span className="text-xs text-muted-foreground">
                                {product.stock} in stock
                              </span>
                            </div>
                          </div>
                          
                          <div className="text-right">
                            <p className="text-lg font-bold">
                              ${Number(product.price).toFixed(2)}
                            </p>
                            <Button variant="ghost" size="sm" className="h-6 px-2 text-xs opacity-0 group-hover:opacity-100 transition-opacity">
                              <Eye className="w-3 h-3 mr-1" />
                              View
                            </Button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </CardContent>
          </Card>

    

          {/* Recent Activity */}
          <Card>
            <CardHeader>
              <CardTitle>Recent Activity</CardTitle>
              <CardDescription>Latest actions in your store</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {isActivityLoading ? (
                <div className="flex items-center justify-center h-[200px]">
                  <Loader2 className="h-6 w-6 animate-spin" />
                </div>
              ) : recentActivity.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <div className="w-16 h-16 mx-auto mb-4 bg-muted rounded-full flex items-center justify-center">
                    <Calendar className="w-8 h-8" />
                  </div>
                  <h3 className="font-medium mb-2">No Recent Activity</h3>
                  <p className="text-sm">Activity will appear here as it happens</p>
                </div>
              ) : (
                (recentActivity || []).map((activity, i) => {
                  const getActivityIcon = (type) => {
                    switch (type) {
                      case 'order': return ShoppingCart;
                      case 'product': return Package;
                      case 'review': return Star;
                      case 'stock': return AlertTriangle;
                      case 'user': return HelpCircle;
                      default: return Calendar;
                    }
                  };
                  
                  const IconComponent = getActivityIcon(activity.type);
                  
                  const getActivityPath = (type) => {
                    switch (type) {
                      case 'order': return '/orders';
                      case 'product': return '/products';
                      case 'review': return '/reviews';
                      case 'stock': return '/products?filter=low_stock';
                      case 'user': return '/customers';
                      default: return '/logs';
                    }
                  };
                  
                  return (
                    <div 
                      key={activity.id || i} 
                      className="flex items-start gap-4 p-3 rounded-lg hover:bg-muted/50 transition-colors cursor-pointer"
                      onClick={() => navigate(getActivityPath(activity.type))}
                    >
                      <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
                        <IconComponent className="w-4 h-4 text-primary" />
                      </div>
                      <div className="flex-1">
                        <p className="text-sm">{activity.description || activity.text}</p>
                        <p className="text-xs text-muted-foreground">{activity.time_ago || activity.time}</p>
                      </div>
                    </div>
                  );
                })
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;