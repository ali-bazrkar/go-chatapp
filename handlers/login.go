package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aliBazrkar/go-chatapp/auth"
)

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if !auth.IsUsernameValid(req.Username) || !auth.IsPasswordValid(req.Password) {
		http.Error(w, "Invalid Username/Password", http.StatusNotAcceptable)
		return
	}

	user, err := h.db.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid Username/Password", http.StatusNotAcceptable)
		return
	}

	if !auth.CheckPassword(req.Password, user.Password) {
		http.Error(w, "Invalid Username/Password", http.StatusNotAcceptable)
		return
	}

	sessionToken, csrfToken, expiresAt, err := h.sm.CreateSession(user.ID, auth.TokenLength)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, sessionToken, expiresAt)

	respondJSON(w, map[string]any{
		"message":    "Logged in successfully",
		"username":   user.Username,
		"csrf_token": csrfToken,
	})
}
