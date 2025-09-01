package taskx

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/idx"
)

// NewSplitter 新建拆分策略
func NewSplitter(limit int) *Splitter {
	return &Splitter{
		limit: limit,
	}
}

// Splitter 拆分策略
type Splitter struct {
	limit int // 每批执行任务数
}

func (s *Splitter) Execute(ctx context.Context, total int, execute func(ctx context.Context, start, end, batch int) error) error {
	var offset, batch int
	for offset < total {
		batch++
		next := offset + s.limit
		if next > total {
			next = total
		}
		if err := execute(ctx, offset, next, batch); err != nil {
			return errorx.Wrap(err, "cursor execute error")
		}
		offset = next
	}
	return nil
}

// NewSplitterTask 新建拆分任务
func NewSplitterTask[T any](limit int, name ...string) *SplitterTask[T] {
	var name_ string
	if len(name) > 0 {
		name_ = name[0]
	} else {
		name_ = idx.SnowFlake().String()
	}
	return &SplitterTask[T]{name: name_, splitter: NewSplitter(limit), tasks: make([]T, 0)}
}

// SplitterTask 多任务拆分执行
type SplitterTask[T any] struct {
	name         string
	splitter     *Splitter       // 并发策略
	tasks        []T             // 所有待处理任务
	batchExecute BatchExecute[T] // 批量执行函数
}

// AddTask 添加任务
func (t *SplitterTask[T]) AddTask(tasks ...T) *SplitterTask[T] {
	if len(tasks) > 0 {
		t.tasks = append(t.tasks, tasks...)
	}
	return t
}

// SetBatchExecute 设置批量执行函数
func (t *SplitterTask[T]) SetBatchExecute(execute BatchExecute[T]) *SplitterTask[T] {
	t.batchExecute = execute
	return t
}

func (t *SplitterTask[T]) GetUnique() string {
	return t.name
}

func (t *SplitterTask[T]) Execute(ctx context.Context) error {
	logger := log.WithField("task_type", "splitter").WithField("task_name", t.name)
	if len(t.tasks) == 0 {
		return errorx.New("splitter tasks is nil")
	}
	if t.batchExecute == nil {
		return errorx.New("splitter batch execute is nil")
	}

	if t.splitter == nil {
		if err := t.batchExecute(ctx, t.tasks); err != nil {
			logger.WithError(err).Error("batch execute error")
			return err
		}
		logger.Info("batch execute success")
		return nil
	}
	// 并发执行子任务
	return t.splitter.Execute(ctx, len(t.tasks), func(ctx context.Context, start, end, batch int) error {
		logger = logger.WithField("start", start).WithField("end", end)
		if err := t.batchExecute(ctx, t.tasks[start:end]); err != nil {
			logger.WithError(err).Error("batch execute error")
			return err
		}
		logger.Info("batch execute success")
		return nil
	})
}
