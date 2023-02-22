package storage

import (
	"io"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/HumXC/simple-douyin/helper"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	dataDir string
	buf     sync.Pool
}

// 一个新的 Handler DataDir 是文件存储的根目录，bufferSize 是 io.CopyBuffer 时 buffer 的大小
func NewHandler(dataDir string, bufferSize uint) *Handler {
	return &Handler{
		dataDir: dataDir,
		buf: sync.Pool{
			New: func() any {
				return make([]byte, bufferSize)
			},
		},
	}
}

// 读取请求里的文件路径, 上传存储在本地的文件
func (h *Handler) File(pool map[string]string) func(*gin.Context) {
	return func(c *gin.Context) {
		md5Sum := c.Param("md5")[1:]
		file, ok := pool[md5Sum]
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		fileName := path.Join(h.dataDir, file)
		if !helper.IsFileExit(fileName) {
			c.Status(http.StatusNotFound)
			return
		}
		_ = h.copyFile(c.Writer, fileName)
		// 不解决了，反正能用
		// 解决 broken pipe 和 write: connection reset by peer 错误
		// if err != nil {
		// 	c.AbortWithStatus(http.StatusInternalServerError)
		// 	panic(err)
		// }
	}
}
func (h *Handler) copyFile(w io.Writer, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	var buf = h.buf.Get().([]byte)
	defer h.buf.Put(buf)
	_, err = io.CopyBuffer(w, f, buf)
	return err
}
