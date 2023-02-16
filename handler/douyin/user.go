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
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

type UserInfoResponse struct {
	Response
	User User `json:"user"`
}

func (h *Handler) User(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		CommonResponseError(c, "user_id解析失败")
		return
	}
	user, err := h.user(userID)
	if err != nil {
		CommonResponseError(c, err.Error())
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

// 通过 User ID 从数据库获取一个 User 实例
func (h *Handler) user(id int64) (User, error) {
	u := model.User{}
	err := h.DB.User.QueryUserByUserId(id, &u)
	if err != nil {
		return User{}, err
	}
	return User{
		Id:            u.Id,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		IsFollow:      false,
		Name:          u.Name,
	}, nil
}

func (h *Handler) UserLogin(c *gin.Context) {
	userMan := h.DB.User
	username := c.Query("username")
	password := c.Query("password")

	//用户不存在
	if ok := userMan.IsUserExistByName(username); !ok {
		CommonResponseError(c, "用户不存在")
		return
	}
	//验证用户名和密码
	if err := userMan.CheckNameAndPwd(username, password); err != nil {
		CommonResponseError(c, err.Error())
		return
	}
	//获取user_id
	userId := userMan.GetIdByName(username)
	//生成token
	token, err := helper.GenerateToken(userId)
	if err != nil {
		CommonResponseError(c, err.Error())
		return
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
		CommonResponseError(c, "密码加密失败")
		return
	}
	password := inputPwd.(string)

	//用户名存在
	if ok := userMan.IsUserExistByName(username); ok {
		CommonResponseError(c, "用户已存在")
		return
	}
	//用户名为空
	if username == "" {
		CommonResponseError(c, "用户名为空")
		return
	}
	//用户名长度超出限制
	if utf8.RuneCount([]byte(username)) > 32 {
		CommonResponseError(c, "用户名长度超过限制")
		return
	}
	//保存到数据库
	user := model.User{
		Name:     username,
		Password: password,
	}
	if err := userMan.AddUser(&user); err != nil {
		CommonResponseError(c, err.Error())
		return
	}
	//生成user_id和token
	userId := user.Id
	token, err := helper.GenerateToken(userId)
	if err != nil {
		CommonResponseError(c, err.Error())
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

func CommonResponseError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		StatusCode: -1,
		StatusMsg:  msg,
	})
}
