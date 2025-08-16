package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/platform/auth"
	"github.com/theotruvelot/catchook/internal/platform/http/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
)

// setupRoutes configures all application routes
func (s *Server) setupRoutes() {
	// API routes
	api := s.app.Group("/api/v1")

	// Health check
	api.Get("/health", s.healthHandler.HealthCheck)

	// Auth routes
	s.setupAuthRoutes(api)

	// User routes
	s.setupUserRoutes(api)

	// Setup routes
	s.setupSetupRoutes(api)

	// Admin routes
	s.setupAdminRoutes(api)

	// Source routes
	s.setupSourceRoutes(api)

	// Destination routes
	s.setupDestinationRoutes(api)

	// 404 handler
	s.app.Use(func(c *fiber.Ctx) error {
		return response.NotFound(c, "Route not found")
	})
}

// setupAuthRoutes configures authentication routes
func (s *Server) setupAuthRoutes(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/login", s.authHandler.Login)
	auth.Post("/refresh", s.authHandler.RefreshSession)
	auth.Post("/logout", s.authHandler.Logout)
}

// setupUserRoutes configures user management routes
func (s *Server) setupUserRoutes(api fiber.Router) {
	users := api.Group("/users")

	users.Use(middleware.SessionAuth(s.container.Session))

	users.Get("/me", s.userHandler.GetMe)
	users.Get("/profile/:id", s.userHandler.GetProfile)
	users.Put("/profile/:id", s.userHandler.UpdateProfile)

}

func (s *Server) setupAdminRoutes(api fiber.Router) {
	admin := api.Group("/admin")

	admin.Use(middleware.SessionAuth(s.container.Session))
	admin.Use(middleware.RequireAdmin()) // Nouveau système plus élégant

	//users management routes
	admin.Get("/users", s.userHandler.ListUsers)
	admin.Post("/users", s.userHandler.CreateUser)
	//admin.Put("/users/:id", s.handleUpdateUser) //TODO: Implement user update
	//admin.Delete("/users/:id", s.handleDeleteUser) //TODO: Implement user delete
}

func (s *Server) setupSetupRoutes(api fiber.Router) {
	setup := api.Group("/setup")

	setup.Post("/", s.setupHandler.Setup)
}

func (s *Server) setupSourceRoutes(api fiber.Router) {
	sources := api.Group("/sources")
	sources.Use(middleware.SessionAuth(s.container.Session))

	sources.Post("/", middleware.RequirePermission(auth.PermissionWrite), s.sourceHandler.CreateSource)
	sources.Get("/:id", s.sourceHandler.GetSource)
	sources.Get("/", s.sourceHandler.ListSources)
	sources.Put("/:id", middleware.RequireOwnershipOrAdmin("id"), s.sourceHandler.UpdateSource)
	sources.Delete("/:id", middleware.RequireOwnershipOrAdmin("id"), s.sourceHandler.DeleteSource)
}

func (s *Server) setupDestinationRoutes(api fiber.Router) {
	destinations := api.Group("/destinations")
	destinations.Use(middleware.SessionAuth(s.container.Session))

	destinations.Post("/", middleware.RequirePermission(auth.PermissionWrite), s.destinationHandler.CreateDestination)
	destinations.Get("/:id", s.destinationHandler.GetDestination)
	destinations.Get("/", s.destinationHandler.ListDestinations)
	destinations.Put("/:id", middleware.RequireOwnershipOrAdmin("id"), s.destinationHandler.UpdateDestination)
	destinations.Delete("/:id", middleware.RequireOwnershipOrAdmin("id"), s.destinationHandler.DeleteDestination)
}
