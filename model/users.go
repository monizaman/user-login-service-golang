package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserObject struct {
	Email     string `json:"email"`
	Fullname  string `json:"fullname,omitempty"`
	Telephone string `json:"telephone,omitempty"`
	GoogleId  string `json:"google_id,omitempty"`
}

type User struct {
	gorm.Model
	ID        uint   `json:"id" gorm:"unique"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password,omitempty"`
	Fullname  string `json:"fullname,omitempty"`
	Telephone string `json:"telephone,omitempty"`
	GoogleId  string `json:"google_id,omitempty"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
