package errorx

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

// New 创建error
func New(v any) error {
	var err = &Error{stack: getStack()}
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

// Newf 格式化创建error
func Newf(format string, a ...interface{}) error {
	return &Error{
		msg:   fmt.Sprintf(format, a...),
		stack: getStack(),
	}
}

// Wrap 包装error
func Wrap(v any, msg string) error {
	var err = &Error{msg: msg}
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

// Unwrap 解包装
func Unwrap(err error) error {
	if t, ok := err.(interface {
		Unwrap() error
	}); ok {
		return t.Unwrap()
	} else {
		return nil
	}
}

// Panic 恐慌
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

// Error 通用error
type Error struct {
	source error  // 源error
	msg    string // 报错信息
	stack  stack  // 调用栈
}

// 报错信息（用以实现error接口）
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

// Format fmt打印实现
func (err *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		_, _ = io.WriteString(s, err.Error())
		err.stack.Format(s, verb)
	case 's':
		_, _ = io.WriteString(s, err.Error())
	}
}

// MarshalJSON 序列化json实现
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

// Format 打印调用栈信息（fmt实现）
func (s *stack) Format(f fmt.State, verb rune) {
	if verb == 'v' {
		i, frames := 1, runtime.CallersFrames(*s)
		for {
			if pc, more := frames.Next(); more && i <= 5 {
				// 使用strings.Builder减少字符串拼接的内存分配
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
