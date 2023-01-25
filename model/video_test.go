package model_test

import (
	"testing"

	"github.com/HumXC/simple-douyin/model"
)

func TestPutAndGetByID(t *testing.T) {
	videoMan := model.VideoMan
	newVideo := model.Video{
		UserID: "111222333",
		Hash:   "testvideo",
		Title:  "这是一个测试视频",
	}

	t.Run("Put", func(t *testing.T) {
		err := videoMan.Put(newVideo)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("GetByID", func(t *testing.T) {
		v, err := videoMan.GetByID("1")
		if err != nil {
			t.Fatal(err)
		}
		if v.Hash != newVideo.Hash {
			t.Fatal("视频不匹配")
		}
	})
}
