package douyin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var DemoVideos = []Video{
	{
		Id:            1,
		Author:        User{},
		PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

func (h *Handler) Feed(num int) func(*gin.Context) {
	return func(c *gin.Context) {
		type Resp struct {
			Response
			VideoList []Video `json:"video_list,omitempty"`
			NextTime  int64   `json:"next_time,omitempty"`
		}

		var httpStatusCode = http.StatusOK
		resp := Resp{
			Response: Response{
				StatusMsg:  "OK",
				StatusCode: StatusOK,
			},
		}
		defer func() {
			c.JSON(httpStatusCode, resp)
		}()
		latestTime, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
		if err != nil {
			resp.StatusCode = StatusOtherError
			resp.StatusMsg = "未知错误"
			httpStatusCode = http.StatusInternalServerError
			return
		}
		videos, err := h.DB.Video.GetFeed(latestTime, num)
		if err != nil {
			resp.StatusCode = StatusOtherError
			resp.StatusMsg = "未知错误"
			httpStatusCode = http.StatusInternalServerError
			return
		}

		resp.VideoList = make([]Video, len(videos))
		for i := 0; i < len(videos); i++ {
			user, err := h.user(videos[i].UserID)
			if err != nil {
				resp.StatusCode = StatusOtherError
				resp.StatusMsg = "未知错误"
				httpStatusCode = http.StatusInternalServerError
				return
			}
			resp.VideoList[i].Author = user
			resp.VideoList[i].CommentCount = videos[i].CommentCount
			resp.VideoList[i].FavoriteCount = videos[i].FavoriteCount
			// FIXME 获取正确的 IsFavorite
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
