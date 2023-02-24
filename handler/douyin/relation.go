package douyin

import (
	"net/http"
	"strconv"

	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RelationAction(c *gin.Context) {
	resp := BaseResponse()
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()

	inputId, ok := c.Get("user_id")
	if !ok {
		resp.Status(StatusAuthFailed)
		return
	}
	userId := inputId.(int64)
	//解析需要关注的id
	followId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		resp.Status(StatusOtherError)
		return
	}
	//解析action_type
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil {
		resp.Status(StatusOtherError)
		return
	}

	userMan := h.DB.User
	//关注用户不存在
	if !userMan.IsExistWithId(followId) {
		resp.Status(StatusUserNotFound)
		return
	}
	//未定义操作 1-关注，2-取消关注
	if actionType != 1 && actionType != 2 {
		resp.Status(StatusOtherError)
		return
	}
	//自己不能关注自己
	if userId == followId {
		resp.Status(StatusCanNotLikeSelf)
		return
	}

	switch actionType {
	case 1:
		err := userMan.Follow(userId, followId)
		if err != nil {
			resp.Status(StatusOtherError)
			return
		}
	case 2:
		err := userMan.CancelFollow(userId, followId)
		if err != nil {
			resp.Status(StatusOtherError)
			return
		}
	}
}

func (h *Handler) FollowList(c *gin.Context) {
	type Resp struct {
		Response
		UserList []User `json:"user_list"`
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		resp.Status(StatusInvalidParams)
		return
	}
	userMan := h.DB.User
	//用户不存在
	if !userMan.IsExistWithId(userId) {
		resp.Status(StatusUserNotFound)
		return
	}
	//在数据库查询关注的用户信息
	follows := userMan.FollowList(userId)
	userList := h.ConvertUsers(follows, false)
	for i := 0; i < len(*userList); i++ {
		(*userList)[i].IsFollow = h.DB.User.IsFollow(userId, (*userList)[i].Id)
	}
	resp.UserList = *userList
}

func (h *Handler) FollowerList(c *gin.Context) {
	type Resp struct {
		Response
		UserList []User `json:"user_list"`
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		resp.Status(StatusInvalidParams)
		return
	}
	userMan := h.DB.User
	//用户不存在
	if !userMan.IsExistWithId(userId) {
		resp.Status(StatusUserNotFound)
		return
	}
	//在数据库查询粉丝信息
	followers := userMan.FollowerList(userId)
	userList := h.ConvertUsers(followers, false)
	for i := 0; i < len(*userList); i++ {
		(*userList)[i].IsFollow = h.DB.User.IsFollow(userId, (*userList)[i].Id)
	}
	resp.UserList = *userList
}

func (h *Handler) FriendList(c *gin.Context) {
	type Resp struct {
		Response
		UserList []User `json:"user_list"`
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()

	inputId, ok := c.Get("user_id")
	if !ok {
		resp.Status(StatusAuthFailed)
		return
	}
	userId := inputId.(int64)

	userMan := h.DB.User
	//用户不存在
	if !userMan.IsExistWithId(userId) {
		resp.Status(StatusUserNotFound)
		return
	}
	//在数据库查询朋友信息
	userList, err := h.friends(userId)
	if err != nil {
		resp.Status(StatusOtherError)
		return
	}
	resp.UserList = userList
}

func (h *Handler) friends(id int64) ([]User, error) {
	friends := []model.User{}
	if err := h.DB.User.QueryFriendsById(id, &friends); err != nil {
		return nil, err
	}
	UserList := []User{}
	for _, v := range friends {
		message := model.Message{}
		h.DB.Message.QueryNewMsg(id, v.ID, &message)
		//该最新消息为当前请求用户发送的消息
		msgType := 0
		if message.FromUserId == id {
			msgType = 1
		}
		follower := User{
			Id:            v.ID,
			Name:          v.Name,
			FollowCount:   h.DB.User.CountFollow(v.ID),
			FollowerCount: h.DB.User.CountFollower(v.ID),
			IsFollow:      true,
			Avatar:        h.StorageClient.GetURL("avatars", v.Avatar),
			Background:    h.StorageClient.GetURL("backgrounds", v.Background),
			Message:       message.Content,
			MsgType:       int64(msgType),
		}
		UserList = append(UserList, follower)
	}
	return UserList, nil
}
