package taskx

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
)

// 重试配置常量
const (
	defaultRetryTimes      = 3                // 最大重试次数
	defaultRetryInterval   = 10 * time.Second // 初始重试间隔
	defaultRetryMultiplier = 1.5              // 重试间隔倍增系数
)

// NewRetry 创建重试任务
func NewRetry(task Task, strategy *RetryStrategy) *Retry {
	return &Retry{
		task:     task,
		strategy: strategy,
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
	err := t.strategy.Execute(ctx, t.task)
	if err != nil {
		logger.WithError(err).Error("retry task execute error")
		return errorx.Wrap(err, "retry task execute error")
	}
	logger.Info("retry task execute success")
	return nil
}

// DefaultRetryStrategy 默认重试策略
func DefaultRetryStrategy() *RetryStrategy {
	return NewRetryStrategy(3, 5*time.Second, 1.5)
}

// NewRetryStrategy 创建重试策略
func NewRetryStrategy(times int, interval time.Duration, multiplier float64) *RetryStrategy {
	if times <= 0 {
		times = defaultRetryTimes
	}
	if interval <= 0 {
		interval = defaultRetryInterval
	}
	if multiplier <= 0 {
		multiplier = defaultRetryMultiplier
	}
	return &RetryStrategy{
		times:      times,
		interval:   interval,
		multiplier: multiplier,
	}
}

// RetryStrategy 重试策略
type RetryStrategy struct {
	times      int           // 重试次数
	interval   time.Duration // 重试间隔
	multiplier float64       // 重试间隔倍数
}

func (s RetryStrategy) Execute(ctx context.Context, task Task) error {
	logger := log.WithField("id", task.GetID())

	start, interval := time.Now(), s.interval
	for times := 0; times < s.times; times++ {
		if times > 0 {
			logger.WithField("retry_times", times).
				WithField("next_retry_interval", interval).
				WithField("elapsed", time.Since(start)).
				Error("task retry failed, retrying...")
			select {
			case <-ctx.Done():
				return errorx.Wrap(ctx.Err(), "retry cancelled")
			case <-time.After(interval):
			}
			interval = time.Duration(float64(interval) * s.multiplier)
		}
		if err := task.Execute(ctx); err == nil {
			if times > 0 {
				logger.WithField("retry_times", times).
					WithField("elapsed", time.Since(start)).
					Info("task retry success")
			}
			return nil
		}
	}
	return errorx.New(fmt.Sprintf("task retry failed after %d retries, elapsed: %s", s.times, time.Since(start)))
}
