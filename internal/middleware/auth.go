package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/theotruvelot/catchook/pkg/jwt"
	"github.com/theotruvelot/catchook/pkg/response"
)

func JWTAuth(jwtManager jwt.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Missing authorization header")
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Unauthorized(c, "Invalid authorization format")
		}

		token := authHeader[7:]
		if token == "" {
			return response.Unauthorized(c, "Missing token")
		}

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return response.Unauthorized(c, "Invalid token")
		}

		if claims.TokenType != "access" {
			return response.Unauthorized(c, "Invalid token type")
		}

		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("token", token)

		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) (int, bool) {
	userID := c.Locals("userID")
	if userID == nil {
		return 0, false
	}

	if id, ok := userID.(int); ok {
		return id, true
	}

	return 0, false
}

func GetToken(c *fiber.Ctx) (string, bool) {
	token := c.Locals("token")
	if token == nil {
		return "", false
	}

	if t, ok := token.(string); ok {
		return t, true
	}

	return "", false
}
