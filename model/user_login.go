package model

import (
	"errors"

	"gorm.io/gorm"
)

type UserLogin struct {
	Id       int64  `json:"id,omitempty"`
	UserId   int64	`json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty" gorm:"username"` //用户名
	Password string `json:"password,omitempty"`                 //密码
}

type UserLoginMan struct {
	db *gorm.DB
}

func (u *UserLoginMan) FindUserByUsername(username string) bool {
	var userlogin UserLogin
	err := u.db.Where("user_name=?", username).First(&userlogin).Error
	return err == nil
}

func (u *UserLoginMan) CheckNameAndPwd(username, password string, userLogin *UserLogin) error {
	if userLogin == nil {
		return errors.New("CheckNameAndPwd userLogin空指针")
	}
	err := u.db.Where("user_name=? and password=?", username, password).First(userLogin).Error
	if err != nil {
		return errors.New("用户名或密码错误")
	}
	return nil
}