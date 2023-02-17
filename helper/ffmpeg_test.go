package helper_test

import (
	"os"
	"testing"

	"github.com/HumXC/simple-douyin/helper"
)

const testVideo = "../test/video.mp4"

func TestCutVideoWithFfmpeg(t *testing.T) {
	// 没报错就是成功
	output, err := helper.CutVideoWithFfmpeg(testVideo)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(output)
}

func BenchmarkCutVideoWithFfmpeg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		output, err := helper.CutVideoWithFfmpeg(testVideo)
		if err != nil {
			b.Fatal(err)
		}
		defer os.Remove(output)
	}
}

func TestSmallVideoWithFfmpeg(t *testing.T) {
	// 没报错就是成功
	output, err := helper.SmallVideoWithFfmpeg(testVideo)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(output)
}

func BenchmarkSmallVideoWithFfmpeg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		output, err := helper.SmallVideoWithFfmpeg(testVideo)
		if err != nil {
			b.Fatal(err)
		}
		defer os.Remove(output)
	}
}
