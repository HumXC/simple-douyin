package model

import (
	"github.com/gin-gonic/gin"
	"path"
	"testing"
)

func TestFavorite(t *testing.T) {
	douyinDB, err := NewDouyinDB(path.Join("douyin.db"))
	var c *gin.Context
	if err != nil {
		t.Fatal(err)
	}
	ThumbsUpMan := douyinDB.ThumbsUp
	thumbsUp1 := &ThumbsUp{
		UserId:     1,
		VideoId:    1,
		ActionType: 2,
	}
	thumbsUp2 := &ThumbsUp{
		UserId:     2,
		VideoId:    2,
		ActionType: 1,
	}
	t.Run("addAction", func(t *testing.T) {
		err := ThumbsUpMan.ActionTypeChange(c, int(thumbsUp1.UserId), int(thumbsUp1.VideoId))
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("QueryCommentById", func(t *testing.T) {
		err := ThumbsUpMan.ActionTypeAdd(c, int(thumbsUp2.UserId), int(thumbsUp2.VideoId))
		if err != nil {
			t.Fatal(err)
		}
	})
}
