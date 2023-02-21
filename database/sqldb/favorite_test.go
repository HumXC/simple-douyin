package sqldb_test

import (
	"path"
	"testing"

	"github.com/HumXC/simple-douyin/database/sqldb"
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

func TestThumbsUp(t *testing.T) {
	douyinDB, err := sqldb.NewDouyinDB("sqlite", path.Join(TEST_DIR, "douyin.db"))
	if err != nil {
		t.Fatal(err)
	}
	tb := douyinDB.ThumbsUp
	// 用户给视频点赞，user_id,video_id,action_type
	actions := [][]int64{
		{1, 1, 1},    // 视频[1] 被用户[1]点赞
		{2, 1, 1},    // 视频[2] 被用户[1]点赞
		{1, 1, 2},    // 视频[1] 被用户[1]取消点赞
		{1, 2, 2},    // 视频[1] 被用户[2]取消点赞
		{3, 2, 0},    // 视频[3] 不知道被用户[2]干了什么
		{-1, -1, -1}, // 视频[-1?] 不知道被用户[-1?]干了什么
	}
	for _, v := range actions {
		vid := v[0]
		uid := v[1]
		act := int32(v[2])
		err = tb.Action(vid, uid, act)
		if err != nil {
			t.Errorf("点赞数据库异常: %s", err)
		}
	}

	// 不想写数据库测试了，自己拿工具查数据库对数据
	_ = [][]int64{
		{1, 1, 2},    // 视频[1] 被用户[1]取消点赞
		{1, 2, 2},    // 视频[1] 被用户[2]取消点赞
		{2, 1, 1},    // 视频[2] 被用户[1]点赞
		{3, 2, 0},    // 视频[3] 不知道被用户[2]干了什么
		{-1, -1, -1}, // 视频[-1?] 不知道被用户[-1?]干了什么
	}

}
func TestFavorite(t *testing.T) {
	// FIXME 传入 Redis 实例
	douyinDB, err := sqldb.NewDouyinDB("sqlite", path.Join(TEST_DIR, "douyin.db"))
	var c *gin.Context
	if err != nil {
		t.Fatal(err)
	}
	ThumbsUpMan := douyinDB.ThumbsUp
	thumbsUp1 := &model.ThumbsUp{
		UserId:     1,
		VideoId:    1,
		ActionType: 2,
	}
	thumbsUp2 := &model.ThumbsUp{
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
