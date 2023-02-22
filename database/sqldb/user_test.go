package sqldb_test

import (
	"path"
	"testing"

	"github.com/HumXC/simple-douyin/database/sqldb"
	"github.com/HumXC/simple-douyin/model"
)

func TestIsFollow(t *testing.T) {
	douyinDB, err := sqldb.NewDouyinDB("sqlite", path.Join(TEST_DIR, "douyin.db"))
	if err != nil {
		t.Fatal(err)
	}
	u := douyinDB.User
	u.AddUser(&model.User{})
	u.AddUser(&model.User{})
	u.AddUser(&model.User{})
	u.Follow(1, 2)
	u.Follow(1, 3)
	list := [][]int64{
		{1, 2},
		{1, 3},
		{2, 1},
		{1, 0},
	}
	want := []bool{true, true, false, false}
	for i, w := range want {
		g := u.IsFollow(list[i][0], list[i][1])
		if w != g {
			t.Errorf("关系不正确: got: %t want:  %t", g, w)
		}

	}
}
func TestCountFollowAndCountFollower(t *testing.T) {
	douyinDB, err := sqldb.NewDouyinDB("sqlite", path.Join(TEST_DIR, "douyin.db"))
	if err != nil {
		t.Fatal(err)
	}
	u := douyinDB.User
	u.AddUser(&model.User{})
	u.AddUser(&model.User{})
	u.AddUser(&model.User{})
	u.Follow(1, 2)
	u.Follow(1, 3)
	u.Follow(2, 1)
	u.Follow(3, 2)
	t.Run("关注数", func(t *testing.T) {
		var want int64 = 2
		got := u.CountFollow(1)
		if want != got {
			t.Errorf("关注数不正确: got: %d want:  %d", got, want)
		}
	})
	t.Run("粉丝数", func(t *testing.T) {
		var want int64 = 1
		got := u.CountFollower(1)
		if want != got {
			t.Errorf("粉丝数不正确: got: %d want:  %d", got, want)
		}
		want = 0
		got = u.CountFollower(5)
		if want != got {
			t.Errorf("粉丝数不正确: got: %d want:  %d", got, want)
		}
	})
	t.Run("关注者", func(t *testing.T) {
		var want []model.User = []model.User{
			{ID: 2},
			{ID: 3},
		}
		got := *u.QueryFollows(1)
		if len(got) == 0 {
			t.Errorf("居然找不到关注者: want:  %v", want)
			return
		}
		for _, g := range got {
			ok := false
			for _, w := range want {
				if w.ID == g.ID {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("没有找到关注者: got: %v want:  %v", g, want)
			}
		}
	})
	t.Run("粉丝", func(t *testing.T) {
		var want []model.User = []model.User{
			{ID: 1},
			{ID: 3},
		}

		got := *u.QueryFollowers(2)
		if len(got) == 0 {
			t.Errorf("居然找不到粉丝: want:  %v", want)
			return
		}
		for _, g := range got {
			ok := false
			for _, w := range want {
				if w.ID == g.ID {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("没有找到粉丝: got: %v want:  %v", g, want)
			}
		}
	})
}
