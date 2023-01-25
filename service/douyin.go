package service

import (
	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/gin-gonic/gin"
)

type DouYin struct {
	engine *gin.Engine
}

func NewDouyin(g *gin.Engine) *DouYin { // 初始化 douyin
	handler := douyin.Handler{}
	g.GET("hello", handler.Hello)

	douyinGroup := g.Group("douyin")
	douyinGroup.GET("feed", handler.Feed)

	return &DouYin{
		engine: g,
	}
}
