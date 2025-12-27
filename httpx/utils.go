package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/stringx"
)

const (
	HTTP  = "http"  // http协议
	HTTPS = "https" // https协议
)

// NewJsonReader 获取json读取器
func NewJsonReader(body any) (io.Reader, string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, "", errorx.Wrap(err, "json marshal error")
	}
	return bytes.NewReader(b), ApplicationJSON, nil
}

// NewFormReader 获取表单读取器
func NewFormReader(form url.Values) (io.Reader, string, error) {
	return strings.NewReader(form.Encode()), ApplicationForm, nil
}

// NewFileReader 获取文件读取器
func NewFileReader(files []*File) (io.Reader, string, error) {
	reader := &bytes.Buffer{}
	writer := multipart.NewWriter(reader)
	for _, file := range files {
		wf, err := writer.CreateFormFile(file.Field, file.Name)
		if err != nil {
			return nil, "", errorx.Wrap(err, "create form file error")
		}
		if _, err = wf.Write(file.Data); err != nil {
			return nil, "", errorx.Wrap(err, "write form file error")
		}
		if file.Params != nil && len(file.Params) > 0 {
			for k, v := range file.Params {
				if err = writer.WriteField(k, v); err != nil {
					return nil, "", errorx.Wrap(err, fmt.Sprintf("write file params [%s:%s] error", k, v))
				}
			}
		}
	}
	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return nil, "", errorx.Wrap(err, "close multipart writer error")
	}
	return reader, contentType, nil
}

// ParseResponseBody 解析响应体
func ParseResponseBody(resp io.ReadCloser, v any) error {
	body, err := io.ReadAll(resp)
	if err != nil {
		return errorx.Wrap(err, "read response body error")
	}
	return json.Unmarshal(body, v)
}

// HasProtocol 检查url是否包含协议头
func HasProtocol(url string) (string, bool) {
	protocols := []string{HTTPS, HTTP}
	for _, protocol := range protocols {
		if strings.HasPrefix(url, protocol) {
			return protocol, true
		}
	}
	return "", false
}

// AddProtocol 添加协议头
func AddProtocol(url string, protocol ...string) string {
	if _, ok := HasProtocol(url); !ok {
		prot := stringx.Default(HTTP, protocol...)
		return fmt.Sprintf("%s://%s", prot, url)
	}
	return url
}

// ParseHost 解析host
func ParseHost(host_ string) (string, string, int) {
	if host_ = strings.TrimSpace(host_); host_ == "" {
		return "", "", 0
	}
	protocol, _ := HasProtocol(host_)
	_, host_ = stringx.Cut(host_, "://")
	host_, _ = stringx.Cut(host_, "/")
	host, port := stringx.Cut(host_, ":")
	return protocol, host, stringx.ParseInt(port)
}

// DownloadFile 下载文件
func DownloadFile(url string) ([]byte, error) {
	resp, err := NewClient().Get(url)
	if err != nil {
		return nil, errorx.Wrap(err, "http request error")
	}
	defer errorx.Close(resp.Body)

	var data []byte
	if data, err = io.ReadAll(resp.Body); err != nil {
		return nil, errorx.Wrap(err, "read response body error")
	}
	return data, nil
}

// DownloadFileTo 下载文件到指定路径
func DownloadFileTo(url string, path string) error {
	if data, err := DownloadFile(url); err != nil {
		return errorx.Wrap(err, "download file error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}
