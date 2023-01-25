package model_test

import (
	"os"
	"testing"

	"github.com/HumXC/simple-douyin/model"
)

// 该包所有的测试文件都存在 "../test/model" 里
const TEST_DIR = "../test/model"

// 测试用 sqlite 数据库的文件位置
const DB_NAME = "../test/model/test.db"

// 清理测试环境
func TestMain(m *testing.M) {
	_ = os.RemoveAll(TEST_DIR)
	_ = os.MkdirAll(TEST_DIR, 0755)
	err := model.InitDB(DB_NAME)
	if err != nil {
		panic(err)
	}
	m.Run()
}
