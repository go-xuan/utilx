package taskx

import (
	"context"
)

// Task 任务接口
type Task interface {
	GetID() string                 // 任务ID
	Execute(context.Context) error // 执行函数
}

// Execute 任务执行函数
type Execute func(context.Context) error

// Wrap 任务执行函数包装器
type Wrap func(Execute) Execute

// Error 任务错误
type Error struct {
	ID  string `json:"id"`  // 任务唯一标识
	Msg string `json:"msg"` // 任务错误信息
}

func (e *Error) Error() string {
	return e.Msg
}
