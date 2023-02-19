package douyin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Action 赞操作
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

// List 喜欢列表
func (h *Handler) List(c *gin.Context) {
	// user_id 可以直接 c.GetInt64() 获取，因为 /relation 接口添加了 NeedLogin 中间件，所以不用担心 user_id 为 0
	userId, _ := strconv.Atoi(c.Query("user_id"))
	if userId == 0 {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: -1,
			StatusMsg:  "InvalidParams",
		})
		return
	}

}
