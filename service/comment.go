package service

import (
	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/model"
	"path"
	"time"
)

/**
 * @Description 评论操作-业务层
 * @Author xyc
 * @Date 2023/1/28 10:34
 **/

var douyinDB, _ = model.NewDouyinDB(path.Join("../db", "douyin.db"))
var commentMan = douyinDB.Comment

func AddComment(comment *model.Comment) (douyin.Comment, error) {
	err := commentMan.AddCommentAndUpdateCommentCount(comment)
	if err != nil {
		return douyin.Comment{}, err
	}
	//TODO 根据user_id获取user对象

	//封装返回数据
	commentData := douyin.Comment{
		Id:         int64(comment.Model.ID),
		User:       douyin.User{}, //TODO 此处设置评论的user对象,暂时置为空
		Content:    comment.Content,
		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
	}
	return commentData, nil
}

func DelComment(commentId, videoId int64) error {
	var comment = model.Comment{}
	err := commentMan.QueryCommentById(commentId, &comment)
	if err != nil {
		return err
	}
	err = commentMan.DeleteCommentAndUpdateCountById(commentId, videoId)
	if err != nil {
		return err
	}
	return nil
}
