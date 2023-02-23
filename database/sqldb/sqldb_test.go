package sqldb_test

import (
	"os"
	"path"
	"testing"

	"github.com/HumXC/simple-douyin/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 该包所有的测试文件都存在 "../test/database/sqldb" 里
const TEST_DIR = "../../test/database/sqldb"

func NewDB() *gorm.DB {
	db, err := gorm.Open(
		sqlite.Open(path.Join(TEST_DIR, "data.db")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.Comment{},
		&model.VideoJob{},
		&model.Message{})
	return db
}

// 清理测试环境
func TestMain(m *testing.M) {
	_ = os.RemoveAll(TEST_DIR)
	_ = os.MkdirAll(TEST_DIR, 0755)
	m.Run()
}
