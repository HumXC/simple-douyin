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
		// host 是服务端主机的 ip, 如果想要运行正常就得自行替换host的内容
		// 例如	"http://192.168.90.148"
		// 提交时请勿修改此值
		URLPrefix: "http://host" + ServeAddr,
	})
	_ = service.NewDouyin(engine, db, storage)
	panic(engine.Run(ServeAddr))
}
