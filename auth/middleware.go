package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/aliBazrkar/go-chatapp/db"
)

type contextKey string

const UserContextKey contextKey = "user"

func (sm *SessionManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, err := sm.ValidateSession(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) (*db.User, bool) {
	user, ok := r.Context().Value(UserContextKey).(*db.User)
	return user, ok

}

// In auth/session.go - add new method
func (sm *SessionManager) WebSocketMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session token from cookie
		cookie, err := r.Cookie(SessionCookieName)
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get CSRF from query param (WebSocket can't send headers)
		csrfToken := r.URL.Query().Get("csrf_token")
		if csrfToken == "" {
			http.Error(w, "Missing CSRF token", http.StatusUnauthorized)
			return
		}

		// Find token in database
		token, err := sm.db.GetTokenBySession(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate CSRF
		if token.CSRFToken != csrfToken {
			http.Error(w, "Invalid CSRF token", http.StatusUnauthorized)
			return
		}

		// Check expiration
		if time.Now().After(token.ExpiresAt) {
			sm.DeleteSession(cookie.Value)
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), UserContextKey, token.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
