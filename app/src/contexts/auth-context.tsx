"use client"

import { useValidationErrors, ValidationErrors } from '@/hooks/useValidationErrors';
import { authAPI } from '@/lib/api';
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
  
  // Vérification immédiate et synchrone des tokens
  const authToken = Cookies.get('authToken');
  const refreshToken = Cookies.get('refreshToken');
  const hasTokens = !!(authToken && refreshToken);
  const isPublicRoute = PUBLIC_ROUTES.includes(pathname);

  // Redirection immédiate basée sur les tokens seulement
  useEffect(() => {
    // Si connecté (tokens présents) et sur page publique → redirect vers home
    if (hasTokens && isPublicRoute) {
      router.replace('/');
      return;
    }

    // Si pas connecté (pas de tokens) et sur page protégée → redirect vers login  
    if (!hasTokens && !isPublicRoute) {
      router.replace('/login');
      return;
    }

    // Si on arrive ici, pas de redirection nécessaire
    setIsReady(true);
  }, [hasTokens, isPublicRoute, router]);

  // Récupérer les données utilisateur si tokens présents
  useEffect(() => {
    if (!hasTokens || !isReady) return;

    setIsLoading(true);
    authAPI.getCurrentUser()
      .then(setUser)
      .catch(() => {
        // Si l'API échoue, tokens probablement expirés
        Cookies.remove('authToken');
        Cookies.remove('refreshToken');
        router.replace('/login');
      })
      .finally(() => setIsLoading(false));
  }, [hasTokens, isReady, router]);

  const login = async (credentials: LoginCredentials): Promise<void> => {
    setIsLoggingIn(true);
    setLoginError(null);
    clearErrors();
    
    try {
      const data = await authAPI.login(credentials);
      Cookies.set('authToken', data.tokens.access_token, { expires: 7 });
      Cookies.set('refreshToken', data.tokens.refresh_token, { expires: 30 });
      setUser(data.user);
      router.replace('/');
    } catch (error) {
      // Essayer d'extraire les erreurs de validation
      const hasValidationErrors = setErrorsFromException(error);
      
      // Si pas d'erreurs de validation spécifiques, afficher l'erreur générale
      if (!hasValidationErrors) {
        setLoginError(error instanceof Error ? error.message : 'Erreur de connexion');
      }
      
      throw error;
    } finally {
      setIsLoggingIn(false);
    }
  };

  const logout = () => {
    Cookies.remove('authToken');
    Cookies.remove('refreshToken');
    setUser(undefined);
    router.replace('/login');
  };

  // Ne pas rendre si redirection en cours
  if (!isReady) return null;

  const value: AuthContextType = {
    user,
    isAuthenticated: hasTokens,
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