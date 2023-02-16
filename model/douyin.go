package model

// 数据库相关
import (
	"errors"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
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
func NewDouyinDB(dbType string, dsn string, rdb *redis.Client) (*DouyinDB, error) {
	var db *gorm.DB
	switch dbType {
	case "sqlite":
		_db, err := gorm.Open(
			sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		db = _db
	case "mysql":
		_db, err := gorm.Open(
			mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		db = _db
	default:
		return nil, errors.New("不支持的数据库类型: " + dbType)
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
