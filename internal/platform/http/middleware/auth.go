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

		// Créer l'utilisateur et l'injecter dans le contexte
		authUser := &auth.AuthUser{
			ID:   session.UserID,
			Role: session.Role,
		}

		ctx := auth.WithUser(c.Context(), authUser)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

// Helper pour récupérer l'utilisateur depuis le contexte Fiber
func GetAuthUser(c *fiber.Ctx) (*auth.AuthUser, error) {
	return auth.GetUser(c.Context())
}

func GetAuthUserID(c *fiber.Ctx) (string, error) {
	return auth.GetUserID(c.Context())
}

// RequirePermission middleware qui vérifie une permission spécifique
func RequirePermission(permission auth.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := auth.RequirePermission(c.Context(), permission); err != nil {
			return response.Forbidden(c, "Insufficient permissions")
		}
		return c.Next()
	}
}

// RequireAdmin middleware qui requiert les droits admin
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := auth.RequireAdmin(c.Context()); err != nil {
			return response.Forbidden(c, "Admin permissions required")
		}
		return c.Next()
	}
}

// RequireOwnership middleware qui vérifie la propriété d'une resource
// Le paramètre paramName indique dans quel paramètre URL trouver l'ID de la resource
func RequireOwnership(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		resourceID := c.Params(paramName)
		if resourceID == "" {
			return response.BadRequest(c, "Resource ID is required", nil)
		}

		user, err := auth.GetUser(c.Context())
		if err != nil {
			return response.Unauthorized(c, "Authentication required")
		}

		// Pour l'instant, on assume que l'ID de la resource est l'ID du propriétaire
		// Dans un vrai système, il faudrait faire une query pour récupérer le propriétaire
		if !user.CanManageResource(resourceID) {
			return response.Forbidden(c, "You can only manage your own resources")
		}

		return c.Next()
	}
}

// RequireOwnershipOrAdmin combine ownership + droits admin
func RequireOwnershipOrAdmin(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := auth.GetUser(c.Context())
		if err != nil {
			return response.Unauthorized(c, "Authentication required")
		}

		// Les admins passent toujours
		if user.IsAdmin() {
			return c.Next()
		}

		// Sinon vérifier la propriété
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
