package config

import (
	"errors"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServeAddr string  `yaml:"serve-addr"`
	Douyin    Douyin  `yaml:"douyin"`
	Storage   Storage `yaml:"storage"`
}
type Douyin struct {
	SQL   SQL   `yaml:"sql"`
	Redis Redis `yaml:"redis"`
	// feed 流一次推送的视频个数，最大 30
	FeedNum int `yaml:"feed-num"`
	// 视频屠夫最大同时压缩视频的数量
	VideoButCherMaxJob int `yaml:"video-butcher-max-job"`
	// 头像和背景
	Avatars     []string `yaml:"avatars"`
	Backgrounds []string `yaml:"backgrounds"`
}

type SQL struct {
	// 数据库驱动类型，只支持 mysql 和 sqlite
	Type string `yaml:"type"`
	DSN  string `yaml:"dsn"`
}
type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
type Storage struct {
	DataDir string `yaml:"data-dir"`
	PreURL  string `yaml:"serve-addr"`
	Token   string `yaml:"token"`
}

// 默认配置
func defaultConfig() *Config {
	return &Config{
		Douyin: Douyin{
			SQL: SQL{
				Type: "sqlite",
				DSN:  "data.db",
			},
			FeedNum:            1,
			VideoButCherMaxJob: 2,
		},
		Storage: Storage{
			DataDir: "Data",
			Token:   "Kitty",
		},
	}
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
	if err != nil {
		return nil, err
	}
	if err := verify(config); err != nil {
		return nil, err
	}
	return config, err
}

// 创建一个新文件，初始化一个空的配置
func New(fileName string) error {
	// 设置默认值
	config := defaultConfig()
	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, out, 0755)
	return err
}

func verify(c *Config) error {
	sqlType := c.Douyin.SQL.Type
	if !(sqlType == "sqlite" || sqlType == "mysql") {
		return errors.New("配置验证失败: Douyin.SQL.Type 只能为 sqlite 和 mysql, 而不是 \"" + sqlType + "\"")
	}
	if c.Douyin.VideoButCherMaxJob < 1 {
		return errors.New("配置验证失败: Douyin.Douyin.VideoButCherMaxJob 必须大于 1, 而不是 \"" +
			strconv.Itoa(c.Douyin.VideoButCherMaxJob) + "\"")
	}
	return nil
}
