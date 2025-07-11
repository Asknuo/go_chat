package config

type Config struct {
	Mysql Mysql `json:"mysql" yaml:"mysql"` // MySQL数据库配置
}
