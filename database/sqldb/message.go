package sqldb

import (
	"errors"

	"github.com/HumXC/simple-douyin/model"
	"gorm.io/gorm"
)

type MessageMan struct {
	DB *gorm.DB
}

func (m *MessageMan) AddMessage(message *model.Message) error {
	if message == nil {
		return errors.New("AddMessage message空指针")
	}
	//添加消息
	if err := m.DB.Create(message).Error; err != nil {
		return err
	}
	return nil
}

func (m *MessageMan) QueryChat(fromUserId int64, toUserId int64, time string, messages *[]model.Message) error {
	if messages == nil {
		return errors.New("QueryMessageRecord messages空指针")
	}
	if err := m.DB.Model(&model.Message{}).Where("created_at > ? and ((from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?))", time, fromUserId, toUserId, toUserId, fromUserId).Find(messages).Error; err != nil {
		return err
	}
	return nil
}

func (m *MessageMan) QueryNewMsg(userId1 int64, userId2 int64, message *model.Message) error {
	if message == nil {
		return errors.New("QueryNewMsg message空指针")
	}
	if err := m.DB.Model(&model.Message{}).Order("created_at DESC").Where("((from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?))", userId1, userId2, userId2, userId1).Limit(1).Find(message).Error; err != nil {
		return err
	}
	return nil
}
