

import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import DashboardHeader from './components/layout/Header';
import DashboardSidebar from './components/layout/Sidebar';
import TenantsHome from './pages/Tenants';
import SubscriptionsHome from './pages/Subscriptions';
import BillingHome from './pages/Billing';
import UsersHome from './pages/UsersPermission/User';
import AnalyticsHome from './pages/Analytics';
import SupportHome from './pages/Support';
import SettingsHome from './pages/Settings';
import PaymentsHome from './pages/Payments';
import Dashboard from './pages/Dashboard';
import TenantDetail from './pages/Tenants/Detail';
import UserDetail from './pages/UsersPermission/User/create';
import PaymentDetail from './pages/Payments/Detail';
import SubscriptionDetail from './pages/Subscriptions/Detail';
import SupportDetail from './pages/Support/Detail';
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
import RoleManagement from './pages/UsersPermission/Role';
import EditRole from './pages/UsersPermission/Role/edit';
import CreateRole from './pages/UsersPermission/Role/create';

function App() {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  return (
    <Router>
      <div className="min-h-screen bg-background flex">
        <DashboardSidebar />
        <div className="flex-1 flex flex-col">
          <DashboardHeader sidebarCollapsed={sidebarCollapsed} onToggleSidebar={() => setSidebarCollapsed(!sidebarCollapsed)} />
          <main className="flex-1 p-6 pt-8 md:pt-10 lg:pt-12 xl:pt-14" style={{ marginTop: '4rem' }}>
            <Routes>
              <Route index element={<Dashboard />} />
              <Route path="dashboard" element={<Dashboard />} />
              <Route path="tenants" element={<TenantsHome />} />
              <Route path="tenants/:id" element={<TenantDetail />} />
              <Route path="tenants/:status" element={<TenantsHome />} />
              <Route path="subscriptions" element={<SubscriptionsHome />} />
              <Route path="subscriptions/:plan" element={<SubscriptionsHome />} />
              <Route path="subscriptions/:id" element={<SubscriptionDetail />} />
              <Route path="billing" element={<BillingHome />} />
              <Route path="users" element={<UsersHome />} />
              <Route path="users/:id" element={<UserDetail />} />
              <Route path="users/permissions" element={<RoleManagement />} />
              <Route path="users/permissions/role/create" element={<CreateRole />} />
              <Route path="users/permissions/:id" element={<EditRole />} />
              <Route path="users/permissions/create" element={<CreateRole />} />
              <Route path="users/permissions/edit/:id" element={<EditRole />} />
              <Route path="analytics" element={<AnalyticsHome />} />
              <Route path="support" element={<SupportHome />} />
              <Route path="support/:id" element={<SupportDetail />} />
              <Route path="support/knowledge-base" element={<KnowledgeBase />} />
              <Route path="support/communication" element={<Communication />} />
              <Route path="support/announcements" element={<Announcements />} />
              <Route path="settings" element={<SettingsHome />} />
              <Route path="settings/system-health" element={<SystemHealth />} />
              <Route path="settings/database" element={<Database />} />
              <Route path="settings/api" element={<APIManagement />} />
              <Route path="settings/backups" element={<Backups />} />
              <Route path="settings/logs" element={<AuditLogs />} />
              <Route path="payments" element={<PaymentsHome />} />
              <Route path="payments/:id" element={<PaymentDetail />} />
              <Route path="analytics" element={<AnalyticsHome />} />
              <Route path="analytics/revenue" element={<RevenueAnalytics />} />
              <Route path="analytics/tenants" element={<TenantAnalytics />} />
              <Route path="analytics/performance" element={<PerformanceAnalytics />} />
              <Route path="analytics/custom" element={<CustomReports />} />
            </Routes>
          </main>
        </div>
      </div>
    </Router>
  )
}

export default App;
