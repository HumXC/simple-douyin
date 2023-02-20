package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ThumbsUp struct {
	gorm.Model
	UserId     int64
	VideoId    int64
	ActionType int32
}

type thumbsUpMan struct {
	db  *gorm.DB
	rdb *redis.Client
}

// 点赞操作，此函数不会对行参做任何验证
func (t *thumbsUpMan) Action(videoID, userID int64, actionType int32) error {
	data := ThumbsUp{
		VideoId: videoID,
		UserId:  userID,
	}
	err := t.db.Where("video_id = ? AND user_id = ?", videoID, userID).
		Find(&data).Error
	if err != nil {
		return err
	}
	data.ActionType = actionType
	return t.db.Save(&data).Error
}

// ActionTypeChange 取消点赞
func (t *thumbsUpMan) ActionTypeChange(c *gin.Context, videoId int, userId int) error {
	//从redis里查
	actionType, rdbErr := t.rdb.Get(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
	//从sqlite里查是否存在数据
	data := ThumbsUp{}
	dbErr := t.db.Where("video_id = ? and user_id = ?", videoId, userId).Find(&data).Error
	actionTypeInt, _ := strconv.Atoi(actionType)
	//redis或者sqlite中数据不存在
	if rdbErr == redis.Nil || dbErr == gorm.ErrRecordNotFound {
		//不存在redis，存在sqlite
		if rdbErr == redis.Nil && dbErr != gorm.ErrRecordNotFound {
			t.rdb.Set(c, strconv.Itoa(videoId)+strconv.Itoa(userId), actionType, time.Hour*24)
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
	_, err := t.rdb.Incr(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
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

// ActionTypeAdd 添加一条点赞信息
func (t *thumbsUpMan) ActionTypeAdd(c *gin.Context, videoId int, userId int) error {
	actionType, rdbErr := t.rdb.Get(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
	data := ThumbsUp{}
	dbErr := t.db.Where("video_id = ? and user_id = ?", videoId, userId).Find(&data).Error
	//redis里没有数据
	if rdbErr == redis.Nil && dbErr == gorm.ErrRecordNotFound {
		t.rdb.Set(c, strconv.Itoa(videoId)+strconv.Itoa(userId), actionType, time.Hour*24)
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
