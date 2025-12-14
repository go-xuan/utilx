package httpx

import (
	"time"
)

// NewConfig 新建配置
func NewConfig(opts ...Option) *Config {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// defaultConfig 默认配置
func defaultConfig() *Config {
	return &Config{
		Timeout: 30 * time.Second,
	}
}

// Config http客户端配置
type Config struct {
	Timeout time.Duration `json:"timeout"` // 超时时间(秒)，默认10秒
	Crt     string        `json:"crt"`     // https证书路径
	Proxy   string        `json:"proxy"`   // 代理地址
}

// Option 配置选项
type Option func(*Config)

// SetTimeout 设置超时时间
func SetTimeout(timeout time.Duration) Option {
	return func(config *Config) {
		config.Timeout = timeout
	}
}

// SetCrt 设置https证书路径
func SetCrt(crt string) Option {
	return func(config *Config) {
		config.Crt = crt
	}
}

// SetProxy 设置代理地址
func SetProxy(proxy string) Option {
	return func(config *Config) {
		config.Proxy = proxy
	}
}
