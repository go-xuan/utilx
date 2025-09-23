package taskx

import (
	"fmt"
	"sync/atomic"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/marshalx"
	log "github.com/sirupsen/logrus"
)

// ResultHook 任务结果处理钩子函数
type ResultHook func(Result)

// Result 任务结果接口
type Result interface {
	GetTask() Task   // 获取任务
	GetError() error // 获取任务执行错误
}

// TaskResult 任务执行结果
type TaskResult struct {
	task  Task  // 任务
	error error // 任务执行错误
}

func (r *TaskResult) GetTask() Task {
	return r.task
}

func (r *TaskResult) GetError() error {
	return r.error
}

// ErrorsCollect 任务错误数据收集器
type ErrorsCollect struct {
	Ids    []string `json:"ids"`    // 任务ID
	Errors []*Error `json:"errors"` // 失败任务错误
}

func (c *ErrorsCollect) ResultHook(result Result) {
	id := result.GetTask().GetID()
	c.Ids = append(c.Ids, id)
	if err := result.GetError(); err != nil {
		c.Errors = append(c.Errors, &Error{
			ID:  id,
			Msg: err.Error(),
		})
	}
}

// Save 保存任务错误数据
func (c *ErrorsCollect) Save(path string) error {
	if c == nil || path == "" {
		return nil
	}
	if err := marshalx.Apply(path).Write(path, c); err != nil {
		return errorx.Wrap(err, "save tasks collect")
	}
	return nil
}

// ErrorCount 统计错误次数
func ErrorCount(count *int64) ResultHook {
	return func(result Result) {
		if result.GetError() != nil {
			atomic.AddInt64(count, 1)
		}
	}
}

// LogResult 任务结果记录日志
func LogResult(result Result) {
	logger := log.WithField("task_id", result.GetTask().GetID())
	if err := result.GetError(); err != nil {
		logger.WithField("error", err.Error()).Error("task execute error")
	} else {
		logger.Info("task execute success")
	}
}

// PrintResult 打印任务执行结果（最好在并发数==1时使用，日志输出比较直观）
func PrintResult(result Result) {
	fmt.Println("=====================================================")
	if err := result.GetError(); err != nil {
		fmt.Printf("task[%s] execute error: %s\n", result.GetTask().GetID(), result.GetError().Error())
	} else {
		fmt.Printf("task[%s] execute success\n", result.GetTask().GetID())
	}
	fmt.Println("=====================================================")
}
