package storage

import (
	"io"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/gin-gonic/gin"
)

var videoBuf = sync.Pool{
	New: func() any {
		return make([]byte, 512)
	},
}

// 读取请求里的 hash 参数, 上传存储的本地视频
func Fetch(dataDir string) func(*gin.Context) {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		dir := c.Param("dir")
		videoName := path.Join(dataDir, dir, hash)
		_, err := os.Stat(videoName)
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
		}
		f, err := os.Open(videoName)
		// TODO: 将错误写入日志
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		var buf = videoBuf.Get().([]byte)
		io.CopyBuffer(c.Writer, f, buf)
		videoBuf.Put(buf)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
