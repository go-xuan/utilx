package taskx

import (
	"context"
	"fmt"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// NewRetryScheduler 创建重试任务调度器
func NewRetryScheduler(task Task, retry *Retry) *RetryScheduler {
	return &RetryScheduler{
		retry: retry,
		task:  task,
	}
}

// RetryScheduler 重试任务调度器
type RetryScheduler struct {
	retry *Retry // 重试策略
	task  Task   // 任务执行函数
}

func (t *RetryScheduler) GetID() string {
	return fmt.Sprintf("retry:%s", t.task.GetID())
}

func (t *RetryScheduler) Execute(ctx context.Context) error {
	logger := log.WithField("task_id", t.task.GetID()).
		WithField("task_type", "retry_scheduler")

	if t.task == nil {
		logger.Error("retry scheduler task is nil")
		return errorx.New("retry scheduler task is nil")
	}
	if t.retry == nil {
		return errorx.New("retry scheduler retry is nil")
	}
	if err := t.retry.Execute(ctx, t.task); err != nil {
		logger.WithField("error", err.Error()).
			Error("retry scheduler execute error")
		return errorx.Wrap(err, "retry scheduler execute error")
	}
	logger.Info("retry scheduler execute success")
	return nil
}
