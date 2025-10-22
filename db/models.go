package db

import "time"

type Hub struct {
	ID       uint16    `gorm:"primaryKey"`
	Name     string    `gorm:"not null"`
	Address  string    `gorm:"uniqueIndex;not null"`
	Messages []Message `gorm:"foreignKey:HubID"`
}

type User struct {
	ID        uint32    `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Messages  []Message `gorm:"foreignKey:UserID"`
	Tokens    []Token   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Message struct {
	ID        uint32    `gorm:"primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	UserID    uint32    `gorm:"not null;index"`
	User      *User     `gorm:"foreignKey:UserID"`
	HubID     uint16    `gorm:"not null;index"`
	Hub       *Hub      `gorm:"foreignKey:HubID"`
	Timestamp time.Time `gorm:"index"`
}

type Token struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint32    `gorm:"index;not null"`
	User         *User     `gorm:"foreignKey:UserID"`
	SessionToken string    `gorm:"uniqueIndex;not null"`
	CSRFToken    string    `gorm:"not null"`
	ExpiresAt    time.Time `gorm:"index;not null"`
	CreatedAt    time.Time
}
