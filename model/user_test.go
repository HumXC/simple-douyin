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
