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
		_client = DefaultClient()
	}
	return _client
}

// DefaultClient 默认http客户端
func DefaultClient() *Client {
	settings := defaultSettings()
	client := settings.NewClient()
	return &Client{
		mu:     sync.RWMutex{},
		client: client,
		clients: map[string]*http.Client{
			settings.Unique(): client,
		},
	}
}

// Client httpx客户端
type Client struct {
	mu      sync.RWMutex            // 读写锁
	client  *http.Client            // 默认客户端
	clients map[string]*http.Client // 客户端缓存
}

// GetHttpClient 获取http客户端
func (c *Client) GetHttpClient(options ...SettingsOption) *http.Client {
	if len(options) > 0 && options[0] != nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		settings := NewSettings(options...)
		unique := settings.Unique()
		if client, ok := c.clients[unique]; ok {
			return client
		} else {
			client = settings.NewClient()
			c.clients[unique] = client
			return client
		}
	}
	return c.client
}

// Do 执行http请求
func (c *Client) Do(request *http.Request, option ...SettingsOption) (*Response, error) {
	resp, err := c.GetHttpClient(option...).Do(request)
	if err != nil {
		return nil, errorx.Wrap(err, "http request error")
	}
	defer resp.Body.Close()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, errorx.Wrap(err, "http response Body read error")
	}
	return &Response{
		Status: resp.StatusCode,
		Body:   body,
		Header: resp.Header,
	}, nil
}
