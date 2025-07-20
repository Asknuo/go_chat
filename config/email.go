package config

type Email struct {
	Host     string `json:"host" yaml:"host"`         // 邮件服务器的地址
	Port     int    `json:"port" yaml:"port"`         // 邮件服务器
	From     string `json:"from" yaml:"from"`         // 发件人邮箱地址
	Nickname string `json:"nickname" yaml:"nickname"` // 发件人昵称
	Secret   string `json:"secret" yaml:"secret"`     // 邮箱授权码或密码
	IsSSL    bool   `json:"is_ssl" yaml:"is_ssl"`     // 是否
}
