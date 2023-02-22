package douyin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Favorite(c *gin.Context) {
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		panic(err)
	}
	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		panic(err)
	}
	userId := c.GetInt64("user_id")
	resp := BaseResponse()
	defer c.JSON(http.StatusOK, &resp)
	if videoId == 0 || actionType == 0 {
		resp.Status(InvalidParams)
		return
	}
	if actionType != 1 && actionType != 2 {
		resp.Status(InvalidParams)
		return
	}
	err = h.RDB.User.Favorite(userId, videoId, int32(actionType))
	if err != nil {
		resp.Status(StatusOtherError)
		panic(fmt.Errorf("喜欢/取消喜欢错误: %w", err))
	}
}

func (h *Handler) FavoriteList(c *gin.Context) {
	type Resp struct {
		Response
		VideoList []Video `json:"video_list"`
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	userID := c.GetInt64("user_id")
	vsID := h.RDB.User.FavoriteList(userID)
	vs := h.DB.Video.GetByIDs(vsID)
	resp.VideoList = *h.ConvertVideos(&vs, userID, nil)
	c.JSON(http.StatusOK, &resp)
}

// Deprecated Action 赞操作
func (h *Handler) Action(c *gin.Context) {
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	actionType, _ := strconv.Atoi(c.Query("action_type"))
	userId := c.GetInt64("user_id")
	action := h.DB.ThumbsUp
	resp := BaseResponse()
	defer c.JSON(http.StatusOK, &resp)
	if videoId == 0 || actionType == 0 {
		resp.Status(InvalidParams)
		return
	}
	//取消点赞操作
	if actionType == 2 {
		err := action.ActionTypeChange(c, videoId, int(userId))
		if err != nil {
			resp.Status(StatusOtherError)
			panic(fmt.Errorf("取消点赞错误: %w", err))
		}
		return
	}
	//点赞
	err := action.ActionTypeAdd(c, videoId, int(userId))
	if err != nil {
		resp.Status(StatusOtherError)
		panic(fmt.Errorf("点赞错误: %w", err))
	}
}
