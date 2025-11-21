package taskx

import (
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/marshalx"
)

// Error 任务错误
type Error struct {
	ID  string `json:"id"`  // 任务唯一标识
	Msg string `json:"msg"` // 任务错误信息
}

// Error 任务错误实现error接口
func (e *Error) Error() string {
	return e.Msg
}

// ErrorsCollect 任务错误数据收集器
type ErrorsCollect struct {
	Ids    []string `json:"ids"`    // 任务ID
	Errors []*Error `json:"errors"` // 失败任务错误
}

// ResultHook 任务结果处理钩子函数
func (c *ErrorsCollect) ResultHook(result IResult) {
	id := result.GetID()
	c.Ids = append(c.Ids, id)
	if err := result.GetError(); err != nil {
		c.Errors = append(c.Errors, &Error{
			ID:  id,
			Msg: err.Error(),
		})
	}
}

// Save 保存任务错误数据
func (c *ErrorsCollect) Save(path string) error {
	if c == nil || path == "" {
		return nil
	}
	if err := marshalx.Apply(path).Write(path, c); err != nil {
		return errorx.Wrap(err, "save tasks collect")
	}
	return nil
}
