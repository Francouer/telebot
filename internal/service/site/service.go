package site

import (
	"context"
	"fmt"
	"log"
	"time"

	"telebot/telebot/CA/internal/domain/entity"
	"telebot/telebot/CA/internal/domain/repository"
	"telebot/telebot/CA/internal/domain/service"
)

type Service struct {
	SiteRepository repository.SiteRepository
	UserService    service.UserService
}

func (svc Service) CheckNotExistByNameAndUserID(ctx context.Context, name string, userID uint) error {
	_, err := svc.SiteRepository.FindByUserIdAndName(
		ctx,
		userID,
		name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (svc Service) GetByUserID(ctx context.Context, userID uint) ([]entity.Site, error) {
	site, err := svc.SiteRepository.FindByUserID(ctx, userID)
	if err != nil {
		log.Printf("Check GetByUserID's FindByUserID for %v", err)
		return nil, err
	}
	log.Printf("Site, err - %v, %e", site, err)
	return site, nil
}

//Create - services creating incoming url in the database's table "sites"
func (svc *Service) Create(ctx context.Context, params service.CreateSiteParams) (*entity.Site, error) {

	if params.UserID == 0 {
		return nil, fmt.Errorf("user id must be set")
	}
	if params.Name == "" {
		return nil, fmt.Errorf("name must be set")
	}

	_, err := svc.UserService.CheckExistByUserID(ctx, params.UserID)
	if err != nil {
		return nil, err
	}

	if err := svc.CheckNotExistByNameAndUserID(ctx, params.Name, params.UserID); err != nil {
		return nil, err
	}

	var site = &entity.Site{
		UserID: params.UserID,
		Name:   params.Name,
	}

	if err := fillEntity(site, params.SiteParams); err != nil {
		return nil, err
	}

	if err := svc.SiteRepository.Save(ctx, site); err != nil {
		return nil, err
	}

	return site, nil
}

func (svc Service) GetIdByUserIdAndName(ctx context.Context, name string, userID uint) (*entity.Site, error) {
	site, err := svc.SiteRepository.FindByUserIdAndName(ctx, userID, name)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (svc Service) GetIdByUserIdAndUrl(ctx context.Context, name string, userID uint) (*entity.Site, error) {
	site, err := svc.SiteRepository.FindByUserIdAndURL(ctx, userID, name)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (svc Service) GetByIDAndUserID(ctx context.Context, id, userID uint) (*entity.Site, error) {

	site, err := svc.SiteRepository.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		log.Printf("Check GetByIDAndUserID's FindByIDAndUserID for %v", err)
		return nil, err
	}
	if site == nil {
		return nil, service.ErrSiteNotFound
	}

	return site, nil
}

func (svc Service) Update(ctx context.Context, params service.UpdateSiteParams) error {
	if params.UserID == 0 {
		return fmt.Errorf("user id must be set")
	}
	if params.ID == 0 {
		return fmt.Errorf("id must be set")
	}
	site, err := svc.GetByIDAndUserID(ctx, params.ID, params.UserID)
	if err != nil {
		return err
	}

	if params.Name != nil {
		if *params.Name == "" {
			return fmt.Errorf("name must be set")
		}
		if err := svc.CheckNotExistByNameAndUserID(ctx, *params.Name, params.UserID); err != nil {
			return err
		}

		site.Name = *params.Name

	}

	if err := fillEntity(site, params.SiteParams); err != nil {
		return err
	}

	if err := svc.SiteRepository.Update(ctx, site); err != nil {
		return err
	}

	return nil
}

func (svc Service) DeleteByIDAndUserID(ctx context.Context, id, userID uint) error {
	if userID == 0 {
		return fmt.Errorf("user id must be set")
	}
	if id == 0 {
		return fmt.Errorf("id must be set")
	}

	site, err := svc.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}

	return svc.SiteRepository.Delete(ctx, site.ID)
}

//Functions
//fillSiteEntityWithParams returns new site entity with site's params.
func fillSiteEntityWithParams(ctx context.Context, params service.CreateSiteParams) entity.Site {
	site := entity.Site{
		UserID:         params.UserID,
		Name:           params.Name,
		URL:            *params.URL,
		RequestTimeout: *params.RequestTimeout,
		ResponseStatus: *params.ResponseStatus,
		Description:    *params.Description,
		CreatedAt:      time.Now(),
	}
	return site
}

func fillEntity(site *entity.Site, params service.SiteParams) error {
	if params.URL != nil {
		if len(*params.URL) > 1000 {
			return fmt.Errorf("url must be less than 1000")
		}
		site.URL = *params.URL
	}
	return nil
}
