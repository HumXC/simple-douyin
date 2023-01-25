package model

import "gorm.io/gorm"

// 保存用户上传的视频
type Video struct {
	gorm.Model
	Hash   string
	Title  string
	UserID string // 上传用户
}

type videoMan struct{}

var VideoMan videoMan = videoMan{}

// 通过 id 获取一个视频记录
func (v *videoMan) GetByID(id string) (Video, error) {
	var video Video
	tx := db.Model(&Video{}).Where("id = ?", id).Find(&video)
	return video, tx.Error
}

// 在数据库里添加一条记录
func (v *videoMan) Put(video Video) error {
	tx := db.Create(&video)
	return tx.Error
}
