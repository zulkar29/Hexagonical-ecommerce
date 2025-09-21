import { useState, useEffect } from "react";
import { Outlet } from "react-router-dom";
import Header from "@/components/dashboard/Header";
import Sidebar from "@/components/dashboard/Sidebar";
import { cn } from "@/lib/utils";
import ErrorBoundary from "@/components/ErrorBoundary";

const DashboardLayout = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkSize = () => {
      const width = window.innerWidth;
      setIsMobile(width < 768);
      setSidebarOpen(width >= 768);
      setIsCollapsed(width < 1024);
    };

    checkSize();
    window.addEventListener("resize", checkSize);
    return () => window.removeEventListener("resize", checkSize);
  }, []);

  return (
    <div className="min-h-screen bg-background">
      <div className="flex h-screen overflow-hidden">
        {/* Sidebar */}
        <div
          className={cn(
            "fixed inset-y-0 left-0 z-50 md:relative transform transition-all duration-200 ease-in-out",
            sidebarOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0",
            isCollapsed ? "w-14 sm:w-16 md:w-20" : "w-52 sm:w-56 md:w-64",
            isMobile && "backdrop-blur-sm"
          )}
        >
          <Sidebar isCollapsed={isCollapsed} />
        </div>

        {/* Main Content */}
        <div className="flex-1 flex flex-col overflow-hidden">
          <Header 
            onToggleSidebar={() => setSidebarOpen(!sidebarOpen)}
            onToggleCollapse={() => setIsCollapsed(!isCollapsed)}
            isCollapsed={isCollapsed}
            isMobile={isMobile}
          />
          <main className="flex-1 overflow-y-auto bg-muted/10">
            <ErrorBoundary>
              <Outlet />
            </ErrorBoundary>
          </main>
        </div>

        {/* Mobile Overlay */}
        {isMobile && sidebarOpen && (
          <div
            className="fixed inset-0 bg-black/50 z-40"
            onClick={() => setSidebarOpen(false)}
          />
        )}
      </div>
    </div>
  );
};

export default DashboardLayout;