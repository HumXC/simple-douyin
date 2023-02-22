package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/go-redis/redis/v8"
)

func (r *RCache) Action(videoID, userID int64, actionType int32) error {
	key := strconv.FormatInt(videoID, 10) + "." + strconv.FormatInt(userID, 10)
	getAction, err := r.rdb.Get(r.ctx, key).Result()
	if err == redis.Nil {
		r.rdb.Set(r.ctx, key, 1, time.Hour*24).Err()
		return nil
	}
	intAction, err := strconv.ParseInt(getAction, 10, 32)
	if int32(intAction) == actionType {
		return fmt.Errorf("重复操作: %d", actionType)
	}
	if err != redis.Nil && err != nil {
		return fmt.Errorf("喜欢错误： %w", err)
	}
	// 点赞还是取消点赞
	switch actionType {
	case 1:
		r.rdb.Incr(r.ctx, key)
	case 2:
		r.rdb.Decr(r.ctx, key)
	}
	return nil
}

// 实现
func (c *RCache) Count(vdeoID int64) int64 {
	return 0
}
func (c *RCache) Sync(duration time.Duration, syncFunc func() error) {
	// 未完成的函数
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
	ctx context.Context
}

func NewDouyinRDB(rdb *redis.Client) *douyin.RDBMan {
	return &douyin.RDBMan{
		Favorite: &RCache{
			rdb: rdb,
			ctx: context.TODO(),
		},
	}
}
