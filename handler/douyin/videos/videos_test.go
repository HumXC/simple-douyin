package videos_test

import (
	"os"
	"testing"
)

const TestDir = "../../../test/videos"

func TestMain(m *testing.M) {
	_ = os.MkdirAll(TestDir, 0755)
	m.Run()
}
