package entity

import "time"

//User -- contain all users entity information(ID, UserID, UserName, ChatID, FirstName, LastName,
//LanguageCode, IsBot, CreateAt)
type User struct {
	ID           uint
	UserID       uint
	UserName     string
	ChatID       int64
	FirstName    string
	LastName     string
	LanguageCode string
	IsBot        bool
	CreatedAt    time.Time
}
