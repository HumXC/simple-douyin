package service_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/HumXC/simple-douyin/service"
	"github.com/gin-gonic/gin"
)

func init() {
	os.Mkdir(path.Join(TEST_DIR, "storage"), 0755)
}

// 向 r 发出请求
func DoRequest(r http.Handler, method, path string, data url.Values) (*http.Response, error) {
	reqbody := strings.NewReader(data.Encode())
	req, err := http.NewRequest(method, "/storage"+path, reqbody)
	if err != nil {
		return nil, err
	}
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec.Result(), nil
}
func TestUpload(t *testing.T) {
	dataDir := path.Join(TEST_DIR, "storage")
	s := service.NewStorage(gin.Default(), service.StorageOption{
		DataDir: dataDir,
	})
	want := "3cf571d4cf2a4c4b2df823a27852a7d5"
	dir := "videos"
	file := "../test/video.mp4"
	// 测试新文件
	hash, err := s.Upload(file, dir)
	if err != nil {
		t.Error(err)
		return
	}
	if hash != want {
		t.Errorf("md5 值不匹配 got: %s,want: %s", hash, want)
		return
	}

	// 测试已经存在的文件
	hash, err = s.Upload(file, dir)
	if err != nil {
		t.Error(err)
		return
	}
	if hash != want {
		t.Errorf("md5 值不匹配 got: %s,want: %s", hash, want)
	}
}
func TestVideo(t *testing.T) {
	dataDir := path.Join(TEST_DIR, "storage")
	r := gin.Default()
	_ = service.NewStorage(r, service.StorageOption{
		DataDir: dataDir,
	})

	// 在文件夹里存储一个视频
	videoData := "假装这是一个视频文件"
	hash := "testvideo"
	_ = os.MkdirAll(path.Join(dataDir, "videos"), 0755)
	err := os.WriteFile(path.Join(dataDir, "videos", hash), []byte(videoData), 0755)
	if err != nil {
		t.Fatal("无法写入视频文件", err)
	}

	result, err := DoRequest(r, http.MethodGet, "/video/errvideo", url.Values{})
	if err != nil {
		t.Fatal(err)
	}
	if result.StatusCode != 404 {
		t.Errorf("错误的状态码: got:%d want:%d", result.StatusCode, 404)
	}

	result, err = DoRequest(r, http.MethodGet, "/video/"+hash, url.Values{})
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != videoData {
		t.Errorf("文件内容不符合预期: got:%s want:%s", string(body), videoData)
	}
}
