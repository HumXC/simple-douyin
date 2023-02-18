package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/HumXC/simple-douyin/config"
	"github.com/HumXC/simple-douyin/handler/storage"
	"github.com/gin-gonic/gin"
)

// 该文件实现了一个文件服务器, 用于给用户存储/发送视频

type Storage struct {
	DataDir   string
	hashToURL func(dir, hash string) string
}
type StorageOption struct {
	DataDir   string
	URLPrefix string
}

// 将文件保存到本地存储，保存完成返回文件的 MD5 hash 值
// dir 是需要保存的目录
// 如果 dir="videos", 那么上传的文件就会保存在 [DataDir]/videos 目录
func (s *Storage) Upload(file, dir string) (string, error) {
	// 创建文件夹
	makeErr := func(err error) error {
		return fmt.Errorf("上传文件失败 [%s] 到 [%s]: %w", file, dir, err)
	}
	fullDir := path.Join(s.DataDir, dir)
	_, err := os.Stat(fullDir)
	if err != nil || os.IsNotExist(err) {
		_ = os.Mkdir(fullDir, 0755)
	}
	f, err := os.CreateTemp(s.DataDir, "upload_video")
	if err != nil {
		return "", makeErr(err)
	}
	defer f.Close()
	defer os.Remove(f.Name())
	b := bytes.Buffer{}
	src, err := os.Open(file)
	if err != nil {
		return "", makeErr(err)
	}
	defer src.Close()
	_, err = b.ReadFrom(src)
	if err != nil {
		return "", makeErr(err)
	}
	sum := md5.Sum(b.Bytes())
	hashStr := hex.EncodeToString(sum[:])
	fileName := path.Join(s.DataDir, dir, hashStr)
	// 如果已有同名文件，则删除新创建的文件
	_, err = os.Stat(fileName)
	if err == nil || os.IsExist(err) {
		return hashStr, nil
	}
	_, err = b.WriteTo(f)
	if err != nil {
		f.Close()
		return "", makeErr(err)
	}
	f.Close()
	err = os.Rename(f.Name(), fileName)
	if err != nil {
		return "", makeErr(err)
	}
	return hashStr, nil
}

func (s *Storage) GetURLWithHash(dir, hash string) string {
	return s.hashToURL(dir, hash)
}
func NewStorage(g *gin.Engine, conf config.Storage) *Storage {
	s := &Storage{
		DataDir: conf.DataDir,
	}
	_, err := os.Stat(conf.DataDir)
	if os.IsNotExist(err) {
		os.MkdirAll(conf.DataDir, 0755)
	}
	storageGroup := g.Group("storage")

	storageGroup.GET(":dir/:hash", storage.Fetch(s.DataDir))

	s.hashToURL = func(dir, hash string) string {
		return "http://" + conf.ServeAddr + "/storage/" + dir + "/" + hash
	}
	return s
}
