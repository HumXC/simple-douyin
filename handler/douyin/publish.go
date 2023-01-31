package douyin

import (
	"io"
	"net/http"

	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

type UploadFunc = func(r io.Reader, dir string) (string, error)

func (h *Handler) PublishAction(c *gin.Context) {
	resp := Response{
		StatusMsg: "投稿成功",
	}
	// TODO: 通过 token 获取用户 id
	// TODO: 日志
	title := c.PostForm("title")
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	defer f.Close()

	// 保存文件
	hash, err := h.UploadFunc(f, "videos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// 将视频信息写入数据库
	err = h.DB.Video.Put(model.Video{
		Hash:  hash,
		Title: title,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}
