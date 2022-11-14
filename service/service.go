package service

import (
	"user-management-api/model"
)

// UserService DB holds the functions for database access
type UserService interface {
	GetUser(id uint) (user model.User, err error)
	GetUserByEmail(email string) (user model.User, err error)
	GetUserByGoogleId(googleId string) (user model.User, err error)
	CreateUser(user *model.User) (model.User, error)
	UpdateUser(user *model.User, condition model.User) (model.User, error)
}
