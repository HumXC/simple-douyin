package model

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	Id             int64      `json:"id,omitempty"`
	Name           string     `json:"name,omitempty"`
	Password	   string  	  `json:"password,omitempty"`
	FollowCount    int64      `json:"follow_count,omitempty"`
	FollowerCount  int64      `json:"follower_count,omitempty"`
	IsFollow       bool       `json:"is_follow,omitempty"`
	TotalFavorited int64      `json:"total_favorited,omitempty"`
	FavoriteCount  int64      `json:"favorite_count,omitempty"`
}

type userMan struct {
	db *gorm.DB
}

func (u *userMan) GetUserIdByName(name string) (userId int64) {
	user := User{}
	u.db.Where("name=?", name).First(&user)
	return user.Id
}

func (u *userMan) UserIsExistByName(name string) bool {
	err := u.db.Where("name=?", name).First(&User{}).Error
	return err == nil
}

func (u *userMan) CheckNameAndPwd(name string, password string) error {
	err := u.db.Where("name=? and password=?", name, password).First(&User{}).Error
	return err
}

func (u *userMan) AddUser(user *User) error {
	if user == nil {
		return errors.New("AddUser user空指针")
	}
	//注册用户
	return u.db.Create(user).Error
}

func (u *userMan) QueryUserInfoByUserId(userId int64, user *User) error {
	if user == nil {
		return errors.New("AddUser user空指针")
	}
	return u.db.Where("id=?", userId).First(user).Error
}


