import React, { useState, useEffect } from 'react';
import {
  Store,
  Bell,
  Settings,
  User,
  LogOut,
  Shield,
  Moon,
  Sun,
  Search,
  Menu,
  X,
  ChevronDown,
  Activity,
  AlertCircle,
  CheckCircle,
  Info,
  Zap,
  RefreshCw,
  Users,
  DollarSign,
  TrendingUp,
  Globe
} from 'lucide-react';
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
  DropdownMenuShortcut,
} from '@/components/ui/dropdown-menu';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

const AdminHeader = ({ 
  user = { name: 'Admin User', email: 'admin@shopowner.com', role: 'Super Admin', avatar: null },
  onMenuToggle,
  isMobileMenuOpen = false,
  theme = 'light',
  onThemeToggle,
  onLogout,
  platformStats = {
    totalTenants: 156,
    activeUsers: 1240,
    monthlyRevenue: 245000,
    systemHealth: 'operational'
  }
}) => {
  const [currentTime, setCurrentTime] = useState(new Date());
  const [searchQuery, setSearchQuery] = useState('');
  const [notifications, setNotifications] = useState([
    {
      id: 1,
      type: 'alert',
      title: 'High Priority Ticket',
      message: 'Rahman Electronics reported POS system issues',
      time: '5 minutes ago',
      read: false,
      icon: AlertCircle,
      color: 'text-red-500'
    },
    {
      id: 2,
      type: 'success',
      title: 'New Tenant Registered',
      message: 'Modern Pharmacy successfully onboarded',
      time: '15 minutes ago',
      read: false,
      icon: CheckCircle,
      color: 'text-green-500'
    },
    {
      id: 3,
      type: 'info',
      title: 'Monthly Report Ready',
      message: 'June revenue report is now available',
      time: '1 hour ago',
      read: true,
      icon: Info,
      color: 'text-blue-500'
    },
    {
      id: 4,
      type: 'warning',
      title: 'Server Performance',
      message: 'Database response time slightly elevated',
      time: '2 hours ago',
      read: true,
      icon: Activity,
      color: 'text-yellow-500'
    },
    {
      id: 5,
      type: 'success',
      title: 'Payment Received',
      message: 'Tech Solutions - ৳10,000 subscription payment',
      time: '3 hours ago',
      read: true,
      icon: DollarSign,
      color: 'text-green-500'
    }
  ]);

  const [systemStatus, setSystemStatus] = useState({
    api: 'operational',
    database: 'operational',
    payments: 'operational',
    notifications: 'degraded'
  });

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;
  
  const unreadNotifications = notifications.filter(n => !n.read).length;

  const getStatusColor = (status) => {
    switch (status) {
      case 'operational': return 'text-green-500';
      case 'degraded': return 'text-yellow-500';
      case 'outage': return 'text-red-500';
      default: return 'text-gray-500';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'operational': return CheckCircle;
      case 'degraded': return AlertCircle;
      case 'outage': return X;
      default: return Activity;
    }
  };

  const markNotificationAsRead = (notificationId) => {
    setNotifications(prev => 
      prev.map(n => 
        n.id === notificationId ? { ...n, read: true } : n
      )
    );
  };

  const markAllAsRead = () => {
    setNotifications(prev => prev.map(n => ({ ...n, read: true })));
  };

  return (
    <header className="fixed top-0 left-0 right-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-16 items-center px-4 md:px-6">
        {/* Mobile Menu Toggle */}
        <Button
          variant="ghost"
          size="sm"
          className="mr-2 md:hidden"
          onClick={onMenuToggle}
        >
          {isMobileMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
        </Button>

        {/* Logo and Brand */}
        <div className="flex items-center space-x-2 md:space-x-4">
          <div className="flex items-center space-x-2">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary">
              <Store className="h-6 w-6 text-primary-foreground" />
            </div>
            <div className="hidden sm:block">
              <h1 className="text-lg font-semibold text-foreground">Shop Owner SaaS</h1>
              <p className="text-xs text-muted-foreground">Admin Control Panel</p>
            </div>
          </div>
        </div>

        {/* Search Bar - Hidden on mobile */}
        <div className="mx-4 flex-1 max-w-md hidden md:block">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Search tenants, payments, tickets..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 pr-4"
            />
          </div>
        </div>

        {/* Right Side Items */}
        <div className="flex items-center space-x-2 md:space-x-4 ml-auto">
          {/* Current Time */}
          <div className="hidden lg:block text-right">
            <p className="text-sm font-medium text-foreground">
              {currentTime.toLocaleDateString('en-US', { 
                weekday: 'short', 
                month: 'short', 
                day: 'numeric' 
              })}
            </p>
            <p className="text-xs text-muted-foreground">
              {currentTime.toLocaleTimeString('en-US', { 
                hour12: true, 
                hour: 'numeric', 
                minute: '2-digit' 
              })}
            </p>
          </div>

          {/* Quick Stats - Hidden on mobile */}
          <div className="hidden xl:flex items-center space-x-4">
            <Separator orientation="vertical" className="h-8" />
            <div className="flex items-center space-x-4 text-sm">
              <div className="flex items-center space-x-1">
                <Users className="h-4 w-4 text-muted-foreground" />
                <span className="font-medium">{platformStats.totalTenants}</span>
                <span className="text-muted-foreground">tenants</span>
              </div>
              <div className="flex items-center space-x-1">
                <TrendingUp className="h-4 w-4 text-green-500" />
                <span className="font-medium">{formatCurrency(platformStats.monthlyRevenue)}</span>
                <span className="text-muted-foreground">MRR</span>
              </div>
            </div>
          </div>

          {/* System Status */}
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" size="sm" className="relative">
                <Activity className={`h-4 w-4 ${getStatusColor(platformStats.systemHealth)}`} />
                <span className="sr-only">System Status</span>
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80" align="end">
              <div className="space-y-4">
                <div className="space-y-2">
                  <h4 className="font-medium">System Status</h4>
                  <p className="text-sm text-muted-foreground">
                    Current status of all services
                  </p>
                </div>
                <div className="space-y-3">
                  {Object.entries(systemStatus).map(([service, status]) => {
                    const StatusIcon = getStatusIcon(status);
                    return (
                      <div key={service} className="flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                          <StatusIcon className={`h-4 w-4 ${getStatusColor(status)}`} />
                          <span className="text-sm capitalize">{service}</span>
                        </div>
                        <Badge 
                          variant={status === 'operational' ? 'default' : status === 'degraded' ? 'secondary' : 'destructive'}
                          className="text-xs"
                        >
                          {status}
                        </Badge>
                      </div>
                    );
                  })}
                </div>
                <Separator />
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Last updated</span>
                  <span className="text-sm">
                    {currentTime.toLocaleTimeString('en-US', { 
                      hour: 'numeric', 
                      minute: '2-digit' 
                    })}
                  </span>
                </div>
              </div>
            </PopoverContent>
          </Popover>

          {/* Notifications */}
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" size="sm" className="relative">
                <Bell className="h-4 w-4" />
                {unreadNotifications > 0 && (
                  <Badge className="absolute -top-1 -right-1 h-5 w-5 p-0 text-xs bg-red-500 text-white">
                    {unreadNotifications}
                  </Badge>
                )}
                <span className="sr-only">Notifications</span>
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-96" align="end">
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <h4 className="font-medium">Notifications</h4>
                  {unreadNotifications > 0 && (
                    <Button variant="ghost" size="sm" onClick={markAllAsRead}>
                      Mark all read
                    </Button>
                  )}
                </div>
                <div className="space-y-3 max-h-96 overflow-y-auto">
                  {notifications.length === 0 ? (
                    <div className="text-center py-6">
                      <Bell className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
                      <p className="text-sm text-muted-foreground">No notifications</p>
                    </div>
                  ) : (
                    notifications.map((notification) => {
                      const IconComponent = notification.icon;
                      return (
                        <div
                          key={notification.id}
                          className={`p-3 rounded-lg border cursor-pointer transition-colors hover:bg-muted/50 ${
                            !notification.read ? 'bg-muted/30' : ''
                          }`}
                          onClick={() => markNotificationAsRead(notification.id)}
                        >
                          <div className="flex items-start space-x-3">
                            <IconComponent className={`h-4 w-4 mt-0.5 ${notification.color}`} />
                            <div className="flex-1 space-y-1">
                              <div className="flex items-center justify-between">
                                <p className="text-sm font-medium">{notification.title}</p>
                                {!notification.read && (
                                  <div className="h-2 w-2 bg-blue-500 rounded-full" />
                                )}
                              </div>
                              <p className="text-xs text-muted-foreground">
                                {notification.message}
                              </p>
                              <p className="text-xs text-muted-foreground">
                                {notification.time}
                              </p>
                            </div>
                          </div>
                        </div>
                      );
                    })
                  )}
                </div>
                <Separator />
                <div className="flex justify-center">
                  <Button variant="ghost" size="sm" className="w-full">
                    View all notifications
                  </Button>
                </div>
              </div>
            </PopoverContent>
          </Popover>

          {/* Theme Toggle */}
          <Button
            variant="ghost"
            size="sm"
            onClick={onThemeToggle}
          >
            {theme === 'light' ? (
              <Moon className="h-4 w-4" />
            ) : (
              <Sun className="h-4 w-4" />
            )}
            <span className="sr-only">Toggle theme</span>
          </Button>

          {/* Settings */}
          <Button variant="ghost" size="sm">
            <Settings className="h-4 w-4" />
            <span className="sr-only">Settings</span>
          </Button>

          <Separator orientation="vertical" className="h-6" />

          {/* User Menu */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="relative h-10 w-10 rounded-full">
                <Avatar className="h-10 w-10">
                  <AvatarImage src={user.avatar} alt={user.name} />
                  <AvatarFallback className="bg-primary text-primary-foreground">
                    {user.name.split(' ').map(n => n[0]).join('').toUpperCase()}
                  </AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-80" align="end" forceMount>
              <DropdownMenuLabel className="font-normal">
                <div className="flex flex-col space-y-2">
                  <div className="flex items-center space-x-3">
                    <Avatar className="h-12 w-12">
                      <AvatarImage src={user.avatar} alt={user.name} />
                      <AvatarFallback className="bg-primary text-primary-foreground">
                        {user.name.split(' ').map(n => n[0]).join('').toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex-1">
                      <p className="text-sm font-medium leading-none">{user.name}</p>
                      <p className="text-xs leading-none text-muted-foreground mt-1">
                        {user.email}
                      </p>
                      <Badge variant="secondary" className="mt-1 text-xs">
                        <Shield className="h-3 w-3 mr-1" />
                        {user.role}
                      </Badge>
                    </div>
                  </div>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              
              {/* Quick Platform Stats */}
              <div className="p-2">
                <div className="grid grid-cols-2 gap-3 mb-2">
                  <div className="text-center p-2 bg-muted/30 rounded">
                    <p className="text-lg font-semibold">{platformStats.totalTenants}</p>
                    <p className="text-xs text-muted-foreground">Active Tenants</p>
                  </div>
                  <div className="text-center p-2 bg-muted/30 rounded">
                    <p className="text-lg font-semibold">{formatCurrency(platformStats.monthlyRevenue)}</p>
                    <p className="text-xs text-muted-foreground">Monthly Revenue</p>
                  </div>
                </div>
              </div>
              
              <DropdownMenuSeparator />
              
              <DropdownMenuItem>
                <User className="mr-2 h-4 w-4" />
                <span>Profile Settings</span>
                <DropdownMenuShortcut>⇧⌘P</DropdownMenuShortcut>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Settings className="mr-2 h-4 w-4" />
                <span>Admin Settings</span>
                <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Activity className="mr-2 h-4 w-4" />
                <span>System Logs</span>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Globe className="mr-2 h-4 w-4" />
                <span>Platform Status</span>
              </DropdownMenuItem>
              
              <DropdownMenuSeparator />
              
              <DropdownMenuItem onClick={onLogout} className="text-red-600">
                <LogOut className="mr-2 h-4 w-4" />
                <span>Log out</span>
                <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {/* Mobile Search Bar */}
      <div className="px-4 pb-3 md:hidden">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            type="search"
            placeholder="Search..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10 pr-4"
          />
        </div>
      </div>
    </header>
  );
};

export default AdminHeader;