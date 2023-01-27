package storage

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

// 读取请求里的 hash 参数, 上传存储的本地视频
func Video(dataDir string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		hash := ctx.Param("hash")
		videoName := path.Join(dataDir, "videos", hash)
		_, err := os.Stat(videoName)
		if os.IsNotExist(err) {
			ctx.AbortWithStatus(http.StatusNotFound)
		}
		f, err := os.Open(videoName)
		// TODO: 将错误写入日志
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		_, err = io.Copy(ctx.Writer, f)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
