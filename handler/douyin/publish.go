package douyin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

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
	resp := BaseResponse()
	userID := c.GetInt64("user_id")
	title := c.PostForm("title")
	data, _, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	defer data.Close()
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()
	// 保存文件到本地先，file 作为临时文件，最后要删除
	file, err := os.CreateTemp("", "video")
	if err != nil {
		resp.Status(StatusOtherError)
		panic(fmt.Errorf("无法创建临时文件: %w", err))
	}
	defer file.Close()

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
		resp.Status(StatusUploadNotAVideo)
		return
	}
	// 走到这就已经返回 “发布成功了，其实还有额外的工作在 Butcher 进行”
	h.VideoButcher.Add(file.Name(), title, userID)
	return
}

func (h *Handler) PublishList(c *gin.Context) {
	type Resp struct {
		Response
		VideoList []Video `json:"video_list"`
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()
	userID := c.GetInt64("user_id")
	u := model.User{}
	err := h.DB.User.QueryById(userID, &u)
	if err != nil {
		resp.Status(StatusOtherError)
		panic(fmt.Errorf("无法获取用户 [%d] : %w", userID, err))
	}
	if u.Id == 0 {
		resp.Status(StatusUserNotFound)
		return
	}
	user := h.ConvertUser(u, false)

	videos, err := h.DB.Video.GetByUser(userID)
	if err != nil {
		resp.Status(StatusOtherError)
		panic(fmt.Errorf("无法获取用户发布的视频 [%d] : %w", userID, err))
	}
	resp.VideoList = make([]Video, len(videos))
	for i := 0; i < len(videos); i++ {
		resp.VideoList[i].Author = user
		resp.VideoList[i].CommentCount = videos[i].CommentCount
		resp.VideoList[i].FavoriteCount = videos[i].FavoriteCount
		resp.VideoList[i].IsFavorite = false
		resp.VideoList[i].Id = videos[i].ID
		resp.VideoList[i].CoverUrl = h.StorageClient.GetURLWithHash("covers", videos[i].Cover)
		resp.VideoList[i].PlayUrl = h.StorageClient.GetURLWithHash("videos", videos[i].Video)
	}

}
