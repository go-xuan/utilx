package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
)

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
	ApplicationForm = "application/x-www-form-urlencoded"
)

// NewRequest 新建请求
func NewRequest(method string, url_ string) *Request {
	return &Request{
		method:     method,
		url:        url_,
		headers:    make(map[string]string),
		cookies:    make([]*http.Cookie, 0),
		form:       make(url.Values),
		files:      make([]*File, 0),
		decorators: make([]Decorator, 0),
	}
}

// File 上传的文件
type File struct {
	Field  string            // 表单字段名
	Name   string            // 文件名
	Data   []byte            // 文件内容
	Params map[string]string // 文件参数
}

// Request 请求器
type Request struct {
	method     string            // 请求方法
	url        string            // 请求URL
	trace      string            // 跟踪ID
	debug      bool              // 是否开启调试模式
	headers    map[string]string // 请求头
	cookies    []*http.Cookie    // 请求cookie
	form       url.Values        // 请求表单参数
	files      []*File           // 请求文件
	body       any               // 请求体
	decorators []Decorator       // 请求装饰器
}

// SetMethod 设置请求方法
func (r *Request) SetMethod(method string) *Request {
	if method != "" {
		r.method = method
	}
	return r
}

// SetURL 设置请求URL
func (r *Request) SetURL(url_ string) *Request {
	if url_ != "" {
		r.url = url_
	}
	return r
}

// AddTrace 添加跟踪ID
func (r *Request) AddTrace(trace string) *Request {
	r.trace = trace
	return r
}

// Debug 开启调试模式
func (r *Request) Debug() *Request {
	r.debug = true
	return r
}

// AddParams 添加查询参数
func (r *Request) AddParams(params map[string]string) *Request {
	if len(params) == 0 {
		return r
	}
	// 解析现有URL查询参数
	parse, err := url.Parse(r.url)
	if err != nil {
		return r
	}
	values := parse.Query()
	for k, v := range params {
		values.Add(k, v)
	}
	parse.RawQuery = values.Encode()
	r.url = parse.String()
	return r
}

// AddBody 添加请求体
func (r *Request) AddBody(body any) *Request {
	r.body = body
	return r
}

// AddForm 添加请求表单参数
func (r *Request) AddForm(form url.Values) *Request {
	if len(form) > 0 {
		for key, vals := range form {
			for _, val := range vals {
				r.form.Add(key, val)
			}
		}
	}
	return r
}

// AddFile 添加请求文件
func (r *Request) AddFile(file *File) *Request {
	if file != nil {
		r.files = append(r.files, file)
	}
	return r
}

// AddHeader 添加请求头
func (r *Request) AddHeader(key, value string) *Request {
	if key != "" && value != "" {
		r.headers[key] = value
	}
	return r
}

// AddHeaders 添加请求头
func (r *Request) AddHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.AddHeader(k, v)
	}
	return r
}

// AddCookie 添加请求cookie
func (r *Request) AddCookie(cookie *http.Cookie) *Request {
	if cookie != nil {
		r.cookies = append(r.cookies, cookie)
	}
	return r
}

// AddCookies 添加请求cookie
func (r *Request) AddCookies(cookies []*http.Cookie) *Request {
	if len(cookies) > 0 {
		r.cookies = append(r.cookies, cookies...)
	}
	return r
}

// AddDecorator 添加请求装饰器
func (r *Request) AddDecorator(decorators ...Decorator) *Request {
	r.decorators = append(r.decorators, decorators...)
	return r
}

// NewHttpRequest 创建http请求
func (r *Request) NewHttpRequest() (*http.Request, error) {
	body, err := r.getBodyReader()
	if err != nil {
		return nil, errorx.Wrap(err, "get body reader error")
	}

	// 创建请求
	var request *http.Request
	if request, err = http.NewRequest(r.method, r.url, body); err != nil {
		return nil, errorx.Wrap(err, "new request error")
	}

	// 添加请求预设
	r.decorate(request)
	// 添加请求头
	r.setRequestHeaders(request)
	// 添加cookie
	r.setRequestCookie(request)

	return request, nil
}

// getBodyReader 获取请求体读取器
func (r *Request) getBodyReader() (io.Reader, error) {
	if r.body != nil {
		r.AddHeader(ContentType, ApplicationJSON)
		b, err := json.Marshal(r.body)
		if err != nil {
			return nil, errorx.Wrap(err, "marshal body error")
		}
		return bytes.NewReader(b), nil
	} else if len(r.form) > 0 {
		r.AddHeader(ContentType, ApplicationForm)
		return strings.NewReader(r.form.Encode()), nil
	} else if len(r.files) > 0 {
		buffer := &bytes.Buffer{}
		writer := multipart.NewWriter(buffer)

		for _, file := range r.files {
			wf, err := writer.CreateFormFile(file.Field, file.Name)
			if err != nil {
				return nil, errorx.Wrap(err, "create form file error")
			}
			if _, err = wf.Write(file.Data); err != nil {
				return nil, errorx.Wrap(err, "write form file error")
			}
			if file.Params != nil && len(file.Params) > 0 {
				for k, v := range file.Params {
					if err = writer.WriteField(k, v); err != nil {
						return nil, errorx.Wrap(err, fmt.Sprintf("write file params [%s:%s] error", k, v))
					}
				}
			}
		}

		r.AddHeader(ContentType, writer.FormDataContentType())
		if err := writer.Close(); err != nil {
			return nil, errorx.Wrap(err, "close multipart writer error")
		}
		return buffer, nil
	}
	return nil, nil
}

// 请求装饰
func (r *Request) decorate(request *http.Request) {
	for _, decorator := range r.decorators {
		if decorator != nil {
			decorator.Decorate(request)
		}
	}
}

// 设置请求头
func (r *Request) setRequestHeaders(request *http.Request) {
	if r.headers != nil && len(r.headers) > 0 {
		for key, value := range r.headers {
			request.Header.Set(key, value)
		}
	}
}

// 设置cookie
func (r *Request) setRequestCookie(request *http.Request) {
	if len(r.cookies) > 0 {
		for _, cookie := range r.cookies {
			request.AddCookie(cookie)
		}
	}
}

// Send 发送请求
func (r *Request) Send(client ...*http.Client) (*Response, error) {
	// 发送请求
	var response *Response
	if request, err := r.NewHttpRequest(); err != nil {
		return nil, errorx.Wrap(err, "new http request error")
	} else if response, err = SendHttpRequest(request, client...); err != nil {
		return response, errorx.Wrap(err, "send http request error")
	}
	// 添加trace
	response.AddTrace(r.trace)
	// 打印调试信息
	if r.debug {
		logger := log.WithField("http_debug", true)
		if trace := r.trace; trace != "" {
			logger = logger.WithField("trace", trace)
		}
		logger.Printf("http_url: %s", r.url)
		logger.Printf("http_body: %s", string(response.body))
	}
	// 关联trace
	return response, nil
}
