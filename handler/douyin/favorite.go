package douyin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ActionRequest struct {
	Token      string `json:"token"`
	VideoId    int64  `json:"video_id"`
	ActionType int32  `json:"action_type"`
}

type ActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type ListResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	VideoList  []ListVideo `json:"video_list"`
}

type ListVideo struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

type ListUser struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

func ToThumbsUp(c *gin.Context) {
	videoId := c.PostForm("video_id")
	actionType := c.PostForm("action_type")
	if videoId == "" || actionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "InvalidParams",
		})
		return
	}
}
