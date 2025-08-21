package taskx

import (
	"context"
	"fmt"
	"time"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// NewRetry 新建重试策略
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
func (r *Retry) Execute(ctx context.Context, execute Execute) error {
	if err := execute(ctx); err != nil {
		if r.times > 0 {
			curr := 1
			for curr <= r.times {
				log.WithField("retry_times", curr).Error("retry failed, retrying...")
				if err = execute(ctx); err == nil {
					log.WithField("retry_times", curr).Info("retry success finally")
					return nil
				}
				time.Sleep(r.interval)
				curr++
			}
			return errorx.Wrap(err, fmt.Sprintf("retry failed after %d retries", r.times))
		}
		return errorx.Wrap(err, fmt.Sprintf("retry failed, no retry"))
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
	Retry           // 重试策略
	execute Execute // 任务实例
}

func (s *RetryTask) Execute(ctx context.Context) error {
	if s.execute == nil {
		return errorx.New("retry execute is nil")
	}
	if err := s.Retry.Execute(ctx, s.execute); err != nil {
		return errorx.Wrap(err, "retry execute failed")
	}
	return nil
}

// AddExecute 设置任务执行函数
func (s *RetryTask) AddExecute(execute Execute) *RetryTask {
	s.execute = execute
	return s
}
