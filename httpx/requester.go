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

func NewRequest(method string, url_ string) *Request {
	return &Request{
		method:  method,
		url:     url_,
		headers: make(map[string]string),
		cookies: make([]*http.Cookie, 0),
		form:    make(url.Values),
		files:   make([]*File, 0),
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
	method  string
	url     string
	trace   string
	debug   bool
	headers map[string]string
	cookies []*http.Cookie
	form    url.Values
	files   []*File
	body    any
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

func (r *Request) AddBody(body any) *Request {
	r.body = body
	return r
}

func (r *Request) AddForm(form url.Values) *Request {
	r.form = form
	return r
}

func (r *Request) AddFile(file *File) *Request {
	r.files = append(r.files, file)
	return r
}

func (r *Request) AddHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

func (r *Request) AddHeader(key, value string) *Request {
	r.headers[key] = value
	return r
}

func (r *Request) Debug() *Request {
	r.debug = true
	return r
}

func (r *Request) Trace(trace string) *Request {
	r.trace = trace
	return r
}

func (r *Request) GetBodyReader() (io.Reader, error) {
	if r.body != nil {
		r.headers["Content-Type"] = "application/json"
		marshal, err := json.Marshal(r.body)
		if err != nil {
			return nil, errorx.Wrap(err, "marshal body error")
		}
		return bytes.NewReader(marshal), nil
	} else if len(r.form) > 0 {
		r.headers["Content-Type"] = "application/x-www-form-urlencoded"
		return strings.NewReader(r.form.Encode()), nil
	} else if len(r.files) > 0 {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

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

		if err := writer.Close(); err != nil {
			return nil, errorx.Wrap(err, "close multipart writer error")
		}

		r.headers["Content-Type"] = writer.FormDataContentType()
		return body, nil
	}
	return nil, nil
}

// Send 发送请求
func (r *Request) Send(options ...SettingsOption) (*Response, error) {
	body, err := r.GetBodyReader()
	if err != nil {
		return nil, errorx.Wrap(err, "get body reader error")
	}
	var request *http.Request
	if request, err = http.NewRequest(r.method, r.url, body); err != nil {
		return nil, errorx.Wrap(err, "new request error")
	}
	// 设置请求头
	if r.headers != nil && len(r.headers) > 0 {
		for key, val := range r.headers {
			request.Header.Set(key, val)
		}
	}
	// 添加cookie
	if len(r.cookies) > 0 {
		for _, cookie := range r.cookies {
			request.AddCookie(cookie)
		}
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
			logger = logger.WithField("trace", trace)
		}
		logger.Printf("http_url: %s", r.url)
		logger.Printf("http_body: %s", string(response.Body()))
	}
	return response, nil
}
