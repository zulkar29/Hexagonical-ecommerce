import React, { useState } from 'react';
import { useLocation, Link } from 'react-router-dom';
import {
  LayoutDashboard,
  Store,
  Users,
  CreditCard,
  LifeBuoy,
  BarChart3,
  Settings,
  UserPlus,
  Receipt,
  MessageSquare,
  Shield,
  Database,
  Globe,
  Bell,
  FileText,
  Package,
  TrendingUp,
  Activity,
  AlertTriangle,
  CheckCircle,
  Clock,
  PieChart,
  DollarSign,
  Mail,
  Phone,
  Calendar,
  Archive,
  Download,
  Upload,
  Zap,
  Target,
  Award,
  BookOpen,
  HelpCircle,
  ChevronDown,
  ChevronRight,
  Minimize2,
  Maximize2,
  LogOut,
  Search,
  Filter,
  Plus
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
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

const AdminSidebar = ({ 
  isCollapsed = false, 
  onToggleCollapse,
  isMobile = false,
  onClose,
  className = "",
  user = { name: 'Admin User', role: 'Super Admin' }
}) => {
  const location = useLocation();
  const [expandedMenus, setExpandedMenus] = useState(['dashboard', 'management']);

  const toggleMenu = (menuId) => {
    setExpandedMenus(prev => 
      prev.includes(menuId) 
        ? prev.filter(id => id !== menuId)
        : [...prev, menuId]
    );
  };

  const formatCurrency = (amount) => `৳${amount.toLocaleString()}`;

  // Quick stats for sidebar
  const quickStats = [
    { label: 'Active Tenants', value: '156', icon: Store, color: 'text-blue-500' },
    { label: 'Monthly Revenue', value: formatCurrency(245000), icon: DollarSign, color: 'text-green-500' },
    { label: 'Open Tickets', value: '12', icon: AlertTriangle, color: 'text-red-500' },
    { label: 'New Signups', value: '24', icon: UserPlus, color: 'text-purple-500' }
  ];

  // Navigation menu structure
  const navigationMenu = [
    {
      id: 'dashboard',
      title: 'Dashboard',
      icon: LayoutDashboard,
      href: '/dashboard',
    },
    {
      title: 'Tenants',
       href: '/tenants',
       icon: Store
    },
    {
     title: 'Users and Permissions',
      icon: Users,
      subItems: [
        { id: 'users', title: 'Users', href: '/users', icon: Users },
        { id: 'permissions', title: 'Permissions', href: '/users/permissions', icon: Shield }
      ]
    },
    {
      id: 'billing',
      title: 'Billing & Payments',
      icon: CreditCard,
      subItems: [
        { title: 'Payments', href: '/payments', icon: Receipt },
        { title: 'Invoices', href: '/billing/invoices', icon: FileText },
        { title: 'Transactions', href: '/billing/transactions', icon: DollarSign },
        { title: 'Refunds', href: '/billing/refunds', icon: Archive }
      ]
    },
    {
      id: 'support',
      title: 'Support & Help',
      icon: LifeBuoy,
      subItems: [
        { title: 'Support Tickets', href: '/support', icon: MessageSquare },
        { title: 'Knowledge Base', href: '/support/knowledge-base', icon: BookOpen },
        { title: 'Communication', href: '/support/communication', icon: Mail },
        { title: 'Announcements', href: '/support/announcements', icon: Bell }
      ]
    },
    {
      id: 'reports',
      title: 'Reports & Analytics',
      icon: BarChart3,
      subItems: [
        { title: 'Revenue Analytics', href: '/analytics/revenue', icon: TrendingUp },
        { title: 'Tenant Analytics', href: '/analytics/tenants', icon: PieChart },
        { title: 'Performance', href: '/analytics/performance', icon: Activity },
        { title: 'Custom Reports', href: '/analytics/custom', icon: Target }
      ]
    },
    {
      id: 'system',
      title: 'System & Settings',
      icon: Settings,
      subItems: [
        { title: 'Platform Settings', href: '/settings', icon: Settings },
        { title: 'System Health', href: '/settings/system-health', icon: Shield },
        { title: 'Database', href: '/settings/database', icon: Database },
        { title: 'API Management', href: '/settings/api', icon: Globe },
        { title: 'Backups', href: '/settings/backups', icon: Download },
        { title: 'Audit Logs', href: '/settings/logs', icon: FileText }
      ]
    }
  ];

  const isActive = (href) => {
    if (href === '/dashboard' || href === '/') {
      return location.pathname === '/' || location.pathname === '/dashboard';
    }
    return location.pathname === href || location.pathname.startsWith(href + '/');
  };

  const hasActiveChild = (items) => {
    return items?.some(item => 
      isActive(item.href) || (item.subItems && hasActiveChild(item.subItems))
    );
  };

  const renderMenuItem = (item, level = 0) => {
    const Icon = item.icon;
    const isItemActive = isActive(item.href);
    const hasChildren = item.subItems && item.subItems.length > 0;
    const isExpanded = expandedMenus.includes(item.id || item.title);
    const hasActiveChildren = hasChildren && hasActiveChild(item.subItems);

    const menuContent = (
      <div className={cn(
        "flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-all hover:bg-accent",
        level > 0 && "ml-4 pl-6",
        isItemActive && "bg-accent text-accent-foreground font-medium",
        hasActiveChildren && !isItemActive && "text-accent-foreground"
      )}>
        {Icon && (
          <Icon className={cn(
            "h-4 w-4 shrink-0",
            isItemActive && "text-primary",
            hasActiveChildren && !isItemActive && "text-primary/70"
          )} />
        )}
        {!isCollapsed && (
          <>
            <span className="truncate">{item.title}</span>
            {item.badge && (
              <Badge 
                variant={item.badge.variant || 'secondary'} 
                className="ml-auto text-xs"
              >
                {item.badge.text}
              </Badge>
            )}
            {hasChildren && (
              <ChevronRight 
                className={cn(
                  "h-3 w-3 ml-auto shrink-0 transition-transform",
                  isExpanded && "rotate-90"
                )}
              />
            )}
          </>
        )}
      </div>
    );

    if (hasChildren && !item.href) {
      return (
        <Collapsible 
          key={item.id || item.title}
          open={isExpanded}
          onOpenChange={() => toggleMenu(item.id || item.title)}
        >
          <CollapsibleTrigger asChild>
            {isCollapsed ? (
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button variant="ghost" className="w-full justify-start p-0">
                      {menuContent}
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent side="right">
                    {item.title}
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            ) : (
              <Button variant="ghost" className="w-full justify-start p-0">
                {menuContent}
              </Button>
            )}
          </CollapsibleTrigger>
          {!isCollapsed && (
            <CollapsibleContent className="space-y-1">
              {item.subItems.map(subItem => renderMenuItem(subItem, level + 1))}
            </CollapsibleContent>
          )}
        </Collapsible>
      );
    }

    const linkContent = item.href ? (
      <Link 
        to={item.href} 
        className="block w-full"
        onClick={isMobile ? onClose : undefined}
      >
        {menuContent}
      </Link>
    ) : menuContent;

    return (
      <div key={item.href || item.title}>
        {isCollapsed ? (
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                {linkContent}
              </TooltipTrigger>
              <TooltipContent side="right">
                {item.title}
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        ) : (
          linkContent
        )}
      </div>
    );
  };

  const renderSection = (section) => {
    const isExpanded = expandedMenus.includes(section.id);
    const hasActiveChildren = hasActiveChild(section.subItems);

    return (
      <div key={section.id} className="space-y-1">
        {!isCollapsed && (
          <Collapsible 
            open={isExpanded}
            onOpenChange={() => toggleMenu(section.id)}
          >
            <CollapsibleTrigger asChild>
              <Button 
                variant="ghost" 
                className={cn(
                  "w-full justify-start px-3 py-2 text-sm font-medium text-muted-foreground hover:text-foreground",
                  hasActiveChildren && "text-foreground"
                )}
              >
                <section.icon className="h-4 w-4 mr-3" />
                {section.title}
                <ChevronRight 
                  className={cn(
                    "h-3 w-3 ml-auto transition-transform",
                    isExpanded && "rotate-90"
                  )}
                />
              </Button>
            </CollapsibleTrigger>
            <CollapsibleContent className="space-y-1 pl-2">
              {section.subItems.map(item => renderMenuItem(item))}
            </CollapsibleContent>
          </Collapsible>
        )}
        {isCollapsed && section.subItems.map(item => renderMenuItem(item))}
      </div>
    );
  };

  return (
    <TooltipProvider>
      <div className={cn(
        "flex h-full flex-col border-r bg-background pt-16",
        isCollapsed ? "w-16" : "w-64",
        isMobile && "fixed inset-y-0 left-0 z-50 shadow-lg",
        className
      )}>
        {/* Sidebar Header */}
        <div className={cn(
          "flex items-center gap-2 border-b px-3 py-4",
          isCollapsed && "justify-center"
        )}>
          {!isCollapsed && (
            <>
              <div className="flex items-center gap-2 flex-1">
                <Avatar className="h-8 w-8">
                  <AvatarFallback className="bg-primary text-primary-foreground text-xs">
                    {user.name.split(' ').map(n => n[0]).join('').toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium truncate">{user.name}</p>
                  <p className="text-xs text-muted-foreground truncate">{user.role}</p>
                </div>
              </div>
            </>
          )}
          
          {/* Collapse Toggle */}
          <Button
            variant="ghost"
            size="sm"
            onClick={onToggleCollapse}
            className={cn(
              "h-8 w-8 p-0",
              isMobile && "hidden"
            )}
          >
            {isCollapsed ? (
              <Maximize2 className="h-4 w-4" />
            ) : (
              <Minimize2 className="h-4 w-4" />
            )}
          </Button>
        </div>

        {/* Quick Stats */}
        {!isCollapsed && (
          <div className="border-b p-3">
            <div className="grid grid-cols-2 gap-2">
              {quickStats.map((stat, index) => (
                <div key={index} className="bg-muted/50 rounded-lg p-2">
                  <div className="flex items-center gap-1">
                    <stat.icon className={cn("h-3 w-3", stat.color)} />
                    <span className="text-xs font-medium">{stat.value}</span>
                  </div>
                  <p className="text-xs text-muted-foreground truncate mt-1">
                    {stat.label}
                  </p>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Navigation Menu */}
        <nav className="flex-1 overflow-y-auto p-3 space-y-2">
          {navigationMenu.map((item) => (
            item.isSection ? renderSection(item) : renderMenuItem(item)
          ))}
        </nav>

        {/* Quick Actions */}
        {!isCollapsed && (
          <div className="border-t p-3 space-y-2">
            <p className="text-xs font-medium text-muted-foreground mb-2">Quick Actions</p>
            <div className="space-y-1">
              <Button variant="outline" size="sm" className="w-full justify-start">
                <Plus className="h-4 w-4 mr-2" />
                Add Tenant
              </Button>
              <Button variant="outline" size="sm" className="w-full justify-start">
                <Search className="h-4 w-4 mr-2" />
                Global Search
              </Button>
            </div>
          </div>
        )}

        {/* Footer */}
        <div className={cn(
          "border-t p-3",
          isCollapsed && "flex justify-center"
        )}>
          {isCollapsed ? (
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                  <HelpCircle className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="right">
                Help & Support
              </TooltipContent>
            </Tooltip>
          ) : (
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="h-2 w-2 bg-green-500 rounded-full animate-pulse" />
                <span className="text-xs text-muted-foreground">System Operational</span>
              </div>
              <Button variant="ghost" size="sm" className="h-6 w-6 p-0">
                <HelpCircle className="h-3 w-3" />
              </Button>
            </div>
          )}
        </div>

        {/* Mobile overlay */}
        {isMobile && (
          <div 
            className="fixed inset-0 bg-background/80 backdrop-blur-sm z-40 md:hidden" 
            onClick={onClose}
          />
        )}
      </div>
    </TooltipProvider>
  );
};

export default AdminSidebar;