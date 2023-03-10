package douyin

import (
	"github.com/HumXC/simple-douyin/handler/douyin/videos"
	"github.com/HumXC/simple-douyin/model"
)

// 用于管理 SQL 数据库
type DBMan struct {
	Video    VideoMan
	User     UserMan
	Comment  CommentMan
	Message  MessageMan
	VideoJob videos.VideoJobMan
}
type UserMan interface {
	// 获取发布的视频数量
	CountPublished(userID int64) int64
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
	FollowList(userID int64) *[]model.User
	// 获取粉丝用户
	FollowerList(userID int64) *[]model.User
	QueryFriendsById(userId int64, users *[]model.User) error
}
type VideoMan interface {
	// 通过 id 获取一堆视频
	GetByIDs(ids []int64) []model.Video
	// 通过 id 获取一个视频记录
	GetByID(id int64) model.Video
	// 通过 user_id 获取一个用户发布所有的视频
	GetByUser(userID int64) []model.Video
	// 按上传时间倒序获取视频, 从latesTime 开始，最多 30 个
	GetFeed(latestTime int64, num int) []model.Video
	// 在数据库里添加一条记录
	Put(video model.Video)
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
	//查询createAt大于time的所有二人聊天记录
	QueryChat(fromUserId int64, toUserId int64, time string, messages *[]model.Message) error
	//查询最新消息记录
	QueryNewMsg(userId1 int64, userId2 int64, message *model.Message) error
}

// 用于管理 Redis
type RDBMan struct {
	User  RUserMan
	Video RVideoMan
}

type RVideoMan interface {
	// 获取某个视频的点赞数量
	CountFavorite(videoID int64) int64
}
type RUserMan interface {
	// 获取 user 喜欢视频的数量
	CountFavorite(userID int64) int64
	// 获取 user 喜欢的视频列表
	FavoriteList(userID int64) []int64
	// user 是否喜欢了 video,userID 可能为0
	IsFavorite(userID, videoID int64) bool
	// 没有数据则会创建记录，永远不会删除记录
	Favorite(userID, videoID int64, actionType int32) error
}
