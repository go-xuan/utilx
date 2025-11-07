package httpx

import "net/http"

// Decorator 装饰器
type Decorator interface {
	Decorate(request *http.Request)
}

// HeaderDecorator 请求头装饰
type HeaderDecorator struct {
	Key   string
	Value string
}

func (d *HeaderDecorator) Decorate(request *http.Request) {
	if d.Key != "" && d.Value != "" {
		request.Header.Set(d.Key, d.Value)
	}
}

// CookieDecorator cookie装饰
type CookieDecorator struct {
	Key   string
	Value string
}

func (d *CookieDecorator) Decorate(request *http.Request) {
	if d.Key != "" && d.Value != "" {
		request.AddCookie(&http.Cookie{
			Name:  d.Key,
			Value: d.Value,
		})
	}
}
