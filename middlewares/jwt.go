package middlewares

import (
	"net/http"
	"time"

	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/helper"
	"github.com/gin-gonic/gin"
)

// JWTMiddleWare 鉴权中间件，鉴权并设置user_id
func JWTMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		// token 为空时需要放行，因为像 Feed 接口可以登录访问也可以不登录访问
		// 对于必须登录才能访问的接口，在添加使用下面的 NeedLogin 中间件
		if tokenStr == "" {
			c.Next()
			return
		}
		//验证token
		tokenStruck, err := helper.AnalyseToken(tokenStr)
		if err != nil {
			resp := douyin.BaseResponse()
			resp.Status(douyin.StatusAuthFailed)
			c.JSON(http.StatusOK, resp)
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			resp := douyin.BaseResponse()
			resp.Status(douyin.StatusAuthKeyTimeout)
			c.JSON(http.StatusOK, resp)
			c.Abort() //阻止执行
			return
		}
		c.Set("user_id", tokenStruck.UserId)
		c.Next()
	}
}

// 拦截未登录的请求，强制返回
func NeedLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")
		if userID == 0 {
			resp := douyin.BaseResponse()
			resp.Status(douyin.StatusNeedLogin)
			c.JSON(http.StatusOK, resp)
			c.Abort()
		}
	}
}
