package db

import "time"

type UserDB struct {
	ID       uint        `gorm:"primaryKey"`
	Username string      `gorm:"uniqueIndex;not null"`
	Password string      `gorm:"not null"`
	Messages []MessageDB `gorm:"foreignKey:UserID"`
	// Tokens   []Token   `gorm:"foreignKey:UserID"`    // For session management
}

type HubDB struct {
	ID       uint        `gorm:"primaryKey"`
	Name     string      `gorm:"not null;"`
	Address  string      `gorm:"uniqueIndex;not null"`
	Messages []MessageDB `gorm:"foreignKey:HubID"`
}

type MessageDB struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	UserID    uint      `gorm:"not null;index"`
	UserDB    UserDB    `gorm:"foreignKey:UserID"`
	HubID     uint      `gorm:"not null;index"`
	HubDB     HubDB     `gorm:"foreignKey:HubID"`
	Timestamp time.Time `gorm:"index"`
}
