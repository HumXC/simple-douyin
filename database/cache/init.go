package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/go-redis/redis/v8"
)

//	type RVideoMan interface {
//		// 获取某个视频的点赞数量
//		CountFavorite(vdeoID int64) int64
//	}
//
//	type RUserMan interface {
//		// 获取 user 喜欢视频的数量
//		CountFavorite(userID int64) int64
//		// 获取 user 喜欢的视频列表
//		FavoriteList(userID int64) []int64
//		// user 是否喜欢了 video,userID 可能为0
//		IsFavorite(userID, videoID int64) bool
//		// 没有数据则会创建记录，永远不会删除记录
//		Favorite(userID, videoID int64, actionType int32) error
//	}
type User struct {
	rdb *redis.Client
	ctx context.Context
}

func (r *User) Favorite(userID, videoID int64, actionType int32) error {
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

// TODO
func (c *User) IsFavorite(userID, videoID int64) bool {
	return false
}

// TODO
func (c *User) CountFavorite(vdeoID int64) int64 {
	return 0
}

// TODO
func (c *User) FavoriteList(userID int64) []int64 {
	return nil
}

type Video struct {
	rdb *redis.Client
	ctx context.Context
}

// TODO
func (v *Video) CountFavorite(vdeoID int64) int64 {
	return 0
}

func NewDouyinRDB(addr, pwd string, db int) (*douyin.RDBMan, error) {
	c := context.TODO()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
	return &douyin.RDBMan{
		User: &User{
			rdb: rdb,
			ctx: c,
		},
		Video: &Video{
			rdb: rdb,
			ctx: c,
		},
	}, nil
}
