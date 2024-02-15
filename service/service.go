package service

import (
	"github.com/SawitProRecruitment/UserService/repository"
	_ "github.com/lib/pq"
)

type service struct {
	userRepository repository.RepositoryInterface
}

type NewServiceOption struct {
	UserRepository repository.RepositoryInterface
}

func NewService(opts NewServiceOption) ServiceInterface {
	return &service{
		userRepository: opts.UserRepository,
	}
}
