package service

import (
	"context"
	"fmt"

	"telebot/telebot/CA/internal/domain/entity"
)

type CreateUserParams struct {
	UserID uint

	UserParams
}

type UpdateUserParams struct {
	ID     uint
	UserID uint

	UserParams
}

type UserParams struct {
	UserName     *string
	ChatID       *int64
	FirstName    *string
	LastName     *string
	LanguageCode *string
	IsBot        *bool
}

type UserService interface {
	CheckNotExistByUserID(ctx context.Context, userID uint) error

	CheckExistByUserID(ctx context.Context, userID uint) (bool, error)

	Create(ctx context.Context, params CreateUserParams) (*entity.User, error)
	GetByIDAndUserID(ctx context.Context, id, userID uint) (*entity.User, error)
	Update(ctx context.Context, params UpdateUserParams) error
}

var (
	ErrUserAlreadyExist = fmt.Errorf("User is already registred.")
)
