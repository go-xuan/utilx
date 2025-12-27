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
	if len(opts) > 0 {
		cfg := NewConfig(opts...)
		httpClientSetTimeout(client, cfg.Timeout)
		httpClientSetCrt(client, cfg.Crt)
		httpClientSetProxy(client, cfg.Proxy)
	}
	return client
}

func httpClientSetTimeout(client *http.Client, timeout time.Duration) {
	if timeout > 0 {
		client.Timeout = timeout
	}
}

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
