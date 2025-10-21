package handlers

import (
	"log"
	"net/http"

	"github.com/aliBazrkar/go-chatapp/chat"
	"github.com/aliBazrkar/go-chatapp/db"
)

// TODO
// this is a temporary holder for username
// as soon as auth and db are set it will update
var id uint = 0

// this is a design choice
// the reason why hub is passed through function
// is allowing scaling for multi-hub creation later
func Setup(dbConn *db.Database, hub *chat.Hub, mux *http.ServeMux) {

	http.HandleFunc("/chat", chatEndpoint)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginEndpoint(dbConn, w, r)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		registerEndpoint(dbConn, w, r)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logoutEndpoint(dbConn, w, r)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(hub, w, r)
	})
}

func chatEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/chat.html")
}

func registerEndpoint(db *db.Database, w http.ResponseWriter, r *http.Request) {}

func loginEndpoint(db *db.Database, w http.ResponseWriter, r *http.Request) {}

func logoutEndpoint(db *db.Database, w http.ResponseWriter, r *http.Request) {}

func wsEndpoint(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := chat.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}

	// TODO - make sure to replace with user
	var idtemp *uint = &id
	*idtemp = *idtemp + 1

	client := chat.NewClient(string(rune(*idtemp)), hub, conn) // fix
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
