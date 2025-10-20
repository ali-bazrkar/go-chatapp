package db

import "time"

type User struct {
	ID       uint      `gorm:"primaryKey"`
	Username string    `gorm:"uniqueIndex;not null"`
	Password string    `gorm:"not null"`
	Messages []Message `gorm:"foreignKey:UserID"`
	// Tokens   []Token   `gorm:"foreignKey:UserID"`    // For session management
}

type Hub struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"not null;"`
	Address  string    `gorm:"uniqueIndex;not null"`
	Messages []Message `gorm:"foreignKey:HubID"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	UserID    uint      `gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserID"`
	HubID     uint      `gorm:"not null;index"`
	Hub       Hub       `gorm:"foreignKey:HubID"`
	Timestamp time.Time `gorm:"index"`
}
