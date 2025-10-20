package db

import (
	"log"
	"time"
)

func (db *Database) FetchAfter(hubID int, afterTime time.Time) ([]struct {
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}, error) {

	var messages []struct {
		Username  string    `json:"username"`
		Content   string    `json:"content"`
		Timestamp time.Time `json:"timestamp"`
	}

	result := db.Gorm.
		Model(&Message{}).
		Select("users.username, messages.content, messages.timestamp").
		Joins("User").
		Where("messages.hub_id = ? AND messages.timestamp > ?", hubID, afterTime).
		Order("messages.timestamp ASC").
		Scan(&messages)

	if result.Error != nil {
		log.Printf("Message Fetch Failed: %v\n", result.Error)
		return nil, result.Error
	}

	return messages, nil
}

func (db *Database) FetchRecent(hubID int, limit int) ([]struct {
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}, error) {

	var messages []struct {
		Username  string    `json:"username"`
		Content   string    `json:"content"`
		Timestamp time.Time `json:"timestamp"`
	}

	// Subquery: get latest N messages DESC
	subQuery := db.Gorm.
		Model(&Message{}).
		Select("id").
		Where("hub_id = ?", hubID).
		Order("timestamp DESC").
		Limit(limit)

	// Main query: fetch selected columns for those IDs, ordered ASC
	err := db.Gorm.
		Model(&Message{}).
		Select("users.username, messages.content, messages.timestamp").
		Joins("JOIN users ON users.id = messages.user_id").
		Where("messages.id IN (?)", subQuery).
		Order("messages.timestamp ASC").
		Scan(&messages).Error

	if err != nil {
		return nil, err
	}
	return messages, nil
}
