package douyin

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
)

/**
 * @Description 评论操作-控制层
 * @Author xyc
 * @Date 2023/1/27 11:53
 **/

func (h *Handler) CommentAction(c *gin.Context) {
	type Resp struct {
		Response         //通用字段
		Comment  Comment `json:"comment,omitempty"` //评论信息
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()

	commentMan := h.DB.Comment
	userMan := h.DB.User
	userId := c.GetInt64("user_id")
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	actionType, _ := strconv.Atoi(c.Query("action_type")) //操作类型 1-发布评论，2-删除评论

	//1-发布评论
	if actionType == 1 {
		commentText := c.Query("comment_text")
		var comment model.Comment
		comment.UserID = userId
		comment.VideoId = int64(videoId)
		comment.Content = commentText
		err := commentMan.AddCommentAndUpdateCommentCount(&comment)
		//发布评论失败
		if err != nil {
			resp.Status(StatusFailedPostComment)
			log.Println("发布评论失败:", err.Error())
			return
		}
		//发布评论成功
		var u model.User
		_ = userMan.QueryById(userId, &u)
		userInfo := h.ConvertUser(u, false)
		commentData := Comment{
			Id:         int64(comment.Model.ID),
			User:       userInfo,
			Content:    comment.Content,
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
		}
		resp.Comment = commentData //发布评论成功
		return
	} else {
		//2-删除评论
		commentId, _ := strconv.Atoi(c.Query("comment_id")) //评论id
		var comment = model.Comment{}
		err := commentMan.QueryCommentById(int64(commentId), &comment)
		if err != nil {
			resp.Status(StatusCommentNotFound)
			log.Println("该评论不存在", err.Error())
			return
		}
		err = commentMan.DeleteCommentAndUpdateCountById(int64(commentId), int64(videoId))
		if err != nil {
			resp.Status(StatusFailedDelComment)
			log.Println("删除评论失败", err.Error())
			return
		}
	}
}

func (h *Handler) CommentList(c *gin.Context) {
	type Resp struct {
		Response              //通用字段
		CommentList []Comment `json:"comment_list,omitempty"` //评论集合
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()

	commentMan := h.DB.Comment
	userMan := h.DB.User
	videoId, _ := strconv.Atoi(c.Query("video_id"))

	//获取该视频的所有评论
	var comments []model.Comment
	err := commentMan.QueryCommentListByVideoId(int64(videoId), &comments)
	if err != nil { //获取评论列表失败
		resp.Status(StatusFailedCommentList)
		log.Println("拉取评论列表失败", err.Error())
		return
	}

	//评论列表为空
	if comments == nil {
		resp.Status(StatusVideoHasNoComment)
		return
	}

	commentList := make([]Comment, len(comments)) //定义切片大小
	idx := 0
	for _, comment := range comments {
		var u model.User
		_ = userMan.QueryById(comment.UserID, &u)
		userInfo := h.ConvertUser(u, false)
		commentData := Comment{
			Id:         int64(comment.Model.ID),
			User:       userInfo,
			Content:    comment.Content,
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
		}
		commentList[idx] = commentData
		idx = idx + 1
	}
	resp.CommentList = commentList
}
