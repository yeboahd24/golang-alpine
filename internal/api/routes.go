package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/yeboahd24/authentication/internal/middleware"
)

func (s *Server) setupRoutes() {
	// Web routes
	s.router.Get("/", s.handleHome)
	s.router.Get("/login", s.handleLogin)
	s.router.Get("/register", s.handleRegister)
	s.router.Get("/dashboard", middleware.RequireAuth(s.handleDashboard))

	// API routes
	s.router.Route("/api", func(r chi.Router) {
		r.Post("/login", s.loginUser)
		r.Post("/register", s.registerUser)
		r.Get("/check-username", s.checkUsername)
		r.Get("/check-email", s.checkEmail)
		r.Post("/check-password", s.checkPasswordStrength)
	})
}
