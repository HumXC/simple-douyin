package model

// 数据库相关
import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// 初始化数据库，只支持 sqlite，fileName 是数据库文件的文件名
// 例如 InitDB("./data.db")
func InitDB(fileName string) error {
	var err error
	db, err = gorm.Open(
		sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		return err
	}
	// db.AutoMigrate()
	return nil
}
