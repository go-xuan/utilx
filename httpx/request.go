package httpx

import (
	"bytes"
	"encoding/json"
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
		method:  method,
		url:     url_,
		headers: make(map[string]string),
		cookies: make([]*http.Cookie, 0),
		form:    make(url.Values),
		files:   make([]*File, 0),
		presets: make([]RequestPreset, 0),
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
	method  string            // 请求方法
	url     string            // 请求URL
	trace   string            // 跟踪ID
	debug   bool              // 是否开启调试模式
	headers map[string]string // 请求头
	cookies []*http.Cookie    // 请求cookie
	form    url.Values        // 请求表单参数
	files   []*File           // 请求文件
	body    any               // 请求体
	presets []RequestPreset   // 请求预设
}

// Params 添加查询参数
func (r *Request) Params(params map[string]string) *Request {
	if len(params) == 0 {
		return r
	}

	// 解析现有URL查询参数
	parsedURL, err := url.Parse(r.url)
	if err != nil {
		return r
	}

	query := parsedURL.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	parsedURL.RawQuery = query.Encode()
	r.url = parsedURL.String()

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
	r.files = append(r.files, file)
	return r
}

// AddHeaders 添加请求头
func (r *Request) AddHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

// AddHeader 添加请求头
func (r *Request) AddHeader(key, value string) *Request {
	r.headers[key] = value
	return r
}

// AddPreset 添加请求预设
func (r *Request) AddPreset(presets ...RequestPreset) *Request {
	r.presets = append(r.presets, presets...)
	return r
}

// Debug 开启调试模式
func (r *Request) Debug() *Request {
	r.debug = true
	return r
}

// Trace 添加跟踪ID
func (r *Request) Trace(trace string) *Request {
	r.trace = trace
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
	r.setRequestPreset(request)
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
		marshal, err := json.Marshal(r.body)
		if err != nil {
			return nil, errorx.Wrap(err, "marshal body error")
		}
		return bytes.NewReader(marshal), nil
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
					_ = writer.WriteField(k, v)
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

// 设置请求预设
func (r *Request) setRequestPreset(request *http.Request) {
	for _, preset := range r.presets {
		preset.Preset(request)
	}
}

// 设置请求头
func (r *Request) setRequestHeaders(request *http.Request) {
	if r.headers != nil && len(r.headers) > 0 {
		for key, val := range r.headers {
			request.Header.Set(key, val)
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
func (r *Request) Send(options ...SettingsOption) (*Response, error) {
	request, err := r.NewHttpRequest()
	if err != nil {
		return nil, errorx.Wrap(err, "new request error")
	}
	// 发送请求
	var response *Response
	if response, err = GetClient().Do(request, options...); err != nil {
		return response, errorx.Wrap(err, "do request error")
	}

	// 打印调试信息
	if r.debug {
		logger := log.WithField("http_debug", true)
		if trace := r.trace; trace != "" {
			logger = logger.WithField("Trace", trace)
		}
		logger.Printf("http_url: %s", r.url)
		logger.Printf("http_body: %s", string(response.Body))
	}
	// 关联trace
	response.Trace = r.trace
	return response, nil
}
