package db

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliBazrkar/go-chatapp/model"
)

func (db *Database) MessageWriter() {
	for message := range db.WriteQueue {
		if err := db.Gorm.Create(message).Error; err != nil {
			log.Printf("Massage Storage Failed: %v", err)
		}
	}
}

func (db *Database) FetchAfter(hubID uint16, afterTime time.Time) ([]*model.Message, error) {
	var messages []*model.Message
	return messages,
		db.Gorm.
			Model(&Message{}).
			Select("users.username, messages.content, messages.timestamp").
			Joins("JOIN users ON users.id = messages.user_id").
			Where("messages.hub_id = ? AND messages.timestamp > ?", hubID, afterTime).
			Order("messages.timestamp ASC").
			Scan(&messages).Error
}

func (db *Database) FetchRecent(hubID uint16, limit int) ([]*model.Message, error) {

	var messages []*model.Message

	subQuery := db.Gorm.
		Model(&Message{}).
		Select("id").
		Where("hub_id = ?", hubID).
		Order("timestamp DESC").
		Limit(limit)

	return messages,
		db.Gorm.
			Model(&Message{}).
			Select("users.username, messages.content, messages.timestamp").
			Joins("JOIN users ON users.id = messages.user_id").
			Where("messages.id IN (?)", subQuery).
			Order("messages.timestamp ASC").
			Scan(&messages).Error
}

func (db *Database) CreateHub(name string, address string) (*Hub, error) {
	var hub = Hub{Name: name, Address: address}
	return &hub, db.Gorm.Create(&hub).Error
}

func (db *Database) CreateUser(username string, password string) (*User, error) {
	var user = User{Username: strings.ToLower(username), Password: password}
	return &user, db.Gorm.Create(&user).Error
}

func (db *Database) UserExists(username string) (bool, error) {
	var count int64
	return count > 0,
		db.Gorm.
			Model(&User{}).
			Where("username = ?", strings.ToLower(username)).
			Count(&count).Error
}

// usage example :
// dbConn.HubExists("address" | "id", "string" | uint16(int))
func (db *Database) HubExists(field string, value any) (bool, error) {
	var count int64
	return count > 0,
		db.Gorm.
			Model(&Hub{}).
			Where(fmt.Sprintf("%s = ?", field), value).
			Count(&count).Error
}

func (db *Database) GetUserByUsername(username string) (*User, error) {
	var user User
	return &user,
		db.Gorm.
			Model(&User{}).
			Select("id, username, password").
			Where("username = ?", strings.ToLower(username)).
			First(&user).Error
}

// DEV NOTE : never put client input as parameters
// client input can lead to SQL Injection
func (db *Database) GetHub(field string, value any) (*Hub, error) {
	var hub Hub
	return &hub,
		db.Gorm.
			Model(&Hub{}).
			Where(fmt.Sprintf("%s = ?", field), value).
			First(&hub).Error
}

func (db *Database) CreateToken(userID uint32, sessionToken, csrfToken string, expiresAt time.Time) (*Token, error) {
	token := &Token{
		UserID:       userID,
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		ExpiresAt:    expiresAt,
	}
	return token, db.Gorm.Create(token).Error
}

func (db *Database) GetTokenBySession(sessionToken string) (*Token, error) {
	var token Token
	err := db.Gorm.
		Preload("User").
		Where("session_token = ?", sessionToken).
		First(&token).Error
	return &token, err
}

func (db *Database) DeleteToken(sessionToken string) error {
	return db.Gorm.Where("session_token = ?", sessionToken).Delete(&Token{}).Error
}

func (db *Database) CleanupExpiredTokens() error {
	return db.Gorm.Where("expires_at < ?", time.Now()).Delete(&Token{}).Error
}

// DEV Reminder:
// Save() updates all fields, Updates() can update specific fields
