package chat

import (
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *Message
	Register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run(dbConn *gorm.DB) {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.Printf("%v connected", client.username)

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
