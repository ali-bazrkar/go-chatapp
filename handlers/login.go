package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aliBazrkar/go-chatapp/auth"
)

// Login authenticates a user and creates a session
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input format
	if !auth.IsUsernameValid(req.Username) || !auth.IsPasswordValid(req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	user, err := h.db.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPassword(req.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionToken, csrfToken, expiresAt, err := h.sm.CreateSession(user.ID, auth.TokenLength)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, sessionToken, expiresAt)

	respondJSON(w, map[string]interface{}{
		"message":    "Logged in successfully",
		"username":   user.Username,
		"csrf_token": csrfToken,
	})
}
