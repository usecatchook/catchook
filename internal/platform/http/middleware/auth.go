package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/platform/auth"
	"github.com/theotruvelot/catchook/internal/platform/session"
	"github.com/theotruvelot/catchook/pkg/response"
)

func SessionAuth(sessionManager session.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Get("Authorization")
		if sessionID == "" {
			return response.Unauthorized(c, "Missing authorization header")
		}

		session, err := sessionManager.ValidateSession(c.Context(), sessionID)
		if err != nil {
			return response.Unauthorized(c, "Invalid or expired session")
		}

		authUser := &auth.AuthUser{
			ID:   session.UserID,
			Role: session.Role,
		}

		c.Locals("auth_user", authUser)

		return c.Next()
	}
}

func GetAuthUser(c *fiber.Ctx) (*auth.AuthUser, error) {
	user := c.Locals("auth_user")
	if user == nil {
		return nil, auth.ErrUserNotInContext
	}

	authUser, ok := user.(*auth.AuthUser)
	if !ok {
		return nil, auth.ErrInvalidUserType
	}

	return authUser, nil
}

func GetAuthUserID(c *fiber.Ctx) (string, error) {
	user, err := GetAuthUser(c)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func RequirePermission(permission auth.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return response.Unauthorized(c, "Authentication required")
		}

		if !user.HasPermission(permission) {
			return response.Forbidden(c, "Insufficient permissions")
		}

		return c.Next()
	}
}

func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return response.Unauthorized(c, "Authentication required")
		}

		if !user.IsAdmin() {
			return response.Forbidden(c, "Admin permissions required")
		}

		return c.Next()
	}
}

func RequireOwnership(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		resourceID := c.Params(paramName)
		if resourceID == "" {
			return response.BadRequest(c, "Resource ID is required", nil)
		}

		user, err := GetAuthUser(c)
		if err != nil {
			return response.Unauthorized(c, "Authentication required")
		}

		if !user.CanManageResource(resourceID) {
			return response.Forbidden(c, "You can only manage your own resources")
		}

		return c.Next()
	}
}

func RequireOwnershipOrAdmin(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return response.Unauthorized(c, "Authentication required")
		}

		if user.IsAdmin() {
			return c.Next()
		}

		resourceID := c.Params(paramName)
		if resourceID == "" {
			return response.BadRequest(c, "Resource ID is required", nil)
		}

		if !user.CanManageResource(resourceID) {
			return response.Forbidden(c, "Insufficient permissions")
		}

		return c.Next()
	}
}
