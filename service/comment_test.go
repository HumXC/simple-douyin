package service_test

import (
	"fmt"
	"github.com/HumXC/simple-douyin/model"
	"github.com/HumXC/simple-douyin/service"
	"testing"
)

/**
 * @Description
 * @Author xyc
 * @Date 2023/1/28 20:17
 **/

func TestAddComment(t *testing.T) {
	comment := &model.Comment{
		UserID:  1,
		VideoId: 1,
		Content: "测试评论3",
	}
	commentData, err := service.AddComment(comment)
	if err != nil {
		fmt.Println("AddComment ERR:", err)
		return
	}
	fmt.Println(commentData)
}
