package sqldb

import (
	"errors"

	"github.com/HumXC/simple-douyin/model"
	"gorm.io/gorm"
)

type messageMan struct {
	db *gorm.DB
}

func (m *messageMan) AddMessage(message *model.Message) error {
	if message == nil {
		return errors.New("AddMessage message空指针")
	}
	//添加消息
	if err := m.db.Create(message).Error; err != nil {
		return err
	}
	return nil
}

func (m *messageMan) QueryMessageRecord(fromUserId int64, toUserId int64, messages *[]model.Message) error {
	if messages == nil {
		return errors.New("QueryMessageRecord messages空指针")
	}
	if err := m.db.Model(&model.Message{}).Where("from_user_id = ? and to_user_id = ?", fromUserId, toUserId).Find(messages).Error; err != nil {
		return err
	}
	return nil
}
