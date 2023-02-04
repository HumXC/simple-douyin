package model

// 数据库相关
import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var RDB = InitRedisDB()

type DouyinDB struct {
	User 	  *UserMan
	Video     *videoMan
	ThumbsUp  *ThumbsUpMan
	Comment   *CommentMan
	UserLogin *UserLoginMan
}

// 初始化一个用于 douyin 业务的数据库，只支持 sqlite，fileName 是数据库文件的文件名
// 例如 NewDouyinDB("./data.db")
func NewDouyinDB(fileName string) (*DouyinDB, error) {
	db, err := gorm.Open(
		sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Video{}, &Comment{}, &ThumbsUp{})
	return &DouyinDB{
		User:	  &UserMan{db: db},
		Video:    &videoMan{db: db},
		ThumbsUp: &ThumbsUpMan{db: db},
		Comment:  &CommentMan{db: db},
	}, nil
}

func InitRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
