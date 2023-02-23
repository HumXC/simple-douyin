package sqldb

import (
	"errors"

	"github.com/HumXC/simple-douyin/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserMan struct {
	DB *gorm.DB
}

func (u *UserMan) CountPublished(userID int64) (result int64) {
	u.DB.Table("videos").Where("user_id=?", userID).Count(&result)
	return
}
func (u *UserMan) IsFavorite(userID, videoID int64) bool {
	hasID := u.DB.Model(&model.User{ID: userID}).
		Where("video_id=?", videoID).
		Association("Favorites").Count()
	return hasID != 0
}
func (u *UserMan) Favorite(userID, videoID int64) error {
	if u.IsFavorite(userID, videoID) {
		return nil
	}
	u.DB.Model(&model.User{
		ID: userID,
	}).Association("Favorites").Append(&model.Video{ID: videoID})
	return nil
}

func (u *UserMan) FavoriteList(userID int64) (v []model.Video) {
	u.DB.Model(&model.User{
		ID: userID,
	}).Association("Follows").Find(&v)
	return
}

// 返回粉丝数量
func (u *UserMan) CountFollower(userID int64) int64 {
	var count int64 = 0
	u.DB.Table("follows").Where("follow_id=?", userID).Select("follow_id").Count(&count)
	return count
}

// 返回关注数量
func (u *UserMan) CountFollow(userID int64) int64 {
	return u.DB.Model(&model.User{
		ID: userID,
	}).Select("id").Association("Follows").Count()
}

// 返回 user1 是否关注了 user2
// 如果 user1 关注了 user2，返回 true
func (u *UserMan) IsFollow(user1, user2 int64) bool {
	if user1 == 0 {
		return false
	}
	id := 0
	_ = u.DB.Model(&model.User{
		ID: user1,
	}).Select("id").Where("id=?", user2).Association("Follows").Find(&id)
	return id != 0
}
func (u *UserMan) GetIdByName(name string) (userId int64) {
	var id int64
	u.DB.Model(&model.User{}).Select("id").Where("name=?", name).Find(&id)
	return id
}

func (u *UserMan) IsExistWithName(name string) bool {
	var id int64
	err := u.DB.Model(&model.User{}).Select("id").Where("name=?", name).Find(&id).Error
	return err == nil && id != 0
}

func (u *UserMan) IsExistWithId(id int64) bool {
	var _id int64
	err := u.DB.Model(&model.User{}).Select("id").Where("id=?", id).First(&_id).Error
	return err == nil && _id != 0
}

func (u *UserMan) CheckNameAndPwd(name string, password string) error {
	pwd := ""
	u.DB.Model(&model.User{}).Select("password").Where("name=?", name).Find(&pwd)
	err := PwdVerify(pwd, password)
	return err
}

func (u *UserMan) AddUser(user *model.User) error {
	if user == nil {
		return errors.New("AddUser user 空指针")
	}
	//注册用户
	return u.DB.Create(user).Error
}

func (u *UserMan) QueryById(userId int64, user *model.User) error {
	if user == nil {
		return errors.New("QueryById user空指针")
	}
	// 此处 Omit 是因为外部从来都不需要 password
	err := u.DB.Omit("password").Where("id=?", userId).Find(user).Error
	return err
}

func (u *UserMan) Follow(userID, followId int64) error {
	hasID := u.DB.Model(&model.User{ID: userID}).
		Where("follow_id=?", followId).
		Association("Follows").Count()
	if hasID != 0 {
		return nil
	}
	u.DB.Model(&model.User{
		ID: userID,
	}).Association("Follows").Append(&model.User{ID: followId})
	return nil
}

func (u *UserMan) CancelFollow(userId, followId int64) error {
	type follow struct {
		UserID   int64
		FollowID int64
	}
	return u.DB.Table("follows").Where("user_id=? AND follow_id=?", userId, followId).Delete(&follow{}).Error
}

// 获取关注者用户
func (u *UserMan) FollowList(userID int64) *[]model.User {
	result := make([]model.User, 0)
	if userID == 0 {
		return &result
	}
	u.DB.Model(&model.User{
		ID: userID,
	}).Omit("password").Association("Follows").Find(&result)
	return &result
}

// 获取粉丝用户
func (u *UserMan) FollowerList(userID int64) *[]model.User {
	result := make([]model.User, 0)
	// 能用就行
	subQuery := u.DB.Table("follows").Where("follow_id=?", userID).Select("user_id")
	u.DB.Model(&model.User{}).Omit("password").Where("id IN (?)", subQuery).Find(&result)
	return &result
}

// Deprecated: 使用 FollowList
func (u *UserMan) QueryFollowsById(userId int64, users *[]model.User) error {
	return u.DB.Table("users as u").
		Select([]string{"id", "name", "follow_count", "follower_count"}).
		Joins("left join relations as r on u.id = r.follow_id").
		Where("r.user_id=?", userId).Find(users).Error
}

// Deprecated: 使用 FollowerList
func (u *UserMan) QueryFollowersById(userId int64, users *[]model.User) error {
	return u.DB.Table("users as u").
		Select([]string{"id", "name", "follow_count", "follower_count"}).
		Joins("left join relations as r on u.id = r.user_id").
		Where("r.follow_id=?", userId).Find(users).Error
}

func (u *UserMan) QueryFriendsById(userId int64, users *[]model.User) error {
	subQuery := u.DB.Raw("select a.follow_id from follows as a join follows b on a.user_id=b.follow_id and a.follow_id=b.user_id and a.user_id=?", userId)
	return u.DB.Model(&model.User{}).Where("id IN (?)", subQuery).Find(users).Error
}

func PwdVerify(hashPassword, password string) error {
	//核对密码,比较用户输入的明文和和数据库取出的的密码解析后是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return errors.New("用户名或密码错误")
	}
	return nil
}
