package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/SawitProRecruitment/UserService/lib/errors"
	"github.com/SawitProRecruitment/UserService/repository"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type component struct {
	ctx        context.Context
	repository *repository.MockRepositoryInterface
	service    ServiceInterface
	mockedErr  error
}

func setupService(t *testing.T) *component {
	g := gomock.NewController(t)

	repository := repository.NewMockRepositoryInterface(g)
	service := NewService(NewServiceOption{
		UserRepository: repository,
	})

	return &component{
		ctx:        context.Background(),
		repository: repository,
		service:    service,
		mockedErr:  fmt.Errorf("mocked error"),
	}
}

func TestUserService_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("something went wrong when getting user", func(t *testing.T) {
		s := setupService(t)
		s.repository.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(nil, s.mockedErr)

		result, err := s.service.GetByID(s.ctx, 1)
		assert.Nil(t, result)
		assert.Equal(t, s.mockedErr, err)
	})

	t.Run("user not found", func(t *testing.T) {
		s := setupService(t)
		s.repository.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(nil, nil)

		result, err := s.service.GetByID(s.ctx, 1)
		assert.Nil(t, result)
		assert.Equal(t, errors.NewNotFoundError("user not found"), err)
	})

	t.Run("successfully get user", func(t *testing.T) {
		s := setupService(t)
		mockedResult := &repository.User{
			Id:    1,
			Name:  "rotan",
			Phone: "+628123456789",
		}
		expected := &User{
			Id:    1,
			Name:  "rotan",
			Phone: "+628123456789",
		}
		s.repository.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(mockedResult, nil)

		result, err := s.service.GetByID(s.ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestUserService_Login(t *testing.T) {
	t.Parallel()

	t.Run("error getting user by phone", func(t *testing.T) {
		s := setupService(t)
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(nil, s.mockedErr)

		result, err := s.service.Login(s.ctx, PayloadLogin{
			Phone:    "+628123456789",
			Password: "password",
		})
		assert.Nil(t, result)
		assert.Equal(t, s.mockedErr, err)
	})

	t.Run("phone not found", func(t *testing.T) {
		s := setupService(t)
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(nil, nil)

		result, err := s.service.Login(s.ctx, PayloadLogin{
			Phone:    "+628123456789",
			Password: "password",
		})
		assert.Nil(t, result)
		assert.Equal(t, errors.NewBadRequestError("invalid phone or password"), err)
	})

	t.Run("invalid password", func(t *testing.T) {
		s := setupService(t)
		user := &repository.User{
			Id:       1,
			Name:     "rotan",
			Phone:    "+628123456789",
			Password: "hashed_password",
		}
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(user, nil)

		result, err := s.service.Login(s.ctx, PayloadLogin{
			Phone:    "+628123456789",
			Password: "wrong_password",
		})
		assert.Nil(t, result)
		assert.Equal(t, errors.NewBadRequestError("invalid phone or password"), err)
	})

	t.Run("error inserting token", func(t *testing.T) {
		s := setupService(t)
		password := "password"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := &repository.User{
			Id:       1,
			Name:     "rotan",
			Phone:    "+628123456789",
			Password: string(hashedPassword),
		}
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(user, nil)
		s.repository.EXPECT().GetUserToken(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.repository.EXPECT().InsertToken(gomock.Any(), gomock.Any()).Return(s.mockedErr)

		result, err := s.service.Login(s.ctx, PayloadLogin{
			Phone:    "+628123456789",
			Password: password,
		})
		assert.Nil(t, result)
		assert.Equal(t, s.mockedErr, err)
	})

	t.Run("error update token", func(t *testing.T) {
		s := setupService(t)
		password := "password"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := &repository.User{
			Id:       1,
			Name:     "rotan",
			Phone:    "+628123456789",
			Password: string(hashedPassword),
		}
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(user, nil)
		s.repository.EXPECT().GetUserToken(gomock.Any(), gomock.Any()).Return(&repository.UserToken{}, nil)
		s.repository.EXPECT().UpdateToken(gomock.Any(), gomock.Any()).Return(s.mockedErr)

		result, err := s.service.Login(s.ctx, PayloadLogin{
			Phone:    "+628123456789",
			Password: password,
		})
		assert.Nil(t, result)
		assert.Equal(t, s.mockedErr, err)
	})

	t.Run("successfully login", func(t *testing.T) {
		s := setupService(t)
		password := "password"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := &repository.User{
			Id:       1,
			Name:     "rotan",
			Phone:    "+628123456789",
			Password: string(hashedPassword),
		}
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(user, nil)
		s.repository.EXPECT().GetUserToken(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.repository.EXPECT().InsertToken(gomock.Any(), gomock.Any()).Return(nil)

		result, err := s.service.Login(s.ctx, PayloadLogin{
			Phone:    "+628123456789",
			Password: password,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestUserService_UpdateProfile(t *testing.T) {
	t.Parallel()

	t.Run("error getting user by phone", func(t *testing.T) {
		s := setupService(t)
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(nil, s.mockedErr)

		err := s.service.UpdateProfile(s.ctx, PayloadUpdate{
			Id:    1,
			Name:  "rotan",
			Phone: "+628123456789",
		})
		assert.Equal(t, s.mockedErr, err)
	})

	t.Run("phone number already used", func(t *testing.T) {
		s := setupService(t)
		user := &repository.User{
			Id:    2,
			Name:  "other_user",
			Phone: "+628123456789",
		}
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(user, nil)

		err := s.service.UpdateProfile(s.ctx, PayloadUpdate{
			Id:    1,
			Name:  "rotan",
			Phone: "+628123456789",
		})
		assert.Equal(t, errors.NewConflictError("phone number already used"), err)
	})

	t.Run("successfully update profile", func(t *testing.T) {
		s := setupService(t)
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.repository.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(nil)

		err := s.service.UpdateProfile(s.ctx, PayloadUpdate{
			Id:    1,
			Name:  "rotan",
			Phone: "+628123456789",
		})
		assert.NoError(t, err)
	})
}
func TestUserService_InsertUser(t *testing.T) {
	t.Parallel()

	t.Run("error getting user by phone", func(t *testing.T) {
		s := setupService(t)
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(nil, s.mockedErr)

		result, err := s.service.InsertUser(s.ctx, PayloadInsert{
			Name:     "rotan",
			Phone:    "+628123456789",
			Password: "password",
		})
		assert.Nil(t, result)
		assert.Equal(t, s.mockedErr, err)
	})

	t.Run("phone number already used", func(t *testing.T) {
		s := setupService(t)
		user := &repository.User{
			Id:    2,
			Name:  "other_user",
			Phone: "+628123456789",
		}
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(user, nil)

		result, err := s.service.InsertUser(s.ctx, PayloadInsert{
			Name:     "rotan",
			Phone:    "+628123456789",
			Password: "password",
		})
		assert.Nil(t, result)
		assert.Equal(t, errors.NewConflictError("phone number already used"), err)
	})

	t.Run("successfully insert user", func(t *testing.T) {
		s := setupService(t)
		id := int64(1)
		s.repository.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.repository.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&id, nil)

		result, err := s.service.InsertUser(s.ctx, PayloadInsert{
			Name:     "rotan",
			Phone:    "+628123456789",
			Password: "password",
		})
		assert.NoError(t, err)
		assert.Equal(t, id, *result)
	})
}
