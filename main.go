package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/HumXC/simple-douyin/config"
	"github.com/HumXC/simple-douyin/model"
	"github.com/HumXC/simple-douyin/service"
	"github.com/gin-gonic/gin"
)

const ConfigFile = "./config.yaml"

func main() {
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
	// 初始化数据库
	db, err := model.NewDouyinDB(conf.Douyin.SQL.DSN)
	if err != nil {
		panic(err)
	}

	// 初始化 gin engine
	douyinEngine := gin.Default()
	storageEngine := gin.Default()

	storage := service.NewStorage(storageEngine, conf.Storage)
	_ = service.NewDouyin(douyinEngine, db, storage)

	go func(s *gin.Engine, serveAddr string) {
		panic(s.Run(serveAddr))
	}(storageEngine, conf.Storage.ServeAddr)

	panic(douyinEngine.Run(conf.Douyin.ServeAddr))
}
