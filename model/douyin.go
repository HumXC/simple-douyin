package model

// 数据库相关
import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DouyinDB struct {
	Video *videoMan
}

// 初始化一个用于 douyin 业务的数据库，只支持 sqlite，fileName 是数据库文件的文件名
// 例如 NewDouyinDB("./data.db")
func NewDouyinDB(fileName string) (*DouyinDB, error) {
	db, err := gorm.Open(
		sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Video{})
	return &DouyinDB{
		Video: &videoMan{db: db},
	}, nil
}
