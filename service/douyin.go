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

func NewDouyin(g *gin.Engine, db *model.DouyinDB, uploadFunc douyin.UploadFunc) *DouYin { // 初始化 douyin
	handler := douyin.Handler{
		DB:         db,
		UploadFunc: uploadFunc,
	}
	douyinGroup := g.Group("douyin")
	douyinGroup.GET("feed", handler.Feed)
	douyinGroup.POST("user/register/", handler.UserRegister)
	douyinGroup.POST("user/login/", handler.UserLogin)
	douyinGroup.GET("user/", middlewares.JWTMiddleWare(), handler.User)

	publish := douyinGroup.Group("publish")
	publish.POST("action/", handler.PublishAction)
	return &DouYin{
		engine: g,
	}
}
