package model

import "gorm.io/gorm"

type ThumbsUp struct {
	gorm.Model
	UserId  int64 `gorm:"user_id;type:integer();" json:"user_id"`
	VideoId int64 `gorm:"video_id;type:integer()" json:"video_id"`
}
