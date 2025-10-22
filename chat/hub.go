package chat

import (
	"encoding/json"
	"log"
	"time"

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
	limit      int
}

func NewHub(id uint16, name string, address string, limit int) *Hub {
	return &Hub{
		id:         id,
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
			log.Printf("%v Connected — Hub Address: %v", client.username, h.address)

			messages, err := h.FetchInitialMessages(client.lastTS, dbConn)
			if err != nil {
				log.Printf("Hub Addr: %v — Client: %v — History Fetch Error: %v", h.address, client.username, err)
			}

			data, err := json.Marshal(messages)
			if err != nil {
				log.Printf("Hub Addr: %v — Marshal Error: %v", h.address, err)
			}

			client.send <- data

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
				log.Printf("%v Disconnected — Hub Address: %v", client.username, h.address)
			}

		case data := <-h.broadcast:

			select {
			case dbConn.WriteQueue <- &db.Message{
				Content:   data.Content,
				UserID:    data.UserID,
				HubID:     h.id,
				Timestamp: data.Timestamp,
			}:
			default:
				log.Printf("Hub Addr: %v — Accessing WriteQueue Failed", h.address)
			}

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
			}
		}
	}
}

func (h *Hub) FetchInitialMessages(lastTS time.Time, dbConn *db.Database) ([]*model.Message, error) {
	var messages []*model.Message
	var err error

	if lastTS.IsZero() {
		messages, err = dbConn.FetchRecent(h.id, h.limit)
	} else {
		messages, err = dbConn.FetchAfter(h.id, lastTS)
	}
	if err != nil {
		return nil, err
	}

	return messages, nil
}
