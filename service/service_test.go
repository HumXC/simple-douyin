package service_test

import (
	"os"
	"testing"
)

const TEST_DIR = "../test/service"

func TestMain(m *testing.M) {
	_ = os.MkdirAll(TEST_DIR, 0755)
	m.Run()
}
