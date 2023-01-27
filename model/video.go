package model

import "gorm.io/gorm"

// 保存用户上传的视频
type Video struct {
	gorm.Model
	Hash          string
	Title         string // 视频标题
	UserID        int64  // 视频上传用户ID
	CommentCount  int64  // 视频评论数(用户评论该值加一)
	FavoriteCount int64  // 视频点赞数(用户点赞该值加一)
	// TODO 还有需要完善的字段 PlayUrl CoverUrl
}

// 这个结构体的存在是为了隔离对不同表的操作
// 让 videoMan 只能操作 videos 数据表
type videoMan struct {
	db *gorm.DB
}

// 通过 id 获取一个视频记录
func (v *videoMan) GetByID(id string) (Video, error) {
	var video Video
	tx := v.db.Model(&Video{}).Where("id = ?", id).Find(&video)
	return video, tx.Error
}

// 在数据库里添加一条记录
func (v *videoMan) Put(video Video) error {
	tx := v.db.Create(&video)
	return tx.Error
}
