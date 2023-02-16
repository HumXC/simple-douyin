package service

import (
	"github.com/HumXC/simple-douyin/config"
	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/middlewares"
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

type DouYin struct {
	engine *gin.Engine
}

func NewDouyin(g *gin.Engine, conf config.Douyin, db *model.DouyinDB, storageClient douyin.StorageClient) *DouYin { // 初始化 douyin
	handler := douyin.Handler{
		DB:            db,
		StorageClient: storageClient,
	}
	douyin := g.Group("douyin")
	douyin.GET("feed", handler.Feed(conf.FeedNum))
	douyin.POST("user/register/", middlewares.PwdHashMiddleWare(), handler.UserRegister)
	douyin.POST("user/login/", handler.UserLogin)
	douyin.GET("user/", middlewares.JWTMiddleWare(), handler.User)
	douyin.POST("comment/action/", middlewares.JWTMiddleWare(), handler.CommentAction)
	douyin.GET("comment/list/", handler.CommentList)

	publish := douyin.Group("publish")
	publish.Use(middlewares.JWTMiddleWare())
	publish.POST("action/", handler.PublishAction)
	publish.GET("list/", handler.PublishList)

	relation := douyin.Group("relation")
	relation.POST("action/", middlewares.JWTMiddleWare(), handler.RelationAction)

	return &DouYin{
		engine: g,
	}
}
