package taskx

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/idx"
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

func (r *Retry) Execute(ctx context.Context, execute Execute) error {
	if r.times <= 0 {
		return errorx.New("retry times must be greater than 0")
	}
	times := 1
	if err := execute(ctx); err != nil {
		for times < r.times {
			log.WithField("retry_times", times).Error("retry failed, retrying...")
			if err = execute(ctx); err == nil {
				log.WithField("retry_times", times).Info("retry success finally")
				return nil
			}
			time.Sleep(r.interval)
			times++
		}
		return errorx.Wrap(err, fmt.Sprintf("retry failed after %d retries", times))
	}
	log.WithField("retry_times", times).Info("retry success finally")
	return nil
}

// ResultHook 结果处理钩子函数
func (r *Retry) ResultHook(result Result) {
	if err := result.GetError(); err != nil {
		log.WithError(err).Error("task execute error, retrying...")
		if err = r.Execute(context.Background(), result.GetExecute()); err != nil {
			log.WithError(err).Error("task retry error")
		}
	}
}

// NewRetryTask 创建重试任务调度器
func NewRetryTask(times int, interval time.Duration, name ...string) *RetryTask {
	var name_ string
	if len(name) > 0 {
		name_ = name[0]
	} else {
		name_ = idx.SnowFlake().String()
	}
	return &RetryTask{retry: NewRetry(times, interval), name: name_}
}

// RetryTask 重试任务调度器
type RetryTask struct {
	name    string  // 任务名
	retry   *Retry  // 重试策略
	execute Execute // 任务执行函数
}

func (s *RetryTask) GetUnique() string {
	return s.name
}

func (s *RetryTask) Execute(ctx context.Context) error {
	logger := log.WithField("task_type", "retry").WithField("task_name", s.name)
	if s.execute == nil {
		logger.Error("retry execute is nil")
		return errorx.New("retry execute is nil")
	}
	if s.retry == nil {
		return errorx.New("retry is nil")
	}
	if err := s.retry.Execute(ctx, s.execute); err != nil {
		logger.WithError(err).Error("retry task execute error")
		return errorx.Wrap(err, "retry task execute error")
	}
	logger.Info("retry task execute success")
	return nil
}

// AddExecute 设置任务执行函数
func (s *RetryTask) AddExecute(execute Execute) *RetryTask {
	s.execute = execute
	return s
}
