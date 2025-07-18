package config

type Config struct {
	Mysql  Mysql   `json:"mysql" yaml:"mysql"`   // MySQL数据库配置
	Zap    Zap     `json:"zap" yaml:"zap"`       // 日志配置
	System System  `json:"system" yaml:"system"` // 系统配置
	Mongo  MongoDB `json:"mongo" yaml:"mongo"`   // MongoDB配置
	Redis  Redis   `json:"redis" yaml:"redis"`   // Redis配置
}
