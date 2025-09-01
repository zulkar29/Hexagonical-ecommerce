// TODO: Implement API service during development

import axios from 'axios'

// Create API client
const api = axios.create({
  // TODO: Configure base URL and defaults
  // baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1',
  // timeout: 10000,
})

// TODO: Add request interceptor for auth token
// api.interceptors.request.use((config) => {
//   const token = localStorage.getItem('token')
//   if (token) {
//     config.headers.Authorization = `Bearer ${token}`
//   }
//   
//   // Add tenant context
//   const tenant = getTenantContext()
//   if (tenant) {
//     config.headers['X-Tenant-ID'] = tenant.id
//   }
//   
//   return config
// })

// TODO: Add response interceptor for error handling
// api.interceptors.response.use(
//   (response) => response,
//   (error) => {
//     if (error.response?.status === 401) {
//       // Redirect to login
//       window.location.href = '/login'
//     }
//     return Promise.reject(error)
//   }
// )

// API services
export const authApi = {
  // TODO: Implement auth API calls
  // login: (email, password) => api.post('/auth/login', { email, password }),
  // logout: () => api.post('/auth/logout'),
  // me: () => api.get('/auth/me'),
  // refresh: () => api.post('/auth/refresh'),
}

export const productsApi = {
  // TODO: Implement products API calls
  // getAll: (params) => api.get('/products', { params }),
  // create: (product) => api.post('/products', product),
  // update: (id, product) => api.put(`/products/${id}`, product),
  // delete: (id) => api.delete(`/products/${id}`),
}

export const ordersApi = {
  // TODO: Implement orders API calls  
  // getAll: (params) => api.get('/orders', { params }),
  // getById: (id) => api.get(`/orders/${id}`),
  // updateStatus: (id, status) => api.patch(`/orders/${id}/status`, { status }),
}

export const customersApi = {
  // TODO: Implement customers API calls
  // getAll: (params) => api.get('/customers', { params }),
  // getById: (id) => api.get(`/customers/${id}`),
}

export const analyticsApi = {
  // TODO: Implement analytics API calls
  // getSalesMetrics: (period) => api.get(`/analytics/sales?period=${period}`),
  // getTopProducts: () => api.get('/analytics/top-products'),
}

export default api