package httpx

import "net/http"

func NewError(status int, code int, msg string, data any) *Error {
	return &Error{
		status: status,
		code:   code,
		msg:    msg,
		data:   data,
	}
}

// Error 错误结构体
type Error struct {
	status int    // http 状态码
	code   int    // 业务状态码
	msg    string // 错误信息
	data   any    // 错误数据
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Status() int {
	return e.status
}

func (e *Error) SetStatus(status int) *Error {
	e.status = status
	return e
}

func (e *Error) StatusOK() bool {
	return e.status == http.StatusOK
}
