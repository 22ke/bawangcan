package config

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// Config 结构体定义了你的配置项
type Config struct {
	Dper     string `yaml:"dper"`
	Cityname string `yaml:"cityname"`
	Menu     []int  `yaml:"menu"`
	Internal int    `yaml:"internaltime""`
	Debug    bool   `yaml:"debug"`
}

func Makeconfig() Config {
	// 读取 YAML 文件内容
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// 创建 Config 对象
	var config Config

	// 解析 YAML 数据到 Config 结构体
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	// 打印解析后的配置
	login(&config)
	fmt.Printf("%+v\n", config)
	return config
}
