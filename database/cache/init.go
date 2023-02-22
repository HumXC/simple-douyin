package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/go-redis/redis/v8"
)

func (r *Favorite) Action(videoID, userID int64, actionType int32) error {
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
func (c *Favorite) Count(vdeoID int64) int64 {
	return 0
}
func (c *Favorite) Sync(duration time.Duration, syncFunc func() error) {
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

type Favorite struct {
	rdb *redis.Client
	ctx context.Context
}

func NewDouyinRDB(addr, pwd string, db int) (*douyin.RDBMan, error) {
	c := context.TODO()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
	return &douyin.RDBMan{
		Favorite: &Favorite{
			rdb: rdb,
			ctx: c,
		},
	}, nil
}
