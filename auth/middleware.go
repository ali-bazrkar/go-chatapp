package auth

import (
	"context"
	"net/http"

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
