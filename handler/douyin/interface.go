package douyin

import (
	"github.com/HumXC/simple-douyin/handler/douyin/videos"
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

// 用于管理 SQL 数据库
type DBMan struct {
	Video    VideoMan
	User     UserMan
	Comment  CommentMan
	Message  MessageMan
	ThumbsUp ThumbsUpMan
	VideoJob videos.VideoJobMan
}
type UserMan interface {
	// 返回粉丝数量
	CountFollower(userID int64) int64
	// 返回关注数量
	CountFollow(userID int64) int64
	// 返回 user1 是否关注了 user2
	// 如果 user1 关注了 user2，返回 true
	IsFollow(user1, user2 int64) bool
	GetIdByName(name string) (userId int64)
	IsExistWithName(name string) bool
	IsExistWithId(id int64) bool
	CheckNameAndPwd(name string, password string) error
	AddUser(user *model.User) error
	QueryById(userId int64, user *model.User) error
	Follow(userId, followId int64) error
	CancelFollow(userId, followId int64) error
	// 获取关注者用户
	QueryFollows(userID int64) *[]model.User
	// 获取粉丝用户
	QueryFollowers(userID int64) *[]model.User
	QueryFriendsById(userId int64, users *[]model.User) error
}
type VideoMan interface {
	// 通过 id 获取一个视频记录
	GetByID(id int) (model.Video, error)
	// 通过 user_id 获取一个用户发布所有的视频
	GetByUser(userID int64) ([]model.Video, error)
	// 按上传时间倒序获取视频, 从latesTime 开始，最多 30 个
	GetFeed(latestTime int64, num int) ([]model.Video, error)
	// 在数据库里添加一条记录
	Put(video model.Video) error
}
type CommentMan interface {
	AddComment(comment *model.Comment) error
	AddCommentAndUpdateCommentCount(comment *model.Comment) error
	DeleteCommentAndUpdateCountById(commentId, videoId int64) error
	QueryCommentById(id int64, comment *model.Comment) error
	QueryCommentListByVideoId(videoId int64, comments *[]model.Comment) error
}
type MessageMan interface {
	AddMessage(message *model.Message) error
	QueryMessageRecord(fromUserId int64, toUserId int64, messages *[]model.Message) error
}
type ThumbsUpMan interface {
	Action(videoID, userID int64, actionType int32) error
	ActionTypeChange(c *gin.Context, videoId int, userId int) error
	ActionTypeAdd(c *gin.Context, videoId int, userId int) error
}

// 用于管理 Redis
type RDBMan struct {
	Favorite FavoriteMan
}

type FavoriteMan interface {
	// 点赞操作
	// 没有数据则会创建记录，永远不会删除记录
	Action(videoID, userID int64, actionType int32) error
	// 获取某个视频的点赞数量
	Count(vdeoID int64) int64
	// 每隔 duration 的时间从缓存拉取数据存放传入 syncFunc
	// TODO ThumbsUpSync(duration time.Duration, syncFunc func() error)
}
