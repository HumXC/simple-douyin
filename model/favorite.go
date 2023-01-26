package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ThumbsUp struct {
	gorm.Model
	UserId     int64 `gorm:"user_id;type:integer();" json:"user_id"`
	VideoId    int64 `gorm:"video_id;type:integer()" json:"video_id"`
	ActionType int   `gorm:"action_type;type:integer()" json:"action_type"`
}

type ThumbsUpMan struct {
	db *gorm.DB
}

func (t *ThumbsUpMan) ActionTypeChange(c *gin.Context, videoId int, userId int) error {
	//从redis里查
	actionType, rdbErr := RDB.Get(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
	//从sqlite里查是否存在数据
	data := ThumbsUp{}
	dbErr := t.db.Where("video_id = ? and user_id = ?", videoId, userId).Find(&data).Error
	actionTypeInt, _ := strconv.Atoi(actionType)
	//redis或者sqlite中数据不存在
	if rdbErr == redis.Nil || dbErr == gorm.ErrRecordNotFound {
		//不存在redis，存在sqlite
		if rdbErr == redis.Nil && dbErr != gorm.ErrRecordNotFound {
			RDB.Set(c, strconv.Itoa(videoId)+strconv.Itoa(userId), actionType, time.Hour*24)
			data.ActionType = data.ActionType + 1
			err := t.db.Save(&data).Error
			if err != nil {
				return err
			}
			return nil
			//不存在数据
		} else if rdbErr == redis.Nil && dbErr == gorm.ErrRecordNotFound {
			return errors.New("date err")
		}
	}
	//数据在redis中
	if actionTypeInt == 2 || data.ActionType == 2 {
		return errors.New("date err")
	}
	//redis中type-1
	_, err := RDB.Incr(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
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
func (t *ThumbsUpMan) ActionTypeAdd(c *gin.Context, videoId int, userId int) error {
	actionType, rdbErr := RDB.Get(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
	data := ThumbsUp{}
	dbErr := t.db.Where("video_id = ? and user_id = ?", videoId, userId).Find(&data).Error
	//redis里没有数据
	if rdbErr == redis.Nil && dbErr == gorm.ErrRecordNotFound {
		RDB.Set(c, strconv.Itoa(videoId)+strconv.Itoa(userId), actionType, time.Hour*24)
		data.ActionType = data.ActionType + 1
		err := t.db.Save(&data).Error
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("date err")
	}
}
