import { atom } from 'jotai';
import { authApi, apiClient } from '@/lib/api';

// Auth atoms
export const userAtom = atom(null);
export const isAuthenticatedAtom = atom(false);
export const isLoadingAtom = atom(false);

// Derived atom for auth status
export const authStatusAtom = atom((get) => ({
  user: get(userAtom),
  isAuthenticated: get(isAuthenticatedAtom),
  isLoading: get(isLoadingAtom)
}));

// Auth actions atom
export const authActionsAtom = atom(
  null,
  (_get, set, action) => {
    switch (action.type) {
      case 'SET_USER':
        set(userAtom, action.payload);
        set(isAuthenticatedAtom, !!action.payload);
        if (action.payload) {
          localStorage.setItem('user', JSON.stringify(action.payload));
        }
        break;
      case 'SET_LOADING':
        set(isLoadingAtom, action.payload);
        break;
      case 'LOGOUT':
        set(userAtom, null);
        set(isAuthenticatedAtom, false);
        apiClient.setAuthToken(null);
        localStorage.removeItem('auth_token');
        localStorage.removeItem('user');
        break;
      case 'LOGIN_SUCCESS':
        set(userAtom, action.payload.user);
        set(isAuthenticatedAtom, true);
        apiClient.setAuthToken(action.payload.token);
        localStorage.setItem('auth_token', action.payload.token);
        localStorage.setItem('user', JSON.stringify(action.payload.user));
        break;
      default:
        break;
    }
  }
);

// Initialize auth state from localStorage
export const initAuthAtom = atom(
  null,
  (get, set) => {
    const currentUser = get(userAtom);
    const isCurrentlyAuthenticated = get(isAuthenticatedAtom);
    
    // Don't re-initialize if already authenticated
    if (isCurrentlyAuthenticated && currentUser) {
      return;
    }
    
    const token = localStorage.getItem('auth_token');
    const storedUser = localStorage.getItem('user');
    
    if (token && storedUser) {
      try {
        const userData = JSON.parse(storedUser);
        apiClient.setAuthToken(token);
        set(userAtom, userData);
        set(isAuthenticatedAtom, true);
        console.log('Auth restored from localStorage:', userData.email || userData.name || 'User');
      } catch (parseError) {
        console.error('Failed to parse stored user data:', parseError);
        localStorage.removeItem('user');
        localStorage.removeItem('auth_token');
        apiClient.setAuthToken(null);
        set(userAtom, null);
        set(isAuthenticatedAtom, false);
      }
    } else {
      // No token or user data, ensure state is clean
      localStorage.removeItem('user');
      localStorage.removeItem('auth_token');
      apiClient.setAuthToken(null);
      set(userAtom, null);
      set(isAuthenticatedAtom, false);
    }
  }
);

// Login action atom
export const loginAtom = atom(
  null,
  async (_get, set, credentials) => {
    set(isLoadingAtom, true);
    
    try {
      const response = await authApi.login(credentials);
      if (response.success) {
        set(authActionsAtom, { 
          type: 'LOGIN_SUCCESS', 
          payload: response.data 
        });
        return { success: true, data: response.data };
      }
      return { success: false, message: response.message };
    } catch (error) {
      console.error('Login failed:', error);
      return { 
        success: false, 
        message: error.message || 'Login failed' 
      };
    } finally {
      set(isLoadingAtom, false);
    }
  }
);

// Logout action atom
export const logoutAtom = atom(
  null,
  async (_get, set) => {
    try {
      await authApi.logout();
    } catch (error) {
      console.error('Logout API call failed:', error);
    } finally {
      set(authActionsAtom, { type: 'LOGOUT' });
    }
  }
);