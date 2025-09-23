package taskx

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
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

func (r *Retry) Execute(ctx context.Context, task Task) error {
	if r.times <= 0 {
		return errorx.New("retry times must be greater than 0")
	}
	logger := log.WithField("task_id", task.GetID()).
		WithField("task_type", "retry")

	times := 1
	if err := task.Execute(ctx); err != nil {
		for times < r.times {
			logger.WithField("retry_times", times).Error("retry failed, retrying...")
			if err = task.Execute(ctx); err == nil {
				logger.WithField("retry_times", times).Info("retry success finally")
				return nil
			}
			time.Sleep(r.interval)
			times++
		}
		return errorx.Wrap(err, fmt.Sprintf("retry failed after %d retries", times))
	}
	logger.WithField("retry_times", times).Info("retry success finally")
	return nil
}

func (r *Retry) ResultHook(result Result) {
	if err := result.GetError(); err != nil {
		log.WithField("error", err.Error()).Error("task execute error, retrying...")
		if err = r.Execute(context.Background(), result.GetTask()); err != nil {
			log.WithField("error", err.Error()).Error("task concurrency error")
		}
	}
}
