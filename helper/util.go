package helper

import (
	"bytes"
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
)

// 调用系统的 ffmpeg 截取视频的第 1 秒的帧输出为一张图片
// output 参数应该明确后缀名 (例如 .jpg), 否则 ffmpeg 会报错
// 示例: CutVideoWithFfmpeg("a.mp4", "output.jpg")
func CutVideoWithFfmpeg(video string) (output string, err error) {
	// https://trac.ffmpeg.org/wiki/Create%20a%20thumbnail%20image%20every%20X%20seconds%20of%20the%20video
	// ffmpeg -i input.flv -ss 00:00:14.435 -frames:v 1 out.png
	output = path.Join(os.TempDir(), "ffmpeg-"+strconv.Itoa(rand.Int())+".jpg")
	c := exec.Command("ffmpeg", "-i", video, "-ss", "00:00:1", "-vcodec", "mjpeg", "-frames:v", "1", output, "-y")
	var stdErr bytes.Buffer
	c.Stderr = &stdErr
	_, err = c.Output()
	if err != nil {
		err = errors.New(stdErr.String())
		return
	}
	return
}
