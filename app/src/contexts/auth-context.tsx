"use client"

import { useValidationErrors, ValidationErrors } from '@/hooks/use-validation-errors';
import { authAPI } from '@/lib/api';
import type { ApiError } from '@/types/api';
import type { LoginCredentials } from '@/types/auth';
import { User } from "@/types/user";
import Cookies from 'js-cookie';
import { usePathname, useRouter } from 'next/navigation';
import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';

interface AuthContextType {
  user: User | undefined;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => void;
  isLoggingIn: boolean;
  loginError: string | null;
  validationErrors: ValidationErrors;
  getFieldError: (fieldName: string) => string | undefined;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

const PUBLIC_ROUTES = ['/login', '/setup'];

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const router = useRouter();
  const pathname = usePathname();
  
  const [user, setUser] = useState<User | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(false);
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const [loginError, setLoginError] = useState<string | null>(null);
  const [isReady, setIsReady] = useState(false);
  
  const { validationErrors, setErrorsFromException, getFieldError, clearErrors } = useValidationErrors();
  
  // Vérification de la session
  const sessionId = Cookies.get('session_id');
  const hasSession = !!sessionId;
  const isPublicRoute = PUBLIC_ROUTES.includes(pathname);

  // Redirection basée sur la session
  useEffect(() => {
    // Si connecté (session présente) et sur page publique → redirect vers home
    if (hasSession && isPublicRoute) {
      router.replace('/');
      return;
    }

    // Si pas connecté (pas de session) et sur page protégée → redirect vers login  
    if (!hasSession && !isPublicRoute) {
      router.replace('/login');
      return;
    }

    // Si on arrive ici, pas de redirection nécessaire
    setIsReady(true);
  }, [hasSession, isPublicRoute, router]);

  // Récupérer les données utilisateur si session présente
  useEffect(() => {
    if (!hasSession || !isReady) return;

    setIsLoading(true);
    authAPI.getCurrentUser()
      .then(setUser)
      .catch(() => {
        // Si l'API échoue, session probablement expirée
        Cookies.remove('session_id');
        Cookies.remove('session_timestamp');
        router.replace('/login');
      })
      .finally(() => setIsLoading(false));
  }, [hasSession, isReady, router]);

  const login = async (credentials: LoginCredentials): Promise<void> => {
    setIsLoggingIn(true);
    setLoginError(null);
    clearErrors();
    
    try {
      const data = await authAPI.login(credentials);
      setUser(data.user);
      router.replace('/');
    } catch (error) {
      // Essayer d'extraire les erreurs de validation
      const hasValidationErrors =
        error && typeof error === 'object' && 'message' in error
          ? setErrorsFromException(error as ApiError)
          : false;
      // Si pas d'erreurs de validation spécifiques, afficher l'erreur générale
      if (!hasValidationErrors) {
        setLoginError(error instanceof Error ? error.message : 'Erreur de connexion');
      }
      throw error;
    } finally {
      setIsLoggingIn(false);
    }
  };

  const logout = async () => {
    try {
      await authAPI.logout();
    } catch (error) {
      // Même si le logout API échoue, on nettoie côté client
      console.error('Logout API failed:', error);
    } finally {
      setUser(undefined);
      router.replace('/login');
    }
  };

  // Ne pas rendre si redirection en cours
  if (!isReady) return null;

  const value: AuthContextType = {
    user,
    isAuthenticated: hasSession,
    isLoading,
    login,
    logout,
    isLoggingIn,
    loginError,
    validationErrors,
    getFieldError,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}