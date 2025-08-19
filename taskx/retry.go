package taskx

import (
	"context"
	"fmt"
	"time"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// NewRetry 创建重试
func NewRetry(times int, interval time.Duration) *Retry {
	return &Retry{
		times:    times,
		interval: interval,
	}
}

// Retry 重试策略
type Retry struct {
	times    int           // 重试次数
	interval time.Duration // 重试间隔
}

// Execute 重试执行
func (r *Retry) Execute(ctx context.Context, task Task) error {
	if err := task.Execute(ctx); err != nil {
		if r.times > 0 {
			curr := 1
			for curr <= r.times {
				log.WithField("retry_times", curr).Error("task retry failed, retrying...")
				if err = task.Execute(ctx); err == nil {
					log.WithField("retry_times", curr).Info("task retry success finally")
					return nil
				}
				time.Sleep(r.interval)
				curr++
			}
			return errorx.Wrap(err, fmt.Sprintf("task retry failed after %d retries", r.times))
		}
		return errorx.Wrap(err, fmt.Sprintf("task retry failed, no retry"))
	}
	return nil
}

// NewRetryTask 创建重试任务调度器
func NewRetryTask(times int, interval time.Duration) *RetryTask {
	return &RetryTask{
		Retry: Retry{
			times:    times,
			interval: interval,
		},
	}
}

// RetryTask 重试任务调度器
type RetryTask struct {
	Retry      // 重试策略
	task  Task // 任务实例
}

func (s *RetryTask) Execute(ctx context.Context) error {
	if s.task == nil {
		return errorx.New("retry task is nil")
	}
	if err := s.Retry.Execute(ctx, s.task); err != nil {
		return errorx.Wrap(err, "retry task execute failed")
	}
	return nil
}

// AddTask 设置任务执行
func (s *RetryTask) AddTask(task Task) *RetryTask {
	s.task = task
	return s
}
