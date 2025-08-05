export interface ApiResponse<TData = unknown> {
    data: TData;
    message?: string;
    error?: {
        code?: string;
        message?: string;
    };
    success: boolean;
}

export interface ValidationErrorResponse {
    success: false;
    message: string;
    errors: Record<string, string>;
    timestamp: string;
}

export interface PaginatedResponse<TData = unknown> {
    data: TData[];
    pagination: {
        currentPage: number;
        totalPages: number;
        total: number;
        limit: number;
        hasNext: boolean;
        hasPrev: boolean;
    };
    success: boolean;
}