import { healthAPI } from '@/lib/api';
import { HealthCheckResponse } from '@/types/health';
import { useQuery } from '@tanstack/react-query';
import { usePathname, useRouter } from 'next/navigation';
import { useEffect } from 'react';



export function useHealth() {
  const router = useRouter();
  const pathname = usePathname();
  
  const { data: health, isLoading, error } = useQuery({
    queryKey: ['health'],
    queryFn: async (): Promise<HealthCheckResponse> => {
      return await healthAPI.getHealth();
    },
    retry: 2,
    refetchOnWindowFocus: false,
    retryDelay: 10,
    staleTime: 0,
    gcTime: 0,
    refetchOnMount: 'always',
  });

  useEffect(() => {
    if (health && health.status !== 'unknown' && health.is_first_time_setup && pathname !== '/setup') {
      router.replace('/setup');
    }
  }, [health, pathname, router]);

  return { health, isLoading, error };
}