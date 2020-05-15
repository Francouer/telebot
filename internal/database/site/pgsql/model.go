package pgsql

import (
	"database/sql"
	"time"

	"telebot/telebot/CA/internal/domain/entity"
)

type SiteModel struct {
	ID             uint
	UserID         uint
	Name           sql.NullString
	URL            sql.NullString
	RequestTimeout int64
	ResponseStatus int64
	Description    sql.NullString
	CreatedAt      time.Time
}

func NewSiteModel(site entity.Site) SiteModel {
	newSite := SiteModel{
		ID:             site.ID,
		UserID:         site.UserID,
		Name:           NullableString(site.Name),
		URL:            NullableString(site.URL),
		RequestTimeout: site.RequestTimeout,
		ResponseStatus: site.ResponseStatus,
		Description:    NullableString(site.Description),
		CreatedAt:      site.CreatedAt,
	}
	return newSite
}

func ConvertIntoSiteEntity(m SiteModel) entity.Site {
	newEntity := entity.Site{
		ID:             m.ID,
		UserID:         m.UserID,
		Name:           NotNullString(m.Name),
		URL:            NotNullString(m.URL),
		RequestTimeout: m.RequestTimeout,
		ResponseStatus: m.ResponseStatus,
		Description:    NotNullString(m.Description),
		CreatedAt:      m.CreatedAt,
	}
	return newEntity
}

func ConvertSliceOfSiteModels(m []SiteModel) []entity.Site {
	var newSliceEntity []entity.Site
	for _, v := range m {
		newEntity := entity.Site{
			ID:             v.ID,
			UserID:         v.UserID,
			Name:           NotNullString(v.Name),
			URL:            NotNullString(v.URL),
			RequestTimeout: v.RequestTimeout,
			ResponseStatus: v.ResponseStatus,
			Description:    NotNullString(v.Description),
			CreatedAt:      v.CreatedAt,
		}
		newSliceEntity = append(newSliceEntity, newEntity)
	}
	return newSliceEntity
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
