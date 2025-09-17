import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || '';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
}

export interface AuthError {
  message: string;
}

const setToken = (token: string): void => {
  localStorage.setItem('auth_token', token);
};

export const getToken = (): string | null => {
  return localStorage.getItem('auth_token');
};

export const removeToken = (): void => {
  localStorage.removeItem('auth_token');
};

export const isAuthenticated = (): boolean => {
  return !!getToken();
};

export const setupAxiosInterceptors = (): void => {
  axios.interceptors.request.use(
    (config) => {
      const token = getToken();
      if (token) {
        config.headers.Authorization = token;
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );
};

export const login = async (credentials: LoginRequest): Promise<string> => {
  try {
    const response = await axios.post<LoginResponse>(
      `${API_URL}/api/v1/auth/login`,
      credentials
    );

    const { token } = response.data;
    setToken(token);
    setupAxiosInterceptors();
    return token;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      const message = error.response?.data?.error || 'Login failed';
      throw new Error(message);
    }
    throw new Error('Login failed. Please try again.');
  }
};

// Handle logout
export const logout = (): void => {
  removeToken();
  // Optional: Redirect to login page
  window.location.href = '/login';
};
