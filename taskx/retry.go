package taskx

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
)

// NewRetryTask 创建重试任务
func NewRetryTask(task Task, times int, interval time.Duration) *RetryTask {
	return &RetryTask{
		retry: NewRetry(times, interval),
		task:  task,
	}
}

// RetryTask 重试任务包装
type RetryTask struct {
	retry *Retry // 重试策略
	task  Task   // 任务执行函数
}

func (t *RetryTask) GetID() string {
	return fmt.Sprintf("retry@%s", t.task.GetID())
}

func (t *RetryTask) Execute(ctx context.Context) error {
	logger := log.WithField("id", t.GetID())

	if t.task == nil {
		logger.Error("retry task is nil")
		return errorx.New("retry task is nil")
	}

	if t.retry == nil {
		return errorx.New("retry retry is nil")
	}

	if times, err := t.retry.Execute(ctx, t.task); err != nil {
		logger.WithField("error", err.Error()).Error("retry task execute error")
		return errorx.Wrap(err, "retry task execute error")
	} else {
		logger.WithField("retry_times", times).Info("retry task execute success")
	}
	return nil
}

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

func (r *Retry) Execute(ctx context.Context, task Task) (int, error) {
	if r.times <= 0 {
		return 0, errorx.New("retry times must be greater than 0")
	}
	logger := log.WithField("id", task.GetID())

	times := 1
	if err := task.Execute(ctx); err != nil {
		for times < r.times {
			logger.WithField("retry_times", times).Error("task retry failed, retrying...")
			if err = task.Execute(ctx); err == nil {
				logger.WithField("retry_times", times).Info("task retry success")
				return times, nil
			}
			time.Sleep(r.interval)
			times++
		}
		return times, errorx.Wrap(err, fmt.Sprintf("task retry failed after %d retries", times))
	}
	logger.WithField("retry_times", times).Info("task retry success")
	return times, nil
}

func (r *Retry) ResultHook(result Result) {
	if err := result.GetError(); err != nil {
		log.WithField("error", err.Error()).Error("task execute error, retrying...")
		if _, err = r.Execute(context.Background(), result.GetTask()); err != nil {
			log.WithField("error", err.Error()).Error("task concurrency error")
		}
	}
}
