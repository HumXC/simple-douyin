package hander

import "github.com/gin-gonic/gin"

// 示例
func Hello(c *gin.Context) {
	c.Writer.Write([]byte("Hello"))
}
