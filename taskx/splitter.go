package taskx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
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
func NewSplitterTask[T any](limit int) *SplitterTask[T] {
	return &SplitterTask[T]{
		Splitter: *NewSplitter(limit),
		list:     make([]T, 0),
	}
}

// SplitterTask 多任务拆分执行
type SplitterTask[T any] struct {
	Splitter                 // 并发策略
	list     []T             // 所有待处理任务
	execute  BatchExecute[T] // 批量执行函数
}

func (t *SplitterTask[T]) Add(v T) *SplitterTask[T] {
	t.list = append(t.list, v)
	return t
}

func (t *SplitterTask[T]) SetList(list []T) *SplitterTask[T] {
	t.list = list
	return t
}

func (t *SplitterTask[T]) SetExecute(execute BatchExecute[T]) *SplitterTask[T] {
	t.execute = execute
	return t
}

func (t *SplitterTask[T]) Execute(ctx context.Context) error {
	if len(t.list) == 0 {
		return errorx.New("splitter list is nil")
	}
	if t.execute == nil {
		return errorx.New("splitter execute is nil")
	}

	// 并发执行子任务
	return t.Splitter.Execute(ctx, len(t.list), func(ctx context.Context, start, end, batch int) error {
		logger := log.WithField("batch", batch).WithField("start", start).WithField("end", end)
		if err := t.execute(ctx, t.list[start:end]); err != nil {
			logger.WithError(err).Error("batch execute error")
			return err
		}
		logger.Info("batch execute success")
		return nil
	})
}
