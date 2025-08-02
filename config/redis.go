package config

type Redis struct {
	Host  string `json:"host" yaml:"host"` // Redis服务器的地址
	Port  int    `json:"port" yaml:"port"` // Redis服务器的端口号
	Table string `json:"table" yaml:"table"`
	Name  string `json:"name"  yaml:"name"`
}
