package douyin

import (
	"github.com/gin-gonic/gin"
)

/**
 * @Description 评论操作-控制层
 * @Author xyc
 * @Date 2023/1/27 11:53
 **/

const (
	CREATE = 1 //添加评论
	DELETE = 2 //删除评论
)

type CommentActionRequest struct {
	Token       string `json:"token"`        //用户鉴权token
	VideoId     int64  `json:"video_id"`     //评论的视频ID
	ActionType  int32  `json:"action_type"`  //操作类型  1-发布评论，2-删除评论
	CommentText string `json:"comment_text"` //用户填写的评论内容，在action_type=1的时候使用
	CommentId   string `json:"comment_id"`   //要删除的评论id，在action_type=2的时候使用
}

type CommentListRequest struct {
	Token   string `json:"token"`    //用户鉴权token
	VideoId int64  `json:"video_id"` //评论的视频ID
}

type CommentActionResponse struct {
	Response         //通用字段
	Comment  Comment `json:"comment,omitempty"` //评论信息
}

type CommentListResponse struct {
	Response              //通用字段
	CommentList []Comment `json:"comment_list,omitempty"` //评论集合
}

func CommentAction(c *gin.Context) {

}

func CommentList(c *gin.Context) {

}
