// Types for the application

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user_id: string;
  expires_at: number;
}

export interface Employee {
  id: string;
  name: string;
  email: string;
  created_at: string;
}

export interface CreateEmployeeRequest {
  name: string;
  email: string;
  password: string;
}

export interface ApiError {
  message: string;
  status: number;
}
