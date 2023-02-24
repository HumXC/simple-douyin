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
		resp.Status(StatusInvalidParams)
		return
	}
	if actionType != 1 && actionType != 2 {
		resp.Status(StatusInvalidParams)
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
