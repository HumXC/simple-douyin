package douyin

import (
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/HumXC/simple-douyin/helper"
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

type StorageClient interface {
	// 上传一个文件到存储
	// TODO 将 fileName 改成 io.Reader 实现
	Upload(fileName, dir string) (string, error)
	GetURLWithHash(dir, hash string) string
}

var buf = sync.Pool{
	New: func() any {
		return make([]byte, 512)
	},
}

func (h *Handler) PublishAction(c *gin.Context) {
	var httpStatusCode = http.StatusOK
	resp := Response{
		StatusMsg:  "投稿成功",
		StatusCode: StatusOK,
	}
	defer func() {
		c.JSON(httpStatusCode, resp)
	}()
	userID := c.GetInt64("user_id")
	// TODO: 日志
	title := c.PostForm("title")
	data, _, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	defer data.Close()
	// 保存文件到本地先，file 作为临时文件，最后要删除
	file, err := os.CreateTemp("", "video")
	if err != nil {
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "服务器错误"
		return
	}
	defer file.Close()
	defer os.Remove(file.Name())

	// 将 data 写到 file
	_, _ = io.Copy(file, data)
	_, _ = file.Seek(0, io.SeekStart)
	// 判断文件是否为视频
	b := buf.Get().([]byte)
	_, _ = file.Read(b)
	mimeType := http.DetectContentType(b)
	buf.Put(b)
	// 上传的文件不是视频
	if !strings.HasPrefix(mimeType, "video") {
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "上传的文件不是视频"
		return
	}

	// 获取视频封面
	cover, err := helper.CutVideoWithFfmpeg(file.Name())
	if err != nil {
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "无法获取视频封面"
		return
	}
	defer os.Remove(cover)

	// 上传视频和封面
	vHash, err := h.StorageClient.Upload(file.Name(), "videos")
	if err != nil {
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "视频上传失败"
		return
	}
	cHash, err := h.StorageClient.Upload(cover, "covers")
	if err != nil {
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "封面上传失败"
		return
	}

	// 将视频信息写入数据库
	err = h.DB.Video.Put(model.Video{
		Video:  vHash,
		Cover:  cHash,
		Title:  title,
		UserID: userID,
	})
	if err != nil {
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "服务器错误"
		return
	}
}

func (h *Handler) PublishList(c *gin.Context) {
	type Resp struct {
		Response
		VideoList []Video `json:"video_list"`
	}
	var httpStatusCode = http.StatusOK
	resp := Resp{
		Response: Response{
			StatusMsg:  "成功",
			StatusCode: StatusOK,
		},
	}
	defer func() {
		c.JSON(httpStatusCode, resp)
	}()
	userID := c.GetInt64("user_id")
	user, err := h.user(userID)
	if err != nil {
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "未知错误"
		httpStatusCode = http.StatusInternalServerError
		return
	}
	if user.Id == 0 {
		resp.StatusCode = StatusUserNotFound
		resp.StatusMsg = "用户不存在"
		return
	}

	videos, err := h.DB.Video.GetByUser(userID)
	if err != nil {
		httpStatusCode = http.StatusInternalServerError
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "无法获取视频"
	}
	resp.VideoList = make([]Video, len(videos))
	for i := 0; i < len(videos); i++ {
		resp.VideoList[i].Author = user
		resp.VideoList[i].CommentCount = videos[i].CommentCount
		resp.VideoList[i].FavoriteCount = videos[i].FavoriteCount
		// FIXME 获取正确的 IsFavorite
		resp.VideoList[i].IsFavorite = false
		resp.VideoList[i].Id = videos[i].ID
		resp.VideoList[i].CoverUrl = h.StorageClient.GetURLWithHash("covers", videos[i].Cover)
		resp.VideoList[i].PlayUrl = h.StorageClient.GetURLWithHash("videos", videos[i].Video)
	}

}
