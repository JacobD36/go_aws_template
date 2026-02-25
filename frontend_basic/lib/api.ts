// API client

import { API_BASE_URL, ENDPOINTS } from './constants';
import { getAuthToken } from './auth';
import type {
  LoginCredentials,
  LoginResponse,
  Employee,
  CreateEmployeeRequest,
  ApiError,
} from '@/types';

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      const errorText = await response.text();
      const error: ApiError = {
        message: errorText || 'An error occurred',
        status: response.status,
      };
      throw error;
    }
    return response.json();
  }

  private getHeaders(includeAuth: boolean = false): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (includeAuth) {
      const token = getAuthToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    return headers;
  }

  async login(credentials: LoginCredentials): Promise<LoginResponse> {
    const response = await fetch(`${this.baseUrl}${ENDPOINTS.LOGIN}`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(credentials),
    });

    return this.handleResponse<LoginResponse>(response);
  }

  async getEmployees(): Promise<Employee[]> {
    const response = await fetch(`${this.baseUrl}${ENDPOINTS.EMPLOYEES}`, {
      method: 'GET',
      headers: this.getHeaders(true),
    });

    return this.handleResponse<Employee[]>(response);
  }

  async createEmployee(data: CreateEmployeeRequest): Promise<Employee> {
    const response = await fetch(`${this.baseUrl}${ENDPOINTS.EMPLOYEES}`, {
      method: 'POST',
      headers: this.getHeaders(true),
      body: JSON.stringify(data),
    });

    return this.handleResponse<Employee>(response);
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
