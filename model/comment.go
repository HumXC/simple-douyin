package model

import (
	"errors"
	"gorm.io/gorm"
)

/**
 * @Description
 * @Author xyc
 * @Date 2023/1/26 21:23
 **/

type Comment struct {
	UserID     int64  //用户ID
	VideoId    int64  //视频ID
	Content    string //评论内容
	gorm.Model        //通用字段
}

type commentMan struct {
	db *gorm.DB
}

// AddComment 添加评论方法
func (c *commentMan) AddComment(comment *Comment) error {
	if comment == nil {
		return errors.New("AddComment comment空指针")
	}
	//添加评论
	if err := c.db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

// QueryCommentById 根据评论ID查询评论信息方法
func (c *commentMan) QueryCommentById(id int64, comment *Comment) error {
	if comment == nil {
		return errors.New("QueryCommentById comment 空指针")
	}
	return c.db.Where("id=?", id).First(comment).Error
}

// QueryCommentListByVideoId 根据视频ID查询评论列表
func (c *commentMan) QueryCommentListByVideoId(videoId int64, comments *[]Comment) error {
	if comments == nil {
		return errors.New("QueryCommentListByVideoId comments空指针")
	}
	if err := c.db.Model(&Comment{}).Where("video_id=?", videoId).Find(comments).Error; err != nil {
		return err
	}
	return nil
}
