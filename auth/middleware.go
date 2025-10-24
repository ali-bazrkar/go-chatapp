package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/aliBazrkar/go-chatapp/db"
)

type contextKey string

const UserContextKey contextKey = "user"

/*
the reason why we define a different middleware for
websocket in this case is that websocket unlike HTTP
methods do not need CSRF tokens and can not pass HTTP
reuest either, it is feasible to send CSRF token
through the /ws URL but it is not needed. we are using
CSRF tokens to make sure the request is coming from
same origin, in our case our "Upgrader" already is
set up to check the origin (defined in "auth" package)

so we use Middleware() for HTTP request that are
vulnerable to CSRF attacks and WebSocketMiddleware()
for other request activities.
*/

func (sm *SessionManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := sm.ValidateSession(w, r)
		if err != nil {
			log.Println("Session Validation Failed:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, token.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (sm *SessionManager) WebSocketMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := sm.CheckAuth(w, r)
		if err != nil {
			log.Println("Websocket Session Validation Failed:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, token.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) (*db.User, bool) {
	user, ok := r.Context().Value(UserContextKey).(*db.User)
	return user, ok

}
