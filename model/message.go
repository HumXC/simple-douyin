package model

import (
	"errors"
	"gorm.io/gorm"
)

type Message struct {
	FromUserId int64 //发送者ID
	ToUserId   int64 //接收者ID
	Content    string
	gorm.Model //通用字段
}

type messageMan struct {
	db *gorm.DB
}

func (m *messageMan) AddMessage(message *Message) error {
	if message == nil {
		return errors.New("AddMessage message空指针")
	}
	//添加消息
	if err := m.db.Create(message).Error; err != nil {
		return err
	}
	return nil
}

func (m *messageMan) QueryMessageRecord(fromUserId int64, toUserId int64, messages *[]Message) error {
	if messages == nil {
		return errors.New("QueryMessageRecord messages空指针")
	}
	if err := m.db.Model(&Message{}).Where("from_user_id = ? and to_user_id = ?", fromUserId, toUserId).Find(messages).Error; err != nil {
		return err
	}
	return nil
}
