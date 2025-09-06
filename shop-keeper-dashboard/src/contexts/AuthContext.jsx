import React, { createContext, useContext, useEffect } from 'react';
import { useAtom } from 'jotai';
import {
  authStatusAtom,
  authActionsAtom,
  initAuthAtom,
  loginAtom,
  logoutAtom
} from '../store/auth';

const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [authStatus] = useAtom(authStatusAtom);
  const [, setAuthAction] = useAtom(authActionsAtom);
  const [, initAuth] = useAtom(initAuthAtom);
  const [, login] = useAtom(loginAtom);
  const [, logout] = useAtom(logoutAtom);

  // Initialize auth state on mount
  useEffect(() => {
    initAuth();
  }, [initAuth]);

  const value = {
    user: authStatus.user,
    isAuthenticated: authStatus.isAuthenticated,
    isLoading: authStatus.isLoading,
    login,
    logout,
    setUser: (user) => setAuthAction({ type: 'SET_USER', payload: user }),
    setLoading: (loading) => setAuthAction({ type: 'SET_LOADING', payload: loading })
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;