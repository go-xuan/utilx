package httpx

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

// DefaultSettings 默认配置
func DefaultSettings() *Settings {
	return &Settings{Timeout: 10, Crt: "", Proxy: ""}
}

// Settings 客户端配置
type Settings struct {
	Timeout int
	Crt     string
	Proxy   string
}

func (s *Settings) UniqueId() string {
	return fmt.Sprintf("%d_%s_%s", s.Timeout, s.Crt, s.Proxy)
}

// NewClient 新建http客户端
func (s *Settings) NewClient() *http.Client {
	client := http.DefaultClient
	SetHttpClientTimeout(s.Timeout)(client)
	SetHttpClientCrt(s.Crt)(client)
	SetHttpClientProxy(s.Proxy)(client)
	return client
}

// SettingsOption 客户端配置选项
type SettingsOption func(*Settings)

// SetTimeout 设置超时时间
func SetTimeout(timeout int) SettingsOption {
	return func(s *Settings) {
		s.Timeout = timeout
	}
}

// SetCrt 设置证书
func SetCrt(crt string) SettingsOption {
	return func(s *Settings) {
		s.Crt = crt
	}
}

// SetProxy 设置代理
func SetProxy(proxy string) SettingsOption {
	return func(s *Settings) {}
}

// HttpClientOption 客户端选项
type HttpClientOption func(*http.Client)

// SetHttpClientTimeout 设置超时时间
func SetHttpClientTimeout(timeout int) HttpClientOption {
	return func(client *http.Client) {
		if timeout > 0 {
			client.Timeout = time.Duration(timeout) * time.Second
		}
	}
}

// SetHttpClientCrt 设置证书
func SetHttpClientCrt(crt string) HttpClientOption {
	return func(client *http.Client) {
		if crt != "" {
			if pem, err := os.ReadFile(crt); err == nil {
				pool := x509.NewCertPool()
				if !pool.AppendCertsFromPEM(pem) {
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
}

// SetHttpClientProxy 设置代理
func SetHttpClientProxy(proxy string) HttpClientOption {
	return func(client *http.Client) {
		if proxy != "" {
			if u, err := url.Parse(proxy); err == nil {
				transport := DefaultTransport()
				transport.Proxy = http.ProxyURL(u)
				client.Transport = transport
			}
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
