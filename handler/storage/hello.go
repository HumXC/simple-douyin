package storage

import "github.com/gin-gonic/gin"

func (h *Handler) Hello(c *gin.Context) {
	c.Writer.Write([]byte("Hello! storage! 文件保存在" + h.DataDir))
}
