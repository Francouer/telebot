package user

import (
	"context"
	"fmt"
	"log"

	"telebot/telebot/CA/internal/domain/service"

	"github.com/BorisBorshevsky/timemock"

	"telebot/telebot/CA/internal/domain/entity"
	"telebot/telebot/CA/internal/domain/repository"
)

type Service struct {
	UserRepository repository.UserRepository
}

func (svc Service) CheckExistByUserID(ctx context.Context, userID uint) (bool, error) {
	ok, err := svc.UserRepository.FindByUserID(
		ctx,
		userID,
	)
	if err != nil {
		log.Printf("Check CheckExistByUserID's FindByUserID for %v", err)
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("User is not exist")
	}
	return ok, nil
}

func (svc Service) CheckNotExistByUserID(ctx context.Context, userID uint) error {
	ok, err := svc.UserRepository.FindByUserID(
		ctx,
		userID,
	)
	if err != nil {
		log.Printf("Check CheckNotExistByUserID's FindByUserID for %v", err)
		return err
	}
	if ok {
		return service.ErrUserAlreadyExist
	}
	return nil
}

func (svc *Service) Create(ctx context.Context, params service.CreateUserParams) (*entity.User, error) {
	if params.UserID == 0 {
		return nil, service.ErrUserIDMustBeSet
	}
	if err := svc.CheckNotExistByUserID(ctx, params.UserID); err != nil {
		return nil, err
	}

	var user = &entity.User{
		UserID:    params.UserID,
		CreatedAt: timemock.Now(),
	}

	if err := fillEntity(user, params.UserParams); err != nil {
		return nil, err
	}

	if err := svc.UserRepository.SaveUser(ctx, user); err != nil {
		log.Printf("Check Create's SaveUser for %v", err)
		return nil, err
	}

	log.Println("User was created.")
	return user, nil
}

func (svc Service) GetByIDAndUserID(ctx context.Context, id, userID uint) (*entity.User, error) {
	user, err := svc.UserRepository.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("site not found")
	}
	return user, nil
}

func (svc Service) Update(ctx context.Context, params service.UpdateUserParams) error {
	if params.UserID == 0 {
		return service.ErrUserIDMustBeSet
	}
	if params.ID == 0 {
		return fmt.Errorf("id must be set")
	}
	user, err := svc.GetByIDAndUserID(ctx, params.ID, params.UserID)
	if err != nil {
		return err
	}
	if err := fillEntity(user, params.UserParams); err != nil {
		return err
	}
	if params.ChatID != nil {
		if *params.ChatID == 0 {
			return fmt.Errorf("chat id must be set")
		}
		if err := svc.CheckNotExistByUserID(ctx, params.UserID); err != nil {
			return err
		}
		user.ChatID = *params.ChatID
	}
	if err := svc.UserRepository.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func fillEntity(user *entity.User, params service.UserParams) error {
	if params.ChatID != nil {
		if *params.ChatID > 0 {
			return fmt.Errorf("chatid must be negative number")
		}
		user.ChatID = *params.ChatID
	}
	return nil
}
