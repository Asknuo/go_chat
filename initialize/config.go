package initialize

import (
	"gochat/config"
	"gochat/utlis"

	"gopkg.in/yaml.v3"
)

func InitConfig() *config.Config {
	c := &config.Config{}
	yamlConf, err := utlis.LoadYaml()
	if err != nil {
		panic("加载配置文件失败: " + err.Error())
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		panic("解析配置文件失败: " + err.Error())
	}
	return c
}
