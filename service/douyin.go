package service

import (
	"sync"

	"github.com/HumXC/simple-douyin/config"
	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/handler/douyin/videos"
	"github.com/HumXC/simple-douyin/handler/middlewares"
	"github.com/gin-gonic/gin"
)

type DouYin struct {
}

func NewDouyin(g *gin.Engine, conf config.Douyin, db *douyin.DBMan, rdb *douyin.RDBMan, storageClient douyin.StorageClient) *DouYin { // 初始化 douyin
	handler := douyin.Handler{
		DB:            db,
		RDB:           rdb,
		StorageClient: storageClient,
		Avatars:       conf.Avatars,
		Backgrounds:   conf.Backgrounds,
	}
	handler.VideoButcher = videos.NewButcher(
		db.VideoJob,
		conf.VideoButCherMaxJob,
		douyin.VideoButcherFinishFunc(&handler),
	)
	handler.Buf = sync.Pool{
		New: func() any {
			return make([]byte, 512)
		},
	}

	douyin := g.Group("douyin")
	douyin.Use(middlewares.JWTMiddleWare())
	douyin.GET("feed", handler.Feed(conf.FeedNum))
	douyin.POST("user/register/", middlewares.PwdHashMiddleWare(), handler.UserRegister)
	douyin.POST("user/login/", handler.UserLogin)
	douyin.GET("user/", middlewares.NeedLogin(), handler.User)

	comment := douyin.Group("comment")
	comment.Use(middlewares.NeedLogin())
	comment.POST("action/", handler.CommentAction)
	comment.GET("list/", handler.CommentList)

	message := douyin.Group("message")
	message.Use(middlewares.NeedLogin())
	message.POST("action/", handler.MessageAction)
	message.GET("chat/", handler.MessageChatListAction)

	publish := douyin.Group("publish")
	publish.Use(middlewares.NeedLogin())
	publish.POST("action/", handler.PublishAction)
	publish.GET("list/", handler.PublishList)

	relation := douyin.Group("relation")
	relation.Use(middlewares.NeedLogin())
	relation.POST("action/", handler.RelationAction)
	relation.GET("/follow/list/", handler.FollowList)
	relation.GET("/follower/list/", handler.FollowerList)
	relation.GET("/friend/list/", handler.FriendList)

	favorite := douyin.Group("favorite")
	favorite.Use(middlewares.NeedLogin())
	favorite.POST("action/", handler.Favorite)
	favorite.GET("list/", handler.FavoriteList)

	return &DouYin{}
}
