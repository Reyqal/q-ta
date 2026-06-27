import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import apiClient from '../lib/apiClient';

export interface User {
  id: number;
  name: string;
  phone_number: string;
  role: 'admin' | 'penghuni';
  room_id?: number;
  room_number?: string;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isAdmin: boolean;
  isTenant: boolean;
  isLoading: boolean;
  login: (phone_number: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem('qta_token'));
  const [isLoading, setIsLoading] = useState(true);

  const fetchUser = useCallback(async () => {
    if (!token) {
      setIsLoading(false);
      return;
    }
    try {
      const response = await apiClient.get('/users/me');
      if (response.data.success) {
        setUser(response.data.data);
        localStorage.setItem('qta_user', JSON.stringify(response.data.data));
      }
    } catch {
      setToken(null);
      setUser(null);
      localStorage.removeItem('qta_token');
      localStorage.removeItem('qta_user');
    } finally {
      setIsLoading(false);
    }
  }, [token]);

  useEffect(() => {
    fetchUser();
  }, [fetchUser]);

  const login = async (phone_number: string, password: string) => {
    const response = await apiClient.post('/auth/login', { phone_number, password });
    if (response.data.success) {
      const { token: newToken, user: newUser } = response.data.data;
      setToken(newToken);
      setUser(newUser);
      localStorage.setItem('qta_token', newToken);
      localStorage.setItem('qta_user', JSON.stringify(newUser));
    } else {
      throw new Error(response.data.message || 'Login gagal');
    }
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('qta_token');
    localStorage.removeItem('qta_user');
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated: !!token && !!user,
        isAdmin: user?.role === 'admin',
        isTenant: user?.role === 'penghuni',
        isLoading,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth harus digunakan di dalam AuthProvider');
  }
  return context;
}
