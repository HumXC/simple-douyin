package sqldb

import (
	"errors"

	"github.com/HumXC/simple-douyin/model"
	"gorm.io/gorm"
)

/**
 * @Description 评论操作-持久层
 * @Author xyc
 * @Date 2023/1/26 21:23
 **/

type CommentMan struct {
	DB *gorm.DB
}

func (c *CommentMan) AddComment(comment *model.Comment) error {
	if comment == nil {
		return errors.New("AddComment comment空指针")
	}
	//添加评论
	if err := c.DB.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func (c *CommentMan) AddCommentAndUpdateCommentCount(comment *model.Comment) error {
	if comment == nil {
		return errors.New("AddCommentAndUpdateCount comment空指针")
	}
	//执行事务
	return c.DB.Transaction(func(tx *gorm.DB) error {
		//添加评论数据
		if err := tx.Create(comment).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		//增加count
		if err := tx.Exec("UPDATE videos  SET comment_count = comment_count+1 WHERE id=?", comment.VideoId).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

func (c *CommentMan) DeleteCommentAndUpdateCountById(commentId, videoId int64) error {
	//执行事务
	return c.DB.Transaction(func(tx *gorm.DB) error {
		//删除评论
		if err := tx.Exec("DELETE FROM comments WHERE id = ?", commentId).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		//视频评论数减一
		if err := tx.Exec("UPDATE videos SET comment_count = comment_count-1 WHERE id=? AND comment_count>0", videoId).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

func (c *CommentMan) QueryCommentById(id int64, comment *model.Comment) error {
	if comment == nil {
		return errors.New("QueryCommentById comment 空指针")
	}
	return c.DB.Where("id=?", id).First(comment).Error
}

func (c *CommentMan) QueryCommentListByVideoId(videoId int64, comments *[]model.Comment) error {
	if comments == nil {
		return errors.New("QueryCommentListByVideoId comments空指针")
	}
	if err := c.DB.Model(&model.Comment{}).Where("video_id=?", videoId).Order("created_at DESC").Find(comments).Error; err != nil {
		return err
	}
	return nil
}
