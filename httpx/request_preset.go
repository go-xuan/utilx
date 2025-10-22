package httpx

import "net/http"

// RequestPreset 请求预设
type RequestPreset interface {
	Preset(request *http.Request)
}

// HeaderPreset 请求头预设
type HeaderPreset struct {
	Key   string
	Value string
}

func (t *HeaderPreset) Preset(request *http.Request) {
	if t.Key != "" && t.Value != "" {
		request.Header.Set(t.Key, t.Value)
	}
}

// CookiePreset 请求cookie预设
type CookiePreset struct {
	Key   string
	Value string
}

func (c *CookiePreset) Preset(request *http.Request) {
	if c.Key != "" && c.Value != "" {
		request.AddCookie(&http.Cookie{
			Name:  c.Key,
			Value: c.Value,
		})
	}
}
