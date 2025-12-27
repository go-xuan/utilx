package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-xuan/utilx/errorx"
)

// HttpRequest 发送HTTP请求并返回响应
func HttpRequest(client *http.Client, request *http.Request) (*Response, error) {
	resp, err := client.Do(request)
	if err != nil {
		return nil, errorx.Wrap(err, "http client do request error")
	}
	defer errorx.Close(resp.Body)

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, errorx.Wrap(err, "http response body read error")
	}
	return &Response{
		status: resp.StatusCode,
		body:   body,
		header: resp.Header,
	}, nil
}

// Response 表示 HTTP 请求的响应结构
type Response struct {
	trace  string      // 追踪标识，用于日志或调试
	status int         // HTTP状态码
	body   []byte      // 响应体
	header http.Header // 响应头
}

// GetTrace 获取响应追踪标识
func (r *Response) GetTrace() string {
	return r.trace
}

// GetStatus 获取响应状态码
func (r *Response) GetStatus() int {
	return r.status
}

// GetBody 获取响应体
func (r *Response) GetBody() []byte {
	return r.body
}

// AddTrace 添加响应trace
func (r *Response) AddTrace(trace string) {
	if trace != "" {
		r.trace = trace
	}
}

// StatusOK 检查响应状态码是否为 200 OK
func (r *Response) StatusOK() bool {
	return r.status == http.StatusOK
}

// StatusMatch 响应状态码匹配
func (r *Response) StatusMatch(status ...int) bool {
	if len(status) == 0 {
		return r.StatusOK()
	}
	for _, v := range status {
		if v == r.status {
			return true
		}
	}
	return false
}

// Unmarshal 将响应体解析到指定的结构体中
func (r *Response) Unmarshal(v any) error {
	if v != nil {
		return nil
	}
	if len(r.body) == 0 {
		return errorx.New("response body is empty")
	}
	if err := json.Unmarshal(r.body, v); err != nil {
		return errorx.Wrap(err, "response body unmarshal error")
	}
	return nil
}

// UnmarshalField 将响应体解析到指定的结构体中
func (r *Response) UnmarshalField(field string, v any) error {
	if v == nil {
		return nil
	}
	if len(r.body) == 0 {
		return errorx.New("response body is empty")
	}
	result := make(map[string]json.RawMessage)
	if err := json.Unmarshal(r.body, &result); err != nil {
		return errorx.Wrap(err, "response body unmarshal error")
	}
	if data, ok := result[field]; !ok {
		return errorx.Sprintf("field %s not found", field)
	} else if err := json.Unmarshal(data, v); err != nil {
		return errorx.Wrap(err, fmt.Sprintf("field %s unmarshal error", field))
	}
	return nil
}

// Write 将响应体写入 io.Writer
func (r *Response) Write(w io.Writer) error {
	if body := r.body; body != nil {
		if _, err := w.Write(body); err != nil {
			return errorx.Wrap(err, "response body write error")
		}
	}
	return nil
}

// GetCookies 从响应头中提取所有 Cookie
func (r *Response) GetCookies() []*http.Cookie {
	values := r.header.Values("Set-Cookie")
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
	values := r.header.Values("Set-Cookie")
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
