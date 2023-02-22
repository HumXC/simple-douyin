package douyin

import (
	"net/http"
	"strconv"

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
		videos := h.DB.Video.GetFeed(latestTime, num)
		userID := c.GetInt64("user_id")
		resp.VideoList = *h.ConvertVideos(&videos, userID, nil)
		if len(videos) != 0 {
			resp.NextTime = videos[len(videos)-1].Time.Unix()
		}
	}
}
