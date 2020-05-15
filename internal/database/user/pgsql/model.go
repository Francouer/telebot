package pgsql

import (
	"database/sql"
	"time"

	"telebot/telebot/CA/internal/domain/entity"
)

type UserModel struct {
	ID           uint
	UserID       uint
	UserName     sql.NullString
	ChatID       int64
	FirstName    sql.NullString
	LastName     sql.NullString
	LanguageCode sql.NullString
	IsBot        bool
	CreatedAt    time.Time
}

func NewUserModel(user entity.User) UserModel {
	newUser := UserModel{
		ID:           user.ID,
		UserID:       user.UserID,
		UserName:     NullableString(user.UserName),
		ChatID:       user.ChatID,
		FirstName:    NullableString(user.FirstName),
		LastName:     NullableString(user.LastName),
		LanguageCode: NullableString(user.LanguageCode),
		IsBot:        user.IsBot,
		CreatedAt:    user.CreatedAt,
	}
	return newUser
}

func ConvertIntoUserEntity(m UserModel) entity.User {
	newEntity := entity.User{
		ID:           m.ID,
		UserID:       m.UserID,
		UserName:     NotNullString(m.UserName),
		ChatID:       m.ChatID,
		FirstName:    NotNullString(m.FirstName),
		LastName:     NotNullString(m.LastName),
		LanguageCode: NotNullString(m.LanguageCode),
		IsBot:        m.IsBot,
		CreatedAt:    m.CreatedAt,
	}
	return newEntity
}

func NullableString(s string) sql.NullString {
	return sql.NullString{s, s != ""}
}

func NotNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	} else {
		return ""
	}
}
