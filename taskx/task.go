package taskx

import "context"

// Task 任务接口
type Task interface {
	GetID() string                 // 任务ID
	Execute(context.Context) error // 带上下文的执行函数
}
