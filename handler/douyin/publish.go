package douyin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/HumXC/simple-douyin/handler/douyin/videos"
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

type StorageClient interface {
	// 上传一个文件到存储
	// TODO 将 fileName 改成 io.Reader 实现
	Upload(fileName, dir string) (string, error)
	GetURL(dir, file string) string
}

// 服务启动后开始异步压缩未处理的视频
func VideoButcherFinishFunc(h *Handler) videos.ButcherFinidhFunc {
	return func(job videos.Job, video, cover string, err error) (delete bool) {
		if err != nil {
			fmt.Println("视频任务失败: " + err.Error())
			return false
		}
		vHash, err := h.StorageClient.Upload(video, "videos")
		if err != nil {
			fmt.Println("视频任务失败: " + err.Error())
			return false
		}
		cHash, err := h.StorageClient.Upload(cover, "covers")
		if err != nil {
			fmt.Println("视频任务失败: " + err.Error())
			return false
		}
		// 将视频信息写入数据库
		h.DB.Video.Put(model.Video{
			Video:  vHash,
			Cover:  cHash,
			Title:  job.Title,
			UserID: job.UserID,
			Time:   time.Now(),
		})
		return true
	}
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
	if u.ID == 0 {
		resp.Status(StatusUserNotFound)
		return
	}
	user := h.ConvertUser(u, false)
	videos := h.DB.Video.GetByUser(userID)
	resp.VideoList = *h.ConvertVideos(&videos, 0, &user)
}
