package sqldb_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/HumXC/simple-douyin/database/sqldb"
	"github.com/HumXC/simple-douyin/model"
)

/**
 * @Description
 * @Author xyc
 * @Date 2023/1/26 21:31
 **/

func TestComment(t *testing.T) {
	douyinDB, err := sqldb.NewDouyinDB("sqlite", path.Join(TEST_DIR, "douyin.db"))
	if err != nil {
		t.Fatal(err)
	}
	commentMan := douyinDB.Comment
	comment1 := &model.Comment{
		UserID:  1,
		VideoId: 1,
		Content: "测试评论1",
	}
	comment2 := &model.Comment{
		UserID:  2,
		VideoId: 1,
		Content: "测试评论2",
	}

	t.Run("AddComment", func(t *testing.T) {
		err := commentMan.AddComment(comment1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("QueryCommentById", func(t *testing.T) {
		var comment model.Comment
		err := commentMan.QueryCommentById(2, &comment)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(comment)
	})

	t.Run("QueryCommentListByVideoId", func(t *testing.T) {
		var comments []model.Comment
		err := commentMan.QueryCommentListByVideoId(1, &comments)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(comments)
	})

	t.Run("AddCommentAndUpdateCommentCount", func(t *testing.T) {
		err := commentMan.AddCommentAndUpdateCommentCount(comment1)
		err = commentMan.AddCommentAndUpdateCommentCount(comment2)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DeleteCommentAndUpdateCountById", func(t *testing.T) {
		err := commentMan.DeleteCommentAndUpdateCountById(1, 1)
		if err != nil {
			t.Fatal(err)
		}
	})
}
