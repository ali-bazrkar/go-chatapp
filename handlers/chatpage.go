package handlers

import "net/http"

func (h *Handler) ChatPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/chat.html")
}
