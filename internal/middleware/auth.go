package middleware

import (
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// TODO: Add JWT token validation here
		// token := cookie.Value
		// if !isValidToken(token) {
		//     http.Redirect(w, r, "/login", http.StatusSeeOther)
		//     return
		// }

		next(w, r)
	}
}
