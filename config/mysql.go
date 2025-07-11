package config

import (
	"strconv"
	"strings"

	"gorm.io/gorm/logger"
)

type Mysql struct {
	Host         string `json:"host" yaml:"host"`                     // 数据库服务器的地址
	Port         int    `json:"port" yaml:"port"`                     // 数据库服务器的端口号
	Config       string `json:"config" yaml:"config"`                 // 数据库连接的配置参数，如驱动、字符集等
	DBName       string `json:"db_name" yaml:"db_name"`               // 要连接的数据库名称
	Username     string `json:"username" yaml:"username"`             // 用于连接数据库的用户名
	Password     string `json:"password" yaml:"password"`             // 用于连接数据库的密码
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"` // 数据库连接池中的最大空闲连接数
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"` // 数据库的最大打开连接数
	LogMode      string `json:"log_mode" yaml:"log_mode"`             // 日志模式，例如 "info" 或 "silent"，用于控制日志输出
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/" + m.DBName + "?" + m.Config
}

func (m Mysql) LogLevel() logger.LogLevel {
	switch strings.ToLower(m.LogMode) {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}
