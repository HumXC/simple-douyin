package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id             int64      `json:"id"`
	Name           string     `json:"name"`
	Password	   string  	  `json:"-"`
	FollowCount    int64      `json:"follow_count"`
	FollowerCount  int64      `json:"follower_count"`
	IsFollow       bool       `json:"is_follow"`
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
	user := User{}
	u.db.Where("name=?", name).First(&user)
	err := PwdVerify(user.Password, password)
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
	return u.db.Select("id", "name", "follow_count", "follower_count", "is_follow").
		Where("id=?", userId).First(user).Error
}

func PwdVerify(hashPassword, password string) error {
	//核对密码,比较用户输入的明文和和数据库取出的的密码解析后是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword),[]byte(password))
	if err != nil {
		return errors.New("用户名或密码错误")
	}
	return nil
}


