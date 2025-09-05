import React, { useState } from 'react';
import DashboardHeader from './components/layout/Header';
import DashboardSidebar from './components/layout/Sidebar';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Inventory from './pages/Inventory';
import POS from './pages/POS';
import Customers from './pages/Customers';
import Sales from './pages/Sales';
import Reports from './pages/Reports';
import CustomerDetails from './pages/CustomerDetails';
import SettingsPage from './pages/Settings';
import HelpPage from './pages/Help';
import SalesReportPage from './pages/SalesReport';
import Products from './pages/Products';
import Categories from './pages/Categories';
import Suppliers from './pages/Suppliers';
import StockAlerts from './pages/StockAlerts';
import Orders from './pages/Orders';
import InventoryReport from './pages/InventoryReport';
import CustomerInsights from './pages/CustomerInsights';
import ProductDetails from './pages/ProductDetails';
import ProductCreate from './pages/ProductCreate';
import SupplierDetails from './pages/SupplierDetails';
import OrderDetails from './pages/OrderDetails';
import AccountsPayable from './pages/AccountsPayable';
import AccountsReceivable from './pages/AccountsReceivable';
import PettyCash from './pages/PettyCash';


function App() {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  return (
    <Router basename="/mini-erp">
      <div className="min-h-screen bg-background flex">
        <DashboardSidebar />
        <div className={sidebarCollapsed ? "ml-16 flex-1 flex flex-col" : "ml-64 flex-1 flex flex-col"}>
          <DashboardHeader sidebarCollapsed={sidebarCollapsed} onToggleSidebar={() => setSidebarCollapsed(!sidebarCollapsed)} />
          <main className="flex-1 p-6 pt-8 md:pt-10 lg:pt-12 xl:pt-14">
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/customer-details/:id" element={<CustomerDetails />} />
              <Route path="/settings" element={<SettingsPage />} />
              <Route path="/help" element={<HelpPage />} />
              <Route path="/inventory" element={<Inventory />} />
              <Route path="/pos" element={<POS />} />
              <Route path="/customers" element={<Customers />} />
              <Route path="/sales" element={<Sales />} />
              <Route path="/reports" element={<Reports />} />
              <Route path="/inventory/products" element={<Products />} />
              <Route path="/inventory/products/create" element={<ProductCreate />} />
              <Route path="/inventory/products/:id" element={<ProductDetails />} />
              <Route path="/inventory/categories" element={<Categories />} />
              <Route path="/inventory/suppliers" element={<Suppliers />} />
              <Route path="/inventory/suppliers/:id" element={<SupplierDetails />} />
              <Route path="/inventory/alerts" element={<StockAlerts />} />
              <Route path="/orders" element={<Orders />} />
              <Route path="/orders/:id" element={<OrderDetails />} />
              <Route path="/reports/sales" element={<SalesReportPage />} />
              <Route path="/reports/inventory" element={<InventoryReport />} />
              <Route path="/reports/customers" element={<CustomerInsights />} />
              <Route path="/finance/payable" element={<AccountsPayable />} />
              <Route path="/finance/receivable" element={<AccountsReceivable />} />
              <Route path="/finance/petty-cash" element={<PettyCash />} />
            </Routes>
          </main>
        </div>
      </div>
    </Router>
  )
}

export default App
