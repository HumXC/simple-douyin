package douyin

import (
	"net/http"
	"strconv"

	"github.com/HumXC/simple-douyin/helper"
	"github.com/gin-gonic/gin"
)

type ActionRequest struct {
	Token      string `json:"token"`
	VideoId    int64  `json:"video_id"`
	ActionType int32  `json:"action_type"`
}

type ListResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	VideoList  []VideoList `json:"video_list"`
}

type VideoList struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

// Action 赞操作
func (h *Handler) Action(c *gin.Context) {

	videoId, _ := strconv.Atoi(c.Query("video_id"))
	actionType, _ := strconv.Atoi(c.Query("action_type"))
	token := c.Query("token")
	//解析toke
	userClaim, _ := helper.AnalyseToken(token)
	userId := userClaim.UserId
	//
	action := h.DB.ThumbsUp
	if videoId == 0 || actionType == 0 {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: -1,
			StatusMsg:  "InvalidParams",
		})
		return
	}
	//取消点赞操作
	if actionType == 2 {
		err := action.ActionTypeChange(c, videoId, int(userId))
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				StatusCode: 403,
				StatusMsg:  err.Error(),
			})
			return
		}
	}
	//点赞
	err := action.ActionTypeAdd(c, videoId, int(userId))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 403,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 200,
		StatusMsg:  "ok",
	})
	return
}

// List 喜欢列表
func (h *Handler) List(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	if userId == 0 {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: -1,
			StatusMsg:  "InvalidParams",
		})
		return
	}

}
