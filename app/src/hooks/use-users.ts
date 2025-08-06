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
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFiltersState] = useState<UserFilters>({
    page: 1,
    limit: 10,
    order_by: 'created_at',
    order: 'desc',
  });
  const [hasInitialized, setHasInitialized] = useState(false);

  const fetchUsers = useCallback(async (currentFilters: UserFilters) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response: PaginatedResponse<User> = await usersAPI.getUsers(currentFilters);
      setUsers(response.data);
      setTotalCount(response.pagination.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch users');
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Initial fetch
  useEffect(() => {
    if (!hasInitialized) {
      fetchUsers(filters);
      setHasInitialized(true);
    }
  }, [hasInitialized, fetchUsers, filters]);

  const setFilters = useCallback((newFilters: Partial<UserFilters>) => {
    setFiltersState(prev => {
      const updatedFilters = {
        ...prev,
        ...newFilters,
        // Reset to first page when filters change
        page: newFilters.page || 1,
      };
      // Only fetch if we've already initialized
      if (hasInitialized) {
        fetchUsers(updatedFilters);
      }
      return updatedFilters;
    });
  }, [fetchUsers, hasInitialized]);

  const refetch = useCallback(() => {
    fetchUsers(filters);
  }, [fetchUsers, filters]);

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