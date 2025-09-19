import { atom } from 'jotai'

// User/Shop session state
export const userAtom = atom(null)
export const shopAtom = atom(null)

// Theme state
export const themeAtom = atom('light')

// Sidebar/Navigation state
export const sidebarCollapsedAtom = atom(false)
export const mobileMenuOpenAtom = atom(false)

// Global loading state
export const globalLoadingAtom = atom(false)

// Selected date range for analytics
export const dateRangeAtom = atom({
  from: null,
  to: null,
  preset: 'last7days'
})

// Current shop settings/subscription info
export const subscriptionAtom = atom(null)

// Notification/Toast state
export const notificationsAtom = atom([])

// Search/Filter state for products
export const productFiltersAtom = atom({
  search: '',
  category: 'all',
  status: 'all',
  sortBy: 'name',
  sortOrder: 'asc'
})

// Search/Filter state for orders
export const orderFiltersAtom = atom({
  search: '',
  status: 'all',
  dateRange: 'all',
  sortBy: 'created_at',
  sortOrder: 'desc'
})

// Search/Filter state for customers
export const customerFiltersAtom = atom({
  search: '',
  status: 'all',
  sortBy: 'name',
  sortOrder: 'asc'
})

// Selected items for bulk actions
export const selectedProductsAtom = atom([])
export const selectedOrdersAtom = atom([])
export const selectedCustomersAtom = atom([])

// Cart/Draft order state (for manual order creation)
export const draftOrderAtom = atom({
  items: [],
  customer: null,
  shipping: null,
  discount: null,
  notes: ''
})