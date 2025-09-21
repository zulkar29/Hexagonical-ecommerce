import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { Toaster } from 'sonner';
import { AuthProvider } from './contexts/AuthContext';
import { ComponentProvider } from './pages/shop-customizer/contexts/ComponentContext';
import Login from './pages/login';
import Dashboard from './pages/dashboard';
import StoreDesigner from './pages/shop-customizer';
import ThemePreview from './pages/shop-customizer/pages/ThemePreview';
import StandaloneThemePreview from './pages/shop-customizer/pages/StandaloneThemePreview';

// Import all page components
import Products from './pages/products';
import Category from './pages/category';
import CategoryCreate from './pages/category/create';
import CategoryEdit from './pages/category/edit';
import Orders from './pages/orders';
import OrderCreate from './pages/orders/create';
import OrderDetails from './pages/orders/details';
import ProductCreate from './pages/products/create';
import ProductEdit from './pages/products/edit';
import Customers from './pages/customers';
import CustomerCreate from './pages/customers/create';
import CustomerEdit from './pages/customers/edit';
import CustomerDetails from './pages/customers/details';
import Reviews from './pages/reviews';
import PettyCash from './pages/petty-cash';
import PettyCashNew from './pages/petty-cash/new';
import PettyCashEdit from './pages/petty-cash/edit';
import Affiliates from './pages/affiliates';
import AffiliateDetails from './pages/affiliates/details';
import AffiliateRegister from './pages/affiliates/register';
import Support from './pages/support';
import SupportCreate from './pages/support/create';
import SupportDetails from './pages/support/details';
import Settings from './pages/settings';
import UserRoles from './pages/user-roles';
import UserRoleCreate from './pages/user-roles/create';
import UserRoleEdit from './pages/user-roles/edit';
// import Logs from './pages/logs'; // TODO: Create logs page
import Reports from './pages/reports';
import ReportCreate from './pages/reports/create';
import Discounts from './pages/discounts';
import DiscountCreate from './pages/discounts/create';
import DiscountEdit from './pages/discounts/edit';
import Banners from './pages/banners';
import BannerCreate from './pages/banners/create';
import BannerEdit from './pages/banners/edit';
import DashboardLayout from './layouts/DashboardLayout';
import ProtectedRoute from './components/ProtectedRoute';
import './App.css';

// Create a client
const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <ComponentProvider>
            <Router>
              <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/standalone-preview/:themeId" element={<StandaloneThemePreview />} />

                <Route path="/dashboard" element={
                  <ProtectedRoute>
                    <DashboardLayout />
                  </ProtectedRoute>
                }>
                  <Route index element={<Dashboard />} />
                  <Route path="theme-preview/:themeId" element={<ThemePreview />} />
                </Route>
                <Route path="/" element={
                  <ProtectedRoute>
                    <DashboardLayout />
                  </ProtectedRoute>
                }>
                  <Route index element={<Navigate to="/dashboard" replace />} />
                  <Route path="products" element={<Products />} />
                  <Route path="products/create" element={<ProductCreate />} />
                  <Route path="products/:id/edit" element={<ProductEdit />} />
                  <Route path="category" element={<Category />} />
                  <Route path="category/create" element={<CategoryCreate />} />
                  <Route path="category/:id/edit" element={<CategoryEdit />} />
                  <Route path="orders" element={<Orders />} />
                  <Route path="orders/create" element={<OrderCreate />} />
                  <Route path="orders/:id/details" element={<OrderDetails />} />
                  <Route path="customers" element={<Customers />} />
                  <Route path="customers/create" element={<CustomerCreate />} />
                  <Route path="customers/edit/:id" element={<CustomerEdit />} />
                  <Route path="customers/:id/details" element={<CustomerDetails />} />
                  <Route path="reviews" element={<Reviews />} />
                  <Route path="petty-cash" element={<PettyCash />} />
                  <Route path="petty-cash/new" element={<PettyCashNew />} />
                  <Route path="petty-cash/:id/edit" element={<PettyCashEdit />} />
                  <Route path="affiliates" element={<Affiliates />} />
                  <Route path="affiliates/register" element={<AffiliateRegister />} />
                  <Route path="affiliates/:id/details" element={<AffiliateDetails />} />
                  <Route path="support" element={<Support />} />
                  <Route path="support/create" element={<SupportCreate />} />
                  <Route path="support/:id/details" element={<SupportDetails />} />
                  <Route path="settings" element={<Settings />} />
                  <Route path="user-roles" element={<UserRoles />} />
                  <Route path="user-roles/create" element={<UserRoleCreate />} />
                  <Route path="user-roles/:id/edit" element={<UserRoleEdit />} />
                  {/* <Route path="logs" element={<Logs />} /> */}
                  <Route path="reports" element={<Reports />} />
                  <Route path="reports/create" element={<ReportCreate />} />
                  <Route path="discounts" element={<Discounts />} />
                  <Route path="discounts/create" element={<DiscountCreate />} />
                  <Route path="discounts/:id/edit" element={<DiscountEdit />} />
                  <Route path="banners" element={<Banners />} />
                  <Route path="banners/create" element={<BannerCreate />} />
                  <Route path="banners/:id/edit" element={<BannerEdit />} />
                  <Route path="customize" element={<StoreDesigner />} />
                </Route>
              </Routes>
            </Router>
            <Toaster />
            <ReactQueryDevtools initialIsOpen={false} />
        </ComponentProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
