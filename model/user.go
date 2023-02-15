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
	Follows		   []User	  `json:"-" gorm:"many2many:relations"`
}

type userMan struct {
	db *gorm.DB
}

func (u *userMan) GetUserIdByName(name string) (userId int64) {
	user := User{}
	u.db.Where("name=?", name).First(&user)
	return user.Id
}

func (u *userMan) IsUserExistByName(name string) bool {
	err := u.db.Where("name=?", name).First(&User{}).Error
	return err == nil
}

func (u *userMan) IsUserExistById(id int64) bool {
	err := u.db.Where("id=?", id).First(&User{}).Error
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

func (u *userMan) AddUserFollow(userId, followId int64) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE users SET follow_count=follow_count+1 WHERE id = ?", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET follower_count=follower_count+1 WHERE id = ?", followId).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO `relations` (`user_id`,`follow_id`) VALUES (?,?)", userId, followId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (u *userMan) CancelUserFollow(userId, followId int64) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE users SET follow_count=follow_count-1 WHERE id = ? AND follow_count>0", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET follower_count=follower_count-1 WHERE id = ? AND follower_count>0", followId).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM `relations` WHERE user_id=? AND follow_id=?", userId, followId).Error; err != nil {
			return err
		}
		return nil
	})
}

func PwdVerify(hashPassword, password string) error {
	//核对密码,比较用户输入的明文和和数据库取出的的密码解析后是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword),[]byte(password))
	if err != nil {
		return errors.New("用户名或密码错误")
	}
	return nil
}


