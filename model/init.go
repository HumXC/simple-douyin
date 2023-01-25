package model

// 数据库相关
import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var RDB = InitRedisDB()

// 初始化数据库，只支持 sqlite，fileName 是数据库文件的文件名
// 例如 InitDB("./data.db")
func InitDB(fileName string) error {
	var err error
	db, err = gorm.Open(
		sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&Video{}, &ThumbsUp{})
	return nil
}

func InitRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
