import React, { useState } from 'react';
import { 
  Search, 
  Bell, 
  Settings, 
  User, 
  LogOut,
  Menu,
  Sun,
  Moon,
  ChevronDown,
  Calculator,
  Calendar,
  Clock,
  Wifi,
  WifiOff,
  Plus,
  ShoppingCart,
  Package,
  Users,
  TrendingUp,
  Globe,
  HelpCircle
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import i18n from '@/i18n';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

const DashboardHeader = ({ sidebarCollapsed = false, onToggleSidebar }) => {
  const [isOnline, setIsOnline] = useState(true);
  const [theme, setTheme] = useState('light');
  const [searchOpen, setSearchOpen] = useState(false);
  const [notificationsOpen, setNotificationsOpen] = useState(false);
  const navigate = useNavigate();
  const { t } = useTranslation();

  // Mock data for search suggestions
  const searchSuggestions = [
    { type: 'product', label: 'Rice 5kg', icon: Package, action: () => navigate('/inventory/products/1') },
    { type: 'customer', label: 'Fatima Khan', icon: Users, action: () => navigate('/customers') },
    { type: 'product', label: 'Dal 1kg', icon: Package, action: () => navigate('/inventory/products/2') },
    { type: 'customer', label: 'Ahmed Ali', icon: Users, action: () => navigate('/customers') },
    { type: 'action', label: 'Add New Product', icon: Plus, action: () => navigate('/inventory/products/create') },
    { type: 'action', label: 'View Products', icon: Package, action: () => navigate('/inventory/products') },
    { type: 'action', label: 'View Sales Report', icon: TrendingUp, action: () => navigate('/reports/sales') },
    { type: 'action', label: 'Create Order', icon: ShoppingCart, action: () => navigate('/orders') },
  ];

  // Mock notifications
  const notifications = [
    {
      id: 1,
      type: 'alert',
      title: 'Low Stock Alert',
      message: 'Rice 5kg is running low (5 items left)',
      time: '2 min ago',
      unread: true
    },
    {
      id: 2,
      type: 'order',
      title: 'New Order',
      message: 'Order #1234 received from Fatima Khan',
      time: '5 min ago',
      unread: true
    },
    {
      id: 3,
      type: 'payment',
      title: 'Payment Received',
      message: '৳2,500 received via bKash',
      time: '10 min ago',
      unread: false
    },
    {
      id: 4,
      type: 'system',
      title: 'Daily Backup Complete',
      message: 'Your data has been backed up successfully',
      time: '1 hour ago',
      unread: false
    }
  ];

  const unreadCount = notifications.filter(n => n.unread).length;

  // Current time and date
  const [currentTime, setCurrentTime] = useState(new Date());
  
  React.useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  const formatTime = (date) => {
    return date.toLocaleTimeString('en-US', { 
      hour12: true, 
      hour: 'numeric', 
      minute: '2-digit' 
    });
  };

  const formatDate = (date) => {
    return date.toLocaleDateString('en-US', { 
      weekday: 'long',
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    });
  };

  return (
    <header className={cn(
      "fixed top-0 left-0 right-0 z-40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 border-b border-border transition-all duration-300",
      sidebarCollapsed ? "ml-16" : "ml-64"
    )}>
      <div className="flex h-16 items-center px-6">
        
        {/* Mobile Menu Toggle */}
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="ghost" size="sm" className="md:hidden mr-2">
              <Menu className="h-5 w-5" />
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-64 p-0">
            <SheetHeader className="p-4 border-b">
              <SheetTitle className="flex items-center space-x-2">
                <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
                  <Package className="h-4 w-4 text-primary-foreground" />
                </div>
                <span>ShopOS</span>
              </SheetTitle>
              <SheetDescription>
                Business Management System
              </SheetDescription>
            </SheetHeader>
            {/* Mobile navigation would go here */}
          </SheetContent>
        </Sheet>

        {/* Page Title and Breadcrumb */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center space-x-2">
            <h1 className="text-xl font-semibold text-foreground truncate">
              {t('Dashboard')}
            </h1>
            <span className="text-muted-foreground">•</span>
            <span className="text-sm text-muted-foreground">
              {t('Overview')}
            </span>
          </div>
          <p className="text-sm text-muted-foreground hidden sm:block">
            {t("Welcome back! Here's what's happening in your store today.")}
          </p>
        </div>

        {/* Quick Actions */}
        <div className="flex items-center space-x-2 mr-4">
          <Button size="sm" className="hidden sm:flex">
            <Plus className="h-4 w-4 mr-2" />
            Quick Sale
          </Button>
          
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="hidden md:flex">
                <Plus className="h-4 w-4 mr-2" />
                Add New
                <ChevronDown className="h-4 w-4 ml-2" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
              <DropdownMenuItem onClick={() => navigate('/inventory/products/create')}>
                <Package className="h-4 w-4 mr-2" />
                Add Product
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => navigate('/customers')}>
                <Users className="h-4 w-4 mr-2" />
                Add Customer
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => navigate('/sales/orders')}>
                <ShoppingCart className="h-4 w-4 mr-2" />
                Create Order
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        {/* Search */}
        <div className="relative mr-4">
          <Popover open={searchOpen} onOpenChange={setSearchOpen}>
            <PopoverTrigger asChild>
              <Button
                variant="outline"
                className="w-64 justify-start text-sm text-muted-foreground hidden lg:flex"
              >
                <Search className="h-4 w-4 mr-2" />
                Search products, customers...
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80 p-0" align="start">
              <Command>
                <CommandInput placeholder="Search anything..." />
                <CommandList>
                  <CommandEmpty>No results found.</CommandEmpty>
                  <CommandGroup heading="Suggestions">
                    {searchSuggestions.map((item, index) => (
                      <CommandItem 
                        key={index} 
                        onSelect={() => {
                          if (item.action) {
                            item.action();
                            setSearchOpen(false);
                          }
                        }}
                      >
                        <item.icon className="h-4 w-4 mr-2" />
                        <span>{item.label}</span>
                        <Badge variant="secondary" className="ml-auto text-xs">
                          {item.type}
                        </Badge>
                      </CommandItem>
                    ))}
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
          
          {/* Mobile Search */}
          <Button variant="ghost" size="sm" className="lg:hidden">
            <Search className="h-5 w-5" />
          </Button>
        </div>

        {/* Date & Time */}
        <div className="hidden xl:flex items-center space-x-4 mr-4 px-3 py-2 bg-muted/50 rounded-lg">
          <div className="flex items-center space-x-2">
            <Calendar className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm font-medium">
              {formatDate(currentTime)}
            </span>
          </div>
          <div className="flex items-center space-x-2">
            <Clock className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm font-mono">
              {formatTime(currentTime)}
            </span>
          </div>
        </div>

        {/* Connection Status */}
        <div className="flex items-center space-x-2 mr-4">
          <div className={cn(
            "flex items-center space-x-1 px-2 py-1 rounded-full text-xs",
            isOnline 
              ? "bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300" 
              : "bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300"
          )}>
            {isOnline ? (
              <Wifi className="h-3 w-3" />
            ) : (
              <WifiOff className="h-3 w-3" />
            )}
            <span className="hidden sm:inline">
              {isOnline ? 'Online' : 'Offline'}
            </span>
          </div>
        </div>

        {/* Notifications */}
        <DropdownMenu open={notificationsOpen} onOpenChange={setNotificationsOpen}>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="sm" className="relative mr-2">
              <Bell className="h-5 w-5" />
              {unreadCount > 0 && (
                <Badge className="absolute -top-1 -right-1 h-5 w-5 flex items-center justify-center p-0 text-xs">
                  {unreadCount}
                </Badge>
              )}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-80">
            <DropdownMenuLabel className="flex items-center justify-between">
              Notifications
              <Badge variant="secondary">{unreadCount} new</Badge>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <div className="max-h-64 overflow-y-auto">
              {notifications.map((notification) => (
                <DropdownMenuItem key={notification.id} className="flex-col items-start p-3">
                  <div className="flex items-start justify-between w-full">
                    <div className="flex-1">
                      <p className={cn(
                        "text-sm font-medium",
                        notification.unread && "text-primary"
                      )}>
                        {notification.title}
                      </p>
                      <p className="text-sm text-muted-foreground mt-1">
                        {notification.message}
                      </p>
                      <p className="text-xs text-muted-foreground mt-1">
                        {notification.time}
                      </p>
                    </div>
                    {notification.unread && (
                      <div className="w-2 h-2 bg-primary rounded-full mt-1"></div>
                    )}
                  </div>
                </DropdownMenuItem>
              ))}
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="text-center text-sm text-primary">
              View All Notifications
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        {/* Settings & Theme */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="sm" className="mr-2">
              <Settings className="h-5 w-5" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}>
              {theme === 'light' ? (
                <Moon className="h-4 w-4 mr-2" />
              ) : (
                <Sun className="h-4 w-4 mr-2" />
              )}
              {theme === 'light' ? t('Dark Mode') : t('Light Mode')}
            </DropdownMenuItem>
            <DropdownMenuItem>
              <Calculator className="h-4 w-4 mr-2" />
              {t('Calculator')}
            </DropdownMenuItem>
            <DropdownMenuItem>
              <Globe className="h-4 w-4 mr-2" />
              {t('Language')}
              <Select
                value={i18n.language}
                onValueChange={lng => i18n.changeLanguage(lng)}
                className="ml-2 w-24"
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="en">{t('English')}</SelectItem>
                  <SelectItem value="bn">{t('Bengali')}</SelectItem>
                </SelectContent>
              </Select>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => navigate('/help')}>
              <HelpCircle className="h-4 w-4 mr-2" />
              {t('Help & Support')}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        {/* User Profile */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="relative h-10 w-10 rounded-full">
              <Avatar className="h-10 w-10">
                <AvatarImage src={import.meta.env.BASE_URL + "api/placeholder/40/40"} alt="Mohammad Rahman" />
                <AvatarFallback className="bg-primary text-primary-foreground">
                  MR
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="w-56" align="end" forceMount>
            <DropdownMenuLabel className="font-normal">
              <div className="flex flex-col space-y-1">
                <p className="text-sm font-medium leading-none">Mohammad Rahman</p>
                <p className="text-xs leading-none text-muted-foreground">
                  rahman.shop@gmail.com
                </p>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => navigate('/settings')}>
              <User className="h-4 w-4 mr-2" />
              Profile Settings
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => navigate('/settings')}>
              <Settings className="h-4 w-4 mr-2" />
              Account Settings
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => navigate('/help')}>
              <HelpCircle className="h-4 w-4 mr-2" />
              Help & Support
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="text-red-600">
              <LogOut className="h-4 w-4 mr-2" />
              Sign Out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
};

export default DashboardHeader;