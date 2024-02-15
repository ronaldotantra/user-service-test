package service

import "github.com/SawitProRecruitment/UserService/repository"

type PayloadInsert struct {
	Name     string `json:"name" validate:"required,min=3,max=60"`
	Phone    string `json:"phone" validate:"required,customPhone"`
	Password string `json:"password" validate:"required,customPassword"`
}

type PayloadLogin struct {
	Phone    string `json:"phone" validate:"required,customPhone"`
	Password string `json:"password" validate:"required,customPassword"`
}

type PayloadUpdate struct {
	Id    int64
	Name  string `json:"name" validate:"required,min=3,max=60"`
	Phone string `json:"phone" validate:"required,customPhone"`
}

type User struct {
	Id    int64
	Name  string
	Phone string
}

func ParseUser(userRepo *repository.User) *User {
	return &User{
		Id:    userRepo.Id,
		Name:  userRepo.Name,
		Phone: userRepo.Phone,
	}
}

type ResponseLogin struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
