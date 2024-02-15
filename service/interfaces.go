package service

import "context"

//go:generate mockgen -source interfaces.go -destination interfaces.mock.gen.go -package=service
type ServiceInterface interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	Login(ctx context.Context, payload PayloadLogin) (*ResponseLogin, error)
	UpdateProfile(ctx context.Context, payload PayloadUpdate) error
	InsertUser(ctx context.Context, payload PayloadInsert) (*int64, error)
}
