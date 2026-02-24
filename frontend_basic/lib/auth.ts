// Authentication utilities

import { STORAGE_KEYS } from './constants';
import type { LoginResponse } from '@/types';

export const saveAuthData = (data: LoginResponse): void => {
  if (typeof window !== 'undefined') {
    localStorage.setItem(STORAGE_KEYS.TOKEN, data.token);
    localStorage.setItem(STORAGE_KEYS.USER_ID, data.user_id);
    localStorage.setItem(STORAGE_KEYS.EXPIRES_AT, data.expires_at.toString());
  }
};

export const getAuthToken = (): string | null => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem(STORAGE_KEYS.TOKEN);
  }
  return null;
};

export const clearAuthData = (): void => {
  if (typeof window !== 'undefined') {
    localStorage.removeItem(STORAGE_KEYS.TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER_ID);
    localStorage.removeItem(STORAGE_KEYS.EXPIRES_AT);
  }
};

export const isAuthenticated = (): boolean => {
  if (typeof window === 'undefined') return false;
  
  const token = localStorage.getItem(STORAGE_KEYS.TOKEN);
  const expiresAt = localStorage.getItem(STORAGE_KEYS.EXPIRES_AT);
  
  if (!token || !expiresAt) return false;
  
  // Check if token is expired
  const now = Date.now() / 1000; // Convert to seconds
  return parseInt(expiresAt, 10) > now;
};

export const validatePassword = (password: string): string | null => {
  if (password.length < 8) {
    return 'La contraseña debe tener al menos 8 caracteres';
  }
  
  if (!/[A-Z]/.test(password)) {
    return 'La contraseña debe contener al menos una letra mayúscula';
  }
  
  if (!/[0-9]/.test(password)) {
    return 'La contraseña debe contener al menos un número';
  }
  
  if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?~]/.test(password)) {
    return 'La contraseña debe contener al menos un carácter especial';
  }
  
  return null;
};
