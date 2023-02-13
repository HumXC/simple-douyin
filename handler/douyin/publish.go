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

type UploadFunc = func(fileName, dir string) (string, error)

var buf = sync.Pool{
	New: func() any {
		return make([]byte, 512)
	},
}

func (h *Handler) PublishAction(c *gin.Context) {
	var httpStatusCode = http.StatusOK
	resp := Response{
		StatusMsg:  "投稿成功",
		StatusCode: 0,
	}
	defer func() {
		c.JSON(httpStatusCode, resp)
	}()
	// TODO: 通过 token 获取用户 id
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
	vHash, err := h.UploadFunc(file.Name(), "videos")
	if err != nil {
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "视频上传失败"
		return
	}
	cHash, err := h.UploadFunc(cover, "covers")
	if err != nil {
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "封面上传失败"
		return
	}

	// 将视频信息写入数据库
	err = h.DB.Video.Put(model.Video{
		Video: vHash,
		Cover: cHash,
		Title: title,
	})
	if err != nil {
		resp.StatusCode = StatusOtherError
		resp.StatusMsg = "服务器错误"
		return
	}
}
