import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Export the api client as apiClient for compatibility
export const apiClient = api;

// Auth API for compatibility with existing code
export const authApi = {
  login: (credentials) => api.post('/auth/login', credentials),
  logout: () => api.post('/auth/logout'),
  getProfile: () => api.get('/auth/me'),
};

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Single API object with all shop endpoints
export const shopApi = {
  // Authentication
  login: (credentials) => api.post('/auth/login', credentials),
  logout: () => api.post('/auth/logout'),
  getProfile: () => api.get('/auth/me'),

  // Dashboard Analytics
  getDashboard: () => api.get('/analytics/dashboard'),
  getRealtimeStats: () => api.get('/analytics/realtime'),
  getSalesStats: () => api.get('/analytics/sales'),
  getTrafficStats: () => api.get('/analytics/traffic'),
  getTopProducts: () => api.get('/analytics/top/products'),

  // Products Management
  getProducts: (params = {}) => api.get('/products', { params }),
  getProduct: (id) => api.get(`/products/${id}`),
  createProduct: (data) => api.post('/products', data),
  updateProduct: (id, data) => api.put(`/products/${id}`, data),
  deleteProduct: (id) => api.delete(`/products/${id}`),
  getProductStats: (params = {}) => api.get('/products', { params: { type: 'stats', ...params } }),
  getLowStockProducts: (params = {}) => api.get('/products', { params: { type: 'low-stock', ...params } }),
  searchProducts: (query, params = {}) => api.get('/products', { params: { type: 'search', q: query, ...params } }),
  updateInventory: (id, data) => api.post(`/products/${id}`, { ...data, action: 'update_inventory' }),
  updateProductStatus: (id, status) => api.post(`/products/${id}`, { status, action: 'status' }),
  duplicateProduct: (id) => api.post(`/products/${id}`, { action: 'duplicate' }),
  bulkUpdateProducts: (data) => api.post('/products', { ...data, operation: 'bulk-update' }),
  importProducts: (data) => api.post('/products', { ...data, operation: 'import' }),
  exportProducts: (params = {}) => api.get('/products', { params: { operation: 'export', ...params } }),

  // Orders Management
  getOrders: (params = {}) => api.get('/orders', { params }),
  getOrder: (id) => api.get(`/orders/${id}`),
  createOrder: (data) => api.post('/orders', data),
  updateOrder: (id, data) => api.put(`/orders/${id}`, data),
  deleteOrder: (id) => api.delete(`/orders/${id}`),
  getOrderStats: (params = {}) => api.get('/orders', { params: { type: 'stats', ...params } }),
  getMyOrders: (params = {}) => api.get('/orders', { params: { type: 'my-orders', ...params } }),
  trackOrder: (trackingId, params = {}) => api.get('/orders', { params: { type: 'track', tracking_id: trackingId, ...params } }),
  exportOrders: (params = {}) => api.get('/orders', { params: { type: 'export', ...params } }),
  bulkUpdateOrders: (data) => api.post('/orders', { ...data, operation: 'bulk-update' }),
  importOrders: (data) => api.post('/orders', { ...data, operation: 'import' }),

  // Customers Management
  getCustomers: (params = {}) => api.get('/user', { params }),
  getCustomer: (id) => api.get(`/user/${id}`),
  createCustomer: (data) => api.post('/user', data),
  updateCustomer: (id, data) => api.put(`/user/${id}`, data),
  deleteCustomer: (id) => api.delete(`/user/${id}`),

  // Categories Management  
  getCategories: (params = {}) => api.get('/category', { params }),
  getCategory: (id) => api.get(`/category/${id}`),
  createCategory: (data) => api.post('/category', data),
  updateCategory: (id, data) => api.put(`/category/${id}`, data),
  deleteCategory: (id) => api.delete(`/category/${id}`),

  // Reviews Management
  getReviews: (params = {}) => api.get('/reviews', { params }),
  getReview: (id) => api.get(`/reviews/${id}`),
  updateReviewStatus: (id, data) => api.put(`/reviews/${id}`, data),
  deleteReview: (id) => api.delete(`/reviews/${id}`),

  // Discounts Management
  getDiscounts: (params = {}) => api.get('/discount', { params }),
  getDiscount: (id) => api.get(`/discount/${id}`),
  createDiscount: (data) => api.post('/discount', data),
  updateDiscount: (id, data) => api.put(`/discount/${id}`, data),
  deleteDiscount: (id) => api.delete(`/discount/${id}`),

  // Shipping Management
  getShippingMethods: (params = {}) => api.get('/shipping', { params }),
  getShippingMethod: (id) => api.get(`/shipping/${id}`),
  createShippingMethod: (data) => api.post('/shipping', data),
  updateShippingMethod: (id, data) => api.put(`/shipping/${id}`, data),
  deleteShippingMethod: (id) => api.delete(`/shipping/${id}`),

  // Payments
  getPayments: (params = {}) => api.get('/payment', { params }),
  getPayment: (id) => api.get(`/payment/${id}`),
  getPaymentMethods: (params = {}) => api.get('/payment/methods', { params }),

  // Notifications
  getNotifications: (params = {}) => api.get('/notification', { params }),
  markNotificationRead: (id) => api.put(`/notification/${id}`, { read: true }),
  markAllNotificationsRead: () => api.put('/notification/mark-all-read'),

  // Wishlist
  getWishlist: (params = {}) => api.get('/wishlist', { params }),
  addToWishlist: (data) => api.post('/wishlist', data),
  removeFromWishlist: (id) => api.delete(`/wishlist/${id}`),

  // Returns
  getReturns: (params = {}) => api.get('/returns', { params }),
  getReturn: (id) => api.get(`/returns/${id}`),
  createReturn: (data) => api.post('/returns', data),
  updateReturnStatus: (id, data) => api.put(`/returns/${id}`, data),

  // Marketing
  getMarketingCampaigns: (params = {}) => api.get('/marketing/campaigns', { params }),
  getCampaign: (id) => api.get(`/marketing/campaigns/${id}`),
  createCampaign: (data) => api.post('/marketing/campaigns', data),
  updateCampaign: (id, data) => api.put(`/marketing/campaigns/${id}`, data),
  deleteCampaign: (id) => api.delete(`/marketing/campaigns/${id}`),

  // Settings
  getSettings: () => api.get('/settings'),
  updateSettings: (data) => api.put('/settings', data),

  // Contact
  getContacts: (params = {}) => api.get('/contact', { params }),
  getContact: (id) => api.get(`/contact/${id}`),
  createContact: (data) => api.post('/contact', data),
  updateContact: (id, data) => api.put(`/contact/${id}`, data),
  deleteContact: (id) => api.delete(`/contact/${id}`),

  // Search
  search: (query, params = {}) => api.get('/search', { params: { q: query, ...params } }),
  getSearchSuggestions: (query) => api.get('/search/suggestions', { params: { q: query } }),

  // Additional missing functions
  getRecentOrders: (limit = 5) => api.get(`/orders?limit=${limit}&sort=created_at&order=desc`),
  getLowStock: (limit = 10) => api.get(`/products?type=low-stock&limit=${limit}`),
  getRevenueAnalytics: (period = 'monthly') => api.get(`/analytics/revenue?period=${period}`),
}

// Add setAuthToken method to apiClient
apiClient.setAuthToken = (token) => {
  if (token) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    localStorage.setItem('auth_token', token);
  } else {
    delete api.defaults.headers.common['Authorization'];
    localStorage.removeItem('auth_token');
  }
};

export default api