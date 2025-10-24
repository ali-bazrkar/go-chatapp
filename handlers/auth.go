package handlers

import "net/http"

/*
while we did have middlewares set, this is a helper function
unlike middleware this api is going to be called only by the
frontend in needed cases. since our work involves using the
websocket connection, we might in future need CSRF-less
user authorization. this api is here to satisfy that need.

then again i am still learning, and there might be better
ways to handle these situation (say LocalStorage?) and clearly
my approaches used for front-end interaction might not be the
best possible solution out there, but it seems like a working
and safe method to implement with my current knowledge
*/
func (h *Handler) checkAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	token, err := h.sm.CheckAuth(w, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	respondJSON(w, map[string]any{
		"authenticated": true,
		"username":      token.User.Username,
		"csrf_token":    token.CSRFToken,
	})
}
