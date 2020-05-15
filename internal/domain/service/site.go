package service

import (
	"context"
	"fmt"

	"telebot/telebot/CA/internal/domain/entity"
)

type CreateSiteParams struct {
	UserID uint
	Name   string

	SiteParams
}

type UpdateSiteParams struct {
	ID     uint
	UserID uint

	Name *string
	SiteParams
}

type SiteParams struct {
	URL            *string
	RequestTimeout *int64
	ResponseStatus *int64
	Description    *string
}

type SiteService interface {
	CheckNotExistByNameAndUserID(ctx context.Context, name string, userID uint) error

	Create(ctx context.Context, params CreateSiteParams) (*entity.Site, error)

	GetByUserID(ctx context.Context, userID uint) ([]entity.Site, error)

	GetIdByUserIdAndName(ctx context.Context, name string, userID uint) (*entity.Site, error)
	GetIdByUserIdAndUrl(ctx context.Context, url string, userID uint) (*entity.Site, error)
	GetByIDAndUserID(ctx context.Context, id, userID uint) (*entity.Site, error)

	Update(ctx context.Context, params UpdateSiteParams) error
	DeleteByIDAndUserID(ctx context.Context, id, userID uint) error
}

var (
	ErrSiteExistWithSameName = fmt.Errorf("site exists with same name")
	ErrSiteNotFound          = fmt.Errorf("site not found")
	ErrUserIDMustBeSet       = fmt.Errorf("userID must be set")
)
