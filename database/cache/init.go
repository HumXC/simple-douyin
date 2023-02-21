package cache

import (
	"fmt"
	"strconv"
	"time"

	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type FavoriteMan struct {
	thumbsUpActions map[string]int32
	thumbsUpNum     map[int64]int64
}

func (r *RCache) Action(c *gin.Context, videoID, userID int64, actionType int32) error {
	key := strconv.FormatInt(videoID, 10) + "." + strconv.FormatInt(userID, 10)
	getAction, err := r.rdb.Get(c, key).Result()
	if err == redis.Nil {
		r.rdb.Set(c, key, 1, time.Hour*24).Err()
		return nil
	}
	// 点赞情况没有变化，什么时候会出现这种情况呢
	intAction, err := strconv.ParseInt(getAction, 10, 32)
	if int32(intAction) == actionType {
		return fmt.Errorf("点赞数量没有变化: %d", actionType)
	}
	if err != redis.Nil && err != nil {
		panic(fmt.Errorf("点赞错误"))
	}
	// 点赞还是取消点赞
	switch actionType {
	case 1:
		r.rdb.Incr(c, key)
	case 2:
		r.rdb.Decr(c, key)
	}
	return nil
}
func (c *FavoriteMan) Count(vdeoID int64) int64 {
	return c.thumbsUpNum[vdeoID]
}
func (c *FavoriteMan) Sync(duration time.Duration, syncFunc func() error) {
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

func NewDouyinRDB(rdb *redis.Client) (*douyin.RDBMan, error) {
	if rdb == nil {
		return &douyin.RDBMan{
			// Favorite: &FavoriteMan{
			// 	thumbsUpActions: make(map[string]int32),
			// 	thumbsUpNum:     make(map[int64]int64),
			// },
		}, nil
	}
	return nil, nil
}
