package user

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/BorisBorshevsky/timemock"
	"telebot/telebot/CA/internal/domain/entity"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"telebot/telebot/CA/internal/domain/mock"
	"telebot/telebot/CA/internal/domain/service"
)

func newUserParams(chatID int64) service.UserParams {
	pointer := &chatID
	params := service.UserParams{
		ChatID: pointer,
	}

	return params
}

func New(id, userID uint) *entity.User {
	return &entity.User{
		ID:     id,
		UserID: userID,
	}
}
func NewCreate(userID uint) *entity.User {
	timemock.Freeze(timemock.Now())
	return &entity.User{
		UserID:    userID,
		CreatedAt: timemock.Now(),
	}
}

func TestService_CheckNotExistByUserID(t *testing.T) {
	var (
		sqlError = fmt.Errorf("sql error")
	)
	var cases = []struct {
		name   string
		userID uint

		userRepositoryOK    bool
		userRepositoryError error

		expError error
	}{
		{
			name:   "user is already registered",
			userID: 10,

			userRepositoryOK:    true,
			userRepositoryError: nil,

			expError: service.ErrUserAlreadyExist,
		},
		{
			name:   "repository returned err",
			userID: 0,

			userRepositoryOK:    false,
			userRepositoryError: sqlError,

			expError: sqlError,
		},
		{
			name:   "return nil",
			userID: 10,

			userRepositoryOK:    false,
			userRepositoryError: nil,

			expError: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, c := range cases {

		t.Run(c.name, func(t *testing.T) {
			UserRepository := mock.NewMockUserRepository(ctrl)
			UserRepository.EXPECT().FindByUserID(gomock.Any(), gomock.Any()).
				Return(c.userRepositoryOK, c.userRepositoryError)

			svc := Service{UserRepository}
			err := svc.CheckNotExistByUserID(context.Background(), c.userID)

			assert.Equal(t, c.expError, err)
		})
	}
}

func TestService_Create(t *testing.T) {
	dummyTime := time.Unix(1522549800, 0)
	timemock.Freeze(dummyTime)

	cases := []struct {
		name   string
		userID service.CreateUserParams

		userRepositoryError       error
		saveUserInRepositoryError error
		userRepositoryOk          bool

		expUser  *entity.User
		expError error
	}{
		{
			name: "userID must be set",
			userID: service.CreateUserParams{
				UserID: 0,
			},

			userRepositoryError: nil,
			userRepositoryOk:    false,

			expUser:  nil,
			expError: service.ErrUserIDMustBeSet,
		},
		{
			name: "user is already exist",
			userID: service.CreateUserParams{
				UserID: 1,
			},

			userRepositoryError: nil,
			userRepositoryOk:    true,

			expUser:  nil,
			expError: service.ErrUserAlreadyExist,
		},
		{
			name: "SaveUser error",
			userID: service.CreateUserParams{
				UserID:     1,
				UserParams: newUserParams(0),
			},

			userRepositoryError:       nil,
			saveUserInRepositoryError: fmt.Errorf("SaveUser error"),
			userRepositoryOk:          false,

			expUser:  nil,
			expError: fmt.Errorf("SaveUser error"),
		},
		{
			name: "fillEntity error with ChatID > 0",
			userID: service.CreateUserParams{
				UserID:     1,
				UserParams: newUserParams(2),
			},

			userRepositoryError:       nil,
			saveUserInRepositoryError: nil,
			userRepositoryOk:          false,

			expUser:  nil,
			expError: fmt.Errorf("chatid must be negative number"),
		},
		{
			name: "user created",
			userID: service.CreateUserParams{
				UserID: 1,
			},

			userRepositoryError:       nil,
			saveUserInRepositoryError: nil,
			userRepositoryOk:          false,

			expUser:  NewCreate(1),
			expError: nil,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, c := range cases {
		if c.name == "user created" {
			t.Run(c.name, func(t *testing.T) {
				UserRepository := mock.NewMockUserRepository(ctrl)

				UserRepository.EXPECT().FindByUserID(context.Background(), gomock.Any()).AnyTimes().Return(
					c.userRepositoryOk, c.userRepositoryError)
				UserRepository.EXPECT().SaveUser(context.Background(), gomock.Any()).AnyTimes().Return(
					c.saveUserInRepositoryError)

				svc := Service{UserRepository}
				usr, err := svc.Create(context.Background(), c.userID)

				assert.Equal(t, c.expUser, usr, "Unexpected user")
				assert.Equal(t, c.expError, err, "Unexpected error")
				timemock.Return()
			})
		} else {
			UserRepository := mock.NewMockUserRepository(ctrl)

			UserRepository.EXPECT().FindByUserID(context.Background(), gomock.Any()).AnyTimes().Return(
				c.userRepositoryOk, c.userRepositoryError)
			UserRepository.EXPECT().SaveUser(context.Background(), gomock.Any()).AnyTimes().Return(
				c.saveUserInRepositoryError)

			svc := Service{UserRepository}
			usr, err := svc.Create(context.Background(), c.userID)

			assert.Equal(t, c.expUser, usr, "Unexpected user")
			assert.Equal(t, c.expError, err, "Unexpected error")
		}
	}
}

func TestService_Update(t *testing.T) {
	var cases = []struct {
		name string

		repoUser *entity.User
		params   service.UpdateUserParams

		userParams service.UserParams

		userRepositoryError       error
		userRepositoryByUserIdErr error
		userRepositoryUpdErr      error
		userRepositoryOk          bool

		expError error
	}{
		{
			name: "userID must be set",

			repoUser: nil,
			params: service.UpdateUserParams{
				UserID: 0,
				ID:     1,
			},

			userRepositoryError: nil,
			expError:            service.ErrUserIDMustBeSet,
		},
		{
			name: "ID must be set",

			repoUser: nil,
			params: service.UpdateUserParams{
				UserID: 1,
				ID:     0,
			},

			userRepositoryError: nil,
			expError:            fmt.Errorf("id must be set"),
		},
		{
			name: "user repository updated",

			repoUser: nil,
			params: service.UpdateUserParams{
				UserID: 1,
				ID:     1,
			},

			userRepositoryError: nil,
			expError:            fmt.Errorf("site not found"),
		},
		{
			name: "chat id set 0",

			params: service.UpdateUserParams{
				ID:         1,
				UserID:     2,
				UserParams: newUserParams(0),
			},
			repoUser: New(1, 2),

			userRepositoryError: nil,

			expError: fmt.Errorf("chat id must be set"),
		},
		{
			name: "userRepositoryError",

			params: service.UpdateUserParams{
				ID:     1,
				UserID: 2,
			},

			repoUser: New(1, 2),

			userRepositoryError: fmt.Errorf("userRepositoryError"),

			expError: fmt.Errorf("userRepositoryError"),
		},
		{
			name: "CheckNotExistByUserID",

			params: service.UpdateUserParams{
				ID:         1,
				UserID:     2,
				UserParams: newUserParams(-1),
			},

			repoUser: New(1, 2),

			userRepositoryError:       nil,
			userRepositoryByUserIdErr: fmt.Errorf("FindByIDError"),
			userRepositoryOk:          false,

			expError: fmt.Errorf("FindByIDError"),
		},
		{
			name: "repository update test",

			params: service.UpdateUserParams{
				ID:     1,
				UserID: 2,
			},

			repoUser: New(1, 2),

			userRepositoryError:       nil,
			userRepositoryByUserIdErr: nil,
			userRepositoryUpdErr:      fmt.Errorf("Update error"),
			userRepositoryOk:          false,

			expError: fmt.Errorf("Update error"),
		},
		{
			name: "fillEntity",

			params: service.UpdateUserParams{
				ID:         1,
				UserID:     2,
				UserParams: newUserParams(1),
			},

			repoUser: New(1, 2),

			userRepositoryError:       nil,
			userRepositoryByUserIdErr: nil,
			userRepositoryOk:          false,

			expError: fmt.Errorf("chatid must be negative number"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, c := range cases {

		t.Run(c.name, func(t *testing.T) {
			UserRepository := mock.NewMockUserRepository(ctrl)
			UserRepository.EXPECT().FindByIDAndUserID(context.Background(), gomock.Any(), gomock.Any()).AnyTimes().Return(
				c.repoUser,
				c.userRepositoryError)
			UserRepository.EXPECT().Update(context.Background(), gomock.Any()).AnyTimes().Return(
				c.userRepositoryUpdErr)
			UserRepository.EXPECT().FindByUserID(context.Background(), gomock.Any()).AnyTimes().Return(
				c.userRepositoryOk,
				c.userRepositoryByUserIdErr)

			svc := Service{UserRepository}

			err := svc.Update(context.Background(), c.params)

			assert.Equal(t, c.expError, err, "Unexpected error")
		})
	}
}

func TestService_GetByIDAndUserID(t *testing.T) {
	var cases = []struct {
		name string

		id     uint
		userID uint

		repoUser            *entity.User
		userRepositoryError error

		expUser  *entity.User
		expError error
	}{
		{
			name:   "site not found",
			id:     1,
			userID: 2,

			repoUser:            nil,
			userRepositoryError: nil,

			expUser:  nil,
			expError: fmt.Errorf("site not found"),
		},
		{
			name:   "userRepositoryError not nil",
			id:     2,
			userID: 3,

			repoUser:            nil,
			userRepositoryError: fmt.Errorf("error"),

			expUser:  nil,
			expError: fmt.Errorf("error"),
		},
		{
			name:   "Error is nil, user not nil",
			id:     2,
			userID: 3,

			repoUser:            New(2, 3),
			userRepositoryError: nil,

			expUser:  New(2, 3),
			expError: nil,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			UserRepository := mock.NewMockUserRepository(ctrl)
			UserRepository.EXPECT().FindByIDAndUserID(context.Background(), gomock.Any(), gomock.Any()).Return(
				c.repoUser, c.userRepositoryError)

			svc := Service{UserRepository}
			usr, err := svc.GetByIDAndUserID(context.Background(), c.id, c.userID)

			assert.Equal(t, c.expUser, usr, "Expected user not equals to inc user")

			assert.Equal(t, c.expError, err, "Expected error not equals to incoming error")
		})
	}

}
