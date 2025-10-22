package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aliBazrkar/go-chatapp/auth"
	"github.com/aliBazrkar/go-chatapp/chat"
	"github.com/aliBazrkar/go-chatapp/db"
)

type Handler struct {
	sm  *auth.SessionManager
	hub *chat.Hub
	db  *db.Database
}

func NewHandler(sm *auth.SessionManager, hub *chat.Hub, database *db.Database) *Handler {
	return &Handler{
		sm:  sm,
		hub: hub,
		db:  database,
	}
}

func (h *Handler) Setup(mux *http.ServeMux) {
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/chat", h.ChatPage)
	mux.Handle("/logout", h.sm.Middleware(http.HandlerFunc(h.Logout)))
	mux.Handle("/ws", h.sm.Middleware(http.HandlerFunc(h.WebSocket)))
}

func respondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
