package repository

import (
	"context"

	"telebot/telebot/CA/internal/domain/entity"
)

type SiteRepository interface {
	Save(ctx context.Context, site *entity.Site) error

	FindByUserID(ctx context.Context, userID uint) ([]entity.Site, error)

	FindByUserIdAndURL(ctx context.Context, userID uint, URL string) (*entity.Site, error)
	FindByUserIdAndName(ctx context.Context, userID uint, Name string) (*entity.Site, error)

	FindByIDAndUserID(ctx context.Context, id, userID uint) (*entity.Site, error)
	Update(ctx context.Context, site *entity.Site) error
	Delete(ctx context.Context, id uint) error
}
