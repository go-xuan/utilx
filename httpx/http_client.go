package httpx

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

// NewClient 新建http客户端
func NewClient(opts ...Option) *http.Client {
	client := http.DefaultClient
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// Option 配置选项
type Option func(client *http.Client)

// SetTimeout 设置超时时间
func SetTimeout(timeout time.Duration) Option {
	return func(client *http.Client) {
		httpClientSetTimeout(client, timeout)
	}
}

// SetCrt 设置https证书路径
func SetCrt(crt string) Option {
	return func(client *http.Client) {
		httpClientSetCrt(client, crt)
	}
}

// SetProxy 设置代理地址
func SetProxy(proxy string) Option {
	return func(client *http.Client) {
		httpClientSetProxy(client, proxy)
	}
}

// 设置http客户端超时时间
func httpClientSetTimeout(client *http.Client, timeout time.Duration) {
	if timeout > 0 {
		client.Timeout = timeout
	}
}

// 设置http客户端证书
func httpClientSetCrt(client *http.Client, crt string) {
	if crt != "" {
		if pem, err := os.ReadFile(crt); err == nil {
			if pool := x509.NewCertPool(); !pool.AppendCertsFromPEM(pem) {
				transport := DefaultTransport()
				transport.TLSClientConfig = &tls.Config{
					RootCAs:            pool,
					InsecureSkipVerify: true,
				}
				client.Transport = transport
			}
		}
	}
}

// 设置http客户端代理
func httpClientSetProxy(client *http.Client, proxy string) {
	if proxy != "" {
		if u, err := url.Parse(proxy); err == nil {
			transport := DefaultTransport()
			transport.Proxy = http.ProxyURL(u)
			client.Transport = transport
		}
	}
}

// DefaultTransport 默认transport
func DefaultTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout:   time.Second * 10,
		KeepAlive: time.Second * 10,
	}
	return &http.Transport{
		DialContext:           dialer.DialContext,
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   50,
		IdleConnTimeout:       time.Second * 90,
		TLSHandshakeTimeout:   time.Second * 10,
		ExpectContinueTimeout: time.Second,
	}
}
