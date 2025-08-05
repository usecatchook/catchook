// lib/api.ts
import { ApiResponse, PaginatedResponse } from "@/types/api";
import { AuthResponse, LoginCredentials } from "@/types/auth";
import { HealthCheckResponse } from '@/types/health';
import { SetupAdminUserRequest } from "@/types/setup";
import { CreateUserRequest, UpdateUserRequest, User, UserFilters } from "@/types/user";
import axios, { AxiosResponse } from 'axios';
import Cookies from 'js-cookie';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const apiClient = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000,
    headers: { 'Content-Type': 'application/json' },
});

// Intercepteur avec auto-refresh token
apiClient.interceptors.request.use((config) => {
    const token = Cookies.get('authToken');
    if (token) config.headers.Authorization = `Bearer ${token}`;
    return config;
});

apiClient.interceptors.response.use(
    (response: AxiosResponse) => response,
    async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            try {
                const refreshToken = Cookies.get('refreshToken');
                if (refreshToken) {
                    const response = await axios.post<ApiResponse<AuthResponse>>(
                        `${API_BASE_URL}/auth/refresh`,
                        { refreshToken }
                    );

                    const { access_token } = response.data.data.tokens;
                    Cookies.set('authToken', access_token);

                    originalRequest.headers.Authorization = `Bearer ${access_token}`;
                    return apiClient(originalRequest);
                }
            } catch {
                Cookies.remove('authToken');
                Cookies.remove('refreshToken');
                window.location.href = '/login';
            }
        }

        if (error.response?.data?.error?.message) {
            error.message = error.response.data.error.message;
        } else if (error.response?.data?.errors) {
            const errors = error.response.data.errors;
            if (typeof errors === 'object') {
                const errorMessages = Object.values(errors).filter(msg => typeof msg === 'string');
                error.message = errorMessages.join(', ');
                error.validationErrors = errors;
            }
        }

        return Promise.reject(error);
    }
);

// API Methods typés
export const authAPI = {
    login: async (credentials: LoginCredentials): Promise<AuthResponse> => {
        const { data } = await apiClient.post<ApiResponse<AuthResponse>>('/auth/login', credentials);
        return data.data;
    },

    getCurrentUser: async (): Promise<User> => {
        const { data } = await apiClient.get<ApiResponse<User>>('/users/me');
        return data.data;
    },

    logout: async (): Promise<void> => {
        await apiClient.post<ApiResponse<null>>('/auth/logout');
    },
};

export const healthAPI = {
    getHealth: async (): Promise<HealthCheckResponse> => {
        const { data } = await apiClient.get<HealthCheckResponse>('/health');
        return data;
    },
};

export const usersAPI = {
    getUsers: async (filters: UserFilters = {}): Promise<PaginatedResponse<User>> => {
        const { data } = await apiClient.get<PaginatedResponse<User>>('/users', { params: filters });
        return data;
    },

    getUser: async (id: number): Promise<User> => {
        const { data } = await apiClient.get<ApiResponse<User>>(`/users/${id}`);
        return data.data;
    },

    createUser: async (userData: CreateUserRequest): Promise<User> => {
        const { data } = await apiClient.post<ApiResponse<User>>('/users', userData);
        return data.data;
    },

    updateUser: async (userData: UpdateUserRequest): Promise<User> => {
        const { id, ...updateData } = userData;
        const { data } = await apiClient.put<ApiResponse<User>>(`/users/${id}`, updateData);
        return data.data;
    },

    deleteUser: async (id: number): Promise<void> => {
        await apiClient.delete<ApiResponse<null>>(`/users/${id}`);
    },
};

export const setupAPI = {
    createAdminUser: async (userData: SetupAdminUserRequest): Promise<void> => {
        await apiClient.post<ApiResponse<void>>('/setup', userData);
    },
};