// Constants for the application

export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export const ENDPOINTS = {
  LOGIN: '/auth/login',
  EMPLOYEES: '/employees',
} as const;

export const STORAGE_KEYS = {
  TOKEN: 'auth_token',
  USER_ID: 'user_id',
  EXPIRES_AT: 'expires_at',
} as const;

export const PASSWORD_REQUIREMENTS = {
  MIN_LENGTH: 8,
  MESSAGE: 'La contraseña debe tener mínimo 8 caracteres, una letra mayúscula, un número y un carácter especial',
};
