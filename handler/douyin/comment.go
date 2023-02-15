package douyin

import (
	"github.com/HumXC/simple-douyin/helper"
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

/**
 * @Description 评论操作-控制层
 * @Author xyc
 * @Date 2023/1/27 11:53
 **/

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

func (h *Handler) CommentAction(c *gin.Context) {
	commentMan := h.DB.Comment
	userMan := h.DB.User
	token := c.Query("token")
	//解析token
	userClaim, _ := helper.AnalyseToken(token)
	userId := userClaim.UserId
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
			c.JSON(http.StatusOK, Response{
				StatusCode: StatusOtherError,
				StatusMsg:  "发布评论失败!",
			})
			log.Println("发布评论失败:", err.Error())
			return
		}
		//发布评论成功
		var user model.User
		userMan.QueryUserInfoByUserId(userId, &user)
		userInfo := User{
			Id:             user.Id,
			Name:           user.Name,
			FollowCount:    user.FollowCount,
			FollowerCount:  user.FollowerCount,
			IsFollow:       user.IsFollow,
			TotalFavorited: user.TotalFavorited,
			FavoriteCount:  user.FavoriteCount,
		}
		commentData := Comment{
			Id:         int64(comment.Model.ID),
			User:       userInfo,
			Content:    comment.Content,
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: StatusOK,
				StatusMsg:  "发布评论成功!",
			},
			Comment: commentData,
		})
		return
	} else {
		//2-删除评论
		commentId, _ := strconv.Atoi(c.Query("comment_id")) //评论id
		var comment = model.Comment{}
		err := commentMan.QueryCommentById(int64(commentId), &comment)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: StatusOtherError,
				StatusMsg:  "该评论不存在!",
			})
			log.Println("该评论不存在", err.Error())
			return
		}
		err = commentMan.DeleteCommentAndUpdateCountById(int64(commentId), int64(videoId))
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: StatusOtherError,
				StatusMsg:  "删除评论失败!",
			})
			log.Println("删除评论失败", err.Error())
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: StatusOK,
				StatusMsg:  "删除评论成功!",
			},
		})
		return
	}
}

func (h *Handler) CommentList(c *gin.Context) {
	commentMan := h.DB.Comment
	userMan := h.DB.User
	//token := c.Query("token")
	////解析token
	//userClaim, _ := helper.AnalyseToken(token)
	//_ = userClaim.UserId
	videoId, _ := strconv.Atoi(c.Query("video_id"))

	//获取该视频的所有评论
	var comments []model.Comment
	err := commentMan.QueryCommentListByVideoId(int64(videoId), &comments)
	if err != nil { //获取评论列表失败
		c.JSON(http.StatusOK, Response{
			StatusCode: StatusOtherError,
			StatusMsg:  "获取评论列表失败",
		})
		log.Println("获取评论列表失败", err.Error())
		return
	}

	//评论列表为空
	if comments == nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{
				StatusCode: StatusOK,
				StatusMsg:  "该视频暂无评论!",
			},
		})
		return
	}

	commentList := make([]Comment, len(comments)) //定义切片大小
	idx := 0
	for _, comment := range comments {
		var user model.User //每个评论的用户信息
		userMan.QueryUserInfoByUserId(comment.UserID, &user)
		userInfo := User{
			Id:             user.Id,
			Name:           user.Name,
			FollowCount:    user.FollowCount,
			FollowerCount:  user.FollowerCount,
			IsFollow:       user.IsFollow,
			TotalFavorited: user.TotalFavorited,
			FavoriteCount:  user.FavoriteCount,
		}
		commentData := Comment{
			Id:         int64(comment.Model.ID),
			User:       userInfo,
			Content:    comment.Content,
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
		}
		commentList[idx] = commentData
		idx = idx + 1
	}

	//TODO 评论切片待排序,按时间倒序
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: StatusOK,
			StatusMsg:  "查询评论列表成功!",
		},
		CommentList: commentList,
	})
}
