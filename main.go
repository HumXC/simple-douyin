package main

import (
	"github.com/HumXC/simple-douyin/model"
	"github.com/HumXC/simple-douyin/service"
	"github.com/gin-gonic/gin"
)

func main() {
	const ServeAddr = ":11451"
	// 初始化数据库
	db, err := model.NewDouyinDB("./data.db")
	if err != nil {
		panic(err)
	}

	// 初始化 gin
	engine := gin.Default()

	// 以下两个服务可以使用同一个 gin.Engine, 也可以使用两个不同的 gin.Engine
	storage := service.NewStorage(engine, service.StorageOption{
		DataDir: "./Data",
	})
	_ = service.NewDouyin(engine, db, storage.Upload)
	panic(engine.Run(ServeAddr))
}
