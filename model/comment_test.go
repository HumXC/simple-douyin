package model_test

import (
	"github.com/HumXC/simple-douyin/model"
	"path"
	"testing"
)

/**
 * @Description
 * @Author xyc
 * @Date 2023/1/26 21:31
 **/

func TestComment(t *testing.T) {
	douyinDB, err := model.NewDouyinDB(path.Join(TEST_DIR, "douyin.db"))
	if err != nil {
		t.Fatal(err)
	}
	commentMan := douyinDB.Comment
	comment := &model.Comment{
		UserID:  1,
		VideoId: 1,
		Content: "测试评论1",
	}

	t.Run("添加评论测试", func(t *testing.T) {
		err := commentMan.AddComment(comment)
		if err != nil {
			t.Fatal(err)
		}
	})
}
