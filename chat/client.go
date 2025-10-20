package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
}

// Reminder: in case i have decided to add image uploads
// maxMessageSize should change to 64 * 1024
// send channel preferably to 1024
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8 * 1024
)

func NewClient(username string, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		username: username,
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 512),
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
				log.Printf("Username: %v — Unexpected Error: %v\n", c.username, err)
			}
			break
		}

		log.Printf("Username: %v — Message Sent: %v\n", c.username, string(text))

		var msg Message
		err = json.Unmarshal(text, &msg)
		if err != nil {
			log.Printf("Username: %v — Unmarshal Error: %v\n", c.username, err)
		}

		c.hub.broadcast <- &Message{Username: c.username, Content: msg.Content, Timestamp: time.Now()}
	}
}

// write messages from hub to WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod) // a timer that sends a signal every N seconds (pingPeriod)

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

			// I have decided to use NextWriter over normal WriteMessage
			// In my case, it helps to drain channel as fast as possible
			// it is also more optimized for sending messages as batches
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Username: %v — NextWriter Error: %v\n", c.username, err)
				return
			}

			w.Write(msg)

			for i := 0; i < len(c.send); i++ {
				if _, err := w.Write(<-c.send); err != nil {
					log.Printf("Username: %v — Write Error: %v\n", c.username, err)
					return
				}
			}

			if err := w.Close(); err != nil {
				log.Printf("Username: %v — Close Writer Error: %v\n", c.username, err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Username: %v — Ping Error: %v\n", c.username, err)
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
