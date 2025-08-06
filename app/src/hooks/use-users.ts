import { usersAPI } from '@/lib/api';
import { PaginatedResponse } from '@/types/api';
import { User, UserFilters } from '@/types/user';
import { useCallback, useEffect, useState } from 'react';

interface UseUsersReturn {
  users: User[];
  totalCount: number;
  isLoading: boolean;
  error: string | null;
  filters: UserFilters;
  setFilters: (filters: Partial<UserFilters>) => void;
  refetch: () => void;
}

export function useUsers(): UseUsersReturn {
  const [users, setUsers] = useState<User[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFiltersState] = useState<UserFilters>({
    page: 1,
    limit: 10,
    order_by: 'created_at',
    order: 'desc',
  });

  const fetchUsers = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response: PaginatedResponse<User> = await usersAPI.getUsers(filters);
      setUsers(response.data);
      setTotalCount(response.pagination.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch users');
    } finally {
      setIsLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const setFilters = useCallback((newFilters: Partial<UserFilters>) => {
    setFiltersState(prev => ({
      ...prev,
      ...newFilters,
      // Reset to first page when filters change
      page: newFilters.page || 1,
    }));
  }, []);

  const refetch = useCallback(() => {
    fetchUsers();
  }, [fetchUsers]);

  return {
    users,
    totalCount,
    isLoading,
    error,
    filters,
    setFilters,
    refetch,
  };
} 