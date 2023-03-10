package sqldb_test

import (
	"path"
	"testing"

	"github.com/HumXC/simple-douyin/database/sqldb"
	"github.com/HumXC/simple-douyin/model"
)

func TestDouyinPutAndGetByID(t *testing.T) {
	douyinDB, err := sqldb.NewDouyinDB("sqlite", path.Join(TEST_DIR, "douyin.db"))
	if err != nil {
		t.Fatal(err)
	}
	videoMan := douyinDB.Video
	newVideo := model.Video{
		UserID:        1,
		Video:         "testvideo",
		Title:         "这是一个测试视频",
		CommentCount:  0,
		FavoriteCount: 0,
	}

	t.Run("Put", func(t *testing.T) {
		videoMan.Put(newVideo)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("GetByID", func(t *testing.T) {
		v := videoMan.GetByID(1)
		if err != nil {
			t.Fatal(err)
		}
		if v.Video != newVideo.Video {
			t.Fatal("视频不匹配")
		}
	})
}
