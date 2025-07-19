package config

type System struct {
	Host           string `json:"host" yaml:"host"`                       // 系统监听的主机地址
	Port           int    `json:"port" yaml:"port"`                       // 系统监听的端口号
	Env            string `json:"env" yaml:"env"`                         // 系统运行环境
	RouterPrefix   string `json:"router_prefix" yaml:"router_prefix"`     // 路由前缀
	UseMultipoint  bool   `json:"use_multipoint" yaml:"use_multipoint"`   // 是否使用多点登录
	SessionsSecret string `json:"sessions_secret" yaml:"sessions_secret"` // 会话密	钥
	OssType        string `json:"oss_type" yaml:"oss_type"`               // 对
}
