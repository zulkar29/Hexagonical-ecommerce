import { Bell, User, LogOut, Store, Menu, PanelLeftClose, PanelLeftOpen, Plus, Package, Users, ShoppingCart, BarChart3, Zap } from "lucide-react";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAtom, useSetAtom } from "jotai";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel,
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { authStatusAtom, logoutAtom } from "@/store/auth";
import { toast } from "sonner";
import { cn } from "@/lib/utils";

const Header = ({ onToggleSidebar, onToggleCollapse, isCollapsed, isMobile }) => {
  const navigate = useNavigate();
  const [{ user }] = useAtom(authStatusAtom);
  const logout = useSetAtom(logoutAtom);
  
  const notifications = [
    { id: 1, text: "New order #1234 received", time: "2m ago", type: "order", path: "/orders" },
    { id: 2, text: "Stock alert: 'iPhone 13 Pro' low stock", time: "1h ago", type: "inventory", path: "/products?filter=low_stock" },
    { id: 3, text: "New review from John Doe", time: "3h ago", type: "review", path: "/reviews" },
  ];

  const todayOrders = 12;
  
  // Quick action shortcuts - only essential ones with Alt key to avoid browser conflicts
  const shortcuts = [
    {
      id: 'add-product',
      label: 'Add Product',
      icon: Plus,
      path: '/products/create',
      color: 'text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-950',
      shortcut: 'Alt+P'
    },
    {
      id: 'create-order',
      label: 'Create Order',
      icon: ShoppingCart,
      path: '/orders/create',
      color: 'text-emerald-600 hover:bg-emerald-50 dark:hover:bg-emerald-950',
      shortcut: 'Alt+O'
    },
    {
      id: 'products',
      label: 'Products',
      icon: Package,
      path: '/products',
      color: 'text-green-600 hover:bg-green-50 dark:hover:bg-green-950'
    },
    {
      id: 'orders',
      label: 'Orders',
      icon: ShoppingCart,
      path: '/orders',
      color: 'text-purple-600 hover:bg-purple-50 dark:hover:bg-purple-950'
    },
    {
      id: 'customers',
      label: 'Customers',
      icon: Users,
      path: '/customers',
      color: 'text-orange-600 hover:bg-orange-50 dark:hover:bg-orange-950'
    },
    {
      id: 'analytics',
      label: 'Analytics',
      icon: BarChart3,
      path: '/reports',
      color: 'text-indigo-600 hover:bg-indigo-50 dark:hover:bg-indigo-950'
    }
  ];
  
  // Handle keyboard shortcuts - only essential ones to avoid browser conflicts
  useEffect(() => {
    const handleKeyPress = (e) => {
      if (e.altKey) {
        switch (e.key.toLowerCase()) {
          case 'p':
            e.preventDefault();
            navigate('/products/create');
            break;
          case 'o':
            e.preventDefault();
            navigate('/orders/create');
            break;
        }
      }
    };
    
    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [navigate]);
  const userName = user?.name || "Admin User";
  const userRole = user?.role?.display_name || "Administrator";

  const handleLogout = async () => {
    try {
      await logout();
      toast.success('Logged out successfully');
      navigate('/login');
    } catch (error) {
      console.error('Logout error:', error);
      toast.error('Error logging out');
    }
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-16 items-center justify-between px-2 sm:px-4 gap-2 sm:gap-4">
        <div className="flex items-center gap-2 sm:gap-4">
          {isMobile ? (
            <Button variant="ghost" size="icon" onClick={onToggleSidebar} className="lg:hidden cursor-pointer">
              <Menu className="h-5 w-5" />
            </Button>
          ) : (
            <Button 
              variant="ghost" 
              size="icon"
              onClick={onToggleCollapse}
              className="hidden lg:flex cursor-pointer"
              title={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
            >
              {isCollapsed ? (
                <PanelLeftOpen className="h-5 w-5" />
              ) : (
                <PanelLeftClose className="h-5 w-5" />
              )}
            </Button>
          )}

          {/* Page Title with Quick Stats */}
          <div className="flex items-center gap-2 sm:gap-4">
            <div 
              className="flex flex-col cursor-pointer hover:text-primary transition-colors"
              onClick={() => navigate('/')}
            >
              <h1 className="text-base sm:text-lg font-semibold hidden xs:block">Dashboard</h1>
              <span className="text-xs text-muted-foreground hidden lg:block">
                Welcome back, {userName.split(' ')[0]}!
              </span>
            </div>
            <Badge 
              variant="secondary" 
              className="hidden md:flex items-center gap-1 cursor-pointer hover:bg-secondary/80 transition-colors"
              onClick={() => navigate('/orders')}
            >
              <span className="font-normal text-xs">Today:</span>
              <span className="font-medium">{todayOrders} orders</span>
            </Badge>
          </div>
        </div>

        {/* Quick Actions and Shortcuts */}
        <div className="flex items-center gap-1 sm:gap-2">

          {/* Quick Action Shortcuts */}
          <TooltipProvider>
            {/* Desktop Quick Actions */}
            <div className="hidden xl:flex items-center gap-1 mr-2">
              {shortcuts.filter(s => s.shortcut).map((shortcut) => {
                const IconComponent = shortcut.icon;
                return (
                  <Tooltip key={shortcut.id}>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => navigate(shortcut.path)}
                        className={cn(
                          "h-9 w-9 transition-all duration-200 hover:scale-105 cursor-pointer", 
                          shortcut.color
                        )}
                      >
                        <IconComponent className="h-4 w-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent side="bottom" className="font-medium">
                      <div className="flex flex-col items-center gap-1">
                        <span>{shortcut.label}</span>
                        <span className="text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded">
                          {shortcut.shortcut}
                        </span>
                      </div>
                    </TooltipContent>
                  </Tooltip>
                );
              })}
            </div>
            
            {/* Tablet Quick Actions */}
            <div className="hidden lg:flex xl:hidden items-center gap-1 mr-2">
              {shortcuts.filter(s => s.shortcut).map((shortcut) => {
                const IconComponent = shortcut.icon;
                return (
                  <Tooltip key={shortcut.id}>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => navigate(shortcut.path)}
                        className={cn("h-9 w-9 transition-colors cursor-pointer", shortcut.color)}
                      >
                        <IconComponent className="h-4 w-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent side="bottom">
                      <div className="flex flex-col items-center gap-1">
                        <span>{shortcut.label}</span>
                        {shortcut.shortcut && (
                          <span className="text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded">
                            {shortcut.shortcut}
                          </span>
                        )}
                      </div>
                    </TooltipContent>
                  </Tooltip>
                );
              })}
              
              {/* More Actions for Tablet */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button variant="ghost" size="icon" className="h-9 w-9 cursor-pointer">
                        <Zap className="h-4 w-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent side="bottom">
                      More Actions
                    </TooltipContent>
                  </Tooltip>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56">
                  <DropdownMenuLabel className="flex items-center gap-2">
                    <Zap className="h-4 w-4" />
                    Quick Actions
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  {shortcuts.map((shortcut) => {
                    const IconComponent = shortcut.icon;
                    return (
                      <DropdownMenuItem
                        key={shortcut.id}
                        onClick={() => navigate(shortcut.path)}
                        className="cursor-pointer group"
                      >
                        <IconComponent className="mr-3 h-4 w-4 group-hover:scale-110 transition-transform" />
                        <div className="flex-1">
                          <div className="font-medium">{shortcut.label}</div>
                          {shortcut.shortcut && (
                            <div className="text-xs text-muted-foreground">{shortcut.shortcut}</div>
                          )}
                        </div>
                      </DropdownMenuItem>
                    );
                  })}
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
            
            {/* Mobile Quick Actions */}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="lg:hidden relative cursor-pointer">
                  <Zap className="h-5 w-5" />
                  <span className="absolute -top-1 -right-1 h-2 w-2 bg-primary rounded-full" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuLabel className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Zap className="h-4 w-4" />
                    Quick Actions
                  </div>
                  <Badge variant="secondary" className="text-xs">⌨️</Badge>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                {shortcuts.map((shortcut) => {
                  const IconComponent = shortcut.icon;
                  return (
                    <DropdownMenuItem
                      key={shortcut.id}
                      onClick={() => navigate(shortcut.path)}
                      className="cursor-pointer group py-3"
                    >
                      <IconComponent className="mr-3 h-4 w-4 group-hover:scale-110 transition-transform" />
                      <div className="flex-1">
                        <div className="font-medium">{shortcut.label}</div>
                        {shortcut.shortcut && (
                          <div className="text-xs text-muted-foreground">{shortcut.shortcut}</div>
                        )}
                      </div>
                    </DropdownMenuItem>
                  );
                })}
              </DropdownMenuContent>
            </DropdownMenu>
          </TooltipProvider>
          
          {/* Notifications */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="relative cursor-pointer">
                <Bell className="h-5 w-5" />
                {notifications.length > 0 && (
                  <span className="absolute top-1 right-1 h-2 w-2 rounded-full bg-destructive" />
                )}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-[280px] sm:w-80">
              <DropdownMenuLabel className="flex items-center justify-between">
                <span>Notifications</span>
                <div className="flex items-center gap-2">
                  <Badge variant="secondary" className="ml-2">{notifications.length}</Badge>
                  <Button 
                    variant="ghost" 
                    size="sm" 
                    className="h-6 px-2 text-xs cursor-pointer"
                    onClick={() => navigate('/notifications')}
                  >
                    View All
                  </Button>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              {notifications.length === 0 ? (
                <div className="p-4 text-center text-sm text-muted-foreground">
                  No new notifications
                </div>
              ) : (
                notifications.map((notification) => (
                  <DropdownMenuItem 
                    key={notification.id} 
                    className="p-3 cursor-pointer hover:bg-accent/50"
                    onClick={() => navigate(notification.path)}
                  >
                    <div className="flex flex-col gap-1">
                      <span className="font-medium text-sm">{notification.text}</span>
                      <div className="flex items-center justify-between">
                        <span className="text-xs text-muted-foreground">{notification.time}</span>
                        <Badge variant="secondary" className="text-xs">
                          {notification.type}
                        </Badge>
                      </div>
                    </div>
                  </DropdownMenuItem>
                ))
              )}
            </DropdownMenuContent>
          </DropdownMenu>

          {/* User Menu */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="relative h-9 w-9 rounded-full cursor-pointer">
                <Avatar className="h-9 w-9">
                  <AvatarImage src="/avatar-placeholder.jpg" alt={userName} />
                  <AvatarFallback>{userName[0]}</AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel className="font-normal">
                <div className="flex flex-col space-y-1">
                  <p className="text-sm font-medium">{userName}</p>
                  <p className="text-xs text-muted-foreground">{userRole}</p>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => navigate('/profile')} className="cursor-pointer">
                <User className="mr-2 h-4 w-4" />
                Profile
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => navigate('/settings')} className="cursor-pointer">
                <Store className="mr-2 h-4 w-4" />
                Store Settings
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => navigate('/help')} className="cursor-pointer">
                <Bell className="mr-2 h-4 w-4" />
                Help & Support
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="text-destructive cursor-pointer" onClick={handleLogout}>
                <LogOut className="mr-2 h-4 w-4" />
                Log out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

    </header>
  );
};

export default Header;