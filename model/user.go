package model

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	Id             int64      `json:"id,omitempty"`
	Name           string     `json:"name,omitempty"`
	FollowCount    int64      `json:"follow_count,omitempty"`
	FollowerCount  int64      `json:"follower_count,omitempty"`
	IsFollow       bool       `json:"is_follow,omitempty"`
	TotalFavorited int64      `json:"total_favorited,omitempty"`
	FavoriteCount  int64      `json:"favorite_count,omitempty"`
	UserLogin      *UserLogin `json:"user_login" gorm:"foreignkey:UserId"`
}

type UserMan struct {
	db *gorm.DB
}

func (u *UserMan) AddUser(user *User) error {
	if user == nil {
		return errors.New("AddUser user空指针")
	}
	//注册用户
	err := u.db.Create(user).Error
	return err
}