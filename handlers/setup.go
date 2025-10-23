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

	mux.HandleFunc("/", h.serveApp)

	mux.HandleFunc("/api/login", h.login)
	mux.HandleFunc("/api/register", h.register)
	mux.HandleFunc("/api/chat", h.chatApp)

	mux.Handle("/api/logout", h.sm.Middleware(http.HandlerFunc(h.logout)))
	mux.Handle("/ws", h.sm.WebSocketMiddleware(http.HandlerFunc(h.webSocket)))
	mux.HandleFunc("/api/check-auth", h.checkAuth)
}

func (h *Handler) serveApp(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func respondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
