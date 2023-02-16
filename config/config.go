package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Douyin  Douyin  `yaml:"douyin"`
	Storage Storage `yaml:"storage"`
}
type Douyin struct {
	ServeAddr string `yaml:"serve-addr"`
	SQL       SQL    `yaml:"sql"`
	Redis     Redis  `yaml:"redis"`
	// feed 流一次推送的视频个数，最大 30
	FeedNum int `yaml:"feed-num"`
}

type SQL struct {
	DSN string `yaml:"dsn"`
}
type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
type Storage struct {
	DataDir   string `yaml:"data-dir"`
	ServeAddr string `yaml:"serve-addr"`
}

// 从 fileName 获取配置
// 如果文件不存在不会创建文件，并返回错误。
// 文件不存在应调用 New()
func Get(filName string) (*Config, error) {
	config := &Config{}
	yamlFile, err := os.ReadFile(filName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, config)
	return config, err
}

// 创建一个新文件，初始化一个空的配置
func New(fileName string) error {
	// 设置默认值
	config := &Config{
		Douyin: Douyin{
			SQL: SQL{
				DSN: "./data.db",
			},
		},
		Storage: Storage{
			DataDir:   "./Data",
			ServeAddr: "由于要拼接 URL, 所以此处不能写 ':port', 要写确切的 'ip:port'",
		},
	}
	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, out, 0755)
	return err
}
