package handlers

import (
	"log"
	"net/http"

	"github.com/aliBazrkar/go-chatapp/auth"
)

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(auth.SessionCookieName)
	if err == nil && cookie.Value != "" {
		if err := h.sm.DeleteSession(cookie.Value); err != nil {
			log.Printf("Session Token Deletage Error: %v", err)
		}
	}

	auth.ClearSessionCookie(w)

	respondJSON(w, map[string]string{"message": "Logged out successfully"})
}
