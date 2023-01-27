package douyin

import (
	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.Writer.Write([]byte("Hello! douyin!"))
}
