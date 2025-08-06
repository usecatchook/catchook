"use client"

import { useAuth } from '@/contexts/auth-context';
import { useDebounce } from '@/hooks/use-debounce';
import { useUsers } from '@/hooks/use-users';
import { UserRole } from '@/types/user';
import {
  IconChevronDown,
  IconChevronUp,
  IconDotsVertical,
  IconEdit,
  IconSearch,
  IconTrash,
  IconUserCheck,
  IconUserX
} from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { usersAPI } from '@/lib/api';

const roleLabels: Record<UserRole, string> = {
  [UserRole.VIEWER]: 'Viewer',
  [UserRole.DEVELOPER]: 'Developer',
  [UserRole.ADMIN]: 'Admin',
};

const roleColors: Record<UserRole, string> = {
  [UserRole.VIEWER]: 'bg-blue-50 text-blue-700',
  [UserRole.DEVELOPER]: 'bg-green-50 text-green-700',
  [UserRole.ADMIN]: 'bg-purple-50 text-purple-700',
};

export function UsersTable() {
  const { user: currentUser } = useAuth();
  const [searchTerm, setSearchTerm] = useState('');
  const [roleFilter, setRoleFilter] = useState<UserRole | ''>('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'inactive'>('all');
  const [sortBy, setSortBy] = useState('created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  const debouncedSearch = useDebounce(searchTerm, 300);

  const { users, totalCount, isLoading, error, filters, setFilters, refetch } = useUsers();
  const [isInitialLoad, setIsInitialLoad] = useState(true);

  // Update filters when search or other filters change
  useEffect(() => {
    setFilters({
      search: debouncedSearch || undefined,
      role: roleFilter || undefined,
      is_active: statusFilter === 'all' ? undefined : statusFilter === 'active',
      order_by: sortBy,
      order: sortOrder,
    });
  }, [debouncedSearch, roleFilter, statusFilter, sortBy, sortOrder, setFilters]);

  // Track initial load vs filter changes
  useEffect(() => {
    if (!isLoading && isInitialLoad) {
      setIsInitialLoad(false);
    }
  }, [isLoading, isInitialLoad]);

  const handleDeleteUser = async (userId: number) => {
    if (!confirm('Are you sure you want to delete this user?')) return;
    
    try {
      await usersAPI.deleteUser(userId);
      toast.success('User deleted successfully');
      refetch();
    } catch (error) {
      toast.error('Failed to delete user');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  if (error) {
    console.error(error);
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <p className="text-red-600 mb-2">Error loading users</p>
          <Button onClick={refetch} variant="outline">
            Try Again
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="relative flex-1">
          <IconSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
          <Input
            placeholder="Search users..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-10"
          />
        </div>
        
        <Select value={roleFilter || "all"} onValueChange={(value) => setRoleFilter(value === "all" ? "" : value as UserRole)}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="All roles" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All roles</SelectItem>
            {Object.entries(roleLabels).map(([role, label]) => (
              <SelectItem key={role} value={role}>
                {label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Select value={statusFilter} onValueChange={(value) => setStatusFilter(value as 'all' | 'active' | 'inactive')}>
          <SelectTrigger className="w-[140px]">
            <SelectValue placeholder="All status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All status</SelectItem>
            <SelectItem value="active">Active</SelectItem>
            <SelectItem value="inactive">Inactive</SelectItem>
          </SelectContent>
        </Select>


      </div>

      {/* Table */}
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setSortBy('first_name');
                    setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
                  }}
                  className="h-auto p-0 font-medium"
                >
                  Name
                  {sortBy === 'first_name' && (
                    sortOrder === 'asc' ? <IconChevronUp className="ml-1 h-4 w-4" /> : <IconChevronDown className="ml-1 h-4 w-4" />
                  )}
                </Button>
              </TableHead>
              <TableHead>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setSortBy('email');
                    setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
                  }}
                  className="h-auto p-0 font-medium"
                >
                  Email
                  {sortBy === 'email' && (
                    sortOrder === 'asc' ? <IconChevronUp className="ml-1 h-4 w-4" /> : <IconChevronDown className="ml-1 h-4 w-4" />
                  )}
                </Button>
              </TableHead>
              <TableHead>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setSortBy('role');
                    setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
                  }}
                  className="h-auto p-0 font-medium"
                >
                  Role
                  {sortBy === 'role' && (
                    sortOrder === 'asc' ? <IconChevronUp className="ml-1 h-4 w-4" /> : <IconChevronDown className="ml-1 h-4 w-4" />
                  )}
                </Button>
              </TableHead>
              <TableHead>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setSortBy('is_active');
                    setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
                  }}
                  className="h-auto p-0 font-medium"
                >
                  Status
                  {sortBy === 'is_active' && (
                    sortOrder === 'asc' ? <IconChevronUp className="ml-1 h-4 w-4" /> : <IconChevronDown className="ml-1 h-4 w-4" />
                  )}
                </Button>
              </TableHead>
              <TableHead>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setSortBy('created_at');
                    setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
                  }}
                  className="h-auto p-0 font-medium"
                >
                  Created
                  {sortBy === 'created_at' && (
                    sortOrder === 'asc' ? <IconChevronUp className="ml-1 h-4 w-4" /> : <IconChevronDown className="ml-1 h-4 w-4" />
                  )}
                </Button>
              </TableHead>
              <TableHead className="w-[50px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading && isInitialLoad ? (
              // Loading skeletons for initial load
              Array.from({ length: 5 }).map((_, i) => (
                                                  <TableRow key={i}>
                   <TableCell>
                     <Skeleton className="h-4 w-[150px]" />
                   </TableCell>
                   <TableCell>
                     <Skeleton className="h-4 w-[200px]" />
                   </TableCell>
                   <TableCell>
                     <Skeleton className="h-6 w-[80px]" />
                   </TableCell>
                   <TableCell>
                     <Skeleton className="h-6 w-[60px]" />
                   </TableCell>
                   <TableCell>
                     <Skeleton className="h-4 w-[100px]" />
                   </TableCell>
                   <TableCell>
                     <Skeleton className="h-8 w-8" />
                   </TableCell>
                 </TableRow>
              ))
                                        ) : isLoading ? (
                 // Show existing data with loading overlay for filter changes
                 users.map((user) => (
                   <TableRow key={user.id} className="opacity-50">
                     <TableCell>
                       <div>
                         <p className="font-medium">{user.first_name} {user.last_name}</p>
                       </div>
                     </TableCell>
                     <TableCell>
                       <p className="text-sm text-muted-foreground">{user.email}</p>
                     </TableCell>
                     <TableCell>
                       <Badge className={roleColors[user.role]}>
                         {roleLabels[user.role]}
                       </Badge>
                     </TableCell>
                   <TableCell>
                     <div className="flex items-center space-x-2">
                       {user.is_active ? (
                         <IconUserCheck className="h-4 w-4 text-green-600" />
                       ) : (
                         <IconUserX className="h-4 w-4 text-red-600" />
                       )}
                       <Badge variant={user.is_active ? 'default' : 'secondary'}>
                         {user.is_active ? 'Active' : 'Inactive'}
                       </Badge>
                     </div>
                   </TableCell>
                   <TableCell>
                     <span className="text-sm text-muted-foreground">
                       {formatDate(user.created_at)}
                     </span>
                   </TableCell>
                   <TableCell>
                     <DropdownMenu>
                       <DropdownMenuTrigger asChild>
                         <Button variant="ghost" size="sm">
                           <IconDotsVertical className="h-4 w-4" />
                         </Button>
                       </DropdownMenuTrigger>
                       <DropdownMenuContent align="end">
                         <DropdownMenuItem>
                           <IconEdit className="h-4 w-4 mr-2" />
                           Edit
                         </DropdownMenuItem>
                         {currentUser?.id !== user.id && (
                           <DropdownMenuItem
                             onClick={() => handleDeleteUser(user.id)}
                             className="text-red-600"
                           >
                             <IconTrash className="h-4 w-4 mr-2" />
                             Delete
                           </DropdownMenuItem>
                         )}
                       </DropdownMenuContent>
                     </DropdownMenu>
                   </TableCell>
                 </TableRow>
               ))
                          ) : users.length === 0 ? (
               <TableRow>
                 <TableCell colSpan={6} className="text-center py-8">
                   <div className="flex flex-col items-center space-y-2">
                     <p className="text-muted-foreground">No users found</p>
                   </div>
                 </TableCell>
               </TableRow>
             ) : (
               users.map((user) => (
                 <TableRow key={user.id}>
                   <TableCell>
                     <div>
                       <p className="font-medium">{user.first_name} {user.last_name}</p>
                     </div>
                   </TableCell>
                   <TableCell>
                     <p className="text-sm text-muted-foreground">{user.email}</p>
                   </TableCell>
                   <TableCell>
                     <Badge className={roleColors[user.role]}>
                       {roleLabels[user.role]}
                     </Badge>
                   </TableCell>
                  <TableCell>
                    <div className="flex items-center space-x-2">
                      {user.is_active ? (
                        <IconUserCheck className="h-4 w-4 text-green-600" />
                      ) : (
                        <IconUserX className="h-4 w-4 text-red-600" />
                      )}
                      <Badge variant={user.is_active ? 'default' : 'secondary'}>
                        {user.is_active ? 'Active' : 'Inactive'}
                      </Badge>
                    </div>
                  </TableCell>
                  <TableCell>
                    <span className="text-sm text-muted-foreground">
                      {formatDate(user.created_at)}
                    </span>
                  </TableCell>
                  <TableCell>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="sm">
                          <IconDotsVertical className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem>
                          <IconEdit className="h-4 w-4 mr-2" />
                          Edit
                        </DropdownMenuItem>
                        {currentUser?.id !== user.id && (
                          <DropdownMenuItem
                            onClick={() => handleDeleteUser(user.id)}
                            className="text-red-600"
                          >
                            <IconTrash className="h-4 w-4 mr-2" />
                            Delete
                          </DropdownMenuItem>
                        )}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Pagination */}
      {!isLoading && users.length > 0 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            Showing {users.length} of {totalCount} users
          </p>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              disabled={filters.page === 1}
              onClick={() => setFilters({ page: (filters.page || 1) - 1 })}
            >
              Previous
            </Button>
            <span className="text-sm">
              Page {filters.page || 1}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={users.length < (filters.limit || 10)}
              onClick={() => setFilters({ page: (filters.page || 1) + 1 })}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
} 