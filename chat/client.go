package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aliBazrkar/go-chatapp/model"
	"github.com/gorilla/websocket"
)

type Client struct {
	userID   uint32
	username string
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	lastTS   time.Time
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: in production i should put domain name here
	},
	ReadBufferSize:  1024 * 4,
	WriteBufferSize: 1024 * 8,
}

// Reminder: in case i have decided to add image uploads
// maxMessageSize should change to 64 * 1024
// send channel -> preferably to 1024
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8 * 1024
	sendCapacity   = 1024 / 2
)

func NewClient(userID uint32, username string, hub *Hub, conn *websocket.Conn, lastTS time.Time) *Client {
	return &Client{
		userID:   userID,
		username: username,
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, sendCapacity),
		lastTS:   lastTS,
	}
}

// read messages from WebSocket and send to hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, text, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Username: %v — Hub Addr: %v — Unexpected Error: %v", c.username, c.hub.address, err)
			}
			break
		}

		log.Printf("Username: %v — Hub Addr: %v — Message Sent: %v", c.username, c.hub.address, string(text))

		var msg model.Message
		err = json.Unmarshal(text, &msg)
		if err != nil {
			log.Printf("Username: %v — Hub Addr: %v — Unmarshal Error: %v", c.username, c.hub.address, err)
		}

		c.hub.broadcast <- &model.Message{
			UserID:    c.userID,
			Username:  c.username,
			Content:   msg.Content,
			Timestamp: time.Now(),
		}
	}
}

// write messages from hub to WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("Username: %v — Hub Addr: %v — Write Error: %v", c.username, c.hub.address, err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Username: %v — Hub Addr: %v — Ping Error: %v", c.username, c.hub.address, err)
				return
			}
		}
	}
}

// ticker := time.NewTicker(1 * time.Minute)
// defer ticker.Stop()

// for {
//     select {
//     case <-ticker.C:
//         if tokenExpired(client.UserToken) {
//             client.Conn.WriteMessage(websocket.CloseMessage, []byte("token expired"))
//             client.Hub.Unregister <- client
//             client.Conn.Close()
//             return
//         }
//     case msg := <-client.Send:
//         ...
//     }
// }
