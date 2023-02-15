package middlewares

import (
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func PwdHashMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		password := c.Query("password")

		//验证密码规范性
		if utf8.RuneCountInString(password) > 32 {
			c.JSON(http.StatusBadRequest, douyin.Response{
				StatusCode: 1,
				StatusMsg:  errors.New("密码长度超出限制").Error(),
			})
			c.Abort()
			return
		}

		//密码加密
		hash, err := PwdHash(password)
		if err != nil {
			c.JSON(http.StatusBadRequest, douyin.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("hash_password", hash)
		c.Next()
	}
}

func PwdHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}