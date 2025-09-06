const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost/api/v1';

class ApiError extends Error {
  constructor(message, status, errors = null) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.errors = errors;
  }
}

class ApiClient {
  constructor() {
    this.baseURL = API_BASE_URL;
    this.token = null; // Will be set via setAuthToken or getAuthHeaders
    this.authCallbacks = new Set();
  }

  setAuthToken(token) {
    this.token = token;
    if (token) {
      localStorage.setItem('auth_token', token);
    } else {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user');
    }
  }

  addAuthCallback(callback) {
    this.authCallbacks.add(callback);
  }

  removeAuthCallback(callback) {
    this.authCallbacks.delete(callback);
  }

  notifyAuthCallbacks(type, data = null) {
    this.authCallbacks.forEach(callback => {
      try {
        callback(type, data);
      } catch (error) {
        console.error('Auth callback error:', error);
      }
    });
  }

  getAuthHeaders() {
    const headers = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };

    // Always get fresh token from localStorage in case it was updated
    const currentToken = this.token || localStorage.getItem('auth_token');
    if (currentToken) {
      headers.Authorization = `Bearer ${currentToken}`;
    }

    return headers;
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    
    const config = {
      headers: this.getAuthHeaders(),
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        if (response.status === 401) {
          this.handleUnauthorized();
          throw new ApiError('Unauthorized', 401);
        }

        let data;
        try {
          data = await response.json();
        } catch {
          data = { message: 'An error occurred' };
        }

        throw new ApiError(
          data.message || 'An error occurred',
          response.status,
          data.errors || null
        );
      }

      const data = await response.json();
      return data;
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }
      
      throw new ApiError('Network error occurred', 0);
    }
  }

  handleUnauthorized() {
    console.warn('Received 401 Unauthorized response');
    this.setAuthToken(null);
    this.notifyAuthCallbacks('LOGOUT');
    
    // Don't auto-redirect, let the ProtectedRoute handle it
  }

  // GET request
  async get(endpoint, params = {}) {
    // Filter out undefined, null, and empty string values
    const filteredParams = Object.entries(params).reduce((acc, [key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        acc[key] = value;
      }
      return acc;
    }, {});
    
    const queryString = new URLSearchParams(filteredParams).toString();
    const url = queryString ? `${endpoint}?${queryString}` : endpoint;
    
    return this.request(url, {
      method: 'GET',
    });
  }

  // POST request
  async post(endpoint, data = {}) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // PUT request
  async put(endpoint, data = {}) {
    return this.request(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  // DELETE request
  async delete(endpoint) {
    return this.request(endpoint, {
      method: 'DELETE',
    });
  }
}

// Create API client instance
const apiClient = new ApiClient();

// Authentication API
export const authApi = {
  login: (credentials) => apiClient.post('/auth/login', credentials),
  register: (userData) => apiClient.post('/auth/register', userData),
  logout: () => apiClient.post('/auth/logout'),
  me: () => apiClient.get('/auth/me'),
  updateProfile: (data) => apiClient.put('/auth/profile', data),
  changePassword: (data) => apiClient.put('/auth/change-password', data),
  forgotPassword: (data) => apiClient.post('/auth/forgot-password', data),
  resetPassword: (data) => apiClient.post('/auth/reset-password', data),
};

// Products API
export const productsApi = {
  getAll: (params) => apiClient.get('/products', params),
  getById: (id) => apiClient.get(`/products/${id}`),
  create: (data) => apiClient.post('/admin/products', data),
  update: (id, data) => apiClient.put(`/admin/products/${id}`, data),
  delete: (id) => apiClient.delete(`/admin/products/${id}`),
  updateStock: (id, stock) => apiClient.put(`/admin/products/${id}/stock`, { stock }),
  getByCategory: (categoryId, params) => apiClient.get(`/categories/${categoryId}/products`, params),
};

// Categories API
export const categoriesApi = {
  getAll: (params) => apiClient.get('/categories', params),
  getById: (id) => apiClient.get(`/categories/${id}`),
  create: (data) => apiClient.post('/admin/categories', data),
  update: (id, data) => apiClient.put(`/admin/categories/${id}`, data),
  delete: (id) => apiClient.delete(`/admin/categories/${id}`),
};

// Orders API
export const ordersApi = {
  getAll: (params) => apiClient.get('/orders', params),
  getAdminAll: (params) => apiClient.get('/admin/orders', params),
  getById: (id) => apiClient.get(`/orders/${id}`),
  getAdminById: (id) => apiClient.get(`/admin/orders/${id}`),
  create: (data) => apiClient.post('/orders', data),
  updateStatus: (id, status) => apiClient.put(`/admin/orders/${id}/status`, { status }),
  cancel: (id) => apiClient.put(`/orders/${id}/cancel`),
  getAnalytics: () => apiClient.get('/admin/orders/analytics'),
};

// Users API
export const usersApi = {
  getAll: (params) => apiClient.get('/admin/users', params),
  getById: (id) => apiClient.get(`/admin/users/${id}`),
  create: (data) => apiClient.post('/admin/users', data),
  update: (id, data) => apiClient.put(`/admin/users/${id}`, data),
  updateStatus: (id, isActive) => apiClient.put(`/admin/users/${id}/status`, { is_active: isActive }),
  delete: (id) => apiClient.delete(`/admin/users/${id}`),
  getProfile: () => apiClient.get('/users/profile'),
  updateProfile: (data) => apiClient.put('/users/profile', data),
};

// Coupons API
export const couponsApi = {
  getAll: (params) => apiClient.get('/admin/coupons', params),
  getAvailable: () => apiClient.get('/coupons/available'),
  getById: (id) => apiClient.get(`/admin/coupons/${id}`),
  create: (data) => apiClient.post('/admin/coupons', data),
  update: (id, data) => apiClient.put(`/admin/coupons/${id}`, data),
  delete: (id) => apiClient.delete(`/admin/coupons/${id}`),
  validate: (code, orderAmount) => apiClient.post('/coupons/validate', { code, order_amount: orderAmount }),
  getUsageStats: (id) => apiClient.get(`/admin/coupons/${id}/usage`),
};

// Reviews API
export const reviewsApi = {
  getAll: (params) => apiClient.get('/admin/reviews', params),
  getUserReviews: () => apiClient.get('/reviews/my-reviews'),
  getProductReviews: (productId, params) => apiClient.get(`/products/${productId}/reviews`, params),
  create: (data) => apiClient.post('/reviews', data),
  update: (id, data) => apiClient.put(`/reviews/${id}`, data),
  delete: (id) => apiClient.delete(`/reviews/${id}`),
  approve: (id) => apiClient.put(`/admin/reviews/${id}/approve`),
  reject: (id) => apiClient.put(`/admin/reviews/${id}/reject`),
  adminDelete: (id) => apiClient.delete(`/admin/reviews/${id}`),
};

// Payments API
export const paymentsApi = {
  getAll: (params) => apiClient.get('/admin/payments', params),
  create: (data) => apiClient.post('/payments', data),
  getOrderPayments: (orderId) => apiClient.get(`/payments/orders/${orderId}`),
  getAnalytics: () => apiClient.get('/admin/payments/analytics'),
};

// Dashboard API
export const dashboardApi = {
  getStats: () => apiClient.get('/dashboard/stats'),
  getRecentOrders: () => apiClient.get('/admin/dashboard/recent-orders'),
  getAnalytics: () => apiClient.get('/dashboard/analytics'),
  getAdminDashboard: () => apiClient.get('/admin/dashboard'),
  getRevenueAnalytics: (period) => apiClient.get('/admin/dashboard/revenue-analytics', { period }),
  getRecentActivity: () => apiClient.get('/admin/dashboard/recent-activity'),
  getLowStockProducts: () => apiClient.get('/admin/dashboard/low-stock-products'),
  getQuickActionsData: () => apiClient.get('/admin/dashboard/quick-actions'),
};

// Affiliates API
export const affiliatesApi = {
  // Admin endpoints
  getAll: (params) => apiClient.get('/admin/affiliates', params),
  getById: (id) => apiClient.get(`/admin/affiliates/${id}`),
  create: (data) => apiClient.post('/admin/affiliates', data),
  update: (id, data) => apiClient.put(`/admin/affiliates/${id}`, data),
  updateStatus: (id, status) => apiClient.put(`/admin/affiliates/${id}/status`, { status }),
  delete: (id) => apiClient.delete(`/admin/affiliates/${id}`),
  getCommissions: (params) => apiClient.get('/admin/affiliate-commissions', params),
  getPayouts: (params) => apiClient.get('/admin/affiliate-payouts', params),
  createPayout: (data) => apiClient.post('/admin/affiliate-payouts', data),
  
  // Affiliate user endpoints
  register: (data) => apiClient.post('/affiliates/register', data),
  getDashboard: () => apiClient.get('/affiliates/dashboard'),
  getProfile: () => apiClient.get('/affiliates/profile'),
  updateProfile: (data) => apiClient.put('/affiliates/profile', data),
  getLinks: (params) => apiClient.get('/affiliates/links', params),
  createLink: (data) => apiClient.post('/affiliates/links', data),
  updateLink: (id, data) => apiClient.put(`/affiliates/links/${id}`, data),
  deleteLink: (id) => apiClient.delete(`/affiliates/links/${id}`),
  getCommissions: (params) => apiClient.get('/affiliates/commissions', params),
  getPayouts: (params) => apiClient.get('/affiliates/payouts', params),
  getAnalytics: (params) => apiClient.get('/affiliates/analytics', params),
  requestPayout: (data) => apiClient.post('/affiliates/payout-request', data),
};

// Petty Cash API
export const pettyCashApi = {
  // Get petty cash overview/stats
  getOverview: () => apiClient.get('/admin/petty-cash/overview'),
  
  // Transaction management
  getTransactions: (params) => apiClient.get('/admin/petty-cash/transactions', params),
  getTransaction: (id) => apiClient.get(`/admin/petty-cash/transactions/${id}`),
  createTransaction: (data) => apiClient.post('/admin/petty-cash/transactions', data),
  updateTransaction: (id, data) => apiClient.put(`/admin/petty-cash/transactions/${id}`, data),
  deleteTransaction: (id) => apiClient.delete(`/admin/petty-cash/transactions/${id}`),
  
  // Transaction approval
  approveTransaction: (id) => apiClient.post(`/admin/petty-cash/transactions/${id}/approve`),
  rejectTransaction: (id, reason) => apiClient.post(`/admin/petty-cash/transactions/${id}/reject`, { reason }),
  
  // Balance management
  getCurrentBalance: () => apiClient.get('/admin/petty-cash/balance'),
  replenishCash: (data) => apiClient.post('/admin/petty-cash/replenish', data),
  
  // Categories
  getCategories: () => apiClient.get('/admin/petty-cash/categories'),
  createCategory: (data) => apiClient.post('/admin/petty-cash/categories', data),
  updateCategory: (id, data) => apiClient.put(`/admin/petty-cash/categories/${id}`, data),
  deleteCategory: (id) => apiClient.delete(`/admin/petty-cash/categories/${id}`),
  
  // Reports
  getReports: (params) => apiClient.get('/admin/petty-cash/reports', params),
  exportTransactions: (params) => apiClient.get('/admin/petty-cash/export', params),
  
  // Receipt management
  uploadReceipt: (transactionId, file) => {
    const formData = new FormData();
    formData.append('receipt', file);
    return apiClient.post(`/admin/petty-cash/transactions/${transactionId}/receipt`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
  },
  deleteReceipt: (transactionId) => apiClient.delete(`/admin/petty-cash/transactions/${transactionId}/receipt`),
};

export { apiClient, ApiError };