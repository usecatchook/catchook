package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/platform/storage/postgres/generated"

	"github.com/theotruvelot/catchook/internal/platform/session"
	"github.com/theotruvelot/catchook/pkg/response"
)

const UserContextKey = "user"

type User struct {
	ID   string             `json:"id"`
	Role generated.UserRole `json:"role"`
}

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

		user := &User{
			ID:   session.UserID,
			Role: session.Role,
		}

		c.Locals(UserContextKey, user)
		return c.Next()
	}
}

func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUser(c)
		if user == nil {
			return response.Unauthorized(c, "Authentication required")
		}

		for _, role := range roles {
			if user.Role == generated.UserRole(role) {
				return c.Next()
			}
		}

		return response.Forbidden(c, "Insufficient permissions")
	}
}

func GetUser(c *fiber.Ctx) *User {
	user := c.Locals(UserContextKey)
	if user == nil {
		return nil
	}

	if u, ok := user.(*User); ok {
		return u
	}

	return nil
}

func GetUserID(c *fiber.Ctx) (string, bool) {
	user := GetUser(c)
	if user == nil {
		return "", false
	}
	return user.ID, true
}
