package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
)

// setupRoutes configures all application routes
func (s *Server) setupRoutes() {
	// API routes
	api := s.app.Group("/api/v1")

	// Health check
	api.Get("/health", s.handleHealthCheck)

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

	// 404 handler
	s.app.Use(func(c *fiber.Ctx) error {
		return response.NotFound(c, "Route not found")
	})
}

// setupAuthRoutes configures authentication routes
func (s *Server) setupAuthRoutes(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/login", s.handleLogin)
	auth.Post("/refresh", s.handleRefreshSession)
	auth.Post("/logout", s.handleLogout)
}

// setupUserRoutes configures user management routes
func (s *Server) setupUserRoutes(api fiber.Router) {
	users := api.Group("/users")

	users.Use(middleware.SessionAuth(s.container.Session))

	users.Get("/me", s.handleGetMe)
	users.Get("/profile/:id", s.handleGetProfile)
	users.Put("/profile/:id", s.handleUpdateProfile)

}

func (s *Server) setupAdminRoutes(api fiber.Router) {
	admin := api.Group("/admin")

	admin.Use(middleware.SessionAuth(s.container.Session))
	admin.Use(middleware.RequireRoles("admin"))

	//users management routes
	admin.Get("/users", s.handleListUsers)
	admin.Post("/users", s.handleCreateUser)
	//admin.Put("/users/:id", s.handleUpdateUser) //TODO: Implement user update
	//admin.Delete("/users/:id", s.handleDeleteUser) //TODO: Implement user delete
}

func (s *Server) setupSetupRoutes(api fiber.Router) {
	setup := api.Group("/setup")

	setup.Post("/", s.handleSetup)
}

func (s *Server) setupSourceRoutes(api fiber.Router) {
	sources := api.Group("/sources")

	sources.Use(middleware.SessionAuth(s.container.Session))

	sources.Post("/", s.handleCreateSource)
	sources.Get(":id", s.handleGetSource)
	//TODO UDL
}
