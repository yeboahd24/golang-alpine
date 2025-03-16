package api

import (
	"encoding/json"
	"log"
	"net/http"

	db "github.com/yeboahd24/authentication/internal/db/sqlc"
	"github.com/yeboahd24/authentication/internal/service"
)

const (
	WeakPassword       = 0
	MediumPassword     = 1
	StrongPassword     = 2
	VeryStrongPassword = 3
)

func calculatePasswordStrength(password string) int {
	// Initialize score
	score := 0

	// Check length
	if len(password) >= 8 {
		score++
	}

	// Check for numbers
	hasNumber := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasNumber = true
			break
		}
	}
	if hasNumber {
		score++
	}

	// Check for special characters
	hasSpecial := false
	for _, char := range password {
		if (char >= '!' && char <= '/') || (char >= ':' && char <= '@') || (char >= '[' && char <= '`') || (char >= '{' && char <= '~') {
			hasSpecial = true
			break
		}
	}
	if hasSpecial {
		score++
	}

	// Check for mixed case
	hasUpper := false
	hasLower := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		}
		if char >= 'a' && char <= 'z' {
			hasLower = true
		}
	}
	if hasUpper && hasLower {
		score++
	}

	return score
}

func getPasswordFeedback(strength int) string {
	switch strength {
	case WeakPassword:
		return "Password is too weak. Use at least 8 characters with numbers, special characters, and mixed case."
	case MediumPassword:
		return "Password could be stronger. Try adding special characters or mixed case."
	case StrongPassword:
		return "Good password strength."
	case VeryStrongPassword:
		return "Excellent password strength!"
	default:
		return "Invalid password strength"
	}
}

func (s *Server) checkUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	exists, err := s.db.CheckUsernameExists(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if exists {
		w.Write([]byte(`<p class="mt-1 text-sm text-red-600">Username already taken</p>`))
	}
}

func (s *Server) checkEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	exists, err := s.db.CheckEmailExists(r.Context(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if exists {
		w.Write([]byte(`<p class="mt-1 text-sm text-red-600">Email already registered</p>`))
	}
}

func (s *Server) checkPasswordStrength(w http.ResponseWriter, r *http.Request) {
	var password struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Implement password strength checking logic
	strength := calculatePasswordStrength(password.Password)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"strength": strength,
		"feedback": getPasswordFeedback(strength),
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

func (s *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from database by email
	user, err := s.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	passwordConfig := service.NewPasswordConfig()
	valid, err := passwordConfig.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil || !valid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := s.jwtMaker.CreateToken(user.ID.String(), user.Username, s.jwtConfig.TokenDuration)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Set token as HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // Set to true if using HTTPS
		MaxAge:   int(s.jwtConfig.TokenDuration.Seconds()),
		SameSite: http.SameSiteLaxMode,
	})

	// Create response
	response := LoginResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Token:    token, // You might want to remove this from the response since we're using cookies
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) registerUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Username == "" || req.Password == "" {
		http.Error(w, "Email, username and password are required", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	emailExists, err := s.db.CheckEmailExists(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if emailExists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Check if username already exists
	usernameExists, err := s.db.CheckUsernameExists(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if usernameExists {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	// Check password strength
	strength := calculatePasswordStrength(req.Password)
	if strength < MediumPassword {
		http.Error(w, "Password too weak", http.StatusBadRequest)
		return
	}

	// Hash password
	passwordConfig := service.NewPasswordConfig()
	hashedPassword, err := passwordConfig.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create user
	user, err := s.db.CreateUser(r.Context(), db.CreateUserParams{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Dashboard",
		"Content": "dashboard", // This tells the layout which content template to use
	}

	err := s.templates.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
// 	data := map[string]interface{}{
// 		"Title":   "Register",
// 		"Content": "register",
// 	}
// 	if err := s.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
// 		log.Printf("Template execution error: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }

// func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
// 	data := map[string]interface{}{
// 		"Title":   "Home",
// 		"Content": "home",
// 	}
// 	if err := s.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
// 		log.Printf("Template execution error: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }
