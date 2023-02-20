package sqldb

import "time"

type DBCache interface {
	// 点赞操作
	// 没有数据则会创建记录，永远不会删除记录
	ThumbsUp(videoID, userID int64, actionType int32) error
	// 获取某个视频的点赞数量
	CountThumbsUp(vdeoID int64) int64
	// 每隔 duration 的时间从缓存拉取数据存放传入 syncFunc
	ThumbsUpSync(duration time.Duration, syncFunc func() error)
}
