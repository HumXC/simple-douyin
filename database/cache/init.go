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

func (c *User) Favorite(videoID, userID int64, actionType int32) error {
	key := strconv.FormatInt(videoID, 10) + "." + strconv.FormatInt(userID, 10)
	getAction, err := c.rdb.Get(c.ctx, key).Result()
	if err == redis.Nil {
		c.rdb.Set(c.ctx, key, 1, time.Hour*12)
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
		c.rdb.Incr(c.ctx, key)
	case 2:
		c.rdb.Decr(c.ctx, key)
	}
	return nil
}

// 实现
func (c *Video) CountFavorite(videoID int64) int64 {
	var count int
	key := strconv.FormatInt(videoID, 10)
	for {
		keys, cursor, err := c.rdb.Scan(c.ctx, 0, key+"*", 0).Result()
		if err != nil {
			panic(fmt.Errorf(err.Error()))
		}
		count += len(keys)
		if cursor == 0 {
			break
		}
	}
	return int64(count)
}

func (c *User) CountFavorite(userID int64) int64 {
	var count int
	key := strconv.FormatInt(userID, 10)
	for {
		keys, cursor, err := c.rdb.Scan(c.ctx, 0, "*"+key, 0).Result()
		if err != nil {
			panic(fmt.Errorf(err.Error()))
		}
		count += len(keys)
		if cursor == 0 {
			break
		}
	}
	return int64(count)
}

func (r *User) FavoriteList(userID int64) []int64 {
	List := make([]int64, 0, 10)
	key := strconv.FormatInt(userID, 10)
	for {
		keys, cursor, err := r.rdb.Scan(r.ctx, 0, "*"+key, 0).Result()
		if err != nil {
			panic(fmt.Errorf(err.Error()))
		}
		for _, kkey := range keys {
			svalue, err := r.rdb.Get(r.ctx, kkey).Result()
			ivalue, err := strconv.ParseInt(svalue, 10, 64)
			List = append(List, ivalue)
			if err != nil {
				panic(fmt.Errorf(err.Error()))
			}
		}
		if cursor == 0 {
			break
		}
	}
	return List
}

func (c *User) IsFavorite(userID, videoID int64) bool {
	key := strconv.FormatInt(videoID, 10) + "." + strconv.FormatInt(userID, 10)
	_, err := c.rdb.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	return true
}

type Video struct {
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
