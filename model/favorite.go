package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
)

type ThumbsUp struct {
	gorm.Model
	UserId     int64 `gorm:"user_id;type:integer();" json:"user_id"`
	VideoId    int64 `gorm:"video_id;type:integer()" json:"video_id"`
	ActionType int   `gorm:"action_type;type:integer()" json:"action_type"`
}

type thumbsUpMan struct {
	db *gorm.DB
}

func (t *thumbsUpMan) ActionTypeChange(c *gin.Context, videoId int, userId int) error {
	//从redis里查
	actionType, err := RDB.Get(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
	//从sqlite里查是否存在数据
	data := ThumbsUp{}
	err = t.db.Where("video_id = ? and user_id = ?", videoId, userId).Find(&data).Error
	actionTypeInt, _ := strconv.Atoi(actionType)
	//数据不存在
	if err == redis.Nil || err == gorm.ErrRecordNotFound {
		return err
		////数据追加到redis,同时将已经保存到sql里更新现有的数据
		//RDB.Set(c, strconv.Itoa(videoId)+strconv.Itoa(userId), actionType, time.Hour*24)
		//DB.Model(&ThumbsUp{}).Save(&s)
	}
	//数据在redis中
	if actionTypeInt == 2 || data.ActionType == 2 {
		return errors.New("date err")
	}
	//redis中type-1
	_, err = RDB.Incr(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
	if err != nil {
		return err
	}
	//存到mysql中更新
	data.ActionType = data.ActionType + 1
	err = t.db.Save(&data).Error
	if err != nil {
		return err
	}
	return nil
}

//需要token鉴权完善及sql
//redis存储key可以考虑userId+videoId

// 添加一条点赞信息
func (*ThumbsUp) ActionTypeAdd(c *gin.Context, videoId int, userId int) error {
	_, err := RDB.Get(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
	//redis里没有数据
	if err == redis.Nil {
		return err
	}
	return nil
}
