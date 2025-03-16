package api

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yeboahd24/authentication/internal/config"
	db "github.com/yeboahd24/authentication/internal/db/sqlc"
	"github.com/yeboahd24/authentication/internal/service"
)

type Server struct {
	router       *chi.Mux
	db           *db.Queries
	jwtMaker     *service.JWTMaker
	emailService *service.EmailService
	templates    *template.Template
	jwtConfig    config.JWTConfig
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func NewServer(db *db.Queries, jwtMaker *service.JWTMaker, emailService *service.EmailService) *Server {
	server := &Server{
		router:       chi.NewRouter(),
		db:           db,
		jwtMaker:     jwtMaker,
		emailService: emailService,
	}

	// Load templates
	templateFiles, err := filepath.Glob("web/templates/*.html")
	if err != nil {
		log.Fatalf("Failed to find templates: %v", err)
	}

	templates := template.New("").Funcs(template.FuncMap{
		// Add any custom template functions here
	})

	templates, err = templates.ParseFiles(templateFiles...)
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	server.templates = templates

	// Middleware
	server.router.Use(middleware.Logger)
	server.router.Use(middleware.Recoverer)

	// Setup routes
	server.setupRoutes()

	return server
}

// func (s *Server) setupRoutes() {
// 	// Web routes
// 	s.router.Get("/", s.handleHome)
// 	s.router.Get("/login", s.handleLogin)
// 	s.router.Get("/register", s.handleRegister)

// 	// API routes
// 	s.router.Route("/api", func(r chi.Router) {
// 		r.Post("/register", s.registerUser)
// 		r.Post("/login", s.loginUser)
// 		r.Get("/check-username", s.checkUsername)
// 		r.Get("/check-email", s.checkEmail)
// 		r.Post("/check-password", s.checkPasswordStrength)
// 	})
// }

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	err := s.templates.ExecuteTemplate(w, "layout.html", nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Login",
		"Content": "login", // This tells the layout which template to use
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Register",
		"Content": "register", // This tells the layout which template to use
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
