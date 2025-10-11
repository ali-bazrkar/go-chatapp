package chat

import "github.com/gorilla/websocket"

type Client struct {
	Username string
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
}
