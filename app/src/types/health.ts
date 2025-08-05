export interface HealthCheckResponse {
    status: string;
    version: string;
    message: string;
    services: {
        database: string;
        redis: string;
    };
    is_first_time_setup: boolean;
}