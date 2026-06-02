package errorx

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
)

// New 创建 error
func New(v any) error {
	err := &Error{stack: getStack()}
	switch e := v.(type) {
	case error:
		err.source = e
		err.msg = e.Error()
	case string:
		err.msg = e
	default:
		err.msg = fmt.Sprint(e)
	}
	return err
}

// Sprintf 格式化创建 error
func Sprintf(format string, a ...interface{}) error {
	return &Error{
		msg:   fmt.Sprintf(format, a...),
		stack: getStack(),
	}
}

// Wrap 包装 error
func Wrap(v any, msg string) error {
	err := &Error{msg: msg}
	switch e := v.(type) {
	case *Error:
		err.source = e
		err.stack = e.stack
	case error:
		err.source = e
		err.stack = getStack()
	default:
		err.source = New(e)
		err.stack = getStack()
	}
	return err
}

// Wrapf 格式化包装 error
func Wrapf(v any, format string, a ...interface{}) error {
	return Wrap(v, fmt.Sprintf(format, a...))
}

// Unwrap 解包装
func Unwrap(err error) error {
	if in, ok := err.(interface {
		Unwrap() error
	}); ok {
		return in.Unwrap()
	}
	return nil
}

// Is 判断 error 是否匹配（Go 1.13 errors.Is 兼容）
func Is(err, target error) bool {
	if err == nil && target == nil {
		return true
	}
	if err == nil || target == nil {
		return false
	}

	// 检查是否相等
	if errors.Is(err, target) {
		return true
	}

	// 检查 target 是否是 err 的包装链中的一环
	for err != nil {
		if errors.Is(err, target) {
			return true
		}
		err = Unwrap(err)
	}
	return false
}

// As 判断 error 链中是否包含指定类型（Go 1.13 errors.As 兼容）
func As(err error, target interface{}) bool {
	if err == nil || target == nil {
		return false
	}

	for err != nil {
		if match(err, target) {
			return true
		}
		err = Unwrap(err)
	}
	return false
}

func match(err error, target interface{}) bool {
	switch t := target.(type) {
	case *error:
		*t = err
		return true
	case interface{ SetError(error) }:
		t.SetError(err)
		return true
	default:
		return false
	}
}

// Panic 断言 error 是否为 nil，若不为 nil 则 panic
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

// Error 通用 error
type Error struct {
	source error  // 源 error
	msg    string // 报错信息
	stack  stack  // 调用栈
}

// Error 实现 error 接口
func (err *Error) Error() string {
	sb := new(strings.Builder)
	sb.WriteString(err.msg)
	if err.source != nil {
		sb.WriteString(" | ")
		sb.WriteString(err.source.Error())
	}
	return sb.String()
}

// Unwrap 解包装
func (err *Error) Unwrap() error { return err.source }

// Format fmt 打印实现
func (err *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 118: // verb == 'v'
		_, _ = io.WriteString(s, err.Error())
		err.stack.Format(s, verb)
	default:
		_, _ = io.WriteString(s, err.Error())
	}
}

// MarshalJSON 序列化 json 实现
func (err *Error) MarshalJSON() ([]byte, error) {
	msg := err.Error()
	bytes := make([]byte, 0, len(msg)+2)
	bytes = append(bytes, 34)
	bytes = append(bytes, []byte(msg)...)
	bytes = append(bytes, 34)
	return bytes, nil
}

// 调用栈
type stack []uintptr

// Format 打印调用栈信息（fmt 实现）
func (s *stack) Format(f fmt.State, verb rune) {
	if verb == 118 { // verb == 'v'
		i, frames := 1, runtime.CallersFrames(*s)
		for {
			if pc, more := frames.Next(); more && i <= 5 {
				// 使用 strings.Builder 减少字符串拼接的内存分配
				var builder strings.Builder
				builder.Grow(100) // 预估容量
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("%d", i))
				builder.WriteString(" : ")
				builder.WriteString(pc.Function)
				builder.WriteString(" >> ")
				builder.WriteString(pc.File)
				builder.WriteString(":")
				builder.WriteString(fmt.Sprintf("%d", pc.Line))
				_, _ = io.WriteString(f, builder.String())
				i++
			} else {
				break
			}
		}
	}
}

// getStack 获取调用栈
func getStack() []uintptr {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[:n-1]
}
