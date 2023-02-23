package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/HumXC/simple-douyin/config"
	"github.com/HumXC/simple-douyin/database/cache"
	"github.com/HumXC/simple-douyin/database/sqldb"
	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/service"
	"github.com/gin-gonic/gin"
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

	var engine *gin.Engine = NewEngine()
	storageClient := Storage(engine, conf.Storage)
	Douyin(engine, conf.Douyin, storageClient)
	panic(engine.Run(conf.ServeAddr))
}

func NewEngine() *gin.Engine {
	return gin.Default()
}
func Douyin(engine *gin.Engine, c config.Douyin, storage douyin.StorageClient) {
	db, err := sqldb.NewDouyinDB(c.SQL.Type, c.SQL.DSN)
	if err != nil {
		panic(err)
	}
	rdb, err := cache.NewDouyinRDB(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	if err != nil {
		panic(err)
	}
	_ = service.NewDouyin(engine, c, db, rdb, storage)
}

func Storage(engine *gin.Engine, c config.Storage) *service.Storage {
	storage := service.NewStorage(engine, c)
	return storage
}
