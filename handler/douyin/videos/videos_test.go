package videos_test

import (
	"os"
	"testing"
)

const TestDir = "../../../test"

func TestMain(m *testing.M) {
	_ = os.MkdirAll(TestDir, 0755)
	m.Run()
}
