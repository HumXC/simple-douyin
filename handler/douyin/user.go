package douyin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

// FIXME： 并未实现登录功能，只是为了方便测试
func (h *Handler) User(c *gin.Context) {
	type resp struct {
		Response
		User
	}
	c.JSON(200, resp{
		Response: Response{
			StatusCode: StatusOK,
			StatusMsg:  "OK",
		},
		User: User{},
	})
}
func (h *Handler) UserLogin(c *gin.Context) {
	resp := UserLoginResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "登录成功",
		},
		UserId: 0,
		Token:  "testUser",
	}
	c.JSON(http.StatusOK, resp)
}
