package repository

import (
	"context"

	"telebot/telebot/CA/internal/domain/entity"
)

//go:generate mockgen -destination=../mock/user_repository.go -package=mock telebot/telebot/CA/internal/domain/repository UserRepository
type UserRepository interface {
	SaveUser(ctx context.Context, user *entity.User) error

	FindByUserID(ctx context.Context, userID uint) (bool, error)
	FindByIDAndUserID(ctx context.Context, id, userID uint) (*entity.User, error)

	Update(ctx context.Context, user *entity.User) error

	Delete(ctx context.Context, user *entity.User) error
}
