import { UserFilters } from "@/types/user";

export const queryKeys = {
    users: {
        all: ['users'] as const,
        lists: () => [...queryKeys.users.all, 'list'] as const,
        list: (filters: UserFilters) => [...queryKeys.users.lists(), filters] as const,
        details: () => [...queryKeys.users.all, 'detail'] as const,
        detail: (id: number) => [...queryKeys.users.details(), id] as const,
    },
    setup: {
        all: ['setup'] as const,
        createAdminUser: () => [...queryKeys.setup.all, 'createAdminUser'] as const,
    },
} as const;