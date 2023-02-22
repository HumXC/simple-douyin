package sqldb

// 数据库相关
import (
	"errors"

	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

// 初始化一个用于 douyin 业务的数据库，只支持 sqlite，fileName 是数据库文件的文件名
// 例如 NewDouyinDB("./data.db")
func NewDouyinDB(dbType string, dsn string) (*douyin.DBMan, error) {
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

	db.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.Comment{},
		&model.ThumbsUp{},
		&model.VideoJob{},
		&model.Message{})

	return &douyin.DBMan{
		User:  &UserMan{DB: db},
		Video: &VideoMan{DB: db},
		ThumbsUp: &ThumbsUpMan{
			DB: db,
		},
		Comment: &CommentMan{DB: db},
		VideoJob: &videoJobMan{
			db: db,
		},
		Message: &MessageMan{
			DB: db,
		},
	}, nil
}
