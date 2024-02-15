// This file contains types that are used in the repository layer.
package repository

type User struct {
	Id       int64
	Name     string
	Phone    string
	Password string
}

type UserToken struct {
	Id         int64
	UserId     int64
	Token      string
	CountLogin int
}

type TokenPayloadInsert struct {
	UserId int64
	Token  string
}

type TokenPayloadUpdate struct {
	Id    int64
	Token string
}
