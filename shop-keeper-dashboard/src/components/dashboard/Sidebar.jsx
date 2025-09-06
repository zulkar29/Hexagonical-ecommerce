import { Link, useLocation, useNavigate } from "react-router-dom";
import { 
  LayoutDashboard, 
  ShoppingBag, 
  Users, 
  ClipboardList, 
  LogOut,
  Tags,
  MessageSquare,
  Gift,
  Image,
  Settings2,
  HelpCircle,
  ChevronDown,
  ChevronRight,
  UserCheck,
  Wallet,
  Activity,
  FileText,
  Palette
} from "lucide-react";
import { Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { useState } from "react";
import { useSetAtom } from 'jotai';
import { authActionsAtom } from '@/store/auth';

const Sidebar = ({ isCollapsed }) => {
  const location = useLocation();
  const navigate = useNavigate();
  const setAuthAction = useSetAtom(authActionsAtom);
  const handleLogout = () => {
    setAuthAction({ type: 'LOGOUT' });
    navigate('/login');
  };
  const [expanded, setExpanded] = useState({
    store: true,
    marketing: false,
    analytics: false,
    settings: false,
    });

  const toggleSection = (section) => {
    if (!isCollapsed) {
      setExpanded((prev) => ({ ...prev, [section]: !prev[section] }));
    }
  };

  const storeItems = [
    { name: "Dashboard", path: "/", icon: <LayoutDashboard className="w-5 h-5" /> },
    { name: "Store Designer", path: "/builder", icon: <Palette className="w-5 h-5" /> },
    { name: "Products", path: "/products", icon: <ShoppingBag className="w-5 h-5" /> },
    { name: "Categories", path: "/category", icon: <Tags className="w-5 h-5" /> },
    { name: "Orders", path: "/orders", icon: <ClipboardList className="w-5 h-5" />, badge: 3 },
    { name: "Customers", path: "/customers", icon: <Users className="w-5 h-5" /> },
    { name: "Reviews", path: "/reviews", icon: <MessageSquare className="w-5 h-5" /> },
    { name: "Petty Cash", path: "/petty-cash", icon: <Wallet className="w-5 h-5" /> },
    { name: "Affiliates", path: "/affiliates", icon: <UserCheck className="w-5 h-5" /> },
    { name: "Support", path: "/support", icon: <HelpCircle className="w-5 h-5" /> },
    { name: "Settings", path: "/settings", icon: <Settings2 className="w-5 h-5" /> },
    { name: "User Roles", path: "/user-roles", icon: <Users className="w-5 h-5" /> },
    { name: "Logs", path: "/logs", icon: <Activity className="w-5 h-5" /> },
    { name: "Reports", path: "/reports", icon: <FileText className="w-5 h-5" /> },
  ];

  const marketingItems = [
    { name: "Discounts", path: "/discounts", icon: <Gift className="w-5 h-5" /> },
    { name: "Banners", path: "/banners", icon: <Image className="w-5 h-5" /> },
  ]; 


  const NavItem = ({ item }) => {
    const isActive = location.pathname === item.path;
    
    const content = (
      <Link
        to={item.path}
        className={cn(
          "flex items-center justify-between px-3 py-2 rounded-lg text-sm font-medium transition-colors cursor-pointer",
          "hover:bg-accent hover:text-accent-foreground",
          isActive ? "bg-accent text-accent-foreground" : "text-muted-foreground"
        )}
      >
        <div className="flex items-center gap-3">
          {item.icon}
          {!isCollapsed && <span>{item.name}</span>}
        </div>
        {!isCollapsed && item.badge && (
          <Badge variant="secondary" className="ml-auto">
            {item.badge}
          </Badge>
        )}
      </Link>
    );

    return isCollapsed ? (
      <TooltipProvider>
        <Tooltip delayDuration={0}>
          <TooltipTrigger asChild>
            {content}
          </TooltipTrigger>
          <TooltipContent side="right" className="font-medium">
            {item.name}
            {item.badge && ` (${item.badge})`}
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    ) : content;
  };

  const renderSection = (title, sectionKey, items) => (
    <div className="mb-2">
      {!isCollapsed && (
        <Button
          variant="ghost"
          size="sm"
          className="w-full cursor-pointer flex items-center justify-between px-3 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider hover:text-primary"
          onClick={() => toggleSection(sectionKey)}
          type="button"
        >
          {title}
          {expanded[sectionKey] ? (
            <ChevronDown className="w-4 h-4" />
          ) : (
            <ChevronRight className="w-4 h-4" />
          )}
        </Button>
      )}
      <div className={cn(
        "space-y-1",
        !isCollapsed && "mt-1 px-2",
        !isCollapsed && !expanded[sectionKey] && "hidden"
      )}>
        {items.map((item) => (
          <NavItem key={item.path} item={item} />
        ))}
      </div>
    </div>
  );

  return (
    <aside className={cn(
      "h-screen bg-background border-r flex flex-col",
      isCollapsed ? "w-14 sm:w-16 md:w-20" : "w-52 sm:w-56 md:w-64",
      "transition-all duration-200 ease-in-out"
    )}>
      <div className={cn(
        "flex flex-col items-center py-4 sm:py-6 border-b",
        isCollapsed && "py-3 sm:py-4"
      )}>
        <Avatar className={cn(
          "mb-2 transition-all duration-200",
          isCollapsed ? "w-6 h-6 sm:w-8 sm:h-8 md:w-10 md:h-10" : "w-10 h-10 sm:w-12 sm:h-12 md:w-14 md:h-14"
        )}>
          <img src="/vite.svg" alt="Store Logo" />
        </Avatar>
        {!isCollapsed && (
          <>
            <h1 className="text-lg sm:text-xl md:text-2xl font-extrabold text-primary tracking-tight">ShopVendor</h1>
            <span className="text-xs text-muted-foreground">Admin Panel</span>
          </>
        )}
      </div>
      
      <nav className={cn(
        "flex-1 overflow-y-auto cursor-pointer",
        isCollapsed ? "px-1 sm:px-2" : "px-2 sm:px-3",
        "transition-all duration-200 py-2"
      )}>
        {renderSection("Store Management", "store", storeItems)}
        {renderSection("Marketing", "marketing", marketingItems)}
      </nav>
      
      <div className={cn(
        "p-2 sm:p-3 md:p-4 border-t",
        isCollapsed && "p-1 sm:p-2"
      )}>
        {isCollapsed ? (
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="w-full h-8 sm:h-10 cursor-pointer" onClick={handleLogout}>
                  <LogOut className="w-4 h-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="right">Logout</TooltipContent>
            </Tooltip>
          </TooltipProvider>
        ) : (
          <>
            <Button variant="ghost" className="w-full cursor-pointer justify-start text-muted-foreground hover:text-destructive text-xs sm:text-sm" size="sm" onClick={handleLogout}>
              <LogOut className="w-4 h-4 mr-2" />
              Logout
            </Button>
            <div className="mt-3 sm:mt-4 text-center text-xs text-muted-foreground">
              <p>ShopVendor v2.0.0</p>
              <p className="mt-1">© 2025 All rights reserved</p>
            </div>
          </>
        )}
      </div>
    </aside>
  );
};

export default Sidebar;