// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

//go:generate mockgen -source interfaces.go -destination interfaces.mock.gen.go
type RepositoryInterface interface {
	GetUserById(ctx context.Context, id int64) (*User, error)
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	UpdateProfile(ctx context.Context, user User) error
	InsertUser(ctx context.Context, user User) (*int64, error)

	GetUserToken(ctx context.Context, id int64) (*UserToken, error)
	InsertToken(ctx context.Context, payload TokenPayloadInsert) error
	UpdateToken(ctx context.Context, payload TokenPayloadUpdate) error
}
