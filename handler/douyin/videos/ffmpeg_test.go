package videos_test

import (
	"os"
	"path"
	"testing"

	"github.com/HumXC/simple-douyin/handler/douyin/videos"
)

var testVideo = path.Join(TestDir, "video.mp4")

func TestCutVideoWithFfmpeg(t *testing.T) {
	// 没报错就是成功
	output, err := videos.CutVideoWithFfmpeg(testVideo)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(output)
}

func BenchmarkCutVideoWithFfmpeg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		output, err := videos.CutVideoWithFfmpeg(testVideo)
		if err != nil {
			b.Fatal(err)
		}
		defer os.Remove(output)
	}
}

func TestSmallVideoWithFfmpeg(t *testing.T) {
	// 没报错就是成功
	output, err := videos.SmallVideoWithFfmpeg(testVideo)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(output)
}

func BenchmarkSmallVideoWithFfmpeg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		output, err := videos.SmallVideoWithFfmpeg(testVideo)
		if err != nil {
			b.Fatal(err)
		}
		defer os.Remove(output)
	}
}
