package taskx

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// ResultCallback 任务结果回调函数
type ResultCallback func(Result)

// ErrorLogCallback 错误日志回调
func ErrorLogCallback(result Result) {
	if err := result.GetError(); err != nil {
		log.WithError(err).Error("任务执行失败")
	}
}

// RetryCallback 重试失败任务
func RetryCallback(times int, interval time.Duration) ResultCallback {
	retry := NewRetry(times, interval)
	return func(result Result) {
		if err := result.GetError(); err != nil {
			log.WithError(err).Error("任务执行失败")
			// 任务重试
			if err = retry.Execute(context.Background(), result.GetExecute()); err != nil {
				log.WithError(err).Error("任务重试失败")
			}
		}
	}
}

// Result 任务结果接口
type Result interface {
	GetExecute() Execute // 获取任务执行函数
	GetError() error     // 获取任务执行错误
}

// ExecuteResult 函数执行结果
type ExecuteResult struct {
	execute Execute
	error   error
}

func (r *ExecuteResult) GetExecute() Execute {
	return r.execute
}

func (r *ExecuteResult) GetError() error {
	return r.error
}

// TaskResult 函数执行结果
type TaskResult struct {
	task  Task
	error error
}

func (r *TaskResult) GetUnique() string {
	return r.task.GetUnique()
}

func (r *TaskResult) GetExecute() Execute {
	return r.task.Execute
}

func (r *TaskResult) GetError() error {
	return r.error
}
