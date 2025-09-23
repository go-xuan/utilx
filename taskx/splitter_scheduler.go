package taskx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// NewSplitterScheduler 新建拆分任务
func NewSplitterScheduler[T any](id string, splitter *Splitter, tasks ...T) *SplitterScheduler[T] {
	return &SplitterScheduler[T]{
		id:       id,
		splitter: splitter,
		tasks:    tasks,
	}
}

// SplitterScheduler 多任务拆分执行
type SplitterScheduler[T any] struct {
	id       string                           // 任务名
	splitter *Splitter                        // 并发策略
	tasks    []T                              // 所有待处理任务
	execute  func(context.Context, []T) error // 批量执行函数
}

func (t *SplitterScheduler[T]) GetID() string {
	return t.id
}

func (t *SplitterScheduler[T]) Execute(ctx context.Context) error {
	logger := log.WithField("task_id", t.GetID()).
		WithField("task_type", "splitter_scheduler")

	if len(t.tasks) == 0 {
		return errorx.New("splitter scheduler tasks is nil")
	}
	if t.execute == nil {
		return errorx.New("splitter scheduler execute function is nil")
	}

	if t.splitter == nil {
		if err := t.execute(ctx, t.tasks); err != nil {
			logger.WithField("error", err.Error()).
				Error("splitter scheduler execute error")
			return err
		}
		logger.Info("splitter scheduler execute success")
		return nil
	}
	// 并发执行子任务
	return t.splitter.Execute(ctx, len(t.tasks), func(ctx context.Context, start, end, times int) error {
		logger = logger.WithField("start", start).
			WithField("end", end).
			WithField("times", times)
		if err := t.execute(ctx, t.tasks[start:end]); err != nil {
			logger.WithField("error", err.Error()).
				Error("splitter execute error")
			return err
		}
		logger.Info("splitter execute success")
		return nil
	})
}

// AddTask 添加任务
func (t *SplitterScheduler[T]) AddTask(tasks ...T) *SplitterScheduler[T] {
	if len(tasks) > 0 {
		t.tasks = append(t.tasks, tasks...)
	}
	return t
}

// SetBatchExecute 设置批量执行函数
func (t *SplitterScheduler[T]) SetBatchExecute(execute func(context.Context, []T) error) *SplitterScheduler[T] {
	t.execute = execute
	return t
}
