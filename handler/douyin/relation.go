package douyin

import (
	"net/http"
	"strconv"

	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

type FollowListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

func (h *Handler) RelationAction(c *gin.Context) {
	inputId, ok := c.Get("user_id")
	if !ok {
		CommonResponseError(c, "user_id解析失败")
		return
	}
	userId := inputId.(int64)
	//解析需要关注的id
	followId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		CommonResponseError(c, "to_user_id解析失败")
		return
	}
	//解析action_type
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil {
		CommonResponseError(c, "action_type解析失败")
		return
	}

	userMan := h.DB.User
	//关注用户不存在
	if !userMan.IsExistWithId(followId) {
		CommonResponseError(c, "关注用户不存在")
		return
	}
	//未定义操作 1-关注，2-取消关注
	if actionType != 1 && actionType != 2 {
		CommonResponseError(c, "未定义操作")
		return
	}
	//自己不能关注自己
	if userId == followId {
		CommonResponseError(c, "自己不能关注自己")
		return
	}

	switch actionType {
	case 1:
		err := userMan.Follow(userId, followId)
		if err != nil {
			CommonResponseError(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "关注成功",
		})
	case 2:
		err := userMan.CancelFollow(userId, followId)
		if err != nil {
			CommonResponseError(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "取关成功",
		})
	}
}

func (h *Handler) FollowList(c *gin.Context) {
	userMan := h.DB.User

	inputId, ok := c.Get("user_id")
	if !ok {
		CommonResponseError(c, "user_id解析失败")
		return
	}
	userId := inputId.(int64)

	//用户不存在
	if !userMan.IsExistWithId(userId) {
		CommonResponseError(c, "用户不存在")
		return
	}
	//在数据库查询关注的用户信息
	follows := userMan.QueryFollows(userId)
	userList := h.ConvertUsers(follows, true)
	c.JSON(0, FollowListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		UserList: *userList,
	})
}

func (h *Handler) FollowerList(c *gin.Context) {
	userMan := h.DB.User

	inputId, ok := c.Get("user_id")
	if !ok {
		CommonResponseError(c, "user_id解析失败")
		return
	}
	userId := inputId.(int64)

	//用户不存在
	if !userMan.IsExistWithId(userId) {
		CommonResponseError(c, "用户不存在")
		return
	}
	//在数据库查询粉丝信息
	followers := userMan.QueryFollowers(userId)
	userList := h.ConvertUsers(followers, true)
	c.JSON(0, FollowListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		UserList: *userList,
	})
}

func (h *Handler) FriendList(c *gin.Context) {
	userMan := h.DB.User

	inputId, ok := c.Get("user_id")
	if !ok {
		CommonResponseError(c, "user_id解析失败")
		return
	}
	userId := inputId.(int64)
	//用户不存在
	if !userMan.IsExistWithId(userId) {
		CommonResponseError(c, "用户不存在")
		return
	}
	//在数据库查询朋友信息
	userList, err := h.friends(userId)
	if err != nil {
		CommonResponseError(c, err.Error())
		return
	}

	c.JSON(0, FollowListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		UserList: userList,
	})

}

func (h *Handler) friends(id int64) ([]User, error) {
	friends := []model.User{}
	if err := h.DB.User.QueryFriendsById(id, &friends); err != nil {
		return nil, err
	}
	UserList := []User{}
	for _, v := range friends {
		follower := User{
			Id:            v.Id,
			Name:          v.Name,
			FollowCount:   h.DB.User.CountFollow(v.Id),
			FollowerCount: h.DB.User.CountFollower(v.Id),
			IsFollow:      true,
		}
		UserList = append(UserList, follower)
	}
	return UserList, nil
}
