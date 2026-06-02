package taskx

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

// IResult 任务结果接口
type IResult interface {
	GetID() string   // 获取任务ID
	GetTask() Task   // 获取带上下文的执行任务
	GetError() error // 获取任务执行错误
}

// Result 任务执行结果
type Result struct {
	task  Task  // 带上下文的执行任务
	error error // 任务执行错误
}

// GetID 获取任务ID
func (r *Result) GetID() string {
	if r.task != nil {
		return r.task.GetID()
	}
	return ""
}

func (r *Result) GetTask() Task {
	return r.task
}

func (r *Result) GetError() error {
	return r.error
}

// ResultHook 任务结果处理钩子函数
type ResultHook func(IResult)

// ErrorCount 错误统计次数
func ErrorCount(count *int64) ResultHook {
	return func(result IResult) {
		if err := result.GetError(); err != nil {
			atomic.AddInt64(count, 1)
		}
	}
}

// ErrorRetry 失败重试
func ErrorRetry(times int, interval time.Duration, multiplier float64) ResultHook {
	strategy := NewRetryStrategy(times, interval, multiplier)
	return func(result IResult) {
		if err := result.GetError(); err != nil {
			log.WithError(err).Error("retrying...")
			if err = strategy.Execute(context.Background(), result.GetTask()); err != nil {
				log.WithError(err).Error("task retry error")
			}
		}
	}
}

// LogResult 任务结果记录日志
func LogResult(result IResult) {
	logger := log.WithField("task_id", result.GetID())
	if err := result.GetError(); err != nil {
		logger.WithError(err).Error("task execute error")
	} else {
		logger.Info("task execute success")
	}
}

// PrintResult 打印任务执行结果（最好在并发数==1时使用，日志输出比较直观）
func PrintResult(result IResult) {
	fmt.Println("=====================================================")
	if err := result.GetError(); err != nil {
		fmt.Printf("task[%s] execute error: %s\n", result.GetID(), result.GetError().Error())
	} else {
		fmt.Printf("task[%s] execute success\n", result.GetID())
	}
	fmt.Println("=====================================================")
}
