

import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import DashboardHeader from './components/layout/Header';
import DashboardSidebar from './components/layout/Sidebar';
import TenantsHome from './pages/Tenants';
import SubscriptionsHome from './pages/Subscriptions';
import BillingHome from './pages/Billing';
import UsersHome from './pages/Users';
import AnalyticsHome from './pages/Analytics';
import SupportHome from './pages/Support';
import SettingsHome from './pages/Settings';
import PaymentsHome from './pages/Payments';
import Dashboard from './pages/Dashboard';
import TenantDetail from './pages/Tenants/Detail';
import UserCreate from './pages/Users/create';
import UserEdit from './pages/Users/edit';
import SupportCreate from './pages/Support/Create';
import SupportEdit from './pages/Support/Edit';
import BillingCreate from './pages/Billing/Create';
import BillingEdit from './pages/Billing/Edit';
import PaymentDetail from './pages/Payments/Detail';
import SubscriptionDetail from './pages/Subscriptions/Detail';
import CreateSubscription from './pages/Subscriptions/Create';
import ModifySubscription from './pages/Subscriptions/Modify';
import SupportDetail from './pages/Support/Detail';
import TenantOnboard from './pages/Tenants/Onboard';
import TenantManage from './pages/Tenants/Manage';
import SystemHealth from './pages/Settings/SystemHealth';
import Database from './pages/Settings/Database';
import APIManagement from './pages/Settings/API';
import Backups from './pages/Settings/Backups';
import AuditLogs from './pages/Settings/Logs';
import KnowledgeBase from './pages/Support/KnowledgeBase';
import Communication from './pages/Support/Communication';
import Announcements from './pages/Support/Announcements';
import RevenueAnalytics from './pages/Analytics/Revenue';
import TenantAnalytics from './pages/Analytics/Tenants';
import PerformanceAnalytics from './pages/Analytics/Performance';
import CustomReports from './pages/Analytics/Custom';
import RoleManagement from './pages/Users/Roles';
import EditRole from './pages/Users/EditRole';
import CreateRole from './pages/Users/CreateRole';

function App() {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  return (
    <Router>
      <div className="min-h-screen bg-background flex">
        <DashboardSidebar />
        <div className="flex-1 flex flex-col min-w-0">
          <DashboardHeader sidebarCollapsed={sidebarCollapsed} onToggleSidebar={() => setSidebarCollapsed(!sidebarCollapsed)} />
          <main className="flex-1 pt-20 overflow-auto">
            <Routes>
              {/* Dashboard Routes */}
              <Route index element={<Dashboard />} />
              <Route path="dashboard" element={<Dashboard />} />
              
              {/* Tenant Management Routes */}
              <Route path="tenants" element={<TenantsHome />} />
              <Route path="tenants/onboard" element={<TenantOnboard />} />
              <Route path="tenants/status/:status" element={<TenantsHome />} />
              <Route path="tenants/:id" element={<TenantDetail />} />
              <Route path="tenants/:id/manage" element={<TenantManage />} />
              
              {/* Subscription Management Routes */}
              <Route path="subscriptions" element={<SubscriptionsHome />} />
              <Route path="subscriptions/create" element={<CreateSubscription />} />
              <Route path="subscriptions/plan/:plan" element={<SubscriptionsHome />} />
              <Route path="subscriptions/:id" element={<SubscriptionDetail />} />
              <Route path="subscriptions/:id/modify" element={<ModifySubscription />} />
              
              {/* Billing & Payments Routes */}
              <Route path="billing" element={<BillingHome />} />
              <Route path="billing/create" element={<BillingCreate />} />
              <Route path="billing/:id/edit" element={<BillingEdit />} />
              <Route path="payments" element={<PaymentsHome />} />
              <Route path="payments/:id" element={<PaymentDetail />} />
              
              {/* User & Permission Management Routes */}
              <Route path="users" element={<UsersHome />} />
              <Route path="users/create" element={<UserCreate />} />
              <Route path="users/:id/edit" element={<UserEdit />} />
              <Route path="users/permissions" element={<RoleManagement />} />
              <Route path="users/permissions/create" element={<CreateRole />} />
              <Route path="users/permissions/:id/edit" element={<EditRole />} />
              
              {/* Analytics Routes */}
              <Route path="analytics" element={<AnalyticsHome />} />
              <Route path="analytics/revenue" element={<RevenueAnalytics />} />
              <Route path="analytics/tenants" element={<TenantAnalytics />} />
              <Route path="analytics/performance" element={<PerformanceAnalytics />} />
              <Route path="analytics/custom" element={<CustomReports />} />
              
              {/* Support Routes */}
              <Route path="support" element={<SupportHome />} />
              <Route path="support/create" element={<SupportCreate />} />
              <Route path="support/:id/edit" element={<SupportEdit />} />
              <Route path="support/knowledge-base" element={<KnowledgeBase />} />
              <Route path="support/communication" element={<Communication />} />
              <Route path="support/announcements" element={<Announcements />} />
              <Route path="support/:id" element={<SupportDetail />} />
              
              {/* Settings Routes */}
              <Route path="settings" element={<SettingsHome />} />
              <Route path="settings/system-health" element={<SystemHealth />} />
              <Route path="settings/database" element={<Database />} />
              <Route path="settings/api" element={<APIManagement />} />
              <Route path="settings/backups" element={<Backups />} />
              <Route path="settings/logs" element={<AuditLogs />} />
            </Routes>
          </main>
        </div>
      </div>
    </Router>
  )
}

export default App;
