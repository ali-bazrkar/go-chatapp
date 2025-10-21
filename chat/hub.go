package chat

import (
	"encoding/json"
	"log"

	"github.com/aliBazrkar/go-chatapp/db"
	"github.com/aliBazrkar/go-chatapp/model"
)

type Hub struct {
	id         uint16
	name       string
	address    string
	clients    map[*Client]bool
	broadcast  chan *model.Message
	Register   chan *Client
	unregister chan *Client
	limit      uint16
}

func NewHub(name string, address string, limit uint16) *Hub {
	return &Hub{
		id:         1,
		name:       name,
		address:    address,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *model.Message),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		limit:      limit,
	}
}

func (h *Hub) Run(dbConn *db.Database) {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.Printf("%v connected", client.username)
			// message retrievement ?

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("%v disconnected", client.username)
			}

		case data := <-h.broadcast:

			msg, err := json.Marshal(data)
			if err != nil {
				log.Println("Marshal Error:", err)
				continue
			}

			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}

				// select {} -> TODO: DB CHAN STORAGE

			}
		}
	}
}
