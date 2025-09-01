package taskx

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// Execute 任务执行函数
type Execute func(context.Context) error

// Wrap 任务执行函数包装器
type Wrap func(Execute) Execute

// BatchExecute 批量执行函数
type BatchExecute[T any] func(ctx context.Context, tasks []T) error

// Task 任务接口
type Task interface {
	GetUnique() string             // 获取任务唯一标识
	Execute(context.Context) error // 执行函数
}

// Result 任务结果接口
type Result interface {
	GetUnique() string   // 获取任务唯一标识
	GetExecute() Execute // 获取任务执行函数
	GetError() error     // 获取任务执行错误
}

// ResultHook 任务结果处理钩子函数
type ResultHook func(Result)

// TaskResult 任务执行结果
type TaskResult struct {
	unique  string  // 任务唯一标识
	execute Execute // 任务执行函数
	error   error   // 任务执行错误
}

func (r *TaskResult) GetUnique() string {
	return r.unique
}

func (r *TaskResult) GetExecute() Execute {
	return r.execute
}

func (r *TaskResult) GetError() error {
	return r.error
}

// ErrorLogHook 错误日志钩子函数
func ErrorLogHook(result Result) {
	if err := result.GetError(); err != nil {
		log.WithError(err).Error("task execute error")
	}
}
