import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { getToken, isAuthenticated, login as authLogin, logout as authLogout, setupAxiosInterceptors, LoginRequest } from '../services/authService';

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
  const [token, setToken] = useState<string | null>(getToken());
  const [loading, setLoading] = useState<boolean>(true);

  // Set up axios interceptors when the component mounts
  useEffect(() => {
    setupAxiosInterceptors();
    setLoading(false);
  }, []);

  const login = async (credentials: LoginRequest): Promise<void> => {
    setLoading(true);
    try {
      const newToken = await authLogin(credentials);
      setToken(newToken);
    } finally {
      setLoading(false);
    }
  };

  const logout = (): void => {
    authLogout();
    setToken(null);
  };

  const value = {
    isAuthenticated: !!token,
    token,
    login,
    logout,
    loading
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
