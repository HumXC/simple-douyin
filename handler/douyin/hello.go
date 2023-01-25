package douyin

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) Hello(c *gin.Context) {
	c.Writer.Write([]byte("Hello! douyin!"))
}
