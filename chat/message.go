package chat

import "time"

type Message struct {
	ID       string // UUID, primary key
	UserID   string // FK to User
	Username string // cached username for easy retrieval
	Text     string
	SentAt   time.Time
}
