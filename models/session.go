package models

import "time"

type Session struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"index;not null"`
	SessionToken string `gorm:"uniqueIndex"`
	ExpiresAt    time.Time
}
