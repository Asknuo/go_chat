package config

type Config struct {
	Mysql Mysql `json:"mysql" yaml:"mysql"` // MySQL数据库配置
	Zap   Zap   `json:"zap" yaml:"zap"`     // 日志配置
}
