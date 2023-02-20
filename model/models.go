package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	UserID     int64  //用户ID
	VideoId    int64  //视频ID
	Content    string //评论内容
	gorm.Model        //通用字段
}

type ThumbsUp struct {
	gorm.Model
	UserId     int64
	VideoId    int64
	ActionType int32
}

type Message struct {
	FromUserId int64 //发送者ID
	ToUserId   int64 //接收者ID
	Content    string
	gorm.Model //通用字段
}

type User struct {
	Id             int64 `gorm:"primarykey"`
	Name           string
	Password       string
	TotalFavorited int64
	FavoriteCount  int64
	Avatar         string
	Background     string
	Follows        []User `gorm:"many2many:relations"`
}

type VideoJob struct {
	gorm.Model
	Src    string // 待处理视频的源文件
	Title  string
	UserID int64
}

// 保存用户上传的视频
type Video struct {
	ID            int64 `gorm:"primarykey"`
	Time          time.Time
	Video         string // 视频文件的 hash
	Cover         string // 视频封面的 hash
	Title         string // 视频标题
	UserID        int64  // 视频上传用户ID
	CommentCount  int64  // 视频评论数(用户评论该值加一)
	FavoriteCount int64  // 视频点赞数(用户点赞该值加一)
}
