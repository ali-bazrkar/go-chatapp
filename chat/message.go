package chat

import (
	"time"
)

type Message struct {
	Username  uint      `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
