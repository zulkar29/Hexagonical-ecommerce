import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Package, 
  ShoppingCart, 
  Users, 
  BarChart3, 
  FileText, 
  Settings, 
  LogOut,
  ChevronDown,
  ChevronRight,
  Store,
  TrendingUp,
  Bell,
  User,
  CreditCard,
  Truck,
  Tags,
  Calendar,
  HelpCircle,
  DollarSign,
  Receipt,
  Wallet
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useTranslation } from 'react-i18next';

const DashboardSidebar = () => {
  const [activeItem, setActiveItem] = useState('dashboard');
  const [expandedItems, setExpandedItems] = useState(['inventory', 'finance']);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const { t } = useTranslation();

  const toggleExpanded = (item) => {
    setExpandedItems(prev => 
      prev.includes(item) 
        ? prev.filter(i => i !== item)
        : [...prev, item]
    );
  };

  const menuItems = [
    {
      id: 'dashboard',
      label: t('Dashboard'),
      icon: LayoutDashboard,
      href: '/dashboard',
      badge: null
    },
    {
      id: 'pos',
      label: t('Point of Sale'),
      icon: ShoppingCart,
      href: '/pos',
      badge: null,
      highlight: true
    },
    {
      id: 'inventory',
      label: t('Inventory'),
      icon: Package,
      expandable: true,
      badge: '5',
      children: [
        { id: 'products', label: t('Products'), href: '/inventory/products' },
        { id: 'categories', label: t('Categories'), href: '/inventory/categories' },
        { id: 'suppliers', label: t('Suppliers'), href: '/inventory/suppliers' },
        { id: 'stock-alerts', label: t('Stock Alerts'), href: '/inventory/alerts', badge: '3' }
      ]
    },
    {
      id: 'customers',
      label: t('Customers'),
      icon: Users,
      href: '/customers',
      badge: null
    },
    {
      id: 'orders',
      label: t('Orders'),
      icon: CreditCard,
      href: '/orders',
      badge: null
    },
    {
      id: 'finance',
      label: t('Financial Management'),
      icon: DollarSign,
      expandable: true,
      children: [
        { id: 'accounts-payable', label: t('Accounts Payable'), href: '/finance/payable', icon: Receipt },
        { id: 'accounts-receivable', label: t('Accounts Receivable'), href: '/finance/receivable', icon: FileText },
        { id: 'petty-cash', label: t('Petty Cash'), href: '/finance/petty-cash', icon: Wallet }
      ]
    },
    {
      id: 'reports',
      label: t('Reports & Analytics'),
      icon: BarChart3,
      expandable: true,
      children: [
        { id: 'sales-report', label: t('Sales Report'), href: '/reports/sales' },
        { id: 'inventory-report', label: t('Inventory Report'), href: '/reports/inventory' },
        { id: 'customer-insights', label: t('Customer Insights'), href: '/reports/customers' }
      ]
    }
  ];

  const bottomMenuItems = [
    {
      id: 'settings',
      label: t('Settings'),
      icon: Settings,
      href: '/settings'
    },
    {
      id: 'help',
      label: t('Help & Support'),
      icon: HelpCircle,
      href: '/help'
    }
  ];

  const navigate = useNavigate();
  const MenuItem = ({ item, level = 0 }) => {
    const isActive = activeItem === item.id;
    const isExpanded = expandedItems.includes(item.id);

    const handleClick = () => {
      if (item.expandable) {
        toggleExpanded(item.id);
      } else {
        setActiveItem(item.id);
        if (item.href) navigate(item.href);
      }
    };

    return (
      <div>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant={isActive ? "secondary" : "ghost"}
                className={cn(
                  "w-full justify-start h-10 px-3 text-left font-normal",
                  level > 0 && "ml-6 w-auto",
                  isActive && "bg-primary/10 text-primary border-r-2 border-primary",
                  item.highlight && !isActive && "bg-accent/20 hover:bg-accent/30",
                  isCollapsed && "justify-center px-2"
                )}
                onClick={handleClick}
              >
                <item.icon className={cn("h-4 w-4", !isCollapsed && "mr-3")} />
                {!isCollapsed && (
                  <>
                    <span className="flex-1 truncate">{item.label}</span>
                    {item.badge && (
                      <Badge variant="secondary" className="ml-auto h-5 text-xs">
                        {item.badge}
                      </Badge>
                    )}
                    {item.expandable && (
                      isExpanded ? <ChevronDown className="h-4 w-4 ml-2" /> : <ChevronRight className="h-4 w-4 ml-2" />
                    )}
                  </>
                )}
              </Button>
            </TooltipTrigger>
            {isCollapsed && (
              <TooltipContent side="right">
                <p>{item.label}</p>
              </TooltipContent>
            )}
          </Tooltip>
        </TooltipProvider>

        {item.expandable && item.children && (
          <Collapsible open={isExpanded}>
            <CollapsibleContent className="space-y-1">
              {item.children.map((child) => {
                const ChildIcon = child.icon;
                return (
                  <Button
                    key={child.id}
                    variant={activeItem === child.id ? "secondary" : "ghost"}
                    className={cn(
                      "w-full justify-start h-9 ml-6 text-sm font-normal",
                      activeItem === child.id && "bg-primary/10 text-primary",
                      isCollapsed && "hidden"
                    )}
                    onClick={() => {
                      setActiveItem(child.id);
                      if (child.href) navigate(child.href);
                    }}
                  >
                    {ChildIcon && <ChildIcon className="h-4 w-4 mr-2" />}
                    <span className="flex-1 truncate">{child.label}</span>
                    {child.badge && (
                      <Badge variant="destructive" className="ml-auto h-4 text-xs">
                        {child.badge}
                      </Badge>
                    )}
                  </Button>
                );
              })}
            </CollapsibleContent>
          </Collapsible>
        )}
      </div>
    );
  };

  return (
    <div className={cn(
      "fixed left-0 top-0 h-full bg-background border-r border-border flex flex-col transition-all duration-300 z-50",
      isCollapsed ? "w-16" : "w-64"
    )}>
      {/* Header */}
      <div className="p-4 border-b border-border">
        <div className="flex items-center space-x-3">
          <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <Store className="h-5 w-5 text-primary-foreground" />
          </div>
          {!isCollapsed && (
            <div>
              <h2 className="font-semibold text-lg text-foreground">ShopOS</h2>
              <p className="text-xs text-muted-foreground">Business Dashboard</p>
            </div>
          )}
        </div>
      </div>

      {/* User Profile */}
      {!isCollapsed && (
        <div className="p-4 border-b border-border">
          <div className="flex items-center space-x-3">
            <Avatar className="h-10 w-10">
              <AvatarImage src={import.meta.env.BASE_URL + "api/placeholder/40/40"} />
              <AvatarFallback className="bg-primary text-primary-foreground">
                MR
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-foreground truncate">
                Mohammad Rahman
              </p>
              <p className="text-xs text-muted-foreground truncate">
                rahman.shop@gmail.com
              </p>
            </div>
            <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
              <Bell className="h-4 w-4" />
            </Button>
          </div>
        </div>
      )}

      {/* Quick Stats */}
      {!isCollapsed && (
        <div className="p-4 border-b border-border">
          <div className="grid grid-cols-2 gap-3">
            <div className="bg-primary/5 rounded-lg p-3">
              <div className="flex items-center space-x-2">
                <TrendingUp className="h-4 w-4 text-primary" />
                <div>
                  <p className="text-xs text-muted-foreground">Today</p>
                  <p className="text-sm font-semibold text-foreground">৳15,420</p>
                </div>
              </div>
            </div>
            <div className="bg-accent/10 rounded-lg p-3">
              <div className="flex items-center space-x-2">
                <Package className="h-4 w-4 text-accent-foreground" />
                <div>
                  <p className="text-xs text-muted-foreground">Stock</p>
                  <p className="text-sm font-semibold text-foreground">1,247</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Navigation */}
      <nav className="flex-1 p-3 space-y-1 overflow-y-auto">
        {menuItems.map((item) => (
          <MenuItem key={item.id} item={item} />
        ))}
        
        <Separator className="my-4" />
        
        {bottomMenuItems.map((item) => (
          <MenuItem key={item.id} item={item} />
        ))}
      </nav>

      {/* Footer */}
      <div className="p-3 border-t border-border">
        <div className="flex items-center justify-between">
          {!isCollapsed && (
            <div className="flex items-center space-x-2">
              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
              <span className="text-xs text-muted-foreground">Online</span>
            </div>
          )}
          <div className="flex space-x-1">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0"
                    onClick={() => setIsCollapsed(!isCollapsed)}
                  >
                    <ChevronRight className={cn(
                      "h-4 w-4 transition-transform",
                      isCollapsed ? "rotate-0" : "rotate-180"
                    )} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent side="right">
                  <p>{isCollapsed ? 'Expand' : 'Collapse'}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
            
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="ghost" size="sm" className="h-8 w-8 p-0 text-destructive">
                    <LogOut className="h-4 w-4" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent side="right">
                  <p>Sign Out</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardSidebar;