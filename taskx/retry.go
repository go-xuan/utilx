package taskx

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
)

// NewRetry 创建重试任务
func NewRetry(task Task, times int, interval time.Duration) *Retry {
	return &Retry{
		task:     task,
		strategy: NewRetryStrategy(times, interval),
	}
}

// Retry 重试任务包装
type Retry struct {
	task     Task           // 任务执行函数
	strategy *RetryStrategy // 重试策略
}

func (t *Retry) GetID() string {
	return fmt.Sprintf("retry@%s", t.task.GetID())
}

func (t *Retry) Execute(ctx context.Context) error {
	logger := log.WithField("id", t.GetID())

	if t.task == nil {
		logger.Error("retry task is nil")
		return errorx.New("retry task is nil")
	} else if t.strategy == nil {
		logger.Error("retry strategy is nil")
		return errorx.New("retry strategy is nil")
	}

	times, err := t.strategy.Execute(ctx, t.task)
	if err != nil {
		logger.WithError(err).Error("retry task execute error")
		return errorx.Wrap(err, "retry task execute error")
	}
	logger.WithField("retry_times", times).Info("retry task execute success")
	return nil
}

// NewRetryStrategy 创建重试策略
func NewRetryStrategy(times int, interval time.Duration) *RetryStrategy {
	return &RetryStrategy{
		times:    times,
		interval: interval,
	}
}

// RetryStrategy 重试策略
type RetryStrategy struct {
	times    int           // 重试次数
	interval time.Duration // 重试间隔
}

func (s RetryStrategy) Execute(ctx context.Context, task Task) (int, error) {
	if s.times <= 0 {
		return 0, errorx.New("retry times must be greater than 0")
	}
	logger := log.WithField("id", task.GetID())

	times := 1
	if err := task.Execute(ctx); err != nil {
		for times < s.times {
			logger.WithField("retry_times", times).Error("task retry failed, retrying...")
			if err = task.Execute(ctx); err == nil {
				logger.WithField("retry_times", times).Info("task retry success")
				return times, nil
			}
			time.Sleep(s.interval)
			times++
		}
		return times, errorx.Wrap(err, fmt.Sprintf("task retry failed after %d retries", times))
	}
	logger.WithField("retry_times", times).Info("task retry success")
	return times, nil
}
