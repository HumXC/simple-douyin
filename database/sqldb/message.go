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

func (m *MessageMan) QueryMessageRecord(fromUserId int64, toUserId int64, messages *[]model.Message) error {
	if messages == nil {
		return errors.New("QueryMessageRecord messages空指针")
	}
	if err := m.DB.Model(&model.Message{}).Where("from_user_id = ? and to_user_id = ?", fromUserId, toUserId).Find(messages).Error; err != nil {
		return err
	}
	return nil
}
