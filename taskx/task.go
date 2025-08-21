package taskx

import (
	"context"
)

// Task 任务接口
type Task interface {
	GetUnique() string             // 获取任务唯一标识
	Execute(context.Context) error // 执行函数
}

// Execute 任务执行函数
type Execute func(context.Context) error

// BatchExecute 批量执行函数
type BatchExecute[T any] func(context.Context, []T) error
