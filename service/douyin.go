package service

import (
	"github.com/HumXC/simple-douyin/config"
	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/middlewares"
	"github.com/HumXC/simple-douyin/model"
	"github.com/HumXC/simple-douyin/videos"
	"github.com/gin-gonic/gin"
)

type DouYin struct {
}

func NewDouyin(g *gin.Engine, conf config.Douyin, db *model.DouyinDB, storageClient douyin.StorageClient) *DouYin { // 初始化 douyin
	handler := douyin.Handler{
		DB:            db,
		StorageClient: storageClient,
		Avatars:       conf.Avatars,
		Backgrounds:   conf.Backgrounds,
	}
	handler.VideoButcher = videos.NewButcher(
		db.VideoJob,
		conf.VideoButCherMaxJob,
		douyin.VideoButcherFinishFunc(&handler),
	)

	douyin := g.Group("douyin")
	douyin.Use(middlewares.JWTMiddleWare())
	douyin.GET("feed", handler.Feed(conf.FeedNum))
	douyin.POST("user/register/", middlewares.PwdHashMiddleWare(), handler.UserRegister)
	douyin.POST("user/login/", handler.UserLogin)
	douyin.GET("user/",
		middlewares.NeedLogin(),
		handler.User,
	)
	douyin.POST("comment/action/", handler.CommentAction)
	douyin.GET("comment/list/", handler.CommentList)

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

	return &DouYin{}
}
