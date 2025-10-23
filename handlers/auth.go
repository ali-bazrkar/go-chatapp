package handlers

import "net/http"

func (h *Handler) checkAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	user, err := h.sm.ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	respondJSON(w, map[string]any{
		"authenticated": true,
		"username":      user.Username,
	})
}
