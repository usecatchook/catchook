export enum UserRole {
    VIEWER = "viewer",
    DEVELOPER = "developer",
    ADMIN = "admin",
}

export interface User {
    id: number;
    email: string;
    role: UserRole;
    first_name: string;
    last_name?: string;
    full_name: string;
    is_active: boolean;
    created_at: string;
    updated_at?: string;
}


export interface UserFilters {
    page?: number;
    limit?: number;
    search?: string;
    role?: UserRole;
    is_active?: boolean;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
}

export interface CreateUserRequest {
    first_name: string;
    last_name: string;
    email: string;
    role?: User['role'];
}

export interface UpdateUserRequest {
    id: string;
    first_name: string;
    last_name: string;
    role?: User['role'];
    isActive?: boolean;
}