package douyin

import (
	"net/http"
	"unicode/utf8"

	"github.com/HumXC/simple-douyin/helper"
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserId 	  int64  `json:"user_id,omitempty"`
	Token  	  string `json:"token,omitempty"`
}

type UserInfoResponse struct {
	Response
	User model.User `json:"user"` 	
}

func (h *Handler) User(c *gin.Context) {
	userMan := h.DB.User
	inputId, ok := c.Get("user_id")
	if !ok {
		ResponseError(c, "user_id解析失败")
		return
	}
	userId := inputId.(int64)

	user := model.User{}
	err := userMan.QueryUserInfoByUserId(userId, &user)
	if err != nil {
		ResponseError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, UserInfoResponse{
		Response: Response{
			StatusCode: StatusOK,
			StatusMsg:  "OK",
		},
		User: user,
	})
}

func (h *Handler) UserLogin(c *gin.Context) {
	userMan := h.DB.User
	username := c.Query("username")
	password := c.Query("password")

	//用户不存在
	if ok := userMan.UserIsExistByName(username); !ok {
		ResponseError(c, "用户不存在")
		return
	}
	//验证用户名和密码
	if err := userMan.CheckNameAndPwd(username, password); err != nil {
		ResponseError(c, err.Error())
		return
	}
	//获取user_id
	userId := userMan.GetUserIdByName(username)
	//生成token
	token ,err := helper.GenerateToken(userId)
	if err != nil {
		ResponseError(c, err.Error())
	}
	
	resp := UserLoginResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "登录成功",
		},
		UserId: userId,
		Token:  token,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UserRegister(c *gin.Context) {
	userMan := h.DB.User
	username := c.Query("username")
	inputPwd, ok := c.Get("hash_password")
	if !ok {
		ResponseError(c, "密码加密失败")
	}
	password := inputPwd.(string)
	
	//用户名存在
	if ok := userMan.UserIsExistByName(username); ok {
		ResponseError(c, "用户已存在")
		return
	}
	//用户名为空
	if username == "" {
		ResponseError(c, "用户名为空")
		return
	}
	//用户名长度超出限制
	if utf8.RuneCount([]byte(username)) > 32 {
		ResponseError(c, "用户名长度超过限制")
		return
	}
	//保存到数据库
	user := model.User{
		Name: username,
		Password: password,
	}
	if err := userMan.AddUser(&user); err != nil {
		ResponseError(c, err.Error())
		return
	}
	//生成user_id和token
	userId := user.Id
	token, err := helper.GenerateToken(userId)
	if err != nil {
		ResponseError(c, err.Error())
		return
	}

	resp := UserLoginResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "注册成功",
		},
		UserId: userId,
		Token:  token,
	}
	c.JSON(http.StatusOK, resp)
}

func ResponseError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{
			StatusCode: -1,
			StatusMsg: msg,
		},
	})
}
