package auth

import (
	"context"
	"errors"

	"github.com/theotruvelot/catchook/internal/platform/storage/postgres/generated"
)

// Clés pour le contexte
type contextKey string

const (
	userContextKey contextKey = "auth:user"
)

var (
	ErrUserNotInContext = errors.New("user not found in context")
	ErrInvalidUserType  = errors.New("invalid user type in context")
)

type AuthUser struct {
	ID   string             `json:"id"`
	Role generated.UserRole `json:"role"`
}

type Permission string

const (
	PermissionRead   Permission = "read"
	PermissionWrite  Permission = "write"
	PermissionDelete Permission = "delete"

	PermissionManageUsers   Permission = "manage:users"
	PermissionManageSources Permission = "manage:sources"
	PermissionViewAnalytics Permission = "view:analytics"

	PermissionAdminAll Permission = "admin:*"
)

var RolePermissions = map[generated.UserRole][]Permission{
	generated.UserRoleAdmin: {
		PermissionAdminAll,
		PermissionRead,
		PermissionWrite,
		PermissionDelete,
		PermissionManageUsers,
		PermissionManageSources,
		PermissionViewAnalytics,
	},
	generated.UserRoleDeveloper: {
		PermissionRead,
		PermissionWrite,
		PermissionDelete,
		PermissionManageSources,
		PermissionViewAnalytics,
	},
	generated.UserRoleViewer: {
		PermissionRead,
	},
}

func (u *AuthUser) HasPermission(permission Permission) bool {
	permissions, exists := RolePermissions[u.Role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == PermissionAdminAll || p == permission {
			return true
		}
	}

	return false
}

func (u *AuthUser) IsAdmin() bool {
	return u.Role == generated.UserRoleAdmin
}

func (u *AuthUser) CanManageUser(targetUserID string) bool {
	if u.IsAdmin() {
		return true
	}

	return u.ID == targetUserID
}

func (u *AuthUser) CanManageResource(resourceOwnerID string) bool {
	if u.IsAdmin() {
		return true
	}
	return u.ID == resourceOwnerID
}

func WithUser(ctx context.Context, user *AuthUser) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUser(ctx context.Context) (*AuthUser, error) {
	user := ctx.Value(userContextKey)
	if user == nil {
		return nil, ErrUserNotInContext
	}

	authUser, ok := user.(*AuthUser)
	if !ok {
		return nil, ErrInvalidUserType
	}

	return authUser, nil
}

func GetUserID(ctx context.Context) (string, error) {
	user, err := GetUser(ctx)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func RequirePermission(ctx context.Context, permission Permission) error {
	user, err := GetUser(ctx)
	if err != nil {
		return err
	}

	if !user.HasPermission(permission) {
		return errors.New("insufficient permissions")
	}

	return nil
}

func RequireAdmin(ctx context.Context) error {
	user, err := GetUser(ctx)
	if err != nil {
		return err
	}

	if !user.IsAdmin() {
		return errors.New("admin permissions required")
	}

	return nil
}

func RequireOwnershipOrAdmin(ctx context.Context, resourceOwnerID string) error {
	user, err := GetUser(ctx)
	if err != nil {
		return err
	}

	if !user.CanManageResource(resourceOwnerID) {
		return errors.New("insufficient permissions: must be owner or admin")
	}

	return nil
}
