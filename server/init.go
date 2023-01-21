package server

import (
	"github.com/HumXC/simple-douyin/hander"
	"github.com/gin-gonic/gin"
)

type DouYin struct {
	engine *gin.Engine
}

func (d *DouYin) Run(addr string) error {
	return d.engine.Run(addr)
}

func InitDouyin(engine *gin.Engine) *DouYin { // 初始化 douyin
	douyin := engine.Group("douyin")
	initDouyinRoute(douyin)
	return &DouYin{
		engine: engine,
	}
}

// 初始化路由组, 所有的 hander 都在此函数分配
func initDouyinRoute(douyin *gin.RouterGroup) {
	routerHello := douyin.Group("hello")
	routerHello.GET("/", hander.Hello)
}
