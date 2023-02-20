package sqldb

import (
	"github.com/HumXC/simple-douyin/model"
	"gorm.io/gorm"
)

// 这个结构体的存在是为了隔离对不同表的操作
// 让 videoMan 只能操作 videos 数据表
type videoMan struct {
	db *gorm.DB
}

// 通过 id 获取一个视频记录
func (v *videoMan) GetByID(id int) (model.Video, error) {
	var video model.Video
	tx := v.db.Model(&model.Video{}).Where("id = ?", id).Find(&video)
	return video, tx.Error
}

// 通过 user_id 获取一个用户发布所有的视频
func (v *videoMan) GetByUser(userID int64) ([]model.Video, error) {
	videos := make([]model.Video, 0, 128)
	tx := v.db.Model(&model.Video{}).Where("user_id = ?", userID).Find(&videos)
	return videos, tx.Error
}

// 按上传时间倒序获取视频, 从latesTime 开始，最多 30 个
func (v *videoMan) GetFeed(latestTime int64, num int) ([]model.Video, error) {
	if num > 30 {
		num = 30
	}
	videos := make([]model.Video, 0, num)
	tx := v.db.Order("time DESC").Where("time>?", latestTime).Debug().Find(&videos)
	return videos, tx.Error
}

// 在数据库里添加一条记录
func (v *videoMan) Put(video model.Video) error {
	tx := v.db.Create(&video)
	return tx.Error
}
