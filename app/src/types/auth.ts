import { User } from "./user";

export interface AuthTokens {
    access_token: string;
    refresh_token: string;
}
  
export interface LoginCredentials {
    email: string;
    password: string;
}
  
export interface AuthResponse {
    user: User;
    tokens: AuthTokens;
}