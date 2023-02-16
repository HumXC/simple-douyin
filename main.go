package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/HumXC/simple-douyin/config"
	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/model"
	"github.com/HumXC/simple-douyin/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const ConfigFile = "./config.yaml"

func main() {
	// 读取配置文件，没有配置文件就创建一个然后退出
	conf, err := config.Get(ConfigFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := config.New(ConfigFile)
			if err != nil {
				panic(err)
			}
			fmt.Println("已经创建新的配置文件，修改后重新运行: " + ConfigFile)
			os.Exit(0)
		}
		panic(err)
	}

	// 初始化两个独立的服务，其中 douyin 服务依赖 storage 服务
	storageEngine, storageClient := Storage(conf.Storage)
	douyinEngine := Douyin(conf.Douyin, storageClient)

	// 另开个协程抛跑 storage 服务
	go func() {
		panic(storageEngine.Run(conf.Storage.ServeAddr))
	}()

	panic(douyinEngine.Run(conf.Douyin.ServeAddr))
}

func Douyin(c config.Douyin, storage douyin.StorageClient) *gin.Engine {
	engin := gin.Default()
	// 初始化数据库
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})
	db, err := model.NewDouyinDB(c.SQL.DSN, rdb)
	if err != nil {
		panic(err)
	}
	_ = service.NewDouyin(engin, c, db, storage)
	return engin
}

func Storage(c config.Storage) (*gin.Engine, *service.Storage) {
	engin := gin.Default()
	storage := service.NewStorage(engin, c)
	return engin, storage
}
