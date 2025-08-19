package taskx

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// Task 任务通用接口
type Task interface {
	Execute(ctx context.Context) error
}

// Result 任务结果接口
type Result interface {
	Task() Task
	Error() error
}

// Wrapper 任务包装器
type Wrapper func(Task) Task

// ResultCallback 任务结果回调函数
type ResultCallback func(Result)

// LogResult 日志记录任务结果
func LogResult(result Result) {
	if err := result.Error(); err != nil {
		log.WithError(err).Error("任务执行失败")
	}
}

// RetryCallback 重试失败任务
func RetryCallback(times int, interval time.Duration) ResultCallback {
	retry := NewRetry(times, interval)
	return func(result Result) {
		if err := result.Error(); err != nil {
			log.WithError(err).Error("任务执行失败")
			// 任务重试
			if err = retry.Execute(context.Background(), result.Task()); err != nil {
				log.WithError(err).Error("任务重试失败")
			}
		}
	}
}
