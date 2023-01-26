package service

import (
	"github.com/HumXC/simple-douyin/handler/storage"
	"github.com/gin-gonic/gin"
)

// 该文件实现了一个文件服务器, 用于给用户存储/发送视频

type Storage struct {
	engine *gin.Engine
}
type StorageOption struct {
	DataDir string
}

func NewStorage(g *gin.Engine, option StorageOption) *Storage {
	handler := storage.Handler{
		DataDir: option.DataDir,
	}
	g.GET("hello", handler.Hello)
	storageGroup := g.Group("storage")

	videoGroup := storageGroup.Group("video")
	videoGroup.GET("/:hash", handler.Video)

	return &Storage{
		engine: g,
	}
}
