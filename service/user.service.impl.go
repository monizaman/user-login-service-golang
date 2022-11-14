package service

import (
	"context"
	"gorm.io/gorm"
	"user-management-api/model"
)


type UserServiceImpl struct {
	userCollection *gorm.DB
	ctx            context.Context
}

func NewUserService(userCollection *gorm.DB, ctx context.Context) UserService {
	return &UserServiceImpl{userCollection, ctx}
}

func (u UserServiceImpl) GetUser(id uint) (user model.User,  err error) {
	usr := model.User{}
	u.userCollection.Where(&model.User{ID: id}).First(&usr)
	return usr, nil
}

func (u UserServiceImpl) CreateUser(user *model.User) (model.User, error) {
	usr := model.User{}
	res := u.userCollection.Model(&user).Create(user)
	if res.Error != nil {
		return usr, res.Error
	}
	return usr, nil
}

func (u UserServiceImpl) UpdateUser(user *model.User, userWhereClaus model.User) (model.User, error) {
	usr := model.User{}
	res := u.userCollection.Where(&userWhereClaus).Updates(&user)
	if res.Error != nil {
		return usr, res.Error
	}
	usr, err := u.GetUserByEmail(user.Email)
	return usr, err
}

func (u UserServiceImpl) GetUserByEmail(email string) (model.User, error) {
	user := model.User{}
	res := u.userCollection.Where(&model.User{Email: email}).First(&user)
	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}

func (u UserServiceImpl) GetUserByGoogleId(googleId string) (model.User, error) {
	user := model.User{}
	res := u.userCollection.Where(&model.User{GoogleId: googleId}).First(&user)
	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}