package hander

import "github.com/gin-gonic/gin"

// 示例
func Hello(ctx *gin.Context) {
	ctx.Writer.Write([]byte("Hello"))
}
