package douyin

import (
	"math/rand"
	"net/http"
	"unicode/utf8"

	"github.com/HumXC/simple-douyin/handler/ginx"
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) User(c *gin.Context) {
	type Resp struct {
		Response
		User User `json:"user"`
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()

	userID := c.GetInt64("user_id")
	if userID == 0 {
		resp.Status(StatusAuthFailed)
		return
	}
	u := model.User{}
	err := h.DB.User.QueryById(userID, &u)
	if err != nil {
		resp.Status(StatusUserNotFound)
		return
	}
	user := h.ConvertUser(u, false)
	resp.User = user
}

func (h *Handler) UserLogin(c *gin.Context) {
	type Resp struct {
		Response
		UserId int64  `json:"user_id,omitempty"`
		Token  string `json:"token,omitempty"`
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()

	userMan := h.DB.User
	username := c.Query("username")
	password := c.Query("password")

	//用户不存在
	if ok := userMan.IsExistWithName(username); !ok {
		resp.Status(StatusAuthFailed)
		return
	}
	//验证用户名和密码
	if err := userMan.CheckNameAndPwd(username, password); err != nil {
		resp.Status(StatusOtherError)
		return
	}
	//获取user_id
	userId := userMan.GetIdByName(username)
	//生成token
	token, err := ginx.GenerateToken(userId)
	if err != nil {
		resp.Status(StatusOtherError)
		return
	}
	//登录成功，包装resp
	resp.UserId = userId
	resp.Token = token
}

// 从切片里随机选择一个元素
func PickOne[T any](list []T) T {
	result := *new(T)
	if len(list) == 0 {
		return result
	}
	i := rand.Intn(len(list))
	result = list[i]
	return result
}

func (h *Handler) UserRegister(c *gin.Context) {
	type Resp struct {
		Response
		UserId int64  `json:"user_id,omitempty"`
		Token  string `json:"token,omitempty"`
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()

	userMan := h.DB.User
	username := c.Query("username")
	inputPwd, ok := c.Get("hash_password")
	if !ok {
		resp.Status(StatusOtherError)
		return
	}
	password := inputPwd.(string)

	//用户名存在
	if ok := userMan.IsExistWithName(username); ok {
		resp.Status(StatusOtherError)
		return
	}
	//用户名为空
	if username == "" {
		resp.Status(StatusOtherError)
		return
	}
	//用户名长度超出限制
	if utf8.RuneCount([]byte(username)) > 32 {
		resp.Status(StatusOtherError)
		return
	}
	//保存到数据库
	user := model.User{
		Name:       username,
		Password:   password,
		Avatar:     PickOne(h.Avatars),
		Background: PickOne(h.Backgrounds),
	}
	if err := userMan.AddUser(&user); err != nil {
		resp.Status(StatusOtherError)
		return
	}
	//生成user_id和token
	userId := user.ID
	token, err := ginx.GenerateToken(userId)
	if err != nil {
		resp.Status(StatusOtherError)
		return
	}
	//注册成功，包装resp
	resp.UserId = userId
	resp.Token = token
}
