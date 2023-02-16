package config

import (
	"gopkg.in/yaml.v2"
	_ "gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Server  `yaml:"server"`
	SQLLite `yaml:"sqlLite"`
	Redis   `yaml:"redis"`
}
type Server struct {
	Host string
	Port int
}

type SQLLite struct {
	Path string
}

type Redis struct {
	Host     string
	Port     int
	Database int
}

// GetConf 获取配置信息,如 conf.Server.Host
func GetConf() (*Config, error) {
	config := &Config{}
	yamlFile, err := os.Open("./config.yaml")
	if err != nil {
		log.Printf("yaml open err is : \n %v \n", err)
		return nil, err
	}
	yaml.NewDecoder(yamlFile).Decode(config)
	if err != nil {
		log.Printf("unmarshel err is \n %v \n", err)
		return nil, err
	}
	return config, nil
}
