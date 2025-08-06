import { User } from "./user";

export interface AuthSession {
    session_id: string;
}
  
export interface LoginCredentials {
    email: string;
    password: string;
}
  
export interface AuthResponse {
    user: User;
    session: AuthSession;
}