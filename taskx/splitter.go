package taskx

import (
	"context"
	"fmt"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// NewSplitter 新建拆分任务
func NewSplitter[T any](limit int, tasks ...T) *Splitter[T] {
	return &Splitter[T]{
		strategy: NewSplitterStrategy(limit),
		tasks:    tasks,
	}
}

// Splitter 多任务拆分执行
type Splitter[T any] struct {
	tasks     []T                              // 所有待处理任务
	batchFunc func(context.Context, []T) error // 批量执行函数
	strategy  *SplitterStrategy                // 并发策略
}

func (s *Splitter[T]) GetID() string {
	return fmt.Sprintf("splitter@%d", s.strategy.limit)
}

func (s *Splitter[T]) Execute(ctx context.Context) error {
	logger := log.WithField("task_id", s.GetID())
	if len(s.tasks) == 0 {
		return errorx.New("splitter tasks is nil")
	} else if s.batchFunc == nil {
		return errorx.New("splitter batch function is nil")
	}
	if s.strategy == nil {
		// 无策略，直接执行全部任务
		if err := s.batchFunc(ctx, s.tasks); err != nil {
			logger.WithError(err).Error("splitter execute error")
			return err
		}
		logger.Info("splitter execute success")
		return nil
	}
	return s.strategy.Execute(ctx, len(s.tasks), func(ctx context.Context, start, end, times int) (bool, error) {
		logger = logger.WithField("start", start).WithField("end", end).WithField("times", times)
		if err := s.batchFunc(ctx, s.tasks[start:end]); err != nil {
			logger.WithError(err).Error("batch function execute error")
			return true, err
		}
		logger.Info("batch function execute success")
		return false, nil
	})
}

// SetBatchExecute 设置批量执行函数
func (s *Splitter[T]) SetBatchExecute(batchFunc func(context.Context, []T) error) *Splitter[T] {
	s.batchFunc = batchFunc
	return s
}

// AddTask 添加任务
func (s *Splitter[T]) AddTask(tasks ...T) *Splitter[T] {
	if len(tasks) > 0 {
		s.tasks = append(s.tasks, tasks...)
	}
	return s
}

// NewSplitterStrategy 新建拆分策略
func NewSplitterStrategy(limit int) *SplitterStrategy {
	return &SplitterStrategy{
		limit: limit,
	}
}

// SplitterStrategy 拆分策略
type SplitterStrategy struct {
	limit int // 每批执行任务数
}

// Execute 执行拆分任务
// start 开始索引
// end 结束索引
// times 执行次数
func (s *SplitterStrategy) Execute(ctx context.Context, total int, execute func(ctx context.Context, start, end, times int) (bool, error)) error {
	var offset, times int
	for offset < total {
		times++
		next := offset + s.limit
		if next > total {
			next = total
		}
		if stop, err := execute(ctx, offset, next, times); err != nil {
			return errorx.Wrap(err, "batch execute error")
		} else if stop {
			break
		}
		offset = next
	}
	return nil
}
