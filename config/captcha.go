package config

type Captcha struct {
	Height   int     `json:"height" yaml:"height"`       // 验证码图片高度
	Width    int     `json:"width" yaml:"width"`         // 验证
	Length   int     `json:"length" yaml:"length"`       // 验证码长度
	MaxSkew  float64 `json:"max_skew" yaml:"max_skew"`   // 最大偏斜度
	DotCount int     `json:"dot_count" yaml:"dot_count"` // 点的数量
}
