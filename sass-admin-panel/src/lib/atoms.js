import { atom } from 'jotai'

// User session state
export const userAtom = atom(null)

// Theme state
export const themeAtom = atom('light')

// Sidebar collapsed state
export const sidebarCollapsedAtom = atom(false)

// Global loading state
export const globalLoadingAtom = atom(false)