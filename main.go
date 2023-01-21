package main

import (
	"github.com/HumXC/simple-douyin/server"
	"github.com/gin-gonic/gin"
)

func main() {
	const ServeAddr = "localhost:11451"
	// 初始化数据库
	// ...

	// 初始化 gin
	engine := gin.Default()
	douyin := server.InitServer(engine)

	panic(douyin.Run(ServeAddr))
}
