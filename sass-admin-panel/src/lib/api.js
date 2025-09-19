import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken')
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
      localStorage.removeItem('authToken')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Single API object with all admin endpoints
export const adminApi = {
  // Authentication
  login: (credentials) => api.post('/auth/login', credentials),
  logout: () => api.post('/auth/logout'),
  getProfile: () => api.get('/auth/me'),

  // Dashboard Analytics
  getDashboard: () => api.get('/analytics/dashboard'),
  getDashboardStats: () => api.get('/analytics/dashboard'),
  getRealtimeStats: () => api.get('/analytics/realtime'),
  getTrafficStats: () => api.get('/analytics/traffic'),
  getSalesStats: () => api.get('/analytics/sales'),

  // Tenants Management
  getTenants: (params) => api.get('/tenant', { params }),
  getTenant: (id) => api.get(`/tenant/${id}`),
  createTenant: (data) => api.post('/tenant', data),
  updateTenant: (id, data) => api.put(`/tenant/${id}`, data),
  deleteTenant: (id) => api.delete(`/tenant/${id}`),
  getTenantStats: (id) => api.get(`/tenant/${id}/stats`),
  
  // Analytics & Reports
  getRevenueAnalytics: (period, filters) => api.get('/analytics/sales', { params: { period, ...filters } }),
  getTenantAnalytics: (period, filters) => api.get('/analytics/dashboard', { params: { period, ...filters } }),
  getPerformanceAnalytics: (period) => api.get('/analytics/realtime', { params: { period } }),
  getTopProducts: () => api.get('/analytics/top/products'),
  getTopPages: () => api.get('/analytics/top/pages'),
  getTopReferrers: () => api.get('/analytics/top/referrers'),
  getCohortAnalysis: () => api.get('/analytics/advanced/cohorts'),
  getFunnelAnalysis: () => api.get('/analytics/advanced/funnel'),
  getCustomerLifetimeValue: () => api.get('/analytics/advanced/clv'),
  getRetentionRate: () => api.get('/analytics/advanced/retention'),

  // Users Management  
  getUsers: (params) => api.get('/user', { params }),
  getUser: (id) => api.get(`/user/${id}`),
  createUser: (data) => api.post('/user', data),
  updateUser: (id, data) => api.put(`/user/${id}`, data),
  deleteUser: (id) => api.delete(`/user/${id}`),

  // Billing & Subscriptions
  getBillingPlans: () => api.get('/billing/plans'),
  getBillingPlan: (id) => api.get(`/billing/plans/${id}`),
  createBillingPlan: (data) => api.post('/billing/plans', data),
  updateBillingPlan: (id, data) => api.put(`/billing/plans/${id}`, data),
  deleteBillingPlan: (id) => api.delete(`/billing/plans/${id}`),
  getSubscriptions: (params) => api.get('/billing/subscriptions', { params }),
  getSubscription: (id) => api.get(`/billing/subscriptions/${id}`),
  createSubscription: (data) => api.post('/billing/subscriptions', data),
  updateSubscription: (id, data) => api.put(`/billing/subscriptions/${id}`, data),
  cancelSubscription: (id) => api.delete(`/billing/subscriptions/${id}`),
  upgradePlan: (data) => api.post('/billing/subscriptions/upgrade', data),
  downgradePlan: (data) => api.post('/billing/subscriptions/downgrade', data),
  recordUsage: (data) => api.post('/billing/usage', data),
  getUsageSummary: () => api.get('/billing/usage'),

  // Payments
  getPayments: (params) => api.get('/payment', { params }),
  getPayment: (id) => api.get(`/payment/${id}`),
  getPaymentMethods: (params) => api.get('/payment/methods', { params }),

  // Support & Communications
  getSupportTickets: (params) => api.get('/support/tickets', { params }),
  getSupportTicket: (id) => api.get(`/support/tickets/${id}`),
  createSupportTicket: (data) => api.post('/support/tickets', data),
  updateSupportTicket: (id, data) => api.put(`/support/tickets/${id}`, data),
  getCommunications: () => api.get('/notification/communications'),
  getAnnouncements: () => api.get('/notification/announcements'),

  // Settings
  getSettings: () => api.get('/settings'),
  updateSettings: (data) => api.put('/settings', data),

  // Reports
  generateReport: (data) => api.post('/analytics/reports/generate', data),
  scheduleReport: (data) => api.post('/analytics/reports/schedule', data),
  getScheduledReports: () => api.get('/analytics/reports/scheduled'),
  deleteScheduledReport: (id) => api.delete(`/analytics/reports/scheduled/${id}`),
}

// Cache invalidation helper
export const invalidateCache = {
  dashboard: (queryClient) => {
    queryClient.invalidateQueries({ queryKey: ['dashboard'] })
  },
  tenants: (queryClient, tenantId) => {
    queryClient.invalidateQueries({ queryKey: ['tenants'] })
    if (tenantId) queryClient.invalidateQueries({ queryKey: ['tenants', tenantId] })
    queryClient.invalidateQueries({ queryKey: ['dashboard'] })
  },
  users: (queryClient, userId) => {
    queryClient.invalidateQueries({ queryKey: ['users'] })
    if (userId) queryClient.invalidateQueries({ queryKey: ['users', userId] })
  },
  subscriptions: (queryClient, subscriptionId) => {
    queryClient.invalidateQueries({ queryKey: ['subscriptions'] })
    if (subscriptionId) queryClient.invalidateQueries({ queryKey: ['subscriptions', subscriptionId] })
    queryClient.invalidateQueries({ queryKey: ['dashboard'] })
  },
  billing: (queryClient) => {
    queryClient.invalidateQueries({ queryKey: ['billing'] })
    queryClient.invalidateQueries({ queryKey: ['dashboard'] })
  },
  support: (queryClient) => {
    queryClient.invalidateQueries({ queryKey: ['support'] })
  },
}

export default api