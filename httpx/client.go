package httpx

import (
	"io"
	"net/http"
	"sync"

	"github.com/go-xuan/utilx/errorx"
)

var _client *Client // http客户端

// GetClient 获取httpx客户端
func GetClient() *Client {
	if _client == nil {
		_client = &Client{
			mu:      sync.Mutex{},
			clients: make(map[string]*http.Client),
		}
		settings := DefaultSettings()
		client := settings.NewClient()
		_client.client = client
		_client.clients[settings.UniqueId()] = client
	}
	return _client
}

// Client httpx客户端
type Client struct {
	mu      sync.Mutex              // 互斥锁
	client  *http.Client            // 默认客户端
	clients map[string]*http.Client // 客户端缓存
}

// HttpClient 获取http客户端
func (c *Client) HttpClient(options ...SettingsOption) *http.Client {
	var client = new(http.Client)
	if len(options) > 0 && options[0] != nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		settings := DefaultSettings()
		for _, option := range options {
			option(settings)
		}
		var ok bool
		if client, ok = c.clients[settings.UniqueId()]; !ok {
			client = settings.NewClient()
			// 缓存客户端, 后续相同配置的请求直接从缓存中获取
			c.clients[settings.UniqueId()] = client
		}
	}
	return c.client
}

// Do 执行http请求
func (c *Client) Do(request *http.Request, option ...SettingsOption) (*Response, error) {
	resp, err := c.HttpClient(option...).Do(request)
	if err != nil {
		return nil, errorx.Wrap(err, "http request error")
	}
	defer resp.Body.Close()

	response := &Response{
		status:  resp.StatusCode,
		cookies: resp.Cookies(),
	}
	if response.body, err = io.ReadAll(resp.Body); err != nil {
		return response, errorx.Wrap(err, "http response body read error")
	}
	return response, nil
}
