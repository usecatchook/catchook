// lib/api.ts
import { ApiResponse, PaginatedResponse } from "@/types/api";
import { AuthResponse, LoginCredentials } from "@/types/auth";
import { HealthCheckResponse } from '@/types/health';
import { SetupAdminUserRequest } from "@/types/setup";
import { CreateUserRequest, ListUsersResponse, UpdateUserRequest, User, UserFilters } from "@/types/user";
import axios, { AxiosResponse } from 'axios';
import Cookies from 'js-cookie';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Durée de session en millisecondes (24h)
const SESSION_DURATION = 24 * 60 * 60 * 1000;
// Marge de sécurité pour refresh (5 minutes avant expiration)
const REFRESH_MARGIN = 5 * 60 * 1000;

export const apiClient = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000,
    headers: { 'Content-Type': 'application/json' },
});

// Fonction pour vérifier si la session doit être refreshée
const shouldRefreshSession = (): boolean => {
    const sessionTimestamp = Cookies.get('session_timestamp');
    if (!sessionTimestamp) return false;
    
    const timestamp = parseInt(sessionTimestamp);
    const now = Date.now();
    const timeUntilExpiry = timestamp + SESSION_DURATION - now;
    
    return timeUntilExpiry <= REFRESH_MARGIN;
};

// Fonction pour refresh la session
const refreshSession = async (): Promise<string | null> => {
    try {
        const sessionId = Cookies.get('session_id');
        if (!sessionId) return null;

        const response = await axios.post<ApiResponse<AuthResponse>>(
            `${API_BASE_URL}/auth/refresh`,
            {},
            {
                headers: {
                    'Authorization': sessionId
                }
            }
        );

        const { session_id } = response.data.data.session;
        const timestamp = Date.now();
        
        Cookies.set('session_id', session_id);
        Cookies.set('session_timestamp', timestamp.toString());
        
        return session_id;
    } catch (error) {
        console.error('Failed to refresh session:', error);
        // Si le refresh échoue, on supprime les cookies et on redirige
        Cookies.remove('session_id');
        Cookies.remove('session_timestamp');
        window.location.href = '/login';
        return null;
    }
};

// Intercepteur pour ajouter le session_id dans Authorization pour toutes les requêtes
apiClient.interceptors.request.use(async (config) => {
    const sessionId = Cookies.get('session_id');
    
    if (sessionId) {
        // Vérifier si on doit refresh la session
        if (shouldRefreshSession()) {
            const newSessionId = await refreshSession();
            if (newSessionId) {
                config.headers['Authorization'] = newSessionId;
            }
        } else {
            config.headers['Authorization'] = sessionId;
        }
    }
    
    return config;
});

apiClient.interceptors.response.use(
    (response: AxiosResponse) => response,
    async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            try {
                // Essayer de refresh la session
                const newSessionId = await refreshSession();
                if (newSessionId) {
                    originalRequest.headers['Authorization'] = newSessionId;
                    return apiClient(originalRequest);
                }
            } catch {
                // Si le refresh échoue, on supprime les cookies et on redirige
                Cookies.remove('session_id');
                Cookies.remove('session_timestamp');
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

// API Methods typés - Gestion technique uniquement
export const authAPI = {
    login: async (credentials: LoginCredentials): Promise<AuthResponse> => {
        const { data } = await apiClient.post<ApiResponse<AuthResponse>>('/auth/login', credentials);
        
        // Stocker le session_id et le timestamp lors du login
        const { session_id } = data.data.session;
        const timestamp = Date.now();
        
        Cookies.set('session_id', session_id);
        Cookies.set('session_timestamp', timestamp.toString());
        
        return data.data;
    },

    getCurrentUser: async (): Promise<User> => {
        const { data } = await apiClient.get<ApiResponse<User>>('/users/me');
        return data.data;
    },

    logout: async (): Promise<void> => {
        await apiClient.post<ApiResponse<null>>('/auth/logout');
        // Supprimer les cookies lors du logout
        Cookies.remove('session_id');
        Cookies.remove('session_timestamp');
    },

    // Méthode pour refresh manuellement la session
    refreshSession: async (): Promise<AuthResponse> => {
        const sessionId = Cookies.get('session_id');
        if (!sessionId) {
            throw new Error('No session ID found');
        }

        const { data } = await apiClient.post<ApiResponse<AuthResponse>>('/auth/refresh', {}, {
            headers: {
                'Authorization': sessionId
            }
        });
        
        const { session_id } = data.data.session;
        const timestamp = Date.now();
        
        Cookies.set('session_id', session_id);
        Cookies.set('session_timestamp', timestamp.toString());
        
        return data.data;
    },
};

export const healthAPI = {
    getHealth: async (): Promise<HealthCheckResponse> => {
        const { data } = await apiClient.get<HealthCheckResponse>('/health');
        return data;
    },
};

export const usersAPI = {
    getUsers: async (filters: UserFilters = {}): Promise<PaginatedResponse<ListUsersResponse>> => {
        const { data } = await apiClient.get<PaginatedResponse<ListUsersResponse>>('/users', { params: filters });
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