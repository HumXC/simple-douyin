package douyin

import (
	"github.com/HumXC/simple-douyin/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

// action 赞操作
func Action(c *gin.Context) {
	videoId, _ := strconv.Atoi(c.PostForm("video_id"))
	actionType, _ := strconv.Atoi(c.PostForm("action_type"))
	//解析toke
	claims, _ := jwt.ParseWithClaims()
	userId := claims.Valid
	//
	var action model.ThumbsUpMan
	if videoId == 0 || actionType == 0 {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: -1,
			StatusMsg:  "InvalidParams",
		})
		return
	}
	//取消点赞操作
	if actionType == 2 {
		err := action.ActionTypeChange(c, videoId, userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				StatusCode: 403,
				StatusMsg:  err.Error(),
			})
			return
		}
	}
	//点赞
	err := action.ActionTypeAdd(c, videoId, userId)
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

// list 喜欢列表
func List(c *gin.Context) {

}
