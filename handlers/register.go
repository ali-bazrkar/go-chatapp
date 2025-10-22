package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aliBazrkar/go-chatapp/auth"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
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

	if !auth.IsUsernameValid(req.Username) {
		http.Error(w, "Invalid username. Must be 3-20 alphanumeric characters or underscores", http.StatusBadRequest)
		return
	}

	if !auth.IsPasswordValid(req.Password) {
		http.Error(w, "Invalid password. Must be at least 8 characters with letters and digits", http.StatusBadRequest)
		return
	}

	exists, err := h.db.UserExists(req.Username)
	if err != nil {
		log.Printf("User Existance Check Error: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Hashing Password Error: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	user, err := h.db.CreateUser(req.Username, hashedPassword)
	if err != nil {
		log.Printf("User Creation Error: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
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
		"message":    "User registered successfully",
		"username":   user.Username,
		"csrf_token": csrfToken,
	})
}
