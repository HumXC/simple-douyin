package service

import (
	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/middlewares"
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

type DouYin struct {
	engine *gin.Engine
}

func NewDouyin(g *gin.Engine, db *model.DouyinDB, storageClient douyin.StorageClient) *DouYin { // 初始化 douyin
	handler := douyin.Handler{
		DB:            db,
		StorageClient: storageClient,
	}
	douyinGroup := g.Group("douyin")
	douyinGroup.GET("feed", handler.Feed)
	douyinGroup.POST("user/register/", middlewares.PwdHashMiddleWare(), handler.UserRegister)
	douyinGroup.POST("user/login/", handler.UserLogin)
	douyinGroup.GET("user/", middlewares.JWTMiddleWare(), handler.User)
	douyinGroup.POST("comment/action/", middlewares.JWTMiddleWare(), handler.CommentAction)
	douyinGroup.GET("comment/list/", handler.CommentList)

	publish := douyinGroup.Group("publish")
	publish.Use(middlewares.JWTMiddleWare())
	publish.POST("action/", handler.PublishAction)
	publish.GET("list/", handler.PublishList)
	relation := douyinGroup.Group("relation")
	relation.POST("action/", middlewares.JWTMiddleWare(), handler.RelationAction)

	return &DouYin{
		engine: g,
	}
}
