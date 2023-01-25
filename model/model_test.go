package model_test

import (
	"os"
	"testing"
)

// 该包所有的测试文件都存在 "../test/model" 里
const TEST_DIR = "../test/model"

// 清理测试环境
func TestMain(m *testing.M) {
	_ = os.RemoveAll(TEST_DIR)
	_ = os.MkdirAll(TEST_DIR, 0755)
	m.Run()
}
