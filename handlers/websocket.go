package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/aliBazrkar/go-chatapp/auth"
	"github.com/aliBazrkar/go-chatapp/chat"
)

func (h *Handler) webSocket(w http.ResponseWriter, r *http.Request) {

	user, ok := auth.GetUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := chat.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}

	// retrieving last messafe timestamp is essential.
	// it helps clients fetch correct messages from
	// database when they are reconnected automatically.

	var lastTS time.Time
	queryStr := r.URL.Query().Get("last_ts")
	if queryStr != "" {
		lastTS, _ = time.Parse(time.RFC3339, queryStr)
	}

	client := chat.NewClient(user.ID, user.Username, h.hub, conn, lastTS)

	h.hub.Register <- client
	go client.WritePump()
	go client.ReadPump()
}
