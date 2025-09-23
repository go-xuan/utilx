package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/go-xuan/utilx/errorx"
)

// Response 表示 HTTP 请求的响应结构
type Response struct {
	Trace  string      // 追踪标识，用于日志或调试
	Status int         // HTTP状态码
	Body   []byte      // 响应体
	Header http.Header // 响应头
}

// StatusOK 检查响应状态码是否为 200 OK
func (r *Response) StatusOK() bool {
	return r.Status == http.StatusOK
}

// StatusIn 检查响应状态码是否在指定状态码列表中，如果未提供状态码列表，则默认检查status是否为200
func (r *Response) StatusIn(status ...int) bool {
	if len(status) == 0 {
		return r.StatusOK()
	}
	for _, v := range status {
		if v == r.Status {
			return true
		}
	}
	return false
}

// Unmarshal 将响应体解析到指定的结构体中
func (r *Response) Unmarshal(v any) error {
	if len(r.Body) == 0 {
		return errorx.New("response Body is empty, cannot unmarshal")
	}
	if err := json.Unmarshal(r.Body, v); err != nil {
		return errorx.Wrap(err, "json unmarshal error")
	}
	return nil
}

// GetCookies 从响应头中提取所有 Cookie
func (r *Response) GetCookies() []*http.Cookie {
	values := r.Header.Values("Set-Cookie")
	if len(values) == 0 {
		return nil
	}
	var cookies []*http.Cookie
	for _, value := range values {
		if cookie, err := http.ParseSetCookie(value); err == nil {
			cookies = append(cookies, cookie)
		}
	}
	return cookies
}

// GetCookie 从响应头中提取指定名称的 Cookie
func (r *Response) GetCookie(name string) *http.Cookie {
	values := r.Header.Values("Set-Cookie")
	if len(values) == 0 {
		return nil
	}
	for _, value := range values {
		if cookie, err := http.ParseSetCookie(value); err == nil && cookie.Name == name {
			return cookie
		}
	}
	return nil
}
