package httpx

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
	ApplicationForm = "application/x-www-form-urlencoded"
)

// Get 发送GET请求
func Get(url string) (*Response, error) {
	return NewRequest(http.MethodGet, url).Send()
}

// NewRequest 新建请求
func NewRequest(method string, url_ string) *Request {
	return &Request{
		client: NewClient(),
		method: method,
		url:    url_,
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
	client     *http.Client      // http客户端
	method     string            // 请求方法
	url        string            // 请求url
	trace      string            // 跟踪id
	debug      bool              // 是否开启调试模式
	headers    map[string]string // 请求头
	cookies    []*http.Cookie    // 请求cookie
	form       url.Values        // 请求表单参数
	files      []*File           // 请求文件
	body       any               // 请求体
	decorators []Decorator       // 请求装饰器
}

// Send 发送请求
func (r *Request) Send() (*Response, error) {
	// 发送请求
	var resp *http.Response
	if request, err := r.Build(); err != nil {
		return nil, errorx.Wrap(err, "new http request error")
	} else if resp, err = r.GetClient().Do(request); err != nil {
		return nil, errorx.Wrap(err, "http client do error")
	}
	defer errorx.Close(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errorx.Wrap(err, "http response body read error")
	}
	// 打印调试信息
	if r.debug {
		log.WithField("trace", r.trace).
			WithField("url", r.url).
			WithField("body", string(body)).
			Print("request-debug")
	}
	return &Response{
		trace:  r.trace,
		status: resp.StatusCode,
		body:   body,
		header: resp.Header,
	}, nil
}

// Build 创建http请求
func (r *Request) Build() (*http.Request, error) {
	// 创建请求
	var request *http.Request
	if reader, err := r.getBodyReader(); err != nil {
		return nil, errorx.Wrap(err, "get body reader error")
	} else if request, err = http.NewRequest(r.method, r.url, reader); err != nil {
		return nil, errorx.Wrap(err, "build request error")
	}

	// 添加请求预设
	r.decorate(request)
	// 添加请求头
	r.setRequestHeaders(request)
	// 添加cookie
	r.setRequestCookie(request)

	return request, nil
}

// GetClient 获取http客户端
func (r *Request) GetClient() *http.Client {
	if r.client == nil {
		r.client = NewClient()
	}
	return r.client
}

// SetClient 设置http客户端
func (r *Request) SetClient(client *http.Client) *Request {
	if client != nil {
		r.client = client
	}
	return r
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

// AddParam 添加查询参数
func (r *Request) AddParam(key, value string) *Request {
	if parse, err := url.Parse(r.url); err == nil {
		values := parse.Query()
		values.Add(key, value)
		parse.RawQuery = values.Encode()
		r.url = parse.String()
	}
	return r
}

// AddParams 添加查询参数
func (r *Request) AddParams(params map[string]string) *Request {
	if len(params) == 0 {
		return r
	}
	// 解析现有url查询参数
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
	if body != nil {
		r.body = body
	}
	return r
}

// AddForm 添加请求表单参数
func (r *Request) AddForm(form url.Values) *Request {
	if len(form) > 0 {
		if r.form == nil {
			r.form = make(url.Values)
		}
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
		if r.files == nil {
			r.files = make([]*File, 0)
		}
		r.files = append(r.files, file)
	}
	return r
}

// AddHeader 添加请求头
func (r *Request) AddHeader(key, value string) *Request {
	if key != "" && value != "" {
		if r.headers == nil {
			r.headers = make(map[string]string)
		}
		r.headers[key] = value
	}
	return r
}

// AddHeaders 添加请求头
func (r *Request) AddHeaders(headers map[string]string) *Request {
	if len(headers) > 0 {
		for k, v := range headers {
			r.AddHeader(k, v)
		}
	}
	return r
}

// AddCookie 添加请求cookie
func (r *Request) AddCookie(cookie *http.Cookie) *Request {
	if cookie != nil {
		if r.cookies == nil {
			r.cookies = make([]*http.Cookie, 0)
		}
		r.cookies = append(r.cookies, cookie)
	}
	return r
}

// AddCookies 添加请求cookie
func (r *Request) AddCookies(cookies []*http.Cookie) *Request {
	if len(cookies) > 0 {
		if r.cookies == nil {
			r.cookies = make([]*http.Cookie, 0)
		}
		r.cookies = append(r.cookies, cookies...)
	}
	return r
}

// AddDecorator 添加请求装饰器
func (r *Request) AddDecorator(decorators ...Decorator) *Request {
	if len(decorators) > 0 {
		if r.decorators == nil {
			r.decorators = make([]Decorator, 0)
		}
		r.decorators = append(r.decorators, decorators...)
	}
	return r
}

// getBodyReader 获取请求体读取器
func (r *Request) getBodyReader() (io.Reader, error) {
	if r.body != nil {
		reader, contentType, err := NewJsonReader(r.body)
		if err != nil {
			return nil, errorx.Wrap(err, "get form reader error")
		}
		r.AddHeader(ContentType, contentType)
		return reader, nil
	}
	if len(r.form) > 0 {
		reader, contentType, err := NewFormReader(r.form)
		if err != nil {
			return nil, errorx.Wrap(err, "get form reader error")
		}
		r.AddHeader(ContentType, contentType)
		return reader, nil
	}
	if len(r.files) > 0 {
		reader, contentType, err := NewFileReader(r.files)
		if err != nil {
			return nil, errorx.Wrap(err, "get form reader error")
		}
		r.AddHeader(ContentType, contentType)
		return reader, nil
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
