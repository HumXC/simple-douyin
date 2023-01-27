package storage

import "github.com/gin-gonic/gin"

func Hello(dataDir string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Write([]byte("Hello! storage! 文件保存在" + dataDir))
	}
}
