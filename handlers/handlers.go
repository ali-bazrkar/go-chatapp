package handlers

import (
	"log"
	"net/http"

	"github.com/aliBazrkar/go-chatapp/chat"
	"github.com/gorilla/websocket"
)

// TODO
// this is a temporary holder for username
// as soon as auth and db are set it will update
var id uint = 0

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: in production i should put domain name here
	},
	ReadBufferSize:  1024 * 4,
	WriteBufferSize: 1024 * 8,
}

// this is a design choice
// the reason why hub is passed through function
// is allowing scaling for multi-hub creation later
func Setup(hub *chat.Hub, mux *http.ServeMux) {
	http.HandleFunc("/chat", chatEndpoint)
	http.HandleFunc("/login", loginEndpoint)
	http.HandleFunc("/register", registerEndpoint)
	http.HandleFunc("/logout", logoutEndpoint)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(hub, w, r)
	})
}

func chatEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/chat" {
		status := http.StatusNotFound
		http.Error(w, "not found", status)
		return
	}

	if r.Method != http.MethodGet {
		status := http.StatusNotFound
		http.Error(w, "Invalid Method", status)
		return
	}

	http.ServeFile(w, r, "templates/chat.html")
}

func registerEndpoint(w http.ResponseWriter, r *http.Request) {}

func loginEndpoint(w http.ResponseWriter, r *http.Request) {}

func logoutEndpoint(w http.ResponseWriter, r *http.Request) {}

func wsEndpoint(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}

	// TODO - make sure to replace with user
	var idtemp *uint = &id
	*idtemp = *idtemp + 1

	client := chat.NewClient(*idtemp, hub, conn)
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
