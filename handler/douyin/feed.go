package douyin

import (
	"net/http"
	"strconv"

	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Feed(num int) func(*gin.Context) {
	return func(c *gin.Context) {
		type Resp struct {
			Response
			VideoList []Video `json:"video_list,omitempty"`
			NextTime  int64   `json:"next_time,omitempty"`
		}
		resp := Resp{
			Response: BaseResponse(),
		}
		defer func() {
			c.JSON(http.StatusOK, resp)
		}()
		latestTime, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
		if err != nil {
			resp.Status(StatusOtherError)
			return
		}
		videos, err := h.DB.Video.GetFeed(latestTime, num)
		if err != nil {
			resp.Status(StatusOtherError)
			return
		}
		userID := c.GetInt64("user_id")
		resp.VideoList = make([]Video, len(videos))
		for i := 0; i < len(videos); i++ {
			u := model.User{}
			err := h.DB.User.QueryById(videos[i].UserID, &u)
			if err != nil {
				resp.Status(StatusOtherError)
				return
			}
			user := h.ConvertUser(u, h.DB.User.IsFollow(userID, u.Id))
			resp.VideoList[i].Author = user
			resp.VideoList[i].CommentCount = videos[i].CommentCount
			resp.VideoList[i].FavoriteCount = videos[i].FavoriteCount
			// FIXME 获取正确的值
			resp.VideoList[i].IsFavorite = false
			resp.VideoList[i].Id = videos[i].ID
			resp.VideoList[i].CoverUrl = h.StorageClient.GetURLWithHash("covers", videos[i].Cover)
			resp.VideoList[i].PlayUrl = h.StorageClient.GetURLWithHash("videos", videos[i].Video)
		}
		if len(videos) != 0 {
			resp.NextTime = videos[len(videos)-1].Time.Unix()
		}
	}
}
