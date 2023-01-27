package middlewares

import (
	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/helper"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Query("token")
		log.Println("用户token => {}", auth)
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, douyin.Response{
				StatusCode: http.StatusUnauthorized,
				StatusMsg:  "Unauthorized",
			})
			return
		}
		c.Set("user", userClaim)
		c.Set("user_id", userClaim.UserId)
		c.Next()
	}
}
