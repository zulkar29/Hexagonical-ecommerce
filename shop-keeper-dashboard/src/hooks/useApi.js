import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  productsApi, 
  categoriesApi, 
  ordersApi, 
  usersApi, 
  couponsApi, 
  reviewsApi, 
  paymentsApi,
  dashboardApi,
  authApi
} from '../lib/api';

// Query Keys
export const QUERY_KEYS = {
  PRODUCTS: 'products',
  CATEGORIES: 'categories',
  ORDERS: 'orders',
  USERS: 'users',
  COUPONS: 'coupons',
  REVIEWS: 'reviews',
  PAYMENTS: 'payments',
  DASHBOARD: 'dashboard',
  AUTH: 'auth',
};

// Authentication Hooks - Note: Login/logout are handled by Jotai atoms
// These hooks are kept for React Query integration where needed

export const useMe = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.AUTH, 'me'],
    queryFn: authApi.me,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Products Hooks
export const useProducts = (params = {}) => {
  return useQuery({
    queryKey: [QUERY_KEYS.PRODUCTS, params],
    queryFn: () => productsApi.getAll(params),
    staleTime: 2 * 60 * 1000, // 2 minutes
    placeholderData: (previousData) => previousData, // Keep previous data while loading
  });
};

export const useProduct = (id) => {
  return useQuery({
    queryKey: [QUERY_KEYS.PRODUCTS, id],
    queryFn: () => productsApi.getById(id),
    enabled: !!id,
  });
};

export const useCreateProduct = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: productsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.PRODUCTS] });
    },
  });
};

export const useUpdateProduct = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }) => productsApi.update(id, data),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.PRODUCTS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.PRODUCTS, variables.id] });
    },
  });
};

export const useDeleteProduct = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: productsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.PRODUCTS] });
    },
  });
};

export const useUpdateProductStock = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, stock }) => productsApi.updateStock(id, stock),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.PRODUCTS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.PRODUCTS, variables.id] });
    },
  });
};

// Categories Hooks
export const useCategories = (params = {}) => {
  return useQuery({
    queryKey: [QUERY_KEYS.CATEGORIES, params],
    queryFn: () => categoriesApi.getAll(params),
    staleTime: 5 * 60 * 1000, // 5 minutes
    placeholderData: (previousData) => previousData,
  });
};

export const useCategory = (id) => {
  return useQuery({
    queryKey: [QUERY_KEYS.CATEGORIES, id],
    queryFn: () => categoriesApi.getById(id),
    enabled: !!id,
  });
};

export const useCreateCategory = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: categoriesApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.CATEGORIES] });
    },
  });
};

export const useUpdateCategory = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }) => categoriesApi.update(id, data),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.CATEGORIES] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.CATEGORIES, variables.id] });
    },
  });
};

export const useDeleteCategory = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: categoriesApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.CATEGORIES] });
    },
  });
};

// Orders Hooks
export const useOrders = (params = {}, isAdmin = false) => {
  return useQuery({
    queryKey: [QUERY_KEYS.ORDERS, params, isAdmin ? 'admin' : 'user'],
    queryFn: () => isAdmin ? ordersApi.getAdminAll(params) : ordersApi.getAll(params),
    staleTime: 1 * 60 * 1000, // 1 minute
    placeholderData: (previousData) => previousData,
  });
};

export const useOrder = (id, isAdmin = true) => {
  return useQuery({
    queryKey: [QUERY_KEYS.ORDERS, id, isAdmin ? 'admin' : 'user'],
    queryFn: () => isAdmin ? ordersApi.getAdminById(id) : ordersApi.getById(id),
    enabled: !!id,
  });
};

export const useCreateOrder = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ordersApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.ORDERS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.DASHBOARD] });
    },
  });
};

export const useUpdateOrderStatus = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, status }) => ordersApi.updateStatus(id, status),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.ORDERS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.ORDERS, variables.id] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.DASHBOARD] });
    },
  });
};

export const useCancelOrder = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ordersApi.cancel,
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.ORDERS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.ORDERS, variables] });
    },
  });
};

export const useOrderAnalytics = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.ORDERS, 'analytics'],
    queryFn: ordersApi.getAnalytics,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Users Hooks
export const useUsers = (params = {}) => {
  return useQuery({
    queryKey: [QUERY_KEYS.USERS, params],
    queryFn: () => usersApi.getAll(params),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};

export const useUser = (id) => {
  return useQuery({
    queryKey: [QUERY_KEYS.USERS, id],
    queryFn: () => usersApi.getById(id),
    enabled: !!id,
  });
};

export const useUpdateUser = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }) => usersApi.update(id, data),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.USERS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.USERS, variables.id] });
    },
  });
};

export const useUpdateUserStatus = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, isActive }) => usersApi.updateStatus(id, isActive),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.USERS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.USERS, variables.id] });
    },
  });
};

export const useDeleteUser = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: usersApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.USERS] });
    },
  });
};

// Coupons Hooks
export const useCoupons = (params = {}) => {
  return useQuery({
    queryKey: [QUERY_KEYS.COUPONS, params],
    queryFn: () => couponsApi.getAll(params),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};

export const useAvailableCoupons = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.COUPONS, 'available'],
    queryFn: couponsApi.getAvailable,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useCoupon = (id) => {
  return useQuery({
    queryKey: [QUERY_KEYS.COUPONS, id],
    queryFn: () => couponsApi.getById(id),
    enabled: !!id,
  });
};

export const useCreateCoupon = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: couponsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.COUPONS] });
    },
  });
};

export const useUpdateCoupon = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }) => couponsApi.update(id, data),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.COUPONS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.COUPONS, variables.id] });
    },
  });
};

export const useDeleteCoupon = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: couponsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.COUPONS] });
    },
  });
};

export const useValidateCoupon = () => {
  return useMutation({
    mutationFn: ({ code, orderAmount }) => couponsApi.validate(code, orderAmount),
  });
};

// Dashboard Hooks
export const useDashboardStats = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.DASHBOARD, 'stats'],
    queryFn: dashboardApi.getStats,
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchInterval: 5 * 60 * 1000, // Refetch every 5 minutes
  });
};

export const useAdminDashboard = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.DASHBOARD, 'admin'],
    queryFn: dashboardApi.getAdminDashboard,
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchInterval: 5 * 60 * 1000, // Refetch every 5 minutes
  });
};

export const useDashboardAnalytics = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.DASHBOARD, 'analytics'],
    queryFn: dashboardApi.getAnalytics,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useRecentOrders = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.DASHBOARD, 'recent-orders'],
    queryFn: dashboardApi.getRecentOrders,
    staleTime: 1 * 60 * 1000, // 1 minute
    refetchInterval: 2 * 60 * 1000, // Refetch every 2 minutes
  });
};

// Revenue Analytics Hook
export const useRevenueAnalytics = (period = 'monthly') => {
  return useQuery({
    queryKey: [QUERY_KEYS.DASHBOARD, 'revenue-analytics', period],
    queryFn: () => dashboardApi.getRevenueAnalytics(period),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Recent Activity Hook
export const useRecentActivity = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.DASHBOARD, 'recent-activity'],
    queryFn: dashboardApi.getRecentActivity,
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchInterval: 3 * 60 * 1000, // 3 minutes
  });
};

// Low Stock Products Hook
export const useLowStockProducts = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.PRODUCTS, 'low-stock'],
    queryFn: dashboardApi.getLowStockProducts,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Quick Actions Data Hook
export const useQuickActionsData = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.DASHBOARD, 'quick-actions'],
    queryFn: dashboardApi.getQuickActionsData,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Reviews Hooks
export const useReviews = (params = {}) => {
  return useQuery({
    queryKey: [QUERY_KEYS.REVIEWS, params],
    queryFn: () => reviewsApi.getAll(params),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};

export const useProductReviews = (productId, params = {}) => {
  return useQuery({
    queryKey: [QUERY_KEYS.REVIEWS, 'product', productId, params],
    queryFn: () => reviewsApi.getProductReviews(productId, params),
    enabled: !!productId,
  });
};

export const useCreateReview = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: reviewsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.REVIEWS] });
    },
  });
};

export const useApproveReview = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: reviewsApi.approve,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.REVIEWS] });
    },
  });
};

export const useRejectReview = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: reviewsApi.reject,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.REVIEWS] });
    },
  });
};

// Payments Hooks
export const usePayments = (params = {}) => {
  return useQuery({
    queryKey: [QUERY_KEYS.PAYMENTS, params],
    queryFn: () => paymentsApi.getAll(params),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};

export const usePaymentAnalytics = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.PAYMENTS, 'analytics'],
    queryFn: paymentsApi.getAnalytics,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};