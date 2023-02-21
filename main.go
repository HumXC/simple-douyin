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

	// 初始化两个服务，其中 douyin 服务依赖 storage 服务
	// 如果配置里两个服务的 ServeAddr 相同，就用同一个 gin.Engine,否则使用两个 gin.Engine 实例
	var storageEngine *gin.Engine = NewEngine()
	var douyinEngine *gin.Engine = storageEngine
	if conf.Douyin.ServeAddr != conf.Storage.ServeAddr {
		douyinEngine = NewEngine()
	}

	storageClient := Storage(storageEngine, conf.Storage)
	Douyin(douyinEngine, conf.Douyin, storageClient)

	// 另开个协程抛跑 storage 服务
	go func() {
		panic(storageEngine.Run(conf.Storage.ServeAddr))
	}()

	panic(douyinEngine.Run(conf.Douyin.ServeAddr))
}

func NewEngine() *gin.Engine {
	return gin.Default()
}
func Douyin(engine *gin.Engine, c config.Douyin, storage douyin.StorageClient) {
	db, err := sqldb.NewDouyinDB(c.SQL.Type, c.SQL.DSN)
	if err != nil {
		panic(err)
	}
	rdb, err := cache.NewDouyinRDB(nil)
	if err != nil {
		panic(err)
	}
	_ = service.NewDouyin(engine, c, db, rdb, storage)
}

func Storage(engine *gin.Engine, c config.Storage) *service.Storage {
	storage := service.NewStorage(engine, c)
	return storage
}
