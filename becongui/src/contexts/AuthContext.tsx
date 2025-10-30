import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { getToken, login as authLogin, logout as authLogout, setupAxiosInterceptors } from '../services/authService';
import type { LoginRequest } from '../services/authService';

interface AuthContextType {
  isAuthenticated: boolean;
  token: string | null;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const bypassAuth = process.env.REACT_APP_BYPASS_AUTH === 'true'; // For testing
  const [token, setToken] = useState<string | null>(bypassAuth ? 'test-token' : getToken());
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    if (!bypassAuth) {
      setupAxiosInterceptors();
    }
    setLoading(false);
  }, [bypassAuth]);

  const login = async (credentials: LoginRequest): Promise<void> => {
    if (bypassAuth) {
      setToken('test-token');
      setLoading(false);
      return;
    }
    setLoading(true);
    try {
      const newToken = await authLogin(credentials);
      setToken(newToken);
    } finally {
      setLoading(false);
    }
  };

  const logout = (): void => {
    if (!bypassAuth) {
      authLogout();
      setToken(null);
    }
  };

  const value = {
    isAuthenticated: bypassAuth || !!token, // Derive from token
    token,
    login,
    logout,
    loading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
