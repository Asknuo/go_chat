package config

type MongoDB struct {
	Host string `json:"host" yaml:"host"` // MongoDB服务器的地址
	Port int    `json:"port" yaml:"port"` // MongoDB服务器的端口号
	Name string `json:"name" yaml:"name"`
}
