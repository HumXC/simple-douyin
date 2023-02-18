package model_test

import (
	"path"
	"testing"

	"github.com/HumXC/simple-douyin/model"
)

func TestIsFollow(t *testing.T) {
	douyinDB, err := model.NewDouyinDB("sqlite", path.Join(TEST_DIR, "douyin.db"), nil)
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
	douyinDB, err := model.NewDouyinDB("sqlite", path.Join(TEST_DIR, "douyin.db"), nil)
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
			{Id: 2},
			{Id: 3},
		}
		var got []model.User = make([]model.User, 0)
		u.QueryFollows(1, &got)
		if len(got) == 0 {
			t.Errorf("居然找不到关注者: want:  %v", want)
			return
		}
		for _, g := range got {
			ok := false
			for _, w := range want {
				if w.Id == g.Id {
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
			{Id: 1},
			{Id: 3},
		}
		var got []model.User = make([]model.User, 0)
		u.QueryFollowers(2, &got)
		if len(got) == 0 {
			t.Errorf("居然找不到粉丝: want:  %v", want)
			return
		}
		for _, g := range got {
			ok := false
			for _, w := range want {
				if w.Id == g.Id {
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