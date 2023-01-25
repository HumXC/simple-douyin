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

func (*ThumbsUp) ActionTypeChange(c *gin.Context, videoId int, userId int) error {
	actionType, err := RDB.Get(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
	sql, err := DB.Find(videoId)
	actionTypeInt, _ := strconv.Atoi(actionType)
	//
	//redis里没找到
	if err == redis.Nil || err == sql.Nil {
		//sql查找是否点赞
		sql, err := DB.Find(videoId).Err
		if err != nil {
			return err
		}
		if sql.find == 0 || sql == 2 {
			//sql里没找到
			return errors.New("date err")
		}
		t := ThumbsUp{
			UserId:     int64(userId),
			VideoId:    int64(videoId),
			ActionType: actionTypeInt,
		}
		//数据追加到redis,同时将已经保存到sql里更新现有的数据
		RDB.Set(c, strconv.Itoa(videoId)+strconv.Itoa(userId), actionType, time.Hour*24)
		DB.Model(&ThumbsUp{}).Save(&s)
	}
	//数据在redis中
	//syscode :=RDB.Get(c,strconv.Itoa(videoId))
	//syscodeInt,err:=strconv.Atoi(syscode.String())
	if actionTypeInt == 2 {
		return errors.New("date err")
	}
	//redis中type-1
	_, err = RDB.Incr(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
	if err != nil {
		return err
	}
	//存到mysql中更新

	//
	return nil
}

//需要token鉴权完善及sql
//redis存储key可以考虑userId+videoId

// 添加一条点赞信息
func (*ThumbsUp) ActionTypeAdd(c *gin.Context, videoId int, userId int) error {
	actionType, err := RDB.Get(c, strconv.Itoa(videoId)+strconv.Itoa(userId)).Result()
	//redis里没有数据
	if err == redis.Nil {

	}
}
