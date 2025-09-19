import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { Toaster } from 'sonner';
import { AuthProvider } from './contexts/AuthContext';
import { ThemeProvider } from './contexts/ThemeContext';
import { ComponentProvider } from './contexts/ComponentContext';
import Login from './pages/login';
import Dashboard from './pages/dashboard';
import ComponentBuilder from './pages/component-builder';
import ThemePreview from './pages/ThemePreview';
import StandaloneThemePreview from './pages/StandaloneThemePreview';

// Import all page components
import Products from './pages/products';
import Category from './pages/category';
import Orders from './pages/orders';
import Customers from './pages/customers';
import Reviews from './pages/reviews';
import PettyCash from './pages/petty-cash';
import Affiliates from './pages/affiliates';
import Support from './pages/support';
import Settings from './pages/settings';
import UserRoles from './pages/user-roles';
// import Logs from './pages/logs'; // TODO: Create logs page
import Reports from './pages/reports';
import Discounts from './pages/discounts';
import Banners from './pages/banners';
import DashboardLayout from './layouts/DashboardLayout';
import ProtectedRoute from './components/ProtectedRoute';
import './App.css';

// Create a client
const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <ThemeProvider>
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
                  <Route path="component-builder" element={<ComponentBuilder />} />
                  <Route path="theme-preview/:themeId" element={<ThemePreview />} />
                </Route>
                <Route path="/" element={
                  <ProtectedRoute>
                    <DashboardLayout />
                  </ProtectedRoute>
                }>
                  <Route index element={<Navigate to="/dashboard" replace />} />
                  <Route path="products" element={<Products />} />
                  <Route path="category" element={<Category />} />
                  <Route path="orders" element={<Orders />} />
                  <Route path="customers" element={<Customers />} />
                  <Route path="reviews" element={<Reviews />} />
                  <Route path="petty-cash" element={<PettyCash />} />
                  <Route path="affiliates" element={<Affiliates />} />
                  <Route path="support" element={<Support />} />
                  <Route path="settings" element={<Settings />} />
                  <Route path="user-roles" element={<UserRoles />} />
                  {/* <Route path="logs" element={<Logs />} /> */}
                  <Route path="reports" element={<Reports />} />
                  <Route path="discounts" element={<Discounts />} />
                  <Route path="banners" element={<Banners />} />
                  <Route path="builder" element={<ComponentBuilder />} />
                </Route>
              </Routes>
            </Router>
            <Toaster />
            <ReactQueryDevtools initialIsOpen={false} />
          </ComponentProvider>
        </ThemeProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
