package service

import (
	"context"

	"github.com/SawitProRecruitment/UserService/lib/errors"
	"github.com/SawitProRecruitment/UserService/lib/jwt"
	"github.com/SawitProRecruitment/UserService/repository"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) GetByID(ctx context.Context, id int64) (*User, error) {
	user, err := s.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	return ParseUser(user), nil
}

func (s *service) Login(ctx context.Context, payload PayloadLogin) (*ResponseLogin, error) {
	user, err := s.userRepository.GetUserByPhone(ctx, payload.Phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NewBadRequestError("invalid phone or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return nil, errors.NewBadRequestError("invalid phone or password")
	}

	token, err := jwt.GenerateToken(jwt.User{
		ID:    user.Id,
		Name:  user.Name,
		Phone: user.Phone,
	})
	if err != nil {
		return nil, err
	}

	userToken, err := s.userRepository.GetUserToken(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	if userToken == nil {
		err = s.userRepository.InsertToken(ctx, repository.TokenPayloadInsert{
			UserId: user.Id,
			Token:  *token,
		})
	} else {
		err = s.userRepository.UpdateToken(ctx, repository.TokenPayloadUpdate{
			Id:    userToken.Id,
			Token: *token,
		})
	}
	if err != nil {
		return nil, err
	}

	return &ResponseLogin{
		UserId: user.Id,
		Token:  *token,
	}, nil
}

func (s *service) UpdateProfile(ctx context.Context, payload PayloadUpdate) error {
	user, err := s.userRepository.GetUserByPhone(ctx, payload.Phone)
	if err != nil {
		return err
	}
	if user != nil && user.Id != payload.Id {
		return errors.NewConflictError("phone number already used")
	}
	return s.userRepository.UpdateProfile(ctx, repository.User{
		Id:    payload.Id,
		Name:  payload.Name,
		Phone: payload.Phone,
	})
}

func (s *service) InsertUser(ctx context.Context, payload PayloadInsert) (*int64, error) {
	user, err := s.userRepository.GetUserByPhone(ctx, payload.Phone)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, errors.NewConflictError("phone number already used")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	id, err := s.userRepository.InsertUser(ctx, repository.User{
		Name:     payload.Name,
		Phone:    payload.Phone,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}
	return id, nil
}
