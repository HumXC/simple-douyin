package model

// 数据库相关
import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DouyinDB struct {
	User     *userMan
	Video    *videoMan
	ThumbsUp *thumbsUpMan
	Comment  *commentMan
}

// 初始化一个用于 douyin 业务的数据库，只支持 sqlite，fileName 是数据库文件的文件名
// 例如 NewDouyinDB("./data.db")
func NewDouyinDB(fileName string, rdb *redis.Client) (*DouyinDB, error) {
	db, err := gorm.Open(
		sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&User{}, &Video{}, &Comment{}, &ThumbsUp{})
	return &DouyinDB{
		User:  &userMan{db: db},
		Video: &videoMan{db: db},
		ThumbsUp: &thumbsUpMan{
			db:  db,
			rdb: rdb},
		Comment: &commentMan{db: db},
	}, nil
}
