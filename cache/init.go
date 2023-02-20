package cache

import (
	"fmt"
	"strconv"
	"time"

	"github.com/HumXC/simple-douyin/model"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	thumbsUpActions map[string]int32
	thumbsUpNum     map[int64]int64
}

func (c *Cache) ThumbsUp(videoID, userID int64, actionType int32) error {
	key := strconv.FormatInt(videoID, 10) + "." + strconv.FormatInt(userID, 10)
	t, ok := c.thumbsUpActions[key]
	if !ok {
		c.thumbsUpActions[key] = actionType
		c.thumbsUpNum[videoID] = c.thumbsUpNum[videoID] + 1
		return nil
	}
	// 点赞情况没有变化，什么时候会出现这种情况呢
	if t == actionType {
		return fmt.Errorf("点赞数量没有变化: %d", actionType)
	}
	c.thumbsUpActions[key] = actionType
	// 点赞还是取消点赞
	switch actionType {
	case 1:
		c.thumbsUpNum[videoID] = c.thumbsUpNum[videoID] + 1
	case 2:
		c.thumbsUpNum[videoID] = c.thumbsUpNum[videoID] - 1
	}
	return nil
}
func (c *Cache) CountThumbsUp(vdeoID int64) int64 {
	return c.thumbsUpNum[vdeoID]
}
func (c *Cache) ThumbsUpSync(duration time.Duration, syncFunc func() error) {
	t := time.NewTimer(duration)
	go func(t *time.Timer) {
		for {
			_ = <-t.C
			syncFunc()
			t.Reset(duration)
		}
	}(t)
}

type RCache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) model.DBCache {
	if rdb == nil {
		return &Cache{
			thumbsUpActions: make(map[string]int32),
			thumbsUpNum:     make(map[int64]int64),
		}
	}
	return nil
}
